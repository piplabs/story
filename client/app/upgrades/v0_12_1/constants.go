//nolint:revive,stylecheck // version underscores
package v0_12_1

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/piplabs/story/client/app/upgrades"
)

const UpgradeName = "v0.12.1"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}
