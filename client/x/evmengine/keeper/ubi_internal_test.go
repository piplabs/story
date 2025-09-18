package keeper

import (
	"testing"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
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

// const (
//	dummyAddressHex = "0x1398C32A45Bc409b6C652E25bb0a3e702492A4ab"
//)
//
// var (
//	dummyContractAddress = common.HexToAddress(dummyAddressHex)
//	dummyHash            = common.HexToHash("0x1398C32A45Bc409b6C652E25bb0a3e702492A4ab")
//)

func TestKeeper_ProcessUBIPercentageSet(t *testing.T) {
	keeper, ctx, ctrl, dk := setupUBITestEnvironment(t)
	t.Cleanup(ctrl.Finish)

	tcs := []struct {
		name        string
		ev          func() *bindings.UBIPoolUBIPercentageSet
		setupMock   func()
		expectedErr string
	}{
		{
			name: "pass: valid UBI percentage set event",
			ev: func() *bindings.UBIPoolUBIPercentageSet {
				return &bindings.UBIPoolUBIPercentageSet{
					Percentage: 2000,
				}
			},
			setupMock: func() {
				dk.EXPECT().SetUbi(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		// Fail cases: The following test cases simulate basic error scenarios.
		// Since a mocked distribution keeper is used, not all error cases can be tested here.
		// Comprehensive error testing would require the real distribution keeper, which is beyond the scope of this unit test.
		{
			name: "fail: invalid UBI percentage set event - value is 0",
			ev: func() *bindings.UBIPoolUBIPercentageSet {
				return &bindings.UBIPoolUBIPercentageSet{
					Percentage: 11000,
				}
			},
			setupMock: func() {
				dk.EXPECT().SetUbi(gomock.Any(), gomock.Any()).Return(errors.New("ubi too large"))
			},
			expectedErr: "set new UBI percentage",
		},
		{
			name: "fail: invalid UBI set request",
			ev: func() *bindings.UBIPoolUBIPercentageSet {
				return &bindings.UBIPoolUBIPercentageSet{
					Percentage: 2000,
				}
			},
			setupMock: func() {
				dk.EXPECT().SetUbi(gomock.Any(), gomock.Any()).Return(sdkerrors.ErrInvalidRequest)
			},
			expectedErr: "invalid request",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			err := keeper.ProcessUBIPercentageSet(ctx, tc.ev())
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_ProcessUBIEvents(t *testing.T) {
	keeper, ctx, ctrl, dk := setupUBITestEnvironment(t)
	t.Cleanup(ctrl.Finish)

	ubiAbi, err := bindings.UBIPoolMetaData.GetAbi()
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
			name: "pass: one valid UBI percentage set event",
			evmEvents: func() []*types.EVMEvent {
				data, err := ubiAbi.Events["UBIPercentageSet"].Inputs.NonIndexed().Pack(uint32(2000))
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.UBIPercentageSetEvent.ID.Bytes()},
						Data:    data,
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dk.EXPECT().SetUbi(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		// Failed but pass cases: The following test cases simulate basic error scenarios.
		// Since a mocked distribution keeper is used, not all error cases can be tested here.
		// Comprehensive error testing would require the real distribution keeper, which is beyond the scope of this unit test.
		{
			name: "pass(failed but continue): invalid UBI percentage set event - too large value",
			evmEvents: func() []*types.EVMEvent {
				data, err := ubiAbi.Events["UBIPercentageSet"].Inputs.NonIndexed().Pack(uint32(11000))
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.UBIPercentageSetEvent.ID.Bytes()},
						Data:    data,
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dk.EXPECT().SetUbi(gomock.Any(), gomock.Any()).Return(errors.New("ubi too large"))
			},
		},
		{
			name: "pass(failed but continue): invalid UBI percentage set event",
			evmEvents: func() []*types.EVMEvent {
				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{types.UBIPercentageSetEvent.ID.Bytes(), dummyHash.Bytes()},
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
		},

		// Fail case: When given EVMEvent is not valid
		{
			name: "fail: invalid EVMEvent - invalid address",
			evmEvents: func() []*types.EVMEvent {
				return []*types.EVMEvent{
					{
						Address: []byte("invalid address"),
						Topics:  [][]byte{types.UBIPercentageSetEvent.ID.Bytes()},
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

			ethLogs := make([]*ethtypes.Log, 0, len(tc.evmEvents()))
			for _, evmEvent := range tc.evmEvents() {
				ethLog, err := evmEvent.ToEthLog()
				require.NoError(t, err)
				ethLogs = append(ethLogs, &ethLog)
			}

			err := keeper.ProcessUBIEvents(cachedCtx, 1, ethLogs)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func setupUBITestEnvironment(t *testing.T) (*Keeper, sdk.Context, *gomock.Controller, *moduletestutil.MockDistrKeeper) {
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

	return keeper, ctx, ctrl, dk
}
