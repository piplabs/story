package virgil

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name for the virgil upgrade.
	UpgradeName = "virgil"

	// AeneidUpgradeHeight defines the block height at which virgil upgrade is triggered on Aeneid.
	AeneidUpgradeHeight = 345158
	// StoryUpgradeHeight defines the block height at which virgil upgrade is triggered on Story.
	StoryUpgradeHeight = 677886
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName: UpgradeName,
	UpgradeInfo: "upgrade to change the duration of the short staking period during the singularity period",
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
