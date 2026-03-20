package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"slices"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGDealing handles the dealing phase event.
// shouldDeal must be pre-computed by the caller while the SDK context is available,
// because async goroutines cannot access the KV store.
func (k *Keeper) handleDKGDealing(ctx context.Context, dkgNetwork *types.DKGNetwork, shouldDeal bool) {
	log.Info(ctx, "Handling DKG dealing",
		"round", dkgNetwork.Round,
		"should_deal", shouldDeal,
	)

	if !tryAcquireDKGSvc(dkgNetwork.Round) {
		log.Info(ctx, "DKG service already running for this round; skipping dealing",
			"round", dkgNetwork.Round,
		)

		return
	}
	defer releaseDKGSvc(dkgNetwork.Round)

	if dkgNetwork.Stage != types.DKGStageDealing {
		log.Info(ctx, "DKG Dealing is skipped because the current network stage is not dealing stage")

		return
	}

	if !shouldDeal {
		log.Debug(ctx, "Skip dealing")

		return
	}

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session", err)

		return
	}

	if session.Phase != types.PhaseInitialized {
		log.Warn(ctx, "Session not in initialized phase, skipping generate deals", nil,
			"current_phase", session.Phase.String())
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	// For upgrade resharing, the dealer uses the old binary's kernel client
	// because the old key shares are sealed by the old binary.
	dealerCC := session.CodeCommitment
	if dkgNetwork.IsUpgrade && len(session.OldCodeCommitment) > 0 {
		dealerCC = session.OldCodeCommitment
	}

	var resp *types.GenerateDealsResponse

	if err := retry(ctx, func(ctx context.Context) error {
		log.Info(ctx, "GenerateDeals call to kernel client",
			"round", session.Round,
			"is_upgrade", dkgNetwork.IsUpgrade,
		)

		req := &types.GenerateDealsRequest{
			CodeCommitment: dealerCC,
			Round:          session.Round,
			IsResharing:    session.IsResharing,
		}
		client, cErr := k.kernelRouter.GetClient(dealerCC)
		if cErr != nil {
			return errors.Wrap(cErr, "no kernel client for session")
		}

		resp, err = client.GenerateDeals(ctx, req)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Error(ctx, "Failed to generate deals", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	session.Phase = types.PhaseDealing

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update session after generating deals", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	k.EnqueueDeals(resp.GetDeals())

	log.Info(ctx, "DKG deals are generated successfully",
		"round", session.Round,
	)
}

// handleDKGProcessDeals handles the deals from other committee members.
func (k *Keeper) handleDKGProcessDeals(ctx context.Context, dkgNetwork *types.DKGNetwork, deals []types.Deal) {
	// Serialize kernel DKG operations to prevent concurrent DistKeyGenerator mutation.
	dkgKernelMu.Lock()
	defer dkgKernelMu.Unlock()

	log.Info(ctx, "Handling DKG process deals",
		"round", dkgNetwork.Round,
		"num_deals", len(deals),
	)

	if !slices.Contains(dkgNetwork.ActiveValSet, k.validatorEVMAddr) {
		log.Info(ctx, "Skip processing deals as the validator is not in current round set")

		return
	}

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session", err)

		return
	}

	// Accept deals in both Initialized and Dealing phases.
	// Deals from other validators arrive via vote extensions as soon as the DKG stage
	// transitions to Dealing. However, the local session only transitions from
	// Initialized to Dealing after GenerateDeals completes (which can take ~30s).
	// If we reject deals during Initialized phase, they are permanently lost because
	// vote extensions deliver each deal exactly once.
	if session.Phase != types.PhaseDealing && session.Phase != types.PhaseInitialized {
		log.Warn(ctx, "Session not in dealing or initialized phase, skipping process deals", nil,
			"current_phase", session.Phase.String(),
		)

		return
	}

	var resp *types.ProcessDealsResponse

	if err := retry(ctx, func(ctx context.Context) error {
		log.Info(ctx, "ProcessDeals call to kernel client",
			"round", session.Round,
			"num_deals", len(deals),
		)

		req := &types.ProcessDealsRequest{
			CodeCommitment: session.CodeCommitment,
			Round:          session.Round,
			Deals:          []types.Deal{},
			IsResharing:    session.IsResharing,
		}

		for _, deal := range deals {
			// RecipientIndex is 0-based (Kyber), session.Index is 1-based (on-chain).
			// Guard against unset index (0) to avoid uint32 underflow.
			if session.Index > 0 && deal.RecipientIndex == session.Index-1 {
				req.Deals = append(req.Deals, deal)
			}
		}

		if len(req.Deals) == 0 {
			log.Info(ctx, "No deals to process. Skip to request")

			return nil
		}

		client, cErr := k.kernelRouter.GetClient(session.CodeCommitment)
		if cErr != nil {
			return errors.Wrap(cErr, "no kernel client for session")
		}

		resp, err = client.ProcessDeals(ctx, req)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		// Cache unprocessed deals in memory for retry when kernel recovers.
		// Not persisted to disk — process restart loses them (round will fail and retry).
		cached := cachePendingIncomingDeals(deals, session.Index)

		log.Error(ctx, "Failed to process deals; cached for retry", err,
			"round", dkgNetwork.Round,
			"cached_deals", cached,
		)

		return
	}

	k.EnqueueResponses(resp.GetResponses())

	log.Info(ctx, "Process deals complete",
		"round", session.Round,
	)
}

// cachePendingIncomingDeals saves deals that failed kernel processing for later retry.
// Only deals addressed to this validator (matching recipientIndex) are cached.
// Returns the number of deals cached.
func cachePendingIncomingDeals(allDeals []types.Deal, sessionIndex uint32) int {
	pendingIncomingDealsMu.Lock()
	defer pendingIncomingDealsMu.Unlock()

	cached := 0
	for _, deal := range allDeals {
		if sessionIndex > 0 && deal.RecipientIndex == sessionIndex-1 {
			if len(pendingIncomingDeals) >= maxPendingIncoming {
				return cached
			}

			pendingIncomingDeals = append(pendingIncomingDeals, deal)
			cached++
		}
	}

	return cached
}

// handleDKGProcessResponses handles the responses of processDeals from other committee members.
// shouldProcess must be pre-computed by the caller while the SDK context is available.
func (k *Keeper) handleDKGProcessResponses(ctx context.Context, dkgNetwork *types.DKGNetwork, responses []types.Response, shouldProcess bool) {
	// Serialize kernel DKG operations to prevent concurrent DistKeyGenerator mutation.
	dkgKernelMu.Lock()
	defer dkgKernelMu.Unlock()

	log.Info(ctx, "Handling DKG process responses",
		"round", dkgNetwork.Round,
		"num_responses", len(responses),
		"should_process", shouldProcess,
	)

	if !shouldProcess {
		log.Info(ctx, "Skip processing of responses")

		return
	}

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session", err)

		return
	}

	// Accept responses in both Initialized and Dealing phases (same reasoning as deals).
	if session.Phase != types.PhaseDealing && session.Phase != types.PhaseInitialized {
		log.Warn(ctx, "Session not in dealing or initialized phase, skipping process responses", nil,
			"current_phase", session.Phase.String(),
		)

		return
	}

	// During upgrade resharing, validators in both old and new sets must send
	// ProcessResponses to BOTH binaries: the old binary (dealer role) and the new binary
	// (recipient role). Each binary maintains its own DKG state that needs updating.
	ccsToProcess := [][]byte{session.CodeCommitment}
	if dkgNetwork.IsUpgrade && len(session.OldCodeCommitment) > 0 && !bytes.Equal(session.CodeCommitment, session.OldCodeCommitment) {
		ccsToProcess = append(ccsToProcess, session.OldCodeCommitment)
	}

	filteredResponses := make([]types.Response, 0, len(responses))
	for _, resp := range responses {
		// VssResponse.Index is 0-based (Kyber), session.Index is 1-based (on-chain).
		// If session.Index is 0 (unset), include all responses (no self-filtering).
		if session.Index == 0 || resp.VssResponse.Index != session.Index-1 {
			filteredResponses = append(filteredResponses, resp)
		}
	}

	if len(filteredResponses) == 0 {
		log.Info(ctx, "No responses to process. Skip to request")

		return
	}

	for _, cc := range ccsToProcess {
		var processResp *types.ProcessResponsesResponse

		if err := retry(ctx, func(ctx context.Context) error {
			log.Info(ctx, "ProcessResponses call to kernel client",
				"round", session.Round,
				"num_responses", len(filteredResponses),
				"code_commitment", hex.EncodeToString(cc),
			)

			req := &types.ProcessResponsesRequest{
				CodeCommitment: cc,
				Round:          session.Round,
				Responses:      filteredResponses,
				IsResharing:    session.IsResharing,
			}

			client, cErr := k.kernelRouter.GetClient(cc)
			if cErr != nil {
				return errors.Wrap(cErr, "no kernel client for session")
			}

			if processResp, err = client.ProcessResponses(ctx, req); err != nil {
				return err
			}

			return nil
		}); err != nil {
			log.Error(ctx, "Failed to process responses", err,
				"code_commitment", hex.EncodeToString(cc),
			)

			// Cache unprocessed responses for retry when kernel recovers.
			cachePendingIncomingResponses(filteredResponses)

			log.Info(ctx, "Cached responses for retry",
				"round", session.Round,
				"cached_responses", len(filteredResponses),
			)

			continue
		}

		// Enqueue any justifications returned by story-kernel for broadcast via Vote Extension.
		// Justifications are produced when a complaint response (status=false) is processed
		// and the dealer needs to reveal the plaintext deal to prove its validity.
		if processResp != nil && len(processResp.GetJustifications()) > 0 {
			k.EnqueueJustifications(processResp.GetJustifications())

			log.Info(ctx, "Enqueued justifications for broadcast",
				"code_commitment", hex.EncodeToString(cc),
				"round", session.Round,
				"num_justifications", len(processResp.GetJustifications()),
			)
		}
	}

	log.Info(ctx, "Process responses complete",
		"round", session.Round,
	)
}

// handleDKGProcessJustifications sends already-verified justifications to the kernel.
// Justifications have passed Schnorr signature and VSS verification in ProcessJustifications
// (FinalizeBlock context). This function only forwards them to the kernel via gRPC.
// Also used for replaying cached justifications that failed the kernel gRPC call.
func (k *Keeper) handleDKGProcessJustifications(ctx context.Context, dkgNetwork *types.DKGNetwork, justifications []types.Justification) {
	dkgKernelMu.Lock()
	defer dkgKernelMu.Unlock()

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session for justification processing", err)

		return
	}

	ccsToProcess := [][]byte{session.CodeCommitment}
	if dkgNetwork.IsUpgrade && len(session.OldCodeCommitment) > 0 && !bytes.Equal(session.CodeCommitment, session.OldCodeCommitment) {
		ccsToProcess = append(ccsToProcess, session.OldCodeCommitment)
	}

	for _, cc := range ccsToProcess {
		if err := retry(ctx, func(ctx context.Context) error {
			req := &types.ProcessJustificationRequest{
				CodeCommitment: cc,
				Round:          session.Round,
				Justifications: justifications,
				IsResharing:    session.IsResharing,
			}

			client, cErr := k.kernelRouter.GetClient(cc)
			if cErr != nil {
				return errors.Wrap(cErr, "no kernel client for session")
			}

			if _, err := client.ProcessJustification(ctx, req); err != nil {
				return err
			}

			return nil
		}); err != nil {
			log.Error(ctx, "Failed to process justifications via kernel", err,
				"code_commitment", hex.EncodeToString(cc),
			)

			cachePendingIncomingJustifications(justifications)

			continue
		}
	}
}

func (k *Keeper) shouldDeal(ctx context.Context, dkgNetwork *types.DKGNetwork) (bool, error) {
	inCurSet := slices.Contains(dkgNetwork.ActiveValSet, k.validatorEVMAddr)

	prevActive, err := k.getLatestActiveDKGNetwork(ctx)
	if err != nil {
		return false, err
	}

	if prevActive == nil {
		// First DKG round (no previous active round): current set of validators deal
		return inCurSet, nil
	}

	// Resharing round: only previous set of validators deal (they hold the existing key shares)
	inPrevSet := slices.Contains(prevActive.ActiveValSet, k.validatorEVMAddr)

	return inPrevSet, nil
}

// cachePendingIncomingResponses saves responses that failed kernel processing for later retry.
func cachePendingIncomingResponses(filteredResponses []types.Response) {
	pendingIncomingResponsesMu.Lock()
	defer pendingIncomingResponsesMu.Unlock()

	remaining := maxPendingIncoming - len(pendingIncomingResponses)
	if remaining <= 0 {
		return
	}

	if len(filteredResponses) > remaining {
		filteredResponses = filteredResponses[:remaining]
	}

	pendingIncomingResponses = append(pendingIncomingResponses, filteredResponses...)
}

func cachePendingIncomingJustifications(justifications []types.Justification) {
	pendingIncomingJustificationsMu.Lock()
	defer pendingIncomingJustificationsMu.Unlock()

	remaining := maxPendingIncoming - len(pendingIncomingJustifications)
	if remaining <= 0 {
		return
	}

	if len(justifications) > remaining {
		justifications = justifications[:remaining]
	}

	pendingIncomingJustifications = append(pendingIncomingJustifications, justifications...)
}

// reprocessPendingIncomingData retries cached deals, responses, and justifications when the kernel recovers.
// Deals are processed FIRST (kyber requires deals before responses can be accepted).
// Called from BeginBlocker via the DKG service loop.
func (k *Keeper) reprocessPendingIncomingData(dkgNetwork *types.DKGNetwork) {
	hasPendingDeals := func() bool {
		pendingIncomingDealsMu.Lock()
		defer pendingIncomingDealsMu.Unlock()

		return len(pendingIncomingDeals) > 0
	}()

	hasPendingResponses := func() bool {
		pendingIncomingResponsesMu.Lock()
		defer pendingIncomingResponsesMu.Unlock()

		return len(pendingIncomingResponses) > 0
	}()

	hasPendingJustifications := func() bool {
		pendingIncomingJustificationsMu.Lock()
		defer pendingIncomingJustificationsMu.Unlock()

		return len(pendingIncomingJustifications) > 0
	}()

	if !hasPendingDeals && !hasPendingResponses && !hasPendingJustifications {
		return
	}

	// Only retry during Dealing stage (pending data is stale after stage transitions).
	if dkgNetwork.Stage != types.DKGStageDealing {
		flushPendingIncoming()

		return
	}

	if k.kernelRouter == nil || !k.kernelRouter.HasClients() {
		return // kernel still unavailable, wait
	}

	asyncCtx, cancel := dkgAsyncContext()

	go func() {
		defer cancel()

		// Process deals FIRST — kyber returns ErrNoDealBeforeResponse if
		// a response arrives for a dealer whose deal hasn't been processed.
		if hasPendingDeals {
			drainedDeals := drainPendingIncomingDeals()

			log.Info(asyncCtx, "Replaying cached deals after kernel recovery",
				"round", dkgNetwork.Round,
				"num_deals", len(drainedDeals),
			)

			k.handleDKGProcessDeals(asyncCtx, dkgNetwork, drainedDeals)
		}

		// Then process responses.
		if hasPendingResponses {
			drainedResponses := drainPendingIncomingResponses()

			log.Info(asyncCtx, "Replaying cached responses after kernel recovery",
				"round", dkgNetwork.Round,
				"num_responses", len(drainedResponses),
			)

			k.handleDKGProcessResponses(asyncCtx, dkgNetwork, drainedResponses, true)
		}

		// Finally process justifications (already verified before caching,
		// so forward directly to kernel without re-verification).
		if hasPendingJustifications {
			drainedJustifications := drainPendingIncomingJustifications()

			log.Info(asyncCtx, "Replaying cached justifications after kernel recovery",
				"round", dkgNetwork.Round,
				"num_justifications", len(drainedJustifications),
			)

			k.handleDKGProcessJustifications(asyncCtx, dkgNetwork, drainedJustifications)
		}
	}()
}

func drainPendingIncomingDeals() []types.Deal {
	pendingIncomingDealsMu.Lock()
	defer pendingIncomingDealsMu.Unlock()

	out := pendingIncomingDeals
	pendingIncomingDeals = nil

	return out
}

func drainPendingIncomingResponses() []types.Response {
	pendingIncomingResponsesMu.Lock()
	defer pendingIncomingResponsesMu.Unlock()

	out := pendingIncomingResponses
	pendingIncomingResponses = nil

	return out
}

func drainPendingIncomingJustifications() []types.Justification {
	pendingIncomingJustificationsMu.Lock()
	defer pendingIncomingJustificationsMu.Unlock()

	out := pendingIncomingJustifications
	pendingIncomingJustifications = nil

	return out
}

func flushPendingIncoming() {
	pendingIncomingDealsMu.Lock()
	pendingIncomingDeals = nil
	pendingIncomingDealsMu.Unlock()

	pendingIncomingResponsesMu.Lock()
	pendingIncomingResponses = nil
	pendingIncomingResponsesMu.Unlock()

	pendingIncomingJustificationsMu.Lock()
	pendingIncomingJustifications = nil
	pendingIncomingJustificationsMu.Unlock()
}

func (k *Keeper) shouldProcessResponses(ctx context.Context, dkgNetwork *types.DKGNetwork) (bool, error) {
	inCurSet := slices.Contains(dkgNetwork.ActiveValSet, k.validatorEVMAddr)

	prevActive, err := k.getLatestActiveDKGNetwork(ctx)
	if err != nil {
		return false, err
	}

	if prevActive == nil {
		// First DKG round (no previous active round): current set of validators process
		return inCurSet, nil
	}

	// Resharing round: both previous and current set of validators process deals or responses
	inPrevSet := slices.Contains(prevActive.ActiveValSet, k.validatorEVMAddr)

	return inCurSet || inPrevSet, nil
}
