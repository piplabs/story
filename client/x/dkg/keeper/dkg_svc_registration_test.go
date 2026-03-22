package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
)

// --- resolveRegistrationKernelClient ---

func TestResolveRegistrationKernelClient_NormalRound_WithCC(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte("test-cc-normal")

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockClient)

	k := &Keeper{kernelRouter: router}

	client, resolvedCC, err := k.resolveRegistrationKernelClient(false, cc)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, cc, resolvedCC)
}

func TestResolveRegistrationKernelClient_NormalRound_NilCC_FirstClient(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte("first-client-cc")

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockClient)

	k := &Keeper{kernelRouter: router}

	// No previous CC — should fall back to first connected client
	client, resolvedCC, err := k.resolveRegistrationKernelClient(false, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, cc, resolvedCC)
}

func TestResolveRegistrationKernelClient_NormalRound_NoClients(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}

	_, _, err := k.resolveRegistrationKernelClient(false, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no kernel clients available")
}

func TestResolveRegistrationKernelClient_Upgrade_FindsNewBinary(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	oldCC := []byte("old-binary-cc")
	newCC := []byte("new-binary-cc")

	oldClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	newClient := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(oldCC, oldClient)
	router.RegisterClient(newCC, newClient)

	k := &Keeper{kernelRouter: router}

	// Upgrade: pass old CC, should find the new binary
	client, resolvedCC, err := k.resolveRegistrationKernelClient(true, oldCC)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotEqual(t, oldCC, resolvedCC, "should return new binary CC, not old")
	require.Equal(t, newCC, resolvedCC)
}

func TestResolveRegistrationKernelClient_Upgrade_NilOldCC(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}

	_, _, err := k.resolveRegistrationKernelClient(true, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "old code commitment required")
}

func TestResolveRegistrationKernelClient_Upgrade_NoNewBinary(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	oldCC := []byte("only-old-cc")
	oldClient := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(oldCC, oldClient)

	k := &Keeper{kernelRouter: router}

	// Only old binary connected — should fail
	_, _, err := k.resolveRegistrationKernelClient(true, oldCC)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no new kernel client found for upgrade")
}

// --- getClientWithReconnect ---

func TestGetClientWithReconnect_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	cc := []byte("reconnect-cc")
	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockClient)

	k := &Keeper{kernelRouter: router}

	client, err := k.getClientWithReconnect(cc)
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestGetClientWithReconnect_NotFound(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}

	// No clients registered — should fail even after reconnection attempt
	_, err := k.getClientWithReconnect([]byte("missing-cc"))
	require.Error(t, err)
}
