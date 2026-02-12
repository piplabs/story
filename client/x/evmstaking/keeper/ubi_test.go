package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	moduletestutil "github.com/piplabs/story/client/x/evmstaking/testutil"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"

	"go.uber.org/mock/gomock"
)

func TestProcessUbiWithdrawal(t *testing.T) {
	tcs := []struct {
		name           string
		setupMocks     func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, dkgk *moduletestutil.MockDKGKeeper)
		expectedErr    string
		expectedResult *types.Withdrawal
	}{
		{
			name: "fail: claim DKG settlement balance fails",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, dkgk *moduletestutil.MockDKGKeeper) {
				dkgk.EXPECT().ClaimSettlementBalance(gomock.Any(), types.ModuleName).Return(math.ZeroInt(), errors.New("claim failed"))
			},
			expectedErr: "claim DKG settlement balance",
		},
		{
			name: "fail: get UBI balance by denom from distribution keeper",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, dkgk *moduletestutil.MockDKGKeeper) {
				dkgk.EXPECT().ClaimSettlementBalance(gomock.Any(), types.ModuleName).Return(math.ZeroInt(), nil)
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.Int{}, errors.New("failed to get UBI balance by denom"))
			},
			expectedErr: "get ubi balance by denom",
		},
		{
			name: "pass(skip): no settlement and balance below threshold",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, dkgk *moduletestutil.MockDKGKeeper) {
				dkgk.EXPECT().ClaimSettlementBalance(gomock.Any(), types.ModuleName).Return(math.ZeroInt(), nil)
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(1), nil)
			},
		},
		{
			name: "pass: settlement only, balance below threshold",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, dkgk *moduletestutil.MockDKGKeeper) {
				dkgk.EXPECT().ClaimSettlementBalance(gomock.Any(), types.ModuleName).Return(math.NewInt(500), nil)
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(1), nil)
				bk.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
			},
			expectedResult: &types.Withdrawal{
				ExecutionAddress: "",
				Amount:           uint64(500),
				WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_UBI,
			},
		},
		{
			name: "fail: withdraw UBI from distribution keeper to evmstaking keeper",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, dkgk *moduletestutil.MockDKGKeeper) {
				dkgk.EXPECT().ClaimSettlementBalance(gomock.Any(), types.ModuleName).Return(math.ZeroInt(), nil)
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(8000000001), nil)
				dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{}, errors.New("failed to withdraw UBI by denom to module"))
			},
			expectedErr: "withdraw ubi by denom to module",
		},
		{
			name: "fail: distribute DKG committee rewards fails",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, dkgk *moduletestutil.MockDKGKeeper) {
				dkgk.EXPECT().ClaimSettlementBalance(gomock.Any(), types.ModuleName).Return(math.ZeroInt(), nil)
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(8000000001), nil)
				dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(8000000001)), nil)
				dkgk.EXPECT().DistributeRewardsToActiveCommittee(gomock.Any(), types.ModuleName, math.NewInt(8000000001)).Return(math.ZeroInt(), errors.New("distribute failed"))
			},
			expectedErr: "distribute DKG committee rewards",
		},
		{
			name: "fail: burn coins from bank keeper",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, dkgk *moduletestutil.MockDKGKeeper) {
				dkgk.EXPECT().ClaimSettlementBalance(gomock.Any(), types.ModuleName).Return(math.ZeroInt(), nil)
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(8000000001), nil)
				dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(8000000001)), nil)
				dkgk.EXPECT().DistributeRewardsToActiveCommittee(gomock.Any(), types.ModuleName, math.NewInt(8000000001)).Return(math.NewInt(1000), nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to burn coins"))
			},
			expectedErr: "burn ubi coins",
		},
		{
			name: "pass: regular withdrawal with DKG committee rewards",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, dkgk *moduletestutil.MockDKGKeeper) {
				dkgk.EXPECT().ClaimSettlementBalance(gomock.Any(), types.ModuleName).Return(math.ZeroInt(), nil)
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(8000000001), nil)
				dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(8000000001)), nil)
				dkgk.EXPECT().DistributeRewardsToActiveCommittee(gomock.Any(), types.ModuleName, math.NewInt(8000000001)).Return(math.NewInt(1000), nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResult: &types.Withdrawal{
				ExecutionAddress: "",
				Amount:           uint64(8000000001 - 1000),
				WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_UBI,
			},
		},
		{
			name: "pass: regular withdrawal plus settlement balance",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper, dkgk *moduletestutil.MockDKGKeeper) {
				dkgk.EXPECT().ClaimSettlementBalance(gomock.Any(), types.ModuleName).Return(math.NewInt(500), nil)
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(8000000001), nil)
				dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(8000000001)), nil)
				dkgk.EXPECT().DistributeRewardsToActiveCommittee(gomock.Any(), types.ModuleName, math.NewInt(8000000001)).Return(math.NewInt(1000), nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResult: &types.Withdrawal{
				ExecutionAddress: "",
				Amount:           uint64(8000000001 - 1000 + 500),
				WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_UBI,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, bk, dk, _, _, dkgk, esk := createKeeperWithMockStaking(t)

			if tc.setupMocks != nil {
				tc.setupMocks(bk, dk, dkgk)
			}

			cachedCtx, _ := ctx.CacheContext()

			// initialize withdrawal queue
			require.NoError(t, esk.WithdrawalQueue.Initialize(cachedCtx))

			err := esk.ProcessUbiWithdrawal(cachedCtx)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				require.Equal(t, uint64(0), esk.WithdrawalQueue.Len(cachedCtx))
			} else {
				require.NoError(t, err)

				if tc.expectedResult != nil {
					ubiWithdrawal, err := esk.WithdrawalQueue.Peek(cachedCtx)
					require.NoError(t, err)
					require.Equal(t, *tc.expectedResult, ubiWithdrawal)
					require.Equal(t, uint64(1), esk.WithdrawalQueue.Len(cachedCtx))
				} else {
					require.Equal(t, uint64(0), esk.WithdrawalQueue.Len(cachedCtx))
				}
			}
		})
	}
}
