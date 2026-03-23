package keeper

import (
	"context"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/log"
)

// handleDKGComplete handles the DKG completion event.
func (k *Keeper) handleDKGComplete(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG completion",
		"round", dkgNetwork.Round,
	)

	if !tryAcquireDKGSvc(dkgNetwork.Round) {
		log.Info(ctx, "DKG service already running for this round; skipping completion",
			"round", dkgNetwork.Round,
		)

		return
	}
	defer releaseDKGSvc(dkgNetwork.Round)

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session", err)

		return
	}

	if session.Phase == types.PhaseCompleted && session.IsFinalized {
		log.Info(ctx, "DKG network already completed")
		// Ensure the decrypt worker is running even if completion was already processed
		// (e.g., after node restart or if the worker exited due to a transient error).
		k.StartDecryptWorker()

		return
	}

	if session.Phase != types.PhaseFinalized {
		log.Error(ctx, "Session is not finalized yet", nil)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	session.UpdatePhase(types.PhaseCompleted)
	session.IsFinalized = true // ready to start DKG threshold encryption/decryption

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update completed session", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	log.Info(ctx, "DKG process completed successfully",
		"round", session.Round,
		"validator_evm_address", k.validatorEVMAddr,
	)

	k.StartDecryptWorker()
}
