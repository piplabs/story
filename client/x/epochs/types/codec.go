package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterInterfaces(registrar cdctypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil))
}
