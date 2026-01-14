package keeper

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/cast"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
	"slices"
)

// handleDKGNetworkSet handles the network set event (after registration period).
func (k *Keeper) handleDKGNetworkSet(ctx context.Context, dkgNetwork *types.DKGNetwork) {
	log.Info(ctx, "Handling DKG network set",
		"mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave),
		"round", dkgNetwork.Round,
		"total", dkgNetwork.Total,
		"threshold", dkgNetwork.Threshold,
	)

	if !dkgSvcRunning.CompareAndSwap(false, true) {
		log.Info(ctx, "DKG service already running; skipping network setting")

		return
	}
	defer dkgSvcRunning.Store(false)

	if dkgNetwork.Stage != types.DKGStageNetworkSet {
		log.Info(ctx, "DKG NetworkSet is skipped because the current network stage is not in network set stage")

		return
	}

	isInCurRoundSet := slices.Contains(dkgNetwork.ActiveValSet, k.validatorEVMAddr)
	if !isInCurRoundSet {
		log.Info(ctx, "Skip setting DKG network as the validator is not in current round set")

		return
	}

	session, err := k.stateManager.GetSession(dkgNetwork.Mrenclave, dkgNetwork.Round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	if session.Phase != types.PhaseInitialized {
		log.Warn(ctx, "Session not in initialized phase, skipping network set", nil,
			"current_phase", session.Phase.String())
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	session.Total = dkgNetwork.Total
	session.Threshold = dkgNetwork.Threshold
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update session total and threshold", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	if err := k.callTEESetupDKGNetwork(ctx, session); err != nil {
		log.Error(ctx, "Failed to setup DKG network", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	if err := k.callContractSetNetwork(ctx, session); err != nil {
		log.Error(ctx, "Failed to call setNetwork method", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	session.UpdatePhase(types.PhaseDealing)
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to update session after calling setNetwork method", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	log.Info(ctx, "DKG Network set complete",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"total", session.Total,
		"threshold", session.Threshold,
	)

	return
}

func (k *Keeper) callTEESetupDKGNetwork(ctx context.Context, session *types.DKGSession) error {
	log.Info(ctx, "SetupDKGNetwork call to TEE client",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"total", session.Total,
		"threshold", session.Threshold,
	)

	if len(session.SigSetupNetwork) > 0 {
		log.Info(ctx, "DKG network already set in TEE client, skipping call SetupDKGNetwork request")

		return nil
	}

	mrenclave, err := cast.ToBytes32(session.Mrenclave)
	if err != nil {
		log.Error(ctx, "Failed to convert mrenclave to bytes32", err)

		return nil
	}

	var resp *types.SetupDKGNetworkResponse
	if err := retry(ctx, func(ctx context.Context) error {
		regs, err := k.getDKGRegistrationsByStatus(ctx, mrenclave, session.Round, types.DKGRegStatusVerified)
		if err != nil {
			return errors.Wrap(err, "failed to get verified dkg registrations from x/dkg module")
		}

		req := &types.SetupDKGNetworkRequest{
			Mrenclave:     session.Mrenclave,
			Round:         session.Round,
			Address:       k.validatorEVMAddr,
			Total:         session.Total,
			Threshold:     session.Threshold,
			Registrations: regs,
		}

		resp, err = k.teeClient.SetupDKGNetwork(ctx, req)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "TEE client SetupDKGNetwork request failed")
	}

	session.SigSetupNetwork = resp.GetSignature()
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to update session after calling SetupDKGNetwork on the TEE client")
	}

	return nil
}

func (k *Keeper) callContractSetNetwork(ctx context.Context, session *types.DKGSession) error {
	log.Info(ctx, "SetNetwork contract call",
		"mrenclave", session.GetMrenclaveString(),
		"round", session.Round,
		"total", session.Total,
		"threshold", session.Threshold,
		"signature_len", len(session.SigSetupNetwork),
	)

	validatorAddr := common.HexToAddress(k.validatorEVMAddr)
	isNetworkSet, err := k.contractClient.IsNetworkSet(ctx, session.Round, session.Mrenclave, validatorAddr)
	if err != nil {
		return err
	}

	if isNetworkSet {
		log.Info(ctx, "Already DKG network set on chain, skipping call setNetwork method")

		return nil
	}

	if _, err := k.contractClient.SetNetwork(
		ctx,
		session.Round,
		session.Total,
		session.Threshold,
		session.Mrenclave,
		session.SigSetupNetwork,
	); err != nil {
		return err
	}

	return nil
}
