package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the upgrade types on the provided
// LegacyAmino codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgScheduleUpgrade{}, URLMsgScheduleUpgrade, nil)
}

// RegisterInterfaces registers the upgrade module types on the provided
// registry.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgScheduleUpgrade{})
	msgservice.RegisterMsgServiceDesc(registry, &_MsgService_serviceDesc)
}
