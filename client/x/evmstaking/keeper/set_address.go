package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
)

func (k Keeper) ProcessSetWithdrawalAddress(ctx context.Context, ev *bindings.IPTokenStakingSetWithdrawalAddress) error {
	depositorPubkey, err := k1util.PubKeyBytesToCosmos(ev.DelegatorCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "depositor pubkey to cosmos")
	}

	depositorAddr := sdk.AccAddress(depositorPubkey.Address().Bytes())
	executionAddr := common.BytesToAddress(ev.ExecutionAddress[:])

	if err := k.DelegatorMap.Set(ctx, depositorAddr.String(), executionAddr.String()); err != nil {
		return errors.Wrap(err, "delegator map set")
	}

	return nil
}
