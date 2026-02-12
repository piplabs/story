package keeper

import (
	"context"
	"testing"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"
	"go.uber.org/mock/gomock"
)

func TestKeeper_RegistrationInitialized(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testValidator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testCodeCommitment := [32]byte{0x12, 0x34, 0x56, 0x78}
	testRound := uint32(1)
	testStartBlockHeight := uint64(100)
	testStartBlockHash := [32]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB}
	testDkgPubKey := []byte("test-dkg-pubkey")
	testCommPubKey := []byte("test-comm-pubkey")
	testRawQuote := []byte("test-raw-quote")

	validDKGNetwork := &types.DKGNetwork{
		CodeCommitment:   testCodeCommitment[:],
		Round:            testRound,
		StartBlockHeight: int64(testStartBlockHeight),
		StartBlockHash:   testStartBlockHash[:],
		ActiveValSet:     []string{testValidator.Hex()},
		Total:            5,
		Threshold:        3,
		Stage:            types.DKGStageRegistration,
	}
	require.NoError(t, k.setDKGNetwork(ctx, validDKGNetwork))

	tcs := []struct {
		name             string
		msgSender        common.Address
		codeCommitment   [32]byte
		round            uint32
		startBlockHeight uint64
		startBlockHash   [32]byte
		dkgPubKey        []byte
		commPubKey       []byte
		rawQuote         []byte
		setupNetwork     func()
		expectedErr      string
		expectedRegData  *types.DKGRegistration
	}{
		{
			name:             "pass: successful registration initialization",
			msgSender:        testValidator,
			codeCommitment:   testCodeCommitment,
			round:            testRound,
			startBlockHeight: testStartBlockHeight,
			startBlockHash:   testStartBlockHash,
			dkgPubKey:        testDkgPubKey,
			commPubKey:       testCommPubKey,
			rawQuote:         testRawQuote,
			setupNetwork: func() {
				// Network already set up in test setup
			},
			expectedRegData: &types.DKGRegistration{
				Round:         testRound,
				ValidatorAddr: testValidator.Hex(),
				Index:         1,
				DkgPubKey:     testDkgPubKey,
				CommPubKey:    testCommPubKey,
				RawQuote:      testRawQuote,
				Status:        types.DKGRegStatusVerified,
			},
		},
		{
			name:             "fail: DKG network not found",
			msgSender:        testValidator,
			codeCommitment:   [32]byte{0x99, 0x99, 0x99, 0x99},
			round:            testRound,
			startBlockHeight: testStartBlockHeight,
			startBlockHash:   testStartBlockHash,
			dkgPubKey:        testDkgPubKey,
			commPubKey:       testCommPubKey,
			rawQuote:         testRawQuote,
			expectedErr:      "dkg network not found",
		},
		{
			name:             "fail: start block height mismatch",
			msgSender:        testValidator,
			codeCommitment:   testCodeCommitment,
			round:            testRound,
			startBlockHeight: 999, // Wrong height
			startBlockHash:   testStartBlockHash,
			dkgPubKey:        testDkgPubKey,
			commPubKey:       testCommPubKey,
			rawQuote:         testRawQuote,
			setupNetwork: func() {
				// Network already set up with height=100
			},
			expectedErr: "start block height mismatch",
		},
		{
			name:             "fail: start block hash mismatch",
			msgSender:        testValidator,
			codeCommitment:   testCodeCommitment,
			round:            testRound,
			startBlockHeight: testStartBlockHeight,
			startBlockHash:   [32]byte{0xFF, 0xFF, 0xFF, 0xFF}, // Wrong hash
			dkgPubKey:        testDkgPubKey,
			commPubKey:       testCommPubKey,
			rawQuote:         testRawQuote,
			setupNetwork: func() {
				// Network already set up with different hash
			},
			expectedErr: "start block hash mismatch",
		},
		{
			name:             "fail: round not in registration stage",
			msgSender:        testValidator,
			codeCommitment:   testCodeCommitment,
			round:            testRound,
			startBlockHeight: testStartBlockHeight,
			startBlockHash:   testStartBlockHash,
			dkgPubKey:        testDkgPubKey,
			commPubKey:       testCommPubKey,
			rawQuote:         testRawQuote,
			setupNetwork: func() {
				networkWithDifferentStage := &types.DKGNetwork{
					CodeCommitment:   testCodeCommitment[:],
					Round:            testRound,
					StartBlockHeight: int64(testStartBlockHeight),
					StartBlockHash:   testStartBlockHash[:],
					ActiveValSet:     []string{testValidator.Hex()},
					Total:            5,
					Threshold:        3,
					Stage:            types.DKGStageDealing,
				}
				require.NoError(t, k.setDKGNetwork(ctx, networkWithDifferentStage))
			},
			expectedErr: "round is not in registration stage",
		},
		{
			name:             "fail: validator not in active set",
			msgSender:        common.HexToAddress("0x9999999999999999999999999999999999999999"),
			codeCommitment:   testCodeCommitment,
			round:            testRound,
			startBlockHeight: testStartBlockHeight,
			startBlockHash:   testStartBlockHash,
			dkgPubKey:        testDkgPubKey,
			commPubKey:       testCommPubKey,
			rawQuote:         testRawQuote,
			setupNetwork: func() {
				networkInRegistrationStage := &types.DKGNetwork{
					CodeCommitment:   testCodeCommitment[:],
					Round:            testRound,
					StartBlockHeight: int64(testStartBlockHeight),
					StartBlockHash:   testStartBlockHash[:],
					ActiveValSet:     []string{testValidator.Hex()},
					Total:            5,
					Threshold:        3,
					Stage:            types.DKGStageRegistration,
				}
				require.NoError(t, k.setDKGNetwork(ctx, networkInRegistrationStage))
			},
			expectedErr: "msg sender is not in the active validator set",
		},
		{
			name:             "pass: second registration gets incremented index",
			msgSender:        testValidator,
			codeCommitment:   testCodeCommitment,
			round:            testRound,
			startBlockHeight: testStartBlockHeight,
			startBlockHash:   testStartBlockHash,
			dkgPubKey:        []byte("second-dkg-pubkey"),
			commPubKey:       []byte("second-comm-pubkey"),
			rawQuote:         []byte("second-raw-quote"),
			setupNetwork: func() {
				anotherValidator := common.HexToAddress("0xAABBCCDDEEFF112233445566778899AABBCCDDEE")
				networkWithMultipleValidators := &types.DKGNetwork{
					CodeCommitment:   testCodeCommitment[:],
					Round:            testRound,
					StartBlockHeight: int64(testStartBlockHeight),
					StartBlockHash:   testStartBlockHash[:],
					ActiveValSet:     []string{testValidator.Hex(), anotherValidator.Hex()},
					Total:            5,
					Threshold:        3,
					Stage:            types.DKGStageRegistration,
				}
				require.NoError(t, k.setDKGNetwork(ctx, networkWithMultipleValidators))

				firstReg := &types.DKGRegistration{
					Round:         testRound,
					ValidatorAddr: anotherValidator.Hex(),
					Index:         1,
					DkgPubKey:     []byte("first-dkg-pubkey"),
					CommPubKey:    []byte("first-comm-pubkey"),
					RawQuote:      []byte("first-raw-quote"),
					Status:        types.DKGRegStatusVerified,
				}
				require.NoError(t, k.setDKGRegistration(ctx, testCodeCommitment, anotherValidator, firstReg))
			},
			expectedRegData: &types.DKGRegistration{
				Round:         testRound,
				ValidatorAddr: testValidator.Hex(),
				Index:         2,
				DkgPubKey:     []byte("second-dkg-pubkey"),
				CommPubKey:    []byte("second-comm-pubkey"),
				RawQuote:      []byte("second-raw-quote"),
				Status:        types.DKGRegStatusVerified,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupNetwork != nil {
				tc.setupNetwork()
			}

			err := k.RegistrationInitialized(ctx, tc.msgSender, tc.codeCommitment, tc.round, tc.startBlockHeight, tc.startBlockHash, tc.dkgPubKey, tc.commPubKey, tc.rawQuote)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)

				if tc.expectedRegData != nil {
					storedReg, err := k.getDKGRegistration(ctx, tc.codeCommitment, tc.round, tc.msgSender)
					require.NoError(t, err)
					require.Equal(t, tc.expectedRegData.Round, storedReg.Round)
					require.Equal(t, tc.expectedRegData.ValidatorAddr, storedReg.ValidatorAddr)
					require.Equal(t, tc.expectedRegData.Index, storedReg.Index)
					require.Equal(t, tc.expectedRegData.DkgPubKey, storedReg.DkgPubKey)
					require.Equal(t, tc.expectedRegData.CommPubKey, storedReg.CommPubKey)
					require.Equal(t, tc.expectedRegData.RawQuote, storedReg.RawQuote)
					require.Equal(t, tc.expectedRegData.Status, storedReg.Status)
				}
			}
		})
	}
}

func TestKeeper_Finalized(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testValidator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testCodeCommitment := [32]byte{0x12, 0x34, 0x56, 0x78}
	testParticipantsRoot := [32]byte{0x12, 0x34, 0x56, 0x78}
	testRound := uint32(1)
	testSignature := []byte("test-signature")

	testReg := &types.DKGRegistration{
		Round:         testRound,
		ValidatorAddr: testValidator.Hex(),
		Index:         1,
		DkgPubKey:     []byte("test-dkg-pubkey"),
		CommPubKey:    []byte("test-comm-pubkey"),
		RawQuote:      []byte("test-raw-quote"),
		Status:        types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, testCodeCommitment, testValidator, testReg))

	tcs := []struct {
		name             string
		round            uint32
		msgSender        common.Address
		codeCommitment   [32]byte
		participantsRoot [32]byte
		globalPubKey     []byte
		signature        []byte
		expectedErr      string
	}{
		{
			name:             "pass: successful finalization",
			round:            testRound,
			msgSender:        testValidator,
			codeCommitment:   testCodeCommitment,
			participantsRoot: testParticipantsRoot,
			signature:        testSignature,
		},
		{
			name:             "fail: registration not found",
			round:            testRound,
			msgSender:        common.HexToAddress("0x9999999999999999999999999999999999999999"),
			codeCommitment:   testCodeCommitment,
			participantsRoot: testParticipantsRoot,
			signature:        testSignature,
			expectedErr:      "dkg registration not found",
		},
		// TODO: add tc for setting global pub key
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.Finalized(ctx, tc.round, tc.msgSender, tc.codeCommitment, tc.participantsRoot, tc.signature, tc.globalPubKey, [][]byte{}) // TODO mock public coeffs

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)

				updatedReg, err := k.getDKGRegistration(ctx, tc.codeCommitment, tc.round, tc.msgSender)
				require.NoError(t, err)
				require.Equal(t, types.DKGRegStatusFinalized, updatedReg.Status)
			}
		})
	}
}

func TestKeeper_UpgradeScheduled(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testActivationHeight := uint32(1000)
	testCodeCommitment := [32]byte{0x12, 0x34, 0x56, 0x78}

	tcs := []struct {
		name             string
		activationHeight uint32
		codeCommitment   [32]byte
		expectedErr      string
	}{
		{
			name:             "pass: successful upgrade scheduling",
			activationHeight: testActivationHeight,
			codeCommitment:   testCodeCommitment,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.UpgradeScheduled(ctx, tc.activationHeight, tc.codeCommitment)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_RemoteAttestationProcessedOnChain(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testValidator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testCodeCommitment := [32]byte{0x12, 0x34, 0x56, 0x78}
	testRound := uint32(1)
	testChalStatus := 1 // ChallengeStatus.Invalidated

	testReg := &types.DKGRegistration{
		Round:         testRound,
		ValidatorAddr: testValidator.Hex(),
		Index:         1,
		DkgPubKey:     []byte("test-dkg-pubkey"),
		CommPubKey:    []byte("test-comm-pubkey"),
		RawQuote:      []byte("test-raw-quote"),
		Status:        types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, testCodeCommitment, testValidator, testReg))

	tcs := []struct {
		name           string
		validator      common.Address
		chalStatus     int
		round          uint32
		codeCommitment [32]byte
		expectedErr    string
	}{
		{
			name:           "pass: successful remote attestation processing",
			validator:      testValidator,
			chalStatus:     testChalStatus,
			round:          testRound,
			codeCommitment: testCodeCommitment,
		},
		{
			name:           "fail: registration not found",
			validator:      common.HexToAddress("0x9999999999999999999999999999999999999999"),
			chalStatus:     testChalStatus,
			round:          testRound,
			codeCommitment: testCodeCommitment,
			expectedErr:    "dkg registration not found",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.RemoteAttestationProcessedOnChain(ctx, tc.validator, tc.chalStatus, tc.round, tc.codeCommitment)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_DealComplaintsSubmitted(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testIndex := uint32(1)
	testComplainIndexes := []uint32{1, 2, 3}
	testRound := uint32(1)
	testCodeCommitment := [32]byte{0x12, 0x34, 0x56, 0x78}

	tcs := []struct {
		name            string
		index           uint32
		complainIndexes []uint32
		round           uint32
		codeCommitment  [32]byte
		expectedErr     string
	}{
		{
			name:            "pass: successful deal complaints submission",
			index:           testIndex,
			complainIndexes: testComplainIndexes,
			round:           testRound,
			codeCommitment:  testCodeCommitment,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.DealComplaintsSubmitted(ctx, tc.index, tc.complainIndexes, tc.round, tc.codeCommitment)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_DealVerified(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testIndex := uint32(1)
	testRecipientIndex := uint32(2)
	testRound := uint32(1)
	testCodeCommitment := [32]byte{0x12, 0x34, 0x56, 0x78}

	tcs := []struct {
		name           string
		index          uint32
		recipientIndex uint32
		round          uint32
		codeCommitment [32]byte
		expectedErr    string
	}{
		{
			name:           "pass: successful deal verification",
			index:          testIndex,
			recipientIndex: testRecipientIndex,
			round:          testRound,
			codeCommitment: testCodeCommitment,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.DealVerified(ctx, tc.index, tc.recipientIndex, tc.round, tc.codeCommitment)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_InvalidDeal(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testIndex := uint32(1)
	testRound := uint32(1)
	testCodeCommitment := [32]byte{0x12, 0x34, 0x56, 0x78}

	tcs := []struct {
		name           string
		index          uint32
		round          uint32
		codeCommitment [32]byte
		expectedErr    string
	}{
		{
			name:           "pass: successful invalid deal processing",
			index:          testIndex,
			round:          testRound,
			codeCommitment: testCodeCommitment,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.InvalidDeal(ctx, tc.index, tc.round, tc.codeCommitment)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// setupDKGKeeper creates a test DKG keeper with necessary dependencies.
func setupDKGKeeper(t *testing.T) (*Keeper, context.Context) {
	t.Helper()

	k, _, _, testCtx := setupDKGKeeperWithMocks(t)

	return k, testCtx
}

// setupDKGKeeperWithMocks creates a test DKG keeper and returns the mock keepers for fine-grained control.
func setupDKGKeeperWithMocks(t *testing.T) (*Keeper, *dkgtestutil.MockBankKeeper, *dkgtestutil.MockDistributionKeeper, context.Context) {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig()

	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	ak := dkgtestutil.NewMockAccountKeeper(ctrl)
	sk := dkgtestutil.NewMockStakingKeeper(ctrl)
	bk := dkgtestutil.NewMockBankKeeper(ctrl)
	dk := dkgtestutil.NewMockDistributionKeeper(ctrl)

	ak.EXPECT().AddressCodec().Return(authcodec.NewBech32Codec("story")).AnyTimes()
	ak.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.AccAddress{}).AnyTimes()

	var valStore baseapp.ValidatorStore = nil

	mockTEEClient := dkgtestutil.NewMockTEEClient(ctrl)

	k := NewKeeper(
		encCfg.Codec,
		storeService,
		ak,
		bk,
		dk,
		sk,
		valStore,
		mockTEEClient,
		nil, // TODO: mock contract client for integration test
		"story1hmjw3pvkjtndpg8wqppwdn8udd835qpan4hm0y",
	)

	require.NoError(t, k.SetParams(testCtx.Ctx, types.DefaultParams()))

	return k, bk, dk, testCtx.Ctx
}
