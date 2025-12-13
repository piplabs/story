package keeper

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	moduletestutil "github.com/piplabs/story/client/x/evmengine/testutil"
	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient"
	"github.com/piplabs/story/lib/ethclient/mock"
	"github.com/piplabs/story/lib/expbackoff"
	"github.com/piplabs/story/lib/tutil"

	"go.uber.org/mock/gomock"
)

var maxAddr = common.MaxAddress

func Test_proposalServer_ExecutionPayload(t *testing.T) {
	t.Parallel()
	cdc := getCodec(t)
	txConfig := authtx.NewTxConfig(cdc, nil)

	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockClient(ctrl)
	ak := moduletestutil.NewMockAccountKeeper(ctrl)
	esk := moduletestutil.NewMockEvmStakingKeeper(ctrl)
	uk := moduletestutil.NewMockUpgradeKeeper(ctrl)
	dk := moduletestutil.NewMockDistrKeeper(ctrl)

	ctx, storeKey, storeService := setupCtxStore(t, &cmtproto.Header{AppHash: tutil.RandomHash().Bytes()})
	ctx = ctx.WithExecMode(sdk.ExecModeFinalize)
	evmLogProc := mockLogProvider{deliverErr: errors.New("test error")}
	mockEngine, err := newMockEngineAPI(storeKey, 0)
	require.NoError(t, err)
	keeper, err := NewKeeper(cdc, storeService, &mockEngine, mockClient, txConfig, ak, esk, uk, dk)
	require.NoError(t, err)
	populateGenesisHead(ctx, t, keeper)
	propSrv := NewProposalServer(keeper)

	keeper.SetValidatorAddress(common.BytesToAddress([]byte("test")))

	createValidPayload := func(c context.Context, withWithdrawals bool) (*etypes.Block, engine.PayloadID, []byte) {
		// get latest block to build on top
		latestBlock, err := mockEngine.HeaderByType(c, ethclient.HeadLatest)
		require.NoError(t, err)

		latestHeight := latestBlock.Number.Uint64()

		sdkCtx := sdk.UnwrapSDKContext(c)
		appHash := common.BytesToHash(sdkCtx.BlockHeader().AppHash)

		block, execPayload := mockEngine.nextBlock(t, latestHeight+1, uint64(time.Now().Unix()), latestBlock.Hash(), keeper.validatorAddr, &appHash)

		if withWithdrawals {
			execPayload.Withdrawals = []*etypes.Withdrawal{
				{
					Index:     uint64(0),
					Validator: uint64(0),
					Address:   common.MaxAddress,
					Amount:    uint64(0),
				},
			}
		}

		payloadID, err := ethclient.MockPayloadID(execPayload, &appHash)
		require.NoError(t, err)

		// Create execution payload message
		payloadData, err := json.Marshal(execPayload)
		require.NoError(t, err)

		return block, payloadID, payloadData
	}

	createRandomEvents := func(c context.Context, blkHash common.Hash) []*types.EVMEvent {
		events, err := evmLogProc.Prepare(c, blkHash)
		require.NoError(t, err)

		return events
	}

	tcs := []struct {
		name                    string
		setup                   func(c context.Context) sdk.Context
		createPayload           func(c context.Context, withWithdrawals bool) (*etypes.Block, engine.PayloadID, []byte)
		withWithdrawal          bool
		createPrevPayloadEvents func(c context.Context, blkHash common.Hash) []*types.EVMEvent
		expectedError           string
		postCheck               func(c context.Context, block *etypes.Block, payloadID engine.PayloadID)
		afterTest               func(c context.Context)
	}{
		{
			name: "pass: valid payload",
			setup: func(c context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(c).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(c, gomock.Any()).Return(nil, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(c, gomock.Any()).Return(nil, nil)

				mockEngine.filterLogsFunc = func(ctx context.Context, query ethereum.FilterQuery) ([]etypes.Log, error) {
					return []etypes.Log{
						{
							Address: maxAddr,
							Topics:  []common.Hash{common.MaxHash},
							TxHash:  common.MaxHash,
						},
					}, nil
				}

				return sdk.UnwrapSDKContext(c)
			},
			createPayload: createValidPayload,
			createPrevPayloadEvents: func(c context.Context, blkHash common.Hash) []*types.EVMEvent {
				return []*types.EVMEvent{
					{
						Address: maxAddr[:],
						Topics:  [][]byte{common.MaxHash[:]},
						TxHash:  common.MaxHash[:],
					},
				}
			},
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
			name: "fail: execution engine is syncing",
			setup: func(c context.Context) sdk.Context {
				keeper.SetIsExecEngSyncing(true)

				return sdk.UnwrapSDKContext(c)
			},
			expectedError: "execution engine is syncing",
			afterTest: func(c context.Context) {
				keeper.SetIsExecEngSyncing(false)
			},
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
			createPayload: func(ctx context.Context, withWithdrawal bool) (*etypes.Block, engine.PayloadID, []byte) {
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
			expectedError: "parse and verify payload",
		},
		{
			name: "fail: MaxWithdrawalPerBlock error",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), errors.New("failed to get max withdrawal per block"))
				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload: createValidPayload,
			expectedError: "get max withdrawals per block",
		},
		{
			name: "fail: PeekEligibleWithdrawals error",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(nil, errors.New("failed to peek"))

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload: createValidPayload,
			expectedError: "peek withdrawals",
		},
		{
			name: "fail: peeked withdrawals are greater than withdrawals in payload",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(etypes.Withdrawals{
					&etypes.Withdrawal{
						Index:     uint64(0),
						Validator: uint64(0),
						Address:   common.MaxAddress,
						Amount:    uint64(0),
					},
				}, nil)

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload: createValidPayload,
			expectedError: "expected withdrawals 1 should not greater than actual withdrawals 0",
		},
		{
			name: "fail: PeekEligibleRewardWithdrawals error",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, errors.New("failed to dequeue"))

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload: createValidPayload,
			expectedError: "peek reward withdrawals",
		},
		{
			name: "fail: peeked total withdrawals are greater than total withdrawals in payload",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(etypes.Withdrawals{
					&etypes.Withdrawal{
						Index:     uint64(0),
						Validator: uint64(0),
						Address:   common.MaxAddress,
						Amount:    uint64(0),
					},
				}, nil)

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload: createValidPayload,
			expectedError: "expected total withdrawals 1 should equal to actual withdrawals 0",
		},
		{
			name: "fail: withdrawals index mismatch",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(etypes.Withdrawals{
					&etypes.Withdrawal{
						Index:     uint64(1),
						Validator: uint64(0),
						Address:   common.MaxAddress,
						Amount:    uint64(0),
					},
				}, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload:  createValidPayload,
			withWithdrawal: true,
			expectedError:  "invalid withdrawal index",
		},
		{
			name: "fail: withdrawals type mismatch",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(etypes.Withdrawals{
					&etypes.Withdrawal{
						Index:     uint64(0),
						Validator: uint64(1),
						Address:   common.MaxAddress,
						Amount:    uint64(0),
					},
				}, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload:  createValidPayload,
			withWithdrawal: true,
			expectedError:  "invalid withdrawal type",
		},
		{
			name: "fail: withdrawals address mismatch",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(etypes.Withdrawals{
					&etypes.Withdrawal{
						Index:     uint64(0),
						Validator: uint64(0),
						Address:   common.Address{},
						Amount:    uint64(0),
					},
				}, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload:  createValidPayload,
			withWithdrawal: true,
			expectedError:  "invalid withdrawal address",
		},
		{
			name: "fail: withdrawals amount mismatch",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(etypes.Withdrawals{
					&etypes.Withdrawal{
						Index:     uint64(0),
						Validator: uint64(0),
						Address:   common.MaxAddress,
						Amount:    uint64(1),
					},
				}, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload:  createValidPayload,
			withWithdrawal: true,
			expectedError:  "invalid withdrawal amount",
		},
		{
			name: "fail: reward withdrawals index mismatch",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(etypes.Withdrawals{
					&etypes.Withdrawal{
						Index:     uint64(1),
						Validator: uint64(0),
						Address:   common.MaxAddress,
						Amount:    uint64(0),
					},
				}, nil)

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload:  createValidPayload,
			withWithdrawal: true,
			expectedError:  "invalid withdrawal index",
		},
		{
			name: "fail: reward withdrawals type mismatch",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(etypes.Withdrawals{
					&etypes.Withdrawal{
						Index:     uint64(0),
						Validator: uint64(1),
						Address:   common.MaxAddress,
						Amount:    uint64(0),
					},
				}, nil)

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload:  createValidPayload,
			withWithdrawal: true,
			expectedError:  "invalid withdrawal type",
		},
		{
			name: "fail: reward withdrawals address mismatch",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(etypes.Withdrawals{
					&etypes.Withdrawal{
						Index:     uint64(0),
						Validator: uint64(0),
						Address:   common.Address{},
						Amount:    uint64(0),
					},
				}, nil)

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload:  createValidPayload,
			withWithdrawal: true,
			expectedError:  "invalid withdrawal address",
		},
		{
			name: "fail: reward withdrawals amount mismatch",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(etypes.Withdrawals{
					&etypes.Withdrawal{
						Index:     uint64(0),
						Validator: uint64(0),
						Address:   common.MaxAddress,
						Amount:    uint64(1),
					},
				}, nil)

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload:  createValidPayload,
			withWithdrawal: true,
			expectedError:  "invalid withdrawal amount",
		},
		{
			name: "fail: NewPayloadV3 returns status invalid",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)

				mockEngine.forceInvalidNewPayloadV3 = true

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload: createValidPayload,
			expectedError: "payload invalid",
		},
		{
			name: "fail: NewPayloadV3 returns status syncing",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				mockEngine.newPayloadV3Func = func(ctx context.Context, data engine.ExecutableData, hashes []common.Hash, hash *common.Hash) (engine.PayloadStatusV1, error) {
					return engine.PayloadStatusV1{
						Status: engine.SYNCING,
					}, nil
				}

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload: createValidPayload,
			expectedError: "execution engine is syncing",
		},
		{
			name: "fail: invalid event",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				mockEngine.newPayloadV3Func = func(ctx context.Context, data engine.ExecutableData, hashes []common.Hash, hash *common.Hash) (engine.PayloadStatusV1, error) {
					return engine.PayloadStatusV1{
						Status:          engine.VALID,
						LatestValidHash: nil,
						ValidationError: nil,
					}, nil
				}
				mockEngine.filterLogsFunc = func(ctx context.Context, query ethereum.FilterQuery) ([]etypes.Log, error) {
					return []etypes.Log{{
						Address: common.MaxAddress,
					}}, nil
				}

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload: createValidPayload,
			expectedError: "verify event",
		},
		{
			name: "fail: event count mismatch",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				mockEngine.newPayloadV3Func = func(ctx context.Context, data engine.ExecutableData, hashes []common.Hash, hash *common.Hash) (engine.PayloadStatusV1, error) {
					return engine.PayloadStatusV1{
						Status:          engine.VALID,
						LatestValidHash: nil,
						ValidationError: nil,
					}, nil
				}
				mockEngine.filterLogsFunc = func(ctx context.Context, query ethereum.FilterQuery) ([]etypes.Log, error) {
					return []etypes.Log{
						{
							Address: common.MaxAddress,
							Topics:  []common.Hash{common.MaxHash},
						},
						{
							Address: common.MaxAddress,
							Topics:  []common.Hash{common.HexToHash("0x1")},
						},
					}, nil
				}

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload:           createValidPayload,
			createPrevPayloadEvents: createRandomEvents,
			expectedError:           "count mismatch",
		},
		{
			name: "fail: event log mismatch",
			setup: func(ctx context.Context) sdk.Context {
				esk.EXPECT().MaxWithdrawalPerBlock(ctx).Return(uint32(0), nil)
				esk.EXPECT().PeekEligibleWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				esk.EXPECT().PeekEligibleRewardWithdrawals(ctx, gomock.Any()).Return(nil, nil)
				mockEngine.newPayloadV3Func = func(ctx context.Context, data engine.ExecutableData, hashes []common.Hash, hash *common.Hash) (engine.PayloadStatusV1, error) {
					return engine.PayloadStatusV1{
						Status:          engine.VALID,
						LatestValidHash: nil,
						ValidationError: nil,
					}, nil
				}
				mockEngine.filterLogsFunc = func(ctx context.Context, query ethereum.FilterQuery) ([]etypes.Log, error) {
					return []etypes.Log{{
						Address: common.MaxAddress,
						Topics:  []common.Hash{common.MaxHash},
					}}, nil
				}

				return sdk.UnwrapSDKContext(ctx)
			},
			createPayload:           createValidPayload,
			createPrevPayloadEvents: createRandomEvents,
			expectedError:           "log mismatch",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var (
				payloadData []byte
				payloadID   engine.PayloadID
				block       *etypes.Block
				events      []*types.EVMEvent
			)

			cachedCtx, _ := ctx.CacheContext()
			if tc.setup != nil {
				cachedCtx = tc.setup(cachedCtx)
			}

			if tc.createPayload != nil {
				block, payloadID, payloadData = tc.createPayload(cachedCtx, tc.withWithdrawal)
			}

			if tc.createPrevPayloadEvents != nil {
				events = tc.createPrevPayloadEvents(cachedCtx, block.Hash())
			}

			resp, err := propSrv.ExecutionPayload(cachedCtx, &types.MsgExecutionPayload{
				Authority:         authtypes.NewModuleAddress(types.ModuleName).String(),
				ExecutionPayload:  payloadData,
				PrevPayloadEvents: events,
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

			if tc.afterTest != nil {
				tc.afterTest(cachedCtx)
			}
		})
	}
}

func fastBackoffForT() {
	backoffFuncMu.Lock()
	defer backoffFuncMu.Unlock()

	backoffFunc = func(context.Context, ...func(*expbackoff.Config)) func() {
		return func() {}
	}
}
