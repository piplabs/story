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
		log.Info(ctx, "[DKG] No active DKG round, starting the first round")

		return k.initiateDKGRound(ctx)
	}

	nextStage, shouldTransition := k.shouldTransitionStage(currentHeight, latestRound, params)
	if shouldTransition {
		// Update the stage of this round before emitting events
		latestRound.Stage = nextStage
		if err := k.SetDKGNetwork(ctx, latestRound); err != nil {
			return err
		}

		// Emit appropriate events for stage transitions
		switch nextStage {
		case types.DKGStageRegistration:
			// round = DKGStageRegistration if either
			// 1. it's the initial (first) round, OR
			// 2. the active stage of the previous round has ended, so DKG needs to reshare deals
			return k.initiateDKGRound(ctx)
		case types.DKGStageNetworkSet:
			// TODO: Update current round's DKG Network data's total and threshold (based on registrations from the active validators)
			// TODO: check if there's enough number of registrations to set the network (and start dealing).
			// Use a DKG module parameter (minDKGMemberAmount). If the amount is less, we need to restart the round.

			return k.emitBeginDKGNetworkSet(ctx, latestRound)
		case types.DKGStageDealing:
			return k.emitBeginDKGDealing(ctx, latestRound)
		case types.DKGStageFinalization:
			return k.emitBeginDKGFinalization(ctx, latestRound)
		case types.DKGStageActive:
			return k.emitDKGFinalized(ctx, latestRound)
		case types.DKGStageUnspecified:
			// This round should not happen since we always have a valid stage (1 to 5) and unspecified is stage 0
			return nil
		}
	}

	return nil
}
