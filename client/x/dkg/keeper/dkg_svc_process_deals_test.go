package keeper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

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
