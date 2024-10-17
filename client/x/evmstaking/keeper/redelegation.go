package keeper

import (
	"context"

	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
)

func (k Keeper) ProcessRedelegate(ctx context.Context, ev *bindings.IPTokenStakingRedelegate) error {
	isInSingularity, err := k.IsSingularity(ctx)
	if err != nil {
		return errors.Wrap(err, "check if it is singularity")
	}

	if isInSingularity {
		log.Debug(ctx, "Relegation event detected, but it is not processed since current block is singularity")
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

	valSrcCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.ValidatorUncmpSrcPubkey)
	if err != nil {
		return errors.Wrap(err, "compress src validator pubkey")
	}
	validatorSrcPubkey, err := k1util.PubKeyBytesToCosmos(valSrcCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "src validator pubkey to cosmos")
	}

	valDstCmpPubkey, err := UncmpPubKeyToCmpPubKey(ev.ValidatorUncmpDstPubkey)
	if err != nil {
		return errors.Wrap(err, "compress dst validator pubkey")
	}
	validatorDstPubkey, err := k1util.PubKeyBytesToCosmos(valDstCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "dst validator pubkey to cosmos")
	}

	depositorAddr := sdk.AccAddress(depositorPubkey.Address().Bytes())
	validatorSrcAddr := sdk.ValAddress(validatorSrcPubkey.Address().Bytes())
	validatorDstAddr := sdk.ValAddress(validatorDstPubkey.Address().Bytes())

	delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(depositorPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "deledator pubkey to evm address")
	}
	valSrcEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(validatorSrcPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "src validator pubkey to evm address")
	}
	valDstEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(validatorDstPubkey.Bytes())
	if err != nil {
		return errors.Wrap(err, "dst validator pubkey to evm address")
	}

	// redelegateOnBehalf txn, need to check if it's from the operator
	if delEvmAddr.String() != ev.OperatorAddress.String() {
		operatorAddr, err := k.DelegatorOperatorAddress.Get(ctx, depositorAddr.String())
		if errors.Is(err, collections.ErrNotFound) {
			return errors.New("invalid redelegateOnBehalf txn, no operator for delegator")
		} else if err != nil {
			return errors.Wrap(err, "get delegator's operator address failed")
		}
		if operatorAddr != ev.OperatorAddress.String() {
			return errors.New("invalid redelegateOnBehalf txn, not from operator")
		}
	}

	amountCoin, _ := IPTokenToBondCoin(ev.Amount)

	log.Debug(ctx, "EVM staking relegation detected",
		"del_story", depositorAddr.String(),
		"val_src_story", validatorSrcAddr.String(),
		"val_dst_story", validatorDstAddr.String(),
		"del_evm_addr", delEvmAddr.String(),
		"val_src_evm_addr", valSrcEvmAddr.String(),
		"val_dst_evm_addr", valDstEvmAddr.String(),
		"amount_coin", amountCoin.String(),
	)

	msg := stypes.NewMsgBeginRedelegate(
		depositorAddr.String(), validatorSrcAddr.String(), validatorDstAddr.String(),
		ev.DelegationId.String(), amountCoin,
	)
	_, err = skeeper.NewMsgServerImpl(k.stakingKeeper.(*skeeper.Keeper)).BeginRedelegate(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "failed to begin redelegation")
	}

	return nil
}

func (k Keeper) ParseRedelegateLog(ethLog ethtypes.Log) (*bindings.IPTokenStakingRedelegate, error) {
	return k.ipTokenStakingContract.ParseRedelegate(ethLog)
}
