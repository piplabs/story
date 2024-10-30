package keeper

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

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
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient"
	"github.com/piplabs/story/lib/ethclient/mock"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/tutil"

	"go.uber.org/mock/gomock"
)

func Test_msgServer_ExecutionPayload(t *testing.T) {
	t.Parallel()
	fastBackoffForT()

	cdc := getCodec(t)
	txConfig := authtx.NewTxConfig(cdc, nil)

	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockClient(ctrl)
	ak := moduletestutil.NewMockAccountKeeper(ctrl)
	esk := moduletestutil.NewMockEvmStakingKeeper(ctrl)
	uk := moduletestutil.NewMockUpgradeKeeper(ctrl)
	dk := moduletestutil.NewMockDistrKeeper(ctrl)

	cmtAPI := newMockCometAPI(t, nil)
	// set the header and proposer so we have the correct next proposer
	header := cmtproto.Header{Height: 1, AppHash: tutil.RandomHash().Bytes()}
	header.ProposerAddress = cmtAPI.validatorSet.Validators[0].Address
	nxtAddr, err := k1util.PubKeyToAddress(cmtAPI.validatorSet.Validators[1].PubKey)
	require.NoError(t, err)

	ctx, storeKey, storeService := setupCtxStore(t, &header)
	ctx = ctx.WithExecMode(sdk.ExecModeFinalize)
	// evmLogProc := mockLogProvider{deliverErr: errors.New("test error")}
	mockEngine, err := newMockEngineAPI(storeKey, 2)
	require.NoError(t, err)
	keeper, err := NewKeeper(cdc, storeService, &mockEngine, mockClient, txConfig, ak, esk, uk, dk)
	require.NoError(t, err)
	keeper.SetCometAPI(cmtAPI)
	keeper.SetValidatorAddress(nxtAddr)
	populateGenesisHead(ctx, t, keeper)

	msgSrv := NewMsgServerImpl(keeper)
	createValidPayload := func(c context.Context) (*etypes.Block, engine.PayloadID, []byte) {
		// get latest block to build on top
		latestBlock, err := mockEngine.HeaderByType(c, ethclient.HeadLatest)
		require.NoError(t, err)
		latestHeight := latestBlock.Number.Uint64()

		sdkCtx := sdk.UnwrapSDKContext(c)
		appHash := common.BytesToHash(sdkCtx.BlockHeader().AppHash)

		block, execPayload := mockEngine.nextBlock(t, latestHeight+1, uint64(time.Now().Unix()), latestBlock.Hash(), keeper.validatorAddr, &appHash)
		payloadID, err := ethclient.MockPayloadID(execPayload, &appHash)
		require.NoError(t, err)

		// Create execution payload message
		payloadData, err := json.Marshal(execPayload)
		require.NoError(t, err)

		return block, payloadID, payloadData
	}
	// createRandomEvents := func(c context.Context, blkHash common.Hash) []*types.EVMEvent {
	//	events, err := evmLogProc.Prepare(c, blkHash)
	//	require.NoError(t, err)
	//
	//	return events
	//}

	tcs := []struct {
		name          string
		setup         func(c context.Context) sdk.Context
		createPayload func(c context.Context) (*etypes.Block, engine.PayloadID, []byte)
		expectedError string
		postCheck     func(c context.Context, block *etypes.Block, payloadID engine.PayloadID)
	}{
		{
			name: "pass: valid payload",
			setup: func(c context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(c).Return(uint32(0), nil)
				esk.EXPECT().DequeueEligibleWithdrawals(c, gomock.Any()).Return(nil, nil)
				esk.EXPECT().DequeueEligibleRewardWithdrawals(c, gomock.Any()).Return(nil, nil)
				esk.EXPECT().ProcessStakingEvents(c, gomock.Any(), gomock.Any()).Return(nil)

				return sdk.UnwrapSDKContext(c)
			},
			createPayload: createValidPayload,
			postCheck: func(c context.Context, block *etypes.Block, payloadID engine.PayloadID) {
				gotPayload, err := mockEngine.GetPayloadV3(c, payloadID)
				require.NoError(t, err)
				require.Equal(t, gotPayload.ExecutionPayload.Number, block.Header().Number.Uint64())
				require.Equal(t, gotPayload.ExecutionPayload.BlockHash, block.Hash())
				require.Equal(t, gotPayload.ExecutionPayload.FeeRecipient, keeper.validatorAddr)
				require.Empty(t, gotPayload.ExecutionPayload.Withdrawals)
			},
		},
		{
			name: "fail: sdk exec mode is not finalize",
			setup: func(c context.Context) sdk.Context {
				return sdk.UnwrapSDKContext(c).WithExecMode(sdk.ExecModeCheck)
			},
			expectedError: "only allowed in finalize mode",
		},
		{
			name: "fail: no execution head",
			setup: func(c context.Context) sdk.Context {
				head, err := keeper.headTable.Get(c, executionHeadID)
				require.NoError(t, err)
				require.NoError(t, keeper.headTable.Delete(c, head))

				return sdk.UnwrapSDKContext(c)
			},
			createPayload: createValidPayload,
			expectedError: "not found",
		},
		{
			name: "fail: invalid payload - wrong payload number",
			createPayload: func(ctx context.Context) (*etypes.Block, engine.PayloadID, []byte) {
				latestBlock, err := mockEngine.HeaderByType(ctx, ethclient.HeadLatest)
				require.NoError(t, err)
				latestHeight := latestBlock.Number.Uint64()
				wrongNextHeight := latestHeight + 2

				sdkCtx := sdk.UnwrapSDKContext(ctx)
				appHash := common.BytesToHash(sdkCtx.BlockHeader().AppHash)

				block, execPayload := mockEngine.nextBlock(t, wrongNextHeight, uint64(time.Now().Unix()), latestBlock.Hash(), keeper.validatorAddr, &appHash)
				payloadID, err := ethclient.MockPayloadID(execPayload, &appHash)
				require.NoError(t, err)

				// Create execution payload message
				payloadData, err := json.Marshal(execPayload)
				require.NoError(t, err)

				return block, payloadID, payloadData
			},
			expectedError: "invalid proposed payload number",
		},
		{
			name: "fail: DequeueEligibleWithdrawals error",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().DequeueEligibleWithdrawals(ctx, gomock.Any()).Return(nil, errors.New("failed to dequeue"))

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload: createValidPayload,
			expectedError: "error on withdrawals dequeue",
		},
		{
			name: "fail: NewPayloadV3 returns status invalid",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().DequeueEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().DequeueEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				mockEngine.forceInvalidNewPayloadV3 = true

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload: createValidPayload,
			expectedError: "payload invalid",
		},
		{
			name: "fail: ForkchoiceUpdatedV3 returns status invalid",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().DequeueEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().DequeueEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				mockEngine.forceInvalidForkchoiceUpdatedV3 = true

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload: createValidPayload,
			expectedError: "payload invalid",
		},
		{
			name: "fail: ProcessStakingEvents error",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().DequeueEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().DequeueEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().ProcessStakingEvents(ctx, gomock.Any(), gomock.Any()).Return(errors.New("failed to process staking events"))

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload: createValidPayload,
			expectedError: "deliver staking-related event logs",
		},
		//{
		//	name: "fail: ProcessUpgradeEvents error",
		//	setup: func(ctx context.Context) sdk.Context {
		//		esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
		//		esk.EXPECT().DequeueEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
		//		esk.EXPECT().DequeueEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)
		//		esk.EXPECT().ProcessStakingEvents(ctx, gomock.Any(), gomock.Any()).Return(nil)
		//
		//		return sdk.UnwrapSDKContext(ctx)
		//	},
		//	createPayload: createValidPayload,
		//	createPrevPayloadEvents: func(_ context.Context, _ common.Hash) []*types.EVMEvent {
		//		// crate invalid upgrade event to trigger ProcessUpgradeEvents failure
		//		upgradeAbi, err := bindings.UpgradeEntrypointMetaData.GetAbi()
		//		require.NoError(t, err, "failed to load ABI")
		//		data, err := upgradeAbi.Events["SoftwareUpgrade"].Inputs.NonIndexed().Pack("test-upgrade", int64(0), "test-info")
		//		require.NoError(t, err)
		//
		//		return []*types.EVMEvent{{
		//			Address: nil, // nil address
		//			Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
		//			Data:    data,
		//		}}
		//	},
		//	expectedError: "deliver upgrade-related event logs",
		// },
	}

	for _, tc := range tcs {
		//nolint:tparallel // currently, we can't run the tests in parallel due to the shared mockEngine. don't know how to fix it yet, just disable parallel for now.
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			var payloadData []byte
			var payloadID engine.PayloadID
			var block *etypes.Block

			cachedCtx, _ := ctx.CacheContext()
			if tc.setup != nil {
				cachedCtx = tc.setup(cachedCtx)
			}
			if tc.createPayload != nil {
				block, payloadID, payloadData = tc.createPayload(cachedCtx)
			}

			resp, err := msgSrv.ExecutionPayload(cachedCtx, &types.MsgExecutionPayload{
				Authority:        authtypes.NewModuleAddress(types.ModuleName).String(),
				ExecutionPayload: payloadData,
			})
			if tc.expectedError != "" {
				require.ErrorContains(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				if tc.postCheck != nil {
					tc.postCheck(cachedCtx, block, payloadID)
				}
			}
		})
	}

	// now lets run optimistic flow
	// ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Second))
	// newPayload(ctx)
	// keeper.SetBuildOptimistic(true)
	// assertExecutionPayload(ctx)
}

// populateGenesisHead inserts the mock genesis execution head into the database.
func populateGenesisHead(ctx context.Context, t *testing.T, keeper *Keeper) {
	t.Helper()
	genesisBlock, err := ethclient.MockGenesisBlock()
	require.NoError(t, err)

	require.NoError(t, keeper.InsertGenesisHead(ctx, genesisBlock.Hash().Bytes()))
}

func Test_pushPayload(t *testing.T) {
	t.Parallel()

	newPayload := func(ctx context.Context, mockEngine mockEngineAPI, address common.Address) (engine.ExecutableData, engine.PayloadID) {
		// get latest block to build on top
		latestBlock, err := mockEngine.HeaderByType(ctx, ethclient.HeadLatest)
		require.NoError(t, err)
		latestHeight := latestBlock.Number.Uint64()

		sdkCtx := sdk.UnwrapSDKContext(ctx)
		appHash := common.BytesToHash(sdkCtx.BlockHeader().AppHash)

		_, execPayload := mockEngine.nextBlock(t, latestHeight+1, uint64(time.Now().Unix()), latestBlock.Hash(), address, &appHash)
		payloadID, err := ethclient.MockPayloadID(execPayload, &appHash)
		require.NoError(t, err)

		return execPayload, payloadID
	}
	type args struct {
		transformPayload func(*engine.ExecutableData)
		newPayloadV3Func func(context.Context, engine.ExecutableData, []common.Hash, *common.Hash) (engine.PayloadStatusV1, error)
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus string
	}{
		{
			name: "new payload error",
			args: args{
				newPayloadV3Func: func(context.Context, engine.ExecutableData, []common.Hash, *common.Hash) (engine.PayloadStatusV1, error) {
					return engine.PayloadStatusV1{}, errors.New("error")
				},
			},
			wantErr:    true,
			wantStatus: "",
		},
		{
			name: "new payload invalid",
			args: args{
				newPayloadV3Func: func(context.Context, engine.ExecutableData, []common.Hash, *common.Hash) (engine.PayloadStatusV1, error) {
					return engine.PayloadStatusV1{
						Status:          engine.INVALID,
						LatestValidHash: nil,
						ValidationError: nil,
					}, nil
				},
			},
			wantErr:    false,
			wantStatus: engine.INVALID,
		},
		{
			name: "new payload invalid val err",
			args: args{
				newPayloadV3Func: func(context.Context, engine.ExecutableData, []common.Hash, *common.Hash) (engine.PayloadStatusV1, error) {
					return engine.PayloadStatusV1{
						Status:          engine.INVALID,
						LatestValidHash: nil,
						ValidationError: func() *string { s := "error"; return &s }(),
					}, nil
				},
			},
			wantErr:    false,
			wantStatus: engine.INVALID,
		},
		{
			name: "new payload syncing",
			args: args{
				newPayloadV3Func: func(context.Context, engine.ExecutableData, []common.Hash, *common.Hash) (engine.PayloadStatusV1, error) {
					return engine.PayloadStatusV1{
						Status:          engine.SYNCING,
						LatestValidHash: nil,
						ValidationError: nil,
					}, nil
				},
			},
			wantErr:    false,
			wantStatus: engine.SYNCING,
		},
		{
			name: "new payload accepted",
			args: args{
				newPayloadV3Func: func(context.Context, engine.ExecutableData, []common.Hash, *common.Hash) (engine.PayloadStatusV1, error) {
					return engine.PayloadStatusV1{
						Status:          engine.ACCEPTED,
						LatestValidHash: nil,
						ValidationError: nil,
					}, nil
				},
			},
			wantErr:    false,
			wantStatus: engine.ACCEPTED,
		},
		{
			name:       "valid payload",
			args:       args{},
			wantErr:    false,
			wantStatus: engine.VALID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			appHash := tutil.RandomHash()
			ctx, storeKey := ctxWithAppHash(t, appHash)

			mockEngine, err := newMockEngineAPI(storeKey, 0)
			require.NoError(t, err)

			mockEngine.newPayloadV3Func = tt.args.newPayloadV3Func
			payload, payloadID := newPayload(ctx, mockEngine, common.Address{})
			if tt.args.transformPayload != nil {
				tt.args.transformPayload(&payload)
			}

			status, err := pushPayload(ctx, &mockEngine, payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("pushPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.wantStatus, status.Status)

			if status.Status == engine.VALID {
				want, err := mockEngine.GetPayloadV3(ctx, payloadID)
				require.NoError(t, err)
				if !reflect.DeepEqual(payload, *want.ExecutionPayload) {
					t.Errorf("pushPayload() got = %v, want %v", payload, want)
				}
			}
		})
	}
}
