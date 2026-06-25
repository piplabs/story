package keeper

import (
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/cosmos/iavl"
	iavldb "github.com/cosmos/iavl/db"
	ics23 "github.com/cosmos/ics23/go"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// TestDealtDealersCollection_NoHashImpact proves that declaring the DealtDealers Map on
// the schema (as NewKeeper does) does not change the committed app hash on its own: the
// hash only changes once something is actually written to the prefix, and returns to the
// original once the entry is pruned. This is what guarantees consensus safety for chains
// that ran before DealtDealers existed (e.g. v1.7.0 Aeneid).
func TestDealtDealersCollection_NoHashImpact(t *testing.T) {
	cdc := moduletestutil.MakeTestEncodingConfig().Codec

	// commitHash builds a fresh store, registers the DKGNetworks map and (optionally) the
	// DealtDealers Map, writes a fixed DKGNetwork, optionally writes/prunes a dealt mark,
	// commits, and returns the resulting multistore commit hash.
	commitHash := func(t *testing.T, registerDealt, writeDealt, pruneDealt bool) []byte {
		t.Helper()

		key := storetypes.NewKVStoreKey(types.StoreKey)
		tkey := storetypes.NewTransientStoreKey("transient_test")
		testCtx := testutil.DefaultContextWithDB(t, key, tkey)
		sb := collections.NewSchemaBuilder(runtime.NewKVStoreService(key))

		networks := collections.NewMap(sb, types.DKGNetworkKey, "dkg_networks",
			collections.StringKey, codec.CollValue[types.DKGNetwork](cdc))

		var dealt collections.Map[string, bool]
		if registerDealt {
			dealt = collections.NewMap(sb, types.DealtDealersKey, "dealt_dealers",
				collections.StringKey, collections.BoolValue)
		}

		_, err := sb.Build()
		require.NoError(t, err)

		ctx := testCtx.Ctx

		// Identical "real" state written in every case.
		require.NoError(t, networks.Set(ctx, "1", types.DKGNetwork{Round: 1, Total: 3, Threshold: 2}))

		const dealerAddr = "0x0000000000000000000000000000000000000002"
		if writeDealt {
			require.NoError(t, dealt.Set(ctx, dealtDealerKey(1, dealerAddr), true))
			if pruneDealt {
				require.NoError(t, dealt.Remove(ctx, dealtDealerKey(1, dealerAddr)))
			}
		}

		return testCtx.CMS.Commit().Hash
	}

	hWithout := commitHash(t, false, false, false)   // Map not declared at all
	hRegistered := commitHash(t, true, false, false) // declared, never written
	hWritten := commitHash(t, true, true, false)     // declared and written
	hPruned := commitHash(t, true, true, true)       // written then pruned

	require.Equal(t, hWithout, hRegistered,
		"declaring the DealtDealers Map without writing must not change the app hash")
	require.NotEqual(t, hWithout, hWritten,
		"writing a dealt mark must change the app hash (sanity: the test can detect a real change)")
	require.Equal(t, hWithout, hPruned,
		"writing then pruning a dealt mark must restore the original app hash")
}

// TestDealtDealersCollection_DeterministicReplay models a v1.9.0 binary (with the new
// DealtDealers Map declared on the keeper) reprocessing v1.7.0-era history that predates
// DealtDealers. Replaying the same pre-DealtDealers writes under the new code must yield
// the same app hash as the original processing, i.e. merely declaring the new collection
// prefix does not perturb replay of old state.
func TestDealtDealersCollection_DeterministicReplay(t *testing.T) {
	cdc := moduletestutil.MakeTestEncodingConfig().Codec

	// commitOldEraState writes a fixed set of pre-DealtDealers dkg-module keys, NEVER
	// touches DealtDealers, commits, and returns the commit hash. declareDealt controls
	// whether the new Map is declared on the schema (i.e. old binary vs new binary).
	commitOldEraState := func(t *testing.T, declareDealt bool) []byte {
		t.Helper()

		key := storetypes.NewKVStoreKey(types.StoreKey)
		tkey := storetypes.NewTransientStoreKey("transient_test")
		testCtx := testutil.DefaultContextWithDB(t, key, tkey)
		sb := collections.NewSchemaBuilder(runtime.NewKVStoreService(key))

		networks := collections.NewMap(sb, types.DKGNetworkKey, "dkg_networks",
			collections.StringKey, codec.CollValue[types.DKGNetwork](cdc))
		votes := collections.NewMap(sb, types.GlobalPubKeyVotesKey, "dkg_global_pub_key_votes",
			collections.StringKey, collections.Uint32Value)

		// v1.9.0 binary additionally declares the new Map; v1.7.0 binary does not.
		if declareDealt {
			collections.NewMap(sb, types.DealtDealersKey, "dealt_dealers",
				collections.StringKey, collections.BoolValue)
		}

		_, err := sb.Build()
		require.NoError(t, err)

		ctx := testCtx.Ctx

		// Pre-DealtDealers state: networks + votes, identical across both binaries.
		require.NoError(t, networks.Set(ctx, "1", types.DKGNetwork{Round: 1, Total: 3, Threshold: 2}))
		require.NoError(t, networks.Set(ctx, "2", types.DKGNetwork{Round: 2, Total: 5, Threshold: 3}))
		require.NoError(t, votes.Set(ctx, "1_abc", 2))
		require.NoError(t, votes.Set(ctx, "2_def", 4))

		return testCtx.CMS.Commit().Hash
	}

	h1 := commitOldEraState(t, false) // v1.7.0 binary processing v1.7.0-era state
	h2 := commitOldEraState(t, true)  // v1.9.0 binary replaying the same state

	require.Equal(t, h1, h2,
		"a v1.9.0 node declaring DealtDealers must replay pre-DealtDealers history to the same app hash")
}

// TestDealtDealers_BoolICS23Provable proves that the bool value makes a DealtDealers leaf
// ICS23-provable. It reproduces the original failure mode at the IAVL layer: an empty-value
// leaf fails ics23 leaf-op verification ("leaf op needs value"), which broke non-membership
// proofs whose absent key has a DealtDealers right neighbor. With the bool marker both the
// membership proof for the leaf and a non-membership proof for an absent key bracketed by
// bool leaves verify against the tree root.
func TestDealtDealers_BoolICS23Provable(t *testing.T) {
	// dealtKey builds the full in-store key for a DealtDealers entry: the collection
	// prefix (0x0e) followed by the StringKey-encoded logical key. ICS23 proves against
	// these raw store keys, so the test must match the on-disk layout.
	dealtKey := func(logical string) []byte {
		out := append([]byte{}, types.DealtDealersKey.Bytes()...)
		return append(out, []byte(logical)...)
	}

	// boolVal is the on-disk value bytes for storing true via collections.BoolValue, which
	// is a single non-empty byte (0x01) and is what keeps the leaf ICS23-provable.
	boolVal, err := collections.BoolValue.Encode(true)
	require.NoError(t, err)
	require.Equal(t, []byte{0x01}, boolVal, "BoolValue must encode true to a single 0x01 byte")

	const (
		loKey = "1_0x0000000000000000000000000000000000000001"
		hiKey = "1_0x0000000000000000000000000000000000000009"
		// absentKey sorts strictly between loKey and hiKey so its right neighbor in a
		// non-membership proof is the hiKey DealtDealers leaf.
		absentKey = "1_0x0000000000000000000000000000000000000005"
	)

	t.Run("bool leaf verifies", func(t *testing.T) {
		tree := iavl.NewMutableTree(iavldb.NewMemDB(), 0, false, iavl.NewNopLogger())

		_, err := tree.Set(dealtKey(loKey), boolVal)
		require.NoError(t, err)
		_, err = tree.Set(dealtKey(hiKey), boolVal)
		require.NoError(t, err)
		_, _, err = tree.SaveVersion()
		require.NoError(t, err)

		imm, err := tree.GetImmutable(tree.Version())
		require.NoError(t, err)
		root := imm.Hash()

		// Membership proof for a bool leaf verifies (non-empty value passes LeafOp).
		memProof, err := imm.GetMembershipProof(dealtKey(hiKey))
		require.NoError(t, err)
		require.True(t, ics23.VerifyMembership(ics23.IavlSpec, root, memProof, dealtKey(hiKey), boolVal),
			"bool-valued DealtDealers leaf must be ICS23-provable")

		// Non-membership proof for an absent key whose right neighbor is a bool leaf
		// verifies. This is the exact path that failed when leaves had empty values.
		nonMemProof, err := imm.GetNonMembershipProof(dealtKey(absentKey))
		require.NoError(t, err)
		require.True(t, ics23.VerifyNonMembership(ics23.IavlSpec, root, nonMemProof, dealtKey(absentKey)),
			"non-membership proof bracketed by bool leaves must verify")
	})

	t.Run("empty-value leaf is not ICS23-provable", func(t *testing.T) {
		tree := iavl.NewMutableTree(iavldb.NewMemDB(), 0, false, iavl.NewNopLogger())

		// Empty value reproduces the old KeySet encoding.
		_, err := tree.Set(dealtKey(hiKey), []byte{})
		require.NoError(t, err)
		_, _, err = tree.SaveVersion()
		require.NoError(t, err)

		imm, err := tree.GetImmutable(tree.Version())
		require.NoError(t, err)
		root := imm.Hash()

		memProof, err := imm.GetMembershipProof(dealtKey(hiKey))
		require.NoError(t, err)
		// ics23 LeafOp.Apply rejects the empty value, so verification fails.
		require.False(t, ics23.VerifyMembership(ics23.IavlSpec, root, memProof, dealtKey(hiKey), []byte{}),
			"empty-value leaf must NOT be ICS23-provable (the bug being fixed)")
	})
}
