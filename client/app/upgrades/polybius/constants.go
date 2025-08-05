package polybius

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
)

const (
	UpgradeName = "polybius"

	// AeneidUpgradeHeight defines the block height at which Polybius upgrade is triggered on Aeneid.
	AeneidUpgradeHeight = 6008000

	// StoryUpgradeHeight defines the block height at which Polybius upgrade is triggered on Story mainnet.
	StoryUpgradeHeight = 8270000

	// DevnetUpgradeHeight defines the block height at which Polybius upgrade is triggered on internal-devnet.
	DevnetUpgradeHeight = 4300

	NewMaxValidators = 80
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

var Fork = upgrades.Fork{
	UpgradeName: UpgradeName,
	UpgradeInfo: "upgrade to increase max number of validator to 80",
	// UpgradeHeight is set in `scheduleForkUpgrade`
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}

func GetUpgradeHeight(chainID string) (int64, bool) {
	switch chainID {
	case upgrades.AeneidChainID:
		return AeneidUpgradeHeight, true
	case upgrades.StoryChainID:
		return StoryUpgradeHeight, true
	case upgrades.DevnetChainID:
		return DevnetUpgradeHeight, true
	default:
		return 0, false
	}
}
