package keeper

import (
	"encoding/hex"
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
	testDkgPubKey := []byte("test-dkg-pubkey")
	testCommPubKey := []byte("test-comm-pubkey")
	testRawQuote := []byte("test-raw-quote")
	testSignature := []byte("test-signature")
	testGlobalPubKey := []byte("test-global-pubkey")
	testPublicCoeffs := []byte("test-public-coeffs")
	testIndex := uint32(1)
	testComplainIndexes := []uint32{1, 2, 3}
	testRecipientIndex := uint32(2)
	testActivationHeight := uint32(100)

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
				data, err := dkgAbi.Events["DKGInitialized"].Inputs.NonIndexed().Pack(
					testCodeCommitment, testRound, testDkgPubKey, testCommPubKey, testRawQuote)
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGInitializedEvent.ID.Bytes(),
							common.LeftPadBytes(testValidator.Bytes(), 32), // indexed msgSender
						},
						Data:   data,
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dkgk.EXPECT().RegistrationInitialized(gomock.Any(), testValidator, testCodeCommitment, testRound, testDkgPubKey, testCommPubKey, testRawQuote).Return(nil)
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
				data, err := dkgAbi.Events["DKGFinalized"].Inputs.NonIndexed().Pack(
					testRound, testCodeCommitment, testParticipantsRoot, testGlobalPubKey, testPublicCoeffs, testSignature)
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
				dkgk.EXPECT().Finalized(gomock.Any(), testRound, testValidator, testCodeCommitment, testParticipantsRoot, testSignature, testGlobalPubKey, testPublicCoeffs).Return(nil)
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
					testActivationHeight, testCodeCommitment)
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
				dkgk.EXPECT().UpgradeScheduled(gomock.Any(), testActivationHeight, testCodeCommitment).Return(nil)
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
						require.Equal(t, hex.EncodeToString(testCodeCommitment[:]), attrs[2].Value)

						break
					}
				}
				require.True(t, found, "Expected DKGUpgradeScheduledSuccess event to be emitted for %s", testName)
			},
		},
		{
			name: "pass: RemoteAttestationProcessedOnChain event",
			evmEvents: func() []*types.EVMEvent {
				data, err := dkgAbi.Events["RemoteAttestationProcessedOnChain"].Inputs.NonIndexed().Pack(
					testValidator, uint8(1), testRound, testCodeCommitment) // ChallengeStatus.Invalidated = 1
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGRemoteAttestationProcessedOnChainEvent.ID.Bytes(),
						},
						Data:   data,
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dkgk.EXPECT().RemoteAttestationProcessedOnChain(gomock.Any(), testValidator, int(1), testRound, testCodeCommitment).Return(nil)
			},
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				// Should emit DKGRemoteAttestationProcessedOnChainSuccess event
				t.Helper()
				found := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGRemoteAttestationProcessedOnChainSuccess {
						found = true
						// Check attributes
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, testValidator.Hex(), attrs[1].Value)
						require.Equal(t, "1", attrs[2].Value) // ChalStatus = 1
						require.Equal(t, strconv.FormatUint(uint64(testRound), 10), attrs[3].Value)
						require.Equal(t, hex.EncodeToString(testCodeCommitment[:]), attrs[4].Value)

						break
					}
				}
				require.True(t, found, "Expected DKGRemoteAttestationProcessedOnChainSuccess event to be emitted for %s", testName)
			},
		},
		{
			name: "pass: DealComplaintsSubmitted event",
			evmEvents: func() []*types.EVMEvent {
				data, err := dkgAbi.Events["DealComplaintsSubmitted"].Inputs.NonIndexed().Pack(
					testIndex, testComplainIndexes, testRound, testCodeCommitment)
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGDealComplaintsSubmittedEvent.ID.Bytes(),
						},
						Data:   data,
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dkgk.EXPECT().DealComplaintsSubmitted(gomock.Any(), testIndex, testComplainIndexes, testRound, testCodeCommitment).Return(nil)
			},
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				// Should emit DKGDealComplaintsSubmittedSuccess event
				t.Helper()
				found := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGDealComplaintsSubmittedSuccess {
						found = true
						// Check attributes
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, strconv.FormatUint(uint64(testIndex), 10), attrs[1].Value)
						require.Equal(t, "1,2,3", attrs[2].Value) // ComplainIndexes joined
						require.Equal(t, strconv.FormatUint(uint64(testRound), 10), attrs[3].Value)
						require.Equal(t, hex.EncodeToString(testCodeCommitment[:]), attrs[4].Value)

						break
					}
				}
				require.True(t, found, "Expected DKGDealComplaintsSubmittedSuccess event to be emitted for %s", testName)
			},
		},
		{
			name: "pass: DealVerified event",
			evmEvents: func() []*types.EVMEvent {
				data, err := dkgAbi.Events["DealVerified"].Inputs.NonIndexed().Pack(
					testIndex, testRecipientIndex, testRound, testCodeCommitment)
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGDealVerifiedEvent.ID.Bytes(),
						},
						Data:   data,
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dkgk.EXPECT().DealVerified(gomock.Any(), testIndex, testRecipientIndex, testRound, testCodeCommitment).Return(nil)
			},
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				// Should emit DKGDealVerifiedSuccess event
				t.Helper()
				found := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGDealVerifiedSuccess {
						found = true
						// Check attributes
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, strconv.FormatUint(uint64(testIndex), 10), attrs[1].Value)
						require.Equal(t, strconv.FormatUint(uint64(testRecipientIndex), 10), attrs[2].Value)
						require.Equal(t, strconv.FormatUint(uint64(testRound), 10), attrs[3].Value)
						require.Equal(t, hex.EncodeToString(testCodeCommitment[:]), attrs[4].Value)

						break
					}
				}
				require.True(t, found, "Expected DKGDealVerifiedSuccess event to be emitted for %s", testName)
			},
		},
		{
			name: "pass: InvalidDeal event",
			evmEvents: func() []*types.EVMEvent {
				data, err := dkgAbi.Events["InvalidDeal"].Inputs.NonIndexed().Pack(
					testIndex, testRound, testCodeCommitment)
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGInvalidDealEvent.ID.Bytes(),
						},
						Data:   data,
						TxHash: dummyHash.Bytes(),
					},
				}
			},
			setupMock: func() {
				dkgk.EXPECT().InvalidDeal(gomock.Any(), testIndex, testRound, testCodeCommitment).Return(nil)
			},
			verifyEvents: func(t *testing.T, events sdk.Events, testName string) {
				// Should emit DKGInvalidDealSuccess event
				t.Helper()
				found := false
				for _, event := range events {
					if event.Type == types.EventTypeDKGInvalidDealSuccess {
						found = true
						// Check attributes
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, strconv.FormatUint(uint64(testIndex), 10), attrs[1].Value)
						require.Equal(t, strconv.FormatUint(uint64(testRound), 10), attrs[2].Value)
						require.Equal(t, hex.EncodeToString(testCodeCommitment[:]), attrs[3].Value)

						break
					}
				}
				require.True(t, found, "Expected DKGInvalidDealSuccess event to be emitted for %s", testName)
			},
		},
		{
			name: "pass: multiple DKG events",
			evmEvents: func() []*types.EVMEvent {
				// DKGInitialized event
				initData, err := dkgAbi.Events["DKGInitialized"].Inputs.NonIndexed().Pack(
					testCodeCommitment, testRound, testDkgPubKey, testCommPubKey, testRawQuote)
				require.NoError(t, err)

				// DKGFinalized event
				finalizedData, err := dkgAbi.Events["DKGFinalized"].Inputs.NonIndexed().Pack(
					testRound, testCodeCommitment, testParticipantsRoot, testGlobalPubKey, testPublicCoeffs, testSignature)
				require.NoError(t, err)

				return []*types.EVMEvent{
					{
						Address: dummyContractAddress.Bytes(),
						Topics: [][]byte{
							types.DKGInitializedEvent.ID.Bytes(),
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
				dkgk.EXPECT().RegistrationInitialized(gomock.Any(), testValidator, testCodeCommitment, testRound, testDkgPubKey, testCommPubKey, testRawQuote).Return(nil)
				dkgk.EXPECT().Finalized(gomock.Any(), testRound, testValidator, testCodeCommitment, testParticipantsRoot, testSignature, testGlobalPubKey, testPublicCoeffs).Return(nil)
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
							types.DKGInitializedEvent.ID.Bytes(),
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
	testDkgPubKey := []byte("test-dkg-pubkey")
	testCommPubKey := []byte("test-comm-pubkey")
	testRawQuote := []byte("test-raw-quote")

	dkgContract := &bindings.DKG{}
	keeper.dkgContract = dkgContract

	// Create mock log with proper structure
	mockLog := &ethtypes.Log{
		Address: dummyContractAddress,
		Topics: []common.Hash{
			types.DKGInitializedEvent.ID,
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
	data, err := dkgAbi.Events["DKGInitialized"].Inputs.NonIndexed().Pack(
		testCodeCommitment, testRound, testDkgPubKey, testCommPubKey, testRawQuote)
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
				dkgk.EXPECT().RegistrationInitialized(gomock.Any(), testValidator, testCodeCommitment, testRound, testDkgPubKey, testCommPubKey, testRawQuote).Return(nil)
			},
		},
		{
			name: "fail: DKG keeper returns error",
			setupMock: func() {
				dkgk.EXPECT().RegistrationInitialized(gomock.Any(), testValidator, testCodeCommitment, testRound, testDkgPubKey, testCommPubKey, testRawQuote).Return(
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

			err := keeper.ProcessDKGInitialized(ctx, mockLog)
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
						// Check attributes
						attrs := event.Attributes
						require.NotEmpty(t, attrs)
						require.Equal(t, strconv.FormatUint(uint64(testRound), 10), attrs[1].Value)
						require.Equal(t, testValidator.Hex(), attrs[2].Value)
						require.Equal(t, hex.EncodeToString(testCodeCommitment[:]), attrs[3].Value)

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
						require.Equal(t, hex.EncodeToString(testDkgPubKey), attrs[4].Value)
						require.Equal(t, hex.EncodeToString(testCommPubKey), attrs[5].Value)
						require.Equal(t, hex.EncodeToString(testRawQuote), attrs[6].Value)

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
