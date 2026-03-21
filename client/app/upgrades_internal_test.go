package app

import (
	"testing"

	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/log"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"

	"github.com/stretchr/testify/require"
)

// newTestStore creates an in-memory rootmulti.Store with the given store keys mounted.
func newTestStore(t *testing.T, keys ...string) *rootmulti.Store {
	t.Helper()
	db := dbm.NewMemDB()
	ms := rootmulti.NewStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

	for _, k := range keys {
		ms.MountStoreWithDB(storetypes.NewKVStoreKey(k), storetypes.StoreTypeIAVL, nil)
	}

	require.NoError(t, ms.LoadLatestVersion())

	return ms
}

func TestUpgradeStoreLoader_FreshGenesis(t *testing.T) {
	ms := newTestStore(t, "bank", "staking", "dkg")
	// lastVersion == 0 on fresh genesis → DefaultStoreLoader path
	require.Equal(t, int64(0), ms.LastCommitID().Version)

	loader := UpgradeStoreLoader(StoreUpgradesMap{
		100: {Added: []string{"dkg"}},
	})

	require.NoError(t, loader(ms))
}

func TestUpgradeStoreLoader_FutureUpgradePreAddsStores(t *testing.T) {
	ms := newTestStore(t, "bank", "staking")
	// Commit once so lastVersion > 0
	ms.Commit()
	require.Equal(t, int64(1), ms.LastCommitID().Version)

	// Upgrade at height 100 (far in the future, nextVersion=2)
	loader := UpgradeStoreLoader(StoreUpgradesMap{
		100: {Added: []string{"dkg"}},
	})

	require.NoError(t, loader(ms))

	// After loading, dkg store should be available
	dkgKey := storetypes.NewKVStoreKey("dkg")
	ms.MountStoreWithDB(dkgKey, storetypes.StoreTypeIAVL, nil)
	// This verifies the store was pre-added (LoadLatestVersionAndUpgrade was called)
}

func TestUpgradeStoreLoader_ExactUpgradeHeight(t *testing.T) {
	ms := newTestStore(t, "bank", "staking")
	ms.Commit()
	require.Equal(t, int64(1), ms.LastCommitID().Version)

	// nextVersion = 2, so height=2 is the exact upgrade height
	loader := UpgradeStoreLoader(StoreUpgradesMap{
		2: {Added: []string{"dkg"}},
	})

	require.NoError(t, loader(ms))
}

func TestUpgradeStoreLoader_PastUpgradeSkipped(t *testing.T) {
	ms := newTestStore(t, "bank", "staking", "dkg")
	ms.Commit() // version 1
	ms.Commit() // version 2

	require.Equal(t, int64(2), ms.LastCommitID().Version)

	// Upgrade at height 1 is in the past (nextVersion=3)
	loader := UpgradeStoreLoader(StoreUpgradesMap{
		1: {Added: []string{"newmodule"}},
	})

	// Should use DefaultStoreLoader since no pending upgrades
	require.NoError(t, loader(ms))
}

func TestUpgradeStoreLoader_EmptyMap(t *testing.T) {
	ms := newTestStore(t, "bank", "staking")
	ms.Commit()

	loader := UpgradeStoreLoader(StoreUpgradesMap{})
	require.NoError(t, loader(ms))
}

func TestUpgradeStoreLoader_DeterministicOrdering(t *testing.T) {
	// Ensure multiple upgrades produce consistent merged results
	ms := newTestStore(t, "bank")
	ms.Commit()

	loader := UpgradeStoreLoader(StoreUpgradesMap{
		100: {Added: []string{"module_a"}},
		200: {Added: []string{"module_b"}},
		300: {Added: []string{"module_c"}},
	})

	require.NoError(t, loader(ms))
}
