package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/log"
	"slices"
)

// handleDKGResharing handles DKG resharing for new epochs.
func (k *Keeper) handleDKGResharing(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG resharing event",
		"mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave),
		"round", dkgNetwork.Round,
	)

	if dkgNetwork.Stage != types.DKGStageRegistration {
		log.Info(ctx, "DKG resharing is skipped because the current network stage is not in the registration stage")
	}

	isParticipant := slices.Contains(dkgNetwork.ActiveValSet, k.validatorEVMAddr)
	if !isParticipant {
		log.Info(ctx, "Validator is not part of current active validator set", "validator_evm_address", k.validatorEVMAddr)

		return
	}

	session := types.NewDKGSession(dkgNetwork.Mrenclave, dkgNetwork.Round, dkgNetwork.ActiveValSet)
	if err := k.stateManager.CreateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to create resharing DKG session", err)

		return
	}

	if err := k.callTEEGenerateAndSealKey(ctx, session); err != nil {
		log.Error(ctx, "Failed to generate the sealed key", err)

		return
	}

	if err := k.callContractInitializeDKG(ctx, session); err != nil {
		log.Error(ctx, "Failed to call initializeDKG method", err)

		return
	}

	session.UpdatePhase(types.PhaseInitialized)
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update session after calling initializeDKG method", err)

		return
	}

	log.Info(ctx, "DKG resharing initialization complete",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
	)

	return
}
