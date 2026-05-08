package keeper

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"
)

// resetPendingIncoming clears all pending incoming queues for test isolation.
func resetPendingIncoming() {
	pendingIncomingDealsMu.Lock()
	pendingIncomingDeals = nil
	pendingIncomingDealsMu.Unlock()

	pendingIncomingResponsesMu.Lock()
	pendingIncomingResponses = nil
	pendingIncomingResponsesMu.Unlock()

	pendingIncomingJustificationsMu.Lock()
	pendingIncomingJustifications = nil
	pendingIncomingJustificationsMu.Unlock()
}

// --- cachePendingDeals ---

func TestCachePendingDeals(t *testing.T) {
	tests := []struct {
		name     string
		items    []pendingDeal
		preload  int
		expected int // expected total after caching
	}{
		{
			name: "cache deals when empty",
			items: []pendingDeal{
				{deal: types.Deal{Index: 0}},
				{deal: types.Deal{Index: 1}},
			},
			preload:  0,
			expected: 2,
		},
		{
			name: "cache with partial capacity remaining",
			items: []pendingDeal{
				{deal: types.Deal{Index: 0}},
				{deal: types.Deal{Index: 1}},
				{deal: types.Deal{Index: 2}},
			},
			preload:  maxPendingIncoming - 2,
			expected: maxPendingIncoming,
		},
		{
			name: "cache at full capacity (no-op)",
			items: []pendingDeal{
				{deal: types.Deal{Index: 0}},
			},
			preload:  maxPendingIncoming,
			expected: maxPendingIncoming,
		},
		{
			name:     "empty items list",
			items:    nil,
			preload:  0,
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetPendingIncoming()
			defer resetPendingIncoming()

			// Pre-load
			pendingIncomingDealsMu.Lock()
			for i := range tc.preload {
				pendingIncomingDeals = append(pendingIncomingDeals, pendingDeal{deal: types.Deal{Index: uint32(i + 100)}})
			}
			pendingIncomingDealsMu.Unlock()

			cachePendingDeals(context.Background(), tc.items)

			pendingIncomingDealsMu.Lock()
			got := len(pendingIncomingDeals)
			pendingIncomingDealsMu.Unlock()
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestCachePendingDeals_PreservesRetryCount(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	items := []pendingDeal{
		{deal: types.Deal{Index: 1}, retryCount: 3},
		{deal: types.Deal{Index: 2}, retryCount: 5},
	}

	cachePendingDeals(context.Background(), items)

	pendingIncomingDealsMu.Lock()
	defer pendingIncomingDealsMu.Unlock()

	require.Len(t, pendingIncomingDeals, 2)
	require.Equal(t, 3, pendingIncomingDeals[0].retryCount, "retryCount must be preserved")
	require.Equal(t, 5, pendingIncomingDeals[1].retryCount, "retryCount must be preserved")
}

// TestCachePending_DropsAtCapacity verifies that overflow items are dropped
// (not appended past capacity) when the queue is saturated. Earlier items
// in the input slice are kept; later ones are discarded — the queue acts as
// a fixed-size buffer that does NOT silently grow without bound.
func TestCachePending_DropsAtCapacity(t *testing.T) {
	t.Run("deals partial overflow", func(t *testing.T) {
		resetPendingIncoming()
		defer resetPendingIncoming()

		pendingIncomingDealsMu.Lock()
		for i := range maxPendingIncoming - 2 {
			pendingIncomingDeals = append(pendingIncomingDeals, pendingDeal{deal: types.Deal{Index: uint32(i + 100)}})
		}
		pendingIncomingDealsMu.Unlock()

		// Try to cache 5 items when only 2 slots remain.
		items := []pendingDeal{
			{deal: types.Deal{Index: 1}},
			{deal: types.Deal{Index: 2}},
			{deal: types.Deal{Index: 3}},
			{deal: types.Deal{Index: 4}},
			{deal: types.Deal{Index: 5}},
		}
		cachePendingDeals(context.Background(), items)

		pendingIncomingDealsMu.Lock()
		defer pendingIncomingDealsMu.Unlock()
		require.Equal(t, maxPendingIncoming, len(pendingIncomingDeals))
		// Earlier items are kept, later ones dropped.
		last := pendingIncomingDeals[len(pendingIncomingDeals)-2:]
		require.Equal(t, uint32(1), last[0].deal.Index)
		require.Equal(t, uint32(2), last[1].deal.Index)
	})

	t.Run("deals full overflow", func(t *testing.T) {
		resetPendingIncoming()
		defer resetPendingIncoming()

		pendingIncomingDealsMu.Lock()
		for i := range maxPendingIncoming {
			pendingIncomingDeals = append(pendingIncomingDeals, pendingDeal{deal: types.Deal{Index: uint32(i + 100)}})
		}
		pendingIncomingDealsMu.Unlock()

		// Queue is at capacity — every item must be dropped.
		items := []pendingDeal{
			{deal: types.Deal{Index: 1}},
			{deal: types.Deal{Index: 2}},
		}
		cachePendingDeals(context.Background(), items)

		pendingIncomingDealsMu.Lock()
		defer pendingIncomingDealsMu.Unlock()
		require.Equal(t, maxPendingIncoming, len(pendingIncomingDeals))
		for i, p := range pendingIncomingDeals {
			require.Equal(t, uint32(i+100), p.deal.Index, "pre-existing items must not be evicted")
		}
	})

	t.Run("responses full overflow", func(t *testing.T) {
		resetPendingIncoming()
		defer resetPendingIncoming()

		pendingIncomingResponsesMu.Lock()
		for i := range maxPendingIncoming {
			pendingIncomingResponses = append(pendingIncomingResponses, pendingResponse{response: types.Response{Index: uint32(i + 100)}})
		}
		pendingIncomingResponsesMu.Unlock()

		cachePendingResponses(context.Background(), []pendingResponse{
			{response: types.Response{Index: 1}},
		})

		pendingIncomingResponsesMu.Lock()
		defer pendingIncomingResponsesMu.Unlock()
		require.Equal(t, maxPendingIncoming, len(pendingIncomingResponses))
	})

	t.Run("justifications full overflow", func(t *testing.T) {
		resetPendingIncoming()
		defer resetPendingIncoming()

		pendingIncomingJustificationsMu.Lock()
		for i := range maxPendingIncoming {
			pendingIncomingJustifications = append(pendingIncomingJustifications, pendingJustification{justification: types.Justification{Index: uint32(i + 100)}})
		}
		pendingIncomingJustificationsMu.Unlock()

		cachePendingJustifications(context.Background(), []pendingJustification{
			{justification: types.Justification{Index: 1}},
		})

		pendingIncomingJustificationsMu.Lock()
		defer pendingIncomingJustificationsMu.Unlock()
		require.Equal(t, maxPendingIncoming, len(pendingIncomingJustifications))
	})
}

// --- cachePendingResponses ---

func TestCachePendingResponses(t *testing.T) {
	tests := []struct {
		name     string
		items    []pendingResponse
		preload  int // number of pre-loaded responses
		expected int // expected total after caching
	}{
		{
			name:     "cache responses when empty",
			items:    []pendingResponse{{response: types.Response{Index: 1}}, {response: types.Response{Index: 2}}},
			preload:  0,
			expected: 2,
		},
		{
			name:     "cache with partial capacity remaining",
			items:    []pendingResponse{{response: types.Response{Index: 1}}, {response: types.Response{Index: 2}}, {response: types.Response{Index: 3}}},
			preload:  maxPendingIncoming - 2,
			expected: maxPendingIncoming,
		},
		{
			name:     "cache at full capacity (no-op)",
			items:    []pendingResponse{{response: types.Response{Index: 1}}},
			preload:  maxPendingIncoming,
			expected: maxPendingIncoming,
		},
		{
			name:     "empty items list",
			items:    nil,
			preload:  0,
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetPendingIncoming()
			defer resetPendingIncoming()

			// Pre-load
			pendingIncomingResponsesMu.Lock()
			for i := range tc.preload {
				pendingIncomingResponses = append(pendingIncomingResponses, pendingResponse{response: types.Response{Index: uint32(i + 100)}})
			}
			pendingIncomingResponsesMu.Unlock()

			cachePendingResponses(context.Background(), tc.items)

			pendingIncomingResponsesMu.Lock()
			got := len(pendingIncomingResponses)
			pendingIncomingResponsesMu.Unlock()
			require.Equal(t, tc.expected, got)
		})
	}
}

// --- cachePendingJustifications ---

func TestCachePendingJustifications(t *testing.T) {
	tests := []struct {
		name     string
		items    []pendingJustification
		preload  int
		expected int
	}{
		{
			name:     "cache justifications when empty",
			items:    []pendingJustification{{justification: types.Justification{Index: 1}}, {justification: types.Justification{Index: 2}}},
			preload:  0,
			expected: 2,
		},
		{
			name:     "cache with partial capacity",
			items:    []pendingJustification{{justification: types.Justification{Index: 1}}, {justification: types.Justification{Index: 2}}, {justification: types.Justification{Index: 3}}},
			preload:  maxPendingIncoming - 1,
			expected: maxPendingIncoming,
		},
		{
			name:     "cache at full capacity (no-op)",
			items:    []pendingJustification{{justification: types.Justification{Index: 1}}},
			preload:  maxPendingIncoming,
			expected: maxPendingIncoming,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetPendingIncoming()
			defer resetPendingIncoming()

			pendingIncomingJustificationsMu.Lock()
			for i := range tc.preload {
				pendingIncomingJustifications = append(pendingIncomingJustifications, pendingJustification{justification: types.Justification{Index: uint32(i + 100)}})
			}
			pendingIncomingJustificationsMu.Unlock()

			cachePendingJustifications(context.Background(), tc.items)

			pendingIncomingJustificationsMu.Lock()
			got := len(pendingIncomingJustifications)
			pendingIncomingJustificationsMu.Unlock()
			require.Equal(t, tc.expected, got)
		})
	}
}

// --- drain functions ---

func TestDrainPendingIncomingDeals(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	// Pre-populate
	pendingIncomingDealsMu.Lock()
	pendingIncomingDeals = []pendingDeal{{deal: types.Deal{Index: 1}}, {deal: types.Deal{Index: 2}}}
	pendingIncomingDealsMu.Unlock()

	drained := drainPendingIncomingDeals()
	require.Len(t, drained, 2)
	require.Equal(t, uint32(1), drained[0].deal.Index)
	require.Equal(t, uint32(2), drained[1].deal.Index)

	// After drain, queue should be empty
	pendingIncomingDealsMu.Lock()
	require.Nil(t, pendingIncomingDeals)
	pendingIncomingDealsMu.Unlock()
}

func TestDrainPendingIncomingDeals_Empty(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	drained := drainPendingIncomingDeals()
	require.Nil(t, drained)
}

func TestDrainPendingIncomingResponses(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	pendingIncomingResponsesMu.Lock()
	pendingIncomingResponses = []pendingResponse{{response: types.Response{Index: 3}}, {response: types.Response{Index: 4}}}
	pendingIncomingResponsesMu.Unlock()

	drained := drainPendingIncomingResponses()
	require.Len(t, drained, 2)

	pendingIncomingResponsesMu.Lock()
	require.Nil(t, pendingIncomingResponses)
	pendingIncomingResponsesMu.Unlock()
}

func TestDrainPendingIncomingJustifications(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	pendingIncomingJustificationsMu.Lock()
	pendingIncomingJustifications = []pendingJustification{{justification: types.Justification{Index: 5}}}
	pendingIncomingJustificationsMu.Unlock()

	drained := drainPendingIncomingJustifications()
	require.Len(t, drained, 1)

	pendingIncomingJustificationsMu.Lock()
	require.Nil(t, pendingIncomingJustifications)
	pendingIncomingJustificationsMu.Unlock()
}

// --- flushPendingIncoming ---

func TestFlushPendingIncoming(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	// Pre-populate all queues
	pendingIncomingDealsMu.Lock()
	pendingIncomingDeals = []pendingDeal{{deal: types.Deal{Index: 1}}}
	pendingIncomingDealsMu.Unlock()

	pendingIncomingResponsesMu.Lock()
	pendingIncomingResponses = []pendingResponse{{response: types.Response{Index: 2}}}
	pendingIncomingResponsesMu.Unlock()

	pendingIncomingJustificationsMu.Lock()
	pendingIncomingJustifications = []pendingJustification{{justification: types.Justification{Index: 3}}}
	pendingIncomingJustificationsMu.Unlock()

	flushPendingIncoming()

	pendingIncomingDealsMu.Lock()
	require.Nil(t, pendingIncomingDeals)
	pendingIncomingDealsMu.Unlock()

	pendingIncomingResponsesMu.Lock()
	require.Nil(t, pendingIncomingResponses)
	pendingIncomingResponsesMu.Unlock()

	pendingIncomingJustificationsMu.Lock()
	require.Nil(t, pendingIncomingJustifications)
	pendingIncomingJustificationsMu.Unlock()
}

// --- Tests merged from dkg_svc_dealing_guards_test.go ---

// NOTE: These tests are NOT parallel because they share the package-level dkgSvcRound atomic.

func TestHandleDKGDealing_WrongStage(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageRegistration, // Not dealing
	}

	session := &types.DKGSession{
		Round: 1,
		Phase: types.PhaseInitialized,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	k.handleDKGDealing(ctx, dkgNetwork, true)

	got, err := k.stateManager.GetSession(1)
	require.NoError(t, err)
	require.Equal(t, types.PhaseInitialized, got.Phase, "phase should not change when stage is wrong")
}

func TestHandleDKGDealing_ShouldNotDeal(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	dkgNetwork := &types.DKGNetwork{
		Round: 2,
		Stage: types.DKGStageDealing,
	}

	session := &types.DKGSession{
		Round: 2,
		Phase: types.PhaseInitialized,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	k.handleDKGDealing(ctx, dkgNetwork, false) // shouldDeal=false

	got, err := k.stateManager.GetSession(2)
	require.NoError(t, err)
	require.Equal(t, types.PhaseInitialized, got.Phase, "phase should not change when shouldDeal is false")
}

func TestHandleDKGDealing_WrongPhase_MarksFailed(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	dkgNetwork := &types.DKGNetwork{
		Round: 3,
		Stage: types.DKGStageDealing,
	}

	session := &types.DKGSession{
		Round: 3,
		Phase: types.PhaseFinalized, // Wrong phase — should be Initialized
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	k.handleDKGDealing(ctx, dkgNetwork, true)

	got, err := k.stateManager.GetSession(3)
	require.NoError(t, err)
	require.Equal(t, types.PhaseFailed, got.Phase, "non-initialized phase should be marked failed")
}

func TestHandleDKGDealing_NoSession(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	dkgNetwork := &types.DKGNetwork{
		Round: 99,
		Stage: types.DKGStageDealing,
	}

	// Should not panic when session doesn't exist
	k.handleDKGDealing(ctx, dkgNetwork, true)
}

func TestHandleDKGDealing_DuplicateAcquire(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Pre-acquire lock for round 4
	dkgSvcRound.Store(4)

	dkgNetwork := &types.DKGNetwork{
		Round: 4,
		Stage: types.DKGStageDealing,
	}

	session := &types.DKGSession{
		Round: 4,
		Phase: types.PhaseInitialized,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	k.handleDKGDealing(ctx, dkgNetwork, true)

	got, err := k.stateManager.GetSession(4)
	require.NoError(t, err)
	require.Equal(t, types.PhaseInitialized, got.Phase, "lock dedup should prevent processing")
}

// --- Tests merged from dkg_svc_process_deals_test.go ---

// --- handleDKGProcessDeals ---

func TestHandleDKGProcessDeals_NotInValSet(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	k.validatorEVMAddr = testValidatorAddr

	dkgNetwork := &types.DKGNetwork{
		Round:        1,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{"0xother1"}, // testValidator not in set
	}

	deals := []types.Deal{{Index: 0, RecipientIndex: 0}}

	// Should skip because validator is not in current round set
	k.handleDKGProcessDeals(ctx, dkgNetwork, wrapDeals(deals))
}

func TestHandleDKGProcessDeals_NoSession(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	k.validatorEVMAddr = testValidatorAddr

	dkgNetwork := &types.DKGNetwork{
		Round:        1,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{testValidatorAddr},
	}

	deals := []types.Deal{{Index: 0, RecipientIndex: 0}}

	// Should log error but not panic
	k.handleDKGProcessDeals(ctx, dkgNetwork, wrapDeals(deals))
}

func TestHandleDKGProcessDeals_WrongPhase(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	k.validatorEVMAddr = testValidatorAddr

	session := &types.DKGSession{
		Round: 2,
		Phase: types.PhaseFinalized, // Wrong phase
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round:        2,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{testValidatorAddr},
	}

	deals := []types.Deal{{Index: 0, RecipientIndex: 0}}

	// Should skip — wrong phase
	k.handleDKGProcessDeals(ctx, dkgNetwork, wrapDeals(deals))

	got, err := k.stateManager.GetSession(2)
	require.NoError(t, err)
	require.Equal(t, types.PhaseFinalized, got.Phase, "phase should remain unchanged")
}

func TestHandleDKGProcessDeals_NoMatchingDeals(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	k.validatorEVMAddr = testValidatorAddr
	router := NewKernelRouter(nil, nil)
	k.kernelRouter = router

	session := &types.DKGSession{
		Round:          3,
		Phase:          types.PhaseDealing,
		Index:          5,
		CodeCommitment: []byte("test-cc"),
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round:        3,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{testValidatorAddr},
	}

	// Deals addressed to other validators (recipientIndex != sessionIndex-1 = 4)
	deals := []types.Deal{
		{Index: 0, RecipientIndex: 0}, // not for index 4
		{Index: 1, RecipientIndex: 1}, // not for index 4
	}

	// Should succeed without calling kernel (no matching deals)
	k.handleDKGProcessDeals(ctx, dkgNetwork, wrapDeals(deals))
}

// --- handleDKGProcessResponses ---

func TestHandleDKGProcessResponses_ShouldNotProcess(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageDealing,
	}

	responses := []types.Response{{Index: 1}}

	// shouldProcess=false → skip
	k.handleDKGProcessResponses(ctx, dkgNetwork, wrapResponses(responses), false)
}

func TestHandleDKGProcessResponses_NoSession(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageDealing,
	}

	responses := []types.Response{{Index: 1}}

	// Should log error but not panic
	k.handleDKGProcessResponses(ctx, dkgNetwork, wrapResponses(responses), true)
}

func TestHandleDKGProcessResponses_WrongPhase(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	session := &types.DKGSession{
		Round: 2,
		Phase: types.PhaseFinalized, // Wrong phase
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round: 2,
		Stage: types.DKGStageDealing,
	}

	responses := []types.Response{{Index: 1}}

	k.handleDKGProcessResponses(ctx, dkgNetwork, wrapResponses(responses), true)

	got, err := k.stateManager.GetSession(2)
	require.NoError(t, err)
	require.Equal(t, types.PhaseFinalized, got.Phase, "phase should remain unchanged")
}

func TestHandleDKGProcessResponses_AllSelfResponses_Filtered(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	router := NewKernelRouter(nil, nil)
	k.kernelRouter = router

	session := &types.DKGSession{
		Round:          3,
		Phase:          types.PhaseDealing,
		Index:          2, // 1-based
		CodeCommitment: []byte("test-cc"),
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round: 3,
		Stage: types.DKGStageDealing,
	}

	// Response from self (VssResponse.Index = session.Index - 1 = 1, which is 0-based)
	responses := []types.Response{
		{Index: 1, VssResponse: &types.VSSResponse{Index: 1}}, // self-response
	}

	// All responses are from self → filtered out → no kernel call
	k.handleDKGProcessResponses(ctx, dkgNetwork, wrapResponses(responses), true)
}

// --- Tests merged from dkg_svc_justifications_test.go ---

func TestHandleDKGProcessJustifications_NoSession(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	dkgNetwork := &types.DKGNetwork{
		Round: 99,
		Stage: types.DKGStageDealing,
	}

	justifications := []types.Justification{{Index: 1}}

	// Should log error but not panic when session doesn't exist
	k.handleDKGProcessJustifications(ctx, dkgNetwork, wrapJustifications(justifications))
}

func TestHandleDKGProcessJustifications_NoKernelClient(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	router := NewKernelRouter(nil, nil)
	k.kernelRouter = router

	session := &types.DKGSession{
		Round:          1,
		Phase:          types.PhaseDealing,
		CodeCommitment: []byte("nonexistent-cc"),
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageDealing,
	}

	justifications := []types.Justification{{Index: 1}}

	// Should cache justifications when kernel client not found
	resetPendingIncoming()
	defer resetPendingIncoming()

	k.handleDKGProcessJustifications(ctx, dkgNetwork, wrapJustifications(justifications))

	// Justifications should be cached for retry
	pendingIncomingJustificationsMu.Lock()
	cached := len(pendingIncomingJustifications)
	pendingIncomingJustificationsMu.Unlock()
	require.Equal(t, 1, cached, "justifications should be cached when kernel client unavailable")
}

// --- Tests merged from dkg_svc_should_deal_test.go ---

const testValidatorAddr = "0xabcdef1234567890abcdef1234567890abcdef12"

// --- shouldDeal ---

func TestShouldDeal(t *testing.T) {
	tests := []struct {
		name             string
		validatorAddr    string
		activeValSet     []string
		prevActive       *types.DKGNetwork // nil = no previous active round
		expectedShouldDl bool
	}{
		{
			name:             "first round: validator in current set should deal",
			validatorAddr:    testValidatorAddr,
			activeValSet:     []string{testValidatorAddr, "0xother"},
			prevActive:       nil,
			expectedShouldDl: true,
		},
		{
			name:             "first round: validator not in current set should not deal",
			validatorAddr:    testValidatorAddr,
			activeValSet:     []string{"0xother1", "0xother2"},
			prevActive:       nil,
			expectedShouldDl: false,
		},
		{
			name:          "resharing: validator in previous set should deal",
			validatorAddr: testValidatorAddr,
			activeValSet:  []string{"0xnewval1", "0xnewval2"},
			prevActive: &types.DKGNetwork{
				Round:        1,
				Stage:        types.DKGStageActive,
				ActiveValSet: []string{testValidatorAddr, "0xoldval"},
			},
			expectedShouldDl: true,
		},
		{
			name:          "resharing: validator not in previous set should not deal",
			validatorAddr: testValidatorAddr,
			activeValSet:  []string{testValidatorAddr, "0xnewval"},
			prevActive: &types.DKGNetwork{
				Round:        1,
				Stage:        types.DKGStageActive,
				ActiveValSet: []string{"0xoldval1", "0xoldval2"},
			},
			expectedShouldDl: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			k, _, _, ctx := setupDKGKeeperWithMocks(t)
			k.validatorEVMAddr = tc.validatorAddr

			// Set up previous active round if present
			if tc.prevActive != nil {
				require.NoError(t, k.setDKGNetwork(ctx, tc.prevActive))
				require.NoError(t, k.setLatestActiveRound(ctx, tc.prevActive))
			}

			dkgNetwork := &types.DKGNetwork{
				Round:        2,
				Stage:        types.DKGStageDealing,
				ActiveValSet: tc.activeValSet,
			}

			shouldDl, err := k.shouldDeal(ctx, dkgNetwork)
			require.NoError(t, err)
			require.Equal(t, tc.expectedShouldDl, shouldDl)
		})
	}
}

// --- shouldProcessResponses ---

func TestShouldProcessResponses(t *testing.T) {
	tests := []struct {
		name           string
		validatorAddr  string
		activeValSet   []string
		prevActive     *types.DKGNetwork
		expectedResult bool
	}{
		{
			name:           "first round: validator in current set should process",
			validatorAddr:  testValidatorAddr,
			activeValSet:   []string{testValidatorAddr, "0xother"},
			prevActive:     nil,
			expectedResult: true,
		},
		{
			name:           "first round: validator not in current set should not process",
			validatorAddr:  testValidatorAddr,
			activeValSet:   []string{"0xother1", "0xother2"},
			prevActive:     nil,
			expectedResult: false,
		},
		{
			name:          "resharing: validator in current set only",
			validatorAddr: testValidatorAddr,
			activeValSet:  []string{testValidatorAddr, "0xnewval"},
			prevActive: &types.DKGNetwork{
				Round:        1,
				Stage:        types.DKGStageActive,
				ActiveValSet: []string{"0xoldval1", "0xoldval2"},
			},
			expectedResult: true,
		},
		{
			name:          "resharing: validator in previous set only",
			validatorAddr: testValidatorAddr,
			activeValSet:  []string{"0xnewval1", "0xnewval2"},
			prevActive: &types.DKGNetwork{
				Round:        1,
				Stage:        types.DKGStageActive,
				ActiveValSet: []string{testValidatorAddr, "0xoldval"},
			},
			expectedResult: true,
		},
		{
			name:          "resharing: validator in both sets",
			validatorAddr: testValidatorAddr,
			activeValSet:  []string{testValidatorAddr, "0xnewval"},
			prevActive: &types.DKGNetwork{
				Round:        1,
				Stage:        types.DKGStageActive,
				ActiveValSet: []string{testValidatorAddr, "0xoldval"},
			},
			expectedResult: true,
		},
		{
			name:          "resharing: validator in neither set",
			validatorAddr: testValidatorAddr,
			activeValSet:  []string{"0xnewval1", "0xnewval2"},
			prevActive: &types.DKGNetwork{
				Round:        1,
				Stage:        types.DKGStageActive,
				ActiveValSet: []string{"0xoldval1", "0xoldval2"},
			},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			k, _, _, ctx := setupDKGKeeperWithMocks(t)
			k.validatorEVMAddr = tc.validatorAddr

			if tc.prevActive != nil {
				require.NoError(t, k.setDKGNetwork(ctx, tc.prevActive))
				require.NoError(t, k.setLatestActiveRound(ctx, tc.prevActive))
			}

			dkgNetwork := &types.DKGNetwork{
				Round:        2,
				Stage:        types.DKGStageDealing,
				ActiveValSet: tc.activeValSet,
			}

			result, err := k.shouldProcessResponses(ctx, dkgNetwork)
			require.NoError(t, err)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

// --- Tests merged from dkg_svc_reprocess_test.go ---

func TestReprocessPendingIncomingData_NoPending(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	k := setupKeeperWithStateManager(t)
	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageDealing,
	}

	// Should be a no-op when no pending data exists
	k.reprocessPendingIncomingData(dkgNetwork)
}

func TestReprocessPendingIncomingData_WrongStage_FlushesData(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	// Pre-populate pending data
	pendingIncomingDealsMu.Lock()
	pendingIncomingDeals = []pendingDeal{{deal: types.Deal{Index: 1}}}
	pendingIncomingDealsMu.Unlock()

	pendingIncomingResponsesMu.Lock()
	pendingIncomingResponses = []pendingResponse{{response: types.Response{Index: 2}}}
	pendingIncomingResponsesMu.Unlock()

	k := setupKeeperWithStateManager(t)
	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageFinalization, // Not dealing
	}

	k.reprocessPendingIncomingData(dkgNetwork)

	// Pending data should be flushed when stage is not dealing
	pendingIncomingDealsMu.Lock()
	require.Nil(t, pendingIncomingDeals, "deals should be flushed when stage is not dealing")
	pendingIncomingDealsMu.Unlock()

	pendingIncomingResponsesMu.Lock()
	require.Nil(t, pendingIncomingResponses, "responses should be flushed when stage is not dealing")
	pendingIncomingResponsesMu.Unlock()
}

func TestReprocessPendingIncomingData_NilKernelRouter(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	pendingIncomingDealsMu.Lock()
	pendingIncomingDeals = []pendingDeal{{deal: types.Deal{Index: 1}}}
	pendingIncomingDealsMu.Unlock()

	k := &Keeper{kernelRouter: nil}
	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageDealing,
	}

	// Should be a no-op when kernel router is nil
	k.reprocessPendingIncomingData(dkgNetwork)

	// Data should still be pending (not flushed, not processed)
	pendingIncomingDealsMu.Lock()
	require.Len(t, pendingIncomingDeals, 1, "data should remain when kernel router is nil")
	pendingIncomingDealsMu.Unlock()
}

func TestReprocessPendingIncomingData_NoClients(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	pendingIncomingDealsMu.Lock()
	pendingIncomingDeals = []pendingDeal{{deal: types.Deal{Index: 1}}}
	pendingIncomingDealsMu.Unlock()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}
	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageDealing,
	}

	// No clients connected, TryReconnect should be called but since
	// no endpoints are configured, nothing happens
	k.reprocessPendingIncomingData(dkgNetwork)

	// Data should still be pending
	pendingIncomingDealsMu.Lock()
	require.Len(t, pendingIncomingDeals, 1, "data should remain when no clients available")
	pendingIncomingDealsMu.Unlock()
}

// --- handleDKGDealing: upgrade round uses OldCodeCommitment (gap 13) ---

// TestHandleDKGDealing_UpgradeUsesOldCC verifies that for upgrade resharing,
// handleDKGDealing uses the session's OldCodeCommitment (not CodeCommitment) for GenerateDeals.
func TestHandleDKGDealing_UpgradeUsesOldCC(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	oldCC := []byte("old-deal-cc")
	newCC := []byte("new-deal-cc")
	mockOldKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockNewKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(oldCC, mockOldKernel)
	router.RegisterClient(newCC, mockNewKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	session := &types.DKGSession{
		Round:             5,
		Phase:             types.PhaseInitialized,
		CodeCommitment:    newCC,
		OldCodeCommitment: oldCC,
		IsUpgrade:         true,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Upgrade resharing: GenerateDeals should be called on the OLD kernel binary.
	mockOldKernel.EXPECT().GenerateDeals(gomock.Any(), gomock.Any()).Return(
		&types.GenerateDealsResponse{Deals: []types.Deal{{Index: 0}}}, nil,
	)
	// mockNewKernel should NOT be called for deal generation.

	dkgNetwork := &types.DKGNetwork{
		Round:     5,
		Stage:     types.DKGStageDealing,
		IsUpgrade: true,
	}

	k.FlushAllQueues()
	k.handleDKGDealing(ctx, dkgNetwork, true)

	got, err := sm.GetSession(5)
	require.NoError(t, err)
	require.Equal(t, types.PhaseDealing, got.Phase, "upgrade dealing should advance to PhaseDealing")
}

// TestHandleDKGDealing_GenerateDealsError_MarksFailed verifies that handleDKGDealing
// marks the session as failed when GenerateDeals returns an error.
func TestHandleDKGDealing_GenerateDealsError_MarksFailed(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("fail-deal-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	session := &types.DKGSession{
		Round:          6,
		Phase:          types.PhaseInitialized,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Kernel returns error on all retry attempts.
	mockKernel.EXPECT().GenerateDeals(gomock.Any(), gomock.Any()).Return(nil, errSentinel).Times(retryAttempts)

	dkgNetwork := &types.DKGNetwork{Round: 6, Stage: types.DKGStageDealing}

	k.handleDKGDealing(ctx, dkgNetwork, true)

	got, err := sm.GetSession(6)
	require.NoError(t, err)
	require.Equal(t, types.PhaseFailed, got.Phase, "failed GenerateDeals should mark session as failed")
}

// --- handleDKGProcessDeals: session.Index > 0 && deal.RecipientIndex filtering (gap 14) ---

// TestHandleDKGProcessDeals_SessionIndexFiltering verifies that handleDKGProcessDeals
// only forwards deals where RecipientIndex == session.Index - 1 (0-based kyber index).
func TestHandleDKGProcessDeals_SessionIndexFiltering(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("process-deal-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		validatorEVMAddr: testValidatorAddr,
	}

	session := &types.DKGSession{
		Round:          7,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
		Index:          2, // 1-based: recipientIndex to accept = 1 (0-based)
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Deal addressed to this validator (RecipientIndex=1 == Index-1=1)
	myDeal := types.Deal{Index: 0, RecipientIndex: 1}
	// Deal addressed to a different validator
	otherDeal := types.Deal{Index: 0, RecipientIndex: 0}

	// ProcessDeals should only forward myDeal (one deal addressed to us).
	mockKernel.EXPECT().ProcessDeals(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, req *types.ProcessDealsRequest, _ ...grpc.CallOption) (*types.ProcessDealsResponse, error) {
			require.Len(t, req.Deals, 1, "should only forward deals addressed to this validator")
			require.Equal(t, uint32(1), req.Deals[0].RecipientIndex)
			return &types.ProcessDealsResponse{Responses: []types.Response{{Index: 0}}}, nil
		},
	)

	dkgNetwork := &types.DKGNetwork{
		Round:        7,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{testValidatorAddr},
	}

	k.handleDKGProcessDeals(ctx, dkgNetwork, wrapDeals([]types.Deal{myDeal, otherDeal}))
}

// TestHandleDKGProcessDeals_ProcessDealsError_CachesDeals verifies that failed
// ProcessDeals calls cache the unprocessed deals for retry.
func TestHandleDKGProcessDeals_ProcessDealsError_CachesDeals(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("fail-process-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		validatorEVMAddr: testValidatorAddr,
	}

	session := &types.DKGSession{
		Round:          8,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
		Index:          1, // accept RecipientIndex=0
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	myDeal := types.Deal{Index: 0, RecipientIndex: 0}

	// Kernel returns error on all retries.
	mockKernel.EXPECT().ProcessDeals(gomock.Any(), gomock.Any()).Return(nil, errSentinel).Times(retryAttempts)

	dkgNetwork := &types.DKGNetwork{
		Round:        8,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{testValidatorAddr},
	}

	resetPendingIncoming()
	defer resetPendingIncoming()

	k.handleDKGProcessDeals(ctx, dkgNetwork, wrapDeals([]types.Deal{myDeal}))

	pendingIncomingDealsMu.Lock()
	cached := len(pendingIncomingDeals)
	pendingIncomingDealsMu.Unlock()
	require.Equal(t, 1, cached, "failed deal should be cached for retry")
}

// --- handleDKGProcessResponses: upgrade CC filtering (gap 15) ---

// TestHandleDKGProcessResponses_UpgradeCCFiltering verifies that for upgrade resharing,
// handleDKGProcessResponses sends responses to BOTH old and new kernel binaries.
func TestHandleDKGProcessResponses_UpgradeCCFiltering(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	oldCC := []byte("old-resp-cc")
	newCC := []byte("new-resp-cc")
	mockOldKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockNewKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(oldCC, mockOldKernel)
	router.RegisterClient(newCC, mockNewKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	session := &types.DKGSession{
		Round:             9,
		Phase:             types.PhaseDealing,
		CodeCommitment:    newCC,
		OldCodeCommitment: oldCC,
		IsUpgrade:         true,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	resp := types.Response{Index: 0, VssResponse: &types.VSSResponse{Index: 1}}

	// Both old and new kernel should receive ProcessResponses.
	mockOldKernel.EXPECT().ProcessResponses(gomock.Any(), gomock.Any()).Return(
		&types.ProcessResponsesResponse{}, nil,
	)
	mockNewKernel.EXPECT().ProcessResponses(gomock.Any(), gomock.Any()).Return(
		&types.ProcessResponsesResponse{}, nil,
	)

	dkgNetwork := &types.DKGNetwork{
		Round:     9,
		Stage:     types.DKGStageDealing,
		IsUpgrade: true,
	}

	k.handleDKGProcessResponses(ctx, dkgNetwork, wrapResponses([]types.Response{resp}), true)
}

// TestHandleDKGProcessResponses_ResponseIndexFiltering verifies that responses
// from self (VssResponse.Index == session.Index - 1) are filtered out.
func TestHandleDKGProcessResponses_ResponseIndexFiltering(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("filter-resp-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	session := &types.DKGSession{
		Round:          10,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
		Index:          2, // self index: 1 (0-based) = VssResponse.Index == 1 should be filtered
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	selfResp := types.Response{Index: 0, VssResponse: &types.VSSResponse{Index: 1}}  // self response (filtered)
	otherResp := types.Response{Index: 1, VssResponse: &types.VSSResponse{Index: 0}} // other response

	// After filtering self, only otherResp remains → ProcessResponses called with 1 response.
	mockKernel.EXPECT().ProcessResponses(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, req *types.ProcessResponsesRequest, _ ...grpc.CallOption) (*types.ProcessResponsesResponse, error) {
			require.Len(t, req.Responses, 1, "self response should be filtered out")
			require.Equal(t, uint32(0), req.Responses[0].VssResponse.Index)
			return &types.ProcessResponsesResponse{}, nil
		},
	)

	dkgNetwork := &types.DKGNetwork{Round: 10, Stage: types.DKGStageDealing}
	k.handleDKGProcessResponses(ctx, dkgNetwork, wrapResponses([]types.Response{selfResp, otherResp}), true)
}

// TestHandleDKGProcessResponses_Error_CachesResponses verifies that failed
// ProcessResponses calls cache the responses for retry.
func TestHandleDKGProcessResponses_Error_CachesResponses(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("fail-resp-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	session := &types.DKGSession{
		Round:          11,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
		Index:          0, // unset → all responses included
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Kernel returns error on all retries.
	mockKernel.EXPECT().ProcessResponses(gomock.Any(), gomock.Any()).Return(nil, errSentinel).Times(retryAttempts)

	dkgNetwork := &types.DKGNetwork{Round: 11, Stage: types.DKGStageDealing}

	resetPendingIncoming()
	defer resetPendingIncoming()

	resp := types.Response{Index: 0, VssResponse: &types.VSSResponse{Index: 1}}
	k.handleDKGProcessResponses(ctx, dkgNetwork, wrapResponses([]types.Response{resp}), true)

	pendingIncomingResponsesMu.Lock()
	cached := len(pendingIncomingResponses)
	pendingIncomingResponsesMu.Unlock()
	require.Equal(t, 1, cached, "failed response should be cached for retry")
}

// --- handleDKGProcessJustifications: upgrade CC filtering and error path (gap 16) ---

// TestHandleDKGProcessJustifications_UpgradeCCFiltering verifies that for upgrade
// resharing, handleDKGProcessJustifications sends justifications to BOTH binaries.
func TestHandleDKGProcessJustifications_UpgradeCCFiltering(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	oldCC := []byte("old-just-cc")
	newCC := []byte("new-just-cc")
	mockOldKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockNewKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(oldCC, mockOldKernel)
	router.RegisterClient(newCC, mockNewKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	session := &types.DKGSession{
		Round:             12,
		Phase:             types.PhaseDealing,
		CodeCommitment:    newCC,
		OldCodeCommitment: oldCC,
		IsUpgrade:         true,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	just := types.Justification{Index: 0}

	// Both old and new kernel should receive ProcessJustification.
	mockOldKernel.EXPECT().ProcessJustification(gomock.Any(), gomock.Any()).Return(
		&types.ProcessJustificationResponse{}, nil,
	)
	mockNewKernel.EXPECT().ProcessJustification(gomock.Any(), gomock.Any()).Return(
		&types.ProcessJustificationResponse{}, nil,
	)

	dkgNetwork := &types.DKGNetwork{Round: 12, Stage: types.DKGStageDealing, IsUpgrade: true}
	k.handleDKGProcessJustifications(ctx, dkgNetwork, wrapJustifications([]types.Justification{just}))
}

// TestHandleDKGProcessJustifications_Error_CachesJustifications verifies that
// failed ProcessJustification kernel calls cache the justifications for retry.
func TestHandleDKGProcessJustifications_Error_CachesJustifications(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("fail-just-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	session := &types.DKGSession{
		Round:          13,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Kernel returns error on all retries.
	mockKernel.EXPECT().ProcessJustification(gomock.Any(), gomock.Any()).Return(nil, errSentinel).Times(retryAttempts)

	dkgNetwork := &types.DKGNetwork{Round: 13, Stage: types.DKGStageDealing}

	resetPendingIncoming()
	defer resetPendingIncoming()

	just := types.Justification{Index: 0}
	k.handleDKGProcessJustifications(ctx, dkgNetwork, wrapJustifications([]types.Justification{just}))

	pendingIncomingJustificationsMu.Lock()
	cached := len(pendingIncomingJustifications)
	pendingIncomingJustificationsMu.Unlock()
	require.Equal(t, 1, cached, "failed justification should be cached for retry")
}

// TestHandleDKGProcessJustifications_SuccessPath verifies the success path of
// handleDKGProcessJustifications forwards justifications to the kernel.
func TestHandleDKGProcessJustifications_SuccessPath(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("success-just-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	session := &types.DKGSession{
		Round:          14,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Kernel successfully processes justification.
	mockKernel.EXPECT().ProcessJustification(gomock.Any(), gomock.Any()).Return(
		&types.ProcessJustificationResponse{}, nil,
	)

	dkgNetwork := &types.DKGNetwork{Round: 14, Stage: types.DKGStageDealing}

	resetPendingIncoming()
	defer resetPendingIncoming()

	just := types.Justification{Index: 0}
	k.handleDKGProcessJustifications(ctx, dkgNetwork, wrapJustifications([]types.Justification{just}))

	pendingIncomingJustificationsMu.Lock()
	cached := len(pendingIncomingJustifications)
	pendingIncomingJustificationsMu.Unlock()
	require.Equal(t, 0, cached, "successful processing should not cache justifications")
}

// --- reprocessPendingIncomingData: async goroutine (gap 17) ---

// TestReprocessPendingIncomingData_SuccessPath verifies reprocessPendingIncomingData
// drains and replays cached deals/responses/justifications when kernel is connected.
func TestReprocessPendingIncomingData_SuccessPath(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("reprocess-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		validatorEVMAddr: testValidatorAddr,
	}

	// Create a session for the reprocess calls.
	session := &types.DKGSession{
		Round:          15,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
		Index:          1, // accept RecipientIndex=0
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Pre-populate pending queues.
	resetPendingIncoming()
	defer resetPendingIncoming()

	pendingIncomingDealsMu.Lock()
	pendingIncomingDeals = []pendingDeal{{deal: types.Deal{Index: 0, RecipientIndex: 0}}}
	pendingIncomingDealsMu.Unlock()

	pendingIncomingResponsesMu.Lock()
	pendingIncomingResponses = []pendingResponse{{response: types.Response{Index: 0, VssResponse: &types.VSSResponse{Index: 1}}}}
	pendingIncomingResponsesMu.Unlock()

	pendingIncomingJustificationsMu.Lock()
	pendingIncomingJustifications = []pendingJustification{{justification: types.Justification{Index: 0}}}
	pendingIncomingJustificationsMu.Unlock()

	dkgNetwork := &types.DKGNetwork{
		Round:        15,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{testValidatorAddr},
	}

	// Kernel should receive ProcessDeals, ProcessResponses, ProcessJustification.
	mockKernel.EXPECT().ProcessDeals(gomock.Any(), gomock.Any()).Return(
		&types.ProcessDealsResponse{Responses: []types.Response{{Index: 0}}}, nil,
	)
	mockKernel.EXPECT().ProcessResponses(gomock.Any(), gomock.Any()).Return(
		&types.ProcessResponsesResponse{}, nil,
	)
	mockKernel.EXPECT().ProcessJustification(gomock.Any(), gomock.Any()).Return(
		&types.ProcessJustificationResponse{}, nil,
	)

	k.reprocessPendingIncomingData(dkgNetwork)

	// Give the goroutine time to complete.
	// We wait briefly since the goroutine is async.
	time.Sleep(200 * time.Millisecond)

	// After reprocessing, queues should be drained.
	pendingIncomingDealsMu.Lock()
	dealsLeft := len(pendingIncomingDeals)
	pendingIncomingDealsMu.Unlock()
	require.Equal(t, 0, dealsLeft, "pending deals should be drained after reprocessing")
}

// --- Full path tests merged from dkg_svc_full_path_test.go ---

// TestHandleDKGDealing_FullPath tests the full dealing path with mocked kernel client.

func TestHandleDKGDealing_FullPath(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("deal-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager: sm,
		kernelRouter: router,
	}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	session := &types.DKGSession{
		Round:          4,
		Phase:          types.PhaseInitialized,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Kernel GenerateDeals returns deals
	testDeals := []types.Deal{
		{Index: 0, RecipientIndex: 1, Deal: types.EncryptedDeal{DhKey: []byte("deal1")}},
	}
	mockKernel.EXPECT().GenerateDeals(gomock.Any(), gomock.Any()).Return(
		&types.GenerateDealsResponse{Deals: testDeals}, nil,
	)

	dkgNetwork := &types.DKGNetwork{
		Round: 4,
		Stage: types.DKGStageDealing,
	}

	// Flush queues first
	k.FlushAllQueues()
	k.handleDKGDealing(ctx, dkgNetwork, true)

	got, err := sm.GetSession(4)
	require.NoError(t, err)
	require.Equal(t, types.PhaseDealing, got.Phase)

	// Verify deals were enqueued
	dequeued := k.DequeueDeals(10)
	require.Len(t, dequeued, 1)
	require.Equal(t, uint32(0), dequeued[0].Index)
}

// --- retry limit tests: handleDKGProcessDeals increments retryCount and drops exceeded ---

func TestHandleDKGProcessDeals_RetryLimit_IncrementsAndDrops(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("retry-limit-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		validatorEVMAddr: testValidatorAddr,
	}

	session := &types.DKGSession{
		Round:          20,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
		Index:          1, // accept RecipientIndex=0
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Kernel returns error on all retries.
	mockKernel.EXPECT().ProcessDeals(gomock.Any(), gomock.Any()).Return(nil, errSentinel).Times(retryAttempts)

	dkgNetwork := &types.DKGNetwork{
		Round:        20,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{testValidatorAddr},
	}

	resetPendingIncoming()
	defer resetPendingIncoming()

	// Item with retryCount at max — should be dropped after increment.
	// Item with retryCount=0 — should survive with retryCount=1.
	pending := []pendingDeal{
		{deal: types.Deal{Index: 0, RecipientIndex: 0}, retryCount: maxReprocessAttempts},
		{deal: types.Deal{Index: 1, RecipientIndex: 0}, retryCount: 0},
	}

	k.handleDKGProcessDeals(ctx, dkgNetwork, pending)

	pendingIncomingDealsMu.Lock()
	cached := pendingIncomingDeals
	pendingIncomingDealsMu.Unlock()

	require.Len(t, cached, 1, "only the item below max should be cached")
	require.Equal(t, uint32(1), cached[0].deal.Index)
	require.Equal(t, 1, cached[0].retryCount, "retryCount should be incremented to 1")
}

func TestHandleDKGProcessResponses_RetryLimit_IncrementsAndDrops(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("retry-limit-resp-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	session := &types.DKGSession{
		Round:          21,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
		Index:          0, // unset: include all responses
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Kernel returns error on all retries.
	mockKernel.EXPECT().ProcessResponses(gomock.Any(), gomock.Any()).Return(nil, errSentinel).Times(retryAttempts)

	dkgNetwork := &types.DKGNetwork{Round: 21, Stage: types.DKGStageDealing}

	resetPendingIncoming()
	defer resetPendingIncoming()

	pending := []pendingResponse{
		{response: types.Response{Index: 0, VssResponse: &types.VSSResponse{Index: 1}}, retryCount: maxReprocessAttempts},
		{response: types.Response{Index: 1, VssResponse: &types.VSSResponse{Index: 2}}, retryCount: 0},
	}

	k.handleDKGProcessResponses(ctx, dkgNetwork, pending, true)

	pendingIncomingResponsesMu.Lock()
	cached := pendingIncomingResponses
	pendingIncomingResponsesMu.Unlock()

	require.Len(t, cached, 1, "only the item below max should be cached")
	require.Equal(t, uint32(1), cached[0].response.Index)
	require.Equal(t, 1, cached[0].retryCount, "retryCount should be incremented to 1")
}

func TestHandleDKGProcessJustifications_RetryLimit_IncrementsAndDrops(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("retry-limit-just-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	session := &types.DKGSession{
		Round:          22,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Kernel returns error on all retries.
	mockKernel.EXPECT().ProcessJustification(gomock.Any(), gomock.Any()).Return(nil, errSentinel).Times(retryAttempts)

	dkgNetwork := &types.DKGNetwork{Round: 22, Stage: types.DKGStageDealing}

	resetPendingIncoming()
	defer resetPendingIncoming()

	pending := []pendingJustification{
		{justification: types.Justification{Index: 0}, retryCount: maxReprocessAttempts},
		{justification: types.Justification{Index: 1}, retryCount: 1},
	}

	k.handleDKGProcessJustifications(ctx, dkgNetwork, pending)

	pendingIncomingJustificationsMu.Lock()
	cached := pendingIncomingJustifications
	pendingIncomingJustificationsMu.Unlock()

	require.Len(t, cached, 1, "only the item below max should be cached")
	require.Equal(t, uint32(1), cached[0].justification.Index)
	require.Equal(t, 2, cached[0].retryCount, "retryCount should be incremented to 2")
}

// --- wrap helper tests ---

func TestWrapDeals(t *testing.T) {
	t.Parallel()

	deals := []types.Deal{
		{Index: 0, RecipientIndex: 1},
		{Index: 1, RecipientIndex: 2},
	}

	wrapped := wrapDeals(deals)
	require.Len(t, wrapped, 2)
	require.Equal(t, uint32(0), wrapped[0].deal.Index)
	require.Equal(t, uint32(1), wrapped[1].deal.Index)
	require.Equal(t, 0, wrapped[0].retryCount, "new items should have retryCount=0")
	require.Equal(t, 0, wrapped[1].retryCount, "new items should have retryCount=0")
}

func TestWrapResponses(t *testing.T) {
	t.Parallel()

	responses := []types.Response{{Index: 5}, {Index: 6}}
	wrapped := wrapResponses(responses)
	require.Len(t, wrapped, 2)
	require.Equal(t, uint32(5), wrapped[0].response.Index)
	require.Equal(t, 0, wrapped[0].retryCount)
}

func TestWrapJustifications(t *testing.T) {
	t.Parallel()

	justifications := []types.Justification{{Index: 7}}
	wrapped := wrapJustifications(justifications)
	require.Len(t, wrapped, 1)
	require.Equal(t, uint32(7), wrapped[0].justification.Index)
	require.Equal(t, 0, wrapped[0].retryCount)
}

// --- per-item retry tests (kernel reports rejected items via rejected_* fields) ---

// TestRejectedKey_DealUsesOuterOnly verifies that for deals the inner part of
// rejectedKey stays zero (Deal.Index is unique within a filtered batch
// because the CL pre-filters by RecipientIndex).
func TestRejectedKey_DealUsesOuterOnly(t *testing.T) {
	t.Parallel()
	require.Equal(t, rejectedKey{outer: 5}, rejectedDealKey(types.Deal{Index: 5, RecipientIndex: 2}))
}

// TestRejectedKey_ResponseDisambiguatesByInner verifies that responses with
// the same dealer index but different complainer indices map to distinct keys.
// This is the regression test for the multi-complainer-per-dealer case.
func TestRejectedKey_ResponseDisambiguatesByInner(t *testing.T) {
	t.Parallel()
	a := rejectedResponseKey(types.Response{Index: 7, VssResponse: &types.VSSResponse{Index: 1}})
	b := rejectedResponseKey(types.Response{Index: 7, VssResponse: &types.VSSResponse{Index: 2}})
	require.NotEqual(t, a, b, "same dealer index + different complainer indices must produce distinct keys")
}

// TestHandleDKGProcessDeals_PerItemRetry_RequeuesOnlyFailed verifies that when
// the kernel returns success with failed_items populated, only those items are
// requeued (not the whole batch).
func TestHandleDKGProcessDeals_PerItemRetry_RequeuesOnlyFailed(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("per-item-deals-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router, validatorEVMAddr: testValidatorAddr}
	defer k.FlushAllQueues()

	session := &types.DKGSession{
		Round:          30,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
		Index:          1, // accept RecipientIndex=0
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Kernel succeeds overall, reports item Index=1 as rejected.
	mockKernel.EXPECT().ProcessDeals(gomock.Any(), gomock.Any()).Return(
		&types.ProcessDealsResponse{
			CodeCommitment: cc,
			Round:          30,
			Responses:      []types.Response{{Index: 0}},
			RejectedDeals:  []types.Deal{{Index: 1, RecipientIndex: 0}},
		}, nil,
	).Times(1)

	dkgNetwork := &types.DKGNetwork{
		Round:        30,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{testValidatorAddr},
	}

	resetPendingIncoming()
	defer resetPendingIncoming()

	pending := []pendingDeal{
		{deal: types.Deal{Index: 0, RecipientIndex: 0}, retryCount: 0},
		{deal: types.Deal{Index: 1, RecipientIndex: 0}, retryCount: 0},
	}

	k.handleDKGProcessDeals(ctx, dkgNetwork, pending)

	pendingIncomingDealsMu.Lock()
	cached := pendingIncomingDeals
	pendingIncomingDealsMu.Unlock()

	require.Len(t, cached, 1, "only the failed item should be requeued")
	require.Equal(t, uint32(1), cached[0].deal.Index)
	require.Equal(t, 1, cached[0].retryCount)
}

// TestHandleDKGProcessDeals_NewKernel_AllSuccess verifies that when the new kernel
// returns success with empty failed_items, nothing is requeued.
func TestHandleDKGProcessDeals_NewKernel_AllSuccess_NoRequeue(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("new-kernel-all-success-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router, validatorEVMAddr: testValidatorAddr}
	defer k.FlushAllQueues()

	session := &types.DKGSession{
		Round:          31,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
		Index:          1,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	mockKernel.EXPECT().ProcessDeals(gomock.Any(), gomock.Any()).Return(
		&types.ProcessDealsResponse{
			CodeCommitment: cc,
			Round:          31,
			Responses:      []types.Response{{Index: 0}, {Index: 1}},
			RejectedDeals:  nil,
		}, nil,
	).Times(1)

	dkgNetwork := &types.DKGNetwork{
		Round:        31,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{testValidatorAddr},
	}

	resetPendingIncoming()
	defer resetPendingIncoming()

	pending := []pendingDeal{
		{deal: types.Deal{Index: 0, RecipientIndex: 0}, retryCount: 0},
		{deal: types.Deal{Index: 1, RecipientIndex: 0}, retryCount: 0},
	}

	k.handleDKGProcessDeals(ctx, dkgNetwork, pending)

	pendingIncomingDealsMu.Lock()
	defer pendingIncomingDealsMu.Unlock()
	require.Empty(t, pendingIncomingDeals, "no items should be requeued when failed_items is empty")
}

// TestHandleDKGProcessDeals_OldKernel_Fallback verifies that an old kernel
// (without rejected_deals support) leaves the new client in legacy batch
// mode: RPC success with non-empty responses and empty rejected_deals →
// no requeue.
func TestHandleDKGProcessDeals_OldKernel_Fallback(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("old-kernel-fallback-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router, validatorEVMAddr: testValidatorAddr}
	defer k.FlushAllQueues()

	session := &types.DKGSession{
		Round:          32,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
		Index:          1,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Old kernel: responds with responses only, no rejected_deals field set.
	mockKernel.EXPECT().ProcessDeals(gomock.Any(), gomock.Any()).Return(
		&types.ProcessDealsResponse{
			CodeCommitment: cc,
			Round:          32,
			Responses:      []types.Response{{Index: 0}},
			// RejectedDeals intentionally omitted — simulates pre-feature kernel.
		}, nil,
	).Times(1)

	dkgNetwork := &types.DKGNetwork{
		Round:        32,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{testValidatorAddr},
	}

	resetPendingIncoming()
	defer resetPendingIncoming()

	pending := []pendingDeal{
		{deal: types.Deal{Index: 0, RecipientIndex: 0}, retryCount: 0},
		{deal: types.Deal{Index: 1, RecipientIndex: 0}, retryCount: 0},
	}

	k.handleDKGProcessDeals(ctx, dkgNetwork, pending)

	pendingIncomingDealsMu.Lock()
	defer pendingIncomingDealsMu.Unlock()
	require.Empty(t, pendingIncomingDeals,
		"old-kernel fallback: success RPC with empty rejected_deals must not requeue anything")
}

// TestHandleDKGProcessResponses_PerItemRetry_RequeuesOnlyFailed is the responses
// counterpart of the deals per-item retry test.
func TestHandleDKGProcessResponses_PerItemRetry_RequeuesOnlyFailed(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("per-item-resp-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	session := &types.DKGSession{
		Round:          40,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
		Index:          0, // unset: include all responses
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	mockKernel.EXPECT().ProcessResponses(gomock.Any(), gomock.Any()).Return(
		&types.ProcessResponsesResponse{
			Justifications: nil,
			// Composite key (outer=1, inner=2) matches the second pending response.
			RejectedResponses: []types.Response{{Index: 1, VssResponse: &types.VSSResponse{Index: 2}}},
		}, nil,
	).Times(1)

	dkgNetwork := &types.DKGNetwork{Round: 40, Stage: types.DKGStageDealing}

	resetPendingIncoming()
	defer resetPendingIncoming()

	pending := []pendingResponse{
		{response: types.Response{Index: 0, VssResponse: &types.VSSResponse{Index: 1}}, retryCount: 0},
		{response: types.Response{Index: 1, VssResponse: &types.VSSResponse{Index: 2}}, retryCount: 0},
	}

	k.handleDKGProcessResponses(ctx, dkgNetwork, pending, true)

	pendingIncomingResponsesMu.Lock()
	cached := pendingIncomingResponses
	pendingIncomingResponsesMu.Unlock()

	require.Len(t, cached, 1)
	require.Equal(t, uint32(1), cached[0].response.Index)
	require.Equal(t, 1, cached[0].retryCount)
}

// TestHandleDKGProcessJustifications_PerItemRetry_RequeuesOnlyFailed is the
// justifications counterpart of the deals per-item retry test.
func TestHandleDKGProcessJustifications_PerItemRetry_RequeuesOnlyFailed(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("per-item-just-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	session := &types.DKGSession{
		Round:          50,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	mockKernel.EXPECT().ProcessJustification(gomock.Any(), gomock.Any()).Return(
		&types.ProcessJustificationResponse{
			RejectedJustifications: []types.Justification{{Index: 0, VssJustification: &types.VSSJustification{Index: 3}}},
		}, nil,
	).Times(1)

	dkgNetwork := &types.DKGNetwork{Round: 50, Stage: types.DKGStageDealing}

	resetPendingIncoming()
	defer resetPendingIncoming()

	pending := []pendingJustification{
		{justification: types.Justification{Index: 0, VssJustification: &types.VSSJustification{Index: 3}}, retryCount: 0},
		{justification: types.Justification{Index: 1, VssJustification: &types.VSSJustification{Index: 4}}, retryCount: 0},
	}

	k.handleDKGProcessJustifications(ctx, dkgNetwork, pending)

	pendingIncomingJustificationsMu.Lock()
	cached := pendingIncomingJustifications
	pendingIncomingJustificationsMu.Unlock()

	require.Len(t, cached, 1)
	require.Equal(t, uint32(0), cached[0].justification.Index)
	require.Equal(t, 1, cached[0].retryCount)
}

// --- RETRY-1 regression: multiple responses targeting the same dealer ---

// TestHandleDKGProcessResponses_MultiComplainerPerDealer_OnlyFailedInnerRequeued
// verifies that when multiple complainers produce Responses with the SAME
// Response.Index (dealer), a kernel rejected_response targeting one specific
// (dealer, complainer) pair via VssResponse.Index requeues ONLY that pair.
// Previous design collapsed duplicate outer Index into a single map entry and
// requeued every response with that dealer index — the typed rejected slice
// fixes this because the entire Response is included in the wire.
func TestHandleDKGProcessResponses_MultiComplainerPerDealer_OnlyFailedInnerRequeued(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("multi-complainer-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}
	defer k.FlushAllQueues()

	session := &types.DKGSession{
		Round:          60,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Kernel reports only (dealer=7, complainer=1) as rejected.
	// Previously this would have requeued all three items (all have Response.Index=7).
	mockKernel.EXPECT().ProcessResponses(gomock.Any(), gomock.Any()).Return(
		&types.ProcessResponsesResponse{
			RejectedResponses: []types.Response{
				{Index: 7, VssResponse: &types.VSSResponse{Index: 1}},
			},
		}, nil,
	).Times(1)

	dkgNetwork := &types.DKGNetwork{Round: 60, Stage: types.DKGStageDealing}

	resetPendingIncoming()
	defer resetPendingIncoming()

	// Three responses all targeting dealer 7 from different complainers (0, 1, 2).
	pending := []pendingResponse{
		{response: types.Response{Index: 7, VssResponse: &types.VSSResponse{Index: 0}}},
		{response: types.Response{Index: 7, VssResponse: &types.VSSResponse{Index: 1}}},
		{response: types.Response{Index: 7, VssResponse: &types.VSSResponse{Index: 2}}},
	}

	k.handleDKGProcessResponses(ctx, dkgNetwork, pending, true)

	pendingIncomingResponsesMu.Lock()
	cached := pendingIncomingResponses
	pendingIncomingResponsesMu.Unlock()

	require.Len(t, cached, 1, "only the (7, 1) pair should be requeued")
	require.Equal(t, uint32(7), cached[0].response.Index)
	require.Equal(t, uint32(1), cached[0].response.VssResponse.Index)
	require.Equal(t, 1, cached[0].retryCount)
}

// --- RETRY-2 regression: upgrade resharing must not double-increment retryCount ---

// TestHandleDKGProcessResponses_UpgradeResharing_BothCCsFail_SingleIncrement
// verifies that during upgrade resharing, when BOTH code commitments report the
// same item as failed, the retryCount is incremented exactly once (not twice
// per ccsToProcess iteration). Previous implementation shared the filtered
// slice across iterations and mutated retryCount in-place, effectively halving
// the retry budget.
func TestHandleDKGProcessResponses_UpgradeResharing_BothCCsFail_SingleIncrement(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	newCC := []byte("upgrade-new-cc")
	oldCC := []byte("upgrade-old-cc")
	mockKernelNew := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockKernelOld := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(newCC, mockKernelNew)
	router.RegisterClient(oldCC, mockKernelOld)

	k := &Keeper{stateManager: sm, kernelRouter: router}
	defer k.FlushAllQueues()

	session := &types.DKGSession{
		Round:             70,
		Phase:             types.PhaseDealing,
		CodeCommitment:    newCC,
		OldCodeCommitment: oldCC,
		IsResharing:       true,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Both kernels report the SAME item as rejected.
	rejectedResp := []types.Response{{Index: 3, VssResponse: &types.VSSResponse{Index: 5}}}
	mockKernelNew.EXPECT().ProcessResponses(gomock.Any(), gomock.Any()).Return(
		&types.ProcessResponsesResponse{RejectedResponses: rejectedResp}, nil,
	).Times(1)
	mockKernelOld.EXPECT().ProcessResponses(gomock.Any(), gomock.Any()).Return(
		&types.ProcessResponsesResponse{RejectedResponses: rejectedResp}, nil,
	).Times(1)

	dkgNetwork := &types.DKGNetwork{Round: 70, Stage: types.DKGStageDealing, IsUpgrade: true}

	resetPendingIncoming()
	defer resetPendingIncoming()

	pending := []pendingResponse{
		{response: types.Response{Index: 3, VssResponse: &types.VSSResponse{Index: 5}}, retryCount: 0},
	}

	k.handleDKGProcessResponses(ctx, dkgNetwork, pending, true)

	pendingIncomingResponsesMu.Lock()
	cached := pendingIncomingResponses
	pendingIncomingResponsesMu.Unlock()

	require.Len(t, cached, 1, "item must be cached exactly once (not once per CC)")
	require.Equal(t, 1, cached[0].retryCount, "retryCount must increment exactly once (not twice)")
}

// TestHandleDKGProcessJustifications_UpgradeResharing_BothCCsFail_SingleIncrement
// is the justifications counterpart.
func TestHandleDKGProcessJustifications_UpgradeResharing_BothCCsFail_SingleIncrement(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	newCC := []byte("upgrade-just-new-cc")
	oldCC := []byte("upgrade-just-old-cc")
	mockKernelNew := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockKernelOld := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(newCC, mockKernelNew)
	router.RegisterClient(oldCC, mockKernelOld)

	k := &Keeper{stateManager: sm, kernelRouter: router}
	defer k.FlushAllQueues()

	session := &types.DKGSession{
		Round:             71,
		Phase:             types.PhaseDealing,
		CodeCommitment:    newCC,
		OldCodeCommitment: oldCC,
		IsResharing:       true,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	rejectedJust := []types.Justification{{Index: 2, VssJustification: &types.VSSJustification{Index: 4}}}
	mockKernelNew.EXPECT().ProcessJustification(gomock.Any(), gomock.Any()).Return(
		&types.ProcessJustificationResponse{RejectedJustifications: rejectedJust}, nil,
	).Times(1)
	mockKernelOld.EXPECT().ProcessJustification(gomock.Any(), gomock.Any()).Return(
		&types.ProcessJustificationResponse{RejectedJustifications: rejectedJust}, nil,
	).Times(1)

	dkgNetwork := &types.DKGNetwork{Round: 71, Stage: types.DKGStageDealing, IsUpgrade: true}

	resetPendingIncoming()
	defer resetPendingIncoming()

	pending := []pendingJustification{
		{justification: types.Justification{Index: 2, VssJustification: &types.VSSJustification{Index: 4}}, retryCount: 0},
	}

	k.handleDKGProcessJustifications(ctx, dkgNetwork, pending)

	pendingIncomingJustificationsMu.Lock()
	cached := pendingIncomingJustifications
	pendingIncomingJustificationsMu.Unlock()

	require.Len(t, cached, 1, "item must be cached exactly once (not once per CC)")
	require.Equal(t, 1, cached[0].retryCount, "retryCount must increment exactly once (not twice)")
}
