package keeper

import (
	"context"

	"github.com/piplabs/story/client/x/mint/types"
)

func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	return k.Params.Get(ctx)
}

func (k Keeper) SetParams(ctx context.Context, value types.Params) error {
	if err := value.Validate(); err != nil {
		return err
	}

	return k.Params.Set(ctx, value)
}
