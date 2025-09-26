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

	testutil2 "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/ethclient/mock"

	"go.uber.org/mock/gomock"
)

func TestKeeper_RegistrationInitialized(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testValidator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testMrenclave := [32]byte{0x12, 0x34, 0x56, 0x78}
	testRound := uint32(1)
	testDkgPubKey := []byte("test-dkg-pubkey")
	testCommPubKey := []byte("test-comm-pubkey")
	testRawQuote := []byte("test-raw-quote")

	validDKGNetwork := &types.DKGNetwork{
		Mrenclave:    testMrenclave[:],
		Round:        testRound,
		StartBlock:   100,
		ActiveValSet: []string{testValidator.Hex()},
		Total:        5,
		Threshold:    3,
		Stage:        types.DKGStageRegistration,
	}
	require.NoError(t, k.SetDKGNetwork(ctx, validDKGNetwork))

	tcs := []struct {
		name            string
		msgSender       common.Address
		mrenclave       [32]byte
		round           uint32
		dkgPubKey       []byte
		commPubKey      []byte
		rawQuote        []byte
		setupNetwork    func()
		expectedErr     string
		expectedRegData *types.DKGRegistration
	}{
		{
			name:       "pass: successful registration initialization",
			msgSender:  testValidator,
			mrenclave:  testMrenclave,
			round:      testRound,
			dkgPubKey:  testDkgPubKey,
			commPubKey: testCommPubKey,
			rawQuote:   testRawQuote,
			setupNetwork: func() {
				// Network already set up in test setup
			},
			expectedRegData: &types.DKGRegistration{
				Round:      testRound,
				MsgSender:  testValidator.Hex(),
				Index:      1,
				DkgPubKey:  testDkgPubKey,
				CommPubKey: testCommPubKey,
				RawQuote:   testRawQuote,
				Status:     types.DKGRegStatusVerified,
			},
		},
		{
			name:        "fail: DKG network not found",
			msgSender:   testValidator,
			mrenclave:   [32]byte{0x99, 0x99, 0x99, 0x99},
			round:       testRound,
			dkgPubKey:   testDkgPubKey,
			commPubKey:  testCommPubKey,
			rawQuote:    testRawQuote,
			expectedErr: "dkg network not found",
		},
		{
			name:       "fail: round not in registration stage",
			msgSender:  testValidator,
			mrenclave:  testMrenclave,
			round:      testRound,
			dkgPubKey:  testDkgPubKey,
			commPubKey: testCommPubKey,
			rawQuote:   testRawQuote,
			setupNetwork: func() {
				networkWithDifferentStage := &types.DKGNetwork{
					Mrenclave:    testMrenclave[:],
					Round:        testRound,
					StartBlock:   100,
					ActiveValSet: []string{testValidator.Hex()},
					Total:        5,
					Threshold:    3,
					Stage:        types.DKGStageNetworkSet,
				}
				require.NoError(t, k.SetDKGNetwork(ctx, networkWithDifferentStage))
			},
			expectedErr: "round is not in registration stage",
		},
		{
			name:       "fail: validator not in active set",
			msgSender:  common.HexToAddress("0x9999999999999999999999999999999999999999"),
			mrenclave:  testMrenclave,
			round:      testRound,
			dkgPubKey:  testDkgPubKey,
			commPubKey: testCommPubKey,
			rawQuote:   testRawQuote,
			setupNetwork: func() {
				networkInRegistrationStage := &types.DKGNetwork{
					Mrenclave:    testMrenclave[:],
					Round:        testRound,
					StartBlock:   100,
					ActiveValSet: []string{testValidator.Hex()},
					Total:        5,
					Threshold:    3,
					Stage:        types.DKGStageRegistration,
				}
				require.NoError(t, k.SetDKGNetwork(ctx, networkInRegistrationStage))
			},
			expectedErr: "msg sender is not in the active validator set",
		},
		{
			name:       "pass: second registration gets incremented index",
			msgSender:  testValidator,
			mrenclave:  testMrenclave,
			round:      testRound,
			dkgPubKey:  []byte("second-dkg-pubkey"),
			commPubKey: []byte("second-comm-pubkey"),
			rawQuote:   []byte("second-raw-quote"),
			setupNetwork: func() {
				anotherValidator := common.HexToAddress("0xAABBCCDDEEFF112233445566778899AABBCCDDEE")
				networkWithMultipleValidators := &types.DKGNetwork{
					Mrenclave:    testMrenclave[:],
					Round:        testRound,
					StartBlock:   100,
					ActiveValSet: []string{testValidator.Hex(), anotherValidator.Hex()},
					Total:        5,
					Threshold:    3,
					Stage:        types.DKGStageRegistration,
				}
				require.NoError(t, k.SetDKGNetwork(ctx, networkWithMultipleValidators))

				firstReg := &types.DKGRegistration{
					Round:      testRound,
					MsgSender:  anotherValidator.Hex(),
					Index:      1,
					DkgPubKey:  []byte("first-dkg-pubkey"),
					CommPubKey: []byte("first-comm-pubkey"),
					RawQuote:   []byte("first-raw-quote"),
					Status:     types.DKGRegStatusVerified,
				}
				require.NoError(t, k.setDKGRegistration(ctx, testMrenclave, anotherValidator, firstReg))
			},
			expectedRegData: &types.DKGRegistration{
				Round:      testRound,
				MsgSender:  testValidator.Hex(),
				Index:      2,
				DkgPubKey:  []byte("second-dkg-pubkey"),
				CommPubKey: []byte("second-comm-pubkey"),
				RawQuote:   []byte("second-raw-quote"),
				Status:     types.DKGRegStatusVerified,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupNetwork != nil {
				tc.setupNetwork()
			}

			err := k.RegistrationInitialized(ctx, tc.msgSender, tc.mrenclave, tc.round, tc.dkgPubKey, tc.commPubKey, tc.rawQuote)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)

				if tc.expectedRegData != nil {
					storedReg, err := k.getDKGRegistration(ctx, tc.mrenclave, tc.round, tc.msgSender)
					require.NoError(t, err)
					require.Equal(t, tc.expectedRegData.Round, storedReg.Round)
					require.Equal(t, tc.expectedRegData.MsgSender, storedReg.MsgSender)
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

func TestKeeper_NetworkSet(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testValidator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testMrenclave := [32]byte{0x12, 0x34, 0x56, 0x78}
	testRound := uint32(1)
	testTotal := uint32(5)
	testThreshold := uint32(3)
	testSignature := []byte("test-signature")

	testReg := &types.DKGRegistration{
		Round:      testRound,
		MsgSender:  testValidator.Hex(),
		Index:      1,
		DkgPubKey:  []byte("test-dkg-pubkey"),
		CommPubKey: []byte("test-comm-pubkey"),
		RawQuote:   []byte("test-raw-quote"),
		Status:     types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, testMrenclave, testValidator, testReg))

	tcs := []struct {
		name        string
		msgSender   common.Address
		mrenclave   [32]byte
		round       uint32
		total       uint32
		threshold   uint32
		signature   []byte
		expectedErr string
	}{
		{
			name:      "pass: successful network set",
			msgSender: testValidator,
			mrenclave: testMrenclave,
			round:     testRound,
			total:     testTotal,
			threshold: testThreshold,
			signature: testSignature,
		},
		{
			name:        "fail: registration not found",
			msgSender:   common.HexToAddress("0x9999999999999999999999999999999999999999"),
			mrenclave:   testMrenclave,
			round:       testRound,
			total:       testTotal,
			threshold:   testThreshold,
			signature:   testSignature,
			expectedErr: "dkg registration not found",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.NetworkSet(ctx, tc.msgSender, tc.mrenclave, tc.round, tc.total, tc.threshold, tc.signature)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)

				updatedReg, err := k.getDKGRegistration(ctx, tc.mrenclave, tc.round, tc.msgSender)
				require.NoError(t, err)
				require.Equal(t, types.DKGRegStatusNetworkSet, updatedReg.Status)
			}
		})
	}
}

func TestKeeper_Finalized(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testValidator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testMrenclave := [32]byte{0x12, 0x34, 0x56, 0x78}
	testRound := uint32(1)
	testSignature := []byte("test-signature")

	testReg := &types.DKGRegistration{
		Round:      testRound,
		MsgSender:  testValidator.Hex(),
		Index:      1,
		DkgPubKey:  []byte("test-dkg-pubkey"),
		CommPubKey: []byte("test-comm-pubkey"),
		RawQuote:   []byte("test-raw-quote"),
		Status:     types.DKGRegStatusNetworkSet,
	}
	require.NoError(t, k.setDKGRegistration(ctx, testMrenclave, testValidator, testReg))

	tcs := []struct {
		name        string
		round       uint32
		msgSender   common.Address
		mrenclave   [32]byte
		signature   []byte
		expectedErr string
	}{
		{
			name:      "pass: successful finalization",
			round:     testRound,
			msgSender: testValidator,
			mrenclave: testMrenclave,
			signature: testSignature,
		},
		{
			name:        "fail: registration not found",
			round:       testRound,
			msgSender:   common.HexToAddress("0x9999999999999999999999999999999999999999"),
			mrenclave:   testMrenclave,
			signature:   testSignature,
			expectedErr: "dkg registration not found",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.Finalized(ctx, tc.round, tc.msgSender, tc.mrenclave, tc.signature)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)

				updatedReg, err := k.getDKGRegistration(ctx, tc.mrenclave, tc.round, tc.msgSender)
				require.NoError(t, err)
				require.Equal(t, types.DKGRegStatusFinalized, updatedReg.Status)
			}
		})
	}
}

func TestKeeper_UpgradeScheduled(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testActivationHeight := uint32(1000)
	testMrenclave := [32]byte{0x12, 0x34, 0x56, 0x78}

	tcs := []struct {
		name             string
		activationHeight uint32
		mrenclave        [32]byte
		expectedErr      string
	}{
		{
			name:             "pass: successful upgrade scheduling",
			activationHeight: testActivationHeight,
			mrenclave:        testMrenclave,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.UpgradeScheduled(ctx, tc.activationHeight, tc.mrenclave)

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
	testMrenclave := [32]byte{0x12, 0x34, 0x56, 0x78}
	testRound := uint32(1)
	testChalStatus := 1 // ChallengeStatus.Invalidated

	testReg := &types.DKGRegistration{
		Round:      testRound,
		MsgSender:  testValidator.Hex(),
		Index:      1,
		DkgPubKey:  []byte("test-dkg-pubkey"),
		CommPubKey: []byte("test-comm-pubkey"),
		RawQuote:   []byte("test-raw-quote"),
		Status:     types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, testMrenclave, testValidator, testReg))

	tcs := []struct {
		name        string
		validator   common.Address
		chalStatus  int
		round       uint32
		mrenclave   [32]byte
		expectedErr string
	}{
		{
			name:       "pass: successful remote attestation processing",
			validator:  testValidator,
			chalStatus: testChalStatus,
			round:      testRound,
			mrenclave:  testMrenclave,
		},
		{
			name:        "fail: registration not found",
			validator:   common.HexToAddress("0x9999999999999999999999999999999999999999"),
			chalStatus:  testChalStatus,
			round:       testRound,
			mrenclave:   testMrenclave,
			expectedErr: "dkg registration not found",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.RemoteAttestationProcessedOnChain(ctx, tc.validator, tc.chalStatus, tc.round, tc.mrenclave)

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
	testMrenclave := [32]byte{0x12, 0x34, 0x56, 0x78}

	tcs := []struct {
		name            string
		index           uint32
		complainIndexes []uint32
		round           uint32
		mrenclave       [32]byte
		expectedErr     string
	}{
		{
			name:            "pass: successful deal complaints submission",
			index:           testIndex,
			complainIndexes: testComplainIndexes,
			round:           testRound,
			mrenclave:       testMrenclave,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.DealComplaintsSubmitted(ctx, tc.index, tc.complainIndexes, tc.round, tc.mrenclave)

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
	testMrenclave := [32]byte{0x12, 0x34, 0x56, 0x78}

	tcs := []struct {
		name           string
		index          uint32
		recipientIndex uint32
		round          uint32
		mrenclave      [32]byte
		expectedErr    string
	}{
		{
			name:           "pass: successful deal verification",
			index:          testIndex,
			recipientIndex: testRecipientIndex,
			round:          testRound,
			mrenclave:      testMrenclave,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.DealVerified(ctx, tc.index, tc.recipientIndex, tc.round, tc.mrenclave)

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
	testMrenclave := [32]byte{0x12, 0x34, 0x56, 0x78}

	tcs := []struct {
		name        string
		index       uint32
		round       uint32
		mrenclave   [32]byte
		expectedErr string
	}{
		{
			name:      "pass: successful invalid deal processing",
			index:     testIndex,
			round:     testRound,
			mrenclave: testMrenclave,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := k.InvalidDeal(ctx, tc.index, tc.round, tc.mrenclave)

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

	encCfg := moduletestutil.MakeTestEncodingConfig()

	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	ak := testutil2.NewMockAccountKeeper(ctrl)
	sk := testutil2.NewMockStakingKeeper(ctrl)

	ak.EXPECT().AddressCodec().Return(authcodec.NewBech32Codec("story")).AnyTimes()
	ak.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.AccAddress{}).AnyTimes()

	ethClient := mock.NewMockClient(ctrl)

	var skeeper baseapp.ValidatorStore = nil

	k := NewKeeper(
		encCfg.Codec,
		storeService,
		ak,
		sk,
		skeeper,
		"story1hmjw3pvkjtndpg8wqppwdn8udd835qpan4hm0y",
		ethClient,
	)

	require.NoError(t, k.SetParams(testCtx.Ctx, types.DefaultParams()))

	return &k, testCtx.Ctx
}
