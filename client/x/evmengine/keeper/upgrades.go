package keeper

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	clog "github.com/piplabs/story/lib/log"
)

func (k *Keeper) ProcessUpgradeEvents(ctx context.Context, height uint64, logs []*types.EVMEvent) error {
	for _, evmLog := range logs {
		if err := evmLog.Verify(); err != nil {
			return errors.Wrap(err, "verify log [BUG]") // This shouldn't happen
		}
		ethlog, err := evmLog.ToEthLog()
		if err != nil {
			return err
		}

		//nolint:gocritic,revive // more cases later
		switch ethlog.Topics[0] {
		case types.SoftwareUpgradeEvent.ID:
			ev, err := k.upgradeContract.ParseSoftwareUpgrade(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse SubmitProposal log", err)
				continue
			}
			if err = k.ProcessSoftwareUpgrade(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process submit proposal", err)
				continue
			}
		}
	}

	clog.Debug(ctx, "Processed governance events", "height", height, "count", len(logs))

	return nil
}

func (k *Keeper) ProcessSoftwareUpgrade(ctx context.Context, ev *bindings.UpgradeEntrypointSoftwareUpgrade) error {
	err := k.upgradeKeeper.ScheduleUpgrade(ctx, upgradetypes.Plan{
		Name:   ev.Name,
		Info:   ev.Info,
		Height: ev.Height,
	})
	if err != nil {
		return errors.Wrap(err, "process software upgrade: schedule upgrade")
	}

	return nil
}
