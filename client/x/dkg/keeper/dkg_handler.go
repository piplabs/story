package keeper

import (
	"context"
	"encoding/hex"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// RegistrationInitialized handles DKG registration initialization event. These verified DKG registrations will be used
// by the DKG module & service to set the DKG network and perform further steps such as dealing.
func (k *Keeper) RegistrationInitialized(ctx context.Context, msgSender common.Address, mrenclave [32]byte, round uint32, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) error {
	dkgNetwork, err := k.getDKGNetwork(ctx, mrenclave, round)
	if err != nil {
		return errors.Wrap(err, "failed to get dkg network")
	}

	if dkgNetwork.Stage != types.DKGStageRegistration {
		return errors.New("round is not in registration stage")
	}

	if !slices.Contains(dkgNetwork.ActiveValSet, strings.ToLower(msgSender.Hex())) {
		return errors.New("msg sender is not in the active validator set")
	}

	index, err := k.getNextDKGRegistrationIndex(ctx, mrenclave, round)
	if err != nil {
		return errors.Wrap(err, "failed to get next dkg registration index")
	}

	dkgReg := &types.DKGRegistration{
		Mrenclave:  mrenclave[:],
		Round:      round,
		MsgSender:  msgSender.Hex(),
		Index:      index,
		DkgPubKey:  dkgPubKey,
		CommPubKey: commPubKey,
		RawQuote:   rawQuote,
		Status:     types.DKGRegStatusVerified,
	}

	if err := k.setDKGRegistration(ctx, mrenclave, msgSender, dkgReg); err != nil {
		log.Error(ctx, "Failed to store DKG registration", err,
			"mrenclave", hex.EncodeToString(mrenclave[:]),
			"round", round,
			"validator_address", msgSender.Hex(),
			"next_index", index,
		)

		return errors.Wrap(err, "failed to store dkg registration")
	}

	log.Info(ctx, "DKG registration stored successfully",
		"mrenclave", hex.EncodeToString(mrenclave[:]),
		"round", round,
		"status", "verified",
		"validator_address", msgSender.Hex(),
		"next_index", index,
		"dkg_pubkey_len", len(dkgPubKey),
		"comm_pubkey_len", len(commPubKey),
		"raw_quote_len", len(rawQuote),
	)

	return nil
}

// NetworkSet handles DKG network set event.
func (k *Keeper) NetworkSet(ctx context.Context, msgSender common.Address, mrenclave [32]byte, round uint32, total uint32, threshold uint32, signature []byte) error {
	log.Info(ctx, "DKG NetworkSet event received",
		"round", round,
		"msg_sender", msgSender.Hex(),
		"mrenclave", hex.EncodeToString(mrenclave[:]),
		"total", total,
		"threshold", threshold,
		"signature_len", len(signature),
	)

	err := k.updateDKGRegistrationStatus(ctx, mrenclave, round, msgSender, types.DKGRegStatusNetworkSet)
	if err != nil {
		return errors.Wrap(err, "failed to update dkg registration status")
	}

	// TODO: Implement remaining network set logic

	return nil
}

// Finalized handles DKG finalization event.
func (k *Keeper) Finalized(ctx context.Context, round uint32, msgSender common.Address, mrenclave [32]byte, signature []byte) error {
	index, err := k.getDKGRegistrationIndex(ctx, mrenclave, round, msgSender)
	if err != nil {
		return errors.Wrap(err, "failed to get dkg registration index")
	}

	log.Info(ctx, "DKG Finalized event received",
		"round", round,
		"msg_sender", msgSender.Hex(),
		"mrenclave", hex.EncodeToString(mrenclave[:]),
		"signature_len", len(signature),
		"index", index,
	)

	err = k.updateDKGRegistrationStatus(ctx, mrenclave, round, msgSender, types.DKGRegStatusFinalized)
	if err != nil {
		return errors.Wrap(err, "failed to update dkg registration status")
	}

	// TODO: Implement remaining finalization logic

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
	index, err := k.getDKGRegistrationIndex(ctx, mrenclave, round, validator)
	if err != nil {
		return errors.Wrap(err, "failed to get dkg registration index")
	}

	log.Info(ctx, "DKG RemoteAttestationProcessedOnChain event received",
		"index", index,
		"validator", validator.Hex(),
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
