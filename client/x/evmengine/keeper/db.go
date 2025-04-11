package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/lib/errors"
)

func (h *ExecutionHead) Hash() common.Hash {
	return common.BytesToHash(h.GetBlockHash())
}

// executionHeadID is the ID of the singleton execution head row in the database.
const executionHeadID = 1

// InsertGenesisHead inserts the genesis execution head into the database.
func (k *Keeper) InsertGenesisHead(ctx context.Context, executionBlockHash []byte) error {
	id, err := k.headTable.InsertReturningId(ctx, &ExecutionHead{
		CreatedHeight: 0, // genesis
		BlockHeight:   0, // genesis
		BlockHash:     executionBlockHash,
		BlockTime:     0, // Timestamp isn't critical, skip it in genesis.
	})
	if err != nil {
		return errors.Wrap(err, "insert genesis head")
	} else if id != executionHeadID {
		return errors.New("unexpected genesis head id", "id", id)
	}

	return nil
}

// GetExecutionHead returns the current execution head.
func (k *Keeper) GetExecutionHead(ctx context.Context) (*ExecutionHead, error) {
	head, err := k.headTable.Get(ctx, executionHeadID)
	if err != nil {
		return nil, errors.Wrap(err, "get execution head")
	}

	return head, nil
}

// UpdateExecutionHead updates the execution head with the given payload.
func (k *Keeper) UpdateExecutionHead(ctx context.Context, payload engine.ExecutableData) error {
	return k.updateExecutionHead(ctx, payload.Number, payload.BlockHash, payload.Timestamp)
}

func (k *Keeper) UpdateExecutionHeadWithBlock(ctx context.Context, blockHeight uint64, blockHash common.Hash, blockTime uint64) error {
	return k.updateExecutionHead(ctx, blockHeight, blockHash, blockTime)
}

func (k *Keeper) updateExecutionHead(
	ctx context.Context, blockHeight uint64, blockHash common.Hash, blockTime uint64,
) error {
	head := &ExecutionHead{
		Id:            executionHeadID,
		CreatedHeight: uint64(sdk.UnwrapSDKContext(ctx).BlockHeight()),
		BlockHeight:   blockHeight,
		BlockHash:     blockHash.Bytes(),
		BlockTime:     blockTime,
	}

	err := k.headTable.Update(ctx, head)
	if err != nil {
		return errors.Wrap(err, "update execution head")
	}

	return nil
}
