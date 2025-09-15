package module

import (
	"fmt"
	"sort"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/piplabs/story/client/api/story/mint/v1/module"
	"github.com/piplabs/story/client/x/epochs/keeper"
	"github.com/piplabs/story/client/x/epochs/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

//nolint:gochecknoinits // Cosmos-style
func init() {
	appconfig.RegisterModule(&module.Module{},
		appconfig.Provide(ProvideModule),
		appconfig.Invoke(InvokeSetHooks),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *module.Module
	Cdc          codec.Codec
	StoreService store.KVStoreService
}

type ModuleOutputs struct {
	depinject.Out

	EpochKeeper keeper.Keeper
	Module      appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(in.StoreService, in.Cdc)
	m := NewAppModule(k)

	return ModuleOutputs{EpochKeeper: k, Module: m}
}

func InvokeSetHooks(keeper *keeper.Keeper, hooks map[string]types.EpochHooksWrapper) error {
	if keeper == nil || hooks == nil {
		return nil
	}

	// Default ordering is lexical by module name.
	// Explicit ordering can be added to the module config if required.
	var modNames []string
	for modName := range hooks {
		modNames = append(modNames, modName)
	}
	sort.Strings(modNames)
	var multiHooks types.MultiEpochHooks
	for _, modName := range modNames {
		hook, ok := hooks[modName]
		if !ok {
			return fmt.Errorf("can't find epoch hooks for module %s", modName)
		}
		multiHooks = append(multiHooks, hook)
	}

	keeper.SetHooks(multiHooks)

	return nil
}
