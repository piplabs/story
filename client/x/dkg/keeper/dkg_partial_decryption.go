package keeper

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

var ErrDuplicatePartialDecryptionSubmission = errors.New("partial decryption submission already exists")

func dkgPartialDecryptKey(requesterPubKey []byte, label []byte, round uint32, validator common.Address) string {
	requesterHash := crypto.Keccak256(requesterPubKey)
	return fmt.Sprintf(
		"%s_%s_%d_%s",
		hex.EncodeToString(requesterHash),
		hex.EncodeToString(label),
		round,
		validator.Hex(),
	)
}

func dkgPartialDecryptPrefix(requesterPubKey []byte, label []byte) string {
	requesterHash := crypto.Keccak256(requesterPubKey)
	return fmt.Sprintf("%s_%s_", hex.EncodeToString(requesterHash), hex.EncodeToString(label))
}

func (k *Keeper) setPartialDecryptionSubmission(
	ctx context.Context,
	validator common.Address,
	round uint32,
	pid uint32,
	encryptedPartial []byte,
	ephemeralPubKey []byte,
	pubShare []byte,
	requesterPubKey []byte,
	label []byte,
) error {
	key := dkgPartialDecryptKey(requesterPubKey, label, round, validator)
	exists, err := k.DKGPartialDecrypt.Has(ctx, key)
	if err != nil {
		return errors.Wrap(err, "check partial decryption submission")
	}
	if exists {
		return ErrDuplicatePartialDecryptionSubmission
	}

	bz, err := json.Marshal(types.DKGPartialDecryptionSubmission{
		Validator:        validator.Hex(),
		Round:            round,
		Pid:              pid,
		EncryptedPartial: encryptedPartial,
		EphemeralPubKey:  ephemeralPubKey,
		PubShare:         pubShare,
		Label:            label,
	})
	if err != nil {
		return errors.Wrap(err, "marshal partial decryption submission")
	}

	if err := k.DKGPartialDecrypt.Set(ctx, key, bz); err != nil {
		return errors.Wrap(err, "set partial decryption submission")
	}

	return nil
}

func decodePartialDecryptionSubmission(bz []byte) (*types.DKGPartialDecryptionSubmission, error) {
	var submission types.DKGPartialDecryptionSubmission
	if err := json.Unmarshal(bz, &submission); err != nil {
		return nil, errors.Wrap(err, "unmarshal partial decryption submission")
	}

	return &submission, nil
}
