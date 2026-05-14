package keeper

import (
	"context"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/upgrades/seneca"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/netconf"
	"github.com/piplabs/story/lib/promutil"
)

// EndBlock query staking module's UnbondingDelegation (UBD Queue) to get the matured unbonding delegations. Then,
// insert the matured unbonding delegations into the withdrawal queue.
func (k *Keeper) EndBlock(ctx context.Context) (abci.ValidatorUpdates, error) {
	log.Debug(ctx, "EndBlock.evmstaking")

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Apply deferred MaxValidators reduction BEFORE singularity check and
	// staking EndBlocker. This mirrors gov.EndBlocker → staking.EndBlocker
	// ordering in standard Cosmos SDK.
	if err := k.applyDeferredMaxValidatorsChange(ctx); err != nil {
		return nil, errors.Wrap(err, "apply deferred max validators change")
	}

	isSingularity, err := k.IsSingularity(ctx)
	if err != nil {
		return nil, err
	}

	if isSingularity {
		return nil, nil
	}

	valUpdates, unbondedEntries, err := k.stakingKeeper.EndBlockerWithUnbondedEntries(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "process staking EndBlocker")
	}

	if err := k.ProcessUnstakeWithdrawals(ctx, unbondedEntries); err != nil {
		return nil, errors.Wrap(err, "process unstake withdrawals")
	}

	if err := k.ProcessRewardWithdrawals(ctx); err != nil {
		return nil, errors.Wrap(err, "process reward withdrawals")
	}

	if err := k.ProcessUbiWithdrawal(ctx); err != nil {
		return nil, errors.Wrap(err, "process ubi withdrawal")
	}

	// set metrics
	promutil.EVMStakingWithdrawalQueueDepth.Set(float64(k.WithdrawalQueue.Len(ctx)))
	promutil.EVMStakingRewardQueueDepth.Set(float64(k.RewardWithdrawalQueue.Len(ctx)))

	return valUpdates, nil
}

// applyDeferredMaxValidatorsChange sets MaxValidators at the exact seneca
// upgrade height. This is a no-op for chains that have not registered Seneca
// in their upgrade history.
func (k *Keeper) applyDeferredMaxValidatorsChange(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	senecaHeight, err := netconf.GetUpgradeHeight(sdkCtx.ChainID(), netconf.Seneca)
	if err != nil {
		return nil // Seneca not registered for this chain — skip
	}

	if sdkCtx.BlockHeight() != senecaHeight {
		return nil
	}

	params, err := k.stakingKeeper.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "get staking params for max validators change")
	}

	if params.MaxValidators <= seneca.NewMaxValidators {
		return nil
	}

	params.MaxValidators = seneca.NewMaxValidators
	if err := k.stakingKeeper.SetParams(ctx, params); err != nil {
		return errors.Wrap(err, "set staking params for max validators change")
	}

	log.Info(ctx, "Applied deferred MaxValidators reduction", "max_validators", seneca.NewMaxValidators)

	return nil
}
