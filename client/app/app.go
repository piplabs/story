package app

import (
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"

	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/client/comet"
	evmenginekeeper "github.com/piplabs/story/client/x/evmengine/keeper"
	evmstakingkeeper "github.com/piplabs/story/client/x/evmstaking/keeper"
	mintkeeper "github.com/piplabs/story/client/x/mint/keeper"
	"github.com/piplabs/story/lib/buildinfo"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient"

	_ "cosmossdk.io/api/cosmos/tx/config/v1"                // import for side-effects
	_ "cosmossdk.io/x/evidence"                             // import for side-effects
	_ "cosmossdk.io/x/upgrade"                              // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth"                 // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config"       // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/bank"                 // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/consensus"            // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/distribution"         // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/genutil"              // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/gov"                  // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/slashing"             // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/staking"              // import for side-effects
	_ "github.com/piplabs/story/client/x/evmengine/module"  // import for side-effects
	_ "github.com/piplabs/story/client/x/evmstaking/module" // import for side-effects
	_ "github.com/piplabs/story/client/x/mint/module"       // import for side-effects
)

const Name = "story"

var (
	_ runtime.AppI            = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*runtime.App

	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	Keepers keepers.Keepers
}

// newApp returns a reference to an initialized App.
func newApp(
	logger log.Logger,
	db dbm.DB,
	engineCl ethclient.EngineClient,
	baseAppOpts ...func(*baseapp.BaseApp),
) (*App, error) {
	depCfg := depinject.Configs(
		DepConfig(),
		depinject.Supply(
			logger, engineCl,
		),
	)

	var (
		app        = new(App)
		appBuilder = new(runtime.AppBuilder)
	)
	if err := depinject.Inject(depCfg,
		&appBuilder,
		&app.appCodec,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.Keepers.AccountKeeper,
		&app.Keepers.BankKeeper,
		&app.Keepers.StakingKeeper,
		&app.Keepers.SlashingKeeper,
		&app.Keepers.DistrKeeper,
		&app.Keepers.EvidenceKeeper,
		&app.Keepers.ConsensusParamsKeeper,
		&app.Keepers.GovKeeper,
		&app.Keepers.UpgradeKeeper,
		&app.Keepers.EvmStakingKeeper,
		&app.Keepers.EVMEngKeeper,
		&app.Keepers.MintKeeper,
	); err != nil {
		return nil, errors.Wrap(err, "dep inject")
	}

	baseAppOpts = append(baseAppOpts, func(bapp *baseapp.BaseApp) {
		// Use evm engine to create block proposals.
		// Note that we do not check MaxTxBytes since all EngineEVM transaction MUST be included since we cannot
		// postpone them to the next block. Nit: we could drop some vote extensions though...?
		bapp.SetPrepareProposal(app.Keepers.EVMEngKeeper.PrepareProposal)

		// Route proposed messages to keepers for verification and external state updates.
		bapp.SetProcessProposal(makeProcessProposalHandler(makeProcessProposalRouter(app), app.txConfig))
	})

	app.App = appBuilder.Build(db, nil, baseAppOpts...)

	// Override the preblockers with custom PreBlocker function, which handles forks.
	{
		app.ModuleManager.SetOrderPreBlockers(preBlockers...)
		app.SetPreBlocker(app.PreBlocker)
	}

	// Set "OrderEndBlockers" directly instead of using "SetOrderEndBlockers," which will panic since the staking module
	// is missing in the "endBlockers", which is an intended behavior in Story. The panic message is:
	// `panic: all modules must be defined when setting SetOrderEndBlockers, missing: [staking]`
	{
		app.ModuleManager.OrderEndBlockers = endBlockers
		app.SetEndBlocker(app.EndBlocker)
	}

	// Need to manually set the module version map, otherwise dep inject will NOT call `SetModuleVersionMap` for
	// whatever reason that needs to be investigated. Since `SetModuleVersionMap` is not called, `fromVM` will have
	// no entries (i.e. does not know about each module's consensus version) and will try to "add" modules during an
	// upgrade. Specifically, the upgrade module will try to add all modules as new from version 0 to the latest version
	// of each module since `fromVM` is empty on the very first upgrade.
	app.SetInitChainer(func(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
		err := app.Keepers.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap())
		if err != nil {
			return nil, errors.Wrap(err, "set module version map")
		}

		return app.App.InitChainer(ctx, req)
	})

	app.setupUpgradeHandlers()
	app.setupUpgradeStoreLoaders()
	app.SetVersion(buildinfo.Version())

	if err := app.Load(true); err != nil {
		return nil, errors.Wrap(err, "load app")
	}

	return app, nil
}

// PreBlocker application updates every pre block.
func (a *App) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	// All forks should be executed at their planned upgrade heights before any modules.
	a.scheduleForkUpgrade(ctx)

	shouldUpgrade, plan := a.Keepers.EVMEngKeeper.ShouldUpgrade(ctx)
	if shouldUpgrade {
		a.BaseApp.Logger().Info("upgrading app", "upgrade_name", plan.Name, "upgrade_height", plan.Height)

		if err := a.Keepers.UpgradeKeeper.ScheduleUpgrade(ctx, plan); err != nil {
			return nil, errors.Wrap(err, "failed to schedule upgrade")
		}

		if err := a.Keepers.EVMEngKeeper.ResetPendingUpgrade(ctx); err != nil {
			return nil, errors.Wrap(err, "failed to reset pending upgrade")
		}
	}

	res, err := a.ModuleManager.PreBlock(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "module manager preblocker")
	}

	return res, nil
}

func (App) LegacyAmino() *codec.LegacyAmino {
	return nil
}

func (App) ExportAppStateAndValidators(_ bool, _, _ []string) (servertypes.ExportedApp, error) {
	return servertypes.ExportedApp{}, errors.New("not implemented")
}

func (App) SimulationManager() *module.SimulationManager {
	return nil
}

// SetCometAPI sets the comet API client.
func (a App) SetCometAPI(api comet.API) {
	a.Keepers.EVMEngKeeper.SetCometAPI(api)
}

func (a App) GetEVMEngineKeeper() *evmenginekeeper.Keeper {
	return a.Keepers.EVMEngKeeper
}

func (a App) GetEvmStakingKeeper() *evmstakingkeeper.Keeper {
	return a.Keepers.EvmStakingKeeper
}

func (a App) GetStakingKeeper() *stakingkeeper.Keeper {
	return a.Keepers.StakingKeeper
}

func (a App) GetSlashingKeeper() slashingkeeper.Keeper {
	return a.Keepers.SlashingKeeper
}

func (a App) GetAccountKeeper() authkeeper.AccountKeeper {
	return a.Keepers.AccountKeeper
}

func (a App) GetBankKeeper() bankkeeper.Keeper {
	return a.Keepers.BankKeeper
}

func (a App) GetDistrKeeper() distrkeeper.Keeper {
	return a.Keepers.DistrKeeper
}

func (a App) GetUpgradeKeeper() *upgradekeeper.Keeper {
	return a.Keepers.UpgradeKeeper
}

func (a App) GetMintKeeper() mintkeeper.Keeper {
	return a.Keepers.MintKeeper
}
