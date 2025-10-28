//nolint:revive // Test upgrade for v1.4
package test_v1_4

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

const UpgradeName = "test_v1_4"

func CreateUpgradeHandler(
	_ *module.Manager,
	_ module.Configurator,
	keep *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		log.Info(ctx, "Starting test_v1_4 upgrade")

		// Get current EVM staking parameters
		oldParams, err := keep.EvmStakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, err
		}

		log.Info(ctx, "Current MaxSweepPerBlock", "value", oldParams.MaxSweepPerBlock)

		// Simple parameter change: increase MaxSweepPerBlock by 16
		newParams := oldParams
		newParams.MaxSweepPerBlock = oldParams.MaxSweepPerBlock + 16

		// Validate the new parameters
		if err := newParams.Validate(); err != nil {
			return vm, err
		}

		// Apply the parameter change
		if err := keep.EvmStakingKeeper.SetParams(ctx, newParams); err != nil {
			return vm, err
		}

		// Verify the change was applied
		verifyParams, err := keep.EvmStakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, err
		}

		if verifyParams.MaxSweepPerBlock != newParams.MaxSweepPerBlock {
			err := errors.New("MaxSweepPerBlock was not updated correctly")
			log.Error(ctx, "upgrade failed", err)
			return vm, err
		}

		log.Info(ctx, "test_v1_4 upgrade completed successfully",
			"old_max_sweep_per_block", oldParams.MaxSweepPerBlock,
			"new_max_sweep_per_block", verifyParams.MaxSweepPerBlock)

		return vm, nil
	}
}
