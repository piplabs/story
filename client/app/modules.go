package app

import (
	"fmt"

	sdkmodule "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/piplabs/story/client/app/module"
	evmenginemodule "github.com/piplabs/story/client/x/evmengine/module"
	evmenginetypes "github.com/piplabs/story/client/x/evmengine/types"
	evmstakingmodule "github.com/piplabs/story/client/x/evmstaking/module"
	evmstakingtypes "github.com/piplabs/story/client/x/evmstaking/types"
	mint "github.com/piplabs/story/client/x/mint/module"
	minttypes "github.com/piplabs/story/client/x/mint/types"
	signalmodule "github.com/piplabs/story/client/x/signal/module"
	signaltypes "github.com/piplabs/story/client/x/signal/types"
)

var (
	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration and genesis verification.
	ModuleBasics = sdkmodule.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distribution.AppModuleBasic{},
		gov.AppModuleBasic{},
		params.AppModuleBasic{},
		slashing.AppModule{},
		vesting.AppModuleBasic{},
		// Story modules
		evmenginemodule.AppModuleBasic{},
		evmstakingmodule.AppModuleBasic{},
		signalmodule.AppModuleBasic{},
	)

	// ModuleEncodingRegisters keeps track of all the module methods needed to
	// register interfaces and specific type to encoding config.
	ModuleEncodingRegisters = extractRegisters(ModuleBasics)
)

func (a *App) setupModuleManager() error {
	//evmstakingStakingkeeper := a.Keepers.StakingKeeper.(*evmstakingtypes.StakingKeeper)
	//signalStakingkeeper := a.Keepers.StakingKeeper.(*signaltypes.StakingKeeper)

	var err error
	a.ModuleManager, err = module.NewManager([]module.VersionedModule{
		{
			Module:      genutil.NewAppModule(a.Keepers.AccountKeeper, a.Keepers.StakingKeeper, a, a.txConfig),
			FromVersion: v1, ToVersion: v1,
		},
		{
			Module:      auth.NewAppModule(a.appCodec, a.Keepers.AccountKeeper, nil, nil),
			FromVersion: v1, ToVersion: v1,
		},
		{
			Module:      vesting.NewAppModule(a.Keepers.AccountKeeper, a.Keepers.BankKeeper),
			FromVersion: v1, ToVersion: v1,
		},
		{
			Module:      bank.NewAppModule(a.appCodec, a.Keepers.BankKeeper, a.Keepers.AccountKeeper, nil),
			FromVersion: v1, ToVersion: v1,
		},
		{
			Module:      gov.NewAppModule(a.appCodec, a.Keepers.GovKeeper, a.Keepers.AccountKeeper, a.Keepers.BankKeeper, nil),
			FromVersion: v1, ToVersion: v1,
		},
		{
			Module:      slashing.NewAppModule(a.appCodec, a.Keepers.SlashingKeeper, a.Keepers.AccountKeeper, a.Keepers.BankKeeper, a.Keepers.StakingKeeper, nil, nil),
			FromVersion: v1, ToVersion: v1,
		},
		{
			Module:      distribution.NewAppModule(a.appCodec, a.Keepers.DistrKeeper, a.Keepers.AccountKeeper, a.Keepers.BankKeeper, a.Keepers.StakingKeeper, nil),
			FromVersion: v1, ToVersion: v1,
		},
		{
			Module:      staking.NewAppModule(a.appCodec, a.Keepers.StakingKeeper, a.Keepers.AccountKeeper, a.Keepers.BankKeeper, nil),
			FromVersion: v1, ToVersion: v1,
		},
		// Story modules
		{
			Module:      mint.NewAppModule(a.appCodec, a.Keepers.MintKeeper, a.Keepers.AccountKeeper, nil),
			FromVersion: v1, ToVersion: v1,
		},
		{
			Module:      evmenginemodule.NewAppModule(a.appCodec, a.Keepers.EVMEngKeeper),
			FromVersion: v1, ToVersion: v1,
		},
		//{
		//	Module:      evmstakingmodule.NewAppModule(a.appCodec, a.Keepers.EvmStakingKeeper, a.Keepers.AccountKeeper, a.Keepers.BankKeeper, a.Keepers.SlashingKeeper, evmstakingStakingkeeper),
		//	FromVersion: v1, ToVersion: v1,
		//},
		//{
		//	Module:      signalmodule.NewAppModule(a.appCodec, &a.Keepers.SignalKeeper, a.Keepers.AccountKeeper, signalStakingkeeper),
		//	FromVersion: v1, ToVersion: v1,
		//},
	})
	if err != nil {
		return err
	}
	return a.ModuleManager.AssertMatchingModules(ModuleBasics)
}

func (a *App) setModuleOrder() {
	// Set "OrderEndBlockers" directly instead of using "SetOrderEndBlockers," which will panic since the staking module
	// is missing in the "endBlockers", which is an intended behavior in Story. The panic message is:
	// `panic: all modules must be defined when setting SetOrderEndBlockers, missing: [staking]`
	a.ModuleManager.OrderEndBlockers = endBlockers
	a.SetEndBlocker(a.EndBlocker)
}

func allStoreKeys() []string {
	return []string{
		authtypes.StoreKey,
		banktypes.StoreKey,
		distrtypes.StoreKey,
		govtypes.StoreKey,
		minttypes.StoreKey,
		paramstypes.StoreKey,
		slashingtypes.StoreKey,
		stakingtypes.StoreKey,
		// Story modules
		evmenginetypes.StoreKey,
		evmstakingtypes.StoreKey,
		signaltypes.StoreKey,
	}
}

// versionedStoreKeys returns the store keys for each app version.
func versionedStoreKeys() map[uint64][]string {
	return map[uint64][]string{
		1: {
			authtypes.StoreKey,
			banktypes.StoreKey,
			distrtypes.StoreKey,
			govtypes.StoreKey,
			minttypes.StoreKey,
			paramstypes.StoreKey,
			slashingtypes.StoreKey,
			stakingtypes.StoreKey,
			// Story modules
			evmenginetypes.StoreKey,
			evmstakingtypes.StoreKey,
			signaltypes.StoreKey,
		},
	}
}

// assertAllKeysArePresent performs a couple sanity checks on startup to ensure each versions key names have
// a key and that all versions supported by the module manager have a respective versioned key.
func (a *App) assertAllKeysArePresent() {
	supportedAppVersions := a.SupportedVersions()
	supportedVersionsMap := make(map[uint64]bool, len(supportedAppVersions))
	for _, version := range supportedAppVersions {
		supportedVersionsMap[version] = false
	}

	for appVersion, keys := range a.keyVersions {
		if _, exists := supportedVersionsMap[appVersion]; exists {
			supportedVersionsMap[appVersion] = true
		} else {
			panic(fmt.Sprintf("keys %v for app version %d are not supported by the module manager", keys, appVersion))
		}
		for _, key := range keys {
			if _, ok := a.keys[key]; !ok {
				panic(fmt.Sprintf("key %s is not present", key))
			}
		}
	}
	for appVersion, supported := range supportedVersionsMap {
		if !supported {
			panic(fmt.Sprintf("app version %d is supported by the module manager but has no keys", appVersion))
		}
	}
}

// extractRegisters returns the encoding module registers from the basic manager.
func extractRegisters(manager sdkmodule.BasicManager) (modules []ModuleRegister) {
	for _, m := range manager {
		modules = append(modules, m)
	}
	return modules
}
