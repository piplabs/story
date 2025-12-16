package types_test

import (
	"context"
	"testing"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/mint/types"
)

func TestNewGenesisState(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name                 string
		params               types.Params
		expectedGenesisState *types.GenesisState
		expectedInflation    math.LegacyDec
	}{
		{
			name:   "default params",
			params: types.DefaultParams(),
			expectedGenesisState: &types.GenesisState{
				Params: types.DefaultParams(),
			},
			expectedInflation: math.LegacyMustNewDecFromStr("3901595812.102314497933936675"),
		},
		{
			name: "custom params",
			params: types.NewParams(
				"test",
				math.LegacyNewDec(100_000_000),
				uint64(60*60*24*365),
			),
			expectedGenesisState: &types.GenesisState{
				Params: types.NewParams(
					"test",
					math.LegacyNewDec(100_000_000),
					uint64(60*60*24*365),
				),
			},
			expectedInflation: math.LegacyMustNewDecFromStr("3.170979198376458650"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// check genesis state
			got := types.NewGenesisState(tc.params)
			require.Equal(t, tc.expectedGenesisState, got)

			// check inflation
			calculated := types.DefaultInflationCalculationFn(context.Background(), tc.params, math.LegacyDec{})
			require.Equal(t, tc.expectedInflation, calculated)
		})
	}
}

func TestValidateGenesis(t *testing.T) {
	tcs := []struct {
		name        string
		gs          types.GenesisState
		expectedErr string
	}{
		{
			name: "fail: blank mint denom",
			gs: types.GenesisState{
				Params: types.Params{
					MintDenom:         "",
					InflationsPerYear: math.LegacyNewDec(24625000000000000.000000000000000000),
					BlocksPerYear:     uint64(60 * 60 * 8766 / 5),
				},
			},
			expectedErr: "mint denom cannot be blank",
		},
		{
			name: "fail: invalid mint denom",
			gs: types.GenesisState{
				Params: types.Params{
					MintDenom:         "abc#123",
					InflationsPerYear: math.LegacyNewDec(24625000000000000.000000000000000000),
					BlocksPerYear:     uint64(60 * 60 * 8766 / 5),
				},
			},
			expectedErr: "mint denom is invalid",
		},
		{
			name: "fail: nil inflations per year",
			gs: types.GenesisState{
				Params: types.Params{
					MintDenom:     "test",
					BlocksPerYear: uint64(60 * 60 * 8766 / 5),
				},
			},
			expectedErr: "inflations per year cannot be ni",
		},
		{
			name: "fail: negative inflations per year",
			gs: types.GenesisState{
				Params: types.Params{
					MintDenom:         "test",
					InflationsPerYear: math.LegacyNewDec(-1),
					BlocksPerYear:     uint64(60 * 60 * 8766 / 5),
				},
			},
			expectedErr: "inflations per year cannot be negative",
		},
		{
			name: "fail: zero blocks per year",
			gs: types.GenesisState{
				Params: types.Params{
					MintDenom:         "test",
					InflationsPerYear: math.LegacyNewDec(24625000000000000.000000000000000000),
					BlocksPerYear:     uint64(0),
				},
			},
			expectedErr: "blocks per year must be positive",
		},
		{
			name: "pass: valid genesis state",
			gs: types.GenesisState{
				Params: types.DefaultParams(),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := types.ValidateGenesis(tc.gs)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
