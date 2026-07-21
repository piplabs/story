package keeper

import (
	"bytes"
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

	// A round just activated, so bound session growth by pruning rounds well below it.
	// session.Round is the latest active round and is always retained.
	k.stateManager.PruneOldSessions(ctx, session.Round)

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

// recoverActiveSessionKeyMaterial recomputes and seals a keyless active-round session's
// share via the kernel (from persisted deals) and completes it. It deliberately does not
// submit an on-chain finalize vote (the round already finalized); on failure the session is
// left failed so the node is visibly broken rather than silently "completed".
func (k *Keeper) recoverActiveSessionKeyMaterial(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	// recoverKeyMaterial holds the per-round lock for the kernel finalize; it is
	// released before handleDKGComplete re-acquires it, avoiding self-deadlock.
	if !k.recoverKeyMaterial(ctx, dkgNetwork) {
		return
	}

	k.handleDKGComplete(ctx, dkgNetwork)
}

// recoverKeyMaterial re-runs node-local TEE finalization for a keyless active-round
// session under the per-round DKG service lock, returning true only if the session now
// holds key material matching the on-chain network key. Node-local only; no consensus state.
func (k *Keeper) recoverKeyMaterial(ctx context.Context, dkgNetwork *types.DKGNetwork) bool {
	if !tryAcquireDKGSvc(dkgNetwork.Round) {
		log.Info(ctx, "DKG service already running for this round; skipping recovery",
			"round", dkgNetwork.Round,
		)

		return false
	}
	defer releaseDKGSvc(dkgNetwork.Round)

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session for active-round recovery", err)

		return false
	}

	if err := k.callTEEFinalizeDKG(ctx, session); err != nil {
		log.Error(ctx, "Failed to recover key material for active round; leaving session failed", err,
			"round", dkgNetwork.Round,
		)
		k.stateManager.MarkFailed(ctx, session)

		return false
	}

	if len(session.GlobalPubKey) == 0 || len(session.PubKeyShare) == 0 {
		log.Error(ctx, "Kernel finalization returned no key material for active round; leaving session failed", nil,
			"round", dkgNetwork.Round,
		)
		k.stateManager.MarkFailed(ctx, session)

		return false
	}

	// The locally recomputed key must match the on-chain consensus network key;
	// otherwise this node would submit partial decryptions under a divergent key.
	if !bytes.Equal(session.GlobalPubKey, dkgNetwork.GlobalPublicKey) {
		log.Error(ctx, "Recovered global public key does not match the on-chain network key; leaving session failed", nil,
			"round", dkgNetwork.Round,
		)
		k.stateManager.MarkFailed(ctx, session)

		return false
	}

	session.UpdatePhase(types.PhaseFinalized)

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to persist recovered session before completion", err)
		k.stateManager.MarkFailed(ctx, session)

		return false
	}

	return true
}
