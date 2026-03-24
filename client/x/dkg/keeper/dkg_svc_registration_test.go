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

// --- resolveRegistrationKernelClient ---

func TestResolveRegistrationKernelClient_NormalRound_WithCC(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte("test-cc-normal")

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockClient)

	k := &Keeper{kernelRouter: router}

	client, resolvedCC, err := k.resolveRegistrationKernelClient(false, cc)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, cc, resolvedCC)
}

func TestResolveRegistrationKernelClient_NormalRound_NilCC_FirstClient(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	cc := []byte("first-client-cc")

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockClient)

	k := &Keeper{kernelRouter: router}

	// No previous CC — should fall back to first connected client
	client, resolvedCC, err := k.resolveRegistrationKernelClient(false, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, cc, resolvedCC)
}

func TestResolveRegistrationKernelClient_NormalRound_NoClients(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}

	_, _, err := k.resolveRegistrationKernelClient(false, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no kernel clients available")
}

func TestResolveRegistrationKernelClient_Upgrade_FindsNewBinary(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	oldCC := []byte("old-binary-cc")
	newCC := []byte("new-binary-cc")

	oldClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	newClient := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(oldCC, oldClient)
	router.RegisterClient(newCC, newClient)

	k := &Keeper{kernelRouter: router}

	// Upgrade: pass old CC, should find the new binary
	client, resolvedCC, err := k.resolveRegistrationKernelClient(true, oldCC)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotEqual(t, oldCC, resolvedCC, "should return new binary CC, not old")
	require.Equal(t, newCC, resolvedCC)
}

func TestResolveRegistrationKernelClient_Upgrade_NilOldCC(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}

	_, _, err := k.resolveRegistrationKernelClient(true, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "old code commitment required")
}

func TestResolveRegistrationKernelClient_Upgrade_NoNewBinary(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	oldCC := []byte("only-old-cc")
	oldClient := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(oldCC, oldClient)

	k := &Keeper{kernelRouter: router}

	// Only old binary connected — should fail
	_, _, err := k.resolveRegistrationKernelClient(true, oldCC)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no new kernel client found for upgrade")
}

// --- getClientWithReconnect ---

func TestGetClientWithReconnect_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	cc := []byte("reconnect-cc")
	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockClient)

	k := &Keeper{kernelRouter: router}

	client, err := k.getClientWithReconnect(cc)
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestGetClientWithReconnect_NotFound(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}

	// No clients registered — should fail even after reconnection attempt
	_, err := k.getClientWithReconnect([]byte("missing-cc"))
	require.Error(t, err)
}

// --- getRegistrationKernelClient ---

func TestGetRegistrationKernelClient_FirstAttemptSuccess(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	cc := []byte("reg-client-cc")
	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockClient)

	k := &Keeper{kernelRouter: router}

	client, resolvedCC, err := k.getRegistrationKernelClient(false, cc)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, cc, resolvedCC)
}

func TestGetRegistrationKernelClient_FallbackToFirstClient(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	cc := []byte("first-cc")
	mockClient := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockClient)

	k := &Keeper{kernelRouter: router}

	// Normal round, nil CC -> fallback to first client
	client, resolvedCC, err := k.getRegistrationKernelClient(false, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, cc, resolvedCC)
}

func TestGetRegistrationKernelClient_NoClientsAvailable(t *testing.T) {
	t.Parallel()

	router := NewKernelRouter(nil, nil)
	k := &Keeper{kernelRouter: router}

	// No clients, no endpoints -> reconnection attempt still fails
	_, _, err := k.getRegistrationKernelClient(false, nil)
	require.Error(t, err)
}

// --- Tests merged from dkg_svc_registration_guards_test.go ---

// NOTE: These tests are NOT parallel because they share the package-level dkgSvcRound atomic.

func TestHandleDKGRegistration_WrongStage(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        1,
		Stage:        types.DKGStageDealing, // Not registration
		ActiveValSet: []string{testValidatorAddr},
	}

	k.handleDKGRegistration(ctx, dkgNetwork, nil, false)

	// No session should be created since the stage check fails before session creation
	_, err := k.stateManager.GetSession(1)
	require.Error(t, err, "no session should exist when stage is wrong")
}

func TestHandleDKGRegistration_DuplicateAcquire(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Pre-acquire lock for round 2
	dkgSvcRound.Store(2)

	k.validatorEVMAddr = testValidatorAddr
	dkgNetwork := &types.DKGNetwork{
		Round:        2,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{testValidatorAddr},
	}

	k.handleDKGRegistration(ctx, dkgNetwork, nil, false)

	// No session created since lock dedup prevented processing
	_, err := k.stateManager.GetSession(2)
	require.Error(t, err, "no session should exist when lock is already held")
}

func TestHandleDKGRegistration_NotInCurRoundSet_SessionCreated(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	k.enclaveType = [32]byte{0x01}

	dkgNetwork := &types.DKGNetwork{
		Round:        3,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{"0xother1", "0xother2"}, // Validator not in set
	}

	k.handleDKGRegistration(ctx, dkgNetwork, nil, false)

	// Session should be created but skip key generation
	// (old members still need a session for dealing/finalization)
	session, err := k.stateManager.GetSession(3)
	require.NoError(t, err)
	require.NotNil(t, session)
	// Phase stays at Initializing because no key gen for non-members
	require.Equal(t, types.PhaseInitializing, session.Phase)
}

func TestHandleDKGRegistration_NotInCurSet_NoKeyGen(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	k.enclaveType = [32]byte{0x01}

	dkgNetwork := &types.DKGNetwork{
		Round:        4,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{"0xnewval"}, // testValidator not in current set
		IsResharing:  true,
	}

	// No kernelRouter needed since old-only members skip key generation
	k.handleDKGRegistration(ctx, dkgNetwork, nil, false)

	session, err := k.stateManager.GetSession(4)
	require.NoError(t, err)
	// Old-only member skips key gen; phase stays at Initializing
	require.Equal(t, types.PhaseInitializing, session.Phase)
}

func TestHandleDKGRegistration_UpgradeRound_SetsOldCC(t *testing.T) {
	k := setupKeeperWithStateManager(t)
	ctx := context.Background()

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	k.validatorEVMAddr = testValidatorAddr
	k.enclaveType = [32]byte{0x02}

	oldCC := []byte("old-code-commitment")
	dkgNetwork := &types.DKGNetwork{
		Round:        5,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{"0xother"}, // Validator not in current set (old-only member)
		IsResharing:  true,
		IsUpgrade:    true,
	}

	k.handleDKGRegistration(ctx, dkgNetwork, oldCC, false)

	session, err := k.stateManager.GetSession(5)
	require.NoError(t, err)
	require.NotNil(t, session)
	require.Equal(t, oldCC, session.OldCodeCommitment, "upgrade round should have old CC")
	// Old-only member (not in current set) should have CodeCommitment set to oldCC
	require.Equal(t, oldCC, session.CodeCommitment, "old-only member should use old CC as code commitment")
	require.True(t, session.IsUpgrade)
}

// --- Full path tests merged from dkg_svc_full_path_test.go ---

// TestHandleDKGRegistration_FullPath tests the full registration path with
// mocked kernel and contract clients, exercising callTEEGenerateAndSealKey
// and callContractRegister.

func TestHandleDKGRegistration_FullPath(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	cc := []byte("test-code-commitment")
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		contractClient:   mockContract,
		validatorEVMAddr: testValidatorAddr,
		enclaveType:      [32]byte{0x01},
	}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Setup kernel mock: GenerateAndSealKey returns test data
	mockKernel.EXPECT().GenerateAndSealKey(gomock.Any(), gomock.Any()).Return(
		&types.GenerateAndSealKeyResponse{
			CodeCommitment:   cc,
			DkgPubKey:        []byte("dkg-pub"),
			CommPubKey:       []byte("comm-pub"),
			EnclaveReport:    []byte("report"),
			StartBlockHeight: 100,
			StartBlockHash:   make([]byte, 32),
		}, nil,
	)

	// Setup contract mock: Register returns a successful receipt
	mockContract.EXPECT().Register(
		gomock.Any(),
		uint32(1),
		gomock.Any(), // enclaveType
		gomock.Any(), // startBlockHeight
		gomock.Any(), // startBlockHash
		gomock.Any(), // dkgPubKey
		gomock.Any(), // commPubKey
		gomock.Any(), // enclaveReport
	).Return(&ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}, nil)

	dkgNetwork := &types.DKGNetwork{
		Round:        1,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{testValidatorAddr},
	}

	k.handleDKGRegistration(ctx, dkgNetwork, nil, false)

	session, err := sm.GetSession(1)
	require.NoError(t, err)
	require.Equal(t, types.PhaseInitialized, session.Phase)
	require.Equal(t, cc, session.CodeCommitment)
	require.Equal(t, []byte("dkg-pub"), session.DKGPubKey)
}

// TestHandleDKGRegistration_AlreadyRegistered_SkipsContract tests that when
// alreadyRegistered=true, the contract Register call is skipped.

func TestHandleDKGRegistration_AlreadyRegistered(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	cc := []byte("test-code-commitment")
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		contractClient:   mockContract,
		validatorEVMAddr: testValidatorAddr,
		enclaveType:      [32]byte{0x01},
	}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	mockKernel.EXPECT().GenerateAndSealKey(gomock.Any(), gomock.Any()).Return(
		&types.GenerateAndSealKeyResponse{
			CodeCommitment:   cc,
			DkgPubKey:        []byte("dkg-pub"),
			CommPubKey:       []byte("comm-pub"),
			EnclaveReport:    []byte("report"),
			StartBlockHeight: 100,
			StartBlockHash:   make([]byte, 32),
		}, nil,
	)

	// Contract.Register should NOT be called when alreadyRegistered=true
	// (no mock expectation set = gomock will fail if called)

	dkgNetwork := &types.DKGNetwork{
		Round:        2,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{testValidatorAddr},
	}

	k.handleDKGRegistration(ctx, dkgNetwork, nil, true)

	session, err := sm.GetSession(2)
	require.NoError(t, err)
	require.Equal(t, types.PhaseInitialized, session.Phase)
}

// TestHandleDKGRegistration_GenerateAndSealKeyError_MarksFailed verifies that
// when callTEEGenerateAndSealKey fails (GenerateAndSealKey returns error on all
// retries), the session is marked as PhaseFailed.
func TestHandleDKGRegistration_GenerateAndSealKeyError_MarksFailed(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	cc := []byte("error-gen-cc")
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		validatorEVMAddr: testValidatorAddr,
		enclaveType:      [32]byte{0x03},
	}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Kernel returns error for all retry attempts (retryAttemts calls total)
	mockKernel.EXPECT().GenerateAndSealKey(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("kernel unavailable")).
		Times(retryAttemts)

	dkgNetwork := &types.DKGNetwork{
		Round:        10,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{testValidatorAddr},
	}

	k.handleDKGRegistration(ctx, dkgNetwork, nil, false)

	// Session should exist but be marked as PhaseFailed
	session, getErr := sm.GetSession(10)
	require.NoError(t, getErr)
	require.Equal(t, types.PhaseFailed, session.Phase, "session should be marked failed when GenerateAndSealKey fails")
}

// TestHandleDKGRegistration_ContractRegisterError_MarksFailed verifies that
// when callContractRegister fails (Register returns an error), the session is
// marked as PhaseFailed.
func TestHandleDKGRegistration_ContractRegisterError_MarksFailed(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	cc := []byte("error-register-cc")
	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager:     sm,
		kernelRouter:     router,
		contractClient:   mockContract,
		validatorEVMAddr: testValidatorAddr,
		enclaveType:      [32]byte{0x04},
	}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	// Kernel succeeds with GenerateAndSealKey
	mockKernel.EXPECT().GenerateAndSealKey(gomock.Any(), gomock.Any()).Return(
		&types.GenerateAndSealKeyResponse{
			CodeCommitment:   cc,
			DkgPubKey:        []byte("dkg-pub"),
			CommPubKey:       []byte("comm-pub"),
			EnclaveReport:    []byte("report"),
			StartBlockHeight: 50,
			StartBlockHash:   make([]byte, 32),
		}, nil,
	)

	// Contract Register returns an error
	mockContract.EXPECT().Register(
		gomock.Any(),
		gomock.Any(), // round
		gomock.Any(), // enclaveType
		gomock.Any(), // startBlockHeight
		gomock.Any(), // startBlockHash
		gomock.Any(), // dkgPubKey
		gomock.Any(), // commPubKey
		gomock.Any(), // enclaveReport
	).Return(nil, errors.New("contract call failed"))

	dkgNetwork := &types.DKGNetwork{
		Round:        11,
		Stage:        types.DKGStageRegistration,
		ActiveValSet: []string{testValidatorAddr},
	}

	k.handleDKGRegistration(ctx, dkgNetwork, nil, false)

	// Session should be marked PhaseFailed
	session, getErr := sm.GetSession(11)
	require.NoError(t, getErr)
	require.Equal(t, types.PhaseFailed, session.Phase, "session should be marked failed when Register contract fails")
}
