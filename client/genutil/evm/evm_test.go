package evm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/storyprotocol/iliad/client/genutil/evm"
	"github.com/storyprotocol/iliad/lib/netconf"
	"github.com/storyprotocol/iliad/lib/tutil"

	_ "github.com/storyprotocol/iliad/client/app" // To init SDK config.
)

//go:generate go test . -golden -clean

func TestMakeGenesis(t *testing.T) {
	t.Parallel()

	genesis, err := evm.MakeGenesis(netconf.Staging)
	require.NoError(t, err)
	tutil.RequireGoldenJSON(t, genesis)
}
