package keeper

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
	"slices"
)

// handleDKGFinalization handles the finalization phase event.
func (k *Keeper) handleDKGFinalization(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG finalization",
		"mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave),
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

	session, err := k.stateManager.GetSession(dkgNetwork.Mrenclave, dkgNetwork.Round)
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
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
	)

	return
}

func (k *Keeper) callTEEFinalizeDKG(ctx context.Context, session *types.DKGSession) error {
	log.Info(ctx, "FinalizeDKG call to TEE client",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
	)

	if len(session.GlobalPubKey) > 0 && len(session.SigFinalizeNetwork) > 0 {
		log.Info(ctx, "DKG network already finalized in TEE client, skipping call FinalizeDKG request")

		return nil
	}

	var (
		resp *types.FinalizeDKGResponse
		err  error
	)
	if err := retry(ctx, func(ctx context.Context) error {
		req := &types.FinalizeDKGRequest{
			Mrenclave: session.Mrenclave,
			Round:     session.Round,
		}

		resp, err = k.teeClient.FinalizeDKG(ctx, req)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "TEE client FinalizeDKG request failed")
	}

	session.GlobalPubKey = resp.GetGlobalPubKey()
	session.SigFinalizeNetwork = resp.GetSignature()
	session.PublicCoeffs = resp.GetPublicCoeffs()
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to update session after calling FinalizeDKG on the TEE client")
	}

	return nil
}

func (k *Keeper) callContractFinalizeDKG(ctx context.Context, session *types.DKGSession) error {
	log.Info(ctx, "FinalizeDKG contract call",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"global_pub_key", hex.EncodeToString(session.GlobalPubKey),
		"signature_len", len(session.SigFinalizeNetwork),
	)

	validatorAddr := common.HexToAddress(k.validatorEVMAddr)
	isFinalized, err := k.contractClient.IsFinalized(ctx, session.Round, session.Mrenclave, validatorAddr)
	if err != nil {
		return err
	}

	if isFinalized {
		log.Info(ctx, "Already finalized DKG on chain, skipping call finalizeDKG method")

		return nil
	}

	if _, err := k.contractClient.FinalizeDKG(
		ctx,
		session.Round,
		session.Mrenclave,
		session.GlobalPubKey,
		session.PublicCoeffs,
		session.SigFinalizeNetwork,
	); err != nil {
		return err
	}

	return nil
}
