package keeper

import (
	"context"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
)

func (k Keeper) MaxWithdrawalPerBlock(ctx context.Context) (uint32, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return 0, err
	}

	return params.MaxWithdrawalPerBlock, nil
}

func (k Keeper) MaxSweepPerBlock(ctx context.Context) (uint32, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return 0, err
	}

	return params.MaxSweepPerBlock, nil
}

func (k Keeper) MinPartialWithdrawalAmount(ctx context.Context) (uint64, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return 0, err
	}

	return params.MinPartialWithdrawalAmount, nil
}

// This method performs no validation of the parameters.
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return errors.Wrap(err, "marshal params")
	}

	err = store.Set(types.ParamsKey, bz)
	if err != nil {
		return errors.Wrap(err, "set params")
	}

	return nil
}

func (k Keeper) GetParams(ctx context.Context) (params types.Params, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ParamsKey)
	if err != nil {
		return params, errors.Wrap(err, "get params")
	}

	if bz == nil {
		return params, nil
	}

	err = k.cdc.Unmarshal(bz, &params)
	if err != nil {
		return params, errors.Wrap(err, "unmarshal params")
	}

	return params, nil
}

func (k Keeper) SetValidatorSweepIndex(ctx context.Context, nextValIndex sdk.IntProto, nextValDelIndex sdk.IntProto) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&types.ValidatorSweepIndex{
		NextValIndex:    nextValIndex.Int.Uint64(),
		NextValDelIndex: nextValDelIndex.Int.Uint64(),
	})
	if err != nil {
		return errors.Wrap(err, "marshal validator sweep index")
	}

	err = store.Set(types.ValidatorSweepIndexKey, bz)
	if err != nil {
		return errors.Wrap(err, "set validator sweep index")
	}

	return nil
}

func (k Keeper) GetValidatorSweepIndex(ctx context.Context) (nextValIndex sdk.IntProto, nextValDelIndex sdk.IntProto, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ValidatorSweepIndexKey)
	if err != nil {
		return nextValIndex, nextValDelIndex, errors.Wrap(err, "get validator sweep index")
	}

	if bz == nil {
		return sdk.IntProto{Int: math.NewInt(0)}, sdk.IntProto{Int: math.NewInt(0)}, nil
	}

	var sweepIndex types.ValidatorSweepIndex
	err = k.cdc.Unmarshal(bz, &sweepIndex)
	if err != nil {
		return nextValIndex, nextValDelIndex, errors.Wrap(err, "unmarshal validator sweep index")
	}

	return nextValIndex, nextValDelIndex, nil
}

func (k Keeper) GetOldValidatorSweepIndex(ctx context.Context) (nextValIndex sdk.IntProto, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ValidatorSweepIndexKey)
	if err != nil {
		return nextValIndex, errors.Wrap(err, "get next validator sweep index")
	}

	if bz == nil {
		return sdk.IntProto{Int: math.NewInt(0)}, nil
	}

	err = k.cdc.Unmarshal(bz, &nextValIndex)
	if err != nil {
		return nextValIndex, errors.Wrap(err, "unmarshal next validator sweep index")
	}

	return nextValIndex, nil
}
