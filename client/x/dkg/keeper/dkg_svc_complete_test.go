package keeper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// setupKeeperWithStateManager creates a minimal Keeper with a real StateManager for testing.
func setupKeeperWithStateManager(t *testing.T) *Keeper {
	t.Helper()

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	return &Keeper{
		stateManager: sm,
	}
}

func TestHandleDKGComplete_AlreadyCompleted(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	session := &types.DKGSession{
		Round:       1,
		Phase:       types.PhaseCompleted,
		IsFinalized: true,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{Round: 1}

	// Reset to ensure acquire succeeds
	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.handleDKGComplete(ctx, dkgNetwork)

	// Session should remain completed
	got, err := k.stateManager.GetSession(1)
	require.NoError(t, err)
	require.Equal(t, types.PhaseCompleted, got.Phase)
	require.True(t, got.IsFinalized)
}

func TestHandleDKGComplete_FromFinalized(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	session := &types.DKGSession{
		Round:       2,
		Phase:       types.PhaseFinalized,
		IsFinalized: false,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{Round: 2}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.handleDKGComplete(ctx, dkgNetwork)

	got, err := k.stateManager.GetSession(2)
	require.NoError(t, err)
	require.Equal(t, types.PhaseCompleted, got.Phase)
	require.True(t, got.IsFinalized)
}

func TestHandleDKGComplete_NotFinalized_MarksFailed(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	session := &types.DKGSession{
		Round:       3,
		Phase:       types.PhaseDealing,
		IsFinalized: false,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{Round: 3}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.handleDKGComplete(ctx, dkgNetwork)

	got, err := k.stateManager.GetSession(3)
	require.NoError(t, err)
	require.Equal(t, types.PhaseFailed, got.Phase)
}

func TestHandleDKGComplete_NoSession(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	dkgNetwork := &types.DKGNetwork{Round: 99}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Should not panic when session doesn't exist
	k.handleDKGComplete(ctx, dkgNetwork)
}

func TestHandleDKGComplete_DuplicateAcquire(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	session := &types.DKGSession{
		Round:       4,
		Phase:       types.PhaseFinalized,
		IsFinalized: false,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{Round: 4}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Pre-acquire the lock for the same round
	dkgSvcRound.Store(4)

	// Should be a no-op due to lock dedup
	k.handleDKGComplete(ctx, dkgNetwork)

	got, err := k.stateManager.GetSession(4)
	require.NoError(t, err)
	// Phase should remain unchanged since we couldn't acquire the lock
	require.Equal(t, types.PhaseFinalized, got.Phase)
}
