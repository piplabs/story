package keeper

import (
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"math/big"
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
	testValidator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testCodeCommitment := [32]byte{0x12, 0x34, 0x56, 0x78}
	testRound := uint32(1)
	testStartBlockHeight := big.NewInt(100)
	testStartBlockHash := [32]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB}
	testEnclaveType := [32]byte{0x01}
	testDkgPubKey := []byte("test-dkg-pubkey")
	testCommPubKey := []byte("test-comm-pubkey")
	testEnclaveReport := []byte("test-enclave-report")

	tcs := []struct {
		name             string
		msgSender        common.Address
		codeCommitment   [32]byte
		round            uint32
		startBlockHeight *big.Int
		startBlockHash   [32]byte
		enclaveType      [32]byte
		dkgPubKey        []byte
		commPubKey       []byte
		enclaveReport    []byte
		setupNetwork     func(t *testing.T, k *Keeper, ctx context.Context)
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
			enclaveType:      testEnclaveType,
			dkgPubKey:        testDkgPubKey,
			commPubKey:       testCommPubKey,
			enclaveReport:    testEnclaveReport,
			setupNetwork:     nil,
			expectedRegData: &types.DKGRegistration{
				Round:          testRound,
				ValidatorAddr:  testValidator.Hex(),
				Index:          1,
				DkgPubKey:      testDkgPubKey,
				CommPubKey:     testCommPubKey,
				EnclaveReport:  testEnclaveReport,
				Status:         types.DKGRegStatusVerified,
				CodeCommitment: testCodeCommitment[:],
				EnclaveType:    testEnclaveType[:],
			},
		},
		{
			name:             "fail: start block height mismatch",
			msgSender:        testValidator,
			codeCommitment:   testCodeCommitment,
			round:            testRound,
			startBlockHeight: big.NewInt(999), // Wrong height
			startBlockHash:   testStartBlockHash,
			dkgPubKey:        testDkgPubKey,
			commPubKey:       testCommPubKey,
			enclaveReport:    testEnclaveReport,
			setupNetwork:     nil,
			expectedErr:      "start block height mismatch",
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
			enclaveReport:    testEnclaveReport,
			setupNetwork:     nil,
			expectedErr:      "start block hash mismatch",
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
			enclaveReport:    testEnclaveReport,
			setupNetwork: func(t *testing.T, k *Keeper, ctx context.Context) {
				networkWithDifferentStage := &types.DKGNetwork{
					Round:            testRound,
					StartBlockHeight: testStartBlockHeight.Int64(),
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
			enclaveReport:    testEnclaveReport,
			setupNetwork: func(t *testing.T, k *Keeper, ctx context.Context) {
				networkInRegistrationStage := &types.DKGNetwork{
					Round:            testRound,
					StartBlockHeight: testStartBlockHeight.Int64(),
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
			name:             "fail: validator already registered for round",
			msgSender:        testValidator,
			codeCommitment:   testCodeCommitment,
			round:            testRound + 1,
			startBlockHeight: testStartBlockHeight,
			startBlockHash:   testStartBlockHash,
			enclaveType:      testEnclaveType,
			dkgPubKey:        testDkgPubKey,
			commPubKey:       testCommPubKey,
			enclaveReport:    testEnclaveReport,
			setupNetwork: func(t *testing.T, k *Keeper, ctx context.Context) {
				network := &types.DKGNetwork{
					Round:            testRound + 1,
					StartBlockHeight: testStartBlockHeight.Int64(),
					StartBlockHash:   testStartBlockHash[:],
					ActiveValSet:     []string{testValidator.Hex()},
					Total:            5,
					Threshold:        3,
					Stage:            types.DKGStageRegistration,
				}
				require.NoError(t, k.setDKGNetwork(ctx, network))

				firstReg := &types.DKGRegistration{
					Round:          testRound + 1,
					ValidatorAddr:  testValidator.Hex(),
					Index:          1,
					DkgPubKey:      []byte("first-dkg-pubkey"),
					CommPubKey:     []byte("first-comm-pubkey"),
					EnclaveReport:  []byte("first-enclave-report"),
					Status:         types.DKGRegStatusVerified,
					CodeCommitment: testCodeCommitment[:],
					EnclaveType:    testEnclaveType[:],
				}
				require.NoError(t, k.setDKGRegistration(ctx, testValidator, firstReg))
			},
			expectedErr: "validator already registered for this round",
		},
		{
			name:             "pass: second registration gets incremented index",
			msgSender:        testValidator,
			codeCommitment:   testCodeCommitment,
			round:            testRound,
			startBlockHeight: testStartBlockHeight,
			startBlockHash:   testStartBlockHash,
			enclaveType:      testEnclaveType,
			dkgPubKey:        []byte("second-dkg-pubkey"),
			commPubKey:       []byte("second-comm-pubkey"),
			enclaveReport:    []byte("second-enclave-report"),
			setupNetwork: func(t *testing.T, k *Keeper, ctx context.Context) {
				anotherValidator := common.HexToAddress("0xAABBCCDDEEFF112233445566778899AABBCCDDEE")
				networkWithMultipleValidators := &types.DKGNetwork{
					Round:            testRound,
					StartBlockHeight: testStartBlockHeight.Int64(),
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
					EnclaveReport: []byte("first-enclave-report"),
					Status:        types.DKGRegStatusVerified,
				}
				require.NoError(t, k.setDKGRegistration(ctx, anotherValidator, firstReg))
			},
			// Index is 2 because setupNetwork adds anotherValidator (index 1) in this subtest only.
			expectedRegData: &types.DKGRegistration{
				Round:          testRound,
				ValidatorAddr:  testValidator.Hex(),
				Index:          2,
				DkgPubKey:      []byte("second-dkg-pubkey"),
				CommPubKey:     []byte("second-comm-pubkey"),
				EnclaveReport:  []byte("second-enclave-report"),
				Status:         types.DKGRegStatusVerified,
				CodeCommitment: testCodeCommitment[:],
				EnclaveType:    testEnclaveType[:],
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			k, ctx := setupDKGKeeper(t)
			if tc.setupNetwork != nil {
				tc.setupNetwork(t, k, ctx)
			} else {
				validDKGNetwork := &types.DKGNetwork{
					Round:            tc.round,
					StartBlockHeight: testStartBlockHeight.Int64(),
					StartBlockHash:   testStartBlockHash[:],
					ActiveValSet:     []string{testValidator.Hex()},
					Total:            5,
					Threshold:        3,
					Stage:            types.DKGStageRegistration,
				}
				require.NoError(t, k.setDKGNetwork(ctx, validDKGNetwork))
			}

			err := k.Registered(ctx, tc.msgSender, tc.codeCommitment, tc.round, tc.startBlockHeight, tc.startBlockHash, tc.enclaveType, tc.dkgPubKey, tc.commPubKey, tc.enclaveReport)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)

				if tc.expectedRegData != nil {
					storedReg, err := k.getDKGRegistration(ctx, tc.round, tc.msgSender)
					require.NoError(t, err)
					require.Equal(t, tc.expectedRegData.Round, storedReg.Round)
					require.Equal(t, tc.expectedRegData.ValidatorAddr, storedReg.ValidatorAddr)
					require.Equal(t, tc.expectedRegData.Index, storedReg.Index)
					require.Equal(t, tc.expectedRegData.DkgPubKey, storedReg.DkgPubKey)
					require.Equal(t, tc.expectedRegData.CommPubKey, storedReg.CommPubKey)
					require.Equal(t, tc.expectedRegData.EnclaveReport, storedReg.EnclaveReport)
					require.Equal(t, tc.expectedRegData.Status, storedReg.Status)
					require.Equal(t, tc.expectedRegData.CodeCommitment, storedReg.CodeCommitment)
					require.Equal(t, tc.expectedRegData.EnclaveType, storedReg.EnclaveType)
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
				sigKey, commPubKey, pRoot := setupFinalizedState(t, k, ctx, testValidator, testRound)
				setVerifiedRegistration(t, k, ctx, testValidator, testRound, 1, commPubKey)
				sig := signFinalizationData(t, sigKey, testCodeCommitment, testRound, pRoot, testGlobalPubKey, testPublicCoeffs, testPubKeyShare)

				return finalizedArgs{k, ctx, testRound, testValidator, testCodeCommitment, pRoot, sig, testGlobalPubKey, testPublicCoeffs, testPubKeyShare}
			},
			postCheck: func(t *testing.T, args finalizedArgs) {
				t.Helper()
				reg, err := args.k.getDKGRegistration(args.ctx, args.round, args.msgSender)
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
					Round: 99,
					Total: 5, Threshold: 3, Stage: types.DKGStageFinalization,
				}
				require.NoError(t, k.setDKGNetwork(ctx, network))

				return finalizedArgs{k: k, ctx: ctx, round: testRound, msgSender: testValidator, codeCommitment: testCodeCommitment, globalPubKey: testGlobalPubKey, publicCoeffs: testPublicCoeffs, pubKeyShare: testPubKeyShare}
			},
			expectedErr: "round mismatch",
		},
		{
			name: "fail: stage not finalization",
			setup: func(t *testing.T) finalizedArgs {
				t.Helper()
				k, ctx := setupDKGKeeper(t)
				network := &types.DKGNetwork{
					Round:        testRound,
					ActiveValSet: []string{testValidator.Hex()},
					Total:        5, Threshold: 3, Stage: types.DKGStageDealing,
				}
				require.NoError(t, k.setDKGNetwork(ctx, network))
				setVerifiedRegistration(t, k, ctx, testValidator, testRound, 1, []byte("comm-key"))
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
					Round:        testRound,
					ActiveValSet: []string{testValidator.Hex()},
					Total:        5, Threshold: 3, Stage: types.DKGStageFinalization,
				}
				require.NoError(t, k.setDKGNetwork(ctx, network))

				// Register testValidator as verified (for participants root validation to pass)
				setVerifiedRegistration(t, k, ctx, testValidator, testRound, 1, []byte("comm-key"))
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
				_, commPubKey, pRoot := setupFinalizedState(t, k, ctx, testValidator, testRound)
				setVerifiedRegistration(t, k, ctx, testValidator, testRound, 1, commPubKey)

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
					Round:        testRound,
					ActiveValSet: []string{testValidator.Hex(), otherValidator.Hex()},
					Total:        2, Threshold: 2, Stage: types.DKGStageFinalization,
				}
				require.NoError(t, k.setDKGNetwork(ctx, network))

				// Both validators are Verified
				setVerifiedRegistration(t, k, ctx, testValidator, testRound, 1, commPubKey)
				setVerifiedRegistration(t, k, ctx, otherValidator, testRound, 2, []byte("other-comm-key-padding-to-64-bytes-1234567890123456789012345678"))

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
				sigKey, commPubKey, pRoot := setupFinalizedState(t, k, ctx, testValidator, testRound)
				// Override threshold to 1 so a single vote triggers global key set
				network := &types.DKGNetwork{
					Round:        testRound,
					ActiveValSet: []string{testValidator.Hex()},
					Total:        1, Threshold: 1, Stage: types.DKGStageFinalization,
				}
				require.NoError(t, k.setDKGNetwork(ctx, network))
				setVerifiedRegistration(t, k, ctx, testValidator, testRound, 1, commPubKey)
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
func setupFinalizedState(t *testing.T, k *Keeper, ctx context.Context, validator common.Address, round uint32) (*ecdsa.PrivateKey, []byte, [32]byte) {
	t.Helper()

	sigKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	commPubKey := crypto.FromECDSAPub(&sigKey.PublicKey)[1:]

	network := &types.DKGNetwork{
		Round:        round,
		ActiveValSet: []string{validator.Hex()},
		Total:        5, Threshold: 3, Stage: types.DKGStageFinalization,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	return sigKey, commPubKey, computeParticipantsRoot(validator)
}

// setVerifiedRegistration creates a DKG registration in Verified status.
func setVerifiedRegistration(t *testing.T, k *Keeper, ctx context.Context, validator common.Address, round uint32, index uint32, commPubKey []byte) {
	t.Helper()

	reg := &types.DKGRegistration{
		Round: round, ValidatorAddr: validator.Hex(), Index: index,
		DkgPubKey: []byte("dkg-key"), CommPubKey: commPubKey,
		EnclaveReport: []byte("enclave-report"), Status: types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, validator, reg))
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
		for range 100 {
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

	round := uint32(7)
	ciphertext := []byte("ciphertext")
	encryptedPartial := []byte("encrypted-partial")
	ephemeralPubKey := []byte("ephemeral-pub-key")
	pubShare := []byte("pub-share")

	validSig := signPartialDecryptionData(t, sigKey, round, ciphertext, encryptedPartial, ephemeralPubKey, pubShare)

	tcs := []struct {
		name             string
		commPubKey       []byte
		round            uint32
		ciphertext       []byte
		encryptedPartial []byte
		ephemeralPubKey  []byte
		pubShare         []byte
		signature        []byte
		expectedErr      string
	}{
		{
			name:             "pass: valid signature",
			commPubKey:       commPubKey,
			round:            round,
			ciphertext:       ciphertext,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        validSig,
		},
		{
			name:             "fail: tampered encryptedPartial",
			commPubKey:       commPubKey,
			round:            round,
			ciphertext:       ciphertext,
			encryptedPartial: []byte("tampered-partial"),
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        validSig,
			expectedErr:      "partial decryption signature address mismatch",
		},
		{
			name:             "fail: tampered ciphertext",
			commPubKey:       commPubKey,
			round:            round,
			ciphertext:       []byte("tampered-ciphertext"),
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        validSig,
			expectedErr:      "partial decryption signature address mismatch",
		},
		{
			name:             "fail: tampered ephemeralPubKey",
			commPubKey:       commPubKey,
			round:            round,
			ciphertext:       ciphertext,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  []byte("tampered-ephemeral"),
			pubShare:         pubShare,
			signature:        validSig,
			expectedErr:      "partial decryption signature address mismatch",
		},
		{
			name:             "fail: wrong commPubKey (64 bytes but different key)",
			commPubKey:       make([]byte, 64),
			round:            round,
			ciphertext:       ciphertext,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        validSig,
			expectedErr:      "partial decryption signature address mismatch",
		},
		{
			name:             "fail: commPubKey too short",
			commPubKey:       []byte("short"),
			round:            round,
			ciphertext:       ciphertext,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        validSig,
			expectedErr:      "invalid commPubKey length",
		},
		{
			name:             "fail: commPubKey 65 bytes (with 0x04 prefix)",
			commPubKey:       append([]byte{0x04}, commPubKey...),
			round:            round,
			ciphertext:       ciphertext,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        validSig,
			expectedErr:      "invalid commPubKey length",
		},
		{
			name:             "fail: invalid signature bytes",
			commPubKey:       commPubKey,
			round:            round,
			ciphertext:       ciphertext,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        []byte("too-short"),
			expectedErr:      "invalid signature length",
		},
		{
			name:             "fail: corrupted 65-byte signature",
			commPubKey:       commPubKey,
			round:            round,
			ciphertext:       ciphertext,
			encryptedPartial: encryptedPartial,
			ephemeralPubKey:  ephemeralPubKey,
			pubShare:         pubShare,
			signature:        make([]byte, 65),
			expectedErr:      "failed to recover public key from signature",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := verifyPartialDecryptionSignature(tc.commPubKey, tc.round, tc.ciphertext, tc.encryptedPartial, tc.ephemeralPubKey, tc.pubShare, tc.signature)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func signPartialDecryptionData(t *testing.T, key *ecdsa.PrivateKey, round uint32, ciphertext []byte, encryptedPartial, ephemeralPubKey, pubShare []byte) []byte {
	t.Helper()

	roundBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(roundBytes, round)

	encoded := make([]byte, 0, 4+len(ciphertext)+len(encryptedPartial)+len(ephemeralPubKey)+len(pubShare))
	encoded = append(encoded, roundBytes...)
	encoded = append(encoded, ciphertext...)
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
	tcs := []struct {
		name             string
		activationHeight int64
		upgradeVersion   string
		setupExisting    bool
		expectedErr      string
	}{
		{
			name:             "pass: successful upgrade scheduling",
			activationHeight: 1000,
			upgradeVersion:   "v1.0.0",
		},
		{
			name:             "fail: empty upgrade version",
			activationHeight: 1000,
			upgradeVersion:   "",
			expectedErr:      "upgrade version cannot be empty",
		},
		{
			name:             "fail: pending upgrade already exists",
			activationHeight: 2000,
			upgradeVersion:   "v2.0.0",
			setupExisting:    true,
			expectedErr:      "pending upgrade already exists",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			k, ctx := setupDKGKeeper(t)

			if tc.setupExisting {
				// Schedule an existing upgrade first
				require.NoError(t, k.UpgradeScheduled(ctx, 1000, "v1.0.0"))
			}

			err := k.UpgradeScheduled(ctx, tc.activationHeight, tc.upgradeVersion)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)

				// Verify the upgrade info is stored
				info, err := k.GetPendingUpgrade(ctx)
				require.NoError(t, err)
				require.NotNil(t, info)
				require.Equal(t, tc.upgradeVersion, info.UpgradeVersion)
				require.Equal(t, tc.activationHeight, info.ActivationHeight)
			}
		})
	}
}

func TestKeeper_UpgradeCancelled(t *testing.T) {
	tcs := []struct {
		name           string
		upgradeVersion string
		setupUpgrade   bool
		expectedErr    string
	}{
		{
			name:           "pass: successful upgrade cancellation",
			upgradeVersion: "v1.0.0",
			setupUpgrade:   true,
		},
		{
			name:           "fail: no upgrade found for version",
			upgradeVersion: "v999.0.0",
			setupUpgrade:   false,
			expectedErr:    "kernel upgrade info not found",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			k, ctx := setupDKGKeeper(t)

			if tc.setupUpgrade {
				require.NoError(t, k.UpgradeScheduled(ctx, 1000, tc.upgradeVersion))
			}

			err := k.UpgradeCancelled(ctx, tc.upgradeVersion)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)

				// Verify the upgrade info is deleted
				info, err := k.GetPendingUpgrade(ctx)
				require.NoError(t, err)
				require.Nil(t, info)
			}
		})
	}
}

func TestKeeper_HasPendingUpgradeActivation(t *testing.T) {
	tcs := []struct {
		name          string
		currentHeight int64
		setupUpgrade  *types.KernelUpgradeInfo
		expectNil     bool
	}{
		{
			name:          "no pending upgrade returns nil",
			currentHeight: 100,
			setupUpgrade:  nil,
			expectNil:     true,
		},
		{
			name:          "pending upgrade not yet at activation height returns nil",
			currentHeight: 50,
			setupUpgrade: &types.KernelUpgradeInfo{
				UpgradeVersion:   "v1.0.0",
				ActivationHeight: 100,
			},
			expectNil: true,
		},
		{
			name:          "pending upgrade at exact activation height returns info",
			currentHeight: 100,
			setupUpgrade: &types.KernelUpgradeInfo{
				UpgradeVersion:   "v1.0.0",
				ActivationHeight: 100,
			},
			expectNil: false,
		},
		{
			name:          "pending upgrade past activation height returns info",
			currentHeight: 200,
			setupUpgrade: &types.KernelUpgradeInfo{
				UpgradeVersion:   "v1.0.0",
				ActivationHeight: 100,
			},
			expectNil: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			k, ctx := setupDKGKeeper(t)

			if tc.setupUpgrade != nil {
				require.NoError(t, k.SetKernelUpgradeInfo(ctx, tc.setupUpgrade))
			}

			result, err := k.hasPendingUpgradeActivation(ctx, tc.currentHeight)
			require.NoError(t, err)

			if tc.expectNil {
				require.Nil(t, result)
			} else {
				require.NotNil(t, result)
				require.Equal(t, tc.setupUpgrade.UpgradeVersion, result.UpgradeVersion)
				require.Equal(t, tc.setupUpgrade.ActivationHeight, result.ActivationHeight)
			}
		})
	}
}

func TestKeeper_HasFinalizedRegistration(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	testValidator := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testRound := uint32(1)

	tcs := []struct {
		name           string
		setup          func()
		round          uint32
		validatorAddr  common.Address
		expectedResult bool
		expectedErr    string
	}{
		{
			name: "pass: returns true when registration is finalized",
			setup: func() {
				err := k.setDKGRegistration(ctx, testValidator, &types.DKGRegistration{
					Round:         testRound,
					ValidatorAddr: testValidator.Hex(),
					Index:         1,
					Status:        types.DKGRegStatusFinalized,
				})
				require.NoError(t, err)
			},
			round:          testRound,
			validatorAddr:  testValidator,
			expectedResult: true,
		},
		{
			name: "pass: returns false when registration is verified (not finalized)",
			setup: func() {
				err := k.setDKGRegistration(ctx, testValidator, &types.DKGRegistration{
					Round:         testRound,
					ValidatorAddr: testValidator.Hex(),
					Index:         1,
					Status:        types.DKGRegStatusVerified,
				})
				require.NoError(t, err)
			},
			round:          testRound,
			validatorAddr:  testValidator,
			expectedResult: false,
		},
		{
			name: "pass: returns false when registration is unspecified",
			setup: func() {
				err := k.setDKGRegistration(ctx, testValidator, &types.DKGRegistration{
					Round:         testRound,
					ValidatorAddr: testValidator.Hex(),
					Index:         1,
					Status:        types.DKGRegStatusUnspecified,
				})
				require.NoError(t, err)
			},
			round:          testRound,
			validatorAddr:  testValidator,
			expectedResult: false,
		},
		{
			name:           "pass: returns false when no registration exists",
			setup:          func() {},
			round:          uint32(999), // non-existent round
			validatorAddr:  common.HexToAddress("0x0000000000000000000000000000000000000001"),
			expectedResult: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			result, err := k.HasFinalizedRegistration(ctx, tc.round, tc.validatorAddr)

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

	mockKernelServiceClient := dkgtestutil.NewMockKernelServiceClient(ctrl)

	// Wrap mock TEE client in a KernelRouter for testing
	kernelRouter := NewKernelRouter(nil, nil)
	kernelRouter.RegisterClient([]byte("test"), mockKernelServiceClient)

	k := NewKeeper(
		encCfg.Codec,
		storeService,
		ak,
		bk,
		dk,
		sk,
		valStore,
		kernelRouter,
		nil, // TODO: mock contract client for integration test
		"story1hmjw3pvkjtndpg8wqppwdn8udd835qpan4hm0y",
	)

	require.NoError(t, k.SetParams(testCtx.Ctx, types.DefaultParams()))

	return k, bk, dk, testCtx.Ctx
}
