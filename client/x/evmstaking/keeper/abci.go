package keeper

import (
	"context"
	"time"

	"cosmossdk.io/math"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

type UnbondedEntry struct {
	validatorAddress string
	delegatorAddress string
	amount           math.Int
}

// Query staking module's UnbondingDelegation (UBD Queue) to get the matured unbonding delegations. Then,
// insert the matured unbonding delegations into the withdrawal queue.
// TODO: check if unbonded delegations in staking module must be distinguished based on source of generation, CL or EL.
func (k *Keeper) EndBlock(ctx context.Context) (abci.ValidatorUpdates, error) {
	log.Debug(ctx, "EndBlock.evmstaking")
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	ctxTime := sdkCtx.BlockHeader().Time
	blockHeight := sdkCtx.BlockHeader().Height

	matureUnbonds, err := k.GetMatureUnbondedDelegations(ctx)
	log.Debug(ctx, "Processing mature unbonding delegations", "count", len(matureUnbonds))
	if err != nil {
		return nil, err
	}

	// make an array with each entry being the validator address, delegator address, and the amount
	var unbondedEntries []UnbondedEntry

	for _, dvPair := range matureUnbonds {
		validatorAddr, err := k.validatorAddressCodec.StringToBytes(dvPair.ValidatorAddress)
		if err != nil {
			return nil, errors.Wrap(err, "validator address from bech32")
		}

		delegatorAddr, err := k.authKeeper.AddressCodec().StringToBytes(dvPair.DelegatorAddress)
		if err != nil {
			return nil, errors.Wrap(err, "delegator address from bech32")
		}

		ubd, err := (k.stakingKeeper).GetUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
		if err != nil {
			return nil, err
		}

		// TODO: parameterized bondDenom
		bondDenom := sdk.DefaultBondDenom

		// loop through all the entries and process unbonding mature entries
		for i := range len(ubd.Entries) {
			entry := ubd.Entries[i]
			if entry.IsMature(ctxTime) && !entry.OnHold() {
				// track undelegation only when remaining or truncated shares are non-zero
				if !entry.Balance.IsZero() {
					amt := sdk.NewCoin(bondDenom, entry.Balance)
					// TODO: check if it's possible to add a double entry in the unbondedEntries array
					unbondedEntries = append(unbondedEntries, UnbondedEntry{
						validatorAddress: dvPair.ValidatorAddress,
						delegatorAddress: dvPair.DelegatorAddress,
						amount:           amt.Amount,
					})
				}
			}
		}
	}

	valUpdates, err := k.stakingKeeper.EndBlocker(ctx)
	if err != nil {
		return nil, err
	}

	for _, entry := range unbondedEntries {
		log.Debug(ctx, "Adding undelegation to withdrawal queue",
			"delegator", entry.delegatorAddress,
			"validator", entry.validatorAddress,
			"amount", entry.amount.String())

		delegatorAddr, err := k.authKeeper.AddressCodec().StringToBytes(entry.delegatorAddress)
		if err != nil {
			return nil, errors.Wrap(err, "delegator address from bech32")
		}
		// Burn tokens from the delegator
		_, coins := IPTokenToBondCoin(entry.amount.BigInt())
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
		delEvmAddr, err := k.DelegatorMap.Get(ctx, entry.delegatorAddress)
		if err != nil {
			return nil, errors.Wrap(err, "map delegator pubkey to evm address")
		}

		// push the undelegation to the withdrawal queue
		err = k.AddWithdrawalToQueue(ctx, types.NewWithdrawal(
			uint64(blockHeight),
			entry.delegatorAddress,
			entry.validatorAddress,
			delEvmAddr,
			entry.amount.Uint64(),
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
