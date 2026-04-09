package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// TestProposalServer_AddVote_Authorized verifies that the proposal server
// AddVote returns success when called with the correct authority.
func TestProposalServer_AddVote_Authorized(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	srv := NewProposalServer(k)

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote:      &types.Vote{},
	}

	resp, err := srv.AddVote(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestProposalServer_AddVote_Unauthorized verifies that the proposal server
// AddVote returns unauthorized error for a wrong authority.
func TestProposalServer_AddVote_Unauthorized(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	srv := NewProposalServer(k)

	msg := &types.MsgAddDkgVote{
		Authority: "story1wrongaddress1111111111111111111111111",
		Vote:      &types.Vote{},
	}

	_, err := srv.AddVote(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized")
}

// TestNewProposalServer_ReturnsNonNil verifies that NewProposalServer returns
// a non-nil server implementing MsgServiceServer.
func TestNewProposalServer_ReturnsNonNil(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)
	srv := NewProposalServer(k)
	require.NotNil(t, srv)

	var _ types.MsgServiceServer = srv
}
