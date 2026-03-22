package keeper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// --- handleDecryptRequest ---

func TestHandleDecryptRequest_NilKernelRouter(t *testing.T) {
	t.Parallel()

	k := &Keeper{kernelRouter: nil}
	session := &types.DKGSession{
		Round: 1,
		Index: 1,
	}
	req := types.DecryptRequest{
		Ciphertext: []byte("encrypted"),
		Label:      make([]byte, 32),
	}

	err := k.handleDecryptRequest(context.Background(), session, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "kernel client not configured")
}

func TestHandleDecryptRequest_ZeroIndex(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}
	session := &types.DKGSession{
		Round: 1,
		Index: 0, // Not set
	}
	req := types.DecryptRequest{
		Ciphertext: []byte("encrypted"),
		Label:      make([]byte, 32),
	}

	err := k.handleDecryptRequest(context.Background(), session, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "session index not set")
}

func TestHandleDecryptRequest_MissingGlobalPubKey(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}
	session := &types.DKGSession{
		Round:        1,
		Index:        1,
		GlobalPubKey: nil, // Missing
	}
	req := types.DecryptRequest{
		Ciphertext: []byte("encrypted"),
		Label:      make([]byte, 32),
	}

	err := k.handleDecryptRequest(context.Background(), session, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing global public key")
}

func TestHandleDecryptRequest_NoKernelClient(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}
	session := &types.DKGSession{
		Round:          1,
		Index:          1,
		GlobalPubKey:   []byte("pubkey"),
		CodeCommitment: []byte("nonexistent-cc"),
	}
	req := types.DecryptRequest{
		Ciphertext: []byte("encrypted"),
		Label:      make([]byte, 32),
	}

	err := k.handleDecryptRequest(context.Background(), session, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no kernel client for session")
}

// --- processDecryptQueue ---

func TestProcessDecryptQueue_NoSessions(t *testing.T) {
	t.Parallel()

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	k := &Keeper{stateManager: sm}

	// Should not panic with no sessions
	k.processDecryptQueue(context.Background())
}

func TestProcessDecryptQueue_SessionWithNoRequests(t *testing.T) {
	t.Parallel()

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	ctx := context.Background()

	session := &types.DKGSession{
		Round:           1,
		Phase:           types.PhaseCompleted,
		DecryptRequests: nil, // No pending requests
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{stateManager: sm}

	// Should be a no-op
	k.processDecryptQueue(ctx)
}

func TestProcessDecryptQueue_FailedRequestsRetained(t *testing.T) {
	t.Parallel()

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	ctx := context.Background()

	router := NewKernelRouter(nil, nil)
	session := &types.DKGSession{
		Round:          1,
		Phase:          types.PhaseCompleted,
		Index:          1,
		GlobalPubKey:   []byte("pubkey"),
		CodeCommitment: []byte("nonexistent-cc"),
		DecryptRequests: []types.DecryptRequest{
			{
				Ciphertext:      []byte("encrypted1"),
				Label:           make([]byte, 32),
				RequesterPubKey: []byte("pub1"),
			},
		},
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	k := &Keeper{
		stateManager: sm,
		kernelRouter: router,
	}

	// All requests will fail (no kernel client for CC)
	k.processDecryptQueue(ctx)

	// Failed requests should be retained for retry
	got, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Len(t, got.GetDecryptRequests(), 1, "failed requests should be retained")
}
