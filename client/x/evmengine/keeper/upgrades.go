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
	"github.com/piplabs/story/lib/netconf"
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

	clog.Info(ctx, "SoftwareUpgradeEvent detected")

	isV140, err := netconf.IsV140(cachedCtx.ChainID(), cachedCtx.BlockHeight())
	if err != nil {
		return errors.Wrap(err, "failed to check v1.4.0 upgrade height")
	}

	if err = k.ScheduleUpgrade(cachedCtx, ev, isV140); err != nil {
		return errors.Wrap(err, "failed to schedule upgrade")
	}

	if err := k.upgradeKeeper.DumpUpgradeInfoToDisk(ev.Height, upgradetypes.Plan{
		Name:   ev.Name,
		Height: ev.Height,
		Info:   ev.Info,
	}); err != nil {
		clog.Warn(ctx, "Fail to dump upgrade info to disk after schedule upgrade. please set upgrade-info.json manually for using cosmovisor", err)
	}

	clog.Info(ctx, "Upgrade is scheduled successfully", "upgrade_name", ev.Name, "upgrade_height", ev.Height, "upgrade_info", ev.Info)

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

	clog.Info(ctx, "SoftwareUpgradeCancelEvent detected")

	isV140, err := netconf.IsV140(cachedCtx.ChainID(), cachedCtx.BlockHeight())
	if err != nil {
		return errors.Wrap(err, "failed to check v1.4.0 upgrade height")
	}

	if err = k.CancelUpgrade(cachedCtx, isV140); err != nil {
		return errors.Wrap(err, "failed to cancel the upgrade")
	}

	clog.Info(ctx, "Upgrade is canceled")

	return nil
}

func (k *Keeper) ScheduleUpgrade(ctx sdk.Context, ev *bindings.UpgradeEntrypointSoftwareUpgrade, isV140 bool) (err error) {
	plan := upgradetypes.Plan{
		Name:   ev.Name,
		Info:   ev.Info,
		Height: ev.Height,
	}
	if isV140 {
		if err = k.SetPendingUpgrade(ctx, plan); errors.Is(err, types.ErrUpgradePending) {
			return errors.WrapErrWithCode(errors.PendingUpgradeExists, err)
		} else if err != nil {
			return errors.Wrap(err, "failed to set pending upgrade")
		}

		return nil
	}

	if err = k.upgradeKeeper.ScheduleUpgrade(ctx, plan); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "failed to schedule upgrade in upgrade module")
	}

	return nil
}

func (k *Keeper) CancelUpgrade(ctx sdk.Context, isV140 bool) (err error) {
	if isV140 {
		if err = k.ResetPendingUpgrade(ctx); err != nil {
			return errors.Wrap(err, "failed to reset upgrade")
		}

		return nil
	}

	if err = k.upgradeKeeper.ClearUpgradePlan(ctx); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "failed to clear upgrade plan")
	}

	return nil
}

// ShouldUpgrade returns whether the signaling mechanism has concluded that the
// network is ready to upgrade. It returns false and nil if no upgrade is scheduled or has not reached to the upgrade
// height.
func (k *Keeper) ShouldUpgrade(ctx sdk.Context) (bool, upgradetypes.Plan) {
	upgradePlan, err := k.getPendingUpgrade(ctx)
	if err != nil {
		return false, upgradetypes.Plan{}
	}

	hasUpgradeHeightBeenReached := ctx.BlockHeight() >= upgradePlan.Height
	if hasUpgradeHeightBeenReached {
		return true, upgradePlan
	}

	return false, upgradetypes.Plan{}
}

func (k *Keeper) PendingUpgrade(ctx context.Context) (upgradetypes.Plan, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	return k.getPendingUpgrade(sdkCtx)
}

// SetPendingUpgrade sets the upgrade if there is no pending upgrade.
func (k *Keeper) SetPendingUpgrade(ctx sdk.Context, plan upgradetypes.Plan) error {
	if err := plan.ValidateBasic(); err != nil {
		return errors.Wrap(err, "invalid plan")
	}

	if plan.Height < ctx.BlockHeight() {
		return errors.WrapErrWithCode(errors.InvalidRequest, sdkerrors.ErrInvalidRequest)
	}

	if k.IsUpgradePending(ctx) {
		return types.ErrUpgradePending
	}

	if err := k.setPendingUpgrade(ctx, plan); err != nil {
		return errors.Wrap(err, "failed to set upgrade plan")
	}

	return nil
}

// ResetPendingUpgrade resets the upgrade. It is called after an upgrade is scheduled successfully in upgrade keeper or when
// CancelUpgrade transaction is processed.
func (k *Keeper) ResetPendingUpgrade(ctx sdk.Context) error {
	stores := k.storeService.OpenKVStore(ctx)

	err := stores.Delete(types.PendingUpgradeKey)
	if err != nil {
		return errors.Wrap(err, "failed to delete pending upgrade")
	}

	return nil
}

// IsUpgradePending returns true if the chain should upgrade at the upgrade height. While the
// keeper has an upgrade pending new upgrade will be rejected. To schedule a new upgrade, the existing one should be
// canceled first.
func (k *Keeper) IsUpgradePending(ctx sdk.Context) bool {
	_, err := k.getPendingUpgrade(ctx)

	return err == nil
}

// getPendingUpgrade returns the current upgrade information from the store.
// If an upgrade is found, it returns the Plan object and nil error.
func (k *Keeper) getPendingUpgrade(ctx sdk.Context) (upgradetypes.Plan, error) {
	stores := k.storeService.OpenKVStore(ctx)
	value, err := stores.Get(types.PendingUpgradeKey)
	if err != nil {
		return upgradetypes.Plan{}, errors.Wrap(err, "failed to get upgrade")
	}
	if value == nil {
		return upgradetypes.Plan{}, types.ErrUpgradeNotFound
	}

	var upgrade upgradetypes.Plan
	if err = k.cdc.Unmarshal(value, &upgrade); err != nil {
		return upgradetypes.Plan{}, errors.Wrap(err, "failed to unmarshal")
	}

	return upgrade, nil
}

// setPendingUpgrade sets the upgrade plan in the store.
func (k *Keeper) setPendingUpgrade(ctx sdk.Context, upgrade upgradetypes.Plan) error {
	stores := k.storeService.OpenKVStore(ctx)
	value := k.cdc.MustMarshal(&upgrade)

	err := stores.Set(types.PendingUpgradeKey, value)
	if err != nil {
		return errors.Wrap(err, "failed to set pending upgrade")
	}

	return nil
}
