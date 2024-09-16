package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
)

func (k Keeper) SetEpochNumber(ctx context.Context, epochNumber uint64) error {
	store := k.storeService.OpenKVStore(ctx)
	epochNumberBytes := sdk.Uint64ToBigEndian(epochNumber)
	if err := store.Set(types.EpochNumberKey, epochNumberBytes); err != nil {
		return errors.Wrap(err, "set epoch num")
	}

	return nil
}

func (k Keeper) GetEpochNumber(ctx context.Context) (uint64, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.EpochNumberKey)
	if err != nil {
		return 0, errors.Wrap(err, "get epoch number")
	}
	if bz == nil {
		return 0, errors.New("epoch number not found")
	}
	epochNumber := sdk.BigEndianToUint64(bz)

	return epochNumber, nil
}

// IsNextEpoch checks if the next epoch has started.
func (k Keeper) IsNextEpoch(ctx context.Context) (bool, error) {
	// get epoch identifier of evmstaking module
	params, err := k.GetParams(ctx)
	if err != nil {
		return false, errors.Wrap(err, "get evmstaking params")
	}
	epochIdentifier := params.EpochIdentifier

	// get epoch info from epochs keeper
	epoch, err := k.epochsKeeper.GetEpochInfo(ctx, epochIdentifier)
	if err != nil {
		return false, errors.Wrap(err, "get current epoch info from epochs keeper", "epoch_identifier", epochIdentifier)
	}

	// get current epoch number of evmstaking keeper
	currentEpoch, err := k.GetEpochNumber(ctx)
	if err != nil {
		return false, errors.Wrap(err, "get epoch num")
	}

	return epoch.CurrentEpoch > int64(currentEpoch), nil
}

func (k Keeper) IncCurEpochNumber(ctx context.Context) error {
	curEpochNumber, err := k.GetEpochNumber(ctx)
	if err != nil {
		return errors.Wrap(err, "get current epoch number")
	}

	inc := curEpochNumber + 1

	return k.SetEpochNumber(ctx, inc)
}
