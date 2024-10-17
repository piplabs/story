package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterCodec(_ *codec.LegacyAmino) {}

// RegisterInterfaces registers the x/staking interfaces types with the interface registry.
func RegisterInterfaces(_ cdctypes.InterfaceRegistry) {}
