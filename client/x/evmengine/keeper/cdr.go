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
		// TODO: Add to handle VaultAllocated, VaultWritten, and EncryptedPartialDecryptionSubmitted events
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
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(ev.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGRequester, ev.Requester.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGCiphertextLen, strconv.Itoa(len(ev.Ciphertext))),
				sdk.NewAttribute(types.AttributeKeyDKGLabelLen, strconv.Itoa(len(label))),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ethlog.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.ThresholdDecryptRequested(cachedCtx, ev.Round, ev.RequesterPubKey, ev.Ciphertext, label[:], uint64(sdkCtx.BlockHeight())); errors.Is(err, sdkerrors.ErrInvalidRequest) {
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
				sdk.NewAttribute(types.AttributeKeyDKGLabelLen, strconv.Itoa(len(label))),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ethlog.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.PartialDecryptionSubmitted(
		cachedCtx,
		ev.Validator,
		ev.Round,
		ev.Pid,
		ev.EncryptedPartial,
		ev.EphemeralPubKey,
		ev.PubShare,
		ev.RequesterPubKey,
		label[:],
		ev.Signature,
	); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "handle PartialDecryptionSubmitted")
	}

	return nil
}
