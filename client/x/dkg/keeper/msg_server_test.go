package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

const testAuthority = "story1hmjw3pvkjtndpg8wqppwdn8udd835qpan4hm0y"

// TestMsgServer_AddVote_UnauthorizedCaller verifies that AddVote returns an
// error when the authority in the message does not match the keeper's authority.
func TestMsgServer_AddVote_UnauthorizedCaller(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	srv := NewMsgServerImpl(k)

	msg := &types.MsgAddDkgVote{
		Authority: "story1wrongauthority111111111111111111111111",
		Vote:      &types.Vote{},
	}

	_, err := srv.AddVote(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized")
}

// TestMsgServer_AddVote_EmptyVote_NoLatestRound verifies that AddVote succeeds
// when no latest DKG round exists (nothing to process).
func TestMsgServer_AddVote_EmptyVote_NoLatestRound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	srv := NewMsgServerImpl(k)

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote:      &types.Vote{},
	}

	resp, err := srv.AddVote(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestMsgServer_AddVote_RoundNotInDealingStage verifies that AddVote with a
// non-dealing stage round does not attempt to process deals/responses.
func TestMsgServer_AddVote_RoundNotInDealingStage(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	srv := NewMsgServerImpl(k)

	// Set up a round in Registration stage (not Dealing)
	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageRegistration,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote: &types.Vote{
			Deals: []types.Deal{
				{Index: 1, RecipientIndex: 2},
			},
		},
	}

	// Should succeed without error even though stage is not Dealing
	resp, err := srv.AddVote(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestMsgServer_NewMsgServerImpl verifies that NewMsgServerImpl returns a
// non-nil server that implements types.MsgServiceServer.
func TestMsgServer_NewMsgServerImpl(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)
	srv := NewMsgServerImpl(k)
	require.NotNil(t, srv)

	// Verify it satisfies the interface
	var _ types.MsgServiceServer = srv
}

// TestMsgServer_AddVote_NilVote verifies that AddVote handles a nil Vote field
// without panicking (should fail gracefully at authority check if authority is wrong).
func TestMsgServer_AddVote_NilVote_Unauthorized(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	srv := NewMsgServerImpl(k)

	msg := &types.MsgAddDkgVote{
		Authority: "story1badauthority111111111111111111111111111",
		Vote:      nil, // nil vote, but authority check comes first
	}

	_, err := srv.AddVote(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized")
}

// TestMsgServer_AddVote_EmptyDeals_DealingStage verifies that AddVote with
// empty deals in dealing stage does not cause errors.
func TestMsgServer_AddVote_EmptyDeals_DealingStage(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	srv := NewMsgServerImpl(k)

	// Set up a round in Dealing stage
	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote:      &types.Vote{
			// Empty deals, responses, justifications
		},
	}

	resp, err := srv.AddVote(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestMsgServer_AddVote_GetLatestRoundError verifies that AddVote propagates
// an error from GetLatestDKGRound (e.g. corrupted store state).
func TestMsgServer_AddVote_GetLatestRoundError(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	srv := NewMsgServerImpl(k)

	// Simulate corrupted state: a LatestDKGNetwork pointer pointing to a
	// non-existent network entry. GetLatestDKGRound will return a not-found error.
	require.NoError(t, k.LatestDKGNetwork.Set(ctx, "9999"))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote:      &types.Vote{},
	}

	_, err := srv.AddVote(ctx, msg)
	require.Error(t, err, "AddVote should propagate GetLatestDKGRound error")
	require.Contains(t, err.Error(), "not found")
}

// TestMsgServer_GetAuthority verifies that the keeper returns the correct
// authority address.
func TestMsgServer_GetAuthority(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	authority := k.GetAuthority()
	require.Equal(t, testAuthority, authority)
}

// TestMsgServer_AddVote_WithDeals_DealingStage verifies that AddVote processes
// deals when the round is in Dealing stage with non-empty deals.
func TestMsgServer_AddVote_WithDeals_DealingStage(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	srv := NewMsgServerImpl(k)

	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote: &types.Vote{
			Deals: []types.Deal{
				{Index: 1, RecipientIndex: 2},
				{Index: 1, RecipientIndex: 3},
			},
		},
	}

	resp, err := srv.AddVote(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestMsgServer_AddVote_WithResponses_DealingStage verifies that AddVote processes
// responses when the round is in Dealing stage with non-empty responses.
func TestMsgServer_AddVote_WithResponses_DealingStage(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	srv := NewMsgServerImpl(k)

	network := &types.DKGNetwork{
		Round:     2,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote: &types.Vote{
			Responses: []types.Response{
				{Index: 1},
				{Index: 2},
			},
		},
	}

	resp, err := srv.AddVote(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestMsgServer_AddVote_WithAllItems_DealingStage verifies that AddVote processes
// deals, responses, and justifications when all are present in a Dealing stage round.
func TestMsgServer_AddVote_WithAllItems_DealingStage(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	srv := NewMsgServerImpl(k)

	network := &types.DKGNetwork{
		Round:     3,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote: &types.Vote{
			Deals: []types.Deal{
				{Index: 1, RecipientIndex: 2},
			},
			Responses: []types.Response{
				{Index: 1},
			},
			Justifications: []types.Justification{
				// Empty justification — no valid sig, no registration → dropped at sig check
				{Index: 0, VssJustification: nil},
			},
		},
	}

	// Should succeed even if justification processing logs errors
	resp, err := srv.AddVote(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
