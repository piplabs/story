package keeper

import (
	"context"

	"cosmossdk.io/math"

	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	clog "github.com/piplabs/story/lib/log"
)

func (k *Keeper) ProcessUbiEvents(ctx context.Context, height uint64, logs []*types.EVMEvent) error {
	for _, evmLog := range logs {
		if err := evmLog.Verify(); err != nil {
			return errors.Wrap(err, "verify log [BUG]")
		}
		ethlog, err := evmLog.ToEthLog()
		if err != nil {
			return err
		}

		//nolint:gocritic,revive // more cases later
		switch ethlog.Topics[0] {
		case types.UBIPercentageSetEvent.ID:
			ev, err := k.ubiContract.ParseUBIPercentageSet(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse UBIPercentageSet log", err)
				continue
			}
			if err = k.ProcessUBIPercentageSet(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process UBI percentage set", err)
				continue
			}
		}
	}

	clog.Debug(ctx, "Processed UBIPool events", "height", height, "count", len(logs))

	return nil
}

func (k *Keeper) ProcessUBIPercentageSet(ctx context.Context, ev *bindings.UBIPoolUBIPercentageSet) error {
	newUBI := math.LegacyNewDecWithPrec(int64(ev.Percentage), 2)
	if err := k.distrKeeper.SetUbi(ctx, newUBI); err != nil {
		return errors.Wrap(err, "set new UBI percentage")
	}

	return nil
}
