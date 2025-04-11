//nolint:revive,stylecheck // versioning
package v_1_3_0

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/ethereum/go-ethereum/common"

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
		log.Info(ctx, "Start upgrade v1.3.0")

		// When this handler is triggered at the beginning of block X, we must process the events from block X-1.
		// Otherwise, events from block X-1 will be discarded as the upgrade changes the logic to process events of the
		// current block once finalized.

		// Get the block hash and events of the previous block.
		head, err := keepers.EVMEngKeeper.GetExecutionHead(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "get execution head")
		}
		lastBlockHeight := head.GetBlockHeight()
		lastBlockHash := common.BytesToHash(head.GetBlockHash())

		events, err := keepers.EVMEngKeeper.EvmEvents(ctx, lastBlockHash)
		if err != nil {
			return nil, errors.Wrap(err, "fetch evm event logs")
		}

		// Deliver all the payload log events of the block X-1
		if err := keepers.EvmStakingKeeper.ProcessStakingEvents(ctx, lastBlockHeight, events); err != nil {
			return nil, errors.Wrap(err, "deliver staking-related event logs")
		}
		if err := keepers.EVMEngKeeper.ProcessUpgradeEvents(ctx, lastBlockHeight, events); err != nil {
			return nil, errors.Wrap(err, "deliver upgrade-related event logs")
		}
		if err := keepers.EVMEngKeeper.ProcessUbiEvents(ctx, lastBlockHeight, events); err != nil {
			return nil, errors.Wrap(err, "deliver ubi-related event logs")
		}

		if err := keepers.EVMEngKeeper.UpdateExecutionHeadWithBlock(ctx, lastBlockHeight, lastBlockHash, head.GetBlockTime()); err != nil {
			return nil, errors.Wrap(err, "update execution head")
		}

		log.Info(ctx, "Upgrade v1.3.0 complete")

		return vm, nil
	}
}
