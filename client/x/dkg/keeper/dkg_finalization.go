package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func (k *Keeper) BeginFinalization(ctx context.Context, latestRound *types.DKGNetwork) error {
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

	// Drop finalized validators whose finalize vote diverged from the consensus
	// polynomial before counting.
	if k.isV190Round(ctx, latestRound) {
		if err := k.invalidateNonConsensusFinalizations(ctx, latestRound); err != nil {
			return errors.Wrap(err, "failed to invalidate non-consensus finalizations")
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

// pruneRoundFinalizeVotes removes all recorded finalize votes for a round, used when the
// round is abandoned via SkipToNextRound (the success path prunes them in
// invalidateNonConsensusFinalizations).
func (k *Keeper) pruneRoundFinalizeVotes(ctx context.Context, round uint32) error {
	regs, err := k.getDKGRegistrationsByRound(ctx, round)
	if err != nil {
		return errors.Wrap(err, "failed to fetch DKG registrations")
	}

	for i := range regs {
		key := finalizeVoteStoreKey(round, common.HexToAddress(regs[i].ValidatorAddr))
		if err := k.FinalizeVotes.Remove(ctx, key); err != nil {
			return errors.Wrap(err, "failed to prune finalize vote", "validator", regs[i].ValidatorAddr)
		}
	}

	return nil
}

// invalidateNonConsensusFinalizations invalidates finalized validators whose finalize
// vote (globalPubKey + publicCoeffs) does not match the consensus polynomial. Because
// the kernel derives globalPubKey, publicCoeffs and pubKeyShare from the same key share,
// matching the consensus vote implies the validator's share lies on the consensus
// polynomial. The recorded votes are pruned here.
func (k *Keeper) invalidateNonConsensusFinalizations(ctx context.Context, latestRound *types.DKGNetwork) error {
	consensusKey := globalPubKeyVoteKey(latestRound.Round, latestRound.GlobalPublicKey, latestRound.PublicCoeffs)

	finalizedRegs, err := k.getDKGRegistrationsByStatus(ctx, latestRound.Round, types.DKGRegStatusFinalized)
	if err != nil {
		return errors.Wrap(err, "failed to fetch finalized DKG registrations")
	}

	for i := range finalizedRegs {
		reg := finalizedRegs[i]
		voteStoreKey := finalizeVoteStoreKey(latestRound.Round, common.HexToAddress(reg.ValidatorAddr))

		submitted, err := k.FinalizeVotes.Get(ctx, voteStoreKey)
		if err != nil && !errors.Is(err, collections.ErrNotFound) {
			return errors.Wrap(err, "failed to get finalize vote", "validator", reg.ValidatorAddr)
		}

		if err := k.FinalizeVotes.Remove(ctx, voteStoreKey); err != nil {
			return errors.Wrap(err, "failed to prune finalize vote", "validator", reg.ValidatorAddr)
		}

		if submitted == consensusKey {
			continue
		}

		reg.Status = types.DKGRegStatusInvalidated
		if err := k.setDKGRegistration(ctx, common.HexToAddress(reg.ValidatorAddr), &reg); err != nil {
			return errors.Wrap(err, "failed to invalidate non-consensus finalization",
				"validator", reg.ValidatorAddr,
				"index", reg.Index,
			)
		}

		log.Warn(ctx, "Invalidated finalized validator: finalize vote does not match consensus",
			nil,
			"round", latestRound.Round,
			"validator", reg.ValidatorAddr,
			"index", reg.Index,
		)
	}

	return nil
}
