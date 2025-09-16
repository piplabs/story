package keeper_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmstaking/keeper"
	"github.com/piplabs/story/client/x/evmstaking/types"
)

var rewardWithdrawals = []types.Withdrawal{
	types.NewWithdrawal(1, evmAddr.String(), 100, types.WithdrawalType_WITHDRAWAL_TYPE_REWARD, valEVMAddr),
	types.NewWithdrawal(2, evmAddr.String(), 200, types.WithdrawalType_WITHDRAWAL_TYPE_REWARD, valEVMAddr),
	types.NewWithdrawal(3, evmAddr.String(), 300, types.WithdrawalType_WITHDRAWAL_TYPE_REWARD, valEVMAddr),
}

func initRewardWithdrawalQueue(t *testing.T, ctx sdk.Context, esk *keeper.Keeper) {
	t.Helper()

	require.NoError(t, esk.RewardWithdrawalQueue.Initialize(ctx))
	require.Equal(t, uint64(0), esk.RewardWithdrawalQueue.Len(ctx))
}

func addRewardWithdrawals(t *testing.T, ctx sdk.Context, esk *keeper.Keeper, withdrawals []types.Withdrawal) {
	t.Helper()

	for _, w := range withdrawals {
		require.NoError(t, esk.AddRewardWithdrawalToQueue(ctx, w))
	}
}

func TestAddRewardWithdrawalToQueue(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

	// initialize reward withdrawal queue
	initRewardWithdrawalQueue(t, ctx, esk)

	// Add a reward withdrawal to the queue
	withdrawal := types.NewWithdrawal(1, evmAddr.String(), 100, types.WithdrawalType_WITHDRAWAL_TYPE_REWARD, valAddr)
	addRewardWithdrawals(t, ctx, esk, []types.Withdrawal{withdrawal})

	// Check the withdrawal is in the queue
	require.Equal(t, uint64(1), esk.RewardWithdrawalQueue.Len(ctx))
	elem, err := esk.RewardWithdrawalQueue.Get(ctx, 0)
	require.NoError(t, err)
	require.Equal(t, withdrawal, elem)
}

func TestDequeueEligibleRewardWithdrawals(t *testing.T) {
	tcs := []struct {
		name        string
		maxDequeue  uint32
		expectedLen int
		expected    []types.Withdrawal
	}{
		{
			name:        "Dequeue 1 reward withdrawal (have: 3, cap: 1)",
			maxDequeue:  1,
			expectedLen: 1,
			expected:    rewardWithdrawals[:1],
		},
		{
			name:        "Dequeue 2 reward withdrawals (have: 3, cap: 2)",
			maxDequeue:  2,
			expectedLen: 2,
			expected:    rewardWithdrawals[:2],
		},
		{
			name:        "Dequeue all reward withdrawals (have: 3, cap: 10)",
			maxDequeue:  10,
			expectedLen: len(rewardWithdrawals),
			expected:    rewardWithdrawals,
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
			initRewardWithdrawalQueue(t, ctx, esk)

			if !strings.Contains(tc.name, "Dequeue with empty queue") {
				addRewardWithdrawals(t, ctx, esk, rewardWithdrawals)
			}

			queueLen := esk.RewardWithdrawalQueue.Len(ctx)

			// Dequeue the reward withdrawals
			result, err := esk.DequeueEligibleRewardWithdrawals(ctx, tc.maxDequeue)
			require.NoError(t, err)
			require.Len(t, result, tc.expectedLen)

			// Check Queue length is decreased by the number of dequeued reward withdrawals
			require.Equal(t,
				queueLen-uint64(tc.expectedLen),
				esk.RewardWithdrawalQueue.Len(ctx),
			)

			// Validate the content of the dequeued reward withdrawals
			for i, w := range result {
				require.Equal(t, tc.expected[i].ExecutionAddress, w.Address.String())
				require.Equal(t, tc.expected[i].Amount, w.Amount)
			}
		})
	}
}

func TestPeekEligibleRewardWithdrawals(t *testing.T) {
	tcs := []struct {
		name        string
		maxPeek     uint32
		expectedLen int
		expected    []types.Withdrawal
	}{
		{
			name:        "Peek 1 reward withdrawal (have: 3, cap: 1)",
			maxPeek:     1,
			expectedLen: 1,
			expected:    rewardWithdrawals[:1],
		},
		{
			name:        "Peek 2 reward withdrawals (have: 3, cap: 2)",
			maxPeek:     2,
			expectedLen: 2,
			expected:    rewardWithdrawals[:2],
		},
		{
			name:        "Peek more than available (have: 3, cap: 10)",
			maxPeek:     10,
			expectedLen: len(rewardWithdrawals),
			expected:    rewardWithdrawals,
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
			initRewardWithdrawalQueue(t, ctx, esk)
			if !strings.Contains(tc.name, "Peek with empty queue") {
				addRewardWithdrawals(t, ctx, esk, rewardWithdrawals)
			}

			queueLen := esk.RewardWithdrawalQueue.Len(ctx)

			// Peek the reward withdrawals
			result, err := esk.PeekEligibleRewardWithdrawals(ctx, tc.maxPeek)
			require.NoError(t, err)
			require.Len(t, result, tc.expectedLen)

			// Peek does not change the queue length
			require.Equal(t, queueLen, esk.RewardWithdrawalQueue.Len(ctx))

			// Validate the content of the dequeued reward withdrawals
			for i, w := range result {
				require.Equal(t, tc.expected[i].ExecutionAddress, w.Address.String())
				require.Equal(t, tc.expected[i].Amount, w.Amount)
			}
		})
	}
}

func TestGetAllRewardWithdrawals(t *testing.T) {
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
			name: "Single reward withdrawal",
			setupQueue: func(c sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.AddRewardWithdrawalToQueue(c, rewardWithdrawals[0]))
			},
			expectedLength: 1,
			expectedResult: []types.Withdrawal{rewardWithdrawals[0]},
		},
		{
			name: "Multiple reward withdrawals",
			setupQueue: func(c sdk.Context, esk *keeper.Keeper) {
				for _, w := range rewardWithdrawals {
					require.NoError(t, esk.AddRewardWithdrawalToQueue(c, w))
				}
			},
			expectedLength: 3,
			expectedResult: rewardWithdrawals,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
			initRewardWithdrawalQueue(t, ctx, esk)

			cachedCtx, _ := ctx.CacheContext()
			tt.setupQueue(cachedCtx, esk)

			result, err := esk.GetAllRewardWithdrawals(cachedCtx)
			require.NoError(t, err)
			require.Len(t, result, tt.expectedLength)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGetRewardWithdrawals(t *testing.T) {
	testCases := []struct {
		name        string
		maxRetrieve uint32
		expectedLen int
	}{
		{
			name:        "retrieve all",
			maxRetrieve: uint32(len(rewardWithdrawals)),
			expectedLen: len(rewardWithdrawals),
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
			expectedLen: len(rewardWithdrawals),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

			initRewardWithdrawalQueue(t, ctx, esk)
			addRewardWithdrawals(t, ctx, esk, rewardWithdrawals)

			result, err := esk.GetRewardWithdrawals(ctx, tc.maxRetrieve)
			require.NoError(t, err)
			require.Len(t, result, tc.expectedLen)
			// check contents
			for i := range result[:tc.expectedLen] {
				require.Equal(t, rewardWithdrawals[i], result[i])
			}
		})
	}
}
