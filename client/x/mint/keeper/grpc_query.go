package keeper

import (
	"context"

	"github.com/piplabs/story/client/x/mint/types"
)

var _ types.QueryServer = Querier{}

func NewQuerier(k Keeper) types.QueryServer {
	return Querier{k}
}

type Querier struct {
	k Keeper
}

// Params returns params of the mint module.
func (q Querier) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := q.k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{Params: params}, nil
}
