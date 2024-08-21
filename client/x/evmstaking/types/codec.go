package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(_ *codec.LegacyAmino) {
	// cdc.RegisterConcrete(&MsgAddWithdrawal{}, "evmstaking/MsgAddWithdrawal", nil)
	// cdc.RegisterConcrete(&MsgRemoveWithdrawal{}, "evmstaking/MsgRemoveWithdrawal", nil)
}

// RegisterInterfaces registers the x/staking interfaces types with the interface registry.
func RegisterInterfaces(registrar cdctypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddWithdrawal{},
		&MsgRemoveWithdrawal{},
	)

	msgservice.RegisterMsgServiceDesc(registrar, &_MsgService_serviceDesc)
}
