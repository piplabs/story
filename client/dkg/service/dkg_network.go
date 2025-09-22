package service

import (
	"context"

	dkgpb "github.com/piplabs/story/client/dkg/pb/v1"
	"github.com/piplabs/story/client/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGNetworkSet handles the network set event (after registration period).
//
// TODO: before setting the network, ensure that `session.Registrations` is updated to only DKG validators who are verified.
func (s *Service) handleDKGNetworkSet(ctx context.Context, event *types.DKGEventData) error {
	log.Info(ctx, "Handling DKG network set event",
		"mrenclave", event.Mrenclave,
		"round", event.Round,
		"total", event.Total,
		"threshold", event.Threshold,
	)

	mrenclave, err := event.ParseMrenclave()
	if err != nil {
		return errors.Wrap(err, "failed to parse mrenclave")
	}

	session, err := s.stateManager.GetSession(mrenclave, event.Round)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG session")
	}

	dkgRegistrations, err := s.queryVerifiedDKGRegistrations(ctx, mrenclave, event.Round)
	if err != nil {
		return errors.Wrap(err, "failed to get verified dkg registrations from x/dkg module")
	}

	req := &dkgpb.SetupDKGNetworkRequest{
		Mrenclave:     session.Mrenclave,
		Round:         session.Round,
		Total:         event.Total,
		Threshold:     event.Threshold,
		Registrations: dkgRegistrations,
	}
	resp, err := s.teeClient.SetupDKGNetwork(ctx, req)
	if err != nil {
		session.UpdatePhase(types.PhaseFailed)
		if updateErr := s.stateManager.UpdateSession(ctx, session); updateErr != nil {
			log.Error(ctx, "Failed to update session after TEE error", updateErr)
		}

		return errors.Wrap(err, "failed to create deals")
	}

	session.ActiveValidators = event.ActiveValidators
	session.Total = event.Total
	session.Threshold = event.Threshold
	session.Index = resp.GetIndex()
	session.Commitments = resp.GetCommitments()

	session.UpdatePhase(types.PhaseChallenging)
	if err := s.stateManager.UpdateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to update session")
	}

	log.Info(ctx, "DKG network set complete, entering dealing phase",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"index", session.Index,
		"total", session.Total,
		"threshold", session.Threshold,
	)

	return nil
}
