package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/storyprotocol/iliad/contracts/bindings"
	"github.com/storyprotocol/iliad/lib/errors"
	"github.com/storyprotocol/iliad/lib/k1util"
)

func (k Keeper) ProcessSetWithdrawalAddress(ctx context.Context, ev *bindings.IPTokenStakingSetWithdrawalAddress) error {
	depositorPubkey, err := k1util.PubKeyBytesToCosmos(ev.DepositorPubkey)
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
