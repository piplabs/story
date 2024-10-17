package keeper

import (
	"context"

	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
)

func (k Keeper) ProcessUnjail(ctx context.Context, ev *bindings.IPTokenStakingUnjail) error {
	valCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.ValidatorUncmpPubkey)
	if err != nil {
		return errors.Wrap(err, "compress validator pubkey")
	}
	validatorPubkey, err := k1util.PubKeyBytesToCosmos(valCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to cosmos")
	}

	valAddr := sdk.ValAddress(validatorPubkey.Address().Bytes())
	valDelAddr := sdk.AccAddress(validatorPubkey.Address().Bytes())

	valEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(valCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to evm address")
	}

	// unjailOnBehalf txn, need to check if it's from the operator
	if valEvmAddr.String() != ev.Unjailer.String() {
		operatorAddr, err := k.DelegatorOperatorAddress.Get(ctx, valDelAddr.String())
		if errors.Is(err, collections.ErrNotFound) {
			return errors.New("invalid unjailOnBehalf txn, no operator for delegator")
		} else if err != nil {
			return errors.Wrap(err, "get validator's operator address failed")
		}
		if operatorAddr != ev.Unjailer.String() {
			return errors.New("invalid unjailOnBehalf txn, not from operator")
		}
	}

	err = k.slashingKeeper.Unjail(ctx, valAddr)
	if err != nil {
		return errors.Wrap(err, "unjail")
	}

	return nil
}

func (k Keeper) ParseUnjailLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingUnjail, error) {
	return k.ipTokenStakingContract.ParseUnjail(ethlog)
}
