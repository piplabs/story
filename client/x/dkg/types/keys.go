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

// CDRFeePoolName is the module account name for the CDR reward pool.
// Fees collected by the CDR contract are minted and held here,
// then distributed to validators at each DKG round end.
const CDRFeePoolName = "cdr-fee-pool"

// KVStore key prefixes.
var (
	ParamsKey                 = collections.NewPrefix(0)
	DKGNetworkKey             = collections.NewPrefix(1)
	LatestDKGNetworkKey       = collections.NewPrefix(2)
	DKGRegistrationKey        = collections.NewPrefix(3)
	LatestActiveRoundKey      = collections.NewPrefix(4)
	GlobalPubKeyVotesKey      = collections.NewPrefix(5)
	KernelUpgradeInfoKey      = collections.NewPrefix(6)
	SettlementBalanceKey      = collections.NewPrefix(7)
	DKGPartialDecryptKey      = collections.NewPrefix(8)
	DecryptRequestRegistryKey = collections.NewPrefix(9)
	CDRPartialSubmitCountKey  = collections.NewPrefix(10)
	CDRFeePoolBalanceKey      = collections.NewPrefix(11)
)
