package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/collections"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

// decryptRequestRegistryKey builds the key for the DecryptRequestRegistry map:
// hex(codeCommitment)_round_hex(sha256(label))
func decryptRequestRegistryKey(round uint32, label []byte) string {
	return fmt.Sprintf("%d_%s", round, hex.EncodeToString(label))
}

// setDecryptRequestHeight records the block height at which a threshold decrypt
// request was registered. Called by all consensus nodes in ThresholdDecryptRequested.
func (k *Keeper) setDecryptRequestHeight(ctx context.Context, round uint32, label []byte, blockHeight uint64) error {
	key := decryptRequestRegistryKey(round, label)
	if err := k.DecryptRequestRegistry.Set(ctx, key, blockHeight); err != nil {
		return errors.Wrap(err, "set decrypt request registry")
	}
	return nil
}

// getDecryptRequestHeight retrieves the block height stored for a decrypt request.
// Returns (height, true, nil) if found, or (0, false, nil) if not found.
func (k *Keeper) getDecryptRequestHeight(ctx context.Context, round uint32, label []byte) (uint64, bool, error) {
	key := decryptRequestRegistryKey(round, label)
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
func (k *Keeper) deleteDecryptRequestHeight(ctx context.Context, round uint32, label []byte) error {
	key := decryptRequestRegistryKey(round, label)
	if err := k.DecryptRequestRegistry.Remove(ctx, key); err != nil {
		return errors.Wrap(err, "delete decrypt request registry entry")
	}
	return nil
}

// pruneTimedOutDecryptRequests iterates all registry entries and removes any whose
// stored block height is older than PartialDecryptionTimeoutBlocks relative to currentHeight.
// Called from BeginBlocker when the background cleanup worker signals.
func (k *Keeper) pruneTimedOutDecryptRequests(ctx context.Context, currentHeight uint64) error {
	iter, err := k.DecryptRequestRegistry.Iterate(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "iterate decrypt request registry for pruning")
	}
	defer iter.Close()

	type entry struct {
		key    string
		height uint64
	}
	var expired []entry
	for ; iter.Valid(); iter.Next() {
		key, err := iter.Key()
		if err != nil {
			return errors.Wrap(err, "iterate decrypt request registry key")
		}
		height, err := iter.Value()
		if err != nil {
			return errors.Wrap(err, "iterate decrypt request registry value")
		}
		if currentHeight > height && currentHeight-height > types.PartialDecryptionTimeoutBlocks {
			expired = append(expired, entry{key, height})
		}
	}

	for _, e := range expired {
		if err := k.DecryptRequestRegistry.Remove(ctx, e.key); err != nil {
			return errors.Wrap(err, "remove expired decrypt request registry entry")
		}
	}

	return nil
}
