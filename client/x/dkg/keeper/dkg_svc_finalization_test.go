package keeper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// NOTE: These tests are NOT parallel because they share the package-level
// dkgSvcRound atomic.

func TestHandleDKGFinalization_WrongStage(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        1,
		Stage:        types.DKGStageRegistration, // Not finalization
		ActiveValSet: []string{testValidatorAddr},
	}

	session := &types.DKGSession{
		Round: 1,
		Phase: types.PhaseDealing,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	k.handleDKGFinalization(ctx, dkgNetwork)

	got, err := k.stateManager.GetSession(1)
	require.NoError(t, err)
	require.Equal(t, types.PhaseDealing, got.Phase, "phase should not change when stage is wrong")
}

func TestHandleDKGFinalization_NotInCurRoundSet(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        2,
		Stage:        types.DKGStageFinalization,
		ActiveValSet: []string{"0xother1", "0xother2"},
	}

	session := &types.DKGSession{
		Round: 2,
		Phase: types.PhaseDealing,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	k.handleDKGFinalization(ctx, dkgNetwork)

	got, err := k.stateManager.GetSession(2)
	require.NoError(t, err)
	require.Equal(t, types.PhaseDealing, got.Phase, "phase should not change when validator not in set")
}

func TestHandleDKGFinalization_WrongPhase_MarksFailed(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        3,
		Stage:        types.DKGStageFinalization,
		ActiveValSet: []string{testValidatorAddr},
	}

	session := &types.DKGSession{
		Round: 3,
		Phase: types.PhaseInitialized,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	k.handleDKGFinalization(ctx, dkgNetwork)

	got, err := k.stateManager.GetSession(3)
	require.NoError(t, err)
	require.Equal(t, types.PhaseFailed, got.Phase, "non-dealing phase should be marked failed")
}

func TestHandleDKGFinalization_NoSession(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        99,
		Stage:        types.DKGStageFinalization,
		ActiveValSet: []string{testValidatorAddr},
	}

	// Should not panic when session doesn't exist
	k.handleDKGFinalization(ctx, dkgNetwork)
}

func TestHandleDKGFinalization_DuplicateAcquire(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Pre-acquire lock for round 5
	dkgSvcRound.Store(5)

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        5,
		Stage:        types.DKGStageFinalization,
		ActiveValSet: []string{testValidatorAddr},
	}

	session := &types.DKGSession{
		Round: 5,
		Phase: types.PhaseDealing,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	k.handleDKGFinalization(ctx, dkgNetwork)

	got, err := k.stateManager.GetSession(5)
	require.NoError(t, err)
	require.Equal(t, types.PhaseDealing, got.Phase, "lock dedup should prevent processing")
}
