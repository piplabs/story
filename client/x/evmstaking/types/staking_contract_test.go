package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/contracts/bindings"
)

func Test_mustGetABI_Nopanic(t *testing.T) {
	require.NotPanics(t, func() {
		mustGetABI(bindings.IPTokenStakingMetaData)
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

	require.Panics(t, func() {
		mustGetEvent(abi, "NonExistentEvent")
	})
}

func Test_mustGetEvent_PanicsOnInvalidEvent(t *testing.T) {
	abi := mustGetABI(bindings.IPTokenStakingMetaData)
	require.Panics(t, func() {
		mustGetEvent(abi, "NonExistentEvent")
	})
}
