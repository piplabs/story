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
	"github.com/piplabs/story/lib/netconf"

	"go.uber.org/mock/gomock"
)

// consensusGPK / consensusCoeffs are arbitrary fixed consensus key material used by
// the finalize-vote tests (only their byte identity matters, not their curve validity).
var (
	consensusGPK    = []byte{0xaa, 0xbb, 0xcc}
	consensusCoeffs = [][]byte{{0x01}, {0x02}}
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
		emptyGlobalKey  bool
	}{
		{
			name:            "skip: global public key not set",
			finalizedCount:  4,
			minReqFinalized: 3,
			threshold:       4,
			total:           5,
			expectSkip:      true,
			emptyGlobalKey:  true,
		},
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

			// Set up DKG network (consensus key material must be non-empty to proceed)
			var (
				globalPubKey []byte
				publicCoeffs [][]byte
			)
			if !tc.emptyGlobalKey {
				globalPubKey = []byte("global-pub-key")
				publicCoeffs = [][]byte{[]byte("coeff")}
			}
			latestRound := &types.DKGNetwork{
				Round:           testRound,
				ActiveValSet:    activeValSet,
				Total:           tc.total,
				Threshold:       tc.threshold,
				Stage:           types.DKGStageFinalization,
				GlobalPublicKey: globalPubKey,
				PublicCoeffs:    publicCoeffs,
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

// TestFinalizeDKGRound_UpgradeRound verifies that when IsUpgrade=true, FinalizeDKGRound
// deletes the activated upgrade info after a successful round finalization.
func TestFinalizeDKGRound_UpgradeRound(t *testing.T) {
	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	params := types.DefaultParams()
	params.MinReqFinalizedParticipants = 1
	params.DkgCommitteeRewardPortion = math.LegacyZeroDec()
	require.NoError(t, k.SetParams(ctx, params))

	latestRound := &types.DKGNetwork{
		Round:           1,
		ActiveValSet:    []string{},
		Total:           1,
		Threshold:       1,
		Stage:           types.DKGStageFinalization,
		IsUpgrade:       true, // upgrade resharing round
		GlobalPublicKey: []byte("global-pub-key"),
		PublicCoeffs:    [][]byte{[]byte("coeff")},
	}
	require.NoError(t, k.setDKGNetwork(sdkCtx, latestRound))

	val := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	setRegistration(t, k, ctx, latestRound, val, types.DKGRegStatusFinalized)

	// Store an activated upgrade info so FinalizeDKGRound can delete it
	upgradeInfo := &types.KernelUpgradeInfo{
		UpgradeVersion:   "v2.0.0",
		ActivationHeight: 100,
		IsActivated:      true,
	}
	require.NoError(t, k.SetKernelUpgradeInfo(ctx, upgradeInfo))

	err := k.FinalizeDKGRound(ctx, latestRound)
	require.NoError(t, err)

	// Verify upgrade info was deleted after successful upgrade round
	info, err := k.GetPendingUpgrade(ctx)
	require.NoError(t, err)
	require.Nil(t, info, "upgrade info should be deleted after successful upgrade round")

	// Verify the round is now active
	activeRound, err := k.GetLatestActiveRound(ctx)
	require.NoError(t, err)
	require.NotNil(t, activeRound)
	require.Equal(t, uint32(1), activeRound.Round)
}

// TestFinalizeDKGRound_DKGSvcEnabled verifies that FinalizeDKGRound spawns
// an async goroutine when isDKGSvcEnabled=true (without waiting for it to complete).
// The stateManager must be initialized so the goroutine (handleDKGComplete) does not panic.
func TestFinalizeDKGRound_DKGSvcEnabled(t *testing.T) {
	// Not parallel: modifies global DKG service state
	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	k.setIsDKGSvcEnabled()
	initTestStateManager(t, k)

	params := types.DefaultParams()
	params.MinReqFinalizedParticipants = 1
	params.DkgCommitteeRewardPortion = math.LegacyZeroDec()
	require.NoError(t, k.SetParams(ctx, params))

	latestRound := &types.DKGNetwork{
		Round:           9,
		ActiveValSet:    []string{},
		Total:           1,
		Threshold:       1,
		Stage:           types.DKGStageFinalization,
		IsUpgrade:       false,
		GlobalPublicKey: []byte("global-pub-key"),
		PublicCoeffs:    [][]byte{[]byte("coeff")},
	}
	require.NoError(t, k.setDKGNetwork(sdkCtx, latestRound))

	val := common.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
	setRegistration(t, k, ctx, latestRound, val, types.DKGRegStatusFinalized)

	// FinalizeDKGRound should succeed and the goroutine (handleDKGComplete) should be spawned.
	// The goroutine will log an error (session not found for round 9) but will not panic.
	err := k.FinalizeDKGRound(ctx, latestRound)
	require.NoError(t, err)

	// Verify round was set as active
	activeRound, err := k.GetLatestActiveRound(sdkCtx)
	require.NoError(t, err)
	require.NotNil(t, activeRound)
	require.Equal(t, uint32(9), activeRound.Round)
}

// TestDistributeCDRFee verifies that distributeCDRFee distributes the CDR fee
// pool proportionally to partial submission counts and clears state afterwards.
func TestDistributeCDRFee(t *testing.T) {
	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	// distributeCDRFee requires an active round to exist.
	prevActive := createTestDKGNetwork(t, k, ctx, 1)
	require.NoError(t, k.setLatestActiveRound(ctx, prevActive))

	val1 := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	val2 := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val1), 3))
	require.NoError(t, k.CDRPartialSubmitCount.Set(ctx, cdrSubmitCountKey(val2), 1))
	require.NoError(t, k.CDRFeePoolBalance.Set(ctx, "100"))

	var sent []int64
	bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.CDRFeePoolName, gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ interface{}, _ string, _ sdk.AccAddress, coins sdk.Coins) error {
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

// --- Tests merged from dkg_process_test.go ---

// TestBeginFinalization_DKGSvcDisabled verifies that BeginFinalization emits
// the event and returns nil when DKG service is disabled.

func TestBeginFinalization_DKGSvcDisabled(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageFinalization,
	}

	err := k.BeginFinalization(ctx, network)
	require.NoError(t, err)
}

// TestBeginFinalization_DKGSvcEnabled verifies that BeginFinalization emits
// the event even when the DKG service is enabled (async goroutine is launched).

func TestBeginFinalization_DKGSvcEnabled(t *testing.T) {
	// Not parallel: modifies global DKG service state
	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()

	network := &types.DKGNetwork{
		Round:     22,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageFinalization,
	}

	err := k.BeginFinalization(ctx, network)
	require.NoError(t, err)
}

// TestInvalidateNonConsensusFinalizations verifies the core sweep: finalized validators
// whose recorded finalize vote matches the consensus key are kept; those that diverge
// (or have no recorded vote) are invalidated, and all recorded votes are pruned.
func TestInvalidateNonConsensusFinalizations(t *testing.T) {
	const (
		round = uint32(7)
		n     = 5
	)

	validators := make([]common.Address, n)
	for i := range n {
		validators[i] = common.BytesToAddress([]byte{byte(i + 1)})
	}

	// Index divergentIdx voted for a different polynomial; index missingIdx has no vote.
	const divergentIdx, missingIdx = 2, 3

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx)

	latestRound := &types.DKGNetwork{
		Round:           round,
		Total:           n,
		Threshold:       3,
		Stage:           types.DKGStageFinalization,
		GlobalPublicKey: consensusGPK,
		PublicCoeffs:    consensusCoeffs,
	}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))

	consensusKey := globalPubKeyVoteKey(round, consensusGPK, consensusCoeffs)
	divergentKey := globalPubKeyVoteKey(round, consensusGPK, [][]byte{{0x09}}) // different coeffs

	for i := range n {
		reg := &types.DKGRegistration{
			Round:         round,
			ValidatorAddr: validators[i].Hex(),
			Index:         uint32(i + 1),
			Status:        types.DKGRegStatusFinalized,
		}
		require.NoError(t, k.setDKGRegistration(ctx, validators[i], reg))

		switch i {
		case missingIdx:
			// no recorded finalize vote
		case divergentIdx:
			require.NoError(t, k.FinalizeVotes.Set(ctx, finalizeVoteStoreKey(round, validators[i]), divergentKey))
		default:
			require.NoError(t, k.FinalizeVotes.Set(ctx, finalizeVoteStoreKey(round, validators[i]), consensusKey))
		}
	}

	require.NoError(t, k.invalidateNonConsensusFinalizations(ctx, latestRound))

	for i := range n {
		reg, err := k.getDKGRegistration(ctx, round, validators[i])
		require.NoError(t, err)

		if i == divergentIdx || i == missingIdx {
			require.Equal(t, types.DKGRegStatusInvalidated, reg.Status,
				"participant %d (non-consensus vote) should be invalidated", i+1)
		} else {
			require.Equal(t, types.DKGRegStatusFinalized, reg.Status,
				"participant %d (consensus vote) should remain finalized", i+1)
		}

		// Recorded votes are pruned.
		has, err := k.FinalizeVotes.Has(ctx, finalizeVoteStoreKey(round, validators[i]))
		require.NoError(t, err)
		require.False(t, has, "finalize vote for participant %d should be pruned", i+1)
	}
}

// TestSkipToNextRound_PrunesFinalizeVotes verifies that abandoning a v1.9.0 round that
// recorded finalize votes (but never reached consensus) prunes those votes, so they do
// not leak across rounds.
func TestSkipToNextRound_PrunesFinalizeVotes(t *testing.T) {
	const round = uint32(14)

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx).WithChainID(netconf.TestChainID)

	val := common.BytesToAddress([]byte{0x01})
	latestRound := &types.DKGNetwork{
		Round:            round,
		Total:            1,
		Threshold:        2, // never reached -> GlobalPublicKey stays empty -> guard skip
		Stage:            types.DKGStageFinalization,
		StartBlockHeight: 400, // v1.9.0 active
		ActiveValSet:     []string{val.Hex()},
	}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))

	// A validator finalized (recorded a vote) but consensus did not reach threshold.
	require.NoError(t, k.setDKGRegistration(ctx, val, &types.DKGRegistration{
		Round:         round,
		ValidatorAddr: val.Hex(),
		Index:         1,
		Status:        types.DKGRegStatusFinalized,
	}))
	require.NoError(t, k.FinalizeVotes.Set(ctx, finalizeVoteStoreKey(round, val), "some-vote-key"))

	// SkipToNextRound -> InitiateDKGRound needs stakingKeeper.GetAllValidators.
	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	sk.EXPECT().GetAllValidators(ctx).Return(nil, nil).Times(1)

	// GlobalPublicKey empty -> FinalizeDKGRound takes the guard skip path (no sweep).
	require.NoError(t, k.FinalizeDKGRound(ctx, latestRound))

	has, err := k.FinalizeVotes.Has(ctx, finalizeVoteStoreKey(round, val))
	require.NoError(t, err)
	require.False(t, has, "finalize vote must be pruned when the round is abandoned")
}

// TestFinalizeDKGRound_NoPublicCoeffs verifies that a round with no consensus public
// coefficients skips to the next round rather than activating.
func TestFinalizeDKGRound_NoPublicCoeffs(t *testing.T) {
	const round = uint32(8)

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx)

	latestRound := &types.DKGNetwork{
		Round:           round,
		Total:           1,
		Threshold:       1,
		Stage:           types.DKGStageFinalization,
		GlobalPublicKey: []byte("global-pub-key"),
		PublicCoeffs:    nil, // consensus polynomial missing
	}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))

	// SkipToNextRound -> InitiateDKGRound needs stakingKeeper.GetAllValidators.
	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	sk.EXPECT().GetAllValidators(ctx).Return(nil, nil).Times(1)

	require.NoError(t, k.FinalizeDKGRound(ctx, latestRound))

	net, err := k.getLatestDKGNetwork(ctx)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageRegistration, net.Stage, "a new round should have been initiated")
	require.Equal(t, round+1, net.Round)
}

// TestFinalizeDKGRound_V190Gate verifies that the non-consensus sweep only runs once the
// v1.9.0 upgrade height is reached, so pre-upgrade blocks replay deterministically.
func TestFinalizeDKGRound_V190Gate(t *testing.T) {
	const (
		round     = uint32(9)
		n         = 4
		threshold = 5 // intentionally > n so FinalizeDKGRound always takes the skip path
	)

	const divergentIdx = 1

	// setup builds a keeper with one divergent finalize vote in a round that started at
	// roundStartHeight (the gate is on the round's start height, not the current height).
	setup := func(t *testing.T, roundStartHeight int64) (*Keeper, sdk.Context, []common.Address) {
		t.Helper()

		k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
		ctx := sdk.UnwrapSDKContext(baseCtx).WithChainID(netconf.TestChainID)

		params := types.DefaultParams()
		params.MinReqFinalizedParticipants = 1
		require.NoError(t, k.SetParams(ctx, params))

		validators := make([]common.Address, n)
		activeValSet := make([]string, n)
		for i := range n {
			validators[i] = common.BytesToAddress([]byte{byte(i + 1)})
			activeValSet[i] = validators[i].Hex()
		}

		latestRound := &types.DKGNetwork{
			Round:            round,
			ActiveValSet:     activeValSet,
			Total:            n,
			Threshold:        threshold,
			Stage:            types.DKGStageFinalization,
			StartBlockHeight: roundStartHeight, // gate is on the round's start height
			GlobalPublicKey:  consensusGPK,
			PublicCoeffs:     consensusCoeffs,
		}
		require.NoError(t, k.setDKGNetwork(ctx, latestRound))

		consensusKey := globalPubKeyVoteKey(round, consensusGPK, consensusCoeffs)
		divergentKey := globalPubKeyVoteKey(round, consensusGPK, [][]byte{{0x09}})

		for i := range n {
			reg := &types.DKGRegistration{
				Round:         round,
				ValidatorAddr: validators[i].Hex(),
				Index:         uint32(i + 1),
				Status:        types.DKGRegStatusFinalized,
			}
			require.NoError(t, k.setDKGRegistration(ctx, validators[i], reg))

			voteKey := consensusKey
			if i == divergentIdx {
				voteKey = divergentKey
			}
			require.NoError(t, k.FinalizeVotes.Set(ctx, finalizeVoteStoreKey(round, validators[i]), voteKey))
		}

		// threshold (5) > finalized count, so FinalizeDKGRound skips -> InitiateDKGRound.
		sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
		sk.EXPECT().GetAllValidators(gomock.Any()).Return(nil, nil).Times(1)

		return k, ctx, validators
	}

	t.Run("gated off: round started below v1.9.0 height, divergent reg stays finalized", func(t *testing.T) {
		k, ctx, validators := setup(t, 399) // TestChainID V190 height is 400

		latestRound, err := k.getLatestDKGNetwork(ctx)
		require.NoError(t, err)
		require.NoError(t, k.FinalizeDKGRound(ctx, latestRound))

		reg, err := k.getDKGRegistration(ctx, round, validators[divergentIdx])
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusFinalized, reg.Status,
			"below v1.9.0 height the sweep must not run")
	})

	t.Run("gated on: round started at v1.9.0 height, divergent reg invalidated", func(t *testing.T) {
		k, ctx, validators := setup(t, 400)

		latestRound, err := k.getLatestDKGNetwork(ctx)
		require.NoError(t, err)
		require.NoError(t, k.FinalizeDKGRound(ctx, latestRound))

		divergent, err := k.getDKGRegistration(ctx, round, validators[divergentIdx])
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusInvalidated, divergent.Status,
			"at v1.9.0 height the divergent vote must be invalidated")

		// A validator that voted for consensus is untouched.
		good, err := k.getDKGRegistration(ctx, round, validators[0])
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusFinalized, good.Status)
	})
}

// TestMarkDealersDealt verifies that 0-based dealer indices [0, Total-1] are resolved to
// the dealer's validator address via the committee's registrations, and out-of-range
// indices are skipped. Non-resharing: the dealer committee is the current round.
func TestMarkDealersDealt(t *testing.T) {
	const (
		round = uint32(12)
		n     = 3
	)

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx)

	validators := make([]common.Address, n)
	latestRound := &types.DKGNetwork{Round: round, Total: n} // valid 0-based indices {0,1,2}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))
	for i := range n {
		validators[i] = common.BytesToAddress([]byte{byte(i + 1)})
		require.NoError(t, k.setDKGRegistration(ctx, validators[i], &types.DKGRegistration{
			Round:         round,
			ValidatorAddr: validators[i].Hex(),
			Index:         uint32(i + 1), // 1-based
			Status:        types.DKGRegStatusVerified,
		}))
	}

	deals := []types.Deal{
		{Index: 0}, {Index: 0}, // dealer 0 (duplicate) -> validators[0]
		{Index: 2}, // dealer 2 -> validators[2]
		{Index: 3}, // out of range (>= Total=3) -> skipped
	}
	require.NoError(t, k.markDealersDealt(ctx, latestRound, deals))

	for i, want := range map[int]bool{0: true, 1: false, 2: true} {
		has, err := k.DealtDealers.Has(ctx, dealtDealerKey(round, validators[i].Hex()))
		require.NoError(t, err)
		require.Equal(t, want, has, "validator %d", i)
	}
}

// TestMarkDealersDealt_AllDealersRecorded is a regression test for the 0-based/1-based
// off-by-one: when all n dealers deal (0-based indices {0..n-1}), every registration
// index {1..n} — including the highest, n — must be recorded, so invalidateMissingDealers
// keeps them all. A test built with 1-based deal indices would hide this.
func TestMarkDealersDealt_AllDealersRecorded(t *testing.T) {
	const (
		round = uint32(15)
		n     = 4
	)

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx)

	validators := make([]common.Address, n)
	latestRound := &types.DKGNetwork{Round: round, Total: n, Threshold: 2, Stage: types.DKGStageFinalization}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))

	deals := make([]types.Deal, n)
	for i := range n {
		validators[i] = common.BytesToAddress([]byte{byte(i + 1)})
		require.NoError(t, k.setDKGRegistration(ctx, validators[i], &types.DKGRegistration{
			Round:         round,
			ValidatorAddr: validators[i].Hex(),
			Index:         uint32(i + 1), // 1-based
			Status:        types.DKGRegStatusVerified,
		}))
		deals[i] = types.Deal{Index: uint32(i)} // 0-based dealer index
	}

	require.NoError(t, k.markDealersDealt(ctx, latestRound, deals))

	// Every dealer (incl. the highest index n) must be recorded.
	for i := range n {
		has, err := k.DealtDealers.Has(ctx, dealtDealerKey(round, validators[i].Hex()))
		require.NoError(t, err)
		require.True(t, has, "validator %d (highest reg.Index=%d) must be recorded as dealt", i+1, n)
	}

	require.NoError(t, k.invalidateMissingDealers(ctx, latestRound))

	for i := range n {
		reg, err := k.getDKGRegistration(ctx, round, validators[i])
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusVerified, reg.Status,
			"reg.Index %d must remain verified (all dealers dealt)", i+1)
	}
}

// TestInvalidateMissingDealers verifies that verified dealers without a recorded deal
// are invalidated, dealers that dealt are kept and their marks pruned, and already
// invalidated dealers are left untouched.
func TestInvalidateMissingDealers(t *testing.T) {
	const (
		round = uint32(11)
		n     = 5
	)

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx)

	validators := make([]common.Address, n)
	for i := range n {
		validators[i] = common.BytesToAddress([]byte{byte(i + 1)})
	}

	latestRound := &types.DKGNetwork{Round: round, Total: n, Threshold: 3, Stage: types.DKGStageFinalization}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))

	// reg index 5 (i=4) is already invalidated and never dealt; others are verified.
	for i := range n {
		status := types.DKGRegStatusVerified
		if i == 4 {
			status = types.DKGRegStatusInvalidated
		}
		reg := &types.DKGRegistration{
			Round:         round,
			ValidatorAddr: validators[i].Hex(),
			Index:         uint32(i + 1),
			Status:        status,
		}
		require.NoError(t, k.setDKGRegistration(ctx, validators[i], reg))
	}

	// 0-based dealer indices 0 and 2 submitted deals -> reg.Index 1 and 3 dealt; 2,4,5 did not.
	require.NoError(t, k.markDealersDealt(ctx, latestRound, []types.Deal{{Index: 0}, {Index: 2}}))

	require.NoError(t, k.invalidateMissingDealers(ctx, latestRound))

	wantStatus := map[int]types.DKGRegStatus{
		0: types.DKGRegStatusVerified,    // dealt
		1: types.DKGRegStatusInvalidated, // missing deal
		2: types.DKGRegStatusVerified,    // dealt
		3: types.DKGRegStatusInvalidated, // missing deal
		4: types.DKGRegStatusInvalidated, // already invalidated, untouched
	}
	for i := range n {
		reg, err := k.getDKGRegistration(ctx, round, validators[i])
		require.NoError(t, err)
		require.Equal(t, wantStatus[i], reg.Status, "reg index %d", i+1)
	}

	// Marks for dealt dealers are pruned (validators[0] and validators[2]).
	for _, i := range []int{0, 2} {
		has, err := k.DealtDealers.Has(ctx, dealtDealerKey(round, validators[i].Hex()))
		require.NoError(t, err)
		require.False(t, has, "dealt mark for validator %d should be pruned", i)
	}
}

// TestBeginFinalization_MissingDealerGate verifies the missing-dealer sweep only runs
// for rounds started at/after the v1.9.0 height.
func TestBeginFinalization_MissingDealerGate(t *testing.T) {
	const (
		round = uint32(13)
		n     = 3
	)

	// nonDealerIdx (0-based) is verified but never submits a deal.
	const nonDealerIdx = 1

	setup := func(t *testing.T, roundStartHeight int64) (*Keeper, sdk.Context, []common.Address) {
		t.Helper()

		k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
		ctx := sdk.UnwrapSDKContext(baseCtx).WithChainID(netconf.TestChainID)

		validators := make([]common.Address, n)
		for i := range n {
			validators[i] = common.BytesToAddress([]byte{byte(i + 1)})
		}

		latestRound := &types.DKGNetwork{
			Round:            round,
			Total:            n,
			Threshold:        2,
			Stage:            types.DKGStageFinalization,
			StartBlockHeight: roundStartHeight,
		}
		require.NoError(t, k.setDKGNetwork(ctx, latestRound))

		for i := range n {
			reg := &types.DKGRegistration{
				Round:         round,
				ValidatorAddr: validators[i].Hex(),
				Index:         uint32(i + 1),
				Status:        types.DKGRegStatusVerified,
			}
			require.NoError(t, k.setDKGRegistration(ctx, validators[i], reg))
		}

		// Every dealer except nonDealerIdx submitted a deal (0-based dealer index = i).
		var deals []types.Deal
		for i := range n {
			if i != nonDealerIdx {
				deals = append(deals, types.Deal{Index: uint32(i)})
			}
		}
		require.NoError(t, k.markDealersDealt(ctx, latestRound, deals))

		return k, ctx, validators
	}

	t.Run("gated off: round started below v1.9.0 height, missing dealer stays verified", func(t *testing.T) {
		k, ctx, validators := setup(t, 399)

		latestRound, err := k.getLatestDKGNetwork(ctx)
		require.NoError(t, err)
		require.NoError(t, k.BeginFinalization(ctx, latestRound))

		reg, err := k.getDKGRegistration(ctx, round, validators[nonDealerIdx])
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusVerified, reg.Status)
	})

	t.Run("gated on: round started at v1.9.0 height, missing dealer invalidated", func(t *testing.T) {
		k, ctx, validators := setup(t, 400)

		latestRound, err := k.getLatestDKGNetwork(ctx)
		require.NoError(t, err)
		require.NoError(t, k.BeginFinalization(ctx, latestRound))

		missing, err := k.getDKGRegistration(ctx, round, validators[nonDealerIdx])
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusInvalidated, missing.Status)

		dealt, err := k.getDKGRegistration(ctx, round, validators[0])
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusVerified, dealt.Status)
	})
}

// TestInvalidateMissingDealers_Resharing covers the resharing index-space defect: in a
// resharing round the dealers are the PREVIOUS active committee, whose index space is
// permuted differently from the current round. Deals must be attributed by validator
// address, never by raw index. Reproduces the devnet 4->3 node-removal scenario.
func TestInvalidateMissingDealers_Resharing(t *testing.T) {
	const (
		oldRound = uint32(40)
		newRound = uint32(42)
	)

	// Old committee (the dealers), by old reg.Index: val1=1(kyber0), val3=2(kyber1),
	// val4=3(kyber2), val2=4(kyber3). val3 is stopped and submits no deal.
	val1 := common.BytesToAddress([]byte{0x01})
	val2 := common.BytesToAddress([]byte{0x02})
	val3 := common.BytesToAddress([]byte{0x03})
	val4 := common.BytesToAddress([]byte{0x04})
	oldByIndex := []common.Address{val1, val3, val4, val2} // index i -> kyber i

	// Deals from the dealers that are up (old kyber indices 0,2,3); val3 (kyber 1) is down.
	dealtDeals := []types.Deal{{Index: 0}, {Index: 2}, {Index: 3}}

	setupOld := func(t *testing.T, k *Keeper, ctx sdk.Context) {
		t.Helper()
		old := &types.DKGNetwork{
			Round:        oldRound,
			Total:        4,
			ActiveValSet: []string{val1.Hex(), val3.Hex(), val4.Hex(), val2.Hex()},
			Stage:        types.DKGStageActive,
		}
		require.NoError(t, k.setDKGNetwork(ctx, old))
		require.NoError(t, k.setLatestActiveRound(ctx, old))
		for i, addr := range oldByIndex {
			require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
				Round:         oldRound,
				ValidatorAddr: addr.Hex(),
				Index:         uint32(i + 1),
				Status:        types.DKGRegStatusFinalized,
			}))
		}
	}

	t.Run("leaving dealer not invalidated; healthy members kept despite index permutation", func(t *testing.T) {
		k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
		ctx := sdk.UnwrapSDKContext(baseCtx).WithChainID(netconf.TestChainID)
		setupOld(t, k, ctx)

		// New committee (val3 left): val2=1, val1=2, val4=3. Note val1 sits at the index
		// (2) the stopped val3 held in the old committee — the bug invalidated val1 here.
		newRegs := []common.Address{val2, val1, val4}
		newDKG := &types.DKGNetwork{
			Round: newRound, Total: 3, Threshold: 2,
			Stage: types.DKGStageFinalization, IsResharing: true, StartBlockHeight: 400,
		}
		require.NoError(t, k.setDKGNetwork(ctx, newDKG))
		for i, addr := range newRegs {
			require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
				Round:         newRound,
				ValidatorAddr: addr.Hex(),
				Index:         uint32(i + 1),
				Status:        types.DKGRegStatusVerified,
			}))
		}

		require.NoError(t, k.markDealersDealt(ctx, newDKG, dealtDeals))
		require.NoError(t, k.invalidateMissingDealers(ctx, newDKG))

		// All three new members dealt (val1/val4/val2) and stay verified; val3 left and
		// has no new registration to invalidate.
		for _, addr := range newRegs {
			reg, err := k.getDKGRegistration(ctx, newRound, addr)
			require.NoError(t, err)
			require.Equal(t, types.DKGRegStatusVerified, reg.Status, "validator %s must stay verified", addr.Hex())
		}
	})

	t.Run("continuing dealer that skips dealing is invalidated", func(t *testing.T) {
		k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
		ctx := sdk.UnwrapSDKContext(baseCtx).WithChainID(netconf.TestChainID)
		setupOld(t, k, ctx)

		// val3 does NOT leave: it registers for the new round but still submits no deal.
		newRegs := []common.Address{val2, val1, val4, val3}
		newDKG := &types.DKGNetwork{
			Round: newRound, Total: 4, Threshold: 2,
			Stage: types.DKGStageFinalization, IsResharing: true, StartBlockHeight: 400,
		}
		require.NoError(t, k.setDKGNetwork(ctx, newDKG))
		for i, addr := range newRegs {
			require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
				Round:         newRound,
				ValidatorAddr: addr.Hex(),
				Index:         uint32(i + 1),
				Status:        types.DKGRegStatusVerified,
			}))
		}

		require.NoError(t, k.markDealersDealt(ctx, newDKG, dealtDeals))
		require.NoError(t, k.invalidateMissingDealers(ctx, newDKG))

		// val3 was an old-committee dealer that skipped dealing and is still a member -> invalidated.
		reg3, err := k.getDKGRegistration(ctx, newRound, val3)
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusInvalidated, reg3.Status)

		for _, addr := range []common.Address{val1, val2, val4} {
			reg, err := k.getDKGRegistration(ctx, newRound, addr)
			require.NoError(t, err)
			require.Equal(t, types.DKGRegStatusVerified, reg.Status, "validator %s must stay verified", addr.Hex())
		}
	})
}

// TestInvalidateMissingDealers_ReshardingExpansion covers the predicted (in the PR review)
// node-addition failure: in a 3->4 expansion the joining node holds no old share and
// correctly submits no deal. The expected dealer set is the OLD committee, so the joiner is
// never a candidate for invalidation — regardless of which registration index it lands on.
// A raw-index implementation would invalidate whoever sits at index newTotal (here val4).
func TestInvalidateMissingDealers_ReshardingExpansion(t *testing.T) {
	const (
		oldRound = uint32(50)
		newRound = uint32(52)
	)

	// Old committee of 3: val1=1(kyber0), val2=2(kyber1), val3=3(kyber2). All deal.
	val1 := common.BytesToAddress([]byte{0x01})
	val2 := common.BytesToAddress([]byte{0x02})
	val3 := common.BytesToAddress([]byte{0x03})
	val4 := common.BytesToAddress([]byte{0x04}) // joining node, no old share
	oldByIndex := []common.Address{val1, val2, val3}

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx).WithChainID(netconf.TestChainID)

	old := &types.DKGNetwork{
		Round:        oldRound,
		Total:        3,
		ActiveValSet: []string{val1.Hex(), val2.Hex(), val3.Hex()},
		Stage:        types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, old))
	require.NoError(t, k.setLatestActiveRound(ctx, old))
	for i, addr := range oldByIndex {
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:         oldRound,
			ValidatorAddr: addr.Hex(),
			Index:         uint32(i + 1),
			Status:        types.DKGRegStatusFinalized,
		}))
	}

	// New committee of 4, joiner val4 registered LAST (index 4 == newTotal): the index the
	// buggy implementation would have killed.
	newRegs := []common.Address{val1, val2, val3, val4}
	newDKG := &types.DKGNetwork{
		Round: newRound, Total: 4, Threshold: 2,
		Stage: types.DKGStageFinalization, IsResharing: true, StartBlockHeight: 400,
	}
	require.NoError(t, k.setDKGNetwork(ctx, newDKG))
	for i, addr := range newRegs {
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:         newRound,
			ValidatorAddr: addr.Hex(),
			Index:         uint32(i + 1),
			Status:        types.DKGRegStatusVerified,
		}))
	}

	// Only the old committee deals (kyber 0,1,2); the joiner submits nothing.
	require.NoError(t, k.markDealersDealt(ctx, newDKG, []types.Deal{{Index: 0}, {Index: 1}, {Index: 2}}))
	require.NoError(t, k.invalidateMissingDealers(ctx, newDKG))

	// All four members — including the joiner that dealt nothing — stay verified.
	for _, addr := range newRegs {
		reg, err := k.getDKGRegistration(ctx, newRound, addr)
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusVerified, reg.Status, "validator %s must stay verified", addr.Hex())
	}
}

// TestInvalidateMissingDealers_AbsentRejoinerNotExpected: a bonded validator absent from the
// previous DKG round holds no share; on rejoin it cannot deal and must not be expected to
// deal (expected dealers come from the prev round's registrations, not the staking set).
func TestInvalidateMissingDealers_AbsentRejoinerNotExpected(t *testing.T) {
	const (
		oldRound = uint32(60)
		newRound = uint32(62)
	)

	// Old committee of 3 share holders. valX was bonded last round (so it is in the staking
	// ActiveValSet) but never participated in DKG, so it has no old registration/share.
	val1 := common.BytesToAddress([]byte{0x01})
	val2 := common.BytesToAddress([]byte{0x02})
	val3 := common.BytesToAddress([]byte{0x03})
	valX := common.BytesToAddress([]byte{0x09}) // bonded but absent last round
	oldByIndex := []common.Address{val1, val2, val3}

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx).WithChainID(netconf.TestChainID)

	// Previous active round: ActiveValSet includes valX, but only val1/val2/val3 registered.
	old := &types.DKGNetwork{
		Round:        oldRound,
		Total:        3,
		ActiveValSet: []string{val1.Hex(), val2.Hex(), val3.Hex(), valX.Hex()},
		Stage:        types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, old))
	require.NoError(t, k.setLatestActiveRound(ctx, old))
	for i, addr := range oldByIndex {
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:         oldRound,
			ValidatorAddr: addr.Hex(),
			Index:         uint32(i + 1),
			Status:        types.DKGRegStatusFinalized,
		}))
	}

	// New resharing round: the three share holders re-register and valX rejoins.
	newRegs := []common.Address{val1, val2, val3, valX}
	newDKG := &types.DKGNetwork{
		Round: newRound, Total: 4, Threshold: 2,
		Stage: types.DKGStageFinalization, IsResharing: true, StartBlockHeight: 400,
	}
	require.NoError(t, k.setDKGNetwork(ctx, newDKG))
	for i, addr := range newRegs {
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:         newRound,
			ValidatorAddr: addr.Hex(),
			Index:         uint32(i + 1),
			Status:        types.DKGRegStatusVerified,
		}))
	}

	// Only the old share holders can deal (kyber 0,1,2); valX holds no share and deals nothing.
	require.NoError(t, k.markDealersDealt(ctx, newDKG, []types.Deal{{Index: 0}, {Index: 1}, {Index: 2}}))
	require.NoError(t, k.invalidateMissingDealers(ctx, newDKG))

	// valX was not a previous-round share holder, so it is not an expected dealer and must
	// not be invalidated; the continuing members stay verified too.
	for _, addr := range newRegs {
		reg, err := k.getDKGRegistration(ctx, newRound, addr)
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusVerified, reg.Status, "validator %s must stay verified", addr.Hex())
	}
}

// TestInvalidateMissingDealers_NonFinalizedPrevNotExpected: a validator that was in the previous
// committee but only VERIFIED (never FINALIZED) holds no finalized share and cannot send a valid
// reshare deal, so it must not be an expected dealer. Only FINALIZED members are.
func TestInvalidateMissingDealers_NonFinalizedPrevNotExpected(t *testing.T) {
	const (
		oldRound = uint32(64)
		newRound = uint32(66)
	)

	val1 := common.BytesToAddress([]byte{0x01}) // prev FINALIZED (share holder)
	val2 := common.BytesToAddress([]byte{0x02}) // prev FINALIZED (share holder)
	valV := common.BytesToAddress([]byte{0x03}) // prev VERIFIED, never FINALIZED -> no share

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx).WithChainID(netconf.TestChainID)

	old := &types.DKGNetwork{
		Round:        oldRound,
		Total:        3,
		ActiveValSet: []string{val1.Hex(), val2.Hex(), valV.Hex()},
		Stage:        types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, old))
	require.NoError(t, k.setLatestActiveRound(ctx, old))
	oldRegs := []struct {
		addr   common.Address
		status types.DKGRegStatus
	}{
		{val1, types.DKGRegStatusFinalized},
		{val2, types.DKGRegStatusFinalized},
		{valV, types.DKGRegStatusVerified},
	}
	for i, r := range oldRegs {
		require.NoError(t, k.setDKGRegistration(ctx, r.addr, &types.DKGRegistration{
			Round:         oldRound,
			ValidatorAddr: r.addr.Hex(),
			Index:         uint32(i + 1),
			Status:        r.status,
		}))
	}

	newRegs := []common.Address{val1, val2, valV}
	newDKG := &types.DKGNetwork{
		Round: newRound, Total: 3, Threshold: 2,
		Stage: types.DKGStageFinalization, IsResharing: true, StartBlockHeight: 400,
	}
	require.NoError(t, k.setDKGNetwork(ctx, newDKG))
	for i, addr := range newRegs {
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:         newRound,
			ValidatorAddr: addr.Hex(),
			Index:         uint32(i + 1),
			Status:        types.DKGRegStatusVerified,
		}))
	}

	// The two FINALIZED share holders deal (kyber 0,1); valV holds no share and deals nothing.
	require.NoError(t, k.markDealersDealt(ctx, newDKG, []types.Deal{{Index: 0}, {Index: 1}}))
	require.NoError(t, k.invalidateMissingDealers(ctx, newDKG))

	// valV was VERIFIED (not FINALIZED) last round, so it is not an expected dealer and must not
	// be invalidated for not dealing.
	for _, addr := range newRegs {
		reg, err := k.getDKGRegistration(ctx, newRound, addr)
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusVerified, reg.Status, "validator %s must stay verified", addr.Hex())
	}
}

// TestInvalidateMissingDealers_ReshardingDealerGhostGap: an INVALIDATED member in the dealer
// round keeps its committee slot (the committee is not re-compacted), so deal.Index still equals
// reg.Index-1 for every survivor. The high-index dealer (val4) stays at its slot and must be
// credited, not invalidated.
func TestInvalidateMissingDealers_ReshardingDealerGhostGap(t *testing.T) {
	const (
		oldRound = uint32(70)
		newRound = uint32(72)
	)

	// Old round registrations, by registration index: val1=1(FIN), val2=2(FIN),
	// valBad=3(INVALIDATED), val4=4(FIN). The dealer committee keeps valBad's slot and sorts by
	// index -> [val1, val2, valBad, val4], so val4 stays at kyber index 3 (= reg.Index 4 - 1).
	val1 := common.BytesToAddress([]byte{0x01})
	val2 := common.BytesToAddress([]byte{0x02})
	valBad := common.BytesToAddress([]byte{0x08}) // invalidated in the old round; slot preserved
	val4 := common.BytesToAddress([]byte{0x04})   // high-index dealer; must keep kyber index 3

	oldRegs := []struct {
		addr   common.Address
		status types.DKGRegStatus
	}{
		{val1, types.DKGRegStatusFinalized},
		{val2, types.DKGRegStatusFinalized},
		{valBad, types.DKGRegStatusInvalidated},
		{val4, types.DKGRegStatusFinalized},
	}

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx).WithChainID(netconf.TestChainID)

	// Previous active round. ActiveValSet is the set of share holders expected to deal.
	old := &types.DKGNetwork{
		Round:        oldRound,
		Total:        3,
		ActiveValSet: []string{val1.Hex(), val2.Hex(), val4.Hex()},
		Stage:        types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, old))
	require.NoError(t, k.setLatestActiveRound(ctx, old))
	for i, r := range oldRegs {
		require.NoError(t, k.setDKGRegistration(ctx, r.addr, &types.DKGRegistration{
			Round:         oldRound,
			ValidatorAddr: r.addr.Hex(),
			Index:         uint32(i + 1),
			Status:        r.status,
		}))
	}

	// New resharing round: the three healthy share holders continue.
	newRegs := []common.Address{val1, val2, val4}
	newDKG := &types.DKGNetwork{
		Round: newRound, Total: 3, Threshold: 2,
		Stage: types.DKGStageFinalization, IsResharing: true, StartBlockHeight: 400,
	}
	require.NoError(t, k.setDKGNetwork(ctx, newDKG))
	for i, addr := range newRegs {
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:         newRound,
			ValidatorAddr: addr.Hex(),
			Index:         uint32(i + 1),
			Status:        types.DKGRegStatusVerified,
		}))
	}

	// The healthy dealers deal at their committee positions: val1@0, val2@1, val4@3
	// (valBad occupies slot 2 and does not deal).
	require.NoError(t, k.markDealersDealt(ctx, newDKG, []types.Deal{{Index: 0}, {Index: 1}, {Index: 3}}))
	require.NoError(t, k.invalidateMissingDealers(ctx, newDKG))

	// val4 deals at kyber index 3 (its slot is preserved past the invalidated valBad). It must
	// stay verified, as must val1/val2.
	for _, addr := range newRegs {
		reg, err := k.getDKGRegistration(ctx, newRound, addr)
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusVerified, reg.Status, "validator %s must stay verified", addr.Hex())
	}
}

// TestInvalidateMissingDealers_ReshardingGhostGapInvalidatesNonDealer: positive direction of
// the ghost-gap topology — a prev-round holder that continues but submits no deal is still
// correctly invalidated (the committee mapping must not over-correct).
func TestInvalidateMissingDealers_ReshardingGhostGapInvalidatesNonDealer(t *testing.T) {
	const (
		oldRound = uint32(74)
		newRound = uint32(76)
	)

	// Same ghost-gap old round as above: val1=1(FIN), val2=2(FIN), valBad=3(INVALIDATED),
	// val4=4(FIN). Committee keeps all slots -> [val1, val2, valBad, val4] (kyber 0,1,2,3).
	val1 := common.BytesToAddress([]byte{0x01})
	val2 := common.BytesToAddress([]byte{0x02})
	valBad := common.BytesToAddress([]byte{0x08})
	val4 := common.BytesToAddress([]byte{0x04})

	oldRegs := []struct {
		addr   common.Address
		status types.DKGRegStatus
	}{
		{val1, types.DKGRegStatusFinalized},
		{val2, types.DKGRegStatusFinalized},
		{valBad, types.DKGRegStatusInvalidated},
		{val4, types.DKGRegStatusFinalized},
	}

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx).WithChainID(netconf.TestChainID)

	old := &types.DKGNetwork{
		Round:        oldRound,
		Total:        3,
		ActiveValSet: []string{val1.Hex(), val2.Hex(), val4.Hex()},
		Stage:        types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, old))
	require.NoError(t, k.setLatestActiveRound(ctx, old))
	for i, r := range oldRegs {
		require.NoError(t, k.setDKGRegistration(ctx, r.addr, &types.DKGRegistration{
			Round:         oldRound,
			ValidatorAddr: r.addr.Hex(),
			Index:         uint32(i + 1),
			Status:        r.status,
		}))
	}

	newRegs := []common.Address{val1, val2, val4}
	newDKG := &types.DKGNetwork{
		Round: newRound, Total: 3, Threshold: 2,
		Stage: types.DKGStageFinalization, IsResharing: true, StartBlockHeight: 400,
	}
	require.NoError(t, k.setDKGNetwork(ctx, newDKG))
	for i, addr := range newRegs {
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:         newRound,
			ValidatorAddr: addr.Hex(),
			Index:         uint32(i + 1),
			Status:        types.DKGRegStatusVerified,
		}))
	}

	// val1 (kyber 0) and val2 (kyber 1) deal; val4 (kyber 3) submits nothing.
	require.NoError(t, k.markDealersDealt(ctx, newDKG, []types.Deal{{Index: 0}, {Index: 1}}))
	require.NoError(t, k.invalidateMissingDealers(ctx, newDKG))

	// val4 was a real prev-round holder that skipped dealing -> invalidated; val1/val2 kept.
	reg4, err := k.getDKGRegistration(ctx, newRound, val4)
	require.NoError(t, err)
	require.Equal(t, types.DKGRegStatusInvalidated, reg4.Status, "non-dealing holder val4 must be invalidated")
	for _, addr := range []common.Address{val1, val2} {
		reg, err := k.getDKGRegistration(ctx, newRound, addr)
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusVerified, reg.Status, "validator %s must stay verified", addr.Hex())
	}
}

// TestMarkDealersDealt_InitialDKGMidGap: a non-resharing round does NOT drop an invalidated
// member from the committee — the kernel fixed it from all registrations at dealing start, so
// positions equal reg.Index-1 and stay stable even after a mid-round invalidation. D keeps its
// position (reg.Index 4 -> kyber index 3), not a compacted index 2.
func TestMarkDealersDealt_InitialDKGMidGap(t *testing.T) {
	const round = uint32(17)

	// Current round: A=1(VERIFIED), B=2(VERIFIED), BAD=3(INVALIDATED), D=4(VERIFIED).
	// Committee (not filtered) -> [A, B, BAD, D] (kyber 0,1,2,3); D stays at kyber index 3.
	valA := common.BytesToAddress([]byte{0x01})
	valB := common.BytesToAddress([]byte{0x02})
	valBad := common.BytesToAddress([]byte{0x08})
	valD := common.BytesToAddress([]byte{0x04})

	regs := []struct {
		addr   common.Address
		status types.DKGRegStatus
	}{
		{valA, types.DKGRegStatusVerified},
		{valB, types.DKGRegStatusVerified},
		{valBad, types.DKGRegStatusInvalidated},
		{valD, types.DKGRegStatusVerified},
	}

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx).WithChainID(netconf.TestChainID)

	// Non-resharing round (dealer committee == current round).
	latestRound := &types.DKGNetwork{
		Round: round, Total: 4, Threshold: 2,
		Stage: types.DKGStageFinalization, StartBlockHeight: 400,
	}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))
	for i, r := range regs {
		require.NoError(t, k.setDKGRegistration(ctx, r.addr, &types.DKGRegistration{
			Round:         round,
			ValidatorAddr: r.addr.Hex(),
			Index:         uint32(i + 1),
			Status:        r.status,
		}))
	}

	// A (kyber 0) and D (kyber 3) deal; B (kyber 1) and BAD (kyber 2) do not.
	require.NoError(t, k.markDealersDealt(ctx, latestRound, []types.Deal{{Index: 0}, {Index: 3}}))
	require.NoError(t, k.invalidateMissingDealers(ctx, latestRound))

	want := map[common.Address]types.DKGRegStatus{
		valA:   types.DKGRegStatusVerified,    // dealt
		valB:   types.DKGRegStatusInvalidated, // verified member, no deal
		valBad: types.DKGRegStatusInvalidated, // already invalidated, untouched
		valD:   types.DKGRegStatusVerified,    // dealt at kyber index 3 (reg.Index 4 - 1)
	}
	for addr, status := range want {
		reg, err := k.getDKGRegistration(ctx, round, addr)
		require.NoError(t, err)
		require.Equal(t, status, reg.Status, "validator %s", addr.Hex())
	}
}
