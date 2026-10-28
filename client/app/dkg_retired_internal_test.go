package app

// Internal test (package app) so it can read the module tables directly. These
// assertions are what make the v1.8.1 upgrade safe: the chain drops the dkg
// store from its app hash only because this binary no longer mounts it, and it
// stops expecting a DKG vote in every block only because the module is gone
// from these tables. A reintroduced reference would change the app hash of a
// chain that has already retired the module.

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/app/upgrades/v_1_8_1"
	"github.com/piplabs/story/lib/netconf"
)

// mentionsDKG reports whether a module, account or permission name belongs to
// the retired module. It matches on substrings because the stack used several
// names: the module itself ("dkg") and its fee pool ("cdr-fee-pool").
func mentionsDKG(name string) bool {
	lower := strings.ToLower(name)

	return strings.Contains(lower, "dkg") || strings.Contains(lower, "cdr")
}

func TestModuleTablesHaveNoDKGEntry(t *testing.T) {
	t.Parallel()

	tables := map[string][]string{
		"genesisModuleOrder": genesisModuleOrder,
		"preBlockers":        preBlockers,
		"beginBlockers":      beginBlockers,
		"endBlockers":        endBlockers,
		"blockAccAddrs":      blockAccAddrs,
	}

	for table, names := range tables {
		require.NotEmpty(t, names, "%s is empty, so this test would assert nothing", table)

		for _, name := range names {
			require.Falsef(t, mentionsDKG(name), "%s still lists %q", table, name)
		}
	}

	require.NotEmpty(t, moduleAccPerms, "moduleAccPerms is empty, so this test would assert nothing")

	for _, perm := range moduleAccPerms {
		require.Falsef(t, mentionsDKG(perm.GetAccount()), "moduleAccPerms still grants %q", perm.GetAccount())
	}
}

// TestUpgradeAndForkListsHaveNoDKGEntry guards the lists the module tables do
// not cover. A partial revert that put the module back into the depinject
// config or the proposal router while leaving the order tables alone would
// otherwise go unnoticed.
func TestUpgradeAndForkListsHaveNoDKGEntry(t *testing.T) {
	t.Parallel()

	require.NotEmpty(t, Upgrades)
	require.NotEmpty(t, Forks)

	// v1.6.0 and v1.8.1 name the retired module in their own names, so only
	// the handlers that are not about retiring it are checked here.
	retirement := map[string]bool{"v1.6.0": true, "v1.8.1": true}

	for _, upgrade := range Upgrades {
		if retirement[upgrade.UpgradeName] {
			continue
		}

		require.Falsef(t, mentionsDKG(upgrade.UpgradeName), "Upgrades still lists %q", upgrade.UpgradeName)
		require.Emptyf(t, upgrade.StoreUpgrades.Added, "%s adds a store", upgrade.UpgradeName)
	}

	for _, fork := range Forks {
		require.Falsef(t, mentionsDKG(fork.UpgradeName), "Forks still lists %q", fork.UpgradeName)
	}
}

// requiredHandlers lists the upgrade names a node may still be asked for.
//
// x/upgrade refuses to start when it has no handler for the completed upgrade
// with the greatest height, and checks on the first block a process handles
// unless that block is exactly a due plan height. Dropping a name from this
// list makes nodes unable to start; neuter a historical upgrade when its body
// stops applying, do not delete it.
//
// v1.6.0 is deliberately absent — see TestV160HandlerIsDeliberatelyAbsent.
//
//nolint:gochecknoglobals // compile-time table
var requiredHandlers = []string{
	"virgil",
	"v1.2.0",
	"polybius",
	"terence",
	"horace",
	"v1.8.1",
}

func TestEveryCompletedUpgradeStillHasAHandler(t *testing.T) {
	t.Parallel()

	registered := make(map[string]bool, len(Upgrades))
	for _, upgrade := range Upgrades {
		registered[upgrade.UpgradeName] = true
	}

	for _, name := range requiredHandlers {
		require.Truef(t, registered[name],
			"no handler for %q: a chain that completed it cannot start this binary", name)
	}
}

// TestV160HandlerIsDeliberatelyAbsent pins the one exception to the rule above.
//
// v1.6.0 activated the module this release retires, so a chain that completed
// it is a chain whose history this binary cannot replay: every block it
// produced carries a DKG vote message that is no longer in the interface
// registry, and executing one skips the execution payload in the same
// transaction. Refusing to start on such a state is the correct outcome, and
// the missing handler is what produces that refusal — before the node commits
// anything.
//
// Registering a no-op handler here would let a lagging or prematurely swapped
// node start, execute a pre-upgrade block, diverge on app hash and need a
// rollback. The legitimate swap never consults this handler: App.PreBlocker
// schedules the plan before the module manager runs, so at the upgrade height
// the plan is due and x/upgrade skips the check entirely.
func TestV160HandlerIsDeliberatelyAbsent(t *testing.T) {
	t.Parallel()

	for _, upgrade := range Upgrades {
		require.NotEqualf(t, "v1.6.0", upgrade.UpgradeName,
			"a v1.6.0 handler lets a node start on a state it cannot execute; see this test's comment")
	}
}

func TestV181UpgradeIsRegistered(t *testing.T) {
	t.Parallel()

	var found bool

	for _, upgrade := range Upgrades {
		if upgrade.UpgradeName == netconf.V181 {
			found = true

			require.NotNil(t, upgrade.CreateUpgradeHandler, "v1.8.1 has no handler")

			// The store is retired by not being mounted, not by a deletion.
			// See the comment on v_1_8_1.Upgrade for why a deletion cannot be
			// used here.
			require.Empty(t, upgrade.StoreUpgrades.Added, "v1.8.1 adds a store")
			require.Empty(t, upgrade.StoreUpgrades.Deleted, "v1.8.1 deletes a store")
			require.Empty(t, upgrade.StoreUpgrades.Renamed, "v1.8.1 renames a store")
		}
	}

	require.True(t, found, "v1.8.1 is not registered in Upgrades, so no binary can run a chain that applied it")
	require.Equal(t, netconf.V181, v_1_8_1.UpgradeName, "upgrade name drifted from netconf")
}

// The fork is what activates the upgrade, so it has to be registered, and it
// has to resolve a height on the network that ran DKG and on no other. The
// app's resolver hands every unknown name a height for every chain through its
// default branch, so without a case of its own this fork would fire everywhere.
func TestV181ForkIsRegisteredAndScopedToTheRightChains(t *testing.T) {
	t.Parallel()

	var found bool

	for _, fork := range Forks {
		if fork.UpgradeName == netconf.V181 {
			found = true

			require.NotNil(t, fork.BeginForkLogic, "the fork has no begin logic")
		}
	}

	require.True(t, found, "v1.8.1 is not registered in Forks, so no height can activate it")

	ctx := sdk.Context{}

	// Only the chains that carry a height take part.
	for _, chainID := range []string{netconf.AeneidChainID, netconf.StoryLocalnetID} {
		height, ok := GetUpgradeHeight(ctx.WithChainID(chainID), netconf.V181, 0)
		require.Truef(t, ok, "%s does not resolve a v1.8.1 fork height", chainID)
		require.Positivef(t, height, "%s resolved a non-positive fork height", chainID)
	}

	// Mainnet never ran DKG, so it must not resolve a height at all -- a
	// fallback of zero would hide this, so pass a height that would fire.
	height, ok := GetUpgradeHeight(ctx.WithChainID(netconf.StoryChainID), netconf.V181, 1)
	require.False(t, ok, "mainnet resolved a v1.8.1 fork height")
	require.Zero(t, height, "mainnet resolved a v1.8.1 fork height")
}
