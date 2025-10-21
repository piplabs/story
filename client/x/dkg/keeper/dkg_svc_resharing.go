package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"

	"github.com/piplabs/story/lib/log"
)

// handleDKGResharing handles DKG resharing for new epochs.
func (k *Keeper) handleDKGResharing(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG resharing event",
		"mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave),
		"round", dkgNetwork.Round,
	)

	valAddr := k.validatorAddress.Hex()
	isParticipant := false
	for _, addr := range dkgNetwork.ActiveValSet {
		if addr == valAddr {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		log.Debug(ctx, "Validator not part of resharing committee")

		return
	}

	session := types.NewDKGSession(dkgNetwork.Mrenclave, dkgNetwork.Round, dkgNetwork.ActiveValSet)
	if err := k.stateManager.CreateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to create resharing session", err)

		return
	}

	session.UpdatePhase(types.PhaseRegistering)

	if err := retry(ctx, func(ctx context.Context) error {
		log.Info(ctx, "GenerateAndSealKey call to TEE client",
			"mrenclave", session.GetMrenclaveString(),
			"round", session.Round,
			"validator", k.validatorAddress.Hex(),
		)

		req := &types.GenerateAndSealKeyRequest{
			Address:   k.validatorAddress.Hex(),
			Mrenclave: session.Mrenclave,
			Round:     session.Round,
		}

		if _, err := k.teeClient.GenerateAndSealKey(ctx, req); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Error(ctx, "Failed to do DKG registration for resharing", err)

		session.UpdatePhase(types.PhaseFailed)
		if updateErr := k.stateManager.UpdateSession(ctx, session); updateErr != nil {
			log.Error(ctx, "Failed to update session after TEE error", updateErr)
		}

		return
	}

	session.UpdatePhase(types.PhaseNetworkSetting)
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update session after registration", err)

		return
	}

	log.Info(ctx, "DKG resharing initiated successfully",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
	)

	return
}
