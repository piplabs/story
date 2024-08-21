package keeper

import (
	"context"
	"strconv"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/piplabs/story/client/x/evmstaking/types"
)

type msgServer struct {
	*Keeper
	types.UnimplementedMsgServiceServer
}

func NewMsgServerImpl(keeper *Keeper) types.MsgServiceServer {
	return &msgServer{Keeper: keeper}
}

func (s msgServer) AddWithdrawal(ctx context.Context, msg *types.MsgAddWithdrawal) (*types.MsgAddWithdrawalResponse, error) {
	validatorAddr, err := s.validatorAddressCodec.StringToBytes(msg.Withdrawal.ValidatorAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	// delegatorAddr
	_, err = s.authKeeper.AddressCodec().StringToBytes(msg.Withdrawal.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}

	if msg.Withdrawal.Amount <= 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid withdrawal amount")
	}

	validator, err := (s.stakingKeeper).GetValidator(ctx, validatorAddr)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "validator not found")
	}

	if validator.IsJailed() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "validator is jailed")
	}

	// if validator.IsUnbonding() {
	//	return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "validator is unbonding")
	//}

	// TODO: when validator is unbonded and stakes are also unbonded back to delegates, figure out how to allow withdrawal to EL
	if validator.IsUnbonded() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "validator is unbonded")
	}

	// TODO: check bond denom once Amount is math.Int with coin type

	// bondDenom, err := (*s.stakingKeeper).BondDenom(ctx)
	// if err != nil {
	//	return nil, err
	//}
	//
	// if msg.Withdrawal.Amount.Denom != bondDenom {
	//	return nil, errorsmod.Wrapf(
	//		sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Withdrawal.Amount.Denom, bondDenom,
	//	)
	//}

	// TODO: balance check and transfer balance into this module, so that the module can burn the balance when the withdrawal is executed

	err = s.AddWithdrawalToQueue(ctx, types.NewWithdrawalFromMsg(msg))
	if err != nil {
		return nil, errorsmod.Wrap(err, "add withdrawal")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddWithdrawal,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Withdrawal.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.Withdrawal.DelegatorAddress),
			sdk.NewAttribute(types.AttributeKeyExecutionAddress, msg.Withdrawal.ExecutionAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, strconv.FormatUint(msg.Withdrawal.Amount, 10)),
			sdk.NewAttribute(types.AttributeKeyCreationHeight, strconv.FormatUint(msg.Withdrawal.CreationHeight, 10)),
		),
	})

	return &types.MsgAddWithdrawalResponse{
		//RequestIndex: s.,
	}, nil
}
