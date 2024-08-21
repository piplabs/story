package evm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/genutil/evm"
	"github.com/piplabs/story/lib/netconf"
	"github.com/piplabs/story/lib/tutil"

	_ "github.com/piplabs/story/client/app" // To init SDK config.
)

//go:generate go test . -golden -clean

func TestMakeGenesis(t *testing.T) {
	t.Parallel()

	genesis, err := evm.MakeGenesis(netconf.Staging)
	require.NoError(t, err)
	tutil.RequireGoldenJSON(t, genesis)
}
