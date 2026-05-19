package keeper

import (
	"encoding/json"
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/netconf"
)

// newTestOffChainStore returns a PartialDecryptOffChainStore backed by an in-memory DB.
func newTestOffChainStore(t *testing.T) *PartialDecryptOffChainStore {
	t.Helper()
	return NewPartialDecryptOffChainStore(dbm.NewMemDB())
}

// marshalSubmission marshals a DKGPartialDecryptionSubmission to JSON for direct store writes.
func marshalSubmission(t *testing.T, sub types.DKGPartialDecryptionSubmission) []byte {
	t.Helper()
	bz, err := json.Marshal(sub)
	require.NoError(t, err)
	return bz
}

// TestOffChainStore_SetAndHas verifies that Set writes a primary entry and
// Has returns true for that key, and false for an absent key.
func TestOffChainStore_SetAndHas(t *testing.T) {
	t.Parallel()

	s := newTestOffChainStore(t)

	bz := marshalSubmission(t, types.DKGPartialDecryptionSubmission{
		Validator:  testValidator1.Hex(),
		Round:      1,
		Ciphertext: []byte("ct"),
	})

	require.NoError(t, s.Set("key-round1", 1, bz))

	ok, err := s.Has("key-round1")
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = s.Has("key-absent")
	require.NoError(t, err)
	require.False(t, ok)
}

// TestOffChainStore_PrefixIterator verifies that PrefixIterator returns only
// entries whose key matches the given prefix.
func TestOffChainStore_PrefixIterator(t *testing.T) {
	t.Parallel()

	s := newTestOffChainStore(t)

	requesterPubKey := []byte("req-pub-key")
	label := testLabel()
	ciphertext := []byte("cipher")

	// Write 2 entries matching the prefix and 1 that doesn't.
	key1 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator1)
	key2 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 2, testValidator2)
	otherKey := dkgPartialDecryptKey([]byte("other-req"), label, ciphertext, 1, testValidator1)

	bz := marshalSubmission(t, types.DKGPartialDecryptionSubmission{Validator: testValidator1.Hex(), Round: 1, Ciphertext: ciphertext})
	require.NoError(t, s.Set(key1, 1, bz))
	bz2 := marshalSubmission(t, types.DKGPartialDecryptionSubmission{Validator: testValidator2.Hex(), Round: 2, Ciphertext: ciphertext})
	require.NoError(t, s.Set(key2, 2, bz2))
	otherBz := marshalSubmission(t, types.DKGPartialDecryptionSubmission{Validator: testValidator1.Hex(), Round: 1, Ciphertext: ciphertext})
	require.NoError(t, s.Set(otherKey, 1, otherBz))

	prefix := dkgPartialDecryptPrefix(requesterPubKey, label)
	iter, err := s.PrefixIterator(prefix)
	require.NoError(t, err)
	defer iter.Close()

	var count int
	for ; iter.Valid(); iter.Next() {
		count++
	}
	require.Equal(t, 2, count, "only entries matching the prefix should be returned")
}

// TestOffChainStore_PruneBeforeRound verifies that PruneBeforeRound deletes
// entries for rounds <= cutoffRound and leaves later rounds intact.
func TestOffChainStore_PruneBeforeRound(t *testing.T) {
	t.Parallel()

	s := newTestOffChainStore(t)
	ctx := t.Context()

	requesterPubKey := []byte("req-pub-key")
	label := testLabel()
	ciphertext := []byte("cipher")

	for _, round := range []uint32{1, 2, 3, 4} {
		key := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, round, testValidator1)
		bz := marshalSubmission(t, types.DKGPartialDecryptionSubmission{Validator: testValidator1.Hex(), Round: round, Ciphertext: ciphertext})
		require.NoError(t, s.Set(key, round, bz))
	}

	require.NoError(t, s.PruneBeforeRound(ctx, 2))

	// Rounds 1 and 2 should be gone.
	for _, round := range []uint32{1, 2} {
		key := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, round, testValidator1)
		ok, err := s.Has(key)
		require.NoError(t, err)
		require.False(t, ok, "round %d should have been pruned", round)
	}

	// Rounds 3 and 4 should survive.
	for _, round := range []uint32{3, 4} {
		key := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, round, testValidator1)
		ok, err := s.Has(key)
		require.NoError(t, err)
		require.True(t, ok, "round %d should survive pruning", round)
	}
}

// TestOffChainStore_PruneBeforeRound_NothingToDelete verifies that pruning is
// a no-op when all stored rounds are above the cutoff.
func TestOffChainStore_PruneBeforeRound_NothingToDelete(t *testing.T) {
	t.Parallel()

	s := newTestOffChainStore(t)
	ctx := t.Context()

	requesterPubKey := []byte("req-pub-key")
	label := testLabel()
	ciphertext := []byte("cipher")

	key5 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 5, testValidator1)
	bz := marshalSubmission(t, types.DKGPartialDecryptionSubmission{Validator: testValidator1.Hex(), Round: 5, Ciphertext: ciphertext})
	require.NoError(t, s.Set(key5, 5, bz))

	require.NoError(t, s.PruneBeforeRound(ctx, 3))

	ok, err := s.Has(key5)
	require.NoError(t, err)
	require.True(t, ok, "round 5 should not be deleted when cutoff is 3")
}

// TestPruneOldPartialDecryptions_ArchivesToOffChain verifies that
// pruneOldPartialDecryptions copies primary entries to the off-chain store
// before deleting them from the on-chain IAVL store.
func TestPruneOldPartialDecryptions_ArchivesToOffChain(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	ctx = sdk.UnwrapSDKContext(ctx).WithChainID(netconf.TestChainID).WithBlockHeight(400)

	// Attach an in-memory off-chain store.
	k.offChainPartialDecryptStore = newTestOffChainStore(t)
	k.partialDecryptRetentionRounds = 10

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext")

	// Store rounds 1, 2, 3 on-chain.
	for _, round := range []uint32{1, 2, 3} {
		require.NoError(t, k.setPartialDecryptionSubmission(ctx,
			testValidator1, round, 1,
			[]byte("enc"), []byte("eph"), []byte("ps"),
			requesterPubKey, label, ciphertext,
		))
	}

	// Prune rounds <= 1.
	require.NoError(t, k.pruneOldPartialDecryptions(ctx, 1, 3))

	// Round 1 must be gone from on-chain.
	key1 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator1)
	exists, err := k.DKGPartialDecrypt.Has(ctx, key1)
	require.NoError(t, err)
	require.False(t, exists, "round 1 should be pruned from on-chain")

	// Round 1 must be present in off-chain archive.
	ok, err := k.offChainPartialDecryptStore.Has(key1)
	require.NoError(t, err)
	require.True(t, ok, "round 1 should be archived in off-chain store")

	// Rounds 2 and 3 should NOT be in the off-chain store (not yet pruned).
	for _, round := range []uint32{2, 3} {
		key := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, round, testValidator1)
		ok, err := k.offChainPartialDecryptStore.Has(key)
		require.NoError(t, err)
		require.False(t, ok, "round %d should not yet be in off-chain store", round)
	}
}

// TestGetCDRPartialsHistory_NilRequest verifies that a nil request returns
// InvalidArgument.
func TestGetCDRPartialsHistory_NilRequest(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.offChainPartialDecryptStore = newTestOffChainStore(t)

	_, err := k.GetCDRPartialsHistory(ctx, nil)
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

// TestGetCDRPartialsHistory_StoreDisabled verifies that Unavailable is returned
// when no off-chain store has been configured.
func TestGetCDRPartialsHistory_StoreDisabled(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	// offChainPartialDecryptStore is nil by default.

	_, err := k.GetCDRPartialsHistory(ctx, &types.QueryGetCDRPartialsRequest{
		RequesterPubKeyHex: "aabbccdd",
		Uuid:               1,
	})
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unavailable, s.Code())
}

// TestGetCDRPartialsHistory_InvalidHex verifies that an invalid hex requester
// pub key returns InvalidArgument.
func TestGetCDRPartialsHistory_InvalidHex(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.offChainPartialDecryptStore = newTestOffChainStore(t)

	_, err := k.GetCDRPartialsHistory(ctx, &types.QueryGetCDRPartialsRequest{
		RequesterPubKeyHex: "not-hex",
		Uuid:               1,
	})
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

// TestGetCDRPartialsHistory_NotFound verifies that NotFound is returned when
// the off-chain store has no entries for the given requester+uuid.
func TestGetCDRPartialsHistory_NotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.offChainPartialDecryptStore = newTestOffChainStore(t)

	_, err := k.GetCDRPartialsHistory(ctx, &types.QueryGetCDRPartialsRequest{
		RequesterPubKeyHex: "aabbccdd",
		Uuid:               42,
	})
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, s.Code())
}

// TestGetCDRPartialsHistory_Found verifies that entries archived to the off-chain
// store are returned correctly grouped by round and ciphertext.
func TestGetCDRPartialsHistory_Found(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	ctx = sdk.UnwrapSDKContext(ctx).WithChainID(netconf.TestChainID).WithBlockHeight(400)

	k.offChainPartialDecryptStore = newTestOffChainStore(t)
	k.partialDecryptRetentionRounds = 10

	requesterPubKey := []byte("requester-pub-key-bytes")
	var label [32]byte
	label[31] = 0x07 // uuid=7
	ciphertext := []byte("some-ciphertext")

	val1 := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	val2 := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	// Store round 1 submissions on-chain, then prune to archive them.
	for _, tc := range []struct {
		v   common.Address
		pid uint32
	}{{val1, 1}, {val2, 2}} {
		require.NoError(t, k.setPartialDecryptionSubmission(ctx,
			tc.v, 1, tc.pid,
			[]byte("enc"), []byte("eph"), []byte("ps"),
			requesterPubKey, label[:], ciphertext,
		))
	}
	require.NoError(t, k.pruneOldPartialDecryptions(ctx, 1, 3))

	// Query the off-chain history.
	resp, err := k.GetCDRPartialsHistory(ctx, &types.QueryGetCDRPartialsRequest{
		RequesterPubKeyHex: common.Bytes2Hex(requesterPubKey),
		Uuid:               7,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Submissions, 1, "one group (round 1)")
	require.Equal(t, uint32(1), resp.Submissions[0].Round)
	require.Len(t, resp.Submissions[0].Submissions, 2, "two validators submitted")
}

// TestOffChainStore_RealDB_Persistence verifies end-to-end behaviour with a
// real on-disk goleveldb backend: data written by one store instance is still
// readable after the DB is closed and reopened, and prefix iteration returns
// the correct entries on the reopened instance.
func TestOffChainStore_RealDB_Persistence(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	requesterPubKey := []byte("real-db-requester")
	label := testLabel()
	ciphertext := []byte("persist-cipher")

	key1 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator1)
	key2 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 2, testValidator2)
	bz1 := marshalSubmission(t, types.DKGPartialDecryptionSubmission{Validator: testValidator1.Hex(), Round: 1, Ciphertext: ciphertext})
	bz2 := marshalSubmission(t, types.DKGPartialDecryptionSubmission{Validator: testValidator2.Hex(), Round: 2, Ciphertext: ciphertext})

	// --- write phase ---
	{
		db, err := dbm.NewDB("dkg-partials-test", dbm.GoLevelDBBackend, dir)
		require.NoError(t, err)

		s := NewPartialDecryptOffChainStore(db)
		require.NoError(t, s.Set(key1, 1, bz1))
		require.NoError(t, s.Set(key2, 2, bz2))
		require.NoError(t, db.Close())
	}

	// --- read phase: reopen the same directory ---
	{
		db, err := dbm.NewDB("dkg-partials-test", dbm.GoLevelDBBackend, dir)
		require.NoError(t, err)
		defer db.Close()

		s := NewPartialDecryptOffChainStore(db)

		ok, err := s.Has(key1)
		require.NoError(t, err)
		require.True(t, ok, "key1 must survive DB close/reopen")

		ok, err = s.Has(key2)
		require.NoError(t, err)
		require.True(t, ok, "key2 must survive DB close/reopen")

		// Prefix iterator must still return both entries.
		prefix := dkgPartialDecryptPrefix(requesterPubKey, label)
		iter, err := s.PrefixIterator(prefix)
		require.NoError(t, err)
		defer iter.Close()

		var count int
		for ; iter.Valid(); iter.Next() {
			sub, decErr := decodePartialDecryptionSubmission(iter.Value())
			require.NoError(t, decErr)
			require.NotEmpty(t, sub.Validator)
			count++
		}
		require.Equal(t, 2, count, "both entries must be visible after reopen")
	}
}

// TestOffChainStore_RealDB_PruneBeforeRound verifies that PruneBeforeRound
// permanently removes entries from the on-disk goleveldb (not just from an
// in-memory cache), so they are absent after a DB close/reopen.
func TestOffChainStore_RealDB_PruneBeforeRound(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ctx := t.Context()

	requesterPubKey := []byte("real-db-prune-requester")
	label := testLabel()
	ciphertext := []byte("prune-cipher")

	key1 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator1)
	key3 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 3, testValidator2)

	bz1 := marshalSubmission(t, types.DKGPartialDecryptionSubmission{Validator: testValidator1.Hex(), Round: 1, Ciphertext: ciphertext})
	bz3 := marshalSubmission(t, types.DKGPartialDecryptionSubmission{Validator: testValidator2.Hex(), Round: 3, Ciphertext: ciphertext})

	// Write and prune in the first DB session.
	{
		db, err := dbm.NewDB("dkg-partials-prune-test", dbm.GoLevelDBBackend, dir)
		require.NoError(t, err)

		s := NewPartialDecryptOffChainStore(db)
		require.NoError(t, s.Set(key1, 1, bz1))
		require.NoError(t, s.Set(key3, 3, bz3))
		require.NoError(t, s.PruneBeforeRound(ctx, 2)) // prune round <= 2
		require.NoError(t, db.Close())
	}

	// Reopen and verify the prune is durable.
	{
		db, err := dbm.NewDB("dkg-partials-prune-test", dbm.GoLevelDBBackend, dir)
		require.NoError(t, err)
		defer db.Close()

		s := NewPartialDecryptOffChainStore(db)

		ok, err := s.Has(key1)
		require.NoError(t, err)
		require.False(t, ok, "round 1 must remain pruned after DB reopen")

		ok, err = s.Has(key3)
		require.NoError(t, err)
		require.True(t, ok, "round 3 must survive pruning and DB reopen")
	}
}

// TestGetCDRPartialsHistory_MultipleGroups verifies that submissions from
// different rounds and ciphertexts are grouped and sorted correctly.
func TestGetCDRPartialsHistory_MultipleGroups(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	ctx = sdk.UnwrapSDKContext(ctx).WithChainID(netconf.TestChainID).WithBlockHeight(400)

	k.offChainPartialDecryptStore = newTestOffChainStore(t)
	k.partialDecryptRetentionRounds = 10

	requesterPubKey := []byte("multi-group-requester")
	var label [32]byte
	label[31] = 0x03 // uuid=3
	ct1 := []byte("cipher-alpha")
	ct2 := []byte("cipher-beta")

	val1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	val2 := common.HexToAddress("0x2222222222222222222222222222222222222222")

	// Round 1, ct1
	require.NoError(t, k.setPartialDecryptionSubmission(ctx, val1, 1, 1, []byte("e"), []byte("ep"), []byte("ps"), requesterPubKey, label[:], ct1))
	// Round 2, ct2
	require.NoError(t, k.setPartialDecryptionSubmission(ctx, val2, 2, 1, []byte("e"), []byte("ep"), []byte("ps"), requesterPubKey, label[:], ct2))

	// Archive both rounds.
	require.NoError(t, k.pruneOldPartialDecryptions(ctx, 2, 4))

	resp, err := k.GetCDRPartialsHistory(ctx, &types.QueryGetCDRPartialsRequest{
		RequesterPubKeyHex: common.Bytes2Hex(requesterPubKey),
		Uuid:               3,
	})
	require.NoError(t, err)
	require.Len(t, resp.Submissions, 2, "two groups: round1/ct1 and round2/ct2")

	// Groups sorted by round ascending.
	require.Equal(t, uint32(1), resp.Submissions[0].Round)
	require.Equal(t, uint32(2), resp.Submissions[1].Round)
}
