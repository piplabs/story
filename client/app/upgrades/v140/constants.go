package v140

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
	"github.com/piplabs/story/lib/netconf"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          netconf.V140,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName:    netconf.V140,
	UpgradeInfo:    "v140 upgrade",
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
	UpgradeHeight:  100000000, // TODO: set fork height depend on the network
}
