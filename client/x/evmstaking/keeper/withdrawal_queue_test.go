package keeper_test

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/evmstaking/types"
)

var (
	delAddr     = "story1hmjw3pvkjtndpg8wqppwdn8udd835qpan4hm0y"
	valAddr     = "storyvaloper1hmjw3pvkjtndpg8wqppwdn8udd835qpaa6r6y0"
	evmAddr     = common.HexToAddress("0x131D25EDE18178BAc9275b312001a63C081722d2")
	withdrawals = []types.Withdrawal{
		types.NewWithdrawal(1, evmAddr.String(), 100, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE),
		types.NewWithdrawal(2, evmAddr.String(), 200, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE),
		types.NewWithdrawal(3, evmAddr.String(), 300, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE),
	}
)

func (s *TestSuite) initQueue() {
	require := s.Require()
	err := s.EVMStakingKeeper.WithdrawalQueue.Initialize(s.Ctx)
	require.NoError(err)
	require.Equal(uint64(0), s.EVMStakingKeeper.WithdrawalQueue.Len(s.Ctx))
}

func (s *TestSuite) addWithdrawals(withdrawals []types.Withdrawal) {
	require := s.Require()
	for _, w := range withdrawals {
		err := s.EVMStakingKeeper.AddWithdrawalToQueue(s.Ctx, w)
		require.NoError(err)
	}
}

func (s *TestSuite) TestAddWithdrawalToQueue() {
	require := s.Require()
	s.initQueue()

	// Add a withdrawal to the queue
	withdrawal := types.NewWithdrawal(1, evmAddr.String(), 100, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE)
	err := s.EVMStakingKeeper.AddWithdrawalToQueue(s.Ctx, withdrawal)
	require.NoError(err)

	// Check the withdrawal is in the queue
	require.Equal(uint64(1), s.EVMStakingKeeper.WithdrawalQueue.Len(s.Ctx))
	elem, err := s.EVMStakingKeeper.WithdrawalQueue.Get(s.Ctx, 0)
	require.NoError(err)
	require.Equal(withdrawal, elem)
}

func (s *TestSuite) TestDequeueEligibleWithdrawals() {
	require := s.Require()
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
		s.Run(tc.name, func() {
			s.initQueue()

			if !strings.Contains(tc.name, "Dequeue with empty queue") {
				s.addWithdrawals(withdrawals)
			}

			// Set max dequeue parameter
			params, err := s.EVMStakingKeeper.GetParams(s.Ctx)
			require.NoError(err)
			params.MaxWithdrawalPerBlock = tc.maxDequeue
			err = s.EVMStakingKeeper.SetParams(s.Ctx, params)
			require.NoError(err)

			queueLen := s.EVMStakingKeeper.WithdrawalQueue.Len(s.Ctx)

			// Dequeue the withdrawals
			result, err := s.EVMStakingKeeper.DequeueEligibleWithdrawals(s.Ctx, uint32(tc.expectedLen))
			require.NoError(err)
			require.Equal(tc.expectedLen, len(result))

			// Check Queue length is decreased by the number of dequeued withdrawals
			require.Equal(
				queueLen-uint64(tc.expectedLen),
				s.EVMStakingKeeper.WithdrawalQueue.Len(s.Ctx),
			)

			// Validate the content of the dequeued withdrawals
			for i, w := range result {
				require.Equal(tc.expected[i].ExecutionAddress, w.Address.String())
				require.Equal(tc.expected[i].Amount, w.Amount)
			}
		})
	}
}

func (s *TestSuite) TestPeekEligibleWithdrawals() {
	require := s.Require()
	tcs := []struct {
		name        string
		maxDequeue  uint32
		expectedLen int
		expected    []types.Withdrawal
	}{
		{
			name:        "Peek 1 withdrawal (have: 3, cap: 1)",
			maxDequeue:  1,
			expectedLen: 1,
			expected:    withdrawals[:1],
		},
		{
			name:        "Peek 2 withdrawals (have: 3, cap: 2)",
			maxDequeue:  2,
			expectedLen: 2,
			expected:    withdrawals[:2],
		},
		{
			name:        "Peek more than available (have: 3, cap: 10)",
			maxDequeue:  10,
			expectedLen: len(withdrawals),
			expected:    withdrawals,
		},
		{
			name:        "Peek with empty queue (have: 0, cap: 3)",
			maxDequeue:  3,
			expectedLen: 0,
			expected:    []types.Withdrawal{},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			s.initQueue()
			if !strings.Contains(tc.name, "Peek with empty queue") {
				s.addWithdrawals(withdrawals)
			}

			// Set max dequeue parameter
			params, err := s.EVMStakingKeeper.GetParams(s.Ctx)
			require.NoError(err)
			params.MaxWithdrawalPerBlock = tc.maxDequeue
			err = s.EVMStakingKeeper.SetParams(s.Ctx, params)
			require.NoError(err)

			queueLen := s.EVMStakingKeeper.WithdrawalQueue.Len(s.Ctx)

			// Peek the withdrawals
			result, err := s.EVMStakingKeeper.PeekEligibleWithdrawals(s.Ctx, uint32(tc.expectedLen))
			require.NoError(err)
			require.Equal(tc.expectedLen, len(result))

			// Peek does not change the queue length
			require.Equal(queueLen, s.EVMStakingKeeper.WithdrawalQueue.Len(s.Ctx))

			// Validate the content of the dequeued withdrawals
			for i, w := range result {
				require.Equal(tc.expected[i].ExecutionAddress, w.Address.String())
				require.Equal(tc.expected[i].Amount, w.Amount)
			}
		})
	}
}

func (s *TestSuite) TestGetAllWithdrawals() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper
	tests := []struct {
		name           string
		setupQueue     func(c context.Context)
		expectedLength int
		expectedResult []types.Withdrawal
	}{
		{
			name:           "Empty queue",
			setupQueue:     func(_ context.Context) {}, // No setup needed
			expectedLength: 0,
			expectedResult: nil,
		},
		{
			name: "Single withdrawal",
			setupQueue: func(c context.Context) {
				require.NoError(keeper.AddWithdrawalToQueue(c, withdrawals[0]))
			},
			expectedLength: 1,
			expectedResult: []types.Withdrawal{withdrawals[0]},
		},
		{
			name: "Multiple withdrawals",
			setupQueue: func(c context.Context) {
				for _, w := range withdrawals {
					require.NoError(keeper.AddWithdrawalToQueue(c, w))
				}
			},
			expectedLength: 3,
			expectedResult: withdrawals,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.initQueue()
			cachedCtx, _ := ctx.CacheContext()
			tt.setupQueue(cachedCtx)

			result, err := s.EVMStakingKeeper.GetAllWithdrawals(cachedCtx)
			require.NoError(err)
			require.Equal(tt.expectedLength, len(result))
			require.Equal(tt.expectedResult, result)
		})
	}
}

func (s *TestSuite) TestGetWithdrawals() {
	require := s.Require()
	s.initQueue()
	s.addWithdrawals(withdrawals)

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
		s.Run(tc.name, func() {
			result, err := s.EVMStakingKeeper.GetWithdrawals(s.Ctx, tc.maxRetrieve)
			require.NoError(err)
			require.Equal(tc.expectedLen, len(result))
			// check contents
			for i := range result[:tc.expectedLen] {
				require.Equal(withdrawals[i], result[i])
			}
		})
	}
}
