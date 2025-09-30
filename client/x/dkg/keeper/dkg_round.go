package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/lib/cast"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/server/utils"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// GetActiveValidators returns the bonded validators' EVM addresses excluding jailed validators.
func (k *Keeper) GetActiveValidators(ctx context.Context) ([]string, error) {
	validators, err := k.stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all validators")
	}

	var bondedValidators []string
	for _, val := range validators {
		if val.IsBonded() && !val.IsJailed() {
			evmOperatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(val.OperatorAddress)
			if err != nil {
				return nil, errors.Wrap(err, "failed to convert to evm address", "operator_addr", val.OperatorAddress)
			}

			bondedValidators = append(bondedValidators, strings.ToLower(evmOperatorAddress))
		}
	}

	return bondedValidators, nil
}

func (*Keeper) shouldTransitionStage(currentHeight int64, dkgNetwork *types.DKGNetwork, params types.Params) (types.DKGStage, bool) {
	currentStage := dkgNetwork.Stage
	elapsed := currentHeight - dkgNetwork.StartBlock

	registrationEnd := int64(params.RegistrationPeriod)
	networkSetEnd := registrationEnd + int64(params.NetworkSetPeriod)
	dealingEnd := networkSetEnd + int64(params.DealingPeriod)
	finalizationEnd := dealingEnd + int64(params.FinalizationPeriod)
	activeEnd := finalizationEnd + int64(params.ActivePeriod)

	// in switch, we check if the elapsed time is greater than the end of the current stage
	//nolint:exhaustive // skip `types.DKGStageNetworkSet` for now
	switch currentStage {
	case types.DKGStageRegistration:
		if elapsed >= registrationEnd {
			return types.DKGStageNetworkSet, true
		}
	// NOTE: DKGStageNetworkSet
	case types.DKGStageNetworkSet:
		if elapsed >= networkSetEnd {
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
		ActiveValSet: activeValidators, // list of active validators' EVM addresses
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
	)

	return k.emitBeginDKGInitialization(ctx, &dkgNetwork)
}

// getVerifiedDKGValidators returns the count of verified DKG validators (those who are participating and not invalidated).
func (k *Keeper) getVerifiedDKGValidators(ctx context.Context, mrenclave []byte, round uint32) (uint32, error) {
	mrenclave32, err := cast.ToBytes32(mrenclave)
	if err != nil {
		return 0, errors.Wrap(err, "failed to cast to bytes32")
	}

	// Get registrations with status VERIFIED
	verifiedRegs, err := k.getDKGRegistrationsByStatus(ctx, mrenclave32, round, types.DKGRegStatusVerified)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get verified registrations")
	}

	total := uint32(len(verifiedRegs))

	return total, nil
}

func (k *Keeper) updateDKGNetworkTotalAndThreshold(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	verifiedCount, err := k.getVerifiedDKGValidators(ctx, dkgNetwork.Mrenclave, dkgNetwork.Round)
	if err != nil {
		return errors.Wrap(err, "failed to get verified DKG validators count")
	}

	dkgNetwork.Total = verifiedCount
	dkgNetwork.Threshold = k.calculateThreshold(verifiedCount)

	return k.SetDKGNetwork(ctx, dkgNetwork)
}

func (k *Keeper) finalizeDKGRound(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	dkgNetwork.Stage = types.DKGStageActive

	log.Info(ctx, "DKG network setup completed. Now active.", "round", dkgNetwork.Round, "mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave))

	return k.SetDKGNetwork(ctx, dkgNetwork)
}
