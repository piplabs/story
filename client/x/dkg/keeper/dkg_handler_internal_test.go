package keeper

import (
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"slices"
	"strings"
	"testing"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
			name:             "fail: codeCommitment mismatch",
			msgSender:        testValidator,
			codeCommitment:   [32]byte{0x99, 0x99, 0x99, 0x99},
			round:            testRound,
			startBlockHeight: testStartBlockHeight,
			startBlockHash:   testStartBlockHash,
			dkgPubKey:        testDkgPubKey,
			commPubKey:       testCommPubKey,
			rawQuote:         testRawQuote,
			expectedErr:      "codeCommitment mismatch",
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
			// Index is 3 because the first "pass" subtest already created a registration (index 1),
			// and setupNetwork adds anotherValidator (index 1), so getNextDKGRegistrationIndex returns 3.
			expectedRegData: &types.DKGRegistration{
				Round:         testRound,
				ValidatorAddr: testValidator.Hex(),
				Index:         3,
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
	testValidator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testCodeCommitment := [32]byte{0x12, 0x34, 0x56, 0x78}
	testRound := uint32(1)
	testGlobalPubKey := []byte("test-global-pubkey")
	testPublicCoeffs := [][]byte{[]byte("coeff1"), []byte("coeff2")}
	testPubKeyShare := []byte("test-pubkey-share")

	type finalizedArgs struct {
		k                *Keeper
		ctx              context.Context
		round            uint32
		msgSender        common.Address
		codeCommitment   [32]byte
		participantsRoot [32]byte
		signature        []byte
		globalPubKey     []byte
		publicCoeffs     [][]byte
		pubKeyShare      []byte
	}

	tcs := []struct {
		name        string
		setup       func(t *testing.T) finalizedArgs
		expectedErr string
		postCheck   func(t *testing.T, args finalizedArgs)
	}{
		{
			name: "pass: successful finalization",
			setup: func(t *testing.T) finalizedArgs {
				t.Helper()
				k, ctx := setupDKGKeeper(t)
				sigKey, commPubKey, pRoot := setupFinalizedState(t, k, ctx, testValidator, testCodeCommitment, testRound)
				setVerifiedRegistration(t, k, ctx, testCodeCommitment, testValidator, testRound, 1, commPubKey)
				sig := signFinalizationData(t, sigKey, testCodeCommitment, testRound, pRoot, testGlobalPubKey, testPublicCoeffs, testPubKeyShare)
				return finalizedArgs{k, ctx, testRound, testValidator, testCodeCommitment, pRoot, sig, testGlobalPubKey, testPublicCoeffs, testPubKeyShare}
			},
			postCheck: func(t *testing.T, args finalizedArgs) {
				t.Helper()
				reg, err := args.k.getDKGRegistration(args.ctx, args.codeCommitment, args.round, args.msgSender)
				require.NoError(t, err)
				require.Equal(t, types.DKGRegStatusFinalized, reg.Status)
			},
		},
		{
			name: "fail: round mismatch",
			setup: func(t *testing.T) finalizedArgs {
				t.Helper()
				k, ctx := setupDKGKeeper(t)
				network := &types.DKGNetwork{
					CodeCommitment: testCodeCommitment[:], Round: 99,
					Total: 5, Threshold: 3, Stage: types.DKGStageFinalization,
				}
				require.NoError(t, k.setDKGNetwork(ctx, network))
				return finalizedArgs{k: k, ctx: ctx, round: testRound, msgSender: testValidator, codeCommitment: testCodeCommitment, globalPubKey: testGlobalPubKey, publicCoeffs: testPublicCoeffs, pubKeyShare: testPubKeyShare}
			},
			expectedErr: "round mismatch",
		},
		{
			name: "fail: codeCommitment mismatch",
			setup: func(t *testing.T) finalizedArgs {
				t.Helper()
				k, ctx := setupDKGKeeper(t)
				otherCommitment := [32]byte{0xFF, 0xEE, 0xDD}
				network := &types.DKGNetwork{
					CodeCommitment: otherCommitment[:], Round: testRound,
					Total: 5, Threshold: 3, Stage: types.DKGStageFinalization,
				}
				require.NoError(t, k.setDKGNetwork(ctx, network))
				return finalizedArgs{k: k, ctx: ctx, round: testRound, msgSender: testValidator, codeCommitment: testCodeCommitment, globalPubKey: testGlobalPubKey, publicCoeffs: testPublicCoeffs, pubKeyShare: testPubKeyShare}
			},
			expectedErr: "codeCommitment mismatch",
		},
		{
			name: "fail: stage not finalization",
			setup: func(t *testing.T) finalizedArgs {
				t.Helper()
				k, ctx := setupDKGKeeper(t)
				network := &types.DKGNetwork{
					CodeCommitment: testCodeCommitment[:], Round: testRound,
					ActiveValSet: []string{testValidator.Hex()},
					Total:        5, Threshold: 3, Stage: types.DKGStageDealing,
				}
				require.NoError(t, k.setDKGNetwork(ctx, network))
				setVerifiedRegistration(t, k, ctx, testCodeCommitment, testValidator, testRound, 1, []byte("comm-key"))
				pRoot := computeParticipantsRoot(testValidator)
				return finalizedArgs{k: k, ctx: ctx, round: testRound, msgSender: testValidator, codeCommitment: testCodeCommitment, participantsRoot: pRoot, globalPubKey: testGlobalPubKey, publicCoeffs: testPublicCoeffs, pubKeyShare: testPubKeyShare}
			},
			expectedErr: "round is not in network set stage",
		},
		{
			name: "fail: registration not found for signature verification",
			setup: func(t *testing.T) finalizedArgs {
				t.Helper()
				k, ctx := setupDKGKeeper(t)
				unknownSender := common.HexToAddress("0x1111222233334444555566667777888899990000")

				network := &types.DKGNetwork{
					CodeCommitment: testCodeCommitment[:], Round: testRound,
					ActiveValSet: []string{testValidator.Hex()},
					Total:        5, Threshold: 3, Stage: types.DKGStageFinalization,
				}
				require.NoError(t, k.setDKGNetwork(ctx, network))

				// Register testValidator as verified (for participants root validation to pass)
				setVerifiedRegistration(t, k, ctx, testCodeCommitment, testValidator, testRound, 1, []byte("comm-key"))
				pRoot := computeParticipantsRoot(testValidator)

				// Call Finalized with unknownSender who has no registration
				return finalizedArgs{k: k, ctx: ctx, round: testRound, msgSender: unknownSender, codeCommitment: testCodeCommitment, participantsRoot: pRoot, globalPubKey: testGlobalPubKey, publicCoeffs: testPublicCoeffs, pubKeyShare: testPubKeyShare}
			},
			expectedErr: "failed to get DKG registration for signature verification",
		},
		{
			name: "fail: invalid signature (wrong signer)",
			setup: func(t *testing.T) finalizedArgs {
				t.Helper()
				k, ctx := setupDKGKeeper(t)
				_, commPubKey, pRoot := setupFinalizedState(t, k, ctx, testValidator, testCodeCommitment, testRound)
				setVerifiedRegistration(t, k, ctx, testCodeCommitment, testValidator, testRound, 1, commPubKey)

				// Sign with a DIFFERENT key so address won't match commPubKey
				wrongKey, err := crypto.GenerateKey()
				require.NoError(t, err)
				badSig := signFinalizationData(t, wrongKey, testCodeCommitment, testRound, pRoot, testGlobalPubKey, testPublicCoeffs, testPubKeyShare)
				return finalizedArgs{k: k, ctx: ctx, round: testRound, msgSender: testValidator, codeCommitment: testCodeCommitment, participantsRoot: pRoot, signature: badSig, globalPubKey: testGlobalPubKey, publicCoeffs: testPublicCoeffs, pubKeyShare: testPubKeyShare}
			},
			expectedErr: "finalization signature verification failed",
		},
		{
			name: "fail: double finalization rejected",
			setup: func(t *testing.T) finalizedArgs {
				t.Helper()
				k, ctx := setupDKGKeeper(t)

				// Two validators: testValidator will finalize twice, otherValidator stays Verified
				otherValidator := common.HexToAddress("0xAAAABBBBCCCCDDDDEEEEFFFF0000111122223333")

				sigKey, err := crypto.GenerateKey()
				require.NoError(t, err)
				commPubKey := crypto.FromECDSAPub(&sigKey.PublicKey)[1:]

				// Compute participants root with both validators
				pRoot := computeParticipantsRoot(testValidator, otherValidator)

				network := &types.DKGNetwork{
					CodeCommitment: testCodeCommitment[:], Round: testRound,
					ActiveValSet: []string{testValidator.Hex(), otherValidator.Hex()},
					Total:        2, Threshold: 2, Stage: types.DKGStageFinalization,
				}
				require.NoError(t, k.setDKGNetwork(ctx, network))

				// Both validators are Verified
				setVerifiedRegistration(t, k, ctx, testCodeCommitment, testValidator, testRound, 1, commPubKey)
				setVerifiedRegistration(t, k, ctx, testCodeCommitment, otherValidator, testRound, 2, []byte("other-comm-key-padding-to-64-bytes-1234567890123456789012345678"))

				sig := signFinalizationData(t, sigKey, testCodeCommitment, testRound, pRoot, testGlobalPubKey, testPublicCoeffs, testPubKeyShare)

				// First finalization should succeed
				err = k.Finalized(ctx, testRound, testValidator, testCodeCommitment, pRoot, sig, testGlobalPubKey, testPublicCoeffs, testPubKeyShare)
				require.NoError(t, err)

				// After first finalization, testValidator status is Finalized and otherValidator is still Verified.
				// validateParticipantsRoot uses only Verified registrations, so the new root only includes otherValidator.
				newRoot := computeParticipantsRoot(otherValidator)

				// Return args for a second call which should hit the double-finalization guard
				return finalizedArgs{k: k, ctx: ctx, round: testRound, msgSender: testValidator, codeCommitment: testCodeCommitment, participantsRoot: newRoot, signature: sig, globalPubKey: testGlobalPubKey, publicCoeffs: testPublicCoeffs, pubKeyShare: testPubKeyShare}
			},
			expectedErr: "validator has already finalized for this round",
		},
		{
			name: "pass: global pub key set when threshold reached",
			setup: func(t *testing.T) finalizedArgs {
				t.Helper()
				k, ctx := setupDKGKeeper(t)
				sigKey, commPubKey, pRoot := setupFinalizedState(t, k, ctx, testValidator, testCodeCommitment, testRound)
				// Override threshold to 1 so a single vote triggers global key set
				network := &types.DKGNetwork{
					CodeCommitment: testCodeCommitment[:], Round: testRound,
					ActiveValSet: []string{testValidator.Hex()},
					Total:        1, Threshold: 1, Stage: types.DKGStageFinalization,
				}
				require.NoError(t, k.setDKGNetwork(ctx, network))
				setVerifiedRegistration(t, k, ctx, testCodeCommitment, testValidator, testRound, 1, commPubKey)
				sig := signFinalizationData(t, sigKey, testCodeCommitment, testRound, pRoot, testGlobalPubKey, testPublicCoeffs, testPubKeyShare)
				return finalizedArgs{k: k, ctx: ctx, round: testRound, msgSender: testValidator, codeCommitment: testCodeCommitment, participantsRoot: pRoot, signature: sig, globalPubKey: testGlobalPubKey, publicCoeffs: testPublicCoeffs, pubKeyShare: testPubKeyShare}
			},
			postCheck: func(t *testing.T, args finalizedArgs) {
				t.Helper()
				net, err := args.k.getLatestDKGNetwork(args.ctx)
				require.NoError(t, err)
				require.Equal(t, args.globalPubKey, net.GlobalPublicKey)
				require.Equal(t, args.publicCoeffs, net.PublicCoeffs)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			args := tc.setup(t)
			err := args.k.Finalized(args.ctx, args.round, args.msgSender, args.codeCommitment, args.participantsRoot, args.signature, args.globalPubKey, args.publicCoeffs, args.pubKeyShare)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				if tc.postCheck != nil {
					tc.postCheck(t, args)
				}
			}
		})
	}
}

// setupFinalizedState creates a DKG network in Finalization stage and returns the signing key,
// commPubKey, and participantsRoot.
func setupFinalizedState(t *testing.T, k *Keeper, ctx context.Context, validator common.Address, codeCommitment [32]byte, round uint32) (*ecdsa.PrivateKey, []byte, [32]byte) {
	t.Helper()

	sigKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	commPubKey := crypto.FromECDSAPub(&sigKey.PublicKey)[1:]

	network := &types.DKGNetwork{
		CodeCommitment: codeCommitment[:], Round: round,
		ActiveValSet: []string{validator.Hex()},
		Total:        5, Threshold: 3, Stage: types.DKGStageFinalization,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	return sigKey, commPubKey, computeParticipantsRoot(validator)
}

// setVerifiedRegistration creates a DKG registration in Verified status.
func setVerifiedRegistration(t *testing.T, k *Keeper, ctx context.Context, codeCommitment [32]byte, validator common.Address, round uint32, index uint32, commPubKey []byte) {
	t.Helper()
	reg := &types.DKGRegistration{
		Round: round, ValidatorAddr: validator.Hex(), Index: index,
		DkgPubKey: []byte("dkg-key"), CommPubKey: commPubKey,
		RawQuote: []byte("raw-quote"), Status: types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, codeCommitment, validator, reg))
}

// computeParticipantsRoot computes the participants root hash for one or more validators.
// Addresses are sorted in ascending order before hashing, matching validateParticipantsRoot.
func computeParticipantsRoot(validators ...common.Address) [32]byte {
	addrs := make([]string, 0, len(validators))
	for _, v := range validators {
		addrs = append(addrs, strings.ToLower(v.Hex()))
	}
	slices.Sort(addrs)

	buf := make([]byte, 0, common.AddressLength*len(addrs))
	for _, a := range addrs {
		buf = append(buf, common.HexToAddress(a).Bytes()...)
	}

	h := crypto.Keccak256(buf)
	var root [32]byte
	copy(root[:], h)
	return root
}

func TestVerifyFinalizationSignature(t *testing.T) {
	// Generate a secp256k1 key pair for testing
	sigKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	commPubKey := crypto.FromECDSAPub(&sigKey.PublicKey)[1:] // 64 bytes, without 0x04 prefix

	codeCommitment := [32]byte{0x01, 0x02, 0x03, 0x04}
	round := uint32(42)
	participantsRoot := [32]byte{0xAA, 0xBB, 0xCC, 0xDD}
	globalPubKey := []byte("global-pub-key-data")
	publicCoeffs := [][]byte{[]byte("coeff-0"), []byte("coeff-1"), []byte("coeff-2")}
	pubKeyShare := []byte("pub-key-share")

	validSig := signFinalizationData(t, sigKey, codeCommitment, round, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare)

	tcs := []struct {
		name             string
		commPubKey       []byte
		round            uint32
		codeCommitment   [32]byte
		participantsRoot [32]byte
		globalPubKey     []byte
		publicCoeffs     [][]byte
		pubKeyShare      []byte
		signature        []byte
		expectedErr      string
	}{
		{
			name:             "pass: valid signature",
			commPubKey:       commPubKey,
			round:            round,
			codeCommitment:   codeCommitment,
			participantsRoot: participantsRoot,
			globalPubKey:     globalPubKey,
			publicCoeffs:     publicCoeffs,
			pubKeyShare:      pubKeyShare,
			signature:        validSig,
		},
		{
			name:             "fail: tampered codeCommitment",
			commPubKey:       commPubKey,
			round:            round,
			codeCommitment:   [32]byte{0xFF, 0xFF, 0xFF, 0xFF},
			participantsRoot: participantsRoot,
			globalPubKey:     globalPubKey,
			publicCoeffs:     publicCoeffs,
			pubKeyShare:      pubKeyShare,
			signature:        validSig,
			expectedErr:      "finalization signature address mismatch",
		},
		{
			name:             "fail: tampered globalPubKey",
			commPubKey:       commPubKey,
			round:            round,
			codeCommitment:   codeCommitment,
			participantsRoot: participantsRoot,
			globalPubKey:     []byte("tampered-global-pub-key"),
			publicCoeffs:     publicCoeffs,
			pubKeyShare:      pubKeyShare,
			signature:        validSig,
			expectedErr:      "finalization signature address mismatch",
		},
		{
			name:             "fail: wrong commPubKey (64 bytes but different key)",
			commPubKey:       make([]byte, 64), // all zeros, 64 bytes, but wrong key
			round:            round,
			codeCommitment:   codeCommitment,
			participantsRoot: participantsRoot,
			globalPubKey:     globalPubKey,
			publicCoeffs:     publicCoeffs,
			pubKeyShare:      pubKeyShare,
			signature:        validSig,
			expectedErr:      "finalization signature address mismatch",
		},
		{
			name:             "fail: commPubKey too short",
			commPubKey:       []byte("short"),
			round:            round,
			codeCommitment:   codeCommitment,
			participantsRoot: participantsRoot,
			globalPubKey:     globalPubKey,
			publicCoeffs:     publicCoeffs,
			pubKeyShare:      pubKeyShare,
			signature:        validSig,
			expectedErr:      "invalid commPubKey length",
		},
		{
			name:             "fail: commPubKey 65 bytes (with 0x04 prefix)",
			commPubKey:       append([]byte{0x04}, commPubKey...),
			round:            round,
			codeCommitment:   codeCommitment,
			participantsRoot: participantsRoot,
			globalPubKey:     globalPubKey,
			publicCoeffs:     publicCoeffs,
			pubKeyShare:      pubKeyShare,
			signature:        validSig,
			expectedErr:      "invalid commPubKey length",
		},
		{
			name:             "fail: empty public coefficient",
			commPubKey:       commPubKey,
			round:            round,
			codeCommitment:   codeCommitment,
			participantsRoot: participantsRoot,
			globalPubKey:     globalPubKey,
			publicCoeffs:     [][]byte{[]byte("coeff-0"), {}, []byte("coeff-2")},
			pubKeyShare:      pubKeyShare,
			signature:        validSig,
			expectedErr:      "empty public coefficient",
		},
		{
			name:             "fail: invalid pub key share",
			commPubKey:       commPubKey,
			round:            round,
			codeCommitment:   codeCommitment,
			participantsRoot: participantsRoot,
			globalPubKey:     globalPubKey,
			publicCoeffs:     publicCoeffs,
			pubKeyShare:      []byte("invalid-pub-key-share"),
			signature:        validSig,
			expectedErr:      "finalization signature address mismatch",
		},
		{
			name:             "fail: invalid signature bytes",
			commPubKey:       commPubKey,
			round:            round,
			codeCommitment:   codeCommitment,
			participantsRoot: participantsRoot,
			globalPubKey:     globalPubKey,
			publicCoeffs:     publicCoeffs,
			pubKeyShare:      pubKeyShare,
			signature:        []byte("too-short"),
			expectedErr:      "invalid signature length",
		},
		{
			name:             "fail: corrupted 65-byte signature",
			commPubKey:       commPubKey,
			round:            round,
			codeCommitment:   codeCommitment,
			participantsRoot: participantsRoot,
			globalPubKey:     globalPubKey,
			publicCoeffs:     publicCoeffs,
			pubKeyShare:      pubKeyShare,
			signature:        make([]byte, 65), // all zeros
			expectedErr:      "failed to recover public key from signature",
		},
	}

	// Test V=28 signature separately since we need a key that produces recovery ID 1
	t.Run("pass: valid signature with V=28", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			k, err := crypto.GenerateKey()
			require.NoError(t, err)
			cpk := crypto.FromECDSAPub(&k.PublicKey)[1:]
			sig := signFinalizationData(t, k, codeCommitment, round, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare)
			if sig[64] == 28 {
				err := verifyFinalizationSignature(cpk, round, codeCommitment, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare, sig)
				require.NoError(t, err)
				return
			}
		}
		t.Skip("could not generate V=28 signature in 100 attempts")
	})

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := verifyFinalizationSignature(tc.commPubKey, tc.round, tc.codeCommitment, tc.participantsRoot, tc.globalPubKey, tc.publicCoeffs, tc.pubKeyShare, tc.signature)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// signFinalizationData creates a valid finalization signature for testing.
func signFinalizationData(t *testing.T, key *ecdsa.PrivateKey, codeCommitment [32]byte, round uint32, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte) []byte {
	t.Helper()

	encoded := make([]byte, 0)
	encoded = append(encoded, codeCommitment[:]...)
	roundBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(roundBytes, round)
	encoded = append(encoded, roundBytes...)
	encoded = append(encoded, participantsRoot[:]...)
	encoded = append(encoded, globalPubKey...)
	for _, coeff := range publicCoeffs {
		encoded = append(encoded, coeff...)
	}
	encoded = append(encoded, pubKeyShare...)

	msgHash := crypto.Keccak256(encoded)
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	ethHash := crypto.Keccak256(append(prefix, msgHash...))

	sig, err := crypto.Sign(ethHash, key)
	require.NoError(t, err)

	// Convert recovery ID to Ethereum V (add 27)
	sig[64] += 27

	return sig
}

func TestVerifyPartialDecryptionSignature(t *testing.T) {
	sigKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	commPubKey := crypto.FromECDSAPub(&sigKey.PublicKey)[1:]

	codeCommitment := [32]byte{0x0A, 0x0B, 0x0C, 0x0D}
	round := uint32(7)
	encryptedPartial := []byte("encrypted-partial")
	ephemeralPubKey := []byte("ephemeral-pub-key")
	pubShare := []byte("pub-share")

	validSig := signPartialDecryptionData(t, sigKey, codeCommitment, round, encryptedPartial, ephemeralPubKey, pubShare)

	tcs := []struct {
		name             string
		commPubKey       []byte
		codeCommitment   [32]byte
		round            uint32
		encryptedPartial []byte
		ephemeralPubKey  []byte
		pubShare         []byte
		signature        []byte
		expectedErr      string
	}{
		{
			name:             "pass: valid signature",
			commPubKey:       commPubKey,
			codeCommitment:   codeCommitment,
			round:            round,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        validSig,
		},
		{
			name:             "fail: tampered encryptedPartial",
			commPubKey:       commPubKey,
			codeCommitment:   codeCommitment,
			round:            round,
			encryptedPartial: []byte("tampered-partial"),
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        validSig,
			expectedErr:      "partial decryption signature address mismatch",
		},
		{
			name:             "fail: tampered ephemeralPubKey",
			commPubKey:       commPubKey,
			codeCommitment:   codeCommitment,
			round:            round,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  []byte("tampered-ephemeral"),
			pubShare:         pubShare,
			signature:        validSig,
			expectedErr:      "partial decryption signature address mismatch",
		},
		{
			name:             "fail: wrong commPubKey (64 bytes but different key)",
			commPubKey:       make([]byte, 64),
			codeCommitment:   codeCommitment,
			round:            round,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        validSig,
			expectedErr:      "partial decryption signature address mismatch",
		},
		{
			name:             "fail: commPubKey too short",
			commPubKey:       []byte("short"),
			codeCommitment:   codeCommitment,
			round:            round,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        validSig,
			expectedErr:      "invalid commPubKey length",
		},
		{
			name:             "fail: commPubKey 65 bytes (with 0x04 prefix)",
			commPubKey:       append([]byte{0x04}, commPubKey...),
			codeCommitment:   codeCommitment,
			round:            round,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        validSig,
			expectedErr:      "invalid commPubKey length",
		},
		{
			name:             "fail: invalid signature bytes",
			commPubKey:       commPubKey,
			codeCommitment:   codeCommitment,
			round:            round,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        []byte("too-short"),
			expectedErr:      "invalid signature length",
		},
		{
			name:             "fail: corrupted 65-byte signature",
			commPubKey:       commPubKey,
			codeCommitment:   codeCommitment,
			round:            round,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        make([]byte, 65),
			expectedErr:      "failed to recover public key from signature",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := verifyPartialDecryptionSignature(tc.commPubKey, tc.codeCommitment, tc.round, tc.encryptedPartial, tc.ephemeralPubKey, tc.pubShare, tc.signature)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func signPartialDecryptionData(t *testing.T, key *ecdsa.PrivateKey, codeCommitment [32]byte, round uint32, encryptedPartial, ephemeralPubKey, pubShare []byte) []byte {
	t.Helper()

	encoded := make([]byte, 0, len(codeCommitment)+4+len(encryptedPartial)+len(ephemeralPubKey)+len(pubShare))
	encoded = append(encoded, codeCommitment[:]...)
	roundBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(roundBytes, round)
	encoded = append(encoded, roundBytes...)
	encoded = append(encoded, encryptedPartial...)
	encoded = append(encoded, ephemeralPubKey...)
	encoded = append(encoded, pubShare...)

	respHash := crypto.Keccak256(encoded)
	sig, err := crypto.Sign(respHash, key)
	require.NoError(t, err)

	sig[64] += 27

	return sig
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

func TestKeeper_HasFinalizedRegistration(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testValidator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testCodeCommitment := [32]byte{0x12, 0x34, 0x56, 0x78}
	testRound := uint32(1)

	tcs := []struct {
		name           string
		setup          func()
		codeCommitment []byte
		round          uint32
		validatorAddr  common.Address
		expectedResult bool
		expectedErr    string
	}{
		{
			name: "pass: returns true when registration is finalized",
			setup: func() {
				err := k.setDKGRegistration(ctx, testCodeCommitment, testValidator, &types.DKGRegistration{
					Round:         testRound,
					ValidatorAddr: testValidator.Hex(),
					Index:         1,
					Status:        types.DKGRegStatusFinalized,
				})
				require.NoError(t, err)
			},
			codeCommitment: testCodeCommitment[:],
			round:          testRound,
			validatorAddr:  testValidator,
			expectedResult: true,
		},
		{
			name: "pass: returns false when registration is verified (not finalized)",
			setup: func() {
				err := k.setDKGRegistration(ctx, testCodeCommitment, testValidator, &types.DKGRegistration{
					Round:         testRound,
					ValidatorAddr: testValidator.Hex(),
					Index:         1,
					Status:        types.DKGRegStatusVerified,
				})
				require.NoError(t, err)
			},
			codeCommitment: testCodeCommitment[:],
			round:          testRound,
			validatorAddr:  testValidator,
			expectedResult: false,
		},
		{
			name: "pass: returns false when registration is unspecified",
			setup: func() {
				err := k.setDKGRegistration(ctx, testCodeCommitment, testValidator, &types.DKGRegistration{
					Round:         testRound,
					ValidatorAddr: testValidator.Hex(),
					Index:         1,
					Status:        types.DKGRegStatusUnspecified,
				})
				require.NoError(t, err)
			},
			codeCommitment: testCodeCommitment[:],
			round:          testRound,
			validatorAddr:  testValidator,
			expectedResult: false,
		},
		{
			name:           "pass: returns false when no registration exists",
			setup:          func() {},
			codeCommitment: testCodeCommitment[:],
			round:          uint32(999), // non-existent round
			validatorAddr:  common.HexToAddress("0x0000000000000000000000000000000000000001"),
			expectedResult: false,
		},
		{
			name:           "fail: returns error when code commitment has wrong length",
			setup:          func() {},
			codeCommitment: []byte{0x01, 0x02}, // too short
			round:          testRound,
			validatorAddr:  testValidator,
			expectedResult: false,
			expectedErr:    "failed to cast code commitment to bytes32",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			result, err := k.HasFinalizedRegistration(ctx, tc.codeCommitment, tc.round, tc.validatorAddr)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.False(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
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
