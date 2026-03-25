package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// TestKeeper_SetGetParams verifies that params can be stored and retrieved
// via SetParams / GetParams.
func TestKeeper_SetGetParams(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	params := types.NewParams(
		10, // RegistrationPeriod
		20, // DealingPeriod
		30, // FinalizationPeriod
		40, // ActivePeriod
		types.DefaultDkgCommitteeRewardPortion,
		5, // MinReqRegisteredParticipants
		4, // MinReqFinalizedParticipants
		types.DefaultOperationalThreshold,
		types.DefaultDecryptTimeout,
	)

	require.NoError(t, k.SetParams(ctx, params))

	got, err := k.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, params.RegistrationPeriod, got.RegistrationPeriod)
	require.Equal(t, params.DealingPeriod, got.DealingPeriod)
	require.Equal(t, params.FinalizationPeriod, got.FinalizationPeriod)
	require.Equal(t, params.ActivePeriod, got.ActivePeriod)
	require.Equal(t, params.MinReqRegisteredParticipants, got.MinReqRegisteredParticipants)
	require.Equal(t, params.MinReqFinalizedParticipants, got.MinReqFinalizedParticipants)
	require.Equal(t, params.OperationalThreshold, got.OperationalThreshold)
}

// TestKeeper_GetParams_DefaultAfterInit verifies that GetParams returns the
// default params that were set during keeper initialization.
func TestKeeper_GetParams_DefaultAfterInit(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	got, err := k.GetParams(ctx)
	require.NoError(t, err)
	// setupDKGKeeperWithMocks sets DefaultParams
	require.Equal(t, types.DefaultParams().RegistrationPeriod, got.RegistrationPeriod)
	require.Equal(t, types.DefaultParams().DealingPeriod, got.DealingPeriod)
}

// TestKeeper_SetMinReqRegisteredParticipants verifies setting individual
// params field via SetMinReqRegisteredParticipants.
func TestKeeper_SetMinReqRegisteredParticipants(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name        string
		value       uint32
		expectedErr string
	}{
		{
			name:  "pass: valid value",
			value: 5,
		},
		{
			name:        "fail: zero value",
			value:       0,
			expectedErr: "must be greater than zero",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			k, _, _, ctx := setupDKGKeeperWithMocks(t)

			err := k.SetMinReqRegisteredParticipants(ctx, tc.value)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				got, err := k.GetParams(ctx)
				require.NoError(t, err)
				require.Equal(t, tc.value, got.MinReqRegisteredParticipants)
			}
		})
	}
}

// TestKeeper_SetMinReqFinalizedParticipants verifies setting individual
// params field via SetMinReqFinalizedParticipants.
func TestKeeper_SetMinReqFinalizedParticipants(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name        string
		value       uint32
		expectedErr string
	}{
		{
			name:  "pass: valid value",
			value: 3,
		},
		{
			name:        "fail: zero value",
			value:       0,
			expectedErr: "must be greater than zero",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			k, _, _, ctx := setupDKGKeeperWithMocks(t)

			err := k.SetMinReqFinalizedParticipants(ctx, tc.value)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				got, err := k.GetParams(ctx)
				require.NoError(t, err)
				require.Equal(t, tc.value, got.MinReqFinalizedParticipants)
			}
		})
	}
}

// TestKeeper_SetOperationalThreshold verifies setting individual
// params field via SetOperationalThreshold.
func TestKeeper_SetOperationalThreshold(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name        string
		value       uint32
		expectedErr string
	}{
		{
			name:  "pass: valid threshold 667 (66.7%)",
			value: 667,
		},
		{
			name:  "pass: threshold at basis (1000 = 100%)",
			value: 1000,
		},
		{
			name:        "fail: zero value",
			value:       0,
			expectedErr: "must be greater than zero",
		},
		{
			name:        "fail: value exceeds basis 1001",
			value:       1001,
			expectedErr: "must not exceed basis",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			k, _, _, ctx := setupDKGKeeperWithMocks(t)

			err := k.SetOperationalThreshold(ctx, tc.value)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				got, err := k.GetParams(ctx)
				require.NoError(t, err)
				require.Equal(t, tc.value, got.OperationalThreshold)
			}
		})
	}
}

// TestKeeper_SetParams_OverwritesPreviousParams verifies that calling SetParams
// twice correctly replaces the first set of params.
func TestKeeper_SetParams_OverwritesPreviousParams(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	first := types.NewParams(
		10, 20, 30, 40,
		types.DefaultDkgCommitteeRewardPortion,
		3, 3, 667, 200,
	)
	second := types.NewParams(
		100, 200, 300, 400,
		types.DefaultDkgCommitteeRewardPortion,
		10, 9, 800, 200,
	)

	require.NoError(t, k.SetParams(ctx, first))
	require.NoError(t, k.SetParams(ctx, second))

	got, err := k.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, second.RegistrationPeriod, got.RegistrationPeriod)
	require.Equal(t, second.MinReqRegisteredParticipants, got.MinReqRegisteredParticipants)
}
