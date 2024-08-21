package e2e_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/e2e/app/geth"
	"github.com/piplabs/story/lib/ethclient"
)

// TestGethConfig ensure that the geth config is setup correctly.
func TestGethConfig(t *testing.T) {
	t.Parallel()
	testIliadEVM(t, func(t *testing.T, client ethclient.Client) {
		t.Helper()
		ctx := context.Background()

		cfg := geth.MakeGethConfig(geth.Config{})

		block, err := client.BlockByNumber(ctx, big.NewInt(1))
		require.NoError(t, err)

		require.EqualValues(t, int(cfg.Eth.Miner.GasCeil), int(block.GasLimit()))
		require.Equal(t, big.NewInt(0), block.Difficulty())
	})
}
