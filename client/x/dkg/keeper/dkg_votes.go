package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/collections"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/lib/errors"
)

// globalPubKeyVoteKey identifies a (round, globalPubKey, publicCoeffs) vote bucket.
// Validators that agree on the same consensus polynomial produce the same key.
func globalPubKeyVoteKey(round uint32, globalPubKey []byte, publicCoeffs [][]byte) string {
	return fmt.Sprintf(
		"%d_%s_%s",
		round,
		hex.EncodeToString(globalPubKey),
		hex.EncodeToString(hashPublicCoeffs(publicCoeffs)),
	)
}

// finalizeVoteStoreKey is the FinalizeVotes state key for a validator in a round.
func finalizeVoteStoreKey(round uint32, validator common.Address) string {
	return fmt.Sprintf("%d_%s", round, validator.Hex())
}

// AddGlobalPubKeyVote Increase vote for global public key by 1.
func (k *Keeper) AddGlobalPubKeyVote(ctx context.Context, round uint32, globalPubKey []byte, publicCoeffs [][]byte) (uint32, error) {
	key := globalPubKeyVoteKey(round, globalPubKey, publicCoeffs)

	// Read current votes (0 if not found)
	current, err := k.GlobalPubKeyVotes.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			current = 0
		} else {
			return 0, errors.Wrap(err, "failed to get votes for global public key",
				"round", round,
			)
		}
	}

	// Increase and store
	newCount := current + 1
	if err := k.GlobalPubKeyVotes.Set(ctx, key, newCount); err != nil {
		return 0, errors.Wrap(err, "failed to set votes",
			"round", round,
		)
	}

	return newCount, nil
}

func hashPublicCoeffs(coeffs [][]byte) []byte {
	h := sha256.New()
	for _, c := range coeffs {
		h.Write(c)
	}

	return h.Sum(nil)
}
