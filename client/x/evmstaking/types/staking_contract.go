package types

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/storyprotocol/iliad/contracts/bindings"
)

var (
	ipTokenStakingABI    = mustGetABI(bindings.IPTokenStakingMetaData)
	SetWithdrawalAddress = mustGetEvent(ipTokenStakingABI, "SetWithdrawalAddress")
	CreateValidatorEvent = mustGetEvent(ipTokenStakingABI, "CreateValidator")
	DepositEvent         = mustGetEvent(ipTokenStakingABI, "Deposit")
	RedelegateEvent      = mustGetEvent(ipTokenStakingABI, "Redelegate")
	WithdrawEvent        = mustGetEvent(ipTokenStakingABI, "Withdraw")
)

// mustGetABI returns the metadata's ABI as an abi.ABI type.
// It panics on error.
func mustGetABI(metadata *bind.MetaData) *abi.ABI {
	abi, err := metadata.GetAbi()
	if err != nil {
		panic(err)
	}

	return abi
}

// mustGetEvent returns the event with the given name from the ABI.
// It panics if the event is not found.
func mustGetEvent(abi *abi.ABI, name string) abi.Event {
	event, ok := abi.Events[name]
	if !ok {
		panic("event not found")
	}

	return event
}
