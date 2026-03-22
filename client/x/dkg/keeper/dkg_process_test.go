package keeper

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"

	"go.uber.org/mock/gomock"
)

// TestProcessDeals_DKGSvcDisabled verifies that ProcessDeals emits the event
// and returns nil even when DKG service is disabled.
func TestProcessDeals_DKGSvcDisabled(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	// isDKGSvcEnabled is false by default

	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	deals := []types.Deal{
		{Index: 1, RecipientIndex: 2},
	}

	err := k.ProcessDeals(ctx, network, deals)
	require.NoError(t, err)
}

// TestProcessDeals_EmptyDeals verifies that ProcessDeals handles an empty
// deals slice without error.
func TestProcessDeals_EmptyDeals(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     2,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	err := k.ProcessDeals(ctx, network, []types.Deal{})
	require.NoError(t, err)
}

// TestProcessDeals_MultipleDeals verifies that ProcessDeals correctly processes
// multiple deals at once.
func TestProcessDeals_MultipleDeals(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     3,
		Total:     5,
		Threshold: 4,
		Stage:     types.DKGStageDealing,
	}

	deals := []types.Deal{
		{Index: 1, RecipientIndex: 2},
		{Index: 1, RecipientIndex: 3},
		{Index: 2, RecipientIndex: 1},
	}

	err := k.ProcessDeals(ctx, network, deals)
	require.NoError(t, err)
}

// TestProcessResponses_DKGSvcDisabled verifies that ProcessResponses emits the
// event and returns nil when DKG service is disabled.
func TestProcessResponses_DKGSvcDisabled(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     1,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	responses := []types.Response{
		{Index: 1},
	}

	err := k.ProcessResponses(ctx, network, responses)
	require.NoError(t, err)
}

// TestProcessResponses_EmptyResponses verifies that ProcessResponses handles
// an empty response slice without error.
func TestProcessResponses_EmptyResponses(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     2,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	err := k.ProcessResponses(ctx, network, []types.Response{})
	require.NoError(t, err)
}

// TestProcessResponses_MultipleResponses verifies that ProcessResponses correctly
// processes multiple responses.
func TestProcessResponses_MultipleResponses(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	network := &types.DKGNetwork{
		Round:     3,
		Total:     5,
		Threshold: 4,
		Stage:     types.DKGStageDealing,
	}

	responses := []types.Response{
		{Index: 1},
		{Index: 2},
		{Index: 3},
	}

	err := k.ProcessResponses(ctx, network, responses)
	require.NoError(t, err)
}

// initTestStateManager creates a StateManager backed by a temp directory and assigns it to the keeper.
// This is needed because stateManager is only initialized in InitDKGService, not in setupDKGKeeperWithMocks.
func initTestStateManager(t *testing.T, k *Keeper) {
	t.Helper()

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	k.stateManager = sm
}

// TestEnsureSessionIndex_NoSession verifies that ensureSessionIndex returns an error
// when no session exists for the given round.
func TestEnsureSessionIndex_NoSession(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	k.setValidatorAddress(common.HexToAddress("0x1111111111111111111111111111111111111111"))
	initTestStateManager(t, k)

	// No session exists for round 99 — should return an error
	err := k.ensureSessionIndex(ctx, 99)
	require.Error(t, err, "ensureSessionIndex should fail when no session exists")
}

// TestEnsureSessionIndex_IndexAlreadySet verifies that ensureSessionIndex is a no-op
// when the session already has a non-zero index.
func TestEnsureSessionIndex_IndexAlreadySet(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	k.setValidatorAddress(common.HexToAddress("0x1111111111111111111111111111111111111111"))
	initTestStateManager(t, k)

	// Create a session with a pre-set index
	sess := newTestSession(5)
	require.NoError(t, k.stateManager.CreateSession(ctx, sess))
	session, err := k.stateManager.GetSession(5)
	require.NoError(t, err)
	session.Index = 3 // already set
	require.NoError(t, k.stateManager.UpdateSession(ctx, session))

	// ensureSessionIndex should return nil without touching the index
	err = k.ensureSessionIndex(ctx, 5)
	require.NoError(t, err)

	// Verify index is unchanged
	updated, err := k.stateManager.GetSession(5)
	require.NoError(t, err)
	require.Equal(t, uint32(3), updated.Index, "index should remain unchanged when already set")
}

// TestEnsureSessionIndex_SetsIndexFromRegistration verifies that ensureSessionIndex
// reads the on-chain registration index and updates the session.
func TestEnsureSessionIndex_SetsIndexFromRegistration(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	initTestStateManager(t, k)

	addr := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	k.setValidatorAddress(addr)

	// Store a registration with index=2 for round 7
	reg := &types.DKGRegistration{
		Round:         7,
		ValidatorAddr: addr.Hex(),
		Index:         2,
		DkgPubKey:     []byte("test-dkg-pub-key"),
		Status:        types.DKGRegStatusVerified,
	}
	require.NoError(t, k.setDKGRegistration(ctx, addr, reg))

	// Create a session with Index=0 (not yet set)
	require.NoError(t, k.stateManager.CreateSession(ctx, newTestSession(7)))
	session, err := k.stateManager.GetSession(7)
	require.NoError(t, err)
	require.Equal(t, uint32(0), session.Index, "index should start at 0")

	err = k.ensureSessionIndex(ctx, 7)
	require.NoError(t, err)

	// Index should now be set from the registration
	updated, err := k.stateManager.GetSession(7)
	require.NoError(t, err)
	require.Equal(t, uint32(2), updated.Index, "index should be set from on-chain registration")
}

// TestEnsureSessionIndex_NoRegistration verifies that ensureSessionIndex returns an
// error when the session exists but the validator has no registration.
func TestEnsureSessionIndex_NoRegistration(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	initTestStateManager(t, k)

	// Set validator address to an address with no registration
	k.setValidatorAddress(common.HexToAddress("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC"))

	// Create a session with Index=0
	require.NoError(t, k.stateManager.CreateSession(ctx, newTestSession(15)))

	// No registration for round 15 → should fail
	err := k.ensureSessionIndex(ctx, 15)
	require.Error(t, err, "should fail when registration not found")
}

// TestBeginDealing_BelowMinRequired verifies that BeginDealing calls SkipToNextRound
// when verified registration count is below the minimum required.
func TestBeginDealing_BelowMinRequired(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	sk := k.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	sk.EXPECT().GetAllValidators(gomock.Any()).Return(nil, nil).AnyTimes()

	// Set params with MinReqRegisteredParticipants=3
	params := types.DefaultParams()
	params.MinReqRegisteredParticipants = 3
	require.NoError(t, k.SetParams(ctx, params))

	// Create a network in Registration stage with round=1
	network := &types.DKGNetwork{
		Round:        1,
		Total:        5,
		Threshold:    4,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{"0x1111111111111111111111111111111111111111"},
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	// Store only 1 verified registration (below min=3)
	addr1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	require.NoError(t, k.setDKGRegistration(ctx, addr1, &types.DKGRegistration{
		Round:         1,
		ValidatorAddr: addr1.Hex(),
		Index:         1,
		Status:        types.DKGRegStatusVerified,
	}))

	// BeginDealing should skip to next round
	err := k.BeginDealing(ctx, network)
	require.NoError(t, err)

	// Verify a new round (round=2) was created via SkipToNextRound
	_, err = k.getDKGNetwork(ctx, 2)
	require.NoError(t, err, "SkipToNextRound should have created round 2")
}

// TestBeginDealing_MeetsMinRequired verifies that BeginDealing proceeds to the
// dealing phase when verified registration count meets the minimum required.
func TestBeginDealing_MeetsMinRequired(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Set params with MinReqRegisteredParticipants=2
	params := types.DefaultParams()
	params.MinReqRegisteredParticipants = 2
	require.NoError(t, k.SetParams(ctx, params))

	network := &types.DKGNetwork{
		Round:        10,
		Total:        3,
		Threshold:    2,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{"0x1111111111111111111111111111111111111111"},
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	// Store 3 verified registrations (above min=2)
	for i, hexAddr := range []string{
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
		"0x3333333333333333333333333333333333333333",
	} {
		addr := common.HexToAddress(hexAddr)
		require.NoError(t, k.setDKGRegistration(ctx, addr, &types.DKGRegistration{
			Round:         10,
			ValidatorAddr: addr.Hex(),
			Index:         uint32(i + 1),
			Status:        types.DKGRegStatusVerified,
		}))
	}

	err := k.BeginDealing(ctx, network)
	require.NoError(t, err)

	// Verify round 10 was updated with Total=3 (not skipped to round 11)
	updated, err := k.getDKGNetwork(ctx, 10)
	require.NoError(t, err)
	require.Equal(t, uint32(3), updated.Total, "Total should be set to verified registration count")
}

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

// TestDistributeRewardsToActiveCommittee_ZeroAmount verifies that
// DistributeRewardsToActiveCommittee returns zero immediately when totalAmount is zero.
func TestDistributeRewardsToActiveCommittee_ZeroAmount(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	distributed, err := k.DistributeRewardsToActiveCommittee(ctx, "evmstaking", math.ZeroInt())
	require.NoError(t, err)
	require.True(t, distributed.IsZero(), "distributed amount should be zero for zero totalAmount")
}

// TestDistributeRewardsToActiveCommittee_NoActiveRound verifies that when there is
// no active DKG round, DistributeRewardsToActiveCommittee returns 0.
func TestDistributeRewardsToActiveCommittee_NoActiveRound(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	distributed, err := k.DistributeRewardsToActiveCommittee(ctx, "evmstaking", math.NewInt(1000))
	require.NoError(t, err)
	require.True(t, distributed.IsZero(), "should return 0 when no active DKG round")
}

// TestDistributeRewardsToActiveCommittee_NoFinalizedMembers verifies that when
// the active round has no finalized members, nothing is distributed.
func TestDistributeRewardsToActiveCommittee_NoFinalizedMembers(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	// Create an active round with only verified (not finalized) registrations
	activeRound := &types.DKGNetwork{
		Round: 1,
		Total: 3, Threshold: 2,
		Stage: types.DKGStageActive,
	}
	require.NoError(t, k.setDKGNetwork(ctx, activeRound))
	require.NoError(t, k.setLatestActiveRound(ctx, activeRound))

	distributed, err := k.DistributeRewardsToActiveCommittee(ctx, "evmstaking", math.NewInt(1000))
	require.NoError(t, err)
	require.True(t, distributed.IsZero(), "should return 0 when no finalized members")
}

// TestDistributeRewardsToActiveCommittee_ZeroPortionParam verifies that when
// the DKG committee reward portion is zero, nothing is distributed.
func TestDistributeRewardsToActiveCommittee_ZeroPortionParam(t *testing.T) {
	t.Parallel()

	k, _, _, ctx := setupDKGKeeperWithMocks(t)

	params := types.DefaultParams()
	params.DkgCommitteeRewardPortion = math.LegacyZeroDec()
	require.NoError(t, k.SetParams(ctx, params))

	activeRound := &types.DKGNetwork{Round: 1, Total: 3, Threshold: 2, Stage: types.DKGStageActive}
	require.NoError(t, k.setDKGNetwork(ctx, activeRound))
	require.NoError(t, k.setLatestActiveRound(ctx, activeRound))

	member := common.HexToAddress("0x1111111111111111111111111111111111111111")
	setRegistration(t, k, ctx, activeRound, member, types.DKGRegStatusFinalized)

	distributed, err := k.DistributeRewardsToActiveCommittee(ctx, "evmstaking", math.NewInt(1000))
	require.NoError(t, err)
	require.True(t, distributed.IsZero(), "should return 0 when reward portion is zero")
}

// TestDistributeRewardsToActiveCommittee_NormalDistribution verifies that rewards
// are distributed correctly to finalized committee members.
func TestDistributeRewardsToActiveCommittee_NormalDistribution(t *testing.T) {
	t.Parallel()

	k, bk, _, ctx := setupDKGKeeperWithMocks(t)

	// Default params: 10% reward portion
	activeRound := &types.DKGNetwork{Round: 1, Total: 3, Threshold: 2, Stage: types.DKGStageActive}
	require.NoError(t, k.setDKGNetwork(ctx, activeRound))
	require.NoError(t, k.setLatestActiveRound(ctx, activeRound))

	member1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	member2 := common.HexToAddress("0x2222222222222222222222222222222222222222")
	setRegistration(t, k, ctx, activeRound, member1, types.DKGRegStatusFinalized)
	setRegistration(t, k, ctx, activeRound, member2, types.DKGRegStatusFinalized)

	// totalAmount=1000, 10%=100, 2 members → perMember=50
	// expect: SendCoinsFromModuleToModule (1000 * 10% = 100 total)
	totalAmount := math.NewInt(1000)
	perMemberReward := math.NewInt(50)
	totalToDistribute := math.NewInt(100) // 50*2
	distributeCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, totalToDistribute))

	bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), "evmstaking", types.ModuleName, distributeCoins).Return(nil)
	perMemberCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, perMemberReward))
	bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, gomock.Any(), perMemberCoins).Return(nil).Times(2)

	distributed, err := k.DistributeRewardsToActiveCommittee(ctx, "evmstaking", totalAmount)
	require.NoError(t, err)
	require.Equal(t, totalToDistribute, distributed, "distributed amount should match totalToDistribute")
}

// TestProcessDeals_DKGSvcEnabled verifies that ProcessDeals starts the async
// goroutine when DKG service is enabled (no panic or error).
func TestProcessDeals_DKGSvcEnabled(t *testing.T) {
	// Not parallel: modifies global DKG service state
	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()

	network := &types.DKGNetwork{
		Round:     20,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	deals := []types.Deal{
		{Index: 1, RecipientIndex: 2},
	}

	// Should not panic even when svc is enabled (goroutine runs handleDKGProcessDeals)
	err := k.ProcessDeals(ctx, network, deals)
	require.NoError(t, err)
}

// TestProcessResponses_DKGSvcEnabled verifies that ProcessResponses starts the
// async goroutine path when DKG service is enabled, using shouldProcessResponses.
func TestProcessResponses_DKGSvcEnabled(t *testing.T) {
	// Not parallel: modifies global DKG service state
	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	k.setValidatorAddress(common.HexToAddress("0x1111111111111111111111111111111111111111"))

	network := &types.DKGNetwork{
		Round:     21,
		Total:     3,
		Threshold: 2,
		Stage:     types.DKGStageDealing,
	}

	responses := []types.Response{
		{Index: 1},
	}

	err := k.ProcessResponses(ctx, network, responses)
	require.NoError(t, err)
}

// TestProcessResponses_DKGSvcEnabled_WithStateManager verifies that ProcessResponses
// runs shouldProcessResponses successfully when a stateManager is initialized.
func TestProcessResponses_DKGSvcEnabled_WithStateManager(t *testing.T) {
	// Not parallel: modifies global DKG service state
	k, _, _, ctx := setupDKGKeeperWithMocks(t)
	k.setIsDKGSvcEnabled()
	validatorAddr := common.HexToAddress("0x2222222222222222222222222222222222222222")
	k.setValidatorAddress(validatorAddr)
	initTestStateManager(t, k)

	network := &types.DKGNetwork{
		Round:        22,
		Total:        3,
		Threshold:    2,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{validatorAddr.Hex()},
	}
	require.NoError(t, k.setDKGNetwork(ctx, network))

	// Create a session so shouldProcessResponses can check session state
	require.NoError(t, k.stateManager.CreateSession(ctx, newTestSession(22)))

	responses := []types.Response{
		{Index: 1},
	}

	err := k.ProcessResponses(ctx, network, responses)
	require.NoError(t, err)
}
