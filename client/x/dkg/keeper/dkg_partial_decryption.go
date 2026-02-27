package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/piplabs/story/lib/errors"
)

type partialDecryptionSubmission struct {
	Validator        string `json:"validator"`
	Round            uint32 `json:"round"`
	CodeCommitment   []byte `json:"code_commitment"`
	Pid              uint32 `json:"pid"`
	EncryptedPartial []byte `json:"encrypted_partial"`
	EphemeralPubKey  []byte `json:"ephemeral_pub_key"`
	PubShare         []byte `json:"pub_share"`
	Label            []byte `json:"label"`
}

func dkgPartialDecryptKey(codeCommitment [32]byte, round uint32, label []byte) string {
	labelHash := sha256.Sum256(label)
	return fmt.Sprintf("%s_%d_%s", hex.EncodeToString(codeCommitment[:]), round, hex.EncodeToString(labelHash[:]))
}

func (k *Keeper) setPartialDecryptionSubmission(
	ctx context.Context,
	validator common.Address,
	round uint32,
	codeCommitment [32]byte,
	pid uint32,
	encryptedPartial []byte,
	ephemeralPubKey []byte,
	pubShare []byte,
	label []byte,
) error {
	bz, err := json.Marshal(partialDecryptionSubmission{
		Validator:        validator.Hex(),
		Round:            round,
		CodeCommitment:   codeCommitment[:],
		Pid:              pid,
		EncryptedPartial: encryptedPartial,
		EphemeralPubKey:  ephemeralPubKey,
		PubShare:         pubShare,
		Label:            label,
	})
	if err != nil {
		return errors.Wrap(err, "marshal partial decryption submission")
	}

	if err := k.DKGPartialDecrypt.Set(ctx, dkgPartialDecryptKey(codeCommitment, round, label), bz); err != nil {
		return errors.Wrap(err, "set partial decryption submission")
	}

	return nil
}
