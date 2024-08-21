package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient"
	"github.com/piplabs/story/lib/log"
)

type msgServer struct {
	*Keeper
	types.UnimplementedMsgServiceServer
}

// NewMsgServerImpl returns an implementation of the MsgServer interface for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServiceServer {
	return &msgServer{Keeper: keeper}
}

// ExecutionPayload handles a new execution payload included in the current finalized block.
// This is called as part of FinalizeBlock ABCI++ method.
func (s msgServer) ExecutionPayload(ctx context.Context, msg *types.MsgExecutionPayload) (*types.ExecutionPayloadResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if sdkCtx.ExecMode() != sdk.ExecModeFinalize {
		return nil, errors.New("only allowed in finalize mode")
	}

	payload, err := s.parseAndVerifyProposedPayload(ctx, msg)
	if err != nil {
		return nil, err
	}

	// TODO: should we compare and reject in a finalized block?
	//// Ensure that the withdrawals in the payload are from the front indices of the queue.
	// if err := s.compareWithdrawals(ctx, payload.Withdrawals); err != nil {
	//	return nil, errors.Wrap(err, "compare local and received withdrawals")
	//}

	// TODO: We dequeue with assumption that the top items of the queue are the ones that are processed in the block.
	// TODO: We might need to check that the withdrawals in the finalized block are the same as the ones dequeued.
	// But there's no way to reject the block at this point, so we can only log an error.
	log.Debug(ctx, "Dequeueing eligible withdrawals [BEFORE]", "len", len(payload.Withdrawals))
	if _, err := s.evmstakingKeeper.DequeueEligibleWithdrawals(ctx); err != nil {
		return nil, errors.Wrap(err, "error on withdrawals dequeue")
	}
	log.Debug(ctx, "Dequeueing eligible withdrawals [AFTER]", "len", len(payload.Withdrawals))

	err = retryForever(ctx, func(ctx context.Context) (bool, error) {
		status, err := pushPayload(ctx, s.engineCl, payload)
		if err != nil || isUnknown(status) {
			// We need to retry forever on networking errors, but can't easily identify them, so retry all errors.
			log.Warn(ctx, "Processing finalized payload failed: push new payload to evm (will retry)", err,
				"status", status.Status)

			return false, nil // Retry
		} else if invalid, err := isInvalid(status); invalid {
			// This should never happen. This node will stall now.
			log.Error(ctx, "Processing finalized payload failed; payload invalid [BUG]", err)

			return false, err // Don't retry, error out.
		} else if isSyncing(status) {
			log.Warn(ctx, "Processing finalized payload; evm syncing", nil)
		}

		return true, nil // We are done, don't retry
	})
	if err != nil {
		return nil, err
	}

	// CometBFT has instant finality, so head/safe/finalized is latest height.
	fcs := engine.ForkchoiceStateV1{
		HeadBlockHash:      payload.BlockHash,
		SafeBlockHash:      payload.BlockHash,
		FinalizedBlockHash: payload.BlockHash,
	}

	err = retryForever(ctx, func(ctx context.Context) (bool, error) {
		fcr, err := s.engineCl.ForkchoiceUpdatedV3(ctx, fcs, nil)
		if err != nil || isUnknown(fcr.PayloadStatus) {
			// We need to retry forever on networking errors, but can't easily identify them, so retry all errors.
			log.Warn(ctx, "Processing finalized payload failed: evm fork choice update (will retry)", err,
				"status", fcr.PayloadStatus.Status)

			return false, nil // Retry
		} else if isSyncing(fcr.PayloadStatus) {
			log.Warn(ctx, "Processing finalized payload halted while evm syncing (will retry)", nil, "payload_height", payload.Number)

			return false, nil // Retry
		} else if invalid, err := isInvalid(fcr.PayloadStatus); invalid {
			// This should never happen. This node will stall now.
			log.Error(ctx, "Processing finalized payload failed; forkchoice update invalid [BUG]", err,
				"payload_height", payload.Number)

			return false, err // Don't retry
		}

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	// Deliver all the previous payload log events
	if err := s.evmstakingKeeper.ProcessStakingEvents(ctx, payload.Number-1, msg.PrevPayloadEvents); err != nil {
		return nil, errors.Wrap(err, "deliver staking-related event logs")
	}
	if err := s.ProcessUpgradeEvents(ctx, payload.Number-1, msg.PrevPayloadEvents); err != nil {
		return nil, errors.Wrap(err, "deliver upgrade-related event logs")
	}

	if err := s.updateExecutionHead(ctx, payload); err != nil {
		return nil, errors.Wrap(err, "update execution head")
	}

	return &types.ExecutionPayloadResponse{}, nil
}

//nolint:unused // compareWithdrawals compares the given actual withdrawals with the expected withdrawals from the queue.
func (s msgServer) compareWithdrawals(ctx context.Context, actualWithdrawals etypes.Withdrawals) error {
	expectedWithdrawals, err := s.evmstakingKeeper.PeekEligibleWithdrawals(ctx)
	if err != nil {
		return errors.Wrap(err, "peek withdrawals")
	}

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

// pushPayload pushes the given Engine API payload to EL and returns the engine payload status or an error.
func pushPayload(ctx context.Context, engineCl ethclient.EngineClient, payload engine.ExecutableData) (engine.PayloadStatusV1, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	appHash := common.BytesToHash(sdkCtx.BlockHeader().AppHash)
	if appHash == (common.Hash{}) {
		return engine.PayloadStatusV1{}, errors.New("app hash is empty")
	}

	emptyVersionHashes := make([]common.Hash, 0) // Cannot use nil.

	// Push it back to the execution client (mark it as possible new head).
	status, err := engineCl.NewPayloadV3(ctx, payload, emptyVersionHashes, &appHash)
	if err != nil {
		return engine.PayloadStatusV1{}, errors.Wrap(err, "new payload")
	}

	return status, nil
}

var _ types.MsgServiceServer = msgServer{}

func isUnknown(status engine.PayloadStatusV1) bool {
	if status.Status == engine.VALID ||
		status.Status == engine.INVALID ||
		status.Status == engine.SYNCING ||
		status.Status == engine.ACCEPTED {
		return false
	}

	return true
}

func isSyncing(status engine.PayloadStatusV1) bool {
	return status.Status == engine.SYNCING || status.Status == engine.ACCEPTED
}

func isInvalid(status engine.PayloadStatusV1) (bool, error) {
	if status.Status != engine.INVALID {
		return false, nil
	}

	valErr := "nil"
	if status.ValidationError != nil {
		valErr = *status.ValidationError
	}

	hash := "nil"
	if status.LatestValidHash != nil {
		hash = status.LatestValidHash.Hex()
	}

	return true, errors.New("payload invalid", "validation_err", valErr, "last_valid_hash", hash)
}
