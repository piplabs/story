package keeper

import (
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"

	"go.uber.org/mock/gomock"
)

func TestFinalizeDKGRound_ThresholdChecks(t *testing.T) {
	testRound := uint32(1)
	validators := []common.Address{
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		common.HexToAddress("0x2222222222222222222222222222222222222222"),
		common.HexToAddress("0x3333333333333333333333333333333333333333"),
		common.HexToAddress("0x4444444444444444444444444444444444444444"),
		common.HexToAddress("0x5555555555555555555555555555555555555555"),
	}

	activeValSet := make([]string, len(validators))
	for i, v := range validators {
		activeValSet[i] = v.Hex()
	}

	// setFinalizedRegs creates N finalized registrations for the given round.
	setFinalizedRegs := func(t *testing.T, k *Keeper, ctx sdk.Context, n int) {
		t.Helper()

		for i := range n {
			reg := &types.DKGRegistration{
				Round:         testRound,
				ValidatorAddr: validators[i].Hex(),
				Index:         uint32(i + 1),
				DkgPubKey:     []byte("dkg-key"),
				CommPubKey:    []byte("comm-key"),
				EnclaveReport: []byte("enclave-report"),
				Status:        types.DKGRegStatusFinalized,
			}
			require.NoError(t, k.setDKGRegistration(ctx, validators[i], reg))
		}
	}

	tcs := []struct {
		name            string
		finalizedCount  int
		minReqFinalized uint32
		threshold       uint32
		total           uint32
		expectSkip      bool
	}{
		{
			name:            "skip: finalized count below min_req_finalized_participants",
			finalizedCount:  2,
			minReqFinalized: 3,
			threshold:       4,
			total:           5,
			expectSkip:      true,
		},
		{
			name:            "skip: meets min_req but below operational threshold",
			finalizedCount:  3,
			minReqFinalized: 3,
			threshold:       4,
			total:           5,
			expectSkip:      true,
		},
		{
			name:            "pass: meets both min_req and threshold",
			finalizedCount:  4,
			minReqFinalized: 3,
			threshold:       4,
			total:           5,
			expectSkip:      false,
		},
		{
			name:            "pass: finalized count exceeds threshold",
			finalizedCount:  5,
			minReqFinalized: 3,
			threshold:       4,
			total:           5,
			expectSkip:      false,
		},
		{
			name:            "pass: threshold equals min_req, both met exactly",
			finalizedCount:  3,
			minReqFinalized: 3,
			threshold:       3,
			total:           5,
			expectSkip:      false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			k, _, _, ctx := setupDKGKeeperWithMocks(t)
			sdkCtx := sdk.UnwrapSDKContext(ctx)

			// Set params
			params := types.DefaultParams()
			params.MinReqFinalizedParticipants = tc.minReqFinalized
			require.NoError(t, k.SetParams(ctx, params))

			// Set up DKG network
			latestRound := &types.DKGNetwork{
				Round:        testRound,
				ActiveValSet: activeValSet,
				Total:        tc.total,
				Threshold:    tc.threshold,
				Stage:        types.DKGStageFinalization,
			}
			require.NoError(t, k.setDKGNetwork(sdkCtx, latestRound))

			// Create finalized registrations
			setFinalizedRegs(t, k, sdkCtx, tc.finalizedCount)

			if tc.expectSkip {
				// SkipToNextRound -> InitiateDKGRound needs stakingKeeper.GetAllValidators
				sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
				sk.EXPECT().GetAllValidators(sdkCtx).Return(nil, nil).Times(1)

				err := k.FinalizeDKGRound(sdkCtx, latestRound)
				require.NoError(t, err)

				// Verify the round was marked as failed
				net, err := k.getLatestDKGNetwork(sdkCtx)
				require.NoError(t, err)
				require.Equal(t, types.DKGStageRegistration, net.Stage, "a new round should have been initiated")
				require.Equal(t, testRound+1, net.Round, "round number should be incremented")
			} else {
				err := k.FinalizeDKGRound(sdkCtx, latestRound)
				require.NoError(t, err)

				// Verify the round was set as active
				activeRound, err := k.GetLatestActiveRound(sdkCtx)
				require.NoError(t, err)
				require.NotNil(t, activeRound)
				require.Equal(t, testRound, activeRound.Round)
			}
		})
	}
}

func TestFinalizeDKGRound_DistributesCDRFeePool(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	params := types.DefaultParams()
	params.MinReqFinalizedParticipants = 1
	params.DkgCommitteeRewardPortion = math.LegacyZeroDec()
	require.NoError(t, k.SetParams(ctx, params))

	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	latestRound := &types.DKGNetwork{
		Round:        2,
		ActiveValSet: []string{},
		Total:        2,
		Threshold:    1,
		Stage:        types.DKGStageFinalization,
	}
	require.NoError(t, k.setDKGNetwork(sdkCtx, latestRound))

	val1 := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	val2 := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	setRegistration(t, k, ctx, latestRound, val1, types.DKGRegStatusFinalized)
	setRegistration(t, k, ctx, latestRound, val2, types.DKGRegStatusFinalized)

	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val1), 3))
	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val2), 1))
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "100"))

	var sent []int64
	bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.CDRFeePoolName, gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ interface{}, _ string, _ sdk.AccAddress, coins sdk.Coins) error {
			sent = append(sent, coins[0].Amount.Int64())
			return nil
		}).Times(2)

	err := k.FinalizeDKGRound(ctx, latestRound)
	require.NoError(t, err)

	require.ElementsMatch(t, []int64{75, 25}, sent)

	_, err = k.CDRPartialSubmitCount.Get(ctx, cdrSubmitCountKey(val1))
	require.ErrorIs(t, err, collections.ErrNotFound)
	_, err = k.CDRPartialSubmitCount.Get(ctx, cdrSubmitCountKey(val2))
	require.ErrorIs(t, err, collections.ErrNotFound)

	_, err = k.CDRFeePoolBalance.Get(ctx)
	require.ErrorIs(t, err, collections.ErrNotFound)
}
