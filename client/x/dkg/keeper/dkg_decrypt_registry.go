package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/collections"
	"github.com/piplabs/story/lib/errors"
)

// decryptRequestRegistryKey builds the key for the DecryptRequestRegistry map:
// hex(codeCommitment)_round_hex(sha256(label))
func decryptRequestRegistryKey(codeCommitment [32]byte, round uint32, label []byte) string {
	labelHash := sha256.Sum256(label)
	return fmt.Sprintf("%s_%d_%s", hex.EncodeToString(codeCommitment[:]), round, hex.EncodeToString(labelHash[:]))
}

// setDecryptRequestHeight records the block height at which a threshold decrypt
// request was registered. Called by all consensus nodes in ThresholdDecryptRequested.
func (k *Keeper) setDecryptRequestHeight(ctx context.Context, codeCommitment [32]byte, round uint32, label []byte, blockHeight uint64) error {
	key := decryptRequestRegistryKey(codeCommitment, round, label)
	if err := k.DecryptRequestRegistry.Set(ctx, key, blockHeight); err != nil {
		return errors.Wrap(err, "set decrypt request registry")
	}
	return nil
}

// getDecryptRequestHeight retrieves the block height stored for a decrypt request.
// Returns (height, true, nil) if found, or (0, false, nil) if not found.
func (k *Keeper) getDecryptRequestHeight(ctx context.Context, codeCommitment [32]byte, round uint32, label []byte) (uint64, bool, error) {
	key := decryptRequestRegistryKey(codeCommitment, round, label)
	height, err := k.DecryptRequestRegistry.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return 0, false, nil
		}
		return 0, false, errors.Wrap(err, "get decrypt request registry")
	}
	return height, true, nil
}

// deleteDecryptRequestHeight removes a single registry entry by label.
func (k *Keeper) deleteDecryptRequestHeight(ctx context.Context, codeCommitment [32]byte, round uint32, label []byte) error {
	key := decryptRequestRegistryKey(codeCommitment, round, label)
	if err := k.DecryptRequestRegistry.Remove(ctx, key); err != nil {
		return errors.Wrap(err, "delete decrypt request registry entry")
	}
	return nil
}

// sweepDecryptRequestRegistry deletes all registry entries for the given
// (codeCommitment, round) pair. Should be called when a DKG round ends.
func (k *Keeper) sweepDecryptRequestRegistry(ctx context.Context, codeCommitment [32]byte, round uint32) error {
	prefix := fmt.Sprintf("%s_%d_", hex.EncodeToString(codeCommitment[:]), round)

	iter, err := k.DecryptRequestRegistry.Iterate(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "iterate decrypt request registry")
	}
	defer iter.Close()

	var toDelete []string
	for ; iter.Valid(); iter.Next() {
		key, err := iter.Key()
		if err != nil {
			return errors.Wrap(err, "iterate decrypt request registry key")
		}
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			toDelete = append(toDelete, key)
		}
	}

	for _, key := range toDelete {
		if err := k.DecryptRequestRegistry.Remove(ctx, key); err != nil {
			return errors.Wrap(err, "remove decrypt request registry entry")
		}
	}

	return nil
}
