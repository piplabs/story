//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/errors"
	clog "github.com/piplabs/story/lib/log"
)

func (k *Keeper) ProcessDKGEvents(ctx context.Context, height uint64, logs []*ethtypes.Log) error {
	for _, ethlog := range logs {
		switch ethlog.Topics[0] {
		case types.DKGRegisteredEvent.ID:
			if err := k.ProcessDKGRegistered(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGInitialized", err)
				continue
			}

		case types.DKGFinalizedEvent.ID:
			if err := k.ProcessDKGFinalized(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGFinalized", err)
				continue
			}

		case types.DKGMinReqRegisteredParticipantsSetEvent.ID:
			if err := k.ProcessDKGMinReqRegisteredParticipantsSet(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGMinReqRegisteredParticipantsSet", err)
				continue
			}

		case types.DKGMinReqFinalizedParticipantsSetEvent.ID:
			if err := k.ProcessDKGMinReqFinalizedParticipantsSet(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGMinReqFinalizedParticipantsSet", err)
				continue
			}

		case types.DKGOperationalThresholdSetEvent.ID:
			if err := k.ProcessDKGOperationalThresholdSet(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGOperationalThresholdSet", err)
				continue
			}

		case types.DKGUpgradeScheduledEvent.ID:
			if err := k.ProcessDKGUpgradeScheduled(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGUpgradeScheduled", err)
				continue
			}

		case types.DKGUpgradeCancelledEvent.ID:
			if err := k.ProcessDKGUpgradeCancelled(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGUpgradeCancelled", err)
				continue
			}
		}

		clog.Debug(ctx, "Processed DKG events", "height", height, "count", len(logs))

	}
	return nil
}

func (k *Keeper) ProcessDKGRegistered(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseRegistered(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse DKGInitialized log")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeDKGInitializedSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGInitializedFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGValidator, ev.ValidatorAddr.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGCodeCommitment, hex.EncodeToString(ev.CodeCommitment[:])),
				sdk.NewAttribute(types.AttributeKeyDKGStartBlockHeight, ev.StartBlockHeight.String()),
				sdk.NewAttribute(types.AttributeKeyDKGStartBlockHash, hex.EncodeToString(ev.StartBlockHash[:])),
				sdk.NewAttribute(types.AttributeKeyDKGDkgPubKey, hex.EncodeToString(ev.DkgPubKey)),
				sdk.NewAttribute(types.AttributeKeyDKGCommPubKey, hex.EncodeToString(ev.EnclaveCommKey)),
				sdk.NewAttribute(types.AttributeKeyDKGEnclaveReport, hex.EncodeToString(ev.ValidationContext)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.Registered(cachedCtx, ev.ValidatorAddr, ev.CodeCommitment, ev.Round, ev.StartBlockHeight, ev.StartBlockHash, ev.EnclaveType, ev.DkgPubKey, ev.EnclaveCommKey, ev.EnclaveReport); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "initialize DKG")
	}

	return nil
}

func (k *Keeper) ProcessDKGFinalized(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseFinalized(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse DKGFinalized log")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeDKGFinalizedSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGFinalizedFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGValidator, ev.ValidatorAddr.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGCodeCommitment, hex.EncodeToString(ev.CodeCommitment[:])),
				sdk.NewAttribute(types.AttributeKeyDKGParticipantsRoot, hex.EncodeToString(ev.ParticipantsRoot[:])),
				sdk.NewAttribute(types.AttributeKeyDKGSignature, hex.EncodeToString(ev.Signature)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.Finalized(cachedCtx, ev.Round, ev.ValidatorAddr, ev.CodeCommitment, ev.ParticipantsRoot, ev.Signature, ev.GlobalPubKey, ev.PublicCoeffs, ev.PubKeyShare); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "finalize DKG")
	}

	return nil
}

// ProcessDKGMinReqRegisteredParticipantsSet handles MinReqRegisteredParticipantsSet events from DKG.sol.
func (k *Keeper) ProcessDKGMinReqRegisteredParticipantsSet(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseMinReqRegisteredParticipantsSet(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse MinReqRegisteredParticipantsSet log")
	}

	newValue := uint32(ev.NewMinReqRegisteredParticipants.Uint64())

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(types.EventTypeDKGMinReqRegisteredParticipantsSetSuccess)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGMinReqRegisteredParticipantsSetFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGMinReqRegisteredParticipants, strconv.FormatUint(uint64(newValue), 10)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.SetMinReqRegisteredParticipants(cachedCtx, newValue); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "set min req registered participants")
	}

	return nil
}

// ProcessDKGMinReqFinalizedParticipantsSet handles MinReqFinalizedParticipantsSet events from DKG.sol.
func (k *Keeper) ProcessDKGMinReqFinalizedParticipantsSet(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseMinReqFinalizedParticipantsSet(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse MinReqFinalizedParticipantsSet log")
	}

	newValue := uint32(ev.NewMinReqFinalizedParticipants.Uint64())

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(types.EventTypeDKGMinReqFinalizedParticipantsSetSuccess)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGMinReqFinalizedParticipantsSetFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGMinReqFinalizedParticipants, strconv.FormatUint(uint64(newValue), 10)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.SetMinReqFinalizedParticipants(cachedCtx, newValue); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "set min req finalized participants")
	}

	return nil
}

// ProcessDKGOperationalThresholdSet handles OperationalThresholdSet events from DKG.sol.
func (k *Keeper) ProcessDKGOperationalThresholdSet(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseOperationalThresholdSet(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse OperationalThresholdSet log")
	}

	newValue := uint32(ev.NewOperationalThreshold.Uint64())

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(types.EventTypeDKGOperationalThresholdSetSuccess)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGOperationalThresholdSetFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGOperationalThreshold, strconv.FormatUint(uint64(newValue), 10)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.SetOperationalThreshold(cachedCtx, newValue); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "set operational threshold")
	}

	return nil
}

// ProcessDKGUpgradeScheduled handles UpgradeScheduled events emitted by the DKG contract.
func (k *Keeper) ProcessDKGUpgradeScheduled(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseUpgradeScheduled(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse UpgradeScheduled log")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(types.EventTypeDKGUpgradeScheduledSuccess)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGUpgradeScheduledFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGActivationHeight, ev.ActivationHeight.String()),
				sdk.NewAttribute(types.AttributeKeyDKGUpgradeVersion, ev.UpgradeVersion),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	// Validate activationHeight is positive and fits safely within int64 (Cosmos SDK block height type).
	if !ev.ActivationHeight.IsInt64() || ev.ActivationHeight.Int64() <= 0 {
		err = errors.New("activation height must be positive and fit within int64",
			"activation_height", ev.ActivationHeight.String(),
			"max_int64", math.MaxInt64,
		)

		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	}

	// Activation height must be in the future relative to the current block.
	if ev.ActivationHeight.Int64() <= sdkCtx.BlockHeight() {
		err = errors.New("activation height must be greater than current block height",
			"activation_height", ev.ActivationHeight.Int64(),
			"current_block_height", sdkCtx.BlockHeight(),
		)

		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	}

	if err = k.dkgKeeper.UpgradeScheduled(cachedCtx, ev.ActivationHeight.Int64(), ev.UpgradeVersion); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "schedule TEE upgrade")
	}

	return nil
}

// ProcessDKGUpgradeCancelled handles UpgradeCancelled events emitted by the DKG contract.
func (k *Keeper) ProcessDKGUpgradeCancelled(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, parseErr := k.dkgContract.ParseUpgradeCancelled(*ethlog)
	if parseErr != nil {
		return errors.Wrap(parseErr, "parse UpgradeCancelled log")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(types.EventTypeDKGUpgradeCancelledSuccess)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGUpgradeCancelledFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ethlog.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.UpgradeCancelled(cachedCtx, ev.UpgradeVersion); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "cancel TEE upgrade")
	}

	return nil
}

// func (k *Keeper) ProcessDKGDealComplaintsSubmitted(ctx context.Context, ethlog *ethtypes.Log) (err error) {
// 	sdkCtx := sdk.UnwrapSDKContext(ctx)
// 	cachedCtx, writeCache := sdkCtx.CacheContext()

// 	ev, err := k.dkgContract.ParseDealComplaintsSubmitted(*ethlog)
// 	if err != nil {
// 		return errors.Wrap(err, "parse DealComplaintsSubmitted log")
// 	}

// 	defer func() {
// 		if r := recover(); r != nil {
// 			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
// 		}

// 		var e sdk.Event
// 		if err == nil {
// 			writeCache()
// 			e = sdk.NewEvent(
// 				types.EventTypeDKGDealComplaintsSubmittedSuccess,
// 			)
// 		} else {
// 			e = sdk.NewEvent(
// 				types.EventTypeDKGDealComplaintsSubmittedFailure,
// 				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
// 			)
// 		}

// 		// Convert uint32 slice to string slice for attribute
// 		complainIndexesStr := make([]string, len(ev.ComplainIndexes))
// 		for i, idx := range ev.ComplainIndexes {
// 			complainIndexesStr[i] = strconv.FormatUint(uint64(idx), 10)
// 		}

// 		sdkCtx.EventManager().EmitEvents(sdk.Events{
// 			e.AppendAttributes(
// 				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
// 				sdk.NewAttribute(types.AttributeKeyDKGIndex, strconv.FormatUint(uint64(ev.Index), 10)),
// 				sdk.NewAttribute(types.AttributeKeyDKGComplainIndexes, strings.Join(complainIndexesStr, ",")),
// 				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
// 				sdk.NewAttribute(types.AttributeKeyDKGCodeCommitment, hex.EncodeToString(ev.CodeCommitment[:])),
// 				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
// 			),
// 		})
// 	}()

// 	if err = k.dkgKeeper.DealComplaintsSubmitted(cachedCtx, ev.Index, ev.ComplainIndexes, ev.Round, ev.CodeCommitment); errors.Is(err, sdkerrors.ErrInvalidRequest) {
// 		return errors.WrapErrWithCode(errors.InvalidRequest, err)
// 	} else if err != nil {
// 		return errors.Wrap(err, "submit deal complaints")
// 	}

// 	return nil
// }

// func (k *Keeper) ProcessDKGDealVerified(ctx context.Context, ethlog *ethtypes.Log) (err error) {
// 	sdkCtx := sdk.UnwrapSDKContext(ctx)
// 	cachedCtx, writeCache := sdkCtx.CacheContext()

// 	ev, err := k.dkgContract.ParseDealVerified(*ethlog)
// 	if err != nil {
// 		return errors.Wrap(err, "parse DealVerified log")
// 	}

// 	defer func() {
// 		if r := recover(); r != nil {
// 			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
// 		}

// 		var e sdk.Event
// 		if err == nil {
// 			writeCache()
// 			e = sdk.NewEvent(
// 				types.EventTypeDKGDealVerifiedSuccess,
// 			)
// 		} else {
// 			e = sdk.NewEvent(
// 				types.EventTypeDKGDealVerifiedFailure,
// 				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
// 			)
// 		}

// 		sdkCtx.EventManager().EmitEvents(sdk.Events{
// 			e.AppendAttributes(
// 				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
// 				sdk.NewAttribute(types.AttributeKeyDKGIndex, strconv.FormatUint(uint64(ev.Index), 10)),
// 				sdk.NewAttribute(types.AttributeKeyDKGRecipientIndex, strconv.FormatUint(uint64(ev.RecipientIndex), 10)),
// 				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
// 				sdk.NewAttribute(types.AttributeKeyDKGCodeCommitment, hex.EncodeToString(ev.CodeCommitment[:])),
// 				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
// 			),
// 		})
// 	}()

// 	if err = k.dkgKeeper.DealVerified(cachedCtx, ev.Index, ev.RecipientIndex, ev.Round, ev.CodeCommitment); errors.Is(err, sdkerrors.ErrInvalidRequest) {
// 		return errors.WrapErrWithCode(errors.InvalidRequest, err)
// 	} else if err != nil {
// 		return errors.Wrap(err, "verify deal")
// 	}

// 	return nil
// }

// func (k *Keeper) ProcessDKGInvalidDeal(ctx context.Context, ethlog *ethtypes.Log) (err error) {
// 	sdkCtx := sdk.UnwrapSDKContext(ctx)
// 	cachedCtx, writeCache := sdkCtx.CacheContext()

// 	ev, err := k.dkgContract.ParseInvalidDeal(*ethlog)
// 	if err != nil {
// 		return errors.Wrap(err, "parse InvalidDeal log")
// 	}

// 	defer func() {
// 		if r := recover(); r != nil {
// 			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
// 		}

// 		var e sdk.Event
// 		if err == nil {
// 			writeCache()
// 			e = sdk.NewEvent(
// 				types.EventTypeDKGInvalidDealSuccess,
// 			)
// 		} else {
// 			e = sdk.NewEvent(
// 				types.EventTypeDKGInvalidDealFailure,
// 				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
// 			)
// 		}

// 		sdkCtx.EventManager().EmitEvents(sdk.Events{
// 			e.AppendAttributes(
// 				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
// 				sdk.NewAttribute(types.AttributeKeyDKGIndex, strconv.FormatUint(uint64(ev.Index), 10)),
// 				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
// 				sdk.NewAttribute(types.AttributeKeyDKGCodeCommitment, hex.EncodeToString(ev.CodeCommitment[:])),
// 				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
// 			),
// 		})
// 	}()

// 	if err = k.dkgKeeper.InvalidDeal(cachedCtx, ev.Index, ev.Round, ev.CodeCommitment); errors.Is(err, sdkerrors.ErrInvalidRequest) {
// 		return errors.WrapErrWithCode(errors.InvalidRequest, err)
// 	} else if err != nil {
// 		return errors.Wrap(err, "process invalid deal")
// 	}

// 	return nil
// }
