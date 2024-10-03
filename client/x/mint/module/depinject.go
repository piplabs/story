//nolint:revive // just use interface{}
package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/piplabs/story/client/x/mint/keeper"
	"github.com/piplabs/story/client/x/mint/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

//nolint:gochecknoinits // depinject
func init() {
	appmodule.Register(
		&Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	ModuleKey              depinject.OwnModuleKey
	Config                 *Module
	StoreService           store.KVStoreService
	Cdc                    codec.Codec
	InflationCalculationFn types.InflationCalculationFn `optional:"true"`

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
	StakingKeeper types.StakingKeeper
}

type ModuleOutputs struct {
	depinject.Out

	MintKeeper keeper.Keeper
	Module     appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	feeCollectorName := in.Config.GetFeeCollectorName()
	if feeCollectorName == "" {
		feeCollectorName = authtypes.FeeCollectorName
	}

	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.GetAuthority() != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.GetAuthority())
	}

	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.StakingKeeper,
		in.AccountKeeper,
		in.BankKeeper,
		feeCollectorName,
		authority.String(),
	)

	// when no inflation calculation function is provided it will use the default types.DefaultInflationCalculationFn
	m := NewAppModule(in.Cdc, k, in.AccountKeeper, in.InflationCalculationFn)

	return ModuleOutputs{MintKeeper: k, Module: m}
}
