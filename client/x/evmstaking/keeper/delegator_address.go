//nolint:dupl // event log
package keeper

import (
	"context"
	"encoding/hex"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
)

func (k Keeper) ProcessSetWithdrawalAddress(ctx context.Context, ev *bindings.IPTokenStakingSetWithdrawalAddress) (err error) {
	defer func() {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		if err != nil {
			sdkCtx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeSetWithdrawalAddressFailure,
					sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
					sdk.NewAttribute(types.AttributeKeyDelegatorUncmpPubKey, hex.EncodeToString(ev.DelegatorUncmpPubkey)),
					sdk.NewAttribute(types.AttributeKeyRewardAddress, hex.EncodeToString(ev.ExecutionAddress[:])),
				),
			})
		}
	}()

	delCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.DelegatorUncmpPubkey)
	if err != nil {
		return errors.Wrap(err, "compress depositor pubkey")
	}
	depositorPubkey, err := k1util.PubKeyBytesToCosmos(delCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "depositor pubkey to cosmos")
	}

	depositorAddr := sdk.AccAddress(depositorPubkey.Address().Bytes())
	executionAddr := common.BytesToAddress(ev.ExecutionAddress[:])

	if err := k.DelegatorWithdrawAddress.Set(ctx, depositorAddr.String(), executionAddr.String()); err != nil {
		return errors.Wrap(err, "delegator withdraw address map set")
	}

	return nil
}

func (k Keeper) ProcessSetRewardAddress(ctx context.Context, ev *bindings.IPTokenStakingSetRewardAddress) (err error) {
	defer func() {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		if err != nil {
			sdkCtx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeSetRewardAddressFailure,
					sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
					sdk.NewAttribute(types.AttributeKeyDelegatorUncmpPubKey, hex.EncodeToString(ev.DelegatorUncmpPubkey)),
					sdk.NewAttribute(types.AttributeKeyWithdrawalAddress, hex.EncodeToString(ev.ExecutionAddress[:])),
				),
			})
		}
	}()

	delCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.DelegatorUncmpPubkey)
	if err != nil {
		return errors.Wrap(err, "compress depositor pubkey")
	}
	depositorPubkey, err := k1util.PubKeyBytesToCosmos(delCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "depositor pubkey to cosmos")
	}

	depositorAddr := sdk.AccAddress(depositorPubkey.Address().Bytes())
	executionAddr := common.BytesToAddress(ev.ExecutionAddress[:])

	if err := k.DelegatorRewardAddress.Set(ctx, depositorAddr.String(), executionAddr.String()); err != nil {
		return errors.Wrap(err, "delegator reward address map set")
	}

	return nil
}

func (k Keeper) ProcessAddOperator(ctx context.Context, ev *bindings.IPTokenStakingAddOperator) (err error) {
	defer func() {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		if err != nil {
			sdkCtx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeAddOperatorFailure,
					sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
					sdk.NewAttribute(types.AttributeKeyDelegatorUncmpPubKey, hex.EncodeToString(ev.UncmpPubkey)),
					sdk.NewAttribute(types.AttributeKeyOperatorAddress, ev.Operator.Hex()),
				),
			})
		}
	}()

	delCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.UncmpPubkey)
	if err != nil {
		return errors.Wrap(err, "compress depositor pubkey")
	}
	depositorPubkey, err := k1util.PubKeyBytesToCosmos(delCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "depositor pubkey to cosmos")
	}

	depositorAddr := sdk.AccAddress(depositorPubkey.Address().Bytes())

	if err := k.DelegatorOperatorAddress.Set(ctx, depositorAddr.String(), ev.Operator.String()); err != nil {
		return errors.Wrap(err, "delegator operator address map set")
	}

	return nil
}

func (k Keeper) ProcessRemoveOperator(ctx context.Context, ev *bindings.IPTokenStakingRemoveOperator) (err error) {
	defer func() {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		if err != nil {
			sdkCtx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeRemoveOperatorFailure,
					sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
					sdk.NewAttribute(types.AttributeKeyDelegatorUncmpPubKey, hex.EncodeToString(ev.UncmpPubkey)),
					sdk.NewAttribute(types.AttributeKeyOperatorAddress, ev.Operator.Hex()),
				),
			})
		}
	}()

	delCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.UncmpPubkey)
	if err != nil {
		return errors.Wrap(err, "compress depositor pubkey")
	}
	depositorPubkey, err := k1util.PubKeyBytesToCosmos(delCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "depositor pubkey to cosmos")
	}

	depositorAddr := sdk.AccAddress(depositorPubkey.Address().Bytes())

	if err := k.DelegatorOperatorAddress.Remove(ctx, depositorAddr.String()); err != nil {
		return errors.Wrap(err, "delegator operator address map remove")
	}

	return nil
}
