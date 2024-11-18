package keeper_test

//
// import (
//	"math"
//	"math/big"
//	"testing"
//
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
//	"github.com/ethereum/go-ethereum/common"
//	gethtypes "github.com/ethereum/go-ethereum/core/types"
//	"github.com/stretchr/testify/require"
//
//	"github.com/piplabs/story/client/x/evmstaking/types"
//	"github.com/piplabs/story/contracts/bindings"
//	"github.com/piplabs/story/lib/errors"
//	"github.com/piplabs/story/lib/k1util"
//
//	"go.uber.org/mock/gomock"
//)
//
//func (s *TestSuite) TestExpectedUnbondingWithdrawals() {
//	require := s.Require()
//	ctx, evmstakingKeeper, stakingKeeper, distrKeeper := s.Ctx, s.EVMStakingKeeper, s.StakingKeeper, s.DistrKeeper
//
//	pubKeys, accAddrs, valAddrs := createAddresses(3)
//	delAddr := accAddrs[0]
//	delPubKey := pubKeys[0]
//	evmDelAddr, err := k1util.CosmosPubkeyToEVMAddress(delPubKey.Bytes())
//	require.NoError(err)
//	valPubKey := pubKeys[1]
//	valAddr := valAddrs[1]
//	valPubKey2 := pubKeys[2]
//	valAddr2 := valAddrs[2]
//
//	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
//	s.setupValidatorAndDelegation(ctx, valPubKey, delPubKey, valAddr, delAddr, valTokens)
//	// set params as default
//	params := types.DefaultParams()
//	require.NoError(evmstakingKeeper.SetParams(ctx, params))
//	delRewardsAmt := params.MinPartialWithdrawalAmount + 100
//	delRewards := sdk.NewDecCoinsFromCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(delRewardsAmt)))
//
//	// Test cases for ExpectedRewardWithdrawals
//	tcs := []struct {
//		name           string
//		preRun         func(ctx sdk.Context)
//		expectedResult []types.Withdrawal
//		expectedError  string
//	}{
//		{
//			name: "pass",
//			preRun: func(ctx sdk.Context) {
//				// bankKeeper.EXPECT().GetBalance(ctx, gomock.Any(), sdk.DefaultBondDenom).AnyTimes()
//				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(ctx, gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil)
//				distrKeeper.EXPECT().IncrementValidatorPeriod(ctx, gomock.Any()).Return(uint64(0), nil)
//				distrKeeper.EXPECT().CalculateDelegationRewards(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(delRewards, nil)
//			},
//			expectedResult: []types.Withdrawal{
//				{
//					CreationHeight:   0,
//					ExecutionAddress: evmDelAddr.String(),
//					Amount:           delRewardsAmt,
//				},
//			},
//		},
//		{
//			name: "pass: val sweep index is out of bounds, so it should be reset to 0 which is the index of the first validator",
//			preRun: func(_ sdk.Context) {
//				validatorSweepIndex := types.ValidatorSweepIndex{
//					NextValIndex:    uint64(100),
//					NextValDelIndex: uint64(0),
//				}
//				require.NoError(evmstakingKeeper.SetValidatorSweepIndex(ctx, validatorSweepIndex))
//				// bankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), sdk.DefaultBondDenom).AnyTimes()
//				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil)
//				distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
//				distrKeeper.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delRewards, nil)
//			},
//			expectedResult: []types.Withdrawal{
//				{
//					CreationHeight:   0,
//					ExecutionAddress: evmDelAddr.String(),
//					Amount:           delRewardsAmt,
//				},
//			},
//		},
//		{
//			name: "fail: increment validator period",
//			preRun: func(_ sdk.Context) {
//				// bankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), sdk.DefaultBondDenom).AnyTimes()
//				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil)
//				distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), errors.New("failed to increment validator period"))
//			},
//			expectedError: "failed to increment validator period",
//		},
//		{
//			name: "fail: calculate delegation rewards",
//			preRun: func(_ sdk.Context) {
//				// bankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), sdk.DefaultBondDenom).AnyTimes()
//				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil)
//				distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
//				distrKeeper.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delRewards, errors.New("failed to calculate delegation rewards"))
//			},
//			expectedError: "failed to calculate delegation rewards",
//		},
//		{
//			name: "pass: multiple validators",
//			preRun: func(c sdk.Context) {
//				// bankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), sdk.DefaultBondDenom).AnyTimes()
//				s.setupValidatorAndDelegation(c, valPubKey2, delPubKey, valAddr2, delAddr, valTokens)
//				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil).Times(2)
//				distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil).Times(2)
//				distrKeeper.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delRewards, nil).Times(2)
//			},
//			expectedResult: []types.Withdrawal{
//				{
//					CreationHeight:   0,
//					ExecutionAddress: evmDelAddr.String(),
//					Amount:           delRewardsAmt,
//				},
//				{
//					CreationHeight:   0,
//					ExecutionAddress: evmDelAddr.String(),
//					Amount:           delRewardsAmt,
//				},
//			},
//		},
//		{
//			name: "pass: skip jailed validator",
//			preRun: func(c sdk.Context) {
//				s.setupValidatorAndDelegation(c, valPubKey2, delPubKey, valAddr2, delAddr, valTokens)
//				val, err := stakingKeeper.GetValidator(c, valAddr2)
//				require.NoError(err)
//				val.Jailed = true
//				require.NoError(stakingKeeper.SetValidator(c, val))
//				// bankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), sdk.DefaultBondDenom).AnyTimes()
//				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil)
//				distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
//				distrKeeper.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delRewards, nil)
//			},
//			expectedResult: []types.Withdrawal{
//				{
//					CreationHeight:   0,
//					ExecutionAddress: evmDelAddr.String(),
//					Amount:           delRewardsAmt,
//				},
//			},
//		},
//	}
//
//	for _, tc := range tcs {
//		s.Run(tc.name, func() {
//			cached, _ := ctx.CacheContext()
//			if tc.preRun != nil {
//				tc.preRun(cached)
//			}
//			_, unbondedEntries, err := stakingKeeper.EndBlockerWithUnbondedEntries(ctx)
//			require.NoError(err)
//			result, err := evmstakingKeeper.ExpectedUnbondingWithdrawals(cached, unbondedEntries)
//			if tc.expectedError != "" {
//				require.ErrorContains(err, tc.expectedError)
//			} else {
//				require.NoError(err)
//				isEqualWithdrawals(s.T(), tc.expectedResult, result)
//			}
//		})
//	}
//}
//
//// TODO: Test for ProcessUnbondingWithdrawals
//// TODO: Test for ProcessRewardWithdrawals
//// TODO: Test for ProcessEligibleRewardWithdrawal
//// TODO: Test for EnqueueRewardWithdrawal
//
//func (s *TestSuite) TestProcessWithdraw() {
//	require := s.Require()
//	ctx, keeper, stakingKeeper := s.Ctx, s.EVMStakingKeeper, s.StakingKeeper
//
//	pubKeys, accAddrs, valAddrs := createAddresses(4)
//	// delegator-1
//	delPubKey1 := pubKeys[0]
//	delAddr1 := accAddrs[0]
//	// validator
//	valPubKey := pubKeys[2]
//	valAddr := valAddrs[2]
//	// unknown pubkey
//	unknownPubKey := pubKeys[3]
//
//	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
//	s.setupValidatorAndDelegation(ctx, valPubKey, delPubKey1, valAddr, delAddr1, valTokens)
//
//	tcs := []struct {
//		name        string
//		settingMock func()
//		withdraw    *bindings.IPTokenStakingWithdraw
//		expectedErr string
//	}{
//		{
//			name: "pass: valid input",
//			settingMock: func() {
//				// bankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.BondedPoolName, stypes.NotBondedPoolName, gomock.Any())
//			},
//			withdraw: &bindings.IPTokenStakingWithdraw{
//				DelegatorUncmpPubkey: cmpToUncmp(delPubKey1.Bytes()),
//				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes()),
//				StakeAmount:          new(big.Int).SetUint64(1),
//				DelegationId:         big.NewInt(0),
//				OperatorAddress:      cmpToEVM(delPubKey1.Bytes()),
//			},
//		},
//		{
//			name: "fail: invalid delegator pubkey",
//			withdraw: &bindings.IPTokenStakingWithdraw{
//				DelegatorUncmpPubkey: cmpToUncmp(delPubKey1.Bytes())[:16],
//				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes()),
//				StakeAmount:          new(big.Int).SetUint64(1),
//				DelegationId:         big.NewInt(0),
//				OperatorAddress:      cmpToEVM(delPubKey1.Bytes()),
//			},
//			expectedErr: "invalid uncompressed public key length or format",
//		},
//		{
//			name: "fail: invalid validator pubkey",
//			withdraw: &bindings.IPTokenStakingWithdraw{
//				DelegatorUncmpPubkey: cmpToUncmp(delPubKey1.Bytes()),
//				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes())[:16],
//				StakeAmount:          new(big.Int).SetUint64(1),
//				DelegationId:         big.NewInt(0),
//				OperatorAddress:      cmpToEVM(delPubKey1.Bytes()),
//			},
//			expectedErr: "invalid uncompressed public key length or format",
//		},
//		{
//			name: "fail: corrupted delegator pubkey",
//			withdraw: &bindings.IPTokenStakingWithdraw{
//				DelegatorUncmpPubkey: createCorruptedPubKey(cmpToUncmp(delPubKey1.Bytes())),
//				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes()),
//				StakeAmount:          new(big.Int).SetUint64(1),
//				DelegationId:         big.NewInt(0),
//				OperatorAddress:      cmpToEVM(delPubKey1.Bytes()),
//			},
//			expectedErr: "invalid uncompressed public key length or format",
//		},
//		{
//			name: "fail: corrupted validator pubkey",
//			withdraw: &bindings.IPTokenStakingWithdraw{
//				DelegatorUncmpPubkey: cmpToUncmp(delPubKey1.Bytes()),
//				ValidatorUncmpPubkey: createCorruptedPubKey(cmpToUncmp(valPubKey.Bytes())),
//				StakeAmount:          new(big.Int).SetUint64(1),
//				DelegationId:         big.NewInt(0),
//				OperatorAddress:      cmpToEVM(delPubKey1.Bytes()),
//			},
//			expectedErr: "invalid uncompressed public key length or format",
//		},
//		{
//			name: "fail: unknown depositor",
//			settingMock: func() {
//				// accountKeeper.EXPECT().HasAccount(gomock.Any(), sdk.AccAddress(unknownPubKey.Address().Bytes())).Return(false).Times(1)
//			},
//			withdraw: &bindings.IPTokenStakingWithdraw{
//				DelegatorUncmpPubkey: cmpToUncmp(unknownPubKey.Bytes()),
//				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes()),
//				StakeAmount:          new(big.Int).SetUint64(1),
//				DelegationId:         big.NewInt(0),
//				OperatorAddress:      cmpToEVM(unknownPubKey.Bytes()),
//			},
//			expectedErr: "depositor account not found",
//		},
//		{
//			name: "fail: amount to withdraw is greater than the delegation amount",
//			settingMock: func() {
//				// accountKeeper.EXPECT().HasAccount(gomock.Any(), sdk.AccAddress(delPubKey1.Address().Bytes())).Return(true).Times(1)
//			},
//			withdraw: &bindings.IPTokenStakingWithdraw{
//				DelegatorUncmpPubkey: cmpToUncmp(delPubKey1.Bytes()),
//				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes()),
//				StakeAmount:          new(big.Int).SetUint64(math.MaxUint64),
//				DelegationId:         big.NewInt(0),
//				OperatorAddress:      cmpToEVM(delPubKey1.Bytes()),
//			},
//			expectedErr: "invalid shares amount",
//		},
//	}
//
//	for _, tc := range tcs {
//		s.Run(tc.name, func() {
//			if tc.settingMock != nil {
//				tc.settingMock()
//			}
//			cachedCtx, _ := ctx.CacheContext()
//			// check undelegation does not exist
//			_, err := s.StakingKeeper.GetUnbondingDelegation(cachedCtx, delAddr1, valAddr)
//			require.ErrorContains(err, "no unbonding delegation found")
//
//			err = keeper.ProcessWithdraw(cachedCtx, tc.withdraw)
//			if tc.expectedErr != "" {
//				require.ErrorContains(err, tc.expectedErr)
//			} else {
//				require.NoError(err)
//				// check undelegation exists
//				ubd, err := s.StakingKeeper.GetUnbondingDelegation(cachedCtx, delAddr1, valAddr)
//				require.NoError(err)
//				require.NotNil(ubd)
//			}
//		})
//	}
//}
//
//func (s *TestSuite) TestParseWithdraw() {
//	require := s.Require()
//	keeper := s.EVMStakingKeeper
//
//	testCases := []struct {
//		name      string
//		log       gethtypes.Log
//		expectErr bool
//	}{
//		{
//			name: "Unknown Topic",
//			log: gethtypes.Log{
//				Topics: []common.Hash{common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")},
//			},
//			expectErr: true,
//		},
//		{
//			name: "Valid Topic",
//			log: gethtypes.Log{
//				Topics: []common.Hash{types.WithdrawEvent.ID},
//			},
//			expectErr: false,
//		},
//	}
//
//	for _, tc := range testCases {
//		s.Run(tc.name, func() {
//			_, err := keeper.ParseWithdrawLog(tc.log)
//			if tc.expectErr {
//				require.Error(err, "should return error for %s", tc.name)
//			} else {
//				require.NoError(err, "should not return error for %s", tc.name)
//			}
//		})
//	}
//}
//
//// isEqualWithdrawals compares two slices of Withdrawal without considering order.
//func isEqualWithdrawals(t *testing.T, expected, actual []types.Withdrawal) {
//	t.Helper()
//	require.Len(t, actual, len(expected))
//	// compare it without considering order
//	for _, e := range expected {
//		found := false
//		for _, a := range actual {
//			if e.CreationHeight == a.CreationHeight &&
//				e.ExecutionAddress == a.ExecutionAddress &&
//				e.Amount == a.Amount {
//				found = true
//				break
//			}
//		}
//		if !found {
//			t.Errorf("expected %+v not found in %+v", e, actual)
//		}
//	}
//}
