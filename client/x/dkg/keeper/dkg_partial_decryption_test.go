package keeper

import (
	"context"
	"encoding/json"
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/iavl"
	iavldb "github.com/cosmos/iavl/db"
	ics23 "github.com/cosmos/ics23/go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/netconf"
)

var (
	testValidator1 = common.HexToAddress("0x1111111111111111111111111111111111111111")
	testValidator2 = common.HexToAddress("0x2222222222222222222222222222222222222222")
)

// testLabel returns a fixed 32-byte label for use in partial decryption tests.
func testLabel() []byte {
	label := make([]byte, 32)
	label[0] = 0xAB
	label[1] = 0xCD
	return label
}

// TestSetPartialDecryptionSubmission_Success verifies that a valid partial
// decryption submission is stored without error.
func TestSetPartialDecryptionSubmission_Success(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	err := k.setPartialDecryptionSubmission(
		ctx,
		testValidator1,
		1, // round
		1, // pid
		[]byte("encrypted-partial"),
		[]byte("ephemeral-pub-key"),
		[]byte("pub-share"),
		[]byte("requester-pub-key"),
		testLabel(),
		[]byte("ciphertext-data"),
	)
	require.NoError(t, err)
}

// TestSetPartialDecryptionSubmission_DuplicateRejected verifies that submitting
// the same partial decryption twice returns ErrDuplicatePartialDecryptionSubmission.
func TestSetPartialDecryptionSubmission_DuplicateRejected(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext-data")

	// First submission
	err := k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 1, 1,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		requesterPubKey, label, ciphertext,
	)
	require.NoError(t, err)

	// Second submission with same (validator, round, requester, label, ciphertext)
	err = k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 1, 1,
		[]byte("enc-partial-2"),
		[]byte("eph-key-2"),
		[]byte("pub-share-2"),
		requesterPubKey, label, ciphertext,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDuplicatePartialDecryptionSubmission)
}

// TestSetPartialDecryptionSubmission_DifferentValidators verifies that two
// validators can each submit their own partial decryption for the same request.
func TestSetPartialDecryptionSubmission_DifferentValidators(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext-data")

	// Validator 1 submits
	err := k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 1, 1,
		[]byte("enc-partial-v1"),
		[]byte("eph-key-v1"),
		[]byte("pub-share-v1"),
		requesterPubKey, label, ciphertext,
	)
	require.NoError(t, err)

	// Validator 2 submits (different validator → different key → no duplicate)
	err = k.setPartialDecryptionSubmission(
		ctx,
		testValidator2, 1, 2,
		[]byte("enc-partial-v2"),
		[]byte("eph-key-v2"),
		[]byte("pub-share-v2"),
		requesterPubKey, label, ciphertext,
	)
	require.NoError(t, err)
}

// TestSetPartialDecryptionSubmission_DifferentRounds verifies that the same
// validator can submit for different rounds without triggering duplicate check.
func TestSetPartialDecryptionSubmission_DifferentRounds(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext-data")

	// Round 1
	err := k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 1, 1,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		requesterPubKey, label, ciphertext,
	)
	require.NoError(t, err)

	// Round 2 — same validator, different round
	err = k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 2, 1,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		requesterPubKey, label, ciphertext,
	)
	require.NoError(t, err)
}

// TestDkgPartialDecryptKey_DifferentInputsProduceDifferentKeys verifies that
// the key function produces distinct keys for distinct inputs.
func TestDkgPartialDecryptKey_DifferentInputsProduceDifferentKeys(t *testing.T) {
	t.Parallel()

	label := testLabel()
	ciphertext := []byte("cipher")
	requesterPubKey := []byte("req-pub-key")

	key1 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator1)
	key2 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 2, testValidator1)
	key3 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator2)

	require.NotEqual(t, key1, key2, "different round → different key")
	require.NotEqual(t, key1, key3, "different validator → different key")
}

// TestDkgPartialDecryptKey_SameInputsSameKey verifies that the key function is
// deterministic.
func TestDkgPartialDecryptKey_SameInputsSameKey(t *testing.T) {
	t.Parallel()

	label := testLabel()
	ciphertext := []byte("cipher")
	requesterPubKey := []byte("req-pub-key")

	key1 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator1)
	key2 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator1)
	require.Equal(t, key1, key2, "same inputs must produce same key")
}

// TestDkgPartialDecryptPrefix_ContainsRequesterAndLabel verifies that the
// prefix function produces a string containing the requester hash and label.
func TestDkgPartialDecryptPrefix_ContainsRequesterAndLabel(t *testing.T) {
	t.Parallel()

	requesterPubKey := []byte("req-pub-key")
	label := testLabel()

	prefix := dkgPartialDecryptPrefix(requesterPubKey, label)
	require.NotEmpty(t, prefix)
	require.Contains(t, prefix, "_", "prefix should use underscore as separator")
}

// TestDecodePartialDecryptionSubmission_RoundTrip verifies encode-decode
// round-trip for DKGPartialDecryptionSubmission.
func TestDecodePartialDecryptionSubmission_RoundTrip(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext")

	require.NoError(t, k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 5, 2,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		requesterPubKey, label, ciphertext,
	))

	// Retrieve stored bytes and decode
	key := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 5, testValidator1)
	bz, err := k.DKGPartialDecrypt.Get(ctx, key)
	require.NoError(t, err)

	submission, err := decodePartialDecryptionSubmission(bz)
	require.NoError(t, err)
	require.Equal(t, testValidator1.Hex(), submission.Validator)
	require.Equal(t, uint32(5), submission.Round)
	require.Equal(t, uint32(2), submission.Pid)
	require.Equal(t, []byte("enc-partial"), submission.EncryptedPartial)
}

// TestGetPartialDecryptionsByPrefix verifies that iterating with a prefix
// derived from (requesterPubKey, label) returns all submissions regardless of
// round, ciphertext, or validator, and nothing for a different requester.
func TestGetPartialDecryptionsByPrefix(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	requesterPubKey := []byte("requester-pub-key")
	otherRequesterPubKey := []byte("other-requester")
	label := testLabel()
	ciphertext1 := []byte("ciphertext-1")
	ciphertext2 := []byte("ciphertext-2")

	// Store 3 submissions under the same (requesterPubKey, label):
	//   - validator1, round 1, ciphertext1
	//   - validator2, round 1, ciphertext1
	//   - validator1, round 2, ciphertext2
	submissions := []struct {
		validator  common.Address
		round      uint32
		pid        uint32
		ciphertext []byte
	}{
		{testValidator1, 1, 1, ciphertext1},
		{testValidator2, 1, 2, ciphertext1},
		{testValidator1, 2, 1, ciphertext2},
	}
	for _, s := range submissions {
		require.NoError(t, k.setPartialDecryptionSubmission(
			ctx,
			s.validator, s.round, s.pid,
			[]byte("enc-partial"),
			[]byte("eph-key"),
			[]byte("pub-share"),
			requesterPubKey, label, s.ciphertext,
		))
	}

	// Store one submission for a different requester — must NOT appear in results.
	require.NoError(t, k.setPartialDecryptionSubmission(
		ctx,
		testValidator1, 1, 1,
		[]byte("enc-partial"),
		[]byte("eph-key"),
		[]byte("pub-share"),
		otherRequesterPubKey, label, ciphertext1,
	))

	// Iterate by prefix for (requesterPubKey, label).
	prefix := dkgPartialDecryptPrefix(requesterPubKey, label)
	rangePrefix := (&collections.Range[string]{}).Prefix(prefix)
	iter, err := k.DKGPartialDecrypt.Iterate(ctx, rangePrefix)
	require.NoError(t, err)
	defer iter.Close()

	var got []string
	for ; iter.Valid(); iter.Next() {
		bz, err := iter.Value()
		require.NoError(t, err)
		sub, err := decodePartialDecryptionSubmission(bz)
		require.NoError(t, err)
		got = append(got, sub.Validator+":"+string(sub.Ciphertext))
	}

	// Expect exactly the 3 entries stored under requesterPubKey.
	require.Len(t, got, 3)
	require.Contains(t, got, testValidator1.Hex()+":"+string(ciphertext1))
	require.Contains(t, got, testValidator2.Hex()+":"+string(ciphertext1))
	require.Contains(t, got, testValidator1.Hex()+":"+string(ciphertext2))
}

// TestDecodePartialDecryptionSubmission_InvalidJSON verifies that decoding
// invalid JSON returns an error.
func TestDecodePartialDecryptionSubmission_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := decodePartialDecryptionSubmission([]byte("not-json"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unmarshal partial decryption submission")
}

// writePrimaryEntry writes a raw primary DKGPartialDecrypt entry without touching
// the secondary round index. This reproduces the pre-upgrade state for migration tests.
func writePrimaryEntry(t *testing.T, k *Keeper, ctx context.Context, validator common.Address, round, pid uint32, requesterPubKey, label, ciphertext []byte) string {
	t.Helper()
	key := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, round, validator)
	bz, err := json.Marshal(types.DKGPartialDecryptionSubmission{
		Validator:        validator.Hex(),
		Round:            round,
		Pid:              pid,
		EncryptedPartial: []byte("enc-partial"),
		EphemeralPubKey:  []byte("eph-key"),
		PubShare:         []byte("pub-share"),
		Label:            label,
		Ciphertext:       ciphertext,
	})
	require.NoError(t, err)
	require.NoError(t, k.DKGPartialDecrypt.Set(ctx, key, bz))
	return key
}

// TestPruneOldPartialDecryptions_Basic stores entries for rounds 1, 2, 3 then
// prunes with cutoff=1 and asserts only round 1 is deleted.
func TestPruneOldPartialDecryptions_Basic(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	ctx = sdk.UnwrapSDKContext(ctx).WithChainID(netconf.TestChainID).WithBlockHeight(400)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext")

	for _, round := range []uint32{1, 2, 3} {
		require.NoError(t, k.setPartialDecryptionSubmission(ctx,
			testValidator1, round, 1,
			[]byte("enc-partial"), []byte("eph-key"), []byte("pub-share"),
			requesterPubKey, label, ciphertext,
		))
	}

	require.NoError(t, k.pruneOldPartialDecryptions(ctx, 1, 3))

	key1 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator1)
	exists, err := k.DKGPartialDecrypt.Has(ctx, key1)
	require.NoError(t, err)
	require.False(t, exists, "round 1 should be pruned")

	for _, round := range []uint32{2, 3} {
		key := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, round, testValidator1)
		exists, err := k.DKGPartialDecrypt.Has(ctx, key)
		require.NoError(t, err)
		require.True(t, exists, "round %d should survive", round)
	}
}

// TestPruneOldPartialDecryptions_NothingToDelete stores only round 5 and prunes
// with cutoff=3, expecting no deletions.
func TestPruneOldPartialDecryptions_NothingToDelete(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	ctx = sdk.UnwrapSDKContext(ctx).WithChainID(netconf.TestChainID).WithBlockHeight(400)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext")

	require.NoError(t, k.setPartialDecryptionSubmission(ctx,
		testValidator1, 5, 1,
		[]byte("enc-partial"), []byte("eph-key"), []byte("pub-share"),
		requesterPubKey, label, ciphertext,
	))

	require.NoError(t, k.pruneOldPartialDecryptions(ctx, 3, 5))

	key5 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 5, testValidator1)
	exists, err := k.DKGPartialDecrypt.Has(ctx, key5)
	require.NoError(t, err)
	require.True(t, exists, "round 5 should not be deleted when cutoff is 3")
}

// TestPruneOldPartialDecryptions_MultipleRoundsOnCutoff stores rounds 1–4 and
// prunes with cutoff=2, expecting rounds 1 and 2 deleted, 3 and 4 remaining.
func TestPruneOldPartialDecryptions_MultipleRoundsOnCutoff(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	ctx = sdk.UnwrapSDKContext(ctx).WithChainID(netconf.TestChainID).WithBlockHeight(400)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext")

	for _, round := range []uint32{1, 2, 3, 4} {
		require.NoError(t, k.setPartialDecryptionSubmission(ctx,
			testValidator1, round, 1,
			[]byte("enc-partial"), []byte("eph-key"), []byte("pub-share"),
			requesterPubKey, label, ciphertext,
		))
	}

	require.NoError(t, k.pruneOldPartialDecryptions(ctx, 2, 4))

	for _, round := range []uint32{1, 2} {
		key := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, round, testValidator1)
		exists, err := k.DKGPartialDecrypt.Has(ctx, key)
		require.NoError(t, err)
		require.False(t, exists, "round %d should be pruned", round)
	}
	for _, round := range []uint32{3, 4} {
		key := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, round, testValidator1)
		exists, err := k.DKGPartialDecrypt.Has(ctx, key)
		require.NoError(t, err)
		require.True(t, exists, "round %d should survive", round)
	}
}

// TestPruneOldPartialDecryptions_ZeroCutoff stores round 1 and prunes with
// cutoff=0, expecting no deletions (round 0 doesn't exist in practice).
func TestPruneOldPartialDecryptions_ZeroCutoff(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	ctx = sdk.UnwrapSDKContext(ctx).WithChainID(netconf.TestChainID).WithBlockHeight(400)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext")

	require.NoError(t, k.setPartialDecryptionSubmission(ctx,
		testValidator1, 1, 1,
		[]byte("enc-partial"), []byte("eph-key"), []byte("pub-share"),
		requesterPubKey, label, ciphertext,
	))

	require.NoError(t, k.pruneOldPartialDecryptions(ctx, 0, 2))

	key1 := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 1, testValidator1)
	exists, err := k.DKGPartialDecrypt.Has(ctx, key1)
	require.NoError(t, err)
	require.True(t, exists, "round 1 should not be deleted when cutoff is 0")
}

// TestSetPartialDecryptionSubmission_WritesSecondaryIndex verifies that a
// successful submission also writes an entry in the round secondary index.
func TestSetPartialDecryptionSubmission_WritesSecondaryIndex(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	ctx = sdk.UnwrapSDKContext(ctx).WithChainID(netconf.TestChainID).WithBlockHeight(400)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext")

	require.NoError(t, k.setPartialDecryptionSubmission(ctx,
		testValidator1, 3, 1,
		[]byte("enc-partial"), []byte("eph-key"), []byte("pub-share"),
		requesterPubKey, label, ciphertext,
	))

	primaryKey := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, 3, testValidator1)
	indexKey := dkgPartialDecryptRoundIndexKey(3, primaryKey)
	exists, err := k.DKGPartialDecryptRoundIndex.Has(ctx, indexKey)
	require.NoError(t, err)
	require.True(t, exists, "secondary round index entry should be written alongside primary")
}

// TestSetPartialDecryptionSubmission_DuplicateDoesNotDoubleWriteIndex verifies
// that a rejected duplicate submission does not produce a second secondary index entry.
func TestSetPartialDecryptionSubmission_DuplicateDoesNotDoubleWriteIndex(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	ctx = sdk.UnwrapSDKContext(ctx).WithChainID(netconf.TestChainID).WithBlockHeight(400)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext")

	require.NoError(t, k.setPartialDecryptionSubmission(ctx,
		testValidator1, 3, 1,
		[]byte("enc-partial"), []byte("eph-key"), []byte("pub-share"),
		requesterPubKey, label, ciphertext,
	))

	err := k.setPartialDecryptionSubmission(ctx,
		testValidator1, 3, 1,
		[]byte("enc-partial-2"), []byte("eph-key-2"), []byte("pub-share-2"),
		requesterPubKey, label, ciphertext,
	)
	require.ErrorIs(t, err, ErrDuplicatePartialDecryptionSubmission)

	// There should be exactly 1 secondary index entry across the entire index.
	allIter, err := k.DKGPartialDecryptRoundIndex.Iterate(ctx, nil)
	require.NoError(t, err)
	keys, err := allIter.Keys()
	allIter.Close()
	require.NoError(t, err)
	require.Len(t, keys, 1, "duplicate rejection must not write a second secondary index entry")
}

// TestMigratePartialDecryptRoundIndex_BackfillsAllEntries writes raw primary
// entries for rounds 1, 2, 3 (no secondary index) and verifies that migration
// backfills the secondary index for every entry without deleting any primary data.
// Pruning is left to BeginBlocker on the next FinalizeDKGRound.
func TestMigratePartialDecryptRoundIndex_BackfillsAllEntries(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext")

	key1 := writePrimaryEntry(t, k, ctx, testValidator1, 1, 1, requesterPubKey, label, ciphertext)
	key2 := writePrimaryEntry(t, k, ctx, testValidator1, 2, 1, requesterPubKey, label, ciphertext)
	key3 := writePrimaryEntry(t, k, ctx, testValidator1, 3, 1, requesterPubKey, label, ciphertext)

	require.NoError(t, k.MigratePartialDecryptRoundIndex(ctx))

	// All primary entries must survive — migration must not delete anything.
	for round, key := range map[uint32]string{1: key1, 2: key2, 3: key3} {
		exists, err := k.DKGPartialDecrypt.Has(ctx, key)
		require.NoError(t, err)
		require.True(t, exists, "primary entry for round %d must not be deleted by migration", round)
	}

	// Every primary entry must have a corresponding secondary index entry.
	for round, key := range map[uint32]string{1: key1, 2: key2, 3: key3} {
		indexKey := dkgPartialDecryptRoundIndexKey(round, key)
		exists, err := k.DKGPartialDecryptRoundIndex.Has(ctx, indexKey)
		require.NoError(t, err)
		require.True(t, exists, "secondary index entry missing for round %d", round)
	}

	allIter, err := k.DKGPartialDecryptRoundIndex.Iterate(ctx, nil)
	require.NoError(t, err)
	indexKeys, err := allIter.Keys()
	allIter.Close()
	require.NoError(t, err)
	require.Len(t, indexKeys, 3, "secondary index should have exactly 3 entries")
}

// TestMigratePartialDecryptRoundIndex_Idempotent verifies that running migration
// twice produces no error and leaves state identical after the second run.
func TestMigratePartialDecryptRoundIndex_Idempotent(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	requesterPubKey := []byte("requester-pub-key")
	label := testLabel()
	ciphertext := []byte("ciphertext")

	key1 := writePrimaryEntry(t, k, ctx, testValidator1, 1, 1, requesterPubKey, label, ciphertext)
	key3 := writePrimaryEntry(t, k, ctx, testValidator1, 3, 1, requesterPubKey, label, ciphertext)

	require.NoError(t, k.MigratePartialDecryptRoundIndex(ctx))
	require.NoError(t, k.MigratePartialDecryptRoundIndex(ctx), "second migration run must be idempotent")

	// Both primary entries must still exist after two migration runs.
	for round, key := range map[uint32]string{1: key1, 3: key3} {
		exists, err := k.DKGPartialDecrypt.Has(ctx, key)
		require.NoError(t, err)
		require.True(t, exists, "primary entry for round %d must survive idempotent migration", round)
	}

	// Secondary index should have exactly 2 entries (one per primary), not 4.
	allIter, err := k.DKGPartialDecryptRoundIndex.Iterate(ctx, nil)
	require.NoError(t, err)
	indexKeys, err := allIter.Keys()
	allIter.Close()
	require.NoError(t, err)
	require.Len(t, indexKeys, 2, "idempotent migration must not duplicate secondary index entries")
}

// TestPartialDecryptRoundIndexCollection_NoHashImpact proves that declaring the
// DKGPartialDecryptRoundIndex Map on the schema does not change the committed app hash on
// its own: the hash only changes once something is actually written to the prefix, and
// returns to the original once the entry is pruned. The index is unreleased, so this also
// confirms that nodes which never write the index commit the same hash whether or not the
// new collection is declared.
func TestPartialDecryptRoundIndexCollection_NoHashImpact(t *testing.T) {
	cdc := moduletestutil.MakeTestEncodingConfig().Codec

	// commitHash builds a fresh store, registers the DKGNetworks map and (optionally) the
	// round-index Map, writes a fixed DKGNetwork, optionally writes/prunes an index entry,
	// commits, and returns the resulting multistore commit hash.
	commitHash := func(t *testing.T, registerIndex, writeIndex, pruneIndex bool) []byte {
		t.Helper()

		key := storetypes.NewKVStoreKey(types.StoreKey)
		tkey := storetypes.NewTransientStoreKey("transient_test")
		testCtx := testutil.DefaultContextWithDB(t, key, tkey)
		sb := collections.NewSchemaBuilder(runtime.NewKVStoreService(key))

		networks := collections.NewMap(sb, types.DKGNetworkKey, "dkg_networks",
			collections.StringKey, codec.CollValue[types.DKGNetwork](cdc))

		var index collections.Map[string, bool]
		if registerIndex {
			index = collections.NewMap(sb, types.DKGPartialDecryptRoundIndexKey, "dkg_partial_decrypt_round_index",
				collections.StringKey, collections.BoolValue)
		}

		_, err := sb.Build()
		require.NoError(t, err)

		ctx := testCtx.Ctx

		// Identical "real" state written in every case.
		require.NoError(t, networks.Set(ctx, "1", types.DKGNetwork{Round: 1, Total: 3, Threshold: 2}))

		const primaryKey = "reqhash_label_cthash_1_validator"
		if writeIndex {
			indexKey := dkgPartialDecryptRoundIndexKey(1, primaryKey)
			require.NoError(t, index.Set(ctx, indexKey, true))
			if pruneIndex {
				require.NoError(t, index.Remove(ctx, indexKey))
			}
		}

		return testCtx.CMS.Commit().Hash
	}

	hWithout := commitHash(t, false, false, false)   // Map not declared at all
	hRegistered := commitHash(t, true, false, false) // declared, never written
	hWritten := commitHash(t, true, true, false)     // declared and written
	hPruned := commitHash(t, true, true, true)       // written then pruned

	require.Equal(t, hWithout, hRegistered,
		"declaring the round-index Map without writing must not change the app hash")
	require.NotEqual(t, hWithout, hWritten,
		"writing an index entry must change the app hash (sanity: the test can detect a real change)")
	require.Equal(t, hWithout, hPruned,
		"writing then pruning an index entry must restore the original app hash")
}

// TestPartialDecryptRoundIndex_BoolICS23Provable proves that the bool value makes a
// DKGPartialDecryptRoundIndex leaf ICS23-provable. It reproduces the original failure mode
// at the IAVL layer: an empty-value leaf fails ics23 leaf-op verification ("leaf op needs
// value"), which broke any non-membership proof whose absent key has a round-index right
// neighbor. With the bool marker both the membership proof for the leaf and a non-membership
// proof for an absent key bracketed by bool leaves verify against the tree root.
func TestPartialDecryptRoundIndex_BoolICS23Provable(t *testing.T) {
	// indexKey builds the full in-store key for a round-index entry: the collection prefix
	// (0x0c) followed by the StringKey-encoded logical key. ICS23 proves against these raw
	// store keys, so the test must match the on-disk layout.
	indexKey := func(logical string) []byte {
		out := append([]byte{}, types.DKGPartialDecryptRoundIndexKey.Bytes()...)
		return append(out, []byte(logical)...)
	}

	// boolVal is the on-disk value bytes for storing true via collections.BoolValue, a
	// single non-empty byte (0x01) and what keeps the leaf ICS23-provable.
	boolVal, err := collections.BoolValue.Encode(true)
	require.NoError(t, err)
	require.Equal(t, []byte{0x01}, boolVal, "BoolValue must encode true to a single 0x01 byte")

	var (
		loKey = dkgPartialDecryptRoundIndexKey(1, "reqhash_label_cthash_1_validatorA")
		hiKey = dkgPartialDecryptRoundIndexKey(1, "reqhash_label_cthash_1_validatorZ")
		// absentKey sorts strictly between loKey and hiKey so its right neighbor in a
		// non-membership proof is the hiKey round-index leaf.
		absentKey = dkgPartialDecryptRoundIndexKey(1, "reqhash_label_cthash_1_validatorM")
	)

	t.Run("bool leaf verifies", func(t *testing.T) {
		tree := iavl.NewMutableTree(iavldb.NewMemDB(), 0, false, iavl.NewNopLogger())

		_, err := tree.Set(indexKey(loKey), boolVal)
		require.NoError(t, err)
		_, err = tree.Set(indexKey(hiKey), boolVal)
		require.NoError(t, err)
		_, _, err = tree.SaveVersion()
		require.NoError(t, err)

		imm, err := tree.GetImmutable(tree.Version())
		require.NoError(t, err)
		root := imm.Hash()

		// Membership proof for a bool leaf verifies (non-empty value passes LeafOp).
		memProof, err := imm.GetMembershipProof(indexKey(hiKey))
		require.NoError(t, err)
		require.True(t, ics23.VerifyMembership(ics23.IavlSpec, root, memProof, indexKey(hiKey), boolVal),
			"bool-valued round-index leaf must be ICS23-provable")

		// Non-membership proof for an absent key whose right neighbor is a bool leaf
		// verifies. This is the exact path that failed when leaves had empty values.
		nonMemProof, err := imm.GetNonMembershipProof(indexKey(absentKey))
		require.NoError(t, err)
		require.True(t, ics23.VerifyNonMembership(ics23.IavlSpec, root, nonMemProof, indexKey(absentKey)),
			"non-membership proof bracketed by bool leaves must verify")
	})

	t.Run("empty-value leaf is not ICS23-provable", func(t *testing.T) {
		tree := iavl.NewMutableTree(iavldb.NewMemDB(), 0, false, iavl.NewNopLogger())

		// Empty value reproduces the old []byte{} encoding.
		_, err := tree.Set(indexKey(hiKey), []byte{})
		require.NoError(t, err)
		_, _, err = tree.SaveVersion()
		require.NoError(t, err)

		imm, err := tree.GetImmutable(tree.Version())
		require.NoError(t, err)
		root := imm.Hash()

		memProof, err := imm.GetMembershipProof(indexKey(hiKey))
		require.NoError(t, err)
		// ics23 LeafOp.Apply rejects the empty value, so verification fails.
		require.False(t, ics23.VerifyMembership(ics23.IavlSpec, root, memProof, indexKey(hiKey), []byte{}),
			"empty-value leaf must NOT be ICS23-provable (the bug being fixed)")
	})
}
