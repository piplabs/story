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

		return k.InitiateDKGRound(ctx)
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
			return k.InitiateDKGRound(ctx)
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
