package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

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
