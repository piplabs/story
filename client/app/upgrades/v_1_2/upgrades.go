//nolint:revive,stylecheck // versioning
package v_1_2

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/piplabs/story/client/app/keepers"
	estypes "github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func CreateUpgradeHandler(
	_ *module.Manager,
	_ module.Configurator,
	keepers *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		blockHeight := sdkCtx.BlockHeight()
		log.Info(ctx, "Current block height", "Height", blockHeight)

		// Check if the upgrade is needed for current chain
		chainID := sdkCtx.ChainID()
		if _, ok := GetUpgradeHeight(chainID); !ok {
			log.Info(ctx, "Upgrade v1.2 not needed for current chain, skip", "ChainID", chainID)
			return vm, nil
		}
		log.Info(ctx, "Start upgrade v1.2", "ChainID", chainID)

		oldParams, err := keepers.EvmStakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get evm staking params")
		}

		if err = keepers.EvmStakingKeeper.SetParams(ctx, estypes.Params{
			MaxWithdrawalPerBlock:      oldParams.MaxWithdrawalPerBlock,
			MaxSweepPerBlock:           oldParams.MaxSweepPerBlock,
			MinPartialWithdrawalAmount: oldParams.MinPartialWithdrawalAmount,
			RefundFeeBps:               InitialRefundFeeBps,
			RefundPeriod:               InitialRefundPeriod,
		}); err != nil {
			return vm, errors.Wrap(err, "failed to set evm staking params")
		}

		log.Info(ctx, "Upgrade v1.2 complete")

		return vm, nil
	}
}
