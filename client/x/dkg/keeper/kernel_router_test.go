package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/keeper"
	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"

	"go.uber.org/mock/gomock"
)

func TestKernelRouter_RegisterAndGetClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte{0x01, 0x02, 0x03}

	router := keeper.NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockClient)

	// Should find the client by exact code commitment
	client, err := router.GetClient(cc)
	require.NoError(t, err)
	require.Equal(t, mockClient, client)
}

func TestKernelRouter_GetClientNoFallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte{0x01, 0x02, 0x03}

	router := keeper.NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockClient)

	// Should return error when code commitment doesn't match — no fallback
	unknownCC := []byte{0xaa, 0xbb}
	_, err := router.GetClient(unknownCC)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no kernel client for code commitment")
}

func TestKernelRouter_GetClientEmpty(t *testing.T) {
	router := keeper.NewKernelRouter(nil, nil)

	// Should return error when no clients are available
	_, err := router.GetClient([]byte{0x01})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no kernel client for code commitment")
}

func TestKernelRouter_HasClients(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := keeper.NewKernelRouter(nil, nil)

	require.False(t, router.HasClients())

	router.RegisterClient([]byte{0x01}, mockClient)
	require.True(t, router.HasClients())
}

func TestKernelRouter_Disconnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte{0x01, 0x02, 0x03}

	router := keeper.NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockClient)

	require.True(t, router.HasClients())

	router.Disconnect(cc)
	require.False(t, router.HasClients())
}

func TestKernelRouter_GetAllCodeCommitments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock1 := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mock2 := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc1 := []byte{0x01, 0x02}
	cc2 := []byte{0x03, 0x04}

	router := keeper.NewKernelRouter(nil, nil)
	router.RegisterClient(cc1, mock1)
	router.RegisterClient(cc2, mock2)

	ccs := router.GetAllCodeCommitments()
	require.Len(t, ccs, 2)

	// Check both code commitments exist (order is non-deterministic from map)
	found := make(map[string]bool)
	for _, cc := range ccs {
		found[string(cc)] = true
	}

	require.True(t, found[string(cc1)])
	require.True(t, found[string(cc2)])
}

func TestKernelRouter_MultipleClients(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock1 := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mock2 := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc1 := []byte{0x01, 0x02}
	cc2 := []byte{0x03, 0x04}

	router := keeper.NewKernelRouter(nil, nil)
	router.RegisterClient(cc1, mock1)
	router.RegisterClient(cc2, mock2)

	// Should route to the correct client by code commitment
	client1, err := router.GetClient(cc1)
	require.NoError(t, err)
	require.Equal(t, mock1, client1)

	client2, err := router.GetClient(cc2)
	require.NoError(t, err)
	require.Equal(t, mock2, client2)
}

func TestKernelRouter_GetClientWithNilCodeCommitment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := keeper.NewKernelRouter(nil, nil)
	router.RegisterClient([]byte{0x01}, mockClient)

	// Nil code commitment should return error — no fallback
	_, err := router.GetClient(nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty code commitment")
}

func TestKernelRouter_GetClientWithEmptyCodeCommitment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	router := keeper.NewKernelRouter(nil, nil)
	router.RegisterClient([]byte{0x01}, mockClient)

	// Empty code commitment should return error — no fallback
	_, err := router.GetClient([]byte{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty code commitment")
}
