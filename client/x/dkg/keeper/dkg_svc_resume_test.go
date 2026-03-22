package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

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
