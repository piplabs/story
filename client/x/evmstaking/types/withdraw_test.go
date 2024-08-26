package types

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
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
	suite.T().Parallel()

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
	suite.T().Parallel()
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
	suite.T().Parallel()
	ws := Withdrawals{
		Withdrawals: []Withdrawal{
			NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1),
			NewWithdrawal(2, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 2),
		},
	}

	suite.Equal(2, ws.Len())
}

func (suite *WithdrawTestSuite) TestNewWithdrawalFromMsg() {
	suite.T().Parallel()
	withdrawal := NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1)
	msgAddwithdrawal := MsgAddWithdrawal{
		Authority:  "gov",
		Withdrawal: &withdrawal,
	}
	w := NewWithdrawalFromMsg(&msgAddwithdrawal)
	suite.Equal(withdrawal, w, "NewWithdrawalFromMsg should return the same withdrawal")
}

func (suite *WithdrawTestSuite) TestMustMarshalWithdraw() {
	suite.T().Parallel()
	withdrawal := NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1)
	suite.NotPanics(func() {
		marshaled := MustMarshalWithdrawal(suite.encConf.Codec, &withdrawal)
		suite.NotNil(marshaled, "MarshalWithdrawal should not return nil")
	})
}

func (suite *WithdrawTestSuite) TestUnmarshalWithdraw() {
	suite.T().Parallel()
	withdrawal := NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1)
	marshaled := MustMarshalWithdrawal(suite.encConf.Codec, &withdrawal)

	w, err := UnmarshalWithdrawal(suite.encConf.Codec, marshaled)
	suite.NoError(err)
	suite.Equal(withdrawal, w, "UnmarshalWithdrawal should return the same withdrawal")
}

func (suite *WithdrawTestSuite) TestMustUnmarshalWithdraw() {
	suite.T().Parallel()
	withdrawal := NewWithdrawal(1, suite.delAddr, suite.valAddr, suite.evmAddr.String(), 1)
	suite.NotPanics(func() {
		w := MustUnmarshalWithdrawal(
			suite.encConf.Codec,
			MustMarshalWithdrawal(suite.encConf.Codec, &withdrawal),
		)
		require.Equal(suite.T(), withdrawal, w)
	})
}

func TestTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WithdrawTestSuite))
}
