package keeper

import (
	"context"
	"encoding/hex"
	"slices"
	"time"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGFinalization handles the finalization phase event.
func (k *Keeper) handleDKGFinalization(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG finalization",
		"round", dkgNetwork.Round,
	)

	if !tryAcquireDKGSvc(dkgNetwork.Round) {
		log.Info(ctx, "DKG service already running for this round; skipping finalization",
			"round", dkgNetwork.Round,
		)

		return
	}
	defer releaseDKGSvc(dkgNetwork.Round)

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

		return
	}

	if session.GetPhase() != types.PhaseDealing {
		log.Warn(ctx, "Session not in dealing phase, skipping finalize DKG", nil,
			"current_phase", session.GetPhase().String())
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
		"round", session.GetRound(),
	)
}

func (k *Keeper) callTEEFinalizeDKG(ctx context.Context, session *types.DKGSession) error {
	log.Info(ctx, "Finalize call to kernel client",
		"round", session.GetRound(),
	)

	if len(session.GetGlobalPubKey()) > 0 && len(session.GetSigFinalizeNetwork()) > 0 {
		log.Info(ctx, "DKG network already finalized in kernel client, skipping call Finalize request")

		return nil
	}

	var (
		resp *types.FinalizeDKGResponse
		err  error
	)
	codeCommitment := session.GetCodeCommitment()
	start := time.Now()
	retryErr := retry(ctx, func(ctx context.Context) error {
		req := &types.FinalizeDKGRequest{
			CodeCommitment: codeCommitment,
			Round:          session.GetRound(),
			IsResharing:    session.GetIsResharing(),
		}

		client, cErr := k.getClientWithReconnect(codeCommitment)
		if cErr != nil {
			return errors.Wrap(cErr, "no kernel client for session")
		}

		resp, err = client.FinalizeDKG(ctx, req)
		if err != nil {
			return err
		}

		return nil
	})
	observeKernelCall(labelOpFinalizeDKG, start, retryErr)

	if retryErr != nil {
		return errors.Wrap(retryErr, "kernel client Finalize request failed")
	}

	// Store all key-material fields atomically under a single lock so a concurrent reader
	// (e.g. HasKeyMaterial on the ABCI resume path) never observes a torn slice header.
	session.SetKeyMaterial(
		resp.GetParticipantsRoot(),
		resp.GetGlobalPubKey(),
		resp.GetSignature(),
		resp.GetPubKeyShare(),
		resp.GetPublicCoeffs(),
	)

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to update session after calling Finalize on the kernel client")
	}

	return nil
}

func (k *Keeper) callContractFinalizeDKG(ctx context.Context, session *types.DKGSession) error {
	// Read the key-material fields through mutex-guarded getters: SetKeyMaterial writes them
	// under Lock (here or on the active-round recovery path), so a raw read could race that write.
	globalPubKey := session.GetGlobalPubKey()
	sigFinalizeNetwork := session.GetSigFinalizeNetwork()

	log.Info(ctx, "Finalize contract call",
		"round", session.GetRound(),
		"global_pub_key", hex.EncodeToString(globalPubKey),
		"signature_len", len(sigFinalizeNetwork),
	)

	if _, err := k.contractClient.Finalize(
		ctx,
		session.GetRound(),
		session.GetEnclaveType(),
		session.GetParticipantsRoot(),
		globalPubKey,
		session.GetPublicCoeffs(),
		session.GetPubKeyShare(),
		sigFinalizeNetwork,
	); err != nil {
		return err
	}

	return nil
}
