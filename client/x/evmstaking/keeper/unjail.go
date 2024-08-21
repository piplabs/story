package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/storyprotocol/iliad/contracts/bindings"
	"github.com/storyprotocol/iliad/lib/errors"
	"github.com/storyprotocol/iliad/lib/k1util"
)

func (k Keeper) ProcessUnjail(ctx context.Context, ev *bindings.IPTokenSlashingUnjail) error {
	validatorPubkey, err := k1util.PubKeyBytesToCosmos(ev.ValidatorCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to cosmos")
	}

	valAddr := sdk.ValAddress(validatorPubkey.Address().Bytes())
	err = k.slashingKeeper.Unjail(ctx, valAddr)
	if err != nil {
		return errors.Wrap(err, "unjail")
	}

	return nil
}

func (k Keeper) ParseUnjailLog(ethlog ethtypes.Log) (*bindings.IPTokenSlashingUnjail, error) {
	return k.ipTokenSlashingContract.ParseUnjail(ethlog)
}
