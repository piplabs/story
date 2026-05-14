package v_1_9_0

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
	"github.com/piplabs/story/lib/netconf"
)

// UpgradeName is the on-chain name for the v1.9.0 upgrade that introduces
// round-based pruning for DKGPartialDecrypt entries. No new module stores are
// added; the secondary index lives within the existing DKG KV store prefix.
const UpgradeName = netconf.V190

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeInfo:    "add round-based pruning for DKG partial decrypt submissions",
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}
