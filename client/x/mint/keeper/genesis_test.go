//nolint:paralleltest // just for testing
package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/suite"

	"github.com/piplabs/story/client/x/mint/keeper"
	mintmodule "github.com/piplabs/story/client/x/mint/module"
	minttestutil "github.com/piplabs/story/client/x/mint/testutil"
	"github.com/piplabs/story/client/x/mint/types"

	"go.uber.org/mock/gomock"
)

var minterAcc = authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter)

type GenesisTestSuite struct {
	suite.Suite

	sdkCtx        sdk.Context
	keeper        keeper.Keeper
	cdc           codec.BinaryCodec
	accountKeeper types.AccountKeeper
	key           *storetypes.KVStoreKey
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (s *GenesisTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(mintmodule.AppModuleBasic{})

	// gomock initializations
	ctrl := gomock.NewController(s.T())
	s.cdc = codec.NewProtoCodec(encCfg.InterfaceRegistry)
	s.sdkCtx = testCtx.Ctx
	s.key = key

	stakingKeeper := minttestutil.NewMockStakingKeeper(ctrl)
	accountKeeper := minttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := minttestutil.NewMockBankKeeper(ctrl)
	s.accountKeeper = accountKeeper
	accountKeeper.EXPECT().GetModuleAddress(minterAcc.Name).Return(minterAcc.GetAddress())
	accountKeeper.EXPECT().GetModuleAccount(s.sdkCtx, minterAcc.Name).Return(minterAcc)

	s.keeper = keeper.NewKeeper(s.cdc, runtime.NewKVStoreService(key), stakingKeeper, accountKeeper, bankKeeper, "")
}

func (s *GenesisTestSuite) TestImportExportGenesis() {
	genesisState := types.DefaultGenesisState()
	genesisState.Params = types.NewParams(
		"testDenom",
		math.LegacyNewDec(24625000000000000.000000000000000000),
		uint64(60*60*8766/5),
	)

	s.keeper.InitGenesis(s.sdkCtx, s.accountKeeper, genesisState)

	params, err := s.keeper.Params.Get(s.sdkCtx)
	s.Require().Equal(genesisState.Params, params)
	s.Require().NoError(err)

	genesisState2 := s.keeper.ExportGenesis(s.sdkCtx)
	s.Require().Equal(genesisState, genesisState2)
}
