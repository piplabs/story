package keeper

import (
	"context"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func (k Keeper) ProcessUbiWithdrawal(ctx context.Context) error {
	log.Debug(ctx, "Processing eligible ubi withdraw")

	params, err := k.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "get ubi params")
	}

	// Sweep any settlement balance from the DKG module (leftover UBI from FinalizeDKGRound).
	settlementAmount, err := k.dkgKeeper.ClaimSettlementBalance(ctx, types.ModuleName)
	if err != nil {
		return errors.Wrap(err, "claim DKG settlement balance")
	}

	ubiBalance, err := k.distributionKeeper.GetUbiBalanceByDenom(ctx, sdk.DefaultBondDenom)
	if err != nil {
		return errors.Wrap(err, "get ubi balance by denom")
	}

	// Check if UBI balance meets the minimum threshold for withdrawal.
	meetsWithdrawalThreshold := ubiBalance.Uint64() >= params.MinPartialWithdrawalAmount
	withdrawnAmount := math.ZeroInt()

	if meetsWithdrawalThreshold {
		// Withdraw UBI coins to evmstaking module.
		ubiCoin, err := k.distributionKeeper.WithdrawUbiByDenomToModule(
			ctx, sdk.DefaultBondDenom, types.ModuleName,
		)
		if err != nil {
			return errors.Wrap(err, "withdraw ubi by denom to module")
		}

		withdrawnAmount = ubiCoin.Amount

		// Distribute DKG rewards for the current active DKG committee.
		distributed, err := k.dkgKeeper.DistributeRewardsToActiveCommittee(ctx, types.ModuleName, withdrawnAmount)
		if err != nil {
			return errors.Wrap(err, "distribute DKG committee rewards")
		}

		withdrawnAmount = withdrawnAmount.Sub(distributed)
	}

	// Total amount to burn and add to withdrawal queue = withdrawn remainder + settlement.
	totalToBurn := withdrawnAmount.Add(settlementAmount)
	if totalToBurn.IsZero() {
		return nil
	}

	// Burn tokens.
	if err = k.bankKeeper.BurnCoins(
		ctx, types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, totalToBurn)),
	); err != nil {
		return errors.Wrap(err, "burn ubi coins")
	}

	// Add withdrawal entry to the withdrawal queue.
	if err = k.AddWithdrawalToQueue(ctx, types.NewWithdrawal(
		uint64(sdk.UnwrapSDKContext(ctx).BlockHeight()),
		params.UbiWithdrawAddress,
		totalToBurn.Uint64(),
		types.WithdrawalType_WITHDRAWAL_TYPE_UBI,
		"",
	)); err != nil {
		return errors.Wrap(err, "add ubi withdrawal to queue")
	}

	return nil
}
