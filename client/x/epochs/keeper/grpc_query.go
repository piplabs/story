package keeper

import (
	"context"
	"errors"

	"github.com/piplabs/story/client/x/epochs/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Querier{}

// Querier defines a wrapper around the x/epochs keeper providing gRPC method
// handlers.
type Querier struct {
	Keeper
}

// NewQuerier initializes new querier.
func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

// GetEpochInfos provide running epochInfos.
func (q Querier) GetEpochInfos(ctx context.Context, _ *types.GetEpochInfosRequest) (*types.GetEpochInfosResponse, error) {
	epochs, err := q.Keeper.AllEpochInfos(ctx)
	return &types.GetEpochInfosResponse{
		Epochs: epochs,
	}, err
}

// GetEpochInfo provide epoch info of specified identifier.
func (q Querier) GetEpochInfo(ctx context.Context, req *types.GetEpochInfoRequest) (*types.GetEpochInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Identifier == "" {
		return nil, status.Error(codes.InvalidArgument, "identifier is empty")
	}

	info, err := q.Keeper.GetEpochInfo(ctx, req.Identifier)
	if err != nil {
		return nil, errors.New("not available identifier")
	}

	return &types.GetEpochInfoResponse{
		Epoch: info,
	}, nil
}
