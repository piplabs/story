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

	"go.dedis.ch/kyber/v4/group/edwards25519"
	"go.dedis.ch/kyber/v4/share"
	"go.uber.org/mock/gomock"
)

// generateConsensusPolyShares returns the serialized public coefficients and the
// public key share for each of n participants (pubKeyShares[i] = F(i+1), 0-indexed slice).
func generateConsensusPolyShares(t *testing.T, n, threshold int) (publicCoeffs [][]byte, pubKeyShares [][]byte) {
	t.Helper()

	suite := edwards25519.NewBlakeSHA256Ed25519()
	secret := suite.Scalar().Pick(suite.RandomStream())
	priPoly := share.NewPriPoly(suite, threshold, secret, suite.RandomStream())
	pubPoly := priPoly.Commit(suite.Point().Base())

	_, commits := pubPoly.Info()
	publicCoeffs = make([][]byte, len(commits))
	for i, c := range commits {
		bz, err := c.MarshalBinary()
		require.NoError(t, err)

		publicCoeffs[i] = bz
	}

	pubKeyShares = make([][]byte, n)
	for i := range n {
		// PubPoly.Eval(i) evaluates at x = i+1, matching the participant's 1-based index.
		bz, err := pubPoly.Eval(i).V.MarshalBinary()
		require.NoError(t, err)

		pubKeyShares[i] = bz
	}

	return publicCoeffs, pubKeyShares
}

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

// TestInvalidateDivergentShares verifies the core sweep logic: finalized validators
// whose public key share lies on the consensus polynomial are kept, while those whose
// share diverges (or is undecodable) are invalidated.
func TestInvalidateDivergentShares(t *testing.T) {
	const (
		round     = uint32(7)
		n         = 5
		threshold = 3
	)

	// Consensus polynomial and the valid public key share for each participant.
	publicCoeffs, validShares := generateConsensusPolyShares(t, n, threshold)
	// An independent polynomial: a diverged validator's share lies on this instead.
	_, divergentShares := generateConsensusPolyShares(t, n, threshold)

	validators := make([]common.Address, n)
	for i := range n {
		validators[i] = common.BytesToAddress([]byte{byte(i + 1)})
	}

	// Per-participant share to store; indices 2 (divergent) and 3 (malformed) are bad.
	const divergentIdx, malformedIdx = 2, 3

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx)

	latestRound := &types.DKGNetwork{
		Round:           round,
		Total:           n,
		Threshold:       threshold,
		Stage:           types.DKGStageFinalization,
		GlobalPublicKey: publicCoeffs[0], // GPK = F(0) = C_0
		PublicCoeffs:    publicCoeffs,
	}
	require.NoError(t, k.setDKGNetwork(ctx, latestRound))

	for i := range n {
		pubKeyShare := validShares[i]
		switch i {
		case divergentIdx:
			pubKeyShare = divergentShares[i]
		case malformedIdx:
			pubKeyShare = []byte{0x01, 0x02, 0x03} // undecodable point
		}

		reg := &types.DKGRegistration{
			Round:         round,
			ValidatorAddr: validators[i].Hex(),
			Index:         uint32(i + 1), // 1-based
			PubKeyShare:   pubKeyShare,
			Status:        types.DKGRegStatusFinalized,
		}
		require.NoError(t, k.setDKGRegistration(ctx, validators[i], reg))
	}

	require.NoError(t, k.invalidateDivergentShares(ctx, latestRound))

	for i := range n {
		reg, err := k.getDKGRegistration(ctx, round, validators[i])
		require.NoError(t, err)

		if i == divergentIdx || i == malformedIdx {
			require.Equal(t, types.DKGRegStatusInvalidated, reg.Status,
				"participant %d (off-consensus-polynomial) should be invalidated", i+1)
		} else {
			require.Equal(t, types.DKGRegStatusFinalized, reg.Status,
				"participant %d (on consensus polynomial) should remain finalized", i+1)
		}
	}
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

// TestFinalizeDKGRound_V190Gate verifies that the divergent-share sweep only runs once
// the v1.9.0 upgrade height is reached, so pre-upgrade blocks replay deterministically.
func TestFinalizeDKGRound_V190Gate(t *testing.T) {
	const (
		round     = uint32(9)
		n         = 4
		threshold = 5 // intentionally > n so FinalizeDKGRound always takes the skip path
	)

	publicCoeffs, validShares := generateConsensusPolyShares(t, n, threshold)
	_, divergentShares := generateConsensusPolyShares(t, n, threshold)

	const divergentIdx = 1

	// setup builds a keeper with one divergent finalized reg in a round that started at
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
			GlobalPublicKey:  publicCoeffs[0],
			PublicCoeffs:     publicCoeffs,
		}
		require.NoError(t, k.setDKGNetwork(ctx, latestRound))

		for i := range n {
			pubKeyShare := validShares[i]
			if i == divergentIdx {
				pubKeyShare = divergentShares[i]
			}
			reg := &types.DKGRegistration{
				Round:         round,
				ValidatorAddr: validators[i].Hex(),
				Index:         uint32(i + 1),
				PubKeyShare:   pubKeyShare,
				Status:        types.DKGRegStatusFinalized,
			}
			require.NoError(t, k.setDKGRegistration(ctx, validators[i], reg))
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
			"at v1.9.0 height the divergent share must be invalidated")

		// A validator on the consensus polynomial is untouched.
		good, err := k.getDKGRegistration(ctx, round, validators[0])
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusFinalized, good.Status)
	})
}

// TestMarkDealersDealt verifies that only in-range dealer indices (1..Total) are recorded.
func TestMarkDealersDealt(t *testing.T) {
	const round = uint32(12)

	k, _, _, baseCtx := setupDKGKeeperWithMocks(t)
	ctx := sdk.UnwrapSDKContext(baseCtx)

	latestRound := &types.DKGNetwork{Round: round, Total: 3}
	deals := []types.Deal{
		{Index: 1}, {Index: 1}, // duplicate dealer 1
		{Index: 3},             // valid
		{Index: 0}, {Index: 4}, // out of range (Total = 3)
	}
	require.NoError(t, k.markDealersDealt(ctx, latestRound, deals))

	for idx, want := range map[uint32]bool{1: true, 2: false, 3: true, 0: false, 4: false} {
		has, err := k.DealtDealers.Has(ctx, dealtDealerKey(round, idx))
		require.NoError(t, err)
		require.Equal(t, want, has, "index %d", idx)
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

	// Dealers 1 and 3 submitted deals; dealers 2 and 4 did not.
	require.NoError(t, k.markDealersDealt(ctx, latestRound, []types.Deal{{Index: 1}, {Index: 3}}))

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

	// Marks for dealt dealers are pruned.
	for _, idx := range []uint32{1, 3} {
		has, err := k.DealtDealers.Has(ctx, dealtDealerKey(round, idx))
		require.NoError(t, err)
		require.False(t, has, "dealt mark for index %d should be pruned", idx)
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

		// Every dealer except nonDealerIdx submitted a deal.
		var deals []types.Deal
		for i := range n {
			if i != nonDealerIdx {
				deals = append(deals, types.Deal{Index: uint32(i + 1)})
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
