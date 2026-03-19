package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"slices"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// handleDKGRegistration handles the DKG registration.
// oldCC is the pre-computed old code commitment from the previous active DKG round,
// resolved before the goroutine (which requires SDK context for KV store access).
// alreadyRegistered indicates whether this validator already has an on-chain
// registration for this round (pre-computed with SDK context before the goroutine).
// When true, the contract call is skipped and the session advances directly.
func (k *Keeper) handleDKGRegistration(ctx context.Context, dkgNetwork *types.DKGNetwork, oldCC []byte, alreadyRegistered bool) {
	log.Info(ctx, "Handling DKG registration",
		"round", dkgNetwork.Round,
	)

	if !tryAcquireDKGSvc(dkgNetwork.Round) {
		log.Info(ctx, "DKG service already running for this round; skipping registration",
			"round", dkgNetwork.Round,
		)

		return
	}
	defer releaseDKGSvc(dkgNetwork.Round)

	if dkgNetwork.Stage != types.DKGStageRegistration {
		log.Info(ctx, "DKG registration is skipped because the current network stage is not in the registration stage")

		return
	}

	isInCurRoundSet := slices.Contains(dkgNetwork.ActiveValSet, k.validatorEVMAddr)

	// NOTE:
	// Even if this validator is an old member (i.e., not part of the current round set),
	// we still create a session. Old members do not generate new keys, but they still
	// participate in later stages (especially dealing, and finalization).
	// Therefore, a session must exist regardless of key generation eligibility.
	session := types.NewDKGSession(dkgNetwork.Round, dkgNetwork.ActiveValSet, dkgNetwork.IsResharing, k.enclaveType)
	session.IsUpgrade = dkgNetwork.IsUpgrade

	// For upgrade rounds, store the old binary's code commitment so that dealers
	// can route TEE calls to the correct (old) binary for deal generation.
	if dkgNetwork.IsUpgrade && len(oldCC) > 0 {
		session.OldCodeCommitment = oldCC

		// For old-only members (not in current round set), set CodeCommitment
		// to old CC so they have a valid routing target for processing deals/responses.
		if !isInCurRoundSet {
			session.CodeCommitment = oldCC
		}
	}

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

	// only current round members generate a new key
	if !isInCurRoundSet {
		log.Info(ctx, "Skip generating sealed keys for DKG as the validator is not in current round set")

		return
	}

	if err := k.callTEEGenerateAndSealKey(ctx, session, dkgNetwork.IsUpgrade, oldCC); err != nil {
		log.Error(ctx, "Failed to generate the sealed key", err)
		k.stateManager.MarkFailed(ctx, session)

		return
	}

	if alreadyRegistered {
		log.Info(ctx, "Validator already registered on-chain; skipping contract call",
			"round", session.Round,
		)
	} else if err := k.callContractRegister(ctx, session); err != nil {
		log.Error(ctx, "Failed to call register method", err)
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
		"round", session.Round,
	)
}

func (k *Keeper) callTEEGenerateAndSealKey(ctx context.Context, session *types.DKGSession, isUpgrade bool, oldCC []byte) error {
	log.Info(ctx, "GenerateAndSealKey call to kernel client",
		"round", session.Round,
		"validator", k.validatorEVMAddr,
		"is_upgrade", isUpgrade,
	)

	if len(session.DKGPubKey) > 0 && len(session.CommPubKey) > 0 && len(session.EnclaveReport) > 0 {
		log.Info(ctx, "Already generated and sealed the key, skipping call GenerateAndSealKey request")

		return nil
	}

	// For upgrade rounds, route to the NEW binary's kernel client (CC != oldCC).
	// For normal rounds, route to the previous active round's kernel client.
	client, clientCC, cErr := k.getRegistrationKernelClient(isUpgrade, oldCC, session.OldCodeCommitment)
	if cErr != nil {
		return errors.Wrap(cErr, "no kernel client available for registration")
	}

	// Set the code commitment on the session from the resolved kernel client
	// so the GenerateAndSealKey request includes it.
	if len(session.CodeCommitment) == 0 {
		session.CodeCommitment = clientCC
	}

	var (
		resp *types.GenerateAndSealKeyResponse
		err  error
	)
	if err := retry(ctx, func(ctx context.Context) error {
		req := &types.GenerateAndSealKeyRequest{
			Address:        k.validatorEVMAddr,
			CodeCommitment: session.CodeCommitment,
			Round:          session.Round,
		}

		resp, err = client.GenerateAndSealKey(ctx, req)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "kernel client GenerateAndSealKey request failed")
	}

	// Persist the code commitment from the TEE response for future routing
	// (used by dealing, finalization, and threshold decryption)
	session.CodeCommitment = resp.GetCodeCommitment()
	session.DKGPubKey = resp.GetDkgPubKey()
	session.CommPubKey = resp.GetCommPubKey()
	session.EnclaveReport = resp.GetEnclaveReport()
	session.StartBlockHeight = resp.GetStartBlockHeight()

	session.StartBlockHash = resp.GetStartBlockHash()
	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to update session after calling GenerateAndSealKey on the kernel client")
	}

	return nil
}

// getRegistrationKernelClient returns the appropriate kernel client and its code commitment for key generation.
//
// Non-upgrade: looks up the previous active round's code commitment from on-chain state
// and returns the corresponding kernel client. For the very first round (no previous active
// round exists), falls back to the first connected client.
//
// Upgrade: uses the provided upgradeOldCC to find the NEW binary — returns the connected client
// whose code commitment differs from upgradeOldCC.
//
// prevRoundCC is the code commitment from the previous active DKG round (pre-computed
// from SDK context before the async goroutine).
func (k *Keeper) getRegistrationKernelClient(isUpgrade bool, prevRoundCC []byte, upgradeOldCC []byte) (types.KernelServiceClient, []byte, error) {
	if !isUpgrade {
		// Use pre-computed CC from previous active round's registration.
		if len(prevRoundCC) > 0 {
			client, err := k.kernelRouter.GetClient(prevRoundCC)

			return client, prevRoundCC, err
		}

		// First round: no previous active round exists. Use first connected client.
		allCCs := k.kernelRouter.GetAllCodeCommitments()
		if len(allCCs) == 0 {
			return nil, nil, errors.New("no kernel clients available for registration")
		}

		client, err := k.kernelRouter.GetClient(allCCs[0])

		return client, allCCs[0], err
	}

	// Upgrade: upgradeOldCC must be provided so we can find the NEW binary.
	if len(upgradeOldCC) == 0 {
		return nil, nil, errors.New("old code commitment required for upgrade registration")
	}

	allCCs := k.kernelRouter.GetAllCodeCommitments()
	for _, cc := range allCCs {
		if !bytes.Equal(cc, upgradeOldCC) {
			client, err := k.kernelRouter.GetClient(cc)

			return client, cc, err
		}
	}

	return nil, nil, errors.New("no new kernel client found for upgrade; ensure the new binary is running",
		"old_code_commitment", hex.EncodeToString(upgradeOldCC),
		"connected_clients", len(allCCs),
	)
}

// getOldCodeCommitment retrieves the code commitment from this validator's registration
// in the previous active DKG round. Returns nil if no previous active round or registration exists.
func (k *Keeper) getOldCodeCommitment(ctx context.Context) ([]byte, error) {
	prevActive, err := k.getLatestActiveDKGNetwork(ctx)
	if err != nil {
		return nil, err
	}

	if prevActive == nil {
		return nil, nil
	}

	prevReg, err := k.getDKGRegistration(ctx, prevActive.Round, common.HexToAddress(k.validatorEVMAddr))
	if err != nil {
		return nil, err
	}

	return prevReg.CodeCommitment, nil
}

func (k *Keeper) callContractRegister(ctx context.Context, session *types.DKGSession) error {
	log.Info(ctx, "Register contract call",
		"round", session.Round,
		"start_block_height", session.StartBlockHeight,
		"start_block_hash", hex.EncodeToString(session.StartBlockHash),
		"dkg_pub_key", hex.EncodeToString(session.DKGPubKey),
		"comm_pub_key", hex.EncodeToString(session.CommPubKey),
		"raw_quote_len", len(session.EnclaveReport),
	)

	if _, err := k.contractClient.Register(
		ctx,
		session.Round,
		session.EnclaveType,
		uint64(session.StartBlockHeight),
		session.StartBlockHash,
		session.DKGPubKey,
		session.CommPubKey,
		session.EnclaveReport,
	); err != nil {
		return err
	}

	return nil
}
