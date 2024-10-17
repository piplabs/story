package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) IsSingularity(ctx context.Context) (bool, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeader().Height

	if blockHeight < int64(params.SingularityHeight) {
		return false, nil
	}

	return true, nil
}
