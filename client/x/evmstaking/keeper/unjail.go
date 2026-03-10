//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
)

func (k Keeper) ProcessUnjail(ctx context.Context, ev *bindings.IPTokenStakingUnjail) (err error) {
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
				types.EventTypeUnjailSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeUnjailFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyValidatorCmpPubKey, hex.EncodeToString(ev.ValidatorCmpPubkey)),
				sdk.NewAttribute(types.AttributeKeySenderAddress, ev.Unjailer.Hex()),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	validatorPubkey, err := k1util.PubKeyBytesToCosmos(ev.ValidatorCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to cosmos")
	}

	valEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(validatorPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "validator pubkey to evm address")
	}

	valAddr := sdk.ValAddress(valEvmAddr.Bytes())
	valDelAddr := sdk.AccAddress(valEvmAddr.Bytes())

	// unjailOnBehalf txn, need to check if it's from the operator
	if valEvmAddr.String() != ev.Unjailer.String() {
		operatorAddr, err := k.DelegatorOperatorAddress.Get(cachedCtx, valDelAddr.String())
		if errors.Is(err, collections.ErrNotFound) {
			return errors.WrapErrWithCode(
				errors.InvalidOperator,
				errors.New("invalid unjailOnBehalf txn, no operator for delegator"),
			)
		} else if err != nil {
			return errors.Wrap(err, "get validator's operator address failed")
		}

		if operatorAddr != ev.Unjailer.String() {
			return errors.WrapErrWithCode(
				errors.InvalidOperator,
				errors.New("invalid unjailOnBehalf txn, not from operator"),
			)
		}
	}

	log.Debug(cachedCtx, "EVM unjail detected, unjail validator",
		"val_story", valAddr.String(),
		"unjailer_evm_addr", ev.Unjailer.String(),
	)

	if err = k.slashingKeeper.Unjail(cachedCtx, valAddr); errors.Is(err, slashtypes.ErrNoValidatorForAddress) {
		return errors.WrapErrWithCode(errors.ValidatorNotFound, err)
	} else if errors.Is(err, slashtypes.ErrMissingSelfDelegation) {
		return errors.WrapErrWithCode(errors.MissingSelfDelegation, err)
	} else if errors.Is(err, slashtypes.ErrValidatorNotJailed) {
		return errors.WrapErrWithCode(errors.ValidatorNotJailed, err)
	} else if errors.Is(err, slashtypes.ErrValidatorJailed) {
		return errors.WrapErrWithCode(errors.ValidatorStillJailed, err)
	} else if err != nil {
		return errors.Wrap(err, "unjail")
	}

	return nil
}

func (k Keeper) ParseUnjailLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingUnjail, error) {
	return k.ipTokenStakingContract.ParseUnjail(ethlog)
}
