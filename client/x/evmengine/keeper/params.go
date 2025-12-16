package keeper

import (
	"context"

	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/errors"
)

// ExecutionBlockHash returns the genesis execution block hash.
func (k *Keeper) ExecutionBlockHash(ctx context.Context) ([]byte, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return params.ExecutionBlockHash, nil
}

// SetParams set parameters of evmengine keeper.
func (k *Keeper) SetParams(ctx context.Context, params types.Params) error {
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

// GetParams gets parameters of evmengine keeeper.
func (k *Keeper) GetParams(ctx context.Context) (params types.Params, err error) {
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
