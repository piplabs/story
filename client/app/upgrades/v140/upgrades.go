package v140

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/lib/log"
)

func CreateUpgradeHandler(
	_ *module.Manager,
	_ module.Configurator,
	_ *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		log.Info(ctx, "The upgrade to v1.4.0")

		return vm, nil
	}
}
