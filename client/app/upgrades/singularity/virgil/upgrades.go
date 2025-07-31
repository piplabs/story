package virgil

import (
	"context"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

const (
	// NewShortPeriodDuration defines the duration of the new short period.
	NewShortPeriodDuration = time.Second * 7776000 // 90 days
)

func CreateUpgradeHandler(
	_ *module.Manager,
	_ module.Configurator,
	keepers *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		blockHeight := sdkCtx.BlockHeight()
		log.Info(ctx, "Current block height", "Height", blockHeight)

		// Check if the upgrade is needed for current chain
		chainID := sdkCtx.ChainID()
		if _, ok := GetUpgradeHeight(chainID); !ok {
			log.Info(ctx, "Virgil upgrade not needed for current chain, skip", "ChainID", chainID)
			return vm, nil
		}
		log.Info(ctx, "Start Virgil upgrade", "ChainID", chainID)

		log.Info(ctx, "Get current staking params...")
		stakingParams, err := keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get staking params")
		}

		log.Info(ctx, "Update staking periods...")
		newRewardsMultiplier := GetRewardsMultipliers(chainID)
		var oldShortPeriodDuration time.Duration
		for i := range stakingParams.Periods {
			if stakingParams.Periods[i].PeriodType == 1 {
				oldShortPeriodDuration = stakingParams.Periods[i].Duration
				log.Info(ctx, "SKIP: Change short period duration to 90 days (7776000 seconds)")
				
				// stakingParams.Periods[i].Duration = NewShortPeriodDuration
				log.Info(ctx, "Existing short period duration", "Time", stakingParams.Periods[i].Duration.String())
				log.Info(ctx, "Change short period rewards multiplier", "new_multiplier", newRewardsMultiplier.Short.String())
				stakingParams.Periods[i].RewardsMultiplier = newRewardsMultiplier.Short
			} else if stakingParams.Periods[i].PeriodType == 2 {
				log.Info(ctx, "Existing medium period duration", "Time", stakingParams.Periods[i].Duration.String())
				log.Info(ctx, "Change medium period rewards multiplier", "new_multiplier", newRewardsMultiplier.Medium.String())
				stakingParams.Periods[i].RewardsMultiplier = newRewardsMultiplier.Medium
			} else if stakingParams.Periods[i].PeriodType == 3 {
				log.Info(ctx, "Existing long period duration", "Time", stakingParams.Periods[i].Duration.String())
				log.Info(ctx, "Change long period rewards multiplier", "new_multiplier", newRewardsMultiplier.Long.String())
				stakingParams.Periods[i].RewardsMultiplier = newRewardsMultiplier.Long
			}
		}

		log.Info(ctx, "Apply staking param changes...")
		if err := keepers.StakingKeeper.SetParams(ctx, stakingParams); err != nil {
			return vm, errors.Wrap(err, "failed to update staking params")
		}

		log.Info(ctx, "Check new staking params...")
		stakingParams, err = keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get staking params")
		}

		for _, p := range stakingParams.Periods {
			if p.PeriodType == 1 { //nolint:nestif // no issue
				log.Info(ctx, "New short period duration", "Time", p.Duration.String())
				// if p.Duration != NewShortPeriodDuration {
				// 	return vm, errors.New("new short period duration is not correct")
				// }
				if !p.RewardsMultiplier.Equal(newRewardsMultiplier.Short) {
					return vm, errors.New("new short period rewards multiplier is not correct")
				}
			} else if p.PeriodType == 2 {
				log.Info(ctx, "New medium period duration", "Time", p.Duration.String())
				if !p.RewardsMultiplier.Equal(newRewardsMultiplier.Medium) {
					return vm, errors.New("new medium period rewards multiplier is not correct")
				}
			} else if p.PeriodType == 3 {
				log.Info(ctx, "New long period duration", "Time", p.Duration.String())
				if !p.RewardsMultiplier.Equal(newRewardsMultiplier.Long) {
					return vm, errors.New("new long period rewards multiplier is not correct")
				}
			}
		}

		periodDelegations, err := keepers.StakingKeeper.GetAllPeriodDelegations(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get all period delegations")
		}
		log.Info(ctx, "Sweep all delegations and modify short period delegations", "Count", len(periodDelegations))
		for i := range periodDelegations {
			if periodDelegations[i].PeriodType != 1 {
				continue
			}

			log.Info(
				ctx, "Find short period delegation",
				"Delegator", periodDelegations[i].DelegatorAddress,
				"Validator", periodDelegations[i].ValidatorAddress,
				"Shares", periodDelegations[i].Shares.String(),
				"RewardsShares", periodDelegations[i].RewardsShares.String(),
				"PeriodDelegationId", periodDelegations[i].PeriodDelegationId,
				"PeriodType", periodDelegations[i].PeriodType,
				"EndTime", periodDelegations[i].EndTime.String(),
			)

			oldEndTime := periodDelegations[i].EndTime
			newEndTime := oldEndTime.Add(NewShortPeriodDuration - oldShortPeriodDuration)
			log.Info(ctx, "New short period delegation", "EndTime", newEndTime)

			log.Info(ctx, "Set new short period delegation")
			periodDelegations[i].EndTime = newEndTime
			delAddr, err := sdk.AccAddressFromBech32(periodDelegations[i].DelegatorAddress)
			if err != nil {
				return vm, errors.Wrap(err, "failed to get delegator address")
			}
			valAddr, err := sdk.ValAddressFromBech32(periodDelegations[i].ValidatorAddress)
			if err != nil {
				return vm, errors.Wrap(err, "failed to get validator address")
			}
			if err := keepers.StakingKeeper.SetPeriodDelegation(ctx, delAddr, valAddr, periodDelegations[i]); err != nil {
				return vm, errors.Wrap(err, "failed to set period delegation")
			}

			newPeriodDelegation, err := keepers.StakingKeeper.GetPeriodDelegation(ctx, delAddr, valAddr, periodDelegations[i].PeriodDelegationId)
			if err != nil {
				return vm, errors.Wrap(err, "failed to get new period delegation")
			}
			log.Info(
				ctx, "Get new short period delegation",
				"Delegator", periodDelegations[i].DelegatorAddress,
				"Validator", periodDelegations[i].ValidatorAddress,
				"Shares", periodDelegations[i].Shares.String(),
				"RewardsShares", periodDelegations[i].RewardsShares.String(),
				"PeriodDelegationId", periodDelegations[i].PeriodDelegationId,
				"PeriodType", periodDelegations[i].PeriodType,
				"EndTime", periodDelegations[i].EndTime.String(),
			)
			if newPeriodDelegation.EndTime != newEndTime {
				return vm, errors.New("new period delegation end time is not correct")
			}
		}

		log.Info(ctx, "Virgil upgrade complete")

		return vm, nil
	}
}
