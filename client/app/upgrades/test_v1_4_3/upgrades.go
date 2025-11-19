package test_v1_4_3

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

const (
	// NewMinDelegationIP is the new minimum delegation amount (1 IP token = 1e9 stake units)
	NewMinDelegationIP = "1000000000"
)

func CreateUpgradeHandler(
	_ *module.Manager,
	_ module.Configurator,
	keepers *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		log.Info(ctx, "Start test-v1.4.3 upgrade")

		stakingParams, err := keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get staking params")
		}

		oldMinDelegation := stakingParams.MinDelegation
		log.Info(ctx, "Current min_delegation", "value", oldMinDelegation.String())

		newMinDelegation, ok := math.NewIntFromString(NewMinDelegationIP)
		if !ok {
			return vm, errors.New("failed to parse new min_delegation value")
		}

		log.Info(ctx, "Update min_delegation to 1 IP token worth of stake units",
			"old_value", oldMinDelegation.String(),
			"new_value", newMinDelegation.String(),
		)
		stakingParams.MinDelegation = newMinDelegation

		if err := keepers.StakingKeeper.SetParams(ctx, stakingParams); err != nil {
			return vm, errors.Wrap(err, "failed to update staking params")
		}

		updatedParams, err := keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get updated staking params")
		}

		if !updatedParams.MinDelegation.Equal(newMinDelegation) {
			return vm, errors.New("min_delegation update verification failed",
				"expected", newMinDelegation.String(),
				"actual", updatedParams.MinDelegation.String(),
			)
		}

		log.Info(ctx, "Min_delegation successfully updated", "new_value", updatedParams.MinDelegation.String())
		log.Info(ctx, "The test-v1.4.3 upgrade has been completed")

		return vm, nil
	}
}
