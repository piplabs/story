package keeper

import (
	"bytes"
	"context"
	"encoding/binary"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"

	"go.dedis.ch/kyber/v4/group/edwards25519"
	"go.dedis.ch/kyber/v4/sign/schnorr"
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

	// Snapshot the dealer committee once so j.Index resolves stably even as the loop below
	// invalidates registrations.
	dealerRound, err := k.resolveDealerRound(ctx, latestRound)
	if err != nil {
		return errors.Wrap(err, "resolve dealer committee round")
	}

	var committee []types.DKGRegistration
	if dealerRound != nil {
		committee, err = k.dealerCommitteeRegs(ctx, *dealerRound)
		if err != nil {
			return errors.Wrap(err, "build dealer committee")
		}
	}

	dealerPubKeys := dealerPubKeyMap(ctx, committee, suite)

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
			// so they cannot finalize or receive committee rewards. Resolve the dealer index
			// against the committee snapshot taken above.
			if int(j.Index) >= len(committee) {
				log.Error(ctx, "Justification dealer index out of committee range", nil,
					"dealer_index", j.Index,
					"committee_size", len(committee),
				)

				continue
			}

			dealerAddr := common.HexToAddress(strings.TrimSpace(committee[j.Index].ValidatorAddr))

			log.Info(ctx, "Justification VSS verification failed, invalidating dealer",
				"dealer_index", j.Index,
				"validator_addr", dealerAddr.Hex(),
				"round", latestRound.Round,
			)

			if err := k.invalidateDealerByAddr(ctx, latestRound.Round, dealerAddr); err != nil {
				log.Error(ctx, "Failed to invalidate dealer registration", err,
					"validator_addr", dealerAddr.Hex(),
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

// markDealersDealt records which dealers submitted a deal so missing dealers can be invalidated
// at BeginFinalization. deal.Index is the dealer's 0-based position in the dealer committee,
// resolved to a registration (address + dkgPubKey) via dealerCommitteeRegs. Deals arrive via
// vote extensions unauthenticated, so each deal's Schnorr signature is verified against the
// dealer's registered dkgPubKey before recording — a forged deal under another dealer's index
// must not shield a missing dealer from invalidation. Verification is deterministic, so all
// nodes agree.
func (k *Keeper) markDealersDealt(ctx context.Context, latestRound *types.DKGNetwork, deals []types.Deal) error {
	dealerRound, err := k.resolveDealerRound(ctx, latestRound)
	if err != nil {
		return err
	}
	if dealerRound == nil {
		return nil
	}

	// committee[i] is the dealer registration at kyber index i; deal.Index indexes directly into it.
	committee, err := k.dealerCommitteeRegs(ctx, *dealerRound)
	if err != nil {
		return err
	}

	suite := edwards25519.NewBlakeSHA256Ed25519()

	log.Debug(ctx, "MarkDealersDealt: dealer committee",
		"round", latestRound.Round,
		"dealer_round", *dealerRound,
		"committee_size", len(committee),
		"num_deals", len(deals),
		"is_resharing", latestRound.IsResharing,
	)

	recorded := 0
	for _, deal := range deals {
		if deal.Index >= uint32(len(committee)) {
			log.Debug(ctx, "MarkDealersDealt: deal index out of dealer-committee range; skipping",
				"round", latestRound.Round,
				"deal_index", deal.Index,
				"committee_size", len(committee),
			)

			continue
		}

		dealerReg := committee[deal.Index]

		// Skip (do not mark) deals whose signature does not verify against the dealer's dkgPubKey.
		if err := verifyDealSignature(suite, deal, dealerReg.DkgPubKey); err != nil {
			log.Debug(ctx, "MarkDealersDealt: deal signature invalid; skipping",
				"round", latestRound.Round,
				"deal_index", deal.Index,
				"error", err,
			)

			continue
		}

		addr := dealerReg.ValidatorAddr

		log.Debug(ctx, "MarkDealersDealt: dealer submitted deal",
			"round", latestRound.Round,
			"deal_index", deal.Index,
			"dealer", addr,
		)

		if err := k.DealtDealers.Set(ctx, dealtDealerKey(latestRound.Round, addr), true); err != nil {
			return errors.Wrap(err, "failed to record dealt dealer", "round", latestRound.Round, "dealer", addr)
		}
		recorded++
	}

	log.Info(ctx, "MarkDealersDealt: recorded dealers that submitted a deal",
		"round", latestRound.Round,
		"committee_size", len(committee),
		"recorded", recorded,
	)

	return nil
}

// verifyDealSignature verifies a deal's Schnorr signature against the dealer's dkgPubKey.
// The signed message matches kyber dkg.Deal.MarshalBinary: little-endian uint32 Index followed
// by the EncryptedDeal cipher. This is deterministic and safe for the consensus path.
func verifyDealSignature(suite *edwards25519.SuiteEd25519, deal types.Deal, dkgPubKey []byte) error {
	if len(dkgPubKey) == 0 {
		return errors.New("dealer has no dkg pubkey")
	}
	if len(deal.GetSignature()) == 0 {
		return errors.New("empty deal signature")
	}

	dealerPub := suite.Point()
	if err := dealerPub.UnmarshalBinary(dkgPubKey); err != nil {
		return errors.Wrap(err, "unmarshal dealer dkg pubkey")
	}

	var b bytes.Buffer
	if err := binary.Write(&b, binary.LittleEndian, deal.Index); err != nil {
		return errors.Wrap(err, "encode deal index")
	}
	b.Write(deal.Deal.Cipher)

	if err := schnorr.Verify(suite, dealerPub, b.Bytes(), deal.GetSignature()); err != nil {
		return errors.Wrap(err, "schnorr signature verification failed")
	}

	return nil
}

// dealerCommitteeRegs returns the dealer round's committee: ALL of the round's registrations
// (any status) sorted by registration index, matching the kernel's kyber committee where the
// 0-based position is the kyber deal/justification Index. Invalidated members keep their slot so
// the survivors' indices do not shift (kyber resharing requires the position to be preserved).
func (k *Keeper) dealerCommitteeRegs(ctx context.Context, round uint32) ([]types.DKGRegistration, error) {
	committee, err := k.getDKGRegistrationsByRound(ctx, round)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch dealer registrations", "round", round)
	}

	sort.SliceStable(committee, func(i, j int) bool {
		return committee[i].Index < committee[j].Index
	})

	return committee, nil
}

// resolveDealerRound returns the round of the committee expected to deal in latestRound. For
// resharing rounds this is the previous active committee (which holds the existing shares);
// otherwise it is the current round. A nil round means there is no dealer committee to attribute
// deals to (e.g. resharing with no prior active round).
func (k *Keeper) resolveDealerRound(ctx context.Context, latestRound *types.DKGNetwork) (*uint32, error) {
	if !latestRound.IsResharing {
		return &latestRound.Round, nil
	}

	prevActive, err := k.getLatestActiveDKGNetwork(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get previous active round")
	}
	if prevActive == nil {
		return nil, nil
	}

	return &prevActive.Round, nil
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
