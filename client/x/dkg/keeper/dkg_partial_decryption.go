package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/collections"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

var ErrDuplicatePartialDecryptionSubmission = errors.New("partial decryption submission already exists")

func dkgPartialDecryptKey(requesterPubKey []byte, label []byte, ciphertext []byte, round uint32, validator common.Address) string {
	requesterHash := sha256.Sum256(requesterPubKey)
	ciphertextHash := sha256.Sum256(ciphertext)
	return fmt.Sprintf(
		"%s_%s_%s_%d_%s",
		hex.EncodeToString(requesterHash[:]),
		hex.EncodeToString(label),
		hex.EncodeToString(ciphertextHash[:]),
		round,
		validator.Hex(),
	)
}

func dkgPartialDecryptPrefix(requesterPubKey []byte, label []byte) string {
	requesterHash := sha256.Sum256(requesterPubKey)
	return fmt.Sprintf("%s_%s_", hex.EncodeToString(requesterHash[:]), hex.EncodeToString(label))
}

// dkgPartialDecryptRoundIndexKey builds the secondary index key for round-based pruning.
// Format: "{round:010d}_{primaryKey}"
// Round is zero-padded to 10 digits (uint32 max = 4,294,967,295 = 10 digits) so that
// lexicographic prefix scans equal numeric range scans over round values.
func dkgPartialDecryptRoundIndexKey(round uint32, primaryKey string) string {
	return fmt.Sprintf("%010d_%s", round, primaryKey)
}

// dkgPartialDecryptRoundIndexUpperBound returns the exclusive upper bound for an
// EndExclusive range scan over all secondary index entries with round <= cutoffRound.
func dkgPartialDecryptRoundIndexUpperBound(cutoffRound uint32) string {
	return fmt.Sprintf("%010d_", cutoffRound+1)
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
	ciphertext []byte,
) error {
	key := dkgPartialDecryptKey(requesterPubKey, label, ciphertext, round, validator)
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
		Ciphertext:       ciphertext,
	})
	if err != nil {
		return errors.Wrap(err, "marshal partial decryption submission")
	}

	if err := k.DKGPartialDecrypt.Set(ctx, key, bz); err != nil {
		return errors.Wrap(err, "set partial decryption submission")
	}

	// Write secondary round index entry alongside the primary. Value is empty;
	// presence is sufficient for range-delete pruning.
	indexKey := dkgPartialDecryptRoundIndexKey(round, key)
	if err := k.DKGPartialDecryptRoundIndex.Set(ctx, indexKey, []byte{}); err != nil {
		return errors.Wrap(err, "set partial decrypt round index")
	}

	return nil
}

// pruneOldPartialDecryptions removes all DKGPartialDecrypt entries (primary +
// secondary index) for rounds <= cutoffRound. It uses the secondary round index
// to avoid a full table scan: only entries in the range [0, cutoffRound] are visited.
// Called from BeginBlocker after a DKG round successfully becomes Active.
func (k *Keeper) pruneOldPartialDecryptions(ctx context.Context, cutoffRound uint32) error {
	upperBound := dkgPartialDecryptRoundIndexUpperBound(cutoffRound)
	rng := (&collections.Range[string]{}).EndExclusive(upperBound)

	iter, err := k.DKGPartialDecryptRoundIndex.Iterate(ctx, rng)
	if err != nil {
		return errors.Wrap(err, "iterate partial decrypt round index for pruning")
	}

	indexKeys, err := iter.Keys()
	iter.Close()
	if err != nil {
		return errors.Wrap(err, "collect partial decrypt round index keys for pruning")
	}

	for _, indexKey := range indexKeys {
		// Safety guard: index key must be at least "XXXXXXXXXX_" (11 bytes) + 1 char primary key.
		if len(indexKey) <= 11 {
			log.Warn(ctx, "Skipping malformed partial decrypt round index key during pruning", nil, "key", indexKey)
			continue
		}
		primaryKey := indexKey[11:] // strip "{round:010d}_" prefix

		if err := k.DKGPartialDecrypt.Remove(ctx, primaryKey); err != nil {
			return errors.Wrap(err, "remove partial decryption primary entry during pruning")
		}
		if err := k.DKGPartialDecryptRoundIndex.Remove(ctx, indexKey); err != nil {
			return errors.Wrap(err, "remove partial decryption round index entry during pruning")
		}
	}

	return nil
}

// MigratePartialDecryptRoundIndex is called once by the upgrade handler to
// backfill the secondary round index for all existing DKGPartialDecrypt entries.
// It does NOT delete any entries — pruning is handled exclusively by BeginBlocker
// on the next FinalizeDKGRound, ensuring all nodes prune at the same height and
// consensus is preserved.
func (k *Keeper) MigratePartialDecryptRoundIndex(ctx context.Context) error {
	iter, err := k.DKGPartialDecrypt.Iterate(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "migrate partial decrypt round index: iterate primary entries")
	}

	var primaryKeys []string
	for ; iter.Valid(); iter.Next() {
		primaryKey, err := iter.Key()
		if err != nil {
			iter.Close()
			return errors.Wrap(err, "migrate partial decrypt round index: read key")
		}
		primaryKeys = append(primaryKeys, primaryKey)
	}
	iter.Close()

	var indexed int
	for _, primaryKey := range primaryKeys {
		// Primary key format: {reqHash}_{labelHex}_{ciphertextHash}_{round}_{validator}
		// All components are hex or decimal — no underscores within any component.
		parts := strings.SplitN(primaryKey, "_", 5)
		if len(parts) < 5 {
			log.Warn(ctx, "Skipping malformed partial decrypt primary key during migration", nil, "key", primaryKey)
			continue
		}
		roundVal, err := strconv.ParseUint(parts[3], 10, 32)
		if err != nil {
			log.Warn(ctx, "Skipping partial decrypt key with unparseable round during migration", nil, "key", primaryKey, "round_field", parts[3])
			continue
		}
		indexKey := dkgPartialDecryptRoundIndexKey(uint32(roundVal), primaryKey)
		if err := k.DKGPartialDecryptRoundIndex.Set(ctx, indexKey, []byte{}); err != nil {
			return errors.Wrap(err, "migrate partial decrypt round index: write secondary index entry")
		}
		indexed++
	}

	log.Info(ctx, "Migrated partial decrypt round index", "indexed", indexed)

	return nil
}

func decodePartialDecryptionSubmission(bz []byte) (*types.DKGPartialDecryptionSubmission, error) {
	var submission types.DKGPartialDecryptionSubmission
	if err := json.Unmarshal(bz, &submission); err != nil {
		return nil, errors.Wrap(err, "unmarshal partial decryption submission")
	}

	return &submission, nil
}
