package keeper

import (
	"context"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	estypes "github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/promutil"
)

func (k Keeper) ExpectedPartialWithdrawals(ctx context.Context) ([]estypes.Withdrawal, error) {
	// TODO: user more fine-grained cursor with next delegator sweep index.
	nextValSweepIndex, err := k.GetNextValidatorSweepIndex(ctx)
	if err != nil {
		return nil, err
	}
	nextValIndex := nextValSweepIndex.Int.Int64()
	// Get all validators first, and then do a circular sweep
	validatorSet, err := (k.stakingKeeper.(*skeeper.Keeper)).GetAllValidators(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get all validators")
	}

	if nextValIndex >= int64(len(validatorSet)) {
		// TODO: TBD
		log.Warn(
			ctx, "NextValidatorIndex exceeds the validator set size",
			errors.New("nextValidatorIndex overflow"),
			"validator_set", len(validatorSet),
			"next_validator_index", nextValIndex,
		)
		nextValIndex = 0
	}

	// Iterate all validators from `nextValidatorIndex` to find out eligible partial withdrawals.
	var (
		swept       uint32
		withdrawals []estypes.Withdrawal
	)
	// Get sweep limit per block.
	sweepBound, err := k.MaxSweepPerBlock(ctx)
	if err != nil {
		return nil, err
	}
	// Get minimal partial withdrawal amount.
	minPartialWithdrawalAmount, err := k.MinPartialWithdrawalAmount(ctx)
	if err != nil {
		return nil, err
	}
	log.Debug(
		ctx, "partial withdrawal params",
		"min_partial_withdraw_amount", minPartialWithdrawalAmount,
		"max_sweep_per_block", sweepBound,
	)
	// Sweep and get eligible partial withdrawals.
	for range validatorSet {
		if swept > sweepBound {
			break
		}
		if validatorSet[nextValIndex].IsJailed() {
			continue
		}
		// Get validator's address.
		valBz, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(validatorSet[nextValIndex].GetOperator())
		if err != nil {
			return nil, errors.Wrap(err, "validator address from bech32")
		}
		valAddr := sdk.ValAddress(valBz)
		valAccAddr := sdk.AccAddress(valAddr)
		// Get validator commissions.
		valCommission, err := k.distributionKeeper.GetValidatorAccumulatedCommission(ctx, valAddr)
		if err != nil {
			return nil, err
		}
		log.Debug(
			ctx, "Get validator commission",
			"val_addr", valAddr.String(),
			"commission_amount", valCommission.Commission.String(),
		)
		// Get all delegators of the validator.
		delegators, err := (k.stakingKeeper.(*skeeper.Keeper)).GetValidatorDelegations(ctx, valAddr)
		if err != nil {
			return nil, errors.Wrap(err, "get validator delegations")
		}
		swept += uint32(len(delegators))
		log.Debug(
			ctx, "Get all delegators of validator",
			"val_addr", valAddr.String(),
			"delegator_amount", len(delegators),
		)
		// Get delegator rewards.
		for i := range delegators {
			// Get end current period and calculate rewards.
			endingPeriod, err := k.distributionKeeper.IncrementValidatorPeriod(ctx, validatorSet[nextValIndex])
			if err != nil {
				return nil, err
			}
			delRewards, err := k.distributionKeeper.CalculateDelegationRewards(ctx, validatorSet[nextValIndex], delegators[i], endingPeriod)
			if err != nil {
				return nil, err
			}
			if delegators[i].DelegatorAddress == valAccAddr.String() {
				delRewards = delRewards.Add(valCommission.Commission...)
			}
			delRewardsTruncated, _ := delRewards.TruncateDecimal()
			bondDenomAmount := delRewardsTruncated.AmountOf(sdk.DefaultBondDenom).Uint64()

			log.Debug(
				ctx, "Calculate delegator rewards",
				"val_addr", valAddr.String(),
				"del_addr", delegators[i].DelegatorAddress,
				"rewards_amount", bondDenomAmount,
				"ending_period", endingPeriod,
			)

			if bondDenomAmount >= minPartialWithdrawalAmount {
				delEvmAddr, err := k.DelegatorMap.Get(ctx, delegators[i].DelegatorAddress)
				if err != nil {
					return nil, errors.Wrap(err, "map delegator pubkey to evm address")
				}
				withdrawals = append(withdrawals, estypes.NewWithdrawal(
					uint64(sdk.UnwrapSDKContext(ctx).BlockHeight()),
					delegators[i].DelegatorAddress,
					valAddr.String(),
					delEvmAddr,
					bondDenomAmount,
				))

				log.Debug(
					ctx, "Found an eligible partial withdrawal",
					"val_addr", valAddr.String(),
					"del_addr", delegators[i].DelegatorAddress,
					"del_evm_addr", delEvmAddr,
					"rewards_amount", bondDenomAmount,
				)
			}
		}
		nextValIndex = (nextValIndex + 1) % int64(len(validatorSet))
	}
	// Update the nextValidatorSweepIndex.
	if err := k.SetNextValidatorSweepIndex(
		ctx,
		sdk.IntProto{Int: math.NewInt(nextValIndex)},
	); err != nil {
		return nil, err
	}
	log.Debug(
		ctx, "Finish validator sweep for partial withdrawals",
		"next_validator_index", nextValIndex,
		"partial_withdrawals", len(withdrawals),
	)

	return withdrawals, nil
}

func (k Keeper) EnqueueEligiblePartialWithdrawal(ctx context.Context, withdrawals []estypes.Withdrawal) error {
	for i := range withdrawals {
		valAddr, err := sdk.ValAddressFromBech32(withdrawals[i].ValidatorAddress)
		if err != nil {
			return errors.Wrap(err, "validator address from bech32")
		}

		valAccAddr := sdk.AccAddress(valAddr).String()

		// Withdraw delegation rewards.
		delAddr := sdk.MustAccAddressFromBech32(withdrawals[i].DelegatorAddress)
		delRewards, err := k.distributionKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
		if err != nil {
			return err
		}

		// Withdraw commission if it is a self delegation.
		if withdrawals[i].DelegatorAddress == valAccAddr {
			commissionRewards, err := k.distributionKeeper.WithdrawValidatorCommission(ctx, valAddr)
			if errors.Is(err, dtypes.ErrNoValidatorCommission) {
				log.Debug(
					ctx, "No validator commission",
					"validator_addr", withdrawals[i].ValidatorAddress,
					"validator_account_addr", valAccAddr,
				)
			} else if err != nil {
				return err
			} else {
				delRewards = delRewards.Add(commissionRewards...)
			}
		}
		curBondDenomAmount := delRewards.AmountOf(sdk.DefaultBondDenom).Uint64()
		log.Debug(
			ctx, "Withdraw delegator rewards",
			"validator_addr", withdrawals[i].ValidatorAddress,
			"validator_account_addr", valAccAddr,
			"delegator_addr", withdrawals[i].DelegatorAddress,
			"amount_calculate", withdrawals[i].Amount,
			"amount_withdraw", curBondDenomAmount,
		)

		// Burn tokens from the delegator
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, delAddr, estypes.ModuleName, delRewards)
		if err != nil {
			return err
		}
		err = k.bankKeeper.BurnCoins(ctx, estypes.ModuleName, delRewards)
		if err != nil {
			return err
		}

		// Enqueue to the global withdrawal queue.
		if err := k.AddWithdrawalToQueue(ctx, estypes.NewWithdrawal(
			uint64(sdk.UnwrapSDKContext(ctx).BlockHeight()),
			withdrawals[i].DelegatorAddress, withdrawals[i].ValidatorAddress, withdrawals[i].ExecutionAddress,
			curBondDenomAmount,
		)); err != nil {
			return err
		}
	}

	// set metrics
	promutil.EVMStakingQueueDepth.Set(float64(k.WithdrawalQueue.Len(ctx)))

	return nil
}

func (k Keeper) ProcessWithdraw(ctx context.Context, ev *bindings.IPTokenStakingWithdraw) error {
	depositorPubkey, err := k1util.PubKeyBytesToCosmos(ev.DepositorPubkey)
	if err != nil {
		return errors.Wrap(err, "depositor pubkey to cosmos")
	}

	validatorPubkey, err := k1util.PubKeyBytesToCosmos(ev.ValidatorPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to cosmos")
	}

	depositorAddr := sdk.AccAddress(depositorPubkey.Address().Bytes())
	validatorAddr := sdk.ValAddress(validatorPubkey.Address().Bytes())

	valEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(validatorPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "validator pubkey to evm address")
	}
	delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(depositorPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "validator pubkey to evm address")
	}

	amountCoin, _ := IPTokenToBondCoin(ev.Amount)

	log.Debug(ctx, "Processing EVM staking withdraw",
		"del_iliad", depositorAddr.String(),
		"val_iliad", validatorAddr.String(),
		"del_evm_addr", delEvmAddr.String(),
		"val_evm_addr", valEvmAddr.String(),
		"amount", ev.Amount.String(),
	)

	if !k.authKeeper.HasAccount(ctx, depositorAddr) {
		// TODO: gracefully handle when malicious or uninformed user tries to withdraw from non-existent account
		// skip errors.Wrap(err) since err will be nil (since all prev errors were nil to reach this branch)
		return errors.New("depositor account not found")
	}

	msg := stypes.NewMsgUndelegate(depositorAddr.String(), validatorAddr.String(), amountCoin)

	// Undelegate from the validator (validator existence is checked in ValidateUnbondAmount)
	resp, err := skeeper.NewMsgServerImpl(k.stakingKeeper.(*skeeper.Keeper)).Undelegate(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "undelegate")
	}

	log.Info(ctx, "EVM staking withdraw detected, undelegating from validator",
		"delegator", depositorAddr.String(),
		"validator", validatorAddr.String(),
		"amount", resp.Amount.String(),
		"completion_time", resp.CompletionTime)

	return nil
}

func (k Keeper) ParseWithdrawLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingWithdraw, error) {
	return k.ipTokenStakingContract.ParseWithdraw(ethlog)
}
