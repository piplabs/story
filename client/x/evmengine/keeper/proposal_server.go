package keeper

import (
	"context"
	"fmt"

	etypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

type proposalServer struct {
	*Keeper
	types.UnimplementedMsgServiceServer
}

// ExecutionPayload handles a new execution payload proposed in a block.
func (s proposalServer) ExecutionPayload(ctx context.Context, msg *types.MsgExecutionPayload) (*types.ExecutionPayloadResponse, error) {
	payload, err := s.parseAndVerifyProposedPayload(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal payload")
	}

	// Ensure that the withdrawals in the payload are from the front indices of the queue.
	if err := s.compareWithdrawals(ctx, payload.Withdrawals); err != nil {
		return nil, errors.Wrap(err, "compare local and received withdrawals")
	}

	// Push the payload to the EVM.
	err = retryForever(ctx, func(ctx context.Context) (bool, error) {
		status, err := pushPayload(ctx, s.engineCl, payload)
		if err != nil || isUnknown(status) {
			// We need to retry forever on networking errors, but can't easily identify them, so retry all errors.
			log.Warn(ctx, "Verifying proposal failed: push new payload to evm (will retry)", err,
				"status", status.Status)

			return false, nil // Retry
		} else if invalid, err := isInvalid(status); invalid {
			return false, errors.Wrap(err, "invalid payload, rejecting proposal") // Don't retry
		} else if isSyncing(status) {
			// If this is initial sync, we need to continue and set a target head to sync to, so don't retry.
			log.Warn(ctx, "Can't properly verifying proposal: evm syncing", err,
				"payload_height", payload.Number)
		}

		return true, nil // We are done, don't retry.
	})
	if err != nil {
		return nil, err
	}

	return &types.ExecutionPayloadResponse{}, nil
}

// NewProposalServer returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewProposalServer(keeper *Keeper) types.MsgServiceServer {
	return &proposalServer{Keeper: keeper}
}

var _ types.MsgServiceServer = proposalServer{}

// compareWithdrawals compares the local peek and received withdrawals.
func (s proposalServer) compareWithdrawals(ctx context.Context, actualWithdrawals etypes.Withdrawals) error {
	maxWithdrawals, err := s.evmstakingKeeper.MaxWithdrawalPerBlock(ctx)
	if err != nil {
		return errors.Wrap(err, "get max withdrawals per block")
	}

	expectedWithdrawals, err := s.evmstakingKeeper.PeekEligibleWithdrawals(ctx, maxWithdrawals)
	if err != nil {
		return errors.Wrap(err, "peek withdrawals")
	}

	if len(expectedWithdrawals) > len(actualWithdrawals) {
		return fmt.Errorf(
			"expected withdrawals %v should not greater than actual withdrawals %v",
			len(expectedWithdrawals), len(actualWithdrawals),
		)
	}

	maxRewardWithdrawals := maxWithdrawals - uint32(len(expectedWithdrawals))
	expectedRewardWithdrawals, err := s.evmstakingKeeper.PeekEligibleRewardWithdrawals(ctx, maxRewardWithdrawals)
	if err != nil {
		return errors.Wrap(err, "peek reward withdrawals")
	}

	expectedTotalWithdrawals := len(expectedWithdrawals) + len(expectedRewardWithdrawals)
	log.Debug(ctx, "Comparing local and received withdrawals",
		"local", expectedTotalWithdrawals,
		"received", len(actualWithdrawals),
	)
	if expectedTotalWithdrawals != len(actualWithdrawals) {
		return fmt.Errorf(
			"expected total withdrawals %v should equal to actual withdrawals %v",
			expectedTotalWithdrawals, len(actualWithdrawals),
		)
	}

	pos := 0
	for i := range expectedWithdrawals {
		if expectedWithdrawals[i].Index != actualWithdrawals[pos].Index {
			return errors.New("invalid withdrawal index")
		}
		// skip the Validator index equality check (always 0)
		if expectedWithdrawals[i].Address != actualWithdrawals[pos].Address {
			return errors.New("invalid withdrawal address")
		}
		if expectedWithdrawals[i].Amount != actualWithdrawals[pos].Amount {
			return errors.New("invalid withdrawal amount")
		}
		pos++
	}
	for i := range expectedRewardWithdrawals {
		if expectedRewardWithdrawals[i].Index != actualWithdrawals[pos].Index {
			return errors.New("invalid withdrawal index")
		}
		// skip the Validator index equality check (always 0)
		if expectedRewardWithdrawals[i].Address != actualWithdrawals[pos].Address {
			return errors.New("invalid withdrawal address")
		}
		if expectedRewardWithdrawals[i].Amount != actualWithdrawals[pos].Amount {
			return errors.New("invalid withdrawal amount")
		}
		pos++
	}

	return nil
}
