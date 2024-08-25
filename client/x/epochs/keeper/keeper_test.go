package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/core/header"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosmodule "github.com/cosmos/cosmos-sdk/types/module"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	epochskeeper "github.com/piplabs/story/client/x/epochs/keeper"
	"github.com/piplabs/story/client/x/epochs/module"
	"github.com/piplabs/story/client/x/epochs/types"
)

type KeeperTestSuite struct {
	suite.Suite
	Ctx          sdk.Context
	EpochsKeeper *epochskeeper.Keeper
	queryClient  types.QueryClient
}

func (s *KeeperTestSuite) SetupTest() {
	s.Ctx, s.EpochsKeeper = Setup(s.T())

	queryRouter := baseapp.NewGRPCQueryRouter()
	cfg := cosmosmodule.NewConfigurator(nil, nil, queryRouter)
	types.RegisterQueryServer(cfg.QueryServer(), epochskeeper.NewQuerier(*s.EpochsKeeper))
	grpcQueryService := &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: queryRouter,
		Ctx:             s.Ctx,
	}
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})
	grpcQueryService.SetInterfaceRegistry(encCfg.InterfaceRegistry)
	s.queryClient = types.NewQueryClient(grpcQueryService)
}

func Setup(t *testing.T) (sdk.Context, *epochskeeper.Keeper) {
	t.Helper()

	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	eventService := runtime.EventService{}
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithHeaderInfo(header.Info{Time: time.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})

	epochsKeeper := epochskeeper.NewKeeper(
		storeService,
		eventService,
		encCfg.Codec,
	)
	epochsKeeper = epochsKeeper.SetHooks(types.NewMultiEpochHooks())
	ctx.WithHeaderInfo(header.Info{Height: 1, Time: time.Now().UTC(), ChainID: "epochs"})
	err := epochsKeeper.InitGenesis(ctx, *types.DefaultGenesis())
	require.NoError(t, err)
	SetEpochStartTime(ctx, *epochsKeeper)

	return ctx, epochsKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(KeeperTestSuite))
}

func SetEpochStartTime(ctx sdk.Context, epochsKeeper epochskeeper.Keeper) {
	epochs, err := epochsKeeper.AllEpochInfos(ctx)
	if err != nil {
		panic(err)
	}
	for _, epoch := range epochs {
		epoch.StartTime = ctx.BlockTime()
		err := epochsKeeper.EpochInfo.Remove(ctx, epoch.Identifier)
		if err != nil {
			panic(err)
		}
		err = epochsKeeper.AddEpochInfo(ctx, epoch)
		if err != nil {
			panic(err)
		}
	}
}
