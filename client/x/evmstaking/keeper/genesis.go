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

	// The epoch duration must be less than or equal to the unbonding time.
	if err := k.validateEpochDuration(ctx, gs.Params.EpochIdentifier); err != nil {
		panic(err)
	}

	if err := k.WithdrawalQueue.Initialize(ctx); err != nil {
		log.Error(ctx, "InitGenesis.evmstaking not initialized", err)
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

		err = k.DelegatorMap.Set(ctx, delAddr.String(), evmAddr.String())
		if err != nil {
			return errors.Wrap(err, "set delegator map")
		}
	}

	// init message queue
	if err := k.MessageQueue.Initialize(ctx); err != nil {
		return errors.Wrap(err, "initialize message queue")
	}

	// InitEpochNumber
	if err := k.SetEpochNumber(ctx, gs.EpochNumber); err != nil {
		return errors.Wrap(err, "initialize epoch num")
	}

	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	epochNum, err := k.GetEpochNumber(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Params:      params,
		EpochNumber: epochNum,
	}
}

//nolint:revive // TODO: validate genesis
func (k Keeper) ValidateGenesis(gs *types.GenesisState) error {
	if err := types.ValidateMaxWithdrawalPerBlock(gs.Params.MaxWithdrawalPerBlock); err != nil {
		return err
	}

	if err := types.ValidateMaxSweepPerBlock(gs.Params.MaxSweepPerBlock, gs.Params.MaxWithdrawalPerBlock); err != nil {
		return err
	}

	if err := types.ValidateMinPartialWithdrawalAmount(gs.Params.MinPartialWithdrawalAmount); err != nil {
		return err
	}

	return types.ValidateEpochIdentifier(gs.Params.EpochIdentifier)
}

func (k Keeper) validateEpochDuration(ctx context.Context, epochIdentifier string) error {
	unbondingTime, err := k.stakingKeeper.UnbondingTime(ctx)
	if err != nil {
		return errors.Wrap(err, "get unbonding time of staking keeper")
	}

	epoch, err := k.epochsKeeper.GetEpochInfo(ctx, epochIdentifier)
	if err != nil {
		return errors.Wrap(err, "get epoch info from epochs keeper", "epoch_identifier", epochIdentifier)
	}

	if epoch.Duration > unbondingTime {
		return errors.New("epoch duration must be less than or equal to the unbonding time", "unbondind_time", unbondingTime, "epoch_duration", epoch.Duration)
	}

	return nil
}
