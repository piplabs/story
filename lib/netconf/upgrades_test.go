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
			upgradeName: netconf.V140,
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
			upgradeName:    netconf.V140,
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
