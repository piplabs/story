package keeper

import (
	"context"
	"testing"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	moduletestutil "github.com/piplabs/story/client/x/evmengine/testutil"
	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/ethclient/mock"

	"go.uber.org/mock/gomock"
)

func createTestKeeper(t *testing.T) (context.Context, *Keeper) {
	t.Helper()

	cdc := getCodec(t)
	txConfig := authtx.NewTxConfig(cdc, nil)
	mockEngine, err := newMockEngineAPI(0)
	require.NoError(t, err)

	cmtAPI := newMockCometAPI(t, nil)
	header := cmtproto.Header{Height: 1}

	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockClient(ctrl)
	ak := moduletestutil.NewMockAccountKeeper(ctrl)
	esk := moduletestutil.NewMockEvmStakingKeeper(ctrl)
	uk := moduletestutil.NewMockUpgradeKeeper(ctrl)

	ctx, storeService := setupCtxStore(t, &header)

	keeper, err := NewKeeper(cdc, storeService, &mockEngine, mockClient, txConfig, ak, esk, uk)
	require.NoError(t, err)
	keeper.SetCometAPI(cmtAPI)

	return ctx, keeper
}

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
