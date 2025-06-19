package keeper

import (
	"context"

	"github.com/piplabs/story/client/x/evmengine/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = &Keeper{}

func (k *Keeper) GetPendingUpgrade(ctx context.Context, request *types.QueryGetPendingUpgradeRequest) (*types.QueryGetPendingUpgradeResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	pendingUpgrade, err := k.PendingUpgrade(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetPendingUpgradeResponse{
		Plan: pendingUpgrade,
	}, nil
}
