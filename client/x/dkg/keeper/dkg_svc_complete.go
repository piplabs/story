package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"

	"github.com/piplabs/story/lib/log"
)

// handleDKGComplete handles the DKG completion event.
func (k *Keeper) handleDKGComplete(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG completion",
		"mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave),
		"round", dkgNetwork.Round,
	)

	session, err := k.stateManager.GetSession(dkgNetwork.Mrenclave, dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session", err)

		return
	}

	session.UpdatePhase(types.PhaseCompleted)
	session.IsFinalized = true // ready to start DKG threshold encryption/decryption

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update completed session", err)

		return
	}

	log.Info(ctx, "DKG process completed successfully",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"validator", k.validatorAddress.Hex(),
	)

	return
}
