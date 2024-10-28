//nolint:revive,stylecheck // version underscores
package v0_12_1

import (
	"context"
	"time"

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

		clog.Info(ctx, "Decreasing staking periods...")
		stakingParams, err := keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get staking params")
		}
		periods := stakingParams.GetPeriods()

		for _, p := range periods {
			if p.PeriodType == 1 {
				clog.Info(ctx, "Existing short period duration: %d", p.Duration)
				p.Duration = time.Hour * 7 * 24 // one week in hours
			} else if p.PeriodType == 2 {
				clog.Info(ctx, "Existing medium period duration: %d", p.Duration)
				p.Duration = time.Hour * 14 * 24 // two weeks in hours
			} else if p.PeriodType == 3 {
				clog.Info(ctx, "Existing long period duration: %d", p.Duration)
				p.Duration = time.Hour * 21 * 24 // three weeks in hours
			}
		}

		err = keepers.StakingKeeper.SetParams(ctx, stakingParams)
		if err != nil {
			return vm, errors.Wrap(err, "failed to set staking params")
		}

		clog.Info(ctx, "Checking newstaking periods...")
		stakingParams, err = keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get staking params")
		}
		periods = stakingParams.GetPeriods()

		for _, p := range periods {
			if p.PeriodType == 1 {
				clog.Info(ctx, "New short period duration: %d", p.Duration)
			} else if p.PeriodType == 2 {
				clog.Info(ctx, "New medium period duration: %d", p.Duration)
			} else if p.PeriodType == 3 {
				clog.Info(ctx, "New long period duration: %d", p.Duration)
			}
		}

		clog.Info(ctx, "Upgrade v0.12.1 complete")

		return vm, nil
	}
}
