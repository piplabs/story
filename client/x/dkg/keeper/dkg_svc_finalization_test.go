package keeper

import (
	"context"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"
)

// NOTE: These tests are NOT parallel because they share the package-level
// dkgSvcRound atomic.

func TestHandleDKGFinalization_WrongStage(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        1,
		Stage:        types.DKGStageRegistration, // Not finalization
		ActiveValSet: []string{testValidatorAddr},
	}

	session := &types.DKGSession{
		Round: 1,
		Phase: types.PhaseDealing,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	k.handleDKGFinalization(ctx, dkgNetwork)

	got, err := k.stateManager.GetSession(1)
	require.NoError(t, err)
	require.Equal(t, types.PhaseDealing, got.Phase, "phase should not change when stage is wrong")
}

func TestHandleDKGFinalization_NotInCurRoundSet(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        2,
		Stage:        types.DKGStageFinalization,
		ActiveValSet: []string{"0xother1", "0xother2"},
	}

	session := &types.DKGSession{
		Round: 2,
		Phase: types.PhaseDealing,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	k.handleDKGFinalization(ctx, dkgNetwork)

	got, err := k.stateManager.GetSession(2)
	require.NoError(t, err)
	require.Equal(t, types.PhaseDealing, got.Phase, "phase should not change when validator not in set")
}

func TestHandleDKGFinalization_WrongPhase_MarksFailed(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        3,
		Stage:        types.DKGStageFinalization,
		ActiveValSet: []string{testValidatorAddr},
	}

	session := &types.DKGSession{
		Round: 3,
		Phase: types.PhaseInitialized,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	k.handleDKGFinalization(ctx, dkgNetwork)

	got, err := k.stateManager.GetSession(3)
	require.NoError(t, err)
	require.Equal(t, types.PhaseFailed, got.Phase, "non-dealing phase should be marked failed")
}

func TestHandleDKGFinalization_NoSession(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        99,
		Stage:        types.DKGStageFinalization,
		ActiveValSet: []string{testValidatorAddr},
	}

	// Should not panic when session doesn't exist
	k.handleDKGFinalization(ctx, dkgNetwork)
}

func TestHandleDKGFinalization_DuplicateAcquire(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Pre-acquire lock for round 5
	dkgSvcRound.Store(5)

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        5,
		Stage:        types.DKGStageFinalization,
		ActiveValSet: []string{testValidatorAddr},
	}

	session := &types.DKGSession{
		Round: 5,
		Phase: types.PhaseDealing,
	}
	require.NoError(t, k.stateManager.CreateSession(ctx, session))

	k.handleDKGFinalization(ctx, dkgNetwork)

	got, err := k.stateManager.GetSession(5)
	require.NoError(t, err)
	require.Equal(t, types.PhaseDealing, got.Phase, "lock dedup should prevent processing")
}

// --- Full path tests merged from dkg_svc_full_path_test.go ---

// TestHandleDKGFinalization_FullPath tests the full finalization path with
// mocked kernel and contract clients.

func TestHandleDKGFinalization_FullPath(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("finalize-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		contractClient:   mockContract,
		validatorEVMAddr: testValidatorAddr,
		enclaveType:      [32]byte{0x02},
	}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Create session in PhaseDealing
	session := &types.DKGSession{
		Round:          3,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Kernel FinalizeDKG returns finalization data
	mockKernel.EXPECT().FinalizeDKG(gomock.Any(), gomock.Any()).Return(
		&types.FinalizeDKGResponse{
			ParticipantsRoot: make([]byte, 32),
			GlobalPubKey:     []byte("global-pub"),
			Signature:        []byte("sig"),
			PublicCoeffs:     [][]byte{[]byte("c1"), []byte("c2")},
			PubKeyShare:      []byte("share"),
		}, nil,
	)

	// Contract Finalize returns success
	mockContract.EXPECT().Finalize(
		gomock.Any(),
		uint32(3),
		gomock.Any(), // enclaveType
		gomock.Any(), // participantsRoot
		gomock.Any(), // globalPubKey
		gomock.Any(), // publicCoeffs
		gomock.Any(), // pubKeyShare
		gomock.Any(), // signature
	).Return(&ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}, nil)

	dkgNetwork := &types.DKGNetwork{
		Round:        3,
		Stage:        types.DKGStageFinalization,
		ActiveValSet: []string{testValidatorAddr},
	}

	k.handleDKGFinalization(ctx, dkgNetwork)

	got, err := sm.GetSession(3)
	require.NoError(t, err)
	require.Equal(t, types.PhaseFinalized, got.Phase)
	require.Equal(t, []byte("global-pub"), got.GlobalPubKey)
	require.Equal(t, []byte("sig"), got.SigFinalizeNetwork)
}
