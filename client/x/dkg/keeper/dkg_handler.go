package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// RegistrationInitialized handles DKG registration initialization event. These verified DKG registrations will be used
// by the DKG module & service to set the DKG network and perform further steps such as dealing.
func (k *Keeper) RegistrationInitialized(ctx context.Context, validator common.Address, mrenclave [32]byte, round uint32, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) error {
	latest, err := k.getLatestDKGNetwork(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get the latest dkg network")
	}

	if latest.Round != round {
		return errors.New(fmt.Sprintf("round mismatch: expected %d, got %d)", latest.Round, round))
	}

	if !bytes.Equal(latest.Mrenclave, mrenclave[:]) {
		return errors.New(fmt.Sprintf("mrenclave mismatch: expected %s, got %s)", hex.EncodeToString(latest.Mrenclave), hex.EncodeToString(mrenclave[:])))
	}

	if latest.Stage != types.DKGStageRegistration {
		return errors.New("round is not in registration stage")
	}

	if !slices.Contains(latest.ActiveValSet, strings.ToLower(validator.Hex())) {
		return errors.New("msg sender is not in the active validator set")
	}

	index, err := k.getNextDKGRegistrationIndex(ctx, mrenclave, round)
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

	if err := k.setDKGRegistration(ctx, mrenclave, validator, dkgReg); err != nil {
		log.Error(ctx, "Failed to store DKG registration", err,
			"mrenclave", hex.EncodeToString(mrenclave[:]),
			"round", round,
			"validator_address", validator.Hex(),
			"next_index", index,
		)

		return errors.Wrap(err, "failed to store dkg registration")
	}

	log.Info(ctx, "DKG registration stored successfully",
		"mrenclave", hex.EncodeToString(mrenclave[:]),
		"round", round,
		"validator_address", validator.Hex(),
		"index", index,
		"status", types.DKGRegStatus_name[int32(types.DKGRegStatusVerified)],
		"dkg_pubkey", hex.EncodeToString(dkgPubKey),
		"comm_pubkey", hex.EncodeToString(commPubKey),
		"raw_quote_len", len(rawQuote),
	)

	return nil
}

// Finalized handles DKG finalization event.
func (k *Keeper) Finalized(ctx context.Context, round uint32, msgSender common.Address, mrenclave, participantsRoot [32]byte, signature, globalPubKey []byte, publicCoeffs [][]byte) error {
	latest, err := k.getLatestDKGNetwork(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get the latest dkg network")
	}

	if latest.Round != round {
		return errors.New(fmt.Sprintf("round mismatch: expected %d, got %d)", latest.Round, round))
	}

	if !bytes.Equal(latest.Mrenclave, mrenclave[:]) {
		return errors.New(fmt.Sprintf("mrenclave mismatch: expected %s, got %s)", hex.EncodeToString(latest.Mrenclave), hex.EncodeToString(mrenclave[:])))
	}

	if err := k.validateParticipantsRoot(ctx, round, mrenclave, participantsRoot); err != nil {
		return errors.Wrap(err, "failed to validate participants root")
	}

	if latest.Stage != types.DKGStageFinalization {
		return errors.New("round is not in network set stage")
	}

	voteCount, err := k.AddGlobalPubKeyVote(ctx, mrenclave, round, globalPubKey, publicCoeffs)
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

	if err := k.updateDKGRegistrationStatus(ctx, mrenclave, round, msgSender, types.DKGRegStatusFinalized); err != nil {
		return errors.Wrap(err, "failed to update dkg registration status")
	}

	log.Info(ctx, "DKG successfully finalized",
		"mrenclave", hex.EncodeToString(mrenclave[:]),
		"round", round,
		"validator_address", msgSender.Hex(),
		"status", types.DKGRegStatus_name[int32(types.DKGRegStatusFinalized)],
		"signature_len", len(signature),
	)

	return nil
}

// validateParticipantsRoot validates the root hash of the participants
func (k *Keeper) validateParticipantsRoot(ctx context.Context, round uint32, mrenclave, participantsRoot [32]byte) error {
	verifiedRegs, err := k.getDKGRegistrationsByStatus(ctx, mrenclave, round, types.DKGRegStatusVerified)
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
func (*Keeper) UpgradeScheduled(ctx context.Context, activationHeight uint32, mrenclave [32]byte) error {
	log.Info(ctx, "DKG UpgradeScheduled event received",
		"activation_height", activationHeight,
		"mrenclave", hex.EncodeToString(mrenclave[:]),
	)
	// TODO: Implement actual upgrade scheduling logic
	return nil
}

// RemoteAttestationProcessedOnChain handles remote attestation processed event.
func (k *Keeper) RemoteAttestationProcessedOnChain(ctx context.Context, validator common.Address, chalStatus int, round uint32, mrenclave [32]byte) error {
	index, err := k.getDKGRegistrationIndex(ctx, mrenclave, round, strings.ToLower(validator.Hex()))
	if err != nil {
		return errors.Wrap(err, "failed to get dkg registration index")
	}

	log.Info(ctx, "DKG RemoteAttestationProcessedOnChain event received",
		"index", index,
		"validator", strings.ToLower(validator.Hex()),
		"challenge_status", chalStatus,
		"round", round,
		"mrenclave", hex.EncodeToString(mrenclave[:]),
	)
	// TODO: Implement actual remote attestation processing logic
	return nil
}

// DealComplaintsSubmitted handles deal complaints submission event.
func (*Keeper) DealComplaintsSubmitted(ctx context.Context, index uint32, complainIndexes []uint32, round uint32, mrenclave [32]byte) error {
	log.Info(ctx, "DKG DealComplaintsSubmitted event received",
		"index", index,
		"complain_indexes", complainIndexes,
		"round", round,
		"mrenclave", hex.EncodeToString(mrenclave[:]),
	)
	// TODO: Implement actual deal complaints handling logic
	return nil
}

// DealVerified handles deal verification event.
func (*Keeper) DealVerified(ctx context.Context, index uint32, recipientIndex uint32, round uint32, mrenclave [32]byte) error {
	log.Info(ctx, "DKG DealVerified event received",
		"index", index,
		"recipient_index", recipientIndex,
		"round", round,
		"mrenclave", hex.EncodeToString(mrenclave[:]),
	)
	// TODO: Implement actual deal verification logic
	return nil
}

// InvalidDeal handles invalid deal event.
func (*Keeper) InvalidDeal(ctx context.Context, index uint32, round uint32, mrenclave [32]byte) error {
	log.Info(ctx, "DKG InvalidDeal event received",
		"index", index,
		"round", round,
		"mrenclave", hex.EncodeToString(mrenclave[:]),
	)
	// TODO: Implement actual invalid deal handling logic
	return nil
}

// ThresholdDecryptRequested handles TDH2 threshold decryption requests emitted by the contract.
// This is where validators should fetch ciphertext/label and produce partial decryptions (via TEE/TDH2).
func (k *Keeper) ThresholdDecryptRequested(ctx context.Context, requester common.Address, round uint32, mrenclave [32]byte, requesterPubKey []byte, ciphertext []byte, label []byte) error {
	if !k.isDKGSvcEnabled {
		log.Info(ctx, "DKG service disabled; skipping threshold decrypt request")

		return nil
	}

	dkgNetwork, err := k.getDKGNetwork(ctx, mrenclave, round)
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
		"mrenclave", hex.EncodeToString(mrenclave[:]),
		"requester_pubkey_len", len(requesterPubKey),
		"ciphertext_len", len(ciphertext),
		"label_len", len(label),
	)

	session, err := k.stateManager.GetSession(mrenclave[:], round)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG session for decrypt request")
	}

	// Record the request so the off-chain service can pick it up and produce a TDH2 partial decrypt.
	session.AddDecryptRequest(types.DecryptRequest{
		Requester:       requester.Hex(),
		Round:           round,
		Mrenclave:       mrenclave[:],
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
