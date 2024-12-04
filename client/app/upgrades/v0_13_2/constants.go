//nolint:revive,stylecheck // version underscores
package v0_13_2

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/piplabs/story/client/app/upgrades"
)

const UpgradeName = "v0.13.2"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}
