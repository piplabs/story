//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
)

func (k Keeper) ProcessUpdateValidatorCommission(ctx context.Context, ev *bindings.IPTokenStakingUpdateValidatorCommission) (err error) {
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
				types.EventTypeUpdateValidatorCommissionSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeUpdateValidatorCommissionFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyValidatorCmpPubKey, hex.EncodeToString(ev.ValidatorCmpPubkey)),
				sdk.NewAttribute(types.AttributeKeyCommissionRate, strconv.FormatUint(uint64(ev.CommissionRate), 10)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	validatorPubkey, err := k1util.PubKeyBytesToCosmos(ev.ValidatorCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to cosmos")
	}

	// derive validator EL address from given pubkey (as the contract already emits the pubkey)
	valEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(validatorPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "validator pubkey to evm address")
	}

	validatorAddr := sdk.ValAddress(valEvmAddr.Bytes())

	validator, err := k.stakingKeeper.GetValidator(cachedCtx, validatorAddr)
	if errors.Is(err, stypes.ErrNoValidatorFound) {
		return errors.WrapErrWithCode(errors.ValidatorNotFound, err)
	} else if err != nil {
		return errors.Wrap(err, "get validator")
	}

	newComm := math.LegacyNewDecWithPrec(int64(ev.CommissionRate), 4)
	minSelfDelegation := validator.MinSelfDelegation
	msg := stypes.NewMsgEditValidator(
		validatorAddr.String(),
		validator.Description,
		&newComm,
		&minSelfDelegation,
	)

	evmstakingSKeeper, ok := k.stakingKeeper.(*skeeper.Keeper)
	if !ok {
		return errors.New("type assertion failed")
	}

	skeeperMsgServer := skeeper.NewMsgServerImpl(evmstakingSKeeper)

	if _, err := skeeperMsgServer.EditValidator(cachedCtx, msg); err != nil {
		return errors.Wrap(err, "update validator commission rate")
	}

	return nil
}
