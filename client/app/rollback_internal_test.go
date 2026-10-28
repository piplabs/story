package app

// Internal test (package app) so it can call the guard directly.

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"
)

const (
	bankStore = "bank"
	dkgStore  = "dkg"
)

// commitWithStores builds a store that has committed one version holding the
// given stores, then returns a store that mounts only `keep` — the shape a
// binary has after a module's store was retired.
func commitWithStores(t *testing.T, all, keep []string) *rootmulti.Store {
	t.Helper()

	db := dbm.NewMemDB()

	first := rootmulti.NewStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	for _, name := range all {
		first.MountStoreWithDB(storetypes.NewKVStoreKey(name), storetypes.StoreTypeIAVL, nil)
	}

	require.NoError(t, first.LoadLatestVersion())
	first.Commit()

	// A second store over the same database, mounting the surviving set only.
	second := rootmulti.NewStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	for _, name := range keep {
		second.MountStoreWithDB(storetypes.NewKVStoreKey(name), storetypes.StoreTypeIAVL, nil)
	}

	require.NoError(t, second.LoadLatestVersion())
	second.Commit()

	return second
}

func TestRollbackRefusedAcrossARetiredStore(t *testing.T) {
	t.Parallel()

	// Version 1 holds "dkg" and this binary mounts only "bank", so a target of
	// 1 would rewrite that version's commit info without "dkg" and leave a
	// state neither binary can open. Commit further so the target is in range.
	rs := commitWithStores(t, []string{bankStore, dkgStore}, []string{bankStore})
	for range 3 {
		rs.Commit()
	}

	err := assertRollbackTargetIsLoadable(rs, uint64(rs.LastCommitID().Version-1))
	require.Error(t, err, "rollback across the retired store was allowed")
	require.Contains(t, err.Error(), "retired module store")
}

func TestRollbackAllowedWhenNoStoreWasRetired(t *testing.T) {
	t.Parallel()

	rs := commitWithStores(t, []string{bankStore}, []string{bankStore})
	for range 3 {
		rs.Commit()
	}

	require.NoError(t, assertRollbackTargetIsLoadable(rs, 1),
		"an ordinary rollback must not be refused")
}

func TestRollbackRefusedWhenTheCountIsOutOfRange(t *testing.T) {
	t.Parallel()

	// The rollback loop deletes a block per pass and is not atomic, so a count
	// past the current height must be refused rather than left to fail midway
	// with the node's CometBFT state already stripped.
	// Use a store with nothing retired, so the only reason to refuse is range.
	rs := commitWithStores(t, []string{bankStore}, []string{bankStore})

	for _, heights := range []uint64{uint64(rs.LastCommitID().Version), 99, 1000, ^uint64(0)} {
		err := assertRollbackTargetIsLoadable(rs, heights)
		require.Errorf(t, err, "a count of %d past the tip was allowed", heights)
		require.Contains(t, err.Error(), "out of range")
	}
}

func TestRollbackRefusedWhenTheCountIsZero(t *testing.T) {
	t.Parallel()

	rs := commitWithStores(t, []string{bankStore}, []string{bankStore})

	err := assertRollbackTargetIsLoadable(rs, 0)
	require.Error(t, err, "a zero count was allowed")
	require.Contains(t, err.Error(), "out of range")
}
