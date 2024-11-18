package keeper

import (
	"testing"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	moduletestutil "github.com/piplabs/story/client/x/evmengine/testutil"
	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient/mock"
	"github.com/piplabs/story/lib/k1util"
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
	keeper, ctx, ctrl, uk := setupTestEnvironment(t)
	t.Cleanup(ctrl.Finish)

	tcs := []struct {
		name        string
		ev          func() *bindings.UpgradeEntrypointSoftwareUpgrade
		setupMock   func()
		expectedErr string
	}{
		{
			name: "pass: valid software upgrade event",
			ev: func() *bindings.UpgradeEntrypointSoftwareUpgrade {
				return &bindings.UpgradeEntrypointSoftwareUpgrade{
					Name:   "test-upgrade",
					Height: 1,
					Info:   "test-info",
				}
			},
			setupMock: func() {
				uk.EXPECT().ScheduleUpgrade(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		// Fail cases: The following test cases simulate basic error scenarios.
		// Since a mocked upgrade keeper is used, not all error cases can be tested here.
		// Comprehensive error testing would require the real upgrade keeper, which is beyond the scope of this unit test.
		{
			name: "fail: invalid upgrade event - height is 0",
			ev: func() *bindings.UpgradeEntrypointSoftwareUpgrade {
				return &bindings.UpgradeEntrypointSoftwareUpgrade{
					Name:   "test upgrade",
					Height: 0,
					Info:   "test-info",
				}
			},
			setupMock: func() {
				uk.EXPECT().ScheduleUpgrade(gomock.Any(), gomock.Any()).Return(errors.New("height must be greater than 0"))
			},
			expectedErr: "height must be greater than 0",
		},
		{
			name: "fail: invalid upgrade event - name is empty",
			ev: func() *bindings.UpgradeEntrypointSoftwareUpgrade {
				return &bindings.UpgradeEntrypointSoftwareUpgrade{
					Name:   "",
					Height: 1,
					Info:   "test-info",
				}
			},
			setupMock: func() {
				uk.EXPECT().ScheduleUpgrade(gomock.Any(), gomock.Any()).Return(errors.New("name cannot be empty"))
			},
			expectedErr: "name cannot be empty",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			err := keeper.ProcessSoftwareUpgrade(ctx, tc.ev())
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
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

		// Fail case: When given EVMEvent is not valid
		{
			name: "fail: invalid EVMEvent",
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

	return keeper, ctx, ctrl, uk
}
