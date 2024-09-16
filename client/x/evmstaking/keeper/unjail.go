package keeper

import (
	"context"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
)

// HandleUnjailEvent handles Unjail event. It converts event to sdk.Msg and enqueues for epoched staking.
func (k Keeper) HandleUnjailEvent(ctx context.Context, ev *bindings.IPTokenSlashingUnjail) error {
	validatorPubkey, err := k1util.PubKeyBytesToCosmos(ev.ValidatorCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to cosmos")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	header := sdkCtx.BlockHeader()
	txID := tmhash.Sum(sdkCtx.TxBytes())

	msg := stypes.NewMsgUnjail(sdk.ValAddress(validatorPubkey.Address().Bytes()).String())
	qMsg, err := types.NewQueuedMessage(uint64(header.Height), header.Time, txID, msg)
	if err != nil {
		return errors.Wrap(err, "new queued message for Unjail event")
	}

	if err := k.EnqueueMsg(ctx, qMsg); err != nil {
		return errors.Wrap(err, "enqueue Unjail message")
	}

	return nil
}

// ProcessUnjailMsg processes the Unjail message.
func (k Keeper) ProcessUnjailMsg(ctx context.Context, msg *stypes.MsgUnjail) error {
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddr)
	if err != nil {
		return errors.Wrap(err, "val address from bech32", "validator_addr", msg.ValidatorAddr)
	}

	if err := k.slashingKeeper.Unjail(ctx, valAddr); err != nil {
		return errors.Wrap(err, "unjail")
	}

	return nil
}

func (k Keeper) ParseUnjailLog(ethlog ethtypes.Log) (*bindings.IPTokenSlashingUnjail, error) {
	return k.ipTokenSlashingContract.ParseUnjail(ethlog)
}
