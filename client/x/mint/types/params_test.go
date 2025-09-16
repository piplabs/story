package types_test

import (
	"testing"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/mint/types"
)

func TestParams(t *testing.T) {
	tcs := []struct {
		name        string
		params      types.Params
		expectedErr string
	}{
		{
			name: "fail: blank mint denom",
			params: types.NewParams(
				"",
				math.LegacyNewDec(24625000000000000.000000000000000000),
				uint64(60*60*8766/5)),
			expectedErr: "mint denom cannot be blank",
		},
		{
			name: "fail: invalid mint denom",
			params: types.NewParams(
				"abc#123",
				math.LegacyNewDec(24625000000000000.000000000000000000),
				uint64(60*60*8766/5),
			),
			expectedErr: "mint denom is invalid",
		},
		{
			name: "fail: nil inflations per year",
			params: types.Params{
				MintDenom:     "test",
				BlocksPerYear: uint64(60 * 60 * 8766 / 5),
			},
			expectedErr: "inflations per year cannot be ni",
		},
		{
			name: "fail: negative inflations per year",
			params: types.NewParams(
				"test",
				math.LegacyNewDec(-1),
				uint64(60*60*8766/5),
			),
			expectedErr: "inflations per year cannot be negative",
		},
		{
			name: "fail: zero blocks per year",
			params: types.NewParams(
				"test",
				math.LegacyNewDec(24625000000000000.000000000000000000),
				uint64(0),
			),
			expectedErr: "blocks per year must be positive",
		},
		{
			name:   "pass: valid genesis state",
			params: types.DefaultParams(),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.params.Validate()
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
