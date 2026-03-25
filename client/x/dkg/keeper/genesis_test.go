package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// TestKeeper_InitGenesis_ValidState verifies that InitGenesis successfully
// stores valid params into the keeper.
func TestKeeper_InitGenesis_ValidState(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	gs := &types.GenesisState{
		Params: types.NewParams(
			10, 20, 30, 40,
			types.DefaultDkgCommitteeRewardPortion,
			5, 4, 700, 30,
		),
	}

	err := k.InitGenesis(ctx, gs)
	require.NoError(t, err)

	got, err := k.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, gs.Params.RegistrationPeriod, got.RegistrationPeriod)
	require.Equal(t, gs.Params.DealingPeriod, got.DealingPeriod)
	require.Equal(t, gs.Params.FinalizationPeriod, got.FinalizationPeriod)
	require.Equal(t, gs.Params.ActivePeriod, got.ActivePeriod)
	require.Equal(t, gs.Params.MinReqRegisteredParticipants, got.MinReqRegisteredParticipants)
	require.Equal(t, gs.Params.MinReqFinalizedParticipants, got.MinReqFinalizedParticipants)
	require.Equal(t, gs.Params.OperationalThreshold, got.OperationalThreshold)
}

// TestKeeper_InitGenesis_InvalidParams verifies that InitGenesis returns an
// error when the genesis state contains invalid params.
func TestKeeper_InitGenesis_InvalidParams(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name        string
		params      types.Params
		expectedErr string
	}{
		{
			name: "fail: zero registration period",
			params: types.NewParams(
				0, // RegistrationPeriod = 0 is invalid
				20, 30, 40,
				types.DefaultDkgCommitteeRewardPortion,
				3, 3, 667, 30,
			),
			expectedErr: "invalid dkg registration period",
		},
		{
			name: "fail: zero dealing period",
			params: types.NewParams(
				10, 0, // DealingPeriod = 0 is invalid
				30, 40,
				types.DefaultDkgCommitteeRewardPortion,
				3, 3, 667, 30,
			),
			expectedErr: "invalid dkg dealing period",
		},
		{
			name: "fail: zero finalization period",
			params: types.NewParams(
				10, 20, 0, // FinalizationPeriod = 0 is invalid
				40,
				types.DefaultDkgCommitteeRewardPortion,
				3, 3, 667, 30,
			),
			expectedErr: "invalid dkg finalization period",
		},
		{
			name: "fail: zero min_req_registered_participants",
			params: types.NewParams(
				10, 20, 30, 40,
				types.DefaultDkgCommitteeRewardPortion,
				0, // MinReqRegisteredParticipants = 0 is invalid
				3, 667, 30,
			),
			expectedErr: "must be greater than zero",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			k, _, _, ctx := setupDKGKeeperWithMocks(t)

			gs := &types.GenesisState{Params: tc.params}
			err := k.InitGenesis(ctx, gs)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

// TestKeeper_ExportGenesis_ReturnsStoredParams verifies that ExportGenesis
// returns the params that were previously stored via InitGenesis.
func TestKeeper_ExportGenesis_ReturnsStoredParams(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	params := types.NewParams(
		50, 60, 70, 80,
		types.DefaultDkgCommitteeRewardPortion,
		7, 6, 750, 30,
	)

	gs := &types.GenesisState{Params: params}
	require.NoError(t, k.InitGenesis(ctx, gs))

	exported := k.ExportGenesis(sdkCtx)
	require.NotNil(t, exported)
	require.Equal(t, params.RegistrationPeriod, exported.Params.RegistrationPeriod)
	require.Equal(t, params.DealingPeriod, exported.Params.DealingPeriod)
	require.Equal(t, params.MinReqRegisteredParticipants, exported.Params.MinReqRegisteredParticipants)
	require.Equal(t, params.OperationalThreshold, exported.Params.OperationalThreshold)
}

// TestKeeper_ExportGenesis_DefaultParams verifies that ExportGenesis returns
// the default params when no custom params have been set.
func TestKeeper_ExportGenesis_DefaultParams(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// setupDKGKeeperWithMocks already sets DefaultParams
	exported := k.ExportGenesis(sdkCtx)
	require.NotNil(t, exported)
	require.Equal(t, types.DefaultParams().RegistrationPeriod, exported.Params.RegistrationPeriod)
}

// TestKeeper_ValidateGenesis_ValidState verifies that ValidateGenesis passes
// for a valid genesis state.
func TestKeeper_ValidateGenesis_ValidState(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	gs := &types.GenesisState{
		Params: types.DefaultParams(),
	}
	err := k.ValidateGenesis(gs)
	require.NoError(t, err)
}

// TestKeeper_ValidateGenesis_InvalidParams verifies that ValidateGenesis
// returns an error for invalid params.
func TestKeeper_ValidateGenesis_InvalidParams(t *testing.T) {
	t.Parallel()

	k, _, _, _ := setupDKGKeeperWithMocks(t)

	gs := &types.GenesisState{
		Params: types.NewParams(
			0, // RegistrationPeriod = 0 is invalid
			20, 30, 40,
			types.DefaultDkgCommitteeRewardPortion,
			3, 3, 667, 30,
		),
	}
	err := k.ValidateGenesis(gs)
	require.Error(t, err)
}

// TestKeeper_InitGenesis_DefaultParams verifies that InitGenesis with default
// params succeeds and round-trips correctly.
func TestKeeper_InitGenesis_DefaultParams(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	gs := types.DefaultGenesisState()
	require.NoError(t, k.InitGenesis(ctx, gs))

	exported := k.ExportGenesis(sdkCtx)
	require.Equal(t, gs.Params.RegistrationPeriod, exported.Params.RegistrationPeriod)
	require.Equal(t, gs.Params.ActivePeriod, exported.Params.ActivePeriod)
}
