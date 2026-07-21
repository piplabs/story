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

	// completeSessionLocked runs under the round lock we hold here.
	if err := k.completeSessionLocked(ctx, session); err != nil {
		return
	}
}

// completeSessionLocked marks a finalized session completed and starts the decrypt worker.
// Caller MUST hold the per-round DKG service lock (this helper does not acquire it), so the
// recovery path can complete without releasing it. Marks the session failed on error.
func (k *Keeper) completeSessionLocked(ctx context.Context, session *types.DKGSession) error {
	session.UpdatePhase(types.PhaseCompleted)
	session.IsFinalized = true // ready to start DKG threshold encryption/decryption

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update completed session", err)
		k.stateManager.MarkFailed(ctx, session)

		return err
	}

	log.Info(ctx, "DKG process completed successfully",
		"round", session.Round,
		"validator_evm_address", k.validatorEVMAddr,
	)

	k.StartDecryptWorker()

	return nil
}

// recoverActiveSessionKeyMaterial recomputes and seals a keyless active-round session's share
// via the kernel and completes it, all under one continuous hold of the per-round lock (no
// release/re-acquire gap for the stuck-check to force-fail). No on-chain finalize vote (the
// round already finalized); on failure the session is left failed. Node-local only.
func (k *Keeper) recoverActiveSessionKeyMaterial(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	if !tryAcquireDKGSvc(dkgNetwork.Round) {
		log.Info(ctx, "DKG service already running for this round; skipping recovery",
			"round", dkgNetwork.Round,
		)

		return
	}
	defer releaseDKGSvc(dkgNetwork.Round)

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session for active-round recovery", err)

		return
	}

	// Cap check under the lock: reaching attempts >= cap here means the prior cap-th attempt
	// already concluded, so escalation is never premature and fires exactly once.
	attempts := session.GetRecoveryAttempts()
	if attempts >= maxActiveRecoveryAttempts {
		if attempts == maxActiveRecoveryAttempts {
			log.Error(ctx, "Active-round key recovery exhausted after max attempts; manual intervention required", nil,
				"round", dkgNetwork.Round,
				"recovery_attempts", attempts,
				"max_active_recovery_attempts", maxActiveRecoveryAttempts,
			)
			session.SetRecoveryAttempts(maxActiveRecoveryAttempts + 1)

			if err := k.stateManager.UpdateSession(ctx, session); err != nil {
				log.Error(ctx, "Failed to persist exhausted recovery marker", err)
			}
		}

		return
	}

	// Count the attempt only under the lock, so dispatches that lost the tryAcquireDKGSvc race
	// (no-ops) don't burn the retry budget.
	newAttempts := session.IncrementRecoveryAttempts()

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to persist recovery attempt counter", err)

		return
	}

	log.Warn(ctx, "Active-round session is missing key material; attempting local finalization recovery", nil,
		"round", dkgNetwork.Round,
		"recovery_attempt", newAttempts,
		"max_active_recovery_attempts", maxActiveRecoveryAttempts,
	)

	if err := k.callTEEFinalizeDKG(ctx, session); err != nil {
		log.Error(ctx, "Failed to recover key material for active round; leaving session failed", err,
			"round", dkgNetwork.Round,
		)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	if !session.HasKeyMaterial() {
		log.Error(ctx, "Kernel finalization returned no key material for active round; leaving session failed", nil,
			"round", dkgNetwork.Round,
		)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	// The recomputed key must match the on-chain network key, else this node would submit
	// partial decryptions under a divergent key.
	if !bytes.Equal(session.GetGlobalPubKey(), dkgNetwork.GlobalPublicKey) {
		// Divergent is deterministic (retry never helps), so jump straight to the sentinel to
		// stop re-entry; the Error below is this path's escalation alert.
		log.Error(ctx, "Recovered global public key does not match the on-chain network key; leaving session failed", nil,
			"round", dkgNetwork.Round,
		)
		session.SetRecoveryAttempts(maxActiveRecoveryAttempts + 1)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	// Success: finalize then complete inline under the same lock (no gap for the stuck-check).
	session.UpdatePhase(types.PhaseFinalized)

	if err := k.completeSessionLocked(ctx, session); err != nil {
		// completeSessionLocked already marked the session failed on its internal error.
		log.Error(ctx, "Failed to complete recovered session", err,
			"round", dkgNetwork.Round,
		)

		return
	}
}
