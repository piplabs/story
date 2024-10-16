package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
)

func (k Keeper) ProcessSetWithdrawalAddress(ctx context.Context, ev *bindings.IPTokenStakingSetWithdrawalAddress) error {
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

func (k Keeper) ProcessSetRewardAddress(ctx context.Context, ev *bindings.IPTokenStakingSetRewardAddress) error {
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

func (k Keeper) ProcessAddOperator(ctx context.Context, ev *bindings.IPTokenStakingAddOperator) error {
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
