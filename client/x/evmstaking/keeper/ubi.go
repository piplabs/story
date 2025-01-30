package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/piplabs/story/client/server/utils"
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
	ubiBalance, err := k.distributionKeeper.GetUbiBalanceByDenom(ctx, sdk.DefaultBondDenom)
	if err != nil {
		return errors.Wrap(err, "get ubi balance by denom")
	}

	if ubiBalance.Uint64() < params.MinPartialWithdrawalAmount {
		return nil
	}

	moduleEVMAddr, err := utils.Bech32DelegatorAddressToEvmAddress(authtypes.NewModuleAddress(types.ModuleName).String())
	if err != nil {
		return errors.Wrap(err, "convert evmstaking module address to evm address")
	}

	// Withdraw and burn ubi coins.
	ubiCoin, err := k.distributionKeeper.WithdrawUbiByDenomToModule(
		ctx, sdk.DefaultBondDenom, types.ModuleName,
	)
	if err != nil {
		return errors.Wrap(err, "withdraw ubi by denom to module")
	}
	// Burn tokens from the ubi.
	if err = k.bankKeeper.BurnCoins(
		ctx, types.ModuleName,
		sdk.NewCoins(ubiCoin),
	); err != nil {
		return errors.Wrap(err, "burn ubi coins")
	}

	// Add withdrawal entry to the withdrawal queue.
	if err = k.AddWithdrawalToQueue(ctx, types.NewWithdrawal(
		uint64(sdk.UnwrapSDKContext(ctx).BlockHeight()),
		params.UbiWithdrawAddress,
		ubiBalance.Uint64(),
		types.WithdrawalType_WITHDRAWAL_TYPE_UBI,
		moduleEVMAddr,
	)); err != nil {
		return errors.Wrap(err, "add ubi withdrawal to queue")
	}

	return nil
}
