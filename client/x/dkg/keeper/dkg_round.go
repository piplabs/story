package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// GetActiveValidators returns the bonded validators excluding jailed validators
func (k *Keeper) GetActiveValidators(ctx context.Context) ([]string, error) {
	validators, err := k.stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all validators")
	}

	var bondedValidators []string
	for _, val := range validators {
		if val.IsBonded() && !val.IsJailed() {
			bondedValidators = append(bondedValidators, val.OperatorAddress)
		}
	}

	return bondedValidators, nil
}

func (k *Keeper) shouldTransitionStage(currentHeight int64, dkgNetwork *types.DKGNetwork, params types.Params) (types.DKGStage, bool) {
	currentStage := dkgNetwork.Stage
	elapsed := currentHeight - dkgNetwork.StartBlock

	registrationEnd := int64(params.RegistrationPeriod)
	challengeEnd := registrationEnd + int64(params.ChallengePeriod)
	dealingEnd := challengeEnd + int64(params.DealingPeriod)
	finalizationEnd := dealingEnd + int64(params.FinalizationPeriod)
	activeEnd := finalizationEnd + int64(params.ActivePeriod)

	// in switch, we check if the elapsed time is greater than the end of the current stage
	switch currentStage {
	case types.DKGStageRegistration:
		if elapsed >= registrationEnd {
			return types.DKGStageChallenge, true
		}
	case types.DKGStageChallenge:
		if elapsed >= challengeEnd {
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
	}

	return currentStage, false
}

// getVerifiedDKGValidators returns the count of verified DKG validators (those who are participating and not invalidated)
func (k *Keeper) getVerifiedDKGValidators(ctx context.Context, mrenclave []byte, round uint32) (uint32, error) {
	// Get registrations with status VERIFIED
	verifiedRegs, err := k.GetDKGRegistrationsByStatus(ctx, mrenclave, round, types.DKGRegStatusVerified)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get verified registrations")
	}

	// Get registrations with status FINALIZED
	finalizedRegs, err := k.GetDKGRegistrationsByStatus(ctx, mrenclave, round, types.DKGRegStatusFinalized)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get finalized registrations")
	}

	total := uint32(len(verifiedRegs) + len(finalizedRegs))
	return total, nil
}

// updateDKGNetworkTotalAndThreshold updates the total and threshold for a DKG network after the challenge period
func (k *Keeper) updateDKGNetworkTotalAndThreshold(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	verifiedCount, err := k.getVerifiedDKGValidators(ctx, dkgNetwork.Mrenclave, dkgNetwork.Round)
	if err != nil {
		return errors.Wrap(err, "failed to get verified DKG validators count")
	}

	dkgNetwork.Total = verifiedCount
	dkgNetwork.Threshold = k.calculateThreshold(verifiedCount)

	return k.SetDKGNetwork(ctx, dkgNetwork)
}

func (k *Keeper) initiateDKGRound(ctx context.Context) error {
	currentHeight := sdk.UnwrapSDKContext(ctx).BlockHeight()

	activeValidators, err := k.GetActiveValidators(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get active validators")
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get params")
	}

	roundNum := k.getNextRoundNumber(ctx)

	dkgNetwork := types.DKGNetwork{
		Round:        roundNum,
		StartBlock:   currentHeight,
		Mrenclave:    params.Mrenclave, // latest TEE mrenclave
		ActiveValSet: activeValidators,
		Total:        0,
		Threshold:    0,
		Stage:        types.DKGStageRegistration,
	}

	if err := k.SetDKGNetwork(ctx, &dkgNetwork); err != nil {
		return err
	}

	log.Info(ctx, "Initiated new DKG round",
		"round", roundNum,
		"start_block", currentHeight,
		"validators", len(activeValidators),
		"threshold", dkgNetwork.Threshold,
	)

	return k.emitBeginDKGInitialization(ctx, &dkgNetwork)
}
