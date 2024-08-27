package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/contracts/bindings"
)

func Test_mustGetABI_NoPanic(t *testing.T) {
	require.NotPanics(t, func() {
		mustGetABI(bindings.IPTokenStakingMetaData)
	})
}

func Test_mustGetABI_PanicsOnInvalidMetadata(t *testing.T) {
	require.Panics(t, func() {
		mustGetABI(&bind.MetaData{ABI: "invalid ABI"})
	})
}

func Test_mustGetEvent_NoPanic(t *testing.T) {
	abi := mustGetABI(bindings.IPTokenStakingMetaData)
	require.NotPanics(t, func() {
		mustGetEvent(abi, "SetWithdrawalAddress")
		mustGetEvent(abi, "CreateValidator")
		mustGetEvent(abi, "Deposit")
		mustGetEvent(abi, "Redelegate")
		mustGetEvent(abi, "Withdraw")
	})
}

func Test_mustGetEvent_PanicsOnInvalidEvent(t *testing.T) {
	abi := mustGetABI(bindings.IPTokenStakingMetaData)
	require.Panics(t, func() {
		mustGetEvent(abi, "NonExistentEvent")
	})
}
