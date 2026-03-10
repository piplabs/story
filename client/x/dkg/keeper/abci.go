package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

const dkgStartBlock = 10

func (k *Keeper) BeginBlocker(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentHeight := sdkCtx.BlockHeight()

	// TODO: temporal code for delaying the DKG setup
	if currentHeight < dkgStartBlock {
		return nil
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get params")
	}

	latestRound, err := k.GetLatestDKGRound(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get latest DKG round")
	}

	if latestRound == nil {
		// No active DKG round, start the first round
		log.Info(ctx, "No active DKG round, starting the first round")

		return k.InitiateDKGRound(ctx, false)
	}

	// Check for pending kernel upgrade activation before normal stage transitions.
	// If an upgrade is pending and ready to activate, delete the upgrade info and
	// initiate an upgrade resharing round.
	upgradeInfo, err := k.hasPendingUpgradeActivation(ctx, currentHeight)
	if err != nil {
		return err
	}

	if upgradeInfo != nil {
		log.Info(ctx, "Kernel upgrade activated, initiating upgrade resharing round",
			"upgrade_version", upgradeInfo.UpgradeVersion,
			"activation_height", upgradeInfo.ActivationHeight,
		)

		if err := k.DeleteKernelUpgradeInfo(ctx, upgradeInfo.UpgradeVersion); err != nil {
			return errors.Wrap(err, "failed to delete kernel upgrade info after activation")
		}

		return k.InitiateDKGRound(ctx, true)
	}

	// Prune timed-out decrypt request registry entries every N blocks.
	// All nodes evaluate this condition identically, preserving consensus.
	if currentHeight%types.DecryptRequestRegistryCleanupInterval == 0 {
		if err := k.pruneTimedOutDecryptRequests(ctx, uint64(currentHeight)); err != nil {
			log.Error(ctx, "Failed to prune timed-out decrypt request registry", err)
		}
	}

	nextStage, shouldTransition := k.shouldTransitionStage(currentHeight, latestRound, params)
	if shouldTransition {
		// Update the stage of this round before emitting events
		latestRound.Stage = nextStage
		if err := k.setDKGNetwork(ctx, latestRound); err != nil {
			return err
		}

		// Emit appropriate events for stage transitions
		switch nextStage {
		case types.DKGStageRegistration:
			// round = DKGStageRegistration if either
			// 1. it's the initial (first) round, OR
			// 2. the active stage of the previous round has ended, so DKG needs to reshare deals
			return k.InitiateDKGRound(ctx, false)
		case types.DKGStageDealing:
			return k.BeginDealing(ctx, latestRound)
		case types.DKGStageFinalization:
			return k.BeginFinalization(ctx, latestRound)
		case types.DKGStageActive:
			return k.FinalizeDKGRound(ctx, latestRound)
		case types.DKGStageUnspecified:
			// This round should not happen since we always have a valid stage (1 to 5) and unspecified is stage 0
			return nil
		}
	}

	return nil
}

// hasPendingUpgradeActivation checks if there is a pending kernel upgrade that should be
// activated at the current block height. Returns the upgrade info if activation is due,
// or nil if no upgrade needs activation.
func (k *Keeper) hasPendingUpgradeActivation(ctx context.Context, currentHeight int64) (*types.KernelUpgradeInfo, error) {
	upgradeInfo, err := k.GetPendingUpgrade(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check pending upgrade")
	}

	if upgradeInfo == nil || currentHeight < upgradeInfo.ActivationHeight {
		return nil, nil
	}

	return upgradeInfo, nil
}
