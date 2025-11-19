package test_v1_4_3

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/netconf"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          netconf.TestV143,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName:    netconf.TestV143,
	UpgradeInfo:    "test-v1.4.3 upgrade that lowers min delegation to 1 IP",
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}

func GetUpgradeHeight(ctx sdk.Context) (int64, bool) {
	height, err := netconf.GetUpgradeHeight(ctx.ChainID(), netconf.TestV143)
	if err != nil {
		log.Error(ctx, "Failed to get upgrade height", err, "chain_id", ctx.ChainID(), "upgrade_name", netconf.TestV143)

		return 0, false
	}

	return height, true
}
