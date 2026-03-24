package app

import (
	"fmt"
	"sort"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/upgrades"
	"github.com/piplabs/story/client/app/upgrades/horace"
	"github.com/piplabs/story/client/app/upgrades/polybius"
	"github.com/piplabs/story/client/app/upgrades/singularity/virgil"
	"github.com/piplabs/story/client/app/upgrades/terence"
	"github.com/piplabs/story/client/app/upgrades/v_1_2_0"
	"github.com/piplabs/story/client/app/upgrades/v_1_6_0"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/netconf"
)

var (
	// Upgrades defines the upgrade handlers and store loaders for the application.
	// New upgrades should be added to this slice after they are implemented.
	Upgrades = []upgrades.Upgrade{
		virgil.Upgrade,
		v_1_2_0.Upgrade,
		polybius.Upgrade,
		terence.Upgrade,
		horace.Upgrade,
		v_1_6_0.Upgrade,
	}
	// Forks are for hard forks that breaks backward compatibility.
	Forks = []upgrades.Fork{
		virgil.Fork,
		v_1_2_0.Fork,
		polybius.Fork,
		terence.Fork,
		horace.Fork,
		v_1_6_0.Fork,
	}
)

type StoreUpgradesMap map[int64]storetypes.StoreUpgrades

func (a *App) setupUpgradeHandlers() {
	for _, upgrade := range Upgrades {
		a.Keepers.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(a.ModuleManager, a.Configurator(), &a.Keepers),
		)
	}
}

// setupUpgradeStoreLoaders sets custom store loaders to customize the rootMultiStore initialization for software upgrades.
func (a *App) setupUpgradeStoreLoaders() {
	upgradeHistory, err := netconf.GetUpgradeHistory(a.ChainID())
	if err != nil {
		panic(errors.Wrap(err, "failed to get upgrade history"))
	}

	storeUpgradesMap := make(StoreUpgradesMap)

	for name, height := range upgradeHistory {
		if a.Keepers.UpgradeKeeper.IsSkipHeight(height) {
			continue
		}

		for _, upgrade := range Upgrades {
			if name == upgrade.UpgradeName {
				storeUpgradesMap[height] = upgrade.StoreUpgrades
			}
		}
	}

	// For binary-swap upgrades scheduled on-chain via planUpgrade (not in
	// UpgradeHistories), the old binary writes upgrade-info.json to disk
	// before halting. Read it to register the store upgrades at the correct
	// height so the new binary can mount new module stores on startup.
	diskPlan, diskErr := a.Keepers.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if diskErr == nil && diskPlan.Height > 0 && !a.Keepers.UpgradeKeeper.IsSkipHeight(diskPlan.Height) {
		storeUpgrades, err := GetStoreUpgrades(diskPlan.Name)
		if err == nil {
			storeUpgradesMap[diskPlan.Height] = storeUpgrades
		}
	}

	a.SetStoreLoader(UpgradeStoreLoader(storeUpgradesMap))
}

// UpgradeStoreLoader returns store loader including all upgrades.
// For future upgrades (height > current), it pre-adds new stores so the
// binary can start before the upgrade height without a store version
// mismatch. Note: pre-adding stores changes the app hash at the next
// commit, so ALL validators must switch to the new binary together.
func UpgradeStoreLoader(storeUpgradesMap StoreUpgradesMap) baseapp.StoreLoader {
	return func(ms storetypes.CommitMultiStore) error {
		lastVersion := ms.LastCommitID().Version
		nextVersion := lastVersion + 1

		// On a fresh genesis chain (no commits yet), all modules are already
		// registered and their stores are mounted. Applying store upgrades
		// would conflict with existing stores (e.g., trying to add DKG store
		// that already exists from genesis). Skip the upgrade loader entirely.
		if lastVersion == 0 {
			return baseapp.DefaultStoreLoader(ms)
		}

		// Sort heights for deterministic iteration order across all validators.
		heights := make([]int64, 0, len(storeUpgradesMap))
		for h := range storeUpgradesMap {
			heights = append(heights, h)
		}
		sort.Slice(heights, func(i, j int) bool { return heights[i] < heights[j] })

		var merged storetypes.StoreUpgrades
		addedSet := make(map[string]bool)

		for _, height := range heights {
			su := storeUpgradesMap[height]
			if height == nextVersion {
				// Exact upgrade height: apply all operations.
				// Do NOT filter by mountedStores here — the new binary
				// mounts modules at startup (via app_config.go), but the
				// store does not yet exist on disk. We must include it in
				// Added so LoadLatestVersionAndUpgrade creates it at the
				// correct version.
				for _, key := range su.Added {
					if !addedSet[key] {
						merged.Added = append(merged.Added, key)
						addedSet[key] = true
					}
				}
				merged.Deleted = append(merged.Deleted, su.Deleted...)
				merged.Renamed = append(merged.Renamed, su.Renamed...)
			} else if height > nextVersion {
				// Future upgrade: pre-add new stores so the binary can
				// load without crashing on missing stores. Same logic —
				// the store is mounted in code but not yet on disk.
				for _, key := range su.Added {
					if !addedSet[key] {
						merged.Added = append(merged.Added, key)
						addedSet[key] = true
					}
				}
			}
			// Past upgrades (height <= lastVersion) are skipped — those
			// stores were already added during the original upgrade.
		}

		if len(merged.Renamed) > 0 || len(merged.Deleted) > 0 || len(merged.Added) > 0 {
			return ms.LoadLatestVersionAndUpgrade(&merged)
		}

		return baseapp.DefaultStoreLoader(ms)
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
		upgradeHeight, ok := GetUpgradeHeight(ctx, fork.UpgradeName, fork.UpgradeHeight)
		if !ok {
			continue
		}

		if currentBlockHeight == upgradeHeight {
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

// GetUpgradeHeight returns the upgrade height for a given upgrade name,
// or the static fallback height if no dynamic resolution exists.
// If the upgrade is not applicable for this chain, it returns (0, false).
func GetUpgradeHeight(ctx sdk.Context, upgradeName string, fallbackHeight int64) (int64, bool) {
	switch upgradeName {
	case virgil.UpgradeName:
		return virgil.GetUpgradeHeight(ctx.ChainID())

	case v_1_2_0.UpgradeName:
		return v_1_2_0.GetUpgradeHeight(ctx.ChainID())

	case polybius.UpgradeName:
		return polybius.GetUpgradeHeight(ctx.ChainID())

	case netconf.Terence:
		return terence.GetUpgradeHeight(ctx)

	case netconf.Horace:
		return horace.GetUpgradeHeight(ctx)

	default:
		// no dynamic resolver → use fallback (static height)
		return fallbackHeight, true
	}
}

// GetStoreUpgrades returns the store upgrades for a given scheduled upgrade on-chain.
// This is used by the disk-based fallback in setupUpgradeStoreLoaders to
// determine which stores to add when the upgrade height comes from
// upgrade-info.json rather than from hardcoded UpgradeHistories.
func GetStoreUpgrades(upgradeName string) (storetypes.StoreUpgrades, error) {
	switch upgradeName {
	case netconf.V160:
		return v_1_6_0.Upgrade.StoreUpgrades, nil
	default:
		return storetypes.StoreUpgrades{}, errors.New("no matched store upgrades")
	}
}
