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
		setupMocks     func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper)
		expectedErr    string
		expectedResult *types.Withdrawal
	}{
		{
			name: "fail: get UBI balance by denom from distribution keeper",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.Int{}, errors.New("failed to get UBI balance by denom"))
			},
			expectedErr: "get ubi balance by denom",
		},
		{
			name: "pass(skip): balance of UBI is less than minimum partial withdrawal amount",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(1), nil)
			},
		},
		{
			name: "fail: withdraw UBI from distribution keeper to evmstaking keeper",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(8000000001), nil)
				dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.Coin{}, errors.New("failed to withdraw UBI by denom to module"))
			},
			expectedErr: "withdraw ubi by denom to module",
		},
		{
			name: "fail: burn coins from bank keeper",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(8000000001), nil)
				dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.NewCoin("test", math.NewInt(1)), nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to burn coins"))
			},
			expectedErr: "burn ubi coins",
		},
		{
			name: "pass: add UBI withdrawal to the withdrawal queue",
			setupMocks: func(bk *moduletestutil.MockBankKeeper, dk *moduletestutil.MockDistributionKeeper) {
				dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), gomock.Any()).Return(math.NewInt(8000000001), nil)
				dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.NewCoin("test", math.NewInt(1)), nil)
				bk.EXPECT().BurnCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResult: &types.Withdrawal{
				ExecutionAddress: "",
				Amount:           uint64(8000000001),
				WithdrawalType:   types.WithdrawalType_WITHDRAWAL_TYPE_UBI,
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
