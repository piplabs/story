//nolint:revive,stylecheck // version underscores
package v0_12_1

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/lib/errors"
	clog "github.com/piplabs/story/lib/log"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		clog.Info(ctx, "Starting module migrations...")

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, errors.Wrap(err, "run migrations")
		}

		clog.Info(ctx, "Upgrade v0.12.1 complete")

		return vm, nil
	}
}
