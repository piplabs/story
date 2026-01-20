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
	"github.com/ethereum/go-ethereum/common"
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

		case types.DKGNetworkSetEvent.ID:
			if err := k.ProcessDKGNetworkSet(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGNetworkSet", err)
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

		case types.DKGThresholdDecryptRequestedEvent.ID:
			if err := k.ProcessDKGThresholdDecryptRequested(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process DKGThresholdDecryptRequested", err)
				continue
			}
		}
	}

	clog.Debug(ctx, "Processed DKG events", "height", height, "count", len(logs))

	return nil
}

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
				sdk.NewAttribute(types.AttributeKeyDKGValidator, ev.MsgSender.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave[:])),
				sdk.NewAttribute(types.AttributeKeyDKGDkgPubKey, hex.EncodeToString(ev.DkgPubKey)),
				sdk.NewAttribute(types.AttributeKeyDKGCommPubKey, hex.EncodeToString(ev.CommPubKey)),
				sdk.NewAttribute(types.AttributeKeyDKGRawQuote, hex.EncodeToString(ev.RawQuote)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.RegistrationInitialized(cachedCtx, ev.MsgSender, ev.Mrenclave, ev.Round, ev.DkgPubKey, ev.CommPubKey, ev.RawQuote); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "initialize DKG")
	}

	return nil
}

func (k *Keeper) ProcessDKGNetworkSet(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseDKGNetworkSet(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse DKGNetworkSet log")
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeDKGNetworkSetSuccess,
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGNetworkSetFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGValidator, ev.MsgSender.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave[:])),
				sdk.NewAttribute(types.AttributeKeyDKGTotal, strconv.FormatUint(uint64(ev.Total), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGThreshold, strconv.FormatUint(uint64(ev.Threshold), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGSignature, hex.EncodeToString(ev.Signature)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.NetworkSet(cachedCtx, ev.MsgSender, ev.Mrenclave, ev.Round, ev.Total, ev.Threshold, ev.Signature); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "set DKG network")
	}

	return nil
}

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
				sdk.NewAttribute(types.AttributeKeyDKGValidator, ev.MsgSender.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave[:])),
				sdk.NewAttribute(types.AttributeKeyDKGSignature, hex.EncodeToString(ev.Signature)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.Finalized(cachedCtx, ev.Round, ev.MsgSender, ev.Mrenclave, ev.Signature, ev.GlobalPubKey, ev.PublicCoeffs); errors.Is(err, sdkerrors.ErrInvalidRequest) {
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
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave[:])),
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
				sdk.NewAttribute(types.AttributeKeyDKGValidator, ev.Validator.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGChalStatus, strconv.FormatUint(uint64(ev.ChalStatus), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave[:])),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.RemoteAttestationProcessedOnChain(cachedCtx, ev.Validator, int(ev.ChalStatus), ev.Round, ev.Mrenclave); errors.Is(err, sdkerrors.ErrInvalidRequest) {
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
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave[:])),
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
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave[:])),
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
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave[:])),
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

// ProcessDKGThresholdDecryptRequested handles ThresholdDecryptRequested events emitted by the DKG contract.
func (k *Keeper) ProcessDKGThresholdDecryptRequested(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.dkgContract.ParseThresholdDecryptRequested(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse ThresholdDecryptRequested log")
	}

	// requester is indexed address (topic[1])
	requester := common.BytesToAddress(ethlog.Topics[1].Bytes()[12:])

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(types.EventTypeDKGThresholdDecryptRequestedSuccess)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGThresholdDecryptRequestedFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGMrenclave, hex.EncodeToString(ev.Mrenclave[:])),
				sdk.NewAttribute(types.AttributeKeyDKGRequester, requester.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGCiphertextLen, strconv.Itoa(len(ev.Ciphertext))),
				sdk.NewAttribute(types.AttributeKeyDKGLabelLen, strconv.Itoa(len(ev.Label))),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ethlog.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.ThresholdDecryptRequested(cachedCtx, requester, ev.Round, ev.Mrenclave, ev.RequesterPubKey, ev.Ciphertext, ev.Label); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "handle ThresholdDecryptRequested")
	}

	return nil
}
