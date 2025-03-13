package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"

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

	// Since we already checked the withdrawals in the proposal server, we simply check the length here.
	log.Debug(
		ctx, "Dequeueing eligible withdrawals [BEFORE]",
		"total_len", len(payload.Withdrawals),
	)
	maxWithdrawals, err := s.evmstakingKeeper.MaxWithdrawalPerBlock(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error getting max withdrawal per block")
	}
	ws, err := s.evmstakingKeeper.DequeueEligibleWithdrawals(ctx, maxWithdrawals)
	if err != nil {
		return nil, errors.Wrap(err, "error on withdrawals dequeue")
	}
	log.Debug(
		ctx, "Dequeueing eligible withdrawals [AFTER]",
		"total_len", len(payload.Withdrawals),
		"withdrawals_len", len(ws),
	)

	if len(ws) > len(payload.Withdrawals) {
		return nil, fmt.Errorf(
			"dequeued withdrawals %v should not greater than proposed withdrawals %v",
			len(ws), len(payload.Withdrawals),
		)
	}

	log.Debug(
		ctx, "Dequeueing eligible reward withdrawals [BEFORE]",
		"total_len", len(payload.Withdrawals),
		"withdrawals_len", len(ws),
	)
	maxRewardWithdrawals := maxWithdrawals - uint32(len(ws))
	rws, err := s.evmstakingKeeper.DequeueEligibleRewardWithdrawals(ctx, maxRewardWithdrawals)
	if err != nil {
		return nil, errors.Wrap(err, "error on reward withdrawals dequeue")
	}
	log.Debug(
		ctx, "Dequeueing eligible reward withdrawals [AFTER]",
		"total_len", len(payload.Withdrawals),
		"withdrawals_len", len(ws),
		"reward_withdrawals_len", len(rws),
	)

	if totalWithdrawals := len(ws) + len(rws); totalWithdrawals != len(payload.Withdrawals) {
		return nil, fmt.Errorf(
			"dequeued total withdrawals %v should equal to proposed withdrawals %v",
			totalWithdrawals, len(payload.Withdrawals),
		)
	}

	err = retryForever(ctx, func(ctx context.Context) (bool, error) {
		status, err := pushPayload(ctx, s.engineCl, payload)
		if err != nil {
			// We need to retry forever on networking errors, but can't easily identify them, so retry all errors.
			log.Warn(ctx, "Processing finalized payload failed: push new payload to evm (will retry)", err)

			return false, nil // Retry
		} else if invalid, err := isInvalid(status); invalid {
			// This should never happen. This node will stall now.
			log.Error(ctx, "Processing finalized payload failed; payload invalid [BUG]", err)

			return false, err // Don't retry, error out.
		} else if isSyncing(status) {
			log.Warn(ctx, "Processing finalized payload; evm syncing", nil)
		} /* else isValid(status) */

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
		if err != nil {
			// We need to retry forever on networking errors, but can't easily identify them, so retry all errors.
			log.Warn(ctx, "Processing finalized payload failed: evm fork choice update (will retry)", err)

			return false, nil // Retry
		} else if isSyncing(fcr.PayloadStatus) {
			log.Warn(ctx, "Processing finalized payload halted while evm syncing (will retry)", nil, "payload_height", payload.Number)

			return false, nil // Retry
		} else if invalid, err := isInvalid(fcr.PayloadStatus); invalid {
			// This should never happen. This node will stall now.
			log.Error(ctx, "Processing finalized payload failed; forkchoice update invalid [BUG]", err,
				"payload_height", payload.Number)

			return false, err // Don't retry
		} /* else isValid(status) */

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	// get events of the newly finalized block
	events, err := s.evmEvents(ctx, payload.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "fetch evm event logs")
	}

	// Deliver all the payload log events of the newly finalized block
	if err := s.evmstakingKeeper.ProcessStakingEvents(ctx, payload.Number-1, events); err != nil {
		return nil, errors.Wrap(err, "deliver staking-related event logs")
	}
	if err := s.ProcessUpgradeEvents(ctx, payload.Number-1, events); err != nil {
		return nil, errors.Wrap(err, "deliver upgrade-related event logs")
	}
	if err := s.ProcessUbiEvents(ctx, payload.Number-1, events); err != nil {
		return nil, errors.Wrap(err, "deliver ubi-related event logs")
	}

	if err := s.updateExecutionHead(ctx, payload); err != nil {
		return nil, errors.Wrap(err, "update execution head")
	}

	return &types.ExecutionPayloadResponse{}, nil
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
