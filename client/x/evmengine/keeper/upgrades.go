//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
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

		switch ethlog.Topics[0] {
		case types.SoftwareUpgradeEvent.ID:
			ev, err := k.upgradeContract.ParseSoftwareUpgrade(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse SoftwareUpgrade log", err)
				continue
			}
			if err = k.ProcessSoftwareUpgrade(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process software upgrade", err)
				continue
			}
		case types.CancelUpgradeEvent.ID:
			ev, err := k.upgradeContract.ParseCancelUpgrade(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse CancelUpgrade log", err)
				continue
			}
			if err = k.ProcessCancelUpgrade(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process cancel upgrade", err)
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
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeUpgradeSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeUpgradeFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyUpgradeName, ev.Name),
				sdk.NewAttribute(types.AttributeKeyUpgradeHeight, strconv.FormatInt(ev.Height, 10)),
				sdk.NewAttribute(types.AttributeKeyUpgradeInfo, ev.Info),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
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

func (k *Keeper) ProcessCancelUpgrade(ctx context.Context, ev *bindings.UpgradeEntrypointCancelUpgrade) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeCancelUpgradeSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeCancelUpgradeFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.upgradeKeeper.ClearUpgradePlan(cachedCtx); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "process cancel upgrade: clear upgrade plan")
	}

	return nil
}
