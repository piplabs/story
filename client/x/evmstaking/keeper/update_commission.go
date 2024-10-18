package keeper

import (
	"context"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
)

func (k Keeper) ProcessUpdateValidatorCommission(ctx context.Context, ev *bindings.IPTokenStakingUpdateValidatorCommssion) error {
	log.Info(ctx, "[Hans] Update Validator Comm detected")
	valCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.ValidatorUncmpPubkey)
	if err != nil {
		return errors.Wrap(err, "compress validator pubkey")
	}
	validatorPubkey, err := k1util.PubKeyBytesToCosmos(valCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to cosmos")
	}

	validatorAddr := sdk.ValAddress(validatorPubkey.Address().Bytes())
	validator, err := k.stakingKeeper.GetValidator(ctx, validatorAddr)
	if err != nil {
		return errors.Wrap(err, "get validator")
	}

	newComm := math.LegacyNewDec(int64(ev.CommissionRate))
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
