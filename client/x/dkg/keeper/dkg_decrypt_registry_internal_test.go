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

// TestSetAndGetDecryptRequestHeight verifies round-trip set/get and the not-found path.
func TestSetAndGetDecryptRequestHeight(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	label := []byte("label-a")

	// Not found before setting.
	h, found, err := k.getDecryptRequestHeight(ctx, testRequesterPubKey1, label)
	require.NoError(t, err)
	require.False(t, found)
	require.Zero(t, h)

	// Set and retrieve.
	require.NoError(t, k.setDecryptRequestHeight(ctx, testRequesterPubKey1, label, 42))
	h, found, err = k.getDecryptRequestHeight(ctx, testRequesterPubKey1, label)
	require.NoError(t, err)
	require.True(t, found)
	require.EqualValues(t, 42, h)

	// Different requester public key returns not-found.
	_, found, err = k.getDecryptRequestHeight(ctx, testRequesterPubKey2, label)
	require.NoError(t, err)
	require.False(t, found)
}

// TestDeleteDecryptRequestHeight verifies that a single entry can be deleted.
func TestDeleteDecryptRequestHeight(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	label := []byte("label-b")
	require.NoError(t, k.setDecryptRequestHeight(ctx, testRequesterPubKey1, label, 100))

	// Confirm it exists.
	_, found, err := k.getDecryptRequestHeight(ctx, testRequesterPubKey1, label)
	require.NoError(t, err)
	require.True(t, found)

	// Delete.
	require.NoError(t, k.deleteDecryptRequestHeight(ctx, testRequesterPubKey1, label))

	// Confirm it's gone.
	_, found, err = k.getDecryptRequestHeight(ctx, testRequesterPubKey1, label)
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
		require.NoError(t, k.setDecryptRequestHeight(ctx, testRequesterPubKey1, e.label, e.height))
	}

	currentHeight := uint64(300) + timeout
	require.NoError(t, k.pruneTimedOutDecryptRequests(ctx, currentHeight))

	// Expired entries must be gone.
	for _, label := range [][]byte{[]byte("old-1"), []byte("old-2")} {
		_, found, err := k.getDecryptRequestHeight(ctx, testRequesterPubKey1, label)
		require.NoError(t, err)
		require.False(t, found, "expected expired entry to be pruned: %s", label)
	}

	// Fresh entries must remain.
	for _, label := range [][]byte{[]byte("fresh-1"), []byte("fresh-2")} {
		_, found, err := k.getDecryptRequestHeight(ctx, testRequesterPubKey1, label)
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
		require.NoError(t, k.setDecryptRequestHeight(ctx, testRequesterPubKey1, l, 1))
	}

	require.NoError(t, k.pruneTimedOutDecryptRequests(ctx, 9999))

	for _, l := range labels {
		_, found, err := k.getDecryptRequestHeight(ctx, testRequesterPubKey1, l)
		require.NoError(t, err)
		require.False(t, found)
	}
}
