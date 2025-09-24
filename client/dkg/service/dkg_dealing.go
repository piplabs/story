package service

import (
	"context"
	"sync"

	dkgpb "github.com/piplabs/story/client/dkg/pb/v1"
	"github.com/piplabs/story/client/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

var (
	dealsMu   sync.Mutex
	Deals     []*types.Deal
	respsMu   sync.Mutex
	Responses []*types.Response
)

// todo: limit the number popped each time
func PopDeals() []*types.Deal {
	dealsMu.Lock()
	defer dealsMu.Unlock()
	out := Deals
	Deals = nil
	return out
}

// todo: limit the number popped each time
func PopResponses() []*types.Response {
	respsMu.Lock()
	defer respsMu.Unlock()
	out := Responses
	Responses = nil
	return out
}

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

	if session.Phase != types.PhaseChallenging {
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

	// todo: the slice length should be restricted by the max tx size
	dealsMu.Lock()
	defer dealsMu.Unlock()
	Deals = []*types.Deal{}
	for _, deal := range resp.GetDeals() {
		session.Deals[deal.Index] = *deal
		Deals = append(Deals, deal)
	}
	return nil
}

// handleDKGProcessDeals handles the deal verification phase event.
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
		Deals:     event.Deals,
	}
	resp, err := s.teeClient.ProcessDeals(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to process deals")
	}
	respsMu.Lock()
	defer respsMu.Unlock()
	Responses = resp.GetResponses()

	return nil
}

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
		Responses: event.Responses,
	}
	_, err = s.teeClient.ProcessResponses(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to process responses")
	}
	return nil
}
