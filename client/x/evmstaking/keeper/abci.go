package keeper

import (
	"context"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/promutil"
)

// Query staking module's UnbondingDelegation (UBD Queue) to get the matured unbonding delegations. Then,
// insert the matured unbonding delegations into the withdrawal queue.
func (k *Keeper) EndBlock(ctx context.Context) (abci.ValidatorUpdates, error) {
	log.Debug(ctx, "EndBlock.evmstaking")
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

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
