package test_v1_4

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/piplabs/story/client/app/upgrades"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

// Note: This test upgrade does NOT need to be in the Forks array
// because it's meant to be triggered via the UpgradeEntrypoint contract
// for testing rolling upgrades, not as a hard fork.
