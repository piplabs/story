package keeper

import (
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

	// Not found before setting.
	req, found, err := k.getDecryptRequest(ctx, testRequesterPubKey1, label)
	require.NoError(t, err)
	require.False(t, found)
	require.Equal(t, types.DecryptRequest{}, req)

	// Set and retrieve.
	stored := types.DecryptRequest{Round: 7, Ciphertext: []byte("cipher"), Label: label, RequesterPubKey: testRequesterPubKey1, Height: 42}
	require.NoError(t, k.setDecryptRequest(ctx, testRequesterPubKey1, label, stored))
	req, found, err = k.getDecryptRequest(ctx, testRequesterPubKey1, label)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, stored, req)

	// Different requester public key returns not-found.
	_, found, err = k.getDecryptRequest(ctx, testRequesterPubKey2, label)
	require.NoError(t, err)
	require.False(t, found)
}

// TestDeleteDecryptRequest verifies that a single entry can be deleted.
func TestDeleteDecryptRequest(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	label := []byte("label-b")
	require.NoError(t, k.setDecryptRequest(ctx, testRequesterPubKey1, label, types.DecryptRequest{Round: 9, Label: label, Height: 100}))

	// Confirm it exists.
	_, found, err := k.getDecryptRequest(ctx, testRequesterPubKey1, label)
	require.NoError(t, err)
	require.True(t, found)

	// Delete.
	require.NoError(t, k.deleteDecryptRequest(ctx, testRequesterPubKey1, label))

	// Confirm it's gone.
	_, found, err = k.getDecryptRequest(ctx, testRequesterPubKey1, label)
	require.NoError(t, err)
	require.False(t, found)
}

// TestPruneTimedOutDecryptRequests verifies that only expired entries are removed.
func TestPruneTimedOutDecryptRequests(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	const timeout = types.PartialDecryptionTimeoutBlocks

	// Insert entries at various heights.
	entries := []struct {
		label  []byte
		height uint64
	}{
		{[]byte("old-1"), 100},                  // age = 300 > timeout → expired
		{[]byte("old-2"), 50},                   // age = 350 > timeout → expired
		{[]byte("fresh-1"), 300 + timeout - 10}, // age = 10 < timeout → kept
		{[]byte("fresh-2"), 300 + timeout},      // age = 0 → kept (exactly at boundary)
	}
	for _, e := range entries {
		require.NoError(t, k.setDecryptRequest(ctx, testRequesterPubKey1, e.label, types.DecryptRequest{Round: 1, Label: e.label, Height: e.height}))
	}

	currentHeight := uint64(300) + timeout
	require.NoError(t, k.pruneTimedOutDecryptRequests(ctx, currentHeight))

	// Expired entries must be gone.
	for _, label := range [][]byte{[]byte("old-1"), []byte("old-2")} {
		_, found, err := k.getDecryptRequest(ctx, testRequesterPubKey1, label)
		require.NoError(t, err)
		require.False(t, found, "expected expired entry to be pruned: %s", label)
	}

	// Fresh entries must remain.
	for _, label := range [][]byte{[]byte("fresh-1"), []byte("fresh-2")} {
		_, found, err := k.getDecryptRequest(ctx, testRequesterPubKey1, label)
		require.NoError(t, err)
		require.True(t, found, "expected fresh entry to survive pruning: %s", label)
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

	labels := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	for _, l := range labels {
		require.NoError(t, k.setDecryptRequest(ctx, testRequesterPubKey1, l, types.DecryptRequest{Round: 1, Label: l, Height: 1}))
	}

	require.NoError(t, k.pruneTimedOutDecryptRequests(ctx, 9999))

	for _, l := range labels {
		_, found, err := k.getDecryptRequest(ctx, testRequesterPubKey1, l)
		require.NoError(t, err)
		require.False(t, found)
	}
}
