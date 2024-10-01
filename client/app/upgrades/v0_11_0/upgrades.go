//nolint:revive,stylecheck // version underscores
package v0_11_0

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/genutil/evm/predeploys"
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

		clog.Info(ctx, "Setting updated IPTokenSlashing address...")

		// Upgrade to use the fixed slashing contract
		predeploys.UpdatedIPTokenSlashing = "0xEEf1c4fD443965404f13BE2705766988317b3B32"

		clog.Info(ctx, "Upgrade v0.11.0 complete")

		return vm, nil
	}
}
