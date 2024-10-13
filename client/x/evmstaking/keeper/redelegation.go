package keeper

import (
	"context"

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
	depositorPubkey, err := k1util.PubKeyBytesToCosmos(ev.DelegatorCmpPubkey)
	if err != nil {
		return errors.Wrap(err, "depositor pubkey to cosmos")
	}

	validatorSrcPubkey, err := k1util.PubKeyBytesToCosmos(ev.ValidatorSrcPubkey)
	if err != nil {
		return errors.Wrap(err, "src validator pubkey to cosmos")
	}

	validatorDstPubkey, err := k1util.PubKeyBytesToCosmos(ev.ValidatorDstPubkey)
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

	amountCoin, _ := IPTokenToBondCoin(ev.Amount)

	log.Info(ctx, "EVM staking relegation detected",
		"del_story", depositorAddr.String(),
		"val_src_story", validatorSrcAddr.String(),
		"val_dst_story", validatorDstAddr.String(),
		"del_evm_addr", delEvmAddr.String(),
		"val_src_evm_addr", valSrcEvmAddr.String(),
		"val_dst_evm_addr", valDstEvmAddr.String(),
		"amount_coin", amountCoin.String(),
	)

	msg := stypes.NewMsgBeginRedelegate(depositorAddr.String(), validatorSrcAddr.String(), validatorDstAddr.String(), amountCoin)
	_, err = skeeper.NewMsgServerImpl(k.stakingKeeper.(*skeeper.Keeper)).BeginRedelegate(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "failed to begin redelegation")
	}

	return nil
}

func (k Keeper) ParseRedelegateLog(ethLog ethtypes.Log) (*bindings.IPTokenStakingRedelegate, error) {
	return k.ipTokenStakingContract.ParseRedelegate(ethLog)
}
