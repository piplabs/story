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

	for _, deal := range resp.GetDeals() {
		session.Deals[deal.Index] = *deal
	}

	// TODO: dealing logic
	// DKG service would request to process the deals which is given through vote extension to the TEE client.
	// Then, the TEE client will process the deals. If there is any invalid deals, it will return the complaints in ProcessDealResponse (step 5 in DKG dealing section).

	return nil
}
