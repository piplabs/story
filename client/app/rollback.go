package app

import (
	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/piplabs/story/lib/errors"
)

func RollbackCometAndAppState(a *App, cometCfg cmtcfg.Config, rollbackHeights int64, removeBlock bool) (lastHeight int64, lastHash []byte, err error) {
	// TODO: recovery when rollback fails mid-way
	for i := int64(0); i < rollbackHeights; i++ {
		lastHeight, lastHash, err = cmtcmd.RollbackState(&cometCfg, removeBlock)
		if err != nil {
			return lastHeight, lastHash, errors.Wrap(err, "failed to rollback CometBFT state")
		}
	}

	if err = a.CommitMultiStore().RollbackToVersion(lastHeight); err != nil {
		return 0, nil, errors.Wrap(err, "failed to rollback to version")
	}

	return lastHeight, lastHash, nil
}
