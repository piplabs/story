package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/log"
)

// handleDKGDealing handles the dealing phase event.
func (k *Keeper) handleDKGDealing(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG dealing",
		"mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave),
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

	session, err := k.stateManager.GetSession(dkgNetwork.Mrenclave, dkgNetwork.Round)
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
			"mrenclave", session.GetMrenclaveString(),
			"round", session.Round,
		)

		req := &types.GenerateDealsRequest{
			Mrenclave: session.Mrenclave,
			Round:     session.Round,
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
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
	)

	return
}

// handleDKGProcessDeals handles the deals from other committee members.
func (k *Keeper) handleDKGProcessDeals(ctx context.Context, dkgNetwork *types.DKGNetwork, deals []types.Deal) {
	log.Info(ctx, "Handling DKG process deals",
		"mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave),
		"round", dkgNetwork.Round,
		"num_deals", len(deals),
	)

	session, err := k.stateManager.GetSession(dkgNetwork.Mrenclave, dkgNetwork.Round)
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
			"mrenclave", session.GetMrenclaveString(),
			"round", session.Round,
			"num_deals", len(deals),
		)

		req := &types.ProcessDealRequest{
			Mrenclave: session.Mrenclave,
			Round:     session.Round,
			Deals:     []types.Deal{},
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
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
	)

	return
}

// handleDKGProcessResponses handles the responses of processDeals from other committee members.
func (k *Keeper) handleDKGProcessResponses(ctx context.Context, dkgNetwork *types.DKGNetwork, responses []types.Response) {
	log.Info(ctx, "Handling DKG process responses",
		"mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave),
		"round", dkgNetwork.Round,
		"num_responses", len(responses),
	)

	session, err := k.stateManager.GetSession(dkgNetwork.Mrenclave, dkgNetwork.Round)
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
			"mrenclave", session.GetMrenclaveString(),
			"round", session.Round,
			"num_deals", len(responses),
		)

		req := &types.ProcessResponsesRequest{
			Mrenclave: session.Mrenclave,
			Round:     session.Round,
			Responses: []types.Response{},
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
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
	)

	return
}
