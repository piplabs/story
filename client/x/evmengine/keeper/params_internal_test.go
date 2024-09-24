package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmengine/types"
)

func TestKeeper_ExecutionBlockHash(t *testing.T) {
	t.Parallel()
	ctx, keeper := createTestKeeper(t)

	// check existing execution block hash
	execHash, err := keeper.ExecutionBlockHash(ctx)
	require.NoError(t, err)
	require.Nil(t, execHash, "execution block hash should be nil because it is not set yet")

	// set execution block hash
	dummyHash := common.HexToHash("0x047e24c3455107d87c68dffa307b3b7fa1877f3e9d7f30c7ee359f2eff3a75d9")
	require.NoError(t, keeper.SetParams(ctx, types.Params{ExecutionBlockHash: dummyHash.Bytes()}))

	// check execution block hash whether it is set correctly
	execHash, err = keeper.ExecutionBlockHash(ctx)
	require.NoError(t, err)
	require.Equal(t, dummyHash.Bytes(), execHash, "execution block hash should be equal to the dummy hash")
}

func TestKeeper_GetSetParams(t *testing.T) {
	t.Parallel()
	ctx, keeper := createTestKeeper(t)

	// check existing params
	params, err := keeper.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, types.DefaultParams(), params, "params should be equal to the default params")

	// set execution block hash
	dummyHash := common.HexToHash("0x047e24c3455107d87c68dffa307b3b7fa1877f3e9d7f30c7ee359f2eff3a75d9")
	require.NoError(t, keeper.SetParams(ctx, types.Params{ExecutionBlockHash: dummyHash.Bytes()}))

	// check params whether it is set correctly
	params, err = keeper.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, types.Params{ExecutionBlockHash: dummyHash.Bytes()}, params)
}
