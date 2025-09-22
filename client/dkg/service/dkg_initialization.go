package service

import (
	"context"
	"encoding/hex"
	"slices"

	dkgpb "github.com/piplabs/story/client/dkg/pb/v1"
	"github.com/piplabs/story/client/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGInitialization handles the DKG initialization event.
func (s *Service) handleDKGInitialization(ctx context.Context, event *types.DKGEventData) error {
	valAddr := s.validatorAddress.Hex()
	isParticipant := slices.Contains(event.ActiveValidators, valAddr)
	if !isParticipant {
		log.Debug(ctx, "Validator is not part of DKG committee", "validator", valAddr)

		return nil
	}

	mrenclave, err := event.ParseMrenclave()
	if err != nil {
		return errors.Wrap(err, "failed to parse mrenclave")
	}

	session := types.NewDKGSession(mrenclave, event.Round, event.ActiveValidators)
	if err := s.stateManager.CreateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to create DKG session")
	}

	return s.startDKGRegistration(ctx, session)
}

// startDKGRegistration begins the DKG registration process.
func (s *Service) startDKGRegistration(ctx context.Context, session *types.DKGSession) error {
	log.Info(ctx, "Starting DKG registration",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
	)

	session.UpdatePhase(types.PhaseRegistering)

	log.Info(ctx, "GenerateAndSealKey call to TEE client",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"validator", s.validatorAddress.Hex(),
	)

	req := &dkgpb.GenerateAndSealKeyRequest{
		Address:   s.validatorAddress.Hex(),
		Mrenclave: session.Mrenclave,
		Round:     session.Round,
	}
	resp, err := s.teeClient.GenerateAndSealKey(ctx, req)
	if err != nil {
		session.UpdatePhase(types.PhaseFailed)
		if updateErr := s.stateManager.UpdateSession(ctx, session); updateErr != nil {
			log.Error(ctx, "Failed to update session after TEE error", updateErr)
		}

		return errors.Wrap(err, "failed to generate and seal key")
	}

	log.Info(ctx, "InitializeDKG contract call",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"dkg_pub_key", hex.EncodeToString(resp.GetDkgPubKey()),
		"comm_pub_key", hex.EncodeToString(resp.GetCommPubKey()),
		"raw_quote_len", len(resp.GetRawQuote()),
	)

	_, err = s.contractClient.InitializeDKG(ctx, session.Round, session.Mrenclave, resp.GetDkgPubKey(), resp.GetCommPubKey(), resp.GetRawQuote())
	if err != nil {
		return errors.Wrap(err, "failed to call InitializeDKG contract method")
	}

	return nil
}
