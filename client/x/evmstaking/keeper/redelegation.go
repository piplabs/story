//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"encoding/hex"
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

	defer func() {
		if err == nil {
			writeCache()
			return
		}
		sdkCtx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeRedelegateFailure,
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDelegatorUncmpPubKey, hex.EncodeToString(ev.DelegatorUncmpPubkey)),
				sdk.NewAttribute(types.AttributeKeySrcValidatorUncmpPubKey, hex.EncodeToString(ev.ValidatorUncmpSrcPubkey)),
				sdk.NewAttribute(types.AttributeKeyDstValidatorUncmpPubKey, hex.EncodeToString(ev.ValidatorUncmpDstPubkey)),
				sdk.NewAttribute(types.AttributeKeyDelegateID, ev.DelegationId.String()),
				sdk.NewAttribute(types.AttributeKeyAmount, ev.Amount.String()),
				sdk.NewAttribute(types.AttributeKeySenderAddress, ev.OperatorAddress.Hex()),
				sdk.NewAttribute(types.AttributeKeyStatusCode, errors.UnwrapErrCode(err).String()),
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

	delCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.DelegatorUncmpPubkey)
	if err != nil {
		return errors.WrapErrWithCode(errors.InvalidUncmpPubKey, errors.Wrap(err, "compress delegator pubkey"))
	}
	depositorPubkey, err := k1util.PubKeyBytesToCosmos(delCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "depositor pubkey to cosmos")
	}

	valSrcCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.ValidatorUncmpSrcPubkey)
	if err != nil {
		return errors.WrapErrWithCode(errors.InvalidUncmpPubKey, errors.Wrap(err, "compress src validator pubkey"))
	}
	validatorSrcPubkey, err := k1util.PubKeyBytesToCosmos(valSrcCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "src validator pubkey to cosmos")
	}

	valDstCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.ValidatorUncmpDstPubkey)
	if err != nil {
		return errors.WrapErrWithCode(errors.InvalidUncmpPubKey, errors.Wrap(err, "compress dst validator pubkey"))
	}
	validatorDstPubkey, err := k1util.PubKeyBytesToCosmos(valDstCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "dst validator pubkey to cosmos")
	}

	depositorAddr := sdk.AccAddress(depositorPubkey.Address().Bytes())
	validatorSrcAddr := sdk.ValAddress(validatorSrcPubkey.Address().Bytes())
	validatorDstAddr := sdk.ValAddress(validatorDstPubkey.Address().Bytes())

	delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(depositorPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "deledator pubkey to evm address")
	}
	valSrcEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(validatorSrcPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "src validator pubkey to evm address")
	}
	valDstEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(validatorDstPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "dst validator pubkey to evm address")
	}

	// redelegateOnBehalf txn, need to check if it's from the operator
	if delEvmAddr.String() != ev.OperatorAddress.String() {
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

	log.Debug(cachedCtx, "EVM staking relegation detected",
		"del_story", depositorAddr.String(),
		"val_src_story", validatorSrcAddr.String(),
		"val_dst_story", validatorDstAddr.String(),
		"del_evm_addr", delEvmAddr.String(),
		"val_src_evm_addr", valSrcEvmAddr.String(),
		"val_dst_evm_addr", valDstEvmAddr.String(),
		"amount_coin", amountCoin.String(),
	)

	msg := stypes.NewMsgBeginRedelegate(
		depositorAddr.String(), validatorSrcAddr.String(), validatorDstAddr.String(),
		ev.DelegationId.String(), amountCoin,
	)
	_, err = skeeper.NewMsgServerImpl(k.stakingKeeper.(*skeeper.Keeper)).BeginRedelegate(cachedCtx, msg)
	if err != nil {
		return errors.Wrap(err, "failed to begin redelegation")
	}

	return nil
}

func (k Keeper) ParseRedelegateLog(ethLog ethtypes.Log) (*bindings.IPTokenStakingRedelegate, error) {
	return k.ipTokenStakingContract.ParseRedelegate(ethLog)
}
