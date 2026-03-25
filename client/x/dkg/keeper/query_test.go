package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/piplabs/story/client/x/dkg/types"
)

// TestQuery_Params_NilRequest verifies Params returns InvalidArgument for nil request.
func TestQuery_Params_NilRequest(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.Params(ctx, nil)
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

// TestQuery_Params_ReturnsCurrentParams verifies Params returns the currently
// stored params.
func TestQuery_Params_ReturnsCurrentParams(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Override the default params
	params := types.NewParams(
		5, 10, 15, 20,
		types.DefaultDkgCommitteeRewardPortion,
		3, 3, 700, 200,
	)
	require.NoError(t, k.SetParams(ctx, params))

	resp, err := k.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, params.RegistrationPeriod, resp.Params.RegistrationPeriod)
	require.Equal(t, params.DealingPeriod, resp.Params.DealingPeriod)
	require.Equal(t, params.OperationalThreshold, resp.Params.OperationalThreshold)
}

// TestQuery_GetDKGNetwork_NilRequest verifies GetDKGNetwork returns InvalidArgument
// for nil request.
func TestQuery_GetDKGNetwork_NilRequest(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetDKGNetwork(ctx, nil)
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

// TestQuery_GetDKGNetwork_NotFound verifies GetDKGNetwork returns NotFound
// when the round does not exist.
func TestQuery_GetDKGNetwork_NotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetDKGNetwork(ctx, &types.QueryGetDKGNetworkRequest{Round: 999})
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, s.Code())
}

// TestQuery_GetDKGNetwork_Found verifies GetDKGNetwork returns the network
// when it exists.
func TestQuery_GetDKGNetwork_Found(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     1,
		Total:     5,
		Threshold: 4,
		Stage:     types.DKGStageRegistration,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	resp, err := k.GetDKGNetwork(ctx, &types.QueryGetDKGNetworkRequest{Round: 1})
	require.NoError(t, err)
	require.Equal(t, uint32(1), resp.Network.Round)
	require.Equal(t, uint32(5), resp.Network.Total)
}

// TestQuery_GetLatestDKGNetwork_NilRequest verifies GetLatestDKGNetwork returns
// InvalidArgument for nil request.
func TestQuery_GetLatestDKGNetwork_NilRequest(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetLatestDKGNetwork(ctx, nil)
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

// TestQuery_GetLatestDKGNetwork_NotFound verifies GetLatestDKGNetwork returns
// NotFound when no network exists.
func TestQuery_GetLatestDKGNetwork_NotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetLatestDKGNetwork(ctx, &types.QueryGetLatestDKGNetworkRequest{})
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, s.Code())
}

// TestQuery_GetLatestDKGNetwork_Found verifies GetLatestDKGNetwork returns
// the most recently stored network.
func TestQuery_GetLatestDKGNetwork_Found(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     3,
		Total:     7,
		Threshold: 5,
		Stage:     types.DKGStageDealing,
	}
	// setDKGNetwork also updates the latest pointer when round > current latest
	require.NoError(t, k.setDKGNetwork(ctx, network))

	resp, err := k.GetLatestDKGNetwork(ctx, &types.QueryGetLatestDKGNetworkRequest{})
	require.NoError(t, err)
	require.Equal(t, uint32(3), resp.Network.Round)
}

// TestQuery_GetAllDKGNetworks_NilRequest verifies GetAllDKGNetworks returns
// InvalidArgument for nil request.
func TestQuery_GetAllDKGNetworks_NilRequest(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetAllDKGNetworks(ctx, nil)
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

// TestQuery_GetAllDKGNetworks_Empty verifies GetAllDKGNetworks returns empty
// when no networks exist.
func TestQuery_GetAllDKGNetworks_Empty(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	resp, err := k.GetAllDKGNetworks(ctx, &types.QueryGetAllDKGNetworksRequest{})
	require.NoError(t, err)
	require.Empty(t, resp.Networks)
}

// TestQuery_GetAllDKGNetworks_Multiple verifies GetAllDKGNetworks returns all
// stored networks.
func TestQuery_GetAllDKGNetworks_Multiple(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	for _, round := range []uint32{1, 2, 3} {
		net := &types.DKGNetwork{Round: round, Total: 5, Threshold: 4, Stage: types.DKGStageActive}
		require.NoError(t, k.setDKGNetwork(ctx, net))
	}

	resp, err := k.GetAllDKGNetworks(ctx, &types.QueryGetAllDKGNetworksRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Networks, 3)
}

// TestQuery_GetDKGRegistration_NilRequest verifies GetDKGRegistration returns
// InvalidArgument for nil request.
func TestQuery_GetDKGRegistration_NilRequest(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetDKGRegistration(ctx, nil)
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

// TestQuery_GetDKGRegistration_Unimplemented verifies that GetDKGRegistration
// returns Unimplemented for a valid request (not yet implemented).
func TestQuery_GetDKGRegistration_Unimplemented(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetDKGRegistration(ctx, &types.QueryGetDKGRegistrationRequest{})
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unimplemented, s.Code())
}

// TestQuery_GetAllDKGRegistrations_NilRequest verifies GetAllDKGRegistrations
// returns InvalidArgument for nil request.
func TestQuery_GetAllDKGRegistrations_NilRequest(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetAllDKGRegistrations(ctx, nil)
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

// TestQuery_GetAllDKGRegistrations_Empty verifies GetAllDKGRegistrations returns
// empty slice when no registrations exist for the given round.
func TestQuery_GetAllDKGRegistrations_Empty(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	resp, err := k.GetAllDKGRegistrations(ctx, &types.QueryGetAllDKGRegistrationsRequest{Round: 1})
	require.NoError(t, err)
	require.Empty(t, resp.Registrations)
}

// TestQuery_GetAllVerifiedDKGRegistrations_NilRequest verifies that nil request
// returns InvalidArgument.
func TestQuery_GetAllVerifiedDKGRegistrations_NilRequest(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetAllVerifiedDKGRegistrations(ctx, nil)
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

// TestQuery_GetAllVerifiedDKGRegistrations_Found verifies that
// GetAllVerifiedDKGRegistrations returns only verified registrations.
func TestQuery_GetAllVerifiedDKGRegistrations_Found(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	require.NoError(t, k.SetParams(ctx, types.DefaultParams()))

	round := uint32(1)

	// Store one Verified and one Finalized registration
	addrVerified := common.HexToAddress("0x1111111111111111111111111111111111111111")
	addrFinalized := common.HexToAddress("0x2222222222222222222222222222222222222222")

	require.NoError(t, k.setDKGRegistration(ctx, addrVerified, &types.DKGRegistration{
		Round:         round,
		ValidatorAddr: addrVerified.Hex(),
		Index:         1,
		Status:        types.DKGRegStatusVerified,
	}))
	require.NoError(t, k.setDKGRegistration(ctx, addrFinalized, &types.DKGRegistration{
		Round:         round,
		ValidatorAddr: addrFinalized.Hex(),
		Index:         2,
		Status:        types.DKGRegStatusFinalized,
	}))

	resp, err := k.GetAllVerifiedDKGRegistrations(ctx, &types.QueryGetAllVerifiedDKGRegistrationsRequest{Round: round})
	require.NoError(t, err)
	require.Len(t, resp.Registrations, 1, "only verified registrations should be returned")
	require.Equal(t, types.DKGRegStatusVerified, resp.Registrations[0].Status)
}

// TestQuery_GetAllVerifiedDKGRegistrations_EmptyRound verifies that an empty
// slice is returned when no verified registrations exist.
func TestQuery_GetAllVerifiedDKGRegistrations_EmptyRound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	resp, err := k.GetAllVerifiedDKGRegistrations(ctx, &types.QueryGetAllVerifiedDKGRegistrationsRequest{Round: 99})
	require.NoError(t, err)
	require.Empty(t, resp.Registrations)
}

// TestQuery_GetLatestActiveDKGNetwork_NilRequest verifies that nil request
// returns InvalidArgument.
func TestQuery_GetLatestActiveDKGNetwork_NilRequest(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetLatestActiveDKGNetwork(ctx, nil)
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

// TestQuery_GetLatestActiveDKGNetwork_NotFound verifies that NotFound is returned
// when no active network exists.
func TestQuery_GetLatestActiveDKGNetwork_NotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetLatestActiveDKGNetwork(ctx, &types.QueryGetLatestActiveDKGNetworkRequest{})
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, s.Code())
}

// TestQuery_GetLatestActiveDKGNetwork_Found verifies that the latest active
// network is returned when one exists.
func TestQuery_GetLatestActiveDKGNetwork_Found(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Create an Active network and set it as the latest active round
	network := &types.DKGNetwork{
		Round:           5,
		Total:           4,
		Threshold:       3,
		Stage:           types.DKGStageActive,
		GlobalPublicKey: []byte("test-pub-key"),
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))
	// setLatestActiveRound sets LatestActiveRound which getLatestActiveDKGNetwork reads
	require.NoError(t, k.setLatestActiveRound(ctx, network))

	resp, err := k.GetLatestActiveDKGNetwork(ctx, &types.QueryGetLatestActiveDKGNetworkRequest{})
	require.NoError(t, err)
	require.Equal(t, uint32(5), resp.Network.Round)
	require.Equal(t, []byte("test-pub-key"), resp.Network.GlobalPublicKey)
}

// TestQuery_GetCDRPartials_NilRequest verifies GetCDRPartials returns
// InvalidArgument for nil request.
func TestQuery_GetCDRPartials_NilRequest(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetCDRPartials(ctx, nil)
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

// TestQuery_GetCDRPartials_InvalidHex verifies GetCDRPartials returns
// InvalidArgument for invalid pubkey hex string.
func TestQuery_GetCDRPartials_InvalidHex(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	_, err := k.GetCDRPartials(ctx, &types.QueryGetCDRPartialsRequest{
		RequesterPubKeyHex: "not-hex",
		Uuid:               1,
	})
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

// TestQuery_GetCDRPartials_NotFound verifies GetCDRPartials returns NotFound
// when no partial submissions exist.
func TestQuery_GetCDRPartials_NotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Valid hex, but no submissions stored
	resp, err := k.GetCDRPartials(ctx, &types.QueryGetCDRPartialsRequest{
		RequesterPubKeyHex: "aabbccdd",
		Uuid:               42,
	})
	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, s.Code())
	require.Nil(t, resp)
}

// TestQuery_GetCDRPartials_FoundSingleSubmission verifies GetCDRPartials returns
// a grouped result when one partial decryption submission exists.
func TestQuery_GetCDRPartials_FoundSingleSubmission(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Store a DKG network for the round so GetCDRPartials can read the threshold.
	round := uint32(1)
	network := &types.DKGNetwork{
		Round:     round,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	// Parameters for the submission.
	requesterPubKey := []byte("requester-pub-key-bytes")
	var label [32]byte
	// uuid=7 → stored in last 4 bytes of label (big-endian)
	label[28] = 0x00
	label[29] = 0x00
	label[30] = 0x00
	label[31] = 0x07
	ciphertext := []byte("some-ciphertext")
	encryptedPartial := []byte("encrypted-partial")
	ephemeralPubKey := []byte("ephemeral-pub-key")
	pubShare := []byte("pub-share")
	validator := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	// Directly set the partial decryption submission (bypasses signature verification).
	require.NoError(t, k.setPartialDecryptionSubmission(
		ctx,
		validator,
		round,
		1, // pid
		encryptedPartial,
		ephemeralPubKey,
		pubShare,
		requesterPubKey,
		label[:],
		ciphertext,
	))

	// Query using the requester pubkey hex and the uuid.
	resp, err := k.GetCDRPartials(ctx, &types.QueryGetCDRPartialsRequest{
		RequesterPubKeyHex: common.Bytes2Hex(requesterPubKey),
		Uuid:               7,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Submissions, 1, "expected one grouped result")
	require.Equal(t, round, resp.Submissions[0].Round)
	require.Equal(t, uint32(2), resp.Submissions[0].Threshold)
	require.Len(t, resp.Submissions[0].Submissions, 1)
	require.False(t, resp.Submissions[0].ThresholdMet, "threshold not met with only 1 of 2 required")
}

// TestQuery_GetCDRPartials_ThresholdMet verifies GetCDRPartials reports ThresholdMet=true
// when the number of submissions reaches the round threshold.
func TestQuery_GetCDRPartials_ThresholdMet(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	round := uint32(2)
	network := &types.DKGNetwork{
		Round:     round,
		Total:     2,
		Threshold: 2, // need 2 submissions
		Stage:     types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	requesterPubKey := []byte("another-requester-key")
	var label [32]byte
	label[31] = 0x0A // uuid=10
	ciphertext := []byte("ciphertext-abc")
	val1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	val2 := common.HexToAddress("0x2222222222222222222222222222222222222222")

	require.NoError(t, k.setPartialDecryptionSubmission(ctx, val1, round, 1, []byte("ep1"), []byte("eph1"), []byte("ps1"), requesterPubKey, label[:], ciphertext))
	require.NoError(t, k.setPartialDecryptionSubmission(ctx, val2, round, 2, []byte("ep2"), []byte("eph2"), []byte("ps2"), requesterPubKey, label[:], ciphertext))

	resp, err := k.GetCDRPartials(ctx, &types.QueryGetCDRPartialsRequest{
		RequesterPubKeyHex: common.Bytes2Hex(requesterPubKey),
		Uuid:               10,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Submissions, 1)
	require.Equal(t, round, resp.Submissions[0].Round)
	require.Len(t, resp.Submissions[0].Submissions, 2)
	require.True(t, resp.Submissions[0].ThresholdMet, "threshold should be met with 2 of 2")
}
