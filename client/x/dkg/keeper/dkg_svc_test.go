package keeper

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
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

// --- Tests merged from dkg_svc_decrypt_test.go ---

// --- processDecryptRequests / computePartialDecrypt / submitPartialDecryption ---

// TestProcessDecryptQueue_NilKernelRouter verifies that processDecryptQueue returns early
// without draining any session when kernelRouter is not configured.
func TestProcessDecryptQueue_NilKernelRouter(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	// kernelRouter guard fires before BlockNumber — no BlockNumber call expected.

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	req := types.DecryptRequest{Ciphertext: []byte("encrypted"), Label: make([]byte, 32)}
	session := &types.DKGSession{
		Round: 1, Index: 1, GlobalPubKey: []byte("pub"),
		DecryptRequests: []types.PendingDecryptRequest{{DecryptRequest: req}},
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: nil}
	k.processDecryptQueue(ctx)

	// Request must still be in the queue — kernelRouter guard fires before DrainDecryptRequests.
	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Len(t, got.GetDecryptRequests(), 1, "requests must remain queued when kernelRouter is nil")
}

// TestProcessDecryptQueue_ZeroIndex verifies that sessions with unset Index are skipped
// without draining their request queue.
func TestProcessDecryptQueue_ZeroIndex(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(uint64(100), nil)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	req := types.DecryptRequest{Ciphertext: []byte("encrypted"), Label: make([]byte, 32)}
	session := &types.DKGSession{
		Round: 1, Index: 0, GlobalPubKey: []byte("pub"), // Index not set
		DecryptRequests: []types.PendingDecryptRequest{{DecryptRequest: req}},
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: NewKernelRouter(nil, nil)}
	k.processDecryptQueue(ctx)

	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Len(t, got.GetDecryptRequests(), 1, "requests must remain queued when session index is unset")
}

// TestProcessDecryptQueue_MissingGlobalPubKey verifies that sessions without a global
// public key are skipped without draining their request queue.
func TestProcessDecryptQueue_MissingGlobalPubKey(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(uint64(100), nil)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	req := types.DecryptRequest{Ciphertext: []byte("encrypted"), Label: make([]byte, 32)}
	session := &types.DKGSession{
		Round: 1, Index: 1, GlobalPubKey: nil, // Missing
		DecryptRequests: []types.PendingDecryptRequest{{DecryptRequest: req}},
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: NewKernelRouter(nil, nil)}
	k.processDecryptQueue(ctx)

	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Len(t, got.GetDecryptRequests(), 1, "requests must remain queued when global pub key is missing")
}

func TestComputePartialDecrypt_NoKernelClient(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}
	session := &types.DKGSession{
		Round:          1,
		Index:          1,
		GlobalPubKey:   []byte("pubkey"),
		CodeCommitment: []byte("nonexistent-cc"),
	}
	req := types.DecryptRequest{
		Ciphertext: []byte("encrypted"),
		Label:      make([]byte, 32),
	}

	result := k.computePartialDecrypt(context.Background(), session, req)
	require.Error(t, result.err)
	require.Contains(t, result.err.Error(), "no kernel client for session")
}

// --- processDecryptQueue ---

func TestProcessDecryptQueue_NoSessions(t *testing.T) {
	t.Parallel()

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	k := &Keeper{stateManager: sm}

	// Should not panic with no sessions
	k.processDecryptQueue(context.Background())
}

func TestProcessDecryptQueue_SessionWithNoRequests(t *testing.T) {
	t.Parallel()

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	ctx := context.Background()

	session := &types.DKGSession{
		Round:           1,
		Phase:           types.PhaseCompleted,
		DecryptRequests: nil, // No pending requests
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{stateManager: sm}

	// Should be a no-op
	k.processDecryptQueue(ctx)
}

// TestProcessDecryptQueue_SessionReadyButEmpty verifies that a session with valid
// index and pubkey but no pending requests is skipped without error after the drain.
func TestProcessDecryptQueue_SessionReadyButEmpty(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(uint64(100), nil)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	// Session has index and pubkey set (passes pre-drain guards) but no requests.
	session := &types.DKGSession{
		Round: 1, Index: 1, GlobalPubKey: []byte("pub"),
		CodeCommitment:  []byte("cc"),
		DecryptRequests: nil,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: NewKernelRouter(nil, nil)}
	k.processDecryptQueue(ctx) // should be a no-op after drain
}

// TestProcessDecryptQueue_IdleSessionSkippedNoMutation verifies that a non-finalized
// session (Index==0, empty GlobalPubKey) with an EMPTY decrypt queue is skipped at the
// top of the loop before the precondition checks, so no state mutation occurs. This is
// the idle-session log-noise fix: the index/GPK deferral warnings must not fire for
// sessions that have no queued work.
func TestProcessDecryptQueue_IdleSessionSkippedNoMutation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(uint64(100), nil)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	// Non-finalized idle session: Index unset, GlobalPubKey empty, no requests.
	original := time.Now().Add(-time.Hour).UTC()
	session := &types.DKGSession{
		Round: 1, Index: 0, GlobalPubKey: nil,
		LastUpdate:      original,
		DecryptRequests: nil,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: NewKernelRouter(nil, nil)}
	k.processDecryptQueue(ctx)

	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Empty(t, got.GetDecryptRequests(), "idle session must stay empty")
	// LastUpdate is bumped by DrainDecryptRequests/UpdateSession; an unchanged value
	// proves the session was skipped before any draining or persistence occurred.
	require.True(t, got.LastUpdate.Equal(original), "idle session must not be mutated when skipped")
}

// TestProcessDecryptQueue_QueuedRequestReachesPreconditions verifies that a session with
// a non-empty queue still reaches the precondition checks (here: Index==0 deferral), i.e.
// the idle-session guard only skips sessions whose queue is empty.
func TestProcessDecryptQueue_QueuedRequestReachesPreconditions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(uint64(100), nil)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	req := types.DecryptRequest{Ciphertext: []byte("encrypted"), Label: make([]byte, 32)}
	session := &types.DKGSession{
		Round: 1, Index: 0, GlobalPubKey: nil, // not yet finalized
		DecryptRequests: []types.PendingDecryptRequest{{DecryptRequest: req}},
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: NewKernelRouter(nil, nil)}
	k.processDecryptQueue(ctx)

	// The guard did not skip this session; it hit the Index==0 deferral and the request
	// remains queued (deferred, not drained) for the next tick.
	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Len(t, got.GetDecryptRequests(), 1, "queued request must be deferred, not dropped")
}

func TestProcessDecryptQueue_FailedRequestsRetained(t *testing.T) {
	t.Parallel()

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	ctx := context.Background()

	router := NewKernelRouter(nil, nil)
	session := &types.DKGSession{
		Round:          1,
		Phase:          types.PhaseCompleted,
		Index:          1,
		GlobalPubKey:   []byte("pubkey"),
		CodeCommitment: []byte("nonexistent-cc"),
		DecryptRequests: []types.PendingDecryptRequest{
			{DecryptRequest: types.DecryptRequest{
				Ciphertext:      []byte("encrypted1"),
				Label:           make([]byte, 32),
				RequesterPubKey: []byte("pub1"),
			}},
		},
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{
		stateManager: sm,
		kernelRouter: router,
	}

	// All requests will fail (no kernel client for CC)
	k.processDecryptQueue(ctx)

	// Failed requests should be retained for retry
	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Len(t, got.GetDecryptRequests(), 1, "failed requests should be retained")
}

func TestProcessDecryptQueue_BlockNumberError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(uint64(0), errors.New("rpc error"))

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: NewKernelRouter(nil, nil)}
	k.processDecryptQueue(ctx) // should return early without panic
}

func TestProcessDecryptQueue_KernelUnavailable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(uint64(100), nil)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	// Session with requests but no kernel client registered for its code commitment.
	session := &types.DKGSession{
		Round: 1, Index: 1, GlobalPubKey: []byte("pub"),
		CodeCommitment: []byte("unknown-cc"),
		DecryptRequests: []types.PendingDecryptRequest{
			{DecryptRequest: types.DecryptRequest{Ciphertext: []byte("ct"), Label: make([]byte, 32), RequesterPubKey: []byte("rpk")}},
		},
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: NewKernelRouter(nil, nil)}
	k.processDecryptQueue(ctx)

	// Requests are dropped (not re-queued) when kernel binary is unavailable.
	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Empty(t, got.GetDecryptRequests(), "requests must be dropped for unavailable kernel")
}

func TestProcessDecryptQueue_StaleRequestsFiltered(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	cc := []byte("cc-stale")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	// Current height well beyond the stale timeout.
	currentHeight := uint64(types.DefaultDecryptTimeout + 1000)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(currentHeight, nil)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	// Request with Height=0 is far below currentHeight-DefaultDecryptTimeout → stale.
	session := &types.DKGSession{
		Round: 1, Index: 1, GlobalPubKey: []byte("pub"), CodeCommitment: cc,
		DecryptRequests: []types.PendingDecryptRequest{
			{DecryptRequest: types.DecryptRequest{Ciphertext: []byte("ct"), Label: make([]byte, 32), Height: 0}},
		},
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: router}
	k.processDecryptQueue(ctx)

	// Stale requests must be dropped — no kernel call expected (mockKernel has no EXPECT).
	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Empty(t, got.GetDecryptRequests(), "stale requests must be dropped")
}

func TestProcessDecryptQueue_PartialStaleRequests(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	cc := []byte("cc-partial-stale")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	currentHeight := uint64(types.DefaultDecryptTimeout + 1000)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(currentHeight, nil)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	// One stale request (Height=0) and one fresh request (Height=currentHeight-1).
	session := &types.DKGSession{
		Round: 1, Index: 1, GlobalPubKey: []byte("pub"), CodeCommitment: cc,
		DecryptRequests: []types.PendingDecryptRequest{
			{DecryptRequest: types.DecryptRequest{Ciphertext: []byte("stale"), Label: makeLabel(1), Height: 0}},
			{DecryptRequest: types.DecryptRequest{Ciphertext: []byte("fresh"), Label: makeLabel(2), Height: currentHeight - 1}},
		},
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).
		Return(&types.PartialDecryptTDH2Response{
			EncryptedPartialDecryption: []byte("p"), EphemeralPubKey: []byte("e"),
			PubShare: []byte("s"), Signature: []byte("sig"),
		}, nil).Times(1)

	mockContract.EXPECT().SubmitEncryptedPartialDecryptionBatch(
		gomock.Any(), gomock.Any(),
	).Return(&ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}, nil).Times(1)

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: router}
	k.processDecryptQueue(ctx)

	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Empty(t, got.GetDecryptRequests(), "only fresh request processed, nothing re-queued")
}

// ---------- stale session clearing (round <= latestActivated-2) ----------

// Drain requests from sessions >=2 rounds behind the latest activated round.
func TestProcessDecryptQueue_StaleSessionRequestsDropped(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(uint64(100), nil)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	req := types.DecryptRequest{Ciphertext: []byte("ct"), Label: make([]byte, 32)}

	// latestActivated=3 → drain round <=1.
	stale := &types.DKGSession{
		Round: 1, Index: 0, GlobalPubKey: []byte("pub"),
		DecryptRequests: []types.PendingDecryptRequest{{DecryptRequest: req}},
		IsFinalized:     true,
	}
	current := &types.DKGSession{
		Round: 3, Index: 0, GlobalPubKey: []byte("pub"),
		DecryptRequests: []types.PendingDecryptRequest{{DecryptRequest: req}},
		IsFinalized:     true,
	}
	require.NoError(t, sm.CreateSession(ctx, stale))
	require.NoError(t, sm.CreateSession(ctx, current))

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: NewKernelRouter(nil, nil)}
	k.processDecryptQueue(ctx)

	gotStale, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Empty(t, gotStale.GetDecryptRequests())

	gotCurrent, err := sm.GetSession(3)
	require.NoError(t, err)
	require.Len(t, gotCurrent.GetDecryptRequests(), 1)
}

// Boundary at latestActivated-1 is preserved (one-round buffer).
func TestProcessDecryptQueue_StaleSessionBoundaryPreserved(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(uint64(100), nil)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	req := types.DecryptRequest{Ciphertext: []byte("ct"), Label: make([]byte, 32)}

	// latestActivated=3 → drain round <=1; round 2 preserved.
	boundary := &types.DKGSession{
		Round: 2, Index: 0, GlobalPubKey: []byte("pub"),
		DecryptRequests: []types.PendingDecryptRequest{{DecryptRequest: req}},
		IsFinalized:     true,
	}
	latest := &types.DKGSession{
		Round: 3, Index: 0, GlobalPubKey: []byte("pub"),
		DecryptRequests: []types.PendingDecryptRequest{{DecryptRequest: req}},
		IsFinalized:     true,
	}
	require.NoError(t, sm.CreateSession(ctx, boundary))
	require.NoError(t, sm.CreateSession(ctx, latest))

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: NewKernelRouter(nil, nil)}
	k.processDecryptQueue(ctx)

	gotBoundary, err := sm.GetSession(2)
	require.NoError(t, err)
	require.Len(t, gotBoundary.GetDecryptRequests(), 1)
}

// latestActivated<2 guard prevents uint32 underflow.
func TestProcessDecryptQueue_SingleSessionNotCleared(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(uint64(100), nil)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	req := types.DecryptRequest{Ciphertext: []byte("ct"), Label: make([]byte, 32)}
	session := &types.DKGSession{
		Round: 1, Index: 0, GlobalPubKey: []byte("pub"),
		DecryptRequests: []types.PendingDecryptRequest{{DecryptRequest: req}},
		IsFinalized:     true,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: NewKernelRouter(nil, nil)}
	k.processDecryptQueue(ctx)

	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Len(t, got.GetDecryptRequests(), 1)
}

// Issue piplabs/story#826: stuck N+1 + created N+2 must not drain active N.
func TestProcessDecryptQueue_StuckRoundPreservesActive(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(uint64(100), nil)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	req := types.DecryptRequest{Ciphertext: []byte("ct"), Label: make([]byte, 32)}

	// Round 10: chain's latest_active. Index=0 → deferred, not cleared.
	active := &types.DKGSession{
		Round: 10, Index: 0, GlobalPubKey: []byte("pub"),
		DecryptRequests: []types.PendingDecryptRequest{{DecryptRequest: req}},
		IsFinalized:     true,
	}
	// Round 11: stuck mid-finalization.
	stuck := &types.DKGSession{
		Round: 11, Index: 0, GlobalPubKey: []byte("pub"),
		DecryptRequests: []types.PendingDecryptRequest{},
		IsFinalized:     false,
	}
	// Round 12: registration phase just started.
	pending := &types.DKGSession{
		Round: 12, Index: 0, GlobalPubKey: []byte("pub"),
		DecryptRequests: []types.PendingDecryptRequest{},
		IsFinalized:     false,
	}
	require.NoError(t, sm.CreateSession(ctx, active))
	require.NoError(t, sm.CreateSession(ctx, stuck))
	require.NoError(t, sm.CreateSession(ctx, pending))

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: NewKernelRouter(nil, nil)}
	k.processDecryptQueue(ctx)

	gotActive, err := sm.GetSession(10)
	require.NoError(t, err)
	require.Len(t, gotActive.GetDecryptRequests(), 1,
		"active round's queue must survive stuck/created successors")
}

// Fresh node with no activated session must not drain (memory > silent loss).
func TestProcessDecryptQueue_NoActivatedSessionsNoDrain(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	mockContract.EXPECT().BlockNumber(gomock.Any()).Return(uint64(100), nil)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	req := types.DecryptRequest{Ciphertext: []byte("ct"), Label: make([]byte, 32)}

	for _, r := range []uint32{1, 2, 3, 5, 8} {
		s := &types.DKGSession{
			Round: r, Index: 0, GlobalPubKey: []byte("pub"),
			DecryptRequests: []types.PendingDecryptRequest{{DecryptRequest: req}},
			IsFinalized:     false,
		}
		require.NoError(t, sm.CreateSession(ctx, s))
	}

	k := &Keeper{stateManager: sm, contractClient: mockContract, kernelRouter: NewKernelRouter(nil, nil)}
	k.processDecryptQueue(ctx)

	for _, r := range []uint32{1, 2, 3, 5, 8} {
		got, err := sm.GetSession(r)
		require.NoError(t, err)
		require.Len(t, got.GetDecryptRequests(), 1, "round %d", r)
	}
}

// TestComputePartialDecrypt_PanicRecovery verifies that a panic inside the kernel call
// is caught by the deferred recover and returned as an error (not a crash).
func TestComputePartialDecrypt_PanicRecovery(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	cc := []byte("cc-panic")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{kernelRouter: router}

	session := &types.DKGSession{
		Round: 1, Index: 1, GlobalPubKey: []byte("pub"), CodeCommitment: cc,
	}

	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ *types.PartialDecryptTDH2Request, _ ...interface{}) (*types.PartialDecryptTDH2Response, error) {
			panic("simulated kernel panic")
		}).Times(1)

	req := types.DecryptRequest{Ciphertext: []byte("ct"), Label: make([]byte, 32)}
	result := k.computePartialDecrypt(ctx, session, req)
	require.Error(t, result.err)
	require.Contains(t, result.err.Error(), "panic in computePartialDecrypt")
}

// TestSubmitPartialDecryption_ContractError verifies that a contract call failure is
// returned as an error so the caller can re-queue the request.
func TestSubmitPartialDecryption_ContractError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	k := &Keeper{contractClient: mockContract}
	session := &types.DKGSession{Round: 1, Index: 2}

	label := makeLabel(99)
	req := types.DecryptRequest{Ciphertext: []byte("ct"), Label: label, RequesterPubKey: []byte("rpk")}
	resp := &types.PartialDecryptTDH2Response{
		EncryptedPartialDecryption: []byte("p"), EphemeralPubKey: []byte("e"),
		PubShare: []byte("s"), Signature: []byte("sig"),
	}

	mockContract.EXPECT().SubmitEncryptedPartialDecryption(
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(nil, errors.New("contract rejected"))

	err := k.submitPartialDecryption(ctx, session, req, resp)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to submit partial decryption")
}

// --- Tests merged from dkg_svc_decrypt_full_test.go ---

// TestComputePartialDecrypt_FullPath exercises the kernel gRPC leg of the decrypt flow.
func TestComputePartialDecrypt_FullPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("decrypt-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{kernelRouter: router}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	label := make([]byte, 32)
	label[31] = 42

	req := types.DecryptRequest{
		Ciphertext:      []byte("ciphertext"),
		Label:           label,
		RequesterPubKey: []byte("requester-pub"),
	}

	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).Return(
		&types.PartialDecryptTDH2Response{
			EncryptedPartialDecryption: []byte("partial"),
			EphemeralPubKey:            []byte("eph-pub"),
			PubShare:                   []byte("pub-share"),
			Signature:                  []byte("sig"),
		}, nil,
	)

	result := k.computePartialDecrypt(ctx, session, req)
	require.NoError(t, result.err)
	require.NotNil(t, result.resp)
	require.Equal(t, []byte("partial"), result.resp.EncryptedPartialDecryption)
}

// TestSubmitPartialDecryption_FullPath exercises the contract submission leg.
func TestSubmitPartialDecryption_FullPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)
	k := &Keeper{contractClient: mockContract}

	session := &types.DKGSession{
		Round: 1,
		Index: 2,
	}

	label := make([]byte, 32)
	label[31] = 42 // UUID = 42

	req := types.DecryptRequest{
		Ciphertext:      []byte("ciphertext"),
		Label:           label,
		RequesterPubKey: []byte("requester-pub"),
	}

	resp := &types.PartialDecryptTDH2Response{
		EncryptedPartialDecryption: []byte("partial"),
		EphemeralPubKey:            []byte("eph-pub"),
		PubShare:                   []byte("pub-share"),
		Signature:                  []byte("sig"),
	}

	mockContract.EXPECT().SubmitEncryptedPartialDecryption(
		gomock.Any(),
		uint32(1), // round
		uint32(2), // pid
		[]byte("partial"),
		[]byte("eph-pub"),
		[]byte("pub-share"),
		[]byte("requester-pub"),
		[]byte("ciphertext"),
		uint32(42), // uuid from label
		[]byte("sig"),
	).Return(&ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}, nil)

	err := k.submitPartialDecryption(ctx, session, req, resp)
	require.NoError(t, err)
}

// TestSubmitPartialDecryption_InvalidLabel exercises the label-too-short error path.
func TestSubmitPartialDecryption_InvalidLabel(t *testing.T) {
	t.Parallel()

	k := &Keeper{}
	session := &types.DKGSession{Round: 1, Index: 2}

	req := types.DecryptRequest{
		Ciphertext:      []byte("ciphertext"),
		Label:           make([]byte, 20), // < 32 bytes
		RequesterPubKey: []byte("requester-pub"),
	}

	resp := &types.PartialDecryptTDH2Response{
		EncryptedPartialDecryption: []byte("partial"),
		EphemeralPubKey:            []byte("eph-pub"),
		PubShare:                   []byte("pub-share"),
		Signature:                  []byte("sig"),
	}

	err := k.submitPartialDecryption(context.Background(), session, req, resp)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid decrypt request label")
}

// --- processDecryptRequests concurrent scenarios ---

// makeLabel builds a 32-byte label with the given UUID encoded in the last 4 bytes.
func makeLabel(uuid uint32) []byte {
	label := make([]byte, 32)
	label[28] = byte(uuid >> 24)
	label[29] = byte(uuid >> 16)
	label[30] = byte(uuid >> 8)
	label[31] = byte(uuid)
	return label
}

// TestProcessDecryptRequests_AllSucceed verifies that N requests are computed in parallel
// and all submitted serially with no re-queuing.
func TestProcessDecryptRequests_AllSucceed(t *testing.T) {
	t.Parallel()

	const n = 5

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("cc-all-succeed")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{kernelRouter: router, contractClient: mockContract}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	requests := make([]types.PendingDecryptRequest, n)
	for i := range n {
		requests[i] = types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{
			Ciphertext:      []byte("ct"),
			Label:           makeLabel(uint32(i + 1)),
			RequesterPubKey: []byte("rpk"),
		}}
	}

	// Each kernel call may arrive in any order — use AnyTimes with Times(n).
	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).
		Return(&types.PartialDecryptTDH2Response{
			EncryptedPartialDecryption: []byte("partial"),
			EphemeralPubKey:            []byte("eph"),
			PubShare:                   []byte("share"),
			Signature:                  []byte("sig"),
		}, nil).
		Times(n)

	// All n results fit in one batch (n=5 < decryptBatchSize=20).
	mockContract.EXPECT().SubmitEncryptedPartialDecryptionBatch(
		gomock.Any(), gomock.Any(),
	).Return(&ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}, nil).
		Times(1)

	k.processDecryptRequests(ctx, session, requests)

	// No failures — nothing should be re-queued.
	require.Empty(t, session.GetDecryptRequests(), "no requests should be re-queued on full success")
}

// TestProcessDecryptRequests_PartialKernelFailure verifies that requests whose kernel
// call fails are re-queued while successfully computed requests are submitted.
func TestProcessDecryptRequests_PartialKernelFailure(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("cc-partial-fail")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{kernelRouter: router, contractClient: mockContract}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	// 3 requests: first two succeed, third fails at kernel.
	requests := []types.PendingDecryptRequest{
		{DecryptRequest: types.DecryptRequest{Ciphertext: []byte("ct0"), Label: makeLabel(1), RequesterPubKey: []byte("rpk")}},
		{DecryptRequest: types.DecryptRequest{Ciphertext: []byte("ct1"), Label: makeLabel(2), RequesterPubKey: []byte("rpk")}},
		{DecryptRequest: types.DecryptRequest{Ciphertext: []byte("ct2"), Label: makeLabel(3), RequesterPubKey: []byte("rpk")}},
	}

	successResp := &types.PartialDecryptTDH2Response{
		EncryptedPartialDecryption: []byte("partial"),
		EphemeralPubKey:            []byte("eph"),
		PubShare:                   []byte("share"),
		Signature:                  []byte("sig"),
	}

	// Use DoAndReturn to simulate one failure and two successes in any order.
	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, req *types.PartialDecryptTDH2Request, _ ...interface{}) (*types.PartialDecryptTDH2Response, error) {
			// Fail on the request carrying ciphertext "ct2".
			if string(req.Ciphertext) == "ct2" {
				return nil, errors.New("kernel error")
			}
			return successResp, nil
		}).
		Times(3)

	// 2 successful results go into one batch.
	mockContract.EXPECT().SubmitEncryptedPartialDecryptionBatch(
		gomock.Any(), gomock.Any(),
	).Return(&ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}, nil).
		Times(1)

	k.processDecryptRequests(ctx, session, requests)

	// The failed request must be re-queued for retry.
	requeued := session.GetDecryptRequests()
	require.Len(t, requeued, 1, "exactly one failed request should be re-queued")
	require.Equal(t, []byte("ct2"), requeued[0].Ciphertext)
}

// TestProcessDecryptRequests_DropsAfterMaxRetries: a failure past the retry cap drops
// the request instead of re-queuing it.
func TestProcessDecryptRequests_DropsAfterMaxRetries(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("cc-drop")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{kernelRouter: router, contractClient: mockContract}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	// Single request already at the retry cap; the next failure must drop it.
	requests := []types.PendingDecryptRequest{
		{
			DecryptRequest: types.DecryptRequest{Ciphertext: []byte("ct"), Label: makeLabel(1), RequesterPubKey: []byte("rpk")},
			RetryCount:     maxReprocessAttempts,
		},
	}

	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("kernel error")).
		Times(1)

	k.processDecryptRequests(ctx, session, requests)

	require.Empty(t, session.GetDecryptRequests(), "request past the retry cap must be dropped, not re-queued")
}

// TestProcessDecryptRequests_ConcurrentABCIWrite verifies that new requests added by the
// ABCI thread via AddDecryptRequest while processDecryptRequests is running are not lost.
// This is the key race scenario: DrainDecryptRequests atomically clears the queue before
// processing; any request added after the drain must survive in the queue.
func TestProcessDecryptRequests_ConcurrentABCIWrite(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("cc-abci-race")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{kernelRouter: router, contractClient: mockContract}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	// One pre-existing request that processDecryptRequests will process.
	existing := types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{
		Ciphertext:      []byte("existing-ct"),
		Label:           makeLabel(10),
		RequesterPubKey: []byte("rpk"),
	}}

	// While the kernel gRPC is in-flight (simulated by a brief sleep), the ABCI
	// thread concurrently adds a new request.
	abciAdded := types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{
		Ciphertext:      []byte("abci-ct"),
		Label:           makeLabel(11),
		RequesterPubKey: []byte("rpk"),
	}}

	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ *types.PartialDecryptTDH2Request, _ ...interface{}) (*types.PartialDecryptTDH2Response, error) {
			// Simulate kernel latency; ABCI thread writes concurrently during this window.
			go session.AddDecryptRequest(abciAdded)
			time.Sleep(5 * time.Millisecond)
			return &types.PartialDecryptTDH2Response{
				EncryptedPartialDecryption: []byte("partial"),
				EphemeralPubKey:            []byte("eph"),
				PubShare:                   []byte("share"),
				Signature:                  []byte("sig"),
			}, nil
		}).Times(1)

	mockContract.EXPECT().SubmitEncryptedPartialDecryptionBatch(
		gomock.Any(), gomock.Any(),
	).Return(&ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}, nil).Times(1)

	k.processDecryptRequests(ctx, session, []types.PendingDecryptRequest{existing})

	// Give the concurrent AddDecryptRequest goroutine time to complete.
	time.Sleep(20 * time.Millisecond)

	// The ABCI-added request must still be present — processDecryptRequests must not
	// overwrite it when it calls UpdateSession at the end.
	requeued := session.GetDecryptRequests()
	require.Len(t, requeued, 1, "ABCI-added request must not be lost")
	require.Equal(t, []byte("abci-ct"), requeued[0].Ciphertext)
}

// TestProcessDecryptRequests_MultiBatch verifies that when the number of requests
// exceeds decryptBatchSize, batchSubmitConsumer calls flushBatch multiple times:
// once per full batch and once for the remainder.
//
// With n=25 and decryptBatchSize=20 we expect exactly 2 contract calls:
//   - first flush:  20 items (full batch)
//   - second flush:  5 items (remainder)
func TestProcessDecryptRequests_MultiBatch(t *testing.T) {
	t.Parallel()

	const n = 25 // intentionally > decryptBatchSize(20)

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("cc-multi-batch")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{kernelRouter: router, contractClient: mockContract}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	requests := make([]types.DecryptRequest, n)
	for i := range n {
		requests[i] = types.DecryptRequest{
			Ciphertext:      []byte("ct"),
			Label:           makeLabel(uint32(i + 1)),
			RequesterPubKey: []byte("rpk"),
		}
	}

	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).
		Return(&types.PartialDecryptTDH2Response{
			EncryptedPartialDecryption: []byte("partial"),
			EphemeralPubKey:            []byte("eph"),
			PubShare:                   []byte("share"),
			Signature:                  []byte("sig"),
		}, nil).
		Times(n)

	// 25 results → flush at 20 + flush remainder 5 = 2 contract calls.
	mockContract.EXPECT().SubmitEncryptedPartialDecryptionBatch(
		gomock.Any(), gomock.Any(),
	).Return(&ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}, nil).
		Times(2)

	k.processDecryptRequests(ctx, session, wrapDecryptRequests(requests))

	require.Empty(t, session.GetDecryptRequests(), "no requests should be re-queued on full success")
}

// --- Parallel vs sequential throughput comparison ---

// TestProcessDecryptRequests_ParallelSpeedup measures the wall-clock improvement of the
// parallel kernel phase over a naive sequential implementation.
//
// Simulated latencies:
//   - kernelDelay: 100ms per request (proxy for real 6-8s TEE computation)
//   - contractDelay: 20ms per request (proxy for real ~2s tx mining)
//
// With N=5:
//
//	Sequential: 5 × (100ms + 20ms) = 600ms
//	Parallel:   100ms + 5 × 20ms   = 200ms   → ~3× speedup
//
// Extrapolated to production latencies (kernel=7s, contract=2s, N=5):
//
//	Sequential: 5 × 9s = 45s
//	Parallel:   7s + 5 × 2s = 17s  → ~62% reduction
func TestProcessDecryptRequests_ParallelSpeedup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping throughput comparison in short mode")
	}

	const (
		n             = 5
		kernelDelay   = 100 * time.Millisecond
		contractDelay = 20 * time.Millisecond
		// Minimum expected speedup ratio (parallel / sequential elapsed).
		// We use a conservative 2× to avoid flakiness on loaded CI machines.
		minSpeedup = 2.0
	)

	cc := []byte("cc-speedup")

	buildSession := func() *types.DKGSession {
		return &types.DKGSession{
			Round:          1,
			Index:          2,
			GlobalPubKey:   []byte("global-pub"),
			CodeCommitment: cc,
		}
	}

	buildRequests := func() []types.PendingDecryptRequest {
		reqs := make([]types.PendingDecryptRequest, n)
		for i := range n {
			reqs[i] = types.PendingDecryptRequest{DecryptRequest: types.DecryptRequest{
				Ciphertext:      []byte("ct"),
				Label:           makeLabel(uint32(i + 1)),
				RequesterPubKey: []byte("rpk"),
			}}
		}
		return reqs
	}

	kernelResp := &types.PartialDecryptTDH2Response{
		EncryptedPartialDecryption: []byte("partial"),
		EphemeralPubKey:            []byte("eph"),
		PubShare:                   []byte("share"),
		Signature:                  []byte("sig"),
	}

	// --- Sequential baseline ---
	seqCtrl := gomock.NewController(t)
	seqKernel := dkgtestutil.NewMockKernelServiceClient(seqCtrl)
	seqContract := dkgtestutil.NewMockDKGContractClient(seqCtrl)

	seqKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ *types.PartialDecryptTDH2Request, _ ...interface{}) (*types.PartialDecryptTDH2Response, error) {
			time.Sleep(kernelDelay)
			return kernelResp, nil
		}).Times(n)

	seqContract.EXPECT().SubmitEncryptedPartialDecryption(
		gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
	).DoAndReturn(func(_ context.Context, _, _ uint32, _, _, _, _, _ []byte, _ uint32, _ []byte) (*ethtypes.Receipt, error) {
		time.Sleep(contractDelay)
		return &ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}, nil
	}).Times(n)

	seqRouter := NewKernelRouter(nil, nil)
	seqRouter.RegisterClient(cc, seqKernel)
	seqKeeper := &Keeper{kernelRouter: seqRouter, contractClient: seqContract}
	seqSession := buildSession()

	seqStart := time.Now()
	// Sequential: process each request one by one (old approach).
	for _, preq := range buildRequests() {
		result := seqKeeper.computePartialDecrypt(context.Background(), seqSession, preq.DecryptRequest)
		if result.err != nil {
			t.Fatalf("sequential kernel call failed: %v", result.err)
		}
		if err := seqKeeper.submitPartialDecryption(context.Background(), seqSession, result.req, result.resp); err != nil {
			t.Fatalf("sequential submit failed: %v", err)
		}
	}
	seqElapsed := time.Since(seqStart)

	// --- Parallel implementation ---
	parCtrl := gomock.NewController(t)
	parKernel := dkgtestutil.NewMockKernelServiceClient(parCtrl)
	parContract := dkgtestutil.NewMockDKGContractClient(parCtrl)

	parKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ *types.PartialDecryptTDH2Request, _ ...interface{}) (*types.PartialDecryptTDH2Response, error) {
			time.Sleep(kernelDelay)
			return kernelResp, nil
		}).Times(n)

	// All n results fit in one batch (n=5 < decryptBatchSize=20).
	parContract.EXPECT().SubmitEncryptedPartialDecryptionBatch(
		gomock.Any(), gomock.Any(),
	).DoAndReturn(func(_ context.Context, _ interface{}) (*ethtypes.Receipt, error) {
		time.Sleep(contractDelay)
		return &ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}, nil
	}).Times(1)

	parRouter := NewKernelRouter(nil, nil)
	parRouter.RegisterClient(cc, parKernel)
	parKeeper := &Keeper{kernelRouter: parRouter, contractClient: parContract}
	parSession := buildSession()

	parStart := time.Now()
	parKeeper.processDecryptRequests(context.Background(), parSession, buildRequests())
	parElapsed := time.Since(parStart)

	speedup := float64(seqElapsed) / float64(parElapsed)

	t.Logf("n=%d  kernel_delay=%s  contract_delay=%s", n, kernelDelay, contractDelay)
	t.Logf("sequential: %s", seqElapsed.Round(time.Millisecond))
	t.Logf("parallel:   %s", parElapsed.Round(time.Millisecond))
	t.Logf("speedup:    %.2fx", speedup)
	t.Logf("")
	t.Logf("Extrapolated to production (kernel=7s, contract=2s, n=%d):", n)
	t.Logf("  sequential: %s", (time.Duration(n) * (7*time.Second + 2*time.Second)).Round(time.Second))
	t.Logf("  parallel:   %s", (7*time.Second + time.Duration(n)*2*time.Second).Round(time.Second))

	require.GreaterOrEqual(t, speedup, minSpeedup,
		"parallel should be at least %.1fx faster than sequential", minSpeedup)
	require.Empty(t, parSession.GetDecryptRequests(), "no requests should be re-queued")
}

// --- Decrypt retry limit tests ---

// TestProcessDecryptRequests_RetryLimitDrops verifies that a decrypt request at
// maxReprocessAttempts retryCount is dropped (not re-queued) on failure.
func TestProcessDecryptRequests_RetryLimitDrops(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("cc-retry-drop")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{kernelRouter: router}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	// Request already at max retry count — should be dropped on next failure.
	requests := []types.PendingDecryptRequest{
		{
			DecryptRequest: types.DecryptRequest{
				Ciphertext:      []byte("ct-drop"),
				Label:           makeLabel(1),
				RequesterPubKey: []byte("rpk"),
			},
			RetryCount: maxReprocessAttempts,
		},
	}

	// Kernel fails for this request.
	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("kernel unavailable")).Times(1)

	k.processDecryptRequests(ctx, session, requests)

	// The request should be dropped, not re-queued.
	requeued := session.GetDecryptRequests()
	require.Empty(t, requeued, "request at max retries should be dropped, not re-queued")
}

// TestProcessDecryptRequests_RetryCountIncrements verifies that retry count
// increments on each failure and the request is re-queued until the limit.
func TestProcessDecryptRequests_RetryCountIncrements(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("cc-retry-inc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{kernelRouter: router}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	// Request at retryCount=2 — should survive (incremented to 3, still <= 5).
	requests := []types.PendingDecryptRequest{
		{
			DecryptRequest: types.DecryptRequest{
				Ciphertext:      []byte("ct-inc"),
				Label:           makeLabel(1),
				RequesterPubKey: []byte("rpk"),
			},
			RetryCount: 2,
		},
	}

	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("transient error")).Times(1)

	k.processDecryptRequests(ctx, session, requests)

	requeued := session.GetDecryptRequests()
	require.Len(t, requeued, 1, "request below limit should be re-queued")
	require.Equal(t, 3, requeued[0].RetryCount, "retry count should be incremented from 2 to 3")
}

// TestProcessDecryptRequests_SubmitFailureRetryLimit verifies that a decrypt request
// at max retries is dropped when the contract submission fails (not just kernel failure).
func TestProcessDecryptRequests_SubmitFailureRetryLimit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("cc-submit-drop")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{kernelRouter: router, contractClient: mockContract}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	// Request at max retry count — kernel succeeds but submit fails.
	requests := []types.PendingDecryptRequest{
		{
			DecryptRequest: types.DecryptRequest{
				Ciphertext:      []byte("ct-submit-drop"),
				Label:           makeLabel(1),
				RequesterPubKey: []byte("rpk"),
			},
			RetryCount: maxReprocessAttempts,
		},
	}

	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).
		Return(&types.PartialDecryptTDH2Response{
			EncryptedPartialDecryption: []byte("partial"),
			EphemeralPubKey:            []byte("eph"),
			PubShare:                   []byte("share"),
			Signature:                  []byte("sig"),
		}, nil).Times(1)

	mockContract.EXPECT().SubmitEncryptedPartialDecryptionBatch(
		gomock.Any(), gomock.Any(),
	).Return(nil, errors.New("contract rejected")).Times(1)

	k.processDecryptRequests(ctx, session, requests)

	requeued := session.GetDecryptRequests()
	require.Empty(t, requeued, "request at max retries should be dropped on submit failure")
}

// --- Tests merged from dkg_svc_worker_test.go ---

func TestStartDecryptWorker_OnlyOneInstance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping worker test in short mode")
	}

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	k := &Keeper{stateManager: sm}

	// Reset atomic guard
	decryptWorkerRunning.Store(false)

	// Start first worker
	k.StartDecryptWorker()
	require.True(t, decryptWorkerRunning.Load(), "worker should be running")

	// Second call should be a no-op (already running)
	k.StartDecryptWorker()
	// Guard should still be true (no double-start)
	require.True(t, decryptWorkerRunning.Load(), "worker should still be running after duplicate call")
}

// --- Tests merged from dkg_svc_resume_test.go ---

// NOTE: Tests that modify dkgSvcRound are NOT parallel.

// --- ResumeDKGService ---

func TestResumeDKGService_FailedSession_Registration(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	k.stateManager = sm
	k.validatorEVMAddr = testValidatorAddr

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Create a failed session
	session := &types.DKGSession{
		Round: 1,
		Phase: types.PhaseFailed,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// No on-chain registration → alreadyRegistered will be false
	dkgNetwork := &types.DKGNetwork{
		Round:        1,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{"0xother"}, // validator not in set
		IsResharing:  false,
	}

	// ResumeDKGService dispatches to resumeFailedSession which spawns
	// a goroutine for registration. Since validator is not in current set,
	// it will create session and skip key generation.
	k.ResumeDKGService(ctx, dkgNetwork)

	// Give async goroutine time to complete (registration for non-member is fast)
	// We verify the initial dispatch happened by checking the session was updated.
	// After resumeFailedSession with Registration stage and not already registered,
	// the phase should be updated to PhaseInitializing.
	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Equal(t, types.PhaseInitializing, got.Phase)
}

func TestResumeDKGService_FailedSession_AlreadyRegistered(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	k.stateManager = sm
	k.validatorEVMAddr = testValidatorAddr

	// Set up registration so isAlreadyRegistered returns true
	require.NoError(t, k.setDKGRegistration(ctx, common.HexToAddress(testValidatorAddr), &types.DKGRegistration{
		Round:     2,
		DkgPubKey: []byte("some-pub-key"), // non-empty → isAlreadyRegistered returns true
		Status:    types.DKGRegStatusVerified,
	}))

	session := &types.DKGSession{
		Round: 2,
		Phase: types.PhaseFailed,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round:        2,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{testValidatorAddr},
	}

	k.ResumeDKGService(ctx, dkgNetwork)

	got, err := sm.GetSession(2)
	require.NoError(t, err)
	// When already registered, resumeFailedSession updates phase to PhaseInitialized
	require.Equal(t, types.PhaseInitialized, got.Phase)
}

func TestResumeDKGService_StuckSession(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	k.stateManager = sm

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Session in PhaseInitializing during Registration stage = stuck
	session := &types.DKGSession{
		Round: 3,
		Phase: types.PhaseInitializing,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round: 3,
		Stage: types.DKGStageRegistration,
	}

	k.ResumeDKGService(ctx, dkgNetwork)

	got, err := sm.GetSession(3)
	require.NoError(t, err)
	// Stuck session should be marked as failed for recovery on next block
	require.Equal(t, types.PhaseFailed, got.Phase)
}

func TestResumeDKGService_NotStuck(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	k.stateManager = sm

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Session at PhaseInitialized during Registration stage = NOT stuck
	session := &types.DKGSession{
		Round: 4,
		Phase: types.PhaseInitialized,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round: 4,
		Stage: types.DKGStageRegistration,
	}

	k.ResumeDKGService(ctx, dkgNetwork)

	got, err := sm.GetSession(4)
	require.NoError(t, err)
	// Should remain unchanged
	require.Equal(t, types.PhaseInitialized, got.Phase)
}

func TestResumeDKGService_NoSession(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	k.stateManager = sm

	dkgNetwork := &types.DKGNetwork{
		Round: 99,
		Stage: types.DKGStageRegistration,
	}

	// Should not panic when session doesn't exist
	k.ResumeDKGService(ctx, dkgNetwork)
}

// --- resumeFailedSession ---

func TestResumeFailedSession_DealingStage(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	k.stateManager = sm
	k.validatorEVMAddr = testValidatorAddr

	session := &types.DKGSession{
		Round: 10,
		Phase: types.PhaseFailed,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round:        10,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{"0xother"}, // not in set
	}

	// shouldDeal will be false (not in set, no prev active)
	// so the dealing goroutine returns quickly
	k.resumeFailedSession(ctx, session, dkgNetwork)

	got, err := sm.GetSession(10)
	require.NoError(t, err)
	// Phase should be updated to PhaseInitialized (reset for dealing recovery)
	require.Equal(t, types.PhaseInitialized, got.Phase)
}

func TestResumeFailedSession_UnspecifiedStage(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	k.stateManager = sm

	session := &types.DKGSession{
		Round: 11,
		Phase: types.PhaseFailed,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round: 11,
		Stage: types.DKGStageUnspecified,
	}

	// Should be a no-op
	k.resumeFailedSession(ctx, session, dkgNetwork)

	got, err := sm.GetSession(11)
	require.NoError(t, err)
	require.Equal(t, types.PhaseFailed, got.Phase, "unspecified stage should not modify session")
}

// --- Tests merged from dkg_svc_old_cc_test.go ---

func TestGetOldCodeCommitment_NoPreviousActive(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.validatorEVMAddr = testValidatorAddr

	cc, err := k.getOldCodeCommitment(ctx)
	require.NoError(t, err)
	require.Nil(t, cc, "should return nil when no previous active round exists")
}

func TestGetOldCodeCommitment_WithPreviousActive(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.validatorEVMAddr = testValidatorAddr

	// Set up previous active round
	network := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))
	require.NoError(t, k.setLatestActiveRound(ctx, network))

	// Set up registration for validator in previous round
	expectedCC := []byte("previous-code-commitment")
	require.NoError(t, k.setDKGRegistration(ctx, common.HexToAddress(testValidatorAddr), &types.DKGRegistration{
		Round:          1,
		ValidatorAddr:  testValidatorAddr,
		DkgPubKey:      []byte("pub-key"),
		CodeCommitment: expectedCC,
	}))

	cc, err := k.getOldCodeCommitment(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedCC, cc)
}

func TestGetOldCodeCommitment_NoRegistration(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.validatorEVMAddr = testValidatorAddr

	// Set up previous active round but no registration for this validator
	network := &types.DKGNetwork{
		Round: 2,
		Stage: types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))
	require.NoError(t, k.setLatestActiveRound(ctx, network))

	_, err := k.getOldCodeCommitment(ctx)
	require.Error(t, err, "should return error when no registration exists")
}

// --- resumeFailedSession: DKGStageFinalization and DKGStageActive cases (gap 12) ---

// TestResumeFailedSession_FinalizationStage verifies that resumeFailedSession
// dispatches to handleDKGFinalization for DKGStageFinalization stage.
// The session phase is updated to PhaseDealing before spawning the goroutine.
// NOTE: not parallel — uses package-level dkgSvcRound atomic.
func TestResumeFailedSession_FinalizationStage(t *testing.T) {
	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	k.stateManager = sm
	k.validatorEVMAddr = testValidatorAddr

	session := &types.DKGSession{
		Round: 20,
		Phase: types.PhaseFailed,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round:        20,
		Stage:        types.DKGStageFinalization,
		ActiveValSet: []string{"0xother"}, // not in set → handleDKGFinalization returns early
	}

	k.resumeFailedSession(ctx, session, dkgNetwork)

	// Phase is updated synchronously before the goroutine is launched.
	got, err := sm.GetSession(20)
	require.NoError(t, err)
	// resumeFailedSession updates phase to PhaseDealing before launching the goroutine
	require.Equal(t, types.PhaseDealing, got.Phase, "finalization stage should set phase to PhaseDealing")

	// Allow the spawned goroutine to run and release the dkgSvcRound lock
	// so it does not race with the deferred resetDKGSvcRound.
	time.Sleep(50 * time.Millisecond)
}

// TestResumeFailedSession_ActiveStage verifies that a failed session that DOES
// have key material (GlobalPubKey + PubKeyShare) is completed via handleDKGComplete
// for the DKGStageActive stage. This is the unchanged happy path.
// NOTE: not parallel — uses package-level dkgSvcRound atomic.
func TestResumeFailedSession_ActiveStage(t *testing.T) {
	resetDKGSvcRound()
	defer resetDKGSvcRound()
	// Pre-mark the worker as running so handleDKGComplete does not spawn a real
	// background goroutine that would outlive the test and tick against a finished
	// mock controller.
	decryptWorkerRunning.Store(true)
	defer decryptWorkerRunning.Store(false)

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Use os.MkdirTemp instead of t.TempDir() to avoid automatic cleanup race
	// with background goroutines spawned by handleDKGComplete.
	tmpDir, err := os.MkdirTemp("", "dkg-resume-active-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint:errcheck // best-effort cleanup

	sm, err := NewStateManager(tmpDir)
	require.NoError(t, err)
	k.stateManager = sm
	k.validatorEVMAddr = testValidatorAddr

	// Session already holds key material: completion must proceed as before.
	session := &types.DKGSession{
		Round:        21,
		Phase:        types.PhaseFailed,
		GlobalPubKey: []byte("global-pub"),
		PubKeyShare:  []byte("share"),
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round:        21,
		Stage:        types.DKGStageActive,
		ActiveValSet: []string{"0xother"},
	}

	k.resumeFailedSession(ctx, session, dkgNetwork)

	// resumeFailedSession (Active stage) synchronously sets PhaseFinalized,
	// then spawns a goroutine for handleDKGComplete that mutates the same
	// session object. We must wait for the goroutine to finish before reading
	// the session, otherwise the race detector flags the concurrent read/write
	// on session.Phase.
	//
	// Wait strategy: the goroutine acquires dkgSvcRound, does work, then
	// releases it. We first wait for acquisition (non-zero), then for release
	// (back to zero). The initial sleep gives the goroutine time to start.
	time.Sleep(50 * time.Millisecond)
	for deadline := time.Now().Add(3 * time.Second); time.Now().Before(deadline); {
		if dkgSvcRound.Load() == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Goroutine has completed — safe to read session without race.
	got, err := sm.GetSession(21)
	require.NoError(t, err)
	// handleDKGComplete advances phase from Finalized to Completed.
	require.Equal(t, types.PhaseCompleted, got.Phase,
		"active stage: handleDKGComplete should advance phase to PhaseCompleted")
	require.NotEmpty(t, got.GlobalPubKey, "key material must be preserved through completion")
}

// TestResumeFailedSession_ActiveStage_MissingKeyMaterial verifies that a failed
// session WITHOUT key material in the active stage is NOT force-completed with an
// empty key. It routes to local finalization recovery; when the kernel call fails,
// the session is left PhaseFailed (visibly broken) rather than PhaseCompleted.
// NOTE: not parallel — uses package-level dkgSvcRound atomic.
func TestResumeFailedSession_ActiveStage_MissingKeyMaterial(t *testing.T) {
	resetDKGSvcRound()
	defer resetDKGSvcRound()
	decryptWorkerRunning.Store(false)
	defer decryptWorkerRunning.Store(false)

	ctrl := gomock.NewController(t)

	tmpDir, err := os.MkdirTemp("", "dkg-resume-active-nokey-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint:errcheck // best-effort cleanup

	sm, err := NewStateManager(tmpDir)
	require.NoError(t, err)

	cc := []byte("recover-nokey-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		validatorEVMAddr: testValidatorAddr,
	}

	// Failed session with no key material — the recovery path must be taken.
	session := &types.DKGSession{
		Round:          30,
		Phase:          types.PhaseFailed,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(context.Background(), session))

	// Kernel finalize keeps failing → recovery cannot recover the share. The three
	// retries run sequentially in the recovery goroutine; close done on the last one
	// so the test can synchronize deterministically.
	done := make(chan struct{})
	attempts := 0
	mockKernel.EXPECT().FinalizeDKG(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ *types.FinalizeDKGRequest, _ ...grpc.CallOption) (*types.FinalizeDKGResponse, error) {
			attempts++
			if attempts == retryAttempts {
				close(done)
			}

			return nil, errors.New("kernel finalize failed")
		}).
		Times(retryAttempts)

	dkgNetwork := &types.DKGNetwork{
		Round:        30,
		Stage:        types.DKGStageActive,
		ActiveValSet: []string{testValidatorAddr},
	}

	k.resumeFailedSession(context.Background(), session, dkgNetwork)

	// Recovery runs in a goroutine and retries the kernel with backoff. Wait for the
	// final retry, then a short grace period for the goroutine to persist MarkFailed.
	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("kernel finalize was not retried to exhaustion")
	}
	time.Sleep(200 * time.Millisecond)

	got, err := sm.GetSession(30)
	require.NoError(t, err)
	require.Equal(t, types.PhaseFailed, got.Phase,
		"missing key material: session must be left failed, not completed")
	require.NotEqual(t, types.PhaseCompleted, got.Phase,
		"missing key material: session must never reach PhaseCompleted with empty key")
	require.Empty(t, got.GlobalPubKey, "no key material should have been produced")
}

// TestResumeFailedSession_ActiveStage_RecoversKeyMaterial verifies that a failed
// session WITHOUT key material in the active stage recovers the share via a local
// kernel finalization (no on-chain finalize vote) and is then completed.
// NOTE: not parallel — uses package-level dkgSvcRound atomic.
func TestResumeFailedSession_ActiveStage_RecoversKeyMaterial(t *testing.T) {
	resetDKGSvcRound()
	defer resetDKGSvcRound()
	// Pre-mark the worker as running so handleDKGComplete does not spawn a real
	// background goroutine that would outlive the test.
	decryptWorkerRunning.Store(true)
	defer decryptWorkerRunning.Store(false)

	ctrl := gomock.NewController(t)

	tmpDir, err := os.MkdirTemp("", "dkg-resume-active-recover-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint:errcheck // best-effort cleanup

	sm, err := NewStateManager(tmpDir)
	require.NoError(t, err)

	cc := []byte("recover-ok-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		validatorEVMAddr: testValidatorAddr,
	}

	session := &types.DKGSession{
		Round:          31,
		Phase:          types.PhaseFailed,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(context.Background(), session))

	// Kernel recomputes and seals the share for the active round. No contract
	// Finalize is expected — recovery must not submit a stale on-chain vote.
	mockKernel.EXPECT().FinalizeDKG(gomock.Any(), gomock.Any()).Return(
		&types.FinalizeDKGResponse{
			ParticipantsRoot: make([]byte, 32),
			GlobalPubKey:     []byte("global-pub"),
			Signature:        []byte("sig"),
			PublicCoeffs:     [][]byte{[]byte("c1")},
			PubKeyShare:      []byte("share"),
		}, nil,
	)

	dkgNetwork := &types.DKGNetwork{
		Round:        31,
		Stage:        types.DKGStageActive,
		ActiveValSet: []string{testValidatorAddr},
		// On-chain consensus key must match the kernel-recomputed GlobalPubKey so the
		// recovery guard passes and the session completes.
		GlobalPublicKey: []byte("global-pub"),
	}

	k.resumeFailedSession(context.Background(), session, dkgNetwork)

	// Recovery runs in a goroutine (kernel finalize → handleDKGComplete) that mutates
	// the shared session pointer. handleDKGComplete acquires the round lock and releases
	// it on completion. Wait for the lock to drain before reading to avoid a race on the
	// shared session object. The kernel mock returns immediately (no retries/backoff).
	time.Sleep(100 * time.Millisecond)
	for deadline := time.Now().Add(3 * time.Second); time.Now().Before(deadline); {
		if dkgSvcRound.Load() == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	got, err := sm.GetSession(31)
	require.NoError(t, err)
	require.Equal(t, types.PhaseCompleted, got.Phase,
		"recovered session should be completed")
	require.True(t, got.IsFinalized, "completed session should be finalized")
	require.Equal(t, []byte("global-pub"), got.GlobalPubKey,
		"recovered session must be completed WITH the recovered key material")
	require.Equal(t, []byte("share"), got.PubKeyShare)
}

// TestResumeFailedSession_ActiveStage_DivergentKey verifies that when the kernel
// recomputes a GlobalPubKey that does NOT match the on-chain consensus network key,
// the session is left PhaseFailed and never completed, so the node does not submit
// partial decryptions under a divergent key.
// NOTE: not parallel — uses package-level dkgSvcRound atomic.
func TestResumeFailedSession_ActiveStage_DivergentKey(t *testing.T) {
	resetDKGSvcRound()
	defer resetDKGSvcRound()
	decryptWorkerRunning.Store(false)
	defer decryptWorkerRunning.Store(false)

	ctrl := gomock.NewController(t)

	tmpDir, err := os.MkdirTemp("", "dkg-resume-active-divergent-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint:errcheck // best-effort cleanup

	sm, err := NewStateManager(tmpDir)
	require.NoError(t, err)

	cc := []byte("recover-divergent-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		validatorEVMAddr: testValidatorAddr,
	}

	session := &types.DKGSession{
		Round:          32,
		Phase:          types.PhaseFailed,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(context.Background(), session))

	// Kernel recomputes a key that diverges from the on-chain network key.
	mockKernel.EXPECT().FinalizeDKG(gomock.Any(), gomock.Any()).Return(
		&types.FinalizeDKGResponse{
			ParticipantsRoot: make([]byte, 32),
			GlobalPubKey:     []byte("divergent-pub"),
			Signature:        []byte("sig"),
			PublicCoeffs:     [][]byte{[]byte("c1")},
			PubKeyShare:      []byte("share"),
		}, nil,
	).Times(1)

	dkgNetwork := &types.DKGNetwork{
		Round:           32,
		Stage:           types.DKGStageActive,
		ActiveValSet:    []string{testValidatorAddr},
		GlobalPublicKey: []byte("global-pub"), // does not match the kernel-recomputed key
	}

	k.resumeFailedSession(context.Background(), session, dkgNetwork)

	// Recovery runs in a goroutine (kernel finalize immediate, no retries). Wait for the
	// round lock to drain before reading the shared session.
	time.Sleep(100 * time.Millisecond)
	for deadline := time.Now().Add(3 * time.Second); time.Now().Before(deadline); {
		if dkgSvcRound.Load() == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	got, err := sm.GetSession(32)
	require.NoError(t, err)
	require.Equal(t, types.PhaseFailed, got.Phase,
		"divergent key: session must be left failed, not completed")
	require.NotEqual(t, types.PhaseCompleted, got.Phase,
		"divergent key: session must never reach PhaseCompleted")
	require.False(t, decryptWorkerRunning.Load(),
		"divergent key: decrypt worker must not be started")
}

// TestRecoverActiveSessionKeyMaterial_ConcurrentSingleFinalize verifies that when two
// blocks concurrently spawn active-round recovery for the same round, the per-round lock
// serializes them so the kernel FinalizeDKG is invoked at most once and the session is
// completed exactly once (no double-completion / phase flapping). Runs under -race.
// NOTE: not parallel — uses package-level dkgSvcRound atomic.
func TestRecoverActiveSessionKeyMaterial_ConcurrentSingleFinalize(t *testing.T) {
	resetDKGSvcRound()
	defer resetDKGSvcRound()
	// Pre-mark the worker as running so handleDKGComplete does not spawn a real
	// background goroutine that would outlive the test.
	decryptWorkerRunning.Store(true)
	defer decryptWorkerRunning.Store(false)

	ctrl := gomock.NewController(t)

	tmpDir, err := os.MkdirTemp("", "dkg-resume-active-concurrent-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint:errcheck // best-effort cleanup

	sm, err := NewStateManager(tmpDir)
	require.NoError(t, err)

	cc := []byte("recover-concurrent-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		validatorEVMAddr: testValidatorAddr,
	}

	session := &types.DKGSession{
		Round:          33,
		Phase:          types.PhaseFailed,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(context.Background(), session))

	// Slow kernel finalize: the winner holds the round lock for its whole duration, so a
	// concurrent recovery for the same round is deduplicated by tryAcquireDKGSvc. Count
	// invocations to assert it runs at most once.
	var finalizeCalls atomic.Int32
	mockKernel.EXPECT().FinalizeDKG(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ *types.FinalizeDKGRequest, _ ...grpc.CallOption) (*types.FinalizeDKGResponse, error) {
			finalizeCalls.Add(1)
			time.Sleep(200 * time.Millisecond)

			return &types.FinalizeDKGResponse{
				ParticipantsRoot: make([]byte, 32),
				GlobalPubKey:     []byte("global-pub"),
				Signature:        []byte("sig"),
				PublicCoeffs:     [][]byte{[]byte("c1")},
				PubKeyShare:      []byte("share"),
			}, nil
		}).
		MaxTimes(1)

	dkgNetwork := &types.DKGNetwork{
		Round:           33,
		Stage:           types.DKGStageActive,
		ActiveValSet:    []string{testValidatorAddr},
		GlobalPublicKey: []byte("global-pub"),
	}

	// Two blocks fire recovery for the same round concurrently.
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			k.recoverActiveSessionKeyMaterial(context.Background(), dkgNetwork)
		}()
	}
	wg.Wait()

	// Wait for the round lock (held across finalize + handleDKGComplete) to drain.
	for deadline := time.Now().Add(3 * time.Second); time.Now().Before(deadline); {
		if dkgSvcRound.Load() == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	require.Equal(t, int32(1), finalizeCalls.Load(),
		"kernel FinalizeDKG must be invoked exactly once despite two concurrent recoveries")

	got, err := sm.GetSession(33)
	require.NoError(t, err)
	require.Equal(t, types.PhaseCompleted, got.Phase,
		"concurrent recovery must complete the session exactly once")
	require.True(t, got.IsFinalized, "completed session should be finalized")
	require.Equal(t, []byte("global-pub"), got.GlobalPubKey)
}
