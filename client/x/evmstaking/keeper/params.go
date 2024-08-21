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

func (k Keeper) SetNextValidatorSweepIndex(ctx context.Context, nextValIndex sdk.IntProto) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&nextValIndex)
	if err != nil {
		return errors.Wrap(err, "marshal next validator sweep index")
	}

	err = store.Set(types.NextValidatorSweepIndexKey, bz)
	if err != nil {
		return errors.Wrap(err, "set next validator sweep index")
	}

	return nil
}

func (k Keeper) GetNextValidatorSweepIndex(ctx context.Context) (nextValIndex sdk.IntProto, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.NextValidatorSweepIndexKey)
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
