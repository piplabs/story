package keeper

import (
	"context"
	"fmt"

	dbm "github.com/cosmos/cosmos-db"

	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

const (
	// offChainPrimaryPfx prefixes primary data entries: "p|{primaryKey}" → JSON bytes.
	offChainPrimaryPfx = "p|"
	// offChainRoundPfx prefixes round-index entries: "r|{round:010d}|{primaryKey}" → empty.
	offChainRoundPfx = "r|"
)

// PartialDecryptOffChainStore is an off-chain KV store (backed by dbm.DB) that
// mirrors DKGPartialDecrypt submissions with operator-configurable round retention.
// It is write-through alongside the on-chain IAVL store and serves GetCDRPartials
// queries when available, enabling longer history without bloating the app hash.
type PartialDecryptOffChainStore struct {
	db dbm.DB
}

// NewPartialDecryptOffChainStore returns a store backed by db.
func NewPartialDecryptOffChainStore(db dbm.DB) *PartialDecryptOffChainStore {
	return &PartialDecryptOffChainStore{db: db}
}

func (s *PartialDecryptOffChainStore) primaryKey(key string) []byte {
	return []byte(offChainPrimaryPfx + key)
}

func (s *PartialDecryptOffChainStore) roundIndexKey(round uint32, primaryKey string) []byte {
	return []byte(fmt.Sprintf("%s%010d|%s", offChainRoundPfx, round, primaryKey))
}

// Set writes the primary entry and its round index entry atomically.
func (s *PartialDecryptOffChainStore) Set(primaryKey string, round uint32, bz []byte) error {
	batch := s.db.NewBatch()
	defer batch.Close()

	if err := batch.Set(s.primaryKey(primaryKey), bz); err != nil {
		return errors.Wrap(err, "off-chain partial decrypt: set primary")
	}
	if err := batch.Set(s.roundIndexKey(round, primaryKey), []byte{}); err != nil {
		return errors.Wrap(err, "off-chain partial decrypt: set round index")
	}

	return errors.Wrap(batch.Write(), "off-chain partial decrypt: batch write")
}

// Has returns whether the primary entry exists (used for duplicate detection).
func (s *PartialDecryptOffChainStore) Has(primaryKey string) (bool, error) {
	ok, err := s.db.Has(s.primaryKey(primaryKey))
	return ok, errors.Wrap(err, "off-chain partial decrypt: has")
}

// PrefixIterator returns an iterator over all primary entries whose key begins
// with prefix. The caller must call Close() on the returned iterator.
// Returns (nil, nil) when no entries match.
func (s *PartialDecryptOffChainStore) PrefixIterator(prefix string) (dbm.Iterator, error) {
	start := []byte(offChainPrimaryPfx + prefix)
	end := dbPrefixEndBytes(start)

	iter, err := s.db.Iterator(start, end)
	if err != nil {
		return nil, errors.Wrap(err, "off-chain partial decrypt: prefix iterator")
	}

	return iter, nil
}

// PruneBeforeRound deletes all entries whose round <= cutoffRound using the
// secondary round index for an O(deleted entries) scan.
func (s *PartialDecryptOffChainStore) PruneBeforeRound(ctx context.Context, cutoffRound uint32) error {
	start := []byte(offChainRoundPfx)
	end := []byte(fmt.Sprintf("%s%010d|", offChainRoundPfx, cutoffRound+1))

	iter, err := s.db.Iterator(start, end)
	if err != nil {
		return errors.Wrap(err, "off-chain partial decrypt: round index iterator for pruning")
	}

	// Collect keys before deleting to avoid mutating the iterator mid-scan.
	var indexKeys [][]byte
	for ; iter.Valid(); iter.Next() {
		k := make([]byte, len(iter.Key()))
		copy(k, iter.Key())
		indexKeys = append(indexKeys, k)
	}
	iter.Close()

	if len(indexKeys) == 0 {
		return nil
	}

	batch := s.db.NewBatch()
	defer batch.Close()

	// offChainRoundPfx(2) + round(10) + "|"(1) = 13-byte prefix
	const indexPfxLen = 13
	for _, indexKey := range indexKeys {
		if len(indexKey) <= indexPfxLen {
			log.Warn(ctx, "Skipping malformed off-chain round index key during pruning", nil, "key", string(indexKey))
			continue
		}
		pk := string(indexKey[indexPfxLen:])
		if err := batch.Delete(s.primaryKey(pk)); err != nil {
			return errors.Wrap(err, "off-chain partial decrypt: delete primary during pruning")
		}
		if err := batch.Delete(indexKey); err != nil {
			return errors.Wrap(err, "off-chain partial decrypt: delete round index during pruning")
		}
	}

	return errors.Wrap(batch.Write(), "off-chain partial decrypt: pruning batch write")
}

// dbPrefixEndBytes returns the exclusive upper bound for a prefix scan.
// It increments the last byte of prefix, handling overflow by shortening.
func dbPrefixEndBytes(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}
	end := make([]byte, len(prefix))
	copy(end, prefix)
	for i := len(end) - 1; i >= 0; i-- {
		end[i]++
		if end[i] != 0 {
			return end[:i+1]
		}
	}
	return nil // overflow — no upper bound
}
