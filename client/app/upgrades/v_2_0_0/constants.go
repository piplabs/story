package v_2_0_0

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/netconf"
)

// UpgradeName is the on-chain name for the v2.0.0 upgrade that activates the
// DKG module. This is a binary-swap upgrade: the old binary should halt at
// the scheduled height, and operators replace it with the v2.0.0 binary.
//
// the upgrade is scheduled on-chain by calling UpgradeEntrypoint.planUpgrade("v2.0.0", height).
// The old binary writes upgrade-info.json to disk before halting, which the new binary reads
// to configure the store loader (see setupUpgradeStoreLoaders in upgrades.go).
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
	UpgradeInfo:    "activate DKG module on the network",
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}
