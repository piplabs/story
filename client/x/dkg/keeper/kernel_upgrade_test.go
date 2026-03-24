package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// TestSetKernelUpgradeInfo_And_Get verifies that SetKernelUpgradeInfo stores
// the info and GetKernelUpgradeInfo retrieves it correctly.
func TestSetKernelUpgradeInfo_And_Get(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	info := &types.KernelUpgradeInfo{
		UpgradeVersion:   "v2.0.0",
		ActivationHeight: 500,
		IsActivated:      false,
	}

	require.NoError(t, k.SetKernelUpgradeInfo(ctx, info))

	got, err := k.GetKernelUpgradeInfo(ctx, "v2.0.0")
	require.NoError(t, err)
	require.Equal(t, "v2.0.0", got.UpgradeVersion)
	require.Equal(t, int64(500), got.ActivationHeight)
	require.False(t, got.IsActivated)
}

// TestGetKernelUpgradeInfo_NotFound verifies that GetKernelUpgradeInfo returns
// an error when the version does not exist.
func TestGetKernelUpgradeInfo_NotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetKernelUpgradeInfo(ctx, "nonexistent")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

// TestSetKernelUpgradeInfos_Multiple verifies that SetKernelUpgradeInfos stores
// all provided infos and each is retrievable.
func TestSetKernelUpgradeInfos_Multiple(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	infos := []types.KernelUpgradeInfo{
		{UpgradeVersion: "v1.0.0", ActivationHeight: 100, IsActivated: false},
		{UpgradeVersion: "v2.0.0", ActivationHeight: 200, IsActivated: true},
		{UpgradeVersion: "v3.0.0", ActivationHeight: 300, IsActivated: false},
	}

	require.NoError(t, k.SetKernelUpgradeInfos(ctx, infos))

	for _, info := range infos {
		got, err := k.GetKernelUpgradeInfo(ctx, info.UpgradeVersion)
		require.NoError(t, err)
		require.Equal(t, info.ActivationHeight, got.ActivationHeight)
		require.Equal(t, info.IsActivated, got.IsActivated)
	}
}

// TestSetKernelUpgradeInfos_Empty verifies that SetKernelUpgradeInfos with
// empty input succeeds without error.
func TestSetKernelUpgradeInfos_Empty(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	require.NoError(t, k.SetKernelUpgradeInfos(ctx, nil))
}

// TestGetAllKernelUpgradeInfos_Empty verifies that GetAllKernelUpgradeInfos
// returns an empty slice when no infos are stored.
func TestGetAllKernelUpgradeInfos_Empty(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	infos, err := k.GetAllKernelUpgradeInfos(ctx)
	require.NoError(t, err)
	require.Empty(t, infos)
}

// TestGetAllKernelUpgradeInfos_Multiple verifies that GetAllKernelUpgradeInfos
// returns all stored infos.
func TestGetAllKernelUpgradeInfos_Multiple(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	infos := []types.KernelUpgradeInfo{
		{UpgradeVersion: "v1.0.0", ActivationHeight: 100},
		{UpgradeVersion: "v2.0.0", ActivationHeight: 200},
	}
	require.NoError(t, k.SetKernelUpgradeInfos(ctx, infos))

	all, err := k.GetAllKernelUpgradeInfos(ctx)
	require.NoError(t, err)
	require.Len(t, all, 2)
}

// TestGetPendingUpgrade_NoPending verifies that GetPendingUpgrade returns nil
// when all entries are activated or none exist.
func TestGetPendingUpgrade_NoPending(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// No entries at all
	pending, err := k.GetPendingUpgrade(ctx)
	require.NoError(t, err)
	require.Nil(t, pending)
}

// TestGetPendingUpgrade_AllActivated verifies that GetPendingUpgrade returns nil
// when all stored entries are already activated.
func TestGetPendingUpgrade_AllActivated(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	require.NoError(t, k.SetKernelUpgradeInfo(ctx, &types.KernelUpgradeInfo{
		UpgradeVersion: "v1.0.0", IsActivated: true,
	}))

	pending, err := k.GetPendingUpgrade(ctx)
	require.NoError(t, err)
	require.Nil(t, pending, "activated entry should not be returned as pending")
}

// TestGetPendingUpgrade_ReturnsPending verifies that GetPendingUpgrade returns
// the first non-activated entry.
func TestGetPendingUpgrade_ReturnsPending(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	require.NoError(t, k.SetKernelUpgradeInfo(ctx, &types.KernelUpgradeInfo{
		UpgradeVersion: "v3.0.0", ActivationHeight: 999, IsActivated: false,
	}))

	pending, err := k.GetPendingUpgrade(ctx)
	require.NoError(t, err)
	require.NotNil(t, pending)
	require.Equal(t, "v3.0.0", pending.UpgradeVersion)
}

// TestDeleteKernelUpgradeInfo_Success verifies that a stored info can be deleted.
func TestDeleteKernelUpgradeInfo_Success(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	require.NoError(t, k.SetKernelUpgradeInfo(ctx, &types.KernelUpgradeInfo{
		UpgradeVersion: "v5.0.0",
	}))

	require.NoError(t, k.DeleteKernelUpgradeInfo(ctx, "v5.0.0"))

	_, err := k.GetKernelUpgradeInfo(ctx, "v5.0.0")
	require.Error(t, err)
}

// TestDeleteActivatedUpgradeInfo_RemovesActivatedEntries verifies that
// deleteActivatedUpgradeInfo removes all activated entries but keeps pending ones.
func TestDeleteActivatedUpgradeInfo_RemovesActivatedEntries(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	require.NoError(t, k.SetKernelUpgradeInfo(ctx, &types.KernelUpgradeInfo{
		UpgradeVersion: "v-activated", IsActivated: true,
	}))
	require.NoError(t, k.SetKernelUpgradeInfo(ctx, &types.KernelUpgradeInfo{
		UpgradeVersion: "v-pending", IsActivated: false,
	}))

	require.NoError(t, k.deleteActivatedUpgradeInfo(ctx))

	// Activated entry should be deleted
	_, err := k.GetKernelUpgradeInfo(ctx, "v-activated")
	require.Error(t, err)

	// Pending entry should remain
	got, err := k.GetKernelUpgradeInfo(ctx, "v-pending")
	require.NoError(t, err)
	require.Equal(t, "v-pending", got.UpgradeVersion)
}

// TestDeleteActivatedUpgradeInfo_NoActivated verifies that calling
// deleteActivatedUpgradeInfo when no entries are activated is a no-op.
func TestDeleteActivatedUpgradeInfo_NoActivated(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	require.NoError(t, k.SetKernelUpgradeInfo(ctx, &types.KernelUpgradeInfo{
		UpgradeVersion: "v-pending", IsActivated: false,
	}))

	require.NoError(t, k.deleteActivatedUpgradeInfo(ctx))

	// Pending entry should still be there
	_, err := k.GetKernelUpgradeInfo(ctx, "v-pending")
	require.NoError(t, err)
}
