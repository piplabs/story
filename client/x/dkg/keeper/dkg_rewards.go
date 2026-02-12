package keeper

import (
	"context"
	"sort"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/server/utils"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/cast"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// DistributeRewardsToActiveCommittee distributes the configured portion of UBI rewards
// to the current active DKG committee members. The senderModule must hold totalAmount
// coins. Returns the total amount distributed to committee members (0 if no active committee).
//
// This is called by evmstaking during ProcessUbiWithdrawal (EndBlock) to reward
// the currently serving DKG committee on every UBI withdrawal cycle.
func (k *Keeper) DistributeRewardsToActiveCommittee(ctx context.Context, senderModule string, totalAmount math.Int) (math.Int, error) {
	if totalAmount.IsZero() {
		return math.ZeroInt(), nil
	}

	// Get the latest active DKG network (current serving committee).
	activeRound, err := k.getLatestActiveDKGNetwork(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// No active DKG round — no committee to reward.
			return math.ZeroInt(), nil
		}

		return math.ZeroInt(), errors.Wrap(err, "get latest active DKG network")
	}

	return k.distributeRewardsFromModule(ctx, activeRound, senderModule, totalAmount)
}

// distributeRewardsFromModule calculates and distributes DKG committee rewards for the
// given round. Coins are transferred from senderModule to the DKG module, then
// distributed to individual committee members. Returns the total amount distributed.
func (k *Keeper) distributeRewardsFromModule(ctx context.Context, round *types.DKGNetwork, senderModule string, totalAmount math.Int) (math.Int, error) {
	// Read the DKG committee reward portion from params.
	params, err := k.GetParams(ctx)
	if err != nil {
		return math.ZeroInt(), errors.Wrap(err, "get DKG params")
	}

	portion := params.DkgCommitteeRewardPortion
	if portion.IsZero() {
		return math.ZeroInt(), nil
	}

	// Get the finalized committee members.
	codeCommitment32, err := cast.ToBytes32(round.CodeCommitment)
	if err != nil {
		return math.ZeroInt(), errors.Wrap(err, "cast code commitment to bytes32")
	}

	finalizedRegs, err := k.getDKGRegistrationsByStatus(ctx, codeCommitment32, round.Round, types.DKGRegStatusFinalized)
	if err != nil {
		return math.ZeroInt(), errors.Wrap(err, "get finalized DKG registrations")
	}

	if len(finalizedRegs) == 0 {
		return math.ZeroInt(), nil
	}

	// Collect and sort member EVM addresses for deterministic iteration.
	memberAddrs := make([]string, 0, len(finalizedRegs))
	for _, reg := range finalizedRegs {
		memberAddrs = append(memberAddrs, reg.ValidatorAddr)
	}
	sort.Strings(memberAddrs)

	// Calculate DKG reward amounts.
	memberCount := math.NewInt(int64(len(memberAddrs)))
	dkgReward := portion.MulInt(totalAmount).TruncateInt()
	perMemberReward := dkgReward.Quo(memberCount)

	if perMemberReward.IsZero() {
		return math.ZeroInt(), nil
	}

	// Transfer the exact distribution amount from sender module to DKG module.
	totalToDistribute := perMemberReward.Mul(memberCount)
	distributeCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, totalToDistribute))

	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, senderModule, types.ModuleName, distributeCoins); err != nil {
		return math.ZeroInt(), errors.Wrap(err, "transfer coins to DKG module for distribution")
	}

	// Distribute rewards to each committee member.
	perMemberCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, perMemberReward))

	for _, evmAddr := range memberAddrs {
		recipientAddr, err := utils.EvmAddressToBech32AccAddress(evmAddr)
		if err != nil {
			return math.ZeroInt(), errors.Wrap(err, "convert EVM address to bech32",
				"address", evmAddr,
			)
		}

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, perMemberCoins); err != nil {
			return math.ZeroInt(), errors.Wrap(err, "send DKG committee reward",
				"recipient", evmAddr,
				"amount", perMemberReward.String(),
			)
		}
	}

	// Emit event.
	if err := k.emitDKGCommitteeRewarded(ctx, round, uint32(len(memberAddrs)), totalToDistribute, perMemberReward, totalAmount.Sub(totalToDistribute)); err != nil {
		return math.ZeroInt(), errors.Wrap(err, "emit DKG committee rewarded event")
	}

	log.Info(ctx, "Distributed DKG committee rewards",
		"round", round.Round,
		"member_count", len(memberAddrs),
		"total_distributed", totalToDistribute.String(),
		"per_member", perMemberReward.String(),
	)

	return totalToDistribute, nil
}

// ClaimSettlementBalance transfers any remaining UBI settlement balance from the
// DKG module to the specified recipient module. This is called by evmstaking during
// EndBlock to sweep leftover UBI that was not distributed to committee members
// during FinalizeDKGRound. Returns the amount transferred (0 if no settlement pending).
func (k *Keeper) ClaimSettlementBalance(ctx context.Context, recipientModule string) (math.Int, error) {
	balStr, err := k.SettlementBalance.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return math.ZeroInt(), nil
		}

		return math.ZeroInt(), errors.Wrap(err, "get settlement balance")
	}

	balance, ok := math.NewIntFromString(balStr)
	if !ok || balance.IsZero() {
		// Clear invalid or zero entry.
		if removeErr := k.SettlementBalance.Remove(ctx); removeErr != nil {
			return math.ZeroInt(), errors.Wrap(removeErr, "remove settlement balance")
		}

		return math.ZeroInt(), nil
	}

	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, balance))
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, recipientModule, coins); err != nil {
		return math.ZeroInt(), errors.Wrap(err, "transfer settlement balance to recipient module")
	}

	if err := k.SettlementBalance.Remove(ctx); err != nil {
		return math.ZeroInt(), errors.Wrap(err, "remove settlement balance after transfer")
	}

	log.Info(ctx, "Claimed DKG settlement balance",
		"recipient_module", recipientModule,
		"amount", balance.String(),
	)

	return balance, nil
}

// settleRewardsForPreviousCommittee distributes a portion of accrued UBI rewards to the
// members of the previous active DKG committee. The remainder is stored in the
// SettlementBalance state for evmstaking to sweep during EndBlock.
//
// This is called during FinalizeDKGRound, right before the new active round pointer
// is updated, so getLatestActiveDKGNetwork still returns the *previous* active round.
func (k *Keeper) settleRewardsForPreviousCommittee(ctx context.Context) error {
	// Get the previous active round (the one whose committee served).
	prevActive, err := k.getLatestActiveDKGNetwork(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// First round ever — no previous committee to reward.
			return nil
		}

		return errors.Wrap(err, "failed to get previous active DKG network")
	}

	// Read the DKG committee reward portion from params.
	params, err := k.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get DKG params")
	}

	if params.DkgCommitteeRewardPortion.IsZero() {
		return nil
	}

	// Check current UBI balance before attempting withdrawal.
	ubiBalance, err := k.distributionKeeper.GetUbiBalanceByDenom(ctx, sdk.DefaultBondDenom)
	if err != nil {
		return errors.Wrap(err, "failed to get UBI balance")
	}

	if ubiBalance.IsZero() {
		return nil
	}

	// Withdraw ALL UBI into the DKG module account.
	withdrawnCoin, err := k.distributionKeeper.WithdrawUbiByDenomToModule(ctx, sdk.DefaultBondDenom, types.ModuleName)
	if err != nil {
		return errors.Wrap(err, "failed to withdraw UBI to DKG module")
	}

	withdrawnAmount := withdrawnCoin.Amount
	if withdrawnAmount.IsZero() {
		return nil
	}

	// Distribute to the previous active committee from the DKG module account.
	distributed, err := k.distributeFromModuleBalance(ctx, prevActive, withdrawnAmount)
	if err != nil {
		return errors.Wrap(err, "failed to distribute rewards to previous committee")
	}

	// Store the remaining in settlement balance for evmstaking to sweep during EndBlock.
	remaining := withdrawnAmount.Sub(distributed)
	if remaining.IsPositive() {
		if err := k.SettlementBalance.Set(ctx, remaining.String()); err != nil {
			return errors.Wrap(err, "failed to set settlement balance")
		}
	}

	return nil
}

// distributeFromModuleBalance distributes rewards from the DKG module's own balance
// to committee members (used during FinalizeDKGRound where coins are already
// in the DKG module account). Returns the total amount distributed.
func (k *Keeper) distributeFromModuleBalance(ctx context.Context, round *types.DKGNetwork, withdrawnAmount math.Int) (math.Int, error) {
	// Read the DKG committee reward portion from params.
	params, err := k.GetParams(ctx)
	if err != nil {
		return math.ZeroInt(), errors.Wrap(err, "get DKG params")
	}

	portion := params.DkgCommitteeRewardPortion
	if portion.IsZero() {
		return math.ZeroInt(), nil
	}

	// Get the finalized committee members.
	codeCommitment32, err := cast.ToBytes32(round.CodeCommitment)
	if err != nil {
		return math.ZeroInt(), errors.Wrap(err, "cast code commitment to bytes32")
	}

	finalizedRegs, err := k.getDKGRegistrationsByStatus(ctx, codeCommitment32, round.Round, types.DKGRegStatusFinalized)
	if err != nil {
		return math.ZeroInt(), errors.Wrap(err, "get finalized DKG registrations")
	}

	if len(finalizedRegs) == 0 {
		return math.ZeroInt(), nil
	}

	// Collect and sort member EVM addresses for deterministic iteration.
	memberAddrs := make([]string, 0, len(finalizedRegs))
	for _, reg := range finalizedRegs {
		memberAddrs = append(memberAddrs, reg.ValidatorAddr)
	}
	sort.Strings(memberAddrs)

	// Calculate DKG reward amounts.
	memberCount := math.NewInt(int64(len(memberAddrs)))
	dkgReward := portion.MulInt(withdrawnAmount).TruncateInt()
	perMemberReward := dkgReward.Quo(memberCount)

	// Distribute rewards to each committee member directly from DKG module.
	totalDistributed := math.ZeroInt()

	if perMemberReward.IsPositive() {
		perMemberCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, perMemberReward))

		for _, evmAddr := range memberAddrs {
			recipientAddr, err := utils.EvmAddressToBech32AccAddress(evmAddr)
			if err != nil {
				return math.ZeroInt(), errors.Wrap(err, "convert EVM address to bech32",
					"address", evmAddr,
				)
			}

			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, perMemberCoins); err != nil {
				return math.ZeroInt(), errors.Wrap(err, "send DKG committee reward",
					"recipient", evmAddr,
					"amount", perMemberReward.String(),
				)
			}

			totalDistributed = totalDistributed.Add(perMemberReward)
		}
	}

	// Emit event.
	remaining := withdrawnAmount.Sub(totalDistributed)

	if err := k.emitDKGCommitteeRewarded(ctx, round, uint32(len(memberAddrs)), totalDistributed, perMemberReward, remaining); err != nil {
		return math.ZeroInt(), errors.Wrap(err, "emit DKG committee rewarded event")
	}

	log.Info(ctx, "Distributed DKG committee rewards",
		"round", round.Round,
		"member_count", len(memberAddrs),
		"total_distributed", totalDistributed.String(),
		"per_member", perMemberReward.String(),
		"remaining", remaining.String(),
	)

	return totalDistributed, nil
}
