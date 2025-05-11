package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmstaking/types"
)

func TestGetParams(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

	params, err := esk.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, types.DefaultParams(), params)
}

func TestSetParams(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

	newMaxWithdrawalPerBlock := uint32(10)
	newMaxSweepPerBlock := uint32(100)
	newMinPartialWithdrawalAmount := uint64(100_000)

	params, err := esk.GetParams(ctx)
	require.NoError(t, err)

	// check existing params are not equal to new params
	require.NotEqual(t, newMaxWithdrawalPerBlock, params.MaxWithdrawalPerBlock)
	require.NotEqual(t, newMaxSweepPerBlock, params.MaxSweepPerBlock)
	require.NotEqual(t, newMinPartialWithdrawalAmount, params.MinPartialWithdrawalAmount)

	newParams := params
	// set new params
	newParams.MaxWithdrawalPerBlock = newMaxWithdrawalPerBlock
	newParams.MaxSweepPerBlock = newMaxSweepPerBlock
	newParams.MinPartialWithdrawalAmount = newMinPartialWithdrawalAmount
	require.NoError(t, esk.SetParams(ctx, newParams))

	// check new params are set correctly
	params, err = esk.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, newParams, params)
}

func TestMaxWithdrawalPerBlock(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

	// MaxWithdrawalPerBlock check
	maxWithdrawalPerBlock, err := esk.MaxWithdrawalPerBlock(ctx)
	require.NoError(t, err)
	require.Equal(t, types.DefaultMaxWithdrawalPerBlock, maxWithdrawalPerBlock)
}

func TestMaxSweepPerBlock(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

	// MaxSweepPerBlock check
	maxSweepPerBlock, err := esk.MaxSweepPerBlock(ctx)
	require.NoError(t, err)
	require.Equal(t, types.DefaultMaxSweepPerBlock, maxSweepPerBlock)
}

func TestMinPartialWithdrawalAmount(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

	// MinPartialWithdrawalAmount check
	minPartialWithdrawalAmount, err := esk.MinPartialWithdrawalAmount(ctx)
	require.NoError(t, err)
	require.Equal(t, types.DefaultMinPartialWithdrawalAmount, minPartialWithdrawalAmount)
}

func TestSetValidatorSweepIndex(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
	existingSweepIndex, err := esk.GetValidatorSweepIndex(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(0), existingSweepIndex.NextValIndex)
	require.Equal(t, uint64(0), existingSweepIndex.NextValDelIndex)

	// set new sweep index
	newNextValIndex := uint64(10)
	newNextValDelIndex := uint64(100)
	// make sure new value is different from existing value
	require.NotEqual(t, existingSweepIndex.NextValIndex, newNextValIndex)
	require.NotEqual(t, existingSweepIndex.NextValDelIndex, newNextValDelIndex)
	require.NoError(t, esk.SetValidatorSweepIndex(ctx, types.NewValidatorSweepIndex(newNextValIndex, newNextValDelIndex)))

	// check new sweep index is set correctly
	sweepIndex, err := esk.GetValidatorSweepIndex(ctx)
	require.NoError(t, err)
	require.Equal(t, newNextValIndex, sweepIndex.NextValIndex)
	require.Equal(t, newNextValDelIndex, sweepIndex.NextValDelIndex)
}

func TestGetOldValidatorSweepIndex(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

	nextValIndex, err := esk.GetOldValidatorSweepIndex(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(0), nextValIndex)
}
