package keeper

import (
	"context"

	"github.com/cosmos/gogoproto/proto"
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

	// Collect local view of the evm logs from the previous payload.
	evmEvents, err := s.evmEvents(ctx, payload.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "prepare evm event logs")
	}

	// Ensure the proposed evm event logs are equal to the local view.
	if err := evmEventsEqual(evmEvents, msg.PrevPayloadEvents); err != nil {
		return nil, errors.Wrap(err, "verify prev payload events")
	}

	return &types.ExecutionPayloadResponse{}, nil
}

// NewProposalServer returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewProposalServer(keeper *Keeper) types.MsgServiceServer {
	return &proposalServer{Keeper: keeper}
}

var _ types.MsgServiceServer = proposalServer{}

// TODO: benchmark this function, might be adding to overhead. Esp. the for loop. If so, parallelize since array checks are independent.
func evmEventsEqual(a, b []*types.EVMEvent) error {
	if len(a) != len(b) {
		return errors.New("count mismatch")
	}

	for i := range a {
		if !proto.Equal(a[i], b[i]) {
			return errors.New("log mismatch", "index", i)
		}
	}

	return nil
}

// compareWithdrawals compares the local peek and received withdrawals.
func (s proposalServer) compareWithdrawals(ctx context.Context, actualWithdrawals etypes.Withdrawals) error {
	expectedWithdrawals, err := s.evmstakingKeeper.PeekEligibleWithdrawals(ctx)
	if err != nil {
		return errors.Wrap(err, "peek withdrawals")
	}

	log.Debug(ctx, "Comparing local and received withdrawals",
		"local", len(expectedWithdrawals),
		"received", len(actualWithdrawals),
	)
	if len(actualWithdrawals) != len(expectedWithdrawals) {
		return errors.New("invalid withdrawals length")
	}

	for i, withdrawal := range actualWithdrawals {
		if withdrawal.Index != expectedWithdrawals[i].Index {
			return errors.New("invalid withdrawal index")
		}
		// skip the Validator index equality check (always 0)
		if withdrawal.Address != expectedWithdrawals[i].Address {
			return errors.New("invalid withdrawal address")
		}
		if withdrawal.Amount != expectedWithdrawals[i].Amount {
			return errors.New("invalid withdrawal amount")
		}
	}

	return nil
}
