package keeper

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"
	"testing"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	moduletestutil "github.com/piplabs/story/client/x/evmengine/testutil"
	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/ethclient/mock"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/tutil"

	"go.uber.org/mock/gomock"
)

func TestKeeper_ProcessDKGEvents(t *testing.T) {
	keeper, ctx, ctrl, dkgk := setupDKGTestEnvironment(t)
	t.Cleanup(ctrl.Finish)

	dkgAbi, err := bindings.DKGMetaData.GetAbi()
	require.NoError(t, err, "failed to load DKG ABI")

	testValidator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testCodeCommitment := [32]byte{0x12, 0x34}
	testParticipantsRoot := [32]byte{0x56, 0x78}
	testRound := uint32(1)
	testStartBlockHeight := big.NewInt(100)
	testStartBlockHash := [32]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB}
	testDkgPubKey := []byte("test-dkg-pubkey")
	testCommPubKey := []byte("test-comm-pubkey")
	testEnclaveReport := []byte("test-enclave-report")
	testSignature := []byte("test-signature")
	testGlobalPubKey := []byte("test-global-pubkey")
	testPublicCoeffs := [][]byte{[]byte("test-public-coeffs")}
	testPubKeyShare := []byte("test-pubkey-share")
	testEnclaveType := [32]byte{0x01}
	testValidationContext := []byte("")
	testActivationHeight := int64(100)

	tcs := []struct {
		name         string
		evmEvents    func() []*types.EVMEvent
		setupMock    func()
		expectedErr  string
		verifyEvents func(t *testing.T, events sdk.Events, testName string)
	}{
		{
			name:      "pass: nil events - nothing to process",
			evmEvents: func() []*types.EVMEvent { return nil },
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				// Should not emit any DKG events when no events to process
				t.Helper()
				for _, event := range events {
					require.NotContains(t, event.Type, "dkg_", "No DKG events should be emitted for %s", testName)
				}
			},
		},
		{
			name:      "pass: empty events - nothing to process",
			evmEvents: func() []*types.EVMEvent { return []*types.EVMEvent{} },
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				// Should not emit any DKG events when no events to process
				t.Helper()
				for _, event := range events {
					require.NotContains(t, event.Type, "dkg_", "No DKG events should be emitted for %s", testName)
				}
			},
		},
		{
			name: "pass: DKGInitialized event",
			evmEvents: func() []*types.EVMEvent {
				data, err := dkgAbi.Events["Registered"].Inputs.NonIndexed().Pack(
					testEnclaveReport, testRound, testEnclaveType, testCommPubKey, testDkgPubKey, testCodeCommitment, testStartBlockHeight, testStartBlockHash, testValidationContext)
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGRegisteredEvent.ID.Bytes(),
							common.LeftPadBytes(testValidator.Bytes(), 32), // indexed msgSender
						},
						Data:   data,
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dkgk.EXPECT().Registered(gomock.Any(), testValidator, testCodeCommitment, testRound, testStartBlockHeight, testStartBlockHash, testEnclaveType, testDkgPubKey, testCommPubKey, testEnclaveReport).Return(nil)
			},
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				// Should emit DKGInitializedSuccess event
				t.Helper()
				found := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGInitializedSuccess {
						found = true
						// Check attributes
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, strconv.FormatUint(uint64(testRound), 10), attrs[1].Value)
						require.Equal(t, testValidator.Hex(), attrs[2].Value)
						require.Equal(t, hex.EncodeToString(testCodeCommitment[:]), attrs[3].Value)

						break
					}
				}
				require.True(t, found, "Expected DKGInitializedSuccess event to be emitted for %s", testName)
			},
		},
		{
			name: "pass: DKGFinalized event",
			evmEvents: func() []*types.EVMEvent {
				data, err := dkgAbi.Events["Finalized"].Inputs.NonIndexed().Pack(
					testRound, testEnclaveType, testCodeCommitment, testParticipantsRoot, testGlobalPubKey, testPublicCoeffs, testPubKeyShare, testSignature)
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGFinalizedEvent.ID.Bytes(),
							common.LeftPadBytes(testValidator.Bytes(), 32), // indexed msgSender
						},
						Data:   data,
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dkgk.EXPECT().Finalized(gomock.Any(), testRound, testValidator, testCodeCommitment, testParticipantsRoot, testSignature, testGlobalPubKey, testPublicCoeffs, testPubKeyShare).Return(nil)
			},
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				// Should emit DKGFinalizedSuccess event
				t.Helper()
				found := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGFinalizedSuccess {
						found = true
						// Check attributes
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, strconv.FormatUint(uint64(testRound), 10), attrs[1].Value)
						require.Equal(t, testValidator.Hex(), attrs[2].Value)
						require.Equal(t, hex.EncodeToString(testCodeCommitment[:]), attrs[3].Value)

						break
					}
				}
				require.True(t, found, "Expected DKGFinalizedSuccess event to be emitted for %s", testName)
			},
		},
		{
			name: "pass: DKG UpgradeScheduled event",
			evmEvents: func() []*types.EVMEvent {
				data, err := dkgAbi.Events["UpgradeScheduled"].Inputs.NonIndexed().Pack(
					big.NewInt(int64(testActivationHeight)), "v1.0.0")
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGUpgradeScheduledEvent.ID.Bytes(),
						},
						Data:   data,
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dkgk.EXPECT().UpgradeScheduled(gomock.Any(), testActivationHeight, "v1.0.0").Return(nil)
			},
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				// Should emit DKGUpgradeScheduledSuccess event
				t.Helper()
				found := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGUpgradeScheduledSuccess {
						found = true
						// Check attributes
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, strconv.FormatUint(uint64(testActivationHeight), 10), attrs[1].Value)
						require.Equal(t, "v1.0.0", attrs[2].Value)

						break
					}
				}
				require.True(t, found, "Expected DKGUpgradeScheduledSuccess event to be emitted for %s", testName)
			},
		},
		{
			name: "pass: multiple DKG events",
			evmEvents: func() []*types.EVMEvent {
				// Registered event
				initData, err := dkgAbi.Events["Registered"].Inputs.NonIndexed().Pack(
					testEnclaveReport, testRound, testEnclaveType, testCommPubKey, testDkgPubKey, testCodeCommitment, testStartBlockHeight, testStartBlockHash, testValidationContext)
				require.NoError(t, err)

				// Finalized event
				finalizedData, err := dkgAbi.Events["Finalized"].Inputs.NonIndexed().Pack(
					testRound, testEnclaveType, testCodeCommitment, testParticipantsRoot, testGlobalPubKey, testPublicCoeffs, testPubKeyShare, testSignature)
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGRegisteredEvent.ID.Bytes(),
							common.LeftPadBytes(testValidator.Bytes(), 32),
						},
						Data:   initData,
						TxHash: dummyHash.Bytes(),
					},
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGFinalizedEvent.ID.Bytes(),
							common.LeftPadBytes(testValidator.Bytes(), 32),
						},
						Data:   finalizedData,
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dkgk.EXPECT().Registered(gomock.Any(), testValidator, testCodeCommitment, testRound, testStartBlockHeight, testStartBlockHash, testEnclaveType, testDkgPubKey, testCommPubKey, testEnclaveReport).Return(nil)
				dkgk.EXPECT().Finalized(gomock.Any(), testRound, testValidator, testCodeCommitment, testParticipantsRoot, testSignature, testGlobalPubKey, testPublicCoeffs, testPubKeyShare).Return(nil)
			},
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				// Should emit both DKGInitializedSuccess and DKGFinalizedSuccess events
				t.Helper()
				foundInit := false
				foundFinalized := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGInitializedSuccess {
						foundInit = true
					}
					if event.Type == types.EventTypeDKGFinalizedSuccess {
						foundFinalized = true
					}
				}
				require.True(t, foundInit, "Expected DKGInitializedSuccess event to be emitted for %s", testName)
				require.True(t, foundFinalized, "Expected DKGFinalizedSuccess event to be emitted for %s", testName)
			},
		},
		{
			name: "pass: MinReqRegisteredParticipantsSet event",
			evmEvents: func() []*types.EVMEvent {
				data, err := dkgAbi.Events["MinReqRegisteredParticipantsSet"].Inputs.NonIndexed().Pack(
					big.NewInt(5))
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGMinReqRegisteredParticipantsSetEvent.ID.Bytes(),
						},
						Data:   data,
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dkgk.EXPECT().SetMinReqRegisteredParticipants(gomock.Any(), uint32(5)).Return(nil)
			},
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				t.Helper()
				found := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGMinReqRegisteredParticipantsSetSuccess {
						found = true
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, strconv.FormatUint(5, 10), attrs[1].Value)

						break
					}
				}
				require.True(t, found, "Expected DKGMinReqRegisteredParticipantsSetSuccess event to be emitted for %s", testName)
			},
		},
		{
			name: "pass: MinReqFinalizedParticipantsSet event",
			evmEvents: func() []*types.EVMEvent {
				data, err := dkgAbi.Events["MinReqFinalizedParticipantsSet"].Inputs.NonIndexed().Pack(
					big.NewInt(4))
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGMinReqFinalizedParticipantsSetEvent.ID.Bytes(),
						},
						Data:   data,
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dkgk.EXPECT().SetMinReqFinalizedParticipants(gomock.Any(), uint32(4)).Return(nil)
			},
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				t.Helper()
				found := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGMinReqFinalizedParticipantsSetSuccess {
						found = true
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, strconv.FormatUint(4, 10), attrs[1].Value)

						break
					}
				}
				require.True(t, found, "Expected DKGMinReqFinalizedParticipantsSetSuccess event to be emitted for %s", testName)
			},
		},
		{
			name: "pass: OperationalThresholdSet event",
			evmEvents: func() []*types.EVMEvent {
				data, err := dkgAbi.Events["OperationalThresholdSet"].Inputs.NonIndexed().Pack(
					big.NewInt(750))
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGOperationalThresholdSetEvent.ID.Bytes(),
						},
						Data:   data,
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dkgk.EXPECT().SetOperationalThreshold(gomock.Any(), uint32(750)).Return(nil)
			},
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				t.Helper()
				found := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGOperationalThresholdSetSuccess {
						found = true
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, strconv.FormatUint(750, 10), attrs[1].Value)

						break
					}
				}
				require.True(t, found, "Expected DKGOperationalThresholdSetSuccess event to be emitted for %s", testName)
			},
		},
		{
			name: "fail: invalid log data - unrecognized topic",
			evmEvents: func() []*types.EVMEvent {
				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics:  [][]byte{common.Hash{0x99}.Bytes()}, // Unrecognized topic hash
						Data:    []byte("invalid"),
						TxHash:  dummyHash.Bytes(),
					},
				}
			},
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				// Should not emit any DKG events when there's an unrecognized event
				t.Helper()
				for _, event := range events {
					require.NotContains(t, event.Type, "dkg_", "No DKG events should be emitted for unrecognized event %s", testName)
				}
			},
		},
		{
			name: "pass(failed but continue): DKGInitialized with invalid data - parsing error",
			evmEvents: func() []*types.EVMEvent {
				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGRegisteredEvent.ID.Bytes(),
							common.LeftPadBytes(testValidator.Bytes(), 32),
						},
						Data:   []byte("invalid-data"), // Invalid data that will fail parsing
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			// No setupMock since the parsing should fail before reaching the DKG keeper
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				t.Helper()
				// Since parsing fails, ProcessDKGEvents continues but doesn't emit success events
				// It should not emit any DKG success events, but may emit debug/error events
				successEventFound := false
				for _, event := range events {
					if strings.Contains(event.Type, "dkg_") && strings.Contains(event.Type, "_success") {
						successEventFound = true

						break
					}
				}
				require.False(t, successEventFound, "Should not emit DKG success events when parsing fails for %s", testName)
			},
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

			err := keeper.ProcessDKGEvents(cachedCtx, 1, ethLogs)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			// Check if the correct events were emitted
			events := cachedCtx.EventManager().Events()
			tc.verifyEvents(t, events, tc.name)
		})
	}
}

func TestKeeper_ProcessDKGInitialized(t *testing.T) {
	keeper, ctx, ctrl, dkgk := setupDKGTestEnvironment(t)
	t.Cleanup(ctrl.Finish)

	testValidator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testCodeCommitment := [32]byte{0x12, 0x34}
	testRound := uint32(1)
	testStartBlockHeight := big.NewInt(100)
	testStartBlockHash := [32]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB}
	testDkgPubKey := []byte("test-dkg-pubkey")
	testCommPubKey := []byte("test-comm-pubkey")
	testEnclaveReport := []byte("test-enclave-report")
	testEnclaveType := [32]byte{0x01}
	testValidationContext := []byte("")

	// Create mock log with proper structure
	mockLog := &ethtypes.Log{
		Address: dummyContractAddress,
		Topics: []common.Hash{
			types.DKGRegisteredEvent.ID,
			common.BytesToHash(common.LeftPadBytes(testValidator.Bytes(), 32)),
		},
		Data:        []byte{}, // Will be filled by ABI packing
		TxHash:      dummyHash,
		BlockNumber: 1,
		Index:       0,
	}

	// Pack the non-indexed data using the ABI
	dkgAbi, err := bindings.DKGMetaData.GetAbi()
	require.NoError(t, err)
	data, err := dkgAbi.Events["Registered"].Inputs.NonIndexed().Pack(
		testEnclaveReport, testRound, testEnclaveType, testCommPubKey, testDkgPubKey, testCodeCommitment, testStartBlockHeight, testStartBlockHash, testValidationContext)
	require.NoError(t, err)
	mockLog.Data = data

	tcs := []struct {
		name        string
		setupMock   func()
		expectedErr string
	}{
		{
			name: "pass: successful DKG initialization",
			setupMock: func() {
				dkgk.EXPECT().Registered(gomock.Any(), testValidator, testCodeCommitment, testRound, testStartBlockHeight, testStartBlockHash, testEnclaveType, testDkgPubKey, testCommPubKey, testEnclaveReport).Return(nil)
			},
		},
		{
			name: "fail: DKG keeper returns error",
			setupMock: func() {
				dkgk.EXPECT().Registered(gomock.Any(), testValidator, testCodeCommitment, testRound, testStartBlockHeight, testStartBlockHash, testEnclaveType, testDkgPubKey, testCommPubKey, testEnclaveReport).Return(
					sdkerrors.ErrInvalidRequest.Wrap("invalid request"))
			},
			expectedErr: "invalid request",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMock != nil {
				tc.setupMock()
			}

			err := keeper.ProcessDKGRegistered(ctx, mockLog)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			// Check if the correct events were emitted
			events := ctx.EventManager().Events()
			if tc.expectedErr != "" {
				// Should emit failure event
				found := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGInitializedFailure {
						found = true
						// Check attributes (ErrorCode is attrs[0], so other attrs shift by 1)
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, strconv.FormatUint(uint64(testRound), 10), attrs[2].Value)
						require.Equal(t, testValidator.Hex(), attrs[3].Value)
						require.Equal(t, hex.EncodeToString(testCodeCommitment[:]), attrs[4].Value)

						break
					}
				}
				require.True(t, found, "Expected failure event to be emitted")
			} else {
				// Should emit success event
				found := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGInitializedSuccess {
						found = true
						// Check attributes
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, strconv.FormatUint(uint64(testRound), 10), attrs[1].Value)
						require.Equal(t, testValidator.Hex(), attrs[2].Value)
						require.Equal(t, hex.EncodeToString(testCodeCommitment[:]), attrs[3].Value)
						require.Equal(t, testStartBlockHeight.String(), attrs[4].Value)
						require.Equal(t, hex.EncodeToString(testStartBlockHash[:]), attrs[5].Value)
						require.Equal(t, hex.EncodeToString(testDkgPubKey), attrs[6].Value)
						require.Equal(t, hex.EncodeToString(testCommPubKey), attrs[7].Value)
						require.Equal(t, hex.EncodeToString(testValidationContext), attrs[8].Value)

						break
					}
				}
				require.True(t, found, "Expected success event to be emitted")
			}
		})
	}
}

func setupDKGTestEnvironment(t *testing.T) (*Keeper, sdk.Context, *gomock.Controller, *moduletestutil.MockDKGKeeper) {
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
	dkgk := moduletestutil.NewMockDKGKeeper(ctrl)

	ctx, storeKey, storeService := setupCtxStore(t, &header)
	mockEngine, err := newMockEngineAPI(storeKey, 0)
	require.NoError(t, err)

	keeper, err := NewKeeper(cdc, storeService, &mockEngine, mockClient, txConfig, ak, esk, uk, dk, dkgk)
	require.NoError(t, err)
	keeper.SetCometAPI(cmtAPI)
	nxtAddr, err := k1util.PubKeyToAddress(cmtAPI.validatorSet.CopyIncrementProposerPriority(1).Proposer.PubKey)
	require.NoError(t, err)
	keeper.SetValidatorAddress(nxtAddr)
	populateGenesisHead(ctx, t, keeper)

	return keeper, ctx, ctrl, dkgk
}
