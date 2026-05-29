package keeper

import (
	"testing"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// TestDealtDealersCollection_NoHashImpact proves that registering the DealtDealers
// KeySet on the schema (as NewKeeper does) does not change the committed app hash on
// its own: the hash only changes once something is actually written to the prefix, and
// returns to the original once the entry is pruned.
func TestDealtDealersCollection_NoHashImpact(t *testing.T) {
	cdc := moduletestutil.MakeTestEncodingConfig().Codec

	// commitHash builds a fresh store, registers the DKGNetworks map and (optionally) the
	// DealtDealers KeySet, writes a fixed DKGNetwork, optionally writes/prunes a dealt
	// mark, commits, and returns the resulting multistore commit hash.
	commitHash := func(t *testing.T, registerDealt, writeDealt, pruneDealt bool) []byte {
		t.Helper()

		key := storetypes.NewKVStoreKey(types.StoreKey)
		tkey := storetypes.NewTransientStoreKey("transient_test")
		testCtx := testutil.DefaultContextWithDB(t, key, tkey)
		sb := collections.NewSchemaBuilder(runtime.NewKVStoreService(key))

		networks := collections.NewMap(sb, types.DKGNetworkKey, "dkg_networks",
			collections.StringKey, codec.CollValue[types.DKGNetwork](cdc))

		var dealt collections.KeySet[string]
		if registerDealt {
			dealt = collections.NewKeySet(sb, types.DealtDealersKey, "dealt_dealers", collections.StringKey)
		}

		_, err := sb.Build()
		require.NoError(t, err)

		ctx := testCtx.Ctx

		// Identical "real" state written in every case.
		require.NoError(t, networks.Set(ctx, "1", types.DKGNetwork{Round: 1, Total: 3, Threshold: 2}))

		const dealerAddr = "0x0000000000000000000000000000000000000002"
		if writeDealt {
			require.NoError(t, dealt.Set(ctx, dealtDealerKey(1, dealerAddr)))
			if pruneDealt {
				require.NoError(t, dealt.Remove(ctx, dealtDealerKey(1, dealerAddr)))
			}
		}

		return testCtx.CMS.Commit().Hash
	}

	hWithout := commitHash(t, false, false, false)   // KeySet not registered at all
	hRegistered := commitHash(t, true, false, false) // registered, never written
	hWritten := commitHash(t, true, true, false)     // registered and written
	hPruned := commitHash(t, true, true, true)       // written then pruned

	require.Equal(t, hWithout, hRegistered,
		"registering the DealtDealers KeySet without writing must not change the app hash")
	require.NotEqual(t, hWithout, hWritten,
		"writing a dealt mark must change the app hash (sanity: the test can detect a real change)")
	require.Equal(t, hWithout, hPruned,
		"writing then pruning a dealt mark must restore the original app hash")
}
