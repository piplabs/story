package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"slices"

	"cosmossdk.io/collections"
	"go.dedis.ch/kyber/v4/group/edwards25519"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGDealing handles the dealing phase event.
func (k *Keeper) handleDKGDealing(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG dealing",
		"round", dkgNetwork.Round,
	)

	if !dkgSvcRunning.CompareAndSwap(false, true) {
		log.Info(ctx, "DKG service already running; skipping dealing")

		return
	}
	defer dkgSvcRunning.Store(false)

	if dkgNetwork.Stage != types.DKGStageDealing {
		log.Info(ctx, "DKG Dealing is skipped because the current network stage is not dealing stage")

		return
	}

	shouldDeal, err := k.shouldDeal(ctx, dkgNetwork)
	if err != nil {
		log.Error(ctx, "Failed to check whether the validator should deal or not", err)

		return
	}

	if !shouldDeal {
		log.Debug(ctx, "Skip dealing")

		return
	}

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	if session.Phase != types.PhaseDealing {
		log.Warn(ctx, "Session not in dealing phase, skipping generate deals", nil,
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

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update session after generating deals", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	k.EnqueueDeals(resp.GetDeals())

	log.Info(ctx, "DKG deals are generated successfully",
		"round", session.Round,
	)

	return
}

// handleDKGProcessDeals handles the deals from other committee members.
func (k *Keeper) handleDKGProcessDeals(ctx context.Context, dkgNetwork *types.DKGNetwork, deals []types.Deal) {
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

	if session.Phase != types.PhaseDealing {
		log.Warn(ctx, "Session not in dealing phase, skipping process deals", nil,
			"current_phase", session.Phase.String(),
		)

		return
	}

	var resp *types.ProcessDealResponse
	if err := retry(ctx, func(ctx context.Context) error {
		log.Info(ctx, "ProcessDeals call to kernel client",
			"round", session.Round,
			"num_deals", len(deals),
		)

		req := &types.ProcessDealRequest{
			CodeCommitment: session.CodeCommitment,
			Round:          session.Round,
			Deals:          []types.Deal{},
			IsResharing:    session.IsResharing,
		}

		for _, deal := range deals {
			if deal.RecipientIndex == session.Index {
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
		log.Error(ctx, "Failed to process deals", err)

		return
	}

	k.EnqueueResponses(resp.GetResponses())

	log.Info(ctx, "Process deals complete",
		"round", session.Round,
	)

	return
}

// handleDKGProcessResponses handles the responses of processDeals from other committee members.
func (k *Keeper) handleDKGProcessResponses(ctx context.Context, dkgNetwork *types.DKGNetwork, responses []types.Response) {
	log.Info(ctx, "Handling DKG process responses",
		"round", dkgNetwork.Round,
		"num_responses", len(responses),
	)

	shouldProcess, err := k.shouldProcessResponses(ctx, dkgNetwork)
	if err != nil {
		log.Error(ctx, "Failed to check whether the validator should process responses", err)

		return
	}

	if !shouldProcess {
		log.Info(ctx, "Skip processing of responses")

		return
	}

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session", err)

		return
	}

	if session.Phase != types.PhaseDealing {
		log.Warn(ctx, "Session not in dealing phase, skipping process responses", nil,
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
		if resp.VssResponse.Index != session.Index {
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

	return
}

func (k *Keeper) shouldDeal(ctx context.Context, dkgNetwork *types.DKGNetwork) (bool, error) {
	inCurSet := slices.Contains(dkgNetwork.ActiveValSet, k.validatorEVMAddr)

	inPrevSet, err := k.isInPrevActiveValSet(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// First DKG round: only current set of validators deal
			return inCurSet, nil
		}

		return false, err
	}

	// Resharing round: only previous set of validators deal (they hold the existing key shares)
	return inPrevSet, nil
}

// handleDKGProcessJustifications verifies and forwards valid justifications to story-kernel
// for DKG state restoration. This is the off-chain (goroutine) handler that performs:
//  1. Empty check and MaxJustificationsPerBlock cap
//  2. ActiveValSet / session / phase check
//  3. Schnorr signature verification → deduplication → Pedersen VSS verification
//  4. Forward only valid justifications to story-kernel
func (k *Keeper) handleDKGProcessJustifications(ctx context.Context, dkgNetwork *types.DKGNetwork, justifications []types.Justification) {
	log.Info(ctx, "Handling DKG process justifications",
		"round", dkgNetwork.Round,
		"num_justifications", len(justifications),
	)

	if len(justifications) == 0 {
		return
	}

	// Cap justifications per block to prevent resource exhaustion
	if len(justifications) > MaxJustificationsPerBlock {
		log.Warn(ctx, "Justification count exceeds max, truncating", nil,
			"count", len(justifications),
			"max", MaxJustificationsPerBlock,
		)
		justifications = justifications[:MaxJustificationsPerBlock]
	}

	if !slices.Contains(dkgNetwork.ActiveValSet, k.validatorEVMAddr) {
		log.Info(ctx, "Skip processing justifications as the validator is not in current round set")

		return
	}

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session", err)

		return
	}

	if session.Phase != types.PhaseDealing {
		log.Warn(ctx, "Session not in dealing phase, skipping process justifications", nil,
			"current_phase", session.Phase.String(),
		)

		return
	}

	suite := edwards25519.NewBlakeSHA256Ed25519()

	// Build dealer public key map once for all justifications (avoids O(N) registration scan per justification).
	dealerPubKeys, err := k.buildDealerPubKeyMap(ctx, dkgNetwork, suite)
	if err != nil {
		log.Error(ctx, "Failed to build dealer public key map", err)

		return
	}

	// Step 1: Verify Schnorr signature FIRST, filter out unsigned/forged justifications.
	// This MUST happen before deduplication so that an attacker cannot preempt a valid
	// justification by broadcasting an unsigned one with the same (dealerIndex, recipientIndex).
	var signatureVerified []types.Justification
	for _, j := range justifications {
		if err := verifyJustificationSignature(suite, j, dealerPubKeys); err != nil {
			log.Debug(ctx, "Dropping justification with invalid signature",
				"dealer_index", j.Index,
				"error", err,
			)

			continue
		}

		signatureVerified = append(signatureVerified, j)
	}

	// Step 2: Deduplicate by (dealerIndex, recipientIndex) — now only signature-verified entries.
	deduped := deduplicateJustifications(signatureVerified)

	// Step 3: Perform Pedersen VSS verification on each deduplicated justification.
	// Invalid justifications are silently dropped (not errors).
	var validJustifications []types.Justification
	for _, j := range deduped {
		valid, err := verifyJustification(dkgNetwork, j)
		if err != nil {
			log.Warn(ctx, "Justification VSS verification error, dropping", err,
				"dealer_index", j.Index,
			)

			continue
		}

		if !valid {
			log.Info(ctx, "Justification VSS verification failed (deal was invalid), dropping",
				"dealer_index", j.Index,
			)

			continue
		}

		validJustifications = append(validJustifications, j)
	}

	if len(validJustifications) == 0 {
		log.Info(ctx, "No valid justifications after verification, skipping story-kernel call")

		return
	}

	// During upgrade resharing, justifications must be forwarded to BOTH old and new
	// binaries so each can update its DKG state accordingly.
	ccsToProcess := [][]byte{session.CodeCommitment}
	if dkgNetwork.IsUpgrade && len(session.OldCodeCommitment) > 0 && !bytes.Equal(session.CodeCommitment, session.OldCodeCommitment) {
		ccsToProcess = append(ccsToProcess, session.OldCodeCommitment)
	}

	for _, cc := range ccsToProcess {
		if err := retry(ctx, func(ctx context.Context) error {
			log.Info(ctx, "ProcessJustification batch call to story-kernel client",
				"code_commitment", hex.EncodeToString(cc),
				"round", session.Round,
				"num_justifications", len(validJustifications),
			)

			req := &types.ProcessJustificationRequest{
				CodeCommitment: cc,
				Round:          session.Round,
				Justifications: validJustifications,
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
			log.Error(ctx, "Failed to process justifications", err,
				"code_commitment", hex.EncodeToString(cc),
			)

			continue
		}
	}

	log.Info(ctx, "Process justifications complete",
		"round", session.Round,
	)
}

func (k *Keeper) shouldProcessResponses(ctx context.Context, dkgNetwork *types.DKGNetwork) (bool, error) {
	inCurSet := slices.Contains(dkgNetwork.ActiveValSet, k.validatorEVMAddr)

	inPrevSet, err := k.isInPrevActiveValSet(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// First DKG round: only current set of validators process deals or responses
			return inCurSet, nil
		}

		return false, err
	}

	// Resharing round: both previous and current set of validators process deals or responses
	return inCurSet || inPrevSet, nil
}
