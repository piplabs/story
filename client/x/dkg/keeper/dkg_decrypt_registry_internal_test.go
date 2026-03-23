package keeper

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

var (
	testRequesterPubKey1 = []byte("requester-pubkey-1")
	testRequesterPubKey2 = []byte("requester-pubkey-2")
)

// TestSetAndGetDecryptRequest verifies round-trip set/get and the not-found path.
func TestSetAndGetDecryptRequest(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	label := []byte("label-a")
	ciphertext := []byte("cipher-a")
	const round = uint32(7)

	// Not found before setting.
	req, found, err := k.getDecryptRequest(ctx, testRequesterPubKey1, label, round, ciphertext)
	require.NoError(t, err)
	require.False(t, found)
	require.Equal(t, types.DecryptRequest{}, req)

	// Set and retrieve.
	stored := types.DecryptRequest{Round: round, Ciphertext: ciphertext, Label: label, RequesterPubKey: testRequesterPubKey1, Height: 42}
	require.NoError(t, k.setDecryptRequest(ctx, testRequesterPubKey1, label, stored))
	req, found, err = k.getDecryptRequest(ctx, testRequesterPubKey1, label, round, ciphertext)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, stored, req)

	// Different requester public key returns not-found.
	_, found, err = k.getDecryptRequest(ctx, testRequesterPubKey2, label, round, ciphertext)
	require.NoError(t, err)
	require.False(t, found)
}

// TestDeleteDecryptRequest verifies that a single entry can be deleted.
func TestDeleteDecryptRequest(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	label := []byte("label-b")
	ciphertext := []byte("cipher-b")
	const round = uint32(9)
	require.NoError(t, k.setDecryptRequest(ctx, testRequesterPubKey1, label, types.DecryptRequest{Round: round, Ciphertext: ciphertext, Label: label, Height: 100}))

	// Confirm it exists.
	_, found, err := k.getDecryptRequest(ctx, testRequesterPubKey1, label, round, ciphertext)
	require.NoError(t, err)
	require.True(t, found)

	// Delete.
	require.NoError(t, k.deleteDecryptRequest(ctx, testRequesterPubKey1, label, round, ciphertext))

	// Confirm it's gone.
	_, found, err = k.getDecryptRequest(ctx, testRequesterPubKey1, label, round, ciphertext)
	require.NoError(t, err)
	require.False(t, found)
}

// TestPruneTimedOutDecryptRequests verifies that only expired entries are removed.
func TestPruneTimedOutDecryptRequests(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	const timeout = types.PartialDecryptionTimeoutBlocks

	// Insert entries at various heights.
	entries := []struct {
		label      []byte
		ciphertext []byte
		height     uint64
		round      uint32
	}{
		{[]byte("old-1"), []byte("cipher-old-1"), 100, 1},                    // age = 300 > timeout → expired
		{[]byte("old-2"), []byte("cipher-old-2"), 50, 1},                     // age = 350 > timeout → expired
		{[]byte("fresh-1"), []byte("cipher-fresh-1"), 300 + timeout - 10, 1}, // age = 10 < timeout → kept
		{[]byte("fresh-2"), []byte("cipher-fresh-2"), 300 + timeout, 1},      // age = 0 → kept (exactly at boundary)
	}
	for _, e := range entries {
		require.NoError(t, k.setDecryptRequest(ctx, testRequesterPubKey1, e.label, types.DecryptRequest{Round: e.round, Ciphertext: e.ciphertext, Label: e.label, Height: e.height}))
	}

	currentHeight := uint64(300) + timeout
	require.NoError(t, k.pruneTimedOutDecryptRequests(ctx, currentHeight))

	// Expired entries must be gone.
	for _, e := range entries[:2] {
		_, found, err := k.getDecryptRequest(ctx, testRequesterPubKey1, e.label, e.round, e.ciphertext)
		require.NoError(t, err)
		require.False(t, found, "expected expired entry to be pruned: %s", e.label)
	}

	// Fresh entries must remain.
	for _, e := range entries[2:] {
		_, found, err := k.getDecryptRequest(ctx, testRequesterPubKey1, e.label, e.round, e.ciphertext)
		require.NoError(t, err)
		require.True(t, found, "expected fresh entry to survive pruning: %s", e.label)
	}
}

// TestPruneTimedOutDecryptRequests_EmptyRegistry verifies pruning an empty registry is a no-op.
func TestPruneTimedOutDecryptRequests_EmptyRegistry(t *testing.T) {
	k, ctx := setupDKGKeeper(t)
	require.NoError(t, k.pruneTimedOutDecryptRequests(ctx, 9999))
}

// TestPruneTimedOutDecryptRequests_AllExpired verifies all entries are removed when all are expired.
func TestPruneTimedOutDecryptRequests_AllExpired(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	entries := []struct {
		label      []byte
		ciphertext []byte
	}{
		{[]byte("a"), []byte("cipher-a")},
		{[]byte("b"), []byte("cipher-b")},
		{[]byte("c"), []byte("cipher-c")},
	}
	for _, e := range entries {
		require.NoError(t, k.setDecryptRequest(ctx, testRequesterPubKey1, e.label, types.DecryptRequest{Round: 1, Ciphertext: e.ciphertext, Label: e.label, Height: 1}))
	}

	require.NoError(t, k.pruneTimedOutDecryptRequests(ctx, 9999))

	for _, e := range entries {
		_, found, err := k.getDecryptRequest(ctx, testRequesterPubKey1, e.label, 1, e.ciphertext)
		require.NoError(t, err)
		require.False(t, found)
	}
}

func TestHasDecryptRequestQuery(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	requesterPubKey := []byte("requester-pubkey")
	label := make([]byte, 32)
	copy(label, []byte("label-32-bytes"))
	ciphertext := []byte("ciphertext")
	const round = uint32(7)

	resp, err := k.HasDecryptRequest(ctx, &types.QueryHasDecryptRequestRequest{
		Round:              round,
		RequesterPubKeyHex: hex.EncodeToString(requesterPubKey),
		LabelHex:           hex.EncodeToString(label),
		CiphertextHex:      hex.EncodeToString(ciphertext),
	})
	require.NoError(t, err)
	require.False(t, resp.Exists)

	stored := types.DecryptRequest{
		Round:           round,
		Ciphertext:      ciphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
		Height:          42,
	}
	require.NoError(t, k.setDecryptRequest(ctx, requesterPubKey, label, stored))

	resp, err = k.HasDecryptRequest(ctx, &types.QueryHasDecryptRequestRequest{
		Round:              round,
		RequesterPubKeyHex: hex.EncodeToString(requesterPubKey),
		LabelHex:           hex.EncodeToString(label),
		CiphertextHex:      hex.EncodeToString(ciphertext),
	})
	require.NoError(t, err)
	require.True(t, resp.Exists)
}

func TestHasDecryptRequestQuery_InvalidLabel(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	resp, err := k.HasDecryptRequest(ctx, &types.QueryHasDecryptRequestRequest{
		Round:              1,
		RequesterPubKeyHex: hex.EncodeToString([]byte("requester")),
		LabelHex:           hex.EncodeToString([]byte("short")),
		CiphertextHex:      hex.EncodeToString([]byte("ciphertext")),
	})
	require.Error(t, err)
	require.Nil(t, resp)
}
