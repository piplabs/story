package types_test

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/piplabs/story/client/x/evmstaking/types"
)

type WithdrawTestSuite struct {
	suite.Suite

	encConf    testutil.TestEncodingConfig
	delAddr    string
	valAddr    string
	valEVMAddr string
	evmAddr    common.Address
}

func (suite *WithdrawTestSuite) SetupTest() {
	suite.encConf = testutil.MakeTestEncodingConfig()
	// set dummy addresses
	suite.delAddr = "story1hmjw3pvkjtndpg8wqppwdn8udd835qpan4hm0y"
	suite.valAddr = "storyvaloper1hmjw3pvkjtndpg8wqppwdn8udd835qpaa6r6y0"
	suite.valEVMAddr = "0xbee4e8859692e6d0a0ee0042e6ccfc6b4f1a003d"
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
			withdrawal:     types.NewWithdrawal(1, suite.evmAddr.String(), 100, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, suite.valEVMAddr),
			expectedString: fmt.Sprintf("creation_height:1 execution_address:\"%s\" amount:100 validator_address:\"%s\" ", suite.evmAddr.String(), suite.valEVMAddr),
		},
		{
			name:           "Empty addresses",
			withdrawal:     types.NewWithdrawal(1, "", 1, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, suite.valEVMAddr),
			expectedString: fmt.Sprintf("creation_height:1 amount:1 validator_address:\"%s\" ", suite.valEVMAddr),
		},
		{
			name:           "Large amount",
			withdrawal:     types.NewWithdrawal(1, suite.evmAddr.String(), math.MaxUint64, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, suite.valEVMAddr),
			expectedString: fmt.Sprintf("creation_height:1 execution_address:\"%s\" amount:%s validator_address:\"%s\" ", suite.evmAddr.String(), new(big.Int).SetUint64(math.MaxUint64).String(), suite.valEVMAddr),
		},
		{
			name:           "Zero amount",
			withdrawal:     types.NewWithdrawal(1, suite.evmAddr.String(), 0, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, suite.valEVMAddr),
			expectedString: fmt.Sprintf("creation_height:1 execution_address:\"%s\" validator_address:\"%s\" ", suite.evmAddr.String(), suite.valEVMAddr),
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
			types.NewWithdrawal(1, suite.evmAddr.String(), 1, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, suite.valEVMAddr),
			types.NewWithdrawal(2, suite.evmAddr.String(), 2, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, suite.valEVMAddr),
		},
	}

	expectedString := fmt.Sprintf(
		"creation_height:1 execution_address:\"%s\" amount:1 validator_address:\"%s\" \n"+
			"creation_height:2 execution_address:\"%s\" amount:2 validator_address:\"%s\"",
		suite.evmAddr.String(), suite.valEVMAddr,
		suite.evmAddr.String(), suite.valEVMAddr,
	)
	require.Equal(expectedString, ws.String())
}

func (suite *WithdrawTestSuite) TestWithdrawalsLen() {
	require := suite.Require()
	ws := types.Withdrawals{
		Withdrawals: []types.Withdrawal{
			types.NewWithdrawal(1, suite.evmAddr.String(), 1, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, suite.valEVMAddr),
			types.NewWithdrawal(2, suite.evmAddr.String(), 2, types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE, suite.valEVMAddr),
		},
	}

	require.Equal(2, ws.Len())
}

func TestWithdrawalTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WithdrawTestSuite))
}
