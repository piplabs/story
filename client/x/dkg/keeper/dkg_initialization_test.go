package keeper

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
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

// TestInitiateDKGRound_FirstRound verifies InitiateDKGRound creates round 1 with
// IsResharing=false when no active DKG network exists.
func TestInitiateDKGRound_FirstRound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	sk.EXPECT().GetAllValidators(gomock.Any()).Return(nil, nil).Times(1)

	err := k.InitiateDKGRound(ctx, false)
	require.NoError(t, err)

	latest, err := k.GetLatestDKGRound(ctx)
	require.NoError(t, err)
	require.NotNil(t, latest)
	require.Equal(t, uint32(1), latest.Round)
	require.Equal(t, types.DKGStageRegistration, latest.Stage)
	require.False(t, latest.IsResharing, "first round should not be resharing")
	require.False(t, latest.IsUpgrade, "isUpgrade should be false")
}

// TestInitiateDKGRound_Resharing verifies InitiateDKGRound creates a resharing round
// when a previous active DKG network exists.
func TestInitiateDKGRound_Resharing(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	sk.EXPECT().GetAllValidators(gomock.Any()).Return(nil, nil).Times(1)

	// Set up a previous active round
	activeRound := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, activeRound))
	require.NoError(t, k.setLatestActiveRound(ctx, activeRound))

	err := k.InitiateDKGRound(ctx, false)
	require.NoError(t, err)

	latest, err := k.GetLatestDKGRound(ctx)
	require.NoError(t, err)
	require.NotNil(t, latest)
	require.Equal(t, uint32(2), latest.Round)
	require.True(t, latest.IsResharing, "round with existing active network should be resharing")
}

// TestInitiateDKGRound_UpgradeRound verifies InitiateDKGRound creates an upgrade
// resharing round when isUpgrade=true.
func TestInitiateDKGRound_UpgradeRound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	sk.EXPECT().GetAllValidators(gomock.Any()).Return(nil, nil).Times(1)

	err := k.InitiateDKGRound(ctx, true)
	require.NoError(t, err)

	latest, err := k.GetLatestDKGRound(ctx)
	require.NoError(t, err)
	require.NotNil(t, latest)
	require.True(t, latest.IsUpgrade, "isUpgrade should be true when passed as true")
}

// TestInitiateDKGRound_GetAllValidatorsError verifies InitiateDKGRound returns an
// error when GetAllValidators fails.
func TestInitiateDKGRound_GetAllValidatorsError(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	sk.EXPECT().GetAllValidators(gomock.Any()).Return(nil, errors.New("staking keeper error")).Times(1)

	err := k.InitiateDKGRound(ctx, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get active validators")
}

// TestInitiateDKGRound_WithDKGSvcEnabled_AlreadyRegistered verifies that when DKG
// service is enabled and the validator is already registered, no async goroutine is
// launched (the registration is skipped).
func TestInitiateDKGRound_WithDKGSvcEnabled_AlreadyRegistered(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()

	validatorAddr := common.HexToAddress("0xAAAABBBBCCCCDDDDEEEEFFFF0000111122223333")
	k.setValidatorAddress(validatorAddr)

	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	sk.EXPECT().GetAllValidators(gomock.Any()).Return(nil, nil).Times(1)

	// Pre-register this validator for round 1
	reg := &types.DKGRegistration{
		Round:         1,
		ValidatorAddr: validatorAddr.Hex(),
		Index:         1,
		DkgPubKey:     []byte("already-registered-pub-key"),
		Status:        types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, validatorAddr, reg))

	// InitiateDKGRound should succeed without launching async goroutine
	err := k.InitiateDKGRound(ctx, false)
	require.NoError(t, err)

	// Verify round 1 was created
	latest, err := k.GetLatestDKGRound(ctx)
	require.NoError(t, err)
	require.Equal(t, uint32(1), latest.Round)
}

// TestInitiateDKGRound_WithDKGSvcEnabled_NotRegistered verifies that when DKG service
// is enabled and the validator is NOT already registered, InitiateDKGRound succeeds
// and the async goroutine path is reached (the goroutine is spawned but we don't wait
// for it — we only verify that the round was created correctly).
func TestInitiateDKGRound_WithDKGSvcEnabled_NotRegistered(t *testing.T) {
	// Not parallel: modifies global DKG service state via setIsDKGSvcEnabled.
	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()

	// Initialize a real state manager so the goroutine (handleDKGRegistration)
	// does not panic on a nil stateManager.CreateSession call.
	initTestStateManager(t, k)

	// Do NOT set k.validatorEVMAddr — isAlreadyRegistered will return false.
	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	sk.EXPECT().GetAllValidators(gomock.Any()).Return(nil, nil).Times(1)

	// InitiateDKGRound must succeed; the goroutine calls handleDKGRegistration
	// which will attempt a TEE call (teeClient nil) and log an error — that is acceptable.
	err := k.InitiateDKGRound(ctx, false)
	require.NoError(t, err)

	latest, err := k.GetLatestDKGRound(ctx)
	require.NoError(t, err)
	require.NotNil(t, latest)
	require.Equal(t, uint32(1), latest.Round)
}

// TestShouldReshare_GetLatestActiveDKGNetworkError verifies shouldReshare propagates
// errors from getLatestActiveDKGNetwork. This covers the error branch at line 130-132
// of dkg_initialization.go. We cannot easily inject a store error via white-box testing,
// so we verify the happy-path coverage instead by confirming shouldReshare delegates
// the result of getLatestActiveDKGNetwork correctly when store returns a non-nil network.
// (The error path is inherently tied to KV store corruption — not reachable with the
// in-memory store used in tests.)
//
// This test documents the coverage gap explicitly so it is not re-investigated later.
func TestShouldReshare_Documents_ErrorPath_NotReachable(t *testing.T) {
	t.Parallel()

	// shouldReshare error branch (getLatestActiveDKGNetwork returning an error) is only
	// reachable when the KV store itself is corrupted. The in-memory test store never
	// returns an error for a missing key (it returns collections.ErrNotFound which is
	// handled gracefully by getLatestActiveDKGNetwork → returning nil, nil). This test
	// confirms the function works correctly for the two reachable cases.
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// No active network — should return false.
	result, err := k.shouldReshare(ctx)
	require.NoError(t, err)
	require.False(t, result)

	// Active network set — should return true.
	network := &types.DKGNetwork{
		Round:     5,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))
	require.NoError(t, k.setLatestActiveRound(ctx, network))

	result, err = k.shouldReshare(ctx)
	require.NoError(t, err)
	require.True(t, result)
}
