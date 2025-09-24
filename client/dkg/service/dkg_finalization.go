package service

import (
	"context"

	dkgpb "github.com/piplabs/story/client/dkg/pb/v1"
	"github.com/piplabs/story/client/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGFinalization handles the finalization phase event.
func (s *Service) handleDKGFinalization(ctx context.Context, event *types.DKGEventData) error {
	log.Info(ctx, "Handling DKG finalization phase event",
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
		log.Warn(ctx, "Session not in dealing phase, skipping finalization", nil,
			"current_phase", session.Phase.String(),
		)

		return nil
	}

	session.UpdatePhase(types.PhaseFinalizing)

	finalizeReq := &dkgpb.FinalizeDKGRequest{
		Mrenclave: session.Mrenclave,
		Round:     session.Round,
	}

	finalizeResp, err := s.teeClient.FinalizeDKG(ctx, finalizeReq)
	if err != nil {
		session.UpdatePhase(types.PhaseFailed)
		if updateErr := s.stateManager.UpdateSession(ctx, session); updateErr != nil {
			log.Error(ctx, "Failed to update session after TEE finalization error", updateErr)
		}

		return errors.Wrap(err, "failed to finalize DKG")
	}

	if !finalizeResp.GetFinalized() {
		log.Warn(ctx, "DKG finalization failed, marking session as failed", nil)
		session.UpdatePhase(types.PhaseFailed)

		return s.stateManager.UpdateSession(ctx, session)
	}

	session.IsFinalized = true
	if err := s.stateManager.UpdateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to update session")
	}

	if err := s.submitFinalizeDKG(ctx, session, finalizeResp.GetSignature()); err != nil {
		return errors.Wrap(err, "failed to submit DKG finalization to blockchain")
	}

	log.Info(ctx, "DKG finalization phase complete",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"finalized", finalizeResp.GetFinalized(),
	)

	return nil
}

// submitFinalizeDKG submits the DKG finalization transaction to the DKG contract.
func (s *Service) submitFinalizeDKG(ctx context.Context, session *types.DKGSession, signature []byte) error {
	log.Info(ctx, "Submitting FinalizeDKG transaction to contract",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"finalized", session.IsFinalized,
	)

	_, err := s.contractClient.FinalizeDKG(
		ctx,
		session.Round,
		session.Mrenclave,
		session.GlobalPubKey,
		signature,
	)
	if err != nil {
		return errors.Wrap(err, "failed to call FinalizeDKG contract method")
	}

	log.Info(ctx, "FinalizeDKG contract call successful",
		"round", session.Round,
		"finalized", session.IsFinalized,
		"signature_len", len(signature),
	)

	return nil
}
