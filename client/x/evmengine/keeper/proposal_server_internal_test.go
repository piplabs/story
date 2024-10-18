package keeper

import (
	"context"
	"encoding/json"
	"testing"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	moduletestutil "github.com/piplabs/story/client/x/evmengine/testutil"
	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/ethclient"
	"github.com/piplabs/story/lib/ethclient/mock"
	"github.com/piplabs/story/lib/expbackoff"
	"github.com/piplabs/story/lib/tutil"

	"go.uber.org/mock/gomock"
)

func Test_proposalServer_ExecutionPayload(t *testing.T) {
	t.Parallel()
	cdc := getCodec(t)
	txConfig := authtx.NewTxConfig(cdc, nil)

	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockClient(ctrl)
	ak := moduletestutil.NewMockAccountKeeper(ctrl)
	esk := moduletestutil.NewMockEvmStakingKeeper(ctrl)
	sk := moduletestutil.NewMockSignalKeeper(ctrl)
	mk := moduletestutil.NewMockMintKeeper(ctrl)
	esk.EXPECT().PeekEligibleWithdrawals(gomock.Any()).Return(nil, nil).AnyTimes()

	sdkCtx, storeKey, storeService := setupCtxStore(t, &cmtproto.Header{AppHash: tutil.RandomHash().Bytes()})
	sdkCtx = sdkCtx.WithExecMode(sdk.ExecModeFinalize)
	mockEngine, err := newMockEngineAPI(storeKey, 0)
	require.NoError(t, err)
	keeper, err := NewKeeper(cdc, storeService, &mockEngine, mockClient, txConfig, ak, esk, mk, sk)
	require.NoError(t, err)
	populateGenesisHead(sdkCtx, t, keeper)
	propSrv := NewProposalServer(keeper)

	keeper.SetValidatorAddress(common.BytesToAddress([]byte("test")))

	var payloadData []byte
	var payloadID engine.PayloadID
	var latestHeight uint64
	var block *etypes.Block
	newPayload := func(ctx context.Context) {
		// get latest block to build on top
		latestBlock, err := mockEngine.HeaderByType(ctx, ethclient.HeadLatest)
		require.NoError(t, err)
		latestHeight = latestBlock.Number.Uint64()

		sdkCtx := sdk.UnwrapSDKContext(ctx)
		appHash := common.BytesToHash(sdkCtx.BlockHeader().AppHash)

		b, execPayload := mockEngine.nextBlock(
			t,
			latestHeight+1,
			uint64(sdkCtx.BlockHeader().Time.Unix()),
			latestBlock.Hash(),
			keeper.validatorAddr,
			&appHash,
		)
		block = b

		payloadID, err = ethclient.MockPayloadID(execPayload, &appHash)
		require.NoError(t, err)

		// Create execution payload message
		payloadData, err = json.Marshal(execPayload)
		require.NoError(t, err)
	}

	assertExecutionPayload := func(ctx context.Context) {
		resp, err := propSrv.ExecutionPayload(ctx, &types.MsgExecutionPayload{
			Authority:        authtypes.NewModuleAddress(types.ModuleName).String(),
			ExecutionPayload: payloadData,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)

		gotPayload, err := mockEngine.GetPayloadV3(ctx, payloadID)
		require.NoError(t, err)
		require.Equal(t, latestHeight+1, gotPayload.ExecutionPayload.Number)
		require.Equal(t, block.Hash(), gotPayload.ExecutionPayload.BlockHash)
		require.Equal(t, keeper.validatorAddr, gotPayload.ExecutionPayload.FeeRecipient)
		require.Empty(t, gotPayload.ExecutionPayload.Withdrawals)
	}

	newPayload(sdkCtx)
	assertExecutionPayload(sdkCtx)
}

func fastBackoffForT() {
	backoffFuncMu.Lock()
	defer backoffFuncMu.Unlock()
	backoffFunc = func(context.Context, ...func(*expbackoff.Config)) func() {
		return func() {}
	}
}
