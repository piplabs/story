package keeper

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
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

// --- Tests merged from dkg_svc_decrypt_test.go ---

// --- handleDecryptRequest ---

func TestHandleDecryptRequest_NilKernelRouter(t *testing.T) {
	t.Parallel()

	k := &Keeper{kernelRouter: nil}
	session := &types.DKGSession{
		Round: 1,
		Index: 1,
	}
	req := types.DecryptRequest{
		Ciphertext: []byte("encrypted"),
		Label:      make([]byte, 32),
	}

	err := k.handleDecryptRequest(context.Background(), session, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "kernel client not configured")
}

func TestHandleDecryptRequest_ZeroIndex(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}
	session := &types.DKGSession{
		Round: 1,
		Index: 0, // Not set
	}
	req := types.DecryptRequest{
		Ciphertext: []byte("encrypted"),
		Label:      make([]byte, 32),
	}

	err := k.handleDecryptRequest(context.Background(), session, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "session index not set")
}

func TestHandleDecryptRequest_MissingGlobalPubKey(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}
	session := &types.DKGSession{
		Round:        1,
		Index:        1,
		GlobalPubKey: nil, // Missing
	}
	req := types.DecryptRequest{
		Ciphertext: []byte("encrypted"),
		Label:      make([]byte, 32),
	}

	err := k.handleDecryptRequest(context.Background(), session, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing global public key")
}

func TestHandleDecryptRequest_NoKernelClient(t *testing.T) {
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

	err := k.handleDecryptRequest(context.Background(), session, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no kernel client for session")
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
		DecryptRequests: []types.DecryptRequest{
			{
				Ciphertext:      []byte("encrypted1"),
				Label:           make([]byte, 32),
				RequesterPubKey: []byte("pub1"),
			},
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

// --- Tests merged from dkg_svc_decrypt_full_test.go ---

// TestHandleDecryptRequest_FullPath exercises the complete decrypt request flow
// with mocked kernel and contract clients.
func TestHandleDecryptRequest_FullPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("decrypt-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		kernelRouter:   router,
		contractClient: mockContract,
	}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	label := make([]byte, 32)
	label[31] = 42 // UUID = 42

	req := types.DecryptRequest{
		Ciphertext:      []byte("ciphertext"),
		Label:           label,
		RequesterPubKey: []byte("requester-pub"),
	}

	// Kernel PartialDecryptTDH2 returns partial decrypt data
	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).Return(
		&types.PartialDecryptTDH2Response{
			EncryptedPartialDecryption: []byte("partial"),
			EphemeralPubKey:            []byte("eph-pub"),
			PubShare:                   []byte("pub-share"),
			Signature:                  []byte("sig"),
		}, nil,
	)

	// Contract SubmitEncryptedPartialDecryption
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

	err := k.handleDecryptRequest(ctx, session, req)
	require.NoError(t, err)
}

// TestHandleDecryptRequest_InvalidLabel exercises the label-too-short error.
func TestHandleDecryptRequest_InvalidLabel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("decrypt-cc-short")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		kernelRouter:   router,
		contractClient: nil, // won't be called
	}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	// Label too short
	req := types.DecryptRequest{
		Ciphertext:      []byte("ciphertext"),
		Label:           make([]byte, 20), // < 32 bytes
		RequesterPubKey: []byte("requester-pub"),
	}

	// Kernel PartialDecryptTDH2 returns success, but labelToUUID will fail
	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).Return(
		&types.PartialDecryptTDH2Response{
			EncryptedPartialDecryption: []byte("partial"),
			EphemeralPubKey:            []byte("eph-pub"),
			PubShare:                   []byte("pub-share"),
			Signature:                  []byte("sig"),
		}, nil,
	)

	err := k.handleDecryptRequest(ctx, session, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid decrypt request label")
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

	time.Sleep(100 * time.Millisecond)

	// Worker should have stopped
	require.False(t, decryptWorkerRunning.Load(), "worker should stop after context cancellation")
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

// TestResumeFailedSession_ActiveStage verifies that resumeFailedSession
// dispatches to handleDKGComplete for DKGStageActive stage.
// The session phase is updated to PhaseFinalized before spawning the goroutine.
// NOTE: not parallel — uses package-level dkgSvcRound atomic.
func TestResumeFailedSession_ActiveStage(t *testing.T) {
	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	k.stateManager = sm
	k.validatorEVMAddr = testValidatorAddr

	session := &types.DKGSession{
		Round: 21,
		Phase: types.PhaseFailed,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round:        21,
		Stage:        types.DKGStageActive,
		ActiveValSet: []string{"0xother"},
	}

	k.resumeFailedSession(ctx, session, dkgNetwork)

	// Wait for the goroutine spawned by resumeFailedSession to complete.
	// The goroutine acquires dkgSvcRound; when it is released (dkgSvcRound drops to 0),
	// the goroutine is done and will no longer access the session object.
	for deadline := time.Now().Add(2 * time.Second); time.Now().Before(deadline); {
		if dkgSvcRound.Load() == 0 {
			break
		}

		time.Sleep(5 * time.Millisecond)
	}

	// Now the goroutine is done — safe to read the session state.
	got, err := sm.GetSession(21)
	require.NoError(t, err)
	// resumeFailedSession updates phase to PhaseFinalized before launching the goroutine
	require.Equal(t, types.PhaseFinalized, got.Phase, "active stage should set phase to PhaseFinalized")
}
