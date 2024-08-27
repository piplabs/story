package types

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
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
	// Define the test cases
	testCases := []struct {
		name           string
		withdrawal     Withdrawal
		expectedString string
	}{
		{
			name:           "Normal values",
			withdrawal:     NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 100),
			expectedString: fmt.Sprintf("creation_height:1 delegator_address:\"%s\" validator_address:\"%s\" execution_address:\"%s\" amount:100 ", suite.delAddr, suite.valAddr, suite.evmAddr.String()),
		},
		{
			name:           "Empty addresses",
			withdrawal:     NewWithdrawal(1, "", "", "", 1),
			expectedString: "creation_height:1 amount:1 ",
		},
		{
			name:           "Large amount",
			withdrawal:     NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), math.MaxUint64),
			expectedString: fmt.Sprintf("creation_height:1 delegator_address:\"%s\" validator_address:\"%s\" execution_address:\"%s\" amount:%s ", suite.delAddr, suite.valAddr, suite.evmAddr.String(), new(big.Int).SetUint64(math.MaxUint64).String()),
		},
		{
			name:           "Zero amount",
			withdrawal:     NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 0),
			expectedString: fmt.Sprintf("creation_height:1 delegator_address:\"%s\" validator_address:\"%s\" execution_address:\"%s\" ", suite.delAddr, suite.valAddr, suite.evmAddr.String()),
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			result := tc.withdrawal.String()
			suite.Equal(tc.expectedString, result)
		})
	}
}

func (suite *WithdrawTestSuite) TestWithdrawalsString() {
	ws := Withdrawals{
		Withdrawals: []Withdrawal{
			NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1),
			NewWithdrawal(2, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 2),
		},
	}

	expectedString := fmt.Sprintf(
		`creation_height:1 delegator_address:"%s" validator_address:"%s" execution_address:"%s" amount:1 
creation_height:2 delegator_address:"%s" validator_address:"%s" execution_address:"%s" amount:2`,
		suite.delAddr, suite.valAddr, suite.evmAddr.String(),
		suite.delAddr, suite.valAddr, suite.evmAddr.String(),
	)
	suite.Equal(expectedString, ws.String())
}

func (suite *WithdrawTestSuite) TestWithdrawalsLen() {
	ws := Withdrawals{
		Withdrawals: []Withdrawal{
			NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1),
			NewWithdrawal(2, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 2),
		},
	}

	suite.Equal(2, ws.Len())
}

func (suite *WithdrawTestSuite) TestNewWithdrawalFromMsg() {
	withdrawal := NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1)
	msgAddwithdrawal := MsgAddWithdrawal{
		Authority:  "gov",
		Withdrawal: &withdrawal,
	}
	w := NewWithdrawalFromMsg(&msgAddwithdrawal)
	suite.Equal(withdrawal, w, "NewWithdrawalFromMsg should return the same withdrawal")
}

func (suite *WithdrawTestSuite) TestMustMarshalWithdraw() {
	withdrawal := NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1)
	suite.NotPanics(func() {
		marshaled := MustMarshalWithdrawal(suite.encConf.Codec, &withdrawal)
		suite.NotNil(marshaled, "MarshalWithdrawal should not return nil")
	})
}

func (suite *WithdrawTestSuite) TestUnmarshalWithdraw() {
	withdrawal := NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1)

	testCases := []struct {
		name           string
		input          []byte
		expectedResult Withdrawal
		expectError    bool
	}{
		{
			name:           "Unmarshal valid withdrawal bytes",
			input:          MustMarshalWithdrawal(suite.encConf.Codec, &withdrawal),
			expectedResult: NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1),
			expectError:    false,
		},
		{
			name:           "Unmarshal invalid withdrawal bytes",
			input:          []byte{1},
			expectedResult: Withdrawal{}, // Expecting an empty struct since it will fail
			expectError:    true,
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			result, err := UnmarshalWithdrawal(suite.encConf.Codec, tc.input)
			if tc.expectError {
				suite.Error(err)
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedResult, result, "UnmarshalWithdrawal should return the correct withdrawal")
			}
		})
	}
}

func (suite *WithdrawTestSuite) TestMustUnmarshalWithdraw() {
	withdrawal := NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1)

	testCases := []struct {
		name        string
		input       []byte
		expected    Withdrawal
		expectPanic bool
	}{
		{
			name:        "Unmarshal valid withdrawal bytes",
			input:       MustMarshalWithdrawal(suite.encConf.Codec, &withdrawal),
			expected:    NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1),
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
				suite.Panics(func() {
					MustUnmarshalWithdrawal(suite.encConf.Codec, tc.input)
				})
			} else {
				suite.NotPanics(func() {
					result := MustUnmarshalWithdrawal(suite.encConf.Codec, tc.input)
					suite.Equal(tc.expected, result, "MustUnmarshalWithdrawal should return the correct withdrawal")
				})
			}
		})
	}
}

func TestTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WithdrawTestSuite))
}
