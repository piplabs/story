package test_v1_4_2

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
	// NewMinDelegation is the new minimum delegation amount (1 ether = 1e18 wei)
	NewMinDelegation = "1000000000000000000" // 1 ether
)

func CreateUpgradeHandler(
	_ *module.Manager,
	_ module.Configurator,
	keepers *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		log.Info(ctx, "Start test-v1.4.2 upgrade")

		// Get current staking params
		log.Info(ctx, "Get current staking params...")
		stakingParams, err := keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get staking params")
		}

		oldMinDelegation := stakingParams.MinDelegation
		log.Info(ctx, "Current min_delegation", "value", oldMinDelegation.String())

		// Update min_delegation from 1024 ether to 1 ether
		newMinDelegation, ok := math.NewIntFromString(NewMinDelegation)
		if !ok {
			return vm, errors.New("failed to parse new min_delegation value")
		}

		log.Info(ctx, "Update min_delegation", "old_value", oldMinDelegation.String(), "new_value", newMinDelegation.String())
		stakingParams.MinDelegation = newMinDelegation

		// Apply staking param changes
		log.Info(ctx, "Apply staking param changes...")
		if err := keepers.StakingKeeper.SetParams(ctx, stakingParams); err != nil {
			return vm, errors.Wrap(err, "failed to update staking params")
		}

		// Verify the changes
		log.Info(ctx, "Verify new staking params...")
		updatedParams, err := keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get updated staking params")
		}

		if !updatedParams.MinDelegation.Equal(newMinDelegation) {
			return vm, errors.New("min_delegation update verification failed",
				"expected", newMinDelegation.String(),
				"actual", updatedParams.MinDelegation.String())
		}

		log.Info(ctx, "Min_delegation successfully updated", "new_value", updatedParams.MinDelegation.String())

		// Note: Contract upgrade and setMinStakeAmount/setMinUnstakeAmount calls
		// should be done off-chain via governance before or after this hardfork.
		log.Info(ctx, "The test-v1.4.2 upgrade has been completed")
		log.Info(ctx, "NOTE: Contract upgrade and setMinStakeAmount/setMinUnstakeAmount calls must be done via governance")

		return vm, nil
	}
}
