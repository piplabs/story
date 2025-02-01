package v1

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name for the singularity v1 upgrade.
	UpgradeName = "singularity_v1"

	// AeneidUpgradeHeight defines the block height at which the Aeneid singularity v1 upgrade is triggered.
	AeneidUpgradeHeight = 345158
	// StoryUpgradeHeight defines the block height at which the Story singularity v1 upgrade is triggered.
	StoryUpgradeHeight = 677886
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName: UpgradeName,
	UpgradeInfo: "singularity upgrade to change the duration of the short staking period",
	// UpgradeHeight is set in `scheduleForkUpgrade`
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}

func GetUpgradeHeight(chainID string) (int64, bool) {
	switch chainID {
	case upgrades.AeneidChainID:
		return AeneidUpgradeHeight, true
	case upgrades.StoryChainID:
		return StoryUpgradeHeight, true
	default:
		return 0, false
	}
}
