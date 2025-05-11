package keeper_test

import (
	"math/big"
	"strings"
	"testing"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/server/utils"
	"github.com/piplabs/story/client/x/evmstaking/keeper"
	moduletestutil "github.com/piplabs/story/client/x/evmstaking/testutil"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"

	"go.uber.org/mock/gomock"
)

func TestProcessUnstakeWithdrawals(t *testing.T) {
	pubKeys, accAddrs, valAddrs := createAddresses(2)

	// delegator
	delPubKey := pubKeys[0]
	delAccAddr := accAddrs[0]
	delEVMAddr, err := keeper.CmpPubKeyToEVMAddress(delPubKey.Bytes())
	require.NoError(t, err)

	// validators
	valPubKey := pubKeys[1]
	valValAddr := valAddrs[1]
	valEVMAddr, err := keeper.CmpPubKeyToEVMAddress(valPubKey.Bytes())
	require.NoError(t, err)

	getDelegatorAddr := func(isValidator bool) string {
		if isValidator {
			return sdk.AccAddress(valValAddr).String()
		} else {
			return delAccAddr.String()
		}
	}

	type expectedResult struct {
		withdrawals       []types.Withdrawal
		rewardWithdrawals []types.Withdrawal
	}

	tcs := []struct {
		name           string
		setupMocks     func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper)
		setup          func(ctx sdk.Context, esk *keeper.Keeper)
		unbondEntries  []stypes.UnbondedEntry
		expectedErr    string
		expectedResult func(ctx sdk.Context) expectedResult
	}{
		{
			name: "fail: invalid delegator address",
			unbondEntries: []stypes.UnbondedEntry{
				{
					DelegatorAddress: strings.Replace(delAccAddr.String(), "story", "cosmos", 1),
					ValidatorAddress: valValAddr.String(),
					Amount:           math.NewInt(10),
				},
			},
			expectedErr: "delegator address from bech32",
		},
		{
			name: "fail: invalid validator address",
			unbondEntries: []stypes.UnbondedEntry{
				{
					DelegatorAddress: delAccAddr.String(),
					ValidatorAddress: strings.Replace(valValAddr.String(), "story", "cosmos", 1),
					Amount:           math.NewInt(10),
				},
			},
			expectedErr: "validator address from bech32",
		},
		{
			name: "fail: get delegation - unknown error",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(stypes.Delegation{}, errors.New("failed to get delegation"))
			},
			unbondEntries: []stypes.UnbondedEntry{
				{
					DelegatorAddress: delAccAddr.String(),
					ValidatorAddress: valValAddr.String(),
					Amount:           math.NewInt(10),
				},
			},
			expectedErr: "get delegation: failed to get delegation",
		},
		{
			name: "fail: unstake from validator and totally unstaked, but failed to withdraw validator commission",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(stypes.Delegation{}, stypes.ErrNoDelegation)
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).Return(sdk.Coins{}, errors.New("failed to withdraw validator commission"))
			},
			unbondEntries: []stypes.UnbondedEntry{
				{
					DelegatorAddress: delAccAddr.String(),
					ValidatorAddress: sdk.ValAddress(delAccAddr.Bytes()).String(),
					Amount:           math.NewInt(10),
				},
			},
			expectedErr: "withdraw validator commission",
		},
		{
			name: "fail: send coins from account to module",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(stypes.Delegation{}, nil)
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to send coins from account to module"))
			},
			unbondEntries: []stypes.UnbondedEntry{
				{
					DelegatorAddress: delAccAddr.String(),
					ValidatorAddress: valValAddr.String(),
					Amount:           math.NewInt(10),
				},
			},
			expectedErr: "send coins from account to module",
		},
		{
			name: "fail: burn coins",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(stypes.Delegation{}, nil)
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to burn coins"))
			},
			unbondEntries: []stypes.UnbondedEntry{
				{
					DelegatorAddress: delAccAddr.String(),
					ValidatorAddress: valValAddr.String(),
					Amount:           math.NewInt(10),
				},
			},
			expectedErr: "burn coins",
		},
		{
			name: "fail: no withdraw address found",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(stypes.Delegation{}, nil)
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			unbondEntries: []stypes.UnbondedEntry{
				{
					DelegatorAddress: delAccAddr.String(),
					ValidatorAddress: valValAddr.String(),
					Amount:           math.NewInt(10),
				},
			},
			expectedErr: "map delegator pubkey to evm address",
		},
		{
			name: "pass: process withdrawal from delegator without reward",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(stypes.Delegation{}, nil)
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorWithdrawAddress.Set(ctx, getDelegatorAddr(false), delEVMAddr.String()))
			},
			unbondEntries: []stypes.UnbondedEntry{
				{
					DelegatorAddress: getDelegatorAddr(false),
					ValidatorAddress: valValAddr.String(),
					Amount:           math.NewInt(10),
				},
			},
			expectedResult: func(ctx sdk.Context) expectedResult {
				return expectedResult{
					withdrawals: []types.Withdrawal{
						{
							CreationHeight:   uint64(ctx.BlockHeight()),
							ExecutionAddress: delEVMAddr.String(),
							Amount:           10,
							WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE,
							ValidatorAddress: strings.ToLower(valEVMAddr.String()),
						},
					},
				}
			},
		},
		{
			name: "pass: process withdrawal from totally unstaked validator, but no commission",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(stypes.Delegation{}, stypes.ErrNoDelegation)
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				bk.EXPECT().SpendableCoin(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{Amount: math.NewInt(0)})
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorWithdrawAddress.Set(ctx, getDelegatorAddr(true), valEVMAddr.String()))
			},
			unbondEntries: []stypes.UnbondedEntry{
				{
					DelegatorAddress: getDelegatorAddr(true),
					ValidatorAddress: valValAddr.String(),
					Amount:           math.NewInt(10),
				},
			},
			expectedResult: func(ctx sdk.Context) expectedResult {
				return expectedResult{
					withdrawals: []types.Withdrawal{
						{
							CreationHeight:   uint64(ctx.BlockHeight()),
							ExecutionAddress: valEVMAddr.String(),
							Amount:           10,
							WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE,
							ValidatorAddress: strings.ToLower(valEVMAddr.String()),
						},
					},
				}
			},
		},
		{
			name: "process withdrawal from totally unstaked delegator with residue reward",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(stypes.Delegation{}, stypes.ErrNoDelegation)
				bk.EXPECT().SpendableCoin(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{Amount: math.NewInt(20)})
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorWithdrawAddress.Set(ctx, getDelegatorAddr(false), delEVMAddr.String()))
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(false), delEVMAddr.String()))
			},
			unbondEntries: []stypes.UnbondedEntry{
				{
					DelegatorAddress: getDelegatorAddr(false),
					ValidatorAddress: valValAddr.String(),
					Amount:           math.NewInt(10),
				},
			},
			expectedResult: func(ctx sdk.Context) expectedResult {
				return expectedResult{
					withdrawals: []types.Withdrawal{
						{
							CreationHeight:   uint64(ctx.BlockHeight()),
							ExecutionAddress: delEVMAddr.String(),
							Amount:           10,
							WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE,
							ValidatorAddress: strings.ToLower(valEVMAddr.String()),
						},
					},
					rewardWithdrawals: []types.Withdrawal{
						{
							CreationHeight:   uint64(ctx.BlockHeight()),
							ExecutionAddress: delEVMAddr.String(),
							Amount:           10,
							WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
							ValidatorAddress: strings.ToLower(valEVMAddr.String()),
						},
					},
				}
			},
		},
		{
			name: "process multiple withdrawals",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, bk, dk, sk, _, esk := createKeeperWithMockStaking(t)

			cachedCtx, _ := ctx.CacheContext()

			if tc.setupMocks != nil {
				tc.setupMocks(bk, dk, sk)
			}

			if tc.setup != nil {
				tc.setup(cachedCtx, esk)
			}

			// initialize reward withdrawal queue
			require.NoError(t, esk.RewardWithdrawalQueue.Initialize(cachedCtx))

			err := esk.ProcessUnstakeWithdrawals(cachedCtx, tc.unbondEntries)
			//nolint:nestif // nested check
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)

				if tc.expectedResult != nil {
					expected := tc.expectedResult(cachedCtx)

					ws, err := esk.GetAllWithdrawals(cachedCtx)
					require.NoError(t, err)

					rws, err := esk.GetAllRewardWithdrawals(cachedCtx)
					require.NoError(t, err)

					if len(expected.withdrawals) > 0 {
						require.Equal(t, len(expected.withdrawals), len(ws))
						for i, w := range ws {
							require.Equal(t, expected.withdrawals[i], w)
						}
					}

					if len(expected.rewardWithdrawals) > 0 {
						require.Equal(t, len(expected.rewardWithdrawals), len(rws))
						for i, rw := range rws {
							require.Equal(t, expected.rewardWithdrawals[i], rw)
						}
					}
				}
			}
		})
	}
}

func TestProcessRewardWithdrawals(t *testing.T) {
	pubKeys, accAddrs, valAddrs := createAddresses(4)

	// delegator
	delPubKey := pubKeys[0]
	delAccAddr := accAddrs[0]
	delEVMAddr, err := keeper.CmpPubKeyToEVMAddress(delPubKey.Bytes())
	require.NoError(t, err)

	// validators
	val1PubKey := pubKeys[1]
	val1ValAddr := valAddrs[1]
	val1CosmosPubKey, err := k1util.PubKeyToCosmos(val1PubKey)
	require.NoError(t, err)
	val1EVMAddr, err := keeper.CmpPubKeyToEVMAddress(val1PubKey.Bytes())
	require.NoError(t, err)

	val2PubKey := pubKeys[2]
	val2ValAddr := valAddrs[2]
	val2CosmosPubKey, err := k1util.PubKeyToCosmos(val2PubKey)
	require.NoError(t, err)
	val2EVMAddr, err := keeper.CmpPubKeyToEVMAddress(val2PubKey.Bytes())
	require.NoError(t, err)

	// additional delegator
	delAccAddr2 := accAddrs[3]

	delegationRewardDecCoin := sdk.NewDecCoins(sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, math.LegacyNewDec(10)))
	delegationRewardCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(10))}

	getDelegatorAddr := func(isValidator bool) string {
		if isValidator {
			return sdk.AccAddress(val1ValAddr).String()
		} else {
			return delAccAddr.String()
		}
	}

	tcs := []struct {
		name           string
		setupMocks     func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper)
		setup          func(ctx sdk.Context, esk *keeper.Keeper)
		expectedErr    string
		expectedResult func(ctx sdk.Context) []types.Withdrawal
	}{
		{
			name: "fail: get all validators",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{}, errors.New("failed to get all validators"))
			},
			expectedErr: "get all validators",
		},
		{
			name: "fail: invalid operator address of validator",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				valSet := []stypes.Validator{{
					Jailed:          false,
					OperatorAddress: strings.Replace(val1ValAddr.String(), "story", "cosmos", 1),
				}}
				sk.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec("storyvaloper"))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return(valSet, nil)
			},
			expectedErr: "convert validator address from string to bytes",
		},
		{
			name: "fail: get validator delegations",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				valSet := []stypes.Validator{{
					Jailed:          false,
					OperatorAddress: val1ValAddr.String(),
				}}
				sk.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec("storyvaloper"))
				sk.EXPECT().GetAllValidators(gomock.Any()).Return(valSet, nil)
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{}, errors.New("failed to get delegations"))
			},
			expectedErr: "get validator delegations",
		},
		{
			name: "fail: process eligible reward withdrawal due to the failure of increment validator period",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				valSet := []stypes.Validator{{
					Jailed:          false,
					OperatorAddress: val1ValAddr.String(),
				}}
				delegations := []stypes.Delegation{{
					DelegatorAddress: getDelegatorAddr(false),
					ValidatorAddress: val1ValAddr.String(),
				}}
				sk.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec("storyvaloper")).Times(2)
				sk.EXPECT().GetAllValidators(gomock.Any()).Return(valSet, nil)
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return(delegations, nil)
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), errors.New("failed to increment validator period"))
			},
			expectedErr: "process eligible reward withdrawal",
		},
		{
			name: "pass: next validator index overflow",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				// set validator set
				validator1, err := stypes.NewValidator(val1ValAddr.String(), val1CosmosPubKey, stypes.Description{Moniker: "validator1"}, 0)
				require.NoError(t, err)
				validator2, err := stypes.NewValidator(val2ValAddr.String(), val2CosmosPubKey, stypes.Description{Moniker: "validator2"}, 0)
				require.NoError(t, err)
				valSet := []stypes.Validator{validator1, validator2}

				// set delegations
				delegation1 := stypes.NewDelegation(delAccAddr.String(), val1ValAddr.String(), math.LegacyNewDec(10), math.LegacyNewDec(10).Quo(math.LegacyNewDec(2)))
				delegation2 := stypes.NewDelegation(delAccAddr.String(), val2ValAddr.String(), math.LegacyNewDec(10), math.LegacyNewDec(10).Quo(math.LegacyNewDec(2)))

				sk.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec("storyvaloper")).AnyTimes()
				sk.EXPECT().GetAllValidators(gomock.Any()).Return(valSet, nil)
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), val1ValAddr).Return([]stypes.Delegation{delegation1}, nil)
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), val2ValAddr).Return([]stypes.Delegation{delegation2}, nil)
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil).AnyTimes()
				dk.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardDecCoin, nil).AnyTimes()
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardCoins, nil).AnyTimes()
				bk.EXPECT().SpendableCoin(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: math.NewInt(0)}).AnyTimes()
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(false), delEVMAddr.String()))

				// set params
				params := types.Params{
					MaxWithdrawalPerBlock:      2,
					MaxSweepPerBlock:           4,
					MinPartialWithdrawalAmount: 2,
				}
				require.NoError(t, esk.SetParams(ctx, params))

				// set sweep index
				newSweepIndex := types.NewValidatorSweepIndex(3, 0)
				require.NoError(t, esk.SetValidatorSweepIndex(ctx, newSweepIndex))
			},
			expectedResult: func(ctx sdk.Context) []types.Withdrawal {
				return []types.Withdrawal{
					{
						CreationHeight:   uint64(ctx.BlockHeight()),
						ExecutionAddress: delEVMAddr.String(),
						Amount:           10,
						WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
						ValidatorAddress: strings.ToLower(val1EVMAddr.String()),
					},
					{
						CreationHeight:   uint64(ctx.BlockHeight()),
						ExecutionAddress: delEVMAddr.String(),
						Amount:           10,
						WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
						ValidatorAddress: strings.ToLower(val2EVMAddr.String()),
					},
				}
			},
		},
		{
			name: "pass: next validator delegation index overflow",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				// set validator set
				validator1, err := stypes.NewValidator(val1ValAddr.String(), val1CosmosPubKey, stypes.Description{Moniker: "validator1"}, 0)
				require.NoError(t, err)
				validator2, err := stypes.NewValidator(val2ValAddr.String(), val2CosmosPubKey, stypes.Description{Moniker: "validator2"}, 0)
				require.NoError(t, err)
				valSet := []stypes.Validator{validator1, validator2}

				// set delegations
				delegation2 := stypes.NewDelegation(delAccAddr.String(), val2ValAddr.String(), math.LegacyNewDec(10), math.LegacyNewDec(10).Quo(math.LegacyNewDec(2)))

				sk.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec("storyvaloper")).AnyTimes()
				sk.EXPECT().GetAllValidators(gomock.Any()).Return(valSet, nil)
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), val1ValAddr).Return([]stypes.Delegation{}, nil)
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), val2ValAddr).Return([]stypes.Delegation{delegation2}, nil)
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil).AnyTimes()
				dk.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardDecCoin, nil).AnyTimes()
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardCoins, nil).AnyTimes()
				bk.EXPECT().SpendableCoin(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: math.NewInt(0)}).AnyTimes()
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(false), delEVMAddr.String()))

				// set params
				params := types.Params{
					MaxWithdrawalPerBlock:      2,
					MaxSweepPerBlock:           4,
					MinPartialWithdrawalAmount: 2,
				}
				require.NoError(t, esk.SetParams(ctx, params))

				// set sweep index
				newSweepIndex := types.NewValidatorSweepIndex(0, 2)
				require.NoError(t, esk.SetValidatorSweepIndex(ctx, newSweepIndex))
			},
			expectedResult: func(ctx sdk.Context) []types.Withdrawal {
				return []types.Withdrawal{
					{
						CreationHeight:   uint64(ctx.BlockHeight()),
						ExecutionAddress: delEVMAddr.String(),
						Amount:           10,
						WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
						ValidatorAddress: strings.ToLower(val2EVMAddr.String()),
					},
				}
			},
		},
		{
			name: "pass: skip jailed validator and overflowed delegations",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				// set validator set
				validator1, err := stypes.NewValidator(val1ValAddr.String(), val1CosmosPubKey, stypes.Description{Moniker: "validator1"}, 0)
				require.NoError(t, err)
				validator1.Jailed = true
				validator2, err := stypes.NewValidator(val2ValAddr.String(), val2CosmosPubKey, stypes.Description{Moniker: "validator2"}, 0)
				require.NoError(t, err)
				valSet := []stypes.Validator{validator1, validator2}

				// set delegations
				delegation1 := stypes.NewDelegation(delAccAddr2.String(), val2ValAddr.String(), math.LegacyNewDec(10), math.LegacyNewDec(10).Quo(math.LegacyNewDec(2)))
				delegation2 := stypes.NewDelegation(delAccAddr.String(), val2ValAddr.String(), math.LegacyNewDec(10), math.LegacyNewDec(10).Quo(math.LegacyNewDec(2)))

				sk.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec("storyvaloper")).AnyTimes()
				sk.EXPECT().GetAllValidators(gomock.Any()).Return(valSet, nil)
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), val2ValAddr).Return([]stypes.Delegation{delegation1, delegation2}, nil)
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil).AnyTimes()
				dk.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardDecCoin, nil).AnyTimes()
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardCoins, nil).AnyTimes()
				bk.EXPECT().SpendableCoin(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: math.NewInt(0)}).AnyTimes()
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(false), delEVMAddr.String()))
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, delAccAddr2.String(), delEVMAddr.String()))

				// set params
				params := types.Params{
					MaxWithdrawalPerBlock:      2,
					MaxSweepPerBlock:           1,
					MinPartialWithdrawalAmount: 2,
				}
				require.NoError(t, esk.SetParams(ctx, params))

				// set sweep index
				newSweepIndex := types.NewValidatorSweepIndex(0, 0)
				require.NoError(t, esk.SetValidatorSweepIndex(ctx, newSweepIndex))
			},
			expectedResult: func(ctx sdk.Context) []types.Withdrawal {
				return []types.Withdrawal{
					{
						CreationHeight:   uint64(ctx.BlockHeight()),
						ExecutionAddress: delEVMAddr.String(),
						Amount:           10,
						WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
						ValidatorAddress: strings.ToLower(val2EVMAddr.String()),
					},
				}
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, bk, dk, sk, _, esk := createKeeperWithMockStaking(t)

			cachedCtx, _ := ctx.CacheContext()

			if tc.setupMocks != nil {
				tc.setupMocks(bk, dk, sk)
			}

			if tc.setup != nil {
				tc.setup(cachedCtx, esk)
			}

			// initialize reward withdrawal queue
			require.NoError(t, esk.RewardWithdrawalQueue.Initialize(cachedCtx))

			err := esk.ProcessRewardWithdrawals(cachedCtx)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)

				if len(tc.expectedResult(cachedCtx)) != 0 {
					expected := tc.expectedResult(cachedCtx)
					rws, err := esk.GetAllRewardWithdrawals(cachedCtx)
					require.NoError(t, err)
					require.Equal(t, len(expected), len(rws))
					for i, rw := range rws {
						require.Equal(t, expected[i], rw)
					}
				}
			}
		})
	}
}

func TestProcessEligibleRewardWithdrawal(t *testing.T) {
	pubKeys, accAddrs, valAddrs := createAddresses(2)

	// delegator
	delPubKey := pubKeys[0]
	delAccAddr := accAddrs[0]
	delEVMAddr, err := keeper.CmpPubKeyToEVMAddress(delPubKey.Bytes())
	require.NoError(t, err)

	// validator
	valValAddr := valAddrs[1]
	valEVMAddr, err := utils.Bech32ValidatorAddressToEvmAddress(valValAddr.String())
	require.NoError(t, err)

	delegationRewardDecCoin := sdk.NewDecCoins(sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, math.LegacyNewDec(10)))
	delegationRewardCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(10))}

	getDelegatorAddr := func(isValidator bool) string {
		if isValidator {
			return sdk.AccAddress(valValAddr).String()
		} else {
			return delAccAddr.String()
		}
	}

	tcs := []struct {
		name                      string
		setupMocks                func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper)
		setup                     func(ctx sdk.Context, esk *keeper.Keeper)
		delegation                stypes.Delegation
		validator                 stypes.Validator
		minRewardWithdrawalAmount uint64
		expectedErr               string
		expectedResult            func(ctx sdk.Context) *types.Withdrawal
	}{
		{
			name: "fail: invalid validator address",
			validator: stypes.Validator{
				OperatorAddress: strings.Replace(valValAddr.String(), "story", "cosmos", 1),
			},
			expectedErr: "validator address from bech32",
		},
		{
			name: "fail: increase validator period from distribution keeper",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), errors.New("failed to increase validator period"))
			},
			validator: stypes.Validator{
				OperatorAddress: valValAddr.String(),
			},
			expectedErr: "failed to increase validator period",
		},
		{
			name: "fail: calculate delegation rewards",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				dk.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.DecCoins{}, errors.New("failed to calculate delegation rewards"))
			},
			validator: stypes.Validator{
				OperatorAddress: valValAddr.String(),
			},
			expectedErr: "failed to calculate delegation rewards",
		},
		{
			name: "fail: get validator accumulated commission when delegator is validator",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				dk.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.DecCoins{}, nil)
				dk.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{}, errors.New("failed to get validator accumulated commission"))
			},
			delegation: stypes.Delegation{
				DelegatorAddress: getDelegatorAddr(true),
			},
			validator: stypes.Validator{
				OperatorAddress: valValAddr.String(),
			},
			expectedErr: "failed to get validator accumulated commission",
		},
		{
			name: "fail: invalid delegator address",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				dk.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.DecCoins{}, nil)
			},
			delegation: stypes.Delegation{
				DelegatorAddress: strings.Replace(delAccAddr.String(), "story", "cosmos", 1),
			},
			validator: stypes.Validator{
				OperatorAddress: valValAddr.String(),
			},
			expectedErr: "convert acc address from bech32 address",
		},
		{
			name: "fail: enqueue reward withdrawal due to the failure of delegation rewards withdrawal from distribution keeper",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				dk.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardDecCoin, nil)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coins{}, errors.New("failed to withdraw delegation rewards"))
				bk.EXPECT().SpendableCoin(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: math.NewInt(0)})
			},
			delegation: stypes.Delegation{
				DelegatorAddress: getDelegatorAddr(false),
			},
			validator: stypes.Validator{
				OperatorAddress: valValAddr.String(),
			},
			minRewardWithdrawalAmount: 10,
			expectedErr:               "enqueue reward withdrawal",
		},
		{
			name: "pass(skip): total reward is less than minimum reward withdrawal amount",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				dk.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardDecCoin, nil)
				bk.EXPECT().SpendableCoin(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: math.NewInt(0)})
			},
			delegation: stypes.Delegation{
				DelegatorAddress: getDelegatorAddr(false),
			},
			validator: stypes.Validator{
				OperatorAddress: valValAddr.String(),
			},
			minRewardWithdrawalAmount: 11,
		},
		{
			name: "pass: process eligible reward withdrawal - with zero claimed reward",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				dk.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardDecCoin, nil)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardCoins, nil)
				bk.EXPECT().SpendableCoin(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: math.NewInt(0)})
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(false), delEVMAddr.String()))
			},
			delegation: stypes.Delegation{
				DelegatorAddress: getDelegatorAddr(false),
			},
			validator: stypes.Validator{
				OperatorAddress: valValAddr.String(),
			},
			minRewardWithdrawalAmount: 10,
			expectedResult: func(ctx sdk.Context) *types.Withdrawal {
				return &types.Withdrawal{
					CreationHeight:   uint64(ctx.BlockHeight()),
					ExecutionAddress: delEVMAddr.String(),
					Amount:           uint64(10),
					WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
					ValidatorAddress: valEVMAddr,
				}
			},
		},
		{
			name: "pass: process eligible reward withdrawal - with non-zero claimed reward",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				dk.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardDecCoin, nil)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardCoins, nil)
				bk.EXPECT().SpendableCoin(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: math.NewInt(10)})
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(false), delEVMAddr.String()))
			},
			delegation: stypes.Delegation{
				DelegatorAddress: getDelegatorAddr(false),
			},
			validator: stypes.Validator{
				OperatorAddress: valValAddr.String(),
			},
			minRewardWithdrawalAmount: 10,
			expectedResult: func(ctx sdk.Context) *types.Withdrawal {
				return &types.Withdrawal{
					CreationHeight:   uint64(ctx.BlockHeight()),
					ExecutionAddress: delEVMAddr.String(),
					Amount:           uint64(20), // delegation reward + claimed reward
					WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
					ValidatorAddress: valEVMAddr,
				}
			},
		},
		{
			name: "pass: process eligible reward withdrawal - with non-zero claimed reward and commission",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, sk *moduletestutil.MockStakingKeeper) {
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				dk.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardDecCoin, nil)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardCoins, nil)
				dk.EXPECT().GetValidatorAccumulatedCommission(gomock.Any(), gomock.Any()).Return(dtypes.ValidatorAccumulatedCommission{Commission: delegationRewardDecCoin}, nil)
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).Return(delegationRewardCoins, nil)
				bk.EXPECT().SpendableCoin(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: math.NewInt(10)})
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(true), delEVMAddr.String()))
			},
			delegation: stypes.Delegation{
				DelegatorAddress: getDelegatorAddr(true),
			},
			validator: stypes.Validator{
				OperatorAddress: valValAddr.String(),
			},
			minRewardWithdrawalAmount: 10,
			expectedResult: func(ctx sdk.Context) *types.Withdrawal {
				return &types.Withdrawal{
					CreationHeight:   uint64(ctx.BlockHeight()),
					ExecutionAddress: delEVMAddr.String(),
					Amount:           uint64(30), // delegation reward + claimed reward + commission
					WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
					ValidatorAddress: valEVMAddr,
				}
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, bk, dk, sk, _, esk := createKeeperWithMockStaking(t)

			cachedCtx, _ := ctx.CacheContext()

			if tc.setupMocks != nil {
				tc.setupMocks(bk, dk, sk)
			}

			sk.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec("storyvaloper"))

			if tc.setup != nil {
				tc.setup(cachedCtx, esk)
			}

			// initialize reward withdrawal queue
			require.NoError(t, esk.RewardWithdrawalQueue.Initialize(cachedCtx))

			err := esk.ProcessEligibleRewardWithdrawal(cachedCtx, tc.delegation, tc.validator, tc.minRewardWithdrawalAmount)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)

				if tc.expectedResult != nil {
					expectedWithdrawal := tc.expectedResult(cachedCtx)
					actualWithdrawal, err := esk.RewardWithdrawalQueue.Peek(cachedCtx)
					require.NoError(t, err)
					require.Equal(t, *expectedWithdrawal, actualWithdrawal)
				}
			}
		})
	}
}

func TestEnqueueRewardWithdrawal(t *testing.T) {
	_, accAddrs, valAddrs := createAddresses(2)

	// delegator
	delAccAddr := accAddrs[0]

	// validator
	valAccAddr := accAddrs[1]
	valValAddr := valAddrs[1]
	valEVMAddr, err := utils.Bech32ValidatorAddressToEvmAddress(valValAddr.String())
	require.NoError(t, err)

	defaultDelReward := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(10)))

	type args struct {
		delAddrBech32 string
		valAddrBech32 string
		claimedReward uint64
	}

	getDelegatorAddr := func(isValidator bool) string {
		if isValidator {
			return sdk.AccAddress(valValAddr).String()
		} else {
			return delAccAddr.String()
		}
	}

	createInput := func(claimedAmount uint64, isValidator bool) args {
		return args{
			delAddrBech32: getDelegatorAddr(isValidator),
			valAddrBech32: valValAddr.String(),
			claimedReward: claimedAmount,
		}
	}

	tcs := []struct {
		name           string
		setupMocks     func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper)
		setup          func(ctx sdk.Context, esk *keeper.Keeper)
		input          args
		expectedErr    string
		expectedResult func(ctx sdk.Context) types.Withdrawal
	}{
		{
			name: "fail: invalid validator bech32 address",
			input: args{
				delAddrBech32: delAccAddr.String(),
				valAddrBech32: strings.Replace(valAccAddr.String(), "story", "cosmos", 1),
				claimedReward: uint64(0),
			},
			expectedErr: "validator address from bech32",
		},
		{
			name: "fail: invalid delegator bech32 address - prefix",
			input: args{
				delAddrBech32: strings.Replace(delAccAddr.String(), "story", "cosmos", 1),
				valAddrBech32: valValAddr.String(),
				claimedReward: uint64(0),
			},
			expectedErr: "convert acc address from bech32 address",
		},
		{
			name: "fail: withdraw delegation rewards from distribution keeper",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coins{}, errors.New("failed to withdraw delegation rewards"))
			},
			input:       createInput(0, false),
			expectedErr: "failed to withdraw delegation rewards",
		},
		{
			name: "fail: unknown error from withdrawing validator commission from distribution keeper when the delegator is validator",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coins{}, nil)
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).Return(sdk.Coins{}, errors.New("failed to withdraw validator commission"))
			},
			input:       createInput(0, true),
			expectedErr: "failed to withdraw validator commission",
		},
		{
			name: "fail: no delegator reward address",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coins{}, nil)
			},
			input:       createInput(0, false),
			expectedErr: "map delegator bech32 address to evm reward address",
		},
		{
			name: "fail: send coins from account to module in bank keeper",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coins{}, nil)
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to send coins from account to module"))
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(false), delAccAddr.String()))
			},
			input:       createInput(0, false),
			expectedErr: "failed to send coins from account to module",
		},
		{
			name: "fail: burn coins in bank keeper",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coins{}, nil)
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to burn coins"))
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(false), delAccAddr.String()))
			},
			input:       createInput(0, false),
			expectedErr: "failed to burn coins",
		},
		{
			name: "pass: enqueue reward when delegator is not validator, same recipient",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(defaultDelReward, nil)
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(false), delAccAddr.String()))
			},
			input: createInput(0, false),
			expectedResult: func(ctx sdk.Context) types.Withdrawal {
				return types.Withdrawal{
					CreationHeight:   uint64(ctx.BlockHeight()),
					ExecutionAddress: delAccAddr.String(),
					Amount:           defaultDelReward.AmountOf(sdk.DefaultBondDenom).Uint64(),
					WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
					ValidatorAddress: valEVMAddr,
				}
			},
		},
		{
			name: "pass: enqueue reward when delegator is not validator, different recipient",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(defaultDelReward, nil)
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(false), common.MaxAddress.String()))
			},
			input: createInput(0, false),
			expectedResult: func(ctx sdk.Context) types.Withdrawal {
				return types.Withdrawal{
					CreationHeight:   uint64(ctx.BlockHeight()),
					ExecutionAddress: common.MaxAddress.String(),
					Amount:           defaultDelReward.AmountOf(sdk.DefaultBondDenom).Uint64(),
					WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
					ValidatorAddress: valEVMAddr,
				}
			},
		},
		{
			name: "pass: enqueue reward when delegator is validator and zero commission",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(defaultDelReward, nil)
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).Return(sdk.Coins{}, dtypes.ErrNoValidatorCommission)
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(true), delAccAddr.String()))
			},
			input: createInput(0, true),
			expectedResult: func(ctx sdk.Context) types.Withdrawal {
				return types.Withdrawal{
					CreationHeight:   uint64(ctx.BlockHeight()),
					ExecutionAddress: delAccAddr.String(),
					Amount:           defaultDelReward.AmountOf(sdk.DefaultBondDenom).Uint64(),
					WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
					ValidatorAddress: valEVMAddr,
				}
			},
		},
		{
			name: "pass: enqueue reward when delegator is validator and non-zero commission",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(defaultDelReward, nil)
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).Return(defaultDelReward, nil)
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(true), delAccAddr.String()))
			},
			input: createInput(0, true),
			expectedResult: func(ctx sdk.Context) types.Withdrawal {
				return types.Withdrawal{
					CreationHeight:   uint64(ctx.BlockHeight()),
					ExecutionAddress: delAccAddr.String(),
					Amount:           defaultDelReward.AmountOf(sdk.DefaultBondDenom).Uint64() * 2, // delegation reward + commission
					WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
					ValidatorAddress: valEVMAddr,
				}
			},
		},
		{
			name: "pass: enqueue reward when delegator is validator and non-zero commission with non-zero claimed reward",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(defaultDelReward, nil)
				dk.EXPECT().WithdrawValidatorCommission(gomock.Any(), gomock.Any()).Return(defaultDelReward, nil)
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, getDelegatorAddr(true), delAccAddr.String()))
			},
			input: createInput(10, true),
			expectedResult: func(ctx sdk.Context) types.Withdrawal {
				return types.Withdrawal{
					CreationHeight:   uint64(ctx.BlockHeight()),
					ExecutionAddress: delAccAddr.String(),
					Amount:           defaultDelReward.AmountOf(sdk.DefaultBondDenom).Uint64() * 3, // delegation reward + commission + claimed reward
					WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
					ValidatorAddress: valEVMAddr,
				}
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, bk, dk, _, _, esk := createKeeperWithMockStaking(t)

			if tc.setupMocks != nil {
				tc.setupMocks(bk, dk)
			}

			cachedCtx, _ := ctx.CacheContext()

			if tc.setup != nil {
				tc.setup(cachedCtx, esk)
			}

			// initialize reward queue
			require.NoError(t, esk.RewardWithdrawalQueue.Initialize(cachedCtx))

			err := esk.EnqueueRewardWithdrawal(cachedCtx, tc.input.delAddrBech32, tc.input.valAddrBech32, tc.input.claimedReward)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)

				withdraw, err := esk.RewardWithdrawalQueue.Peek(cachedCtx)
				require.NoError(t, err)
				require.Equal(t, uint64(1), esk.RewardWithdrawalQueue.Len(cachedCtx))
				require.Equal(t, tc.expectedResult(cachedCtx), withdraw)
			}
		})
	}
}

func TestProcessWithdraw(t *testing.T) {
	pubKeys, accAddrs, valAddrs := createAddresses(2)

	// delegator
	delPubKey := pubKeys[0]
	delAddr := accAddrs[0]

	// validator
	valPubKey := pubKeys[1]
	valAddr := valAddrs[1]

	invalidPubKey := append([]byte{0x04}, valPubKey.Bytes()[1:]...)

	createWithdraw := func(amount *big.Int) *bindings.IPTokenStakingWithdraw {
		return &bindings.IPTokenStakingWithdraw{
			Delegator:          common.Address(delAddr),
			ValidatorCmpPubkey: valPubKey.Bytes(),
			StakeAmount:        amount,
			DelegationId:       big.NewInt(0),
			OperatorAddress:    common.Address(delAddr),
		}
	}

	tcs := []struct {
		name              string
		mockStakingKeeper bool
		setupMocks        func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper)
		setup             func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context
		withdraw          *bindings.IPTokenStakingWithdraw
		expectedErr       string
		expectedResult    *stypes.UnbondingDelegation
	}{
		{
			name:              "fail: get singularity height from staking keeper",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), errors.New("failed to get singularity height"))
				}
			},
			withdraw:    createWithdraw(new(big.Int).SetUint64(1)),
			expectedErr: "check if it is singularity",
		},
		{
			name:              "pass(skip): before the singularity",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(1), nil)
				}
			},
			withdraw: createWithdraw(new(big.Int).SetUint64(1)),
		},
		{
			name:              "fail: invalid validator public key - length",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				}
			},
			withdraw: &bindings.IPTokenStakingWithdraw{
				Delegator:          common.Address(delAddr),
				ValidatorCmpPubkey: valPubKey.Bytes()[:16],
				StakeAmount:        new(big.Int).SetUint64(1),
				DelegationId:       big.NewInt(0),
				OperatorAddress:    common.Address(delAddr),
			},
			expectedErr: "validator pubkey to cosmos: invalid pubkey length",
		},
		{
			name:              "fail: invalid validator public key - prefix",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				}
			},
			withdraw: &bindings.IPTokenStakingWithdraw{
				Delegator:          common.Address(delAddr),
				ValidatorCmpPubkey: invalidPubKey,
				StakeAmount:        new(big.Int).SetUint64(1),
				DelegationId:       big.NewInt(0),
				OperatorAddress:    common.Address(delAddr),
			},
			expectedErr: "validator pubkey to evm address",
		},
		{
			name:              "fail: different executor and operator not found",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				}
			},
			withdraw: &bindings.IPTokenStakingWithdraw{
				Delegator:          common.Address(delAddr),
				ValidatorCmpPubkey: valPubKey.Bytes(),
				StakeAmount:        new(big.Int).SetUint64(1),
				DelegationId:       big.NewInt(0),
				OperatorAddress:    common.MaxAddress,
			},
			expectedErr: "invalid unstakeOnBehalf txn, no operator",
		},
		{
			name:              "fail: different executor and not allowed operator",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				}
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				require.NoError(t, esk.DelegatorOperatorAddress.Set(ctx, delAddr.String(), cmpToEVM(delPubKey.Bytes()).String()))

				return ctx
			},
			withdraw: &bindings.IPTokenStakingWithdraw{
				Delegator:          common.Address(delAddr),
				ValidatorCmpPubkey: valPubKey.Bytes(),
				StakeAmount:        new(big.Int).SetUint64(1),
				DelegationId:       big.NewInt(0),
				OperatorAddress:    common.MaxAddress,
			},
			expectedErr: "invalid unstakeOnBehalf txn, not from operator",
		},
		{
			name:              "fail: no depositor account",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(false)

				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				}
			},
			withdraw:    createWithdraw(new(big.Int).SetUint64(1)),
			expectedErr: "depositor account not found",
		},
		{
			name:              "fail: get locked token type",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)

				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
					sk.EXPECT().GetLockedTokenType(gomock.Any()).Return(int32(0), errors.New("failed to get locked token type"))
				}
			},
			withdraw:    createWithdraw(new(big.Int).SetUint64(1)),
			expectedErr: "get locked token type",
		},
		{
			name:              "fail: get validator - not found",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)

				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
					sk.EXPECT().GetLockedTokenType(gomock.Any()).Return(int32(0), nil)
					sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(stypes.Validator{}, stypes.ErrNoValidatorFound)
				}
			},
			withdraw:    createWithdraw(new(big.Int).SetUint64(1)),
			expectedErr: stypes.ErrNoValidatorFound.Error(),
		},
		{
			name:              "fail: get validator - unknown error",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)

				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
					sk.EXPECT().GetLockedTokenType(gomock.Any()).Return(int32(0), nil)
					sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(stypes.Validator{}, errors.New("unknown error"))
				}
			},
			withdraw:    createWithdraw(new(big.Int).SetUint64(1)),
			expectedErr: "unknown error",
		},
		{
			name: "fail: undelegate - undelegation amount less than minimum delegation of staking keeper",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				ctx = ctx.WithBlockHeight(1000000)

				// create validator
				createValidator(t, ctx, sk, valPubKey, valAddr, 0)

				return ctx
			},
			withdraw:    createWithdraw(new(big.Int).SetUint64(1)),
			expectedErr: "undelegation amount is less than the minimum undelegation amount",
		},
		{
			name: "fail: undelegate - no delegation",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				ctx = ctx.WithBlockHeight(1000000)

				// create validator
				createValidator(t, ctx, sk, valPubKey, valAddr, 0)

				return ctx
			},
			withdraw:    createWithdraw(new(big.Int).SetUint64(2)),
			expectedErr: stypes.ErrNoDelegation.Error(),
		},
		{
			name: "fail: undelegate - no period delegation",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				ctx = ctx.WithBlockHeight(1000000)

				// create validator
				createValidator(t, ctx, sk, valPubKey, valAddr, 0)

				// make delegation (not period delegation) to source validator
				createDelegation(t, ctx, sk, delAddr, valAddr, 2, false)

				return ctx
			},
			withdraw:    createWithdraw(new(big.Int).SetUint64(2)),
			expectedErr: stypes.ErrNoPeriodDelegation.Error(),
		},
		{
			name: "pass: withdraw from the delegator account",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				ctx = ctx.WithBlockHeight(1000000)

				// create validator
				createValidator(t, ctx, sk, valPubKey, valAddr, 0)

				// make delegation (not period delegation) to source validator
				createDelegation(t, ctx, sk, delAddr, valAddr, 2, true)

				return ctx
			},
			withdraw: createWithdraw(new(big.Int).SetUint64(2)),
			expectedResult: &stypes.UnbondingDelegation{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Entries: []stypes.UnbondingDelegationEntry{
					{
						Balance: math.NewInt(2),
					},
				},
			},
		},
		{
			name: "pass: withdraw amount exceeds the delegation amount, results in max share withdrawal",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				ctx = ctx.WithBlockHeight(1000000)

				// create validator
				createValidator(t, ctx, sk, valPubKey, valAddr, 0)

				// make delegation (not period delegation) to source validator
				createDelegation(t, ctx, sk, delAddr, valAddr, 2, true)

				return ctx
			},
			withdraw: createWithdraw(new(big.Int).SetUint64(3)),
			expectedResult: &stypes.UnbondingDelegation{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Entries: []stypes.UnbondingDelegationEntry{
					{
						Balance: math.NewInt(2),
					},
				},
			},
		},
		{
			name: "pass: withdraw from the allowed operator",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				ctx = ctx.WithBlockHeight(1000000)

				// set operator
				require.NoError(t, esk.DelegatorOperatorAddress.Set(ctx, delAddr.String(), common.MaxAddress.String()))

				// create validator
				createValidator(t, ctx, sk, valPubKey, valAddr, 0)

				// make delegation (not period delegation) to source validator
				createDelegation(t, ctx, sk, delAddr, valAddr, 2, true)

				return ctx
			},
			withdraw: &bindings.IPTokenStakingWithdraw{
				Delegator:          common.Address(delAddr),
				ValidatorCmpPubkey: valPubKey.Bytes(),
				StakeAmount:        new(big.Int).SetUint64(2),
				DelegationId:       big.NewInt(0),
				OperatorAddress:    common.MaxAddress,
			},
			expectedResult: &stypes.UnbondingDelegation{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Entries: []stypes.UnbondingDelegationEntry{
					{
						Balance: math.NewInt(2),
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var (
				ctx    sdk.Context
				ak     *moduletestutil.MockAccountKeeper
				bk     *moduletestutil.MockBankKeeper
				sk     *skeeper.Keeper
				mockSK *moduletestutil.MockStakingKeeper
				esk    *keeper.Keeper
			)

			if tc.mockStakingKeeper {
				ctx, ak, bk, _, mockSK, _, esk = createKeeperWithMockStaking(t)
			} else {
				ctx, ak, bk, sk, _, esk = createKeeperWithRealStaking(t)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(ak, bk, mockSK)
			}

			cachedCtx, _ := ctx.CacheContext()

			if tc.setup != nil {
				cachedCtx = tc.setup(cachedCtx, sk, esk)
			}

			err := esk.ProcessWithdraw(cachedCtx, tc.withdraw)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)

				// check undelegation exists
				if tc.expectedResult != nil {
					ubd, err := sk.GetUnbondingDelegation(cachedCtx, delAddr, valAddr)
					require.NoError(t, err)
					require.Equal(t, tc.expectedResult.DelegatorAddress, ubd.DelegatorAddress)
					require.Equal(t, tc.expectedResult.ValidatorAddress, ubd.ValidatorAddress)
					require.Equal(t, len(tc.expectedResult.Entries), len(ubd.Entries))

					for i, entry := range ubd.Entries {
						require.Equal(t, tc.expectedResult.Entries[i].Balance.Uint64(), entry.Balance.Uint64())
					}
				}
			}
		})
	}
}

func TestParseWithdraw(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			_, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

			_, err := esk.ParseWithdrawLog(tc.log)
			if tc.expectErr {
				require.Error(t, err, "should return error for %s", tc.name)
			} else {
				require.NoError(t, err, "should not return error for %s", tc.name)
			}
		})
	}
}
