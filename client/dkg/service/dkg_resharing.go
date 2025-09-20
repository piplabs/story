package service

import (
	"context"

	"github.com/piplabs/story/client/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGResharing handles DKG resharing for new epochs.
func (s *Service) handleDKGResharing(ctx context.Context, event *types.DKGEventData) error {
	log.Info(ctx, "Handling DKG resharing event",
		"mrenclave", event.Mrenclave,
		"round", event.Round,
	)

	mrenclave, err := event.ParseMrenclave()
	if err != nil {
		return errors.Wrap(err, "failed to parse mrenclave")
	}

	valAddr := s.validatorAddress.Hex()
	isParticipant := false
	for _, addr := range event.ActiveValidators {
		if addr == valAddr {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		log.Debug(ctx, "Validator not part of resharing committee")
		return nil
	}

	session := types.NewDKGSession(mrenclave, event.Round, event.ActiveValidators)

	if err := s.stateManager.CreateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to create resharing session")
	}

	if err := s.startDKGRegistration(ctx, session); err != nil {
		return errors.Wrap(err, "failed to start DKG registration for resharing")
	}

	log.Info(ctx, "DKG resharing initiated successfully",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
	)

	return nil
}
