package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
	"slices"
)

// handleDKGDealing handles the dealing phase event.
func (k *Keeper) handleDKGDealing(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG dealing",
		"code_commitment", hex.EncodeToString(dkgNetwork.CodeCommitment),
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

	session, err := k.stateManager.GetSession(dkgNetwork.CodeCommitment, dkgNetwork.Round)
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

	var resp *types.GenerateDealsResponse
	if err := retry(ctx, func(ctx context.Context) error {
		log.Info(ctx, "GenerateDeals call to TEE client",
			"code_commitment", session.GetCodeCommitmentString(),
			"round", session.Round,
		)

		req := &types.GenerateDealsRequest{
			CodeCommitment: session.CodeCommitment,
			Round:          session.Round,
			IsResharing:    session.IsResharing,
		}
		resp, err = k.teeClient.GenerateDeals(ctx, req)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Error(ctx, "Failed to generate deals", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	for _, deal := range resp.GetDeals() {
		session.Deals[deal.Index] = deal
		session.Index = deal.Index // same for all deals
	}

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update session after generating deals", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	k.EnqueueDeals(resp.GetDeals())

	log.Info(ctx, "DKG deals are generated successfully",
		"code_commitment", session.GetCodeCommitmentString(),
		"round", session.Round,
	)

	return
}

// handleDKGProcessDeals handles the deals from other committee members.
func (k *Keeper) handleDKGProcessDeals(ctx context.Context, dkgNetwork *types.DKGNetwork, deals []types.Deal) {
	log.Info(ctx, "Handling DKG process deals",
		"code_commitment", hex.EncodeToString(dkgNetwork.CodeCommitment),
		"round", dkgNetwork.Round,
		"num_deals", len(deals),
	)

	if !slices.Contains(dkgNetwork.ActiveValSet, k.validatorEVMAddr) {
		log.Info(ctx, "Skip processing deals as the validator is not in current round set")

		return
	}

	session, err := k.stateManager.GetSession(dkgNetwork.CodeCommitment, dkgNetwork.Round)
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
		log.Info(ctx, "ProcessDeals call to TEE client",
			"code_commitment", session.GetCodeCommitmentString(),
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

		resp, err = k.teeClient.ProcessDeals(ctx, req)
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
		"code_commitment", session.GetCodeCommitmentString(),
		"round", session.Round,
	)

	return
}

// handleDKGProcessResponses handles the responses of processDeals from other committee members.
func (k *Keeper) handleDKGProcessResponses(ctx context.Context, dkgNetwork *types.DKGNetwork, responses []types.Response) {
	log.Info(ctx, "Handling DKG process responses",
		"code_commitment", hex.EncodeToString(dkgNetwork.CodeCommitment),
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

	session, err := k.stateManager.GetSession(dkgNetwork.CodeCommitment, dkgNetwork.Round)
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

	if err := retry(ctx, func(ctx context.Context) error {
		log.Info(ctx, "ProcessResponses call to TEE client",
			"code_commitment", session.GetCodeCommitmentString(),
			"round", session.Round,
			"num_deals", len(responses),
		)

		req := &types.ProcessResponsesRequest{
			CodeCommitment: session.CodeCommitment,
			Round:          session.Round,
			Responses:      []types.Response{},
			IsResharing:    session.IsResharing,
		}

		for _, resp := range responses {
			if resp.VssResponse.Index != session.Index {
				req.Responses = append(req.Responses, resp)
			}
		}

		if len(req.Responses) == 0 {
			log.Info(ctx, "No responses to process. Skip to request")

			return nil
		}

		if _, err := k.teeClient.ProcessResponses(ctx, req); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Error(ctx, "Failed to process responses", err)

		return
	}

	log.Info(ctx, "Process responses complete",
		"code_commitment", session.GetCodeCommitmentString(),
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

	// Resharing round: only current set of validators deal
	return inPrevSet, nil
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
