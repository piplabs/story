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
		"old_code_commitment", hex.EncodeToString(oldCC),
		"already_registered", alreadyRegistered,
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

	// Resolve kernel client by code commitment:
	// - Normal rounds: use previous active round's CC to find the same binary.
	// - Upgrade rounds: use old binary's CC to find the NEW binary (by exclusion).
	var targetCC []byte
	if isUpgrade {
		targetCC = session.OldCodeCommitment
	} else {
		targetCC = oldCC
	}

	client, clientCC, cErr := k.getRegistrationKernelClient(isUpgrade, targetCC)
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

// getRegistrationKernelClient returns the appropriate kernel client and its code commitment
// for key generation. The cc parameter is interpreted differently based on isUpgrade:
//   - Normal round: cc is the previous active round's code commitment — used to find the
//     same kernel binary. nil cc (first round) falls back to first connected client.
//   - Upgrade round: cc is the OLD binary's code commitment — used to find the NEW binary
//     by returning a connected client whose CC differs from cc.
func (k *Keeper) getRegistrationKernelClient(isUpgrade bool, cc []byte) (types.KernelServiceClient, []byte, error) {
	client, resolvedCC, err := k.resolveKernelClientForRegistration(isUpgrade, cc)
	if err == nil {
		return client, resolvedCC, nil
	}

	// First attempt failed — try reconnecting to any disconnected endpoints
	// and resolve again. This handles the case where story started before kernel.
	k.kernelRouter.TryReconnect()

	return k.resolveKernelClientForRegistration(isUpgrade, cc)
}

// resolveKernelClientForRegistration looks up the appropriate kernel client without reconnection.
// For normal rounds, it returns the client matching the given code commitment (or the first
// available client when cc is nil, i.e. the very first DKG round).
// For upgrade rounds, cc is the OLD binary's code commitment; the function finds a connected
// client with a DIFFERENT code commitment, which is the new binary to generate keys on.
func (k *Keeper) resolveKernelClientForRegistration(isUpgrade bool, cc []byte) (types.KernelServiceClient, []byte, error) {
	if !isUpgrade {
		// Use pre-computed CC from previous active round's registration.
		if len(cc) > 0 {
			client, err := k.kernelRouter.GetClient(cc)

			return client, cc, err
		}

		// First round: no previous active round exists. Use first connected client.
		allCCs := k.kernelRouter.GetAllCodeCommitments()
		if len(allCCs) == 0 {
			return nil, nil, errors.New("no kernel clients available for registration")
		}

		client, err := k.kernelRouter.GetClient(allCCs[0])

		return client, allCCs[0], err
	}

	// Upgrade: cc is the old binary's CC. Find a connected client with a DIFFERENT CC.
	// If cc is empty, this validator was not in the previous committee
	// (newly joined). In this case, register directly against the new binary
	// since there is no old CC to compare against.
	allCCs := k.kernelRouter.GetAllCodeCommitments()
	if len(cc) == 0 {
		// New validator — use the first available client (should be the new binary).
		if len(allCCs) == 0 {
			return nil, nil, errors.New("no kernel clients available for upgrade registration (new validator)")
		}

		client, err := k.kernelRouter.GetClient(allCCs[0])

		return client, allCCs[0], err
	}

	for _, connCC := range allCCs {
		if !bytes.Equal(connCC, cc) {
			client, err := k.kernelRouter.GetClient(connCC)

			return client, connCC, err
		}
	}

	return nil, nil, errors.New("no new kernel client found for upgrade; ensure the new binary is running",
		"old_code_commitment", hex.EncodeToString(cc),
		"connected_clients", len(allCCs),
	)
}

// getOldCodeCommitment retrieves the code commitment from this validator's registration
// in the previous active DKG round. Returns nil if no previous active round or registration exists.
func (k *Keeper) getOldCodeCommitment(ctx context.Context) ([]byte, error) {
	prevActive, err := k.getLatestActiveDKGNetwork(ctx)
	if err != nil {
		log.Error(ctx, "Failed to get latest active DKG round for old code commitment retrieval", err)
		return nil, err
	}

	if prevActive == nil {
		log.Info(ctx, "No previous active DKG round found; returning nil old code commitment")
		return nil, nil
	}

	log.Info(ctx, "Fetching old code commitment from previous active DKG round",
		"round", prevActive.Round,
	)

	prevReg, err := k.getDKGRegistration(ctx, prevActive.Round, common.HexToAddress(k.validatorEVMAddr))
	if err != nil {
		return nil, err
	}
	log.Info(ctx, "Previous registration found for old code commitment retrieval",
		"code_commitment", hex.EncodeToString(prevReg.CodeCommitment),
		"round", prevReg.Round,
		"status", prevReg.Status.String(),
	)

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
