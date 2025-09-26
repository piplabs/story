package service

import (
	"context"

	dkgpb "github.com/piplabs/story/client/dkg/pb/v1"
	"github.com/piplabs/story/client/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGDealing handles the dealing phase event.
func (s *Service) handleDKGDealing(ctx context.Context, event *types.DKGEventData) error {
	log.Info(ctx, "Handling DKG dealing phase event",
		"mrenclave", event.Mrenclave,
		"round", event.Round,
	)

	mrenclave, err := event.ParseMrenclave()
	if err != nil {
		return errors.Wrap(err, "failed to parse mrenclave")
	}

	session, err := s.stateManager.GetSession(mrenclave, event.Round)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG session")
	}

	if session.Phase != types.PhaseDealing {
		log.Warn(ctx, "Session not in challenge phase, skipping dealing", nil,
			"current_phase", session.Phase.String(),
		)

		return nil
	}

	req := &dkgpb.GenerateDealsRequest{
		Mrenclave: session.Mrenclave,
		Round:     session.Round,
	}
	resp, err := s.teeClient.GenerateDeals(ctx, req)
	if err != nil {
		session.UpdatePhase(types.PhaseFailed)
		if updateErr := s.stateManager.UpdateSession(ctx, session); updateErr != nil {
			log.Error(ctx, "Failed to update session after TEE error", updateErr)
		}

		return errors.Wrap(err, "failed to create deals")
	}

	for _, deal := range resp.GetDeals() {
		session.Deals[deal.Index] = *deal
		s.index = deal.Index // same for all deals
	}

	err = AddDealsFile(resp.GetDeals())
	if err != nil {
		return errors.Wrap(err, "failed to add deals to file")
	}

	return nil
}

// handleDKGProcessDeals handles the deals from other committee members.
func (s *Service) handleDKGProcessDeals(ctx context.Context, event *types.DKGEventData) error {
	log.Info(ctx, "Handling DKG deal verification phase event",
		"mrenclave", event.Mrenclave,
		"round", event.Round,
	)

	mrenclave, err := event.ParseMrenclave()
	if err != nil {
		return errors.Wrap(err, "failed to parse mrenclave")
	}

	session, err := s.stateManager.GetSession(mrenclave, event.Round)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG session")
	}

	if session.Phase != types.PhaseDealing {
		log.Warn(ctx, "Session not in dealing phase, skipping deal verification", nil,
			"current_phase", session.Phase.String(),
		)

		return nil
	}

	req := &dkgpb.ProcessDealRequest{
		Mrenclave: session.Mrenclave,
		Round:     session.Round,
		Index:     session.Index,
		Deals:     []*types.Deal{},
	}

	for _, deal := range event.Deals {
		if deal.RecipientIndex == s.index {
			req.Deals = append(req.Deals, deal)
		}
	}
	log.Info(ctx, "Processing deals", "event_deals", len(event.Deals), "req_deals", len(req.GetDeals()), "index", s.index)

	resp, err := s.teeClient.ProcessDeals(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to process deals")
	}

	err = AddResponsesFile(resp.GetResponses())
	if err != nil {
		return errors.Wrap(err, "failed to add responses to file")
	}

	return nil
}

// handleDKGProcessResponses handles the responses of processDeals from other committee members.
func (s *Service) handleDKGProcessResponses(ctx context.Context, event *types.DKGEventData) error {
	log.Info(ctx, "Handling DKG process responses event",
		"mrenclave", event.Mrenclave,
		"round", event.Round,
	)

	mrenclave, err := event.ParseMrenclave()
	if err != nil {
		return errors.Wrap(err, "failed to parse mrenclave")
	}

	session, err := s.stateManager.GetSession(mrenclave, event.Round)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG session")
	}

	if session.Phase != types.PhaseDealing {
		log.Warn(ctx, "Session not in dealing phase, skipping process responses", nil,
			"current_phase", session.Phase.String(),
		)

		return nil
	}

	req := &dkgpb.ProcessResponsesRequest{
		Mrenclave: session.Mrenclave,
		Round:     session.Round,
		Responses: []*types.Response{},
	}

	for _, resp := range event.Responses {
		if resp.Index != s.index {
			req.Responses = append(req.Responses, resp)
		}
	}
	log.Info(ctx, "Processing responses", "event_responses", len(event.Responses), "req_responses", len(req.GetResponses()), "index", s.index)

	_, err = s.teeClient.ProcessResponses(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to process responses")
	}

	return nil
}
