package keeper

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"

	addcollections "github.com/piplabs/story/client/collections"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/log"
)

// AddWithdrawalToQueue inserts a withdrawal into the queue.
func (k Keeper) AddWithdrawalToQueue(ctx context.Context, withdrawal types.Withdrawal) error {
	return k.WithdrawalQueue.Enqueue(ctx, withdrawal)
}

func (k Keeper) DequeueEligibleWithdrawals(ctx context.Context, maxDequeue uint32) (withdrawals etypes.Withdrawals, err error) {
	// front is the unique monotonically increasing index of a withdrawal in the queue.
	// It's used as the value in etypes.Withdrawal.Index for later validation purposes,
	// when evmengine's msg_server receives withdrawals as part of the execution payload
	// and needs to verify that the received withdrawals are in the correct order from
	// the front of the queue.
	front, err := k.WithdrawalQueue.Front(ctx)
	if err != nil {
		log.Debug(ctx, "Front", "err", err)
		return nil, err
	}

	for i := range uint64(maxDequeue) {
		withdrawal, err := k.WithdrawalQueue.Dequeue(ctx)
		if err != nil {
			// Dequeue will return ErrEmptyQueue if the queue is empty
			if errors.Is(err, addcollections.ErrEmptyQueue) {
				break
			}

			return nil, err
		}
		withdrawals = append(withdrawals, &etypes.Withdrawal{
			Index:     front + i, // increment front by i to get the correct index in the loop
			Validator: 0,         // does not matter for EL
			Address:   common.HexToAddress(withdrawal.ExecutionAddress),
			Amount:    withdrawal.Amount,
		})
	}

	return withdrawals, nil
}

func (k Keeper) PeekEligibleWithdrawals(ctx context.Context, maxPeek uint32) (withdrawals etypes.Withdrawals, err error) {
	if k.WithdrawalQueue.IsEmpty(ctx) {
		return withdrawals, nil
	}

	// front is the unique monotonically increasing index of a withdrawal in the queue.
	// It's used as the value in etypes.Withdrawal.Index for later validation purposes,
	// when evmengine's msg_server receives withdrawals as part of the execution payload
	// and needs to verify that the received withdrawals are in the correct order from
	// the front of the queue.
	front, err := k.WithdrawalQueue.Front(ctx)
	if err != nil {
		return nil, err
	}

	for i := range uint64(maxPeek) {
		// NOTE: Get adjusts the provided index by the front index of the queue
		withdrawal, err := k.WithdrawalQueue.Get(ctx, i)
		if err != nil {
			// Get will return ErrOutOfBoundsQueue if the queue is empty
			if errors.Is(err, addcollections.ErrOutOfBoundsQueue) {
				break
			}

			return nil, err
		}
		withdrawals = append(withdrawals, &etypes.Withdrawal{
			Index:     front + i, // increment front by i to get the correct index in the loop
			Validator: 0,         // does not matter for EL
			Address:   common.HexToAddress(withdrawal.ExecutionAddress),
			Amount:    withdrawal.Amount,
		})
	}

	return withdrawals, nil
}

// GetAllWithdrawals gets the set of all stake withdrawals with no limits.
func (k Keeper) GetAllWithdrawals(ctx context.Context) (withdrawals []types.Withdrawal, err error) {
	iterator, err := k.WithdrawalQueue.Iterate(ctx)
	if err != nil {
		return nil, err
	}

	wdrKvs, err := iterator.KeyValues()
	if err != nil {
		return nil, err
	}

	for _, withdrawal := range wdrKvs {
		withdrawals = append(withdrawals, withdrawal.Value)
	}

	return withdrawals, nil
}

// GetWithdrawals returns at max the requested amount of withdrawals.
func (k Keeper) GetWithdrawals(ctx context.Context, maxRetrieve uint32) (withdrawals []types.Withdrawal, err error) {
	iterator, err := k.WithdrawalQueue.Iterate(ctx)
	if err != nil {
		return nil, err
	}

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		withdrawal, err := iterator.Value()
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdrawal)
		i++
	}

	return withdrawals[:i], nil // trim if the array length < maxRetrieve
}
