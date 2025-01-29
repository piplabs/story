package v1

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name for the Story singularity v1 upgrade.
	UpgradeName = "singularity_v1"

	// UpgradeHeight defines the block height at which the Story singularity v1 upgrade is triggered.
	UpgradeHeight = 1_000_000
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeInfo:    "singularity upgrade to change singularity height and the duration of the short staking period",
	UpgradeHeight:  UpgradeHeight,
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}
