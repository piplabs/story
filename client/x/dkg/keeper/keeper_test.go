package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	grpc1 "google.golang.org/grpc"
)

func TestInitDKGService(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	addr := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	enclaveType := [32]byte{0x01, 0x02}

	err := k.InitDKGService(t.TempDir(), addr, enclaveType)
	require.NoError(t, err)
	require.True(t, k.isDKGSvcEnabled)
	require.Equal(t, "0x1234567890abcdef1234567890abcdef12345678", k.validatorEVMAddr)
	require.Equal(t, enclaveType, k.enclaveType)
	require.NotNil(t, k.stateManager)
}

func TestInitDKGService_InvalidDir(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	addr := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	enclaveType := [32]byte{0x01}

	// Use a path that cannot be created
	err := k.InitDKGService("/dev/null/invalid/path", addr, enclaveType)
	require.Error(t, err)
}

// mockGRPCServer implements gogoproto grpc.Server for testing RegisterProposalService.
type mockGRPCServer struct {
	services []string
}

func (m *mockGRPCServer) RegisterService(sd *grpc1.ServiceDesc, _ interface{}) {
	m.services = append(m.services, sd.ServiceName)
}

func TestRegisterProposalService(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	server := &mockGRPCServer{}
	require.NotPanics(t, func() {
		k.RegisterProposalService(server)
	})

	require.Len(t, server.services, 1)
	require.Contains(t, server.services[0], "MsgService")
}

func TestGetAuthority(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	authority := k.GetAuthority()
	require.NotEmpty(t, authority)
	require.Equal(t, "story1hmjw3pvkjtndpg8wqppwdn8udd835qpan4hm0y", authority)
}
