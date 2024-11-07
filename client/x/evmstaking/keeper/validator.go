//nolint:contextcheck // use cached context
package keeper

import (
	"context"
	"encoding/hex"
	"strconv"

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

//nolint:contextcheck // already inherited new context
func (k Keeper) ProcessCreateValidator(ctx context.Context, ev *bindings.IPTokenStakingCreateValidator) (err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cachedCtx, writeCache := sdkCtx.CacheContext()

	defer func() {
		if err == nil {
			writeCache()
			return
		}
		sdkCtx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeCreateValidatorFailure,
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(sdkCtx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyValidatorUncmpPubKey, hex.EncodeToString(ev.ValidatorUncmpPubkey)),
				sdk.NewAttribute(types.AttributeKeyMoniker, ev.Moniker),
				sdk.NewAttribute(types.AttributeKeyAmount, ev.StakeAmount.String()),
				sdk.NewAttribute(types.AttributeKeyCommissionRate, strconv.FormatUint(uint64(ev.CommissionRate), 10)),
				sdk.NewAttribute(types.AttributeKeyMaxCommissionRate, strconv.FormatUint(uint64(ev.MaxCommissionRate), 10)),
				sdk.NewAttribute(types.AttributeKeyMaxCommissionChangeRate, strconv.FormatUint(uint64(ev.MaxCommissionChangeRate), 10)),
				sdk.NewAttribute(types.AttributeKeyTokenType, strconv.FormatUint(uint64(ev.SupportsUnlocked), 10)),
				sdk.NewAttribute(types.AttributeKeySenderAddress, ev.OperatorAddress.Hex()),
				sdk.NewAttribute(types.AttributeKeyStatusCode, errors.UnwrapErrCode(err).String()),
			),
		})
	}()

	// When creating a validator, it's self-delegation. Thus, validator pubkey is also delegation pubkey.
	valCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.ValidatorUncmpPubkey)
	if err != nil {
		return errors.WrapErrWithCode(errors.InvalidUncmpPubKey, errors.Wrap(err, "compress validator pubkey"))
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
	if !k.authKeeper.HasAccount(cachedCtx, delegatorAddr) {
		acc := k.authKeeper.NewAccountWithAddress(cachedCtx, delegatorAddr)
		k.authKeeper.SetAccount(cachedCtx, acc)
		log.Debug(cachedCtx, "Created account for depositor",
			"address", validatorAddr.String(),
			"evm_address", delEvmAddr.String(),
		)
	}

	log.Debug(cachedCtx, "EVM staking create validator detected",
		"val_story", validatorAddr.String(),
		"val_pubkey", validatorPubkey.String(),
		"del_story", delegatorAddr.String(),
		"del_evm_addr", delEvmAddr.String(),
		"amount_coin", amountCoin.String(),
	)

	// TODO: Check if we can instantiate the msgServer without type assertion
	evmstakingSKeeper, ok := k.stakingKeeper.(*skeeper.Keeper)
	if !ok {
		return errors.New("type assertion failed")
	}
	skeeperMsgServer := skeeper.NewMsgServerImpl(evmstakingSKeeper)

	if _, err = k.stakingKeeper.GetValidator(cachedCtx, validatorAddr); err == nil {
		return errors.WrapErrWithCode(errors.ValidatorAlreadyExists, errors.New("validator already exists"))
	} else if !errors.Is(err, stypes.ErrNoValidatorFound) {
		// Either the validator does not exist, or unknown error.
		return errors.Wrap(err, "get validator")
	}

	moniker := ev.Moniker
	if moniker == "" {
		moniker = validatorAddr.String() // use validator address as moniker if not provided
	}

	minSelfDelegation, err := k.stakingKeeper.MinDelegation(cachedCtx)
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
		int32(ev.SupportsUnlocked),
	)
	if err != nil {
		return errors.Wrap(err, "create validator message")
	}

	// Note that, after minting, we save the mapping between delegator bech32 address and evm address, which will be used in the withdrawal queue.
	// The saving is done regardless of any error below, as the money is already minted and sent to the delegator, who can withdraw the minted amount.
	// TODO: Confirm that bech32 address and evm address can be used interchangeably. Must be one-to-one or many-bech32-to-one-evm.
	// NOTE: Do not overwrite the existing withdraw/reward address set by the validator.
	if exists, err := k.DelegatorWithdrawAddress.Has(cachedCtx, delegatorAddr.String()); err != nil {
		return errors.Wrap(err, "check delegator withdraw address existence")
	} else if !exists {
		if err := k.DelegatorWithdrawAddress.Set(cachedCtx, delegatorAddr.String(), delEvmAddr.String()); err != nil {
			return errors.Wrap(err, "set delegator withdraw address map")
		}
	}
	if exists, err := k.DelegatorRewardAddress.Has(cachedCtx, delegatorAddr.String()); err != nil {
		return errors.Wrap(err, "check delegator reward address existence")
	} else if !exists {
		if err := k.DelegatorRewardAddress.Set(cachedCtx, delegatorAddr.String(), delEvmAddr.String()); err != nil {
			return errors.Wrap(err, "set delegator reward address map")
		}
	}

	if err := k.bankKeeper.MintCoins(cachedCtx, types.ModuleName, amountCoins); err != nil {
		return errors.Wrap(err, "create stake coin for depositor: mint coins")
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(cachedCtx, types.ModuleName, delegatorAddr, amountCoins); err != nil {
		return errors.Wrap(err, "create stake coin for depositor: send coins")
	}

	if _, err = skeeperMsgServer.CreateValidator(cachedCtx, msg); errors.Is(err, stypes.ErrCommissionLTMinRate) {
		return errors.WrapErrWithCode(errors.InvalidCommissionRate, err)
	} else if errors.Is(err, stypes.ErrMinSelfDelegationBelowMinDelegation) {
		return errors.WrapErrWithCode(errors.InvalidMinSelfDelegation, err)
	} else if errors.Is(err, stypes.ErrNoTokenTypeFound) {
		return errors.WrapErrWithCode(errors.InvalidTokenType, err)
	} else if err != nil {
		return errors.Wrap(err, "create validator")
	}

	return nil
}

func (k Keeper) ParseCreateValidatorLog(ethlog ethtypes.Log) (*bindings.IPTokenStakingCreateValidator, error) {
	return k.ipTokenStakingContract.ParseCreateValidator(ethlog)
}
