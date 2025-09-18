//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

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
		case types.DKGInitializedEvent.ID:
			if err := k.ProcessDKGInitialized(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGInitialized", err)
				continue
			}

		case types.DKGCommitmentsUpdatedEvent.ID:
			if err := k.ProcessDKGCommitmentsUpdated(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGCommitmentsUpdated", err)
				continue
			}

		case types.DKGFinalizedEvent.ID:
			if err := k.ProcessDKGFinalized(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGFinalized", err)
				continue
			}

		case types.DKGUpgradeScheduledEvent.ID:
			if err := k.ProcessDKGUpgradeScheduled(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGUpgradeScheduled", err)
				continue
			}

		case types.DKGRegistrationChallengedEvent.ID:
			if err := k.ProcessDKGRegistrationChallenged(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGRegistrationChallenged", err)
				continue
			}

		case types.DKGInvalidDKGInitializationEvent.ID:
			if err := k.ProcessDKGInvalidDKGInitialization(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGInvalidDKGInitialization", err)
				continue
			}

		case types.DKGRemoteAttestationProcessedOnChainEvent.ID:
			if err := k.ProcessDKGRemoteAttestationProcessedOnChain(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGRemoteAttestationProcessedOnChain", err)
				continue
			}

		case types.DKGDealComplaintsSubmittedEvent.ID:
			if err := k.ProcessDKGDealComplaintsSubmitted(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGDealComplaintsSubmitted", err)
				continue
			}

		case types.DKGDealVerifiedEvent.ID:
			if err := k.ProcessDKGDealVerified(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGDealVerified", err)
				continue
			}

		case types.DKGInvalidDealEvent.ID:
			if err := k.ProcessDKGInvalidDeal(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGInvalidDeal", err)
				continue
			}
		}
	}

	clog.Debug(ctx, "Processed DKG events", "height", height, "count", len(logs))

	return nil
}

//nolint:dupl // ProcessDKGInitialized and ProcessDKGFinalized have similar structure but different logic
func (k *Keeper) ProcessDKGInitialized(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseDKGInitialized(*ethlog)
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
				sdk.NewAttribute(types.AttributeKeyDKGIndex, strconv.FormatUint(uint64(ev.Index), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave)),
				sdk.NewAttribute(types.AttributeKeyDKGPubKey, hex.EncodeToString(ev.PubKey)),
				sdk.NewAttribute(types.AttributeKeyDKGRemoteReport, hex.EncodeToString(ev.RemoteReport)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.Initialized(cachedCtx, ev.Mrenclave, ev.Round, ev.Index, ev.PubKey, ev.RemoteReport); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "initialize DKG")
	}

	return nil
}

func (k *Keeper) ProcessDKGCommitmentsUpdated(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseDKGCommitmentsUpdated(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse DKGCommitmentsUpdated log")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeDKGCommitmentsUpdatedSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGCommitmentsUpdatedFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGTotal, strconv.FormatUint(uint64(ev.Total), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGThreshold, strconv.FormatUint(uint64(ev.Threshold), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGIndex, strconv.FormatUint(uint64(ev.Index), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGCommitments, hex.EncodeToString(ev.Commitments)),
				sdk.NewAttribute(types.AttributeKeyDKGSignature, hex.EncodeToString(ev.Signature)),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.CommitmentsUpdated(cachedCtx, ev.Round, ev.Total, ev.Threshold, ev.Index, ev.Commitments, ev.Signature, ev.Mrenclave); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "update DKG commitments")
	}

	return nil
}

//nolint:dupl // ProcessDKGInitialized and ProcessDKGFinalized have similar structure but different logic
func (k *Keeper) ProcessDKGFinalized(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseDKGFinalized(*ethlog)
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
				sdk.NewAttribute(types.AttributeKeyDKGIndex, strconv.FormatUint(uint64(ev.Index), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGFinalized, strconv.FormatBool(ev.Finalized)),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave)),
				sdk.NewAttribute(types.AttributeKeyDKGSignature, hex.EncodeToString(ev.Signature)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.Finalized(cachedCtx, ev.Round, ev.Index, ev.Finalized, ev.Mrenclave, ev.Signature); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "finalize DKG")
	}

	return nil
}

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
			e = sdk.NewEvent(
				types.EventTypeDKGUpgradeScheduledSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGUpgradeScheduledFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGActivationHeight, strconv.FormatUint(uint64(ev.ActivationHeight), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.UpgradeScheduled(cachedCtx, ev.ActivationHeight, ev.Mrenclave); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "schedule DKG upgrade")
	}

	return nil
}

func (k *Keeper) ProcessDKGRegistrationChallenged(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseRegistrationChallenged(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse RegistrationChallenged log")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeDKGRegistrationChallengedSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGRegistrationChallengedFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave)),
				sdk.NewAttribute(types.AttributeKeyDKGChallenger, ev.Challenger.Hex()),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.RegistrationChallenged(cachedCtx, ev.Round, ev.Mrenclave, ev.Challenger); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "challenge DKG registration")
	}

	return nil
}

func (k *Keeper) ProcessDKGInvalidDKGInitialization(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseInvalidDKGInitialization(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse InvalidDKGInitialization log")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeDKGInvalidDKGInitializationSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGInvalidDKGInitializationFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGIndex, strconv.FormatUint(uint64(ev.Index), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGValidator, ev.Validator.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.InvalidDKGInitialization(cachedCtx, ev.Round, ev.Index, ev.Validator, ev.Mrenclave); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "process invalid DKG initialization")
	}

	return nil
}

func (k *Keeper) ProcessDKGRemoteAttestationProcessedOnChain(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseRemoteAttestationProcessedOnChain(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse RemoteAttestationProcessedOnChain log")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeDKGRemoteAttestationProcessedOnChainSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGRemoteAttestationProcessedOnChainFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGIndex, strconv.FormatUint(uint64(ev.Index), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGValidator, ev.Validator.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGChalStatus, strconv.FormatUint(uint64(ev.ChalStatus), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.RemoteAttestationProcessedOnChain(cachedCtx, ev.Index, ev.Validator, int(ev.ChalStatus), ev.Round, ev.Mrenclave); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "process remote attestation on chain")
	}

	return nil
}

func (k *Keeper) ProcessDKGDealComplaintsSubmitted(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseDealComplaintsSubmitted(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse DealComplaintsSubmitted log")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeDKGDealComplaintsSubmittedSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGDealComplaintsSubmittedFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		// Convert uint32 slice to string slice for attribute
		complainIndexesStr := make([]string, len(ev.ComplainIndexes))
		for i, idx := range ev.ComplainIndexes {
			complainIndexesStr[i] = strconv.FormatUint(uint64(idx), 10)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGIndex, strconv.FormatUint(uint64(ev.Index), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGComplainIndexes, strings.Join(complainIndexesStr, ",")),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.DealComplaintsSubmitted(cachedCtx, ev.Index, ev.ComplainIndexes, ev.Round, ev.Mrenclave); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "submit deal complaints")
	}

	return nil
}

func (k *Keeper) ProcessDKGDealVerified(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseDealVerified(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse DealVerified log")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeDKGDealVerifiedSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGDealVerifiedFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGIndex, strconv.FormatUint(uint64(ev.Index), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRecipientIndex, strconv.FormatUint(uint64(ev.RecipientIndex), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.DealVerified(cachedCtx, ev.Index, ev.RecipientIndex, ev.Round, ev.Mrenclave); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "verify deal")
	}

	return nil
}

func (k *Keeper) ProcessDKGInvalidDeal(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseInvalidDeal(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse InvalidDeal log")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeDKGInvalidDealSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGInvalidDealFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGIndex, strconv.FormatUint(uint64(ev.Index), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.InvalidDeal(cachedCtx, ev.Index, ev.Round, ev.Mrenclave); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "process invalid deal")
	}

	return nil
}
