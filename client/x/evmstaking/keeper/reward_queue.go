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

// AddRewardWithdrawalToQueue inserts a reward withdrawal into the queue.
func (k Keeper) AddRewardWithdrawalToQueue(ctx context.Context, withdrawal types.Withdrawal) error {
	return k.RewardWithdrawalQueue.Enqueue(ctx, withdrawal)
}

func (k Keeper) DequeueEligibleRewardWithdrawals(ctx context.Context, maxDequeue uint32) (withdrawals etypes.Withdrawals, err error) {
	// front is the unique monotonically increasing index of a withdrawal in the queue.
	// It's used as the value in etypes.Withdrawal.Index for later validation purposes,
	// when evmengine's msg_server receives withdrawals as part of the execution payload
	// and needs to verify that the received withdrawals are in the correct order from
	// the front of the queue.
	front, err := k.RewardWithdrawalQueue.Front(ctx)
	if err != nil {
		log.Debug(ctx, "Front", "err", err)
		return nil, err
	}

	for i := range maxDequeue {
		withdrawal, err := k.RewardWithdrawalQueue.Dequeue(ctx)
		if err != nil {
			// Dequeue will return ErrEmptyQueue if the queue is empty
			if errors.Is(err, addcollections.ErrEmptyQueue) {
				break
			}

			return nil, err
		}
		withdrawals = append(withdrawals, &etypes.Withdrawal{
			Index: front + uint64(i), // increment front by i to get the correct index in the loop
			// used to distinguish between unstake and reward withdrawals as EL does not care about the value
			Validator: uint64(withdrawal.WithdrawalType),
			Address:   common.HexToAddress(withdrawal.ExecutionAddress),
			Amount:    withdrawal.Amount,
		})
	}

	return withdrawals, nil
}

func (k Keeper) PeekEligibleRewardWithdrawals(ctx context.Context, maxPeek uint32) (withdrawals etypes.Withdrawals, err error) {
	if k.RewardWithdrawalQueue.IsEmpty(ctx) {
		return withdrawals, nil
	}

	// front is the unique monotonically increasing index of a withdrawal in the queue.
	// It's used as the value in etypes.Withdrawal.Index for later validation purposes,
	// when evmengine's msg_server receives withdrawals as part of the execution payload
	// and needs to verify that the received withdrawals are in the correct order from
	// the front of the queue.
	front, err := k.RewardWithdrawalQueue.Front(ctx)
	if err != nil {
		return nil, err
	}

	for i := range maxPeek {
		// NOTE: Get adjusts the provided index by the front index of the queue
		withdrawal, err := k.RewardWithdrawalQueue.Get(ctx, uint64(i))
		if err != nil {
			// Get will return ErrOutOfBoundsQueue if the queue is empty
			if errors.Is(err, addcollections.ErrOutOfBoundsQueue) {
				break
			}

			return nil, err
		}
		withdrawals = append(withdrawals, &etypes.Withdrawal{
			Index: front + uint64(i), // increment front by i to get the correct index in the loop
			// used to distinguish between unstake and reward withdrawals as EL does not care about the value
			Validator: uint64(withdrawal.WithdrawalType),
			Address:   common.HexToAddress(withdrawal.ExecutionAddress),
			Amount:    withdrawal.Amount,
		})
	}

	return withdrawals, nil
}
