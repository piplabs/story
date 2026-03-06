package horace

import (
	"context"
	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/piplabs/story/client/app/keepers"
	lerrors "github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

//nolint:gocyclo,maintidx // many changes
func CreateUpgradeHandler(
	_ *module.Manager,
	_ module.Configurator,
	keepers *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		if err := runHoraceUpgrade(
			ctx,
			keepers.AccountKeeper,
			keepers.StakingKeeper,
			keepers.DistrKeeper,
			keepers.MintKeeper,
		); err != nil {
			return vm, err
		}
		return vm, nil
	}
}

func runHoraceUpgrade(ctx context.Context, aKeeper AccountKeeper, sKeeper StakingKeeper, dKeeper DistributionKeeper, mKeeper MintKeeper) error {
	log.Info(ctx, "Start Horace upgrade")

	// --------------------------------
	// update staking params
	// --------------------------------
	log.Info(ctx, "Start updating staking params...")

	stakingParams, err := sKeeper.GetParams(ctx)
	if err != nil {
		return lerrors.Wrap(err, "get staking params")
	}

	lockedTokenType := stakingParams.LockedTokenType

	// update the locked token multiplier in x/staking module params
	// NOTE: index 0 is the locked token type - see x/staking/types/params.go:DefaultTokenTypes
	if stakingParams.TokenTypes[0].TokenType != lockedTokenType {
		return lerrors.New("locked token type not found in token types", "expected", lockedTokenType, "found", stakingParams.TokenTypes[0].TokenType)
	}
	stakingParams.TokenTypes[0].RewardsMultiplier = NewLockedTokenMultiplier

	if err := sKeeper.SetParams(ctx, stakingParams); err != nil {
		return lerrors.Wrap(err, "update staking params")
	}

	// check the updated min delegation amount & locked token multiplier in x/staking module params
	newStakingParams, err := sKeeper.GetParams(ctx)
	if err != nil {
		return lerrors.Wrap(err, "reload staking params")
	}

	if !newStakingParams.TokenTypes[0].RewardsMultiplier.Equal(NewLockedTokenMultiplier) {
		return lerrors.New(
			"locked token multiplier not updated",
			"expected", NewLockedTokenMultiplier,
			"actual", stakingParams.TokenTypes[0].RewardsMultiplier,
		)
	}

	log.Info(ctx, "Completed adjusting locked staking multiplier...")

	// --------------------------------
	// rescale rewards shares for delegations delegated on locked validators
	// --------------------------------

	log.Info(ctx, "Start updating rewards shares for locked delegations...")

	// Update the rewards shares for locked validators and their rewards tokens.
	// NOTE: In x/distribution's `AllocateTokens` function (keeper.go), calculation of each validator's allocation of
	// total block rewards is based on the validator's rewards tokens at that block. So, updating all locked delegations'
	// rewards shares will shrink the total rewards tokens of the locked validator. This rescales each validator's allocation, ie.
	// unlocked validators receive more rewards and locked validators receive less rewards. Within each validator, each
	// the rewards shares of each delegation and period delegation are updated accordingly. And, rewards stake in starting
	// info in distribution module, which is used for reward distribution, is also updated.

	validators, err := sKeeper.GetAllValidators(ctx)
	if err != nil {
		return lerrors.Wrap(err, "get all validators")
	}

	// for checking later
	oldTotalRewardsTokens := math.LegacyZeroDec()
	oldTotalRewardsTokensOfLockedValidators := math.LegacyZeroDec()
	oldTotalRewardsTokensOfUnlockedValidators := math.LegacyZeroDec()
	for _, val := range validators {
		rt := val.GetRewardsTokens()
		oldTotalRewardsTokens = oldTotalRewardsTokens.Add(rt)
		if val.SupportTokenType == lockedTokenType {
			oldTotalRewardsTokensOfLockedValidators = oldTotalRewardsTokensOfLockedValidators.Add(rt)
		} else {
			oldTotalRewardsTokensOfUnlockedValidators = oldTotalRewardsTokensOfUnlockedValidators.Add(rt)
		}
	}
	if oldTotalRewardsTokens.IsZero() {
		return lerrors.New("old total rewards tokens is zero")
	}
	if !oldTotalRewardsTokensOfLockedValidators.Add(oldTotalRewardsTokensOfUnlockedValidators).Equal(oldTotalRewardsTokens) {
		return lerrors.New("old total rewards tokens of locked validators and unlocked validators do not sum to old total rewards tokens")
	}

	// update the rewards tokens of locked validators and rewards shares of locked delegations
	for _, val := range validators {
		if val.SupportTokenType != lockedTokenType {
			continue
		}

		valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
		if err != nil {
			return lerrors.Wrap(err, "failed to convert validator bech32 address to validator address", "validator_address", val.OperatorAddress)
		}

		// get all delegations
		dels, err := sKeeper.GetValidatorDelegations(ctx, valAddr)
		if err != nil {
			return lerrors.Wrap(err, "failed to get all delegations for the validator")
		}

		newTotalRewardsShares := math.LegacyZeroDec()                                          // used for updating the total rewards shares of each validator
		newRewardsTokens := math.LegacyNewDecFromInt(val.Tokens).Mul(NewLockedTokenMultiplier) // used for updating the rewards tokens of each validator

		// withdraw validator commission before updating rewards tokens of validator
		if _, err := dKeeper.WithdrawValidatorCommission(ctx, valAddr); errors.Is(err, dtypes.ErrNoValidatorCommission) {
			log.Debug(ctx, "Skip withdrawal of commission as there is no commission", "validator_addr", valAddr.String())
		} else if err != nil {
			return lerrors.Wrap(err, "failed to withdraw validator commission")
		}

		// iterate all delegations which is delegated on the locked validator to scale the rewards shares
		for _, del := range dels {
			// withdraw delegation rewards before updating rewards shares
			delAccAddr, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
			if err != nil {
				return lerrors.Wrap(err, "failed to convert delegation delegator address")
			}

			if _, err := dKeeper.WithdrawDelegationRewards(ctx, delAccAddr, valAddr); err != nil {
				return lerrors.Wrap(err, "failed to withdraw delegation rewards")
			}

			shares := del.Shares

			newDelRewardsShares := shares.MulTruncate(NewLockedTokenMultiplier)
			newTotalRewardsShares = newTotalRewardsShares.Add(newDelRewardsShares)

			del.RewardsShares = newDelRewardsShares

			// save updated delegation
			if err := sKeeper.SetDelegation(ctx, del); err != nil {
				return lerrors.Wrap(err, "failed to set fixed delegation")
			}

			// verify after set
			delAddr, err := aKeeper.AddressCodec().StringToBytes(del.DelegatorAddress)
			if err != nil {
				return lerrors.Wrap(err, "failed to convert bech32 delegator address from string")
			}

			updatedDel, err := sKeeper.GetDelegation(ctx, sdk.AccAddress(delAddr), valAddr)
			if err != nil {
				return lerrors.Wrap(err, "failed to verify updated delegation")
			}

			if !updatedDel.RewardsShares.Equal(newDelRewardsShares) {
				return lerrors.New(fmt.Sprintf("delegation rewards_shares fix verification failed: expected=%s got=%s",
					newDelRewardsShares.String(), updatedDel.RewardsShares.String()),
				)
			}

			// update period delegation
			periodDel, err := sKeeper.GetPeriodDelegation(ctx, sdk.AccAddress(delAddr), valAddr, stypes.FlexiblePeriodDelegationID) // only flexible period delegation is allowed for locked validator
			if err != nil {
				return lerrors.Wrap(err, "failed to get period delegation",
					"validator_address", valAddr.String(),
					"delegator_address", sdk.AccAddress(delAddr).String(),
				)
			}

			periodDel.RewardsShares = newDelRewardsShares

			if err := sKeeper.SetPeriodDelegation(ctx, sdk.AccAddress(delAddr), valAddr, periodDel); err != nil {
				return lerrors.Wrap(err, "failed to set updated period delegation",
					"validator_address", valAddr.String(),
					"delegator_address", sdk.AccAddress(delAddr).String(),
				)
			}

			updatedPeriodDel, err := sKeeper.GetPeriodDelegation(ctx, sdk.AccAddress(delAddr), valAddr, stypes.FlexiblePeriodDelegationID)
			if err != nil {
				return lerrors.Wrap(err, "failed to get updated period delegation",
					"validator_address", valAddr.String(),
					"delegator_address", sdk.AccAddress(delAddr).String(),
				)
			}

			if !updatedPeriodDel.RewardsShares.Equal(newDelRewardsShares) {
				return lerrors.New(fmt.Sprintf("period delegation rewards_shares fix verification failed: expected=%s got=%s",
					newDelRewardsShares.String(), updatedPeriodDel.RewardsShares.String()),
				)
			}
		}

		val.DelegatorRewardsShares = newTotalRewardsShares
		val.RewardsTokens = newRewardsTokens

		// set updated validator
		if err := sKeeper.SetValidator(ctx, val); err != nil {
			return lerrors.Wrap(err, "update validator rewards tokens", "validator", val.OperatorAddress)
		}

		// validation
		newVal, err := sKeeper.GetValidator(ctx, valAddr)
		if err != nil {
			return lerrors.Wrap(err, "failed to get validator")
		}

		if !newVal.DelegatorRewardsShares.Equal(newTotalRewardsShares) {
			return lerrors.New("sum of delegators' rewards shares does not equal to the expected one",
				"validator_address", valAddr.String(),
				"expected", newTotalRewardsShares,
				"actual", newVal.DelegatorRewardsShares,
			)
		}
		if !newVal.RewardsTokens.Equal(newRewardsTokens) {
			return lerrors.New("rewards tokens of the validator does not equal to the expected one",
				"validator_address", valAddr.String(),
				"expected", newRewardsTokens,
				"actual", newVal.RewardsTokens,
			)
		}

		// scale delegator starting info

		// get all updated delegations as it should be rescaled based on the updated delegations
		updatedDels, err := sKeeper.GetValidatorDelegations(ctx, valAddr)
		if err != nil {
			return lerrors.Wrap(err, "failed to get all updated delegations for the validator")
		}

		for _, del := range updatedDels {
			delAddr, err := aKeeper.AddressCodec().StringToBytes(del.DelegatorAddress)
			if err != nil {
				return lerrors.Wrap(err, "failed to convert bech32 delegator address from string")
			}

			startingInfo, err := dKeeper.GetDelegatorStartingInfo(ctx, valAddr, sdk.AccAddress(delAddr))
			if err != nil {
				return lerrors.Wrap(err, "failed to get delegator starting info")
			}

			newRewardsStake := val.RewardsTokensFromRewardsSharesTruncated(del.RewardsShares)
			startingInfo.RewardsStake = newRewardsStake
			if err := dKeeper.SetDelegatorStartingInfo(ctx, valAddr, sdk.AccAddress(delAddr), startingInfo); err != nil {
				return lerrors.Wrap(err, "failed to set delegator starting info")
			}

			scaledDelegatorStartingInfo, err := dKeeper.GetDelegatorStartingInfo(ctx, valAddr, sdk.AccAddress(delAddr))
			if err != nil {
				return lerrors.Wrap(err, "failed to get scaled delegator starting info")
			}

			if !scaledDelegatorStartingInfo.RewardsStake.Equal(newRewardsStake) {
				return lerrors.New("delegator rewards stake mismatch",
					"validator_address", valAddr.String(),
					"delegator_address", sdk.AccAddress(delAddr).String(),
					"expected_rewards_stake", newRewardsStake.String(),
					"actual_rewards_stake", scaledDelegatorStartingInfo.RewardsStake.String(),
				)
			}
		}
	}

	// check that each unlocked validator's rewards tokens are not changed
	newValidators, err := sKeeper.GetAllValidators(ctx)
	if err != nil {
		return lerrors.Wrap(err, "get all validators")
	}

	newTotalRewardsTokens := math.LegacyZeroDec()
	newTotalRewardsTokensOfLockedValidators := math.LegacyZeroDec()
	newTotalRewardsTokensOfUnlockedValidators := math.LegacyZeroDec()
	for _, val := range newValidators {
		rt := val.GetRewardsTokens()
		newTotalRewardsTokens = newTotalRewardsTokens.Add(rt)
		if val.SupportTokenType == lockedTokenType {
			newTotalRewardsTokensOfLockedValidators = newTotalRewardsTokensOfLockedValidators.Add(rt)
		} else {
			newTotalRewardsTokensOfUnlockedValidators = newTotalRewardsTokensOfUnlockedValidators.Add(rt)
		}
	}
	if newTotalRewardsTokens.IsZero() {
		return lerrors.New(
			"new total rewards tokens is zero",
		)
	}
	if !newTotalRewardsTokensOfLockedValidators.Add(newTotalRewardsTokensOfUnlockedValidators).Equal(newTotalRewardsTokens) {
		return lerrors.New(
			"sum of new locked and unlocked TotalRewardsTokens does not equal new total rewards tokens",
			"new total", newTotalRewardsTokens.String(),
			"new sum", newTotalRewardsTokensOfLockedValidators.Add(newTotalRewardsTokensOfUnlockedValidators).String(),
		)
	}
	if !newTotalRewardsTokensOfUnlockedValidators.Equal(oldTotalRewardsTokensOfUnlockedValidators) {
		return lerrors.New(
			"new unlocked TotalRewardsTokens does not equal old unlocked TotalRewardsTokens",
			"old", oldTotalRewardsTokensOfUnlockedValidators.String(),
			"new", newTotalRewardsTokensOfUnlockedValidators.String(),
		)
	}

	log.Info(ctx, "Completed updating rewards shares for locked delegations...")

	// --------------------------------
	// update mint params
	// --------------------------------

	log.Info(ctx, "Start updating mint params...")

	mintParams, err := mKeeper.GetParams(ctx)
	if err != nil {
		return lerrors.Wrap(err, "get mint params")
	}

	mintParams.InflationsPerYear = NewAnnualInflationsPerYear
	mintParams.BlocksPerYear = NewBlocksPerYear
	if err := mintParams.Validate(); err != nil {
		return lerrors.Wrap(err, "validate mint params")
	}

	if err := mKeeper.SetParams(ctx, mintParams); err != nil {
		return lerrors.Wrap(err, "update mint params")
	}

	// check that the annual inflations per year is updated
	mintParams, err = mKeeper.GetParams(ctx)
	if err != nil {
		return lerrors.Wrap(err, "reload mint params")
	}
	if !mintParams.InflationsPerYear.Equal(NewAnnualInflationsPerYear) {
		return lerrors.New(
			"inflations_per_year not updated",
			"expected", NewAnnualInflationsPerYear,
			"actual", mintParams.InflationsPerYear,
		)
	}
	if mintParams.BlocksPerYear != NewBlocksPerYear {
		return lerrors.New(
			"blocks_per_year not updated",
			"expected", NewBlocksPerYear,
			"actual", mintParams.BlocksPerYear,
		)
	}

	log.Info(ctx, "Completed updating mint params...")

	log.Info(ctx, "Horace upgrade complete")

	return nil
}
