package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/piplabs/story/lib/errors"
)

// PackInitCode packs the init code for a contract deployment.
func PackInitCode(abi *abi.ABI, bytecodeHex string, params ...any) ([]byte, error) {
	input, err := abi.Pack("", params...)
	if err != nil {
		return nil, errors.Wrap(err, "pack init code")
	}

	bytecode, err := hexutil.Decode(bytecodeHex)
	if err != nil {
		return nil, errors.Wrap(err, "decode bytecode")
	}

	return append(bytecode, input...), nil
}
