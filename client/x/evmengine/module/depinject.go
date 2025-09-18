package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/piplabs/story/client/api/story/evmengine/v1/module"
	"github.com/piplabs/story/client/x/evmengine/keeper"
	"github.com/piplabs/story/client/x/evmengine/types"
	"github.com/piplabs/story/lib/ethclient"
)

//nolint:gochecknoinits // Cosmos-style
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

	StoreService     store.KVStoreService
	Cdc              codec.Codec
	Config           *module.Module
	TXConfig         client.TxConfig
	EngineCl         ethclient.EngineClient
	EthCl            ethclient.Client
	AccountKeeper    types.AccountKeeper
	EvmStakingKeeper types.EvmStakingKeeper
	UpgradeKeeper    types.UpgradeKeeper
	DistrKeeper      types.DistrKeeper
	DKGKeeper        types.DKGKeeper
}

type ModuleOutputs struct {
	depinject.Out

	EngEVMKeeper *keeper.Keeper
	Module       appmodule.AppModule
	Hooks        stakingtypes.StakingHooksWrapper
}

func ProvideModule(in ModuleInputs) (ModuleOutputs, error) {
	k, err := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.EngineCl,
		in.EthCl,
		in.TXConfig,
		in.AccountKeeper,
		in.EvmStakingKeeper,
		in.UpgradeKeeper,
		in.DistrKeeper,
		in.DKGKeeper,
	)
	if err != nil {
		return ModuleOutputs{}, err
	}

	m := NewAppModule(
		in.Cdc,
		k,
	)

	return ModuleOutputs{
		EngEVMKeeper: k,
		Module:       m,
		Hooks:        stakingtypes.StakingHooksWrapper{StakingHooks: keeper.Hooks{}},
	}, nil
}
