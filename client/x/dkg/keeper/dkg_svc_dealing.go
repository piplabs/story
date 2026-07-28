package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"slices"
	"time"

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

	if session.GetPhase() != types.PhaseInitialized {
		log.Warn(ctx, "Session not in initialized phase, skipping generate deals", nil,
			"current_phase", session.GetPhase().String())
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	round := session.GetRound()
	isResharing := session.GetIsResharing()

	// For upgrade resharing, the dealer uses the old binary's kernel client
	// because the old key shares are sealed by the old binary.
	dealerCC := session.GetCodeCommitment()
	if oldCC := session.GetOldCodeCommitment(); dkgNetwork.IsUpgrade && len(oldCC) > 0 {
		dealerCC = oldCC
	}

	var resp *types.GenerateDealsResponse

	start := time.Now()
	retryErr := retry(ctx, func(ctx context.Context) error {
		log.Info(ctx, "GenerateDeals call to kernel client",
			"round", round,
			"is_upgrade", dkgNetwork.IsUpgrade,
		)

		req := &types.GenerateDealsRequest{
			CodeCommitment: dealerCC,
			Round:          round,
			IsResharing:    isResharing,
		}
		client, cErr := k.getClientWithReconnect(dealerCC)
		if cErr != nil {
			return errors.Wrap(cErr, "no kernel client for session")
		}

		resp, err = client.GenerateDeals(ctx, req)
		if err != nil {
			return err
		}

		return nil
	})
	observeKernelCall(labelOpGenerateDeals, start, retryErr)

	if retryErr != nil {
		log.Error(ctx, "Failed to generate deals", retryErr)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	session.UpdatePhase(types.PhaseDealing)

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update session after generating deals", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	k.EnqueueDeals(resp.GetDeals())

	// DEBUG: self_index is this validator's 1-based on-chain registration index
	// (its dealer identity). Logging it next to the enqueued deal count makes clear
	// which validator produced this batch of deals.
	log.Info(ctx, "DKG deals are generated successfully",
		"round", round,
		"self_index", session.GetIndex(),
		"is_resharing", isResharing,
		"enqueued_deals", len(resp.GetDeals()),
	)
}

// handleDKGProcessDeals handles the deals from other committee members.
// Accepts []pendingDeal so that both first-time callers (via wrapDeals) and
// reprocess callers (with preserved retryCount) use the same code path.
func (k *Keeper) handleDKGProcessDeals(ctx context.Context, dkgNetwork *types.DKGNetwork, pending []pendingDeal) {
	// Serialize kernel DKG operations to prevent concurrent DistKeyGenerator mutation.
	dkgKernelMu.Lock()
	defer dkgKernelMu.Unlock()

	log.Info(ctx, "Handling DKG process deals",
		"round", dkgNetwork.Round,
		"num_deals", len(pending),
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
	if phase := session.GetPhase(); phase != types.PhaseDealing && phase != types.PhaseInitialized {
		log.Warn(ctx, "Session not in dealing or initialized phase, skipping process deals", nil,
			"current_phase", phase.String(),
		)

		return
	}

	round := session.GetRound()
	isResharing := session.GetIsResharing()
	codeCommitment := session.GetCodeCommitment()
	selfIndex := session.GetIndex()

	// Filter deals addressed to this validator before entering retry loop.
	// RecipientIndex is 0-based (Kyber), selfIndex is 1-based (on-chain).
	// Guard against unset index (0) to avoid uint32 underflow.
	filtered := make([]pendingDeal, 0, len(pending))
	for _, pd := range pending {
		if selfIndex > 0 && pd.deal.RecipientIndex == selfIndex-1 {
			filtered = append(filtered, pd)
		}
	}

	if len(filtered) == 0 {
		log.Info(ctx, "No deals addressed to this validator; skipping",
			"round", dkgNetwork.Round,
			"total_deals", len(pending),
		)

		return
	}

	// Extract raw deals from filtered pending items for the kernel call.
	rawDeals := make([]types.Deal, len(filtered))
	for i, pd := range filtered {
		rawDeals[i] = pd.deal
	}

	var resp *types.ProcessDealsResponse

	start := time.Now()
	retryErr := retry(ctx, func(ctx context.Context) error {
		log.Info(ctx, "ProcessDeals call to kernel client",
			"round", round,
			"total_deals", len(pending),
			"filtered_deals", len(filtered),
		)

		req := &types.ProcessDealsRequest{
			CodeCommitment: codeCommitment,
			Round:          round,
			Deals:          rawDeals,
			IsResharing:    isResharing,
		}

		client, cErr := k.getClientWithReconnect(codeCommitment)
		if cErr != nil {
			return errors.Wrap(cErr, "no kernel client for session")
		}

		resp, err = client.ProcessDeals(ctx, req)
		if err != nil {
			return err
		}

		return nil
	})
	observeKernelCall(labelOpProcessDeals, start, retryErr)

	if retryErr != nil {
		// Increment retryCount and drop items exceeding maxReprocessAttempts.
		var remaining []pendingDeal
		for i := range filtered {
			filtered[i].retryCount++
			if filtered[i].retryCount > maxReprocessAttempts {
				log.Warn(ctx, "Dropping pending deal after max retry attempts", nil,
					"round", dkgNetwork.Round,
					"deal_index", filtered[i].deal.Index,
					"retry_count", filtered[i].retryCount,
				)

				continue
			}

			remaining = append(remaining, filtered[i])
		}

		if len(remaining) > 0 {
			cachePendingDeals(ctx, remaining)
		}

		log.Error(ctx, "Failed to process deals; cached for retry", retryErr,
			"round", dkgNetwork.Round,
			"cached_deals", len(remaining),
		)

		return
	}

	// Per-item retry: requeue any deals the kernel rejected. Idempotent
	// re-submissions are silently skipped on the kernel side and never appear
	// in rejected_deals; old kernels leave the field absent. Both cases fall
	// through with an empty `rejected` slice.
	rejected := resp.GetRejectedDeals()
	if len(rejected) > 0 {
		// DEBUG: list the rejected deals' sender (dealer) indices. This handler runs
		// in an async context without KV access, so the dealer index -> validator
		// address mapping is NOT resolved here; cross-reference the consensus-path
		// "markDealersDealt: dealer submitted deal" logs for the same round to turn
		// these indices into addresses.
		rejectedIdx := make([]uint32, 0, len(rejected))
		for _, r := range rejected {
			rejectedIdx = append(rejectedIdx, r.Index)
		}
		log.Warn(ctx, "Kernel rejected deals", nil,
			"round", dkgNetwork.Round,
			"self_index", selfIndex,
			"rejected_sender_index", rejectedIdx,
		)

		rejectedSet := make(map[rejectedKey]struct{}, len(rejected))
		for _, r := range rejected {
			rejectedSet[rejectedDealKey(r)] = struct{}{}
		}
		var remaining []pendingDeal
		for i := range filtered {
			if _, isRejected := rejectedSet[rejectedDealKey(filtered[i].deal)]; !isRejected {
				continue
			}
			filtered[i].retryCount++
			if filtered[i].retryCount > maxReprocessAttempts {
				log.Warn(ctx, "Dropping pending deal after max retry attempts", nil,
					"round", dkgNetwork.Round,
					"deal_index", filtered[i].deal.Index,
					"retry_count", filtered[i].retryCount,
				)

				continue
			}
			remaining = append(remaining, filtered[i])
		}
		if len(remaining) > 0 {
			cachePendingDeals(ctx, remaining)
		}
	}

	k.EnqueueResponses(resp.GetResponses())

	log.Info(ctx, "Process deals complete",
		"round", round,
		"submitted_deals", len(rawDeals),
		"responses_generated", len(resp.GetResponses()),
		"rejected_deals", len(resp.GetRejectedDeals()),
	)
}

// rejectedKey identifies a single rejected item within a batch. For deals
// the outer Index (sender) is unique within the filtered batch (the CL
// pre-filters by RecipientIndex). For responses and justifications the
// outer Index is the dealer and can legitimately repeat across complainers,
// so the inner index (VssResponse.Index / VssJustification.Index) is needed
// for disambiguation.
type rejectedKey struct {
	outer uint32
	inner uint32
}

func rejectedDealKey(d types.Deal) rejectedKey {
	return rejectedKey{outer: d.Index}
}

func rejectedResponseKey(r types.Response) rejectedKey {
	if r.VssResponse == nil {
		return rejectedKey{outer: r.Index}
	}

	return rejectedKey{outer: r.Index, inner: r.VssResponse.Index}
}

func rejectedJustificationKey(j types.Justification) rejectedKey {
	if j.VssJustification == nil {
		return rejectedKey{outer: j.Index}
	}

	return rejectedKey{outer: j.Index, inner: j.VssJustification.Index}
}

// cachePendingDeals saves pending deals for later retry, respecting the
// capacity limit. The pending-deals metric is incremented by the number of
// items actually cached (post-capacity-trim). Items dropped because the
// queue is at capacity are logged and counted under op="dropped".
func cachePendingDeals(ctx context.Context, items []pendingDeal) {
	pendingIncomingDealsMu.Lock()
	defer pendingIncomingDealsMu.Unlock()

	remaining := maxPendingIncoming - len(pendingIncomingDeals)
	if remaining <= 0 {
		log.Warn(ctx, "Pending deals queue at capacity; dropping items", nil,
			"dropped", len(items),
			"capacity", maxPendingIncoming,
		)
		incPendingData(labelPendingDeals, labelPendingDropped, len(items))

		return
	}

	if len(items) > remaining {
		dropped := len(items) - remaining
		log.Warn(ctx, "Pending deals queue near capacity; dropping overflow", nil,
			"dropped", dropped,
			"capacity", maxPendingIncoming,
		)
		incPendingData(labelPendingDeals, labelPendingDropped, dropped)
		items = items[:remaining]
	}

	pendingIncomingDeals = append(pendingIncomingDeals, items...)
	incPendingData(labelPendingDeals, labelPendingCached, len(items))
}

// handleDKGProcessResponses handles the responses of processDeals from other committee members.
// shouldProcess must be pre-computed by the caller while the SDK context is available.
// Accepts []pendingResponse so that both first-time callers (via wrapResponses) and
// reprocess callers (with preserved retryCount) use the same code path.
func (k *Keeper) handleDKGProcessResponses(ctx context.Context, dkgNetwork *types.DKGNetwork, pending []pendingResponse, shouldProcess bool) {
	// Serialize kernel DKG operations to prevent concurrent DistKeyGenerator mutation.
	dkgKernelMu.Lock()
	defer dkgKernelMu.Unlock()

	log.Info(ctx, "Handling DKG process responses",
		"round", dkgNetwork.Round,
		"num_responses", len(pending),
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
	if phase := session.GetPhase(); phase != types.PhaseDealing && phase != types.PhaseInitialized {
		log.Warn(ctx, "Session not in dealing or initialized phase, skipping process responses", nil,
			"current_phase", phase.String(),
		)

		return
	}

	round := session.GetRound()
	isResharing := session.GetIsResharing()
	codeCommitment := session.GetCodeCommitment()
	oldCodeCommitment := session.GetOldCodeCommitment()
	selfIndex := session.GetIndex()

	// During upgrade resharing, validators in both old and new sets must send
	// ProcessResponses to BOTH binaries: the old binary (dealer role) and the new binary
	// (recipient role). Each binary maintains its own DKG state that needs updating.
	ccsToProcess := [][]byte{codeCommitment}
	if dkgNetwork.IsUpgrade && len(oldCodeCommitment) > 0 && !bytes.Equal(codeCommitment, oldCodeCommitment) {
		ccsToProcess = append(ccsToProcess, oldCodeCommitment)
	}

	// Filter self-responses: VssResponse.Index is 0-based (Kyber), selfIndex is 1-based.
	// If selfIndex is 0 (unset), include all responses (no self-filtering).
	filtered := make([]pendingResponse, 0, len(pending))
	for _, pr := range pending {
		if selfIndex == 0 || pr.response.VssResponse.Index != selfIndex-1 {
			filtered = append(filtered, pr)
		}
	}

	if len(filtered) == 0 {
		log.Info(ctx, "No responses to process. Skip to request")

		return
	}

	// Extract raw responses from filtered pending items for the kernel call.
	rawResponses := make([]types.Response, len(filtered))
	for i, pr := range filtered {
		rawResponses[i] = pr.response
	}

	var totalJustifications int
	// Union of rejected items across ccsToProcess (up to 2 in upgrade
	// resharing). retryCount is incremented exactly once per filtered item
	// regardless of how many CCs rejected it.
	rejected := make(map[rejectedKey]struct{})

	for _, cc := range ccsToProcess {
		var processResp *types.ProcessResponsesResponse

		start := time.Now()
		retryErr := retry(ctx, func(ctx context.Context) error {
			log.Info(ctx, "ProcessResponses call to kernel client",
				"round", round,
				"num_responses", len(rawResponses),
				"code_commitment", hex.EncodeToString(cc),
			)

			req := &types.ProcessResponsesRequest{
				CodeCommitment: cc,
				Round:          round,
				Responses:      rawResponses,
				IsResharing:    isResharing,
			}

			client, cErr := k.getClientWithReconnect(cc)
			if cErr != nil {
				return errors.Wrap(cErr, "no kernel client for session")
			}

			var rpcErr error
			processResp, rpcErr = client.ProcessResponses(ctx, req)
			if rpcErr != nil {
				return rpcErr
			}

			return nil
		})
		observeKernelCall(labelOpProcessResponses, start, retryErr)

		if retryErr != nil {
			log.Error(ctx, "Failed to process responses", retryErr,
				"code_commitment", hex.EncodeToString(cc),
			)

			// Batch RPC failure: mark every filtered item as rejected (we
			// cannot tell which items the kernel actually accepted).
			for i := range filtered {
				rejected[rejectedResponseKey(filtered[i].response)] = struct{}{}
			}

			continue
		}

		// Enqueue any justifications returned by story-kernel for broadcast via Vote Extension.
		// Justifications are produced when a complaint response (status=false) is processed
		// and the dealer needs to reveal the plaintext deal to prove its validity.
		if processResp != nil && len(processResp.GetJustifications()) > 0 {
			k.EnqueueJustifications(processResp.GetJustifications())
			totalJustifications += len(processResp.GetJustifications())

			log.Info(ctx, "Enqueued justifications for broadcast",
				"code_commitment", hex.EncodeToString(cc),
				"round", round,
				"num_justifications", len(processResp.GetJustifications()),
			)
		}

		// Union the kernel-reported rejections from this CC into the set.
		if processResp != nil {
			for _, r := range processResp.GetRejectedResponses() {
				rejected[rejectedResponseKey(r)] = struct{}{}
			}
		}
	}

	// Apply retry decisions exactly once (critical when ccsToProcess has 2
	// entries — same item rejected by both must increment retryCount once).
	if len(rejected) > 0 {
		var remaining []pendingResponse
		for i := range filtered {
			if _, isRejected := rejected[rejectedResponseKey(filtered[i].response)]; !isRejected {
				continue
			}
			filtered[i].retryCount++
			if filtered[i].retryCount > maxReprocessAttempts {
				log.Warn(ctx, "Dropping pending response after max retry attempts", nil,
					"round", dkgNetwork.Round,
					"response_index", filtered[i].response.Index,
					"retry_count", filtered[i].retryCount,
				)

				continue
			}
			remaining = append(remaining, filtered[i])
		}
		if len(remaining) > 0 {
			cachePendingResponses(ctx, remaining)
		}
	}

	log.Info(ctx, "Process responses complete",
		"round", round,
		"submitted_responses", len(rawResponses),
		"justifications_received", totalJustifications,
		"rejected_responses", len(rejected),
	)
}

// handleDKGProcessJustifications sends already-verified justifications to the kernel.
// Justifications have passed Schnorr signature and VSS verification in ProcessJustifications
// (FinalizeBlock context). This function only forwards them to the kernel via gRPC.
// Also used for replaying cached justifications that failed the kernel gRPC call.
// Accepts []pendingJustification so that both first-time callers (via wrapJustifications) and
// reprocess callers (with preserved retryCount) use the same code path.
func (k *Keeper) handleDKGProcessJustifications(ctx context.Context, dkgNetwork *types.DKGNetwork, pending []pendingJustification) {
	dkgKernelMu.Lock()
	defer dkgKernelMu.Unlock()

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session for justification processing", err)

		return
	}

	round := session.GetRound()
	isResharing := session.GetIsResharing()
	codeCommitment := session.GetCodeCommitment()
	oldCodeCommitment := session.GetOldCodeCommitment()

	ccsToProcess := [][]byte{codeCommitment}
	if dkgNetwork.IsUpgrade && len(oldCodeCommitment) > 0 && !bytes.Equal(codeCommitment, oldCodeCommitment) {
		ccsToProcess = append(ccsToProcess, oldCodeCommitment)
	}

	// Extract raw justifications from pending items for the kernel call.
	rawJustifications := make([]types.Justification, len(pending))
	for i, pj := range pending {
		rawJustifications[i] = pj.justification
	}

	// Union of rejected justifications across ccsToProcess (up to 2 in upgrade
	// resharing). retryCount is incremented exactly once per pending item.
	rejected := make(map[rejectedKey]struct{})

	for _, cc := range ccsToProcess {
		var processResp *types.ProcessJustificationResponse

		start := time.Now()
		retryErr := retry(ctx, func(ctx context.Context) error {
			req := &types.ProcessJustificationRequest{
				CodeCommitment: cc,
				Round:          round,
				Justifications: rawJustifications,
				IsResharing:    isResharing,
			}

			client, cErr := k.getClientWithReconnect(cc)
			if cErr != nil {
				return errors.Wrap(cErr, "no kernel client for session")
			}

			var rpcErr error
			processResp, rpcErr = client.ProcessJustification(ctx, req)
			if rpcErr != nil {
				return rpcErr
			}

			return nil
		})
		observeKernelCall(labelOpProcessJustifications, start, retryErr)

		if retryErr != nil {
			log.Error(ctx, "Failed to process justifications via kernel", retryErr,
				"round", dkgNetwork.Round,
				"code_commitment", hex.EncodeToString(cc),
			)

			// Batch RPC failure: mark every pending item as rejected.
			for i := range pending {
				rejected[rejectedJustificationKey(pending[i].justification)] = struct{}{}
			}

			continue
		}

		// Union the kernel-reported rejections from this CC into the set.
		if processResp != nil {
			for _, j := range processResp.GetRejectedJustifications() {
				rejected[rejectedJustificationKey(j)] = struct{}{}
			}
		}
	}

	// Apply retry decisions exactly once (critical when ccsToProcess has 2
	// entries — same item rejected by both must increment retryCount once).
	if len(rejected) > 0 {
		var remaining []pendingJustification
		for i := range pending {
			if _, isRejected := rejected[rejectedJustificationKey(pending[i].justification)]; !isRejected {
				continue
			}
			pending[i].retryCount++
			if pending[i].retryCount > maxReprocessAttempts {
				log.Warn(ctx, "Dropping pending justification after max retry attempts", nil,
					"round", dkgNetwork.Round,
					"justification_index", pending[i].justification.Index,
					"retry_count", pending[i].retryCount,
				)

				continue
			}
			remaining = append(remaining, pending[i])
		}
		if len(remaining) > 0 {
			cachePendingJustifications(ctx, remaining)
		}
	}

	log.Info(ctx, "Process justifications complete",
		"round", round,
		"submitted_justifications", len(rawJustifications),
		"rejected_justifications", len(rejected),
	)
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

// cachePendingResponses saves pending responses for later retry, respecting
// the capacity limit. The pending-responses metric is incremented by the
// number of items actually cached (post-capacity-trim). Items dropped because
// the queue is at capacity are logged and counted under op="dropped".
func cachePendingResponses(ctx context.Context, items []pendingResponse) {
	pendingIncomingResponsesMu.Lock()
	defer pendingIncomingResponsesMu.Unlock()

	remaining := maxPendingIncoming - len(pendingIncomingResponses)
	if remaining <= 0 {
		log.Warn(ctx, "Pending responses queue at capacity; dropping items", nil,
			"dropped", len(items),
			"capacity", maxPendingIncoming,
		)
		incPendingData(labelPendingResponses, labelPendingDropped, len(items))

		return
	}

	if len(items) > remaining {
		dropped := len(items) - remaining
		log.Warn(ctx, "Pending responses queue near capacity; dropping overflow", nil,
			"dropped", dropped,
			"capacity", maxPendingIncoming,
		)
		incPendingData(labelPendingResponses, labelPendingDropped, dropped)
		items = items[:remaining]
	}

	pendingIncomingResponses = append(pendingIncomingResponses, items...)
	incPendingData(labelPendingResponses, labelPendingCached, len(items))
}

// cachePendingJustifications saves pending justifications for later retry,
// respecting the capacity limit. The pending-justifications metric is
// incremented by the number of items actually cached (post-capacity-trim).
// Items dropped because the queue is at capacity are logged and counted
// under op="dropped".
func cachePendingJustifications(ctx context.Context, items []pendingJustification) {
	pendingIncomingJustificationsMu.Lock()
	defer pendingIncomingJustificationsMu.Unlock()

	remaining := maxPendingIncoming - len(pendingIncomingJustifications)
	if remaining <= 0 {
		log.Warn(ctx, "Pending justifications queue at capacity; dropping items", nil,
			"dropped", len(items),
			"capacity", maxPendingIncoming,
		)
		incPendingData(labelPendingJustifications, labelPendingDropped, len(items))

		return
	}

	if len(items) > remaining {
		dropped := len(items) - remaining
		log.Warn(ctx, "Pending justifications queue near capacity; dropping overflow", nil,
			"dropped", dropped,
			"capacity", maxPendingIncoming,
		)
		incPendingData(labelPendingJustifications, labelPendingDropped, dropped)
		items = items[:remaining]
	}

	pendingIncomingJustifications = append(pendingIncomingJustifications, items...)
	incPendingData(labelPendingJustifications, labelPendingCached, len(items))
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

	if k.kernelRouter == nil {
		return // kernel router not configured
	}

	if !k.kernelRouter.HasClients() {
		// No clients connected — attempt reconnection before giving up.
		k.kernelRouter.TryReconnect()

		if !k.kernelRouter.HasClients() {
			return // kernel still unavailable, wait
		}
	}

	asyncCtx, cancel := dkgAsyncContext()

	go func() {
		defer cancel()

		// Process deals FIRST — kyber returns ErrNoDealBeforeResponse if
		// a response arrives for a dealer whose deal hasn't been processed.
		if hasPendingDeals {
			drained := drainPendingIncomingDeals()
			log.Info(asyncCtx, "Replaying cached deals after kernel recovery",
				"round", dkgNetwork.Round, "num_deals", len(drained))
			k.handleDKGProcessDeals(asyncCtx, dkgNetwork, drained)
			incPendingData(labelPendingDeals, labelPendingReplayed, len(drained))
		}

		// Then process responses.
		if hasPendingResponses {
			drained := drainPendingIncomingResponses()
			log.Info(asyncCtx, "Replaying cached responses after kernel recovery",
				"round", dkgNetwork.Round, "num_responses", len(drained))
			k.handleDKGProcessResponses(asyncCtx, dkgNetwork, drained, true)
			incPendingData(labelPendingResponses, labelPendingReplayed, len(drained))
		}

		// Finally process justifications (already verified before caching,
		// so forward directly to kernel without re-verification).
		if hasPendingJustifications {
			drained := drainPendingIncomingJustifications()
			log.Info(asyncCtx, "Replaying cached justifications after kernel recovery",
				"round", dkgNetwork.Round, "num_justifications", len(drained))
			k.handleDKGProcessJustifications(asyncCtx, dkgNetwork, drained)
			incPendingData(labelPendingJustifications, labelPendingReplayed, len(drained))
		}
	}()
}

func drainPendingIncomingDeals() []pendingDeal {
	pendingIncomingDealsMu.Lock()
	defer pendingIncomingDealsMu.Unlock()

	out := pendingIncomingDeals
	pendingIncomingDeals = nil

	return out
}

func drainPendingIncomingResponses() []pendingResponse {
	pendingIncomingResponsesMu.Lock()
	defer pendingIncomingResponsesMu.Unlock()

	out := pendingIncomingResponses
	pendingIncomingResponses = nil

	return out
}

func drainPendingIncomingJustifications() []pendingJustification {
	pendingIncomingJustificationsMu.Lock()
	defer pendingIncomingJustificationsMu.Unlock()

	out := pendingIncomingJustifications
	pendingIncomingJustifications = nil

	return out
}

// wrapDeals converts raw deals into the first-attempt retry wrappers.
// retryCount is explicitly initialized to 0 for the first attempt.
func wrapDeals(deals []types.Deal) []pendingDeal {
	pd := make([]pendingDeal, len(deals))
	for i, d := range deals {
		pd[i] = pendingDeal{deal: d, retryCount: 0}
	}

	return pd
}

// wrapResponses converts raw responses into the first-attempt retry wrappers.
// retryCount is explicitly initialized to 0 for the first attempt.
func wrapResponses(responses []types.Response) []pendingResponse {
	pr := make([]pendingResponse, len(responses))
	for i, r := range responses {
		pr[i] = pendingResponse{response: r, retryCount: 0}
	}

	return pr
}

// wrapJustifications converts raw justifications into the first-attempt retry
// wrappers. retryCount is explicitly initialized to 0 for the first attempt.
func wrapJustifications(justifications []types.Justification) []pendingJustification {
	pj := make([]pendingJustification, len(justifications))
	for i, j := range justifications {
		pj[i] = pendingJustification{justification: j, retryCount: 0}
	}

	return pj
}

// wrapDecryptRequests converts raw decrypt requests into the first-attempt
// retry wrappers. RetryCount is explicitly initialized to 0 for the first attempt.
func wrapDecryptRequests(requests []types.DecryptRequest) []types.PendingDecryptRequest {
	pd := make([]types.PendingDecryptRequest, len(requests))
	for i, r := range requests {
		pd[i] = types.PendingDecryptRequest{DecryptRequest: r, RetryCount: 0}
	}

	return pd
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
