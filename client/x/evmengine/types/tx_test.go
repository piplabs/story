package types_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/crypto"
	k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/genutil/evm/predeploys"
	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/contracts/bindings"
)

const (
	dummyAddressHex = "0x1398C32A45Bc409b6C652E25bb0a3e702492A4ab"
)

var (
	dummyContractAddress = common.HexToAddress(dummyAddressHex)
	emptyAddr            = common.Address{}
	emptyData            = []byte{}
)

// initializeABI loads the ABI once.
func initializeABI(t *testing.T) *abi.ABI {
	t.Helper()
	upgradeAbi, err := bindings.UpgradeEntrypointMetaData.GetAbi()
	require.NoError(t, err, "failed to load ABI")

	return upgradeAbi
}

func TestEVMEvent_ToEthLog(t *testing.T) {
	t.Parallel()
	upgradeAbi := initializeABI(t)
	data, err := upgradeAbi.Events["SoftwareUpgrade"].Inputs.NonIndexed().Pack("test-upgrade", int64(1), "test-info")
	require.NoError(t, err)

	tcs := []struct {
		name           string
		evmEvent       *types.EVMEvent
		expectedResult ethtypes.Log
	}{
		{
			name: "zero address & empty data",
			evmEvent: &types.EVMEvent{
				Address: emptyAddr.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				Data:    emptyData,
			},
			expectedResult: ethtypes.Log{
				Address: emptyAddr,
				Topics:  []common.Hash{types.SoftwareUpgradeEvent.ID},
				Data:    emptyData,
			},
		},
		{
			name: "empty data",
			evmEvent: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				Data:    emptyData,
			},
			expectedResult: ethtypes.Log{
				Address: dummyContractAddress,
				Topics:  []common.Hash{types.SoftwareUpgradeEvent.ID},
				Data:    emptyData,
			},
		},
		{
			name: "full log",
			evmEvent: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				Data:    data,
			},
			expectedResult: ethtypes.Log{
				Address: dummyContractAddress,
				Topics:  []common.Hash{types.SoftwareUpgradeEvent.ID},
				Data:    data,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := tc.evmEvent.ToEthLog()
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestEVMEvent_Verify(t *testing.T) {
	t.Parallel()
	upgradeAbi := initializeABI(t)
	data, err := upgradeAbi.Events["SoftwareUpgrade"].Inputs.NonIndexed().Pack("test-upgrade", int64(1), "test-info")
	require.NoError(t, err)

	tcs := []struct {
		name        string
		evmEvent    *types.EVMEvent
		expectedErr string
	}{
		{
			name:        "fail: nil",
			evmEvent:    nil,
			expectedErr: "nil log",
		},
		{
			name: "fail: nil address",
			evmEvent: &types.EVMEvent{
				Address: nil,
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
			},
			expectedErr: "nil address",
		},
		{
			name: "fail: empty topics",
			evmEvent: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{},
			},
			expectedErr: "empty topics",
		},
		{
			name: "fail: invalid address length",
			evmEvent: &types.EVMEvent{
				Address: []byte{0x01},
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
			},
			expectedErr: "invalid address length",
		},
		{
			name: "fail: invalid topic length",
			evmEvent: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{{0x01}},
			},
			expectedErr: "invalid topic length",
		},
		{
			name: "pass: valid log",
			evmEvent: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				Data:    data,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.evmEvent.Verify()
			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEthLogToEVMEvent(t *testing.T) {
	t.Parallel()
	upgradeAbi := initializeABI(t)
	data, err := upgradeAbi.Events["SoftwareUpgrade"].Inputs.NonIndexed().Pack("test-upgrade", int64(1), "test-info")
	require.NoError(t, err)

	tcs := []struct {
		name           string
		log            ethtypes.Log
		expectedResult *types.EVMEvent
		expectedErr    string
	}{
		{
			name: "pass: zero address",
			log: ethtypes.Log{
				Address: emptyAddr,
				Topics:  []common.Hash{types.SoftwareUpgradeEvent.ID},
				Data:    data,
			},
			expectedResult: &types.EVMEvent{
				Address: emptyAddr.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				Data:    data,
			},
		},
		{
			name: "pass: empty data",
			log: ethtypes.Log{
				Address: dummyContractAddress,
				Topics:  []common.Hash{types.SoftwareUpgradeEvent.ID},
				Data:    emptyData,
			},
			expectedResult: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				Data:    emptyData,
			},
		},
		{
			name: "pass: zero address & empty data",
			log: ethtypes.Log{
				Address: emptyAddr,
				Topics:  []common.Hash{types.SoftwareUpgradeEvent.ID},
				Data:    emptyData,
			},
			expectedResult: &types.EVMEvent{
				Address: emptyAddr.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				Data:    emptyData,
			},
		},
		{
			name: "pass: full log",
			log: ethtypes.Log{
				Address: dummyContractAddress,
				Topics:  []common.Hash{types.SoftwareUpgradeEvent.ID},
				Data:    data,
			},
			expectedResult: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				Data:    data,
			},
		},
		{
			name: "fail: empty topics",
			log: ethtypes.Log{
				Address: dummyContractAddress,
				Topics:  []common.Hash{},
				Data:    data,
			},
			expectedErr: "verify log: empty topics",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := types.EthLogToEVMEvent(tc.log)
			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestSortEVMEvents(t *testing.T) {
	t.Parallel()
	stakingAbi, err := bindings.IPTokenStakingMetaData.GetAbi()
	require.NoError(t, err, "failed to load ABI")
	stakingAddr := common.HexToAddress(predeploys.IPTokenStaking)
	slashingAbi, err := bindings.IPTokenSlashingMetaData.GetAbi()
	require.NoError(t, err, "failed to load ABI")
	slashingAddr := common.HexToAddress(predeploys.IPTokenSlashing)
	require.Negative(t, bytes.Compare(stakingAddr.Bytes(), slashingAddr.Bytes()), "stakingAddr should be less than slashingAddr")

	withdrawEv := stakingAbi.Events["Withdraw"].ID.Bytes()
	depositEv := stakingAbi.Events["Deposit"].ID.Bytes()
	require.Negative(t, bytes.Compare(withdrawEv, depositEv), "withdrawEv should be less than depositEv")
	unjailEv := slashingAbi.Events["Unjail"].ID.Bytes()

	// prepare data
	pubKeys, _, _ := createAddresses(2)
	delPubKey := pubKeys[0]
	valPubKey := pubKeys[1]
	delSecp256k1PubKey, err := secp256k1.ParsePubKey(delPubKey.Bytes())
	require.NoError(t, err)
	uncompressedDelPubKeyBytes := delSecp256k1PubKey.SerializeUncompressed()
	gwei, exp := big.NewInt(10), big.NewInt(9)
	gwei.Exp(gwei, exp, nil)
	delAmtGwei := new(big.Int).Mul(gwei, new(big.Int).SetUint64(100))

	// staking contract events
	withdrawData, err := stakingAbi.Events["Withdraw"].Inputs.NonIndexed().Pack(
		delPubKey.Bytes(), valPubKey.Bytes(), delAmtGwei,
	)
	require.NoError(t, err)
	depositData, err := stakingAbi.Events["Deposit"].Inputs.NonIndexed().Pack(
		uncompressedDelPubKeyBytes, delPubKey.Bytes(), valPubKey.Bytes(), delAmtGwei,
	)
	require.NoError(t, err)
	require.Negative(t, bytes.Compare(withdrawData, depositData), "withdrawData should be less than depositData")
	cpyDelPubKey := delPubKey.Bytes()
	cpyDelPubKey[0] += 1 // add 1 to the first byte so it should be greater than delPubKey
	depositData2, err := stakingAbi.Events["Deposit"].Inputs.NonIndexed().Pack(
		uncompressedDelPubKeyBytes, delPubKey.Bytes(), valPubKey.Bytes(), delAmtGwei,
	)
	require.NoError(t, err)
	require.Negative(t, bytes.Compare(depositData, depositData2), "depositData should be less than depositData2")
	// slashing contract events
	unjailData, err := slashingAbi.Events["Unjail"].Inputs.NonIndexed().Pack(valPubKey.Bytes())
	require.NoError(t, err)

	tcs := []struct {
		name           string
		evmEvents      []*types.EVMEvent
		expectedResult []*types.EVMEvent
	}{
		{
			name: "pass: single contract and single event",
			evmEvents: []*types.EVMEvent{
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{withdrawEv},
					Data:    withdrawData,
				},
			},
			expectedResult: []*types.EVMEvent{
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{withdrawEv},
					Data:    withdrawData,
				},
			},
		},
		{
			name: "pass: single contract and multiple event - sorted by topics",
			evmEvents: []*types.EVMEvent{
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{depositEv},
					Data:    depositData,
				},
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{withdrawEv},
					Data:    withdrawData,
				},
			},
			expectedResult: []*types.EVMEvent{
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{withdrawEv},
					Data:    withdrawData,
				},
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{depositEv},
					Data:    depositData,
				},
			},
		},
		{
			name: "pass: single contract and multiple event - sorted by data",
			evmEvents: []*types.EVMEvent{
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{depositEv},
					Data:    depositData2,
				},
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{depositEv},
					Data:    depositData,
				},
			},
			expectedResult: []*types.EVMEvent{
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{depositEv},
					Data:    depositData,
				},
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{depositEv},
					Data:    depositData2,
				},
			},
		},
		{
			name: "pass: multiple contract - should be sorted by address",
			evmEvents: []*types.EVMEvent{
				{
					Address: slashingAddr.Bytes(),
					Topics:  [][]byte{unjailEv},
					Data:    unjailData,
				},
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{depositEv},
					Data:    depositData,
				},
			},
			expectedResult: []*types.EVMEvent{
				{
					Address: stakingAddr.Bytes(),
					Topics:  [][]byte{depositEv},
					Data:    depositData,
				},
				{
					Address: slashingAddr.Bytes(),
					Topics:  [][]byte{unjailEv},
					Data:    unjailData,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			types.SortEVMEvents(tc.evmEvents)
			require.Equal(t, tc.expectedResult, tc.evmEvents)
		})
	}
}

func createAddresses(count int) ([]crypto.PubKey, []sdk.AccAddress, []sdk.ValAddress) {
	var pubKeys []crypto.PubKey
	var accAddrs []sdk.AccAddress
	var valAddrs []sdk.ValAddress
	for range count {
		pubKey := k1.GenPrivKey().PubKey()
		accAddr := sdk.AccAddress(pubKey.Address().Bytes())
		valAddr := sdk.ValAddress(pubKey.Address().Bytes())
		pubKeys = append(pubKeys, pubKey)
		accAddrs = append(accAddrs, accAddr)
		valAddrs = append(valAddrs, valAddr)
	}

	return pubKeys, accAddrs, valAddrs
}
