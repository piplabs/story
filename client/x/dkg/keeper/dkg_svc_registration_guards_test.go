package keeper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// NOTE: These tests are NOT parallel because they share the package-level dkgSvcRound atomic.

func TestHandleDKGRegistration_WrongStage(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        1,
		Stage:        types.DKGStageDealing, // Not registration
		ActiveValSet: []string{testValidatorAddr},
	}

	k.handleDKGRegistration(ctx, dkgNetwork, nil, false)

	// No session should be created since the stage check fails before session creation
	_, err := k.stateManager.GetSession(1)
	require.Error(t, err, "no session should exist when stage is wrong")
}

func TestHandleDKGRegistration_DuplicateAcquire(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Pre-acquire lock for round 2
	dkgSvcRound.Store(2)

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        2,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{testValidatorAddr},
	}

	k.handleDKGRegistration(ctx, dkgNetwork, nil, false)

	// No session created since lock dedup prevented processing
	_, err := k.stateManager.GetSession(2)
	require.Error(t, err, "no session should exist when lock is already held")
}

func TestHandleDKGRegistration_NotInCurRoundSet_SessionCreated(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	k.enclaveType = [32]byte{0x01}

	dkgNetwork := &types.DKGNetwork{
		Round:        3,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{"0xother1", "0xother2"}, // Validator not in set
	}

	k.handleDKGRegistration(ctx, dkgNetwork, nil, false)

	// Session should be created but skip key generation
	// (old members still need a session for dealing/finalization)
	session, err := k.stateManager.GetSession(3)
	require.NoError(t, err)
	require.NotNil(t, session)
	// Phase stays at Initializing because no key gen for non-members
	require.Equal(t, types.PhaseInitializing, session.Phase)
}

func TestHandleDKGRegistration_NotInCurSet_NoKeyGen(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	k.enclaveType = [32]byte{0x01}

	dkgNetwork := &types.DKGNetwork{
		Round:        4,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{"0xnewval"}, // testValidator not in current set
		IsResharing:  true,
	}

	// No kernelRouter needed since old-only members skip key generation
	k.handleDKGRegistration(ctx, dkgNetwork, nil, false)

	session, err := k.stateManager.GetSession(4)
	require.NoError(t, err)
	// Old-only member skips key gen; phase stays at Initializing
	require.Equal(t, types.PhaseInitializing, session.Phase)
}

func TestHandleDKGRegistration_UpgradeRound_SetsOldCC(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	k.enclaveType = [32]byte{0x02}

	oldCC := []byte("old-code-commitment")
	dkgNetwork := &types.DKGNetwork{
		Round:        5,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{"0xother"}, // Validator not in current set (old-only member)
		IsResharing:  true,
		IsUpgrade:    true,
	}

	k.handleDKGRegistration(ctx, dkgNetwork, oldCC, false)

	session, err := k.stateManager.GetSession(5)
	require.NoError(t, err)
	require.NotNil(t, session)
	require.Equal(t, oldCC, session.OldCodeCommitment, "upgrade round should have old CC")
	// Old-only member (not in current set) should have CodeCommitment set to oldCC
	require.Equal(t, oldCC, session.CodeCommitment, "old-only member should use old CC as code commitment")
	require.True(t, session.IsUpgrade)
}
