package terence

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/netconf"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          netconf.Terence,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName:    netconf.Terence,
	UpgradeInfo:    "Terence upgrade that supports rolling upgrade via upgrade entrypoint contract with single binary",
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}

func GetUpgradeHeight(ctx sdk.Context) (int64, bool) {
	height, err := netconf.GetUpgradeHeight(ctx.ChainID(), netconf.Terence)
	if err != nil {
		log.Error(ctx, "Failed to get upgrade height", err, "chain_id", ctx.ChainID(), "upgrade_name", netconf.Terence)

		return 0, false
	}

	return height, true
}
