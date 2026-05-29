package keeper

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"go.dedis.ch/kyber/v4/group/edwards25519"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/vss"
)

func (k *Keeper) BeginFinalization(ctx context.Context, latestRound *types.DKGNetwork) error {
	// Invalidate dealers that never submitted a deal during the dealing phase.
	if k.isV190Round(ctx, latestRound) {
		if err := k.invalidateMissingDealers(ctx, latestRound); err != nil {
			return errors.Wrap(err, "failed to invalidate missing dealers")
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
	// Guard: consensus key material must be set before finalizing. If it's empty, the
	// round did not complete key generation and should not be activated.
	if len(latestRound.GlobalPublicKey) == 0 || len(latestRound.PublicCoeffs) == 0 {
		log.Info(ctx, "Consensus key material not set, skipping to next round",
			"round", latestRound.Round,
		)

		return k.SkipToNextRound(ctx, latestRound)
	}

	// Drop finalized validators whose share is off the consensus polynomial before counting.
	if k.isV190Round(ctx, latestRound) {
		if err := k.invalidateDivergentShares(ctx, latestRound); err != nil {
			return errors.Wrap(err, "failed to invalidate divergent shares")
		}
	}

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

	roundsTotal.WithLabelValues(labelRoundCompleted).Inc()
	log.Info(ctx, "DKG network setup completed", "round", latestRound.Round)

	return nil
}

// invalidateDivergentShares invalidates finalized validators whose public key share
// does not lie on the consensus polynomial.
func (k *Keeper) invalidateDivergentShares(ctx context.Context, latestRound *types.DKGNetwork) error {
	finalizedRegs, err := k.getDKGRegistrationsByStatus(ctx, latestRound.Round, types.DKGRegStatusFinalized)
	if err != nil {
		return errors.Wrap(err, "failed to fetch finalized DKG registrations")
	}

	suite := edwards25519.NewBlakeSHA256Ed25519()

	for i := range finalizedRegs {
		reg := finalizedRegs[i]

		// reg.Index is 1-based; the commitment count must equal the threshold.
		onPoly, verifyErr := vss.VerifyPublicKeyShare(suite, reg.PubKeyShare, int(reg.Index), latestRound.PublicCoeffs, latestRound.Threshold)
		if verifyErr == nil && onPoly {
			continue
		}

		reg.Status = types.DKGRegStatusInvalidated
		if err := k.setDKGRegistration(ctx, common.HexToAddress(reg.ValidatorAddr), &reg); err != nil {
			return errors.Wrap(err, "failed to invalidate divergent share",
				"validator", reg.ValidatorAddr,
				"index", reg.Index,
			)
		}

		log.Warn(ctx, "Invalidated finalized validator: share not on consensus polynomial",
			verifyErr,
			"round", latestRound.Round,
			"validator", reg.ValidatorAddr,
			"index", reg.Index,
		)
	}

	return nil
}

// invalidateMissingDealers invalidates verified dealers that never submitted a deal
// during the dealing phase. The dealt set recorded in ProcessDeals is pruned here.
func (k *Keeper) invalidateMissingDealers(ctx context.Context, latestRound *types.DKGNetwork) error {
	regs, err := k.getDKGRegistrationsByRound(ctx, latestRound.Round)
	if err != nil {
		return errors.Wrap(err, "failed to fetch DKG registrations")
	}

	for i := range regs {
		reg := regs[i]
		key := dealtDealerKey(latestRound.Round, reg.Index)

		dealt, err := k.DealtDealers.Has(ctx, key)
		if err != nil {
			return errors.Wrap(err, "failed to check dealt dealer", "index", reg.Index)
		}

		if dealt {
			if err := k.DealtDealers.Remove(ctx, key); err != nil {
				return errors.Wrap(err, "failed to prune dealt dealer", "index", reg.Index)
			}

			continue
		}

		// Only verified dealers can be missing a deal; others (already invalidated) are left as-is.
		if reg.Status != types.DKGRegStatusVerified {
			continue
		}

		reg.Status = types.DKGRegStatusInvalidated
		if err := k.setDKGRegistration(ctx, common.HexToAddress(reg.ValidatorAddr), &reg); err != nil {
			return errors.Wrap(err, "failed to invalidate missing dealer",
				"validator", reg.ValidatorAddr,
				"index", reg.Index,
			)
		}

		log.Warn(ctx, "Invalidated dealer: no deal submitted during dealing phase",
			nil,
			"round", latestRound.Round,
			"validator", reg.ValidatorAddr,
			"index", reg.Index,
		)
	}

	return nil
}
