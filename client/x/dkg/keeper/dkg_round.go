package keeper

import (
	"context"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

func (*Keeper) shouldTransitionStage(currentHeight int64, dkgNetwork *types.DKGNetwork, params types.Params) (types.DKGStage, bool) {
	currentStage := dkgNetwork.Stage
	elapsed := currentHeight - dkgNetwork.StartBlock

	registrationEnd := int64(params.RegistrationPeriod)
	dealingEnd := registrationEnd + int64(params.DealingPeriod)
	finalizationEnd := dealingEnd + int64(params.FinalizationPeriod)
	activeEnd := finalizationEnd + int64(params.ActivePeriod)

	// in switch, we check if the elapsed time is greater than the end of the current stage
	switch currentStage {
	case types.DKGStageRegistration:
		if elapsed >= registrationEnd {
			return types.DKGStageDealing, true
		}
	case types.DKGStageDealing:
		if elapsed >= dealingEnd {
			return types.DKGStageFinalization, true
		}
	case types.DKGStageFinalization:
		if elapsed >= finalizationEnd {
			return types.DKGStageActive, true
		}
	case types.DKGStageActive:
		if elapsed >= activeEnd {
			// Round has ended, should initiate new round (resharing)
			return types.DKGStageRegistration, true
		}
	case types.DKGStageUnspecified:
		return types.DKGStageUnspecified, false
	}

	return currentStage, false
}

func (k *Keeper) SkipToNextRound(ctx context.Context, currentRound *types.DKGNetwork) error {
	// Mark the current round as failed
	currentRound.Stage = types.DKGStageFailed
	if err := k.setDKGNetwork(ctx, currentRound); err != nil {
		return errors.Wrap(err, "failed to mark the current round as failed")
	}

	return k.InitiateDKGRound(ctx)
}
