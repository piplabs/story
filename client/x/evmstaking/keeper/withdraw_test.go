package keeper_test

import (
	"math"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"

	"go.uber.org/mock/gomock"
)

func (s *TestSuite) TestExpectedPartialWithdrawals() {
	require := s.Require()
	ctx, keeper, stakingKeeper, distrKeeper := s.Ctx, s.EVMStakingKeeper, s.StakingKeeper, s.DistrKeeper

	pubKeys, accAddrs, valAddrs := createAddresses(3)
	delAddr := accAddrs[0]
	delPubKey := pubKeys[0]
	evmDelAddr, err := k1util.CosmosPubkeyToEVMAddress(delPubKey.Bytes())
	require.NoError(err)
	valPubKey := pubKeys[1]
	valAddr := valAddrs[1]
	valPubKey2 := pubKeys[2]
	valAddr2 := valAddrs[2]

	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
	s.setupValidatorAndDelegation(ctx, valPubKey, delPubKey, valAddr, delAddr, valTokens)
	// set params as default
	params := types.DefaultParams()
	require.NoError(keeper.SetParams(ctx, params))
	delRewardsAmt := params.MinPartialWithdrawalAmount + 100
	delRewards := sdk.NewDecCoinsFromCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(delRewardsAmt)))

	// Test cases for ExpectedPartialWithdrawals
	tcs := []struct {
		name           string
		preRun         func(ctx sdk.Context)
		expectedResult []types.Withdrawal
		expectedError  string
	}{
		{
			name: "pass",
			preRun: func(_ sdk.Context) {
				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil)
				distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				distrKeeper.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delRewards, nil)
			},
			expectedResult: []types.Withdrawal{
				{
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: valAddr.String(),
					ExecutionAddress: evmDelAddr.String(),
					Amount:           delRewardsAmt,
				},
			},
		},
		{
			name: "pass: val sweep index is out of bounds, so it should be reset to 0 which is the index of the first validator",
			preRun: func(_ sdk.Context) {
				require.NoError(keeper.SetValidatorSweepIndex(ctx, uint64(100), uint64(0)))
				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil)
				distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				distrKeeper.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delRewards, nil)
			},
			expectedResult: []types.Withdrawal{
				{
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: valAddr.String(),
					ExecutionAddress: evmDelAddr.String(),
					Amount:           delRewardsAmt,
				},
			},
		},
		{
			name: "fail: increment validator period",
			preRun: func(_ sdk.Context) {
				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil)
				distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), errors.New("failed to increment validator period"))
			},
			expectedError: "failed to increment validator period",
		},
		{
			name: "fail: calculate delegation rewards",
			preRun: func(_ sdk.Context) {
				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil)
				distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				distrKeeper.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delRewards, errors.New("failed to calculate delegation rewards"))
			},
			expectedError: "failed to calculate delegation rewards",
		},
		{
			name: "pass: multiple validators",
			preRun: func(c sdk.Context) {
				s.setupValidatorAndDelegation(c, valPubKey2, delPubKey, valAddr2, delAddr, valTokens)
				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil).Times(2)
				distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil).Times(2)
				distrKeeper.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delRewards, nil).Times(2)
			},
			expectedResult: []types.Withdrawal{
				{
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: valAddr.String(),
					ExecutionAddress: evmDelAddr.String(),
					Amount:           delRewardsAmt,
				},
				{
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: valAddr2.String(),
					ExecutionAddress: evmDelAddr.String(),
					Amount:           delRewardsAmt,
				},
			},
		},
		{
			name: "pass: skip jailed validator",
			preRun: func(c sdk.Context) {
				s.setupValidatorAndDelegation(c, valPubKey2, delPubKey, valAddr2, delAddr, valTokens)
				val, err := stakingKeeper.GetValidator(c, valAddr2)
				require.NoError(err)
				val.Jailed = true
				require.NoError(stakingKeeper.SetValidator(c, val))

				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil)
				distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				distrKeeper.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delRewards, nil)
			},
			expectedResult: []types.Withdrawal{
				{
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: valAddr.String(),
					ExecutionAddress: evmDelAddr.String(),
					Amount:           delRewardsAmt,
				},
			},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cached, _ := ctx.CacheContext()
			if tc.preRun != nil {
				tc.preRun(cached)
			}
			result, err := keeper.ExpectedPartialWithdrawals(cached)
			if tc.expectedError != "" {
				require.ErrorContains(err, tc.expectedError)
			} else {
				require.NoError(err)
				isEqualWithdrawals(s.T(), tc.expectedResult, result)
			}
		})
	}
}

func (s *TestSuite) TestEnqueueEligiblePartialWithdrawal() {
	require := s.Require()
	ctx, keeper, bankKeeper, distrKeeper := s.Ctx, s.EVMStakingKeeper, s.BankKeeper, s.DistrKeeper

	pubKeys, accAddrs, valAddrs := createAddresses(2)
	// delegator
	delPubKey := pubKeys[0]
	delAddr := accAddrs[0]
	delValAddr := valAddrs[0] // delegator and validator are the same
	delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(delPubKey.Bytes())
	require.NoError(err)
	// validator
	valAddr := valAddrs[1]

	// Test cases for EnqueueEligiblePartialWithdrawal
	tcs := []struct {
		name          string
		settingMock   func(delRewards sdk.Coins)
		input         func() []types.Withdrawal
		expectedError string
	}{
		{
			name:          "fail: empty validator address",
			input:         func() []types.Withdrawal { return []types.Withdrawal{{ValidatorAddress: ""}} },
			expectedError: "validator address from bech32",
		},
		{
			name: "fail: invalid validator address",
			settingMock: func(_ sdk.Coins) {
				distrKeeper.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dtypes.ErrEmptyDelegationDistInfo)
			},
			input: func() []types.Withdrawal {
				return []types.Withdrawal{
					{ValidatorAddress: valAddr.String(), DelegatorAddress: delAddr.String(), Amount: 100},
				}
			},
			expectedError: dtypes.ErrEmptyDelegationDistInfo.Error(),
		},
		{
			name: "fail: validator and delegator are the same, but failed to withdraw commission",
			settingMock: func(delRewards sdk.Coins) {
				distrKeeper.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(delRewards, nil)
				distrKeeper.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).Return(sdk.NewCoins(), errors.New("failed to withdraw commission"))
			},
			input: func() []types.Withdrawal {
				return []types.Withdrawal{
					{
						CreationHeight:   0,
						DelegatorAddress: delAddr.String(),
						ValidatorAddress: delValAddr.String(),
						ExecutionAddress: delEvmAddr.String(),
						Amount:           100,
					},
				}
			},
			expectedError: "failed to withdraw commission",
		},
		{
			name: "pass: valid input",
			settingMock: func(delRewards sdk.Coins) {
				distrKeeper.EXPECT().WithdrawDelegationRewards(gomock.Any(), delAddr, valAddr).Return(delRewards, nil)
				bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), delAddr, types.ModuleName, delRewards).Return(nil)
				bankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, delRewards).Return(nil)
			},
			input: func() []types.Withdrawal {
				return []types.Withdrawal{
					{
						CreationHeight:   0,
						DelegatorAddress: delAddr.String(),
						ValidatorAddress: valAddr.String(),
						ExecutionAddress: delEvmAddr.String(),
						Amount:           100,
					},
				}
			},
		},
		{
			name: "pass: validator and delegator are the same",
			settingMock: func(delRewards sdk.Coins) {
				distrKeeper.EXPECT().WithdrawDelegationRewards(gomock.Any(), delAddr, delValAddr).Return(delRewards, nil)
				distrKeeper.EXPECT().WithdrawValidatorCommission(gomock.Any(), delValAddr).Return(sdk.NewCoins(), nil)
				bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), delAddr, types.ModuleName, delRewards).Return(nil)
				bankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, delRewards).Return(nil)
			},
			input: func() []types.Withdrawal {
				return []types.Withdrawal{
					{
						CreationHeight:   0,
						DelegatorAddress: delAddr.String(),
						ValidatorAddress: delValAddr.String(),
						ExecutionAddress: delEvmAddr.String(),
						Amount:           100,
					},
				}
			},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			input := tc.input()
			coinsExpectedToWithdraw := sdk.NewCoins(
				sdk.NewCoin(
					sdk.DefaultBondDenom,
					sdkmath.NewInt(int64(input[0].Amount)),
				),
			)
			if tc.settingMock != nil {
				tc.settingMock(coinsExpectedToWithdraw)
			}
			cachedCtx, _ := ctx.CacheContext()
			err := keeper.EnqueueEligiblePartialWithdrawal(cachedCtx, tc.input())
			if tc.expectedError != "" {
				require.ErrorContains(err, tc.expectedError)
			} else {
				require.NoError(err)
				withdrawals, err := keeper.GetAllWithdrawals(cachedCtx)
				require.NoError(err)
				isEqualWithdrawals(s.T(), tc.input(), withdrawals)
			}
		})
	}
}

func (s *TestSuite) TestProcessWithdraw() {
	require := s.Require()
	ctx, keeper, accountKeeper, bankKeeper, stakingKeeper := s.Ctx, s.EVMStakingKeeper, s.AccountKeeper, s.BankKeeper, s.StakingKeeper

	pubKeys, accAddrs, valAddrs := createAddresses(4)
	// delegator-1
	delPubKey1 := pubKeys[0]
	delAddr1 := accAddrs[0]
	// validator
	valPubKey := pubKeys[2]
	valAddr := valAddrs[2]
	// unknown pubkey
	unknownPubKey := pubKeys[3]

	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
	s.setupValidatorAndDelegation(ctx, valPubKey, delPubKey1, valAddr, delAddr1, valTokens)

	tcs := []struct {
		name        string
		settingMock func()
		withdraw    *bindings.IPTokenStakingWithdraw
		expectedErr string
	}{
		{
			name: "pass: valid input",
			settingMock: func() {
				accountKeeper.EXPECT().HasAccount(gomock.Any(), delAddr1).Return(true)
				bankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.BondedPoolName, stypes.NotBondedPoolName, gomock.Any())
			},
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorCmpPubkey: delPubKey1.Bytes(),
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Amount:             new(big.Int).SetUint64(1),
			},
		},
		{
			name: "fail: invalid delegator pubkey",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorCmpPubkey: delPubKey1.Bytes()[:16],
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Amount:             new(big.Int).SetUint64(1),
			},
			expectedErr: "invalid pubkey length",
		},
		{
			name: "fail: invalid validator pubkey",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorCmpPubkey: delPubKey1.Bytes(),
				ValidatorCmpPubkey: valPubKey.Bytes()[:16],
				Amount:             new(big.Int).SetUint64(1),
			},
			expectedErr: "invalid pubkey length",
		},
		{
			name: "fail: corrupted delegator pubkey",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorCmpPubkey: createCorruptedPubKey(delPubKey1.Bytes()),
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Amount:             new(big.Int).SetUint64(1),
			},
			expectedErr: "delegator pubkey to evm address",
		},
		{
			name: "fail: corrupted validator pubkey",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorCmpPubkey: delPubKey1.Bytes(),
				ValidatorCmpPubkey: createCorruptedPubKey(valPubKey.Bytes()),
				Amount:             new(big.Int).SetUint64(1),
			},
			expectedErr: "validator pubkey to evm address",
		},
		{
			name: "fail: unknown depositor",
			settingMock: func() {
				accountKeeper.EXPECT().HasAccount(gomock.Any(), sdk.AccAddress(unknownPubKey.Address().Bytes())).Return(false).Times(1)
			},
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorCmpPubkey: unknownPubKey.Bytes(),
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Amount:             new(big.Int).SetUint64(1),
			},
			expectedErr: "depositor account not found",
		},
		{
			name: "fail: amount to withdraw is greater than the delegation amount",
			settingMock: func() {
				accountKeeper.EXPECT().HasAccount(gomock.Any(), sdk.AccAddress(delPubKey1.Address().Bytes())).Return(true).Times(1)
			},
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorCmpPubkey: delPubKey1.Bytes(),
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Amount:             new(big.Int).SetUint64(math.MaxUint64),
			},
			expectedErr: "invalid shares amount",
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			if tc.settingMock != nil {
				tc.settingMock()
			}
			cachedCtx, _ := ctx.CacheContext()
			// check undelegation does not exist
			_, err := s.StakingKeeper.GetUnbondingDelegation(cachedCtx, delAddr1, valAddr)
			require.ErrorContains(err, "no unbonding delegation found")

			err = keeper.ProcessWithdraw(cachedCtx, tc.withdraw)
			if tc.expectedErr != "" {
				require.ErrorContains(err, tc.expectedErr)
			} else {
				require.NoError(err)
				// check undelegation exists
				ubd, err := s.StakingKeeper.GetUnbondingDelegation(cachedCtx, delAddr1, valAddr)
				require.NoError(err)
				require.NotNil(ubd)
			}
		})
	}
}

func (s *TestSuite) TestParseWithdraw() {
	require := s.Require()
	keeper := s.EVMStakingKeeper

	testCases := []struct {
		name      string
		log       gethtypes.Log
		expectErr bool
	}{
		{
			name: "Unknown Topic",
			log: gethtypes.Log{
				Topics: []common.Hash{common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")},
			},
			expectErr: true,
		},
		{
			name: "Valid Topic",
			log: gethtypes.Log{
				Topics: []common.Hash{types.WithdrawEvent.ID},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := keeper.ParseWithdrawLog(tc.log)
			if tc.expectErr {
				require.Error(err, "should return error for %s", tc.name)
			} else {
				require.NoError(err, "should not return error for %s", tc.name)
			}
		})
	}
}

// isEqualWithdrawals compares two slices of Withdrawal without considering order.
func isEqualWithdrawals(t *testing.T, expected, actual []types.Withdrawal) {
	t.Helper()
	require.Len(t, actual, len(expected))
	// compare it without considering order
	for _, e := range expected {
		found := false
		for _, a := range actual {
			if e.DelegatorAddress == a.DelegatorAddress &&
				e.ValidatorAddress == a.ValidatorAddress &&
				e.ExecutionAddress == a.ExecutionAddress &&
				e.Amount == a.Amount {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %+v not found in %+v", e, actual)
		}
	}
}
