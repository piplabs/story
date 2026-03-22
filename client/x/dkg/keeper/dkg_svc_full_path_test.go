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

// TestHandleDKGDealing_FullPath tests the full dealing path with mocked kernel client.
func TestHandleDKGDealing_FullPath(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)

	sm, err := NewStateManager(t.TempDir())
	require.NoError(t, err)

	cc := []byte("deal-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		stateManager: sm,
		kernelRouter: router,
	}

	resetDKGSvcRound()
	defer resetDKGSvcRound()

	session := &types.DKGSession{
		Round:          4,
		Phase:          types.PhaseInitialized,
		CodeCommitment: cc,
	}
	require.NoError(t, sm.CreateSession(ctx, session))

	// Kernel GenerateDeals returns deals
	testDeals := []types.Deal{
		{Index: 0, RecipientIndex: 1, Deal: types.EncryptedDeal{DhKey: []byte("deal1")}},
	}
	mockKernel.EXPECT().GenerateDeals(gomock.Any(), gomock.Any()).Return(
		&types.GenerateDealsResponse{Deals: testDeals}, nil,
	)

	dkgNetwork := &types.DKGNetwork{
		Round: 4,
		Stage: types.DKGStageDealing,
	}

	// Flush queues first
	k.FlushAllQueues()
	k.handleDKGDealing(ctx, dkgNetwork, true)

	got, err := sm.GetSession(4)
	require.NoError(t, err)
	require.Equal(t, types.PhaseDealing, got.Phase)

	// Verify deals were enqueued
	dequeued := k.DequeueDeals(10)
	require.Len(t, dequeued, 1)
	require.Equal(t, uint32(0), dequeued[0].Index)
}
