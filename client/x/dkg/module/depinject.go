package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/piplabs/story/client/api/story/dkg/v1/module"
	"github.com/piplabs/story/client/x/dkg/keeper"
	"github.com/piplabs/story/client/x/dkg/types"
)

//nolint:gochecknoinits // depinject
func init() {
	appmodule.Register(
		&module.Module{},
		appmodule.Provide(
			ProvideModule,
		),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *module.Module
	Cdc          codec.Codec
	StoreService store.KVStoreService

	AccountKeeper types.AccountKeeper
	StakingKeeper types.StakingKeeper
	ValStore      *skeeper.Keeper

	// These two clients can be nil depending on whether the DKG service is enabled.
	DKGTEEClient      types.TEEClient        `optional:"true"`
	DKGContractClient *keeper.ContractClient `optional:"true"`
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
		in.ValStore,
		in.DKGTEEClient,
		in.DKGContractClient,
		authority.String(),
	)

	m := NewAppModule(in.Cdc, &k)

	return ModuleOutputs{Keeper: &k, Module: m}
}
