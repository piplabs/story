package app

import (
	"context"
	"math/big"

	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/lib/errors"
)

func RollbackCometAndAppState(ctx context.Context, a *App, appCfg Config, cometCfg cmtcfg.Config, rollbackHeights int64, removeBlock bool, rollbackEVM bool) (lastHeight int64, lastHash []byte, err error) {
	for range rollbackHeights {
		lastHeight, lastHash, err = cmtcmd.RollbackState(&cometCfg, removeBlock)
		if err != nil {
			return lastHeight, lastHash, errors.Wrap(err, "failed to rollback CometBFT state")
		}
	}

	if err = a.CommitMultiStore().RollbackToVersion(lastHeight); err != nil {
		return lastHeight, lastHash, errors.Wrap(err, "failed to rollback to version")
	}

	if !rollbackEVM || rollbackHeights <= 1 {
		return lastHeight, lastHash, nil
	}

	engineCl, err := newEngineClient(ctx, appCfg)
	if err != nil {
		return lastHeight, lastHash, err
	}

	// Note that EVM block might not match the consensus block height, so we need to calculate the rollback
	// EVM block height separately.
	latestHeight, err := engineCl.BlockNumber(ctx)
	if err != nil {
		return lastHeight, lastHash, err
	}

	latestBlock, err := engineCl.BlockByNumber(ctx, big.NewInt(int64(latestHeight)))
	if err != nil {
		return lastHeight, lastHash, err
	} else if latestBlock.BeaconRoot() == nil {
		return lastHeight, lastHash, errors.New("cannot rollback EVM with nil beacon root", "height", lastHeight)
	}

	// Rollback EVM if latest EVM block built on-top of new rolled-back consensus head.
	if *latestBlock.BeaconRoot() != common.BytesToHash(lastHash) {
		return lastHeight, lastHash, errors.New(
			"cannot rollback EVM, latest EVM block not built on new rolled-back state",
			"evm_height", latestHeight,
			"evm_beacon_root", *latestBlock.BeaconRoot(),
		)
	}

	rollbackEVMHeight := latestHeight - uint64(rollbackHeights)
	if err := engineCl.SetHead(ctx, rollbackEVMHeight); err != nil {
		return lastHeight, lastHash, errors.Wrap(err, "set head")
	}

	return lastHeight, lastHash, nil
}
