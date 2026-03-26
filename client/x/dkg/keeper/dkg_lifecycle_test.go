package keeper

import (
	"context"
	"math/big"
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/netconf"

	"go.uber.org/mock/gomock"
)

// testDKGParams returns DKG params with short stage durations for lifecycle testing.
func testDKGParams() types.Params {
	return types.NewParams(
		5,  // RegistrationPeriod: 5 blocks
		10, // DealingPeriod: 10 blocks
		10, // FinalizationPeriod: 10 blocks
		20, // ActivePeriod: 20 blocks
		types.DefaultDkgCommitteeRewardPortion,
		3, // MinReqRegisteredParticipants
		3, // MinReqFinalizedParticipants
		types.DefaultOperationalThreshold,
	)
}

// dkgLifecycleEnv holds the test environment for DKG lifecycle tests.
type dkgLifecycleEnv struct {
	keeper     *Keeper
	sdkCtx     sdk.Context
	sk         *dkgtestutil.MockStakingKeeper
	dk         *dkgtestutil.MockDistributionKeeper
	validators []common.Address
}

// setupDKGLifecycleEnv creates a test environment with DKG keeper, short stage
// periods, and mock validators configured. The SDK context uses the DKGTestChainID
// with block height at the V160 activation point so BeginBlocker is active.
func setupDKGLifecycleEnv(t *testing.T, numValidators int) *dkgLifecycleEnv {
	t.Helper()

	k, _, dk, ctx := setupDKGKeeperWithMocks(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Use TestChainID where V160 activates at block 110.
	// Set a non-nil HeaderHash so DKG network's StartBlockHash is populated.
	testHeaderHash := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	sdkCtx = sdkCtx.WithChainID(netconf.TestChainID).WithBlockHeight(110).WithHeaderHash(testHeaderHash.Bytes())

	// Set short stage durations for testing
	require.NoError(t, k.SetParams(sdkCtx, testDKGParams()))

	// Generate deterministic validator addresses
	validators := make([]common.Address, numValidators)
	for i := range numValidators {
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		validators[i] = crypto.PubkeyToAddress(key.PublicKey)
	}

	// Find the gomock controller's StakingKeeper for expectations
	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)

	return &dkgLifecycleEnv{
		keeper:     k,
		sdkCtx:     sdkCtx,
		sk:         sk,
		dk:         dk,
		validators: validators,
	}
}

// mockValidators returns a slice of staking validators whose OperatorAddresses
// are derived from the test validators' EVM addresses. Each validator is bonded
// and not jailed.
func (env *dkgLifecycleEnv) mockValidators() []stakingtypes.Validator {
	vals := make([]stakingtypes.Validator, len(env.validators))
	for i, addr := range env.validators {
		vals[i] = stakingtypes.Validator{
			OperatorAddress: sdk.ValAddress(addr.Bytes()).String(),
			Jailed:          false,
			Status:          stakingtypes.Bonded,
		}
	}

	return vals
}

// expectZeroUbiBalance sets up mock expectations for settleRewardsForPreviousCommittee
// to return zero UBI balance, causing it to short-circuit without further calls.
func (env *dkgLifecycleEnv) expectZeroUbiBalance() {
	env.dk.EXPECT().GetUbiBalanceByDenom(gomock.Any(), sdk.DefaultBondDenom).Return(math.ZeroInt(), nil)
}

// advanceToHeight sets the block height to the specified value.
func (env *dkgLifecycleEnv) advanceToHeight(h int64) {
	env.sdkCtx = env.sdkCtx.WithBlockHeight(h)
}

// registerValidators creates verified DKG registrations for all validators in the env.
func (env *dkgLifecycleEnv) registerValidators(t *testing.T, round uint32) {
	t.Helper()

	for i, val := range env.validators {
		reg := &types.DKGRegistration{
			Round:          round,
			ValidatorAddr:  val.Hex(),
			Index:          uint32(i + 1),
			DkgPubKey:      []byte("dkg-pubkey-" + val.Hex()),
			CommPubKey:     []byte("comm-pubkey-" + val.Hex()),
			EnclaveReport:  []byte("enclave-report"),
			Status:         types.DKGRegStatusVerified,
			CodeCommitment: make([]byte, 32),
			EnclaveType:    make([]byte, 32),
		}
		require.NoError(t, env.keeper.setDKGRegistration(env.sdkCtx, val, reg))
	}
}

// finalizeValidators creates finalized DKG registrations for all validators in the env.
func (env *dkgLifecycleEnv) finalizeValidators(t *testing.T, round uint32, globalPubKey []byte) {
	t.Helper()

	for i, val := range env.validators {
		// Mark registration as finalized
		reg := &types.DKGRegistration{
			Round:          round,
			ValidatorAddr:  val.Hex(),
			Index:          uint32(i + 1),
			DkgPubKey:      []byte("dkg-pubkey-" + val.Hex()),
			CommPubKey:     []byte("comm-pubkey-" + val.Hex()),
			EnclaveReport:  []byte("enclave-report"),
			Status:         types.DKGRegStatusFinalized,
			CodeCommitment: make([]byte, 32),
			EnclaveType:    make([]byte, 32),
			PubKeyShare:    []byte("pub-key-share-" + val.Hex()),
		}
		require.NoError(t, env.keeper.setDKGRegistration(env.sdkCtx, val, reg))
	}

	// Add global public key votes to meet threshold
	net, err := env.keeper.getLatestDKGNetwork(env.sdkCtx)
	require.NoError(t, err)
	net.GlobalPublicKey = globalPubKey
	net.PublicCoeffs = [][]byte{globalPubKey} // simplified
	require.NoError(t, env.keeper.setDKGNetwork(env.sdkCtx, net))
}

// TestDKGLifecycle_InitialRound verifies the full lifecycle of the first DKG round:
// BeginBlocker initiates round → Registration → Dealing → Finalization → Active.
func TestDKGLifecycle_InitialRound(t *testing.T) {
	t.Parallel()
	env := setupDKGLifecycleEnv(t, 3)
	params := testDKGParams()

	// Round 1 should not exist yet
	latestRound, err := env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.Nil(t, latestRound, "no DKG round should exist before BeginBlocker")

	// --- Phase 1: First BeginBlocker initiates Round 1 ---
	// Mock GetAllValidators for InitiateDKGRound
	env.sk.EXPECT().GetAllValidators(gomock.Any()).Return(env.mockValidators(), nil).Times(1)

	startHeight := env.sdkCtx.BlockHeight()
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Verify round 1 was created in Registration stage
	round1, err := env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.NotNil(t, round1)
	require.Equal(t, uint32(1), round1.Round)
	require.Equal(t, types.DKGStageRegistration, round1.Stage)
	require.Equal(t, startHeight, round1.StartBlockHeight)
	require.False(t, round1.IsResharing, "initial round should not be resharing")
	require.False(t, round1.IsUpgrade, "initial round should not be upgrade")

	// --- Phase 2: Registration period — simulate validators registering ---
	env.registerValidators(t, 1)

	// Verify registrations are stored
	for _, val := range env.validators {
		reg, err := env.keeper.getDKGRegistration(env.sdkCtx, 1, val)
		require.NoError(t, err)
		require.Equal(t, types.DKGRegStatusVerified, reg.Status)
	}

	// --- Phase 3: Advance to dealing transition ---
	dealingHeight := startHeight + int64(params.RegistrationPeriod)
	env.advanceToHeight(dealingHeight)

	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	round1, err = env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageDealing, round1.Stage, "should transition to Dealing")
	require.Equal(t, uint32(3), round1.Total, "total should match registered validators")
	require.True(t, round1.Threshold > 0, "threshold should be computed")

	// --- Phase 4: Advance to finalization transition ---
	finalizationHeight := startHeight + int64(params.RegistrationPeriod) + int64(params.DealingPeriod)
	env.advanceToHeight(finalizationHeight)

	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	round1, err = env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageFinalization, round1.Stage, "should transition to Finalization")

	// --- Phase 5: Simulate validators finalizing ---
	testGlobalPubKey := []byte("test-global-public-key-32bytes!!")
	env.finalizeValidators(t, 1, testGlobalPubKey)

	// --- Phase 6: Advance to active transition ---
	activeHeight := startHeight + int64(params.RegistrationPeriod) + int64(params.DealingPeriod) + int64(params.FinalizationPeriod)
	env.advanceToHeight(activeHeight)

	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	round1, err = env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageActive, round1.Stage, "should transition to Active")

	// Verify the round is now the latest active round
	activeRound, err := env.keeper.GetLatestActiveRound(env.sdkCtx)
	require.NoError(t, err)
	require.NotNil(t, activeRound)
	require.Equal(t, uint32(1), activeRound.Round)
	require.Equal(t, testGlobalPubKey, activeRound.GlobalPublicKey)
}

// TestDKGLifecycle_ProactiveResharing verifies that after an initial round completes
// and the Active period ends, a new round is automatically initiated as a resharing
// round. The global public key should be preserved across rounds.
func TestDKGLifecycle_ProactiveResharing(t *testing.T) {
	t.Parallel()
	env := setupDKGLifecycleEnv(t, 3)
	params := testDKGParams()

	// --- Setup: Complete Round 1 ---
	startHeight := env.sdkCtx.BlockHeight()
	env.sk.EXPECT().GetAllValidators(gomock.Any()).Return(env.mockValidators(), nil).Times(1)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Register validators
	env.registerValidators(t, 1)

	// Transition to Dealing
	dealingHeight := startHeight + int64(params.RegistrationPeriod)
	env.advanceToHeight(dealingHeight)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Transition to Finalization
	finalizationHeight := startHeight + int64(params.RegistrationPeriod) + int64(params.DealingPeriod)
	env.advanceToHeight(finalizationHeight)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Finalize validators
	testGlobalPubKey := []byte("global-public-key-round1-32byte")
	env.finalizeValidators(t, 1, testGlobalPubKey)

	// Transition to Active
	activeHeight := startHeight + int64(params.RegistrationPeriod) + int64(params.DealingPeriod) + int64(params.FinalizationPeriod)
	env.advanceToHeight(activeHeight)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Verify Round 1 is active
	activeRound, err := env.keeper.GetLatestActiveRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, uint32(1), activeRound.Round)
	require.Equal(t, testGlobalPubKey, activeRound.GlobalPublicKey)

	// --- Phase 1: Active period ends — auto-initiation of Round 2 ---
	resharingHeight := startHeight + int64(params.RegistrationPeriod) + int64(params.DealingPeriod) + int64(params.FinalizationPeriod) + int64(params.ActivePeriod)
	env.advanceToHeight(resharingHeight)

	// InitiateDKGRound will call GetAllValidators again
	env.sk.EXPECT().GetAllValidators(gomock.Any()).Return(env.mockValidators(), nil).Times(1)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Verify Round 2 was created
	round2, err := env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.NotNil(t, round2)
	require.Equal(t, uint32(2), round2.Round)
	require.Equal(t, types.DKGStageRegistration, round2.Stage)
	require.True(t, round2.IsResharing, "Round 2 should be a resharing round")
	require.False(t, round2.IsUpgrade, "proactive resharing should not be upgrade")

	// --- Phase 2: Complete Round 2 lifecycle ---
	round2Start := round2.StartBlockHeight

	// Register validators for round 2
	env.registerValidators(t, 2)

	// Transition to Dealing
	env.advanceToHeight(round2Start + int64(params.RegistrationPeriod))
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	round2, err = env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageDealing, round2.Stage)

	// Transition to Finalization
	env.advanceToHeight(round2Start + int64(params.RegistrationPeriod) + int64(params.DealingPeriod))
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	round2, err = env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageFinalization, round2.Stage)

	// Finalize validators with the SAME global public key (proactive resharing preserves it)
	env.finalizeValidators(t, 2, testGlobalPubKey)

	// Transition to Active — FinalizeDKGRound calls settleRewardsForPreviousCommittee
	// which queries distribution keeper for UBI balance.
	env.expectZeroUbiBalance()
	env.advanceToHeight(round2Start + int64(params.RegistrationPeriod) + int64(params.DealingPeriod) + int64(params.FinalizationPeriod))
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Verify Round 2 is now the active round with the same global public key
	activeRound, err = env.keeper.GetLatestActiveRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, uint32(2), activeRound.Round)
	require.Equal(t, testGlobalPubKey, activeRound.GlobalPublicKey,
		"global public key should be preserved across proactive resharing rounds")

	// Verify Round 1's stage: when Active stage ends, BeginBlocker sets
	// latestRound.Stage = Registration (the transition target) before calling
	// InitiateDKGRound. So Round 1's stored stage is Registration, not Active.
	// endPreviousActiveRound only marks stages == Active as Ended, so Round 1
	// remains at Registration (its stage was already overwritten by the transition).
	round1, err := env.keeper.getDKGNetworkByRound(env.sdkCtx, 1)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageRegistration, round1.Stage,
		"Round 1 stage was overwritten to Registration during Active→Registration transition")
}

// TestDKGLifecycle_UpgradeResharing verifies that when a kernel upgrade is
// scheduled and its activation height is reached, BeginBlocker initiates an
// upgrade resharing round (IsResharing=true, IsUpgrade=true).
func TestDKGLifecycle_UpgradeResharing(t *testing.T) {
	t.Parallel()
	env := setupDKGLifecycleEnv(t, 3)
	params := testDKGParams()

	// --- Setup: Complete Round 1 (initial DKG) ---
	startHeight := env.sdkCtx.BlockHeight()
	env.sk.EXPECT().GetAllValidators(gomock.Any()).Return(env.mockValidators(), nil).Times(1)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	env.registerValidators(t, 1)

	dealingHeight := startHeight + int64(params.RegistrationPeriod)
	env.advanceToHeight(dealingHeight)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	finalizationHeight := startHeight + int64(params.RegistrationPeriod) + int64(params.DealingPeriod)
	env.advanceToHeight(finalizationHeight)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	testGlobalPubKey := []byte("global-public-key-round1-32byte")
	env.finalizeValidators(t, 1, testGlobalPubKey)

	activeHeight := startHeight + int64(params.RegistrationPeriod) + int64(params.DealingPeriod) + int64(params.FinalizationPeriod)
	env.advanceToHeight(activeHeight)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// --- Schedule a kernel upgrade ---
	upgradeActivationHeight := activeHeight + 5 // 5 blocks into Active phase
	require.NoError(t, env.keeper.UpgradeScheduled(env.sdkCtx, upgradeActivationHeight, "v3.0.0"))

	// Verify upgrade info is stored
	upgradeInfo, err := env.keeper.GetPendingUpgrade(env.sdkCtx)
	require.NoError(t, err)
	require.NotNil(t, upgradeInfo)
	require.Equal(t, "v3.0.0", upgradeInfo.UpgradeVersion)
	require.False(t, upgradeInfo.IsActivated)

	// --- Advance to upgrade activation height ---
	env.advanceToHeight(upgradeActivationHeight)
	env.sk.EXPECT().GetAllValidators(gomock.Any()).Return(env.mockValidators(), nil).Times(1)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Verify upgrade resharing round was initiated
	round2, err := env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, uint32(2), round2.Round)
	require.Equal(t, types.DKGStageRegistration, round2.Stage)
	require.True(t, round2.IsResharing, "upgrade resharing round must be resharing")
	require.True(t, round2.IsUpgrade, "upgrade resharing round must be marked as upgrade")

	// Verify upgrade info is marked as activated (not deleted yet).
	// GetPendingUpgrade skips activated entries, so use GetKernelUpgradeInfo directly.
	upgradeInfo, err = env.keeper.GetKernelUpgradeInfo(env.sdkCtx, "v3.0.0")
	require.NoError(t, err)
	require.NotNil(t, upgradeInfo)
	require.True(t, upgradeInfo.IsActivated, "upgrade should be marked as activated")

	// --- Complete the upgrade resharing round ---
	round2Start := round2.StartBlockHeight

	env.registerValidators(t, 2)

	env.advanceToHeight(round2Start + int64(params.RegistrationPeriod))
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	env.advanceToHeight(round2Start + int64(params.RegistrationPeriod) + int64(params.DealingPeriod))
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Finalize with same global pub key
	env.finalizeValidators(t, 2, testGlobalPubKey)

	// Transition to Active — FinalizeDKGRound calls settleRewardsForPreviousCommittee
	env.expectZeroUbiBalance()
	env.advanceToHeight(round2Start + int64(params.RegistrationPeriod) + int64(params.DealingPeriod) + int64(params.FinalizationPeriod))
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Verify upgrade info is deleted after successful upgrade round.
	// GetKernelUpgradeInfo should return "not found" since deleteActivatedUpgradeInfo removed it.
	_, err = env.keeper.GetKernelUpgradeInfo(env.sdkCtx, "v3.0.0")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")

	// Verify round 2 is now active
	activeRound, err := env.keeper.GetLatestActiveRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, uint32(2), activeRound.Round)
}

// TestDKGLifecycle_StageTransitionTiming verifies that stage transitions happen
// at exact block boundaries and not before.
func TestDKGLifecycle_StageTransitionTiming(t *testing.T) {
	t.Parallel()
	env := setupDKGLifecycleEnv(t, 3)
	params := testDKGParams()

	startHeight := env.sdkCtx.BlockHeight()
	env.sk.EXPECT().GetAllValidators(gomock.Any()).Return(env.mockValidators(), nil).Times(1)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Register validators so dealing transition succeeds
	env.registerValidators(t, 1)

	// --- Verify no transition 1 block before dealing ---
	oneBeforeDealing := startHeight + int64(params.RegistrationPeriod) - 1
	env.advanceToHeight(oneBeforeDealing)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	round, err := env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageRegistration, round.Stage,
		"should still be in Registration 1 block before dealing")

	// --- Verify transition at exact dealing boundary ---
	dealingExact := startHeight + int64(params.RegistrationPeriod)
	env.advanceToHeight(dealingExact)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	round, err = env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageDealing, round.Stage,
		"should transition to Dealing at exact boundary")
}

// TestDKGLifecycle_InsufficientRegistrations verifies that when the number of
// verified registrations is below MinReqRegisteredParticipants, the round is
// skipped and a new round is initiated.
func TestDKGLifecycle_InsufficientRegistrations(t *testing.T) {
	t.Parallel()
	env := setupDKGLifecycleEnv(t, 3)
	params := testDKGParams()

	startHeight := env.sdkCtx.BlockHeight()

	// Initiate Round 1
	env.sk.EXPECT().GetAllValidators(gomock.Any()).Return(env.mockValidators(), nil).Times(1)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Only register 2 of 3 validators (below min_req_registered=3)
	for i := 0; i < 2; i++ {
		reg := &types.DKGRegistration{
			Round:          1,
			ValidatorAddr:  env.validators[i].Hex(),
			Index:          uint32(i + 1),
			DkgPubKey:      []byte("dkg-pubkey"),
			CommPubKey:     []byte("comm-pubkey"),
			EnclaveReport:  []byte("enclave-report"),
			Status:         types.DKGRegStatusVerified,
			CodeCommitment: make([]byte, 32),
			EnclaveType:    make([]byte, 32),
		}
		require.NoError(t, env.keeper.setDKGRegistration(env.sdkCtx, env.validators[i], reg))
	}

	// Transition to dealing — should skip round due to insufficient registrations
	dealingHeight := startHeight + int64(params.RegistrationPeriod)
	env.advanceToHeight(dealingHeight)

	// SkipToNextRound -> InitiateDKGRound will call GetAllValidators
	env.sk.EXPECT().GetAllValidators(gomock.Any()).Return(env.mockValidators(), nil).Times(1)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Verify Round 1 was marked as Failed and Round 2 was created
	round1, err := env.keeper.getDKGNetworkByRound(env.sdkCtx, 1)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageFailed, round1.Stage,
		"Round 1 should be failed due to insufficient registrations")

	round2, err := env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, uint32(2), round2.Round)
	require.Equal(t, types.DKGStageRegistration, round2.Stage)
}

// TestDKGLifecycle_InsufficientFinalizations verifies that when the number of
// finalized registrations is below the threshold, the round is skipped.
func TestDKGLifecycle_InsufficientFinalizations(t *testing.T) {
	t.Parallel()
	env := setupDKGLifecycleEnv(t, 3)
	params := testDKGParams()

	startHeight := env.sdkCtx.BlockHeight()

	// Initiate and register Round 1
	env.sk.EXPECT().GetAllValidators(gomock.Any()).Return(env.mockValidators(), nil).Times(1)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))
	env.registerValidators(t, 1)

	// Transition to Dealing
	dealingHeight := startHeight + int64(params.RegistrationPeriod)
	env.advanceToHeight(dealingHeight)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Transition to Finalization
	finalizationHeight := startHeight + int64(params.RegistrationPeriod) + int64(params.DealingPeriod)
	env.advanceToHeight(finalizationHeight)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Only finalize 1 of 3 validators (below min_req_finalized=3)
	reg := &types.DKGRegistration{
		Round:          1,
		ValidatorAddr:  env.validators[0].Hex(),
		Index:          1,
		DkgPubKey:      []byte("dkg-pubkey"),
		CommPubKey:     []byte("comm-pubkey"),
		EnclaveReport:  []byte("enclave-report"),
		Status:         types.DKGRegStatusFinalized,
		CodeCommitment: make([]byte, 32),
		EnclaveType:    make([]byte, 32),
		PubKeyShare:    []byte("share"),
	}
	require.NoError(t, env.keeper.setDKGRegistration(env.sdkCtx, env.validators[0], reg))

	// Transition to Active — should skip due to insufficient finalizations
	activeHeight := startHeight + int64(params.RegistrationPeriod) + int64(params.DealingPeriod) + int64(params.FinalizationPeriod)
	env.advanceToHeight(activeHeight)

	// SkipToNextRound will call GetAllValidators
	env.sk.EXPECT().GetAllValidators(gomock.Any()).Return(env.mockValidators(), nil).Times(1)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	// Verify Round 1 was failed and Round 2 started
	round1, err := env.keeper.getDKGNetworkByRound(env.sdkCtx, 1)
	require.NoError(t, err)
	require.Equal(t, types.DKGStageFailed, round1.Stage)

	round2, err := env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, uint32(2), round2.Round)
	require.Equal(t, types.DKGStageRegistration, round2.Stage)
}

// TestDKGLifecycle_RegistrationValidation verifies the Registered handler's
// validation logic: round mismatch, stage check, active validator set membership,
// and start block verification.
func TestDKGLifecycle_RegistrationValidation(t *testing.T) {
	t.Parallel()
	env := setupDKGLifecycleEnv(t, 3)

	// Initiate round 1
	env.sk.EXPECT().GetAllValidators(gomock.Any()).Return(env.mockValidators(), nil).Times(1)
	require.NoError(t, env.keeper.BeginBlocker(env.sdkCtx))

	round, err := env.keeper.GetLatestDKGRound(env.sdkCtx)
	require.NoError(t, err)

	testCodeCommitment := [32]byte{0x01}
	testEnclaveType := [32]byte{0x01}

	// The test context's header hash is empty (zero bytes), so StartBlockHash
	// stored in the DKG network is also empty. Use the stored value directly.
	var startBlockHash [32]byte
	if len(round.StartBlockHash) == 32 {
		copy(startBlockHash[:], round.StartBlockHash)
	}

	tcs := []struct {
		name             string
		round            uint32
		startBlockHeight *big.Int
		startBlockHash   [32]byte
		validator        common.Address
		expectedErr      string
	}{
		{
			name:             "fail: round mismatch",
			round:            999,
			startBlockHeight: big.NewInt(round.StartBlockHeight),
			startBlockHash:   startBlockHash,
			validator:        env.validators[0],
			expectedErr:      "round mismatch",
		},
		{
			name:             "fail: start block height mismatch",
			round:            1,
			startBlockHeight: big.NewInt(9999),
			startBlockHash:   startBlockHash,
			validator:        env.validators[0],
			expectedErr:      "start block height mismatch",
		},
		{
			name:             "fail: start block hash mismatch",
			round:            1,
			startBlockHeight: big.NewInt(round.StartBlockHeight),
			startBlockHash:   [32]byte{0xFF}, // wrong hash
			validator:        env.validators[0],
			expectedErr:      "start block hash mismatch",
		},
		{
			name:             "fail: validator not in active set",
			round:            1,
			startBlockHeight: big.NewInt(round.StartBlockHeight),
			startBlockHash:   startBlockHash,
			validator:        common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			expectedErr:      "msg sender is not in the active validator set",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := env.keeper.Registered(
				env.sdkCtx,
				tc.validator,
				testCodeCommitment,
				tc.round,
				tc.startBlockHeight,
				tc.startBlockHash,
				testEnclaveType,
				[]byte("dkg-pubkey"),
				[]byte("comm-pubkey"),
				[]byte("enclave-report"),
			)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

// TestDKGLifecycle_UpgradeScheduleAndCancel verifies the upgrade scheduling
// and cancellation flow.
func TestDKGLifecycle_UpgradeScheduleAndCancel(t *testing.T) {
	t.Parallel()
	env := setupDKGLifecycleEnv(t, 3)

	// Schedule an upgrade
	require.NoError(t, env.keeper.UpgradeScheduled(env.sdkCtx, 500, "v3.0.0"))

	info, err := env.keeper.GetPendingUpgrade(env.sdkCtx)
	require.NoError(t, err)
	require.NotNil(t, info)
	require.Equal(t, "v3.0.0", info.UpgradeVersion)
	require.Equal(t, int64(500), info.ActivationHeight)

	// Cannot schedule another while one is pending
	err = env.keeper.UpgradeScheduled(env.sdkCtx, 600, "v4.0.0")
	require.Error(t, err)
	require.Contains(t, err.Error(), "pending upgrade already exists")

	// Cancel the upgrade
	require.NoError(t, env.keeper.UpgradeCancelled(env.sdkCtx, "v3.0.0"))

	info, err = env.keeper.GetPendingUpgrade(env.sdkCtx)
	require.NoError(t, err)
	require.Nil(t, info, "upgrade should be deleted after cancellation")

	// Now can schedule a new one
	require.NoError(t, env.keeper.UpgradeScheduled(env.sdkCtx, 700, "v4.0.0"))

	info, err = env.keeper.GetPendingUpgrade(env.sdkCtx)
	require.NoError(t, err)
	require.Equal(t, "v4.0.0", info.UpgradeVersion)
}

// TestDKGLifecycle_CalculateThreshold verifies threshold calculation from
// total participants and operational threshold basis points.
func TestDKGLifecycle_CalculateThreshold(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name      string
		total     uint32
		threshold uint32
		expected  uint32
	}{
		{name: "3 of 3 at 667/1000", total: 3, threshold: 667, expected: 3},
		{name: "5 at 667/1000 = ceil(3.335) = 4", total: 5, threshold: 667, expected: 4},
		{name: "10 at 667/1000 = ceil(6.67) = 7", total: 10, threshold: 667, expected: 7},
		{name: "zero total", total: 0, threshold: 667, expected: 0},
		{name: "zero threshold", total: 5, threshold: 0, expected: 0},
		{name: "100% threshold", total: 5, threshold: 1000, expected: 5},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := types.CalculateThreshold(tc.total, tc.threshold)
			require.Equal(t, tc.expected, result)
		})
	}
}

// Helper to get a DKG network by specific round number.
func (k *Keeper) getDKGNetworkByRound(ctx context.Context, round uint32) (*types.DKGNetwork, error) {
	return k.getDKGNetwork(ctx, round)
}
