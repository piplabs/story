package module

import (
	"fmt"
	"sort"

	"golang.org/x/exp/maps"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/piplabs/story/client/x/epochs/keeper"
	"github.com/piplabs/story/client/x/epochs/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

//nolint:gochecknoinits // Cosmos-style
func init() {
	appmodule.Register(
		&Module{},
		appmodule.Provide(ProvideModule),
		appmodule.Invoke(InvokeSetHooks),
	)
}

type ModuleInputs struct {
	depinject.In

	StoreService store.KVStoreService
	Cdc          codec.Codec
	Config       *Module
	TXConfig     client.TxConfig

	EventService event.Service
}

type ModuleOutputs struct {
	depinject.Out

	EpochKeeper *keeper.Keeper
	Module      appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(in.StoreService, in.EventService, in.Cdc)
	m := NewAppModule(in.Cdc, k)

	return ModuleOutputs{EpochKeeper: k, Module: m}
}

func InvokeSetHooks(keeper *keeper.Keeper, hooks map[string]types.EpochHooksWrapper) error {
	if hooks == nil {
		return nil
	}

	// Default ordering is lexical by module name.
	// Explicit ordering can be added to the module config if required.
	modNames := maps.Keys(hooks)
	order := modNames
	sort.Strings(order)

	var multiHooks types.MultiEpochHooks
	for _, modName := range order {
		hook, ok := hooks[modName]
		if !ok {
			return fmt.Errorf("can't find epoch hooks for module %s", modName)
		}
		multiHooks = append(multiHooks, hook)
	}

	keeper.SetHooks(multiHooks)

	return nil
}
