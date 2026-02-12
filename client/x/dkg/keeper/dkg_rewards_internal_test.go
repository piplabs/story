package keeper

import (
	"context"
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"

	"go.uber.org/mock/gomock"
)

// ---------- settleRewardsForPreviousCommittee ----------

func TestDistributeDKGCommitteeRewards_NoPreviousActiveRound(t *testing.T) {
	// First round ever — no previous active round set. Should return nil immediately.

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.NoError(t, err)

	// No settlement balance should be set.
	_, err = k.SettlementBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}

func TestDistributeDKGCommitteeRewards_ZeroRewardPortion(t *testing.T) {
	// Reward portion is zero — should return nil without touching UBI.

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	params := types.DefaultParams()
	params.DkgCommitteeRewardPortion = math.LegacyZeroDec()
	require.NoError(t, k.SetParams(ctx, params))

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.NoError(t, err)

	_, err = k.SettlementBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}

func TestDistributeDKGCommitteeRewards_ZeroUbiBalance(t *testing.T) {
	// UBI balance is zero — should skip distribution gracefully.

	k, _, dk, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).Return(math.ZeroInt(), nil)

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.NoError(t, err)

	_, err = k.SettlementBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}

func TestDistributeDKGCommitteeRewards_NoFinalizedMembers(t *testing.T) {
	// Previous active round exists, but no members have DKGRegStatusFinalized.
	// UBI is withdrawn, distributeFromModuleBalance returns 0, all goes to settlement balance.

	k, _, dk, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	member1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	setRegistration(t, k, ctx, prevActive, member1, types.DKGRegStatusVerified)

	ubiBalance := math.NewInt(5000)
	dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).Return(ubiBalance, nil)
	dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), sdk.DefaultBondDenom, types.ModuleName).
		Return(sdk.NewCoin(sdk.DefaultBondDenom, ubiBalance), nil)

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.NoError(t, err)

	// Full amount should be stored as settlement balance.
	balStr, err := k.SettlementBalance.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, "5000", balStr)
}

func TestDistributeDKGCommitteeRewards_NormalDistribution(t *testing.T) {
	// Normal case: 3 finalized members, 10% reward portion, 10000 UBI balance.
	// Expected: perMember = (10000 * 0.10) / 3 = 333, totalDistributed = 999, remaining = 9001.

	k, bk, dk, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	member1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	member2 := common.HexToAddress("0x2222222222222222222222222222222222222222")
	member3 := common.HexToAddress("0x3333333333333333333333333333333333333333")

	setRegistration(t, k, ctx, prevActive, member1, types.DKGRegStatusFinalized)
	setRegistration(t, k, ctx, prevActive, member2, types.DKGRegStatusFinalized)
	setRegistration(t, k, ctx, prevActive, member3, types.DKGRegStatusFinalized)

	ubiBalance := math.NewInt(10000)

	dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).Return(ubiBalance, nil)
	dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), sdk.DefaultBondDenom, types.ModuleName).
		Return(sdk.NewCoin(sdk.DefaultBondDenom, ubiBalance), nil)

	perMemberReward := math.NewInt(333)
	perMemberCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, perMemberReward))
	bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, gomock.Any(), perMemberCoins).
		Return(nil).Times(3)

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.NoError(t, err)

	// remaining = 10000 - 999 = 9001, stored as settlement balance.
	balStr, err := k.SettlementBalance.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, "9001", balStr)
}

func TestDistributeDKGCommitteeRewards_SingleMember(t *testing.T) {
	// Single finalized member gets the full DKG reward portion.

	k, bk, dk, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	member1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	setRegistration(t, k, ctx, prevActive, member1, types.DKGRegStatusFinalized)

	ubiBalance := math.NewInt(10000)

	dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).Return(ubiBalance, nil)
	dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), sdk.DefaultBondDenom, types.ModuleName).
		Return(sdk.NewCoin(sdk.DefaultBondDenom, ubiBalance), nil)

	perMemberReward := math.NewInt(1000)
	perMemberCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, perMemberReward))
	bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, gomock.Any(), perMemberCoins).
		Return(nil).Times(1)

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.NoError(t, err)

	// remaining = 10000 - 1000 = 9000
	balStr, err := k.SettlementBalance.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, "9000", balStr)
}

func TestDistributeDKGCommitteeRewards_DustHandling(t *testing.T) {
	// Test that integer division truncation dust goes to the remaining UBI.
	// 7 members, 10% of 100 = 10, perMember = 10/7 = 1, totalDistributed = 7, remaining = 93.

	k, bk, dk, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	for i := 1; i <= 7; i++ {
		addr := common.BigToAddress(common.Big1.Add(common.Big1, common.Big0).SetInt64(int64(i)))
		setRegistration(t, k, ctx, prevActive, addr, types.DKGRegStatusFinalized)
	}

	ubiBalance := math.NewInt(100)

	dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).Return(ubiBalance, nil)
	dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), sdk.DefaultBondDenom, types.ModuleName).
		Return(sdk.NewCoin(sdk.DefaultBondDenom, ubiBalance), nil)

	perMemberReward := math.NewInt(1)
	perMemberCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, perMemberReward))
	bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, gomock.Any(), perMemberCoins).
		Return(nil).Times(7)

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.NoError(t, err)

	// remaining = 100 - 7 = 93
	balStr, err := k.SettlementBalance.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, "93", balStr)
}

func TestDistributeDKGCommitteeRewards_PerMemberRewardZero(t *testing.T) {
	// When the per-member reward truncates to zero (e.g. 1 UBI, 10%, 100 members),
	// no sends should happen and remaining should equal the full withdrawn amount.

	k, _, dk, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	for i := 1; i <= 100; i++ {
		addr := common.BigToAddress(common.Big0.SetInt64(int64(i)))
		setRegistration(t, k, ctx, prevActive, addr, types.DKGRegStatusFinalized)
	}

	ubiBalance := math.NewInt(1)

	dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).Return(ubiBalance, nil)
	dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), sdk.DefaultBondDenom, types.ModuleName).
		Return(sdk.NewCoin(sdk.DefaultBondDenom, ubiBalance), nil)

	// perMember is 0 (truncation), no SendCoinsFromModuleToAccount expected.
	// remaining = 1 - 0 = 1, stored as settlement balance.

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.NoError(t, err)

	balStr, err := k.SettlementBalance.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, "1", balStr)
}

func TestDistributeDKGCommitteeRewards_WithdrawnAmountZeroAfterWithdraw(t *testing.T) {
	// Edge case: UBI balance check says nonzero, but actual withdrawal returns zero.
	// Should skip distribution gracefully.

	k, _, dk, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	member1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	setRegistration(t, k, ctx, prevActive, member1, types.DKGRegStatusFinalized)

	dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).Return(math.NewInt(100), nil)
	dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), sdk.DefaultBondDenom, types.ModuleName).
		Return(sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()), nil)

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.NoError(t, err)

	_, err = k.SettlementBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}

// ---------- Error propagation tests ----------

func TestDistributeDKGCommitteeRewards_ErrorGetLatestActiveDKGNetwork(t *testing.T) {
	k, _, dk, ctx := setupDKGKeeperWithMocks(t)

	// Set a key pointing to a nonexistent DKG network. getLatestActiveDKGNetwork
	// finds the key but fails to find the network, which wraps with ErrNotFound,
	// so settleRewardsForPreviousCommittee enters the "no previous active" branch.
	require.NoError(t, k.LatestActiveRound.Set(ctx, "nonexistent_key"))

	// Since it falls through to the "no previous active" branch, no mocks needed.
	_ = dk

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.NoError(t, err)
}

func TestDistributeDKGCommitteeRewards_ErrorGetUbiBalance(t *testing.T) {
	k, _, dk, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).
		Return(math.Int{}, errors.New("distribution keeper error"))

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get UBI balance")
}

func TestDistributeDKGCommitteeRewards_ErrorWithdrawUbi(t *testing.T) {
	k, _, dk, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	member1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	setRegistration(t, k, ctx, prevActive, member1, types.DKGRegStatusFinalized)

	dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).Return(math.NewInt(1000), nil)
	dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), sdk.DefaultBondDenom, types.ModuleName).
		Return(sdk.Coin{}, errors.New("withdraw failed"))

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to withdraw UBI to DKG module")
}

func TestDistributeDKGCommitteeRewards_ErrorSendCoinsPartialFailure(t *testing.T) {
	// If SendCoinsFromModuleToAccount fails on the 2nd of 3 members,
	// the function should return an error.

	k, bk, dk, ctx := setupDKGKeeperWithMocks(t)

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	member1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	member2 := common.HexToAddress("0x2222222222222222222222222222222222222222")
	member3 := common.HexToAddress("0x3333333333333333333333333333333333333333")

	setRegistration(t, k, ctx, prevActive, member1, types.DKGRegStatusFinalized)
	setRegistration(t, k, ctx, prevActive, member2, types.DKGRegStatusFinalized)
	setRegistration(t, k, ctx, prevActive, member3, types.DKGRegStatusFinalized)

	ubiBalance := math.NewInt(10000)
	dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).Return(ubiBalance, nil)
	dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), sdk.DefaultBondDenom, types.ModuleName).
		Return(sdk.NewCoin(sdk.DefaultBondDenom, ubiBalance), nil)

	// First call succeeds, second fails.
	gomock.InOrder(
		bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, gomock.Any(), gomock.Any()).
			Return(nil),
		bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, gomock.Any(), gomock.Any()).
			Return(errors.New("bank send failed")),
	)

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "send DKG committee reward")
}

// ---------- ClaimSettlementBalance ----------

func TestClaimSettlementBalance_NoSettlement(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	amount, err := k.ClaimSettlementBalance(ctx, "evmstaking")
	require.NoError(t, err)
	require.True(t, amount.IsZero())
}

func TestClaimSettlementBalance_WithBalance(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	require.NoError(t, k.SettlementBalance.Set(ctx, "5000"))

	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(5000)))
	bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), types.ModuleName, "evmstaking", coins).Return(nil)

	amount, err := k.ClaimSettlementBalance(ctx, "evmstaking")
	require.NoError(t, err)
	require.Equal(t, math.NewInt(5000), amount)

	// Settlement balance should be cleared.
	_, err = k.SettlementBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}

func TestClaimSettlementBalance_ZeroBalance(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	require.NoError(t, k.SettlementBalance.Set(ctx, "0"))

	amount, err := k.ClaimSettlementBalance(ctx, "evmstaking")
	require.NoError(t, err)
	require.True(t, amount.IsZero())

	// Zero balance entry should be cleared.
	_, err = k.SettlementBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}

func TestClaimSettlementBalance_TransferError(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	require.NoError(t, k.SettlementBalance.Set(ctx, "1000"))

	bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), types.ModuleName, "evmstaking", gomock.Any()).
		Return(errors.New("transfer failed"))

	amount, err := k.ClaimSettlementBalance(ctx, "evmstaking")
	require.Error(t, err)
	require.Contains(t, err.Error(), "transfer settlement balance to recipient module")
	require.True(t, amount.IsZero())
}

// ---------- Determinism: sorted member ordering ----------

func TestDistributeDKGCommitteeRewards_DeterministicOrdering(t *testing.T) {
	for run := 0; run < 2; run++ {
		k, bk, dk, ctx := setupDKGKeeperWithMocks(t)

		prevActive := createTestDKGNetwork(t, k, ctx, 1)
		require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

		member1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
		member2 := common.HexToAddress("0x2222222222222222222222222222222222222222")
		member3 := common.HexToAddress("0x3333333333333333333333333333333333333333")

		// Insert in scrambled order.
		setRegistration(t, k, ctx, prevActive, member3, types.DKGRegStatusFinalized)
		setRegistration(t, k, ctx, prevActive, member1, types.DKGRegStatusFinalized)
		setRegistration(t, k, ctx, prevActive, member2, types.DKGRegStatusFinalized)

		ubiBalance := math.NewInt(30000)
		dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).Return(ubiBalance, nil)
		dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), sdk.DefaultBondDenom, types.ModuleName).
			Return(sdk.NewCoin(sdk.DefaultBondDenom, ubiBalance), nil)

		var recipientOrder []common.Address
		bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ interface{}, _ string, recipient sdk.AccAddress, _ sdk.Coins) error {
				recipientOrder = append(recipientOrder, common.BytesToAddress(recipient.Bytes()))
				return nil
			}).Times(3)

		err := k.settleRewardsForPreviousCommittee(ctx)
		require.NoError(t, err)
		require.Len(t, recipientOrder, 3)

		for i := 0; i < len(recipientOrder)-1; i++ {
			hexI := common.BytesToAddress(recipientOrder[i].Bytes()).Hex()
			hexJ := common.BytesToAddress(recipientOrder[i+1].Bytes()).Hex()
			require.True(t, hexI < hexJ,
				"run %d: recipients must be in sorted hex order: %s should come before %s",
				run, hexI, hexJ)
		}
	}
}

// ---------- Large reward portion (100%) ----------

func TestDistributeDKGCommitteeRewards_FullRewardPortion(t *testing.T) {
	// Reward portion = 1.0 (100%). All UBI goes to committee, remaining = 0.

	k, bk, dk, ctx := setupDKGKeeperWithMocks(t)

	params := types.DefaultParams()
	params.DkgCommitteeRewardPortion = math.LegacyOneDec()
	require.NoError(t, k.SetParams(ctx, params))

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	member1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	setRegistration(t, k, ctx, prevActive, member1, types.DKGRegStatusFinalized)

	ubiBalance := math.NewInt(5000)
	dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).Return(ubiBalance, nil)
	dk.EXPECT().WithdrawUbiByDenomToModule(gomock.Any(), sdk.DefaultBondDenom, types.ModuleName).
		Return(sdk.NewCoin(sdk.DefaultBondDenom, ubiBalance), nil)

	perMemberCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(5000)))
	bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, gomock.Any(), perMemberCoins).
		Return(nil)

	err := k.settleRewardsForPreviousCommittee(ctx)
	require.NoError(t, err)

	// remaining = 0, so no settlement balance should be set.
	_, err = k.SettlementBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}

// ---------- Arithmetic invariant: totalDistributed + remaining == withdrawnAmount ----------

func TestDistributeDKGCommitteeRewards_ArithmeticInvariant(t *testing.T) {
	testCases := []struct {
		name        string
		memberCount int
		ubiBalance  int64
		portion     string
	}{
		{"3 members, 10000 balance, 10%", 3, 10000, "0.10"},
		{"7 members, 100 balance, 10%", 7, 100, "0.10"},
		{"1 member, 1 balance, 50%", 1, 1, "0.50"},
		{"5 members, 999 balance, 33%", 5, 999, "0.33"},
		{"10 members, 1 balance, 10%", 10, 1, "0.10"},
		{"2 members, 1000000 balance, 99%", 2, 1000000, "0.99"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			portion := math.LegacyMustNewDecFromStr(tc.portion)
			withdrawn := math.NewInt(tc.ubiBalance)

			dkgReward := portion.MulInt(withdrawn).TruncateInt()
			memberCount := math.NewInt(int64(tc.memberCount))
			perMember := dkgReward.Quo(memberCount)

			totalDistributed := perMember.Mul(math.NewInt(int64(tc.memberCount)))
			remaining := withdrawn.Sub(totalDistributed)

			require.True(t, totalDistributed.Add(remaining).Equal(withdrawn),
				"totalDistributed(%s) + remaining(%s) must equal withdrawn(%s)",
				totalDistributed.String(), remaining.String(), withdrawn.String())

			require.True(t, totalDistributed.GTE(math.ZeroInt()), "totalDistributed must be non-negative")
			require.True(t, remaining.GTE(math.ZeroInt()), "remaining must be non-negative")
			require.True(t, perMember.GTE(math.ZeroInt()), "perMember must be non-negative")
		})
	}
}

// ---------- Test helpers ----------

// createTestDKGNetwork creates and stores a DKG network for testing.
func createTestDKGNetwork(t *testing.T, k *Keeper, ctx context.Context, round uint32) *types.DKGNetwork {
	t.Helper()

	codeCommitment := [32]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF,
		0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF,
		0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF,
		0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF}

	dkgNetwork := &types.DKGNetwork{
		CodeCommitment:   codeCommitment[:],
		Round:            round,
		StartBlockHeight: 100,
		StartBlockHash:   make([]byte, 32),
		ActiveValSet:     []string{},
		Total:            5,
		Threshold:        3,
		Stage:            types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, dkgNetwork))

	return dkgNetwork
}

// setRegistration creates and stores a DKG registration for a given member.
func setRegistration(t *testing.T, k *Keeper, ctx context.Context, network *types.DKGNetwork, addr common.Address, status types.DKGRegStatus) {
	t.Helper()

	reg := &types.DKGRegistration{
		Round:         network.Round,
		ValidatorAddr: addr.Hex(),
		Index:         1,
		DkgPubKey:     []byte("test-dkg-pubkey"),
		CommPubKey:    []byte("test-comm-pubkey"),
		RawQuote:      []byte("test-raw-quote"),
		Status:        status,
	}

	var codeCommitment [32]byte
	copy(codeCommitment[:], network.CodeCommitment)

	require.NoError(t, k.setDKGRegistration(ctx, codeCommitment, addr, reg))
}
