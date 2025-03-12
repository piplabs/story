package keeper_test

import (
	"time"

	"github.com/piplabs/story/client/x/evmstaking/types"
)

func (s *TestSuite) TestGetParams() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper
	params, err := keeper.GetParams(ctx)
	require.NoError(err)
	require.Equal(types.DefaultParams(), params)
}

func (s *TestSuite) TestSetParams() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper
	newMaxWithdrawalPerBlock := uint32(10)
	newMaxSweepPerBlock := uint32(100)
	newMinPartialWithdrawalAmount := uint64(100_000)
	newRefundFeeBps := uint32(10)
	newRefundPeriod := time.Duration(100)

	params, err := keeper.GetParams(ctx)
	require.NoError(err)

	// check existing params are not equal to new params
	require.NotEqual(newMaxWithdrawalPerBlock, params.MaxWithdrawalPerBlock)
	require.NotEqual(newMaxSweepPerBlock, params.MaxSweepPerBlock)
	require.NotEqual(newMinPartialWithdrawalAmount, params.MinPartialWithdrawalAmount)
	require.NotEqual(newRefundFeeBps, params.RefundFeeBps)
	require.NotEqual(newRefundPeriod, params.RefundPeriod)

	newParams := params
	// set new params
	newParams.MaxWithdrawalPerBlock = newMaxWithdrawalPerBlock
	newParams.MaxSweepPerBlock = newMaxSweepPerBlock
	newParams.MinPartialWithdrawalAmount = newMinPartialWithdrawalAmount
	newParams.RefundFeeBps = newRefundFeeBps
	newParams.RefundPeriod = newRefundPeriod
	require.NoError(keeper.SetParams(ctx, newParams))

	// check new params are set correctly
	params, err = keeper.GetParams(ctx)
	require.NoError(err)
	require.Equal(newParams, params)
}

func (s *TestSuite) TestMaxWithdrawalPerBlock() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper

	// params are set default during TestSuite.SetupTest
	params, err := keeper.GetParams(ctx)
	require.NoError(err)
	require.Equal(types.DefaultMaxWithdrawalPerBlock, params.MaxWithdrawalPerBlock)
}

func (s *TestSuite) TestMaxSweepPerBlock() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper

	// params are set default during TestSuite.SetupTest
	params, err := keeper.GetParams(ctx)
	require.NoError(err)
	require.Equal(types.DefaultMaxSweepPerBlock, params.MaxSweepPerBlock)
}

func (s *TestSuite) TestMinPartialWithdrawalAmount() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper

	// params are set default during TestSuite.SetupTest
	params, err := keeper.GetParams(ctx)
	require.NoError(err)
	require.Equal(types.DefaultMinPartialWithdrawalAmount, params.MinPartialWithdrawalAmount)
}

func (s *TestSuite) TestRefundFeeBps() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper

	// params are set default during TestSuite.SetupTest
	params, err := keeper.GetParams(ctx)
	require.NoError(err)
	require.Equal(types.DefaultRefundFeeBps, params.RefundFeeBps)
}

func (s *TestSuite) TestRefundPeriod() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper

	// params are set default during TestSuite.SetupTest
	params, err := keeper.GetParams(ctx)
	require.NoError(err)
	require.Equal(types.DefaultRefundPeriod, params.RefundPeriod)
}

func (s *TestSuite) TestSetValidatorSweepIndex() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper
	existingSweepIndex, err := keeper.GetValidatorSweepIndex(ctx)
	require.NoError(err)

	// set new sweep index
	newNextValIndex := uint64(10)
	newNextValDelIndex := uint64(100)
	// make sure new value is different from existing value
	require.NotEqual(existingSweepIndex.NextValIndex, newNextValIndex)
	require.NotEqual(existingSweepIndex.NextValDelIndex, newNextValDelIndex)
	require.NoError(keeper.SetValidatorSweepIndex(ctx, types.NewValidatorSweepIndex(newNextValIndex, newNextValDelIndex)))

	// check new sweep index is set correctly
	sweepIndex, err := keeper.GetValidatorSweepIndex(ctx)
	require.NoError(err)
	require.Equal(newNextValIndex, sweepIndex.NextValIndex)
	require.Equal(newNextValDelIndex, sweepIndex.NextValDelIndex)
}

func (s *TestSuite) TestGetValidatorSweepIndex() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper
	sweepIndex, err := keeper.GetValidatorSweepIndex(ctx)
	require.NoError(err)
	require.Equal(uint64(0), sweepIndex.NextValIndex)
	require.Equal(uint64(0), sweepIndex.NextValDelIndex)
}

func (s *TestSuite) TestGetOldValidatorSweepIndex() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper
	nextValIndex, err := keeper.GetOldValidatorSweepIndex(ctx)
	require.NoError(err)
	require.Equal(uint64(0), nextValIndex)
}
