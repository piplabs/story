package keeper

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"slices"
	"strings"
	"testing"
	"time"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
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

// signFinalizationData creates a valid finalization signature for testing using RLP encoding.
func signFinalizationData(t *testing.T, key *ecdsa.PrivateKey, codeCommitment [32]byte, round uint32, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte) []byte {
	t.Helper()

	material := finalizationSignatureMaterial{
		CodeCommitment:   codeCommitment[:],
		Round:            round,
		ParticipantsRoot: participantsRoot,
		GlobalPubKey:     globalPubKey,
		PublicCoeffs:     publicCoeffs,
		PubKeyShare:      pubKeyShare,
	}
	encoded, err := rlp.EncodeToBytes(material)
	require.NoError(t, err)

	msgHash := crypto.Keccak256(encoded)

	sig, err := crypto.Sign(msgHash, key)
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

	material := partialDecryptSignatureMaterial{
		Round:            round,
		Ciphertext:       ciphertext,
		EncryptedPartial: encryptedPartial,
		EphemeralPubKey:  ephemeralPubKey,
		PubShare:         pubShare,
	}
	encoded, err := rlp.EncodeToBytes(material)
	require.NoError(t, err)

	respHash := crypto.Keccak256(encoded)
	sig, err := crypto.Sign(respHash, key)
	require.NoError(t, err)

	sig[64] += 27

	return sig
}

func TestKeeper_PartialDecryptionSubmitted_IgnoredCases(t *testing.T) {
	t.Run("unknown request returns not accepted", func(t *testing.T) {
		k, ctx := setupDKGKeeper(t)
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		accepted, err := k.PartialDecryptionSubmitted(
			sdkCtx,
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
			uint32(1),
			uint32(1),
			[]byte("encrypted-partial"),
			[]byte("ephemeral-pub-key"),
			[]byte("pub-share"),
			[]byte("requester-pub-key"),
			[]byte("ciphertext"),
			[]byte("label"),
			[]byte("signature"),
		)
		require.NoError(t, err)
		require.False(t, accepted)
	})

	t.Run("expired request returns not accepted and deletes registry entry", func(t *testing.T) {
		k, ctx := setupDKGKeeper(t)
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		validator := common.HexToAddress("0x1234567890123456789012345678901234567890")
		round := uint32(2)
		pid := uint32(1)
		requesterPubKey := []byte("requester-pub-key")
		ciphertext := []byte("ciphertext")
		label := []byte("label")

		currentHeight := uint64(1000)
		requestHeight := currentHeight - types.DefaultDecryptTimeout - 1
		sdkCtx = sdkCtx.WithBlockHeight(int64(currentHeight))

		req := types.DecryptRequest{
			Round:           round,
			Ciphertext:      ciphertext,
			Label:           label,
			RequesterPubKey: requesterPubKey,
			Height:          requestHeight,
		}
		require.NoError(t, k.setDecryptRequest(sdkCtx, requesterPubKey, label, req))

		accepted, err := k.PartialDecryptionSubmitted(
			sdkCtx,
			validator,
			round,
			pid,
			[]byte("encrypted-partial"),
			[]byte("ephemeral-pub-key"),
			[]byte("pub-share"),
			requesterPubKey,
			ciphertext,
			label,
			[]byte("signature"),
		)
		require.NoError(t, err)
		require.False(t, accepted)

		_, found, err := k.getDecryptRequest(sdkCtx, requesterPubKey, label, round, ciphertext)
		require.NoError(t, err)
		require.False(t, found)
	})

	t.Run("duplicate submission returns not accepted", func(t *testing.T) {
		k, ctx := setupDKGKeeper(t)
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		validator := common.HexToAddress("0x1234567890123456789012345678901234567890")
		round := uint32(3)
		pid := uint32(1)
		requesterPubKey := []byte("requester-pub-key")
		ciphertext := []byte("ciphertext")
		label := []byte("label")
		encryptedPartial := []byte("encrypted-partial")
		ephemeralPubKey := []byte("ephemeral-pub-key")
		pubShare := []byte("pub-share")

		sigKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		commPubKey := crypto.FromECDSAPub(&sigKey.PublicKey)[1:]

		req := types.DecryptRequest{
			Round:           round,
			Ciphertext:      ciphertext,
			Label:           label,
			RequesterPubKey: requesterPubKey,
			Height:          uint64(sdkCtx.BlockHeight()),
		}
		require.NoError(t, k.setDecryptRequest(sdkCtx, requesterPubKey, label, req))

		reg := &types.DKGRegistration{
			Round:         round,
			ValidatorAddr: validator.Hex(),
			Index:         1,
			DkgPubKey:     []byte("dkg-pub-key"),
			CommPubKey:    commPubKey,
			PubKeyShare:   pubShare,
			Status:        types.DKGRegStatusVerified,
		}
		require.NoError(t, k.setDKGRegistration(sdkCtx, validator, reg))

		signature := signPartialDecryptionData(t, sigKey, round, ciphertext, encryptedPartial, ephemeralPubKey, pubShare)

		accepted, err := k.PartialDecryptionSubmitted(
			sdkCtx,
			validator,
			round,
			pid,
			encryptedPartial,
			ephemeralPubKey,
			pubShare,
			requesterPubKey,
			ciphertext,
			label,
			signature,
		)
		require.NoError(t, err)
		require.True(t, accepted)

		accepted, err = k.PartialDecryptionSubmitted(
			sdkCtx,
			validator,
			round,
			pid,
			encryptedPartial,
			ephemeralPubKey,
			pubShare,
			requesterPubKey,
			ciphertext,
			label,
			signature,
		)
		require.NoError(t, err)
		require.False(t, accepted)
	})
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
			upgradeVersion:   "v1.6.0",
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
	mockContractClient := dkgtestutil.NewMockDKGContractClient(ctrl)

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
		mockContractClient,
		"story1hmjw3pvkjtndpg8wqppwdn8udd835qpan4hm0y",
	)
	_ = mockContractClient // available for tests that need to set expectations

	require.NoError(t, k.SetParams(testCtx.Ctx, types.DefaultParams()))

	return k, bk, dk, testCtx.Ctx
}

// --- ThresholdDecryptRequested / PartialDecryptionSubmitted ---
// (Tests for dkg_handler.go threshold decryption functions)

// buildValidPartialDecryptSignature creates a valid ECDSA signature for a partial
// decryption response. It replicates the signPartialDecryptResponse logic from the
// DKG service using RLP encoding, which computes:
//
//	encoded = RLP(round, ciphertext, encryptedPartial, ephPubKey, pubShare)
//	hash    = Keccak256(encoded)
//	sig     = ECDSA.Sign(privKey, hash)
//
// Returns (commPubKey [64 bytes], signature [65 bytes]).
func buildValidPartialDecryptSignature(t *testing.T, round uint32, ciphertext, encryptedPartial, ephemeralPubKey, pubShare []byte) (commPubKey []byte, signature []byte) {
	t.Helper()

	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	material := partialDecryptSignatureMaterial{
		Round:            round,
		Ciphertext:       ciphertext,
		EncryptedPartial: encryptedPartial,
		EphemeralPubKey:  ephemeralPubKey,
		PubShare:         pubShare,
	}
	encoded, err := rlp.EncodeToBytes(material)
	require.NoError(t, err)

	hash := crypto.Keccak256(encoded)

	sig, err := crypto.Sign(hash, privKey)
	require.NoError(t, err)

	// commPubKey is the uncompressed public key without the 0x04 prefix (64 bytes)
	pub := privKey.PublicKey
	commPubKey = make([]byte, 64)
	copy(commPubKey[:32], pub.X.Bytes())
	copy(commPubKey[32:], pub.Y.Bytes())

	return commPubKey, sig
}

// TestThresholdDecryptRequested_DKGSvcDisabled verifies that when DKG service is
// disabled, ThresholdDecryptRequested stores the decrypt request but returns nil.
func TestThresholdDecryptRequested_DKGSvcDisabled(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	// isDKGSvcEnabled is false by default

	requesterPubKey := []byte("req-pub-key")
	ciphertext := []byte("ciphertext")
	label := []byte("label")

	err := k.ThresholdDecryptRequested(ctx, 1, requesterPubKey, ciphertext, label, 100)
	require.NoError(t, err)

	// Verify the decrypt request was stored
	req, found, err := k.getDecryptRequest(ctx, requesterPubKey, label, 1, ciphertext)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(1), req.Round)
	require.Equal(t, uint64(100), req.Height)
}

// TestThresholdDecryptRequested_RoundNotFound verifies that when DKG service is
// enabled but the DKG network for the round does not exist, an error is returned.
func TestThresholdDecryptRequested_RoundNotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	k.setValidatorAddress(common.HexToAddress("0x1111111111111111111111111111111111111111"))

	requesterPubKey := []byte("req-pub-key")
	ciphertext := []byte("ciphertext")
	label := []byte("label")

	// Round 999 does not exist
	err := k.ThresholdDecryptRequested(ctx, 999, requesterPubKey, ciphertext, label, 100)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get dkg network for decrypt request")
}

// TestThresholdDecryptRequested_RoundNotActive verifies that when the DKG round
// is not in Active stage, the function returns nil without error.
func TestThresholdDecryptRequested_RoundNotActive(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	validatorAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	k.setValidatorAddress(validatorAddr)

	// Set up a non-active network
	network := &types.DKGNetwork{
		Round:        5,
		Total:        3,
		Threshold:    2,
		Stage:        types.DKGStageDealing, // not Active
		ActiveValSet: []string{validatorAddr.Hex()},
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	requesterPubKey := []byte("req-pub-key")
	err := k.ThresholdDecryptRequested(ctx, 5, requesterPubKey, []byte("cipher"), []byte("label"), 100)
	require.NoError(t, err, "non-active stage should return nil without error")
}

// TestThresholdDecryptRequested_ValidatorNotInCommittee verifies that when the
// validator is not in the active val set, the function returns nil.
func TestThresholdDecryptRequested_ValidatorNotInCommittee(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	outsiderAddr := common.HexToAddress("0x9999999999999999999999999999999999999999")
	k.setValidatorAddress(outsiderAddr)

	network := &types.DKGNetwork{
		Round:        6,
		Total:        3,
		Threshold:    2,
		Stage:        types.DKGStageActive,
		ActiveValSet: []string{"0x1111111111111111111111111111111111111111"}, // outsider not in set
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	err := k.ThresholdDecryptRequested(ctx, 6, []byte("req-key"), []byte("cipher"), []byte("label"), 100)
	require.NoError(t, err, "validator not in committee should return nil")
}

// TestThresholdDecryptRequested_ValidatorInCommitteeSessionNotFound verifies that when
// the validator is in the active committee and DKG service is enabled but no session
// exists for the round, ThresholdDecryptRequested returns nil (session error is handled
// asynchronously in the goroutine and logged, not returned to the caller).
func TestThresholdDecryptRequested_ValidatorInCommitteeSessionNotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	validatorAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	k.setValidatorAddress(validatorAddr)
	initTestStateManager(t, k)

	network := &types.DKGNetwork{
		Round:        7,
		Total:        3,
		Threshold:    2,
		Stage:        types.DKGStageActive,
		ActiveValSet: []string{strings.ToLower(validatorAddr.Hex())},
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	// No session created for round 7 → handleThresholdDecryptRequest logs error asynchronously.
	// ThresholdDecryptRequested itself returns nil because the async goroutine handles the error.
	err := k.ThresholdDecryptRequested(ctx, 7, []byte("req-key"), []byte("cipher"), []byte("label"), 100)
	require.NoError(t, err, "async dispatch should not return error to caller")
}

// TestThresholdDecryptRequested_ValidatorInCommitteeWithSession verifies the happy
// path: when DKG service is enabled, validator is in committee, and a session exists,
// the request is added to the session asynchronously.
func TestThresholdDecryptRequested_ValidatorInCommitteeWithSession(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	validatorAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	k.setValidatorAddress(validatorAddr)
	initTestStateManager(t, k)

	network := &types.DKGNetwork{
		Round:        8,
		Total:        3,
		Threshold:    2,
		Stage:        types.DKGStageActive,
		ActiveValSet: []string{strings.ToLower(validatorAddr.Hex())},
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	// Create a session for round 8
	require.NoError(t, k.stateManager.CreateSession(ctx, newTestSession(8)))

	err := k.ThresholdDecryptRequested(ctx, 8, []byte("req-key"), []byte("cipher"), []byte("label"), 100)
	require.NoError(t, err)

	// Wait for the async goroutine to complete
	require.Eventually(t, func() bool {
		session, err := k.stateManager.GetSession(8)
		if err != nil {
			return false
		}

		return len(session.GetDecryptRequests()) == 1
	}, 5*time.Second, 50*time.Millisecond, "session should have one pending decrypt request after async processing")
}

// TestPartialDecryptionSubmitted_RequestNotFound verifies that when the
// decrypt request registry doesn't have a matching entry, the function
// returns nil (request cleaned up or never recorded).
func TestPartialDecryptionSubmitted_RequestNotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	validator := common.HexToAddress("0x1111111111111111111111111111111111111111")
	accepted, err := k.PartialDecryptionSubmitted(
		ctx,
		validator,
		1, // round
		1, // pid
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		[]byte("unknown-req-key"), // request not in registry
		[]byte("ciphertext"),
		[]byte("label"),
		[]byte("signature"),
	)
	require.NoError(t, err, "unknown request should be silently ignored")
	require.False(t, accepted)
}

// TestPartialDecryptionSubmitted_CiphertextMismatch verifies that a ciphertext
// mismatch between the submission and the stored request returns an error.
// Note: round is part of the storage key, so we must store the request with
// the same round used in the submission lookup.
func TestPartialDecryptionSubmitted_CiphertextMismatch(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	requesterPubKey := []byte("req-pub-key-cipher-mismatch")
	ciphertext := []byte("cipher-correct")
	label := []byte("label-cipher-mismatch")

	// Store with round=1 and correct ciphertext
	require.NoError(t, k.setDecryptRequest(ctx, requesterPubKey, label, types.DecryptRequest{
		Round:           1,
		Ciphertext:      ciphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
		Height:          1,
	}))

	validator := common.HexToAddress("0x1111111111111111111111111111111111111111")

	// Submit with round=1 but WRONG ciphertext — the key lookup will fail (not found)
	// because the key includes the ciphertext hash. So this tests the not-found path.
	accepted, err := k.PartialDecryptionSubmitted(
		ctx,
		validator,
		1,
		1,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		requesterPubKey,
		[]byte("different-ciphertext"), // different ciphertext → key mismatch
		label,
		[]byte("signature"),
	)
	// Not found → silently ignored (nil), not an error
	require.NoError(t, err)
	require.False(t, accepted)
}

// TestPartialDecryptionSubmitted_NoRegistration verifies that when a matching
// request is found but the validator has no DKG registration, an error is returned.
func TestPartialDecryptionSubmitted_NoRegistration(t *testing.T) {
	t.Parallel()

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	// Use current_height=1, request at height=0 → 1-0=1 ≤ 200 → no timeout
	sdkCtx := sdk.UnwrapSDKContext(baseCtx).WithBlockHeight(1)

	requesterPubKey := []byte("req-pub-key-no-reg")
	ciphertext := []byte("cipher-no-reg")
	label := []byte("label-no-reg")

	require.NoError(t, k.setDecryptRequest(sdkCtx, requesterPubKey, label, types.DecryptRequest{
		Round:           1,
		Ciphertext:      ciphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
		Height:          0, // height=0, current=1 → 1-0=1 ≤ 200
	}))

	validator := common.HexToAddress("0x1111111111111111111111111111111111111111")
	// No DKG registration for this validator/round → should fail

	accepted, err := k.PartialDecryptionSubmitted(
		sdkCtx,
		validator,
		1,
		1,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		requesterPubKey,
		ciphertext,
		label,
		[]byte("signature"),
	)
	require.Error(t, err, "should fail when no DKG registration exists")
	require.Contains(t, err.Error(), "failed to get DKG registration")
	require.False(t, accepted)
}

// TestPartialDecryptionSubmitted_PubShareMismatch verifies that when the submitted
// pubShare does not match the stored registration pubKeyShare, an error is returned.
func TestPartialDecryptionSubmitted_PubShareMismatch(t *testing.T) {
	t.Parallel()

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(baseCtx).WithBlockHeight(1)

	requesterPubKey := []byte("req-pub-key-pubshare-mismatch")
	ciphertext := []byte("cipher-pubshare-mismatch")
	label := []byte("label-pubshare-mismatch")

	require.NoError(t, k.setDecryptRequest(sdkCtx, requesterPubKey, label, types.DecryptRequest{
		Round:           1,
		Ciphertext:      ciphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
		Height:          0,
	}))

	validator := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	// Register the validator with a specific pubKeyShare
	require.NoError(t, k.setDKGRegistration(sdkCtx, validator, &types.DKGRegistration{
		Round:         1,
		ValidatorAddr: validator.Hex(),
		Index:         1,
		DkgPubKey:     []byte("dkg-pub"),
		CommPubKey:    []byte("comm-pub"),
		PubKeyShare:   []byte("correct-pub-share"),
		Status:        types.DKGRegStatusFinalized,
	}))

	accepted, err := k.PartialDecryptionSubmitted(
		sdkCtx,
		validator,
		1,
		1,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("wrong-pub-share"), // does not match stored "correct-pub-share"
		requesterPubKey,
		ciphertext,
		label,
		[]byte("signature"),
	)
	require.Error(t, err, "pubShare mismatch should return an error")
	require.Contains(t, err.Error(), "pubShare mismatch")
	require.False(t, accepted)
}

// TestPartialDecryptionSubmitted_InvalidSignature verifies that when the ECDSA
// signature verification fails (wrong commPubKey or bad signature), an error is returned.
func TestPartialDecryptionSubmitted_InvalidSignature(t *testing.T) {
	t.Parallel()

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(baseCtx).WithBlockHeight(1)

	requesterPubKey := []byte("req-pub-key-invalid-sig")
	ciphertext := []byte("cipher-invalid-sig")
	label := []byte("label-invalid-sig")

	require.NoError(t, k.setDecryptRequest(sdkCtx, requesterPubKey, label, types.DecryptRequest{
		Round:           2,
		Ciphertext:      ciphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
		Height:          0,
	}))

	validator := common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	pubShare := []byte("my-pub-share")

	// Register with a 64-byte commPubKey (but random/wrong one so signature check fails)
	require.NoError(t, k.setDKGRegistration(sdkCtx, validator, &types.DKGRegistration{
		Round:         2,
		ValidatorAddr: validator.Hex(),
		Index:         1,
		DkgPubKey:     []byte("dkg-pub"),
		CommPubKey:    make([]byte, 64), // 64 zero bytes → valid length but wrong key
		PubKeyShare:   pubShare,
		Status:        types.DKGRegStatusFinalized,
	}))

	// Build a 65-byte signature that is not a valid ECDSA sig
	invalidSig := make([]byte, 65)

	accepted, err := k.PartialDecryptionSubmitted(
		sdkCtx,
		validator,
		2,
		1,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		pubShare,
		requesterPubKey,
		ciphertext,
		label,
		invalidSig,
	)
	require.Error(t, err, "invalid signature should return an error")
	require.Contains(t, err.Error(), "partial decryption signature verification failed")
	require.False(t, accepted)
}

// TestPartialDecryptionSubmitted_TimeoutExceeded verifies that when the current
// block height exceeds the timeout window, the decrypt request is cleaned up
// and PartialDecryptionSubmitted returns nil.
func TestPartialDecryptionSubmitted_TimeoutExceeded(t *testing.T) {
	t.Parallel()

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	// Use the same KV store but advance block height beyond the timeout window.
	// DefaultDecryptTimeout = 200; request stored at height=1,
	// current height=300 → 300-1=299 > 200 → timeout path.
	sdkCtx := sdk.UnwrapSDKContext(baseCtx).WithBlockHeight(300)

	requesterPubKey := []byte("req-pub-key-timeout")
	ciphertext := []byte("cipher-timeout")
	label := []byte("label-timeout")

	require.NoError(t, k.setDecryptRequest(sdkCtx, requesterPubKey, label, types.DecryptRequest{
		Round:           1,
		Ciphertext:      ciphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
		Height:          1, // stored at block 1
	}))

	validator := common.HexToAddress("0x1111111111111111111111111111111111111111")

	accepted, err := k.PartialDecryptionSubmitted(
		sdkCtx,
		validator,
		1,
		1,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		requesterPubKey,
		ciphertext,
		label,
		[]byte("signature"),
	)
	// Timeout exceeded → cleanup and return nil
	require.NoError(t, err)
	require.False(t, accepted)

	// Verify the request was cleaned up
	_, found, err := k.getDecryptRequest(sdkCtx, requesterPubKey, label, 1, ciphertext)
	require.NoError(t, err)
	require.False(t, found, "request should have been cleaned up after timeout")
}

// TestPartialDecryptionSubmitted_Success verifies the happy path: a valid submission
// with correct signature is stored successfully.
func TestPartialDecryptionSubmitted_Success(t *testing.T) {
	t.Parallel()

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(baseCtx).WithBlockHeight(1)

	requesterPubKey := []byte("req-pub-key-success")
	ciphertext := []byte("cipher-success")
	label := []byte("label-success")

	encryptedPartial := []byte("enc-partial-data")
	ephemeralPubKey := []byte("eph-pub-key-data")
	pubShare := []byte("pub-share-data")

	// Build a valid ECDSA signature
	commPubKey, sig := buildValidPartialDecryptSignature(t, 3, ciphertext, encryptedPartial, ephemeralPubKey, pubShare)

	// Store the decrypt request
	require.NoError(t, k.setDecryptRequest(sdkCtx, requesterPubKey, label, types.DecryptRequest{
		Round:           3,
		Ciphertext:      ciphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
		Height:          0, // height=0, current=1 → no timeout
	}))

	validator := common.HexToAddress("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")

	// Register the validator with matching commPubKey and pubKeyShare
	require.NoError(t, k.setDKGRegistration(sdkCtx, validator, &types.DKGRegistration{
		Round:         3,
		ValidatorAddr: validator.Hex(),
		Index:         2,
		DkgPubKey:     []byte("dkg-pub"),
		CommPubKey:    commPubKey,
		PubKeyShare:   pubShare,
		Status:        types.DKGRegStatusFinalized,
	}))

	accepted, err := k.PartialDecryptionSubmitted(
		sdkCtx,
		validator,
		3,
		2,
		encryptedPartial,
		ephemeralPubKey,
		pubShare,
		requesterPubKey,
		ciphertext,
		label,
		sig,
	)
	require.NoError(t, err, "valid submission should succeed")
	require.True(t, accepted)
}

// TestPartialDecryptionSubmitted_DuplicateSubmission verifies that a duplicate
// partial decryption submission (same validator, same request) returns nil (silently
// ignored with a log message), not an error.
func TestPartialDecryptionSubmitted_DuplicateSubmission(t *testing.T) {
	t.Parallel()

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(baseCtx).WithBlockHeight(1)

	requesterPubKey := []byte("req-pub-key-dup")
	ciphertext := []byte("cipher-dup")
	label := []byte("label-dup")

	encryptedPartial := []byte("enc-partial-dup")
	ephemeralPubKey := []byte("eph-pub-key-dup")
	pubShare := []byte("pub-share-dup")

	commPubKey, sig := buildValidPartialDecryptSignature(t, 4, ciphertext, encryptedPartial, ephemeralPubKey, pubShare)

	require.NoError(t, k.setDecryptRequest(sdkCtx, requesterPubKey, label, types.DecryptRequest{
		Round:           4,
		Ciphertext:      ciphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
		Height:          0,
	}))

	validator := common.HexToAddress("0xDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")

	require.NoError(t, k.setDKGRegistration(sdkCtx, validator, &types.DKGRegistration{
		Round:         4,
		ValidatorAddr: validator.Hex(),
		Index:         1,
		DkgPubKey:     []byte("dkg-pub"),
		CommPubKey:    commPubKey,
		PubKeyShare:   pubShare,
		Status:        types.DKGRegStatusFinalized,
	}))

	// First submission — should succeed
	accepted, err := k.PartialDecryptionSubmitted(
		sdkCtx, validator, 4, 1,
		encryptedPartial, ephemeralPubKey, pubShare,
		requesterPubKey, ciphertext, label, sig,
	)
	require.NoError(t, err, "first submission should succeed")
	require.True(t, accepted)

	// Build a second valid signature (same data → same sig is valid)
	commPubKey2, sig2 := buildValidPartialDecryptSignature(t, 4, ciphertext, encryptedPartial, ephemeralPubKey, pubShare)
	// Update registration to use the new commPubKey so signature verification passes
	require.NoError(t, k.setDKGRegistration(sdkCtx, validator, &types.DKGRegistration{
		Round:         4,
		ValidatorAddr: validator.Hex(),
		Index:         1,
		DkgPubKey:     []byte("dkg-pub"),
		CommPubKey:    commPubKey2,
		PubKeyShare:   pubShare,
		Status:        types.DKGRegStatusFinalized,
	}))

	// Second submission — duplicate → silently ignored (returns nil)
	accepted, err = k.PartialDecryptionSubmitted(
		sdkCtx, validator, 4, 1,
		encryptedPartial, ephemeralPubKey, pubShare,
		requesterPubKey, ciphertext, label, sig2,
	)
	require.NoError(t, err, "duplicate submission should be silently ignored")
	require.False(t, accepted)
}

// --- Finalized: invalidated dealer branch (gap 7) ---

// TestFinalized_InvalidatedDealerRejected verifies that Finalized returns an error
// when the submitting validator's registration status is DKGRegStatusInvalidated.
func TestFinalized_InvalidatedDealerRejected(t *testing.T) {
	t.Parallel()

	k, ctx := setupDKGKeeper(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx).WithBlockHeight(100)

	validator := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	round := uint32(1)
	codeCommitment := [32]byte{0xAA}
	participantsRoot := [32]byte{0xBB}

	// Set up network in finalization stage
	network := &types.DKGNetwork{
		Round: round, Total: 3, Threshold: 2, Stage: types.DKGStageFinalization,
		ActiveValSet: []string{validator.Hex()},
	}
	require.NoError(t, k.setDKGNetwork(sdkCtx, network))

	// Register with Invalidated status
	reg := &types.DKGRegistration{
		Round:         round,
		ValidatorAddr: validator.Hex(),
		Index:         1,
		Status:        types.DKGRegStatusInvalidated,
		CommPubKey:    make([]byte, 64),
	}
	require.NoError(t, k.setDKGRegistration(sdkCtx, validator, reg))

	err := k.Finalized(sdkCtx, round, validator, codeCommitment, participantsRoot,
		make([]byte, 65), []byte("global-pub"), [][]byte{[]byte("c1")}, []byte("share"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalidated")
}

// --- validateParticipantsRoot: zero registrations branch (gap 8) ---

// TestValidateParticipantsRoot_NoRegistrations verifies that validateParticipantsRoot
// returns an error when no verified or finalized registrations exist.
func TestValidateParticipantsRoot_NoRegistrations(t *testing.T) {
	t.Parallel()

	k, ctx := setupDKGKeeper(t)

	var root [32]byte
	err := k.validateParticipantsRoot(ctx, 1, root)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no verified or finalized DKG registrations found")
}

// TestValidateParticipantsRoot_HashMismatch verifies that validateParticipantsRoot
// returns an error when the computed hash does not match the provided root.
func TestValidateParticipantsRoot_HashMismatch(t *testing.T) {
	t.Parallel()

	k, ctx := setupDKGKeeper(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	validator := common.HexToAddress("0x1111111111111111111111111111111111111111")
	reg := &types.DKGRegistration{
		Round:         1,
		ValidatorAddr: validator.Hex(),
		Index:         1,
		Status:        types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(sdkCtx, validator, reg))

	// Provide a wrong participants root (all zeros)
	var wrongRoot [32]byte
	err := k.validateParticipantsRoot(sdkCtx, 1, wrongRoot)
	require.Error(t, err)
	require.Contains(t, err.Error(), "participants root mismatch")
}

// --- UpgradeCancelled: info == nil branch (gap 9) ---

// TestUpgradeCancelled_NotFound verifies that UpgradeCancelled returns an error
// when no upgrade info exists for the specified version.
// GetKernelUpgradeInfo wraps collections.ErrNotFound, so the error message contains
// "kernel upgrade info not found" rather than the nil-info branch message.
func TestUpgradeCancelled_NotFound(t *testing.T) {
	t.Parallel()

	k, ctx := setupDKGKeeper(t)

	// No upgrade info stored — GetKernelUpgradeInfo returns a not-found error,
	// which UpgradeCancelled wraps with "failed to get kernel upgrade info".
	err := k.UpgradeCancelled(ctx, "v99.0.0")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get kernel upgrade info")
}

// --- PartialDecryptionSubmitted: round mismatch and ciphertext mismatch branches (gap 10) ---

// TestPartialDecryptionSubmitted_RoundMismatch verifies that PartialDecryptionSubmitted
// returns an error when the submission round does not match the stored request round.
// Note: getDecryptRequest uses (requesterPubKey, label, round, ciphertext) as key,
// so a round mismatch means the lookup returns not-found (silently returns nil).
// The "round mismatch" check occurs AFTER lookup — if found, rounds must match.
// In practice, the key includes round so a different round == not found == nil return.
// We test the not-found (stale/unknown request) path here.
func TestPartialDecryptionSubmitted_UnknownRequest(t *testing.T) {
	t.Parallel()

	k, ctx := setupDKGKeeper(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx).WithBlockHeight(100)

	validator := common.HexToAddress("0x2222222222222222222222222222222222222222")

	// Submit without storing a decrypt request — not found → returns nil (skipped)
	accepted, err := k.PartialDecryptionSubmitted(
		sdkCtx, validator, 99, 1,
		[]byte("enc"), []byte("eph"), []byte("share"),
		[]byte("req-pub"), []byte("cipher"), []byte("label"), make([]byte, 65),
	)
	require.NoError(t, err, "unknown request should be silently ignored (not found path)")
	require.False(t, accepted)
}

// TestPartialDecryptionSubmitted_CiphertextMismatch verifies that PartialDecryptionSubmitted
// returns an error when the ciphertext in the submission does not match the stored request.
// Since getDecryptRequest key includes ciphertext hash, a different ciphertext
// means the request is NOT found (treated as unknown). The ciphertext mismatch
// check is only reachable if two requests exist with same requesterPubKey+label+round
// but different ciphertext — which the current key structure prevents.
// This test documents the not-found path triggered by different ciphertext.
func TestPartialDecryptionSubmitted_DifferentCiphertext_NotFound(t *testing.T) {
	t.Parallel()

	k, ctx := setupDKGKeeper(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx).WithBlockHeight(100)

	requesterPubKey := []byte("req-pub-ct-test")
	label := []byte("label-ct-test")
	storedCiphertext := []byte("stored-cipher")
	differentCiphertext := []byte("different-cipher")
	round := uint32(5)

	// Store a decrypt request with storedCiphertext
	require.NoError(t, k.setDecryptRequest(sdkCtx, requesterPubKey, label, types.DecryptRequest{
		Round:           round,
		Ciphertext:      storedCiphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
		Height:          0,
	}))

	validator := common.HexToAddress("0x3333333333333333333333333333333333333333")

	// Submit with differentCiphertext → lookup uses differentCiphertext in key → not found → nil
	accepted, err := k.PartialDecryptionSubmitted(
		sdkCtx, validator, round, 1,
		[]byte("enc"), []byte("eph"), []byte("share"),
		requesterPubKey, differentCiphertext, label, make([]byte, 65),
	)
	require.NoError(t, err, "different ciphertext triggers not-found path, silently ignored")
	require.False(t, accepted)
}
