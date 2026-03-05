package keeper

import (
	"context"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func (k *Keeper) BeginDealing(ctx context.Context, latestRound *types.DKGNetwork) error {
	verifiedRegCount, err := k.countDKGRegistrationsByStatus(ctx, latestRound.CodeCommitment, latestRound.Round, types.DKGRegStatusVerified)
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
		asyncCtx, cancel := dkgAsyncContext()
		go func() {
			defer cancel()
			k.handleDKGDealing(asyncCtx, latestRound)
		}()
	}

	return nil
}

// ProcessJustifications handles justifications from the Vote Extension.
// Like ProcessDeals and ProcessResponses, it emits an event and then delegates
// verification and story-kernel forwarding only when the DKG service is enabled.
//
// This function does NOT affect on-chain state. Invalid justifications are
// silently dropped (logged but not returned as errors) to prevent consensus
// divergence between validators running with different DKG service configurations.
//
// Replay protection is structural, not explicit:
//   - Cross-round: kyber's DistKeyGenerator derives SessionID from (dealer pubkey +
//     verifiers + commitments + threshold). Each round produces fresh keys and thus
//     unique SessionIDs. Since SessionID is embedded in Justification.Hash(), Schnorr
//     signature verification implicitly rejects replayed justifications from other rounds.
//   - Stage gating: justifications are only accepted during DKGStageDealing (checked
//     in msg_server.go) and PhaseDealing (checked in handleDKGProcessJustifications).
//   - Within-block: deduplication by (dealerIndex, recipientIndex) prevents redundant
//     processing of the same justification broadcast by multiple validators.
func (k *Keeper) ProcessJustifications(ctx context.Context, latestRound *types.DKGNetwork, justifications []types.Justification) error {
	if err := k.emitBeginProcessJustifications(ctx, latestRound, justifications); err != nil {
		return errors.Wrap(err, "failed to emit begin process justifications event")
	}

	if k.isDKGSvcEnabled {
		asyncCtx, cancel := dkgAsyncContext()
		go func() {
			defer cancel()
			k.handleDKGProcessJustifications(asyncCtx, latestRound, justifications)
		}()
	}

	return nil
}

func (k *Keeper) ProcessDeals(ctx context.Context, latestRound *types.DKGNetwork, deals []types.Deal) error {
	if err := k.emitBeginProcessDeals(ctx, latestRound, deals); err != nil {
		return errors.Wrap(err, "failed to emit begin process deals event")
	}

	if k.isDKGSvcEnabled {
		asyncCtx, cancel := dkgAsyncContext()
		go func() {
			defer cancel()
			k.handleDKGProcessDeals(asyncCtx, latestRound, deals)
		}()
	}

	return nil
}

func (k *Keeper) ProcessResponses(ctx context.Context, latestRound *types.DKGNetwork, responses []types.Response) error {
	if err := k.emitBeginProcessResponses(ctx, latestRound, responses); err != nil {
		return errors.Wrap(err, "failed to emit begin process responses event")
	}

	if k.isDKGSvcEnabled {
		asyncCtx, cancel := dkgAsyncContext()
		go func() {
			defer cancel()
			k.handleDKGProcessResponses(asyncCtx, latestRound, responses)
		}()
	}

	return nil
}
