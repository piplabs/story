//nolint:staticcheck,revive,nolintlint // versioning
package v_1_8_1

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/piplabs/story/client/app/upgrades"
)

// Upgrade carries no StoreUpgrades on purpose. The dkg store is not deleted,
// it is only no longer mounted, so it drops out of the app hash at the upgrade
// height on every node and its IAVL data stays behind as unreferenced bytes.
//
// StoreUpgrades{Deleted} cannot reclaim them: rootmulti applies deletions only
// to mounted store keys, so the deleting binary would have to keep the module
// registered, and that binary cannot restart once the store has left CommitInfo.
var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}
