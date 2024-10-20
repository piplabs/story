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
	ParamsKey                      = collections.NewPrefix(0)
	ValidatorSweepIndexKey         = collections.NewPrefix(1)
	DelegatorWithdrawAddressMapKey = collections.NewPrefix(2)
	DelegatorRewardAddressMapKey   = collections.NewPrefix(3)
	DelegatorOperatorAddressMapKey = collections.NewPrefix(4)
	WithdrawalQueueKey             = collections.NewPrefix(5)
	RewardWithdrawalQueueKey       = collections.NewPrefix(6)
)
