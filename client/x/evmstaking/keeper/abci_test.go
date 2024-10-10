package keeper_test

import (
	"context"
	"testing"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"

	"go.uber.org/mock/gomock"
)

func (s *TestSuite) TestEndBlock() {
	require := s.Require()
	ctx, keeper, bankKeeper, stakingKeeper, distrKeeper := s.Ctx, s.EVMStakingKeeper, s.BankKeeper, s.StakingKeeper, s.DistrKeeper

	// create addresses
	pubKeys, accAddrs, valAddrs := createAddresses(3)
	delAddr := accAddrs[0]
	delPubKey := pubKeys[0]
	delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(delPubKey.Bytes())
	require.NoError(err)
	// setup two validators and delegations
	valPubKey1 := pubKeys[1]
	valAddr1 := valAddrs[1]
	valPubKey2 := pubKeys[2]
	valAddr2 := valAddrs[2]
	valCosmosPubKey1, err := k1util.PBPubKeyFromBytes(valPubKey1.Bytes())
	require.NoError(err)
	valCosmosPubKey2, err := k1util.PBPubKeyFromBytes(valPubKey2.Bytes())
	require.NoError(err)
	const valCnt = 2
	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
	s.setupValidatorAndDelegation(ctx, valPubKey1, delPubKey, valAddr1, delAddr, valTokens)
	s.setupValidatorAndDelegation(ctx, valPubKey2, delPubKey, valAddr2, delAddr, valTokens)

	// set evmstaking module's params
	minPartialAmt := int64(100)
	params, err := keeper.GetParams(ctx)
	require.NoError(err)
	params.MinPartialWithdrawalAmount = uint64(minPartialAmt)
	params.MaxSweepPerBlock = 100
	require.NoError(keeper.SetParams(ctx, params))

	// set staking module's params
	ubdTime := time.Duration(3600) * time.Second
	stakingParams, err := stakingKeeper.GetParams(ctx)
	require.NoError(err)
	stakingParams.UnbondingTime = ubdTime
	require.NoError(err)
	require.NoError(stakingKeeper.SetParams(ctx, stakingParams))

	mockEnqueueEligiblePartialWithdrawal := func(delAddr sdk.AccAddress, valAddr sdk.ValAddress, rewardAmt int64) {
		rewardsCoins := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, rewardAmt))
		distrKeeper.EXPECT().WithdrawDelegationRewards(gomock.Any(), delAddr, valAddr).Return(rewardsCoins, nil)
		bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), delAddr, types.ModuleName, rewardsCoins).Return(nil)
		bankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, rewardsCoins).Return(nil)
	}

	mockExpectPartialWithdrawals := func(c context.Context, valAddr sdk.ValAddress, exists bool) *types.Withdrawal {
		val, err := stakingKeeper.Validator(c, valAddr)
		require.NoError(err)
		dels, err := stakingKeeper.GetValidatorDelegations(c, valAddr)
		require.NoError(err)
		distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), valAddr).Return(dtypes.ValidatorAccumulatedCommission{}, nil)
		distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), val).Return(uint64(0), nil)
		rewardAmt := int64(0)
		if exists {
			// set reward amount to bigger than minPartialAmt, so that it can be partially withdrawn
			rewardAmt = minPartialAmt + 100
		}
		rewards := sdk.NewDecCoinsFromCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, rewardAmt))
		distrKeeper.EXPECT().CalculateDelegationRewards(gomock.Any(), val, dels[0], gomock.Any()).Return(rewards, nil)

		if exists {
			w := types.NewWithdrawal(0, delAddr.String(), valAddr.String(), delEvmAddr.String(), uint64(rewardAmt))
			mockEnqueueEligiblePartialWithdrawal(delAddr, valAddr, rewardAmt)

			return &w
		}

		return nil
	}

	postStateCheck := func(t *testing.T, c context.Context, expectedWithdrawals []types.Withdrawal) {
		t.Helper()
		withdrawals, err := keeper.GetWithdrawals(c, 100)
		require.NoError(err)

		isEqualWithdrawals(t, expectedWithdrawals, withdrawals)
	}

	tcs := []struct {
		name           string
		setup          func(c context.Context) ([]types.Withdrawal, []abcitypes.ValidatorUpdate)
		postStateCheck func(t *testing.T, c context.Context, expectedWithdrawals []types.Withdrawal)
		expectedError  string
	}{
		{
			name: "pass: no mature unbonded delegations & eligible partial withdrawals",
			setup: func(c context.Context) ([]types.Withdrawal, []abcitypes.ValidatorUpdate) {
				w1 := mockExpectPartialWithdrawals(c, valAddr1, true)
				w2 := mockExpectPartialWithdrawals(c, valAddr2, true)

				return []types.Withdrawal{*w1, *w2}, nil
			},
			postStateCheck: postStateCheck,
		},
		{
			name: "pass: no mature unbonded delegations & no eligible partial withdrawals",
			setup: func(c context.Context) ([]types.Withdrawal, []abcitypes.ValidatorUpdate) {
				mockExpectPartialWithdrawals(c, valAddr1, false)
				mockExpectPartialWithdrawals(c, valAddr2, false)

				return nil, nil
			},
			postStateCheck: postStateCheck,
		},
		{
			name: "pass: mature unbonded delegations & no eligible partial withdrawals",
			setup: func(c context.Context) ([]types.Withdrawal, []abcitypes.ValidatorUpdate) {
				sdkCtx := sdk.UnwrapSDKContext(c)
				s.setupMatureUnbondingDelegation(sdkCtx, delAddr, valAddr1, "15", ubdTime)
				s.setupMatureUnbondingDelegation(sdkCtx, delAddr, valAddr2, "15", ubdTime)

				mockExpectPartialWithdrawals(c, valAddr1, false)
				mockExpectPartialWithdrawals(c, valAddr2, false)

				return []types.Withdrawal{
						types.NewWithdrawal(0, delAddr.String(), valAddr1.String(), delEvmAddr.String(), 15),
						types.NewWithdrawal(0, delAddr.String(), valAddr2.String(), delEvmAddr.String(), 15),
					}, []abcitypes.ValidatorUpdate{
						{
							PubKey: valCosmosPubKey1,
							Power:  9,
						},
						{
							PubKey: valCosmosPubKey2,
							Power:  9,
						},
					}
			},
			postStateCheck: postStateCheck,
		},
		{
			name: "pass: mature unbonded delegations & eligible partial withdrawals",
			setup: func(c context.Context) ([]types.Withdrawal, []abcitypes.ValidatorUpdate) {
				sdkCtx := sdk.UnwrapSDKContext(c)
				s.setupMatureUnbondingDelegation(sdkCtx, delAddr, valAddr1, "10", ubdTime)
				s.setupMatureUnbondingDelegation(sdkCtx, delAddr, valAddr2, "10", ubdTime)

				w1 := mockExpectPartialWithdrawals(c, valAddr1, true)
				w2 := mockExpectPartialWithdrawals(c, valAddr2, true)

				return []types.Withdrawal{
						// withdrawals from unbonding delegations
						types.NewWithdrawal(0, delAddr.String(), valAddr1.String(), delEvmAddr.String(), 10),
						types.NewWithdrawal(0, delAddr.String(), valAddr2.String(), delEvmAddr.String(), 10),
						// partial withdrawals
						*w1,
						*w2,
					}, []abcitypes.ValidatorUpdate{
						{
							PubKey: valCosmosPubKey1,
							Power:  9,
						},
						{
							PubKey: valCosmosPubKey2,
							Power:  9,
						},
					}
			},
			postStateCheck: postStateCheck,
		},
		{
			name: "pass: mature + not matured unbonded delegations & no eligible partial withdrawals",
			setup: func(c context.Context) ([]types.Withdrawal, []abcitypes.ValidatorUpdate) {
				sdkCtx := sdk.UnwrapSDKContext(c)
				s.setupMatureUnbondingDelegation(sdkCtx, delAddr, valAddr1, "10", ubdTime)
				s.setupUnbonding(c, delAddr, valAddr2, "10")

				mockExpectPartialWithdrawals(c, valAddr1, false)
				mockExpectPartialWithdrawals(c, valAddr2, false)

				return []types.Withdrawal{
						types.NewWithdrawal(0, delAddr.String(), valAddr1.String(), delEvmAddr.String(), 10),
					}, []abcitypes.ValidatorUpdate{
						{
							PubKey: valCosmosPubKey1,
							Power:  9,
						},
						{
							PubKey: valCosmosPubKey2,
							Power:  9,
						},
					}
			},
			postStateCheck: postStateCheck,
		},
		{
			name: "pass: mature + not matured unbonded delegations & eligible partial withdrawals",
			setup: func(c context.Context) ([]types.Withdrawal, []abcitypes.ValidatorUpdate) {
				sdkCtx := sdk.UnwrapSDKContext(c)
				s.setupMatureUnbondingDelegation(sdkCtx, delAddr, valAddr1, "10", ubdTime)
				s.setupUnbonding(c, delAddr, valAddr2, "10")

				w1 := mockExpectPartialWithdrawals(c, valAddr1, true)
				w2 := mockExpectPartialWithdrawals(c, valAddr2, true)

				return []types.Withdrawal{
						// withdrawals from unbonding delegations
						types.NewWithdrawal(0, delAddr.String(), valAddr1.String(), delEvmAddr.String(), 10),
						// partial withdrawals
						*w1,
						*w2,
					}, []abcitypes.ValidatorUpdate{
						{
							PubKey: valCosmosPubKey1,
							Power:  9,
						},
						{
							PubKey: valCosmosPubKey2,
							Power:  9,
						},
					}
			},
			postStateCheck: postStateCheck,
		},
		{
			name: "fail: send coins from account to module during processing unbonded entry",
			setup: func(c context.Context) ([]types.Withdrawal, []abcitypes.ValidatorUpdate) {
				sdkCtx := sdk.UnwrapSDKContext(c)
				pastHeader := sdkCtx.BlockHeader()
				pastHeader.Time = pastHeader.Time.Add(-ubdTime).Add(-time.Minute)
				s.setupUnbonding(sdkCtx.WithBlockHeader(pastHeader), delAddr, valAddr1, "10")

				// Mock evmstaking.EndBlocker
				s.BankKeeper.EXPECT().UndelegateCoinsFromModuleToAccount(gomock.Any(), stypes.NotBondedPoolName, delAddr, gomock.Any()).Return(nil)
				s.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), delAddr, types.ModuleName, gomock.Any()).Return(errors.New("failed to send coins to module"))

				return nil, []abcitypes.ValidatorUpdate{
					{
						PubKey: valCosmosPubKey1,
						Power:  9,
					},
					{
						PubKey: valCosmosPubKey2,
						Power:  9,
					},
				}
			},
			expectedError: "failed to send coins to module",
		},
		{
			name: "fail: burn coins during processing unbonded entry",
			setup: func(c context.Context) ([]types.Withdrawal, []abcitypes.ValidatorUpdate) {
				sdkCtx := sdk.UnwrapSDKContext(c)
				pastHeader := sdkCtx.BlockHeader()
				pastHeader.Time = pastHeader.Time.Add(-ubdTime).Add(-time.Minute)
				s.setupUnbonding(sdkCtx.WithBlockHeader(pastHeader), delAddr, valAddr1, "10")

				// Mock evmstaking.EndBlocker
				s.BankKeeper.EXPECT().UndelegateCoinsFromModuleToAccount(gomock.Any(), stypes.NotBondedPoolName, delAddr, gomock.Any()).Return(nil)
				s.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), delAddr, types.ModuleName, gomock.Any()).Return(nil)
				s.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(errors.New("failed to burn coins"))

				return nil, []abcitypes.ValidatorUpdate{
					{
						PubKey: valCosmosPubKey1,
						Power:  9,
					},
					{
						PubKey: valCosmosPubKey2,
						Power:  9,
					},
				}
			},
			expectedError: "failed to burn coins",
		},
		{
			name: "fail: error while processing ExpectedPartialWithdrawals",
			setup: func(c context.Context) ([]types.Withdrawal, []abcitypes.ValidatorUpdate) {
				// Mock failed ExpectedPartialWithdrawals
				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, errors.New("failed to get validator accumulated commission"))

				return nil, nil
			},
			expectedError: "failed to get validator accumulated commission",
		},
		{
			name: "fail: error while processing EnqueueEligiblePartialWithdrawal",
			setup: func(c context.Context) ([]types.Withdrawal, []abcitypes.ValidatorUpdate) {
				rewards := sdk.NewDecCoinsFromCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, minPartialAmt+100))
				// Mock successful ExpectedPartialWithdrawals
				distrKeeper.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, nil).Times(valCnt)
				distrKeeper.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil).Times(valCnt)
				distrKeeper.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(rewards, nil).Times(valCnt)

				// Mock failed EnqueueEligiblePartialWithdrawal
				distrKeeper.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to withdraw delegation rewards"))

				return nil, nil
			},
			expectedError: "failed to withdraw delegation rewards",
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			var expectedWithdrawals []types.Withdrawal
			var expectedValUpdates []abcitypes.ValidatorUpdate
			cachedCtx, _ := ctx.CacheContext()
			if tc.setup != nil {
				expectedWithdrawals, expectedValUpdates = tc.setup(cachedCtx)
			}
			valUpdates, err := keeper.EndBlock(cachedCtx)
			if tc.expectedError != "" {
				require.ErrorContains(err, tc.expectedError)
			} else {
				require.NoError(err)
				if tc.postStateCheck != nil {
					tc.postStateCheck(s.T(), cachedCtx, expectedWithdrawals)
				}
				if expectedValUpdates != nil {
					compareValUpdates(s.T(), expectedValUpdates, valUpdates)
				}
			}
		})
	}
}

// compareValUpdates compares two slices of ValidatorUpdates, ignoring the order.
func compareValUpdates(t *testing.T, expected, actual abcitypes.ValidatorUpdates) {
	t.Helper()
	require.Equal(t, len(expected), len(actual), "the length of expected and actual slices should be equal")

	// Convert both slices to maps for unordered comparison
	expectedMap := make(map[string]abcitypes.ValidatorUpdate)
	actualMap := make(map[string]abcitypes.ValidatorUpdate)

	// Fill the maps using PubKey as the unique key
	for _, exp := range expected {
		expectedMap[exp.PubKey.String()] = exp
	}

	for _, act := range actual {
		actualMap[act.PubKey.String()] = act
	}

	// Compare the maps
	require.Equal(t, expectedMap, actualMap, "the content of expected and actual slices should match regardless of order")
}

// setupMaturedUnbonding creates matured unbondings for testing.
func (s *TestSuite) setupMatureUnbondingDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amt string, duration time.Duration) {
	pastHeader := ctx.BlockHeader()
	pastHeader.Time = pastHeader.Time.Add(-duration).Add(-time.Minute)
	pastCtx := ctx.WithBlockHeader(pastHeader)

	s.setupUnbonding(pastCtx, delAddr, valAddr, amt)

	// Mock staking.EndBlocker
	s.BankKeeper.EXPECT().UndelegateCoinsFromModuleToAccount(gomock.Any(), stypes.NotBondedPoolName, delAddr, gomock.Any()).Return(nil)
	// Mock evmstaking.EndBlocker
	s.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), delAddr, types.ModuleName, gomock.Any()).Return(nil)
	s.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
}
