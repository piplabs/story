package keeper

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
	"time"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// decryptComputeResult holds the outcome of a single parallel kernel partial-decrypt call.
type decryptComputeResult struct {
	idx  int // index into the original requests slice for retry tracking
	req  types.DecryptRequest
	resp *types.PartialDecryptTDH2Response
	err  error
}

// decryptChanBuf limits the number of goroutines running kernel calls concurrently.
const decryptChanBuf = 10

// defaultDecryptBatchSize is the default number of partial decryptions to submit
// in a single batch contract call when no explicit batch size is configured.
// The CDR contract enforces its own maxBatchSize; this client-side value should be ≤ that limit.
// One single submission tx consumes nearly 72k gas, including fixed overhead and per-decryption cost.
// Batching reduces total gas by amortizing the fixed overhead across multiple decryptions.
// With 20 decryptions, we get ~35% gas savings compared to single submissions, while keeping batch size manageable.
const defaultDecryptBatchSize = 20

// dkgSvcRound tracks which DKG round is currently being processed by async
// goroutines. Zero means no round is running. A higher round number always
// supersedes a stale lock from a previous round, preventing cross-round
// blocking where a slow/retrying goroutine from round N blocks round N+1.
var dkgSvcRound atomic.Uint64

// tryAcquireDKGSvc attempts to acquire the DKG service lock for the given round.
// Returns true if acquired. Duplicate calls for the same round return false
// (deduplication). A newer (higher) round always preempts an older one.
func tryAcquireDKGSvc(round uint32) bool {
	const maxCASRetries = 10

	r := uint64(round)
	for range maxCASRetries {
		current := dkgSvcRound.Load()
		if current == r {
			return false // same round already running, deduplicate
		}

		if current == 0 || r > current {
			if dkgSvcRound.CompareAndSwap(current, r) {
				if current > 0 {
					log.Info(context.Background(), "Superseded stale DKG lock from previous round",
						"old_round", current,
						"new_round", r,
					)
				}

				return true
			}

			continue // CAS failed, retry
		}

		return false // round going backwards, skip
	}

	return false // max retries exhausted
}

// releaseDKGSvc releases the lock only if this round still holds it.
// If a newer round has superseded, this is a no-op.
func releaseDKGSvc(round uint32) {
	dkgSvcRound.CompareAndSwap(uint64(round), 0)
}

// dkgAsyncTimeout is the maximum duration for async DKG goroutines that
// communicate with the story-kernel. These goroutines must NOT use the CometBFT
// consensus context because it gets canceled when block processing completes,
// which can abort in-flight gRPC calls to the story-kernel.
//
// The budget must cover a TEE kernel call plus one waitForTransaction call
// (which has its own 60s inner timeout). 2 minutes gives ~60s for the kernel
// operation and ~60s for the on-chain transaction to be mined. If the async
// context expires before waitForTransaction's inner timeout fires, the tx wait
// is cut short — so this value must stay above 60s to be meaningful.
const dkgAsyncTimeout = 2 * time.Minute

// decryptKernelCallTimeout bounds a single kernel RPC made by the decrypt worker.
// The worker runs on a process-lifetime context.Background with no deadline of its
// own, so without a per-call timeout a wedged kernel (accepts the connection but
// never responds) blocks the worker loop forever and decryptWorkerRunning is never
// reset. 60s matches the kernel-operation budget in dkgAsyncTimeout and the
// contract-call timeout in contract_client.go.
const decryptKernelCallTimeout = 60 * time.Second

// dkgAsyncContext creates a new context for async DKG service goroutines with a timeout.
// This replaces the consensus context that would otherwise be canceled after block processing.
func dkgAsyncContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), dkgAsyncTimeout)
}

var decryptWorkerRunning atomic.Bool

// ResumeDKGService reloads unfinished DKG sessions and resumes their execution safely without spawning duplicate goroutines.
func (k *Keeper) ResumeDKGService(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	// The decrypt worker serves the latest ACTIVE round, which may differ from the
	// latest round once the next round has opened. StartDecryptWorker is idempotent.
	if active, err := k.getLatestActiveDKGNetwork(ctx); err != nil {
		log.Warn(ctx, "Failed to get latest active DKG round while resuming decrypt worker", err)
	} else if active != nil {
		k.StartDecryptWorker()
	} else {
		// Debug (not Info): this runs every block; a missing active round is the
		// normal steady state and Info would be too noisy.
		log.Debug(ctx, "No active DKG round; decrypt worker not started")
	}

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session while resuming the DKG service", err)

		return
	}

	// If the session is completed and the DKG round is active, ensure the decrypt
	// worker is running. This covers node restarts and the case where the worker
	// was never started due to context cancellation.
	if session.Phase == types.PhaseCompleted && dkgNetwork.Stage == types.DKGStageActive {
		k.StartDecryptWorker()

		return
	}

	if session.Phase == types.PhaseFailed {
		k.resumeFailedSession(ctx, session, dkgNetwork)

		return
	}

	// Recover sessions stuck in intermediate phases after an unexpected process crash.
	// A session is stuck when its phase has NOT advanced to the expected completion
	// state for the current stage AND no goroutine is actively processing it.
	//
	// Valid (non-stuck) completion states per stage:
	//   Registration → PhaseInitialized (registration done, waiting for dealing)
	//   Dealing      → PhaseDealing     (deals generated, processing responses via VE)
	//   Finalization → PhaseFinalized   (finalization done, waiting for active)
	//   Active       → PhaseCompleted
	if isSessionStuckForStage(session.Phase, dkgNetwork.Stage) {
		if tryAcquireDKGSvc(dkgNetwork.Round) {
			releaseDKGSvc(dkgNetwork.Round)

			log.Info(ctx, "Recovering stuck DKG session: no active goroutine for intermediate phase",
				"round", dkgNetwork.Round,
				"phase", session.Phase.String(),
				"stage", dkgNetwork.Stage.String(),
			)

			k.stateManager.MarkFailed(ctx, session)
			// Will be picked up as PhaseFailed on next block's ResumeDKGService call.
		}
	}
}

// isSessionStuckForStage returns true if the session phase is behind the expected
// completion state for the current DKG stage. A session that reached the valid
// completion phase for its stage is NOT stuck — it completed successfully and is
// waiting for the next stage transition.
func isSessionStuckForStage(phase types.DKGPhase, stage types.DKGStage) bool {
	switch stage {
	case types.DKGStageRegistration:
		// PhaseInitialized = registration done → not stuck
		return phase != types.PhaseInitialized && phase != types.PhaseCompleted
	case types.DKGStageDealing:
		// PhaseDealing = deals generated → not stuck
		return phase != types.PhaseDealing && phase != types.PhaseCompleted
	case types.DKGStageFinalization:
		// PhaseFinalized = finalization done → not stuck
		return phase != types.PhaseFinalized && phase != types.PhaseCompleted
	case types.DKGStageActive:
		return phase != types.PhaseCompleted
	default:
		return false
	}
}

// resumeFailedSession dispatches a PhaseFailed session to the appropriate handler
// based on the current DKG network stage.
func (k *Keeper) resumeFailedSession(ctx context.Context, session *types.DKGSession, dkgNetwork *types.DKGNetwork) {
	sessionRecoveryTotal.Inc()
	switch dkgNetwork.Stage {
	case types.DKGStageRegistration:
		// Pre-compute registration check while SDK context is available.
		alreadyRegistered := k.isAlreadyRegistered(ctx, dkgNetwork.Round)
		if alreadyRegistered {
			session.UpdatePhase(types.PhaseInitialized)

			if err := k.stateManager.UpdateSession(ctx, session); err != nil {
				log.Error(ctx, "Failed to update session phase after existing registration check", err)
			}

			return
		}

		session.UpdatePhase(types.PhaseInitializing)

		if err := k.stateManager.UpdateSession(ctx, session); err != nil {
			log.Error(ctx, "Failed to update session phase to initializing", err)

			return
		}

		// Pre-compute old code commitment while SDK context is available.
		oldCC, _ := k.getOldCodeCommitment(ctx)

		asyncCtx, cancel := dkgAsyncContext()

		go func() {
			defer cancel()

			k.handleDKGRegistration(asyncCtx, dkgNetwork, oldCC, alreadyRegistered)
		}()
	case types.DKGStageDealing:
		session.UpdatePhase(types.PhaseInitialized)

		if err := k.stateManager.UpdateSession(ctx, session); err != nil {
			log.Error(ctx, "Failed to update session phase to initialized for dealing recovery", err)

			return
		}

		// Set session.Index from on-chain registration for deal/response routing.
		if err := k.ensureSessionIndex(ctx, dkgNetwork.Round); err != nil {
			log.Warn(ctx, "Failed to set session index during resume", err,
				"round", dkgNetwork.Round,
			)
		}

		// Pre-compute shouldDeal while SDK context is available.
		deal, err := k.shouldDeal(ctx, dkgNetwork)
		if err != nil {
			log.Error(ctx, "Failed to check shouldDeal during resume", err)

			return
		}

		asyncCtx, cancel := dkgAsyncContext()

		go func() {
			defer cancel()

			k.handleDKGDealing(asyncCtx, dkgNetwork, deal)
		}()
	case types.DKGStageFinalization:
		session.UpdatePhase(types.PhaseDealing)

		if err := k.stateManager.UpdateSession(ctx, session); err != nil {
			log.Error(ctx, "Failed to update session phase to dealing", err)

			return
		}

		asyncCtx, cancel := dkgAsyncContext()

		go func() {
			defer cancel()

			k.handleDKGFinalization(asyncCtx, dkgNetwork)
		}()
	case types.DKGStageActive:
		// Only complete a session that actually obtained key material. Empty
		// GlobalPubKey/PubKeyShare means this node never produced a share (the round
		// advanced to active via other validators); completing it would contribute
		// zero partial decryptions.
		if len(session.GlobalPubKey) > 0 && len(session.PubKeyShare) > 0 {
			session.UpdatePhase(types.PhaseFinalized)

			if err := k.stateManager.UpdateSession(ctx, session); err != nil {
				log.Error(ctx, "Failed to update session phase to finalized", err)

				return
			}

			asyncCtx, cancel := dkgAsyncContext()

			go func() {
				defer cancel()

				k.handleDKGComplete(asyncCtx, dkgNetwork)
			}()

			return
		}

		// Missing key material: recompute and seal the share locally (node-local only, no
		// on-chain finalize vote for an already-active round).
		log.Warn(ctx, "Active-round session is missing key material; attempting local finalization recovery", nil,
			"round", dkgNetwork.Round,
		)

		asyncCtx, cancel := dkgAsyncContext()

		go func() {
			defer cancel()

			k.recoverActiveSessionKeyMaterial(asyncCtx, dkgNetwork)
		}()
	case types.DKGStageUnspecified:
	}
}

// StartDecryptWorker launches a background loop (non-ABCI) that drains pending decrypt requests
// and performs TDH2 partial decrypts. Only one worker runs.
// The worker uses its own long-lived context (derived from context.Background) because the
// caller's context (dkgAsyncContext) is short-lived and gets cancelled when the parent
// goroutine exits. The decrypt worker must run for the lifetime of the process.
func (k *Keeper) StartDecryptWorker() {
	if !decryptWorkerRunning.CompareAndSwap(false, true) {
		// Debug (not Info): this runs every block; skipping because the worker is
		// already running is the normal steady state and Info would be too noisy.
		log.Debug(context.Background(), "Decrypt worker already running; skipping start")

		return
	}

	// Use a process-lifetime context independent of the caller's short-lived async context.
	workerCtx := context.Background()

	log.Info(workerCtx, "Decrypt worker started")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(workerCtx, "Decrypt worker panicked", errors.New("decrypt worker panic", "value", r))
			}
			decryptWorkerRunning.Store(false)
		}()

		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			k.processDecryptQueue(workerCtx)
		}
	}()
}

// processDecryptQueue scans sessions for queued decrypt requests and starts TDH2 partial decrypt + submission.
func (k *Keeper) processDecryptQueue(ctx context.Context) {
	if k.contractClient == nil {
		log.Error(ctx, "Contract client not configured", nil)
		return
	}

	if k.kernelRouter == nil {
		log.Error(ctx, "Kernel client not configured", nil)
		return
	}

	currentHeight, err := k.contractClient.BlockNumber(ctx)
	if err != nil {
		log.Error(ctx, "Failed to get current block height for decrypt queue processing", err)
		return
	}

	sessions := k.stateManager.ListSessions()

	// Anchor drain on the latest IsFinalized round (issue piplabs/story#826).
	var latestActivated uint32
	for _, s := range sessions {
		if s.IsFinalized && s.Round > latestActivated {
			latestActivated = s.Round
		}
	}

	for _, session := range sessions {
		// Skip sessions with no queued requests so the precondition checks below
		// (index / global public key) don't log misleading deferral warnings for idle sessions.
		if len(session.GetDecryptRequests()) == 0 {
			continue
		}

		// Drop queued requests from sessions that are at least 2 rounds behind the
		// current round — their keys are no longer relevant and the requests would
		// never be processed successfully.
		if latestActivated >= 2 && session.Round <= latestActivated-2 {
			dropped := session.DrainDecryptRequests()
			if len(dropped) > 0 {
				log.Warn(ctx, "Dropping decrypt requests from stale session", nil,
					"session", session.GetSessionKey(),
					"dropped_requests", len(dropped),
					"latest_activated_round", latestActivated,
				)
				incDecryptRequest(labelDecryptStaleDropped, len(dropped))
				if err := k.stateManager.UpdateSession(ctx, session); err != nil {
					log.Error(ctx, "Failed to persist stale session after clearing requests", err,
						"session", session.GetSessionKey(),
					)
				}
			}

			continue
		}

		// Check session-level preconditions before draining so that requests are
		// not removed from the queue only to be re-added immediately.
		if session.Index == 0 {
			log.Warn(ctx, "Session index not set, deferring decrypt requests to next tick", nil,
				"session", session.GetSessionKey(),
			)

			continue
		}

		if len(session.GlobalPubKey) == 0 {
			log.Warn(ctx, "Missing global public key for session, deferring decrypt requests to next tick", nil,
				"session", session.GetSessionKey(),
			)

			continue
		}

		// Atomically drain the queue so that requests added by the ABCI thread
		// during processing are not overwritten when we persist the remaining failures.
		requests := session.DrainDecryptRequests()
		if len(requests) == 0 {
			continue
		}

		// Skip sessions whose kernel binary is no longer connected.
		// This happens when old events are replayed during chain catch-up
		// after a kernel binary change — the sealed keys are unreachable.
		if _, err := k.getClientWithReconnect(session.CodeCommitment); err != nil {
			log.Warn(ctx, "Dropping decrypt requests for session with unavailable kernel", nil,
				"session", session.GetSessionKey(),
				"code_commitment", hex.EncodeToString(session.CodeCommitment),
				"dropped_requests", len(requests),
			)

			if err := k.stateManager.UpdateSession(ctx, session); err != nil {
				log.Error(ctx, "Failed to clear stale decrypt requests", err,
					"session", session.GetSessionKey(),
				)
			}

			continue
		}

		// Filter out stale requests that are past the block timeout window.
		// This is especially important during resync when processing old blocks.
		validRequests := make([]types.PendingDecryptRequest, 0, len(requests))
		staleCount := 0
		for _, req := range requests {
			if currentHeight > types.DefaultDecryptTimeout && req.Height < currentHeight-types.DefaultDecryptTimeout {
				staleCount++
				continue
			}
			validRequests = append(validRequests, req)
		}

		if staleCount > 0 {
			log.Info(ctx, "Filtered out stale decrypt requests past timeout window",
				"session", session.GetSessionKey(),
				"stale_requests", staleCount,
				"current_height", currentHeight,
			)
			incDecryptRequest(labelDecryptStaleDropped, staleCount)
		}

		if len(validRequests) == 0 {
			if err := k.stateManager.UpdateSession(ctx, session); err != nil {
				log.Error(ctx, "Failed to persist session after clearing stale requests", err,
					"session", session.GetSessionKey(),
				)
			}
			continue
		}

		requests = validRequests

		log.Info(ctx, "Processing decrypt queue",
			"session", session.GetSessionKey(),
			"pending_requests", len(requests),
		)

		k.processDecryptRequests(ctx, session, requests)

		if err := k.stateManager.UpdateSession(ctx, session); err != nil {
			log.Error(ctx, "Failed to update session after processing decrypt queue", err,
				"session", session.GetSessionKey(),
			)
		}
	}
}

// processDecryptRequests fans out kernel PartialDecryptTDH2 calls in parallel and
// submits the results to the CDR contract in batches of decryptBatchSize.
//
// Phase 1 and Phase 2 run concurrently:
//
//   - Phase 1 launches kernel goroutines (up to decryptChanBuf in-flight).
//     resultCh is buffered to len(requests) so goroutines never block on send.
//
//   - Phase 2 (batchSubmitConsumer) starts immediately, draining resultCh and
//     submitting batches while Phase 1 is still launching kernel calls. This
//     means the first batch tx can go out before all kernel calls have returned.
//
// The function blocks until Phase 2 completes so the caller can safely persist
// the session (UpdateSession must see all re-queued requests).
func (k *Keeper) processDecryptRequests(ctx context.Context, session *types.DKGSession, requests []types.PendingDecryptRequest) {
	sem := make(chan struct{}, decryptChanBuf)
	// Buffer len(requests) so kernel goroutines never block on send regardless of
	// how fast Phase 2 consumes — this is what enables true Phase 1/2 overlap.
	resultCh := make(chan decryptComputeResult, len(requests))

	// Phase 2: start the batch-submission consumer before launching any kernel
	// goroutines so it is already draining resultCh while Phase 1 is still running.
	// Pass requests so the consumer can look up RetryCount by idx on failure.
	done := make(chan struct{})
	go k.batchSubmitConsumer(ctx, session, requests, resultCh, len(requests), done)

	// Phase 1: launch kernel goroutines. May block on sem when decryptChanBuf slots
	// are exhausted, but Phase 2 runs concurrently and submits batches in parallel.
	for i, preq := range requests {
		sem <- struct{}{}
		go func(req types.DecryptRequest, idx int) {
			defer func() { <-sem }()
			// Pre-initialize so deferred send fires even on runtime.Goexit (t.Fatal).
			result := decryptComputeResult{idx: idx, req: req, err: errors.New("goroutine exited without result")}
			defer func() { resultCh <- result }()
			result = k.computePartialDecrypt(ctx, session, req)
			result.idx = idx // preserve idx since computePartialDecrypt overwrites result
		}(preq.DecryptRequest, i)
	}

	// Wait for Phase 2 to finish before returning so the caller can safely
	// persist the session (UpdateSession must see all re-queued requests).
	<-done
}

// batchSubmitConsumer drains n results from resultCh, accumulates them into
// batches of decryptBatchSize, and submits each full batch (plus any remainder)
// to the CDR contract. Failed kernel results are re-queued; failed batch
// submissions are re-queued so the worker retries on the next tick.
//
// The done channel is closed when all n results have been consumed and submitted,
// signalling processDecryptRequests that it is safe to return.
func (k *Keeper) batchSubmitConsumer(ctx context.Context, session *types.DKGSession, requests []types.PendingDecryptRequest, resultCh <-chan decryptComputeResult, n int, done chan struct{}) {
	defer close(done)
	defer func() {
		if r := recover(); r != nil {
			log.Error(ctx, "Panic in batch submission goroutine", errors.New("panic in batchSubmitConsumer", "value", r),
				"session", session.GetSessionKey(),
			)
			// Re-queuing drained-but-unsent requests is not possible after a
			// panic since we don't know which were already sent. The worker will
			// retry on the next tick via the normal re-queue path.
		}
	}()

	batchSize := k.decryptBatchSize
	if batchSize <= 0 {
		batchSize = defaultDecryptBatchSize
	}
	batch := make([]decryptComputeResult, 0, batchSize)

	flushBatch := func() {
		if len(batch) == 0 {
			return
		}
		if err := k.submitPartialDecryptionBatch(ctx, session, batch); err != nil {
			// WARN: retried below unless a request has exhausted its retries.
			log.Warn(ctx, "Failed to submit partial decryption batch; will retry", err,
				"session", session.GetSessionKey(),
				"batch_size", len(batch),
			)
			incDecryptBatch(labelBatchError, len(batch))
			requeued := 0
			for _, r := range batch {
				preq := requests[r.idx]
				preq.RetryCount++
				if preq.RetryCount > maxReprocessAttempts {
					// Exhausted retries: real failure, alert here.
					log.Error(ctx, "Decrypt request dropped after max retry attempts", err,
						"session", session.GetSessionKey(),
						"round", r.req.Round,
						"retry_count", preq.RetryCount,
						"max_reprocess_attempts", maxReprocessAttempts,
					)
					incDecryptRequest(labelDecryptRetryDropped, 1)
					continue
				}
				session.AddDecryptRequest(preq)
				requeued++
			}
			incDecryptRequest(labelDecryptRequeued, requeued)
		} else {
			log.Info(ctx, "Successfully submitted partial decryption batch",
				"session", session.GetSessionKey(),
				"batch_size", len(batch),
			)
			incDecryptBatch(labelBatchSuccess, len(batch))
			incDecryptRequest(labelDecryptSubmitted, len(batch))
		}
		batch = batch[:0]
	}

	for range n {
		r := <-resultCh

		if r.err != nil {
			incDecryptRequest(labelDecryptKernelFailed, 1)
			preq := requests[r.idx]
			preq.RetryCount++
			if preq.RetryCount > maxReprocessAttempts {
				// Exhausted retries: real failure, alert here.
				log.Error(ctx, "Decrypt request dropped after max retry attempts", r.err,
					"session", session.GetSessionKey(),
					"round", r.req.Round,
					"retry_count", preq.RetryCount,
					"max_reprocess_attempts", maxReprocessAttempts,
				)
				incDecryptRequest(labelDecryptRetryDropped, 1)
				continue
			}
			// Transient (usually light-client lag); self-heals on retry, so WARN only.
			log.Warn(ctx, "Kernel partial decrypt failed; will retry", r.err,
				"session", session.GetSessionKey(),
				"round", r.req.Round,
				"retry_count", preq.RetryCount,
				"max_reprocess_attempts", maxReprocessAttempts,
			)
			session.AddDecryptRequest(preq)
			continue
		}
		batch = append(batch, r)
		if len(batch) >= batchSize {
			flushBatch()
		}
	}
	flushBatch()
}

// computePartialDecrypt calls the kernel for a single TDH2 partial decrypt.
// Any panic is recovered and returned as an error so the parent processDecryptRequests
// goroutine always sends exactly one result to the channel.
func (k *Keeper) computePartialDecrypt(ctx context.Context, session *types.DKGSession, req types.DecryptRequest) (result decryptComputeResult) {
	result.req = req
	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			result.err = errors.New("panic in computePartialDecrypt", "value", r)
		}
	}()

	client, err := k.getClientWithReconnect(session.CodeCommitment)
	if err != nil {
		result.err = errors.Wrap(err, "no kernel client for session")
		return
	}

	// Bound the kernel RPC so a wedged kernel cannot block the worker loop
	// indefinitely. The worker's context has no deadline of its own.
	callCtx, cancel := context.WithTimeout(ctx, decryptKernelCallTimeout)
	defer cancel()

	kernelStart := time.Now()
	resp, err := client.PartialDecryptTDH2(callCtx, &types.PartialDecryptTDH2Request{
		CodeCommitment:  session.CodeCommitment,
		Round:           session.Round,
		Ciphertext:      req.Ciphertext,
		Label:           req.Label,
		GlobalPubKey:    session.GlobalPubKey,
		RequesterPubKey: req.RequesterPubKey,
	})
	kernelDuration := time.Since(kernelStart)
	observeKernelCall(labelOpPartialDecryptTDH2, kernelStart, err)

	if err != nil {
		result.err = errors.Wrap(err, "generating partial decrypt failed")
		// DEBUG only; batchSubmitConsumer owns the retry/drop decision and WARN/ERROR.
		log.Debug(ctx, "Kernel PartialDecryptTDH2 call failed; retry decision deferred to batch consumer",
			"session", session.GetSessionKey(),
			"kernel_duration_ms", kernelDuration.Milliseconds(),
			"err", err,
		)
		return
	}

	log.Info(ctx, "Kernel PartialDecryptTDH2 completed",
		"session", session.GetSessionKey(),
		"kernel_duration_ms", kernelDuration.Milliseconds(),
		"total_duration_ms", time.Since(start).Milliseconds(),
	)

	result.resp = resp

	return
}

// submitPartialDecryption submits a single partial-decrypt response to the CDR contract.
// Used by tests and as a single-item fallback.
func (k *Keeper) submitPartialDecryption(ctx context.Context, session *types.DKGSession, req types.DecryptRequest, resp *types.PartialDecryptTDH2Response) error {
	uuid, err := labelToUUID(req.Label)
	if err != nil {
		return errors.Wrap(err, "invalid decrypt request label")
	}

	if _, err := k.contractClient.SubmitEncryptedPartialDecryption(
		ctx,
		session.Round,
		session.Index,
		resp.EncryptedPartialDecryption,
		resp.EphemeralPubKey,
		resp.PubShare,
		req.RequesterPubKey,
		req.Ciphertext,
		uuid,
		resp.Signature,
	); err != nil {
		return errors.Wrap(err, "failed to submit partial decryption")
	}

	return nil
}

// submitPartialDecryptionBatch submits a batch of kernel partial-decrypt results in a
// single contract call, reducing gas overhead and tx count.
func (k *Keeper) submitPartialDecryptionBatch(ctx context.Context, session *types.DKGSession, results []decryptComputeResult) error {
	batchReqs := make([]bindings.ICDRPartialDecryptionRequest, 0, len(results))

	for _, r := range results {
		uuid, err := labelToUUID(r.req.Label)
		if err != nil {
			return errors.Wrap(err, "invalid decrypt request label in batch")
		}

		batchReqs = append(batchReqs, bindings.ICDRPartialDecryptionRequest{
			Round:            session.Round,
			Pid:              session.Index,
			EncryptedPartial: r.resp.EncryptedPartialDecryption,
			EphemeralPubKey:  r.resp.EphemeralPubKey,
			PubShare:         r.resp.PubShare,
			RequesterPubKey:  r.req.RequesterPubKey,
			Ciphertext:       r.req.Ciphertext,
			Uuid:             uuid,
			Signature:        r.resp.Signature,
		})
	}

	if _, err := k.contractClient.SubmitEncryptedPartialDecryptionBatch(ctx, batchReqs); err != nil {
		return errors.Wrap(err, "failed to submit partial decryption batch")
	}

	return nil
}

func labelToUUID(label []byte) (uint32, error) {
	if len(label) < 32 {
		return 0, errors.New("label must be 32 bytes")
	}
	return binary.BigEndian.Uint32(label[28:]), nil
}

// getClientWithReconnect returns the kernel client for the given code commitment.
// If the initial lookup fails, it attempts to reconnect any disconnected endpoints
// and retries the lookup once. This handles the case where story started before kernel.
func (k *Keeper) getClientWithReconnect(codeCommitment []byte) (types.KernelServiceClient, error) {
	start := time.Now()

	client, err := k.kernelRouter.GetClient(codeCommitment)
	if err == nil {
		kernelClientLookupDuration.WithLabelValues(labelLookupHit).Observe(time.Since(start).Seconds())
		return client, nil
	}

	k.kernelRouter.TryReconnect()
	log.Info(context.Background(), "Retrying kernel client lookup after reconnect", "code_commitment", hex.EncodeToString(codeCommitment))

	client, err = k.kernelRouter.GetClient(codeCommitment)
	if err != nil {
		kernelClientLookupDuration.WithLabelValues(labelLookupError).Observe(time.Since(start).Seconds())
		return nil, err
	}

	kernelClientLookupDuration.WithLabelValues(labelLookupReconnect).Observe(time.Since(start).Seconds())

	return client, nil
}
