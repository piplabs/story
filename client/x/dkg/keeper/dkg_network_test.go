package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// TestSetAndGetDKGNetwork verifies the basic set → get round trip.
func TestSetAndGetDKGNetwork(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     1,
		Total:     5,
		Threshold: 4,
		Stage:     types.DKGStageRegistration,
	}

	require.NoError(t, k.setDKGNetwork(ctx, network))

	got, err := k.getDKGNetwork(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, uint32(1), got.Round)
	require.Equal(t, uint32(5), got.Total)
}

// TestGetDKGNetwork_NotFound verifies that getDKGNetwork returns an error for
// a non-existent round.
func TestGetDKGNetwork_NotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.getDKGNetwork(ctx, 999)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

// TestGetLatestDKGRound_NoNetwork verifies that GetLatestDKGRound returns
// (nil, nil) when no network exists.
func TestGetLatestDKGRound_NoNetwork(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	got, err := k.GetLatestDKGRound(ctx)
	require.NoError(t, err)
	require.Nil(t, got)
}

// TestGetLatestDKGRound_UpdatesOnHigherRound verifies that the latest round
// pointer is updated when a higher-round network is stored.
func TestGetLatestDKGRound_UpdatesOnHigherRound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	net1 := &types.DKGNetwork{Round: 1, Total: 3, Threshold: 2, Stage: types.DKGStageDealing}
	net2 := &types.DKGNetwork{Round: 2, Total: 3, Threshold: 2, Stage: types.DKGStageDealing}
	net3 := &types.DKGNetwork{Round: 3, Total: 3, Threshold: 2, Stage: types.DKGStageDealing}

	require.NoError(t, k.setDKGNetwork(ctx, net1))
	require.NoError(t, k.setDKGNetwork(ctx, net2))
	require.NoError(t, k.setDKGNetwork(ctx, net3))

	latest, err := k.GetLatestDKGRound(ctx)
	require.NoError(t, err)
	require.Equal(t, uint32(3), latest.Round)
}

// TestGetDKGNetworksByRound_Found verifies GetDKGNetworksByRound returns
// the correct network.
func TestGetDKGNetworksByRound_Found(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	net := &types.DKGNetwork{Round: 5, Total: 7, Threshold: 5, Stage: types.DKGStageActive}
	require.NoError(t, k.setDKGNetwork(ctx, net))

	got, err := k.GetDKGNetworksByRound(ctx, 5)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, uint32(5), got[0].Round)
}

// TestGetDKGNetworksByRound_NotFound verifies GetDKGNetworksByRound returns
// empty slice for non-existent round.
func TestGetDKGNetworksByRound_NotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	got, err := k.GetDKGNetworksByRound(ctx, 99)
	require.NoError(t, err)
	require.Empty(t, got)
}

// TestDeleteDKGNetwork verifies that a stored network can be deleted and
// is no longer retrievable.
func TestDeleteDKGNetwork(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	net := &types.DKGNetwork{Round: 10, Total: 3, Threshold: 2, Stage: types.DKGStageDealing}
	require.NoError(t, k.setDKGNetwork(ctx, net))

	require.NoError(t, k.DeleteDKGNetwork(ctx, 10))

	_, err := k.getDKGNetwork(ctx, 10)
	require.Error(t, err)
}

// TestGetAllDKGNetworks_Empty verifies getAllDKGNetworks returns empty when
// no networks are stored.
func TestGetAllDKGNetworks_EmptyInternal(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	networks, err := k.getAllDKGNetworks(ctx)
	require.NoError(t, err)
	require.Empty(t, networks)
}

// TestGetAllDKGNetworks_Multiple verifies getAllDKGNetworks returns all
// stored networks.
func TestGetAllDKGNetworks_MultipleInternal(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	for _, round := range []uint32{1, 2, 3, 4} {
		net := &types.DKGNetwork{Round: round, Total: 5, Threshold: 4, Stage: types.DKGStageActive}
		require.NoError(t, k.setDKGNetwork(ctx, net))
	}

	networks, err := k.getAllDKGNetworks(ctx)
	require.NoError(t, err)
	require.Len(t, networks, 4)
}

// TestGetLatestActiveRound verifies GetLatestActiveRound returns nil when
// no active round exists.
func TestGetLatestActiveRound_NoActive(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	got, err := k.GetLatestActiveRound(ctx)
	require.NoError(t, err)
	require.Nil(t, got)
}

// TestSetAndGetLatestActiveRound verifies that setLatestActiveRound stores
// and GetLatestActiveRound retrieves the active round.
func TestSetAndGetLatestActiveRound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:           7,
		Total:           5,
		Threshold:       4,
		Stage:           types.DKGStageActive,
		GlobalPublicKey: []byte("test-global-key"),
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))
	require.NoError(t, k.setLatestActiveRound(ctx, network))

	got, err := k.GetLatestActiveRound(ctx)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, uint32(7), got.Round)
}
