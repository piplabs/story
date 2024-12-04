//nolint:revive,stylecheck // version underscores
package v0_13_2

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		log.Info(ctx, "Starting module migrations...")

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, errors.Wrap(err, "run migrations")
		}

		log.Info(ctx, "Getting existing UnbondingID...")
		unbondingID, err := keepers.StakingKeeper.GetOldUnbondingID(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get existing UnbondingID")
		}
		log.Info(ctx, "Existing UnbondingID", "UnbondingID", unbondingID)

		log.Info(ctx, "Setting new UnbondingID...")
		if err := keepers.StakingKeeper.SetUnbondingID(ctx, unbondingID); err != nil {
			return vm, errors.Wrap(err, "failed to set new UnbondingID")
		}

		log.Info(ctx, "Getting new UnbondingID...")
		newUnbondingID, err := keepers.StakingKeeper.GetUnbondingID(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get new UnbondingID")
		}
		log.Info(ctx, "New UnbondingID", "UnbondingID", newUnbondingID)

		if newUnbondingID != unbondingID {
			return vm, errors.New("new UnbondingID does not match existing UnbondingID")
		}

		log.Info(ctx, "Removing old UnbondingID...")
		// TODO: Remove old UnbondingID

		log.Info(ctx, "Upgrade v0.13.2 complete")

		return vm, nil
	}
}
