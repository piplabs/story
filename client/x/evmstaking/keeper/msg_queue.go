package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stype "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// EnqueueMsg enqueues a message to the queue of the current epoch.
func (k Keeper) EnqueueMsg(ctx context.Context, msg types.QueuedMessage) error {
	return k.MessageQueue.Enqueue(ctx, msg)
}

// ProcessAllMsgs returns the set of messages queued in a given epoch.
func (k Keeper) ProcessAllMsgs(ctx context.Context) error {
	iterator, err := k.MessageQueue.Iterate(ctx)
	if err != nil {
		return errors.Wrap(err, "message queue iterator")
	}

	for ; iterator.Valid(); iterator.Next() {
		qMsg, err := iterator.Value()
		if err != nil {
			return errors.Wrap(err, "get value of message queue")
		}

		if err := k.ProcessMsg(ctx, &qMsg); err != nil {
			log.Warn(ctx, "Failed to process queued message", err, "tx_id", string(qMsg.TxId))
			continue
		}
	}

	return nil
}

// ProcessMsg processes queues message depending on the type of message.
func (k Keeper) ProcessMsg(ctx context.Context, msg *types.QueuedMessage) error {
	var (
		unwrappedMsgWithType sdk.Msg
		err                  error
	)
	unwrappedMsgWithType, err = msg.UnwrapToSdkMsg()
	if err != nil {
		return errors.Wrap(err, "unwrap queued msg to sdk msg")
	}

	switch unwrappedMsg := unwrappedMsgWithType.(type) {
	case *stype.MsgCreateValidator:
		if err := k.ProcessCreateValidatorMsg(ctx, unwrappedMsg); err != nil {
			return errors.Wrap(err, "process CreateValidator msg")
		}
	case *stype.MsgDelegate:
		if err := k.ProcessDepositMsg(ctx, unwrappedMsg); err != nil {
			return errors.Wrap(err, "process Deposit msg")
		}
	case *stype.MsgBeginRedelegate:
		if err := k.ProcessRedelegateMsg(ctx, unwrappedMsg); err != nil {
			return errors.Wrap(err, "process BeginRedelegate msg")
		}
	case *stype.MsgUndelegate:
		if err := k.ProcessWithdrawMsg(ctx, unwrappedMsg); err != nil {
			return errors.Wrap(err, "process Withdraw msg")
		}
	case *slashingtypes.MsgUnjail:
		if err := k.ProcessUnjailMsg(ctx, unwrappedMsg); err != nil {
			return errors.Wrap(err, "process Unjail msg")
		}
	default:
		return errors.New("invalid type of queued message")
	}

	return nil
}
