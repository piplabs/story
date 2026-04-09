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

// TestIsLatestDKGNetwork_SameRoundOlderBlock verifies that isLatestDKGNetwork returns
// false when same round exists but the new network has an older start block height.
func TestIsLatestDKGNetwork_SameRoundOlderBlock(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Store a network with round=1, startBlock=100
	existing := &types.DKGNetwork{
		Round:            1,
		Total:            3,
		Threshold:        2,
		Stage:            types.DKGStageActive,
		StartBlockHeight: 100,
	}
	require.NoError(t, k.setDKGNetwork(ctx, existing))

	// Candidate with same round=1 but older startBlock=50 → should NOT become latest
	candidate := &types.DKGNetwork{
		Round:            1,
		StartBlockHeight: 50,
	}
	isLatest, err := k.isLatestDKGNetwork(ctx, candidate)
	require.NoError(t, err)
	require.False(t, isLatest, "older start block height on same round should not become latest")
}

// TestIsLatestDKGNetwork_SameRoundNewerBlock verifies that isLatestDKGNetwork returns
// true when the same round has a newer start block height (e.g. TEE binary upgrade).
func TestIsLatestDKGNetwork_SameRoundNewerBlock(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Store a network with round=1, startBlock=50
	existing := &types.DKGNetwork{
		Round:            1,
		Total:            3,
		Threshold:        2,
		Stage:            types.DKGStageActive,
		StartBlockHeight: 50,
	}
	require.NoError(t, k.setDKGNetwork(ctx, existing))

	// Candidate with same round=1 but newer startBlock=200 → should become latest
	candidate := &types.DKGNetwork{
		Round:            1,
		StartBlockHeight: 200,
	}
	isLatest, err := k.isLatestDKGNetwork(ctx, candidate)
	require.NoError(t, err)
	require.True(t, isLatest, "newer start block on same round should become latest")
}

// TestIsLatestDKGNetwork_LowerRound verifies that isLatestDKGNetwork returns
// false when the candidate round is lower than the current latest.
func TestIsLatestDKGNetwork_LowerRound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Store a network with round=5
	existing := &types.DKGNetwork{
		Round:     5,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, existing))

	// Candidate with lower round=3 → should NOT become latest
	candidate := &types.DKGNetwork{Round: 3}
	isLatest, err := k.isLatestDKGNetwork(ctx, candidate)
	require.NoError(t, err)
	require.False(t, isLatest, "lower round should not become latest")
}

// TestSetDKGNetwork_DoesNotUpdateLatestForLowerRound verifies that setDKGNetwork
// does not update the latest pointer when storing a network with a lower round.
func TestSetDKGNetwork_DoesNotUpdateLatestForLowerRound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Store round 10 first (becomes latest)
	net10 := &types.DKGNetwork{Round: 10, Total: 3, Threshold: 2, Stage: types.DKGStageActive}
	require.NoError(t, k.setDKGNetwork(ctx, net10))

	// Now store round 5 — should NOT update the latest pointer
	net5 := &types.DKGNetwork{Round: 5, Total: 3, Threshold: 2, Stage: types.DKGStageRegistration}
	require.NoError(t, k.setDKGNetwork(ctx, net5))

	latest, err := k.GetLatestDKGRound(ctx)
	require.NoError(t, err)
	require.Equal(t, uint32(10), latest.Round, "latest pointer should still point to round 10")
}

// TestEndPreviousActiveRound_NilPrevActive verifies endPreviousActiveRound is a
// no-op when no previous active round exists.
func TestEndPreviousActiveRound_NilPrevActive(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// No active round set → endPreviousActiveRound should be no-op
	err := k.endPreviousActiveRound(ctx, 1)
	require.NoError(t, err)
}

// TestEndPreviousActiveRound_SameRound verifies endPreviousActiveRound is a
// no-op when the previous active round is the same as the current round.
func TestEndPreviousActiveRound_SameRound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{Round: 3, Total: 3, Threshold: 2, Stage: types.DKGStageActive}
	require.NoError(t, k.setDKGNetwork(ctx, network))
	require.NoError(t, k.setLatestActiveRound(ctx, network))

	// Same round as current → should be no-op
	err := k.endPreviousActiveRound(ctx, 3)
	require.NoError(t, err)

	// Stage should remain Active
	updated, err := k.getDKGNetwork(ctx, 3)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageActive, updated.Stage)
}

// TestGetLatestDKGRound_PointerExistsButNetworkMissing verifies GetLatestDKGRound
// returns an error and resets the pointer when the pointer exists but the network
// entry is missing (corrupted state).
func TestGetLatestDKGRound_PointerExistsButNetworkMissing(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Store a pointer to a network key that does not exist in DKGNetworks.
	// This simulates corrupted state (pointer without corresponding data).
	require.NoError(t, k.LatestDKGNetwork.Set(ctx, "99"))

	// GetLatestDKGRound should detect the missing network and return an error.
	// The pointer should also be reset (Remove is called internally).
	got, err := k.GetLatestDKGRound(ctx)
	require.Error(t, err, "should return error when pointer exists but network is missing")
	require.Nil(t, got)
	require.Contains(t, err.Error(), "not found")
}

// TestGetDKGNetworksByRound_StopsWhenRoundExceeded verifies that GetDKGNetworksByRound
// stops iteration early when it encounters a network with a round > target round.
func TestGetDKGNetworksByRound_StopsWhenRoundExceeded(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Store networks for rounds 1, 2, 3 in sorted order.
	// The Walk function iterates in lexicographic key order ("1", "2", "3"),
	// so when round 3 > target 2, iteration stops.
	for _, round := range []uint32{1, 2, 3} {
		net := &types.DKGNetwork{Round: round, Total: 3, Threshold: 2, Stage: types.DKGStageActive}
		require.NoError(t, k.setDKGNetwork(ctx, net))
	}

	// Request round 2 — should find it and stop before round 3.
	got, err := k.GetDKGNetworksByRound(ctx, 2)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, uint32(2), got[0].Round)
}

// TestIsLatestDKGNetwork_NoExistingLatest verifies isLatestDKGNetwork returns true
// when no latest DKG network has been set yet (empty store).
func TestIsLatestDKGNetwork_NoExistingLatest(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Empty store — no latest DKG network pointer set.
	candidate := &types.DKGNetwork{Round: 1, Total: 3, Threshold: 2}
	isLatest, err := k.isLatestDKGNetwork(ctx, candidate)
	require.NoError(t, err)
	require.True(t, isLatest, "first network should always become the latest")
}

// TestEndPreviousActiveRound_PrevNotActive verifies endPreviousActiveRound is a
// no-op when the previous active round is not in the Active stage (e.g. already Ended).
func TestEndPreviousActiveRound_PrevNotActive(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Round 1 is Ended (not Active)
	prev := &types.DKGNetwork{Round: 1, Total: 3, Threshold: 2, Stage: types.DKGStageEnded}
	require.NoError(t, k.setDKGNetwork(ctx, prev))
	require.NoError(t, k.setLatestActiveRound(ctx, prev))

	// endPreviousActiveRound for round 2 should not change round 1's stage
	err := k.endPreviousActiveRound(ctx, 2)
	require.NoError(t, err)

	updated, err := k.getDKGNetwork(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageEnded, updated.Stage, "stage should remain Ended")
}

// TestEndPreviousActiveRound_TransitionsToEnded verifies that endPreviousActiveRound
// transitions a round in the Active stage to Ended when a different current round is specified.
func TestEndPreviousActiveRound_TransitionsToEnded(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Round 3 is Active
	prev := &types.DKGNetwork{Round: 3, Total: 5, Threshold: 4, Stage: types.DKGStageActive}
	require.NoError(t, k.setDKGNetwork(ctx, prev))
	require.NoError(t, k.setLatestActiveRound(ctx, prev))

	// endPreviousActiveRound for round 4 should transition round 3 to Ended
	err := k.endPreviousActiveRound(ctx, 4)
	require.NoError(t, err)

	updated, err := k.getDKGNetwork(ctx, 3)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageEnded, updated.Stage, "active round should be transitioned to ended")
}

// TestGetNextRoundNumber_NoLatest verifies that getNextRoundNumber returns 1
// when no rounds exist.
func TestGetNextRoundNumber_NoLatest(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	nextRound := k.getNextRoundNumber(ctx)
	require.Equal(t, uint32(1), nextRound)
}

// TestGetNextRoundNumber_WithExisting verifies that getNextRoundNumber returns
// latestRound+1 when rounds already exist.
func TestGetNextRoundNumber_WithExisting(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{Round: 7, Total: 3, Threshold: 2, Stage: types.DKGStageActive}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	nextRound := k.getNextRoundNumber(ctx)
	require.Equal(t, uint32(8), nextRound)
}

// TestGetLatestDKGNetwork verifies the getLatestDKGNetwork function retrieves
// the latest DKG network correctly when it exists.
func TestGetLatestDKGNetwork_Success(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	net1 := &types.DKGNetwork{Round: 1, Total: 3, Threshold: 2, Stage: types.DKGStageActive}
	net2 := &types.DKGNetwork{Round: 2, Total: 5, Threshold: 4, Stage: types.DKGStageRegistration}
	require.NoError(t, k.setDKGNetwork(ctx, net1))
	require.NoError(t, k.setDKGNetwork(ctx, net2))

	got, err := k.getLatestDKGNetwork(ctx)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, uint32(2), got.Round)
}

// TestGetLatestDKGNetwork_NoPointer verifies that getLatestDKGNetwork returns
// an error when no latest pointer is set (fresh store).
func TestGetLatestDKGNetwork_NoPointer(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.getLatestDKGNetwork(ctx)
	require.Error(t, err, "should error when no latest DKG network pointer is set")
}
