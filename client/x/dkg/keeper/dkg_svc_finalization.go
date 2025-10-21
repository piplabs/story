package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"

	"github.com/piplabs/story/lib/log"
)

// handleDKGFinalization handles the finalization phase event.
func (k *Keeper) handleDKGFinalization(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG finalization",
		"mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave),
		"round", dkgNetwork.Round,
	)

	session, err := k.stateManager.GetSession(dkgNetwork.Mrenclave, dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session", err)

		return
	}

	if session.Phase != types.PhaseDealing {
		log.Warn(ctx, "Session not in dealing phase, skipping finalization", nil,
			"current_phase", session.Phase.String(),
		)

		return
	}

	session.UpdatePhase(types.PhaseFinalizing)

	var resp *types.FinalizeDKGResponse
	if err := retry(ctx, func(ctx context.Context) error {
		log.Info(ctx, "FinalizeDKG call to TEE client",
			"mrenclave", session.GetMrenclaveString(),
			"round", session.Round,
		)

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
		log.Error(ctx, "Failed to finalize DKG", err)

		session.UpdatePhase(types.PhaseFailed)
		if updateErr := k.stateManager.UpdateSession(ctx, session); updateErr != nil {
			log.Error(ctx, "Failed to update session after TEE finalization error", updateErr)
		}

		return
	}

	session.GlobalPubKey = resp.GlobalPubKey
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "failed to update session", err)

		return
	}

	log.Info(ctx, "FinalizeDKG contract call",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"global_pub_key", hex.EncodeToString(session.GlobalPubKey),
		"signature_len", len(resp.GetSignature()),
	)

	_, err = k.contractClient.FinalizeDKG(
		ctx,
		session.Round,
		session.Mrenclave,
		session.GlobalPubKey,
		resp.GetSignature(),
	)
	if err != nil {
		log.Error(ctx, "failed to call FinalizeDKG contract method", err)

		return
	}

	log.Info(ctx, "DKG finalization phase complete",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
	)

	return
}
