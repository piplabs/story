//nolint:revive,stylecheck // versioning
package v_1_2

import (
	"time"

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
)

const (
	UpgradeName = "v1.2"

	// AeneidUpgradeHeight defines the block height at which the upgrade is triggered on Aeneid.
	AeneidUpgradeHeight = 10_000_000
	// StoryUpgradeHeight defines the block height at which the upgrade is triggered on Story.
	StoryUpgradeHeight = 10_000_000

	// Initial parameters.
	InitialRefundFeeBps = 100                // 1%
	InitialRefundPeriod = 7 * 24 * time.Hour // 7 days
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeInfo:    "v1.2 upgrade",
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}

func GetUpgradeHeight(chainID string) (int64, bool) {
	switch chainID {
	case upgrades.AeneidChainID:
		return AeneidUpgradeHeight, true
	case upgrades.StoryChainID:
		return StoryUpgradeHeight, true
	default:
		return 0, false
	}
}
