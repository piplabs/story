package keeper

import (
	"context"

	"cosmossdk.io/math"

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

func (k Keeper) ProcessCreateValidator(ctx context.Context, ev *bindings.IPTokenStakingCreateValidator) error {
	// When creating a validator, it's self-delegation. Thus, validator pubkey is also delegation pubkey.
	valCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.ValidatorUncmpPubkey)
	if err != nil {
		return errors.Wrap(err, "compress validator pubkey")
	}
	validatorPubkey, err := k1util.PubKeyBytesToCosmos(valCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "validator pubkey to cosmos")
	}

	validatorAddr := sdk.ValAddress(validatorPubkey.Address().Bytes())
	delegatorAddr := sdk.AccAddress(validatorPubkey.Address().Bytes())

	delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(validatorPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "validator pubkey to evm address")
	}

	amountCoin, amountCoins := IPTokenToBondCoin(ev.StakeAmount)

	// Create account if not exists
	if !k.authKeeper.HasAccount(ctx, delegatorAddr) {
		acc := k.authKeeper.NewAccountWithAddress(ctx, delegatorAddr)
		k.authKeeper.SetAccount(ctx, acc)
		log.Debug(ctx, "Created account for depositor",
			"address", validatorAddr.String(),
			"evm_address", delEvmAddr.String(),
		)
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, amountCoins); err != nil {
		return errors.Wrap(err, "create stake coin for depositor: mint coins")
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegatorAddr, amountCoins); err != nil {
		return errors.Wrap(err, "create stake coin for depositor: send coins")
	}

	log.Info(ctx, "EVM staking create validator detected",
		"val_story", validatorAddr.String(),
		"val_pubkey", validatorPubkey.String(),
		"del_story", delegatorAddr.String(),
		"del_evm_addr", delEvmAddr.String(),
		"amount_coin", amountCoin.String(),
	)

	// Note that, after minting, we save the mapping between delegator bech32 address and evm address, which will be used in the withdrawal queue.
	// The saving is done regardless of any error below, as the money is already minted and sent to the delegator, who can withdraw the minted amount.
	// TODO: Confirm that bech32 address and evm address can be used interchangeably. Must be one-to-one or many-bech32-to-one-evm.
	if err := k.DelegatorWithdrawAddress.Set(ctx, delegatorAddr.String(), delEvmAddr.String()); err != nil {
		return errors.Wrap(err, "set delegator withdraw address map")
	}
	if err := k.DelegatorRewardAddress.Set(ctx, delegatorAddr.String(), delEvmAddr.String()); err != nil {
		return errors.Wrap(err, "set delegator reward address map")
	}

	// TODO: Check if we can instantiate the msgServer without type assertion
	evmstakingSKeeper, ok := k.stakingKeeper.(*skeeper.Keeper)
	if !ok {
		return errors.New("type assertion failed")
	}
	skeeperMsgServer := skeeper.NewMsgServerImpl(evmstakingSKeeper)

	if _, err = k.stakingKeeper.GetValidator(ctx, validatorAddr); err == nil {
		// TODO(rayden): refund
		return errors.New("validator already exists")
	} else if !errors.Is(err, stypes.ErrNoValidatorFound) {
		// Either the validator does not exist, or unknown error.
		return errors.Wrap(err, "get validator")
	}

	moniker := ev.Moniker
	if moniker == "" {
		moniker = validatorAddr.String() // use validator address as moniker if not provided
	}

	var tokenType stypes.TokenType
	switch ev.SupportsUnlocked {
	case uint8(stypes.TokenType_LOCKED):
		tokenType = stypes.TokenType_LOCKED
	case uint8(stypes.TokenType_UNLOCKED):
		tokenType = stypes.TokenType_UNLOCKED
	default:
		return errors.New("invalid token type")
	}

	minSelfDelegation, err := k.stakingKeeper.MinDelegation(ctx)
	if err != nil {
		return errors.Wrap(err, "get min self delegation")
	}

	// Validator does not exist, create validator with self-delegation.
	msg, err := stypes.NewMsgCreateValidator(
		validatorAddr.String(),
		validatorPubkey,
		amountCoin,
		stypes.Description{Moniker: moniker},
		stypes.NewCommissionRates(
			// Divide these decimals by 100 to convert from basis points to decimal. Will cut off decimal as the rates are integers.
			math.LegacyNewDec(int64(ev.CommissionRate)).Quo(math.LegacyNewDec(10000)),
			math.LegacyNewDec(int64(ev.MaxCommissionRate)).Quo(math.LegacyNewDec(10000)),
			math.LegacyNewDec(int64(ev.MaxCommissionChangeRate)).Quo(math.LegacyNewDec(10000)),
		),
		minSelfDelegation, // make minimum self delegation align with minimum delegation amount
		tokenType,
	)
	if err != nil {
		return errors.Wrap(err, "create validator message")
	}

	_, err = skeeperMsgServer.CreateValidator(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "create validator")
	}

	return nil
}

func (k Keeper) ParseCreateValidatorLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingCreateValidator, error) {
	return k.ipTokenStakingContract.ParseCreateValidator(ethlog)
}
