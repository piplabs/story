//nolint:dupl,contextcheck // event log
package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
)

func (k Keeper) ProcessSetWithdrawalAddress(ctx context.Context, ev *bindings.IPTokenStakingSetWithdrawalAddress) (err error) {
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
				types.EventTypeSetWithdrawalAddressSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeSetWithdrawalAddressFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDelegatorAddress, ev.Delegator.String()),
				sdk.NewAttribute(types.AttributeKeyRewardAddress, hex.EncodeToString(ev.ExecutionAddress[:])),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	depositorAddr := sdk.AccAddress(ev.Delegator.Bytes())
	executionAddr := common.BytesToAddress(ev.ExecutionAddress[:])

	if err := k.DelegatorWithdrawAddress.Set(cachedCtx, depositorAddr.String(), executionAddr.String()); err != nil {
		return errors.Wrap(err, "delegator withdraw address map set")
	}

	return nil
}

func (k Keeper) ProcessSetRewardAddress(ctx context.Context, ev *bindings.IPTokenStakingSetRewardAddress) (err error) {
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
				types.EventTypeSetRewardAddressSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeSetRewardAddressFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDelegatorAddress, ev.Delegator.String()),
				sdk.NewAttribute(types.AttributeKeyWithdrawalAddress, hex.EncodeToString(ev.ExecutionAddress[:])),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	depositorAddr := sdk.AccAddress(ev.Delegator.Bytes())
	executionAddr := common.BytesToAddress(ev.ExecutionAddress[:])

	if err := k.DelegatorRewardAddress.Set(cachedCtx, depositorAddr.String(), executionAddr.String()); err != nil {
		return errors.Wrap(err, "delegator reward address map set")
	}

	return nil
}

func (k Keeper) ProcessSetOperator(ctx context.Context, ev *bindings.IPTokenStakingSetOperator) (err error) {
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
				types.EventTypeSetOperatorSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeSetOperatorFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDelegatorAddress, ev.Delegator.String()),
				sdk.NewAttribute(types.AttributeKeyOperatorAddress, ev.Operator.String()),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	depositorAddr := sdk.AccAddress(ev.Delegator.Bytes())

	if err := k.DelegatorOperatorAddress.Set(cachedCtx, depositorAddr.String(), ev.Operator.String()); err != nil {
		return errors.Wrap(err, "delegator operator address map set")
	}

	return nil
}

func (k Keeper) ProcessUnsetOperator(ctx context.Context, ev *bindings.IPTokenStakingUnsetOperator) (err error) {
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
				types.EventTypeUnsetOperatorSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeUnsetOperatorFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDelegatorAddress, ev.Delegator.String()),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	depositorAddr := sdk.AccAddress(ev.Delegator.Bytes())

	if err := k.DelegatorOperatorAddress.Remove(cachedCtx, depositorAddr.String()); err != nil {
		return errors.Wrap(err, "delegator operator address map remove")
	}

	return nil
}
