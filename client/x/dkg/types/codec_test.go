package types_test

import (
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

func TestRegisterInterfaces(t *testing.T) {
	t.Parallel()

	registry := cdctypes.NewInterfaceRegistry()

	// Should not panic
	require.NotPanics(t, func() {
		types.RegisterInterfaces(registry)
	})
}
