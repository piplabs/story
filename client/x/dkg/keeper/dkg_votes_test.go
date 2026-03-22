package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAddGlobalPubKeyVote_FirstVote verifies that adding a vote for a new
// (round, globalPubKey, coeffs) tuple starts the count at 1.
func TestAddGlobalPubKeyVote_FirstVote(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	round := uint32(1)
	pubKey := []byte("test-global-pub-key")
	coeffs := [][]byte{[]byte("coeff1"), []byte("coeff2")}

	count, err := k.AddGlobalPubKeyVote(ctx, round, pubKey, coeffs)
	require.NoError(t, err)
	require.Equal(t, uint32(1), count, "first vote should return count=1")
}

// TestAddGlobalPubKeyVote_MultipleVotes verifies that adding multiple votes
// for the same key accumulates correctly.
func TestAddGlobalPubKeyVote_MultipleVotes(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	round := uint32(1)
	pubKey := []byte("global-pub-key-1")
	coeffs := [][]byte{[]byte("c1")}

	for i := uint32(1); i <= 5; i++ {
		count, err := k.AddGlobalPubKeyVote(ctx, round, pubKey, coeffs)
		require.NoError(t, err)
		require.Equal(t, i, count, "vote count should increment with each call")
	}
}

// TestAddGlobalPubKeyVote_DifferentRoundsSameKey verifies that votes for the
// same public key but different rounds are tracked independently.
func TestAddGlobalPubKeyVote_DifferentRoundsSameKey(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	pubKey := []byte("shared-global-pub-key")
	coeffs := [][]byte{[]byte("c")}

	// Round 1: 3 votes
	for i := 0; i < 3; i++ {
		_, err := k.AddGlobalPubKeyVote(ctx, 1, pubKey, coeffs)
		require.NoError(t, err)
	}

	// Round 2: 1 vote — should start fresh at 1
	count, err := k.AddGlobalPubKeyVote(ctx, 2, pubKey, coeffs)
	require.NoError(t, err)
	require.Equal(t, uint32(1), count, "different round should start fresh")
}

// TestAddGlobalPubKeyVote_DifferentKeys verifies that votes for different
// global public keys within the same round are tracked separately.
func TestAddGlobalPubKeyVote_DifferentKeys(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	round := uint32(1)
	coeffs := [][]byte{[]byte("c")}
	pubKey1 := []byte("global-pub-key-alpha")
	pubKey2 := []byte("global-pub-key-beta")

	count1, err := k.AddGlobalPubKeyVote(ctx, round, pubKey1, coeffs)
	require.NoError(t, err)
	require.Equal(t, uint32(1), count1)

	count2, err := k.AddGlobalPubKeyVote(ctx, round, pubKey2, coeffs)
	require.NoError(t, err)
	require.Equal(t, uint32(1), count2, "different key should start fresh count")
}

// TestAddGlobalPubKeyVote_DifferentCoeffsDistinctKeys verifies that votes for
// the same global pubkey but different public coefficients are tracked separately,
// because hashPublicCoeffs produces a distinct key.
func TestAddGlobalPubKeyVote_DifferentCoeffsDistinctKeys(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	round := uint32(1)
	pubKey := []byte("global-pub-key")
	coeffs1 := [][]byte{[]byte("coeff-set-A")}
	coeffs2 := [][]byte{[]byte("coeff-set-B")}

	count1, err := k.AddGlobalPubKeyVote(ctx, round, pubKey, coeffs1)
	require.NoError(t, err)
	require.Equal(t, uint32(1), count1)

	count2, err := k.AddGlobalPubKeyVote(ctx, round, pubKey, coeffs2)
	require.NoError(t, err)
	require.Equal(t, uint32(1), count2, "different coeffs should produce separate vote key")
}

// TestHashPublicCoeffs_Deterministic verifies that hashPublicCoeffs returns
// the same hash for the same input.
func TestHashPublicCoeffs_Deterministic(t *testing.T) {
	t.Parallel()

	coeffs := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	hash1 := hashPublicCoeffs(coeffs)
	hash2 := hashPublicCoeffs(coeffs)
	require.Equal(t, hash1, hash2)
}

// TestHashPublicCoeffs_EmptyInput verifies that hashPublicCoeffs handles
// empty (nil) input without panicking and returns a consistent hash.
func TestHashPublicCoeffs_EmptyInput(t *testing.T) {
	t.Parallel()

	hash := hashPublicCoeffs(nil)
	require.NotEmpty(t, hash)

	// SHA256 of empty input is a known value
	expected := sha256.Sum256(nil)
	require.Equal(t, expected[:], hash)
}

// TestHashPublicCoeffs_OrderMatters verifies that hashPublicCoeffs produces
// different hashes for different orderings of the same coefficients.
func TestHashPublicCoeffs_OrderMatters(t *testing.T) {
	t.Parallel()

	coeffsAB := [][]byte{[]byte("alpha"), []byte("beta")}
	coeffsBA := [][]byte{[]byte("beta"), []byte("alpha")}

	hashAB := hashPublicCoeffs(coeffsAB)
	hashBA := hashPublicCoeffs(coeffsBA)
	require.NotEqual(t, hex.EncodeToString(hashAB), hex.EncodeToString(hashBA),
		"hash should differ when coeff order changes")
}

// TestAddGlobalPubKeyVote_KeyFormat verifies that the storage key format is
// consistent with the expected format (round_pubkeyHex_coeffHash).
func TestAddGlobalPubKeyVote_KeyFormat(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	round := uint32(42)
	pubKey := []byte("test-key")
	coeffs := [][]byte{[]byte("coeff")}

	// Add two votes and verify the count is correctly accumulated
	count1, err := k.AddGlobalPubKeyVote(ctx, round, pubKey, coeffs)
	require.NoError(t, err)
	require.Equal(t, uint32(1), count1)

	count2, err := k.AddGlobalPubKeyVote(ctx, round, pubKey, coeffs)
	require.NoError(t, err)
	require.Equal(t, uint32(2), count2)

	// Verify it is stored under the expected key
	coeffHash := hex.EncodeToString(hashPublicCoeffs(coeffs))
	expectedKey := fmt.Sprintf("%d_%s_%s", round, hex.EncodeToString(pubKey), coeffHash)

	stored, err := k.GlobalPubKeyVotes.Get(ctx, expectedKey)
	require.NoError(t, err)
	require.Equal(t, uint32(2), stored)
}
