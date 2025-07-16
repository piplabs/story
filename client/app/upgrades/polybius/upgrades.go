package polybius

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
		log.Info(ctx, "Start upgrade Polybius")

		params, err := keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get existing params of staking module")
		}

		params.MaxValidators = NewMaxValidators
		if err := keepers.StakingKeeper.SetParams(ctx, params); err != nil {
			return vm, errors.Wrap(err, "failed to set new params of staking module")
		}

		newParams, err := keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get new params of staking module")
		}

		if !newParams.Equal(params) {
			return vm, errors.New("new params mismatch")
		}

		log.Info(ctx, "Upgrade Polybius complete")

		return vm, nil
	}
}
