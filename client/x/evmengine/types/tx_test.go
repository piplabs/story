package types_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/contracts/bindings"
)

const (
	dummyAddressHex = "0x1398C32A45Bc409b6C652E25bb0a3e702492A4ab"
)

var (
	dummyContractAddress = common.HexToAddress(dummyAddressHex)
	dummyTxHash          = common.HexToHash(dummyAddressHex)
	emptyAddr            = common.Address{}
	emptyData            = []byte{}
	emptyTxHash          = common.Hash{}
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
				TxHash:  emptyTxHash.Bytes(),
			},
			expectedResult: ethtypes.Log{
				Address: emptyAddr,
				Topics:  []common.Hash{types.SoftwareUpgradeEvent.ID},
				Data:    emptyData,
				TxHash:  emptyTxHash,
			},
		},
		{
			name: "empty data",
			evmEvent: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				Data:    emptyData,
				TxHash:  emptyTxHash.Bytes(),
			},
			expectedResult: ethtypes.Log{
				Address: dummyContractAddress,
				Topics:  []common.Hash{types.SoftwareUpgradeEvent.ID},
				Data:    emptyData,
				TxHash:  emptyTxHash,
			},
		},
		{
			name: "full log",
			evmEvent: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				Data:    data,
				TxHash:  emptyTxHash.Bytes(),
			},
			expectedResult: ethtypes.Log{
				Address: dummyContractAddress,
				Topics:  []common.Hash{types.SoftwareUpgradeEvent.ID},
				Data:    data,
				TxHash:  emptyTxHash,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := tc.evmEvent.ToEthLog()
			require.NoError(t, err)
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
				TxHash:  dummyTxHash.Bytes(),
			},
			expectedErr: "nil address",
		},
		{
			name: "fail: empty topics",
			evmEvent: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{},
				TxHash:  dummyTxHash.Bytes(),
			},
			expectedErr: "empty topics",
		},
		{
			name: "fail: nil tx hash",
			evmEvent: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				TxHash:  nil,
			},
			expectedErr: "nil tx hash",
		},
		{
			name: "fail: invalid address length",
			evmEvent: &types.EVMEvent{
				Address: []byte{0x01},
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				TxHash:  dummyTxHash.Bytes(),
			},
			expectedErr: "invalid address length",
		},
		{
			name: "fail: invalid topic length",
			evmEvent: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{{0x01}},
				TxHash:  dummyTxHash.Bytes(),
			},
			expectedErr: "invalid topic length",
		},
		{
			name: "fail: invalid tx hash length",
			evmEvent: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				TxHash:  []byte{0x01},
			},
			expectedErr: "invalid tx hash length",
		},
		{
			name: "pass: valid log",
			evmEvent: &types.EVMEvent{
				Address: dummyContractAddress.Bytes(),
				Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
				Data:    data,
				TxHash:  dummyTxHash.Bytes(),
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
