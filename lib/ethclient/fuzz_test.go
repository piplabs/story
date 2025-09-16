package ethclient_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/lib/ethclient"
)

func TestFuzzer(t *testing.T) {
	t.Parallel()
	f := ethclient.NewFuzzer(0)

	var payload engine.ExecutableData
	f.Fuzz(&payload)

	// Ensure the fuzzed payload is valid by converting it to a block.
	_, err := engine.ExecutableDataToBlock(payload, nil, nil, nil)
	require.NoError(t, err)
}
