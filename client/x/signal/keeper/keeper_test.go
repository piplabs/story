package keeper_test

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmversion "github.com/cometbft/cometbft/proto/tendermint/version"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	v1 "github.com/docker/docker/image/spec/specs-go"
	"github.com/stretchr/testify/suite"

	"github.com/piplabs/story/client/x/signal/keeper"
	"github.com/piplabs/story/client/x/signal/module"
	"github.com/piplabs/story/client/x/signal/testutil"
	"github.com/piplabs/story/client/x/signal/types"
	"github.com/piplabs/story/lib/errors"

	"go.uber.org/mock/gomock"
)

type TestSuite struct {
	suite.Suite

	Ctx sdk.Context

	AccountKeeper *testutil.MockAccountKeeper
	StakingKeeper types.StakingKeeper
	UpgradeKeeper *keeper.Keeper
	msgServer     types.MsgServiceServer

	encCfg moduletestutil.TestEncodingConfig
}

func (s *TestSuite) SetupTest(t *testing.T) {
	s.encCfg = moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})
	signalStore := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(signalStore)

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	cms.MountStoreWithDB(signalStore, storetypes.StoreTypeIAVL, db)
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
	accountKeeper := testutil.NewMockAccountKeeper(ctrl)
	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(authtypes.NewModuleAddress(types.ModuleName)).AnyTimes()
	accountKeeper.EXPECT().GetModuleAddress(stypes.ModuleName).Return(authtypes.NewModuleAddress(stypes.ModuleName)).AnyTimes()
	accountKeeper.EXPECT().GetModuleAddress(stypes.BondedPoolName).Return(authtypes.NewModuleAddress(stypes.BondedPoolName)).AnyTimes()
	accountKeeper.EXPECT().GetModuleAddress(stypes.NotBondedPoolName).Return(authtypes.NewModuleAddress(stypes.NotBondedPoolName)).AnyTimes()
	accountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("story")).AnyTimes()

	mockStakingKeeper := newMockStakingKeeper(
		map[string]int64{
			testutil.ValAddrs[0].String(): 40,
			testutil.ValAddrs[1].String(): 1,
			testutil.ValAddrs[2].String(): 59,
			testutil.ValAddrs[3].String(): 20,
		},
	)
	s.StakingKeeper = mockStakingKeeper

	upgradeKeeper := keeper.NewKeeper(
		marshaler,
		storeService,
		accountKeeper,
		mockStakingKeeper,
		authtypes.NewModuleAddress(stypes.ModuleName).String(),
	)
	s.UpgradeKeeper = upgradeKeeper
}

func (s *TestSuite) TestGetVotingPowerThreshold() {
	require := s.Require()

	bigInt := big.NewInt(0)
	bigInt.SetString("23058430092136939509", 10)

	type testCase struct {
		name       string
		validators map[string]int64
		want       sdkmath.Int
	}
	tcs := []testCase{
		{
			name:       "empty validators",
			validators: map[string]int64{},
			want:       sdkmath.NewInt(0),
		},
		{
			name:       "one validator with 6 power returns 5 because the defaultSignalThreshold is 5/6",
			validators: map[string]int64{"a": 6},
			want:       sdkmath.NewInt(5),
		},
		{
			name:       "one validator with 11 power (11 * 5/6 = 9.16666667) so should round up to 10",
			validators: map[string]int64{"a": 11},
			want:       sdkmath.NewInt(10),
		},
		{
			name:       "one validator with voting power of math.MaxInt64",
			validators: map[string]int64{"a": math.MaxInt64},
			want:       sdkmath.NewInt(7686143364045646503),
		},
		{
			name:       "multiple validators with voting power of math.MaxInt64",
			validators: map[string]int64{"a": math.MaxInt64, "b": math.MaxInt64, "c": math.MaxInt64},
			want:       sdkmath.NewIntFromBigInt(bigInt),
		},
	}
	for _, tc := range tcs {
		s.Run(tc.name, func() {
			got := s.UpgradeKeeper.GetVotingPowerThreshold(sdk.Context{})
			require.Equal(tc.want, got, fmt.Sprintf("want %v, got %v", tc.want.String(), got.String()))
		})
	}
}

func (s *TestSuite) TestSignalVersion() {
	require := s.Require()

	ctx, upgradeKeeper := s.Ctx, s.UpgradeKeeper
	goCtx := sdk.WrapSDKContext(ctx)

	s.Run("should return an error if the signal version is less than the current version", func() {
		_, err := upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
			ValidatorAddress: testutil.ValAddrs[0].String(),
			Version:          0,
		})
		require.Error(err)
		require.ErrorIs(err, types.ErrInvalidSignalVersion)
	})

	s.Run("should not return an error if the signal version is greater than the next version", func() {
		_, err := upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
			ValidatorAddress: testutil.ValAddrs[0].String(),
			Version:          3,
		})
		require.NoError(err)
	})

	s.Run("should return an error if the validator was not found", func() {
		_, err := upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
			ValidatorAddress: testutil.ValAddrs[4].String(),
			Version:          2,
		})
		require.Error(err)
		require.ErrorIs(err, stypes.ErrNoValidatorFound)
	})

	s.Run("should not return an error if the signal version and validator are valid", func() {
		_, err := upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
			ValidatorAddress: testutil.ValAddrs[0].String(),
			Version:          2,
		})
		require.NoError(err)

		res, err := upgradeKeeper.VersionTally(goCtx, &types.QueryVersionTallyRequest{
			Version: 2,
		})
		require.NoError(err)
		require.EqualValues(40, res.VotingPower)
		require.EqualValues(100, res.ThresholdPower)
		require.EqualValues(120, res.TotalVotingPower)
	})
}

func (s *TestSuite) TestTallyingLogic(t *testing.T) {
	require := s.Require()

	ctx, upgradeKeeper := s.Ctx, s.UpgradeKeeper
	stakingKeeper := s.StakingKeeper.(*mockStakingKeeper)
	goCtx := sdk.WrapSDKContext(ctx)

	_, err := upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
		ValidatorAddress: testutil.ValAddrs[0].String(),
		Version:          0,
	})
	require.Error(err)
	require.ErrorIs(err, types.ErrInvalidSignalVersion)

	_, err = upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
		ValidatorAddress: testutil.ValAddrs[0].String(),
		Version:          3,
	})
	require.NoError(err) // version 3 is valid because it is greater than the current version

	_, err = upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
		ValidatorAddress: testutil.ValAddrs[0].String(),
		Version:          2,
	})
	require.NoError(err)

	res, err := upgradeKeeper.VersionTally(goCtx, &types.QueryVersionTallyRequest{
		Version: 2,
	})
	require.NoError(err)
	require.EqualValues(40, res.VotingPower)
	require.EqualValues(100, res.ThresholdPower)
	require.EqualValues(120, res.TotalVotingPower)

	_, err = upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
		ValidatorAddress: testutil.ValAddrs[2].String(),
		Version:          2,
	})
	require.NoError(err)

	res, err = upgradeKeeper.VersionTally(goCtx, &types.QueryVersionTallyRequest{
		Version: 2,
	})
	require.NoError(err)
	require.EqualValues(99, res.VotingPower)
	require.EqualValues(100, res.ThresholdPower)
	require.EqualValues(120, res.TotalVotingPower)

	_, err = upgradeKeeper.TryUpgrade(goCtx, &types.MsgTryUpgrade{})
	require.NoError(err)
	shouldUpgrade, version := upgradeKeeper.ShouldUpgrade(ctx)
	require.False(shouldUpgrade)
	require.Equal(uint64(0), version)

	// we now have 101/120
	_, err = upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
		ValidatorAddress: testutil.ValAddrs[1].String(),
		Version:          2,
	})
	require.NoError(err)

	_, err = upgradeKeeper.TryUpgrade(goCtx, &types.MsgTryUpgrade{})
	require.NoError(err)

	shouldUpgrade, version = upgradeKeeper.ShouldUpgrade(ctx)
	require.False(shouldUpgrade) // should be false because upgrade height hasn't been reached.
	require.Equal(uint64(0), version)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + keeper.DefaultUpgradeHeightDelay)

	shouldUpgrade, version = upgradeKeeper.ShouldUpgrade(ctx)
	require.True(shouldUpgrade) // should be true because upgrade height has been reached.
	require.Equal(2, version)

	upgradeKeeper.ResetTally(ctx)

	// update the version to 2
	ctx = ctx.WithBlockHeader(cmtproto.Header{
		Version: tmversion.Consensus{
			Block: 1,
			App:   2,
		},
	})
	goCtx = sdk.WrapSDKContext(ctx)

	_, err = upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
		ValidatorAddress: testutil.ValAddrs[0].String(),
		Version:          3,
	})
	require.NoError(err)
	_, err = upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
		ValidatorAddress: testutil.ValAddrs[1].String(),
		Version:          2,
	})
	require.NoError(err)
	_, err = upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
		ValidatorAddress: testutil.ValAddrs[2].String(),
		Version:          2,
	})
	require.NoError(err)

	res, err = upgradeKeeper.VersionTally(goCtx, &types.QueryVersionTallyRequest{
		Version: 2,
	})
	require.NoError(err)
	require.EqualValues(60, res.VotingPower)
	require.EqualValues(100, res.ThresholdPower)
	require.EqualValues(120, res.TotalVotingPower)

	// remove one of the validators from the set
	delete(stakingKeeper.validators, testutil.ValAddrs[1].String())
	// the validator had 1 voting power, so we deduct it from the total
	stakingKeeper.totalVotingPower = stakingKeeper.totalVotingPower.SubRaw(1)

	res, err = upgradeKeeper.VersionTally(goCtx, &types.QueryVersionTallyRequest{
		Version: 2,
	})
	require.NoError(err)
	require.EqualValues(59, res.VotingPower)
	require.EqualValues(100, res.ThresholdPower)
	require.EqualValues(119, res.TotalVotingPower)

	// That validator should not be able to signal a version
	_, err = upgradeKeeper.SignalVersion(goCtx, &types.MsgSignalVersion{
		ValidatorAddress: testutil.ValAddrs[1].String(),
		Version:          2,
	})
	require.Error(err)

	// resetting the tally should clear other votes
	upgradeKeeper.ResetTally(ctx)
	res, err = upgradeKeeper.VersionTally(goCtx, &types.QueryVersionTallyRequest{
		Version: 2,
	})
	require.NoError(err)
	require.EqualValues(0, res.VotingPower)
}

// TestCanSkipVersion verifies that the signal keeper can upgrade to an app
// version greater than the next app version. Example: if the current version is
// 1, the next version is 2, but the chain can upgrade directly from 1 to 3.
func (s *TestSuite) TestCanSkipVersion(t *testing.T) {
	require := s.Require()

	ctx, upgradeKeeper := s.Ctx, s.UpgradeKeeper
	goCtx := sdk.WrapSDKContext(ctx)

	require.Equal(v1.Version, ctx.BlockHeader().Version.App)

	validators := []sdk.ValAddress{
		testutil.ValAddrs[0],
		testutil.ValAddrs[1],
		testutil.ValAddrs[2],
		testutil.ValAddrs[3],
	}
	// signal version 3 for all validators
	for _, validator := range validators {
		_, err := upgradeKeeper.SignalVersion(sdk.WrapSDKContext(ctx), &types.MsgSignalVersion{
			ValidatorAddress: validator.String(),
			Version:          3,
		})
		require.NoError(err)
	}

	_, err := upgradeKeeper.TryUpgrade(goCtx, &types.MsgTryUpgrade{})
	require.NoError(err)

	isUpgradePending := upgradeKeeper.IsUpgradePending(ctx)
	require.True(isUpgradePending)
}

func (s *TestSuite) TestEmptyStore(t *testing.T) {
	require := s.Require()

	ctx, upgradeKeeper := s.Ctx, s.UpgradeKeeper
	goCtx := sdk.WrapSDKContext(ctx)

	res, err := upgradeKeeper.VersionTally(goCtx, &types.QueryVersionTallyRequest{
		Version: 2,
	})
	require.NoError(err)
	require.EqualValues(0, res.VotingPower)
	// 120 is the summation in voting power of the four validators
	require.EqualValues(120, res.TotalVotingPower)
}

func (s *TestSuite) TestThresholdVotingPower(t *testing.T) {
	require := s.Require()

	ctx, upgradeKeeper := s.Ctx, s.UpgradeKeeper
	stakingKeeper := s.StakingKeeper.(*mockStakingKeeper)

	for _, tc := range []struct {
		total     int64
		threshold int64
	}{
		{total: 1, threshold: 1},
		{total: 2, threshold: 2},
		{total: 3, threshold: 3},
		{total: 6, threshold: 5},
		{total: 59, threshold: 50},
	} {
		stakingKeeper.totalVotingPower = sdkmath.NewInt(tc.total)
		threshold := upgradeKeeper.GetVotingPowerThreshold(ctx)
		require.EqualValues(tc.threshold, threshold.Int64())
	}
}

// TestResetTally verifies that ResetTally resets the VotingPower for all
// versions to 0 and any pending upgrade is cleared.
func (s *TestSuite) TestResetTally(t *testing.T) {
	require := s.Require()

	ctx, upgradeKeeper := s.Ctx, s.UpgradeKeeper

	_, err := upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[0].String(), Version: 2})
	require.NoError(err)
	resp, err := upgradeKeeper.VersionTally(ctx, &types.QueryVersionTallyRequest{Version: 2})
	require.NoError(err)
	require.Equal(uint64(40), resp.VotingPower)

	_, err = upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[1].String(), Version: 3})
	require.NoError(err)
	resp, err = upgradeKeeper.VersionTally(ctx, &types.QueryVersionTallyRequest{Version: 3})
	require.NoError(err)
	require.Equal(uint64(1), resp.VotingPower)

	_, err = upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[2].String(), Version: 2})
	require.NoError(err)
	_, err = upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[3].String(), Version: 2})
	require.NoError(err)

	_, err = upgradeKeeper.TryUpgrade(ctx, &types.MsgTryUpgrade{})
	require.NoError(err)

	require.True(upgradeKeeper.IsUpgradePending(ctx))

	upgradeKeeper.ResetTally(ctx)

	resp, err = upgradeKeeper.VersionTally(ctx, &types.QueryVersionTallyRequest{Version: 2})
	require.NoError(err)
	require.Equal(uint64(0), resp.VotingPower)

	resp, err = upgradeKeeper.VersionTally(ctx, &types.QueryVersionTallyRequest{Version: 3})
	require.NoError(err)
	require.Equal(uint64(0), resp.VotingPower)

	require.False(upgradeKeeper.IsUpgradePending(ctx))
}

func (s *TestSuite) TestTryUpgrade(t *testing.T) {
	require := s.Require()

	ctx, upgradeKeeper := s.Ctx, s.UpgradeKeeper
	goCtx := sdk.WrapSDKContext(ctx)

	t.Run("should return an error if an upgrade is already pending", func(t *testing.T) {
		_, err := upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[0].String(), Version: 2})
		require.NoError(err)
		_, err = upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[1].String(), Version: 2})
		require.NoError(err)
		_, err = upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[2].String(), Version: 2})
		require.NoError(err)
		_, err = upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[3].String(), Version: 2})
		require.NoError(err)

		// This TryUpgrade should succeed.
		_, err = upgradeKeeper.TryUpgrade(goCtx, &types.MsgTryUpgrade{})
		require.NoError(err)

		// This TryUpgrade should fail because an upgrade is pending.
		_, err = upgradeKeeper.TryUpgrade(goCtx, &types.MsgTryUpgrade{})
		require.Error(err)
		require.ErrorIs(err, types.ErrUpgradePending)
	})

	t.Run("should return an error if quorum version is less than or equal to the current version", func(t *testing.T) {
		_, err := upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[0].String(), Version: 1})
		require.NoError(err)
		_, err = upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[1].String(), Version: 1})
		require.NoError(err)
		_, err = upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[2].String(), Version: 1})
		require.NoError(err)
		_, err = upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[3].String(), Version: 1})
		require.NoError(err)

		_, err = upgradeKeeper.TryUpgrade(goCtx, &types.MsgTryUpgrade{})
		require.Error(err)
		require.ErrorIs(err, types.ErrInvalidUpgradeVersion)
	})
}

func (s *TestSuite) TestGetUpgrade(t *testing.T) {
	require := s.Require()

	ctx, upgradeKeeper := s.Ctx, s.UpgradeKeeper
	goCtx := sdk.WrapSDKContext(ctx)

	t.Run("should return an empty upgrade if no upgrade is pending", func(t *testing.T) {
		got, err := upgradeKeeper.GetUpgrade(ctx, &types.QueryGetUpgradeRequest{})
		require.NoError(err)
		require.Nil(got.Upgrade)
	})

	t.Run("should return an upgrade if an upgrade is pending", func(t *testing.T) {
		_, err := upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[0].String(), Version: 2})
		require.NoError(err)
		_, err = upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[1].String(), Version: 2})
		require.NoError(err)
		_, err = upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[2].String(), Version: 2})
		require.NoError(err)
		_, err = upgradeKeeper.SignalVersion(ctx, &types.MsgSignalVersion{ValidatorAddress: testutil.ValAddrs[3].String(), Version: 2})
		require.NoError(err)

		// This TryUpgrade should succeed.
		_, err = upgradeKeeper.TryUpgrade(goCtx, &types.MsgTryUpgrade{})
		require.NoError(err)

		got, err := upgradeKeeper.GetUpgrade(ctx, &types.QueryGetUpgradeRequest{})
		require.NoError(err)
		require.Equal(2, got.Upgrade.AppVersion)
		require.Equal(keeper.DefaultUpgradeHeightDelay, got.Upgrade.UpgradeHeight)
	})
}

var _ types.StakingKeeper = (*mockStakingKeeper)(nil)

type mockStakingKeeper struct {
	totalVotingPower sdkmath.Int
	validators       map[string]int64
}

func newMockStakingKeeper(validators map[string]int64) *mockStakingKeeper {
	totalVotingPower := sdkmath.NewInt(0)
	for _, power := range validators {
		totalVotingPower = totalVotingPower.AddRaw(power)
	}
	return &mockStakingKeeper{
		totalVotingPower: totalVotingPower,
		validators:       validators,
	}
}

func (m *mockStakingKeeper) GetLastTotalPower(ctx context.Context) (sdkmath.Int, error) {
	return m.totalVotingPower, nil
}

func (m *mockStakingKeeper) GetLastValidatorPower(_ context.Context, addr sdk.ValAddress) (power int64, err error) {
	addrStr := addr.String()
	if power, ok := m.validators[addrStr]; ok {
		return power, nil
	}

	return 0, nil
}

func (m *mockStakingKeeper) GetValidator(_ context.Context, addr sdk.ValAddress) (validator stypes.Validator, err error) {
	addrStr := addr.String()
	if _, ok := m.validators[addrStr]; ok {
		return stypes.Validator{Status: stypes.Bonded}, nil
	}

	return stypes.Validator{}, errors.New("not found")
}
