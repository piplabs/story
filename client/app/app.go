package app

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/piplabs/story/client/app/encoding"
	"io"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkmodule "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/app/module"
	"github.com/piplabs/story/client/comet"
	appv1 "github.com/piplabs/story/client/pkg/appconsts/v1"
	evmstakingkeeper "github.com/piplabs/story/client/x/evmstaking/keeper"
	mintkeeper "github.com/piplabs/story/client/x/mint/keeper"
	signalkeeper "github.com/piplabs/story/client/x/signal/keeper"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient"

	_ "cosmossdk.io/api/cosmos/tx/config/v1"          // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth"           // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config" // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/bank"           // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/consensus"      // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/distribution"   // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/genutil"        // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/gov"            // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/slashing"       // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/staking"        // import for side-effects
)

const AppName = "StoryApp"

const (
	v1                    = appv1.Version
	DefaultInitialVersion = v1
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp
	keepers.Keepers

	appCodec          codec.Codec
	legacyAmino       *codec.LegacyAmino
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// keys to access the substores
	keyVersions map[uint64][]string
	keys        map[string]*storetypes.KVStoreKey // keys == storeKeys (in runtime.App)
	tkeys       map[string]*storetypes.TransientStoreKey
	memKeys     map[string]*storetypes.MemoryStoreKey

	// override the runtime baseapp's module manager to use the custom module manager
	manager      *module.Manager
	ModuleBasics module.BasicManager
	sm           *sdkmodule.SimulationManager // TODO: module.SimulationManager
	configurator module.Configurator
	homePath     string
}

// newApp returns a reference to an initialized App.
func newApp(
	logger log.Logger,
	db dbm.DB,
	engineCl ethclient.EngineClient,
	traceStore io.Writer,
	baseAppOpts ...func(*baseapp.BaseApp),
) (*App, error) {
	encodingConfig := encoding.MakeEncodingConfig(ModuleEncodingRegisters...)
	appCodec := encodingConfig.Codec
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	baseApp := baseapp.NewBaseApp(AppName, logger, db, txConfig.TxDecoder(), baseAppOpts...)
	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetVersion(version.Version)
	baseApp.SetInterfaceRegistry(interfaceRegistry)

	// Define what keys will be used in the cosmos-sdk key/value store.
	// Cosmos-SDK modules each have a "key" that allows the application to reference what they've stored on the chain.
	keys := storetypes.NewKVStoreKeys(allStoreKeys()...)
	// Define transient store keys
	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	// MemKeys are for information that is stored only in RAM.
	memKeys := storetypes.NewMemoryStoreKeys()

	app := &App{
		BaseApp:           baseApp,
		appCodec:          appCodec,
		legacyAmino:       cdc,
		interfaceRegistry: interfaceRegistry,
		txConfig:          encodingConfig.TxConfig,
		keyVersions:       versionedStoreKeys(),
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.InitSpecialKeepers(
		appCodec,
		baseApp,
		cdc,
		keys,
		tkeys,
	)
	app.InitNormalKeepers(
		appCodec,
		encodingConfig,
		baseApp,
		moduleAccPerms,
		app.BlockedAddrs(),
		txConfig,
		engineCl,
		keys,
	)

	/****  Module Options ****/

	// TODO: There is a bug here, where we register the govRouter routes in InitNormalKeepers and then
	// call setupHooks afterwards. Therefore, if a gov proposal needs to call a method and that method calls a
	// hook, we will get a nil pointer dereference error due to the hooks in the keeper not being
	// setup yet. I will refrain from creating an issue in the sdk for now until after we unfork to 0.47,
	// because I believe the concept of Routes is going away.
	// https://github.com/osmosis-labs/osmosis/issues/6580
	app.SetupHooks()

	// NOTE: All module / keeper changes should happen prior to this module.NewManager line being called.
	// However in the event any changes do need to happen after this call, ensure that that keeper
	// is only passed in its keeper form (not de-ref'd anywhere)
	//
	// Generally NewAppModule will require the keeper that module defines to be passed in as an exact struct,
	// but should take in every other keeper as long as it matches a certain interface. (So no need to be de-ref'd)
	//
	// Any time a module requires a keeper de-ref'd that's not its native one,
	// its code-smell and should probably change. We should get the staking keeper dependencies fixed.
	//
	// NOTE: Modules can't be modified or else must be passed by reference to the module manager
	if err := app.setupModuleManager(); err != nil {
		panic(err)
	}

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)

	// Tell the app's module manager how to set the order of BeginBlockers, which are run at the beginning of every block.
	app.manager.SetOrderBeginBlockers(beginBlockers...)

	// Tell the app's module manager how to set the order of EndBlockers, which are run at the end of every block.
	// Write directly since we are skipping the staking module, which fails "assertNoForgottenModules" check
	app.manager.OrderEndBlockers = endBlockers

	app.manager.SetOrderInitGenesis(genesisModuleOrder...)

	//app.manager.RegisterInvariants(app.Keepers.CrisisKeeper)

	app.configurator = module.NewConfigurator(app.AppCodec(), app.MsgServiceRouter(), app.GRPCQueryRouter())
	if err := app.manager.RegisterServices(app.configurator); err != nil {
		panic(err)
	}

	// Override the gov ModuleBasic with all the custom proposal handers, otherwise we lose them in the CLI.
	app.ModuleBasics = ModuleBasics

	// Initialize the KV stores for the base modules (e.g. params). The base modules will be included in every app version.
	app.MountKVStores(app.baseKeys()) // versioned keys are mounted in InitChain
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetPrecommiter(app.Precommitter)
	app.SetAnteHandler(nil)

	//app.SetMigrateStoreFn(app.migrateCommitStore)
	//app.SetMigrateModuleFn(app.migrateModules)

	// ABCI handlers
	// Use evm engine to create block proposals.
	// Note that we do not check MaxTxBytes since all EngineEVM transaction MUST be included since we cannot
	// postpone them to the next block. Nit: we could drop some vote extensions though...?
	app.SetPrepareProposal(app.Keepers.EVMEngKeeper.PrepareProposal)

	// Route proposed messages to keepers for verification and external state updates.
	app.SetProcessProposal(makeProcessProposalHandler(makeProcessProposalRouter(app), app.txConfig))

	// assert that keys are present for all supported versions
	app.assertAllKeysArePresent()

	// we don't seal the store until the app version has been initialised
	// this will just initialize the base keys (i.e. the param store)
	if err := app.CommitMultiStore().LoadLatestVersion(); err != nil {
		panic(err)
	}

	return app, nil
}

// PreBlocker application updates before each begin block.
func (a *App) PreBlocker(ctx sdk.Context, _ *abcitypes.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	// Set gas meter to the free gas meter.
	// This is because there is currently non-deterministic gas usage in the
	// pre-blocker, e.g. due to hydration of in-memory data structures.
	//
	// Note that we don't need to reset the gas meter after the pre-blocker
	// because Go is pass by value.
	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	return a.manager.PreBlock(ctx)
}

// BeginBlocker application updates every begin block.
func (a *App) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return a.manager.BeginBlock(ctx)
}

// EndBlocker executes application updates at the end of every block.
func (a *App) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	res, err := a.manager.EndBlock(ctx)
	if err != nil {
		return sdk.EndBlock{}, errors.Wrap(err, "module manager endblocker")
	}

	currentVersion := a.AppVersion()
	if shouldUpgrade, newVersion := a.Keepers.SignalKeeper.ShouldUpgrade(ctx); shouldUpgrade {
		// Version changes must be increasing. Downgrades are not permitted
		if newVersion > currentVersion {
			a.SetProtocolVersion(newVersion)
			// reset tally
		}
	}

	return res, nil
}

// Precommitter application updates before the commital of a block after all transactions have been delivered.
func (a *App) Precommitter(ctx sdk.Context) {
	if err := a.manager.Precommit(ctx); err != nil {
		panic(err)
	}
}

func (App) ExportAppStateAndValidators(_ bool, _, _ []string) (servertypes.ExportedApp, error) {
	return servertypes.ExportedApp{}, errors.New("not implemented")
}

func (App) SimulationManager() *sdkmodule.SimulationManager {
	return nil
}

// SetCometAPI sets the comet API client.
// TODO: Figure out how to use depinject to set this.
func (a App) SetCometAPI(api comet.API) {
	a.Keepers.EVMEngKeeper.SetCometAPI(api)
}

// InitChain implements the ABCI interface. This method is a wrapper around
// baseapp's InitChain so we can take the app version and setup the multicommit
// store.
//
// Side-effect: calls baseapp.Init().
func (a *App) InitChain(req *abcitypes.RequestInitChain) (*abcitypes.ResponseInitChain, error) {
	mreq := setDefaultAppVersion(*req)
	appVersion := mreq.ConsensusParams.Version.App
	// mount the stores for the provided app version if it has not already been mounted
	if a.AppVersion() == 0 && !a.IsSealed() {
		a.mountKeysAndInit(appVersion)
		a.SetProtocolVersion(v1)
	}

	resp, err := a.BaseApp.InitChain(&mreq)
	if err != nil {
		return nil, errors.Wrap(err, "init chain")
	}

	//if appVersion != v1 {
	//	//a.SetInitialAppVersionInConsensusParams(ctx, appVersion)
	//	a.SetProtocolVersion(appVersion)
	//}

	return resp, nil
}

// setDefaultAppVersion sets the default app version in the consensus params if
// it was 0. This is needed because chains (e.x. mocha-4) did not explicitly set
// an app version in genesis.json.
func setDefaultAppVersion(req abcitypes.RequestInitChain) abcitypes.RequestInitChain {
	if req.ConsensusParams == nil {
		panic("no consensus params set")
	}
	if req.ConsensusParams.Version == nil {
		panic("no version set in consensus params")
	}
	if req.ConsensusParams.Version.App == 0 {
		req.ConsensusParams.Version.App = v1
	}

	return req
}

// mountKeysAndInit mounts the keys for the provided app version and then
// invokes baseapp.Init().
func (a *App) mountKeysAndInit(appVersion uint64) {
	a.Logger().Info(fmt.Sprintf("mounting KV stores for app version %v", appVersion))
	a.MountKVStores(a.versionedKeys(appVersion))

	// Invoke load latest version for its side-effect of invoking baseapp.Init()
	if err := a.LoadLatestVersion(); err != nil {
		panic(fmt.Sprintf("loading latest version: %s", err.Error()))
	}
}

// InitChainer is middleware that gets invoked part-way through the baseapp's InitChain invocation.
func (a *App) InitChainer(ctx sdk.Context, req *abcitypes.RequestInitChain) (*abcitypes.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := cmtjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	appVersion := req.ConsensusParams.Version.App

	return a.manager.InitGenesis(ctx, a.appCodec, genesisState, appVersion)
}

// LoadHeight loads a particular height.
func (a *App) LoadHeight(height int64) error {
	return a.LoadVersion(height)
}

// SupportedVersions returns all the state machines that the
// application supports.
func (a *App) SupportedVersions() []uint64 {
	return a.manager.SupportedVersions()
}

// versionedKeys returns a map from moduleName to KV store key for the given app version.
func (a *App) versionedKeys(appVersion uint64) map[string]*storetypes.KVStoreKey {
	output := make(map[string]*storetypes.KVStoreKey)
	if keys, exists := a.keyVersions[appVersion]; exists {
		for _, moduleName := range keys {
			if key, exists := a.keys[moduleName]; exists {
				output[moduleName] = key
			}
		}
	}
	return output
}

// baseKeys returns the base keys that are mounted to every version
func (a *App) baseKeys() map[string]*storetypes.KVStoreKey {
	return map[string]*storetypes.KVStoreKey{
		// we need to know the app version to know what stores to mount
		// thus the paramstore must always be a store that is mounted
		paramstypes.StoreKey: a.keys[paramstypes.StoreKey],
	}
}

// migrateModules performs migrations on existing modules that have registered migrations
// between versions and initializes the state of new modules for the specified app version.
func (a App) migrateModules(ctx sdk.Context, fromVersion, toVersion uint64) error {
	return a.manager.RunMigrations(ctx, a.configurator, fromVersion, toVersion)
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (a *App) LegacyAmino() *codec.LegacyAmino {
	return a.legacyAmino
}

// AppCodec returns the app's appCodec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (a *App) AppCodec() codec.Codec {
	return a.appCodec
}

// InterfaceRegistry returns the app's InterfaceRegistry
func (a *App) InterfaceRegistry() codectypes.InterfaceRegistry {
	return a.interfaceRegistry
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (a *App) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return a.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (a *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return a.memKeys[storeKey]
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (a *App) RegisterAPIRoutes(apiSvr *api.Server, _ config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
}

// RegisterTxService implements the Application.RegisterTxService method.
func (a *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(a.BaseApp.GRPCQueryRouter(), clientCtx, a.BaseApp.Simulate, a.interfaceRegistry)
}

func (a *App) InitializeAppVersion() {
	a.SetProtocolVersion(DefaultInitialVersion)
}

func (a App) GetEvmStakingKeeper() *evmstakingkeeper.Keeper {
	return a.Keepers.EvmStakingKeeper
}

func (a App) GetStakingKeeper() *stakingkeeper.Keeper {
	return a.Keepers.StakingKeeper
}

func (a App) GetSlashingKeeper() *slashingkeeper.Keeper {
	return a.Keepers.SlashingKeeper
}

func (a App) GetAccountKeeper() *authkeeper.AccountKeeper {
	return a.Keepers.AccountKeeper
}

func (a App) GetBankKeeper() bankkeeper.Keeper {
	return a.Keepers.BankKeeper
}

func (a App) GetDistrKeeper() *distrkeeper.Keeper {
	return a.Keepers.DistrKeeper
}

func (a App) GetSignalKeeper() *signalkeeper.Keeper {
	return a.Keepers.SignalKeeper
}

func (a App) GetMintKeeper() *mintkeeper.Keeper {
	return a.Keepers.MintKeeper
}
