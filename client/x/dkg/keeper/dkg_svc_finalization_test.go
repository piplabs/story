package keeper

import (
	"context"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
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

// TestHandleDKGFinalization_TEEFinalizeDKGError_MarksFailed verifies that
// when callTEEFinalizeDKG fails (FinalizeDKG returns error on all retries),
// the session is marked PhaseFailed.
func TestHandleDKGFinalization_TEEFinalizeDKGError_MarksFailed(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("finalize-error-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		validatorEVMAddr: testValidatorAddr,
	}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	session := &types.DKGSession{
		Round:          20,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// FinalizeDKG returns error on all retries (retryAttemts total calls)
	mockKernel.EXPECT().FinalizeDKG(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("kernel finalize failed")).
		Times(retryAttemts)

	dkgNetwork := &types.DKGNetwork{
		Round:        20,
		Stage:        types.DKGStageFinalization,
		ActiveValSet: []string{testValidatorAddr},
	}

	k.handleDKGFinalization(ctx, dkgNetwork)

	got, getErr := sm.GetSession(20)
	require.NoError(t, getErr)
	require.Equal(t, types.PhaseFailed, got.Phase, "session should be marked failed when TEE finalization fails")
}

// TestHandleDKGFinalization_ContractFinalizeError_MarksFailed verifies that
// when callContractFinalizeDKG fails (Finalize returns an error), the session
// is marked PhaseFailed.
func TestHandleDKGFinalization_ContractFinalizeError_MarksFailed(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("contract-finalize-error-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		contractClient:   mockContract,
		validatorEVMAddr: testValidatorAddr,
	}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	session := &types.DKGSession{
		Round:          21,
		Phase:          types.PhaseDealing,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// TEE finalization succeeds
	mockKernel.EXPECT().FinalizeDKG(gomock.Any(), gomock.Any()).Return(
		&types.FinalizeDKGResponse{
			ParticipantsRoot: make([]byte, 32),
			GlobalPubKey:     []byte("global-pub"),
			Signature:        []byte("sig"),
			PublicCoeffs:     [][]byte{[]byte("c1")},
			PubKeyShare:      []byte("share"),
		}, nil,
	)

	// Contract Finalize returns an error
	mockContract.EXPECT().Finalize(
		gomock.Any(),
		uint32(21),
		gomock.Any(), // enclaveType
		gomock.Any(), // participantsRoot
		gomock.Any(), // globalPubKey
		gomock.Any(), // publicCoeffs
		gomock.Any(), // pubKeyShare
		gomock.Any(), // signature
	).Return(nil, errors.New("contract finalize failed"))

	dkgNetwork := &types.DKGNetwork{
		Round:        21,
		Stage:        types.DKGStageFinalization,
		ActiveValSet: []string{testValidatorAddr},
	}

	k.handleDKGFinalization(ctx, dkgNetwork)

	got, getErr := sm.GetSession(21)
	require.NoError(t, getErr)
	require.Equal(t, types.PhaseFailed, got.Phase, "session should be marked failed when contract finalization fails")
}

// TestCallTEEFinalizeDKG_AlreadyFinalized verifies that callTEEFinalizeDKG is a
// no-op when the session already has GlobalPubKey and SigFinalizeNetwork set.
func TestCallTEEFinalizeDKG_AlreadyFinalized(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	// Mock kernel — should NOT be called since we short-circuit early
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte("already-final-cc")
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{stateManager: sm, kernelRouter: router}

	// Session already has finalization data — callTEEFinalizeDKG should skip kernel call
	session := &types.DKGSession{
		Round:              22,
		Phase:              types.PhaseDealing,
		CodeCommitment:     cc,
		GlobalPubKey:       []byte("existing-global-pub"),
		SigFinalizeNetwork: []byte("existing-sig"),
	}

	// callTEEFinalizeDKG should return nil immediately without calling FinalizeDKG
	callErr := k.callTEEFinalizeDKG(ctx, session)
	require.NoError(t, callErr, "should skip finalize when already finalized")
}
