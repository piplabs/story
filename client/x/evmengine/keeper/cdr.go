package keeper

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/errors"
	clog "github.com/piplabs/story/lib/log"
)

func (k *Keeper) ProcessCDREvents(ctx context.Context, height uint64, logs []*ethtypes.Log) error {
	for _, ethlog := range logs {
		switch ethlog.Topics[0] {
		case types.CDRVaultReadEvent.ID:
			if err := k.ProcessCDRVaultRead(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process CDRVaultRead", err)
				continue
			}

		case types.CDREncryptedPartialDecryptionSubmittedEvent.ID:
			if err := k.ProcessDKGPartialDecryptionSubmitted(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process CDREncryptedPartialDecryptionSubmitted", err)
				continue
			}

		case types.CDRFeeCollectedEvent.ID:
			if err := k.ProcessCDRFeeCollected(ctx, ethlog); err != nil {
				clog.Error(ctx, "Failed to process CDRFeeCollected", err)
				continue
			}
		}
	}

	clog.Debug(ctx, "Processed DKG events", "height", height, "count", len(logs))

	return nil
}

// ProcessCDRVaultRead handles VaultRead events emitted by the CDR contract.
func (k *Keeper) ProcessCDRVaultRead(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.cdrContract.ParseVaultRead(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse VaultRead log")
	}

	latestRound, err := k.dkgKeeper.GetLatestActiveRound(cachedCtx)
	if err != nil {
		return errors.Wrap(err, "get latest active round")
	}
	if latestRound == nil {
		return errors.New("no active DKG round")
	}

	round := latestRound.Round

	label := uuidToLabel(ev.Uuid)

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
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRequester, ev.Requester.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGCiphertextLen, strconv.Itoa(len(ev.Ciphertext))),
				sdk.NewAttribute(types.AttributeKeyDKGLabelLen, strconv.Itoa(len(label))),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ethlog.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.ThresholdDecryptRequested(cachedCtx, round, ev.RequesterPubKey, ev.Ciphertext, label[:], uint64(sdkCtx.BlockHeight())); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "handle ThresholdDecryptRequested")
	}

	return nil
}

func uuidToLabel(uuid uint32) [32]byte {
	var label [32]byte
	binary.BigEndian.PutUint32(label[28:], uuid)

	return label
}

// ProcessDKGPartialDecryptionSubmitted handles PartialDecryptionSubmitted events emitted by the DKG contract.
func (k *Keeper) ProcessDKGPartialDecryptionSubmitted(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	ev, err := k.cdrContract.ParseEncryptedPartialDecryptionSubmitted(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse EncryptedPartialDecryptionSubmitted log")
	}

	label := uuidToLabel(ev.Uuid)

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(types.EventTypeDKGPartialDecryptionSubmittedSuccess)
		} else {
			e = sdk.NewEvent(
				types.EventTypeDKGPartialDecryptionSubmittedFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGValidator, ev.Validator.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGPid, strconv.FormatUint(uint64(ev.Pid), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGEncryptedPartLen, strconv.Itoa(len(ev.EncryptedPartial))),
				sdk.NewAttribute(types.AttributeKeyDKGEphemeralKeyLen, strconv.Itoa(len(ev.EphemeralPubKey))),
				sdk.NewAttribute(types.AttributeKeyDKGPubShareLen, strconv.Itoa(len(ev.PubShare))),
				sdk.NewAttribute(types.AttributeKeyDKGRequesterPubKeyLen, strconv.Itoa(len(ev.RequesterPubKey))),
				sdk.NewAttribute(types.AttributeKeyDKGCiphertextLen, strconv.Itoa(len(ev.Ciphertext))),
				sdk.NewAttribute(types.AttributeKeyDKGLabelLen, strconv.Itoa(len(label))),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ethlog.TxHash.Bytes())),
			),
		})
	}()

	partialErr := k.dkgKeeper.PartialDecryptionSubmitted(
		cachedCtx,
		ev.Validator,
		ev.Round,
		ev.Pid,
		ev.EncryptedPartial,
		ev.EphemeralPubKey,
		ev.PubShare,
		ev.RequesterPubKey,
		ev.Ciphertext,
		label[:],
		ev.Signature,
	)

	if partialErr == nil {
		if err := k.dkgKeeper.IncrementCDRPartialSubmitCount(cachedCtx, ev.Validator); err != nil {
			partialErr = errors.Wrap(err, "increment CDR submit count")
		} else if ev.Fee != nil && ev.Fee.Sign() > 0 {
			if err := k.dkgKeeper.RefundCDRFee(cachedCtx, ev.Validator, ev.Fee); err != nil {
				partialErr = errors.Wrap(err, "refund CDR fee")
			}
		}
	}

	if errors.Is(partialErr, sdkerrors.ErrInvalidRequest) {
		err = errors.WrapErrWithCode(errors.InvalidRequest, partialErr)
	} else if partialErr != nil {
		err = errors.Wrap(partialErr, "handle PartialDecryptionSubmitted")
	}

	return err
}

// ProcessCDRFeeCollected handles FeeCollected events, routing all CDR fees to the reward pool.
func (k *Keeper) ProcessCDRFeeCollected(ctx context.Context, ethlog *ethtypes.Log) (err error) {
	ev, err := k.cdrContract.ParseFeeCollected(*ethlog)
	if err != nil {
		return errors.Wrap(err, "parse FeeCollected log")
	}

	if ev.Amount == nil || ev.Amount.Sign() == 0 {
		return nil
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	defer func() {
		var e sdk.Event

		if err == nil {
			writeCache()
			e = sdk.NewEvent(types.EventTypeCDRFeeCollectedSuccess)
		}
		if err != nil {
			e = sdk.NewEvent(
				types.EventTypeCDRFeeCollectedFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ethlog.TxHash.Bytes())),
				sdk.NewAttribute(types.AttributeKeyCDRFeeType, strconv.FormatUint(uint64(ev.FeeType), 10)),
				sdk.NewAttribute(types.AttributeKeyCDRFeeAmount, ev.Amount.String()),
			),
		})
	}()

	if err = k.dkgKeeper.AddCDRFeeToPool(cachedCtx, ev.Amount); err != nil {
		clog.Error(ctx, "Failed to add CDR fee to pool", err, "fee_type", ev.FeeType, "amount", ev.Amount.String())
		return errors.Wrap(err, "add CDR fee to pool")
	}

	return nil
}
