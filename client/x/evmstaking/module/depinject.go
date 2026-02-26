package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/piplabs/story/client/api/story/evmstaking/v1/module"
	"github.com/piplabs/story/client/x/evmstaking/keeper"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/ethclient"
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

	Config                *module.Module
	ValidatorAddressCodec runtime.ValidatorAddressCodec
	EthClient             ethclient.Client
	AccountKeeper         types.AccountKeeper
	BankKeeper            types.BankKeeper
	SlashingKeeper        types.SlashingKeeper
	StakingKeeper         types.StakingKeeper
	DistributionKeeper    types.DistributionKeeper
	DKGKeeper             types.DKGKeeper
	Cdc                   codec.Codec
	StoreService          store.KVStoreService
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
		in.BankKeeper,
		in.SlashingKeeper,
		in.StakingKeeper,
		in.DistributionKeeper,
		in.DKGKeeper,
		authority.String(),
		in.EthClient,
		in.ValidatorAddressCodec,
	)

	m := NewAppModule(in.Cdc, k, in.AccountKeeper, in.BankKeeper, in.SlashingKeeper, &in.StakingKeeper)

	return ModuleOutputs{Keeper: k, Module: m}
}
