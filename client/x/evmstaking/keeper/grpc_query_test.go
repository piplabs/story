package keeper_test

import (
	"context"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/server/utils"
	"github.com/piplabs/story/client/x/evmstaking/types"
)

func TestParams(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
	require.NoError(t, esk.SetParams(ctx, types.DefaultParams()))

	queryClient := createQueryClient(ctx, esk)

	req := &types.QueryParamsRequest{}
	res, err := queryClient.Params(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, types.DefaultParams(), res.GetParams())
}

func TestGetWithdrawalQueue(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
	require.NoError(t, esk.WithdrawalQueue.Initialize(ctx))

	queryClient := createQueryClient(ctx, esk)

	pageReq := &query.PageRequest{
		Key:        nil,
		Limit:      2,
		CountTotal: true,
	}
	req := &types.QueryGetWithdrawalQueueRequest{
		Pagination: pageReq,
	}

	// Query an empty queue
	res, err := queryClient.GetWithdrawalQueue(context.Background(), req)
	require.NoError(t, err)
	require.Empty(t, res.Withdrawals, "expected no withdrawals in the queue yet")

	// Prepare and add three withdrawals to the queue
	valAddr = "storyvaloper1hmjw3pvkjtndpg8wqppwdn8udd835qpaa6r6y0"
	valEVMAddr, err := utils.Bech32ValidatorAddressToEvmAddress(valAddr)
	require.NoError(t, err)
	evmAddr = common.HexToAddress("0x131D25EDE18178BAc9275b312001a63C081722d2")
	withdrawals := []types.Withdrawal{
		types.NewWithdrawal(1, evmAddr.String(), 100, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, valEVMAddr),
		types.NewWithdrawal(2, evmAddr.String(), 200, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, valEVMAddr),
		types.NewWithdrawal(3, evmAddr.String(), 300, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, valEVMAddr),
	}
	require.Len(t, withdrawals, 3)
	for _, w := range withdrawals {
		require.NoError(t, esk.AddWithdrawalToQueue(ctx, w))
	}

	// Query the first page of two withdrawals
	res, err = queryClient.GetWithdrawalQueue(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, res.Withdrawals, 2,
		"expected 2 withdrawals after first page query, but found %d", len(res.Withdrawals))

	// Query the next page for the remaining withdrawal
	nextPage := res.Pagination.NextKey
	require.NotNil(t, nextPage, "expected a next page key to be not nil")

	pageReq.Key = nextPage
	req = &types.QueryGetWithdrawalQueueRequest{
		Pagination: pageReq,
	}
	res, err = queryClient.GetWithdrawalQueue(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, res.Withdrawals, 1, "expected 1 withdrawal after second page query, but found %d", len(res.Withdrawals))
}

func TestGetRewardWithdrawalQueue(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
	require.NoError(t, esk.WithdrawalQueue.Initialize(ctx))

	queryClient := createQueryClient(ctx, esk)

	pageReq := &query.PageRequest{
		Key:        nil,
		Limit:      5,
		CountTotal: true,
	}
	req := &types.QueryGetRewardWithdrawalQueueRequest{
		Pagination: pageReq,
	}

	// Query an empty queue
	res, err := queryClient.GetRewardWithdrawalQueue(context.Background(), req)
	require.NoError(t, err)
	require.Empty(t, res.Withdrawals, "expected no withdrawals in the queue yet")

	// Prepare and add three withdrawals to the queue
	valAddr = "storyvaloper1hmjw3pvkjtndpg8wqppwdn8udd835qpaa6r6y0"
	valEVMAddr, err := utils.Bech32ValidatorAddressToEvmAddress(valAddr)
	require.NoError(t, err)
	evmAddr = common.HexToAddress("0x131D25EDE18178BAc9275b312001a63C081722d2")
	withdrawals := []types.Withdrawal{
		types.NewWithdrawal(1, evmAddr.String(), 100, types.WithdrawalType_WITHDRAWAL_TYPE_REWARD, valEVMAddr),
		types.NewWithdrawal(2, evmAddr.String(), 200, types.WithdrawalType_WITHDRAWAL_TYPE_REWARD, valEVMAddr),
		types.NewWithdrawal(3, evmAddr.String(), 300, types.WithdrawalType_WITHDRAWAL_TYPE_REWARD, valEVMAddr),
	}
	require.Len(t, withdrawals, 3)
	for _, w := range withdrawals {
		require.NoError(t, esk.AddRewardWithdrawalToQueue(ctx, w))
	}

	// Query the first page of two withdrawals
	res, err = queryClient.GetRewardWithdrawalQueue(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, res.Withdrawals, 3,
		"expected 3 withdrawals after first page query, but found %d", len(res.Withdrawals))
}

func TestGetOperatorAddress(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
	require.NoError(t, esk.WithdrawalQueue.Initialize(ctx))

	queryClient := createQueryClient(ctx, esk)

	delEVMAddr, err := utils.Bech32DelegatorAddressToEvmAddress(delAddr)
	require.NoError(t, err)
	evmAddr = common.HexToAddress("0x131D25EDE18178BAc9275b312001a63C081722d2")

	req := &types.QueryGetOperatorAddressRequest{
		Address: delAddr,
	}

	// Query an empty operator
	res, err := queryClient.GetOperatorAddress(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "", res.OperatorAddress, "expected empty operator address yet")

	// Set operator address
	require.NoError(t, esk.DelegatorOperatorAddress.Set(ctx, delAddr, delEVMAddr))

	// Query the operator
	res, err = queryClient.GetOperatorAddress(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, delEVMAddr, res.OperatorAddress, "operator address mismatch")
}

func TestGetWithdrawAddress(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
	require.NoError(t, esk.WithdrawalQueue.Initialize(ctx))

	queryClient := createQueryClient(ctx, esk)

	delEVMAddr, err := utils.Bech32DelegatorAddressToEvmAddress(delAddr)
	require.NoError(t, err)
	evmAddr = common.HexToAddress("0x131D25EDE18178BAc9275b312001a63C081722d2")

	req := &types.QueryGetWithdrawAddressRequest{
		Address: delAddr,
	}

	// Query an empty withdraw address
	res, err := queryClient.GetWithdrawAddress(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "", res.WithdrawAddress, "expected empty withdraw address yet")

	// Set withdraw address
	require.NoError(t, esk.DelegatorWithdrawAddress.Set(ctx, delAddr, delEVMAddr))

	// Query the withdraw address
	res, err = queryClient.GetWithdrawAddress(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, delEVMAddr, res.WithdrawAddress, "withdraw address mismatch")
}

func TestGetRewardAddress(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
	require.NoError(t, esk.WithdrawalQueue.Initialize(ctx))

	queryClient := createQueryClient(ctx, esk)

	delEVMAddr, err := utils.Bech32DelegatorAddressToEvmAddress(delAddr)
	require.NoError(t, err)
	evmAddr = common.HexToAddress("0x131D25EDE18178BAc9275b312001a63C081722d2")

	req := &types.QueryGetRewardAddressRequest{
		Address: delAddr,
	}

	// Query an empty reward address
	res, err := queryClient.GetRewardAddress(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "", res.RewardAddress, "expected empty reward address yet")

	// Set reward address
	require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, delAddr, delEVMAddr))

	// Query the reward address
	res, err = queryClient.GetRewardAddress(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, delEVMAddr, res.RewardAddress, "withdraw address mismatch")
}
