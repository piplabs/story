package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/mint/keeper"
	mintmodule "github.com/piplabs/story/client/x/mint/module"
	minttestutil "github.com/piplabs/story/client/x/mint/testutil"
	"github.com/piplabs/story/client/x/mint/types"

	"go.uber.org/mock/gomock"
)

func createKeeper(t *testing.T) (sdk.Context, *minttestutil.MockAccountKeeper, *minttestutil.MockBankKeeper, *minttestutil.MockStakingKeeper, *keeper.Keeper) {
	t.Helper()
	encCfg := moduletestutil.MakeTestEncodingConfig(mintmodule.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)

	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

	ctrl := gomock.NewController(t)
	accountKeeper := minttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := minttestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := minttestutil.NewMockStakingKeeper(ctrl)

	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.AccAddress{})

	mk := keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		stakingKeeper,
		accountKeeper,
		bankKeeper,
		authtypes.FeeCollectorName,
	)

	// set default params
	require.NoError(t, mk.Params.Set(testCtx.Ctx, types.DefaultParams()))

	return testCtx.Ctx, accountKeeper, bankKeeper, stakingKeeper, &mk
}

func TestNewMintKeeperPanic(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(mintmodule.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)

	ctrl := gomock.NewController(t)
	accountKeeper := minttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := minttestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := minttestutil.NewMockStakingKeeper(ctrl)

	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(nil)

	require.PanicsWithValue(t, "the x/mint module account has not been set", func() {
		_ = keeper.NewKeeper(
			encCfg.Codec,
			storeService,
			stakingKeeper,
			accountKeeper,
			bankKeeper,
			authtypes.FeeCollectorName,
		)
	})
}

func TestAliasFunctions(t *testing.T) {
	ctx, _, bk, sk, mk := createKeeper(t)

	require.Equal(t, ctx.Logger().With("module", "x/"+types.ModuleName),
		mk.Logger(ctx))

	stakingTokenSupply := math.NewIntFromUint64(100000000000)
	sk.EXPECT().StakingTokenSupply(ctx).Return(stakingTokenSupply, nil)
	tokenSupply, err := mk.StakingTokenSupply(ctx)
	require.NoError(t, err)
	require.Equal(t, tokenSupply, stakingTokenSupply)

	bondedRatio := math.LegacyNewDecWithPrec(15, 2)
	sk.EXPECT().BondedRatio(ctx).Return(bondedRatio, nil)
	ratio, err := mk.BondedRatio(ctx)
	require.NoError(t, err)
	require.Equal(t, ratio, bondedRatio)

	coins := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1000000)))
	bk.EXPECT().MintCoins(ctx, types.ModuleName, coins).Return(nil)
	require.NoError(t, mk.MintCoins(ctx, sdk.NewCoins()))
	require.NoError(t, mk.MintCoins(ctx, coins))

	fees := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1000)))
	bk.EXPECT().SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, fees).Return(nil)
	require.NoError(t, mk.AddCollectedFees(ctx, fees))
}
