package types_test

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/piplabs/story/client/x/evmstaking/types"
)

type WithdrawTestSuite struct {
	suite.Suite
	encConf testutil.TestEncodingConfig
	delAddr string
	valAddr string
	evmAddr common.Address
}

func (suite *WithdrawTestSuite) SetupTest() {
	suite.encConf = testutil.MakeTestEncodingConfig()
	// set dummy addresses
	suite.delAddr = "story1hmjw3pvkjtndpg8wqppwdn8udd835qpan4hm0y"
	suite.valAddr = "storyvaloper1hmjw3pvkjtndpg8wqppwdn8udd835qpaa6r6y0"
	suite.evmAddr = common.HexToAddress("0x131D25EDE18178BAc9275b312001a63C081722d2")
}

func (suite *WithdrawTestSuite) TestString() {
	require := suite.Require()
	// Define the test cases
	testCases := []struct {
		name           string
		withdrawal     types.Withdrawal
		expectedString string
	}{
		{
			name:           "Normal values",
			withdrawal:     types.NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 100),
			expectedString: fmt.Sprintf("creation_height:1 delegator_address:\"%s\" validator_address:\"%s\" execution_address:\"%s\" amount:100 ", suite.delAddr, suite.valAddr, suite.evmAddr.String()),
		},
		{
			name:           "Empty addresses",
			withdrawal:     types.NewWithdrawal(1, "", "", "", 1),
			expectedString: "creation_height:1 amount:1 ",
		},
		{
			name:           "Large amount",
			withdrawal:     types.NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), math.MaxUint64),
			expectedString: fmt.Sprintf("creation_height:1 delegator_address:\"%s\" validator_address:\"%s\" execution_address:\"%s\" amount:%s ", suite.delAddr, suite.valAddr, suite.evmAddr.String(), new(big.Int).SetUint64(math.MaxUint64).String()),
		},
		{
			name:           "Zero amount",
			withdrawal:     types.NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 0),
			expectedString: fmt.Sprintf("creation_height:1 delegator_address:\"%s\" validator_address:\"%s\" execution_address:\"%s\" ", suite.delAddr, suite.valAddr, suite.evmAddr.String()),
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			result := tc.withdrawal.String()
			require.Equal(tc.expectedString, result)
		})
	}
}

func (suite *WithdrawTestSuite) TestWithdrawalsString() {
	require := suite.Require()
	ws := types.Withdrawals{
		Withdrawals: []types.Withdrawal{
			types.NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1),
			types.NewWithdrawal(2, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 2),
		},
	}

	expectedString := fmt.Sprintf(
		"creation_height:1 delegator_address:\"%s\" validator_address:\"%s\" execution_address:\"%s\" amount:1 \n"+
			"creation_height:2 delegator_address:\"%s\" validator_address:\"%s\" execution_address:\"%s\" amount:2",
		suite.delAddr, suite.valAddr, suite.evmAddr.String(),
		suite.delAddr, suite.valAddr, suite.evmAddr.String(),
	)
	require.Equal(expectedString, ws.String())

	// Test with spaces, it should trim the spaces
	stringWithSpaces := " " + ws.String() + " " // add leading and trailing spaces
	trimmedString := strings.TrimSpace(stringWithSpaces)
	require.Equal(expectedString, trimmedString, "Withdrawals.String() should trim spaces")
}

func (suite *WithdrawTestSuite) TestWithdrawalsLen() {
	require := suite.Require()
	ws := types.Withdrawals{
		Withdrawals: []types.Withdrawal{
			types.NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1),
			types.NewWithdrawal(2, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 2),
		},
	}

	require.Equal(2, ws.Len())
}

func (suite *WithdrawTestSuite) TestNewWithdrawalFromMsg() {
	require := suite.Require()
	withdrawal := types.NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1)
	msgAddwithdrawal := types.MsgAddWithdrawal{
		Authority:  "gov",
		Withdrawal: &withdrawal,
	}
	w := types.NewWithdrawalFromMsg(&msgAddwithdrawal)
	require.Equal(withdrawal, w, "NewWithdrawalFromMsg should return the same withdrawal")
}

func (suite *WithdrawTestSuite) TestMustMarshalWithdraw() {
	require := suite.Require()
	withdrawal := types.NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1)
	require.NotPanics(func() {
		marshaled := types.MustMarshalWithdrawal(suite.encConf.Codec, &withdrawal)
		require.NotNil(marshaled, "MarshalWithdrawal should not return nil")
	})
}

func (suite *WithdrawTestSuite) TestUnmarshalWithdraw() {
	require := suite.Require()
	withdrawal := types.NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1)

	testCases := []struct {
		name           string
		input          []byte
		expectedResult types.Withdrawal
		expectError    bool
	}{
		{
			name:           "Unmarshal valid withdrawal bytes",
			input:          types.MustMarshalWithdrawal(suite.encConf.Codec, &withdrawal),
			expectedResult: types.NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1),
			expectError:    false,
		},
		{
			name:           "Unmarshal invalid withdrawal bytes",
			input:          []byte{1},
			expectedResult: types.Withdrawal{}, // Expecting an empty struct since it will fail
			expectError:    true,
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			result, err := types.UnmarshalWithdrawal(suite.encConf.Codec, tc.input)
			if tc.expectError {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tc.expectedResult, result, "UnmarshalWithdrawal should return the correct withdrawal")
			}
		})
	}
}

func (suite *WithdrawTestSuite) TestMustUnmarshalWithdraw() {
	require := suite.Require()
	withdrawal := types.NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1)

	testCases := []struct {
		name        string
		input       []byte
		expected    types.Withdrawal
		expectPanic bool
	}{
		{
			name:        "Unmarshal valid withdrawal bytes",
			input:       types.MustMarshalWithdrawal(suite.encConf.Codec, &withdrawal),
			expected:    types.NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1),
			expectPanic: false,
		},
		{
			name:        "Unmarshal invalid withdrawal bytes - panic",
			input:       []byte{1},
			expectPanic: true,
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.expectPanic {
				require.Panics(func() {
					types.MustUnmarshalWithdrawal(suite.encConf.Codec, tc.input)
				})
			} else {
				require.NotPanics(func() {
					result := types.MustUnmarshalWithdrawal(suite.encConf.Codec, tc.input)
					require.Equal(tc.expected, result, "MustUnmarshalWithdrawal should return the correct withdrawal")
				})
			}
		})
	}
}

func TestTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WithdrawTestSuite))
}
