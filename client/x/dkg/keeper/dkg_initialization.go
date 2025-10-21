package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/piplabs/story/client/server/utils"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
	"strings"
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

func (k *Keeper) InitiateDKGRound(ctx context.Context) error {
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

	if err := k.setDKGNetwork(ctx, &dkgNetwork); err != nil {
		return err
	}

	log.Info(ctx, "Initiated new DKG round",
		"round", roundNum,
		"start_block", currentHeight,
	)

	if err := k.emitBeginDKGInitialization(ctx, &dkgNetwork); err != nil {
		return errors.Wrap(err, "failed to emit begin dkg initialization event")
	}

	if k.isDKGSvcEnabled {
		go k.handleDKGInitialization(ctx, &dkgNetwork)
	}

	return nil
}
