package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestKeeper_InsertGenesisHead(t *testing.T) {
	t.Parallel()

	ctx, keeper := createTestKeeper(t)

	// make sure the execution head does not exist
	_, err := keeper.GetExecutionHead(ctx)
	require.Error(t, err, "execution head should not exist")

	// insert genesis head
	dummyBlockHash := []byte("test")
	err = keeper.InsertGenesisHead(ctx, dummyBlockHash)
	require.NoError(t, err)

	// make sure the execution head is set correctly
	head, err := keeper.GetExecutionHead(ctx)
	require.NoError(t, err)
	require.NotNil(t, head, "execution head should exist")
	require.Equal(t, dummyBlockHash, head.GetBlockHash(), "block hash should match")

	// next try should fail because the genesis head already exists
	err = keeper.InsertGenesisHead(ctx, []byte("another hash"))
	require.Error(t, err, "genesis head should already exist")
}

func TestKeeper_updateExecutionHead(t *testing.T) {
	t.Parallel()

	ctx, keeper := createTestKeeper(t)

	// make sure the execution head does not exist
	_, err := keeper.GetExecutionHead(ctx)
	require.Error(t, err, "execution head should not exist")

	// insert genesis head
	dummyBlockHash := []byte("test")
	err = keeper.InsertGenesisHead(ctx, dummyBlockHash)
	require.NoError(t, err)

	// make sure the execution head is set correctly
	head, err := keeper.GetExecutionHead(ctx)
	require.NoError(t, err)
	require.NotNil(t, head, "execution head should exist")

	// update the execution head
	newBlockHash := common.BytesToHash([]byte("new hash"))
	err = keeper.UpdateExecutionHead(ctx, engine.ExecutableData{
		Number:    100,
		BlockHash: newBlockHash,
		Timestamp: 0,
	})
	require.NoError(t, err)

	// make sure the execution head is updated correctly
	head, err = keeper.GetExecutionHead(ctx)
	require.NoError(t, err)
	require.NotNil(t, head, "execution head should exist")
	require.Equal(t, newBlockHash.Bytes(), head.GetBlockHash(), "block hash should match")
	require.Equal(t, uint64(100), head.GetBlockHeight(), "block height should match")
}
