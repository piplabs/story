package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/log"
	"sync/atomic"
)

var dkgSvcRunning atomic.Bool

// ResumeDKGService reloads unfinished DKG sessions and resumes their execution safely without spawning duplicate goroutines.
func (k *Keeper) ResumeDKGService(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	session, err := k.stateManager.GetSession(dkgNetwork.Mrenclave, dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session while resuming the DKG service", err)

		return
	}

	if session.Phase != types.PhaseFailed {
		log.Debug(ctx, "No failed DKG session found; skipping resume process", "mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave), "round", dkgNetwork.Round)

		return
	}

	switch dkgNetwork.Stage {
	case types.DKGStageRegistration:
		session.UpdatePhase(types.PhaseInitializing)
		if err := k.stateManager.UpdateSession(ctx, session); err != nil {
			log.Error(ctx, "Failed to update session phase to initializing", err)

			return
		}

		go k.handleDKGInitialization(ctx, dkgNetwork)
	case types.DKGStageNetworkSet:
		session.UpdatePhase(types.PhaseInitialized)
		if err := k.stateManager.UpdateSession(ctx, session); err != nil {
			log.Error(ctx, "Failed to update session phase to initialized", err)

			return
		}

		go k.handleDKGNetworkSet(ctx, dkgNetwork)
	case types.DKGStageDealing:
		session.UpdatePhase(types.PhaseDealing)
		if err := k.stateManager.UpdateSession(ctx, session); err != nil {
			log.Error(ctx, "Failed to update session phase to dealing", err)

			return
		}

		go k.handleDKGDealing(ctx, dkgNetwork)
	case types.DKGStageFinalization:
		session.UpdatePhase(types.PhaseDealing)
		if err := k.stateManager.UpdateSession(ctx, session); err != nil {
			log.Error(ctx, "Failed to update session phase to dealing", err)

			return
		}

		go k.handleDKGFinalization(ctx, dkgNetwork)
	case types.DKGStageActive:
		session.UpdatePhase(types.PhaseFinalized)
		if err := k.stateManager.UpdateSession(ctx, session); err != nil {
			log.Error(ctx, "Failed to update session phase to finalized", err)

			return
		}

		go k.handleDKGComplete(ctx, dkgNetwork)
	case types.DKGStageUnspecified:
		return
	}

	return
}
