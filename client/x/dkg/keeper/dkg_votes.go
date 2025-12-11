package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"encoding/hex"
	"fmt"
	"github.com/piplabs/story/lib/errors"
)

// AddGlobalPubKeyVote Increase vote for global public key by 1.
func (k *Keeper) AddGlobalPubKeyVote(ctx context.Context, mrenclave [32]byte, round uint32, globalPubKey []byte) (uint32, error) {
	key := fmt.Sprintf("%s_%d_%s", hex.EncodeToString(mrenclave[:]), round, hex.EncodeToString(globalPubKey))

	// Read current votes (0 if not found)
	current, err := k.GlobalPubKeyVotes.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			current = 0
		} else {
			return 0, errors.Wrap(err, "failed to get votes for global public key", "mrenclave", hex.EncodeToString(mrenclave[:]), "round", round)
		}
	}

	// Increase and store
	newCount := current + 1
	if err := k.GlobalPubKeyVotes.Set(ctx, key, newCount); err != nil {
		return 0, errors.Wrap(err, "failed to set votes", "mrenclave", hex.EncodeToString(mrenclave[:]), "round", round)
	}

	return newCount, nil
}
