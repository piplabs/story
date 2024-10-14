package types

import (
	"bytes"
	"slices"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/lib/errors"
)

// ToEthLog converts an EVMEvent to an Ethereum Log.
// Note it assumes that Verify has been called before.
func (l *EVMEvent) ToEthLog() ethtypes.Log {
	topics := make([]common.Hash, 0, len(l.Topics))
	for _, t := range l.Topics {
		topics = append(topics, common.BytesToHash(t))
	}

	return ethtypes.Log{
		Address: common.BytesToAddress(l.Address),
		Topics:  topics,
		Data:    l.Data,
	}
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

	if len(l.Address) != len(common.Address{}) {
		return errors.New("invalid address length")
	}

	for _, t := range l.Topics {
		if len(t) != len(common.Hash{}) {
			return errors.New("invalid topic length")
		}
	}

	return nil
}

// EthLogToEVMEvent converts an Ethereum Log to an EVMEvent.
func EthLogToEVMEvent(ethLog ethtypes.Log) (*EVMEvent, error) {
	topics := make([][]byte, 0, len(ethLog.Topics))
	for _, t := range ethLog.Topics {
		topics = append(topics, t.Bytes())
	}

	evmEvent := &EVMEvent{
		Address: ethLog.Address.Bytes(),
		Topics:  topics,
		Data:    ethLog.Data,
	}
	if err := evmEvent.Verify(); err != nil {
		return nil, errors.Wrap(err, "verify log")
	}

	return evmEvent, nil
}

// SortEVMEvents sorts EVM events by Address > Topics > Data.
func SortEVMEvents(events []*EVMEvent) {
	sort.Slice(events, func(i, j int) bool {
		if cmp := bytes.Compare(events[i].Address, events[j].Address); cmp != 0 {
			return cmp < 0
		}

		topicI := slices.Concat(events[i].Topics...)
		topicJ := slices.Concat(events[j].Topics...)
		if cmp := bytes.Compare(topicI, topicJ); cmp != 0 {
			return cmp < 0
		}

		return bytes.Compare(events[i].Data, events[j].Data) < 0
	})
}

// IsSortedEVMEvents check if the events are sorted by ascending order of address, topics, and data.
func IsSortedEVMEvents(events []*EVMEvent) bool {
	for i := 1; i < len(events); i++ {
		// Compare addresses first
		addressComparison := bytes.Compare(events[i-1].Address, events[i].Address)
		if addressComparison > 0 {
			// it is not sorted by ascending order of address
			return false
		}

		if addressComparison == 0 {
			// If addresses are equal, compare by topics
			previousTopic := slices.Concat(events[i-1].Topics...)
			currentTopic := slices.Concat(events[i].Topics...)
			topicComparison := bytes.Compare(previousTopic, currentTopic)

			if topicComparison > 0 {
				// it is not sorted by ascending order of topics
				return false
			}

			if topicComparison == 0 {
				// If topics are also equal, compare by data
				if bytes.Compare(events[i-1].Data, events[i].Data) > 0 {
					// it is not sorted by ascending order of data
					return false
				}
			}
		}
	}

	return true
}
