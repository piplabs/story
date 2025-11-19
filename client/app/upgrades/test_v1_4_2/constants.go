package test_v1_4_2

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/netconf"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          netconf.TestV142,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName:    netconf.TestV142,
	UpgradeInfo:    "test-v1.4.2 upgrade that lowers min stake and unstake amounts",
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}

func GetUpgradeHeight(ctx sdk.Context) (int64, bool) {
	height, err := netconf.GetUpgradeHeight(ctx.ChainID(), netconf.TestV142)
	if err != nil {
		log.Error(ctx, "Failed to get upgrade height", err, "chain_id", ctx.ChainID(), "upgrade_name", netconf.TestV142)

		return 0, false
	}

	return height, true
}
