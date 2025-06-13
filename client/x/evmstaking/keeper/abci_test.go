package keeper_test

import (
	"strings"
	"testing"

	"cosmossdk.io/math"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmstaking/keeper"
	estestutil "github.com/piplabs/story/client/x/evmstaking/testutil"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"

	"go.uber.org/mock/gomock"
)

func TestEndBlock(t *testing.T) {
	pubKeys, _, valAddrs := createAddresses(3)

	// delegator
	delPubKey := pubKeys[0]
	// delAccAddr := accAddrs[0]
	delEVMAddr, err := keeper.CmpPubKeyToEVMAddress(delPubKey.Bytes())
	require.NoError(t, err)
	delAccAddrFromEVM := sdk.AccAddress(delEVMAddr.Bytes())

	// validator
	valPubKey := pubKeys[1]
	valValAddr := valAddrs[1]
	valCosmosPubKey, err := k1util.PubKeyToCosmos(valPubKey)
	require.NoError(t, err)
	valEVMAddr, err := keeper.CmpPubKeyToEVMAddress(valPubKey.Bytes())
	require.NoError(t, err)
	// invalidPubKey := append([]byte{0x04}, valPubKey.Bytes()[1:]...)

	type expectedResult struct {
		withdrawals       []types.Withdrawal
		rewardWithdrawals []types.Withdrawal
	}

	tcs := []struct {
		name           string
		setupMocks     func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, dk *estestutil.MockDistributionKeeper, sk *estestutil.MockStakingKeeper)
		setup          func(ctx sdk.Context, esk *keeper.Keeper) sdk.Context
		ubdEntries     []stypes.UnbondedEntry
		expectedErr    string
		expectedResult func(ctx sdk.Context) expectedResult
	}{
		{
			name: "fail: get singularity height from staking keeper",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, dk *estestutil.MockDistributionKeeper, sk *estestutil.MockStakingKeeper) {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), errors.New("failed to get singularity height"))
			},
			expectedErr: "failed to get singularity height",
		},
		{
			name: "pass(skip): skip EndBlock within the singularity",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, dk *estestutil.MockDistributionKeeper, sk *estestutil.MockStakingKeeper) {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(10), nil)
			},
		},
		{
			name: "fail: error from EndBlockerWithUnbondedeEntries from staking keeper",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, dk *estestutil.MockDistributionKeeper, sk *estestutil.MockStakingKeeper) {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				sk.EXPECT().EndBlockerWithUnbondedEntries(gomock.Any()).Return([]abci.ValidatorUpdate{}, []stypes.UnbondedEntry{}, errors.New("failed end blocker"))
			},
			expectedErr: "process staking EndBlocker",
		},
		{
			name: "fail: process unstake withdrawals",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, dk *estestutil.MockDistributionKeeper, sk *estestutil.MockStakingKeeper) {
				invalidUBD := []stypes.UnbondedEntry{
					{
						DelegatorAddress: strings.Replace(delAccAddrFromEVM.String(), "story", "cosmos", 1),
						ValidatorAddress: valValAddr.String(),
						Amount:           math.NewInt(10),
					},
				}
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				sk.EXPECT().EndBlockerWithUnbondedEntries(gomock.Any()).Return([]abci.ValidatorUpdate{}, invalidUBD, nil)
			},
			expectedErr: "process unstake withdrawals",
		},
		{
			name: "pass: process unstake withdrawals",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, dk *estestutil.MockDistributionKeeper, sk *estestutil.MockStakingKeeper) {
				ubd := []stypes.UnbondedEntry{
					{
						DelegatorAddress: delAccAddrFromEVM.String(),
						ValidatorAddress: valValAddr.String(),
						Amount:           math.NewInt(10),
					},
				}
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				sk.EXPECT().EndBlockerWithUnbondedEntries(gomock.Any()).Return([]abci.ValidatorUpdate{}, ubd, nil)
				sk.EXPECT().GetDelegation(gomock.Any(), gomock.Any(), gomock.Any()).Return(stypes.Delegation{}, nil)
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{}, nil)
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SpendableCoin(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(0)))
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(0), nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) sdk.Context {
				require.NoError(t, esk.DelegatorWithdrawAddress.Set(ctx, delAccAddrFromEVM.String(), delEVMAddr.String()))

				return ctx
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
			name: "fail: process reward withdrawals",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, dk *estestutil.MockDistributionKeeper, sk *estestutil.MockStakingKeeper) {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				sk.EXPECT().EndBlockerWithUnbondedEntries(gomock.Any()).Return([]abci.ValidatorUpdate{}, []stypes.UnbondedEntry{}, nil)
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{}, errors.New("failed to get all validators"))
			},
			expectedErr: "process reward withdrawals",
		},
		{
			name: "pass: process reward withdrawals",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, dk *estestutil.MockDistributionKeeper, sk *estestutil.MockStakingKeeper) {
				validator, err := stypes.NewValidator(valValAddr.String(), valCosmosPubKey, stypes.Description{Moniker: "test"}, 0)
				require.NoError(t, err)
				delegation := stypes.NewDelegation(delAccAddrFromEVM.String(), valValAddr.String(), math.LegacyNewDec(10), math.LegacyNewDec(10).Quo(math.LegacyNewDec(2)))

				delegationRewardDec := sdk.NewDecCoins(sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, math.LegacyNewDec(10)))
				delegationRewardCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(10))}

				sk.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec("storyvaloper")).AnyTimes()
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				sk.EXPECT().EndBlockerWithUnbondedEntries(gomock.Any()).Return([]abci.ValidatorUpdate{}, []stypes.UnbondedEntry{}, nil)
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{validator}, nil)
				sk.EXPECT().GetValidatorDelegations(gomock.Any(), gomock.Any()).Return([]stypes.Delegation{delegation}, nil)
				dk.EXPECT().IncrementValidatorPeriod(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
				dk.EXPECT().CalculateDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardDec, nil)
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(0), nil)
				dk.EXPECT().WithdrawDelegationRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(delegationRewardCoins, nil)
				bk.EXPECT().SpendableCoin(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{Amount: math.NewInt(0)})
				bk.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) sdk.Context {
				newParams := types.Params{
					MaxWithdrawalPerBlock:      2,
					MaxSweepPerBlock:           4,
					MinPartialWithdrawalAmount: 2,
				}
				require.NoError(t, esk.SetParams(ctx, newParams))

				require.NoError(t, esk.DelegatorRewardAddress.Set(ctx, delAccAddrFromEVM.String(), delEVMAddr.String()))

				return ctx
			},
			expectedResult: func(ctx sdk.Context) expectedResult {
				return expectedResult{
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
			name: "fail: process UBI withdrawal",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, dk *estestutil.MockDistributionKeeper, sk *estestutil.MockStakingKeeper) {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				sk.EXPECT().EndBlockerWithUnbondedEntries(gomock.Any()).Return([]abci.ValidatorUpdate{}, []stypes.UnbondedEntry{}, nil)
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{}, nil)
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(0), errors.New("failed to get UBI balance by denom"))
			},
			expectedErr: "process ubi withdrawal",
		},
		{
			name: "pass: process UBI withdrawal",
			setupMocks: func(ak *estestutil.MockAccountKeeper, bk *estestutil.MockBankKeeper, dk *estestutil.MockDistributionKeeper, sk *estestutil.MockStakingKeeper) {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				sk.EXPECT().EndBlockerWithUnbondedEntries(gomock.Any()).Return([]abci.ValidatorUpdate{}, []stypes.UnbondedEntry{}, nil)
				sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stypes.Validator{}, nil)
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(10), nil)
				dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: math.NewInt(0)}, nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, esk *keeper.Keeper) sdk.Context {
				newParams := types.Params{
					MaxWithdrawalPerBlock:      2,
					MaxSweepPerBlock:           4,
					MinPartialWithdrawalAmount: 2,
					UbiWithdrawAddress:         common.MaxAddress.String(),
				}
				require.NoError(t, esk.SetParams(ctx, newParams))

				return ctx
			},
			expectedResult: func(ctx sdk.Context) expectedResult {
				return expectedResult{
					withdrawals: []types.Withdrawal{
						{
							CreationHeight:   uint64(ctx.BlockHeight()),
							ExecutionAddress: common.MaxAddress.String(),
							Amount:           10,
							WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_UBI,
							ValidatorAddress: "",
						},
					},
				}
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, ak, bk, dk, sk, _, esk := createKeeperWithMockStaking(t)

			if tc.setupMocks != nil {
				tc.setupMocks(ak, bk, dk, sk)
			}

			cachedCtx, _ := ctx.CacheContext()

			if tc.setup != nil {
				cachedCtx = tc.setup(cachedCtx, esk)
			}

			_, err := esk.EndBlock(cachedCtx)
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
