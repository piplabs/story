package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/epochs/types"
	"github.com/piplabs/story/lib/errors"
)

// BeginBlocker of epochs module.
func (k *Keeper) BeginBlocker(ctx context.Context) error {
	start := telemetry.Now()
	defer telemetry.ModuleMeasureSince(types.ModuleName, start, telemetry.MetricKeyBeginBlocker)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockTime := sdkCtx.BlockTime()
	blockHeight := sdkCtx.BlockHeight()

	err := k.EpochInfo.Walk(
		ctx,
		nil,
		func(_ string, epochInfo types.EpochInfo) (stop bool, err error) {
			// If blocktime < initial epoch start time, return
			if blockTime.Before(epochInfo.StartTime) {
				return false, nil
			}
			// if epoch counting hasn't started, signal we need to start.
			shouldInitialEpochStart := !epochInfo.EpochCountingStarted

			epochEndTime := epochInfo.CurrentEpochStartTime.Add(epochInfo.Duration)
			shouldEpochStart := (blockTime.After(epochEndTime)) || shouldInitialEpochStart

			if !shouldEpochStart {
				return false, nil
			}
			epochInfo.CurrentEpochStartHeight = blockHeight

			if shouldInitialEpochStart {
				epochInfo.EpochCountingStarted = true
				epochInfo.CurrentEpoch = 1
				epochInfo.CurrentEpochStartTime = epochInfo.StartTime
				sdkCtx.Logger().Debug(fmt.Sprintf("Starting new epoch with identifier %s epoch number %d", epochInfo.Identifier, epochInfo.CurrentEpoch))
			} else {
				err := sdkCtx.EventManager().EmitTypedEvent(&types.EventEpochEnd{
					EpochNumber: epochInfo.CurrentEpoch,
				})
				if err != nil {
					return false, errors.Wrap(err, "failed to emit epoch end event")
				}

				cacheCtx, writeFn := sdkCtx.CacheContext()
				//nolint:contextcheck // use cached context
				if err := k.AfterEpochEnd(cacheCtx, epochInfo.Identifier, epochInfo.CurrentEpoch); err != nil {
					// purposely ignoring the error here not to halt the chain if the hook fails
					sdkCtx.Logger().Error(fmt.Sprintf("Error after epoch end with identifier %s epoch number %d", epochInfo.Identifier, epochInfo.CurrentEpoch))
				} else {
					writeFn()
				}

				epochInfo.CurrentEpoch++
				epochInfo.CurrentEpochStartTime = epochInfo.CurrentEpochStartTime.Add(epochInfo.Duration)
				sdkCtx.Logger().Debug(fmt.Sprintf("Starting epoch with identifier %s epoch number %d", epochInfo.Identifier, epochInfo.CurrentEpoch))
			}

			// emit new epoch start event, set epoch info, and run BeforeEpochStart hook
			err = sdkCtx.EventManager().EmitTypedEvent(&types.EventEpochStart{
				EpochNumber:    epochInfo.CurrentEpoch,
				EpochStartTime: epochInfo.CurrentEpochStartTime.Unix(),
			})
			if err != nil {
				return false, errors.Wrap(err, "failed to emit epoch start event")
			}
			err = k.EpochInfo.Set(ctx, epochInfo.Identifier, epochInfo)
			if err != nil {
				sdkCtx.Logger().Error(fmt.Sprintf("Error set epoch info with identifier %s epoch number %d", epochInfo.Identifier, epochInfo.CurrentEpoch))
				//nolint:nilerr // return nil per original x/epochs code
				return false, nil
			}

			cacheCtx, writeFn := sdkCtx.CacheContext()
			//nolint:contextcheck // use cached context
			if err := k.BeforeEpochStart(cacheCtx, epochInfo.Identifier, epochInfo.CurrentEpoch); err != nil {
				// purposely ignoring the error here not to halt the chain if the hook fails
				sdkCtx.Logger().Error(fmt.Sprintf("Error before epoch start with identifier %s epoch number %d", epochInfo.Identifier, epochInfo.CurrentEpoch))
			} else {
				writeFn()
			}

			return false, nil
		},
	)

	return err
}
