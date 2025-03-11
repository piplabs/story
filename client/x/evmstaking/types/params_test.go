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
	maxWithdrawalPerBlock, maxSweepPerBlock, minPartialWithdrawalAmount, refundFee, refundPeriod := uint32(1), uint32(2), uint64(3), uint32(4), uint32(5)
	params := types.NewParams(maxWithdrawalPerBlock, maxSweepPerBlock, minPartialWithdrawalAmount, refundFee, refundPeriod)
	// check values are set correctly
	require.Equal(maxWithdrawalPerBlock, params.MaxWithdrawalPerBlock)
	require.Equal(maxSweepPerBlock, params.MaxSweepPerBlock)
	require.Equal(minPartialWithdrawalAmount, params.MinPartialWithdrawalAmount)
	require.Equal(refundFee, params.RefundFeeBps)
	require.Equal(refundPeriod, params.RefundPeriod)
}

func (suite *ParamsTestSuite) TestDefaultParams() {
	require := suite.Require()
	params := types.DefaultParams()
	// check values are set correctly
	require.Equal(types.DefaultMaxWithdrawalPerBlock, params.MaxWithdrawalPerBlock)
	require.Equal(types.DefaultMaxSweepPerBlock, params.MaxSweepPerBlock)
	require.Equal(types.DefaultMinPartialWithdrawalAmount, params.MinPartialWithdrawalAmount)
	require.Equal(types.DefaultRefundFeeBps, params.RefundFeeBps)
	require.Equal(types.DefaultRefundPeriod, params.RefundPeriod)
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

func (suite *ParamsTestSuite) TestValidateRefundFeePercentage() {
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
			expectedErr: "refund fee bps must be positive: 0",
		},
		{
			name:        "invalid value",
			input:       10001,
			expectedErr: "refund fee bps must be less than or equal to 10000bps (100%): 10001",
		},
	}

	for _, tc := range tcs {
		suite.Run(tc.name, func() {
			err := types.ValidateRefundFeeBps(tc.input)
			if tc.expectedErr == "" {
				require.NoError(err)
			} else {
				require.Error(err)
				require.Contains(err.Error(), tc.expectedErr)
			}
		})
	}
}

func (suite *ParamsTestSuite) TestValidateRefundPeriod() {
	require := suite.Require()

	tcs := []struct {
		name        string
		input       uint32
		expectedErr string
	}{
		{
			name:  "valid value",
			input: 24,
		},
		{
			name:        "invalid value",
			input:       0,
			expectedErr: "refund period must be at least 24 hours: 0",
		},
		{
			name:        "invalid value",
			input:       10,
			expectedErr: "refund period must be at least 24 hours: 10",
		},
	}

	for _, tc := range tcs {
		suite.Run(tc.name, func() {
			err := types.ValidateRefundPeriod(tc.input)
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
