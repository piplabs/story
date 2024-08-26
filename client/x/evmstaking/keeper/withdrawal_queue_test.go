package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/evmstaking/types"
)

var (
	delAddr     = "story1hmjw3pvkjtndpg8wqppwdn8udd835qpan4hm0y"
	valAddr     = "storyvaloper1hmjw3pvkjtndpg8wqppwdn8udd835qpaa6r6y0"
	evmAddr     = common.HexToAddress("0x131D25EDE18178BAc9275b312001a63C081722d2")
	withdrawals = []types.Withdrawal{
		types.NewWithdrawal(1, delAddr, valAddr, evmAddr.String(), 100),
		types.NewWithdrawal(2, delAddr, valAddr, evmAddr.String(), 200),
		types.NewWithdrawal(3, delAddr, valAddr, evmAddr.String(), 300),
	}
)

func (s *TestSuite) initQueue() {
	err := s.EVMStakingKeeper.WithdrawalQueue.Initialize(s.Ctx)
	s.NoError(err)
	s.Equal(uint64(0), s.EVMStakingKeeper.WithdrawalQueue.Len(s.Ctx))
}

func (s *TestSuite) addWithdrawals(withdrawals []types.Withdrawal) {
	for _, w := range withdrawals {
		err := s.EVMStakingKeeper.AddWithdrawalToQueue(s.Ctx, w)
		s.NoError(err)
	}
}

func (s *TestSuite) TestAddWithdrawalToQueue() {
	s.initQueue()

	// Add a withdrawal to the queue
	withdrawal := types.NewWithdrawal(1, delAddr, valAddr, evmAddr.String(), 100)
	err := s.EVMStakingKeeper.AddWithdrawalToQueue(s.Ctx, withdrawal)
	s.NoError(err)

	// Check the withdrawal is in the queue
	s.Equal(uint64(1), s.EVMStakingKeeper.WithdrawalQueue.Len(s.Ctx))
	elem, err := s.EVMStakingKeeper.WithdrawalQueue.Get(s.Ctx, 0)
	s.NoError(err)
	s.Equal(withdrawal, elem)
}

func (s *TestSuite) TestDequeueEligibleWithdrawals() {
	tcs := []struct {
		name        string
		maxDequeue  uint32
		expectedLen int
		expected    []types.Withdrawal
	}{
		{
			name:        "Dequeue 1 withdrawal",
			maxDequeue:  1,
			expectedLen: 1,
			expected:    withdrawals[:1],
		},
		{
			name:        "Dequeue 2 withdrawals",
			maxDequeue:  2,
			expectedLen: 2,
			expected:    withdrawals[:2],
		},
		{
			name:        "Dequeue more than available",
			maxDequeue:  10,
			expectedLen: len(withdrawals),
			expected:    withdrawals,
		},
		{
			name:        "Dequeue with empty queue",
			maxDequeue:  3,
			expectedLen: 0,
			expected:    []types.Withdrawal{},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			s.initQueue()

			if tc.name != "Dequeue with empty queue" {
				s.addWithdrawals(withdrawals)
			}

			// Set max dequeue parameter
			params, err := s.EVMStakingKeeper.GetParams(s.Ctx)
			s.NoError(err)
			params.MaxWithdrawalPerBlock = tc.maxDequeue
			err = s.EVMStakingKeeper.SetParams(s.Ctx, params)
			s.NoError(err)

			queueLen := s.EVMStakingKeeper.WithdrawalQueue.Len(s.Ctx)

			// Dequeue the withdrawals
			result, err := s.EVMStakingKeeper.DequeueEligibleWithdrawals(s.Ctx)
			s.NoError(err)
			s.Equal(tc.expectedLen, len(result))

			// Check Queue length is decreased by the number of dequeued withdrawals
			s.Equal(
				queueLen-uint64(tc.expectedLen),
				s.EVMStakingKeeper.WithdrawalQueue.Len(s.Ctx),
			)

			// Validate the content of the dequeued withdrawals
			for i, w := range result {
				s.Equal(tc.expected[i].ExecutionAddress, w.Address.String())
				s.Equal(tc.expected[i].Amount, w.Amount)
			}
		})
	}
}

func (s *TestSuite) TestPeekEligibleWithdrawals() {
	tcs := []struct {
		name        string
		maxDequeue  uint32
		expectedLen int
		expected    []types.Withdrawal
	}{
		{
			name:        "Peek 1 withdrawal",
			maxDequeue:  1,
			expectedLen: 1,
			expected:    withdrawals[:1],
		},
		{
			name:        "Peek 2 withdrawals",
			maxDequeue:  2,
			expectedLen: 2,
			expected:    withdrawals[:2],
		},
		{
			name:        "Peek more than available",
			maxDequeue:  10,
			expectedLen: len(withdrawals),
			expected:    withdrawals,
		},
		{
			name:        "Peek with empty queue",
			maxDequeue:  3,
			expectedLen: 0,
			expected:    []types.Withdrawal{},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			s.initQueue()
			if tc.name != "Peek with empty queue" {
				s.addWithdrawals(withdrawals)
			}

			// Set max dequeue parameter
			params, err := s.EVMStakingKeeper.GetParams(s.Ctx)
			s.NoError(err)
			params.MaxWithdrawalPerBlock = tc.maxDequeue
			err = s.EVMStakingKeeper.SetParams(s.Ctx, params)

			queueLen := s.EVMStakingKeeper.WithdrawalQueue.Len(s.Ctx)

			// Peek the withdrawals
			result, err := s.EVMStakingKeeper.PeekEligibleWithdrawals(s.Ctx)
			s.NoError(err)
			s.Equal(tc.expectedLen, len(result))

			// Peek does not change the queue length
			s.Equal(queueLen, s.EVMStakingKeeper.WithdrawalQueue.Len(s.Ctx))

			// Validate the content of the dequeued withdrawals
			for i, w := range result {
				s.Equal(tc.expected[i].ExecutionAddress, w.Address.String())
				s.Equal(tc.expected[i].Amount, w.Amount)
			}
		})
	}
}

func (s *TestSuite) TestGetAllWithdrawals() {
	s.initQueue()

	// Add a withdrawal to the queue
	withdrawal := types.NewWithdrawal(1, delAddr, valAddr, evmAddr.String(), 100)
	err := s.EVMStakingKeeper.AddWithdrawalToQueue(s.Ctx, withdrawal)
	s.NoError(err)

	// Get all withdrawals
	result, err := s.EVMStakingKeeper.GetAllWithdrawals(s.Ctx)
	s.NoError(err)
	s.Equal(1, len(result))
	s.Equal(withdrawal, result[0])
}

func (s *TestSuite) TestGetWithdrawals() {
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
			s.NoError(err)
			s.Equal(tc.expectedLen, len(result))
			// check contents
			for i := 0; i < tc.expectedLen; i++ {
				s.Equal(withdrawals[i], result[i])
			}
		})
	}
}
