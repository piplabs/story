package keeper

import (
	"testing"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"
)

// TestBeginDealing_BelowMinRegistrations verifies that BeginDealing skips
// to the next round when verified registrations are below the minimum.
func TestBeginDealing_BelowMinRegistrations(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// SkipToNextRound -> InitiateDKGRound -> GetAllValidators
	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	sk.EXPECT().GetAllValidators(gomock.Any()).Return([]stakingtypes.Validator{}, nil).AnyTimes()

	params := types.DefaultParams()
	params.MinReqRegisteredParticipants = 5 // require 5 but we only register 2
	require.NoError(t, k.SetParams(ctx, params))

	round := uint32(1)
	latestRound := &types.DKGNetwork{
		Round: round,
		Stage: types.DKGStageRegistration,
	}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))

	// Register only 2 verified validators (below min of 5)
	for i, addrHex := range []string{
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
	} {
		addr := common.HexToAddress(addrHex)
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:     round,
			Index:     uint32(i + 1),
			DkgPubKey: []byte("pub"),
			Status:    types.DKGRegStatusVerified,
		}))
	}

	err := k.BeginDealing(ctx, latestRound)
	require.NoError(t, err)

	// Verify a new round was created (SkipToNextRound)
	nextRound, err := k.getDKGNetwork(ctx, round+1)
	require.NoError(t, err)
	require.Equal(t, round+1, nextRound.Round)
}

// TestBeginDealing_MeetsMinRegistrations verifies that BeginDealing succeeds
// when verified registrations meet the minimum threshold.
func TestBeginDealing_MeetsMinRegistrations(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	params := types.DefaultParams()
	params.MinReqRegisteredParticipants = 3
	params.OperationalThreshold = 667 // 66.7%
	require.NoError(t, k.SetParams(ctx, params))

	round := uint32(1)
	latestRound := &types.DKGNetwork{
		Round: round,
		Stage: types.DKGStageRegistration,
	}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))

	// Register 3 verified validators
	for i, addrHex := range []string{
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
		"0x3333333333333333333333333333333333333333",
	} {
		addr := common.HexToAddress(addrHex)
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:     round,
			Index:     uint32(i + 1),
			DkgPubKey: []byte("pub"),
			Status:    types.DKGRegStatusVerified,
		}))
	}

	err := k.BeginDealing(ctx, latestRound)
	require.NoError(t, err)

	// Network should be updated with total and threshold
	updated, err := k.getDKGNetwork(ctx, round)
	require.NoError(t, err)
	require.Equal(t, uint32(3), updated.Total)
	require.True(t, updated.Threshold > 0)
}

// TestBeginDealing_WithDKGSvcEnabled verifies that BeginDealing spawns
// an async dealing goroutine when isDKGSvcEnabled is true.
func TestBeginDealing_WithDKGSvcEnabled(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)
	k.stateManager = sm
	k.isDKGSvcEnabled = true
	k.validatorEVMAddr = testValidatorAddr

	params := types.DefaultParams()
	params.MinReqRegisteredParticipants = 2
	params.OperationalThreshold = 667
	require.NoError(t, k.SetParams(ctx, params))

	round := uint32(1)
	latestRound := &types.DKGNetwork{
		Round:        round,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{testValidatorAddr, "0xother"},
	}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))

	for i, addrHex := range []string{testValidatorAddr, "0x2222222222222222222222222222222222222222"} {
		addr := common.HexToAddress(addrHex)
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:     round,
			Index:     uint32(i + 1),
			DkgPubKey: []byte("pub"),
			Status:    types.DKGRegStatusVerified,
		}))
	}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	err = k.BeginDealing(ctx, latestRound)
	require.NoError(t, err)
}
