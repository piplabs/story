package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
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

	if !dkgSvcRunning.CompareAndSwap(false, true) {
		log.Info(ctx, "DKG service already running; skipping initialization")

		return
	}
	defer dkgSvcRunning.Store(false)

	if dkgNetwork.Stage != types.DKGStageRegistration {
		log.Info(ctx, "DKG initialization is skipped because the current network stage is not in the registration stage")

		return
	}

	valAddr := k.validatorAddress.Hex()
	isParticipant := slices.Contains(dkgNetwork.ActiveValSet, strings.ToLower(valAddr))
	if !isParticipant {
		log.Info(ctx, "Validator is not part of current active validator set", "validator", valAddr)

		return
	}

	session := types.NewDKGSession(dkgNetwork.Mrenclave, dkgNetwork.Round, dkgNetwork.ActiveValSet)
	if err := k.stateManager.CreateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to create DKG session", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	if session.Phase != types.PhaseInitializing {
		log.Warn(ctx, "Session not in initializing phase, skipping initialization", nil,
			"current_phase", session.Phase.String())
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	if err := k.callTEEGenerateAndSealKey(ctx, session); err != nil {
		log.Error(ctx, "Failed to generate the sealed key", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	if err := k.callContractInitializeDKG(ctx, session); err != nil {
		log.Error(ctx, "Failed to call initializeDKG method", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	session.UpdatePhase(types.PhaseInitialized)
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update session after calling initializeDKG method", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	log.Info(ctx, "DKG initialization complete",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
	)

	return
}

func (k *Keeper) callTEEGenerateAndSealKey(ctx context.Context, session *types.DKGSession) error {
	log.Info(ctx, "GenerateAndSealKey call to TEE client",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"validator", k.validatorAddress.Hex(),
	)

	if len(session.DKGPubKey) > 0 && len(session.CommPubKey) > 0 && len(session.RawQuote) > 0 {
		log.Info(ctx, "Already generated and sealed the key, skipping call GenerateAndSealKey request")

		return nil
	}

	var (
		resp *types.GenerateAndSealKeyResponse
		err  error
	)
	if err := retry(ctx, func(ctx context.Context) error {
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
		return errors.Wrap(err, "TEE client GenerateAndSealKey request failed")
	}

	session.DKGPubKey = resp.GetDkgPubKey()
	session.CommPubKey = resp.GetCommPubKey()
	session.RawQuote = resp.GetRawQuote()
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to update session after calling GenerateAndSealKey on the TEE client")
	}

	return nil
}

func (k *Keeper) callContractInitializeDKG(ctx context.Context, session *types.DKGSession) error {
	log.Info(ctx, "InitializeDKG contract call",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"dkg_pub_key", hex.EncodeToString(session.DKGPubKey),
		"comm_pub_key", hex.EncodeToString(session.CommPubKey),
		"raw_quote_len", len(session.RawQuote),
	)

	isInitialized, err := k.contractClient.IsInitialized(ctx, session.Round, session.Mrenclave, k.validatorAddress)
	if err != nil {
		return err
	}

	if isInitialized {
		log.Info(ctx, "Already initialized DKG on chain, skipping call initializeDKG method")

		return nil
	}

	if _, err := k.contractClient.InitializeDKG(
		ctx,
		session.Round,
		session.Mrenclave,
		session.DKGPubKey,
		session.CommPubKey,
		session.RawQuote,
	); err != nil {
		return err
	}

	return nil
}
