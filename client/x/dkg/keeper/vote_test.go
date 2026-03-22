package keeper

import (
	"testing"

	storetypes "cosmossdk.io/store/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/netconf"
)

// --- helpers ---

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

// --- parseAndVerifyVoteExtension ---

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

// --- VerifyVoteExtension ---

// TestVerifyVoteExtension_ValidVote verifies that VerifyVoteExtension accepts
// a valid vote extension.
func TestVerifyVoteExtension_ValidVote(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)
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

// --- aggregateVotes ---

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

// --- deduplicateDeals ---

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

// --- deduplicateResponses ---

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

// --- deduplicateJustifications ---

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

// --- votesFromExtension ---

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

// --- ExtendVote ---

// TestExtendVote_ReturnsSerializedVote verifies ExtendVote encodes the
// queued deals/responses/justifications into the vote extension.
// Note: uses package-level globals, so is NOT run in parallel.
func TestExtendVote_ReturnsSerializedVote(t *testing.T) {
	k, _, _, _ := setupDKGKeeperWithMocks(t)
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

// --- PrepareVotes ---

// TestPrepareVotes_NoVEHeight verifies that PrepareVotes returns an empty
// MsgAddDkgVote when vote extensions are not configured (veHeight == 0).
func TestPrepareVotes_NoVEHeight(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)
	sdkCtx := newTestSDKContext(t, "prepare_votes_no_ve")

	commit := abci.ExtendedCommitInfo{}
	msg, err := k.PrepareVotes(sdkCtx, commit, uint64(sdkCtx.BlockHeight()))
	require.NoError(t, err)
	require.NotNil(t, msg, "should return empty MsgAddDkgVote when VE not configured")

	dkgMsg, ok := msg.(*types.MsgAddDkgVote)
	require.True(t, ok)
	require.NotNil(t, dkgMsg.Vote)
	require.Empty(t, dkgMsg.Vote.Deals)
	require.Empty(t, dkgMsg.Vote.Responses)
	require.Empty(t, dkgMsg.Vote.Justifications)
}
