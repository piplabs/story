package keeper

// Internal test file (package keeper, not keeper_test) to access private
// methods directly. Tests the deferred MaxValidators reduction flow:
//   applyDeferredMaxValidatorsChange → staking ApplyAndReturnValidatorSetUpdates
//
// Uses the real staking keeper to verify that ABCI validator updates are
// produced correctly.

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"

	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	estestutil "github.com/piplabs/story/client/x/evmstaking/testutil"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/ethclient"
	"github.com/piplabs/story/lib/netconf"
)

var initBech32Once sync.Once

func initBech32Prefixes() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("story", "storypub")
	cfg.SetBech32PrefixForValidator("storyvaloper", "storyvaloperpub")
	cfg.SetBech32PrefixForConsensusNode("storyvalcons", "storyvalconspub")
}

func TestDeferredMaxValidatorsChange_EndToEnd(t *testing.T) {
	const numValidators = 25

	initBech32Once.Do(initBech32Prefixes)

	key := storetypes.NewKVStoreKey("test")
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()}).WithChainID(netconf.TestChainID)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	stypes.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	ctrl := gomock.NewController(t)
	ak := estestutil.NewMockAccountKeeper(ctrl)
	bk := estestutil.NewMockBankKeeper(ctrl)
	dk := estestutil.NewMockDistributionKeeper(ctrl)
	slk := estestutil.NewMockSlashingKeeper(ctrl)
	dkgk := estestutil.NewMockDKGKeeper(ctrl)

	ak.EXPECT().AddressCodec().Return(authcodec.NewBech32Codec("story")).AnyTimes()
	ak.EXPECT().GetModuleAddress(gomock.Any()).Return(authtypes.NewModuleAddress(types.ModuleName)).AnyTimes()
	bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	sk := skeeper.NewKeeper(cdc, storeService, ak, bk,
		authtypes.NewModuleAddress(stypes.ModuleName).String(),
		address.NewBech32Codec("storyvaloper"),
		address.NewBech32Codec("storyvalcons"))

	params := stypes.NewParams(
		7*24*time.Hour, // unbonding time — long enough
		uint32(numValidators),
		stypes.DefaultMaxEntries,
		stypes.DefaultHistoricalEntries,
		sdk.DefaultBondDenom,
		stypes.DefaultMinCommissionRate,
		math.NewInt(2),
		stypes.DefaultFlexiblePeriodType,
		stypes.DefaultPeriods,
		stypes.DefaultLockedTokenType,
		stypes.DefaultTokenTypes,
		0, // singularity disabled
	)
	require.NoError(t, sk.SetParams(ctx, params))

	ethCl, err := ethclient.NewEngineMock(key)
	require.NoError(t, err)

	esk := NewKeeper(cdc, storeService, ak, bk, slk, sk, dk, dkgk,
		authtypes.NewModuleAddress(types.ModuleName).String(),
		ethCl,
		address.NewBech32Codec("storyvaloper"))
	require.NoError(t, esk.SetParams(ctx, types.DefaultParams()))

	// Create bonded validators
	for i := 0; i < numValidators; i++ {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		valAddr := sdk.ValAddress(pubKey.Address())
		val, err := stypes.NewValidator(valAddr.String(), pubKey, stypes.Description{Moniker: "test"}, 0)
		require.NoError(t, err)
		val.Commission = stypes.NewCommission(
			math.LegacyNewDecWithPrec(1000, 4),
			math.LegacyNewDecWithPrec(5000, 4),
			math.LegacyNewDecWithPrec(1000, 4),
		)
		valTokens := sk.TokensFromConsensusPower(ctx, 10)
		validator, _, _ := val.AddTokensFromDel(valTokens, math.LegacyNewDecFromInt(valTokens).Quo(math.LegacyNewDec(2)))
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		_ = skeeper.TestingUpdateValidator(sk, sdkCtx, validator, true)
		require.NoError(t, sk.SetValidatorByConsAddr(ctx, validator))
	}

	// Verify pre-state
	bonded, err := sk.GetBondedValidatorsByPower(ctx)
	require.NoError(t, err)
	require.Len(t, bonded, numValidators, "pre: %d validators bonded", numValidators)

	// Step 1: applyDeferredMaxValidatorsChange at Seneca height
	ctx = sdk.UnwrapSDKContext(ctx).WithBlockHeight(300) // TestChainID Seneca = 300
	err = esk.applyDeferredMaxValidatorsChange(ctx)
	require.NoError(t, err)

	p, err := sk.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, uint32(21), p.MaxValidators, "MaxValidators should be 21")

	// Step 2: staking ApplyAndReturnValidatorSetUpdates (same as EndBlocker would call)
	updates, err := sk.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)

	expectedRemovals := numValidators - 21
	require.Len(t, updates, expectedRemovals,
		"Apply should return %d removal updates", expectedRemovals)

	for i, u := range updates {
		require.Equal(t, int64(0), u.Power, "update[%d] should be zero-power", i)
	}

	// Step 3: bonded set is now 21
	bonded, err = sk.GetBondedValidatorsByPower(ctx)
	require.NoError(t, err)
	require.Len(t, bonded, 21, "post: exactly 21 bonded")
}

func TestDeferredMaxValidatorsChange_NotUpgradeHeight(t *testing.T) {
	key := storetypes.NewKVStoreKey("test")
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()}).WithChainID(netconf.TestChainID)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	stypes.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	ctrl := gomock.NewController(t)
	ak := estestutil.NewMockAccountKeeper(ctrl)
	bk := estestutil.NewMockBankKeeper(ctrl)

	ak.EXPECT().AddressCodec().Return(authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())).AnyTimes()
	ak.EXPECT().GetModuleAddress(gomock.Any()).Return(authtypes.NewModuleAddress(types.ModuleName)).AnyTimes()

	sk := skeeper.NewKeeper(cdc, storeService, ak, bk,
		authtypes.NewModuleAddress(stypes.ModuleName).String(),
		address.NewBech32Codec("storyvaloper"),
		address.NewBech32Codec("storyvalcons"))
	require.NoError(t, sk.SetParams(ctx, stypes.DefaultParams()))

	esk := &Keeper{stakingKeeper: sk}

	origParams, _ := sk.GetParams(ctx)
	origMax := origParams.MaxValidators

	// Height 100 — not Seneca (300)
	ctx = sdk.UnwrapSDKContext(ctx).WithBlockHeight(100)
	err := esk.applyDeferredMaxValidatorsChange(ctx)
	require.NoError(t, err)

	p, _ := sk.GetParams(ctx)
	require.Equal(t, origMax, p.MaxValidators, "should not change at non-upgrade height")
}

func TestDeferredMaxValidatorsChange_AlreadyBelow(t *testing.T) {
	key := storetypes.NewKVStoreKey("test")
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()}).WithChainID(netconf.TestChainID)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	stypes.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	ctrl := gomock.NewController(t)
	ak := estestutil.NewMockAccountKeeper(ctrl)
	bk := estestutil.NewMockBankKeeper(ctrl)

	ak.EXPECT().AddressCodec().Return(authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())).AnyTimes()
	ak.EXPECT().GetModuleAddress(gomock.Any()).Return(authtypes.NewModuleAddress(types.ModuleName)).AnyTimes()

	sk := skeeper.NewKeeper(cdc, storeService, ak, bk,
		authtypes.NewModuleAddress(stypes.ModuleName).String(),
		address.NewBech32Codec("storyvaloper"),
		address.NewBech32Codec("storyvalcons"))

	params := stypes.DefaultParams()
	params.MaxValidators = 10
	require.NoError(t, sk.SetParams(ctx, params))

	esk := &Keeper{stakingKeeper: sk}

	ctx = sdk.UnwrapSDKContext(ctx).WithBlockHeight(300)
	err := esk.applyDeferredMaxValidatorsChange(ctx)
	require.NoError(t, err)

	p, _ := sk.GetParams(ctx)
	require.Equal(t, uint32(10), p.MaxValidators, "should not increase")
}
