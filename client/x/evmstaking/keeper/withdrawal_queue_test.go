package keeper_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmstaking/keeper"
	"github.com/piplabs/story/client/x/evmstaking/types"
)

var (
	delAddr     = "story1hmjw3pvkjtndpg8wqppwdn8udd835qpan4hm0y"
	valAddr     = "storyvaloper1hmjw3pvkjtndpg8wqppwdn8udd835qpaa6r6y0"
	valEVMAddr  = "0xbee4e8859692e6d0a0ee0042e6ccfc6b4f1a003d"
	evmAddr     = common.HexToAddress("0x131D25EDE18178BAc9275b312001a63C081722d2")
	withdrawals = []types.Withdrawal{
		types.NewWithdrawal(1, evmAddr.String(), 100, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, valEVMAddr),
		types.NewWithdrawal(2, evmAddr.String(), 200, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, valEVMAddr),
		types.NewWithdrawal(3, evmAddr.String(), 300, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, valEVMAddr),
	}
)

func initWithdrawalQueue(t *testing.T, ctx sdk.Context, esk *keeper.Keeper) {
	t.Helper()

	require.NoError(t, esk.WithdrawalQueue.Initialize(ctx))
	require.Equal(t, uint64(0), esk.WithdrawalQueue.Len(ctx))
}

func addWithdrawals(t *testing.T, ctx sdk.Context, esk *keeper.Keeper, withdrawals []types.Withdrawal) {
	t.Helper()

	for _, w := range withdrawals {
		require.NoError(t, esk.AddWithdrawalToQueue(ctx, w))
	}
}

func TestAddWithdrawalToQueue(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

	// initialize withdrawal queue
	initWithdrawalQueue(t, ctx, esk)

	// Add a withdrawal to the queue
	withdrawal := types.NewWithdrawal(1, evmAddr.String(), 100, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, valAddr)
	addWithdrawals(t, ctx, esk, []types.Withdrawal{withdrawal})

	// Check the withdrawal is in the queue
	require.Equal(t, uint64(1), esk.WithdrawalQueue.Len(ctx))
	elem, err := esk.WithdrawalQueue.Get(ctx, 0)
	require.NoError(t, err)
	require.Equal(t, withdrawal, elem)
}

func TestDequeueEligibleWithdrawals(t *testing.T) {
	tcs := []struct {
		name        string
		maxDequeue  uint32
		expectedLen int
		expected    []types.Withdrawal
	}{
		{
			name:        "Dequeue 1 withdrawal (have: 3, cap: 1)",
			maxDequeue:  1,
			expectedLen: 1,
			expected:    withdrawals[:1],
		},
		{
			name:        "Dequeue 2 withdrawals (have: 3, cap: 2)",
			maxDequeue:  2,
			expectedLen: 2,
			expected:    withdrawals[:2],
		},
		{
			name:        "Dequeue all withdrawals (have: 3, cap: 10)",
			maxDequeue:  10,
			expectedLen: len(withdrawals),
			expected:    withdrawals,
		},
		{
			name:        "Dequeue with empty queue (have: 0, cap: 3)",
			maxDequeue:  3,
			expectedLen: 0,
			expected:    []types.Withdrawal{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
			initWithdrawalQueue(t, ctx, esk)

			if !strings.Contains(tc.name, "Dequeue with empty queue") {
				addWithdrawals(t, ctx, esk, withdrawals)
			}

			queueLen := esk.WithdrawalQueue.Len(ctx)

			// Dequeue the withdrawals
			result, err := esk.DequeueEligibleWithdrawals(ctx, tc.maxDequeue)
			require.NoError(t, err)
			require.Len(t, result, tc.expectedLen)

			// Check Queue length is decreased by the number of dequeued withdrawals
			require.Equal(t,
				queueLen-uint64(tc.expectedLen),
				esk.WithdrawalQueue.Len(ctx),
			)

			// Validate the content of the dequeued withdrawals
			for i, w := range result {
				require.Equal(t, tc.expected[i].ExecutionAddress, w.Address.String())
				require.Equal(t, tc.expected[i].Amount, w.Amount)
			}
		})
	}
}

func TestPeekEligibleWithdrawals(t *testing.T) {
	tcs := []struct {
		name        string
		maxPeek     uint32
		expectedLen int
		expected    []types.Withdrawal
	}{
		{
			name:        "Peek 1 withdrawal (have: 3, cap: 1)",
			maxPeek:     1,
			expectedLen: 1,
			expected:    withdrawals[:1],
		},
		{
			name:        "Peek 2 withdrawals (have: 3, cap: 2)",
			maxPeek:     2,
			expectedLen: 2,
			expected:    withdrawals[:2],
		},
		{
			name:        "Peek more than available (have: 3, cap: 10)",
			maxPeek:     10,
			expectedLen: len(withdrawals),
			expected:    withdrawals,
		},
		{
			name:        "Peek with empty queue (have: 0, cap: 3)",
			maxPeek:     3,
			expectedLen: 0,
			expected:    []types.Withdrawal{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
			initWithdrawalQueue(t, ctx, esk)
			if !strings.Contains(tc.name, "Peek with empty queue") {
				addWithdrawals(t, ctx, esk, withdrawals)
			}

			queueLen := esk.WithdrawalQueue.Len(ctx)

			// Peek the withdrawals
			result, err := esk.PeekEligibleWithdrawals(ctx, tc.maxPeek)
			require.NoError(t, err)
			require.Len(t, result, tc.expectedLen)

			// Peek does not change the queue length
			require.Equal(t, queueLen, esk.WithdrawalQueue.Len(ctx))

			// Validate the content of the dequeued withdrawals
			for i, w := range result {
				require.Equal(t, tc.expected[i].ExecutionAddress, w.Address.String())
				require.Equal(t, tc.expected[i].Amount, w.Amount)
			}
		})
	}
}

func TestGetAllWithdrawals(t *testing.T) {
	tests := []struct {
		name           string
		setupQueue     func(c sdk.Context, esk *keeper.Keeper)
		expectedLength int
		expectedResult []types.Withdrawal
	}{
		{
			name:           "Empty queue",
			setupQueue:     func(_ sdk.Context, esk *keeper.Keeper) {}, // No setup needed
			expectedLength: 0,
			expectedResult: nil,
		},
		{
			name: "Single withdrawal",
			setupQueue: func(c sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.AddWithdrawalToQueue(c, withdrawals[0]))
			},
			expectedLength: 1,
			expectedResult: []types.Withdrawal{withdrawals[0]},
		},
		{
			name: "Multiple withdrawals",
			setupQueue: func(c sdk.Context, esk *keeper.Keeper) {
				for _, w := range withdrawals {
					require.NoError(t, esk.AddWithdrawalToQueue(c, w))
				}
			},
			expectedLength: 3,
			expectedResult: withdrawals,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
			initWithdrawalQueue(t, ctx, esk)

			cachedCtx, _ := ctx.CacheContext()
			tt.setupQueue(cachedCtx, esk)

			result, err := esk.GetAllWithdrawals(cachedCtx)
			require.NoError(t, err)
			require.Len(t, result, tt.expectedLength)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGetWithdrawals(t *testing.T) {
	testCases := []struct {
		name        string
		maxRetrieve uint32
		expectedLen int
	}{
		{
			name:        "retrieve all",
			maxRetrieve: uint32(len(withdrawals)),
			expectedLen: len(withdrawals),
		},
		{
			name:        "retrieve 0",
			maxRetrieve: 0,
			expectedLen: 0,
		},
		{
			name:        "retrieve 1",
			maxRetrieve: 1,
			expectedLen: 1,
		},
		{
			name:        "retrieve 2",
			maxRetrieve: 2,
			expectedLen: 2,
		},
		{
			name:        "retrieve more than available",
			maxRetrieve: 10,
			expectedLen: len(withdrawals),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

			initWithdrawalQueue(t, ctx, esk)
			addWithdrawals(t, ctx, esk, withdrawals)

			result, err := esk.GetWithdrawals(ctx, tc.maxRetrieve)
			require.NoError(t, err)
			require.Len(t, result, tc.expectedLen)
			// check contents
			for i := range result[:tc.expectedLen] {
				require.Equal(t, withdrawals[i], result[i])
			}
		})
	}
}
