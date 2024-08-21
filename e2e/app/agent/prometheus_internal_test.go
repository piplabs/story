package agent

import (
	"context"
	"testing"

	e2e "github.com/cometbft/cometbft/test/e2e/pkg"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/e2e/types"
	"github.com/piplabs/story/lib/netconf"
	"github.com/piplabs/story/lib/tutil"
)

//go:generate go test . -golden -clean

func TestPromGen(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		network      netconf.ID
		nodes        []string
		newNodes     []string
		geths        []string
		newGeths     []string
		hostname     string
		agentSecrets bool
	}{
		{
			name:         "manifest1",
			network:      netconf.Devnet,
			nodes:        []string{"validator01", "validator02"},
			hostname:     "localhost",
			newNodes:     []string{"validator01"},
			geths:        []string{"iliad_evm"},
			newGeths:     []string{"iliad_evm"},
			agentSecrets: false,
		},
		{
			name:         "manifest2",
			network:      netconf.Staging,
			nodes:        []string{"validator01", "validator02", "fullnode03"},
			hostname:     "vm",
			newNodes:     []string{"fullnode04"},
			geths:        []string{"validator01_evm", "validator02_evm", "validator03_evm"},
			newGeths:     []string{"fullnode04_evm"},
			agentSecrets: true,
		},
		{
			name:         "manifest3",
			network:      netconf.Devnet,
			nodes:        []string{"validator01", "validator02"},
			hostname:     "localhost",
			newNodes:     []string{"validator01"},
			agentSecrets: false,
		},
		{
			name:         "manifest4",
			network:      netconf.Staging,
			nodes:        []string{"validator01", "validator02", "fullnode03"},
			hostname:     "vm",
			newNodes:     []string{"fullnode04"},
			agentSecrets: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			var nodes []*e2e.Node
			for _, name := range test.nodes {
				nodes = append(nodes, &e2e.Node{Name: name})
			}

			var geths []types.IliadEVM
			for _, name := range test.geths {
				geths = append(geths, types.IliadEVM{InstanceName: name})
			}

			testnet := types.Testnet{
				Network: test.network,
				Testnet: &e2e.Testnet{
					Name:  test.name,
					Nodes: nodes,
				},
				IliadEVMs: geths,
			}

			var agentSecrets Secrets
			if test.agentSecrets {
				agentSecrets = Secrets{
					URL:  "https://grafana.com",
					User: "admin",
					Pass: "password",
				}
			}

			cfg1, err := genPromConfig(ctx, testnet, agentSecrets, test.hostname)
			require.NoError(t, err)

			cfg2 := ConfigForHost(cfg1, test.hostname+"-2", test.newNodes, test.newGeths)

			t.Run("gen", func(t *testing.T) {
				t.Parallel()
				tutil.RequireGoldenBytes(t, cfg1)
			})

			t.Run("update", func(t *testing.T) {
				t.Parallel()
				tutil.RequireGoldenBytes(t, cfg2)
			})
		})
	}
}
