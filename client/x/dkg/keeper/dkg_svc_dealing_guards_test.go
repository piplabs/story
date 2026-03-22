package keeper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

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
