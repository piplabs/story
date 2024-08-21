package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/storyprotocol/iliad/lib/errors"
)

// Modified from https://github.com/cosmos/cosmos-sdk/blob/v0.50.7/x/staking/keeper/delegation.go#L521
func (k Keeper) GetMatureUnbondedDelegations(ctx context.Context) (matureUnbonds []stypes.DVPair, err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currTime := sdkCtx.BlockHeader().Time

	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator, err := (k.stakingKeeper).UBDQueueIterator(ctx, currTime)
	if err != nil {
		return nil, err
	}
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := stypes.DVPairs{}
		value := unbondingTimesliceIterator.Value()
		if err = k.cdc.Unmarshal(value, &timeslice); err != nil {
			return matureUnbonds, errors.Wrap(err, "failed to unmarshal unbonding timeslice")
		}

		matureUnbonds = append(matureUnbonds, timeslice.Pairs...)
	}

	return matureUnbonds, nil
}
