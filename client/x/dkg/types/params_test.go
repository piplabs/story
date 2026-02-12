package types_test

import (
	"testing"

	"cosmossdk.io/math"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/stretchr/testify/require"
)

func TestValidateDkgCommitteeRewardPortion(t *testing.T) {
	tcs := []struct {
		name        string
		portion     math.LegacyDec
		expectedErr string
	}{
		{
			name:    "pass: zero portion",
			portion: math.LegacyZeroDec(),
		},
		{
			name:    "pass: 10% portion",
			portion: math.LegacyMustNewDecFromStr("0.10"),
		},
		{
			name:    "pass: 100% portion",
			portion: math.LegacyOneDec(),
		},
		{
			name:    "pass: 50% portion",
			portion: math.LegacyMustNewDecFromStr("0.50"),
		},
		{
			name:    "pass: very small portion",
			portion: math.LegacyMustNewDecFromStr("0.000000000000000001"),
		},
		{
			name:        "fail: negative portion",
			portion:     math.LegacyMustNewDecFromStr("-0.01"),
			expectedErr: "dkg committee reward portion must not be negative",
		},
		{
			name:        "fail: portion exceeds 1.0",
			portion:     math.LegacyMustNewDecFromStr("1.01"),
			expectedErr: "dkg committee reward portion must not exceed 1.0",
		},
		{
			name:        "fail: large negative portion",
			portion:     math.LegacyMustNewDecFromStr("-100.0"),
			expectedErr: "dkg committee reward portion must not be negative",
		},
		{
			name:        "fail: way over 1.0",
			portion:     math.LegacyMustNewDecFromStr("2.0"),
			expectedErr: "dkg committee reward portion must not exceed 1.0",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateDkgCommitteeRewardPortion(tc.portion)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDefaultParams_IncludesDkgCommitteeRewardPortion(t *testing.T) {
	params := types.DefaultParams()

	// Verify the default reward portion is 10%.
	require.True(t, params.DkgCommitteeRewardPortion.Equal(math.LegacyMustNewDecFromStr("0.10")),
		"default DKG committee reward portion should be 0.10 (10%%)")

	// Note: DefaultParams().Validate() will fail because DefaultParams does not
	// include a valid CodeCommitment (which is required to be 32 bytes). This is
	// expected behavior since code commitment is set separately via governance.
	// We only validate the reward portion field here.
	require.NoError(t, types.ValidateDkgCommitteeRewardPortion(params.DkgCommitteeRewardPortion))
}

func TestNewParams_WithDkgCommitteeRewardPortion(t *testing.T) {
	portion := math.LegacyMustNewDecFromStr("0.25")
	params := types.NewParams(
		types.DefaultDkgRegistrationPeriod,
		types.DefaultDkgDealingPeriod,
		types.DefaultDkgFinalizationPeriod,
		types.DefaultDkgActivePeriod,
		types.DefaultDkgComplaintPeriod,
		types.DefaultMinCommitteeSize,
		portion,
	)

	require.True(t, params.DkgCommitteeRewardPortion.Equal(portion))
	// Note: Validate will fail because CodeCommitment is empty in this test.
	// That's expected since we're only testing the reward portion field.
}

func TestParams_Validate_InvalidRewardPortion(t *testing.T) {
	params := types.DefaultParams()
	params.DkgCommitteeRewardPortion = math.LegacyMustNewDecFromStr("-0.5")

	err := params.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "dkg committee reward portion must not be negative")
}
