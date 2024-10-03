package app

import (
	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	distrmodulev1 "cosmossdk.io/api/cosmos/distribution/module/v1"
	genutilmodulev1 "cosmossdk.io/api/cosmos/genutil/module/v1"
	govmodulev1 "cosmossdk.io/api/cosmos/gov/module/v1"
	mintmodulev1 "cosmossdk.io/api/cosmos/mint/module/v1"
	slashingmodulev1 "cosmossdk.io/api/cosmos/slashing/module/v1"
	stakingmodulev1 "cosmossdk.io/api/cosmos/staking/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	upgrademodulev1 "cosmossdk.io/api/cosmos/upgrade/module/v1"
	"cosmossdk.io/core/appconfig"
	"cosmossdk.io/depinject"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	epochsmodule "github.com/piplabs/story/client/x/epochs/module"
	epochstypes "github.com/piplabs/story/client/x/epochs/types"
	evmenginemodule "github.com/piplabs/story/client/x/evmengine/module"
	evmenginetypes "github.com/piplabs/story/client/x/evmengine/types"
	evmstakingmodule "github.com/piplabs/story/client/x/evmstaking/module"
	evmstakingtypes "github.com/piplabs/story/client/x/evmstaking/types"
	minttypes "github.com/piplabs/story/client/x/mint/types"
)

// Bech32HRP is the human-readable-part of the Bech32 address format.
const (
	Bech32HRP = "story"

	defaultPruningKeep     = 72_000 // Keep 1 day's of application state by default (given period of 1.2s).
	defaultPruningInterval = 300    // Prune every 5 minutes or so.
)

// init initializes the Cosmos SDK configuration.
//
//nolint:gochecknoinits // Cosmos-style
func init() {
	// Set prefixes
	accountPubKeyPrefix := Bech32HRP + "pub"
	validatorAddressPrefix := Bech32HRP + "valoper"
	validatorPubKeyPrefix := Bech32HRP + "valoperpub"
	consNodeAddressPrefix := Bech32HRP + "valcons"
	consNodePubKeyPrefix := Bech32HRP + "valconspub"

	// Set and seal config
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(Bech32HRP, accountPubKeyPrefix)
	cfg.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	cfg.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	cfg.Seal()
}

// DepConfig returns the default app depinject config.
func DepConfig() depinject.Config {
	return depinject.Configs(
		appConfig,
		depinject.Supply(),
	)
}

//nolint:gochecknoglobals // Cosmos-style
var (
	genesisModuleOrder = []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		upgradetypes.ModuleName,
		// Story modules
		epochstypes.ModuleName,
		evmenginetypes.ModuleName,
		evmstakingtypes.ModuleName,
	}

	// NOTE: upgrade module must come first, as upgrades might break state schema.
	preBlockers = []string{
		upgradetypes.ModuleName,
	}

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0.
	beginBlockers = []string{
		epochstypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName, // Note: slashing happens after distr.BeginBlocker
		slashingtypes.ModuleName,
		stakingtypes.ModuleName,
	}

	endBlockers = []string{
		govtypes.ModuleName,
		evmstakingtypes.ModuleName, // Must be before staking module removes mature unbonding delegations & validators.
	}

	// blocked account addresses.
	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
		evmstakingtypes.ModuleName,
		epochstypes.ModuleName,
	}

	moduleAccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: authtypes.FeeCollectorName},
		{Account: minttypes.ModuleName, Permissions: []string{authtypes.Minter}},
		{Account: distrtypes.ModuleName},
		{Account: stakingtypes.BondedPoolName, Permissions: []string{authtypes.Burner, authtypes.Staking}},
		{Account: stakingtypes.NotBondedPoolName, Permissions: []string{authtypes.Burner, authtypes.Staking}},
		{Account: evmstakingtypes.ModuleName, Permissions: []string{authtypes.Burner, authtypes.Minter}},
		{Account: govtypes.ModuleName, Permissions: []string{authtypes.Burner}},
	}

	// appConfig application configuration (used by depinject).
	appConfig = appconfig.Compose(&appv1alpha1.Config{
		Modules: []*appv1alpha1.ModuleConfig{
			{
				Name: runtime.ModuleName,
				Config: appconfig.WrapAny(&runtimev1alpha1.Module{
					AppName: Name,
					// NOTE: "PreBlockers" is set in app.go to override the ABCI++ PreBlocker.
					BeginBlockers: beginBlockers,
					// NOTE: "EndBlockers" is set in app.go since evmstaking endblocker replaces the staking endblocker.
					InitGenesis: genesisModuleOrder,
					OverrideStoreKeys: []*runtimev1alpha1.StoreKeyConfig{
						{
							ModuleName: authtypes.ModuleName,
							KvStoreKey: "acc",
						},
					},
				}),
			},
			{
				Name: authtypes.ModuleName,
				Config: appconfig.WrapAny(&authmodulev1.Module{
					ModuleAccountPermissions: moduleAccPerms,
					Bech32Prefix:             Bech32HRP,
				}),
			},
			{
				Name: "tx",
				Config: appconfig.WrapAny(&txconfigv1.Config{
					SkipAnteHandler: true, // Disable ante handler (since we don't have proper txs).
					SkipPostHandler: true,
				}),
			},
			{
				Name: banktypes.ModuleName,
				Config: appconfig.WrapAny(&bankmodulev1.Module{
					BlockedModuleAccountsOverride: blockAccAddrs,
				}),
			},
			{
				Name:   consensustypes.ModuleName,
				Config: appconfig.WrapAny(&consensusmodulev1.Module{}),
			},
			{
				Name:   distrtypes.ModuleName,
				Config: appconfig.WrapAny(&distrmodulev1.Module{}),
			},
			{
				Name:   slashingtypes.ModuleName,
				Config: appconfig.WrapAny(&slashingmodulev1.Module{}),
			},
			{
				Name:   genutiltypes.ModuleName,
				Config: appconfig.WrapAny(&genutilmodulev1.Module{}),
			},
			{
				Name:   govtypes.ModuleName,
				Config: appconfig.WrapAny(&govmodulev1.Module{}),
			},
			{
				Name:   stakingtypes.ModuleName,
				Config: appconfig.WrapAny(&stakingmodulev1.Module{}),
			},
			{
				Name:   upgradetypes.ModuleName,
				Config: appconfig.WrapAny(&upgrademodulev1.Module{}),
			},
			{
				Name:   epochstypes.ModuleName,
				Config: appconfig.WrapAny(&epochsmodule.Module{}),
			},
			{
				Name:   evmstakingtypes.ModuleName,
				Config: appconfig.WrapAny(&evmstakingmodule.Module{}),
			},
			{
				Name:   evmenginetypes.ModuleName,
				Config: appconfig.WrapAny(&evmenginemodule.Module{}),
			},
			{
				Name:   minttypes.ModuleName,
				Config: appconfig.WrapAny(&mintmodulev1.Module{}),
			},
		},
	})
)
