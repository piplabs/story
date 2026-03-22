package keeper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

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

// --- cachePendingIncomingDeals ---

func TestCachePendingIncomingDeals(t *testing.T) {
	tests := []struct {
		name         string
		deals        []types.Deal
		sessionIndex uint32
		expected     int // expected number of cached deals
	}{
		{
			name: "caches deals addressed to this validator",
			deals: []types.Deal{
				{Index: 0, RecipientIndex: 0}, // sessionIndex=1, recipientIndex=sessionIndex-1=0
				{Index: 1, RecipientIndex: 0},
				{Index: 2, RecipientIndex: 1}, // different recipient, not cached
			},
			sessionIndex: 1,
			expected:     2,
		},
		{
			name: "no deals match this validator",
			deals: []types.Deal{
				{Index: 0, RecipientIndex: 5},
				{Index: 1, RecipientIndex: 6},
			},
			sessionIndex: 1,
			expected:     0,
		},
		{
			name:         "empty deals list",
			deals:        nil,
			sessionIndex: 1,
			expected:     0,
		},
		{
			name: "session index 0 means unset, no deals cached",
			deals: []types.Deal{
				{Index: 0, RecipientIndex: 0},
			},
			sessionIndex: 0, // 0 means unset; guard (sessionIndex > 0) fails
			expected:     0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetPendingIncoming()
			defer resetPendingIncoming()

			cached := cachePendingIncomingDeals(tc.deals, tc.sessionIndex)
			require.Equal(t, tc.expected, cached)
		})
	}
}

func TestCachePendingIncomingDeals_MaxCapEnforced(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	// Pre-fill to max capacity
	pendingIncomingDealsMu.Lock()
	for i := range maxPendingIncoming {
		pendingIncomingDeals = append(pendingIncomingDeals, types.Deal{Index: uint32(i), RecipientIndex: 0})
	}
	pendingIncomingDealsMu.Unlock()

	// Try to add one more
	deals := []types.Deal{{Index: 99, RecipientIndex: 0}}
	cached := cachePendingIncomingDeals(deals, 1)
	require.Equal(t, 0, cached, "should not cache beyond maxPendingIncoming")
}

// --- cachePendingIncomingResponses ---

func TestCachePendingIncomingResponses(t *testing.T) {
	tests := []struct {
		name      string
		responses []types.Response
		preload   int // number of pre-loaded responses
		expected  int // expected total after caching
	}{
		{
			name:      "cache responses when empty",
			responses: []types.Response{{Index: 1}, {Index: 2}},
			preload:   0,
			expected:  2,
		},
		{
			name:      "cache with partial capacity remaining",
			responses: []types.Response{{Index: 1}, {Index: 2}, {Index: 3}},
			preload:   maxPendingIncoming - 2,
			expected:  maxPendingIncoming,
		},
		{
			name:      "cache at full capacity (no-op)",
			responses: []types.Response{{Index: 1}},
			preload:   maxPendingIncoming,
			expected:  maxPendingIncoming,
		},
		{
			name:      "empty response list",
			responses: nil,
			preload:   0,
			expected:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetPendingIncoming()
			defer resetPendingIncoming()

			// Pre-load
			pendingIncomingResponsesMu.Lock()
			for i := range tc.preload {
				pendingIncomingResponses = append(pendingIncomingResponses, types.Response{Index: uint32(i + 100)})
			}
			pendingIncomingResponsesMu.Unlock()

			cachePendingIncomingResponses(tc.responses)

			pendingIncomingResponsesMu.Lock()
			got := len(pendingIncomingResponses)
			pendingIncomingResponsesMu.Unlock()
			require.Equal(t, tc.expected, got)
		})
	}
}

// --- cachePendingIncomingJustifications ---

func TestCachePendingIncomingJustifications(t *testing.T) {
	tests := []struct {
		name           string
		justifications []types.Justification
		preload        int
		expected       int
	}{
		{
			name:           "cache justifications when empty",
			justifications: []types.Justification{{Index: 1}, {Index: 2}},
			preload:        0,
			expected:       2,
		},
		{
			name:           "cache with partial capacity",
			justifications: []types.Justification{{Index: 1}, {Index: 2}, {Index: 3}},
			preload:        maxPendingIncoming - 1,
			expected:       maxPendingIncoming,
		},
		{
			name:           "cache at full capacity (no-op)",
			justifications: []types.Justification{{Index: 1}},
			preload:        maxPendingIncoming,
			expected:       maxPendingIncoming,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetPendingIncoming()
			defer resetPendingIncoming()

			pendingIncomingJustificationsMu.Lock()
			for i := range tc.preload {
				pendingIncomingJustifications = append(pendingIncomingJustifications, types.Justification{Index: uint32(i + 100)})
			}
			pendingIncomingJustificationsMu.Unlock()

			cachePendingIncomingJustifications(tc.justifications)

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
	pendingIncomingDeals = []types.Deal{{Index: 1}, {Index: 2}}
	pendingIncomingDealsMu.Unlock()

	drained := drainPendingIncomingDeals()
	require.Len(t, drained, 2)
	require.Equal(t, uint32(1), drained[0].Index)
	require.Equal(t, uint32(2), drained[1].Index)

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
	pendingIncomingResponses = []types.Response{{Index: 3}, {Index: 4}}
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
	pendingIncomingJustifications = []types.Justification{{Index: 5}}
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
	pendingIncomingDeals = []types.Deal{{Index: 1}}
	pendingIncomingDealsMu.Unlock()

	pendingIncomingResponsesMu.Lock()
	pendingIncomingResponses = []types.Response{{Index: 2}}
	pendingIncomingResponsesMu.Unlock()

	pendingIncomingJustificationsMu.Lock()
	pendingIncomingJustifications = []types.Justification{{Index: 3}}
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
	k.handleDKGProcessDeals(ctx, dkgNetwork, deals)
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
	k.handleDKGProcessDeals(ctx, dkgNetwork, deals)
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
	k.handleDKGProcessDeals(ctx, dkgNetwork, deals)

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
	k.handleDKGProcessDeals(ctx, dkgNetwork, deals)
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
	k.handleDKGProcessResponses(ctx, dkgNetwork, responses, false)
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
	k.handleDKGProcessResponses(ctx, dkgNetwork, responses, true)
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

	k.handleDKGProcessResponses(ctx, dkgNetwork, responses, true)

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
	k.handleDKGProcessResponses(ctx, dkgNetwork, responses, true)
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
	k.handleDKGProcessJustifications(ctx, dkgNetwork, justifications)
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

	k.handleDKGProcessJustifications(ctx, dkgNetwork, justifications)

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
	pendingIncomingDeals = []types.Deal{{Index: 1}}
	pendingIncomingDealsMu.Unlock()

	pendingIncomingResponsesMu.Lock()
	pendingIncomingResponses = []types.Response{{Index: 2}}
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
	pendingIncomingDeals = []types.Deal{{Index: 1}}
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
	pendingIncomingDeals = []types.Deal{{Index: 1}}
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
