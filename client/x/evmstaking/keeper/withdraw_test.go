package keeper_test

import (
	"math"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"

	"github.com/piplabs/story/lib/k1util"
)

// setupValidatorAndDelegation creates a validator and delegation for testing
func (s *TestSuite) setupValidatorAndDelegation(ctx sdk.Context, valPubKey, delPubKey crypto.PubKey, valAddr sdk.ValAddress, delAddr sdk.AccAddress) {
	require := s.Require()
	stakingKeeper := s.StakingKeeper
	keeper := s.EVMStakingKeeper

	// Convert public key to cosmos format
	valCosmosPubKey, err := k1util.PubKeyToCosmos(valPubKey)
	require.NoError(err)

	// Create and update validator
	val := testutil.NewValidator(s.T(), valAddr, valCosmosPubKey)
	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
	validator, _ := val.AddTokensFromDel(valTokens)
	s.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.NotBondedPoolName, stypes.BondedPoolName, gomock.Any())
	_ = skeeper.TestingUpdateValidator(stakingKeeper, ctx, validator, true)

	// Create and set delegation
	delAmt := stakingKeeper.TokensFromConsensusPower(ctx, 100).ToLegacyDec()
	delegation := stypes.NewDelegation(delAddr.String(), valAddr.String(), delAmt)
	require.NoError(stakingKeeper.SetDelegation(ctx, delegation))

	// Map delegator to EVM address
	delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(delPubKey.Bytes())
	require.NoError(err)
	require.NoError(keeper.DelegatorMap.Set(ctx, delAddr.String(), delEvmAddr.String()))

	// Ensure delegation is set correctly
	delegation, err = stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	require.NoError(err)
	require.Equal(delAmt, delegation.GetShares())
}

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

	s.setupValidatorAndDelegation(ctx, valPubKey, delPubKey, valAddr, delAddr)
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
			name: "pass: next val sweep index is out of bounds, so it should be reset to 0 which is the index of the first validator",
			preRun: func(_ sdk.Context) {
				require.NoError(keeper.SetNextValidatorSweepIndex(ctx, sdk.IntProto{Int: sdkmath.NewInt(100)}))
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
				s.setupValidatorAndDelegation(c, valPubKey2, delPubKey, valAddr2, delAddr)
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
				s.setupValidatorAndDelegation(c, valPubKey2, delPubKey, valAddr2, delAddr)
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
		settingMock   func()
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
			settingMock: func() {
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
			name: "pass: valid input",
			settingMock: func() {
				distrKeeper.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
			},
			input: func() []types.Withdrawal {
				return []types.Withdrawal{
					{
						CreationHeight:   1,
						DelegatorAddress: delAddr.String(),
						ValidatorAddress: valAddr.String(),
						ExecutionAddress: delEvmAddr.String(),
						Amount:           100,
					},
				}
			},
			expectedError: "",
		},
		{
			name: "pass: validator and delegator are the same",
			settingMock: func() {
				distrKeeper.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				distrKeeper.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).Return(sdk.NewCoins(), nil)
				bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
			},
			input: func() []types.Withdrawal {
				return []types.Withdrawal{
					{
						CreationHeight:   1,
						DelegatorAddress: delAddr.String(),
						ValidatorAddress: delValAddr.String(),
						ExecutionAddress: delEvmAddr.String(),
						Amount:           100,
					},
				}
			},
			expectedError: "",
		},
		{
			name: "fail: validator and delegator are the same, but failed to withdraw commission",
			settingMock: func() {
				distrKeeper.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				distrKeeper.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).Return(sdk.NewCoins(), errors.New("failed to withdraw commission"))
			},
			input: func() []types.Withdrawal {
				return []types.Withdrawal{
					{
						CreationHeight:   1,
						DelegatorAddress: delAddr.String(),
						ValidatorAddress: delValAddr.String(),
						ExecutionAddress: delEvmAddr.String(),
						Amount:           100,
					},
				}
			},
			expectedError: "failed to withdraw commission",
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			if tc.settingMock != nil {
				tc.settingMock()
			}
			err := keeper.EnqueueEligiblePartialWithdrawal(ctx, tc.input())
			if tc.expectedError != "" {
				require.ErrorContains(err, tc.expectedError)
			} else {
				require.NoError(err)
			}
		})
	}
}

func (s *TestSuite) TestProcessWithdraw() {
	require := s.Require()
	ctx, keeper, accountKeeper, bankKeeper := s.Ctx, s.EVMStakingKeeper, s.AccountKeeper, s.BankKeeper

	pubKeys, accAddrs, valAddrs := createAddresses(4)
	// delegator-1
	delPubKey1 := pubKeys[0]
	delAddr1 := accAddrs[0]
	// delegator-2
	delPubKey2 := pubKeys[1]
	delAddr2 := accAddrs[1]
	// validator
	valPubKey := pubKeys[2]
	valAddr := valAddrs[2]
	// unknown pubkey
	unknownPubKey := pubKeys[3]

	s.setupValidatorAndDelegation(ctx, valPubKey, delPubKey1, valAddr, delAddr1)
	s.setupValidatorAndDelegation(ctx, valPubKey, delPubKey2, valAddr, delAddr2)

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
				Raw:                gethtypes.Log{},
			},
		},
		{
			name: "fail: invalid delegator pubkey",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorCmpPubkey: delPubKey2.Bytes()[:16],
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Amount:             new(big.Int).SetUint64(1),
				Raw:                gethtypes.Log{},
			},
			expectedErr: "invalid pubkey length",
		},
		{
			name: "fail: invalid validator pubkey",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorCmpPubkey: delPubKey2.Bytes(),
				ValidatorCmpPubkey: valPubKey.Bytes()[:16],
				Amount:             new(big.Int).SetUint64(1),
				Raw:                gethtypes.Log{},
			},
			expectedErr: "invalid pubkey length",
		},
		{
			name: "fail: corrupted delegator pubkey",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorCmpPubkey: createCorruptedPubKey(delPubKey1.Bytes()),
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Amount:             new(big.Int).SetUint64(1),
				Raw:                gethtypes.Log{},
			},
			expectedErr: "delegator pubkey to evm address",
		},
		{
			name: "fail: corrupted validator pubkey",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorCmpPubkey: delPubKey1.Bytes(),
				ValidatorCmpPubkey: createCorruptedPubKey(valPubKey.Bytes()),
				Amount:             new(big.Int).SetUint64(1),
				Raw:                gethtypes.Log{},
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
				Raw:                gethtypes.Log{},
			},
			expectedErr: "depositor account not found",
		},
		{
			name: "fail: amount to withdraw is greater than the delegation amount",
			settingMock: func() {
				accountKeeper.EXPECT().HasAccount(gomock.Any(), sdk.AccAddress(delPubKey2.Address().Bytes())).Return(true).Times(1)
			},
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorCmpPubkey: delPubKey2.Bytes(),
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Amount:             new(big.Int).SetUint64(math.MaxUint64),
				Raw:                gethtypes.Log{},
			},
			expectedErr: "invalid shares amount",
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			if tc.settingMock != nil {
				tc.settingMock()
			}
			err := keeper.ProcessWithdraw(ctx, tc.withdraw)
			if tc.expectedErr != "" {
				require.ErrorContains(err, tc.expectedErr)
			} else {
				require.NoError(err)
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

func createCorruptedPubKey(pubKey []byte) []byte {
	corruptedPubKey := append([]byte(nil), pubKey...)
	corruptedPubKey[0] = 0x04
	corruptedPubKey[1] = 0xFF
	return corruptedPubKey
}

// isEqualWithdrawals compares two slices of Withdrawal without considering order
func isEqualWithdrawals(t *testing.T, expected, actual []types.Withdrawal) {
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
