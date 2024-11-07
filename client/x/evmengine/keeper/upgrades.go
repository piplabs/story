//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"strconv"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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

func (k *Keeper) ProcessSoftwareUpgrade(ctx context.Context, ev *bindings.UpgradeEntrypointSoftwareUpgrade) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	defer func() {
		if err == nil {
			writeCache()
			return
		}
		sdkCtx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeUpgradeFailure,
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyUpgradeName, ev.Name),
				sdk.NewAttribute(types.AttributeKeyUpgradeHeight, strconv.FormatInt(ev.Height, 10)),
				sdk.NewAttribute(types.AttributeKeyUpgradeInfo, ev.Info),
				sdk.NewAttribute(types.AttributeKeyStatusCode, errors.UnwrapErrCode(err).String()),
			),
		})
	}()

	if err = k.upgradeKeeper.ScheduleUpgrade(cachedCtx, upgradetypes.Plan{
		Name:   ev.Name,
		Info:   ev.Info,
		Height: ev.Height,
	}); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "process software upgrade: schedule upgrade")
	}

	return nil
}
