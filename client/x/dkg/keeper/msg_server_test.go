package keeper

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
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

// --- Determinism tests: DKG-enabled vs DKG-disabled GasUsed must match ---
//
// These tests verify that the gasless context prevents GasUsed divergence
// between TEE (isDKGSvcEnabled=true) and non-TEE (isDKGSvcEnabled=false)
// nodes, which would otherwise cause LastResultsHash mismatch and consensus
// failure.

// simulateAddVoteTxResult executes AddVote on the given keeper and returns
// an ExecTxResult with the measured GasUsed.
func simulateAddVoteTxResult(t *testing.T, k *Keeper, ctx sdk.Context, msg *types.MsgAddDkgVote) *abci.ExecTxResult {
	t.Helper()

	gasMeter := storetypes.NewInfiniteGasMeter()
	gasCtx := ctx.WithGasMeter(gasMeter)

	srv := NewMsgServerImpl(k)
	resp, err := srv.AddVote(gasCtx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)

	return &abci.ExecTxResult{
		Code:      0,
		GasWanted: 0,
		GasUsed:   int64(gasMeter.GasConsumed()),
	}
}

// TestAddVote_Determinism_Responses verifies that AddVote with responses
// produces identical GasUsed for DKG-enabled and DKG-disabled keepers.
// Responses trigger shouldProcessResponses which calls getLatestActiveDKGNetwork
// (a KV read) only on DKG-enabled nodes.
func TestAddVote_Determinism_Responses(t *testing.T) {
	t.Parallel()

	// Setup two keepers: one with DKG enabled, one without.
	kEnabled, _, _, ctxEnabled := setupDKGKeeperWithMocks(t)
	require.NoError(t, kEnabled.InitDKGService(
		t.TempDir(),
		common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		[32]byte{0x01},
	))

	kDisabled, _, _, ctxDisabled := setupDKGKeeperWithMocks(t)
	kDisabled.isDKGSvcEnabled = false

	// Set up identical DKG networks on both keepers.
	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}
	require.NoError(t, kEnabled.setDKGNetwork(ctxEnabled, network))
	require.NoError(t, kDisabled.setDKGNetwork(ctxDisabled, network))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote: &types.Vote{
			Responses: []types.Response{
				{Index: 1},
				{Index: 2},
			},
		},
	}

	resultEnabled := simulateAddVoteTxResult(t, kEnabled, sdk.UnwrapSDKContext(ctxEnabled), msg)
	resultDisabled := simulateAddVoteTxResult(t, kDisabled, sdk.UnwrapSDKContext(ctxDisabled), msg)

	require.Equal(t, resultEnabled.GasUsed, resultDisabled.GasUsed,
		"GasUsed must be identical between DKG-enabled and DKG-disabled keepers (responses)")
}

// TestAddVote_Determinism_ResponsesWithActiveRound verifies determinism when
// there is an active previous round (getLatestActiveDKGNetwork returns non-nil),
// which triggers additional KV reads on DKG-enabled nodes.
func TestAddVote_Determinism_ResponsesWithActiveRound(t *testing.T) {
	t.Parallel()

	kEnabled, _, _, ctxEnabled := setupDKGKeeperWithMocks(t)
	require.NoError(t, kEnabled.InitDKGService(
		t.TempDir(),
		common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		[32]byte{0x01},
	))

	kDisabled, _, _, ctxDisabled := setupDKGKeeperWithMocks(t)
	kDisabled.isDKGSvcEnabled = false

	// Set up an active previous round and a dealing round on both keepers.
	activeNetwork := &types.DKGNetwork{
		Round:        1,
		Total:        3,
		Threshold:    2,
		Stage:        types.DKGStageActive,
		ActiveValSet: []string{"0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"},
	}
	currentNetwork := &types.DKGNetwork{
		Round:     2,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	for _, setup := range []struct {
		k   *Keeper
		ctx sdk.Context
	}{
		{kEnabled, sdk.UnwrapSDKContext(ctxEnabled)},
		{kDisabled, sdk.UnwrapSDKContext(ctxDisabled)},
	} {
		require.NoError(t, setup.k.setDKGNetwork(setup.ctx, activeNetwork))
		require.NoError(t, setup.k.setLatestActiveRound(setup.ctx, activeNetwork))
		require.NoError(t, setup.k.setDKGNetwork(setup.ctx, currentNetwork))
	}

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote: &types.Vote{
			Responses: []types.Response{
				{Index: 1},
			},
		},
	}

	resultEnabled := simulateAddVoteTxResult(t, kEnabled, sdk.UnwrapSDKContext(ctxEnabled), msg)
	resultDisabled := simulateAddVoteTxResult(t, kDisabled, sdk.UnwrapSDKContext(ctxDisabled), msg)

	require.Equal(t, resultEnabled.GasUsed, resultDisabled.GasUsed,
		"GasUsed must be identical between DKG-enabled and DKG-disabled keepers (responses with active round)")
}

// TestAddVote_Determinism_DealsOnly verifies determinism when the vote
// contains only deals (no divergent KV reads expected, but we verify no
// regressions).
func TestAddVote_Determinism_DealsOnly(t *testing.T) {
	t.Parallel()

	kEnabled, _, _, ctxEnabled := setupDKGKeeperWithMocks(t)
	require.NoError(t, kEnabled.InitDKGService(
		t.TempDir(),
		common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		[32]byte{0x01},
	))

	kDisabled, _, _, ctxDisabled := setupDKGKeeperWithMocks(t)
	kDisabled.isDKGSvcEnabled = false

	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}
	require.NoError(t, kEnabled.setDKGNetwork(ctxEnabled, network))
	require.NoError(t, kDisabled.setDKGNetwork(ctxDisabled, network))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote: &types.Vote{
			Deals: []types.Deal{
				{Index: 1, RecipientIndex: 2},
				{Index: 2, RecipientIndex: 1},
			},
		},
	}

	resultEnabled := simulateAddVoteTxResult(t, kEnabled, sdk.UnwrapSDKContext(ctxEnabled), msg)
	resultDisabled := simulateAddVoteTxResult(t, kDisabled, sdk.UnwrapSDKContext(ctxDisabled), msg)

	require.Equal(t, resultEnabled.GasUsed, resultDisabled.GasUsed,
		"GasUsed must be identical between DKG-enabled and DKG-disabled keepers (deals only)")
}

// TestAddVote_Determinism_EmptyVote verifies determinism with an empty vote.
func TestAddVote_Determinism_EmptyVote(t *testing.T) {
	t.Parallel()

	kEnabled, _, _, ctxEnabled := setupDKGKeeperWithMocks(t)
	require.NoError(t, kEnabled.InitDKGService(
		t.TempDir(),
		common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		[32]byte{0x01},
	))

	kDisabled, _, _, ctxDisabled := setupDKGKeeperWithMocks(t)
	kDisabled.isDKGSvcEnabled = false

	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}
	require.NoError(t, kEnabled.setDKGNetwork(ctxEnabled, network))
	require.NoError(t, kDisabled.setDKGNetwork(ctxDisabled, network))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote:      &types.Vote{},
	}

	resultEnabled := simulateAddVoteTxResult(t, kEnabled, sdk.UnwrapSDKContext(ctxEnabled), msg)
	resultDisabled := simulateAddVoteTxResult(t, kDisabled, sdk.UnwrapSDKContext(ctxDisabled), msg)

	require.Equal(t, resultEnabled.GasUsed, resultDisabled.GasUsed,
		"GasUsed must be identical between DKG-enabled and DKG-disabled keepers (empty vote)")
}
