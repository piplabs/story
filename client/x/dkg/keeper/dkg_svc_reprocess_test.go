package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

func TestReprocessPendingIncomingData_NoPending(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	k := setupKeeperWithStateManager(t)
	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageDealing,
	}

	// Should be a no-op when no pending data exists
	k.reprocessPendingIncomingData(dkgNetwork)
}

func TestReprocessPendingIncomingData_WrongStage_FlushesData(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	// Pre-populate pending data
	pendingIncomingDealsMu.Lock()
	pendingIncomingDeals = []types.Deal{{Index: 1}}
	pendingIncomingDealsMu.Unlock()

	pendingIncomingResponsesMu.Lock()
	pendingIncomingResponses = []types.Response{{Index: 2}}
	pendingIncomingResponsesMu.Unlock()

	k := setupKeeperWithStateManager(t)
	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageFinalization, // Not dealing
	}

	k.reprocessPendingIncomingData(dkgNetwork)

	// Pending data should be flushed when stage is not dealing
	pendingIncomingDealsMu.Lock()
	require.Nil(t, pendingIncomingDeals, "deals should be flushed when stage is not dealing")
	pendingIncomingDealsMu.Unlock()

	pendingIncomingResponsesMu.Lock()
	require.Nil(t, pendingIncomingResponses, "responses should be flushed when stage is not dealing")
	pendingIncomingResponsesMu.Unlock()
}

func TestReprocessPendingIncomingData_NilKernelRouter(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	pendingIncomingDealsMu.Lock()
	pendingIncomingDeals = []types.Deal{{Index: 1}}
	pendingIncomingDealsMu.Unlock()

	k := &Keeper{kernelRouter: nil}
	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageDealing,
	}

	// Should be a no-op when kernel router is nil
	k.reprocessPendingIncomingData(dkgNetwork)

	// Data should still be pending (not flushed, not processed)
	pendingIncomingDealsMu.Lock()
	require.Len(t, pendingIncomingDeals, 1, "data should remain when kernel router is nil")
	pendingIncomingDealsMu.Unlock()
}

func TestReprocessPendingIncomingData_NoClients(t *testing.T) {
	resetPendingIncoming()
	defer resetPendingIncoming()

	pendingIncomingDealsMu.Lock()
	pendingIncomingDeals = []types.Deal{{Index: 1}}
	pendingIncomingDealsMu.Unlock()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}
	dkgNetwork := &types.DKGNetwork{
		Round: 1,
		Stage: types.DKGStageDealing,
	}

	// No clients connected, TryReconnect should be called but since
	// no endpoints are configured, nothing happens
	k.reprocessPendingIncomingData(dkgNetwork)

	// Data should still be pending
	pendingIncomingDealsMu.Lock()
	require.Len(t, pendingIncomingDeals, 1, "data should remain when no clients available")
	pendingIncomingDealsMu.Unlock()
}
