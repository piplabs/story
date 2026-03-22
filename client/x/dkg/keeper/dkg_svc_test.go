package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// resetDKGSvcRound resets the package-level dkgSvcRound atomic for test isolation.
func resetDKGSvcRound() {
	dkgSvcRound.Store(0)
}

// --- tryAcquireDKGSvc ---

func TestTryAcquireDKGSvc(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() // pre-test state
		round    uint32 // round to acquire
		expected bool   // expected return value
		postVal  uint64 // expected dkgSvcRound after call
	}{
		{
			name:     "acquire on idle (round=0)",
			setup:    func() { dkgSvcRound.Store(0) },
			round:    1,
			expected: true,
			postVal:  1,
		},
		{
			name:     "duplicate same round returns false",
			setup:    func() { dkgSvcRound.Store(5) },
			round:    5,
			expected: false,
			postVal:  5,
		},
		{
			name:     "newer round supersedes older round",
			setup:    func() { dkgSvcRound.Store(3) },
			round:    7,
			expected: true,
			postVal:  7,
		},
		{
			name:     "older round cannot preempt newer round",
			setup:    func() { dkgSvcRound.Store(10) },
			round:    5,
			expected: false,
			postVal:  10,
		},
		{
			name:     "round 0 is never acquired (always idle check)",
			setup:    func() { dkgSvcRound.Store(0) },
			round:    0,
			expected: false,
			postVal:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			defer resetDKGSvcRound()

			got := tryAcquireDKGSvc(tc.round)
			require.Equal(t, tc.expected, got)
			require.Equal(t, tc.postVal, dkgSvcRound.Load())
		})
	}
}

// --- releaseDKGSvc ---

func TestReleaseDKGSvc(t *testing.T) {
	tests := []struct {
		name    string
		current uint64 // current dkgSvcRound value
		round   uint32 // round to release
		postVal uint64 // expected dkgSvcRound after release
	}{
		{
			name:    "release matching round resets to 0",
			current: 5,
			round:   5,
			postVal: 0,
		},
		{
			name:    "release non-matching round is a no-op",
			current: 7,
			round:   5,
			postVal: 7,
		},
		{
			name:    "release when already idle is a no-op",
			current: 0,
			round:   5,
			postVal: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dkgSvcRound.Store(tc.current)
			defer resetDKGSvcRound()

			releaseDKGSvc(tc.round)
			require.Equal(t, tc.postVal, dkgSvcRound.Load())
		})
	}
}

// --- isSessionStuckForStage ---

func TestIsSessionStuckForStage(t *testing.T) {
	tests := []struct {
		name     string
		phase    types.DKGPhase
		stage    types.DKGStage
		expected bool
	}{
		// Registration stage
		{
			name:     "registration: PhaseInitialized is not stuck",
			phase:    types.PhaseInitialized,
			stage:    types.DKGStageRegistration,
			expected: false,
		},
		{
			name:     "registration: PhaseCompleted is not stuck",
			phase:    types.PhaseCompleted,
			stage:    types.DKGStageRegistration,
			expected: false,
		},
		{
			name:     "registration: PhaseInitializing is stuck",
			phase:    types.PhaseInitializing,
			stage:    types.DKGStageRegistration,
			expected: true,
		},
		{
			name:     "registration: PhaseDealing is stuck",
			phase:    types.PhaseDealing,
			stage:    types.DKGStageRegistration,
			expected: true,
		},
		{
			name:     "registration: PhaseFailed is stuck",
			phase:    types.PhaseFailed,
			stage:    types.DKGStageRegistration,
			expected: true,
		},

		// Dealing stage
		{
			name:     "dealing: PhaseDealing is not stuck",
			phase:    types.PhaseDealing,
			stage:    types.DKGStageDealing,
			expected: false,
		},
		{
			name:     "dealing: PhaseCompleted is not stuck",
			phase:    types.PhaseCompleted,
			stage:    types.DKGStageDealing,
			expected: false,
		},
		{
			name:     "dealing: PhaseInitialized is stuck",
			phase:    types.PhaseInitialized,
			stage:    types.DKGStageDealing,
			expected: true,
		},
		{
			name:     "dealing: PhaseInitializing is stuck",
			phase:    types.PhaseInitializing,
			stage:    types.DKGStageDealing,
			expected: true,
		},

		// Finalization stage
		{
			name:     "finalization: PhaseFinalized is not stuck",
			phase:    types.PhaseFinalized,
			stage:    types.DKGStageFinalization,
			expected: false,
		},
		{
			name:     "finalization: PhaseCompleted is not stuck",
			phase:    types.PhaseCompleted,
			stage:    types.DKGStageFinalization,
			expected: false,
		},
		{
			name:     "finalization: PhaseDealing is stuck",
			phase:    types.PhaseDealing,
			stage:    types.DKGStageFinalization,
			expected: true,
		},

		// Active stage
		{
			name:     "active: PhaseCompleted is not stuck",
			phase:    types.PhaseCompleted,
			stage:    types.DKGStageActive,
			expected: false,
		},
		{
			name:     "active: PhaseFinalized is stuck",
			phase:    types.PhaseFinalized,
			stage:    types.DKGStageActive,
			expected: true,
		},

		// Unspecified stage
		{
			name:     "unspecified stage: any phase is not stuck",
			phase:    types.PhaseFailed,
			stage:    types.DKGStageUnspecified,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := isSessionStuckForStage(tc.phase, tc.stage)
			require.Equal(t, tc.expected, got)
		})
	}
}

// --- labelToUUID ---

func TestLabelToUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		label     []byte
		expected  uint32
		expectErr bool
	}{
		{
			name:      "valid 32-byte label with UUID in last 4 bytes",
			label:     make32ByteLabel(42),
			expected:  42,
			expectErr: false,
		},
		{
			name:      "valid 32-byte label with zero UUID",
			label:     make([]byte, 32),
			expected:  0,
			expectErr: false,
		},
		{
			name:      "valid 32-byte label with max uint32",
			label:     make32ByteLabel(0xFFFFFFFF),
			expected:  0xFFFFFFFF,
			expectErr: false,
		},
		{
			name:      "label longer than 32 bytes still works",
			label:     make32ByteLabel(100),
			expected:  100,
			expectErr: false,
		},
		{
			name:      "label shorter than 32 bytes returns error",
			label:     make([]byte, 31),
			expected:  0,
			expectErr: true,
		},
		{
			name:      "empty label returns error",
			label:     nil,
			expected:  0,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := labelToUUID(tc.label)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			}
		})
	}
}

// make32ByteLabel creates a 32-byte label with the given uint32 value at bytes [28:32].
func make32ByteLabel(val uint32) []byte {
	label := make([]byte, 32)
	label[28] = byte(val >> 24)
	label[29] = byte(val >> 16)
	label[30] = byte(val >> 8)
	label[31] = byte(val)

	return label
}

// --- dkgAsyncContext ---

func TestDkgAsyncContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := dkgAsyncContext()
	defer cancel()

	require.NotNil(t, ctx)

	deadline, ok := ctx.Deadline()
	require.True(t, ok, "context should have a deadline")
	require.WithinDuration(t, deadline, deadline, dkgAsyncTimeout)
}
