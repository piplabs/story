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

func TestDistributeCDRRewardPool_DistributesAndClears(t *testing.T) {
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

	err := k.distributeCDRRewardPool(ctx)
	require.NoError(t, err)

	require.ElementsMatch(t, []int64{75, 25}, sent)

	_, err = k.CDRPartialSubmitCount.Get(ctx, cdrSubmitCountKey(val1))
	require.ErrorIs(t, err, collections.ErrNotFound)
	_, err = k.CDRPartialSubmitCount.Get(ctx, cdrSubmitCountKey(val2))
	require.ErrorIs(t, err, collections.ErrNotFound)

	_, err = k.CDRFeePoolBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}

func TestDistributeCDRRewardPool_ZeroPoolClearsCounts(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	val := common.HexToAddress("0xcccccccccccccccccccccccccccccccccccccccc")
	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val), 2))
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "0"))

	err := k.distributeCDRRewardPool(ctx)
	require.NoError(t, err)

	_, err = k.CDRPartialSubmitCount.Get(ctx, cdrSubmitCountKey(val))
	require.ErrorIs(t, err, collections.ErrNotFound)
}

func TestDistributeCDRRewardPool_RoundingDoesNotOverDistribute(t *testing.T) {
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

	err := k.distributeCDRRewardPool(ctx)
	require.NoError(t, err)

	require.LessOrEqual(t, total, int64(10))
}

func TestDistributeCDRRewardPool_WithSingleValidator(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	val := common.HexToAddress("0x7777777777777777777777777777777777777777")
	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val), 5))
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "42"))

	bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.CDRFeePoolName, gomock.Any(),
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(42)))).Return(nil)

	err := k.distributeCDRRewardPool(ctx)
	require.NoError(t, err)

	_, err = k.CDRFeePoolBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}

func TestDistributeCDRRewardPool_EmptyCountsNoop(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "10"))

	err := k.distributeCDRRewardPool(ctx)
	require.NoError(t, err)
}
