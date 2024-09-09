//nolint:revive,stylecheck // version underscores
package v0_10_0

import (
	"context"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/lib/errors"
	clog "github.com/piplabs/story/lib/log"
)

const (
	NewMaxSweepPerBlock = 1024
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

		clog.Info(ctx, "Setting NextValidatorDelegationSweepIndex parameter...")
		nextValIndex, err := keepers.EvmStakingKeeper.GetOldValidatorSweepIndex(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "get old validator sweep index")
		}

		nextValDelIndex := sdk.IntProto{Int: math.NewInt(0)}
		if err := keepers.EvmStakingKeeper.SetValidatorSweepIndex(
			ctx,
			nextValIndex,
			nextValDelIndex,
		); err != nil {
			return vm, errors.Wrap(err, "set evmstaking NextValidatorDelegationSweepIndex")
		}

		// Update MaxSweepPerBlock
		clog.Info(ctx, "Updating MaxSweepPerBlock parameter...")
		params, err := keepers.EvmStakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "get evmstaking params")
		}

		params.MaxSweepPerBlock = NewMaxSweepPerBlock
		if err = keepers.EvmStakingKeeper.SetParams(ctx, params); err != nil {
			return vm, errors.Wrap(err, "set evmstaking params")
		}

		clog.Info(ctx, "Upgrade v0.10.0 complete")

		return vm, nil
	}
}
