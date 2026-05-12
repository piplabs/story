package v_1_7_0

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
		log.Info(ctx, "Start v1.7.0 upgrade — backfilling DKG partial decrypt round index")

		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		if err := keepers.DKGKeeper.MigratePartialDecryptRoundIndex(ctx); err != nil {
			return newVM, errors.Wrap(err, "migrate partial decrypt round index")
		}

		log.Info(ctx, "V1.7.0 upgrade complete — DKG partial decrypt round index migrated")

		return newVM, nil
	}
}
