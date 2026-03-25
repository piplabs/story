package keeper

import (
	"context"
	"math/big"
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/server/utils"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"

	"go.uber.org/mock/gomock"
)

func TestAddCDRFeeToPool_MintAndTransfer(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	amount := big.NewInt(100)
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100)))

	bk.EXPECT().MintCoins(ctx, types.ModuleName, coins).Return(nil)
	bk.EXPECT().SendCoinsFromModuleToModule(ctx, types.ModuleName, types.CDRFeePoolName, coins).Return(nil)

	err := k.AddCDRFeeToPool(ctx, amount)
	require.NoError(t, err)

	bal, found, err := k.getCDRFeePoolBalance(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.True(t, bal.Equal(math.NewInt(100)))
}

func TestAddCDRFeeToPool_AccumulatesBalance(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "5"))

	amount := big.NewInt(7)
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(7)))

	bk.EXPECT().MintCoins(ctx, types.ModuleName, coins).Return(nil)
	bk.EXPECT().SendCoinsFromModuleToModule(ctx, types.ModuleName, types.CDRFeePoolName, coins).Return(nil)

	err := k.AddCDRFeeToPool(ctx, amount)
	require.NoError(t, err)

	bal, found, err := k.getCDRFeePoolBalance(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.True(t, bal.Equal(math.NewInt(12)))
}

func TestRefundCDRFee_TransfersAndUpdatesBalance(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	validator := common.HexToAddress("0x1111111111111111111111111111111111111111")
	recipient, err := utils.EvmAddressToBech32AccAddress(validator.Hex())
	require.NoError(t, err)

	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "10"))

	amount := big.NewInt(4)
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(4)))

	bk.EXPECT().SendCoinsFromModuleToAccount(ctx, types.CDRFeePoolName, recipient, coins).Return(nil)

	err = k.RefundCDRFee(ctx, validator, amount)
	require.NoError(t, err)

	bal, found, err := k.getCDRFeePoolBalance(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.True(t, bal.Equal(math.NewInt(6)))
}

func TestRefundCDRFee_ZeroBalanceRemoves(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	validator := common.HexToAddress("0x2222222222222222222222222222222222222222")
	recipient, err := utils.EvmAddressToBech32AccAddress(validator.Hex())
	require.NoError(t, err)

	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "5"))

	amount := big.NewInt(5)
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(5)))

	bk.EXPECT().SendCoinsFromModuleToAccount(ctx, types.CDRFeePoolName, recipient, coins).Return(nil)

	err = k.RefundCDRFee(ctx, validator, amount)
	require.NoError(t, err)

	_, err = k.CDRFeePoolBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}

func TestIncrementCDRPartialSubmitCount(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	validator := common.HexToAddress("0x3333333333333333333333333333333333333333")

	require.NoError(t, k.IncrementCDRPartialSubmitCount(ctx, validator))
	require.NoError(t, k.IncrementCDRPartialSubmitCount(ctx, validator))

	count, err := k.CDRPartialSubmitCount.Get(ctx, cdrSubmitCountKey(validator))
	require.NoError(t, err)
	require.Equal(t, uint64(2), count)
}

func TestDistributeCDRFee_DistributesAndClears(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	val1 := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	val2 := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val1), 3))
	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val2), 1))
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "100"))

	var sent []int64

	bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.CDRFeePoolName, gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ string, _ sdk.AccAddress, coins sdk.Coins) error {
			sent = append(sent, coins[0].Amount.Int64())
			return nil
		}).Times(2)

	err := k.distributeCDRFee(ctx)
	require.NoError(t, err)

	require.ElementsMatch(t, []int64{75, 25}, sent)

	_, err = k.CDRPartialSubmitCount.Get(ctx, cdrSubmitCountKey(val1))
	require.ErrorIs(t, err, collections.ErrNotFound)
	_, err = k.CDRPartialSubmitCount.Get(ctx, cdrSubmitCountKey(val2))
	require.ErrorIs(t, err, collections.ErrNotFound)

	_, err = k.CDRFeePoolBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}

func TestDistributeCDRFee_ZeroPoolClearsCounts(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	val := common.HexToAddress("0xcccccccccccccccccccccccccccccccccccccccc")
	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val), 2))
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "0"))

	err := k.distributeCDRFee(ctx)
	require.NoError(t, err)

	_, err = k.CDRPartialSubmitCount.Get(ctx, cdrSubmitCountKey(val))
	require.ErrorIs(t, err, collections.ErrNotFound)
}

func TestDistributeCDRFee_RoundingDoesNotOverDistribute(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	val1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	val2 := common.HexToAddress("0x2222222222222222222222222222222222222222")
	val3 := common.HexToAddress("0x3333333333333333333333333333333333333333")

	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val1), 1))
	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val2), 1))
	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val3), 1))
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "10"))

	var total int64
	bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.CDRFeePoolName, gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ string, _ sdk.AccAddress, coins sdk.Coins) error {
			total += coins[0].Amount.Int64()
			return nil
		}).Times(3)

	err := k.distributeCDRFee(ctx)
	require.NoError(t, err)

	require.LessOrEqual(t, total, int64(10))
}

func TestDistributeCDRFee_WithSingleValidator(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	val := common.HexToAddress("0x7777777777777777777777777777777777777777")
	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val), 5))
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "42"))

	bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.CDRFeePoolName, gomock.Any(),
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(42)))).Return(nil)

	err := k.distributeCDRFee(ctx)
	require.NoError(t, err)

	_, err = k.CDRFeePoolBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}

func TestDistributeCDRFee_EmptyCountsNoop(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "10"))

	err := k.distributeCDRFee(ctx)
	require.NoError(t, err)
}

// TestParseCDRFeeAmount_NilAmount verifies parseCDRFeeAmount returns (zero, false, nil)
// when the input is nil.
func TestParseCDRFeeAmount_NilAmount(t *testing.T) {
	t.Parallel()

	result, ok, err := parseCDRFeeAmount(nil)
	require.NoError(t, err)
	require.False(t, ok, "nil amount should return ok=false")
	require.True(t, result.IsZero())
}

// TestParseCDRFeeAmount_ZeroAmount verifies parseCDRFeeAmount returns (zero, false, nil)
// when the input is zero.
func TestParseCDRFeeAmount_ZeroAmount(t *testing.T) {
	t.Parallel()

	result, ok, err := parseCDRFeeAmount(big.NewInt(0))
	require.NoError(t, err)
	require.False(t, ok, "zero amount should return ok=false")
	require.True(t, result.IsZero())
}

// TestParseCDRFeeAmount_PositiveAmount verifies parseCDRFeeAmount correctly
// converts a positive big.Int to math.Int.
func TestParseCDRFeeAmount_PositiveAmount(t *testing.T) {
	t.Parallel()

	result, ok, err := parseCDRFeeAmount(big.NewInt(42))
	require.NoError(t, err)
	require.True(t, ok, "positive amount should return ok=true")
	require.Equal(t, int64(42), result.Int64())
}

// TestAddCDRFeeToPool_NilAmount verifies AddCDRFeeToPool is a no-op (no mints)
// when called with a nil amount.
func TestAddCDRFeeToPool_NilAmount(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	// No mock expectations — MintCoins and SendCoinsFromModuleToModule should NOT be called

	err := k.AddCDRFeeToPool(ctx, nil)
	require.NoError(t, err)
}

// TestAddCDRFeeToPool_ZeroAmount verifies AddCDRFeeToPool is a no-op when called
// with a zero amount.
func TestAddCDRFeeToPool_ZeroAmount(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	// No mock expectations — should short-circuit

	err := k.AddCDRFeeToPool(ctx, big.NewInt(0))
	require.NoError(t, err)
}

// TestAddCDRFeeToPool_MintError verifies AddCDRFeeToPool returns an error when
// MintCoins fails.
func TestAddCDRFeeToPool_MintError(t *testing.T) {
	t.Parallel()

	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	amount := big.NewInt(50)
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(50)))

	bk.EXPECT().MintCoins(ctx, types.ModuleName, coins).Return(errors.New("mint failed"))

	err := k.AddCDRFeeToPool(ctx, amount)
	require.Error(t, err)
	require.Contains(t, err.Error(), "mint CDR fee coins")
}

// TestAddCDRFeeToPool_SendError verifies AddCDRFeeToPool returns an error when
// SendCoinsFromModuleToModule fails.
func TestAddCDRFeeToPool_SendError(t *testing.T) {
	t.Parallel()

	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	amount := big.NewInt(50)
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(50)))

	bk.EXPECT().MintCoins(ctx, types.ModuleName, coins).Return(nil)
	bk.EXPECT().SendCoinsFromModuleToModule(ctx, types.ModuleName, types.CDRFeePoolName, coins).Return(errors.New("send failed"))

	err := k.AddCDRFeeToPool(ctx, amount)
	require.Error(t, err)
	require.Contains(t, err.Error(), "transfer CDR fee coins to pool")
}

// TestRefundCDRFee_NilAmount verifies RefundCDRFee is a no-op when called
// with a nil amount.
func TestRefundCDRFee_NilAmount(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	validator := common.HexToAddress("0x4444444444444444444444444444444444444444")

	err := k.RefundCDRFee(ctx, validator, nil)
	require.NoError(t, err)
}

// TestRefundCDRFee_PoolNotFound verifies RefundCDRFee returns an error when
// the CDR fee pool balance has not been set yet.
func TestRefundCDRFee_PoolNotFound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	validator := common.HexToAddress("0x5555555555555555555555555555555555555555")

	// Pool balance not set → getCDRFeePoolBalance returns found=false
	err := k.RefundCDRFee(ctx, validator, big.NewInt(10))
	require.Error(t, err)
	require.Contains(t, err.Error(), "cdr fee pool balance not found")
}

// TestRefundCDRFee_BalanceUnderflow verifies RefundCDRFee returns an error
// when the refund amount exceeds the current pool balance.
func TestRefundCDRFee_BalanceUnderflow(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	validator := common.HexToAddress("0x6666666666666666666666666666666666666666")

	// Pool balance = 5, refund = 100 → underflow
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "5"))

	err := k.RefundCDRFee(ctx, validator, big.NewInt(100))
	require.Error(t, err)
	require.Contains(t, err.Error(), "cdr fee pool balance underflow")
}

// TestRefundCDRFee_SendError verifies RefundCDRFee returns an error when
// SendCoinsFromModuleToAccount fails.
func TestRefundCDRFee_SendError(t *testing.T) {
	t.Parallel()

	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	validator := common.HexToAddress("0x7777777777777777777777777777777777777777")

	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "50"))

	amount := big.NewInt(10)
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(10)))

	bk.EXPECT().SendCoinsFromModuleToAccount(ctx, types.CDRFeePoolName, gomock.Any(), coins).Return(errors.New("send failed"))

	err := k.RefundCDRFee(ctx, validator, amount)
	require.Error(t, err)
	require.Contains(t, err.Error(), "refund CDR fee coins")
}

// TestGetCDRFeePoolBalance_CorruptData verifies getCDRFeePoolBalance returns an error
// when the stored balance string is not a valid integer.
func TestGetCDRFeePoolBalance_CorruptData(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Store an invalid balance string directly to simulate data corruption
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "not-a-number"))

	_, _, err := k.getCDRFeePoolBalance(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid CDR fee pool balance")
}

// TODO_CDR068: Characterization test for audit finding CDR-068.
//
// BUG LOCATION: client/x/dkg/keeper/dkg_cdr_fees.go:181 — distributeCDRFee()
//
//	`for addr, count := range counts` iterates a Go map with non-deterministic
//	order. Different validators produce different state hashes, causing
//	consensus split.
//
// CURRENT BEHAVIOR (BUG): Distribution order depends on map iteration, which
//
//	is non-deterministic. The test documents that addresses are NOT sorted.
//
// EXPECTED BEHAVIOR AFTER FIX: Addresses should be sorted before iteration.
//
// HOW TO UPDATE AFTER FIX:
//  1. Remove TODO_CDR068_ prefix from function name
//  2. Change assertion: verify addresses are distributed in sorted order
func TestTODO_CDR068_DistributeCDRFee_NonDeterministicOrder(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	val1 := common.HexToAddress("0xcccccccccccccccccccccccccccccccccccccccc")
	val2 := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	val3 := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val1), 1))
	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val2), 1))
	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val3), 1))
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "300"))

	var distributionOrder []string
	bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.CDRFeePoolName, gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ string, addr sdk.AccAddress, _ sdk.Coins) error {
			distributionOrder = append(distributionOrder, addr.String())
			return nil
		}).Times(3)

	err := k.distributeCDRFee(ctx)
	require.NoError(t, err)
	require.Len(t, distributionOrder, 3)

	t.Logf("CDR-068 characterization:")
	t.Logf("  Distribution order: %v", distributionOrder)
	t.Logf("  BUG: order depends on map iteration, not sorted")
	t.Logf("  Impact: different validators produce different state hashes -> consensus split")
	t.Logf("  After fix: addresses should be distributed in sorted (deterministic) order")
}

// TestDistributeCDRRewardPool_NoPrevActive verifies distributeCDRFee is a
// no-op when no previous active DKG network exists.
func TestDistributeCDRRewardPool_NoPrevActive(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// No active round set → should be a no-op
	err := k.distributeCDRFee(ctx)
	require.NoError(t, err)
}

// TestDistributeCDRRewardPool_PoolFoundButZeroBalance verifies that when the pool
// balance entry exists but is zero, counts are cleared and the zero balance entry
// is removed.
func TestDistributeCDRRewardPool_PoolFoundButZeroBalance(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	val := common.HexToAddress("0xdddddddddddddddddddddddddddddddddddddddd")
	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val), 3))

	// Set balance to "0" (found but zero)
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "0"))

	err := k.distributeCDRFee(ctx)
	require.NoError(t, err)

	// Balance entry should be removed and count cleared
	_, err = k.CDRFeePoolBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
	_, err = k.CDRPartialSubmitCount.Get(ctx, cdrSubmitCountKey(val))
	require.ErrorIs(t, err, collections.ErrNotFound)
}
