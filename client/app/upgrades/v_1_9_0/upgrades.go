package v_1_9_0

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func CreateUpgradeHandler(
	_ *module.Manager,
	_ module.Configurator,
	keepers *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		log.Info(ctx, "Start v1.9.0 upgrade — backfilling DKG partial decrypt round index")

		if err := keepers.DKGKeeper.MigratePartialDecryptRoundIndex(ctx); err != nil {
			return vm, errors.Wrap(err, "migrate partial decrypt round index")
		}

		log.Info(ctx, "V1.9.0 upgrade complete — DKG partial decrypt round index activated")

		return vm, nil
	}
}
