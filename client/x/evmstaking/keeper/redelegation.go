//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
)

func (k Keeper) ProcessRedelegate(ctx context.Context, ev *bindings.IPTokenStakingRedelegate) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	var actualAmount string

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeRedelegateSuccess,
				sdk.NewAttribute(types.AttributeKeyAmount, actualAmount),
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeRedelegateFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
				sdk.NewAttribute(types.AttributeKeyAmount, ev.Amount.String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDelegatorAddress, ev.Delegator.String()),
				sdk.NewAttribute(types.AttributeKeySrcValidatorCmpPubKey, hex.EncodeToString(ev.ValidatorSrcCmpPubkey)),
				sdk.NewAttribute(types.AttributeKeyDstValidatorCmpPubKey, hex.EncodeToString(ev.ValidatorDstCmpPubkey)),
				sdk.NewAttribute(types.AttributeKeyDelegateID, ev.DelegationId.String()),
				sdk.NewAttribute(types.AttributeKeySenderAddress, ev.OperatorAddress.Hex()),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	isInSingularity, err := k.IsSingularity(cachedCtx)
	if err != nil {
		return errors.Wrap(err, "check if it is singularity")
	}

	if isInSingularity {
		log.Debug(cachedCtx, "Relegation event detected, but it is not processed since current block is singularity")
		return nil
	}

	validatorSrcPubkey, err := k1util.PubKeyBytesToCosmos(ev.ValidatorSrcCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "src validator pubkey to cosmos")
	}

	validatorDstPubkey, err := k1util.PubKeyBytesToCosmos(ev.ValidatorDstCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "dst validator pubkey to cosmos")
	}

	depositorAddr := sdk.AccAddress(ev.Delegator.Bytes())
	validatorSrcAddr := sdk.ValAddress(validatorSrcPubkey.Address().Bytes())
	validatorDstAddr := sdk.ValAddress(validatorDstPubkey.Address().Bytes())

	valSrcEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(validatorSrcPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "src validator pubkey to evm address")
	}
	valDstEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(validatorDstPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "dst validator pubkey to evm address")
	}

	// redelegateOnBehalf txn, need to check if it's from the operator
	if ev.Delegator.String() != ev.OperatorAddress.String() {
		operatorAddr, err := k.DelegatorOperatorAddress.Get(cachedCtx, depositorAddr.String())
		if errors.Is(err, collections.ErrNotFound) {
			return errors.WrapErrWithCode(
				errors.InvalidOperator,
				errors.New("invalid redelegateOnBehalf txn, no operator"),
			)
		} else if err != nil {
			return errors.Wrap(err, "get delegator's operator address failed")
		}
		if operatorAddr != ev.OperatorAddress.String() {
			return errors.WrapErrWithCode(
				errors.InvalidOperator,
				errors.New("invalid redelegateOnBehalf txn, not from operator"),
			)
		}
	}

	amountCoin, _ := IPTokenToBondCoin(ev.Amount)

	log.Debug(cachedCtx, "Processing EVM staking relegation",
		"del_story", depositorAddr.String(),
		"val_src_story", validatorSrcAddr.String(),
		"val_dst_story", validatorDstAddr.String(),
		"del_evm_addr", ev.Delegator.String(),
		"val_src_evm_addr", valSrcEvmAddr.String(),
		"val_dst_evm_addr", valDstEvmAddr.String(),
		"amount", amountCoin.Amount.String(),
	)

	msg := stypes.NewMsgBeginRedelegate(
		depositorAddr.String(), validatorSrcAddr.String(), validatorDstAddr.String(),
		ev.DelegationId.String(), amountCoin,
	)

	resp, err := skeeper.NewMsgServerImpl(k.stakingKeeper.(*skeeper.Keeper)).BeginRedelegate(cachedCtx, msg)
	switch {
	case errors.Is(err, stypes.ErrSelfRedelegation):
		return errors.WrapErrWithCode(errors.SelfRedelegation, err)
	case errors.Is(err, stypes.ErrNoValidatorFound):
		return errors.WrapErrWithCode(errors.ValidatorNotFound, err)
	case errors.Is(err, stypes.ErrTokenTypeMismatch):
		return errors.WrapErrWithCode(errors.TokenTypeMismatch, err)
	case errors.Is(err, stypes.ErrNoDelegation):
		return errors.WrapErrWithCode(errors.DelegationNotFound, err)
	case errors.Is(err, stypes.ErrNoPeriodDelegation):
		return errors.WrapErrWithCode(errors.PeriodDelegationNotFound, err)
	case err != nil:
		return errors.Wrap(err, "failed to begin redelegation")
	}

	actualAmount = resp.Amount.Amount.String()

	log.Debug(cachedCtx, "EVM staking relegation processed",
		"del_story", depositorAddr.String(),
		"val_src_story", validatorSrcAddr.String(),
		"val_dst_story", validatorDstAddr.String(),
		"del_evm_addr", ev.Delegator.String(),
		"val_src_evm_addr", valSrcEvmAddr.String(),
		"val_dst_evm_addr", valDstEvmAddr.String(),
		"actual_amount", actualAmount,
	)

	return nil
}

func (k Keeper) ParseRedelegateLog(ethLog ethtypes.Log) (*bindings.IPTokenStakingRedelegate, error) {
	return k.ipTokenStakingContract.ParseRedelegate(ethLog)
}
