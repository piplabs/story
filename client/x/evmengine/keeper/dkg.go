//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"strconv"

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
				sdk.NewAttribute(types.AttributeKeyDKGRawQuote, hex.EncodeToString(ev.ValidationContext)),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.Registered(cachedCtx, ev.ValidatorAddr, ev.CodeCommitment, ev.Round, ev.StartBlockHeight, ev.StartBlockHash, ev.DkgPubKey, ev.EnclaveCommKey, ev.ValidationContext); errors.Is(err, sdkerrors.ErrInvalidRequest) {
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
				sdk.NewAttribute(types.AttributeKeyDKGCodeCommitment, hex.EncodeToString(ev.CodeCommitment[:])),
				sdk.NewAttribute(types.AttributeKeyDKGRequester, requester.Hex()),
				sdk.NewAttribute(types.AttributeKeyDKGCiphertextLen, strconv.Itoa(len(ev.Ciphertext))),
				sdk.NewAttribute(types.AttributeKeyDKGLabelLen, strconv.Itoa(len(ev.Label))),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ethlog.TxHash.Bytes())),
			),
		})
	}()

	if err = k.dkgKeeper.ThresholdDecryptRequested(cachedCtx, requester, ev.Round, ev.CodeCommitment, ev.RequesterPubKey, ev.Ciphertext, ev.Label); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "handle ThresholdDecryptRequested")
	}

	return nil
}
