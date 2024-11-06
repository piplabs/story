//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"encoding/hex"
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

func (k Keeper) ProcessUpdateValidatorCommission(ctx context.Context, ev *bindings.IPTokenStakingUpdateValidatorCommssion) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	defer func() {
		if err == nil {
			writeCache()
			return
		}
		sdkCtx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeUpdateValidatorCommissionFailure,
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyValidatorUncmpPubKey, hex.EncodeToString(ev.ValidatorUncmpPubkey)),
				sdk.NewAttribute(types.AttributeKeyCommissionRate, strconv.FormatUint(uint64(ev.CommissionRate), 10)),
				sdk.NewAttribute(types.AttributeKeyStatusCode, errors.UnwrapErrCode(err).String()),
			),
		})
	}()

	valCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.ValidatorUncmpPubkey)
	if err != nil {
		return errors.WrapErrWithCode(errors.InvalidUncmpPubKey, errors.Wrap(err, "compress validator pubkey"))
	}
	validatorPubkey, err := k1util.PubKeyBytesToCosmos(valCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to cosmos")
	}

	validatorAddr := sdk.ValAddress(validatorPubkey.Address().Bytes())
	validator, err := k.stakingKeeper.GetValidator(cachedCtx, validatorAddr)
	if err != nil {
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
