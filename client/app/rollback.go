package app

import (
	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	cmtcfg "github.com/cometbft/cometbft/config"

	"github.com/piplabs/story/client/config"
	"github.com/piplabs/story/lib/errors"
)

func RollbackCometAndAppState(a *App, cometCfg cmtcfg.Config, rollbackCfg config.RollbackConfig) (lastHeight int64, lastHash []byte, err error) {
	for range rollbackCfg.RollbackHeights {
		// setting removeBlock true to enable rollback multiple blocks by removing the block data
		lastHeight, lastHash, err = cmtcmd.RollbackState(&cometCfg, true)
		if err != nil {
			return lastHeight, lastHash, errors.Wrap(err, "failed to rollback CometBFT state")
		}
	}

	if err = a.CommitMultiStore().RollbackToVersion(lastHeight); err != nil {
		return lastHeight, lastHash, errors.Wrap(err, "failed to rollback to version")
	}

	return lastHeight, lastHash, nil
}
