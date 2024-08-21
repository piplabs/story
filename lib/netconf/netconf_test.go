package netconf_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/lib/netconf"
)

func TestConsensusSeeds(t *testing.T) {
	t.Parallel()

	require.Len(t, netconf.Testnet.Static().ConsensusSeeds(), 2)
}

func TestExecutionSeeds(t *testing.T) {
	t.Skip("testnet shutdown at the moment")
	t.Parallel()

	seeds := netconf.Testnet.Static().ExecutionSeeds()
	require.Len(t, seeds, 2)
	for _, seed := range seeds {
		node, err := enode.ParseV4(seed)
		require.NoError(t, err)

		require.EqualValues(t, 30303, node.TCP())
		require.EqualValues(t, 30303, node.UDP())
		t.Logf("Seed IP: %s: %s", node.IP(), seed)
		require.NotEmpty(t, node.IP())
	}
}
