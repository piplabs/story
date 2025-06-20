package types_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmengine/types"
)

var (
	uintOne = uint64(1)
)

func TestPayloadToProto_AllErrorCases(t *testing.T) {
	// Case: nil payload
	_, err := types.PayloadToProto(nil)
	require.ErrorContains(t, err, "nil payload")

	// Case: nil BlobGasUsed
	payload := &engine.ExecutableData{
		ExcessBlobGas: &uintOne,
		BaseFeePerGas: new(big.Int).SetUint64(1),
	}
	_, err = types.PayloadToProto(payload)
	require.ErrorContains(t, err, "nil payload BlobGasUsed")

	// Case: nil ExcessBlobGas
	payload.BlobGasUsed = &uintOne
	payload.ExcessBlobGas = nil
	_, err = types.PayloadToProto(payload)
	require.ErrorContains(t, err, "nil payload ExcessBlobGas")

	// Case: has ExecutionWitness
	payload.ExcessBlobGas = &uintOne
	payload.ExecutionWitness = &etypes.ExecutionWitness{}
	_, err = types.PayloadToProto(payload)
	require.ErrorContains(t, err, "payload has ExecutionWitness")

	// Case: withdrawal conversion error
	payload.ExecutionWitness = nil
	payload.Withdrawals = []*etypes.Withdrawal{nil}
	_, err = types.PayloadToProto(payload)
	require.ErrorContains(t, err, "withdrawal to proto")
}

func TestPayloadToProto_Valid(t *testing.T) {
	baseFee := big.NewInt(1000000000)
	withdrawal := &etypes.Withdrawal{
		Index:     1,
		Validator: 2,
		Address:   common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Amount:    1000,
	}

	p := &engine.ExecutableData{
		ParentHash:    common.HexToHash("0xabc"),
		FeeRecipient:  common.HexToAddress("0x456"),
		StateRoot:     common.HexToHash("0xdef"),
		ReceiptsRoot:  common.HexToHash("0xaaa"),
		LogsBloom:     []byte{},
		Random:        common.HexToHash("0xbbb"),
		Number:        10,
		GasLimit:      1000000,
		GasUsed:       900000,
		Timestamp:     123456789,
		ExtraData:     []byte("extra"),
		BaseFeePerGas: baseFee,
		BlockHash:     common.HexToHash("0xccc"),
		Transactions:  [][]byte{[]byte("tx1"), []byte("tx2")},
		Withdrawals:   []*etypes.Withdrawal{withdrawal},
		BlobGasUsed:   &uintOne,
		ExcessBlobGas: &uintOne,
	}

	proto, err := types.PayloadToProto(p)
	require.NoError(t, err)
	require.Equal(t, p.Number, proto.BlockNumber)
	require.Equal(t, p.GasUsed, proto.GasUsed)
	require.Equal(t, p.Transactions, proto.Transactions)
	require.Equal(t, *p.BlobGasUsed, proto.BlobGasUsed)
}

func TestPayloadFromProto_AllErrorCases(t *testing.T) {
	// Case: nil input
	_, err := types.PayloadFromProto(nil)
	require.ErrorContains(t, err, "nil payload")

	// Case: invalid BaseFeePerGas length
	badProto := &types.ExecutionPayloadDeneb{
		BaseFeePerGas: []byte{1, 2, 3}, // wrong length
	}
	_, err = types.PayloadFromProto(badProto)
	require.ErrorContains(t, err, "invalid BaseFeePerGas length")
}

func TestPayloadFromProto_Valid(t *testing.T) {
	baseFee := big.NewInt(1000000000)
	bz := baseFee.FillBytes(make([]byte, 32))

	proto := &types.ExecutionPayloadDeneb{
		ParentHash:    types.Hash(common.HexToHash("0xabc")),
		FeeRecipient:  types.Address(common.HexToAddress("0x456")),
		StateRoot:     types.Hash(common.HexToHash("0xdef")),
		ReceiptsRoot:  types.Hash(common.HexToHash("0xaaa")),
		LogsBloom:     []byte{},
		PrevRandao:    types.Hash(common.HexToHash("0xbbb")),
		BlockNumber:   10,
		GasLimit:      1000000,
		GasUsed:       900000,
		Timestamp:     123456789,
		ExtraData:     []byte("extra"),
		BaseFeePerGas: bz,
		BlockHash:     types.Hash(common.HexToHash("0xccc")),
		Transactions:  nil,
		Withdrawals: []types.WithdrawalEVM{
			{
				Index:          1,
				ValidatorIndex: 2,
				Address:        types.Address(common.HexToAddress("0x123")),
				Amount:         1000,
			},
		},
		BlobGasUsed:   1,
		ExcessBlobGas: 2,
	}

	payload, err := types.PayloadFromProto(proto)
	require.NoError(t, err)
	require.Equal(t, proto.BlockNumber, payload.Number)
	require.Equal(t, uint64(0), payload.GasUsed-900000) // Check correctness
	require.NotNil(t, payload.Transactions)
	require.NotNil(t, payload.Withdrawals)
}

func TestWithdrawalToProto(t *testing.T) {
	// Case: nil input
	_, err := types.WithdrawalToProto(nil)
	require.ErrorContains(t, err, "nil withdrawal")

	// Case: valid input
	with := &etypes.Withdrawal{
		Index:     1,
		Validator: 2,
		Address:   common.HexToAddress("0xabc"),
		Amount:    1000,
	}
	proto, err := types.WithdrawalToProto(with)
	require.NoError(t, err)
	require.Equal(t, with.Index, proto.Index)
}

func TestWithdrawalFromProto(t *testing.T) {
	proto := types.WithdrawalEVM{
		Index:          1,
		ValidatorIndex: 2,
		Address:        types.Address(common.HexToAddress("0xabc")),
		Amount:         1000,
	}
	with := types.WithdrawalFromProto(proto)
	require.Equal(t, proto.Index, with.Index)
	require.Equal(t, proto.Amount, with.Amount)
}
