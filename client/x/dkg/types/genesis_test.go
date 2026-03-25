package types_test

import (
	"testing"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

func TestNewGenesisState(t *testing.T) {
	t.Parallel()

	params := types.NewParams(
		100, 200, 300, 400,
		math.LegacyMustNewDecFromStr("0.15"),
		5, 4, 750, 30,
	)

	gs := types.NewGenesisState(params)
	require.NotNil(t, gs)
	require.Equal(t, params.RegistrationPeriod, gs.Params.RegistrationPeriod)
	require.Equal(t, params.DealingPeriod, gs.Params.DealingPeriod)
	require.Equal(t, params.FinalizationPeriod, gs.Params.FinalizationPeriod)
	require.Equal(t, params.ActivePeriod, gs.Params.ActivePeriod)
	require.True(t, params.DkgCommitteeRewardPortion.Equal(gs.Params.DkgCommitteeRewardPortion))
	require.Equal(t, params.MinReqRegisteredParticipants, gs.Params.MinReqRegisteredParticipants)
	require.Equal(t, params.MinReqFinalizedParticipants, gs.Params.MinReqFinalizedParticipants)
	require.Equal(t, params.OperationalThreshold, gs.Params.OperationalThreshold)
}

func TestDefaultGenesisState(t *testing.T) {
	t.Parallel()

	gs := types.DefaultGenesisState()
	require.NotNil(t, gs)

	// Verify all default values match the expected constants.
	require.Equal(t, types.DefaultDkgRegistrationPeriod, gs.Params.RegistrationPeriod)
	require.Equal(t, types.DefaultDkgDealingPeriod, gs.Params.DealingPeriod)
	require.Equal(t, types.DefaultDkgFinalizationPeriod, gs.Params.FinalizationPeriod)
	require.Equal(t, types.DefaultDkgActivePeriod, gs.Params.ActivePeriod)
	require.True(t, types.DefaultDkgCommitteeRewardPortion.Equal(gs.Params.DkgCommitteeRewardPortion))
	require.Equal(t, types.DefaultMinReqRegisteredParticipants, gs.Params.MinReqRegisteredParticipants)
	require.Equal(t, types.DefaultMinReqFinalizedParticipants, gs.Params.MinReqFinalizedParticipants)
	require.Equal(t, types.DefaultOperationalThreshold, gs.Params.OperationalThreshold)

	// Default genesis should pass validation.
	require.NoError(t, gs.Params.Validate())
}
