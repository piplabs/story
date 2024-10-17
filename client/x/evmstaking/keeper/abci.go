package keeper

import (
	"context"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// Query staking module's UnbondingDelegation (UBD Queue) to get the matured unbonding delegations. Then,
// insert the matured unbonding delegations into the withdrawal queue.
// TODO: check if unbonded delegations in staking module must be distinguished based on source of generation, CL or EL.
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
		return nil, err
	}
	if err := k.ProcessUnbondingWithdrawals(ctx, unbondedEntries); err != nil {
		return nil, err
	}

	partialWithdrawals, err := k.ExpectedPartialWithdrawals(ctx)
	if err != nil {
		return nil, err
	}
	if err := k.EnqueueEligiblePartialWithdrawal(ctx, partialWithdrawals); err != nil {
		return nil, err
	}

	if err := k.ProcessUbiWithdrawal(ctx); err != nil {
		return nil, errors.Wrap(err, "process ubi withdrawal")
	}

	return valUpdates, nil
}
