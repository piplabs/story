package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

const testValidatorAddr = "0xabcdef1234567890abcdef1234567890abcdef12"

// --- shouldDeal ---

func TestShouldDeal(t *testing.T) {
	tests := []struct {
		name             string
		validatorAddr    string
		activeValSet     []string
		prevActive       *types.DKGNetwork // nil = no previous active round
		expectedShouldDl bool
	}{
		{
			name:             "first round: validator in current set should deal",
			validatorAddr:    testValidatorAddr,
			activeValSet:     []string{testValidatorAddr, "0xother"},
			prevActive:       nil,
			expectedShouldDl: true,
		},
		{
			name:             "first round: validator not in current set should not deal",
			validatorAddr:    testValidatorAddr,
			activeValSet:     []string{"0xother1", "0xother2"},
			prevActive:       nil,
			expectedShouldDl: false,
		},
		{
			name:          "resharing: validator in previous set should deal",
			validatorAddr: testValidatorAddr,
			activeValSet:  []string{"0xnewval1", "0xnewval2"},
			prevActive: &types.DKGNetwork{
				Round:        1,
				Stage:        types.DKGStageActive,
				ActiveValSet: []string{testValidatorAddr, "0xoldval"},
			},
			expectedShouldDl: true,
		},
		{
			name:          "resharing: validator not in previous set should not deal",
			validatorAddr: testValidatorAddr,
			activeValSet:  []string{testValidatorAddr, "0xnewval"},
			prevActive: &types.DKGNetwork{
				Round:        1,
				Stage:        types.DKGStageActive,
				ActiveValSet: []string{"0xoldval1", "0xoldval2"},
			},
			expectedShouldDl: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			k, _, _, ctx := setupDKGKeeperWithMocks(t)
			k.validatorEVMAddr = tc.validatorAddr

			// Set up previous active round if present
			if tc.prevActive != nil {
				require.NoError(t, k.setDKGNetwork(ctx, tc.prevActive))
				require.NoError(t, k.setLatestActiveRound(ctx, tc.prevActive))
			}

			dkgNetwork := &types.DKGNetwork{
				Round:        2,
				Stage:        types.DKGStageDealing,
				ActiveValSet: tc.activeValSet,
			}

			shouldDl, err := k.shouldDeal(ctx, dkgNetwork)
			require.NoError(t, err)
			require.Equal(t, tc.expectedShouldDl, shouldDl)
		})
	}
}

// --- shouldProcessResponses ---

func TestShouldProcessResponses(t *testing.T) {
	tests := []struct {
		name           string
		validatorAddr  string
		activeValSet   []string
		prevActive     *types.DKGNetwork
		expectedResult bool
	}{
		{
			name:           "first round: validator in current set should process",
			validatorAddr:  testValidatorAddr,
			activeValSet:   []string{testValidatorAddr, "0xother"},
			prevActive:     nil,
			expectedResult: true,
		},
		{
			name:           "first round: validator not in current set should not process",
			validatorAddr:  testValidatorAddr,
			activeValSet:   []string{"0xother1", "0xother2"},
			prevActive:     nil,
			expectedResult: false,
		},
		{
			name:          "resharing: validator in current set only",
			validatorAddr: testValidatorAddr,
			activeValSet:  []string{testValidatorAddr, "0xnewval"},
			prevActive: &types.DKGNetwork{
				Round:        1,
				Stage:        types.DKGStageActive,
				ActiveValSet: []string{"0xoldval1", "0xoldval2"},
			},
			expectedResult: true,
		},
		{
			name:          "resharing: validator in previous set only",
			validatorAddr: testValidatorAddr,
			activeValSet:  []string{"0xnewval1", "0xnewval2"},
			prevActive: &types.DKGNetwork{
				Round:        1,
				Stage:        types.DKGStageActive,
				ActiveValSet: []string{testValidatorAddr, "0xoldval"},
			},
			expectedResult: true,
		},
		{
			name:          "resharing: validator in both sets",
			validatorAddr: testValidatorAddr,
			activeValSet:  []string{testValidatorAddr, "0xnewval"},
			prevActive: &types.DKGNetwork{
				Round:        1,
				Stage:        types.DKGStageActive,
				ActiveValSet: []string{testValidatorAddr, "0xoldval"},
			},
			expectedResult: true,
		},
		{
			name:          "resharing: validator in neither set",
			validatorAddr: testValidatorAddr,
			activeValSet:  []string{"0xnewval1", "0xnewval2"},
			prevActive: &types.DKGNetwork{
				Round:        1,
				Stage:        types.DKGStageActive,
				ActiveValSet: []string{"0xoldval1", "0xoldval2"},
			},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			k, _, _, ctx := setupDKGKeeperWithMocks(t)
			k.validatorEVMAddr = tc.validatorAddr

			if tc.prevActive != nil {
				require.NoError(t, k.setDKGNetwork(ctx, tc.prevActive))
				require.NoError(t, k.setLatestActiveRound(ctx, tc.prevActive))
			}

			dkgNetwork := &types.DKGNetwork{
				Round:        2,
				Stage:        types.DKGStageDealing,
				ActiveValSet: tc.activeValSet,
			}

			result, err := k.shouldProcessResponses(ctx, dkgNetwork)
			require.NoError(t, err)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}
