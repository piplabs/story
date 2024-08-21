package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
)

func (k Keeper) ProcessDeposit(ctx context.Context, ev *bindings.IPTokenStakingDeposit) error {
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

	amountCoin, amountCoins := IPTokenToBondCoin(ev.Amount)

	// Create account if not exists
	if !k.authKeeper.HasAccount(ctx, depositorAddr) {
		acc := k.authKeeper.NewAccountWithAddress(ctx, depositorAddr)
		k.authKeeper.SetAccount(ctx, acc)
		log.Debug(ctx, "Created account for depositor",
			"address", depositorAddr.String(),
			"evm_address", delEvmAddr.String(),
		)
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, amountCoins); err != nil {
		return errors.Wrap(err, "create stake coin for depositor: mint coins")
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositorAddr, amountCoins); err != nil {
		return errors.Wrap(err, "create stake coin for depositor: send coins")
	}

	log.Info(ctx, "EVM staking deposit detected, delegating to validator",
		"del_iliad", depositorAddr.String(),
		"val_iliad", validatorAddr.String(),
		"del_evm_addr", delEvmAddr.String(),
		"val_evm_addr", valEvmAddr.String(),
		"amount_coin", amountCoin.String(),
	)

	// Note that, after minting, we save the mapping between delegator bech32 address and evm address, which will be used in the withdrawal queue.
	// The saving is done regardless of any error below, as the money is already minted and sent to the delegator, who can withdraw the minted amount.
	// TODO: Confirm that bech32 address and evm address can be used interchangeably. Must be one-to-one or many-bech32-to-one-evm.
	if err := k.DelegatorMap.Set(ctx, depositorAddr.String(), delEvmAddr.String()); err != nil {
		return errors.Wrap(err, "set delegator map")
	}

	// TODO: Check if we can instantiate the msgServer without type assertion
	evmstakingSKeeper, ok := k.stakingKeeper.(*skeeper.Keeper)
	if !ok {
		return errors.New("type assertion failed")
	}
	skeeperMsgServer := skeeper.NewMsgServerImpl(evmstakingSKeeper)

	// Delegation by the depositor on the validator (validator existence is checked in msgServer.Delegate)
	msg := stypes.NewMsgDelegate(depositorAddr.String(), validatorAddr.String(), amountCoin)
	_, err = skeeperMsgServer.Delegate(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "delegate")
	}

	return nil
}

func (k Keeper) ParseDepositLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingDeposit, error) {
	return k.ipTokenStakingContract.ParseDeposit(ethlog)
}
