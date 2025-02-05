package virgil

import (
	"cosmossdk.io/math"
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
	StoryUpgradeHeight = 809988
	// DevnetUpgradeHeight defines the block height at which virgil upgrade is triggered on Internal Devnet.
	DevnetUpgradeHeight = 800
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

type RewardsMultipliers struct {
	Short  math.LegacyDec
	Medium math.LegacyDec
	Long   math.LegacyDec
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

var DefaultRewardsMultiplier = RewardsMultipliers{
	Short:  math.LegacyNewDecWithPrec(1051, 3), // 1.051
	Medium: math.LegacyNewDecWithPrec(116, 2),  // 1.16
	Long:   math.LegacyNewDecWithPrec(134, 2),  // 1.34
}

func GetRewardsMultipliers(chainID string) RewardsMultipliers {
	switch chainID {
	case upgrades.StoryChainID:
		return RewardsMultipliers{
			Short:  math.LegacyNewDecWithPrec(11, 1), // 1.1
			Medium: math.LegacyNewDecWithPrec(15, 1), // 1.5
			Long:   math.LegacyNewDecWithPrec(20, 1), // 2
		}
	case upgrades.DevnetChainID:
		return RewardsMultipliers{
			Short:  math.LegacyNewDecWithPrec(11, 1), // 1.1
			Medium: math.LegacyNewDecWithPrec(15, 1), // 1.5
			Long:   math.LegacyNewDecWithPrec(20, 1), // 2
		}
	default:
		return DefaultRewardsMultiplier
	}
}
