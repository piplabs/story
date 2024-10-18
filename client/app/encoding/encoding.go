package encoding

import (
	"cosmossdk.io/x/tx/signing"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/gogoproto/proto"
)

type ModuleRegister interface {
	RegisterLegacyAminoCodec(*codec.LegacyAmino)
	RegisterInterfaces(codectypes.InterfaceRegistry)
}

// Config specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeConfig returns an encoding config for the app.
func MakeEncodingConfig(moduleRegisters ...ModuleRegister) EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry, err := codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec:          address.NewBech32Codec(AccountAddressPrefix),
			ValidatorAddressCodec: address.NewBech32Codec(ValidatorAddressPrefix),
		},
	})
	if err != nil {
		panic(err)
	}

	// Register the standard types from the Cosmos SDK on interfaceRegistry and amino.
	std.RegisterInterfaces(interfaceRegistry)
	std.RegisterLegacyAminoCodec(amino)

	// Register types from the moduleRegisters on interfaceRegistry and amino.
	for _, moduleRegister := range moduleRegisters {
		moduleRegister.RegisterInterfaces(interfaceRegistry)
		moduleRegister.RegisterLegacyAminoCodec(amino)
	}

	protoCodec := codec.NewProtoCodec(interfaceRegistry)
	txConfig := tx.NewTxConfig(protoCodec, tx.DefaultSignModes)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             protoCodec,
		TxConfig:          txConfig,
		Amino:             amino,
	}
}
