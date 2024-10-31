package app

import (
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/upgrades"
	"github.com/piplabs/story/client/app/upgrades/v0_12_1"
)

var (
	// `Upgrades` defines the upgrade handlers and store loaders for the application.
	// New upgrades should be added to this slice after they are implemented.
	Upgrades = []upgrades.Upgrade{
		v0_12_1.Upgrade,
	}
	// Forks are for hard forks that breaks backward compatibility.
	Forks = []upgrades.Fork{}
)

func (a *App) setupUpgradeHandlers() {
	for _, upgrade := range Upgrades {
		a.Keepers.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(a.ModuleManager, a.Configurator(), &a.Keepers),
		)
	}
}

// setUpgradeStoreLoaders sets custom store loaders to customize the rootMultiStore initialization for software upgrades.
func (a *App) setupUpgradeStoreLoaders() {
	upgradeInfo, err := a.Keepers.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if a.Keepers.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrade := range Upgrades {
		if upgradeInfo.Name == upgrade.UpgradeName {
			a.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &upgrade.StoreUpgrades))
		}
	}
}

// ScheduleForkUpgrade executes any necessary fork logic for based upon the current block height. It sets an upgrade
// plan once the chain reaches the pre-defined upgrade height.
//
// CONTRACT: for this logic to work properly it is required to:
//  1. Release a non-breaking patch version so that the chain can set the scheduled upgrade plan at upgrade-height.
//  2. Release the software defined in the upgrade-info.
func (a *App) scheduleForkUpgrade(ctx sdk.Context) {
	currentBlockHeight := ctx.BlockHeight()
	for _, fork := range Forks {
		if currentBlockHeight == fork.UpgradeHeight {
			upgradePlan := upgradetypes.Plan{
				Height: currentBlockHeight,
				Name:   fork.UpgradeName,
				Info:   fork.UpgradeInfo,
			}

			// schedule the upgrade plan to the current block height, effectively performing
			// a hard fork that uses the upgrade handler to manage the migration.
			if err := a.Keepers.UpgradeKeeper.ScheduleUpgrade(ctx, upgradePlan); err != nil {
				panic(
					//nolint:errorlint // use "%v" to obfuscate the underlying error
					fmt.Errorf(
						"hard fork: failed to schedule upgrade %s during BeginBlock at height %d: %v",
						upgradePlan.Name,
						ctx.BlockHeight(),
						err,
					),
				)
			}
		}
	}
}
