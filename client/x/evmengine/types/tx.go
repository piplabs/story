package types

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/lib/cast"
	"github.com/piplabs/story/lib/errors"
)

/*
ToEthLog converts an EVMEvent to an Ethereum Log.
Note:
1. It assumes that Verify has been called before.
2. It's only used for events from the story contracts.
*/
func (l *EVMEvent) ToEthLog() (ethtypes.Log, error) {
	if l == nil {
		return ethtypes.Log{}, errors.New("nil log")
	} else if len(l.Topics) == 0 {
		// Events of story contracts are expected to have topics.
		return ethtypes.Log{}, errors.New("empty topics")
	}
	topics := make([]common.Hash, 0, len(l.Topics))
	for _, t := range l.Topics {
		topics = append(topics, common.BytesToHash(t))
	}

	addr, err := cast.EthAddress(l.Address)
	if err != nil {
		return ethtypes.Log{}, err
	}

	return ethtypes.Log{
		Address: addr,
		Topics:  topics,
		Data:    l.Data,
		TxHash:  common.BytesToHash(l.TxHash),
	}, nil
}

func (l *EVMEvent) Verify() error {
	if l == nil {
		return errors.New("nil log")
	}

	if l.Address == nil {
		return errors.New("nil address")
	}

	if len(l.Topics) == 0 {
		return errors.New("empty topics")
	}

	if l.TxHash == nil {
		return errors.New("nil tx hash")
	}

	if len(l.Address) != len(common.Address{}) {
		return errors.New("invalid address length")
	}

	for _, t := range l.Topics {
		if len(t) != len(common.Hash{}) {
			return errors.New("invalid topic length")
		}
	}

	if len(l.TxHash) != len(common.Hash{}) {
		return errors.New("invalid tx hash length")
	}

	return nil
}
