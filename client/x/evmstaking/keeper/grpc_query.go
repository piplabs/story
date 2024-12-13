package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/piplabs/story/client/collections"
	"github.com/piplabs/story/client/x/evmstaking/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(ctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

// GetWithdrawalQueue returns the withdrawal queue in pagination.
func (k Keeper) GetWithdrawalQueue(ctx context.Context, request *types.QueryGetWithdrawalQueueRequest) (*types.QueryGetWithdrawalQueueResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	wqStore := prefix.NewStore(store, append(types.WithdrawalQueueKey, collections.QueueElementsPrefixSuffix)) // withdrawal queue store

	withdrawals, pageResp, err := query.GenericFilteredPaginate(k.cdc, wqStore, request.Pagination, func(_ []byte, wit *types.Withdrawal) (*types.Withdrawal, error) {
		return wit, nil
	}, func() *types.Withdrawal {
		return &types.Withdrawal{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var ws []*types.Withdrawal
	for _, w := range withdrawals {
		ws = append(ws, &types.Withdrawal{
			CreationHeight:   w.CreationHeight,
			ExecutionAddress: w.ExecutionAddress,
			Amount:           w.Amount,
		})
	}

	return &types.QueryGetWithdrawalQueueResponse{Withdrawals: ws, Pagination: pageResp}, nil
}
