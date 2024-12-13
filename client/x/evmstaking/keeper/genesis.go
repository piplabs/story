package keeper

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
)

func (k Keeper) InitGenesis(ctx context.Context, gs *types.GenesisState) error {
	if err := k.ValidateGenesis(gs); err != nil {
		return err
	}

	if err := k.SetParams(ctx, gs.Params); err != nil {
		return err
	}

	if err := k.SetValidatorSweepIndex(ctx, gs.ValidatorSweepIndex); err != nil {
		return err
	}

	if err := k.WithdrawalQueue.Initialize(ctx); err != nil {
		log.Error(ctx, "InitGenesis.evmstaking withdrawal queue not initialized", err)
		return err
	}
	if err := k.RewardWithdrawalQueue.Initialize(ctx); err != nil {
		log.Error(ctx, "InitGenesis.evmstaking reward withdrawal queue not initialized", err)
		return err
	}
	vals, err := k.stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return err
	}
	for _, v := range vals {
		pk, ok := v.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
		if !ok {
			return err
		}
		evmAddr, err := k1util.CosmosPubkeyToEVMAddress(pk.Bytes())
		if err != nil {
			return err
		}

		delegatorPubkey, err := k1util.PubKeyBytesToCosmos(pk.Bytes())
		if err != nil {
			return err
		}
		delAddr := sdk.AccAddress(delegatorPubkey.Address().Bytes())

		log.Debug(ctx, "InitGenesis.evmstaking validator",
			"validator", v.GetOperator(),
			"val_op", v.OperatorAddress,
			"pk", pk.String(),
			"pk_addr", pk.Address().String(),
			"evm_addr", evmAddr.String(),
			"del_addr", delAddr.String(),
		)

		if err = k.DelegatorWithdrawAddress.Set(ctx, delAddr.String(), evmAddr.String()); err != nil {
			return errors.Wrap(err, "set delegator withdraw address map")
		}
		if err = k.DelegatorRewardAddress.Set(ctx, delAddr.String(), evmAddr.String()); err != nil {
			return errors.Wrap(err, "set delegator reward address map")
		}
	}

	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	validatorSweepIndex, err := k.GetValidatorSweepIndex(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Params:              params,
		ValidatorSweepIndex: validatorSweepIndex,
	}
}

func (Keeper) ValidateGenesis(gs *types.GenesisState) error {
	if err := gs.Params.Validate(); err != nil {
		return errors.Wrap(err, "validate genesis state params")
	}

	return nil
}
