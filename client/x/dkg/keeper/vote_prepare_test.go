package keeper

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/lib/netconf"
)

// TestPrepareVotes_BeforeV200PlusOne verifies that PrepareVotes returns
// (nil, nil) when the block height is at or before v200Height+1.
// This is the early-return path that keeps compatibility with pre-upgrade validators.
func TestPrepareVotes_BeforeV200PlusOne(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	// TestChainID has V200 activation at height 110.
	// At height 110 (= v200Height), we should get the early nil return.
	sdkCtx := newTestSDKContext(t, "prepare_votes_before").WithBlockHeight(110)

	commit := abci.ExtendedCommitInfo{}
	msg, err := k.PrepareVotes(sdkCtx, commit, 110)
	require.NoError(t, err)
	require.Nil(t, msg, "PrepareVotes should return nil before V200+1")
}

// TestPrepareVotes_AtV200PlusOne verifies that PrepareVotes returns (nil, nil)
// exactly at v200Height+1 (the boundary inclusive check).
func TestPrepareVotes_AtV200PlusOne(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	// V200 = 110 for TestChainID, so height 111 is V200+1 (still nil return)
	sdkCtx := newTestSDKContext(t, "prepare_votes_at_v200p1").WithBlockHeight(111)

	commit := abci.ExtendedCommitInfo{}
	msg, err := k.PrepareVotes(sdkCtx, commit, 111)
	require.NoError(t, err)
	require.Nil(t, msg, "PrepareVotes should return nil at V200+1")
}

// TestPrepareVotes_UnknownChainID verifies that PrepareVotes returns an error
// for an unrecognized chain ID.
func TestPrepareVotes_UnknownChainID(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	// Override chain ID to something unknown
	sdkCtx := newTestSDKContext(t, "prepare_votes_bad_chain").
		WithChainID("unknown-chain-1")

	commit := abci.ExtendedCommitInfo{}
	_, err := k.PrepareVotes(sdkCtx, commit, 100)
	_ = ctx
	require.Error(t, err)
	require.Contains(t, err.Error(), "get v2.0.0 upgrade height")
}

// TestPrepareVotes_ValidNetworkID_ReturnsNilBeforeV200 verifies that
// PrepareVotes on a known network (iliad, odyssey) returns nil before V200.
func TestPrepareVotes_ValidNetworkID_ReturnsNilBeforeV200(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	// iliad chain has a specific V200 height; use height=1 which is always before it.
	sdkCtx := newTestSDKContext(t, "prepare_votes_iliad").
		WithChainID(netconf.TestChainID).
		WithBlockHeight(1)

	commit := abci.ExtendedCommitInfo{}
	msg, err := k.PrepareVotes(sdkCtx, commit, 1)
	require.NoError(t, err)
	require.Nil(t, msg)
}
