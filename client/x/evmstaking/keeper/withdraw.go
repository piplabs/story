//nolint:contextcheck // use cached context
package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
)

func (k Keeper) ProcessUnstakeWithdrawals(ctx context.Context, unbondedEntries []stypes.UnbondedEntry) error {
	log.Debug(ctx, "Processing mature unbonding delegations", "count", len(unbondedEntries))

	for _, entry := range unbondedEntries {
		log.Debug(
			ctx, "Unstake withdrawal of mature unbonding delegation",
			"delegator", entry.DelegatorAddress,
			"validator", entry.ValidatorAddress,
			"amount", entry.Amount.String(),
		)

		// Check if the delegation is total unstaked
		delegatorAddr, err := sdk.AccAddressFromBech32(entry.DelegatorAddress)
		if err != nil {
			return errors.Wrap(err, "delegator address from bech32")
		}
		validatorAddr, err := sdk.ValAddressFromBech32(entry.ValidatorAddress)
		if err != nil {
			return errors.Wrap(err, "validator address from bech32")
		}
		var totallyUnstaked bool
		if _, err := k.stakingKeeper.GetDelegation(ctx, delegatorAddr, validatorAddr); err == nil {
			totallyUnstaked = false
		} else if errors.Is(err, stypes.ErrNoDelegation) {
			totallyUnstaked = true
		} else {
			return errors.Wrap(err, "get delegation")
		}

		// Withdraw commission if validator is totally self-unstaked
		if totallyUnstaked && bytes.Equal(delegatorAddr.Bytes(), validatorAddr.Bytes()) {
			if _, err := k.distributionKeeper.WithdrawValidatorCommission(ctx, validatorAddr); err != nil {
				return errors.Wrap(err, "withdraw validator commission")
			}
		}

		// Burn tokens from the delegator
		coinAmount := entry.Amount.Uint64()
		if totallyUnstaked {
			accountBalance := k.bankKeeper.SpendableCoin(ctx, delegatorAddr, sdk.DefaultBondDenom).Amount.Uint64()
			if accountBalance > coinAmount {
				coinAmount = accountBalance
			}
		}
		_, coins := IPTokenToBondCoin(big.NewInt(int64(coinAmount)))
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, delegatorAddr, types.ModuleName, coins); err != nil {
			return errors.Wrap(err, "send coins from account to module")
		}
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
			return errors.Wrap(err, "burn coins")
		}

		// This should not produce error, as all delegations are done via the evmstaking module via EL.
		// However, we should gracefully handle in case Get fails.
		delEvmAddr, err := k.DelegatorWithdrawAddress.Get(ctx, entry.DelegatorAddress)
		if err != nil {
			return errors.Wrap(err, "map delegator pubkey to evm address")
		}

		// push the undelegation to the withdrawal queue
		if err := k.AddWithdrawalToQueue(ctx, types.NewWithdrawal(
			uint64(sdk.UnwrapSDKContext(ctx).BlockHeight()),
			delEvmAddr,
			entry.Amount.Uint64(),
			types.WithdrawalType_WITHDRAWAL_TYPE_UNSTAKE,
		)); err != nil {
			return errors.Wrap(err, "add unstake withdrawal to queue")
		}

		if coinAmount <= entry.Amount.Uint64() {
			// No residue rewards, skip
			continue
		}

		residueAmount := coinAmount - entry.Amount.Uint64()

		log.Debug(ctx, "Residue rewards of mature unbonding delegation",
			"delegator", entry.DelegatorAddress,
			"validator", entry.ValidatorAddress,
			"amount", residueAmount,
		)

		// Enqueue to the global reward withdrawal queue.
		rewardsEVMAddr, err := k.DelegatorRewardAddress.Get(ctx, entry.DelegatorAddress)
		if err != nil {
			return errors.Wrap(err, "map delegator bech32 address to evm reward address")
		}
		if err := k.AddRewardWithdrawalToQueue(ctx, types.NewWithdrawal(
			uint64(sdk.UnwrapSDKContext(ctx).BlockHeight()),
			rewardsEVMAddr,
			residueAmount,
			types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
		)); err != nil {
			return errors.Wrap(err, "add reward withdrawal to queue")
		}
	}

	return nil
}

func (k Keeper) ProcessRewardWithdrawals(ctx context.Context) error {
	log.Debug(ctx, "Processing reward withdrawals")

	validatorSweepIndex, err := k.GetValidatorSweepIndex(ctx)
	if err != nil {
		return errors.Wrap(err, "get validator sweep index")
	}

	nextValIndex, nextValDelIndex := validatorSweepIndex.NextValIndex, validatorSweepIndex.NextValDelIndex

	// Get all validators first, and then do a circular sweep
	validatorSet, err := k.stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return errors.Wrap(err, "get all validators")
	}

	if nextValIndex >= uint64(len(validatorSet)) {
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
	var swept uint32

	// Get sweep limit per block.
	sweepBound, err := k.MaxSweepPerBlock(ctx)
	if err != nil {
		return errors.Wrap(err, "get max sweep per block")
	}

	// Get minimal reward withdrawal amount.
	minRewardWithdrawalAmount, err := k.MinPartialWithdrawalAmount(ctx)
	if err != nil {
		return errors.Wrap(err, "get minimum partial withdrawal amount")
	}

	// Sweep and get eligible partial withdrawals.
	for range validatorSet {
		if validatorSet[nextValIndex].IsJailed() {
			// nextValIndex should be updated, even if the validator is jailed, to progress to the sweep.
			nextValIndex = (nextValIndex + 1) % uint64(len(validatorSet))
			nextValDelIndex = 0

			continue
		}

		valAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(validatorSet[nextValIndex].GetOperator())
		if err != nil {
			return errors.Wrap(err, "convert validator address from string to bytes")
		}

		// Get all delegators of the validator.
		delegations, err := k.stakingKeeper.GetValidatorDelegations(ctx, sdk.ValAddress(valAddr))
		if err != nil {
			return errors.Wrap(err, "get validator delegations")
		}

		if nextValDelIndex >= uint64(len(delegations)) {
			nextValIndex = (nextValIndex + 1) % uint64(len(validatorSet))
			nextValDelIndex = 0

			continue
		}

		nextDelegations := delegations[nextValDelIndex:]
		var shouldStopPrematurely bool

		// Check if the sweep should stop prematurely as the current delegator loop exceeds the sweep bound while sweeping.
		remainingSweep := sweepBound - swept
		if uint32(len(nextDelegations)) > remainingSweep {
			nextDelegations = nextDelegations[:remainingSweep]
			shouldStopPrematurely = true
		}

		// Iterate on the validator's delegator rewards in the range [nextValDelIndex, len(delegators)].
		for _, delegation := range nextDelegations {
			if err := k.ProcessEligibleRewardWithdrawal(ctx, delegation, validatorSet[nextValIndex], minRewardWithdrawalAmount); err != nil {
				return errors.Wrap(err, "process eligible reward withdrawal")
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
		swept += uint32(len(nextDelegations))
	}

	// Update the validator sweep index.
	if err := k.SetValidatorSweepIndex(ctx, types.NewValidatorSweepIndex(nextValIndex, nextValDelIndex)); err != nil {
		return errors.Wrap(err, "set validator sweep index")
	}

	log.Debug(
		ctx, "Finish validator sweep for partial withdrawals",
		"next_validator_index", nextValIndex,
		"next_validator_delegation_index", nextValDelIndex,
	)

	return nil
}

// ProcessEligibleRewardWithdrawal processes the reward withdrawal of delegation.
func (k Keeper) ProcessEligibleRewardWithdrawal(ctx context.Context, delegation stypes.Delegation, validator stypes.Validator, minRewardWithdrawalAmount uint64) error {
	// Get validator's address.
	valBz, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(validator.GetOperator())
	if err != nil {
		return errors.Wrap(err, "validator address from bech32")
	}
	valAddr := sdk.ValAddress(valBz)
	valAccAddr := sdk.AccAddress(valAddr)

	// Get end current period and calculate rewards.
	endingPeriod, err := k.distributionKeeper.IncrementValidatorPeriod(ctx, validator)
	if err != nil {
		return err
	}

	delRewards, err := k.distributionKeeper.CalculateDelegationRewards(ctx, validator, delegation, endingPeriod)
	if err != nil {
		return err
	}

	// if it is self-delegation, add commission
	if delegation.DelegatorAddress == valAccAddr.String() {
		// Get validator commissions.
		valCommission, err := k.distributionKeeper.GetValidatorAccumulatedCommission(ctx, valAddr)
		if err != nil {
			return err
		}

		delRewards = delRewards.Add(valCommission.Commission...)
	}

	addr, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		return errors.Wrap(err, "convert acc address from bech32 address")
	}

	delRewardsTruncated, _ := delRewards.TruncateDecimal()
	unclaimedReward := delRewardsTruncated.AmountOf(sdk.DefaultBondDenom).Uint64()
	claimedReward := k.bankKeeper.SpendableCoin(ctx, addr, sdk.DefaultBondDenom).Amount.Uint64()
	totalReward := unclaimedReward + claimedReward

	// if total reward is greater than or equal to min reward withdrawal amount, enqueue the reward withdrawal
	if totalReward >= minRewardWithdrawalAmount {
		if err := k.EnqueueRewardWithdrawal(ctx, delegation.DelegatorAddress, valAddr.String(), claimedReward); err != nil {
			return errors.Wrap(err, "enqueue reward withdrawal")
		}
	}

	return nil
}

// EnqueueRewardWithdrawal enqueues the reward withdrawal to mint reward IP token on EL side.
func (k Keeper) EnqueueRewardWithdrawal(ctx context.Context, delAddrBech32, valAddrBech32 string, claimedReward uint64) error {
	valAddr, err := sdk.ValAddressFromBech32(valAddrBech32)
	if err != nil {
		return errors.Wrap(err, "validator address from bech32")
	}

	valAccAddr := sdk.AccAddress(valAddr).String()

	// Withdraw delegation rewards.
	delAddr, err := sdk.AccAddressFromBech32(delAddrBech32)
	if err != nil {
		return errors.Wrap(err, "convert acc address from bech32 address")
	}
	delRewards, err := k.distributionKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	if err != nil {
		return err
	}

	// Withdraw commission if it is a self delegation.
	if delAddrBech32 == valAccAddr {
		commissionRewards, err := k.distributionKeeper.WithdrawValidatorCommission(ctx, valAddr)
		if errors.Is(err, dtypes.ErrNoValidatorCommission) {
			log.Debug(
				ctx, "No validator commission",
				"validator_addr", valAddrBech32,
				"validator_account_addr", valAccAddr,
			)
		} else if err != nil {
			return err
		} else {
			delRewards = delRewards.Add(commissionRewards...)
		}
	}

	if claimedReward > 0 {
		rewardCoins := sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(int64(claimedReward)))
		delRewards = delRewards.Add(rewardCoins)
	}

	delRewardUint64 := delRewards.AmountOf(sdk.DefaultBondDenom).Uint64()

	withdrawalEVMAddr, err := k.DelegatorRewardAddress.Get(ctx, delAddrBech32)
	if err != nil {
		return errors.Wrap(err, "map delegator bech32 address to evm reward address")
	}

	log.Debug(
		ctx, "Withdraw delegator rewards",
		"validator_addr", valAddrBech32,
		"validator_account_addr", valAccAddr,
		"delegator_addr", delAddrBech32,
		"withdrawal_evm_addr", withdrawalEVMAddr,
		"amount_claimed", claimedReward,
		"amount_total_withdrawal", delRewardUint64,
	)

	// Burn tokens from the delegator
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, delAddr, types.ModuleName, delRewards)
	if err != nil {
		return err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, delRewards)
	if err != nil {
		return err
	}

	// Enqueue to the global reward withdrawal queue.
	if err := k.AddRewardWithdrawalToQueue(ctx, types.NewWithdrawal(
		uint64(sdk.UnwrapSDKContext(ctx).BlockHeight()),
		withdrawalEVMAddr,
		delRewardUint64,
		types.WithdrawalType_WITHDRAWAL_TYPE_REWARD,
	)); err != nil {
		return errors.Wrap(err, "add reward withdrawal to queue")
	}

	return nil
}

func (k Keeper) ProcessWithdraw(ctx context.Context, ev *bindings.IPTokenStakingWithdraw) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	var actualAmount string

	defer func() {
		if r := recover(); r != nil {
			err = errors.WrapErrWithCode(errors.UnexpectedCondition, fmt.Errorf("panic caused by %v", r))
		}

		var e sdk.Event
		if err == nil {
			writeCache()
			e = sdk.NewEvent(
				types.EventTypeUndelegateSuccess,
				sdk.NewAttribute(types.AttributeKeyAmount, actualAmount),
			)
		} else {
			e = sdk.NewEvent(
				types.EventTypeUndelegateFailure,
				sdk.NewAttribute(types.AttributeKeyErrorCode, errors.UnwrapErrCode(err).String()),
				sdk.NewAttribute(types.AttributeKeyAmount, ev.StakeAmount.String()),
			)
		}

		sdkCtx.EventManager().EmitEvents(sdk.Events{
			e.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyDelegatorAddress, ev.Delegator.String()),
				sdk.NewAttribute(types.AttributeKeyValidatorCmpPubKey, hex.EncodeToString(ev.ValidatorCmpPubkey)),
				sdk.NewAttribute(types.AttributeKeyDelegateID, ev.DelegationId.String()),
				sdk.NewAttribute(types.AttributeKeyAmount, ev.StakeAmount.String()),
				sdk.NewAttribute(types.AttributeKeySenderAddress, ev.OperatorAddress.Hex()),
				sdk.NewAttribute(types.AttributeKeyTxHash, hex.EncodeToString(ev.Raw.TxHash.Bytes())),
			),
		})
	}()

	isInSingularity, err := k.IsSingularity(cachedCtx)
	if err != nil {
		return errors.Wrap(err, "check if it is singularity")
	}

	if isInSingularity {
		log.Debug(cachedCtx, "Withdraw event detected, but it is not processed since current block is singularity")
		return nil
	}

	validatorPubkey, err := k1util.PubKeyBytesToCosmos(ev.ValidatorCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to cosmos")
	}

	valEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(validatorPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "validator pubkey to evm address")
	}

	depositorAddr := sdk.AccAddress(ev.Delegator.Bytes())
	validatorAddr := sdk.ValAddress(valEvmAddr.Bytes())

	// unstakeOnBehalf txn, need to check if it's from the operator
	if ev.Delegator.String() != ev.OperatorAddress.String() {
		operatorAddr, err := k.DelegatorOperatorAddress.Get(cachedCtx, depositorAddr.String())
		if errors.Is(err, collections.ErrNotFound) {
			return errors.WrapErrWithCode(
				errors.InvalidOperator,
				errors.New("invalid unstakeOnBehalf txn, no operator"),
			)
		} else if err != nil {
			return errors.Wrap(err, "get delegator's operator address failed")
		}
		if operatorAddr != ev.OperatorAddress.String() {
			return errors.WrapErrWithCode(
				errors.InvalidOperator,
				errors.New("invalid unstakeOnBehalf txn, not from operator"),
			)
		}
	}

	amountCoin, _ := IPTokenToBondCoin(ev.StakeAmount)

	log.Debug(cachedCtx, "Processing EVM staking withdraw",
		"del_story", depositorAddr.String(),
		"val_story", validatorAddr.String(),
		"del_evm_addr", ev.Delegator.String(),
		"val_evm_addr", valEvmAddr.String(),
		"amount", ev.StakeAmount.String(),
	)

	if !k.authKeeper.HasAccount(cachedCtx, depositorAddr) {
		return errors.New("depositor account not found")
	}

	lockedTokenType, err := k.stakingKeeper.GetLockedTokenType(cachedCtx)
	if err != nil {
		return errors.Wrap(err, "get locked token type")
	}

	val, err := k.stakingKeeper.GetValidator(cachedCtx, validatorAddr)
	if errors.Is(err, stypes.ErrNoValidatorFound) {
		return errors.WrapErrWithCode(errors.ValidatorNotFound, err)
	} else if err != nil {
		return errors.Wrap(err, "get validator failed")
	}

	// locked tokens only have delegation with flexible period,
	// here we automatically set the delegation id to the flexible period delegation id
	delID := ev.DelegationId.String()
	if val.SupportTokenType == lockedTokenType {
		delID = stypes.FlexiblePeriodDelegationID
	}

	msg := stypes.NewMsgUndelegate(depositorAddr.String(), validatorAddr.String(), delID, amountCoin)

	// Undelegate from the validator (validator existence is checked in ValidateUnbondAmount)
	resp, err := skeeper.NewMsgServerImpl(k.stakingKeeper.(*skeeper.Keeper)).Undelegate(cachedCtx, msg)
	if errors.Is(err, stypes.ErrNoValidatorFound) {
		return errors.WrapErrWithCode(errors.ValidatorNotFound, err)
	} else if errors.Is(err, stypes.ErrNoDelegation) {
		return errors.WrapErrWithCode(errors.DelegationNotFound, err)
	} else if errors.Is(err, stypes.ErrNoPeriodDelegation) {
		return errors.WrapErrWithCode(errors.PeriodDelegationNotFound, err)
	} else if err != nil {
		return errors.Wrap(err, "undelegate")
	}

	actualAmount = resp.Amount.Amount.String()

	log.Debug(cachedCtx, "EVM staking withdraw processed",
		"delegator", depositorAddr.String(),
		"validator", validatorAddr.String(),
		"actual_amount", actualAmount,
		"completion_time", resp.CompletionTime,
	)

	return nil
}

func (k Keeper) ParseWithdrawLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingWithdraw, error) {
	return k.ipTokenStakingContract.ParseWithdraw(ethlog)
}
