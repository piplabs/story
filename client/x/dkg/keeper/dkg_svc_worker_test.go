package keeper

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStartDecryptWorker_OnlyOneInstance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping worker test in short mode")
	}

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	k := &Keeper{stateManager: sm}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Reset atomic guard
	decryptWorkerRunning.Store(false)

	// Start first worker
	k.StartDecryptWorker(ctx)
	require.True(t, decryptWorkerRunning.Load(), "worker should be running")

	// Second call should be a no-op (already running)
	k.StartDecryptWorker(ctx)

	// Cancel context to stop the worker
	cancel()
	time.Sleep(100 * time.Millisecond)

	// Worker should have stopped
	require.False(t, decryptWorkerRunning.Load(), "worker should stop after context cancellation")
}
