package keeper

import (
	"context"
	"encoding/hex"
	"sync/atomic"
	"time"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

var dkgSvcRunning atomic.Bool
var decryptWorkerRunning atomic.Bool

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
}

// StartDecryptWorker launches a background loop (non-ABCI) that drains pending decrypt requests
// and performs TDH2 partial decrypts. Only one worker runs.
func (k *Keeper) StartDecryptWorker(ctx context.Context) {
	if !decryptWorkerRunning.CompareAndSwap(false, true) {
		// already running
		return
	}

	go func() {
		defer decryptWorkerRunning.Store(false)
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				k.processDecryptQueue(ctx)
			}
		}
	}()
}

// processDecryptQueue scans sessions for queued decrypt requests and starts TDH2 partial decrypt + submission.
func (k *Keeper) processDecryptQueue(ctx context.Context) {
	sessions := k.stateManager.ListSessions()
	for _, session := range sessions {
		if len(session.DecryptRequests) == 0 {
			continue
		}

		remaining := make([]types.DecryptRequest, 0, len(session.DecryptRequests))
		for _, req := range session.DecryptRequests {
			if err := k.handleDecryptRequest(ctx, session, req); err != nil {
				log.Error(ctx, "Failed to process decrypt request", err,
					"session", session.GetSessionKey(),
					"requester", req.Requester,
					"round", req.Round,
					"ciphertext_len", len(req.Ciphertext),
					"label_len", len(req.Label),
				)
				remaining = append(remaining, req) // keep for retry
				continue
			}
		}

		session.DecryptRequests = remaining
		if err := k.stateManager.UpdateSession(ctx, session); err != nil {
			log.Error(ctx, "Failed to update session after processing decrypt queue", err,
				"session", session.GetSessionKey(),
				"remaining_requests", len(session.DecryptRequests),
			)
		}
	}
}

// handleDecryptRequest attempts TDH2 partial decrypt for a single request.
func (k *Keeper) handleDecryptRequest(ctx context.Context, session *types.DKGSession, req types.DecryptRequest) error {
	if k.teeClient == nil {
		return errors.New("tee client not configured")
	}
	if k.contractClient == nil {
		return errors.New("contract client not configured")
	}

	pid := session.Index
	if pid == 0 {
		return errors.New("session index not set")
	}

	if len(session.GlobalPubKey) == 0 {
		return errors.New("missing GlobalPubKey for session")
	}

	log.Info(ctx, "Calling TEE PartialDecryptTDH2",
		"session", session.GetSessionKey(),
		"round", session.Round,
		"pid", pid,
		"ciphertext_len", len(req.Ciphertext),
		"label_len", len(req.Label),
	)

	teeCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	resp, err := k.teeClient.PartialDecryptTDH2(teeCtx, &types.PartialDecryptTDH2Request{
		Mrenclave:       session.Mrenclave,
		Round:           session.Round,
		Ciphertext:      req.Ciphertext,
		Label:           req.Label,
		DkgPubKey:       session.GlobalPubKey,
		RequesterPubKey: req.RequesterPubKey,
	})
	if err != nil {
		return errors.Wrap(err, "TEE partial decrypt failed")
	}

	if _, err := k.contractClient.SubmitPartialDecryption(
		ctx,
		session.Round,
		session.Mrenclave,
		resp.EncryptedPartialDecryption,
		resp.EphemeralPubKey,
		resp.PubShare,
		req.Label,
	); err != nil {
		return errors.Wrap(err, "failed to submit partial decryption")
	}

	return nil
}
