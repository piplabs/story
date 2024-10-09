//nolint:revive,stylecheck // version underscores
package v0_10_0

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/piplabs/story/client/app/upgrades"
)

const UpgradeName = "v0.10.0"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}
