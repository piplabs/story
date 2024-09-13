package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	moduletestutil "github.com/piplabs/story/client/x/evmengine/testutil"
	"github.com/piplabs/story/lib/ethclient/mock"
)

func createTestKeeper(t *testing.T) (sdk.Context, *Keeper) {
	t.Helper()
	cdc := getCodec(t)
	txConfig := authtx.NewTxConfig(cdc, nil)

	mockEngine, err := newMockEngineAPI(0)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockClient(ctrl)
	ak := moduletestutil.NewMockAccountKeeper(ctrl)
	esk := moduletestutil.NewMockEvmStakingKeeper(ctrl)
	uk := moduletestutil.NewMockUpgradeKeeper(ctrl)
	ctx, storeService := setupCtxStore(t, nil)
	ctx = ctx.WithExecMode(sdk.ExecModeFinalize)
	keeper, err := NewKeeper(cdc, storeService, &mockEngine, mockClient, txConfig, ak, esk, uk)
	require.NoError(t, err)

	return ctx, keeper
}

func TestKeeper_InsertGenesisHead(t *testing.T) {
	t.Parallel()

	ctx, keeper := createTestKeeper(t)

	// make sure the execution head does not exist
	_, err := keeper.getExecutionHead(ctx)
	require.Error(t, err, "execution head should not exist")

	// insert genesis head
	dummyBlockHash := []byte("test")
	err = keeper.InsertGenesisHead(ctx, dummyBlockHash)
	require.NoError(t, err)

	// make sure the execution head is set correctly
	head, err := keeper.getExecutionHead(ctx)
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
	_, err := keeper.getExecutionHead(ctx)
	require.Error(t, err, "execution head should not exist")

	// insert genesis head
	dummyBlockHash := []byte("test")
	err = keeper.InsertGenesisHead(ctx, dummyBlockHash)
	require.NoError(t, err)

	// make sure the execution head is set correctly
	head, err := keeper.getExecutionHead(ctx)
	require.NoError(t, err)
	require.NotNil(t, head, "execution head should exist")

	// update the execution head
	newBlockHash := common.BytesToHash([]byte("new hash"))
	err = keeper.updateExecutionHead(ctx, engine.ExecutableData{
		Number:    100,
		BlockHash: newBlockHash,
		Timestamp: 0,
	})
	require.NoError(t, err)

	// make sure the execution head is updated correctly
	head, err = keeper.getExecutionHead(ctx)
	require.NoError(t, err)
	require.NotNil(t, head, "execution head should exist")
	require.Equal(t, newBlockHash.Bytes(), head.GetBlockHash(), "block hash should match")
	require.Equal(t, uint64(100), head.GetBlockHeight(), "block height should match")
}
