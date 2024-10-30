package keeper

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/genutil/evm/predeploys"
	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/errors"
	clog "github.com/piplabs/story/lib/log"
)

// evmEvents returns selected EVM log events from the provided block hash.
func (k *Keeper) evmEvents(ctx context.Context, blockHash common.Hash) ([]*types.EVMEvent, error) {
	var logs []ethtypes.Log
	err := retryForever(ctx, func(ctx context.Context) (fetched bool, err error) {
		logs, err = k.engineCl.FilterLogs(ctx, ethereum.FilterQuery{
			BlockHash: &blockHash,
			Addresses: []common.Address{
				common.HexToAddress(predeploys.IPTokenStaking),
				common.HexToAddress(predeploys.UBIPool),
				common.HexToAddress(predeploys.UpgradeEntrypoint),
			},
		})
		if err != nil {
			clog.Warn(ctx, "Failed fetching evm events (will retry)", err)

			return false, nil // Retry
		}

		return true, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "filter logs")
	}

	events := make([]*types.EVMEvent, 0, len(logs))
	for _, l := range logs {
		topics := make([][]byte, 0, len(l.Topics))
		for _, t := range l.Topics {
			topics = append(topics, t.Bytes())
		}
		events = append(events, &types.EVMEvent{
			Address: l.Address.Bytes(),
			Topics:  topics,
			Data:    l.Data,
		})
	}

	for _, event := range events {
		if err := event.Verify(); err != nil {
			return nil, errors.Wrap(err, "verify event")
		}
	}

	return events, nil
}
