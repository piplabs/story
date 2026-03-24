package keeper

import (
	"testing"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/stretchr/testify/require"
)

// TestHandleDKGRegistration_StageNotRegistration verifies IT-REG-08: when handleDKGRegistration
// is called with dkgNetwork.Stage != DKGStageRegistration, it returns immediately without
// creating a session or calling TEE/contract. This branch is defensive; in production it is
// only reachable if a recovery path calls handleDKGRegistration with the current round
// from chain after the round has already transitioned to Dealing or later.
func TestHandleDKGRegistration_StageNotRegistration(t *testing.T) {
	k, ctx := setupDKGKeeper(t)

	// Stage=Dealing: should return at the "Stage != Registration" check.
	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageDealing,
	}
	k.handleDKGRegistration(ctx, dkgNetwork)

	// No session should exist for this round (we never reached CreateSession).
	_, err := k.stateManager.GetSession(dkgNetwork.Round)
	require.Error(t, err)
}
