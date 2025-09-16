package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/math"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	minttestutil "github.com/piplabs/story/client/x/mint/testutil"
	"github.com/piplabs/story/client/x/mint/types"
	"github.com/piplabs/story/lib/errors"

	"go.uber.org/mock/gomock"
)

func TestBeginBlocker(t *testing.T) {
	var expectedMintAmount = types.DefaultInflationCalculationFn(context.Background(), types.DefaultParams(), math.LegacyNewDec(0)).TruncateInt()

	tcs := []struct {
		name           string
		setupMocks     func(bk *minttestutil.MockBankKeeper, sk *minttestutil.MockStakingKeeper)
		expectedErr    string
		expectedResult *sdk.Event
	}{
		{
			name: "fail: get singularity height",
			setupMocks: func(bk *minttestutil.MockBankKeeper, sk *minttestutil.MockStakingKeeper) {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), errors.New("failed to get singularity height"))
			},
			expectedErr: "failed to get singularity height",
		},
		{
			name: "pass(skip): no mint before singularity",
			setupMocks: func(bk *minttestutil.MockBankKeeper, sk *minttestutil.MockStakingKeeper) {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(10), nil)
			},
		},
		{
			name: "fail: mint coins from bank keeper",
			setupMocks: func(bk *minttestutil.MockBankKeeper, sk *minttestutil.MockStakingKeeper) {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to mint coins"))
			},
			expectedErr: "failed to mint coins",
		},
		{
			name: "fail: add collected fees",
			setupMocks: func(bk *minttestutil.MockBankKeeper, sk *minttestutil.MockStakingKeeper) {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to send coins from module to module"))
			},
			expectedErr: "failed to send coins from module to module",
		},
		{
			name: "pass: mint coins successfully",
			setupMocks: func(bk *minttestutil.MockBankKeeper, sk *minttestutil.MockStakingKeeper) {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResult: &sdk.Event{
				Type: types.EventTypeMint,
				Attributes: []abci.EventAttribute{
					{
						Key:   sdk.AttributeKeyAmount,
						Value: expectedMintAmount.String(),
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, bk, sk, mk := createKeeper(t)

			if tc.setupMocks != nil {
				tc.setupMocks(bk, sk)
			}

			cachedCtx, _ := ctx.CacheContext()

			err := mk.BeginBlocker(cachedCtx, types.DefaultInflationCalculationFn)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)

				if tc.expectedResult != nil {
					ev := cachedCtx.EventManager().Events()[0]
					require.Equal(t, *tc.expectedResult, ev)
				}
			}
		})
	}
}
