package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGNetworkSet handles the network set event (after registration period).
//
// TODO: before setting the network, ensure that `session.Registrations` is updated to only DKG validators who are verified.
func (k *Keeper) handleDKGNetworkSet(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG network set",
		"mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave),
		"round", dkgNetwork.Round,
		"total", dkgNetwork.Total,
		"threshold", dkgNetwork.Threshold,
	)

	session, err := k.stateManager.GetSession(dkgNetwork.Mrenclave, dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session", err)

		return
	}

	if session.Phase != types.PhaseNetworkSetting {
		log.Warn(ctx, "Session not in DKG network setting phase, skipping network set", nil,
			"current_phase", session.Phase.String())

		return
	}

	session.Total = dkgNetwork.Total
	session.Threshold = dkgNetwork.Threshold
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update session total and threshold", err)

		return
	}

	var resp *types.SetupDKGNetworkResponse
	if err := retry(ctx, func(ctx context.Context) error {
		rpcResp, err := k.GetVerifiedDKGRegistrations(ctx, &types.QueryGetVerifiedDKGRegistrationsRequest{
			Round:     session.Round,
			Mrenclave: session.Mrenclave,
		})
		if err != nil {
			return errors.Wrap(err, "failed to get verified dkg registrations from x/dkg module")
		}

		log.Info(ctx, "SetupDKGNetwork call to TEE client",
			"mrenclave", session.GetMrenclaveString(),
			"round", session.Round,
			"total", session.Total,
			"threshold", session.Threshold,
			"num_registrations", len(rpcResp.Registrations),
		)

		req := &types.SetupDKGNetworkRequest{
			Mrenclave:     session.Mrenclave,
			Round:         session.Round,
			Total:         session.Total,
			Threshold:     session.Threshold,
			Registrations: rpcResp.Registrations,
		}
		resp, err = k.teeClient.SetupDKGNetwork(ctx, req)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Error(ctx, "Failed to set DKG network", err)

		session.UpdatePhase(types.PhaseFailed)
		if updateErr := k.stateManager.UpdateSession(ctx, session); updateErr != nil {
			log.Error(ctx, "Failed to update session after TEE error", updateErr)
		}

		return
	}

	session.UpdatePhase(types.PhaseDealing)
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update session after network set", err)

		return
	}

	log.Info(ctx, "SetNetwork contract call",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"total", session.Total,
		"threshold", session.Threshold,
		"signature_len", len(resp.GetSignature()),
	)

	_, err = k.contractClient.SetNetwork(
		ctx,
		session.Round,
		session.Total,
		session.Threshold,
		session.Mrenclave,
		resp.GetSignature(),
	)
	if err != nil {
		log.Error(ctx, "Failed to call SetNetwork contract method", err)

		return
	}

	log.Info(ctx, "DKG network set complete",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"total", session.Total,
		"threshold", session.Threshold,
	)

	return
}
