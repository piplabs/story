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

	latestActive, err := k.dkgKeeper.GetLatestActiveRound(ctx)
	if err != nil {
		return errors.Wrap(err, "get latest active round")
	}
	if latestActive == nil {
		return errors.New("no active DKG round available for threshold decryption")
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
				sdk.NewAttribute(types.AttributeKeyDKGRound, strconv.FormatUint(uint64(latestActive.Round), 10)),
				sdk.NewAttribute(types.AttributeKeyDKGCiphertextLen, strconv.Itoa(len(ev.EncryptedData))),
				sdk.NewAttribute(types.AttributeKeyDKGLabelLen, strconv.Itoa(len(label))),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ethlog.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.ThresholdDecryptRequested(cachedCtx, requester, ev.Round, ev.CodeCommitment, ev.RequesterPubKey, ev.Ciphertext, ev.Label, uint64(sdkCtx.BlockHeight())); errors.Is(err, sdkerrors.ErrInvalidRequest) {
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
