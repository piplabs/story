package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/log"
	"slices"
	"strings"
)

// handleDKGInitialization handles the DKG initialization event.
func (k *Keeper) handleDKGInitialization(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG initialization",
		"mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave),
		"round", dkgNetwork.Round,
	)

	valAddr := k.validatorAddress.Hex()
	isParticipant := slices.Contains(dkgNetwork.ActiveValSet, strings.ToLower(valAddr))
	if !isParticipant {
		log.Debug(ctx, "Validator is not part of DKG committee", "validator", valAddr)

		return
	}

	session := types.NewDKGSession(dkgNetwork.Mrenclave, dkgNetwork.Round, dkgNetwork.ActiveValSet)
	if err := k.stateManager.CreateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to create DKG session", err)

		return
	}

	session.UpdatePhase(types.PhaseRegistering)

	var (
		resp *types.GenerateAndSealKeyResponse
		err  error
	)
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
		resp, err = k.teeClient.GenerateAndSealKey(ctx, req)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Error(ctx, "Failed to do DKG registration", err)

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

	log.Info(ctx, "InitializeDKG contract call",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"dkg_pub_key", hex.EncodeToString(resp.GetDkgPubKey()),
		"comm_pub_key", hex.EncodeToString(resp.GetCommPubKey()),
		"raw_quote_len", len(resp.GetRawQuote()),
	)

	_, err = k.contractClient.InitializeDKG(ctx, session.Round, session.Mrenclave, resp.GetDkgPubKey(), resp.GetCommPubKey(), resp.GetRawQuote())
	if err != nil {
		log.Error(ctx, "Failed to call InitializeDKG contract method", err)

		return
	}

	log.Info(ctx, "DKG Initialization complete",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
	)

	return
}
