package keeper

import (
	"bytes"
	"context"
	"slices"
	"sort"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/genutil/evm/predeploys"
	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/errors"
)

// evmEvents returns selected EVM log events from the provided block hash.
func (k *Keeper) evmEvents(ctx context.Context, blockHash common.Hash) ([]*types.EVMEvent, error) {
	var events []*types.EVMEvent

	logs, err := k.engineCl.FilterLogs(ctx, ethereum.FilterQuery{
		BlockHash: &blockHash,
		// only IPTokenStaking contract
		Addresses: []common.Address{
			common.HexToAddress(predeploys.IPTokenStaking),
			common.HexToAddress(predeploys.IPTokenSlashing),
			common.HexToAddress(predeploys.UpgradeEntrypoint),
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "filter logs")
	}

	ll := make([]*types.EVMEvent, 0, len(logs))
	for _, l := range logs {
		topics := make([][]byte, 0, len(l.Topics))
		for _, t := range l.Topics {
			topics = append(topics, t.Bytes())
		}
		ll = append(ll, &types.EVMEvent{
			Address: l.Address.Bytes(),
			Topics:  topics,
			Data:    l.Data,
		})
	}

	for _, log := range ll {
		if err := log.Verify(); err != nil {
			return nil, errors.Wrap(err, "verify log")
		}
	}
	events = append(events, ll...)

	// Sort by Address > Topics > Data
	// This avoids dependency on runtime ordering.
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

	return events, nil
}
