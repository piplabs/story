//nolint:paralleltest // just for testing
package keeper_test

import (
	gocontext "context"
	"testing"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/piplabs/story/client/x/mint/keeper"
	mintmodule "github.com/piplabs/story/client/x/mint/module"
	minttestutil "github.com/piplabs/story/client/x/mint/testutil"
	"github.com/piplabs/story/client/x/mint/types"
)

type MintTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	queryClient types.QueryClient
	mintKeeper  keeper.Keeper
}

func (suite *MintTestSuite) SetupTest() {
	encCfg := moduletestutil.MakeTestEncodingConfig(mintmodule.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(suite.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	suite.ctx = testCtx.Ctx

	// gomock initializations
	ctrl := gomock.NewController(suite.T())
	accountKeeper := minttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := minttestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := minttestutil.NewMockStakingKeeper(ctrl)

	accountKeeper.EXPECT().GetModuleAddress("mint").Return(sdk.AccAddress{})

	suite.mintKeeper = keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		stakingKeeper,
		accountKeeper,
		bankKeeper,
		authtypes.FeeCollectorName,
	)

	err := suite.mintKeeper.Params.Set(suite.ctx, types.DefaultParams())
	suite.Require().NoError(err)
	suite.Require().NoError(suite.mintKeeper.Minter.Set(suite.ctx, types.DefaultInitialMinter()))

	queryHelper := baseapp.NewQueryServerTestHelper(testCtx.Ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServerImpl(suite.mintKeeper))

	suite.queryClient = types.NewQueryClient(queryHelper)
}

func (suite *MintTestSuite) TestGRPCParams() {
	params, err := suite.queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	kparams, err := suite.mintKeeper.Params.Get(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(params.Params, kparams)
}

func TestMintTestSuite(t *testing.T) {
	suite.Run(t, new(MintTestSuite))
}
