package keeper

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// marshalVote serializes a Vote proto message and returns the bytes.
func marshalVote(t *testing.T, v *types.Vote) []byte {
	t.Helper()
	bz, err := proto.Marshal(v)
	require.NoError(t, err)
	return bz
}

// newTestDeal creates a minimal Deal for testing.
func newTestDeal(dealerIdx, recipientIdx uint32) types.Deal {
	return types.Deal{Index: dealerIdx, RecipientIndex: recipientIdx}
}

// TestParseAndVerifyVoteExtension_EmptyExtension verifies that an empty
// vote extension is accepted (no error, no vote returned).
func TestParseAndVerifyVoteExtension_EmptyExtension(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	votes, ok, err := k.parseAndVerifyVoteExtension(nil)
	require.NoError(t, err)
	require.True(t, ok, "empty extension should be accepted")
	require.Nil(t, votes)
}

// TestParseAndVerifyVoteExtension_EmptyBytes verifies that zero-length bytes
// are treated the same as nil (accepted).
func TestParseAndVerifyVoteExtension_EmptyBytes(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	votes, ok, err := k.parseAndVerifyVoteExtension([]byte{})
	require.NoError(t, err)
	require.True(t, ok)
	require.Nil(t, votes)
}

// TestParseAndVerifyVoteExtension_ValidVote verifies that a well-formed
// protobuf-encoded Vote is accepted.
func TestParseAndVerifyVoteExtension_ValidVote(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	vote := &types.Vote{
		Deals: []types.Deal{newTestDeal(1, 2)},
	}

	bz := marshalVote(t, vote)
	votes, ok, err := k.parseAndVerifyVoteExtension(bz)
	require.NoError(t, err)
	require.True(t, ok)
	require.Len(t, votes, 1)
	require.Len(t, votes[0].Deals, 1)
}

// TestParseAndVerifyVoteExtension_TooLarge verifies that a vote extension
// exceeding maxVoteExtensionSize is rejected.
func TestParseAndVerifyVoteExtension_TooLarge(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	// Create a byte slice just over 256 KB
	oversized := make([]byte, maxVoteExtensionSize+1)
	_, _, err := k.parseAndVerifyVoteExtension(oversized)
	require.Error(t, err)
	require.Contains(t, err.Error(), "vote extension too large")
}

// TestParseAndVerifyVoteExtension_InvalidProto verifies that non-protobuf
// bytes cause an error.
func TestParseAndVerifyVoteExtension_InvalidProto(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	// Deliberately malformed protobuf — not a valid wire format
	garbage := []byte{0xFF, 0xFE, 0xFD, 0x00, 0x01, 0x02}
	_, _, err := k.parseAndVerifyVoteExtension(garbage)
	require.Error(t, err)
}

// TestVerifyVoteExtension_ValidVote verifies that VerifyVoteExtension accepts
// a valid vote extension after the V200 upgrade height.
func TestVerifyVoteExtension_ValidVote(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)
	// newTestSDKContext already sets TestChainID + height=200 (past V200=110)
	sdkCtx := newTestSDKContext(t, "vve_valid")

	vote := &types.Vote{
		Deals: []types.Deal{newTestDeal(1, 2)},
	}

	bz := marshalVote(t, vote)
	req := &abci.RequestVerifyVoteExtension{VoteExtension: bz}
	resp, err := k.VerifyVoteExtension(sdkCtx, req)
	require.NoError(t, err)
	require.Equal(t, abci.ResponseVerifyVoteExtension_ACCEPT, resp.Status)
}

// TestVerifyVoteExtension_InvalidVote verifies that VerifyVoteExtension rejects
// a malformed vote extension with REJECT status (not a Go error).
func TestVerifyVoteExtension_InvalidVote(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)
	sdkCtx := newTestSDKContext(t, "vve_invalid")

	oversized := make([]byte, maxVoteExtensionSize+1)
	req := &abci.RequestVerifyVoteExtension{VoteExtension: oversized}
	resp, err := k.VerifyVoteExtension(sdkCtx, req)
	require.NoError(t, err, "VerifyVoteExtension must not return Go error on bad VE")
	require.Equal(t, abci.ResponseVerifyVoteExtension_REJECT, resp.Status)
}

// TestVerifyVoteExtension_EmptyVote verifies that an empty vote extension is
// accepted (ACCEPT status, empty vote is fine).
func TestVerifyVoteExtension_EmptyVote(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)
	sdkCtx := newTestSDKContext(t, "vve_empty")

	req := &abci.RequestVerifyVoteExtension{VoteExtension: nil}
	resp, err := k.VerifyVoteExtension(sdkCtx, req)
	require.NoError(t, err)
	require.Equal(t, abci.ResponseVerifyVoteExtension_ACCEPT, resp.Status)
}

// TestVerifyVoteExtension_BeforeV200 verifies that VerifyVoteExtension accepts
// all vote extensions before V200 (before vote extensions are active).
// Uses newTestSDKContext reusing the vve_before_v200 key name but overriding to height 50.
func TestVerifyVoteExtension_BeforeV200(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	// newTestSDKContext sets height=200 (past V200). Override to before V200 (110) here.
	sdkCtx := newTestSDKContext(t, "vve_before_v200").WithBlockHeight(50)

	garbage := make([]byte, maxVoteExtensionSize+1) // would normally be rejected
	req := &abci.RequestVerifyVoteExtension{VoteExtension: garbage}
	resp, err := k.VerifyVoteExtension(sdkCtx, req)
	require.NoError(t, err)
	// Before V200, all VEs are accepted regardless of content
	require.Equal(t, abci.ResponseVerifyVoteExtension_ACCEPT, resp.Status)
}

// TestAggregateVotes_Empty verifies that aggregateVotes on empty input returns
// an empty Vote struct.
func TestAggregateVotes_Empty(t *testing.T) {
	t.Parallel()

	result := aggregateVotes(nil)
	require.NotNil(t, result)
	require.Empty(t, result.Deals)
	require.Empty(t, result.Responses)
	require.Empty(t, result.Justifications)
}

// TestAggregateVotes_DeduplicatesDeals verifies that duplicate deals (same
// dealerIndex + recipientIndex) are deduplicated.
func TestAggregateVotes_DeduplicatesDeals(t *testing.T) {
	t.Parallel()

	deal1 := newTestDeal(1, 2)
	deal2 := newTestDeal(1, 2) // duplicate key
	deal3 := newTestDeal(1, 3)

	votes := []*types.Vote{
		{Deals: []types.Deal{deal1}},
		{Deals: []types.Deal{deal2, deal3}},
	}

	result := aggregateVotes(votes)
	require.Len(t, result.Deals, 2, "duplicate deal should be removed")
}

// TestAggregateVotes_DeduplicatesResponses verifies that duplicate responses
// (same responderIndex + dealerIndex) are deduplicated.
func TestAggregateVotes_DeduplicatesResponses(t *testing.T) {
	t.Parallel()

	resp1 := types.Response{Index: 1, VssResponse: &types.VSSResponse{Index: 2}}
	resp2 := types.Response{Index: 1, VssResponse: &types.VSSResponse{Index: 2}} // duplicate
	resp3 := types.Response{Index: 2, VssResponse: &types.VSSResponse{Index: 2}}

	votes := []*types.Vote{
		{Responses: []types.Response{resp1, resp2}},
		{Responses: []types.Response{resp3}},
	}

	result := aggregateVotes(votes)
	require.Len(t, result.Responses, 2, "duplicate response should be removed")
}

// TestDeduplicateDeals_AllUnique verifies deduplicateDeals preserves all items
// when there are no duplicates.
func TestDeduplicateDeals_AllUnique(t *testing.T) {
	t.Parallel()

	deals := []types.Deal{
		newTestDeal(1, 1),
		newTestDeal(1, 2),
		newTestDeal(2, 1),
	}

	result := deduplicateDeals(deals)
	require.Len(t, result, 3)
}

// TestDeduplicateDeals_AllDuplicates verifies deduplicateDeals returns only
// the first occurrence when all items are duplicates of the same key.
func TestDeduplicateDeals_AllDuplicates(t *testing.T) {
	t.Parallel()

	deals := []types.Deal{
		{Index: 1, RecipientIndex: 1, Signature: []byte("sig-a")},
		{Index: 1, RecipientIndex: 1, Signature: []byte("sig-b")},
		{Index: 1, RecipientIndex: 1, Signature: []byte("sig-c")},
	}

	result := deduplicateDeals(deals)
	require.Len(t, result, 1)
	require.Equal(t, []byte("sig-a"), result[0].Signature, "first occurrence should be kept")
}

// TestDeduplicateDeals_Empty verifies deduplicateDeals handles empty input.
func TestDeduplicateDeals_Empty(t *testing.T) {
	t.Parallel()

	result := deduplicateDeals(nil)
	require.Empty(t, result)
}

// TestDeduplicateResponses_NilVssResponse verifies deduplicateResponses handles
// nil VssResponse (uses responderIndex + 0 as the dedup key).
func TestDeduplicateResponses_NilVssResponse(t *testing.T) {
	t.Parallel()

	responses := []types.Response{
		{Index: 1, VssResponse: nil},
		{Index: 1, VssResponse: nil}, // duplicate
	}

	result := deduplicateResponses(responses)
	require.Len(t, result, 1)
}

// TestDeduplicateResponses_Empty verifies deduplicateResponses handles empty input.
func TestDeduplicateResponses_Empty(t *testing.T) {
	t.Parallel()

	result := deduplicateResponses(nil)
	require.Empty(t, result)
}

// TestVotesFromExtension_Empty verifies votesFromExtension returns (nil, false, nil)
// for an empty byte slice.
func TestVotesFromExtension_Empty(t *testing.T) {
	t.Parallel()

	vote, ok, err := votesFromExtension(nil)
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, vote)
}

// TestVotesFromExtension_ValidVote verifies votesFromExtension correctly
// unmarshals a valid vote extension.
func TestVotesFromExtension_ValidVote(t *testing.T) {
	t.Parallel()

	expected := &types.Vote{
		Deals: []types.Deal{newTestDeal(3, 5)},
	}

	bz, err := proto.Marshal(expected)
	require.NoError(t, err)

	vote, ok, err := votesFromExtension(bz)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, vote)
	require.Len(t, vote.Deals, 1)
	require.Equal(t, uint32(3), vote.Deals[0].Index)
}

// TestVotesFromExtension_InvalidProto verifies votesFromExtension returns an
// error for malformed bytes.
func TestVotesFromExtension_InvalidProto(t *testing.T) {
	t.Parallel()

	_, _, err := votesFromExtension([]byte{0xFF, 0xFE, 0x00})
	require.Error(t, err)
}

// TestExtendVote_ReturnsSerializedVote verifies ExtendVote encodes the
// queued deals/responses/justifications into the vote extension.
// Note: uses package-level globals, so is NOT run in parallel.
func TestExtendVote_ReturnsSerializedVote(t *testing.T) {
	k, _, _, _ := setupDKGKeeperWithMocks(t)
	// Use TestChainID at height > V200 (110)
	sdkCtx := newTestSDKContext(t, "extend_vote_deal")

	// Flush any state left by previous tests before enqueuing
	k.FlushAllQueues()

	// Pre-enqueue a deal
	k.EnqueueDeals([]types.Deal{newTestDeal(1, 2)})

	resp, err := k.ExtendVote(sdkCtx, &abci.RequestExtendVote{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.VoteExtension)

	// Decode and verify the vote contains our deal
	var vote types.Vote
	require.NoError(t, proto.Unmarshal(resp.VoteExtension, &vote))
	require.Len(t, vote.Deals, 1)
	require.Equal(t, uint32(1), vote.Deals[0].Index)
}
