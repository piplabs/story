package keeper

import (
	"context"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

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

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeader().Height

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

	log.Debug(ctx, "Processing mature unbonding delegations", "count", len(unbondedEntries))

	for _, entry := range unbondedEntries {
		delegatorAddr, err := k.authKeeper.AddressCodec().StringToBytes(entry.DelegatorAddress)
		if err != nil {
			return nil, errors.Wrap(err, "delegator address from bech32")
		}

		log.Debug(ctx, "Adding undelegation to withdrawal queue",
			"delegator", entry.DelegatorAddress,
			"validator", entry.ValidatorAddress,
			"amount", entry.Amount.String())

		// Burn tokens from the delegator
		_, coins := IPTokenToBondCoin(entry.Amount.BigInt())
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, delegatorAddr, types.ModuleName, coins)
		if err != nil {
			return nil, errors.Wrap(err, "send coins from account to module")
		}
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
		if err != nil {
			return nil, errors.Wrap(err, "burn coins")
		}

		// This should not produce error, as all delegations are done via the evmstaking module via EL.
		// However, we should gracefully handle in case Get fails.
		delEvmAddr, err := k.DelegatorWithdrawAddress.Get(ctx, entry.DelegatorAddress)
		if err != nil {
			return nil, errors.Wrap(err, "map delegator pubkey to evm address")
		}

		// push the undelegation to the withdrawal queue
		err = k.AddWithdrawalToQueue(ctx, types.NewWithdrawal(
			uint64(blockHeight),
			entry.DelegatorAddress,
			entry.ValidatorAddress,
			delEvmAddr,
			entry.Amount.Uint64(),
		))
		if err != nil {
			return nil, err
		}
	}

	partialWithdrawals, err := k.ExpectedPartialWithdrawals(ctx)
	if err != nil {
		return nil, err
	}
	if err := k.EnqueueEligiblePartialWithdrawal(ctx, partialWithdrawals); err != nil {
		return nil, err
	}

	return valUpdates, nil
}
