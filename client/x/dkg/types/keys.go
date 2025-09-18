package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name.
	ModuleName = "dkg"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key.
	MemStoreKey = "mem_dkg"
)

// KVStore key prefixes.
var (
	ParamsKey           = collections.NewPrefix(0)
	DKGNetworkKey       = collections.NewPrefix(1)
	LatestDKGNetworkKey = collections.NewPrefix(2)
	DKGRegistrationKey  = collections.NewPrefix(3)
	TEEUpgradeInfoKey   = collections.NewPrefix(4)
)
