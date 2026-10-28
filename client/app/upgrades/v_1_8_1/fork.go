//nolint:staticcheck,revive,nolintlint // versioning
package v_1_8_1

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/upgrades"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/netconf"
)

// UpgradeName is the on-chain name for the v1.8.1 upgrade that retires the DKG
// module and the CDR stack, bringing a network that ran them onto the same
// state machine as the networks that never did.
//
// This is a binary-swap upgrade, the mirror of v1.6.0 which activated the
// module. It is driven by the compiled fork height, not by an on-chain plan:
// at that height a binary without a handler for this name halts, and operators
// replace it with one built from this branch.
//
// The name is consumed permanently once a chain records it as done, because
// x/upgrade refuses to schedule a completed name again.
const UpgradeName = netconf.V181

// Fork drives the switch from the compiled height instead of an on-chain plan.
// scheduleForkUpgrade schedules the plan at exactly that height and x/upgrade
// applies it in the same block, so every node switches on the same block
// without an operator transaction.
var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeInfo:    "retire the DKG module and the CDR stack",
	BeginForkLogic: func(_ sdk.Context, _ *keepers.Keepers) {},
}

// GetUpgradeHeight resolves the fork height per chain, so only a chain that
// registered one takes part. A chain with no entry is not applicable rather
// than broken: returning false there keeps the fork off every other network,
// which the default branch of the app's resolver would not do.
func GetUpgradeHeight(ctx sdk.Context) (int64, bool) {
	height, err := netconf.GetUpgradeHeight(ctx.ChainID(), netconf.V181)
	if err != nil {
		if !errors.Is(err, netconf.ErrUnknownUpgrade) && !errors.Is(err, netconf.ErrUnknownChainID) {
			log.Error(ctx, "Failed to get upgrade height", err,
				"chain_id", ctx.ChainID(), "upgrade_name", netconf.V181)
		}

		return 0, false
	}

	return height, true
}
