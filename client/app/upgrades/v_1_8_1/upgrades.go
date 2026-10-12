//nolint:staticcheck,revive,nolintlint // versioning
package v_1_8_1

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// Module accounts the retired stack owned. They are spelled out here because
// their types package is gone from this binary.
const (
	dkgModuleName  = "dkg"
	cdrFeePoolName = "cdr-fee-pool"
)

// CreateUpgradeHandler retires the DKG module and the CDR stack.
//
// Everything that made the module part of the state machine is gone from this
// binary already: the store is unmounted, its BeginBlocker is unregistered, the
// vote extension handlers fall back to baseapp's no-ops, and the proposal
// router no longer expects a DKG vote in each block. What is left for the
// handler is the balance the stack held, which no code can reach once the
// module is gone.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	k *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return createHandler(mm, configurator, k.BankKeeper, k.DistrKeeper)
}

// createHandler holds the body, taking the narrow keeper slices it needs so the
// sweep can be exercised with stated balances.
func createHandler(
	mm *module.Manager,
	configurator module.Configurator,
	bank BankKeeper,
	distr DistributionKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		log.Info(ctx, "Start v1.8.1 upgrade — retiring DKG module and CDR stack")

		if err := sweepRetiredPools(ctx, bank, distr); err != nil {
			return vm, err
		}

		// RunMigrations iterates this binary's modules rather than the stored
		// map, so the returned map describes this binary's module set and the
		// dkg entry is not carried into it. Every surviving module has
		// fromVersion == toVersion, which makes each migration a no-op, and no
		// module in the manager is missing from the stored map — the module
		// list only shrank — so nothing takes the InitGenesis branch. That
		// precondition stops holding the day a later release adds a module
		// while this handler is still the one being applied.
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, errors.Wrap(err, "run migrations")
		}

		// This does not erase the stored entry. SetModuleVersionMap only writes
		// the keys it is given, so "dkg" stays in the upgrade store for good.
		// Nothing reads it — RunMigrations ignores names it has no module for —
		// but it means the dkg module and store names must never be reused on a
		// chain that applied this upgrade: a module reintroduced under the old
		// name would find a stored version and skip InitGenesis, and a store
		// remounted under the old name would load the abandoned IAVL data.
		if _, ok := newVM[dkgModuleName]; ok {
			return vm, errors.New("dkg module is still registered in this binary")
		}

		log.Info(ctx, "V1.8.1 upgrade complete — DKG module and CDR stack retired",
			"modules", len(newVM),
			"stored_map_had_dkg", hasDKG(vm),
		)

		return newVM, nil
	}
}

// sweepRetiredPools moves whatever the retired module accounts still hold into
// the UBI bucket of the fee pool, which is the only pool left that pays out.
//
// The CDR fee pool is the larger of the two and is expected to be non-empty:
// AddCDRFeeToPool minted a mirror of each fee into it, and distribution left
// the truncation remainder behind on purpose, so a residue accumulates between
// round finalizations. The dkg account itself is normally swept every block and
// should be empty, but it is included so nothing is left stranded.
//
// Moving these to the fee pool creates no supply: the matching IP was burned on
// the execution layer when the fee was collected — CDR's _collectFee sends it
// to the zero address — so the minted coins are that burn's only representation
// on this layer. Leaving them in an unreachable module account would strand
// them. Where they end up is worth stating: FundUbi credits feePool.Ubi, and
// once that bucket passes evmstaking's MinPartialWithdrawalAmount,
// ProcessUbiWithdrawal burns it here and queues an execution-layer withdrawal
// to the configured UbiWithdrawAddress. Below that threshold it simply joins
// ordinary UBI.
func sweepRetiredPools(ctx context.Context, bank BankKeeper, distr DistributionKeeper) error {
	for _, name := range []string{cdrFeePoolName, dkgModuleName} {
		addr := authtypes.NewModuleAddress(name)

		balance := bank.GetAllBalances(ctx, addr)
		if balance.IsZero() {
			log.Info(ctx, "Retired module account is empty", "account", name)

			continue
		}

		if err := distr.FundUbi(ctx, balance, addr); err != nil {
			return errors.Wrap(err, "fund fee pool from retired account", "account", name)
		}

		log.Info(ctx, "Moved retired module account balance to the fee pool",
			"account", name,
			"amount", balance.String(),
		)
	}

	return nil
}

// hasDKG reports whether the stored version map still describes the retired
// module, which tells an operator reading the logs whether this chain ever ran
// it.
func hasDKG(vm module.VersionMap) bool {
	_, ok := vm[dkgModuleName]

	return ok
}
