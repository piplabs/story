package types

const (
	// ModuleName is the name of the signal module.
	ModuleName = "signal"

	// StoreKey is the string store representation.
	StoreKey = ModuleName

	// QuerierRoute is the querier route key for the signal module.
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the signal module.
	RouterKey = ModuleName
)

// KVStore keys.
var (
	// UpgradeKey is the key in the signal store used to persist an upgrade if one is pending.
	UpgradeKey = []byte{0x00}

	// FirstSignalKey is used as a divider to separate the UpgradeKey from all
	// the keys associated with signals from validators. In practice, this key
	// isn't expected to be set or retrieved.
	FirstSignalKey = []byte{0x000}
)