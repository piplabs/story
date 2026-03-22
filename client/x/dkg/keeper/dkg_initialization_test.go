package keeper

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// TestIsAlreadyRegistered_NotRegistered verifies isAlreadyRegistered returns
// false when no registration exists for the round.
func TestIsAlreadyRegistered_NotRegistered(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	// validatorEVMAddr is empty by default (not set via InitDKGService)

	result := k.isAlreadyRegistered(ctx, 1)
	require.False(t, result)
}

// TestIsAlreadyRegistered_RegisteredWithPubKey verifies isAlreadyRegistered
// returns true when a registration with a non-empty DkgPubKey exists.
func TestIsAlreadyRegistered_RegisteredWithPubKey(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	addr := common.HexToAddress("0xAAAABBBBCCCCDDDDEEEEFFFF0000111122223333")
	k.setValidatorAddress(addr)

	// Register with a DkgPubKey
	reg := &types.DKGRegistration{
		Round:         1,
		ValidatorAddr: addr.Hex(),
		Index:         1,
		DkgPubKey:     []byte("my-dkg-pub-key"),
		CommPubKey:    []byte("my-comm-pub-key"),
		Status:        types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, addr, reg))

	result := k.isAlreadyRegistered(ctx, 1)
	require.True(t, result, "validator with registered DkgPubKey should return true")
}

// TestIsAlreadyRegistered_RegisteredWithoutPubKey verifies isAlreadyRegistered
// returns false when the registration has an empty DkgPubKey.
func TestIsAlreadyRegistered_RegisteredWithoutPubKey(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	addr := common.HexToAddress("0xAAAABBBBCCCCDDDDEEEEFFFF0000111122223333")
	k.setValidatorAddress(addr)

	// Register without a DkgPubKey (empty)
	reg := &types.DKGRegistration{
		Round:         1,
		ValidatorAddr: addr.Hex(),
		Index:         1,
		DkgPubKey:     []byte{}, // empty
		Status:        types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, addr, reg))

	result := k.isAlreadyRegistered(ctx, 1)
	require.False(t, result, "registration with empty DkgPubKey should return false")
}

// TestSetValidatorAddress verifies that setValidatorAddress stores the address
// in lowercase hex format (strings.ToLower(addr.Hex())).
func TestSetValidatorAddress(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	addr := common.HexToAddress("0xABCDEF1234567890ABCDEF1234567890ABCDEF12")
	k.setValidatorAddress(addr)

	// setValidatorAddress calls strings.ToLower(addr.Hex()), so the stored
	// value is always lowercase.
	require.Equal(t, strings.ToLower(addr.Hex()), k.validatorEVMAddr,
		"stored address should be lowercase hex")
}

// TestSetIsDKGSvcEnabled verifies that setIsDKGSvcEnabled sets the flag.
func TestSetIsDKGSvcEnabled(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	require.False(t, k.isDKGSvcEnabled)
	k.setIsDKGSvcEnabled()
	require.True(t, k.isDKGSvcEnabled)
}

// TestShouldReshare_NoActiveNetwork verifies shouldReshare returns false
// when no active DKG network exists.
func TestShouldReshare_NoActiveNetwork(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	result, err := k.shouldReshare(ctx)
	require.NoError(t, err)
	require.False(t, result)
}

// TestShouldReshare_WithActiveNetwork verifies shouldReshare returns true
// when an active DKG network exists.
func TestShouldReshare_WithActiveNetwork(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Set up an active DKG network
	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))
	require.NoError(t, k.setLatestActiveRound(ctx, network))

	result, err := k.shouldReshare(ctx)
	require.NoError(t, err)
	require.True(t, result)
}
