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
