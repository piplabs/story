package app

import (
	"context"
	"fmt"
	"slices"

	"cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"

	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	cmtcfg "github.com/cometbft/cometbft/config"

	"github.com/piplabs/story/client/config"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func RollbackCometAndAppState(a *App, cometCfg cmtcfg.Config, rollbackCfg config.RollbackConfig) (lastHeight int64, lastHash []byte, err error) {
	// Refuse before touching anything. RollbackState below mutates CometBFT's
	// state on each pass, so a check made afterwards would already be too late.
	if err := assertRollbackTargetIsLoadable(a.CommitMultiStore(), rollbackCfg.RollbackHeights); err != nil {
		return 0, nil, err
	}

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

// assertRollbackTargetIsLoadable refuses a rollback that would leave a state
// neither this binary nor the previous one can open.
//
// A version committed before a module was retired lists that module's store in
// its CommitInfo. RollbackToVersion rewrites that CommitInfo from the mounted
// stores alone, so rolling back across a retirement changes the app hash of the
// target version: this binary then fails replay with "state.AppHash does not
// match", and the binary that still mounts the store fails to load it with
// "version of store ... mismatch". Neither can start, and the only way back is
// a data directory backup or a snapshot.
//
// Rolling back to the retirement height itself is fine — its CommitInfo was
// already written without the store.
func assertRollbackTargetIsLoadable(cms storetypes.CommitMultiStore, heights uint64) error {
	rs, ok := cms.(*rootmulti.Store)
	if !ok {
		// Not the store type this check understands. The app wires a
		// *rootmulti.Store, so this only happens in a test or after a store
		// rewrite; say so rather than passing silently.
		log.Warn(context.Background(), "Skipping the rollback guard: unexpected multistore type",
			errors.New("unsupported multistore"), "type", fmt.Sprintf("%T", cms))

		return nil
	}

	tip := rs.LastCommitID().Version
	if tip <= 0 {
		return nil
	}

	// Refuse a count that cannot name a version. The loop below is not atomic:
	// each pass deletes a block and rewrites CometBFT's state, so letting an
	// out-of-range count through does not fail cleanly — it strips the node's
	// CometBFT state block by block until the rollback errors out, leaving
	// nothing to roll back to.
	if heights == 0 || heights >= uint64(tip) {
		return errors.New("rollback height count is out of range",
			"heights", heights,
			"current_height", tip,
			"hint", "pick a count below the current height; the rollback deletes a block per pass and cannot be undone",
		)
	}

	target := tip - int64(heights)

	cInfo, err := rs.GetCommitInfo(target)
	if err != nil {
		// Refuse rather than guess. Pruning does not remove commit info, so a
		// target that cannot be read is either out of range or corrupt — and a
		// corrupt one is exactly the case this guard exists for, because
		// RollbackToVersion would overwrite it from the mounted stores.
		return errors.Wrap(err, "cannot read the rollback target's commit info",
			"target_height", target,
		)
	}

	mounted := rs.StoreKeysByName()

	var retired []string

	for _, info := range cInfo.StoreInfos {
		if _, ok := mounted[info.Name]; !ok {
			retired = append(retired, info.Name)
		}
	}

	if len(retired) == 0 {
		return nil
	}

	slices.Sort(retired)

	return errors.New("refusing to roll back across a retired module store",
		"target_height", target,
		"retired_stores", fmt.Sprint(retired),
		"hint", "this binary cannot replay that height; roll back with the binary that still mounts that store and swap again, or restore a backup or state sync",
	)
}
