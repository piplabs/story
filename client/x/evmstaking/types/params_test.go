package types_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/suite"

	"github.com/piplabs/story/client/x/evmstaking/types"
)

type ParamsTestSuite struct {
	suite.Suite
	encConf testutil.TestEncodingConfig
}

func (suite *ParamsTestSuite) SetupTest() {
	suite.encConf = testutil.MakeTestEncodingConfig()
}

func (suite *ParamsTestSuite) TestNewParams() {
	require := suite.Require()
	tcs := []struct {
		name                       string
		maxWithdrawalPerBlock      uint32
		maxSweepPerBlock           uint32
		minPartialWithdrawalAmount uint64
		expectedError              string
	}{
		{
			name:                       "pass: valid params",
			maxWithdrawalPerBlock:      uint32(1),
			maxSweepPerBlock:           uint32(2),
			minPartialWithdrawalAmount: uint64(3),
		},
		{
			name:                       "fail: invalid maxWithdrawalPerBlock",
			maxWithdrawalPerBlock:      uint32(0),
			maxSweepPerBlock:           uint32(2),
			minPartialWithdrawalAmount: uint64(3),
			expectedError:              "max withdrawal per block must be positive: 0",
		},
		{
			name:                       "fail: invalid maxSweepPerBlock - 0 of maxSweepPerBlock",
			maxWithdrawalPerBlock:      uint32(1),
			maxSweepPerBlock:           uint32(0),
			minPartialWithdrawalAmount: uint64(3),
			expectedError:              "max sweep per block must be positive: 0",
		},
		{
			name:                       "fail: invalid maxSweepPerBlock - smaller maxSweepPerBlock than maxWithdrawalPerBlock",
			maxWithdrawalPerBlock:      uint32(2),
			maxSweepPerBlock:           uint32(1),
			minPartialWithdrawalAmount: uint64(3),
			expectedError:              "max sweep per block must be greater than or equal to max withdrawal per block: 1 < 2",
		},
		{
			name:                       "fail: invalid minPartialWithdrawalAmount",
			maxWithdrawalPerBlock:      uint32(1),
			maxSweepPerBlock:           uint32(2),
			minPartialWithdrawalAmount: uint64(0),
			expectedError:              "min partial withdrawal amount must be positive: 0",
		},
	}

	for _, tc := range tcs {
		suite.Run(tc.name, func() {
			params := types.NewParams(tc.maxWithdrawalPerBlock, tc.maxSweepPerBlock, tc.minPartialWithdrawalAmount)

			err := params.Validate()
			if tc.expectedError == "" {
				require.NoError(err)
				require.Equal(tc.maxWithdrawalPerBlock, params.MaxWithdrawalPerBlock)
				require.Equal(tc.maxSweepPerBlock, params.MaxSweepPerBlock)
				require.Equal(tc.minPartialWithdrawalAmount, params.MinPartialWithdrawalAmount)
			} else {
				require.ErrorContains(err, tc.expectedError, "unexpected error from params validation")
			}
		})
	}
}

func (suite *ParamsTestSuite) TestDefaultParams() {
	require := suite.Require()
	params := types.DefaultParams()
	// check values are set correctly
	require.Equal(types.DefaultMaxWithdrawalPerBlock, params.MaxWithdrawalPerBlock)
	require.Equal(types.DefaultMaxSweepPerBlock, params.MaxSweepPerBlock)
	require.Equal(types.DefaultMinPartialWithdrawalAmount, params.MinPartialWithdrawalAmount)
}

func (suite *ParamsTestSuite) TestValidateMaxWithdrawalPerBlock() {
	require := suite.Require()

	tcs := []struct {
		name        string
		input       uint32
		expectedErr string
	}{
		{
			name:  "valid value",
			input: 1,
		},
		{
			name:        "invalid value",
			input:       0,
			expectedErr: "max withdrawal per block must be positive: 0",
		},
	}

	for _, tc := range tcs {
		suite.Run(tc.name, func() {
			err := types.ValidateMaxWithdrawalPerBlock(tc.input)
			if tc.expectedErr == "" {
				require.NoError(err)
			} else {
				require.Error(err)
				require.Contains(err.Error(), tc.expectedErr)
			}
		})
	}
}

func (suite *ParamsTestSuite) TestValidateMaxSweepPerBlock() {
	require := suite.Require()

	tcs := []struct {
		name                  string
		maxSweepPerBlock      uint32
		maxWithdrawalPerBlock uint32
		expectedErr           string
	}{
		{
			name:                  "valid value",
			maxSweepPerBlock:      2,
			maxWithdrawalPerBlock: 1,
		},
		{
			name:                  "valid value",
			maxSweepPerBlock:      1,
			maxWithdrawalPerBlock: 1,
		},
		{
			name:                  "invalid value",
			maxSweepPerBlock:      0,
			maxWithdrawalPerBlock: 2,
			expectedErr:           "max sweep per block must be positive: 0",
		},
		{
			name:                  "invalid value",
			maxSweepPerBlock:      1,
			maxWithdrawalPerBlock: 2,
			expectedErr:           "max sweep per block must be greater than or equal to max withdrawal per block",
		},
	}

	for _, tc := range tcs {
		suite.Run(tc.name, func() {
			err := types.ValidateMaxSweepPerBlock(tc.maxSweepPerBlock, tc.maxWithdrawalPerBlock)
			if tc.expectedErr == "" {
				require.NoError(err)
			} else {
				require.Error(err)
				require.Contains(err.Error(), tc.expectedErr)
			}
		})
	}
}

func (suite *ParamsTestSuite) TestValidateMinPartialWithdrawatAmount() {
	require := suite.Require()

	tcs := []struct {
		name        string
		input       uint64
		expectedErr string
	}{
		{
			name:  "valid value",
			input: 1,
		},
		{
			name:        "invalid value",
			input:       0,
			expectedErr: "min partial withdrawal amount must be positive: 0",
		},
	}

	for _, tc := range tcs {
		suite.Run(tc.name, func() {
			err := types.ValidateMinPartialWithdrawalAmount(tc.input)
			if tc.expectedErr == "" {
				require.NoError(err)
			} else {
				require.Error(err)
				require.Contains(err.Error(), tc.expectedErr)
			}
		})
	}
}

func TestParamsTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ParamsTestSuite))
}
