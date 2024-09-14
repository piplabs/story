package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/epochs/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// BeginBlocker of epochs module.
func (k Keeper) BeginBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	// NOTE(Narangde): Use UnwrapSDKContext instead of Environment's HeaderService
	headerInfo := sdk.UnwrapSDKContext(ctx).HeaderInfo()
	err := k.EpochInfo.Walk(
		ctx,
		nil,
		func(_ string, epochInfo types.EpochInfo) (stop bool, err error) {
			// If blocktime < initial epoch start time, return
			if headerInfo.Time.Before(epochInfo.StartTime) {
				return false, nil
			}
			// if epoch counting hasn't started, signal we need to start.
			shouldInitialEpochStart := !epochInfo.EpochCountingStarted

			epochEndTime := epochInfo.CurrentEpochStartTime.Add(epochInfo.Duration)
			shouldEpochStart := (headerInfo.Time.After(epochEndTime)) || shouldInitialEpochStart

			if !shouldEpochStart {
				return false, nil
			}
			epochInfo.CurrentEpochStartHeight = headerInfo.Height

			if shouldInitialEpochStart {
				epochInfo.EpochCountingStarted = true
				epochInfo.CurrentEpoch = 1
				epochInfo.CurrentEpochStartTime = epochInfo.StartTime
				log.Debug(ctx, "Starting new epoch", "epoch_identifier", epochInfo.Identifier, "current_epoch", epochInfo.CurrentEpoch)
			} else {
				err := k.EventService.EventManager(ctx).Emit(ctx, &types.EventEpochEnd{
					EpochNumber: epochInfo.CurrentEpoch,
				})
				if err != nil {
					return false, errors.Wrap(err, "emit epoch end event")
				}

				if err := k.AfterEpochEnd(ctx, epochInfo.Identifier, epochInfo.CurrentEpoch); err != nil {
					// purposely ignoring the error here not to halt the chain if the hook fails
					log.Error(ctx, "Error after epoch end", err, "epoch_identifier", epochInfo.Identifier, "current_epoch", epochInfo.CurrentEpoch)
				}

				epochInfo.CurrentEpoch++
				epochInfo.CurrentEpochStartTime = epochInfo.CurrentEpochStartTime.Add(epochInfo.Duration)
				log.Debug(ctx, "Starting epoch with", "epoch_identifier", epochInfo.Identifier, "current_epoch", epochInfo.CurrentEpoch)
			}

			// emit new epoch start event, set epoch info, and run BeforeEpochStart hook
			err = k.EventService.EventManager(ctx).Emit(ctx, &types.EventEpochStart{
				EpochNumber:    epochInfo.CurrentEpoch,
				EpochStartTime: epochInfo.CurrentEpochStartTime.Unix(),
			})
			if err != nil {
				return false, errors.Wrap(err, "emit epoch start event")
			}
			err = k.EpochInfo.Set(ctx, epochInfo.Identifier, epochInfo)
			if err != nil {
				log.Error(ctx, "Error set epoch info", err, "epoch_identifier", epochInfo.Identifier, "current_epoch", epochInfo.CurrentEpoch)
				return false, nil
			}
			if err := k.BeforeEpochStart(ctx, epochInfo.Identifier, epochInfo.CurrentEpoch); err != nil {
				// purposely ignoring the error here not to halt the chain if the hook fails
				log.Error(ctx, "Error before epoch start", err, "epoch_identifier", epochInfo.Identifier, "current_epoch", epochInfo.CurrentEpoch)
			}

			return false, nil
		},
	)

	return err
}
