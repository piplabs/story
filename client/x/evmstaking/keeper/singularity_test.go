package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	moduletestutil "github.com/piplabs/story/client/x/evmstaking/testutil"
	"github.com/piplabs/story/lib/errors"

	"go.uber.org/mock/gomock"
)

func TestIsSingularity(t *testing.T) {
	tcs := []struct {
		name           string
		setupMock      func(ctx sdk.Context, sk *moduletestutil.MockStakingKeeper) sdk.Context
		expectedResult bool
		expectedErr    bool
	}{
		{
			name: "fail: get singularity height from staking keeper",
			setupMock: func(ctx sdk.Context, sk *moduletestutil.MockStakingKeeper) sdk.Context {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), errors.New("failed to get the singularity height"))

				return ctx
			},
			expectedErr: true,
		},
		{
			name: "pass: is within the singularity",
			setupMock: func(ctx sdk.Context, sk *moduletestutil.MockStakingKeeper) sdk.Context {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(10), nil)

				return ctx
			},
			expectedResult: true,
		},
		{
			name: "pass: is after the singularity",
			setupMock: func(ctx sdk.Context, sk *moduletestutil.MockStakingKeeper) sdk.Context {
				sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(10), nil)

				ctx = ctx.WithBlockHeight(11)

				return ctx
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, sk, _, esk := createKeeperWithMockStaking(t)
			ctx = tc.setupMock(ctx, sk)

			isSingularity, err := esk.IsSingularity(ctx)

			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, isSingularity)
			}
		})
	}
}
