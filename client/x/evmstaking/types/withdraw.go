package types

import (
	"strings"

	"cosmossdk.io/core/address"
)

// Withdrawals is a collection of Withdrawal.
type Withdrawals struct {
	Withdrawals     []Withdrawal
	WithdrawalCodec address.Codec
}

func (ws Withdrawals) String() (out string) {
	for _, w := range ws.Withdrawals {
		out += w.String() + "\n"
	}

	return strings.TrimSpace(out)
}

func (ws Withdrawals) Len() int {
	return len(ws.Withdrawals)
}

func NewWithdrawal(creationHeight uint64, executionAddr string, amount uint64, withdrawalType WithdrawalType, valEVMAddr string) Withdrawal {
	return Withdrawal{
		CreationHeight:   creationHeight,
		ExecutionAddress: executionAddr,
		Amount:           amount,
		WithdrawalType:   withdrawalType,
		ValidatorAddress: valEVMAddr,
	}
}
