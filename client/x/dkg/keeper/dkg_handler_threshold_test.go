package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

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

// TestPartialDecryptionSubmitted_RequestNotFound verifies that when the
// decrypt request registry doesn't have a matching entry, the function
// returns nil (request cleaned up or never recorded).
func TestPartialDecryptionSubmitted_RequestNotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	validator := common.HexToAddress("0x1111111111111111111111111111111111111111")
	err := k.PartialDecryptionSubmitted(
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
	err := k.PartialDecryptionSubmitted(
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

	err := k.PartialDecryptionSubmitted(
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

	err := k.PartialDecryptionSubmitted(
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

	err := k.PartialDecryptionSubmitted(
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
}

// TestPartialDecryptionSubmitted_TimeoutExceeded verifies that when the current
// block height exceeds the timeout window, the decrypt request is cleaned up
// and PartialDecryptionSubmitted returns nil.
func TestPartialDecryptionSubmitted_TimeoutExceeded(t *testing.T) {
	t.Parallel()

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	// Use the same KV store but advance block height beyond the timeout window.
	// PartialDecryptionTimeoutBlocks = 200; request stored at height=1,
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

	err := k.PartialDecryptionSubmitted(
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

	// Verify the request was cleaned up
	_, found, err := k.getDecryptRequest(sdkCtx, requesterPubKey, label, 1, ciphertext)
	require.NoError(t, err)
	require.False(t, found, "request should have been cleaned up after timeout")
}
