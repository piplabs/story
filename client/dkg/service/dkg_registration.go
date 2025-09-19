package service

import (
	"context"

	"github.com/piplabs/story/client/dkg/types"
	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGRegistrationInitialized handles the DKGInitialized event.
func (s *Service) handleDKGRegistrationInitialized(ctx context.Context, event *types.DKGEventData) error {
	log.Info(ctx, "Handling DKG registration initialized event",
		"mrenclave", event.Mrenclave,
		"round", event.Round,
		"msg_sender", event.ValidatorAddr,
		"dkg_pub_key_len", len(event.DkgPubKey),
		"comm_pub_key_len", len(event.CommPubKey),
		"raw_quote_len", len(event.RawQuote),
	)

	mrenclave, err := event.ParseMrenclave()
	if err != nil {
		return errors.Wrap(err, "failed to parse mrenclave")
	}

	session, err := s.stateManager.GetSession(mrenclave, event.Round)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG session")
	}

	// TODO: ensure no duplicate registration within each session (dup by msgSender)
	session.Registrations = append(session.Registrations, dkgtypes.DKGRegistration{
		Mrenclave:   mrenclave,
		Round:       event.Round,
		MsgSender:   event.ValidatorAddr,
		Index:       event.Index,
		DkgPubKey:   event.DkgPubKey,
		CommPubKey:  event.CommPubKey,
		RawQuote:    event.RawQuote,
		Status:      dkgtypes.DKGRegStatusNotVerified,
		Commitments: []byte{},
	})

	if err := s.stateManager.UpdateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to update session")
	}

	log.Info(ctx, "DKG registration initialized as unverified",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"msg_sender", event.ValidatorAddr,
		"index", session.Index,
		"dkg_pub_key_len", len(event.DkgPubKey),
		"comm_pub_key_len", len(event.CommPubKey),
		"raw_quote_len", len(event.RawQuote),
	)

	return nil
}

// handleDKGRegistrationCommitmentsUpdated handles commitment update event and sets the DKG registration status.
func (s *Service) handleDKGRegistrationCommitmentsUpdated(ctx context.Context, event *types.DKGEventData) error {
	log.Info(ctx, "Handling DKG registration commitments updated event",
		"mrenclave", event.Mrenclave,
		"round", event.Round,
		"index", event.Index,
		"msg_sender", event.ValidatorAddr,
		"commitments_len", len(event.Commitments),
		"signature_len", len(event.Signature),
	)

	mrenclave, err := event.ParseMrenclave()
	if err != nil {
		return errors.Wrap(err, "failed to parse mrenclave")
	}

	// TODO: assume the commitments and signature are valid, and that we can mark the registration as verified
	session, err := s.stateManager.GetSession(mrenclave, event.Round)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG session")
	}

	for i, registration := range session.Registrations {
		if registration.Index == event.Index {
			session.Registrations[i].Status = dkgtypes.DKGRegStatusVerified
			break
		}
	}

	if err := s.stateManager.UpdateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to update session")
	}

	return nil
}
