package app

import (
	"context"
	"github.com/piplabs/story/client/x/dkg/keeper"
	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
	"time"

	"cosmossdk.io/store"
	pruningtypes "cosmossdk.io/store/pruning/types"
	"cosmossdk.io/store/snapshots"
	snapshottypes "cosmossdk.io/store/snapshots/types"
	storetypes "cosmossdk.io/store/types"

	abciserver "github.com/cometbft/cometbft/abci/server"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	rpclocal "github.com/cometbft/cometbft/rpc/client/local"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	sdktelemetry "github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/types/mempool"

	"github.com/piplabs/story/client/comet"
	storycfg "github.com/piplabs/story/client/config"
	apisvr "github.com/piplabs/story/client/server"
	"github.com/piplabs/story/lib/buildinfo"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/tracer"
)

// Config wraps the story (app) and comet (client) configurations.
type Config struct {
	storycfg.Config
	Comet cmtcfg.Config
}

// BackendType returns the story config backend type
// or the comet backend type otherwise.
func (c Config) BackendType() dbm.BackendType {
	if c.Config.BackendType == "" {
		return dbm.BackendType(c.Comet.DBBackend)
	}

	return dbm.BackendType(c.Config.BackendType)
}

// Run runs the story client until the context is canceled.
//
//nolint:contextcheck // Explicit new stop context.
func Run(ctx context.Context, cfg Config) error {
	stopFunc, err := Start(ctx, cfg)
	if err != nil {
		return err
	}

	<-ctx.Done()
	log.Info(ctx, "Shutdown detected, stopping...")

	//nolint:mnd // Use a fresh context for stopping (only allow 5 seconds).
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return stopFunc(stopCtx)
}

// Start starts the story client returning a stop function or an error.
//
// Note that the original context used to start the app must be canceled first
// before calling the stop function and a fresh context should be passed into the stop function.
func Start(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	log.Info(ctx, "Starting story consensus client")

	if err := cfg.Verify(); err != nil {
		return nil, errors.Wrap(err, "verify story config")
	}

	buildinfo.Instrument(ctx)

	tracerIDs := tracer.Identifiers{Network: cfg.Network, Service: "story", Instance: cfg.Comet.Moniker}
	stopTracer, err := tracer.Init(ctx, tracerIDs, cfg.Tracer)
	if err != nil {
		return nil, err
	}

	if err := enableSDKTelemetry(); err != nil {
		return nil, errors.Wrap(err, "enable cosmos-sdk telemetry")
	}

	app, privVal, err := CreateApp(ctx, cfg)
	if err != nil {
		return nil, err
	}

	cmtNode, err := newCometNode(ctx, &cfg.Comet, app, privVal, cfg.WithComet, cfg.Address, cfg.Transport)
	if err != nil {
		return nil, errors.Wrap(err, "create comet node")
	}

	var rpcClient *rpclocal.Local
	if cfg.WithComet {
		n, ok := cmtNode.(*node.Node)
		if !ok {
			return nil, errors.Wrap(err, "convert comet node")
		}
		rpcClient = rpclocal.New(n)
		cmtAPI := comet.NewAPI(rpcClient, app.ChainID())

		app.SetCometAPI(cmtAPI)

		log.Info(ctx, "Starting CometBFT", "listeners", n.Listeners())
	}

	if err := cmtNode.Start(); err != nil {
		return nil, errors.Wrap(err, "start comet node")
	}

	var apiSvr *apisvr.Server
	if cfg.API.Enable {
		log.Info(ctx, "Starting API server",
			"api_address", cfg.API.Address,
			"enable_unsafe_cors", cfg.API.EnableUnsafeCORS,
		)
		apiSvr, err = apisvr.NewServer(&cfg.API, app)
		if err != nil {
			return nil, errors.Wrap(err, "create API server")
		}
		if err := apiSvr.Start(); err != nil {
			return nil, errors.Wrap(err, "start API server")
		}
	}

	// Return the stop function.
	// Note that the original context used to start the app must be canceled first.
	// And a fresh context should be passed into the stop function.
	return func(ctx context.Context) error {
		if err := cmtNode.Stop(); err != nil {
			return errors.Wrap(err, "stop comet node")
		}
		<-cmtNode.Quit()

		// Note that cometBFT doesn't shut down cleanly. It leaves a bunch of goroutines running...

		if cfg.API.Enable {
			if err := apiSvr.Stop(ctx); err != nil {
				return errors.Wrap(err, "stop API server")
			}
		}

		if err := stopTracer(ctx); err != nil {
			return errors.Wrap(err, "stop tracer")
		}

		return nil
	}, nil
}

func CreateApp(ctx context.Context, cfg Config) (*App, *privval.FilePV, error) {
	privVal, err := loadPrivVal(cfg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "load validator key")
	}

	db, err := dbm.NewDB("application", cfg.BackendType(), cfg.DataDir())
	if err != nil {
		return nil, nil, errors.Wrap(err, "create db")
	}

	baseAppOpts, err := makeBaseAppOpts(cfg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "make base app opts")
	}

	engineCl, err := newEngineClient(ctx, cfg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "create engine client")
	}

	var (
		dkgTEEClient      dkgtypes.TEEClient
		dkgContractClient *keeper.ContractClient
	)

	if cfg.DKG.Enable {
		dkgTEEClient, err = keeper.CreateTEEClient(cfg.DKG)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create tee client for DKG")
		}

		dkgContractClient, err = keeper.NewContractClient(ctx, cfg.DKG.EngineRPCEndpoint, cfg.EngineChainID, privVal.Key.PrivKey.Bytes())
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create contract client for DKG")
		}
	}

	//nolint:contextcheck // False positive
	app, err := newApp(
		newSDKLogger(ctx),
		db,
		engineCl,
		dkgTEEClient,
		dkgContractClient,
		baseAppOpts...,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "create app")
	}
	app.Keepers.EVMEngKeeper.SetBuildDelay(cfg.EVMBuildDelay)
	app.Keepers.EVMEngKeeper.SetBuildOptimistic(cfg.EVMBuildOptimistic)

	addr, err := k1util.PubKeyToAddress(privVal.Key.PrivKey.PubKey())
	if err != nil {
		return nil, nil, errors.Wrap(err, "convert validator pubkey to address")
	}
	app.Keepers.EVMEngKeeper.SetValidatorAddress(addr)

	if cfg.DKG.Enable {
		if err := app.Keepers.DKGKeeper.InitDKGService(&cfg.Config, addr); err != nil {
			return nil, nil, errors.Wrap(err, "dkg service is enabled, but failed to init dkg service")
		}
	}

	return app, privVal, nil
}

func newCometNode(ctx context.Context, cfg *cmtcfg.Config, app *App, privVal cmttypes.PrivValidator, withComet bool, addr, transport string,
) (service.Service, error) {
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, errors.Wrap(err, "load or gen node key", "key_file", cfg.NodeKeyFile())
	}

	cmtLog, err := newCmtLogger(ctx, cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	wrapper := newABCIWrapper(
		server.NewCometABCIWrapper(app),
		app.Keepers.EVMEngKeeper.PostFinalize,
		func() storetypes.CacheMultiStore {
			return app.CommitMultiStore().CacheMultiStore()
		},
	)

	var cmtNode service.Service
	if withComet {
		cmtNode, err = node.NewNode(cfg,
			privVal,
			nodeKey,
			proxy.NewLocalClientCreator(wrapper),
			node.DefaultGenesisDocProviderFunc(cfg),
			cmtcfg.DefaultDBProvider,
			node.DefaultMetricsProvider(cfg.Instrumentation),
			cmtLog,
		)
		if err != nil {
			return nil, errors.Wrap(err, "create node")
		}
	} else {
		cmtNode, err = abciserver.NewServer(addr, transport, wrapper)
		if err != nil {
			return nil, errors.Wrap(err, "create abci server")
		}
		cmtNode.SetLogger(cmtLog.With("module", "abci-server"))
	}

	return cmtNode, nil
}

func makeBaseAppOpts(cfg Config) ([]func(*baseapp.BaseApp), error) {
	chainID, err := chainIDFromGenesis(cfg)
	if err != nil {
		return nil, err
	}

	snapshotStore, err := newSnapshotStore(cfg)
	if err != nil {
		return nil, err
	}

	snapshotOptions := snapshottypes.NewSnapshotOptions(cfg.SnapshotInterval, uint32(cfg.SnapshotKeepRecent))

	pruneOpts := pruningtypes.NewPruningOptionsFromString(cfg.PruningOption)
	if cfg.PruningOption == pruningtypes.PruningOptionDefault {
		// Override the default cosmosSDK pruning values with much more aggressive defaults
		// since historical state isn't very important for most use-cases.
		pruneOpts = pruningtypes.NewCustomPruningOptions(defaultPruningKeep, defaultPruningInterval)
	} else if cfg.PruningOption == pruningtypes.PruningOptionCustom {
		pruneOpts = pruningtypes.NewCustomPruningOptions(cfg.PruningKeepRecent, cfg.PruningInterval)
	}

	if err := pruneOpts.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid custom pruning option")
	}

	return []func(*baseapp.BaseApp){
		// baseapp.SetOptimisticExecution(), // TODO: research this feature
		baseapp.SetChainID(chainID),
		baseapp.SetMinRetainBlocks(cfg.MinRetainBlocks),
		baseapp.SetPruning(pruneOpts),
		baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager()),
		baseapp.SetSnapshot(snapshotStore, snapshotOptions),
		baseapp.SetMempool(mempool.NoOpMempool{}),
	}, nil
}

func newSnapshotStore(cfg Config) (*snapshots.Store, error) {
	db, err := dbm.NewDB("metadata", cfg.BackendType(), cfg.SnapshotDir())
	if err != nil {
		return nil, errors.Wrap(err, "create snapshot db")
	}

	ss, err := snapshots.NewStore(db, cfg.SnapshotDir())
	if err != nil {
		return nil, errors.Wrap(err, "create snapshot store")
	}

	return ss, nil
}

func chainIDFromGenesis(cfg Config) (string, error) {
	genDoc, err := node.DefaultGenesisDocProviderFunc(&cfg.Comet)()
	if err != nil {
		return "", errors.Wrap(err, "load genesis doc")
	}

	return genDoc.ChainID, nil
}

// newEngineClient returns a new engine API client.
func newEngineClient(ctx context.Context, cfg Config) (ethclient.EngineClient, error) {
	jwtBytes, err := ethclient.LoadJWTHexFile(cfg.EngineJWTFile)
	if err != nil {
		return nil, errors.Wrap(err, "load engine JWT file")
	}

	engineCl, err := ethclient.NewAuthClient(ctx, cfg.EngineEndpoint, jwtBytes)
	if err != nil {
		return nil, errors.Wrap(err, "create engine client")
	}

	return engineCl, nil
}

// enableSDKTelemetry enables prometheus based cosmos-sdk telemetry.
func enableSDKTelemetry() error {
	const farFuture = time.Hour * 24 * 365 * 10 // 10 years ~= infinity.

	_, err := sdktelemetry.New(sdktelemetry.Config{
		ServiceName:             "cosmos",
		Enabled:                 true,
		PrometheusRetentionTime: int64(farFuture.Seconds()), // Prometheus metrics never expire once created in-app.
	})
	if err != nil {
		return errors.Wrap(err, "enable cosmos-sdk telemetry")
	}

	return nil
}
