package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) IsSingularity(ctx context.Context) (bool, error) {
	singularityHeight, err := k.stakingKeeper.GetSingularityHeight(ctx)
	if err != nil {
		return false, err
	}

	blockHeight := sdk.UnwrapSDKContext(ctx).BlockHeight()

	if blockHeight < int64(singularityHeight) {
		return true, nil
	}

	return false, nil
}
