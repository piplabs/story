package service

import (
	"context"

	"github.com/piplabs/story/client/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGComplete handles the DKG completion event.
func (s *Service) handleDKGComplete(ctx context.Context, event *types.DKGEventData) error {
	log.Info(ctx, "Handling DKG completion event",
		"mrenclave", event.Mrenclave,
		"round", event.Round,
	)

	mrenclave, err := event.ParseMrenclave()
	if err != nil {
		return errors.Wrap(err, "failed to parse mrenclave")
	}

	session, err := s.stateManager.GetSession(mrenclave, event.Round)
	if err != nil {
		log.Debug(ctx, "DKG session not found locally, validator may not have participated")

		return errors.Wrap(err, "failed to get dkg session")
	}

	session.UpdatePhase(types.PhaseCompleted)
	session.IsFinalized = true

	if err := s.stateManager.UpdateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to update completed session")
	}

	log.Info(ctx, "DKG process completed successfully",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"validator", s.validatorAddress.Hex(),
	)

	return nil
}
