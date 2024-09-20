package types_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmengine/types"
)

func TestNewGenesisState(t *testing.T) {
	t.Parallel()
	dummyExecutionHead := common.HexToHash("0x047e24c3455107d87c68dffa307b3b7fa1877f3e9d7f30c7ee359f2eff3a75d9")
	tcs := []struct {
		name           string
		params         types.Params
		expectedResult *types.GenesisState
	}{
		{
			name: "not nil params",
			params: types.Params{
				ExecutionBlockHash: dummyExecutionHead.Bytes(),
			},
			expectedResult: &types.GenesisState{
				Params: types.Params{
					ExecutionBlockHash: dummyExecutionHead.Bytes(),
				},
			},
		},
		{
			name: "nil execution block hash",
			params: types.Params{
				ExecutionBlockHash: nil,
			},
			expectedResult: &types.GenesisState{
				Params: types.Params{
					ExecutionBlockHash: nil,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := types.NewGenesisState(tc.params)
			require.Equal(t, tc.params, result.Params)
		})
	}
}

func TestDefaultGenesisState(t *testing.T) {
	t.Parallel()
	result := types.DefaultGenesisState()
	require.Equal(t, types.DefaultParams(), result.Params)
}
