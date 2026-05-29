package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"slices"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// Registered handles DKG registration initialization event. These verified DKG registrations will be used
// by the DKG module & service to set the DKG network and perform further steps such as dealing.
func (k *Keeper) Registered(ctx context.Context, validator common.Address, codeCommitment [32]byte, round uint32, startBlockHeight *big.Int, startBlockHash, enclaveType [32]byte, dkgPubKey []byte, commPubKey []byte, enclaveReport []byte) error {
	latest, err := k.getLatestDKGNetwork(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get the latest dkg network")
	}

	if latest.Round != round {
		return errors.New(fmt.Sprintf("round mismatch: expected %d, got %d)", latest.Round, round))
	}

	// Verify that startBlockHeight and startBlockHash match the latest DKG network's start block
	if latest.StartBlockHeight != startBlockHeight.Int64() {
		return errors.New(fmt.Sprintf("start block height mismatch: expected %d, got %d", latest.StartBlockHeight, startBlockHeight))
	}

	if !bytes.Equal(latest.StartBlockHash, startBlockHash[:]) {
		return errors.New(fmt.Sprintf("start block hash mismatch: expected %s, got %s", hex.EncodeToString(latest.StartBlockHash), hex.EncodeToString(startBlockHash[:])))
	}

	if latest.Stage != types.DKGStageRegistration {
		return errors.New("round is not in registration stage")
	}

	if !slices.Contains(latest.ActiveValSet, strings.ToLower(validator.Hex())) {
		return errors.New("msg sender is not in the active validator set")
	}

	exists, err := k.hasDKGRegistration(ctx, round, validator)
	if err != nil {
		return errors.Wrap(err, "failed to check existing dkg registration")
	}
	if exists {
		return errors.New("validator already registered for this round", "round", round, "validator", validator.Hex())
	}

	index, err := k.getNextDKGRegistrationIndex(ctx, round)
	if err != nil {
		return errors.Wrap(err, "failed to get next dkg registration index")
	}

	dkgReg := &types.DKGRegistration{
		Round:          round,
		ValidatorAddr:  validator.Hex(),
		Index:          uint32(index),
		DkgPubKey:      dkgPubKey,
		CommPubKey:     commPubKey,
		EnclaveReport:  enclaveReport,
		Status:         types.DKGRegStatusVerified,
		CodeCommitment: codeCommitment[:],
		EnclaveType:    enclaveType[:],
	}

	if err := k.setDKGRegistration(ctx, validator, dkgReg); err != nil {
		log.Error(ctx, "Failed to store DKG registration", err,
			"code_commitment", hex.EncodeToString(codeCommitment[:]),
			"round", round,
			"validator_address", validator.Hex(),
			"next_index", index,
		)

		return errors.Wrap(err, "failed to store dkg registration")
	}

	log.Info(ctx, "DKG registration stored successfully",
		"code_commitment", hex.EncodeToString(codeCommitment[:]),
		"round", round,
		"validator_address", validator.Hex(),
		"index", index,
		"start_block_height", startBlockHeight,
		"start_block_hash", hex.EncodeToString(startBlockHash[:]),
		"status", types.DKGRegStatus_name[int32(types.DKGRegStatusVerified)],
		"dkg_pubkey", hex.EncodeToString(dkgPubKey),
		"comm_pubkey", hex.EncodeToString(commPubKey),
		"raw_quote_len", len(enclaveReport),
	)

	return nil
}

// Finalized handles DKG finalization event.
func (k *Keeper) Finalized(ctx context.Context, round uint32, msgSender common.Address, codeCommitment, participantsRoot [32]byte, signature, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte) error {
	latest, err := k.getLatestDKGNetwork(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get the latest dkg network")
	}

	if latest.Round != round {
		return errors.New(fmt.Sprintf("round mismatch: expected %d, got %d)", latest.Round, round))
	}

	if latest.Stage != types.DKGStageFinalization {
		return errors.New("round is not in network set stage")
	}

	// Retrieve the validator's DKG registration to get commPubKey
	reg, err := k.getDKGRegistration(ctx, round, msgSender)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG registration for signature verification")
	}

	// Prevent double-finalization: reject if this validator already finalized for this round
	if reg.Status == types.DKGRegStatusFinalized {
		return errors.New("validator has already finalized for this round")
	}

	// Reject finalization by invalidated dealers (deal complaint found invalid via VSS verification)
	if reg.Status == types.DKGRegStatusInvalidated {
		return errors.New("dealer has been invalidated and cannot finalize")
	}

	if err := k.validateParticipantsRoot(ctx, round, participantsRoot); err != nil {
		return errors.Wrap(err, "failed to validate participants root")
	}

	if err := verifyFinalizationSignature(reg.CommPubKey, round, codeCommitment, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare, signature); err != nil {
		return errors.Wrap(err, "finalization signature verification failed")
	}

	voteCount, err := k.AddGlobalPubKeyVote(ctx, round, globalPubKey, publicCoeffs)
	if err != nil {
		return errors.Wrap(err, "failed to add vote for global public key")
	}

	if voteCount >= latest.Threshold && len(latest.GlobalPublicKey) == 0 {
		latest.GlobalPublicKey = globalPubKey

		latest.PublicCoeffs = publicCoeffs
		if err := k.setDKGNetwork(ctx, latest); err != nil {
			return errors.Wrap(err, "failed to set dkg network")
		}
	}

	if err := k.finalizeDKGRegistration(ctx, round, msgSender, pubKeyShare); err != nil {
		return errors.Wrap(err, "failed to update dkg registration status")
	}

	log.Info(ctx, "DKG successfully finalized",
		"code_commitment", hex.EncodeToString(codeCommitment[:]),
		"round", round,
		"validator_address", msgSender.Hex(),
		"status", types.DKGRegStatus_name[int32(types.DKGRegStatusFinalized)],
		"signature_len", len(signature),
	)

	return nil
}

// validateParticipantsRoot validates the root hash of the participants.
// It considers both Verified and Finalized registrations because earlier
// finalization events transition registrations from Verified to Finalized,
// so by the time later validators finalize, some registrations are already Finalized.
func (k *Keeper) validateParticipantsRoot(ctx context.Context, round uint32, participantsRoot [32]byte) error {
	verifiedRegs, err := k.getDKGRegistrationsByStatus(ctx, round, types.DKGRegStatusVerified)
	if err != nil {
		return errors.Wrap(err, "failed to get verified DKG registration")
	}

	finalizedRegs, err := k.getDKGRegistrationsByStatus(ctx, round, types.DKGRegStatusFinalized)
	if err != nil {
		return errors.Wrap(err, "failed to get finalized DKG registration")
	}

	allRegs := append(verifiedRegs, finalizedRegs...)

	if len(allRegs) == 0 {
		return errors.New("no verified or finalized DKG registrations found")
	}

	addrs := make([]string, 0, len(allRegs))
	for _, reg := range allRegs {
		addr := strings.ToLower(strings.TrimSpace(reg.ValidatorAddr))
		if !common.IsHexAddress(addr) {
			return errors.New("invalid validator evm address in verified registrations", "validator_addr", reg.ValidatorAddr)
		}

		addrs = append(addrs, addr)
	}

	slices.Sort(addrs)

	buf := make([]byte, 0, common.AddressLength*len(addrs))
	for _, a := range addrs {
		evmAddr := common.HexToAddress(a)
		buf = append(buf, evmAddr.Bytes()...)
	}

	if expected := crypto.Keccak256Hash(buf); expected != common.BytesToHash(participantsRoot[:]) {
		return errors.New(
			"participants root mismatch",
			"round", round,
			"expected_root", expected.Hex(),
			"got_root", common.BytesToHash(participantsRoot[:]).Hex(),
			"num_participants", len(addrs),
		)
	}

	return nil
}

// UpgradeScheduled handles upgrade scheduled event.
// If a pending (non-activated) upgrade already exists, it rejects the new one.
// Cancel the existing upgrade first via UpgradeCancelled, then re-schedule.
func (k *Keeper) UpgradeScheduled(ctx context.Context, activationHeight int64, upgradeVersion string) error {
	log.Info(ctx, "DKG UpgradeScheduled event received",
		"activation_height", activationHeight,
		"upgrade_version", upgradeVersion,
	)

	if len(upgradeVersion) == 0 {
		return errors.New("upgrade version cannot be empty")
	}

	// Reject if a pending upgrade already exists; cancel it first via DKG.cancelUpgrade()
	existing, err := k.GetPendingUpgrade(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to check existing pending upgrade")
	}

	if existing != nil {
		return errors.New("pending upgrade already exists; cancel the existing upgrade first before scheduling a new one",
			"existing_version", existing.UpgradeVersion,
			"existing_activation_height", existing.ActivationHeight,
		)
	}

	upgradeInfo := &types.KernelUpgradeInfo{
		UpgradeVersion:   upgradeVersion,
		ActivationHeight: activationHeight,
	}

	if err := k.SetKernelUpgradeInfo(ctx, upgradeInfo); err != nil {
		return errors.Wrap(err, "failed to store kernel upgrade info")
	}

	return nil
}

// UpgradeCancelled handles upgrade canceled event by removing the specified kernel upgrade info.
// The upgradeVersion parameter ensures the operator explicitly confirms which upgrade to cancel.
func (k *Keeper) UpgradeCancelled(ctx context.Context, upgradeVersion string) error {
	log.Info(ctx, "DKG UpgradeCancelled event received", "upgrade_version", upgradeVersion)

	info, err := k.GetKernelUpgradeInfo(ctx, upgradeVersion)
	if err != nil {
		return errors.Wrap(err, "failed to get kernel upgrade info")
	}

	if info == nil {
		log.Warn(ctx, "No upgrade found for specified version", nil, "upgrade_version", upgradeVersion)
		return errors.New("no upgrade found for specified version", "upgrade_version", upgradeVersion)
	}

	if err := k.DeleteKernelUpgradeInfo(ctx, upgradeVersion); err != nil {
		return errors.Wrap(err, "failed to delete upgrade info")
	}

	log.Info(ctx, "TEE upgrade canceled",
		"upgrade_version", upgradeVersion,
		"activation_height", info.ActivationHeight,
	)

	return nil
}

// finalizationSignatureMaterial holds the fields committed to by the validator
// signature on a DKG finalization response. RLP encoding gives each field an
// unambiguous length prefix, preventing boundary-shift collisions that arise
// from raw concatenation of variable-length byte slices.
type finalizationSignatureMaterial struct {
	CodeCommitment   []byte
	Round            uint32
	ParticipantsRoot [32]byte
	GlobalPubKey     []byte
	PublicCoeffs     [][]byte
	PubKeyShare      []byte
}

// verifyFinalizationSignature verifies the TEE's ECDSA signature over the DKG finalization data.
// It reproduces the message hash signed by the TEE, recovers the signer, and checks it matches
// the expected address derived from the validator's commPubKey.
func verifyFinalizationSignature(commPubKey []byte, round uint32, codeCommitment, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte, signature []byte) error {
	// commPubKey must be 64 bytes (uncompressed secp256k1 public key without 0x04 prefix)
	if len(commPubKey) != 64 {
		return errors.New("invalid commPubKey length", "expected", 64, "got", len(commPubKey))
	}

	for _, coeff := range publicCoeffs {
		if len(coeff) == 0 {
			return errors.New("empty public coefficient")
		}
	}

	// RLP-encode the finalization signature material (must match kernel's hashFinalizeDKGResponse).
	material := finalizationSignatureMaterial{
		CodeCommitment:   codeCommitment[:],
		Round:            round,
		ParticipantsRoot: participantsRoot,
		GlobalPubKey:     globalPubKey,
		PublicCoeffs:     publicCoeffs,
		PubKeyShare:      pubKeyShare,
	}
	encoded, err := rlp.EncodeToBytes(material)
	if err != nil {
		return errors.Wrap(err, "failed to RLP encode finalization signature material")
	}

	msgHash := crypto.Keccak256(encoded)

	if len(signature) != 65 {
		return errors.New("invalid signature length", "expected", 65, "got", len(signature))
	}

	// Make a copy to avoid mutating the original signature
	sig := make([]byte, 65)
	copy(sig, signature)

	// Adjust V value: convert Ethereum V (27/28) to recovery ID (0/1)
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	recoveredPub, err := crypto.SigToPub(msgHash, sig)
	if err != nil {
		return errors.Wrap(err, "failed to recover public key from signature")
	}

	recoveredAddr := crypto.PubkeyToAddress(*recoveredPub)

	// Expected address: address(uint160(uint256(keccak256(commPubKey))))
	expectedAddr := common.BytesToAddress(crypto.Keccak256(commPubKey))

	if recoveredAddr != expectedAddr {
		return errors.New("finalization signature address mismatch",
			"recovered", recoveredAddr.Hex(),
			"expected", expectedAddr.Hex(),
		)
	}

	return nil
}

// partialDecryptSignatureMaterial holds the fields committed to by the validator
// signature on a partial decryption response. RLP encoding gives each field an
// unambiguous length prefix, preventing boundary-shift collisions that arise
// from raw concatenation of variable-length byte slices.
type partialDecryptSignatureMaterial struct {
	Round            uint32
	Ciphertext       []byte
	EncryptedPartial []byte
	EphemeralPubKey  []byte
	PubShare         []byte
}

// verifyPartialDecryptionSignature verifies the TEE's ECDSA signature over the partial decryption response data.
// It reproduces the message hash signed by signPartialDecryptResponse in the DKG server, recovers the signer,
// and checks it matches the expected address derived from the validator's commPubKey.
//
// Both finalization and partial decryption signatures use RLP-encoded keccak256 hashes
// without the Ethereum Signed Message prefix, so kernel and CL are consistent.
func verifyPartialDecryptionSignature(commPubKey []byte, round uint32, ciphertext []byte, encryptedPartial, ephemeralPubKey, pubShare, signature []byte) error {
	if len(commPubKey) != 64 {
		return errors.New("invalid commPubKey length", "expected", 64, "got", len(commPubKey))
	}

	// RLP-encode the partial decryption signature material (must match kernel's signPartialDecryptResponse).
	material := partialDecryptSignatureMaterial{
		Round:            round,
		Ciphertext:       ciphertext,
		EncryptedPartial: encryptedPartial,
		EphemeralPubKey:  ephemeralPubKey,
		PubShare:         pubShare,
	}
	encoded, err := rlp.EncodeToBytes(material)
	if err != nil {
		return errors.Wrap(err, "failed to RLP encode partial decryption signature material")
	}

	respHash := crypto.Keccak256(encoded)

	if len(signature) != 65 {
		return errors.New("invalid signature length", "expected", 65, "got", len(signature))
	}

	sig := make([]byte, 65)
	copy(sig, signature)
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	recoveredPub, err := crypto.SigToPub(respHash, sig)
	if err != nil {
		return errors.Wrap(err, "failed to recover public key from signature")
	}

	recoveredAddr := crypto.PubkeyToAddress(*recoveredPub)
	expectedAddr := common.BytesToAddress(crypto.Keccak256(commPubKey))

	if recoveredAddr != expectedAddr {
		return errors.New("partial decryption signature address mismatch",
			"recovered", recoveredAddr.Hex(),
			"expected", expectedAddr.Hex(),
		)
	}

	return nil
}

// ThresholdDecryptRequested handles TDH2 threshold decryption requests emitted by the contract.
// This is where validators should fetch ciphertext/label and produce partial decryptions (via TEE/TDH2).
func (k *Keeper) ThresholdDecryptRequested(ctx context.Context, round uint32, requesterPubKey []byte, ciphertext []byte, label []byte, blockHeight uint64) error {
	// Consensus-level: all nodes record the request so PartialDecryptionSubmitted can enforce the timeout consistently.
	if err := k.setDecryptRequest(ctx, requesterPubKey, label, types.DecryptRequest{
		Round:           round,
		Ciphertext:      ciphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
		Height:          blockHeight,
	}); err != nil {
		return errors.Wrap(err, "failed to register decrypt request")
	}

	if !k.isDKGSvcEnabled {
		log.Info(ctx, "DKG service disabled; skipping threshold decrypt request")

		return nil
	}

	// Use a gasless context for KV reads inside isDKGSvcEnabled so that
	// DKG-enabled and DKG-disabled nodes produce identical GasUsed.
	gaslessCtx := gaslessSDKContext(ctx)

	dkgNetwork, err := k.getDKGNetwork(gaslessCtx, round)
	if err != nil {
		return errors.Wrap(err, "failed to get dkg network for decrypt request")
	}

	if dkgNetwork.Stage != types.DKGStageActive {
		log.Info(ctx, "Skipping threshold decrypt request; DKG round is not active",
			"round", round,
			"stage", dkgNetwork.Stage.String(),
		)

		return nil
	}

	if !slices.Contains(dkgNetwork.ActiveValSet, k.validatorEVMAddr) {
		log.Info(ctx, "Validator not in active DKG committee; skipping threshold decrypt request",
			"validator_address", k.validatorEVMAddr,
			"round", round,
		)

		return nil
	}

	// Pre-compute values that require SDK context, then dispatch async.
	decryptReq := types.DecryptRequest{
		Round:           round,
		Ciphertext:      ciphertext,
		Label:           label[:],
		RequesterPubKey: requesterPubKey,
		Height:          blockHeight,
	}

	asyncCtx, cancel := dkgAsyncContext()

	go func() {
		defer cancel()

		k.handleThresholdDecryptRequest(asyncCtx, round, decryptReq)
	}()

	return nil
}

// handleThresholdDecryptRequest queues a threshold decrypt request to the DKG session
// for the decrypt worker to pick up. Runs asynchronously to avoid blocking FinalizeBlock.
func (k *Keeper) handleThresholdDecryptRequest(ctx context.Context, round uint32, req types.DecryptRequest) {
	log.Info(ctx, "Handling threshold decrypt request",
		"round", round,
		"label", hex.EncodeToString(req.Label),
		"requester_pubkey_len", len(req.RequesterPubKey),
		"ciphertext_len", len(req.Ciphertext),
		"block_height", req.Height,
	)

	session, err := k.stateManager.GetSession(round)
	if err != nil {
		log.Error(ctx, "Failed to get DKG session for decrypt request", err,
			"round", round,
		)

		return
	}

	session.AddDecryptRequest(types.PendingDecryptRequest{DecryptRequest: req, RetryCount: 0})

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		log.Error(ctx, "Failed to persist decrypt request to session", err,
			"round", round,
			"label", hex.EncodeToString(req.Label),
			"session", session.GetSessionKey(),
		)

		return
	}

	log.Info(ctx, "Queued threshold decrypt request",
		"session", session.GetSessionKey(),
		"round", round,
		"label", hex.EncodeToString(req.Label),
		"pending_requests", len(session.GetDecryptRequests()),
	)
}

// PartialDecryptionSubmitted handles TDH2 partial decrypt submissions emitted by the contract.
// It stores submission payloads for later processing and returns true only when validation succeeds.
func (k *Keeper) PartialDecryptionSubmitted(
	ctx context.Context,
	validator common.Address,
	round uint32,
	pid uint32,
	encryptedPartial []byte,
	ephemeralPubKey []byte,
	pubShare []byte,
	requesterPubKey []byte,
	ciphertext []byte,
	label []byte,
	signature []byte,
) (bool, error) {
	// Enforce timeout: reject partial decryptions submitted too late.
	req, found, err := k.getDecryptRequest(ctx, requesterPubKey, label, round, ciphertext)
	if err != nil {
		return false, errors.Wrap(err, "failed to look up decrypt request registry")
	}
	if !found {
		log.Info(ctx, "Partial decryption submitted for unknown or cleaned-up request",
			"validator", validator.Hex(),
			"round", round,
		)
		return false, nil
	}

	if round != req.Round {
		return false, errors.New("round mismatch between partial decryption submission and decrypt request",
			"validator", validator.Hex(),
			"submission_round", round,
			"request_round", req.Round,
		)
	}

	if !bytes.Equal(ciphertext, req.Ciphertext) {
		return false, errors.New("ciphertext mismatch between partial decryption submission and decrypt request",
			"validator", validator.Hex(),
			"label", hex.EncodeToString(label),
			"round", round,
			"submission_ciphertext", hex.EncodeToString(ciphertext),
			"request_ciphertext", hex.EncodeToString(req.Ciphertext),
		)
	}

	currentHeight := uint64(sdk.UnwrapSDKContext(ctx).BlockHeight())
	if currentHeight-req.Height > types.DefaultDecryptTimeout {
		log.Info(ctx, "Partial decryption submission timeout exceeded; cleaning up registry entry",
			"request_height", req.Height,
			"current_height", currentHeight,
			"timeout_blocks", types.DefaultDecryptTimeout,
			"validator", validator.Hex(),
		)
		if err := k.deleteDecryptRequest(ctx, requesterPubKey, label, round, ciphertext); err != nil {
			return false, errors.Wrap(err, "failed to delete expired decrypt request registry entry")
		}
		return false, nil
	}

	reg, err := k.getDKGRegistration(ctx, req.Round, validator)
	if err != nil {
		return false, errors.Wrap(err, "failed to get DKG registration for signature verification")
	}

	// Reject partials from validators invalidated for this round; their partials are
	// inconsistent with the committee.
	if network, err := k.getDKGNetwork(ctx, req.Round); err == nil &&
		k.isV190Round(ctx, network) && reg.Status == types.DKGRegStatusInvalidated {
		return false, errors.New("validator is invalidated and cannot submit partial decryptions",
			"validator", validator.Hex(),
			"round", req.Round,
		)
	}

	if !bytes.Equal(pubShare, reg.PubKeyShare) {
		return false, errors.New("pubShare mismatch: submitted pubShare does not match stored pubKeyShare",
			"validator", validator.Hex(),
			"round", req.Round,
		)
	}

	if err := verifyPartialDecryptionSignature(reg.CommPubKey, round, ciphertext, encryptedPartial, ephemeralPubKey, pubShare, signature); err != nil {
		return false, errors.Wrap(err, "partial decryption signature verification failed")
	}

	if err := k.setPartialDecryptionSubmission(
		ctx,
		validator,
		round,
		pid,
		encryptedPartial,
		ephemeralPubKey,
		pubShare,
		requesterPubKey,
		label,
		ciphertext,
	); err != nil {
		if errors.Is(err, ErrDuplicatePartialDecryptionSubmission) {
			log.Info(ctx, "Duplicate partial decryption submission received; ignoring",
				"validator", validator.Hex(),
				"round", round,
				"pid", pid,
			)
			return false, nil
		}
		return false, errors.Wrap(err, "failed to store partial decryption submission")
	}

	log.Info(ctx, "DKG PartialDecryptionSubmitted event received",
		"validator", validator.Hex(),
		"round", round,
		"pid", pid,
		"encrypted_partial_len", len(encryptedPartial),
		"ephemeral_pub_key_len", len(ephemeralPubKey),
		"pub_share_len", len(pubShare),
		"requester_pub_key_len", len(requesterPubKey),
		"label_len", len(label),
	)

	return true, nil
}
