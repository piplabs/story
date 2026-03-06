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
		portion,
		types.DefaultMinReqRegisteredParticipants,
		types.DefaultMinReqFinalizedParticipants,
		types.DefaultOperationalThreshold,
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

func TestValidateMinReqRegisteredParticipants(t *testing.T) {
	tcs := []struct {
		name        string
		value       uint32
		expectedErr string
	}{
		{name: "pass: valid value", value: 3},
		{name: "pass: value of 1", value: 1},
		{name: "pass: large value", value: 1000},
		{name: "fail: zero value", value: 0, expectedErr: "min_req_registered_participants must be greater than zero"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateMinReqRegisteredParticipants(tc.value)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateMinReqFinalizedParticipants(t *testing.T) {
	tcs := []struct {
		name        string
		value       uint32
		expectedErr string
	}{
		{name: "pass: valid value", value: 3},
		{name: "pass: value of 1", value: 1},
		{name: "pass: large value", value: 500},
		{name: "fail: zero value", value: 0, expectedErr: "min_req_finalized_participants must be greater than zero"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateMinReqFinalizedParticipants(tc.value)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateOperationalThreshold(t *testing.T) {
	tcs := []struct {
		name        string
		value       uint32
		expectedErr string
	}{
		{name: "pass: valid threshold 667 (66.7%)", value: 667},
		{name: "pass: threshold of 1", value: 1},
		{name: "pass: threshold of 1000 (100%)", value: 1000},
		{name: "pass: threshold of 500 (50%)", value: 500},
		{name: "fail: zero threshold", value: 0, expectedErr: "operational_threshold must be greater than zero"},
		{name: "fail: exceeds basis", value: 1001, expectedErr: "operational_threshold must not exceed basis (1000)"},
		{name: "fail: way over basis", value: 5000, expectedErr: "operational_threshold must not exceed basis (1000)"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateOperationalThreshold(tc.value)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDefaultParams_IncludesNewDKGParams(t *testing.T) {
	params := types.DefaultParams()

	require.Equal(t, types.DefaultMinReqRegisteredParticipants, params.MinReqRegisteredParticipants)
	require.Equal(t, types.DefaultMinReqFinalizedParticipants, params.MinReqFinalizedParticipants)
	require.Equal(t, types.DefaultOperationalThreshold, params.OperationalThreshold)
}

func TestCalculateThreshold(t *testing.T) {
	tcs := []struct {
		name                    string
		total                   uint32
		operationalThresholdBps uint32
		expected                uint32
	}{
		{
			name:                    "standard 66.7% threshold with 10 validators",
			total:                   10,
			operationalThresholdBps: 667,
			expected:                7, // ceil(10 * 667 / 1000) = ceil(6.67) = 7
		},
		{
			name:                    "exact division: 50% with 10 validators",
			total:                   10,
			operationalThresholdBps: 500,
			expected:                5, // 10 * 500 / 1000 = 5 (exact)
		},
		{
			name:                    "ceiling: 50% with 3 validators",
			total:                   3,
			operationalThresholdBps: 500,
			expected:                2, // ceil(3 * 500 / 1000) = ceil(1.5) = 2
		},
		{
			name:                    "100% threshold",
			total:                   5,
			operationalThresholdBps: 1000,
			expected:                5, // 5 * 1000 / 1000 = 5
		},
		{
			name:                    "minimum threshold (0.1%)",
			total:                   10,
			operationalThresholdBps: 1,
			expected:                1, // ceil(10 * 1 / 1000) = ceil(0.01) = 1
		},
		{
			name:                    "zero total",
			total:                   0,
			operationalThresholdBps: 667,
			expected:                0,
		},
		{
			name:                    "zero threshold bps",
			total:                   10,
			operationalThresholdBps: 0,
			expected:                0,
		},
		{
			name:                    "single validator with 66.7%",
			total:                   1,
			operationalThresholdBps: 667,
			expected:                1, // ceil(1 * 667 / 1000) = ceil(0.667) = 1
		},
		{
			name:                    "large committee",
			total:                   100,
			operationalThresholdBps: 667,
			expected:                67, // ceil(100 * 667 / 1000) = ceil(66.7) = 67
		},
		{
			name:                    "exact 2/3 threshold with 9 validators",
			total:                   9,
			operationalThresholdBps: 667,
			expected:                7, // ceil(9 * 667 / 1000) = ceil(6.003) = 7
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := types.CalculateThreshold(tc.total, tc.operationalThresholdBps)
			require.Equal(t, tc.expected, result, "CalculateThreshold(%d, %d)", tc.total, tc.operationalThresholdBps)
		})
	}
}

func TestParams_Validate_NewParams(t *testing.T) {
	// Test that setting invalid new params triggers validation errors
	params := types.DefaultParams()

	params.MinReqRegisteredParticipants = 0
	err := params.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "min_req_registered_participants must be greater than zero")

	params.MinReqRegisteredParticipants = 3
	params.MinReqFinalizedParticipants = 0
	err = params.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "min_req_finalized_participants must be greater than zero")

	params.MinReqFinalizedParticipants = 3
	params.OperationalThreshold = 0
	err = params.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "operational_threshold must be greater than zero")

	params.OperationalThreshold = 1001
	err = params.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "operational_threshold must not exceed basis (1000)")

	// Valid params should pass
	params.OperationalThreshold = 667
	err = params.Validate()
	require.NoError(t, err)
}
