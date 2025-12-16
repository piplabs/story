package keeper

import (
	"context"
	"testing"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	moduletestutil "github.com/piplabs/story/client/x/evmengine/testutil"
	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient/mock"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/netconf"
	"github.com/piplabs/story/lib/tutil"

	"go.uber.org/mock/gomock"
)

const (
	dummyAddressHex = "0x1398C32A45Bc409b6C652E25bb0a3e702492A4ab"
)

var (
	dummyContractAddress = common.HexToAddress(dummyAddressHex)
	dummyHash            = common.HexToHash("0x1398C32A45Bc409b6C652E25bb0a3e702492A4ab")
)

func TestKeeper_ProcessSoftwareUpgrade(t *testing.T) {
	tcs := []struct {
		name           string
		ev             func() *bindings.UpgradeEntrypointSoftwareUpgrade
		setup          func(ctx context.Context, keeper *Keeper) sdk.Context
		setupMock      func(uk *moduletestutil.MockUpgradeKeeper)
		expectedErr    string
		expectedResult *upgradetypes.Plan
	}{
		{
			name: "pass: valid software upgrade event - before Terence upgrade",
			ev: func() *bindings.UpgradeEntrypointSoftwareUpgrade {
				return &bindings.UpgradeEntrypointSoftwareUpgrade{
					Name:   "test-upgrade",
					Height: 1,
					Info:   "test-info",
				}
			},
			setupMock: func(uk *moduletestutil.MockUpgradeKeeper) {
				uk.EXPECT().ScheduleUpgrade(gomock.Any(), gomock.Any()).Return(nil)
				uk.EXPECT().DumpUpgradeInfoToDisk(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "pass: valid software upgrade event - after Terence upgrade",
			ev: func() *bindings.UpgradeEntrypointSoftwareUpgrade {
				return &bindings.UpgradeEntrypointSoftwareUpgrade{
					Name:   "test-upgrade",
					Height: 51,
					Info:   "test-info",
				}
			},
			setup: func(ctx context.Context, _ *Keeper) sdk.Context {
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				sdkCtx = sdkCtx.WithBlockHeight(51)

				return sdkCtx
			},
			setupMock: func(uk *moduletestutil.MockUpgradeKeeper) {
				uk.EXPECT().DumpUpgradeInfoToDisk(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResult: &upgradetypes.Plan{
				Name:   "test-upgrade",
				Height: 51,
				Info:   "test-info",
			},
		},
		// Fail cases: The following test cases simulate basic error scenarios.
		// Since a mocked upgrade keeper is used, not all error cases can be tested here.
		// Comprehensive error testing would require the real upgrade keeper, which is beyond the scope of this unit test.
		{
			name: "fail: check if Terence upgrade is activated or not",
			ev: func() *bindings.UpgradeEntrypointSoftwareUpgrade {
				return &bindings.UpgradeEntrypointSoftwareUpgrade{
					Name:   "test-upgrade",
					Height: 1,
					Info:   "test-info",
				}
			},
			setup: func(ctx context.Context, _ *Keeper) sdk.Context {
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				sdkCtx = sdkCtx.WithChainID("invalid-chain-id")

				return sdkCtx
			},
			expectedErr: "failed to check Terence upgrade height",
		},
		{
			name: "fail: invalid upgrade event before Terence upgrade - height is 0",
			ev: func() *bindings.UpgradeEntrypointSoftwareUpgrade {
				return &bindings.UpgradeEntrypointSoftwareUpgrade{
					Name:   "test-upgrade",
					Height: 0,
					Info:   "test-info",
				}
			},
			setupMock: func(uk *moduletestutil.MockUpgradeKeeper) {
				uk.EXPECT().ScheduleUpgrade(gomock.Any(), gomock.Any()).Return(errors.New("height must be greater than 0"))
			},
			expectedErr: "height must be greater than 0",
		},
		{
			name: "fail: invalid upgrade event before Terence upgrade - name is empty",
			ev: func() *bindings.UpgradeEntrypointSoftwareUpgrade {
				return &bindings.UpgradeEntrypointSoftwareUpgrade{
					Name:   "",
					Height: 1,
					Info:   "test-info",
				}
			},
			setupMock: func(uk *moduletestutil.MockUpgradeKeeper) {
				uk.EXPECT().ScheduleUpgrade(gomock.Any(), gomock.Any()).Return(errors.New("name cannot be empty"))
			},
			expectedErr: "name cannot be empty",
		},
		{
			name: "fail: invalid upgrade event after Terence upgrade - height is 0",
			ev: func() *bindings.UpgradeEntrypointSoftwareUpgrade {
				return &bindings.UpgradeEntrypointSoftwareUpgrade{
					Name:   "test-upgrade",
					Height: 0,
					Info:   "test-info",
				}
			},
			setup: func(ctx context.Context, keeper *Keeper) sdk.Context {
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				sdkCtx = sdkCtx.WithBlockHeight(51)

				return sdkCtx
			},
			expectedErr: "height must be greater than 0",
		},
		{
			name: "fail: invalid upgrade event after Terence upgrade - prior to the current block",
			ev: func() *bindings.UpgradeEntrypointSoftwareUpgrade {
				return &bindings.UpgradeEntrypointSoftwareUpgrade{
					Name:   "test-upgrade",
					Height: 20,
					Info:   "test-info",
				}
			},
			setup: func(ctx context.Context, keeper *Keeper) sdk.Context {
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				sdkCtx = sdkCtx.WithBlockHeight(51)

				return sdkCtx
			},
			expectedErr: "failed to set pending upgrade",
		},
		{
			name: "fail: invalid upgrade event after Terence upgrade - name is empty",
			ev: func() *bindings.UpgradeEntrypointSoftwareUpgrade {
				return &bindings.UpgradeEntrypointSoftwareUpgrade{
					Name:   "",
					Height: 1,
					Info:   "test-info",
				}
			},
			setup: func(ctx context.Context, keeper *Keeper) sdk.Context {
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				sdkCtx = sdkCtx.WithBlockHeight(51)

				return sdkCtx
			},
			expectedErr: "name cannot be empty",
		},
		{
			name: "fail: have already pending upgrade",
			ev: func() *bindings.UpgradeEntrypointSoftwareUpgrade {
				return &bindings.UpgradeEntrypointSoftwareUpgrade{
					Name:   "test-upgrade",
					Height: 51,
					Info:   "test-info",
				}
			},
			setup: func(ctx context.Context, keeper *Keeper) sdk.Context {
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				sdkCtx = sdkCtx.WithBlockHeight(51)

				err := keeper.SetPendingUpgrade(sdkCtx, upgradetypes.Plan{
					Name:   "existing-upgrade",
					Height: 51,
					Info:   "existing-upgrade-info",
				})
				require.NoError(t, err)

				return sdkCtx
			},
			expectedErr: types.ErrUpgradePending.Error(),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			keeper, ctx, ctrl, uk := setupTestEnvironment(t)
			t.Cleanup(ctrl.Finish)

			if tc.setup != nil {
				ctx = tc.setup(ctx, keeper)
			}

			if tc.setupMock != nil {
				tc.setupMock(uk)
			}

			err := keeper.ProcessSoftwareUpgrade(ctx, tc.ev())
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)

				if tc.expectedResult != nil {
					actual, err := keeper.getPendingUpgrade(ctx)
					require.NoError(t, err)

					require.Equal(t, *tc.expectedResult, actual)
				}
			}
		})
	}
}

func TestKeeper_ProcessCancelUpgrade(t *testing.T) {
	tcs := []struct {
		name        string
		setup       func(ctx context.Context, keeper *Keeper) sdk.Context
		setupMock   func(uk *moduletestutil.MockUpgradeKeeper)
		expectedErr string
		postCheck   func(ctx sdk.Context, keeper *Keeper)
	}{
		{
			name: "pass: valid cancel upgrade - before Terence upgrade",
			setupMock: func(uk *moduletestutil.MockUpgradeKeeper) {
				uk.EXPECT().ClearUpgradePlan(gomock.Any()).Return(nil)
			},
		},
		{
			name: "pass: valid cancel upgrade - after Terence upgrade",
			setup: func(ctx context.Context, keeper *Keeper) sdk.Context {
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				sdkCtx = sdkCtx.WithBlockHeight(51)

				return sdkCtx
			},
			postCheck: func(ctx sdk.Context, keeper *Keeper) {
				_, err := keeper.getPendingUpgrade(ctx)
				require.Error(t, err, types.ErrUpgradeNotFound)
			},
		},
		{
			name: "fail: check if Terence upgrade is activated or not - unknown chain ID",
			setup: func(ctx context.Context, keeper *Keeper) sdk.Context {
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				sdkCtx = sdkCtx.WithChainID("unknown-chain-id")

				return sdkCtx
			},
			expectedErr: "failed to check Terence upgrade height",
		},
		{
			name: "fail: clear upgrade plan of upgrade keeper - invalid request",
			setupMock: func(uk *moduletestutil.MockUpgradeKeeper) {
				uk.EXPECT().ClearUpgradePlan(gomock.Any()).Return(sdkerrors.ErrInvalidRequest)
			},
			expectedErr: "failed to cancel the upgrade: invalid_request",
		},
		{
			name: "fail: clear upgrade plan of upgrade keeper - unknown error",
			setupMock: func(uk *moduletestutil.MockUpgradeKeeper) {
				uk.EXPECT().ClearUpgradePlan(gomock.Any()).Return(errors.New("unknown error"))
			},
			expectedErr: "failed to cancel the upgrade: failed to clear upgrade plan: unknown error",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			keeper, ctx, ctrl, uk := setupTestEnvironment(t)
			t.Cleanup(ctrl.Finish)

			if tc.setup != nil {
				ctx = tc.setup(ctx, keeper)
			}

			if tc.setupMock != nil {
				tc.setupMock(uk)
			}

			err := keeper.ProcessCancelUpgrade(ctx, &bindings.UpgradeEntrypointCancelUpgrade{})
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)

				if tc.postCheck != nil {
					tc.postCheck(ctx, keeper)
				}
			}
		})
	}
}

func TestKeeper_ProcessUpgradeEvents(t *testing.T) {
	keeper, ctx, ctrl, uk := setupTestEnvironment(t)
	t.Cleanup(ctrl.Finish)

	upgradeAbi, err := bindings.UpgradeEntrypointMetaData.GetAbi()
	require.NoError(t, err, "failed to load ABI")

	tcs := []struct {
		name        string
		evmEvents   func() []*types.EVMEvent
		setupMock   func()
		expectedErr string
	}{
		{
			name:      "pass: nil events - nothing to process",
			evmEvents: func() []*types.EVMEvent { return nil },
		},
		{
			name:      "pass: empty events - nothing to process",
			evmEvents: func() []*types.EVMEvent { return []*types.EVMEvent{} },
			setupMock: func() {
				uk.EXPECT().DumpUpgradeInfoToDisk(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "pass: one valid upgrade event",
			evmEvents: func() []*types.EVMEvent {
				data, err := upgradeAbi.Events["SoftwareUpgrade"].Inputs.NonIndexed().Pack("test-upgrade", int64(1), "test-info")
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
						Data:    data,
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				uk.EXPECT().ScheduleUpgrade(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "pass: multiple valid upgrade events",
			evmEvents: func() []*types.EVMEvent {
				data1, err := upgradeAbi.Events["SoftwareUpgrade"].Inputs.NonIndexed().Pack("test-upgrade1", int64(2), "test-info")
				require.NoError(t, err)
				data2, err := upgradeAbi.Events["SoftwareUpgrade"].Inputs.NonIndexed().Pack("test-upgrade2", int64(3), "test-info")
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
						Data:    data1,
						TxHash:  dummyHash.Bytes(),
					},
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
						Data:    data2,
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				uk.EXPECT().ScheduleUpgrade(gomock.Any(), gomock.Any()).Return(nil).Times(2)
				uk.EXPECT().DumpUpgradeInfoToDisk(gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
		},
		{
			name: "pass: valid cancel upgrade event",
			evmEvents: func() []*types.EVMEvent {
				data, err := upgradeAbi.Events["CancelUpgrade"].Inputs.NonIndexed().Pack()
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.CancelUpgradeEvent.ID.Bytes()},
						Data:    data,
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				uk.EXPECT().ClearUpgradePlan(gomock.Any()).Return(nil)
			},
		},
		// Failed but pass cases: The following test cases simulate basic error scenarios.
		// Since a mocked upgrade keeper is used, not all error cases can be tested here.
		// Comprehensive error testing would require the real upgrade keeper, which is beyond the scope of this unit test.
		{
			name: "pass(failed but continue): invalid upgrade event - height is 0",
			evmEvents: func() []*types.EVMEvent {
				data, err := upgradeAbi.Events["SoftwareUpgrade"].Inputs.NonIndexed().Pack("test-upgrade", int64(0), "test-info")
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
						Data:    data,
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				uk.EXPECT().ScheduleUpgrade(gomock.Any(), gomock.Any()).Return(errors.New("height must be greater than 0"))
			},
		},
		{
			name: "pass(failed but continue): invalid upgrade event - name is empty",
			evmEvents: func() []*types.EVMEvent {
				data, err := upgradeAbi.Events["SoftwareUpgrade"].Inputs.NonIndexed().Pack("", int64(5), "test-info")
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
						Data:    data,
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				uk.EXPECT().ScheduleUpgrade(gomock.Any(), gomock.Any()).Return(errors.New("name cannot be empty"))
			},
		},
		{
			name: "pass(failed but continue): invalid upgrade event - not an upgrade event, it don't reach ProcessSoftwareUpgrade",
			evmEvents: func() []*types.EVMEvent {
				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes(), dummyHash.Bytes()},
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
		},
		{
			name: "pass(failed but continue): invalid request for scheduling upgrade plan",
			evmEvents: func() []*types.EVMEvent {
				data, err := upgradeAbi.Events["SoftwareUpgrade"].Inputs.NonIndexed().Pack("test-upgrade", int64(1), "test-info")
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
						Data:    data,
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				uk.EXPECT().ScheduleUpgrade(gomock.Any(), gomock.Any()).Return(sdkerrors.ErrInvalidRequest)
			},
		},
		{
			name: "pass(failed but continue): invalid cancel upgrade event - not an upgrade cancel event",
			evmEvents: func() []*types.EVMEvent {
				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.CancelUpgradeEvent.ID.Bytes(), dummyHash.Bytes()},
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
		},
		{
			name: "pass(failed but continue): invalid request for clear upgrade plan",
			evmEvents: func() []*types.EVMEvent {
				data, err := upgradeAbi.Events["CancelUpgrade"].Inputs.NonIndexed().Pack()
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.CancelUpgradeEvent.ID.Bytes()},
						Data:    data,
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				uk.EXPECT().ClearUpgradePlan(gomock.Any()).Return(sdkerrors.ErrInvalidRequest)
			},
		},
		{
			name: "pass(failed but continue): failed to clear upgrade plan",
			evmEvents: func() []*types.EVMEvent {
				data, err := upgradeAbi.Events["CancelUpgrade"].Inputs.NonIndexed().Pack()
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.CancelUpgradeEvent.ID.Bytes()},
						Data:    data,
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				uk.EXPECT().ClearUpgradePlan(gomock.Any()).Return(errors.New("failed to clear upgrade plan"))
			},
		},

		// Fail case: When given EVMEvent is not valid
		{
			name: "fail: invalid EVMEvent - invalid address",
			evmEvents: func() []*types.EVMEvent {
				return []*types.EVMEvent{
					{
						Address: []byte("invalid address"),
						Topics:  [][]byte{types.SoftwareUpgradeEvent.ID.Bytes()},
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			expectedErr: "invalid address length",
		},
		{
			name: "fail: invalid EVMEvent - empty topics",
			evmEvents: func() []*types.EVMEvent {
				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			expectedErr: "empty topics",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMock != nil {
				tc.setupMock()
			}

			cachedCtx, _ := ctx.CacheContext()

			err := keeper.ProcessUpgradeEvents(cachedCtx, 1, tc.evmEvents())
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_ShouldUpgrade(t *testing.T) {
	tcs := []struct {
		name            string
		setup           func(ctx context.Context, keeper *Keeper) sdk.Context
		shouldUpgrade   bool
		expectedUpgrade upgradetypes.Plan
	}{
		{
			name:          "no pending upgrade",
			shouldUpgrade: false,
		},
		{
			name: "upgrade height not reached",
			setup: func(ctx context.Context, keeper *Keeper) sdk.Context {
				sdkCtx := sdk.UnwrapSDKContext(ctx)

				err := keeper.setPendingUpgrade(sdkCtx, upgradetypes.Plan{
					Name:   "test-upgrade",
					Height: 50,
					Info:   "test-info",
				})
				require.NoError(t, err)

				return sdkCtx
			},
			shouldUpgrade: false,
		},
		{
			name: "should upgrade",
			setup: func(ctx context.Context, keeper *Keeper) sdk.Context {
				sdkCtx := sdk.UnwrapSDKContext(ctx)

				err := keeper.setPendingUpgrade(sdkCtx, upgradetypes.Plan{
					Name:   "test-upgrade",
					Height: 51,
					Info:   "test-info",
				})
				require.NoError(t, err)

				sdkCtx = sdkCtx.WithBlockHeight(51)

				return sdkCtx
			},
			shouldUpgrade: true,
			expectedUpgrade: upgradetypes.Plan{
				Name:   "test-upgrade",
				Height: 51,
				Info:   "test-info",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			keeper, ctx, ctrl, _ := setupTestEnvironment(t)
			t.Cleanup(ctrl.Finish)

			if tc.setup != nil {
				ctx = tc.setup(ctx, keeper)
			}

			shouldUpgrade, pendingUpgrade := keeper.ShouldUpgrade(ctx)
			require.Equal(t, tc.shouldUpgrade, shouldUpgrade)

			if tc.shouldUpgrade {
				require.Equal(t, tc.expectedUpgrade, pendingUpgrade)
			}
		})
	}
}

func setupTestEnvironment(t *testing.T) (*Keeper, sdk.Context, *gomock.Controller, *moduletestutil.MockUpgradeKeeper) {
	t.Helper()
	cdc := getCodec(t)
	txConfig := authtx.NewTxConfig(cdc, nil)

	cmtAPI := newMockCometAPI(t, nil)
	header := cmtproto.Header{Height: 1, AppHash: tutil.RandomHash().Bytes(), ProposerAddress: cmtAPI.validatorSet.Validators[0].Address}
	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockClient(ctrl)
	ak := moduletestutil.NewMockAccountKeeper(ctrl)
	esk := moduletestutil.NewMockEvmStakingKeeper(ctrl)
	uk := moduletestutil.NewMockUpgradeKeeper(ctrl)
	dk := moduletestutil.NewMockDistrKeeper(ctrl)

	ctx, storeKey, storeService := setupCtxStore(t, &header)
	mockEngine, err := newMockEngineAPI(storeKey, 0)
	require.NoError(t, err)

	keeper, err := NewKeeper(cdc, storeService, &mockEngine, mockClient, txConfig, ak, esk, uk, dk)
	require.NoError(t, err)
	keeper.SetCometAPI(cmtAPI)
	nxtAddr, err := k1util.PubKeyToAddress(cmtAPI.validatorSet.CopyIncrementProposerPriority(1).Proposer.PubKey)
	require.NoError(t, err)
	keeper.SetValidatorAddress(nxtAddr)
	populateGenesisHead(ctx, t, keeper)

	ctx = ctx.WithChainID(netconf.TestChainID)

	return keeper, ctx, ctrl, uk
}
