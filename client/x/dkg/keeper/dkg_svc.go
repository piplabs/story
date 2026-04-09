package keeper

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
	"time"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// decryptComputeResult holds the outcome of a single parallel kernel partial-decrypt call.
type decryptComputeResult struct {
	req  types.DecryptRequest
	resp *types.PartialDecryptTDH2Response
	err  error
}

// decryptChanBuf is the fixed buffer size for the parallel-decrypt result channel.
// This caps in-flight results and slows producers when the serial submission
// loop is slower than the parallel kernel calls.
const decryptChanBuf = 10

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
const dkgAsyncTimeout = 1 * time.Minute

// dkgAsyncContext creates a new context for async DKG service goroutines with a timeout.
// This replaces the consensus context that would otherwise be canceled after block processing.
func dkgAsyncContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), dkgAsyncTimeout)
}

var decryptWorkerRunning atomic.Bool

// ResumeDKGService reloads unfinished DKG sessions and resumes their execution safely without spawning duplicate goroutines.
func (k *Keeper) ResumeDKGService(ctx context.Context, dkgNetwork *types.DKGNetwork) {
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
		// already running
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
	for _, session := range sessions {
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
		validRequests := make([]types.DecryptRequest, 0, len(requests))
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
// submits the results to the contract serially (nonce-safe).
//
// Phase 1 (parallel): each request is dispatched to the kernel in its own goroutine.
// Results are sent to a fixed-size buffered channel (decryptChanBuf). Goroutines are
// protected by a deferred recover so a panic in one never silently leaks the others —
// the panic is converted to an error result and the goroutine count sent to the channel
// is always exactly len(requests).
//
// Phase 2 (serial): the main goroutine reads results one-by-one and submits each to
// the CDR contract. Because PendingNonceAt is called once per submission without a
// mutex, concurrent submissions would produce nonce collisions, so this phase must
// remain sequential.
func (k *Keeper) processDecryptRequests(ctx context.Context, session *types.DKGSession, requests []types.DecryptRequest) {

	// Phase 1: parallel kernel calls.
	// Each goroutine sends exactly one result, so we read exactly len(requests) times
	// below — no close or WaitGroup needed.
	// Buffer size is capped at decryptChanBuf to limit concurrent in-flight RPCs.
	bufSize := min(len(requests), decryptChanBuf)
	resultCh := make(chan decryptComputeResult, bufSize)

	for _, req := range requests {
		go func(req types.DecryptRequest) {
			// Pre-initialize with an error so that if runtime.Goexit() is called
			// (e.g. t.Fatal from a test mock), the deferred send still fires and
			// the consumer loop is not left blocked on an unreceived value.
			// recover() inside computePartialDecrypt handles actual panics;
			// this defer handles Goexit, which recover() cannot intercept.
			result := decryptComputeResult{req: req, err: errors.New("goroutine exited without result")}
			defer func() { resultCh <- result }()
			result = k.computePartialDecrypt(ctx, session, req)
		}(req)
	}

	// Phase 2: serial contract submission.
	// Read exactly len(requests) results — guaranteed because every goroutine above
	// sends exactly one value via the deferred send (covers both normal return,
	// panic via recover, and runtime.Goexit from test helpers).
	for range requests {
		r := <-resultCh
		if r.err != nil {
			log.Error(ctx, "Kernel partial decrypt failed", r.err,
				"session", session.GetSessionKey(),
				"round", r.req.Round,
			)
			session.AddDecryptRequest(r.req)

			continue
		}

		if err := k.submitPartialDecryption(ctx, session, r.req, r.resp); err != nil {
			log.Error(ctx, "Failed to submit partial decryption", err,
				"session", session.GetSessionKey(),
				"round", r.req.Round,
				"ciphertext_len", len(r.req.Ciphertext),
				"label_len", len(r.req.Label),
				"requester_pub_key_len", len(r.req.RequesterPubKey),
			)
			session.AddDecryptRequest(r.req)

			continue
		}

		log.Info(ctx, "Successfully processed decrypt request",
			"session", session.GetSessionKey(),
			"round", r.req.Round,
		)
	}
}

// computePartialDecrypt calls the kernel for a single TDH2 partial decrypt.
// Any panic is recovered and returned as an error so the parent processDecryptRequests
// goroutine always sends exactly one result to the channel.
func (k *Keeper) computePartialDecrypt(ctx context.Context, session *types.DKGSession, req types.DecryptRequest) (result decryptComputeResult) {
	result.req = req

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

	resp, err := client.PartialDecryptTDH2(ctx, &types.PartialDecryptTDH2Request{
		CodeCommitment:  session.CodeCommitment,
		Round:           session.Round,
		Ciphertext:      req.Ciphertext,
		Label:           req.Label,
		GlobalPubKey:    session.GlobalPubKey,
		RequesterPubKey: req.RequesterPubKey,
	})
	if err != nil {
		result.err = errors.Wrap(err, "generating partial decrypt failed")
		return
	}

	result.resp = resp

	return
}

// submitPartialDecryption submits the kernel's partial-decrypt response to the CDR contract.
// Must be called serially — ContractClient.createTransactOpts fetches PendingNonceAt without
// a mutex, so concurrent calls produce nonce collisions.
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
	client, err := k.kernelRouter.GetClient(codeCommitment)
	if err == nil {
		return client, nil
	}

	k.kernelRouter.TryReconnect()

	return k.kernelRouter.GetClient(codeCommitment)
}
