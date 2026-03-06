package keeper

import (
	"context"
	"encoding/hex"
	"slices"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGFinalization handles the finalization phase event.
func (k *Keeper) handleDKGFinalization(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG finalization",
		"round", dkgNetwork.Round,
	)

	if !dkgSvcRunning.CompareAndSwap(false, true) {
		log.Info(ctx, "DKG service already running; skipping finalization")

		return
	}
	defer dkgSvcRunning.Store(false)

	if dkgNetwork.Stage != types.DKGStageFinalization {
		log.Info(ctx, "DKG Finalization is skipped because the current network stage is not in the finalization stage")

		return
	}

	isInCurRoundSet := slices.Contains(dkgNetwork.ActiveValSet, k.validatorEVMAddr)
	if !isInCurRoundSet {
		log.Info(ctx, "Skip finalizing DKG network as the validator is not in current round set")

		return
	}

	session, err := k.stateManager.GetSession(dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	if session.Phase != types.PhaseDealing {
		log.Warn(ctx, "Session not in dealing phase, skipping finalize DKG", nil,
			"current_phase", session.Phase.String())
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	if err := k.callTEEFinalizeDKG(ctx, session); err != nil {
		log.Error(ctx, "Failed to finalize DKG network", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	if err := k.callContractFinalizeDKG(ctx, session); err != nil {
		log.Error(ctx, "Failed to call finalizeDKG method", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	session.UpdatePhase(types.PhaseFinalized)

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update session after calling finalizeDKG method", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	log.Info(ctx, "DKG finalization phase complete",
		"round", session.Round,
	)
}

func (k *Keeper) callTEEFinalizeDKG(ctx context.Context, session *types.DKGSession) error {
	log.Info(ctx, "Finalize call to kernel client",
		"round", session.Round,
	)

	if len(session.GlobalPubKey) > 0 && len(session.SigFinalizeNetwork) > 0 {
		log.Info(ctx, "DKG network already finalized in kernel client, skipping call Finalize request")

		return nil
	}

	var (
		resp *types.FinalizeDKGResponse
		err  error
	)
	if err := retry(ctx, func(ctx context.Context) error {
		req := &types.FinalizeDKGRequest{
			CodeCommitment: session.CodeCommitment,
			Round:          session.Round,
			IsResharing:    session.IsResharing,
		}

		client, cErr := k.kernelRouter.GetClient(session.CodeCommitment)
		if cErr != nil {
			return errors.Wrap(cErr, "no kernel client for session")
		}

		resp, err = client.FinalizeDKG(ctx, req)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "kernel client Finalize request failed")
	}

	session.ParticipantsRoot = resp.GetParticipantsRoot()
	session.GlobalPubKey = resp.GetGlobalPubKey()
	session.SigFinalizeNetwork = resp.GetSignature()
	session.PublicCoeffs = resp.GetPublicCoeffs()

	session.PubKeyShare = resp.GetPubKeyShare()
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to update session after calling Finalize on the kernel client")
	}

	return nil
}

func (k *Keeper) callContractFinalizeDKG(ctx context.Context, session *types.DKGSession) error {
	log.Info(ctx, "Finalize contract call",
		"round", session.Round,
		"global_pub_key", hex.EncodeToString(session.GlobalPubKey),
		"signature_len", len(session.SigFinalizeNetwork),
	)

	if _, err := k.contractClient.Finalize(
		ctx,
		session.Round,
		session.EnclaveType,
		session.ParticipantsRoot,
		session.GlobalPubKey,
		session.PublicCoeffs,
		session.PubKeyShare,
		session.SigFinalizeNetwork,
	); err != nil {
		return err
	}

	return nil
}
