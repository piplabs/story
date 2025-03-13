//nolint:revive,stylecheck // versioning
package v_1_2_0

import (
	"time"

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
)

const (
	UpgradeName = "v1.2.0"

	// Initial parameters.
	InitialRefundFeeBps = 100                // 1%
	InitialRefundPeriod = 1 * 24 * time.Hour // 1 days
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeInfo:    "v1.2.0 upgrade",
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}
