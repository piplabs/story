package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/netconf"
)

// TestHasPendingUpgradeActivation_AlreadyActivated verifies that hasPendingUpgradeActivation
// returns nil when the upgrade info is marked as already activated (IsActivated=true).
// This prevents re-initialization during a TEE binary swap.
func TestHasPendingUpgradeActivation_AlreadyActivated(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx).WithBlockHeight(200)

	upgradeInfo := &types.KernelUpgradeInfo{
		UpgradeVersion:   "v2.1.0",
		ActivationHeight: 100,
		IsActivated:      true, // already activated → should be filtered out
	}
	require.NoError(t, k.SetKernelUpgradeInfo(ctx, upgradeInfo))

	result, err := k.hasPendingUpgradeActivation(sdkCtx, 200)
	require.NoError(t, err)
	require.Nil(t, result, "already-activated upgrade should be filtered out")
}

// TestHasPendingUpgradeActivation_NotYetReached verifies that hasPendingUpgradeActivation
// returns nil when the current height is below the activation height.
func TestHasPendingUpgradeActivation_NotYetReached(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	upgradeInfo := &types.KernelUpgradeInfo{
		UpgradeVersion:   "v2.1.0",
		ActivationHeight: 500,
		IsActivated:      false,
	}
	require.NoError(t, k.SetKernelUpgradeInfo(ctx, upgradeInfo))

	result, err := k.hasPendingUpgradeActivation(ctx, 100) // 100 < 500
	require.NoError(t, err)
	require.Nil(t, result, "upgrade not yet reached should return nil")
}

// TestHasPendingUpgradeActivation_ReadyToActivate verifies that hasPendingUpgradeActivation
// returns the upgrade info when the current height meets or exceeds the activation height.
func TestHasPendingUpgradeActivation_ReadyToActivate(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	upgradeInfo := &types.KernelUpgradeInfo{
		UpgradeVersion:   "v2.1.0",
		ActivationHeight: 300,
		IsActivated:      false,
	}
	require.NoError(t, k.SetKernelUpgradeInfo(ctx, upgradeInfo))

	result, err := k.hasPendingUpgradeActivation(ctx, 300) // exactly at activation height
	require.NoError(t, err)
	require.NotNil(t, result, "upgrade ready to activate should be returned")
	require.Equal(t, int64(300), result.ActivationHeight)
}

// TestBeginBlocker_PrunesDecryptRequestsAtCleanupInterval verifies that BeginBlocker
// calls pruneTimedOutDecryptRequests when the block height is a multiple of the cleanup
// interval (DecryptRequestRegistryCleanupInterval = 1000).
// The pruning branch is only reached when a DKG round already exists (latestRound != nil).
func TestBeginBlocker_PrunesDecryptRequestsAtCleanupInterval(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx).
		WithChainID(netconf.TestChainID).
		WithBlockHeight(1000)

	// Pre-create a DKG network so BeginBlocker reaches the pruning code path.
	// Without an existing round, BeginBlocker returns immediately after InitiateDKGRound.
	existingRound := &types.DKGNetwork{
		Round: 1,
		Total: 3, Threshold: 2,
		Stage: types.DKGStageRegistration,
	}
	require.NoError(t, k.setDKGNetwork(sdkCtx, existingRound))

	// Store a decrypt request that has timed out: height=1, current=1000 → 999 > 200
	requesterPubKey := []byte("prune-req-key")
	label := []byte("prune-label")
	ciphertext := []byte("prune-cipher")

	require.NoError(t, k.setDecryptRequest(sdkCtx, requesterPubKey, label, types.DecryptRequest{
		Round:           1,
		Ciphertext:      ciphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
		Height:          1, // stored at block 1
	}))

	require.NoError(t, k.SetParams(sdkCtx, types.DefaultParams()))

	err := k.BeginBlocker(sdkCtx)
	require.NoError(t, err)

	// The request should have been pruned (height 1000 - 1 = 999 > 200 timeout)
	_, found, err := k.getDecryptRequest(sdkCtx, requesterPubKey, label, 1, ciphertext)
	require.NoError(t, err)
	require.False(t, found, "timed-out decrypt request should have been pruned")
}

// TestBeginBlocker_NoPruneAtNonCleanupInterval verifies that BeginBlocker does NOT
// prune decrypt requests when the block height is not a multiple of the cleanup interval.
func TestBeginBlocker_NoPruneAtNonCleanupInterval(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	// Height 500 is above V200 activation (110) but not a multiple of 1000
	sdkCtx := sdk.UnwrapSDKContext(ctx).
		WithChainID(netconf.TestChainID).
		WithBlockHeight(500)

	// Pre-create a DKG network so BeginBlocker reaches the pruning code path
	existingRound := &types.DKGNetwork{
		Round: 1,
		Total: 3, Threshold: 2,
		Stage: types.DKGStageRegistration,
	}
	require.NoError(t, k.setDKGNetwork(sdkCtx, existingRound))

	// Store a decrypt request that is NOT timed out at height 500: 500-400=100 ≤ 200
	requesterPubKey := []byte("no-prune-req-key")
	label := []byte("no-prune-label")
	ciphertext := []byte("no-prune-cipher")

	require.NoError(t, k.setDecryptRequest(sdkCtx, requesterPubKey, label, types.DecryptRequest{
		Round:           1,
		Ciphertext:      ciphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
		Height:          400, // 500 - 400 = 100 ≤ 200 → not timed out
	}))

	require.NoError(t, k.SetParams(sdkCtx, types.DefaultParams()))

	err := k.BeginBlocker(sdkCtx)
	require.NoError(t, err)

	// The request should still be present (not timed out, not pruned)
	_, found, err := k.getDecryptRequest(sdkCtx, requesterPubKey, label, 1, ciphertext)
	require.NoError(t, err)
	require.True(t, found, "non-timed-out decrypt request should NOT be pruned")
}

// TestBeginBlocker_DKGSvcEnabledCallsResume verifies that when isDKGSvcEnabled=true
// and a DKG round exists, BeginBlocker calls ResumeDKGService and
// reprocessPendingIncomingData without panicking.
func TestBeginBlocker_DKGSvcEnabledCallsResume(t *testing.T) {
	// Not parallel: modifies global DKG service round state
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Set a pre-existing DKG round so BeginBlocker doesn't return early
	sdkCtx := sdk.UnwrapSDKContext(ctx).
		WithChainID(netconf.TestChainID).
		WithBlockHeight(500) // V200 active (>110), not cleanup interval (not multiple of 1000)

	existingRound := &types.DKGNetwork{
		Round: 1,
		Total: 3, Threshold: 2,
		Stage: types.DKGStageRegistration,
	}
	require.NoError(t, k.setDKGNetwork(sdkCtx, existingRound))

	// Enable DKG service (which triggers the isDKGSvcEnabled branch)
	k.setIsDKGSvcEnabled()

	// ResumeDKGService requires a stateManager — initialize one
	initTestStateManager(t, k)

	require.NoError(t, k.SetParams(sdkCtx, types.DefaultParams()))

	// Should not panic even with isDKGSvcEnabled=true and minimal setup
	err := k.BeginBlocker(sdkCtx)
	require.NoError(t, err)
}

// TestBeginBlocker_UpgradeActivationInitiatesResharing verifies that when a pending
// upgrade reaches its activation height, BeginBlocker marks it as activated and
// calls InitiateDKGRound with isUpgrade=true.
// After BeginBlocker marks the upgrade as activated, GetPendingUpgrade returns nil
// (the entry is filtered because IsActivated=true). We verify via GetKernelUpgradeInfo
// directly using the upgrade version.
func TestBeginBlocker_UpgradeActivationInitiatesResharing(t *testing.T) {
	t.Parallel()

	env := setupDKGLifecycleEnv(t, 3)

	// Initialize an existing DKG round (BeginBlocker won't initiate if latestRound == nil)
	// First BeginBlocker call → creates round 1
	env.sk.EXPECT().GetAllValidators(gomock.Any()).Return(env.mockValidators(), nil).AnyTimes()
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Set a pending upgrade with activation height = current + 1
	currentHeight := env.sdkCtx.BlockHeight()
	upgradeInfo := &types.KernelUpgradeInfo{
		UpgradeVersion:   "v2.1.0",
		ActivationHeight: currentHeight + 1,
		IsActivated:      false,
	}
	require.NoError(t, env.keeper.SetKernelUpgradeInfo(env.sdkCtx, upgradeInfo))

	// Advance to activation height
	env.advanceToHeight(currentHeight + 1)

	// BeginBlocker should detect the pending upgrade, mark it activated, and initiate a new round
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// GetPendingUpgrade returns nil for activated entries (IsActivated=true is filtered out).
	// Verify via GetKernelUpgradeInfo that the entry is marked as activated.
	info, err := env.keeper.GetKernelUpgradeInfo(env.sdkCtx, "v2.1.0")
	require.NoError(t, err)
	require.NotNil(t, info)
	require.True(t, info.IsActivated, "upgrade info should be marked as activated after BeginBlocker")

	// Verify a new round was created (round 2, upgrade resharing)
	latest, err := env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, uint32(2), latest.Round)
	require.True(t, latest.IsUpgrade, "new round should be an upgrade resharing round")
}
