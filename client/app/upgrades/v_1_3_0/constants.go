//nolint:revive,stylecheck // versioning
package v_1_3_0

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
)

const (
	UpgradeName = "v1.3.0"

	// AeneidUpgradeHeight defines the block height at which v1.3.0 upgrade is triggered on Aeneid.
	AeneidUpgradeHeight = 10000000
	// StoryUpgradeHeight defines the block height at which v1.3.0 upgrade is triggered on Mainnet.
	StoryUpgradeHeight = 10000000

	// DevnetUpgradeHeight defines the block height at which virgil upgrade is triggered on Internal Devnet.
	DevnetUpgradeHeight = 5000

	NewMaxValidators = 80
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName: UpgradeName,
	UpgradeInfo: "upgrade to v1.3.0",
	// UpgradeHeight is set in `scheduleForkUpgrade`
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}

func GetUpgradeHeight(chainID string) (int64, bool) {
	switch chainID {
	case upgrades.AeneidChainID:
		return AeneidUpgradeHeight, true
	case upgrades.StoryChainID:
		return StoryUpgradeHeight, true
	case upgrades.DevnetChainID:
		return DevnetUpgradeHeight, true
	default:
		return 0, false
	}
}
