package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/collections"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

// decryptRequestRegistryKey builds the key for the DecryptRequestRegistry map.
// The label is hex-encoded to prevent raw bytes from containing the '_' separator
// and causing key collisions between different (label, round) combinations.
func decryptRequestRegistryKey(requesterPubKey []byte, label []byte, round uint32, ciphertext []byte) string {
	if len(label) == 0 {
		return ""
	}

	requesterHash := sha256.Sum256(requesterPubKey)
	ciphertextHash := sha256.Sum256(ciphertext)
	return fmt.Sprintf(
		"%s_%s_%d_%s",
		hex.EncodeToString(requesterHash[:]),
		hex.EncodeToString(label),
		round,
		hex.EncodeToString(ciphertextHash[:]),
	)
}

// setDecryptRequest registers a threshold decrypt request in the registry.
// Called by all consensus nodes in ThresholdDecryptRequested.
func (k *Keeper) setDecryptRequest(ctx context.Context, requesterPubKey []byte, label []byte, req types.DecryptRequest) error {
	key := decryptRequestRegistryKey(requesterPubKey, label, req.Round, req.Ciphertext)
	if key == "" {
		return errors.New("cannot set decrypt request with empty label")
	}
	if err := k.DecryptRequestRegistry.Set(ctx, key, req); err != nil {
		return errors.Wrap(err, "set decrypt request registry")
	}
	return nil
}

// getDecryptRequest retrieves the stored decrypt request.
// Returns (req, true, nil) if found, or (zero, false, nil) if not found.
func (k *Keeper) getDecryptRequest(ctx context.Context, requesterPubKey []byte, label []byte, round uint32, ciphertext []byte) (types.DecryptRequest, bool, error) {
	key := decryptRequestRegistryKey(requesterPubKey, label, round, ciphertext)
	req, err := k.DecryptRequestRegistry.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.DecryptRequest{}, false, nil
		}
		return types.DecryptRequest{}, false, errors.Wrap(err, "get decrypt request registry")
	}
	return req, true, nil
}

// deleteDecryptRequest removes a single registry entry by label.
func (k *Keeper) deleteDecryptRequest(ctx context.Context, requesterPubKey []byte, label []byte, round uint32, ciphertext []byte) error {
	key := decryptRequestRegistryKey(requesterPubKey, label, round, ciphertext)
	if err := k.DecryptRequestRegistry.Remove(ctx, key); err != nil {
		return errors.Wrap(err, "delete decrypt request registry entry")
	}
	return nil
}

// pruneTimedOutDecryptRequests iterates all registry entries and removes any whose
// stored block height is older than DecryptTimeout (from params) relative to currentHeight.
// Called from BeginBlocker when the background cleanup worker signals.
func (k *Keeper) pruneTimedOutDecryptRequests(ctx context.Context, currentHeight uint64) error {
	params, err := k.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "get DKG params for prune timeout")
	}

	iter, err := k.DecryptRequestRegistry.Iterate(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "iterate decrypt request registry for pruning")
	}
	defer iter.Close()

	type entry struct {
		key string
	}
	var expired []entry
	for ; iter.Valid(); iter.Next() {
		key, err := iter.Key()
		if err != nil {
			return errors.Wrap(err, "iterate decrypt request registry key")
		}
		req, err := iter.Value()
		if err != nil {
			return errors.Wrap(err, "iterate decrypt request registry value")
		}
		if currentHeight > req.Height && currentHeight-req.Height > params.DecryptTimeout {
			expired = append(expired, entry{key})
		}
	}

	for _, e := range expired {
		if err := k.DecryptRequestRegistry.Remove(ctx, e.key); err != nil {
			return errors.Wrap(err, "remove expired decrypt request registry entry")
		}
	}

	return nil
}
