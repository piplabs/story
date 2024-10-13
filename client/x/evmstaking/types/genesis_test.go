package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmstaking/types"
)

var zeroVallidatorSweepIndex = &types.ValidatorSweepIndex{
	NextValIndex:    0,
	NextValDelIndex: 0,
}

func TestNewGenesisState(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name                 string
		params               types.Params
		expectedGenesisState *types.GenesisState
	}{
		{
			name:   "default params",
			params: types.DefaultParams(),
			expectedGenesisState: &types.GenesisState{
				Params:              types.DefaultParams(),
				ValidatorSweepIndex: zeroVallidatorSweepIndex,
			},
		},
		{
			name: "custom params",
			params: types.NewParams(
				10,
				20,
				30,
			),
			expectedGenesisState: &types.GenesisState{
				Params: types.NewParams(
					10,
					20,
					30,
				),
				ValidatorSweepIndex: zeroVallidatorSweepIndex,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := types.NewGenesisState(tc.params)
			require.Equal(t, tc.expectedGenesisState, got)
		})
	}
}

func TestDefaultGenesisState(t *testing.T) {
	t.Parallel()
	expectedGenesisState := &types.GenesisState{
		Params:              types.DefaultParams(),
		ValidatorSweepIndex: zeroVallidatorSweepIndex,
	}
	require.Equal(t, expectedGenesisState, types.DefaultGenesisState())
}
