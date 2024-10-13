package types

import "cosmossdk.io/collections"

const (
	// ModuleName is the name of the evmstaking module.
	ModuleName = "evmstaking"

	// StoreKey is the string store representation.
	StoreKey = ModuleName

	// RouterKey is the msg router key for the evmstaking module.
	RouterKey = ModuleName
)

// KVStore keys.
var (
	ParamsKey              = collections.NewPrefix(0)
	WithdrawalQueueKey     = collections.NewPrefix(1)
	DelegatorMapKey        = collections.NewPrefix(2)
	ValidatorSweepIndexKey = collections.NewPrefix(3)
)
