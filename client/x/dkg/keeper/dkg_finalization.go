package keeper

import (
	"context"
	"os"
	"strconv"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func (k *Keeper) BeginFinalization(ctx context.Context, latestRound *types.DKGNetwork) error {
	// TEST INJECTION: force-invalidate a dealer by index for E2E testing of #717/#719.
	// Placed in BeginFinalization (called every round from BeginBlocker) so it fires
	// even when no justifications exist. Set DKG_TEST_INVALIDATE_INDEX on ALL validators.
	if idxStr := os.Getenv("DKG_TEST_INVALIDATE_INDEX"); idxStr != "" {
		if idx, err := strconv.ParseUint(idxStr, 10, 32); err == nil && idx > 0 {
			if invErr := k.invalidateDealerRegistration(ctx, latestRound, uint32(idx)); invErr != nil {
				log.Warn(ctx, "TEST: force-invalidate dealer failed", invErr, "index", idx)
			} else {
				log.Info(ctx, "TEST: force-invalidated dealer", "index", idx, "round", latestRound.Round)
			}
		}
	}

	if err := k.emitBeginDKGFinalization(ctx, latestRound); err != nil {
		return errors.Wrap(err, "failed to emit begin DKG finalization event")
	}

	if k.isDKGSvcEnabled {
		asyncCtx, cancel := dkgAsyncContext()

		go func() {
			defer cancel()

			k.handleDKGFinalization(asyncCtx, latestRound)
		}()
	}

	return nil
}

func (k *Keeper) FinalizeDKGRound(ctx context.Context, latestRound *types.DKGNetwork) error {
	finalizedCount, err := k.countDKGRegistrationsByStatus(ctx, latestRound.Round, types.DKGRegStatusFinalized)
	if err != nil {
		return errors.Wrap(err, "failed to fetch DKG registrations in Finalized status")
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG params")
	}

	// Check if finalized count meets the minimum required participants
	if finalizedCount < params.MinReqFinalizedParticipants {
		log.Info(ctx, "Finalized registration count below minimum required. Skipping current round.",
			"finalized_count", finalizedCount,
			"min_req_finalized", params.MinReqFinalizedParticipants,
			"current", latestRound.Round,
			"next", latestRound.Round+1,
		)

		return k.SkipToNextRound(ctx, latestRound)
	}

	// Check if finalized count meets the operational threshold
	if finalizedCount < latestRound.Threshold {
		log.Info(ctx, "Finalized registration count below operational threshold. Skipping current round.",
			"finalized_count", finalizedCount,
			"threshold", latestRound.Threshold,
			"current", latestRound.Round,
			"next", latestRound.Round+1,
		)

		return k.SkipToNextRound(ctx, latestRound)
	}

	if err := k.emitDKGFinalized(ctx, latestRound); err != nil {
		return errors.Wrap(err, "failed to emit DKG finalized event")
	}

	// Distribute DKG committee rewards before updating the active round,
	// so that getLatestActiveDKGNetwork still returns the previous active round.
	if err := k.settleRewardsForPreviousCommittee(ctx); err != nil {
		return errors.Wrap(err, "failed to distribute DKG committee rewards")
	}

	if err := k.distributeCDRRewardPool(ctx); err != nil {
		return errors.Wrap(err, "failed to distribute CDR fee pool")
	}

	// End the previous active round's stage before updating the active round pointer,
	// so it no longer undergoes stage transitions.
	if err := k.endPreviousActiveRound(ctx, latestRound.Round); err != nil {
		return errors.Wrap(err, "failed to end previous active round")
	}

	if err := k.setLatestActiveRound(ctx, latestRound); err != nil {
		return errors.Wrap(err, "failed to set the latest active round of DKG")
	}

	// If this was an upgrade round, delete the activated upgrade info and log completion
	if latestRound.IsUpgrade {
		if err := k.deleteActivatedUpgradeInfo(ctx); err != nil {
			return errors.Wrap(err, "failed to delete activated upgrade info after successful upgrade round")
		}

		log.Info(ctx, "Upgrade resharing round completed, new TEE binary is now active",
			"round", latestRound.Round,
		)
	}

	if k.isDKGSvcEnabled {
		asyncCtx, cancel := dkgAsyncContext()

		go func() {
			defer cancel()

			k.handleDKGComplete(asyncCtx, latestRound)
		}()
	}

	log.Info(ctx, "DKG network setup completed", "round", latestRound.Round)

	return nil
}
