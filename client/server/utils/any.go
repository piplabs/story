//nolint:wrapcheck // Internal utils, don't need to wrap it.
package utils

import codectypes "github.com/cosmos/cosmos-sdk/codec/types"

type unpackAny[T any] struct {
	Any *codectypes.Any
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces, where
// UnpackInterfacesMessage is meant to extend protobuf types (which implement
// proto.Message) to support a post-deserialization phase which unpacks
// types packed within Any's using the whitelist provided by AnyUnpacker.
func (t *unpackAny[T]) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var value T
	return unpacker.UnpackAny(t.Any, &value)
}

// WrapTypeAny implements a wrap function for variables with type Any to become UnpackInterfacesMessage,
// which is meant to unpack values packed within Any's using the AnyUnpacker.
func WrapTypeAny[T any](v *codectypes.Any) codectypes.UnpackInterfacesMessage {
	return &unpackAny[T]{Any: v}
}
