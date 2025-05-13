package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/mint/types"
)

func TestParams(t *testing.T) {
	tcs := []struct {
		name        string
		newParams   types.Params
		expectedErr string
	}{
		{
			name: "fail: invalid params",
			newParams: types.NewParams(
				"invalid#",
				math.LegacyNewDec(24625000000000000.000000000000000000),
				uint64(60*60*8766/5),
			),
			expectedErr: "mint denom is invalid",
		},
		{
			name: "pass: valid params",
			newParams: types.NewParams(
				sdk.DefaultBondDenom,
				math.LegacyNewDec(100_000_000),
				uint64(60*60*8766/5),
			),
		},
	}

	for _, tc := range tcs {
		ctx, _, _, mk := createKeeper(t)

		err := mk.SetParams(ctx, tc.newParams)
		if tc.expectedErr != "" {
			require.ErrorContains(t, err, tc.expectedErr)
		} else {
			require.NoError(t, err)

			got, err := mk.GetParams(ctx)
			require.NoError(t, err)
			require.Equal(t, tc.newParams, got)
		}
	}
}
