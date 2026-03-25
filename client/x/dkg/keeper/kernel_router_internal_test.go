package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"

	"go.uber.org/mock/gomock"
)

// TestKernelRouter_DisconnectedEndpoints_AllConnected verifies that
// disconnectedEndpoints returns empty when all endpoints are connected.
func TestKernelRouter_DisconnectedEndpoints_AllConnected(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := NewKernelRouter([]string{"ep1", "ep2"}, nil)

	// Simulate both endpoints connected by setting ccByEP
	r.ccByEP["ep1"] = "deadbeef"
	r.ccByEP["ep2"] = "cafebabe"

	disconnected := r.disconnectedEndpoints()
	require.Empty(t, disconnected)
}

// TestKernelRouter_DisconnectedEndpoints_SomeDisconnected verifies that
// disconnectedEndpoints returns only endpoints without a registered connection.
// maxKernelEndpoints=2, so we use 2 endpoints with only 1 connected.
func TestKernelRouter_DisconnectedEndpoints_SomeDisconnected(t *testing.T) {
	t.Parallel()

	r := NewKernelRouter([]string{"ep1", "ep2"}, nil)

	// Only ep2 has a connection
	r.ccByEP["ep2"] = "deadbeef"

	disconnected := r.disconnectedEndpoints()
	require.Len(t, disconnected, 1)
	require.Contains(t, disconnected, "ep1")
}

// TestKernelRouter_DisconnectedEndpoints_NoneConnected verifies that all
// endpoints are returned as disconnected when none have been connected.
func TestKernelRouter_DisconnectedEndpoints_NoneConnected(t *testing.T) {
	t.Parallel()

	r := NewKernelRouter([]string{"ep1", "ep2"}, nil)

	disconnected := r.disconnectedEndpoints()
	require.Len(t, disconnected, 2)
}

// TestKernelRouter_DisconnectedEndpoints_NoEndpoints verifies that an empty
// endpoint list produces no disconnected entries.
func TestKernelRouter_DisconnectedEndpoints_NoEndpoints(t *testing.T) {
	t.Parallel()

	r := NewKernelRouter(nil, nil)

	disconnected := r.disconnectedEndpoints()
	require.Empty(t, disconnected)
}

// TestKernelRouter_NewKernelRouter_PanicsOnTooManyEndpoints verifies that
// NewKernelRouter panics when more than maxKernelEndpoints are provided.
func TestKernelRouter_NewKernelRouter_PanicsOnTooManyEndpoints(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		NewKernelRouter([]string{"ep1", "ep2", "ep3"}, nil) // maxKernelEndpoints = 2
	})
}

// TestKernelRouter_NewKernelRouter_ExactlyMaxEndpoints verifies that passing
// exactly maxKernelEndpoints endpoints does not panic.
func TestKernelRouter_NewKernelRouter_ExactlyMaxEndpoints(t *testing.T) {
	t.Parallel()

	require.NotPanics(t, func() {
		NewKernelRouter([]string{"ep1", "ep2"}, nil) // exactly 2 = maxKernelEndpoints
	})
}

// TestKernelRouter_RegisterClientForEndpoint_MovesClient verifies that
// RegisterClientForEndpoint moves the client from endpoint key to code commitment key.
func TestKernelRouter_RegisterClientForEndpoint_MovesClient(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte{0xDE, 0xAD, 0xBE, 0xEF}

	r := NewKernelRouter(nil, nil)

	// Simulate a pending client registered by endpoint key (as ConnectAndDiscover
	// would do before calling RegisterClientForEndpoint).
	// For testing we directly set the client under endpoint key.
	r.clients["ep1"] = mockClient

	// RegisterClientForEndpoint should move client to code commitment key
	r.RegisterClientForEndpoint("ep1", cc)

	// Verify old endpoint key is gone
	_, oldKeyExists := r.clients["ep1"]
	require.False(t, oldKeyExists, "client should be moved away from endpoint key")

	// Verify the client is accessible via code commitment
	client, err := r.GetClient(cc)
	require.NoError(t, err)
	require.Equal(t, mockClient, client)
}

// TestKernelRouter_RegisterClientForEndpoint_NoopForMissingEndpoint verifies
// that RegisterClientForEndpoint is a no-op when the endpoint has no pending client.
func TestKernelRouter_RegisterClientForEndpoint_NoopForMissingEndpoint(t *testing.T) {
	t.Parallel()

	r := NewKernelRouter(nil, nil)

	// No client under "ep-unknown" — should not panic or add anything
	require.NotPanics(t, func() {
		r.RegisterClientForEndpoint("ep-unknown", []byte{0x01})
	})

	require.False(t, r.HasClients())
}

// TestKernelRouter_TryReconnect_NoEndpoints verifies that TryReconnect with
// no configured endpoints does nothing (no panic, no side effects).
func TestKernelRouter_TryReconnect_NoEndpoints(t *testing.T) {
	t.Parallel()

	r := NewKernelRouter(nil, nil)

	// Should not panic and should not attempt any connections
	require.NotPanics(t, func() {
		r.TryReconnect()
	})
}

// TestKernelRouter_TryReconnect_AllAlreadyConnected verifies that TryReconnect
// does nothing when all endpoints are already connected.
func TestKernelRouter_TryReconnect_AllAlreadyConnected(t *testing.T) {
	t.Parallel()

	r := NewKernelRouter([]string{"ep1"}, nil)

	// Simulate ep1 already connected
	r.mu.Lock()
	r.ccByEP["ep1"] = "deadbeef"
	r.mu.Unlock()

	// TryReconnect should observe no disconnected endpoints and return quickly
	require.NotPanics(t, func() {
		r.TryReconnect()
	})
}

// TestKernelRouter_Disconnect_NonexistentCodeCommitment verifies that
// Disconnect is safe to call for a code commitment that was never registered.
func TestKernelRouter_Disconnect_NonexistentCodeCommitment(t *testing.T) {
	t.Parallel()

	r := NewKernelRouter(nil, nil)

	require.NotPanics(t, func() {
		r.Disconnect([]byte{0xAA, 0xBB})
	})
}

// TestKernelRouter_RegisterClient_OverwritesPreviousClient verifies that
// registering a new client for the same code commitment replaces the old one.
func TestKernelRouter_RegisterClient_OverwritesPreviousClient(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client1 := dkgtestutil.NewMockKernelServiceClient(ctrl)
	client2 := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte{0x01, 0x02}

	r := NewKernelRouter(nil, nil)
	r.RegisterClient(cc, client1)
	r.RegisterClient(cc, client2) // overwrites client1

	got, err := r.GetClient(cc)
	require.NoError(t, err)
	require.Equal(t, client2, got, "second registration should overwrite first")
}

// mockCloser implements io.Closer for testing.
type mockCloser struct {
	closed bool
}

func (m *mockCloser) Close() error {
	m.closed = true
	return nil
}

// TestKernelRouter_Disconnect_WithCloser verifies that Disconnect calls Close()
// on the underlying connection and removes both the client and closer entries.
func TestKernelRouter_Disconnect_WithCloser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte{0xDE, 0xAD}
	ccHex := "dead"

	closer := &mockCloser{}

	r := NewKernelRouter(nil, nil)
	r.clients[ccHex] = mockClient
	r.closers[ccHex] = closer

	require.True(t, r.HasClients())

	r.Disconnect(cc)

	require.False(t, r.HasClients())
	require.True(t, closer.closed, "closer should have been called")

	_, closerExists := r.closers[ccHex]
	require.False(t, closerExists, "closer entry should be removed")
}

// TODO_CDR078: Characterization test for audit finding CDR-078.
//
// BUG LOCATION: client/x/dkg/keeper/kernel_router.go:162-174 — Disconnect()
//
//	Disconnect removes the client from r.clients and r.closers but does NOT
//	remove the corresponding r.ccByEP[endpoint] entry. This means
//	disconnectedEndpoints() never returns the endpoint again, preventing
//	TryReconnect from reconnecting after a disconnect.
//
// CURRENT BEHAVIOR (BUG): After Disconnect, the endpoint is still in ccByEP,
//
//	so disconnectedEndpoints() thinks it is still connected.
//
// EXPECTED BEHAVIOR AFTER FIX: Disconnect should also remove the ccByEP entry.
//
// HOW TO UPDATE AFTER FIX:
//  1. Remove TODO_CDR078_ prefix from function name
//  2. Assert that disconnectedEndpoints() includes "ep1" after Disconnect
func TestTODO_CDR078_Disconnect_DoesNotRemoveCcByEP(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte{0xDE, 0xAD}
	ccHex := "dead"

	r := NewKernelRouter([]string{"ep1"}, nil)
	r.clients[ccHex] = mockClient
	r.closers[ccHex] = &mockCloser{}
	r.ccByEP["ep1"] = ccHex

	require.True(t, r.HasClients())
	require.Empty(t, r.disconnectedEndpoints(), "ep1 should be connected")

	r.Disconnect(cc)

	require.False(t, r.HasClients(), "client should be removed")

	// BUG: ccByEP entry NOT removed, so endpoint still appears connected
	disconnected := r.disconnectedEndpoints()
	require.Empty(t, disconnected,
		"CDR-078: disconnectedEndpoints() returns empty because ccByEP still present")

	_, ccByEPStillPresent := r.ccByEP["ep1"]
	require.True(t, ccByEPStillPresent,
		"CDR-078: ccByEP['ep1'] still maps to '%s' after Disconnect", ccHex)

	t.Logf("CDR-078 characterization:")
	t.Logf("  After Disconnect: clients=%v, closers=%v, ccByEP=%v", r.clients, r.closers, r.ccByEP)
	t.Logf("  BUG: ccByEP['ep1'] still present -> TryReconnect will never reconnect ep1")
	t.Logf("  After fix: ccByEP entry should be deleted during Disconnect")
}

// TestKernelRouter_Disconnect_ClientWithoutCloser verifies that Disconnect
// removes the client even when no closer was registered (only client in map).
func TestKernelRouter_Disconnect_ClientWithoutCloser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte{0x01}
	ccHex := "01"

	r := NewKernelRouter(nil, nil)
	r.clients[ccHex] = mockClient
	// No closer registered

	r.Disconnect(cc)

	require.False(t, r.HasClients())
}
