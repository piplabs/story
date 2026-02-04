package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/piplabs/story/lib/errors"
)

// AddGlobalPubKeyVote Increase vote for global public key by 1.
func (k *Keeper) AddGlobalPubKeyVote(ctx context.Context, codeCommitment [32]byte, round uint32, globalPubKey []byte, publicCoeffs [][]byte) (uint32, error) {
	coeffHash := hex.EncodeToString(hashPublicCoeffs(publicCoeffs))

	key := fmt.Sprintf(
		"%s_%d_%s_%s",
		hex.EncodeToString(codeCommitment[:]),
		round,
		hex.EncodeToString(globalPubKey),
		coeffHash,
	)

	// Read current votes (0 if not found)
	current, err := k.GlobalPubKeyVotes.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			current = 0
		} else {
			return 0, errors.Wrap(err, "failed to get votes for global public key",
				"code_commitment", hex.EncodeToString(codeCommitment[:]),
				"round", round,
			)
		}
	}

	// Increase and store
	newCount := current + 1
	if err := k.GlobalPubKeyVotes.Set(ctx, key, newCount); err != nil {
		return 0, errors.Wrap(err, "failed to set votes",
			"code_commitment", hex.EncodeToString(codeCommitment[:]),
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
