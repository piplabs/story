package keeper

import (
	"context"

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

	events = make([]*types.EVMEvent, 0, len(logs))
	for _, l := range logs {
		evmEvent, err := types.EthLogToEVMEvent(l)
		if err != nil {
			return nil, errors.Wrap(err, "convert log")
		}

		events = append(events, evmEvent)
	}

	// This avoids dependency on runtime ordering.
	types.SortEVMEvents(events)

	return events, nil
}
