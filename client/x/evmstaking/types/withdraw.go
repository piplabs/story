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

// TODO amount as math.Int.
func NewWithdrawal(creationHeight uint64, delegatorAddr string, validatorAddr string, executionAddr string, amount uint64) Withdrawal {
	return Withdrawal{
		CreationHeight:   creationHeight,
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		ExecutionAddress: executionAddr,
		Amount:           amount,
	}
}

func NewWithdrawalFromMsg(msg *MsgAddWithdrawal) Withdrawal {
	return Withdrawal{
		CreationHeight:   msg.Withdrawal.CreationHeight,
		DelegatorAddress: msg.Withdrawal.DelegatorAddress,
		ValidatorAddress: msg.Withdrawal.ValidatorAddress,
		ExecutionAddress: msg.Withdrawal.ExecutionAddress,
		Amount:           msg.Withdrawal.Amount,
	}
}
