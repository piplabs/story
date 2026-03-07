package v_2_0_0

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/netconf"
)

const UpgradeName = netconf.V200

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{dkgtypes.StoreKey},
	},
}

var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeInfo:    "activate DKG module on the live network",
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}

func GetUpgradeHeight(ctx sdk.Context) (int64, bool) {
	height, err := netconf.GetUpgradeHeight(ctx.ChainID(), UpgradeName)
	if err != nil {
		log.Error(ctx, "Failed to get upgrade height", err, "chain_id", ctx.ChainID(), "upgrade_name", UpgradeName)

		return 0, false
	}

	return height, true
}
