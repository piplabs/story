package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"
	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"

	addcollections "github.com/piplabs/story/client/collections"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
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

	withdrawals, pageResp, err := k.paginateWithdrawalQueue(ctx, types.WithdrawalQueueKey, request.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetWithdrawalQueueResponse{Withdrawals: withdrawals, Pagination: pageResp}, nil
}

// GetRewardWithdrawalQueue returns the withdrawal queue in pagination.
func (k Keeper) GetRewardWithdrawalQueue(ctx context.Context, request *types.QueryGetRewardWithdrawalQueueRequest) (*types.QueryGetRewardWithdrawalQueueResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	withdrawals, pageResp, err := k.paginateWithdrawalQueue(ctx, types.RewardWithdrawalQueueKey, request.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetRewardWithdrawalQueueResponse{Withdrawals: withdrawals, Pagination: pageResp}, nil
}

func (k Keeper) GetOperatorAddress(ctx context.Context, request *types.QueryGetOperatorAddressRequest) (*types.QueryGetOperatorAddressResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	operatorAddress, err := k.DelegatorOperatorAddress.Get(ctx, request.Address)
	if errors.Is(err, collections.ErrNotFound) {
		return &types.QueryGetOperatorAddressResponse{OperatorAddress: ""}, nil
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetOperatorAddressResponse{OperatorAddress: operatorAddress}, nil
}

func (k Keeper) GetWithdrawAddress(ctx context.Context, request *types.QueryGetWithdrawAddressRequest) (*types.QueryGetWithdrawAddressResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	withdrawAddress, err := k.DelegatorWithdrawAddress.Get(ctx, request.Address)
	if errors.Is(err, collections.ErrNotFound) {
		return &types.QueryGetWithdrawAddressResponse{WithdrawAddress: ""}, nil
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetWithdrawAddressResponse{WithdrawAddress: withdrawAddress}, nil
}

func (k Keeper) GetRewardAddress(ctx context.Context, request *types.QueryGetRewardAddressRequest) (*types.QueryGetRewardAddressResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	rewardAddress, err := k.DelegatorRewardAddress.Get(ctx, request.Address)
	if errors.Is(err, collections.ErrNotFound) {
		return &types.QueryGetRewardAddressResponse{RewardAddress: ""}, nil
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetRewardAddressResponse{RewardAddress: rewardAddress}, nil
}

func (k Keeper) paginateWithdrawalQueue(ctx context.Context, queueKey []byte, pagination *query.PageRequest) ([]*types.Withdrawal, *query.PageResponse, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	wqStore := prefix.NewStore(store, append(queueKey, addcollections.QueueElementsPrefixSuffix)) // withdrawal queue store

	withdrawals, pageResp, err := query.GenericFilteredPaginate(k.cdc, wqStore, pagination, func(_ []byte, wit *types.Withdrawal) (*types.Withdrawal, error) {
		return wit, nil
	}, func() *types.Withdrawal {
		return &types.Withdrawal{}
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to paginate withdrawal queue")
	}

	var ws []*types.Withdrawal
	for _, w := range withdrawals {
		ws = append(ws, &types.Withdrawal{
			CreationHeight:   w.CreationHeight,
			ExecutionAddress: w.ExecutionAddress,
			Amount:           w.Amount,
			WithdrawalType:   w.WithdrawalType,
			ValidatorAddress: w.ValidatorAddress,
		})
	}

	return ws, pageResp, nil
}
