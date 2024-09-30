package keeper

import (
	"context"
	"math/big"
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
	ctxTime := sdkCtx.BlockHeader().Time
	blockHeight := sdkCtx.BlockHeader().Height

	matureUnbonds, err := k.GetMatureUnbondedDelegations(ctx)
	log.Debug(ctx, "Processing mature unbonding delegations", "count", len(matureUnbonds))
	if err != nil {
		return nil, err
	}

	// delegatorAddress -> validatorAddress -> unbondingId -> amount
	unbondedEntries := make(map[string]map[string]map[uint64]uint64)
	// Fetch all the mature unbonding entries before processed by staking keeper's EndBlocker.
	for _, dvPair := range matureUnbonds {
		delegatorAddr, err := k.authKeeper.AddressCodec().StringToBytes(dvPair.DelegatorAddress)
		if err != nil {
			return nil, errors.Wrap(err, "delegator address from bech32")
		}
		validatorAddr, err := k.validatorAddressCodec.StringToBytes(dvPair.ValidatorAddress)
		if err != nil {
			return nil, errors.Wrap(err, "validator address from bech32")
		}

		ubd, err := (k.stakingKeeper).GetUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
		if err != nil {
			return nil, err
		}

		// loop through all the entries and process unbonding mature entries
		for i := range len(ubd.Entries) {
			entry := ubd.Entries[i]
			// track undelegation only when remaining or truncated shares are non-zero
			if entry.IsMature(ctxTime) && !entry.OnHold() && !entry.Balance.IsZero() {
				if _, ok := unbondedEntries[dvPair.DelegatorAddress]; !ok {
					unbondedEntries[dvPair.DelegatorAddress] = make(map[string]map[uint64]uint64)
				}
				if _, ok := unbondedEntries[dvPair.DelegatorAddress][dvPair.ValidatorAddress]; !ok {
					unbondedEntries[dvPair.DelegatorAddress][dvPair.ValidatorAddress] = make(map[uint64]uint64)
				}
				unbondedEntries[dvPair.DelegatorAddress][dvPair.ValidatorAddress][entry.UnbondingId] = entry.Balance.Uint64()
			}
		}
	}

	valUpdates, err := k.stakingKeeper.EndBlocker(ctx)
	if err != nil {
		return nil, err
	}

	for delegatorAddr, validatorMap := range unbondedEntries {
		for validatorAddr, unbondingMap := range validatorMap {
			// Double check if there any unprocessed mature unbonding entries.
			delAddr, err := k.authKeeper.AddressCodec().StringToBytes(delegatorAddr)
			if err != nil {
				return nil, errors.Wrap(err, "delegator address from bech32")
			}
			valAddr, err := k.validatorAddressCodec.StringToBytes(validatorAddr)
			if err != nil {
				return nil, errors.Wrap(err, "validator address from bech32")
			}
			ubd, err := (k.stakingKeeper).GetUnbondingDelegation(ctx, delAddr, valAddr)
			if err != nil {
				return nil, err
			}
			for _, entry := range ubd.Entries {
				if entry.IsMature(ctxTime) && !entry.OnHold() {
					// Maps used here should already initialized in the previous loop.
					if _, ok := unbondedEntries[validatorAddr][delegatorAddr][entry.UnbondingId]; !entry.Balance.IsZero() && ok {
						log.Warn(ctx, "Incomplete mature unbonding entry", nil,
							"delegator", delegatorAddr,
							"validator", validatorAddr,
							"unbonding_id", entry.UnbondingId,
							"amount", entry.Balance.Uint64(),
						)
						// NOTE: We should delete the incomplete unbonding entry from the map to avoid duplicate processing later on.
						delete(unbondingMap, entry.UnbondingId)
					}
				}
			}
			// Process the complete mature unbonding entries.
			for _, amount := range unbondingMap {
				log.Debug(ctx, "Adding undelegation to withdrawal queue",
					"delegator", delegatorAddr,
					"validator", validatorAddr,
					"amount", amount,
				)
				// Burn tokens from the delegator
				_, coins := IPTokenToBondCoin(big.NewInt(int64(amount)))
				if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, delAddr, types.ModuleName, coins); err != nil {
					return nil, errors.Wrap(err, "send coins from account to module")
				}
				if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
					return nil, errors.Wrap(err, "burn coins")
				}

				// This should not produce error, as all delegations are done via the evmstaking module via EL.
				// However, we should gracefully handle in case Get fails.
				delEvmAddr, err := k.DelegatorMap.Get(ctx, delegatorAddr)
				if err != nil {
					return nil, errors.Wrap(err, "map delegator pubkey to evm address")
				}

				// push the undelegation to the withdrawal queue
				if err := k.AddWithdrawalToQueue(ctx, types.NewWithdrawal(
					uint64(blockHeight),
					delegatorAddr,
					validatorAddr,
					delEvmAddr,
					amount,
				)); err != nil {
					return nil, err
				}
			}
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
