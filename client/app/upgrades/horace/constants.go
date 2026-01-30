package horace

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/piplabs/story/lib/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
	"github.com/piplabs/story/lib/netconf"
)

var (
	NewAnnualInflationsPerYear        = math.LegacyNewDec(15_315_000_000_000_000) // 15.315M IP
	NewBlocksPerYear           uint64 = 13_140_000                                // 13.14M blocks
	NewLockedTokenMultiplier          = math.LegacyNewDecWithPrec(25, 3)          // 0.025
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          netconf.Horace,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName:    netconf.Horace,
	UpgradeInfo:    "emissions and locked staking multiplier adjustment",
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}

func GetUpgradeHeight(ctx sdk.Context) (int64, bool) {
	height, err := netconf.GetUpgradeHeight(ctx.ChainID(), netconf.Horace)
	if err != nil {
		log.Error(ctx, "Failed to get upgrade height", err, "chain_id", ctx.ChainID(), "upgrade_name", netconf.Horace)

		return 0, false
	}

	return height, true
}
