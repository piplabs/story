package types

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/lib/errors"
)

// NewParams creates a new Params instance.
func NewParams(executionBlockHash []byte) Params {
	return Params{
		ExecutionBlockHash: executionBlockHash,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		nil,
	)
}

func ValidateExecutionBlockHash(executionBlockHash []byte) error {
	if len(executionBlockHash) != common.HashLength {
		return errors.New("invalid execution block hash length", "length", len(executionBlockHash))
	}

	return nil
}
