package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/piplabs/story/client/x/signal/keeper"
	"github.com/piplabs/story/client/x/signal/types"
)

//nolint:gochecknoinits // depinject
func init() {
	appmodule.Register(
		&Module{},
		appmodule.Provide(
			ProvideModule,
		),
	)
}

type ModuleInputs struct {
	depinject.In

	Config        *Module
	AccountKeeper types.AccountKeeper
	StakingKeeper types.StakingKeeper
	Cdc           codec.Codec
	StoreService  store.KVStoreService
}

type ModuleOutputs struct {
	depinject.Out

	Keeper *keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.GetAuthority() != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.GetAuthority())
	}

	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.AccountKeeper,
		in.StakingKeeper,
		authority.String(),
	)

	m := NewAppModule(in.Cdc, k, in.AccountKeeper, &in.StakingKeeper)

	return ModuleOutputs{Keeper: k, Module: m}
}
