package keeper

// This file contains regression tests for the LastResultsHash consensus failure
// observed at block 16366566 on Aeneid testnet.
//
// Root cause:
// CometBFT computes LastResultsHash as the Merkle hash of deterministicExecTxResult
// for each transaction — which keeps only {Code, Data, GasWanted, GasUsed}.
// Events are NOT included.
//
// Cosmos SDK's baseapp.getContextForTx always attaches a fresh InfiniteGasMeter
// to every transaction context (sdk/baseapp/baseapp.go:669), regardless of
// SkipAnteHandler. InfiniteGasMeter is NOT a NoopGasMeter: its GasConsumed()
// returns the actual accumulated amount. gaskv.Store.Get() calls
// gasMeter.ConsumeGas(ReadCostFlat) on every KV access, so KV reads inside
// message handlers accumulate into GasUsed in ExecTxResult.
//
// The bug (fixed):
// ProcessResponses called getLatestActiveDKGNetwork(ctx) only when
// isDKGSvcEnabled=true (TEE nodes). Non-TEE nodes skipped the block entirely,
// consuming 0 gas. The delta of 1003 gas units caused different GasUsed →
// different LastResultsHash → consensus failure when TEE and non-TEE validators
// coexisted.
//
// The fix:
// getLatestActiveDKGNetwork(ctx) is now called unconditionally before the
// isDKGSvcEnabled gate, so all nodes perform the same KV access and produce
// identical GasUsed.

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// TestInfiniteGasMeterTracksKVReads verifies that InfiniteGasMeter (used by
// baseapp.getContextForTx when SkipAnteHandler=true) accumulates gas on KV
// reads. This rules out the hypothesis that "CL has no gas" for system txns.
func TestInfiniteGasMeterTracksKVReads(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Store a DKGNetwork so there is something to read back.
	network := &types.DKGNetwork{Round: 1, Stage: types.DKGStageDealing}
	require.NoError(t, k.setDKGNetwork(sdkCtx, network))

	// Attach a fresh InfiniteGasMeter — same as what getContextForTx does.
	meter := storetypes.NewInfiniteGasMeter()
	require.Equal(t, uint64(0), meter.GasConsumed(), "meter starts at zero")

	sdkCtx = sdkCtx.WithGasMeter(meter)

	// Perform a KV read via the keeper (getDKGNetwork internally calls ctx.KVStore).
	_, err := k.getDKGNetwork(sdkCtx, 1)
	require.NoError(t, err)

	// InfiniteGasMeter must have accumulated gas from the KV read.
	require.Greater(t, meter.GasConsumed(), uint64(0),
		"InfiniteGasMeter should track KV read gas even with SkipAnteHandler=true; "+
			"it is NOT a NoopGasMeter — GasConsumed() returns actual gas consumed")

	t.Logf("gas consumed by one getDKGNetwork KV read: %d", meter.GasConsumed())
}
