package types

import (
	"bytes"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/piplabs/story/lib/errors"
)

// ExecutableData originates from go-ethereum/beacon/engine/gen_ed.go.
// It is designed to disallow unknown fields in executable data.
// However, since engine.ExecutableData implements the MarshalJSON and UnmarshalJSON methods,
// json.Decode.DisallowUnknownFields() cannot be used (see encoding/json/decode.go#L479).
type ExecutableData struct {
	ParentHash    *common.Hash    `json:"parentHash"`
	FeeRecipient  *common.Address `json:"feeRecipient"`
	StateRoot     *common.Hash    `json:"stateRoot"`
	ReceiptsRoot  *common.Hash    `json:"receiptsRoot"`
	LogsBloom     *hexutil.Bytes  `json:"logsBloom"`
	Random        *common.Hash    `json:"prevRandao"`
	Number        *hexutil.Uint64 `json:"blockNumber"`
	GasLimit      *hexutil.Uint64 `json:"gasLimit"`
	GasUsed       *hexutil.Uint64 `json:"gasUsed"`
	Timestamp     *hexutil.Uint64 `json:"timestamp"`
	ExtraData     *hexutil.Bytes  `json:"extraData"`
	BaseFeePerGas *hexutil.Big    `json:"baseFeePerGas"`
	BlockHash     *common.Hash    `json:"blockHash"`
	Transactions  []hexutil.Bytes `json:"transactions"`
	Withdrawals   []*Withdrawal   `json:"withdrawals"`
	BlobGasUsed   *hexutil.Uint64 `json:"blobGasUsed"`
	ExcessBlobGas *hexutil.Uint64 `json:"excessBlobGas"`
}

type Withdrawal struct {
	Index     *hexutil.Uint64 `json:"index"`
	Validator *hexutil.Uint64 `json:"validatorIndex"`
	Address   common.Address  `json:"address"`
	Amount    *hexutil.Uint64 `json:"amount"`
}

// ValidateExecPayload checks the execution payload for any fields that do not match the expected values.
func ValidateExecPayload(msg *MsgExecutionPayload) error {
	decoder := json.NewDecoder(bytes.NewReader(msg.ExecutionPayload))
	decoder.DisallowUnknownFields()

	var payloadTmp ExecutableData
	if err := decoder.Decode(&payloadTmp); err != nil {
		return errors.Wrap(err, "unmarshal payload")
	}

	return nil
}
