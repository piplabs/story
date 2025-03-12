//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/math"

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

func (k Keeper) ProcessDeposit(ctx context.Context, ev *bindings.IPTokenStakingDeposit) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	var isRefund bool
	var refundAmount math.Int
	var completionTime time.Time

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeDelegateSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDelegateFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		e = e.AppendAttributes(
			sdk.NewAttribute(types.AttributeKeyAmount, ev.StakeAmount.String()),
			sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
			sdk.NewAttribute(types.AttributeKeyDelegatorAddress, ev.Delegator.String()),
			sdk.NewAttribute(types.AttributeKeyValidatorCmpPubKey, hex.EncodeToString(ev.ValidatorCmpPubkey)),
			sdk.NewAttribute(types.AttributeKeyDelegateID, ev.DelegationId.String()),
			sdk.NewAttribute(types.AttributeKeyPeriodType, strconv.FormatInt(ev.StakingPeriod.Int64(), 10)),
			sdk.NewAttribute(types.AttributeKeySenderAddress, ev.OperatorAddress.Hex()),
			sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			sdk.NewAttribute(types.AttributeKeyIsRefund, strconv.FormatBool(isRefund)),
		)

		// if it's a refund, add refund attributes
		if isRefund {
			e = e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyRefundAmount, refundAmount.String()),
				sdk.NewAttribute(types.AttributeKeyRefundCompletionTime, completionTime.String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{e})
	}()

	validatorPubkey, err := k1util.PubKeyBytesToCosmos(ev.ValidatorCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to cosmos")
	}

	valEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(validatorPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "validator pubkey to evm address")
	}

	depositorAddr := sdk.AccAddress(ev.Delegator.Bytes())
	validatorAddr := sdk.ValAddress(valEvmAddr.Bytes())

	amountCoin, amountCoins := IPTokenToBondCoin(ev.StakeAmount)

	// Create account if not exists
	if !k.authKeeper.HasAccount(cachedCtx, depositorAddr) {
		acc := k.authKeeper.NewAccountWithAddress(cachedCtx, depositorAddr)
		k.authKeeper.SetAccount(cachedCtx, acc)
		log.Debug(cachedCtx, "Created account for depositor",
			"del_story_addr", depositorAddr.String(),
			"del_evm_addr", ev.Delegator.String(),
			"operator_evm_addr", ev.OperatorAddress.String(),
		)
	}

	log.Debug(cachedCtx, "EVM staking deposit detected, delegating to validator",
		"del_story", depositorAddr.String(),
		"val_story", validatorAddr.String(),
		"del_evm_addr", ev.Delegator.String(),
		"val_evm_addr", valEvmAddr.String(),
		"operator_evm_addr", ev.OperatorAddress.String(),
		"amount", amountCoin.Amount.String(),
	)

	lockedTokenType, err := k.stakingKeeper.GetLockedTokenType(cachedCtx)
	if err != nil {
		return errors.Wrap(err, "get locked token type")
	}

	val, err := k.stakingKeeper.GetValidator(cachedCtx, validatorAddr)

	//nolint:nestif // nested ifs error handling
	if errors.Is(err, stypes.ErrNoValidatorFound) {
		isRefund = true
		log.Info(cachedCtx, "Validator not found, refunding deposit minus the refund fee",
			"val_addr", validatorAddr,
			"val_evm_addr", valEvmAddr.String(),
			"del_addr", depositorAddr.String(),
			"del_evm_addr", ev.Delegator.String(),
		)

		// the min refund fee amount will be `refundFeeBps * minDelegationAmount (1024) / 10_000bps`
		refundFeeBps, err := k.RefundFeeBps(cachedCtx)
		if err != nil {
			return errors.Wrap(err, "get refund fee")
		}
		refundFeeAmount := amountCoin.Amount.Mul(math.NewInt(int64(refundFeeBps))).Quo(math.NewInt(10_000))
		refundAmount = amountCoin.Amount.Sub(refundFeeAmount)

		refundPeriod, err := k.RefundPeriod(cachedCtx)
		if err != nil {
			return errors.Wrap(err, "get refund period")
		}

		completionTime = sdkCtx.BlockTime().Add(refundPeriod)

		// push the refund to the unbonding queue
		ubd, err := k.stakingKeeper.SetUnbondingDelegationEntry(
			cachedCtx,
			depositorAddr,
			validatorAddr,
			sdkCtx.BlockHeight(),
			completionTime,
			refundAmount,
		)
		if err != nil {
			return errors.Wrap(err, "set unbonding delegation entry")
		}

		err = k.stakingKeeper.InsertUBDQueue(cachedCtx, ubd, completionTime)
		if err != nil {
			return errors.Wrap(err, "insert unbonding delegation queue")
		}

		log.Debug(cachedCtx, "Added refund to unbonding queue",
			"del_addr", depositorAddr.String(),
			"val_addr", validatorAddr.String(),
			"amount", refundAmount.String(),
			"completion_time", completionTime,
		)

		return nil // skip delegation logic
	} else if err != nil {
		return errors.Wrap(err, "get validator failed")
	}

	// locked tokens can only be staked with flexible period,
	// here we automatically set the period type to flexible period
	delID := ev.DelegationId.String()
	periodType := int32(ev.StakingPeriod.Int64())
	if val.SupportTokenType == lockedTokenType {
		flexPeriodType, err := k.stakingKeeper.GetFlexiblePeriodType(cachedCtx)
		if err != nil {
			return errors.Wrap(err, "get flexible period type")
		}
		periodType = flexPeriodType
		delID = stypes.FlexiblePeriodDelegationID
	}

	// Note that, after minting, we save the mapping between delegator bech32 address and evm address, which will be used in the withdrawal queue.
	// The saving is done regardless of any error below, as the money is already minted and sent to the delegator, who can withdraw the minted amount.
	// NOTE: Do not overwrite the existing withdraw/reward address set by the delegator.
	if exists, err := k.DelegatorWithdrawAddress.Has(cachedCtx, depositorAddr.String()); err != nil {
		return errors.Wrap(err, "check delegator withdraw address existence")
	} else if !exists {
		if err := k.DelegatorWithdrawAddress.Set(cachedCtx, depositorAddr.String(), ev.Delegator.String()); err != nil {
			return errors.Wrap(err, "set delegator withdraw address map")
		}
	}
	if exists, err := k.DelegatorRewardAddress.Has(cachedCtx, depositorAddr.String()); err != nil {
		return errors.Wrap(err, "check delegator reward address existence")
	} else if !exists {
		if err := k.DelegatorRewardAddress.Set(cachedCtx, depositorAddr.String(), ev.Delegator.String()); err != nil {
			return errors.Wrap(err, "set delegator reward address map")
		}
	}

	if err := k.bankKeeper.MintCoins(cachedCtx, types.ModuleName, amountCoins); err != nil {
		return errors.Wrap(err, "create stake coin for depositor: mint coins")
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(cachedCtx, types.ModuleName, depositorAddr, amountCoins); err != nil {
		return errors.Wrap(err, "create stake coin for depositor: send coins")
	}

	return k.CreateDelegation(cachedCtx, validatorAddr.String(), depositorAddr.String(), amountCoin, delID, periodType)
}

func (k Keeper) CreateDelegation(
	cachedCtx context.Context, validatorAddr, depositorAddr string, amountCoin sdk.Coin, periodDelegationID string,
	periodType int32,
) error {
	evmstakingSKeeper, ok := k.stakingKeeper.(*skeeper.Keeper)
	if !ok {
		return errors.New("type assertion failed")
	}

	skeeperMsgServer := skeeper.NewMsgServerImpl(evmstakingSKeeper)
	// Delegation by the depositor on the validator (validator existence is checked in msgServer.Delegate)
	msg := stypes.NewMsgDelegate(
		depositorAddr, validatorAddr, amountCoin,
		periodDelegationID, periodType,
	)
	if _, err := skeeperMsgServer.Delegate(cachedCtx, msg); errors.Is(err, stypes.ErrDelegationBelowMinimum) {
		return errors.WrapErrWithCode(errors.InvalidDelegationAmount, err)
	} else if errors.Is(err, stypes.ErrNoPeriodTypeFound) {
		return errors.WrapErrWithCode(errors.InvalidPeriodType, err)
	} else if err != nil {
		return errors.Wrap(err, "delegate")
	}

	return nil
}

func (k Keeper) ParseDepositLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingDeposit, error) {
	return k.ipTokenStakingContract.ParseDeposit(ethlog)
}
