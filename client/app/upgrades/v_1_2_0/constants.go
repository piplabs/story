//nolint:revive,stylecheck // versioning
package v_1_2_0

import (
	"time"

	storetypes "cosmossdk.io/store/types"

	"github.com/piplabs/story/client/app/upgrades"
)

const (
	UpgradeName = "v1.2.0"

	// Initial parameters.
	InitialRefundFeeBps uint32        = 100            // 1%
	InitialRefundPeriod time.Duration = 24 * time.Hour // 1 days
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}
