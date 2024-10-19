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
	defer func() {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		if err != nil {
			sdkCtx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeUpdateValidatorCommissionFailure,
					sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
					sdk.NewAttribute(types.AttributeKeyValidatorUncmpPubKey, hex.EncodeToString(ev.ValidatorUncmpPubkey)),
					sdk.NewAttribute(types.AttributeKeyCommissionRate, strconv.FormatUint(uint64(ev.CommissionRate), 10)),
					sdk.NewAttribute(types.AttributeKeyStatusCode, types.UnwrapErrCode(err).String()),
				),
			})
		}
	}()

	valCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.ValidatorUncmpPubkey)
	if err != nil {
		return types.WrapErrWithCode(types.InvalidUncmpPubKey, errors.Wrap(err, "compress validator pubkey"))
	}
	validatorPubkey, err := k1util.PubKeyBytesToCosmos(valCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "compress validator pubkey")
	}

	validatorAddr := sdk.ValAddress(validatorPubkey.Address().Bytes())
	validator, err := k.stakingKeeper.GetValidator(ctx, validatorAddr)
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
	if _, err := skeeperMsgServer.EditValidator(ctx, msg); err != nil {
		return errors.Wrap(err, "update validator commission rate")
	}

	return nil
}
