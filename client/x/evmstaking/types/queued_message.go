package types

import (
	"time"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/piplabs/story/lib/errors"
)

// NewQueuedMessage creates a new QueuedMessage from a wrapped msg
// i.e., wrapped -> unwrapped -> QueuedMessage.
func NewQueuedMessage(blockHeight uint64, blockTime time.Time, txid []byte, msg sdk.Msg) (QueuedMessage, error) {
	// marshal the actual msg (MsgDelegate, MsgBeginRedelegate, MsgUndelegate, MsgCancelUnbondingDelegation) inside isQueuedMessage_Msg
	var qmsg isQueuedMessage_Msg
	var msgBytes []byte
	var err error
	switch msgWithType := msg.(type) {
	case *types.MsgCreateValidator:
		if msgBytes, err = msgWithType.Marshal(); err != nil {
			return QueuedMessage{}, errors.Wrap(err, "marshal CreateValidator msg")
		}
		qmsg = &QueuedMessage_MsgCreateValidator{
			MsgCreateValidator: msgWithType,
		}
	case *types.MsgDelegate:
		if msgBytes, err = msgWithType.Marshal(); err != nil {
			return QueuedMessage{}, errors.Wrap(err, "marshal Delegate msg")
		}
		qmsg = &QueuedMessage_MsgDelegate{
			MsgDelegate: msgWithType,
		}
	case *types.MsgBeginRedelegate:
		if msgBytes, err = msgWithType.Marshal(); err != nil {
			return QueuedMessage{}, errors.Wrap(err, "marshal BeginRedelegate msg")
		}
		qmsg = &QueuedMessage_MsgBeginRedelegate{
			MsgBeginRedelegate: msgWithType,
		}
	case *types.MsgUndelegate:
		if msgBytes, err = msgWithType.Marshal(); err != nil {
			return QueuedMessage{}, errors.Wrap(err, "marshal Undelegate msg")
		}
		qmsg = &QueuedMessage_MsgUndelegate{
			MsgUndelegate: msgWithType,
		}
	case *stypes.MsgUnjail:
		if msgBytes, err = msgWithType.Marshal(); err != nil {
			return QueuedMessage{}, errors.Wrap(err, "marshal Unjail msg")
		}
		qmsg = &QueuedMessage_MsgUnjail{
			MsgUnjail: msgWithType,
		}
	default:
		return QueuedMessage{}, errors.New("invalid message type for queued message")
	}

	queuedMsg := QueuedMessage{
		TxId:        txid,
		MsgId:       tmhash.Sum(msgBytes),
		BlockHeight: blockHeight,
		BlockTime:   &blockTime,
		Msg:         qmsg,
	}

	return queuedMsg, nil
}

func (qm *QueuedMessage) UnwrapToSdkMsg() (sdk.Msg, error) {
	var unwrappedMsgWithType sdk.Msg
	switch unwrappedMsg := qm.Msg.(type) {
	case *QueuedMessage_MsgCreateValidator:
		unwrappedMsgWithType = unwrappedMsg.MsgCreateValidator
	case *QueuedMessage_MsgDelegate:
		unwrappedMsgWithType = unwrappedMsg.MsgDelegate
	case *QueuedMessage_MsgUndelegate:
		unwrappedMsgWithType = unwrappedMsg.MsgUndelegate
	case *QueuedMessage_MsgBeginRedelegate:
		unwrappedMsgWithType = unwrappedMsg.MsgBeginRedelegate
	case *QueuedMessage_MsgUnjail:
		unwrappedMsgWithType = unwrappedMsg.MsgUnjail
	default:
		return nil, errors.New("invalid message type for unwrapping")
	}

	return unwrappedMsgWithType, nil
}
