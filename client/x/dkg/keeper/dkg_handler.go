package keeper

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// RegistrationInitialized handles DKG registration initialization event. These verified DKG registrations will be used
// by the DKG module & service to set the DKG network and perform further steps such as dealing.
func (k *Keeper) RegistrationInitialized(ctx context.Context, validator common.Address, codeCommitment [32]byte, round uint32, startBlockHeight uint64, startBlockHash [32]byte, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) error {
	latest, err := k.getLatestDKGNetwork(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get the latest dkg network")
	}

	if latest.Round != round {
		return errors.New(fmt.Sprintf("round mismatch: expected %d, got %d)", latest.Round, round))
	}

	if !bytes.Equal(latest.CodeCommitment, codeCommitment[:]) {
		return errors.New(fmt.Sprintf("codeCommitment mismatch: expected %s, got %s)", hex.EncodeToString(latest.CodeCommitment), hex.EncodeToString(codeCommitment[:])))
	}

	// Verify that startBlockHeight and startBlockHash match the latest DKG network's start block
	if latest.StartBlockHeight != int64(startBlockHeight) {
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

	index, err := k.getNextDKGRegistrationIndex(ctx, codeCommitment, round)
	if err != nil {
		return errors.Wrap(err, "failed to get next dkg registration index")
	}

	dkgReg := &types.DKGRegistration{
		Round:         round,
		ValidatorAddr: validator.Hex(),
		Index:         uint32(index),
		DkgPubKey:     dkgPubKey,
		CommPubKey:    commPubKey,
		RawQuote:      rawQuote,
		Status:        types.DKGRegStatusVerified,
	}

	if err := k.setDKGRegistration(ctx, codeCommitment, validator, dkgReg); err != nil {
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
		"raw_quote_len", len(rawQuote),
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

	if !bytes.Equal(latest.CodeCommitment, codeCommitment[:]) {
		return errors.New(fmt.Sprintf("codeCommitment mismatch: expected %s, got %s)", hex.EncodeToString(latest.CodeCommitment), hex.EncodeToString(codeCommitment[:])))
	}

	if err := k.validateParticipantsRoot(ctx, round, codeCommitment, participantsRoot); err != nil {
		return errors.Wrap(err, "failed to validate participants root")
	}

	if latest.Stage != types.DKGStageFinalization {
		return errors.New("round is not in network set stage")
	}

	// Retrieve the validator's DKG registration to get commPubKey
	reg, err := k.getDKGRegistration(ctx, codeCommitment, round, msgSender)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG registration for signature verification")
	}

	// Prevent double-finalization: reject if this validator already finalized for this round
	if reg.Status == types.DKGRegStatusFinalized {
		return errors.New("validator has already finalized for this round")
	}

	if err := verifyFinalizationSignature(reg.CommPubKey, round, codeCommitment, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare, signature); err != nil {
		return errors.Wrap(err, "finalization signature verification failed")
	}

	voteCount, err := k.AddGlobalPubKeyVote(ctx, codeCommitment, round, globalPubKey, publicCoeffs)
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

	if err := k.finalizeDKGRegistration(ctx, codeCommitment, round, msgSender, pubKeyShare); err != nil {
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

// validateParticipantsRoot validates the root hash of the participants
func (k *Keeper) validateParticipantsRoot(ctx context.Context, round uint32, codeCommitment, participantsRoot [32]byte) error {
	verifiedRegs, err := k.getDKGRegistrationsByStatus(ctx, codeCommitment, round, types.DKGRegStatusVerified)
	if err != nil {
		return errors.Wrap(err, "failed to get verified DKG registration")
	}

	if len(verifiedRegs) == 0 {
		return errors.New("no verified DKG registrations found")
	}

	addrs := make([]string, 0, len(verifiedRegs))
	for _, reg := range verifiedRegs {
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
func (*Keeper) UpgradeScheduled(ctx context.Context, activationHeight uint32, codeCommitment [32]byte) error {
	log.Info(ctx, "DKG UpgradeScheduled event received",
		"activation_height", activationHeight,
		"code_commitment", hex.EncodeToString(codeCommitment[:]),
	)
	// TODO: Implement actual upgrade scheduling logic
	return nil
}

// RemoteAttestationProcessedOnChain handles remote attestation processed event.
func (k *Keeper) RemoteAttestationProcessedOnChain(ctx context.Context, validator common.Address, chalStatus int, round uint32, codeCommitment [32]byte) error {
	index, err := k.getDKGRegistrationIndex(ctx, codeCommitment, round, strings.ToLower(validator.Hex()))
	if err != nil {
		return errors.Wrap(err, "failed to get dkg registration index")
	}

	log.Info(ctx, "DKG RemoteAttestationProcessedOnChain event received",
		"index", index,
		"validator", strings.ToLower(validator.Hex()),
		"challenge_status", chalStatus,
		"round", round,
		"code_commitment", hex.EncodeToString(codeCommitment[:]),
	)
	// TODO: Implement actual remote attestation processing logic
	return nil
}

// DealComplaintsSubmitted handles deal complaints submission event.
func (*Keeper) DealComplaintsSubmitted(ctx context.Context, index uint32, complainIndexes []uint32, round uint32, codeCommitment [32]byte) error {
	log.Info(ctx, "DKG DealComplaintsSubmitted event received",
		"index", index,
		"complain_indexes", complainIndexes,
		"round", round,
		"code_commitment", hex.EncodeToString(codeCommitment[:]),
	)
	// TODO: Implement actual deal complaints handling logic
	return nil
}

// DealVerified handles deal verification event.
func (*Keeper) DealVerified(ctx context.Context, index uint32, recipientIndex uint32, round uint32, codeCommitment [32]byte) error {
	log.Info(ctx, "DKG DealVerified event received",
		"index", index,
		"recipient_index", recipientIndex,
		"round", round,
		"code_commitment", hex.EncodeToString(codeCommitment[:]),
	)
	// TODO: Implement actual deal verification logic
	return nil
}

// InvalidDeal handles invalid deal event.
func (*Keeper) InvalidDeal(ctx context.Context, index uint32, round uint32, codeCommitment [32]byte) error {
	log.Info(ctx, "DKG InvalidDeal event received",
		"index", index,
		"round", round,
		"code_commitment", hex.EncodeToString(codeCommitment[:]),
	)
	// TODO: Implement actual invalid deal handling logic
	return nil
}

// verifyFinalizationSignature verifies the TEE's ECDSA signature over the DKG finalization data.
// It reproduces the message hash signed by the TEE, recovers the signer, and checks it matches
// the expected address derived from the validator's commPubKey.
func verifyFinalizationSignature(commPubKey []byte, round uint32, codeCommitment, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte, signature []byte) error {
	// commPubKey must be 64 bytes (uncompressed secp256k1 public key without 0x04 prefix)
	if len(commPubKey) != 64 {
		return errors.New("invalid commPubKey length", "expected", 64, "got", len(commPubKey))
	}

	// Compute total size of publicCoeffs for accurate capacity hint
	coeffsLen := 0
	for _, coeff := range publicCoeffs {
		if len(coeff) == 0 {
			return errors.New("empty public coefficient")
		}
		coeffsLen += len(coeff)
	}

	// Construct encoded message: codeCommitment(32B) + round(4B big-endian) + participantsRoot(32B) + globalPubKey + publicCoeffs...
	encoded := make([]byte, 0, 32+4+32+len(globalPubKey)+coeffsLen)
	encoded = append(encoded, codeCommitment[:]...)
	roundBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(roundBytes, round)
	encoded = append(encoded, roundBytes...)
	encoded = append(encoded, participantsRoot[:]...)
	encoded = append(encoded, globalPubKey...)
	for _, coeff := range publicCoeffs {
		encoded = append(encoded, coeff...)
	}
	encoded = append(encoded, pubKeyShare...)

	msgHash := crypto.Keccak256(encoded)

	// Compute Ethereum signed message hash: keccak256("\x19Ethereum Signed Message:\n32" + msgHash)
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	ethHash := crypto.Keccak256(append(prefix, msgHash...))

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

	recoveredPub, err := crypto.SigToPub(ethHash, sig)
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

// verifyPartialDecryptionSignature verifies the TEE's ECDSA signature over the partial decryption response data.
// It reproduces the message hash signed by signPartialDecryptResponse in the DKG server, recovers the signer,
// and checks it matches the expected address derived from the validator's commPubKey.
func verifyPartialDecryptionSignature(commPubKey []byte, codeCommitment [32]byte, round uint32, encryptedPartial, ephemeralPubKey, pubShare, signature []byte) error {
	if len(commPubKey) != 64 {
		return errors.New("invalid commPubKey length", "expected", 64, "got", len(commPubKey))
	}

	// Reconstruct the message exactly as in signPartialDecryptResponse:
	// encoded = codeCommitment || round(4B big-endian) || encryptedPartial || ephPubKey || pubShare
	roundBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(roundBytes, round)

	encoded := make([]byte, 0, len(codeCommitment)+4+len(encryptedPartial)+len(ephemeralPubKey)+len(pubShare))
	encoded = append(encoded, codeCommitment[:]...)
	encoded = append(encoded, roundBytes...)
	encoded = append(encoded, encryptedPartial...)
	encoded = append(encoded, ephemeralPubKey...)
	encoded = append(encoded, pubShare...)

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
func (k *Keeper) ThresholdDecryptRequested(ctx context.Context, requester common.Address, round uint32, codeCommitment [32]byte, requesterPubKey []byte, ciphertext []byte, label []byte) error {
	if !k.isDKGSvcEnabled {
		log.Info(ctx, "DKG service disabled; skipping threshold decrypt request")

		return nil
	}

	dkgNetwork, err := k.getDKGNetwork(ctx, codeCommitment, round)
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

	log.Info(ctx, "DKG ThresholdDecryptRequested event received",
		"requester", requester.Hex(),
		"round", round,
		"code_commitment", hex.EncodeToString(codeCommitment[:]),
		"requester_pubkey_len", len(requesterPubKey),
		"ciphertext_len", len(ciphertext),
		"label_len", len(label),
	)

	session, err := k.stateManager.GetSession(codeCommitment[:], round)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG session for decrypt request")
	}

	// Record the request so the off-chain service can pick it up and produce a TDH2 partial decrypt.
	session.AddDecryptRequest(types.DecryptRequest{
		Requester:       requester.Hex(),
		Round:           round,
		CodeCommitment:  codeCommitment[:],
		Ciphertext:      ciphertext,
		Label:           label,
		RequesterPubKey: requesterPubKey,
	})

	if err := k.stateManager.UpdateSession(ctx, session); err != nil {
		return errors.Wrap(err, "failed to persist decrypt request")
	}

	log.Info(ctx, "Queued threshold decrypt request",
		"session", session.GetSessionKey(),
		"pending_requests", len(session.DecryptRequests),
	)
	return nil
}

// PartialDecryptionSubmitted handles TDH2 partial decrypt submissions emitted by the contract.
// It stores submission payloads for later processing.
func (k *Keeper) PartialDecryptionSubmitted(
	ctx context.Context,
	validator common.Address,
	round uint32,
	codeCommitment [32]byte,
	pid uint32,
	encryptedPartial []byte,
	ephemeralPubKey []byte,
	pubShare []byte,
	label []byte,
	signature []byte,
) error {

	reg, err := k.getDKGRegistration(ctx, codeCommitment, round, validator)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG registration for signature verification")
	}

	if err := verifyPartialDecryptionSignature(reg.CommPubKey, codeCommitment, round, encryptedPartial, ephemeralPubKey, pubShare, signature); err != nil {
		return errors.Wrap(err, "partial decryption signature verification failed")
	}

	if !bytes.Equal(pubShare, reg.PubKeyShare) {
		return errors.New("pubShare mismatch: submitted pubShare does not match stored pubKeyShare",
			"validator", validator.Hex(),
			"round", round,
		)
	}

	if err := k.setPartialDecryptionSubmission(
		ctx,
		validator,
		round,
		codeCommitment,
		pid,
		encryptedPartial,
		ephemeralPubKey,
		pubShare,
		label,
	); err != nil {
		return errors.Wrap(err, "failed to store partial decryption submission")
	}

	log.Info(ctx, "DKG PartialDecryptionSubmitted event received",
		"validator", validator.Hex(),
		"round", round,
		"code_commitment", hex.EncodeToString(codeCommitment[:]),
		"pid", pid,
		"encrypted_partial_len", len(encryptedPartial),
		"ephemeral_pub_key_len", len(ephemeralPubKey),
		"pub_share_len", len(pubShare),
		"label_len", len(label),
	)

	return nil
}
