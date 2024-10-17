package keeper

import (
	"context"

	"cosmossdk.io/collections"

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
	nextValIndex, nextValDelIndex, err := k.GetValidatorSweepIndex(ctx)
	if err != nil {
		return nil, err
	}

	// Get all validators first, and then do a circular sweep
	validatorSet, err := (k.stakingKeeper.(*skeeper.Keeper)).GetAllValidators(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get all validators")
	}

	if nextValIndex >= uint64(len(validatorSet)) {
		// TODO: TBD
		log.Warn(
			ctx, "NextValidatorIndex exceeds the validator set size",
			errors.New("nextValidatorIndex overflow"),
			"validator_set", len(validatorSet),
			"next_validator_index", nextValIndex,
		)
		nextValIndex = 0
		nextValDelIndex = 0
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

	// Sweep and get eligible partial withdrawals.
	for range validatorSet {
		if validatorSet[nextValIndex].IsJailed() {
			// nextValIndex should be updated, even if the validator is jailed, to progress to the sweep.
			nextValIndex = (nextValIndex + 1) % uint64(len(validatorSet))
			nextValDelIndex = 0

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

		// Get all delegators of the validator.
		delegators, err := (k.stakingKeeper.(*skeeper.Keeper)).GetValidatorDelegations(ctx, valAddr)
		if err != nil {
			return nil, errors.Wrap(err, "get validator delegations")
		}

		if nextValDelIndex >= uint64(len(delegators)) {
			nextValIndex = (nextValIndex + 1) % uint64(len(validatorSet))
			nextValDelIndex = 0

			continue
		}

		nextDelegators := delegators[nextValDelIndex:]
		var shouldStopPrematurely bool

		// Check if the sweep should stop prematurely as the current delegator loop exceeds the sweep bound while sweeping.
		remainingSweep := sweepBound - swept
		if uint32(len(nextDelegators)) > remainingSweep {
			nextDelegators = nextDelegators[:remainingSweep]
			shouldStopPrematurely = true
		}

		// Iterate on the validator's delegator rewards in the range [nextValDelIndex, len(delegators)].
		for i := range nextDelegators {
			// Get end current period and calculate rewards.
			endingPeriod, err := k.distributionKeeper.IncrementValidatorPeriod(ctx, validatorSet[nextValIndex])
			if err != nil {
				return nil, err
			}

			delRewards, err := k.distributionKeeper.CalculateDelegationRewards(ctx, validatorSet[nextValIndex], nextDelegators[i], endingPeriod)
			if err != nil {
				return nil, err
			}

			if nextDelegators[i].DelegatorAddress == valAccAddr.String() {
				delRewards = delRewards.Add(valCommission.Commission...)
			}
			delRewardsTruncated, _ := delRewards.TruncateDecimal()
			bondDenomAmount := delRewardsTruncated.AmountOf(sdk.DefaultBondDenom).Uint64()

			if bondDenomAmount >= minPartialWithdrawalAmount {
				delEvmAddr, err := k.DelegatorRewardAddress.Get(ctx, nextDelegators[i].DelegatorAddress)
				if err != nil {
					return nil, errors.Wrap(err, "map delegator pubkey to evm reward address")
				}

				withdrawals = append(withdrawals, estypes.NewWithdrawal(
					uint64(sdk.UnwrapSDKContext(ctx).BlockHeight()),
					nextDelegators[i].DelegatorAddress,
					valAddr.String(),
					delEvmAddr,
					bondDenomAmount,
				))
			}

			nextValDelIndex++
		}

		// If the validator's delegation loop was stopped prematurely, we break from the validator sweep loop.
		if shouldStopPrematurely {
			break
		}

		// Here, we have looped through all delegators of the validator (since we did not prematurely stop in the loop above).
		// Thus, we signal to progress to the next validator by resetting the nextValDelIndex and circularly incrementing the nextValIndex
		nextValIndex = (nextValIndex + 1) % uint64(len(validatorSet))
		nextValDelIndex = 0

		// Increase the total swept amount.
		swept += uint32(len(nextDelegators))
	}

	// Update the validator sweep index.
	if err := k.SetValidatorSweepIndex(
		ctx,
		nextValIndex,
		nextValDelIndex,
	); err != nil {
		return nil, err
	}

	log.Debug(
		ctx, "Finish validator sweep for partial withdrawals",
		"next_validator_index", nextValIndex,
		"next_validator_delegation_index", nextValDelIndex,
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
	isInSingularity, err := k.IsSingularity(ctx)
	if err != nil {
		return errors.Wrap(err, "check if it is singularity")
	}

	if isInSingularity {
		log.Info(ctx, "Withdraw event detected, but it is not processed since current block is singularity")
		return nil
	}

	delCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.DelegatorUncmpPubkey)
	if err != nil {
		return errors.Wrap(err, "compress depositor pubkey")
	}
	depositorPubkey, err := k1util.PubKeyBytesToCosmos(delCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "depositor pubkey to cosmos")
	}

	valCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.ValidatorUnCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "compress validator pubkey")
	}
	validatorPubkey, err := k1util.PubKeyBytesToCosmos(valCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to cosmos")
	}

	depositorAddr := sdk.AccAddress(depositorPubkey.Address().Bytes())
	validatorAddr := sdk.ValAddress(validatorPubkey.Address().Bytes())

	valEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(valCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to evm address")
	}
	delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(delCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "delegator pubkey to evm address")
	}

	// unstakeOnBehalf txn, need to check if it's from the operator
	if delEvmAddr.String() != ev.OperatorAddress.String() {
		operatorAddr, err := k.DelegatorOperatorAddress.Get(ctx, depositorAddr.String())
		if errors.Is(err, collections.ErrNotFound) {
			return errors.New("invalid unstakeOnBehalf txn, no operator for delegator")
		} else if err != nil {
			return errors.Wrap(err, "get delegator's operator address failed")
		}
		if operatorAddr != ev.OperatorAddress.String() {
			return errors.New("invalid unstakeOnBehalf txn, not from operator")
		}
	}

	amountCoin, _ := IPTokenToBondCoin(ev.StakeAmount)

	log.Debug(ctx, "Processing EVM staking withdraw",
		"del_story", depositorAddr.String(),
		"val_story", validatorAddr.String(),
		"del_evm_addr", delEvmAddr.String(),
		"val_evm_addr", valEvmAddr.String(),
		"amount", ev.StakeAmount.String(),
	)

	if !k.authKeeper.HasAccount(ctx, depositorAddr) {
		// TODO: gracefully handle when malicious or uninformed user tries to withdraw from non-existent account
		// skip errors.Wrap(err) since err will be nil (since all prev errors were nil to reach this branch)
		return errors.New("depositor account not found")
	}

	msg := stypes.NewMsgUndelegate(depositorAddr.String(), validatorAddr.String(), ev.DelegationId.String(), amountCoin)

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
