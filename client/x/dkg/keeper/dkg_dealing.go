package keeper

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"

	"go.dedis.ch/kyber/v4/group/edwards25519"
)

func (k *Keeper) BeginDealing(ctx context.Context, latestRound *types.DKGNetwork) error {
	verifiedRegCount, err := k.countDKGRegistrationsByStatus(ctx, latestRound.Round, types.DKGRegStatusVerified)
	if err != nil {
		return errors.Wrap(err, "failed to fetch verified DKG registrations")
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG params")
	}

	// Check if verified registrations meet the minimum required threshold
	if verifiedRegCount < params.MinReqRegisteredParticipants {
		log.Info(ctx, "Verified registration count below minimum required. Skipping current round.",
			"verified_count", verifiedRegCount,
			"min_req_registered", params.MinReqRegisteredParticipants,
			"current", latestRound.Round,
			"next", latestRound.Round+1,
		)

		return k.SkipToNextRound(ctx, latestRound)
	}

	// Update total and threshold based on actual verified registrations
	latestRound.Total = verifiedRegCount
	latestRound.Threshold = types.CalculateThreshold(verifiedRegCount, params.OperationalThreshold)

	if err := k.setDKGNetwork(ctx, latestRound); err != nil {
		return errors.Wrap(err, "failed to update DKG network with total and threshold")
	}

	log.Info(ctx, "DKG dealing phase started",
		"round", latestRound.Round,
		"total", latestRound.Total,
		"threshold", latestRound.Threshold,
		"operational_threshold_bps", params.OperationalThreshold,
	)

	if err := k.emitBeginDKGDealing(ctx, latestRound); err != nil {
		return errors.Wrap(err, "failed to emit begin DKG dealing event")
	}

	if k.isDKGSvcEnabled {
		// Use a gasless context for KV reads inside isDKGSvcEnabled so that
		// DKG-enabled and DKG-disabled nodes produce identical GasUsed.
		gaslessCtx := gaslessSDKContext(ctx)

		// Set session.Index from on-chain registration while SDK context is available.
		// Session.Index stores the 1-based registration index used for deal/response
		// routing (converted to 0-based for Kyber) and threshold decryption PID.
		if err := k.ensureSessionIndex(gaslessCtx, latestRound.Round); err != nil {
			log.Warn(ctx, "Failed to set session index from registration", err,
				"round", latestRound.Round,
			)
		}

		// Pre-compute shouldDeal while SDK context is available.
		// The async goroutine cannot access the KV store.
		deal, err := k.shouldDeal(gaslessCtx, latestRound)
		if err != nil {
			log.Error(ctx, "Failed to check whether the validator should deal", err)

			return nil
		}

		asyncCtx, cancel := dkgAsyncContext()

		go func() {
			defer cancel()

			k.handleDKGDealing(asyncCtx, latestRound, deal)
		}()
	}

	return nil
}

// ProcessJustifications handles justifications from the Vote Extension.
// Like ProcessDeals and ProcessResponses, it emits an event and then delegates
// Verification (signature, dedup, VSS) and dealer invalidation run synchronously
// in FinalizeBlock regardless of whether the DKG service is enabled, because
// invalidation is a consensus state change that all nodes must apply consistently.
// Only the kernel gRPC forwarding is conditional on isDKGSvcEnabled and runs async.
//
// Replay protection is structural, not explicit:
//   - Cross-round: kyber's DistKeyGenerator derives SessionID from (dealer pubkey +
//     verifiers + commitments + threshold). Each round produces fresh keys and thus
//     unique SessionIDs. Since SessionID is embedded in Justification.Hash(), Schnorr
//     signature verification implicitly rejects replayed justifications from other rounds.
//   - Stage gating: justifications are only accepted during DKGStageDealing (checked
//     in msg_server.go).
//   - Within-block: deduplication by (dealerIndex, recipientIndex) prevents redundant
//     processing of the same justification broadcast by multiple validators.
func (k *Keeper) ProcessJustifications(ctx context.Context, latestRound *types.DKGNetwork, justifications []types.Justification) error {
	if err := k.emitBeginProcessJustifications(ctx, latestRound, justifications); err != nil {
		return errors.Wrap(err, "failed to emit begin process justifications event")
	}

	suite := edwards25519.NewBlakeSHA256Ed25519()

	dealerPubKeys, err := k.buildDealerPubKeyMap(ctx, latestRound, suite)
	if err != nil {
		return errors.Wrap(err, "build dealer pub key map")
	}

	// Step 1: Signature verification (deterministic, no kernel needed).
	var signatureVerified []types.Justification

	for _, j := range justifications {
		if err := verifyJustificationSignature(suite, j, dealerPubKeys); err != nil {
			log.Debug(ctx, "Dropping justification with invalid signature",
				"dealer_index", j.Index,
				"error", err,
			)

			continue
		}

		signatureVerified = append(signatureVerified, j)
	}

	// Step 2: VSS verification + dealer invalidation (runs in SDK context so
	// on-chain state can be updated). This MUST run regardless of isDKGSvcEnabled
	// because invalidation is a consensus state change that all nodes must apply.
	var validJustifications []types.Justification

	for _, j := range signatureVerified {
		valid, err := verifyJustification(latestRound, j)
		if err != nil {
			log.Warn(ctx, "Justification VSS verification error, dropping", err,
				"dealer_index", j.Index,
			)

			continue
		}

		if !valid {
			// Deal was genuinely invalid — invalidate the dealer's on-chain registration
			// so they cannot finalize or receive committee rewards.
			log.Info(ctx, "Justification VSS verification failed, invalidating dealer",
				"dealer_index", j.Index,
				"round", latestRound.Round,
			)

			if err := k.invalidateDealerRegistration(ctx, latestRound, j.Index); err != nil {
				log.Error(ctx, "Failed to invalidate dealer registration", err,
					"dealer_index", j.Index,
				)
			}

			continue
		}

		validJustifications = append(validJustifications, j)
	}

	// Step 4: Forward valid justifications to story-kernel for DKG state update.
	// Only when DKG service is enabled (kernel gRPC is async).
	if k.isDKGSvcEnabled && len(validJustifications) > 0 {
		asyncCtx, cancel := dkgAsyncContext()

		go func() {
			defer cancel()

			k.handleDKGProcessJustifications(asyncCtx, latestRound, wrapJustifications(validJustifications))
		}()
	}

	return nil
}

func (k *Keeper) ProcessDeals(ctx context.Context, latestRound *types.DKGNetwork, deals []types.Deal) error {
	// Record which dealers submitted a deal so missing dealers can be invalidated at
	// BeginFinalization.
	if k.isV190Round(ctx, latestRound) {
		if err := k.markDealersDealt(ctx, latestRound, deals); err != nil {
			return errors.Wrap(err, "failed to mark dealers dealt")
		}
	}

	if err := k.emitBeginProcessDeals(ctx, latestRound, deals); err != nil {
		return errors.Wrap(err, "failed to emit begin process deals event")
	}

	if k.isDKGSvcEnabled {
		asyncCtx, cancel := dkgAsyncContext()

		go func() {
			defer cancel()

			k.handleDKGProcessDeals(asyncCtx, latestRound, wrapDeals(deals))
		}()
	}

	return nil
}

// markDealersDealt records the address of each dealer that submitted a deal, so missing
// dealers can be invalidated at BeginFinalization. deal.Index is the dealer's 0-based kyber
// index within the DEALER committee: the current round for the initial DKG, but the previous
// active committee for resharing rounds (they hold the existing shares). The index is resolved
// to a validator address through that committee's registrations, because the dealer and current
// index spaces are permuted differently across rounds.
func (k *Keeper) markDealersDealt(ctx context.Context, latestRound *types.DKGNetwork, deals []types.Deal) error {
	dealerRound, dealerTotal, err := k.dealerCommitteeRound(ctx, latestRound)
	if err != nil {
		return err
	}
	if dealerRound == nil {
		return nil
	}

	dealerRegs, err := k.getDKGRegistrationsByRound(ctx, *dealerRound)
	if err != nil {
		return errors.Wrap(err, "failed to fetch dealer registrations", "round", *dealerRound)
	}

	// Map the dealer committee's 1-based registration index to validator address.
	addrByIndex := make(map[uint32]string, len(dealerRegs))
	for i := range dealerRegs {
		addrByIndex[dealerRegs[i].Index] = dealerRegs[i].ValidatorAddr
	}

	// DEBUG: resolve each incoming deal's (0-based, dealer-committee) index to a
	// validator address so the logs say WHICH validator dealt — the raw kyber
	// index alone is unreadable. dealerRound/dealerTotal describe the committee
	// expected to deal (prev active committee for resharing, current otherwise).
	log.Debug(ctx, "markDealersDealt: dealer committee",
		"round", latestRound.Round,
		"dealer_round", *dealerRound,
		"dealer_total", dealerTotal,
		"num_deals", len(deals),
		"is_resharing", latestRound.IsResharing,
	)

	recorded := 0
	for _, deal := range deals {
		// Bound the 0-based kyber index by the dealer committee's total, then convert to
		// the 1-based registration index.
		if deal.Index >= dealerTotal {
			log.Debug(ctx, "markDealersDealt: deal index out of dealer-committee range; skipping",
				"round", latestRound.Round,
				"deal_index", deal.Index,
				"dealer_total", dealerTotal,
			)

			continue
		}

		addr, ok := addrByIndex[deal.Index+1]
		if !ok {
			log.Debug(ctx, "markDealersDealt: no dealer address for index; skipping",
				"round", latestRound.Round,
				"deal_index", deal.Index,
				"reg_index", deal.Index+1,
			)

			continue
		}

		log.Debug(ctx, "markDealersDealt: dealer submitted deal",
			"round", latestRound.Round,
			"deal_index", deal.Index,
			"dealer", addr,
		)

		if err := k.DealtDealers.Set(ctx, dealtDealerKey(latestRound.Round, addr), true); err != nil {
			return errors.Wrap(err, "failed to record dealt dealer", "round", latestRound.Round, "dealer", addr)
		}
		recorded++
	}

	log.Info(ctx, "markDealersDealt: recorded dealers that submitted a deal",
		"round", latestRound.Round,
		"dealer_total", dealerTotal,
		"recorded", recorded,
	)

	return nil
}

// dealerCommitteeRound returns the round and total of the committee expected to deal in
// latestRound. For resharing rounds this is the previous active committee (which holds the
// existing shares); otherwise it is the current round. A nil round means there is no dealer
// committee to attribute deals to (e.g. resharing with no prior active round).
func (k *Keeper) dealerCommitteeRound(ctx context.Context, latestRound *types.DKGNetwork) (*uint32, uint32, error) {
	if !latestRound.IsResharing {
		return &latestRound.Round, latestRound.Total, nil
	}

	prevActive, err := k.getLatestActiveDKGNetwork(ctx)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get previous active round")
	}
	if prevActive == nil {
		return nil, 0, nil
	}

	return &prevActive.Round, prevActive.Total, nil
}

func (k *Keeper) ProcessResponses(ctx context.Context, latestRound *types.DKGNetwork, responses []types.Response) error {
	if err := k.emitBeginProcessResponses(ctx, latestRound, responses); err != nil {
		return errors.Wrap(err, "failed to emit begin process responses event")
	}

	if k.isDKGSvcEnabled {
		// Use a gasless context for KV reads inside isDKGSvcEnabled so that
		// DKG-enabled and DKG-disabled nodes produce identical GasUsed.
		gaslessCtx := gaslessSDKContext(ctx)

		// Pre-compute shouldProcessResponses while SDK context is available.
		shouldProcess, err := k.shouldProcessResponses(gaslessCtx, latestRound)
		if err != nil {
			log.Error(ctx, "Failed to check shouldProcessResponses", err)

			return nil
		}

		asyncCtx, cancel := dkgAsyncContext()

		go func() {
			defer cancel()

			k.handleDKGProcessResponses(asyncCtx, latestRound, wrapResponses(responses), shouldProcess)
		}()
	}

	return nil
}

// ensureSessionIndex sets the session's 1-based index from the on-chain registration,
// if not already set. Must be called from a context where the KV store is accessible.
// For old-only resharing members who have no registration in the current round,
// the index remains 0 (unset).
func (k *Keeper) ensureSessionIndex(ctx context.Context, round uint32) error {
	session, err := k.stateManager.GetSession(round)
	if err != nil {
		return err
	}

	if session.Index != 0 {
		return nil
	}

	reg, err := k.getDKGRegistration(ctx, round, common.HexToAddress(k.validatorEVMAddr))
	if err != nil {
		return errors.Wrap(err, "registration lookup failed")
	}

	session.Index = reg.Index

	log.Info(ctx, "Session index set from on-chain registration",
		"round", round,
		"index", session.Index,
	)

	return k.stateManager.UpdateSession(ctx, session)
}
