package netconf_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/lib/netconf"
)

func TestGetUpgradeHeight(t *testing.T) {
	tcs := []struct {
		name           string
		chainID        string
		upgradeName    string
		expectedErr    string
		expectedResult int64
	}{
		{
			name:        "unknown chain ID",
			chainID:     "unknown-chain-id",
			upgradeName: netconf.Terence,
			expectedErr: netconf.ErrUnknownChainID.Error(),
		},
		{
			name:        "unknown upgrade name",
			chainID:     netconf.TestChainID,
			upgradeName: "unknown-upgrade",
			expectedErr: netconf.ErrUnknownUpgrade.Error(),
		},
		{
			name:           "known chain ID and upgrade name",
			chainID:        netconf.TestChainID,
			upgradeName:    netconf.Terence,
			expectedResult: 50,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			upgradeHeight, err := netconf.GetUpgradeHeight(tc.chainID, tc.upgradeName)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.Equal(t, tc.expectedResult, upgradeHeight)
			}
		})
	}
}

func TestIsSeneca(t *testing.T) {
	// TestChainID has Seneca registered at height 300.
	tcs := []struct {
		name        string
		chainID     string
		blockNumber int64
		expected    bool
		expectErr   bool
	}{
		{
			name:        "before seneca",
			chainID:     netconf.TestChainID,
			blockNumber: 299,
			expected:    false,
		},
		{
			name:        "at seneca",
			chainID:     netconf.TestChainID,
			blockNumber: 300,
			expected:    true,
		},
		{
			name:        "after seneca",
			chainID:     netconf.TestChainID,
			blockNumber: 1_000_000,
			expected:    true,
		},
		{
			name:        "unknown chain id",
			chainID:     "unknown-chain-id",
			blockNumber: 1,
			expectErr:   true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got, err := netconf.IsSeneca(tc.chainID, tc.blockNumber)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, got)
		})
	}
}
