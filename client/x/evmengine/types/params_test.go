package types_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmengine/types"
)

func TestNewParams(t *testing.T) {
	t.Parallel()

	dummyHash := common.HexToHash("0x047e24c3455107d87c68dffa307b3b7fa1877f3e9d7f30c7ee359f2eff3a75d9")
	tcs := []struct {
		name               string
		executionBlockHash []byte
		expectedResult     types.Params
	}{
		{
			name:               "non-nil execution block hash",
			executionBlockHash: dummyHash.Bytes(),
			expectedResult: types.Params{
				ExecutionBlockHash: dummyHash.Bytes(),
			},
		},
		{
			name:               "nil execution block hash",
			executionBlockHash: nil,
			expectedResult: types.Params{
				ExecutionBlockHash: nil,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := types.NewParams(tc.executionBlockHash)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestDefaultParams(t *testing.T) {
	t.Parallel()

	result := types.DefaultParams()
	require.Equal(t, types.Params{
		ExecutionBlockHash: nil,
	}, result)
}

func TestValidateExecutionBlockHash(t *testing.T) {
	t.Parallel()

	dummyHash := common.HexToHash("0x047e24c3455107d87c68dffa307b3b7fa1877f3e9d7f30c7ee359f2eff3a75d9")
	tcs := []struct {
		name               string
		executionBlockHash []byte
		expectedError      string
	}{
		{
			name:               "pass: valid execution block hash",
			executionBlockHash: dummyHash.Bytes(),
		},
		{
			name:               "fail: invalid execution block hash",
			executionBlockHash: []byte("invalid"),
			expectedError:      "invalid execution block hash length",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := types.ValidateExecutionBlockHash(tc.executionBlockHash)
			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
