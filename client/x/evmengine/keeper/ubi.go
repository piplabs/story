//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"strconv"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	clog "github.com/piplabs/story/lib/log"
)

func (k *Keeper) ProcessUbiEvents(ctx context.Context, height uint64, logs []*types.EVMEvent) error {
	for _, evmLog := range logs {
		if err := evmLog.Verify(); err != nil {
			return errors.Wrap(err, "verify log [BUG]")
		}
		ethlog, err := evmLog.ToEthLog()
		if err != nil {
			return err
		}

		//nolint:gocritic,revive // more cases later
		switch ethlog.Topics[0] {
		case types.UBIPercentageSetEvent.ID:
			ev, err := k.ubiContract.ParseUBIPercentageSet(ethlog)
			if err != nil {
				clog.Error(ctx, "Failed to parse UBIPercentageSet log", err)
				continue
			}
			if err = k.ProcessUBIPercentageSet(ctx, ev); err != nil {
				clog.Error(ctx, "Failed to process UBI percentage set", err)
				continue
			}
		}
	}

	clog.Debug(ctx, "Processed UBIPool events", "height", height, "count", len(logs))

	return nil
}

func (k *Keeper) ProcessUBIPercentageSet(ctx context.Context, ev *bindings.UBIPoolUBIPercentageSet) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	defer func() {
		if err == nil {
			writeCache()
			return
		}
		sdkCtx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeUpdateUbiFailure,
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyUbiPercentage, strconv.FormatUint(uint64(ev.Percentage), 10)),
				sdk.NewAttribute(types.AttributeKeyStatusCode, errors.UnwrapErrCode(err).String()),
			),
		})
	}()

	newUBI := math.LegacyNewDecWithPrec(int64(ev.Percentage), 4)

	if err = k.distrKeeper.SetUbi(cachedCtx, newUBI); errors.Is(err, sdkerrors.ErrInvalidRequest) {
		return errors.WrapErrWithCode(errors.InvalidRequest, err)
	} else if err != nil {
		return errors.Wrap(err, "set new UBI percentage")
	}

	return nil
}
