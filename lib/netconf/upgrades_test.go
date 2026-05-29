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

func TestIsV190(t *testing.T) {
	tcs := []struct {
		name        string
		chainID     string
		blockNumber int64
		expectedErr string
		expected    bool
	}{
		{
			name:        "before upgrade height",
			chainID:     netconf.TestChainID,
			blockNumber: 399,
			expected:    false,
		},
		{
			name:        "at upgrade height",
			chainID:     netconf.TestChainID,
			blockNumber: 400,
			expected:    true,
		},
		{
			name:        "after upgrade height",
			chainID:     netconf.TestChainID,
			blockNumber: 401,
			expected:    true,
		},
		{
			name:        "localnet active from genesis",
			chainID:     netconf.StoryLocalnetID,
			blockNumber: 0,
			expected:    true,
		},
		{
			name:        "unknown chain ID returns error",
			chainID:     "unknown-chain-id",
			blockNumber: 1000,
			expectedErr: netconf.ErrUnknownChainID.Error(),
		},
		{
			name:        "chain without V190 registered returns error",
			chainID:     netconf.StoryChainID,
			blockNumber: 1000,
			expectedErr: netconf.ErrUnknownUpgrade.Error(),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := netconf.IsV190(tc.chainID, tc.blockNumber)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedErr)
				require.False(t, ok)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, ok)
			}
		})
	}
}
