package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"

	"github.com/piplabs/story/client/x/evmstaking/keeper"
	"github.com/piplabs/story/client/x/evmstaking/module"
	estestutil "github.com/piplabs/story/client/x/evmstaking/testutil"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/ethclient"

	"go.uber.org/mock/gomock"
)

var (
	PKs = simtestutil.CreateTestPubKeys(3)
)

type TestSuite struct {
	suite.Suite

	Ctx sdk.Context

	AccountKeeper    *estestutil.MockAccountKeeper
	BankKeeper       *estestutil.MockBankKeeper
	DistrKeeper      *estestutil.MockDistributionKeeper
	StakingKeeper    *skeeper.Keeper
	EVMStakingKeeper *keeper.Keeper
	msgServer        types.MsgServiceServer

	encCfg moduletestutil.TestEncodingConfig
}

func (s *TestSuite) SetupTest() {
	s.encCfg = moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})
	evmstakingKey := storetypes.NewKVStoreKey(types.StoreKey)
	stakingKey := storetypes.NewKVStoreKey(stypes.StoreKey)
	storeService := runtime.NewKVStoreService(evmstakingKey)
	stakingStoreService := runtime.NewKVStoreService(stakingKey)

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	cms.MountStoreWithDB(evmstakingKey, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(stakingKey, storetypes.StoreTypeIAVL, db)
	err := cms.LoadLatestVersion()
	s.Require().NoError(err)

	s.Ctx = sdk.NewContext(cms, cmtproto.Header{Time: time.Now()}, false, log.NewNopLogger())

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	legacyAmino := codec.NewLegacyAmino()
	stypes.RegisterLegacyAminoCodec(legacyAmino)
	stypes.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("story", "storypub")
	cfg.SetBech32PrefixForValidator("storyvaloper", "storyvaloperpub")
	cfg.SetBech32PrefixForConsensusNode("storyvalcons", "storyvalconspub")

	// gomock initializations
	ctrl := gomock.NewController(s.T())

	// mock keepers
	accountKeeper := estestutil.NewMockAccountKeeper(ctrl)
	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(authtypes.NewModuleAddress(types.ModuleName)).AnyTimes()
	accountKeeper.EXPECT().GetModuleAddress(stypes.ModuleName).Return(authtypes.NewModuleAddress(stypes.ModuleName)).AnyTimes()
	accountKeeper.EXPECT().GetModuleAddress(stypes.BondedPoolName).Return(authtypes.NewModuleAddress(stypes.BondedPoolName)).AnyTimes()
	accountKeeper.EXPECT().GetModuleAddress(stypes.NotBondedPoolName).Return(authtypes.NewModuleAddress(stypes.NotBondedPoolName)).AnyTimes()
	accountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("story")).AnyTimes()
	s.AccountKeeper = accountKeeper
	bankKeeper := estestutil.NewMockBankKeeper(ctrl)
	s.BankKeeper = bankKeeper
	distrKeeper := estestutil.NewMockDistributionKeeper(ctrl)
	s.DistrKeeper = distrKeeper
	slashingKeeper := estestutil.NewMockSlashingKeeper(ctrl)

	// staking keeper
	stakingKeeper := skeeper.NewKeeper(
		marshaler,
		stakingStoreService,
		accountKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(stypes.ModuleName).String(),
		address.NewBech32Codec("storyvaloper"),
		address.NewBech32Codec("storyvalcons"),
	)
	s.StakingKeeper = stakingKeeper
	s.Require().NoError(s.StakingKeeper.SetParams(s.Ctx, stypes.DefaultParams()))

	// emvstaking keeper
	ethCl, err := ethclient.NewEngineMock()
	s.Require().NoError(err)
	evmstakingKeeper := keeper.NewKeeper(
		marshaler,
		storeService,
		accountKeeper,
		bankKeeper,
		slashingKeeper,
		stakingKeeper,
		distrKeeper,
		authtypes.NewModuleAddress(types.ModuleName).String(),
		ethCl,
		address.NewBech32Codec("storyvaloper"),
	)
	s.EVMStakingKeeper = evmstakingKeeper
	s.msgServer = keeper.NewMsgServerImpl(evmstakingKeeper)
}

func TestTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TestSuite))
}

func createCorruptedPubKey(pubKey []byte) []byte {
	corruptedPubKey := append([]byte(nil), pubKey...)
	corruptedPubKey[0] = 0x04
	corruptedPubKey[1] = 0xFF

	return corruptedPubKey
}
