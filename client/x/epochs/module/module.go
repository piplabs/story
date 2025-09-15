package module

import (
	"context"
	"encoding/json"

	"cosmossdk.io/core/appmodule"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/piplabs/story/client/x/epochs/keeper"
	"github.com/piplabs/story/client/x/epochs/simulation"
	"github.com/piplabs/story/client/x/epochs/types"
	"github.com/piplabs/story/lib/errors"

	"google.golang.org/grpc"
)

var (
	_ module.AppModuleSimulation = AppModule{}
	_ module.HasGenesis          = AppModule{}

	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
)

const ConsensusVersion = 1

// AppModule implements the AppModule interface for the epochs module.
type AppModule struct {
	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object.
func NewAppModule(keeper keeper.Keeper) AppModule {
	return AppModule{
		keeper: keeper,
	}
}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModule) IsAppModule() {}

// Name returns the epochs module's name.
// Deprecated: kept for legacy reasons.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the epochs module's types for the given codec.
func (AppModule) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

func (AppModule) RegisterInterfaces(_ codectypes.InterfaceRegistry) {}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the epochs module.
func (AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(registrar grpc.ServiceRegistrar) error {
	types.RegisterQueryServer(registrar, keeper.NewQuerier(am.keeper))

	return nil
}

// DefaultGenesis returns the epochs module's default genesis state.
func (AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	data, err := cdc.MarshalJSON(types.DefaultGenesis())
	if err != nil {
		panic(err)
	}

	return data
}

// ValidateGenesis performs genesis state validation for the epochs module.
func (AppModule) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return errors.Wrap(err, "failed to unmarshal "+types.ModuleName+" genesis state")
	}

	return gs.Validate()
}

// InitGenesis performs the epochs module's genesis initialization.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) {
	var gs types.GenesisState
	err := cdc.UnmarshalJSON(bz, &gs)
	if err != nil {
		panic(errors.Wrap(err, "failed to unmarshal "+types.ModuleName+" genesis state"))
	}

	if err := am.keeper.InitGenesis(ctx, gs); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the epochs module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs, err := am.keeper.ExportGenesis(ctx)
	if err != nil {
		panic(err)
	}

	bz, err := cdc.MarshalJSON(gs)
	if err != nil {
		panic(err)
	}

	return bz
}

// ConsensusVersion implements HasConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// BeginBlock executes all ABCI BeginBlock logic respective to the epochs module.
func (am AppModule) BeginBlock(ctx context.Context) error {
	return am.keeper.BeginBlocker(ctx)
}

// AppModuleSimulation functions

// WeightedOperations is a no-op.
func (AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}

// GenerateGenesisState creates a randomized GenState of the epochs module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder for epochs module's types.
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {
	sdr[types.StoreKey] = simtypes.NewStoreDecoderFuncFromCollectionsSchema(am.keeper.Schema)
}

// TODO add when we have collections full support with schema
/*
// ModuleCodec implements schema.HasModuleCodec.
// It allows the indexer to decode the module's KVPairUpdate.
func (am AppModule) ModuleCodec() (schema.ModuleCodec, error) {
	return am.keeper.Schema.ModuleCodec(collections.IndexingOptions{})
}
*/
