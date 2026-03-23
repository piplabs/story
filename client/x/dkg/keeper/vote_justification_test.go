package keeper

import (
	"testing"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/netconf"
)

// newTestSDKContext creates a minimal sdk.Context for testing functions that
// accept sdk.Context but do not use it (e.g., ExtendVote ignores its context).
func newTestSDKContext(t *testing.T, keyName string) sdk.Context {
	t.Helper()

	key := storetypes.NewKVStoreKey(keyName)
	transKey := storetypes.NewTransientStoreKey(keyName + "_transient")
	testCtx := testutil.DefaultContextWithDB(t, key, transKey)

	return testCtx.Ctx.WithChainID(netconf.TestChainID).WithBlockHeight(200)
}

// buildJustificationVote creates a proto-marshaled Vote containing the given
// justifications, suitable for use as a vote extension payload.
func buildJustificationVote(t *testing.T, justifications []types.Justification) []byte {
	t.Helper()

	bz, err := proto.Marshal(&types.Vote{
		Justifications: justifications,
	})
	require.NoError(t, err)

	return bz
}

// TestAggregateVotes_IncludesJustifications verifies that aggregateVotes merges
// all justifications from multiple votes into a single output vote.
func TestAggregateVotes_IncludesJustifications(t *testing.T) {
	t.Parallel()

	j1 := types.Justification{Index: 1, VssJustification: &types.VSSJustification{SessionId: []byte("s1")}}
	j2 := types.Justification{Index: 2, VssJustification: &types.VSSJustification{SessionId: []byte("s2")}}
	j3 := types.Justification{Index: 3, VssJustification: &types.VSSJustification{SessionId: []byte("s3")}}

	vote1 := &types.Vote{Justifications: []types.Justification{j1, j2}}
	vote2 := &types.Vote{Justifications: []types.Justification{j3}}

	result := aggregateVotes([]*types.Vote{vote1, vote2})
	require.Len(t, result.Justifications, 3, "all justifications should be merged")
	require.Equal(t, uint32(1), result.Justifications[0].Index)
	require.Equal(t, uint32(2), result.Justifications[1].Index)
	require.Equal(t, uint32(3), result.Justifications[2].Index)
}

// TestAggregateVotes_EmptyJustifications verifies that aggregateVotes handles
// votes with no justifications and returns an empty slice.
func TestAggregateVotes_EmptyJustifications(t *testing.T) {
	t.Parallel()

	vote1 := &types.Vote{Justifications: []types.Justification{}}
	vote2 := &types.Vote{Justifications: nil}

	result := aggregateVotes([]*types.Vote{vote1, vote2})
	require.Empty(t, result.Justifications, "no justifications across votes should produce an empty slice")
}

// TestAggregateVotes_NilInput verifies that aggregateVotes handles nil input
// without panicking.
func TestAggregateVotes_NilInput(t *testing.T) {
	t.Parallel()

	result := aggregateVotes(nil)
	require.NotNil(t, result)
	require.Empty(t, result.Justifications)
}

// TestAggregateVotes_MixedDealsResponsesJustifications verifies that aggregateVotes
// correctly merges all three data types from multiple votes.
func TestAggregateVotes_MixedDealsResponsesJustifications(t *testing.T) {
	t.Parallel()

	vote1 := &types.Vote{
		Deals:          []types.Deal{{Index: 10, RecipientIndex: 20}},
		Responses:      []types.Response{{Index: 30}},
		Justifications: []types.Justification{{Index: 40}},
	}
	vote2 := &types.Vote{
		Deals:          []types.Deal{{Index: 11}},
		Justifications: []types.Justification{{Index: 41}},
	}

	result := aggregateVotes([]*types.Vote{vote1, vote2})
	require.Len(t, result.Deals, 2, "deals should be merged")
	require.Len(t, result.Responses, 1, "responses should be merged")
	require.Len(t, result.Justifications, 2, "justifications should be merged")
}

// TestVotesFromExtension_WithJustifications verifies that proto-marshaled
// vote extensions containing justifications are parsed correctly.
func TestVotesFromExtension_WithJustifications(t *testing.T) {
	t.Parallel()

	j := types.Justification{
		Index: 99,
		VssJustification: &types.VSSJustification{
			SessionId: []byte("ext-session"),
			Index:     99,
		},
	}

	bz := buildJustificationVote(t, []types.Justification{j})

	vote, ok, err := votesFromExtension(bz)
	require.NoError(t, err)
	require.True(t, ok)
	require.Len(t, vote.Justifications, 1)
	require.Equal(t, uint32(99), vote.Justifications[0].Index)
}

// TestVotesFromExtension_EmptyBytes verifies that nil vote extension bytes
// return (nil, false, nil).
func TestVotesFromExtension_EmptyBytes(t *testing.T) {
	t.Parallel()

	vote, ok, err := votesFromExtension(nil)
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, vote)
}

// TestVotesFromExtension_MalformedBytes verifies that garbage bytes produce an error.
func TestVotesFromExtension_MalformedBytes(t *testing.T) {
	t.Parallel()

	vote, ok, err := votesFromExtension([]byte("this is not proto"))
	require.Error(t, err)
	require.False(t, ok)
	require.Nil(t, vote)
}

// TestExtendVote_IncludesJustifications verifies that ExtendVote dequeues
// justifications and includes them in the vote extension payload.
func TestExtendVote_IncludesJustifications(t *testing.T) {
	drainJustifications()
	defer drainJustifications()

	k, _ := setupDKGKeeper(t)

	// Enqueue two justifications before calling ExtendVote
	k.EnqueueJustifications([]*types.Justification{
		makeTestJustification(55),
		makeTestJustification(66),
	})

	sdkCtx := newTestSDKContext(t, "extend_vote_just")
	resp, err := k.ExtendVote(sdkCtx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Parse the vote extension to verify justifications are included
	vote, ok, err := votesFromExtension(resp.VoteExtension)
	require.NoError(t, err)
	require.True(t, ok)
	require.Len(t, vote.Justifications, 2, "both enqueued justifications should appear in the vote extension")
	require.Equal(t, uint32(55), vote.Justifications[0].Index)
	require.Equal(t, uint32(66), vote.Justifications[1].Index)
}

// TestExtendVote_EmptyQueues verifies that ExtendVote produces a valid (empty)
// vote extension when all queues are empty.
func TestExtendVote_EmptyQueues(t *testing.T) {
	drainJustifications()
	defer drainJustifications()

	k, _ := setupDKGKeeper(t)

	sdkCtx := newTestSDKContext(t, "extend_vote_empty")
	resp, err := k.ExtendVote(sdkCtx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// When all queues are empty, proto.Marshal of an empty Vote produces
	// a zero-length byte slice, so votesFromExtension returns (nil, false, nil).
	vote, ok, err := votesFromExtension(resp.VoteExtension)
	require.NoError(t, err)

	if ok {
		require.Empty(t, vote.Justifications)
		require.Empty(t, vote.Deals)
		require.Empty(t, vote.Responses)
	}
	// ok may be false if all fields are empty (proto3 zero-value marshaling)
}

// TestParseAndVerifyVoteExtension_OversizedPayload verifies that a vote extension
// exceeding maxVoteExtensionSize is rejected before parsing.
func TestParseAndVerifyVoteExtension_OversizedPayload(t *testing.T) {
	t.Parallel()

	k, _ := setupDKGKeeper(t)
	oversized := make([]byte, maxVoteExtensionSize+1)
	_, _, err := k.parseAndVerifyVoteExtension(oversized)
	require.Error(t, err)
	require.Contains(t, err.Error(), "exceeds max")
}

// TestParseAndVerifyVoteExtension_TooManyDeals verifies that a vote with more
// than maxItemsPerVote deals is rejected.
func TestParseAndVerifyVoteExtension_TooManyDeals(t *testing.T) {
	t.Parallel()

	k, _ := setupDKGKeeper(t)
	deals := make([]types.Deal, maxItemsPerVote+1)
	for i := range deals {
		deals[i] = types.Deal{Index: uint32(i)}
	}

	bz, err := proto.Marshal(&types.Vote{Deals: deals})
	require.NoError(t, err)

	_, _, err = k.parseAndVerifyVoteExtension(bz)
	require.Error(t, err)
	require.Contains(t, err.Error(), "deals")
}

// TestParseAndVerifyVoteExtension_TooManyResponses verifies that a vote with more
// than maxItemsPerVote responses is rejected.
func TestParseAndVerifyVoteExtension_TooManyResponses(t *testing.T) {
	t.Parallel()

	k, _ := setupDKGKeeper(t)
	responses := make([]types.Response, maxItemsPerVote+1)
	for i := range responses {
		responses[i] = types.Response{Index: uint32(i)}
	}

	bz, err := proto.Marshal(&types.Vote{Responses: responses})
	require.NoError(t, err)

	_, _, err = k.parseAndVerifyVoteExtension(bz)
	require.Error(t, err)
	require.Contains(t, err.Error(), "responses")
}

// TestParseAndVerifyVoteExtension_TooManyJustifications verifies that a vote with more
// than maxItemsPerVote justifications is rejected.
func TestParseAndVerifyVoteExtension_TooManyJustifications(t *testing.T) {
	t.Parallel()

	k, _ := setupDKGKeeper(t)
	justifications := make([]types.Justification, maxItemsPerVote+1)
	for i := range justifications {
		justifications[i] = types.Justification{Index: uint32(i)}
	}

	bz, err := proto.Marshal(&types.Vote{Justifications: justifications})
	require.NoError(t, err)

	_, _, err = k.parseAndVerifyVoteExtension(bz)
	require.Error(t, err)
	require.Contains(t, err.Error(), "justifications")
}

// TestParseAndVerifyVoteExtension_AtLimits verifies that a vote with exactly
// the maximum allowed counts passes validation.
func TestParseAndVerifyVoteExtension_AtLimits(t *testing.T) {
	t.Parallel()

	k, _ := setupDKGKeeper(t)
	deals := make([]types.Deal, maxItemsPerVote)
	responses := make([]types.Response, maxItemsPerVote)
	justifications := make([]types.Justification, maxItemsPerVote)

	bz, err := proto.Marshal(&types.Vote{
		Deals:          deals,
		Responses:      responses,
		Justifications: justifications,
	})
	require.NoError(t, err)

	votes, _, err := k.parseAndVerifyVoteExtension(bz)
	require.NoError(t, err)
	require.Len(t, votes, 1)
}

// TestDeduplicateDeals verifies that duplicate deals are removed by (dealerIndex, recipientIndex).
func TestDeduplicateDeals(t *testing.T) {
	t.Parallel()

	deals := []types.Deal{
		{Index: 1, RecipientIndex: 2},
		{Index: 1, RecipientIndex: 2}, // duplicate
		{Index: 1, RecipientIndex: 3}, // different recipient
		{Index: 2, RecipientIndex: 2}, // different dealer
	}
	result := deduplicateDeals(deals)
	require.Len(t, result, 3)
}

// TestDeduplicateResponses verifies that duplicate responses are removed by (responderIndex, dealerIndex).
func TestDeduplicateResponses(t *testing.T) {
	t.Parallel()

	responses := []types.Response{
		{Index: 1, VssResponse: &types.VSSResponse{Index: 2}},
		{Index: 1, VssResponse: &types.VSSResponse{Index: 2}}, // duplicate
		{Index: 1, VssResponse: &types.VSSResponse{Index: 3}}, // different dealer
		{Index: 2, VssResponse: &types.VSSResponse{Index: 2}}, // different responder
	}
	result := deduplicateResponses(responses)
	require.Len(t, result, 3)
}

// TestDeduplicateJustifications verifies that duplicate justifications are removed.
func TestDeduplicateJustifications(t *testing.T) {
	t.Parallel()

	justifications := []types.Justification{
		{Index: 1, VssJustification: &types.VSSJustification{PlainDeal: &types.PlainDeal{SecShare: &types.SecShare{I: 2}}}},
		{Index: 1, VssJustification: &types.VSSJustification{PlainDeal: &types.PlainDeal{SecShare: &types.SecShare{I: 2}}}}, // duplicate
		{Index: 1, VssJustification: &types.VSSJustification{PlainDeal: &types.PlainDeal{SecShare: &types.SecShare{I: 3}}}}, // different recipient
	}
	result := deduplicateJustifications(justifications)
	require.Len(t, result, 2)
}
