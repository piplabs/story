package keeper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

func TestHandleDKGProcessJustifications_NoSession(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	dkgNetwork := &types.DKGNetwork{
		Round: 99,
		Stage: types.DKGStageDealing,
	}

	justifications := []types.Justification{{Index: 1}}

	// Should log error but not panic when session doesn't exist
	k.handleDKGProcessJustifications(ctx, dkgNetwork, justifications)
}

func TestHandleDKGProcessJustifications_NoKernelClient(t *testing.T) {
	t.Parallel()

	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	router := NewKernelRouter(nil, nil)
	k.kernelRouter = router

	session := &types.DKGSession{
		Round:          1,
		Phase:          types.PhaseDealing,
		CodeCommitment: []byte("nonexistent-cc"),
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageDealing,
	}

	justifications := []types.Justification{{Index: 1}}

	// Should cache justifications when kernel client not found
	resetPendingIncoming()
	defer resetPendingIncoming()

	k.handleDKGProcessJustifications(ctx, dkgNetwork, justifications)

	// Justifications should be cached for retry
	pendingIncomingJustificationsMu.Lock()
	cached := len(pendingIncomingJustifications)
	pendingIncomingJustificationsMu.Unlock()
	require.Equal(t, 1, cached, "justifications should be cached when kernel client unavailable")
}
