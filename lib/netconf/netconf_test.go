package netconf_test

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/stretchr/testify/require"

	"github.com/storyprotocol/iliad/e2e/app/key"
	"github.com/storyprotocol/iliad/e2e/manifests"
	"github.com/storyprotocol/iliad/e2e/types"
	"github.com/storyprotocol/iliad/lib/k1util"
	"github.com/storyprotocol/iliad/lib/netconf"
	"github.com/storyprotocol/iliad/lib/tutil"
)

//go:generate go test -golden -run=TestGenConsSeeds

// TestGenConsSeeds generates <network>/consensus-seeds.txt by loading e2e manifests and parsing seed* p2p_consensus keys.
func TestGenConsSeeds(t *testing.T) {
	t.Parallel()
	tests := []struct {
		network      netconf.ID
		manifestFunc func() []byte
	}{
		{
			network:      netconf.Testnet,
			manifestFunc: manifests.Testnet,
		},
	}
	for _, test := range tests {
		t.Run(test.network.String(), func(t *testing.T) {
			t.Parallel()
			var manifest types.Manifest
			//nolint:musttag // ignore perturb json
			_, err := toml.Decode(string(test.manifestFunc()), &manifest)
			require.NoError(t, err)

			var peers []string
			for _, node := range sortedKeys(manifest.Keys) {
				if !strings.HasPrefix(node, "seed") {
					continue
				}

				for typ, addr := range manifest.Keys[node] {
					if typ != key.P2PConsensus {
						continue
					}

					peers = append(peers, fmt.Sprintf("%s@%s.%s.storyprotocol.xyz", addr, node, test.network)) // ABCDE123@seed01.staging.storyprotocol.xyz
				}
			}

			seeds := strings.Join(peers, "\n")
			seedsFile := fmt.Sprintf("../%s/consensus-seeds.txt", test.network)
			tutil.RequireGoldenBytes(t, []byte(seeds), tutil.WithFilename(seedsFile))
		})
	}
}

var genExecutionSeeds = flag.Bool("gen-execution-seeds", false, "Enable to generate execution-seeds.txt. Note this requires GCP secret manager read-access")

func TestGenExecutionSeeds(t *testing.T) {
	t.Parallel()
	if !*genExecutionSeeds {
		t.Skip("Skipping since --gen-execution-seeds=false")
		return
	}

	tests := []struct {
		network      netconf.ID
		manifestFunc func() []byte
	}{
		{
			network:      netconf.Testnet,
			manifestFunc: manifests.Testnet,
		},
	}
	for _, test := range tests {
		t.Run(test.network.String(), func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			var manifest types.Manifest
			//nolint:musttag // ignore perturb json
			_, err := toml.Decode(string(test.manifestFunc()), &manifest)
			require.NoError(t, err)

			var peers []string
			for _, node := range sortedKeys(manifest.Keys) {
				if !strings.HasPrefix(node, "seed") {
					continue
				}

				for typ, addr := range manifest.Keys[node] {
					if typ != key.P2PExecution {
						continue
					}

					key, err := key.Download(ctx, test.network, node, typ, addr)
					require.NoError(t, err)

					stdPrivKey, err := key.ECDSA()
					require.NoError(t, err)

					pubkey64 := k1util.PubKeyToBytes64(&stdPrivKey.PublicKey)
					pubkeyHex := hex.EncodeToString(pubkey64)
					nodeName := strings.TrimSuffix(node, "_evm")

					peers = append(peers, fmt.Sprintf("enode://%s@%s.%s.storyprotocol:30303", pubkeyHex, nodeName, test.network))
				}
			}

			seeds := strings.Join(peers, "\n")
			seedsFile := fmt.Sprintf("../%s/execution-seeds.txt", test.network)
			tutil.RequireGoldenBytes(t, []byte(seeds), tutil.WithFilename(seedsFile))
		})
	}
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

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
