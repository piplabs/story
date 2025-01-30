package keeper_test

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/server/utils"
	"github.com/piplabs/story/client/x/evmstaking/types"
)

func (s *TestSuite) TestGetWithdrawalQueue() {
	require := s.Require()
	ctx, keeper, queryClient := s.Ctx, s.EVMStakingKeeper, s.queryClient
	require.NoError(keeper.WithdrawalQueue.Initialize(ctx))

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
	require.NoError(err)
	require.Equal(0, len(res.Withdrawals), "expected no withdrawals in the queue yet")

	// Prepare and add three withdrawals to the queue
	delAddr = "story1hmjw3pvkjtndpg8wqppwdn8udd835qpan4hm0y"
	valAddr = "storyvaloper1hmjw3pvkjtndpg8wqppwdn8udd835qpaa6r6y0"
	valEVMAddr, err := utils.Bech32ValidatorAddressToEvmAddress(valAddr)
	require.NoError(err)
	evmAddr = common.HexToAddress("0x131D25EDE18178BAc9275b312001a63C081722d2")
	withdrawals := []types.Withdrawal{
		types.NewWithdrawal(1, evmAddr.String(), 100, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, valEVMAddr),
		types.NewWithdrawal(2, evmAddr.String(), 200, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, valEVMAddr),
		types.NewWithdrawal(3, evmAddr.String(), 300, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, valEVMAddr),
	}
	require.Len(withdrawals, 3)
	for _, w := range withdrawals {
		err = keeper.AddWithdrawalToQueue(ctx, w)
		require.NoError(err)
	}

	// Query the first page of two withdrawals
	res, err = queryClient.GetWithdrawalQueue(context.Background(), req)
	require.NoError(err)
	require.Equal(2, len(res.Withdrawals),
		"expected 2 withdrawals after first page query, but found %d", len(res.Withdrawals))

	// Query the next page for the remaining withdrawal
	nextPage := res.Pagination.NextKey
	require.NotNil(nextPage, "expected a next page key to be not nil")

	pageReq.Key = nextPage
	req = &types.QueryGetWithdrawalQueueRequest{
		Pagination: pageReq,
	}
	res, err = queryClient.GetWithdrawalQueue(context.Background(), req)
	require.NoError(err)
	require.Equal(1, len(res.Withdrawals), "expected 1 withdrawal after second page query, but found %d", len(res.Withdrawals))
}
