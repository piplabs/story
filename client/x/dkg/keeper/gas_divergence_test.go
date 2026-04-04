package keeper

import (
	"testing"

	storetypes "cosmossdk.io/store/types"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"

	"go.uber.org/mock/gomock"
)

// setupGasTestKeeper creates a keeper with an SDK context backed by a real
// KV store so that gas consumption from store reads is observable.
func setupGasTestKeeper(t *testing.T) (*Keeper, sdk.Context) {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig()

	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	ak := dkgtestutil.NewMockAccountKeeper(ctrl)
	sk := dkgtestutil.NewMockStakingKeeper(ctrl)
	bk := dkgtestutil.NewMockBankKeeper(ctrl)
	dk := dkgtestutil.NewMockDistributionKeeper(ctrl)

	ak.EXPECT().AddressCodec().Return(authcodec.NewBech32Codec("story")).AnyTimes()
	ak.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.AccAddress{}).AnyTimes()

	var valStore baseapp.ValidatorStore = nil

	mockKernelClient := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContractClient := dkgtestutil.NewMockDKGContractClient(ctrl)

	kernelRouter := NewKernelRouter(nil, nil)
	kernelRouter.RegisterClient([]byte("test"), mockKernelClient)

	k := NewKeeper(
		encCfg.Codec,
		storeService,
		ak,
		bk,
		dk,
		sk,
		valStore,
		kernelRouter,
		mockContractClient,
		testAuthority,
	)

	require.NoError(t, k.SetParams(testCtx.Ctx, types.DefaultParams()))

	return k, testCtx.Ctx
}

// simulateAddVoteTxResult executes AddVote inside a fresh InfiniteGasMeter
// context (matching production BaseApp.getContextForTx) and returns the
// ExecTxResult that BaseApp.deliverTx would produce.
func simulateAddVoteTxResult(t *testing.T, k *Keeper, ctx sdk.Context, msg *types.MsgAddDkgVote) *abci.ExecTxResult {
	t.Helper()

	gasMeter := storetypes.NewInfiniteGasMeter()
	gasCtx := ctx.WithGasMeter(gasMeter)

	srv := NewMsgServerImpl(k)
	resp, err := srv.AddVote(gasCtx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)

	return &abci.ExecTxResult{
		Code:      0,
		GasWanted: 0,
		GasUsed:   int64(gasMeter.GasConsumed()),
	}
}

// computeLastResultsHash computes LastResultsHash from ExecTxResults using
// CometBFT's exact logic (types/results.go:13-23).
func computeLastResultsHash(results []*abci.ExecTxResult) []byte {
	return cmttypes.NewResults(results).Hash()
}

// TestAddVote_Determinism_Responses verifies that DKG-enabled and DKG-disabled
// nodes produce identical GasUsed and LastResultsHash when processing a
// MsgAddDkgVote containing responses.
//
// This is the critical regression test for the consensus bug where
// shouldProcessResponses() performed conditional KV reads inside the
// isDKGSvcEnabled block, causing different gas consumption.
func TestAddVote_Determinism_Responses(t *testing.T) {
	t.Parallel()

	kTEE, ctxTEE := setupGasTestKeeper(t)
	kNonTEE, ctxNonTEE := setupGasTestKeeper(t)

	kTEE.isDKGSvcEnabled = true
	kTEE.validatorEVMAddr = "0x1234567890abcdef1234567890abcdef12345678"
	kNonTEE.isDKGSvcEnabled = false

	network := &types.DKGNetwork{
		Round:        2,
		Total:        4,
		Threshold:    3,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{"0xaaaa", "0xbbbb", "0xcccc", "0xdddd"},
	}
	require.NoError(t, kTEE.setDKGNetwork(ctxTEE, network))
	require.NoError(t, kNonTEE.setDKGNetwork(ctxNonTEE, network))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote: &types.Vote{
			Responses: []types.Response{
				{Index: 0}, {Index: 1}, {Index: 2}, {Index: 3},
			},
		},
	}

	resultTEE := simulateAddVoteTxResult(t, kTEE, ctxTEE, msg)
	resultNonTEE := simulateAddVoteTxResult(t, kNonTEE, ctxNonTEE, msg)

	hashTEE := computeLastResultsHash([]*abci.ExecTxResult{resultTEE})
	hashNonTEE := computeLastResultsHash([]*abci.ExecTxResult{resultNonTEE})

	t.Logf("TEE GasUsed=%d, NonTEE GasUsed=%d", resultTEE.GasUsed, resultNonTEE.GasUsed)
	t.Logf("TEE Hash=%X, NonTEE Hash=%X", hashTEE, hashNonTEE)

	require.Equal(t, resultTEE.GasUsed, resultNonTEE.GasUsed,
		"GasUsed must be identical between DKG-enabled and DKG-disabled nodes")
	require.Equal(t, hashTEE, hashNonTEE,
		"LastResultsHash must be identical — no consensus failure")
}

// TestAddVote_Determinism_ResponsesWithActiveRound verifies determinism when
// an active round exists, which triggers additional KV reads in
// shouldProcessResponses (LatestActiveRound + DKGNetworks.Get).
func TestAddVote_Determinism_ResponsesWithActiveRound(t *testing.T) {
	t.Parallel()

	kTEE, ctxTEE := setupGasTestKeeper(t)
	kNonTEE, ctxNonTEE := setupGasTestKeeper(t)

	kTEE.isDKGSvcEnabled = true
	kTEE.validatorEVMAddr = "0x1234567890abcdef1234567890abcdef12345678"
	kNonTEE.isDKGSvcEnabled = false

	activeRound := &types.DKGNetwork{
		Round:        1,
		Total:        4,
		Threshold:    3,
		Stage:        types.DKGStageActive,
		ActiveValSet: []string{"0xaaaa", "0xbbbb", "0xcccc", "0xdddd"},
	}
	require.NoError(t, kTEE.setDKGNetwork(ctxTEE, activeRound))
	require.NoError(t, kNonTEE.setDKGNetwork(ctxNonTEE, activeRound))
	require.NoError(t, kTEE.setLatestActiveRound(ctxTEE, activeRound))
	require.NoError(t, kNonTEE.setLatestActiveRound(ctxNonTEE, activeRound))

	network := &types.DKGNetwork{
		Round:        2,
		Total:        4,
		Threshold:    3,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{"0xaaaa", "0xbbbb", "0xcccc", "0xdddd"},
	}
	require.NoError(t, kTEE.setDKGNetwork(ctxTEE, network))
	require.NoError(t, kNonTEE.setDKGNetwork(ctxNonTEE, network))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote: &types.Vote{
			Responses: []types.Response{{Index: 0}, {Index: 1}},
		},
	}

	resultTEE := simulateAddVoteTxResult(t, kTEE, ctxTEE, msg)
	resultNonTEE := simulateAddVoteTxResult(t, kNonTEE, ctxNonTEE, msg)

	hashTEE := computeLastResultsHash([]*abci.ExecTxResult{resultTEE})
	hashNonTEE := computeLastResultsHash([]*abci.ExecTxResult{resultNonTEE})

	t.Logf("TEE GasUsed=%d, NonTEE GasUsed=%d", resultTEE.GasUsed, resultNonTEE.GasUsed)

	require.Equal(t, resultTEE.GasUsed, resultNonTEE.GasUsed,
		"GasUsed must be identical (with active round)")
	require.Equal(t, hashTEE, hashNonTEE,
		"LastResultsHash must be identical (with active round)")
}

// TestAddVote_Determinism_DealsOnly verifies that deals-only votes produce
// identical LastResultsHash (no conditional KV reads in ProcessDeals).
func TestAddVote_Determinism_DealsOnly(t *testing.T) {
	t.Parallel()

	kTEE, ctxTEE := setupGasTestKeeper(t)
	kNonTEE, ctxNonTEE := setupGasTestKeeper(t)

	kTEE.isDKGSvcEnabled = true
	kTEE.validatorEVMAddr = "0x1234567890abcdef1234567890abcdef12345678"
	kNonTEE.isDKGSvcEnabled = false

	network := &types.DKGNetwork{
		Round:        2,
		Total:        4,
		Threshold:    3,
		Stage:        types.DKGStageDealing,
		ActiveValSet: []string{"0xaaaa", "0xbbbb", "0xcccc", "0xdddd"},
	}
	require.NoError(t, kTEE.setDKGNetwork(ctxTEE, network))
	require.NoError(t, kNonTEE.setDKGNetwork(ctxNonTEE, network))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote: &types.Vote{
			Deals: []types.Deal{
				{Index: 0, RecipientIndex: 1},
				{Index: 0, RecipientIndex: 2},
			},
		},
	}

	resultTEE := simulateAddVoteTxResult(t, kTEE, ctxTEE, msg)
	resultNonTEE := simulateAddVoteTxResult(t, kNonTEE, ctxNonTEE, msg)

	hashTEE := computeLastResultsHash([]*abci.ExecTxResult{resultTEE})
	hashNonTEE := computeLastResultsHash([]*abci.ExecTxResult{resultNonTEE})

	require.Equal(t, hashTEE, hashNonTEE,
		"Deals-only: LastResultsHash must be identical")
}

// TestAddVote_Determinism_EmptyVote verifies that empty votes produce
// identical LastResultsHash.
func TestAddVote_Determinism_EmptyVote(t *testing.T) {
	t.Parallel()

	kTEE, ctxTEE := setupGasTestKeeper(t)
	kNonTEE, ctxNonTEE := setupGasTestKeeper(t)

	kTEE.isDKGSvcEnabled = true
	kTEE.validatorEVMAddr = "0x1234567890abcdef1234567890abcdef12345678"
	kNonTEE.isDKGSvcEnabled = false

	network := &types.DKGNetwork{
		Round:     1,
		Total:     4,
		Threshold: 3,
		Stage:     types.DKGStageDealing,
	}
	require.NoError(t, kTEE.setDKGNetwork(ctxTEE, network))
	require.NoError(t, kNonTEE.setDKGNetwork(ctxNonTEE, network))

	msg := &types.MsgAddDkgVote{
		Authority: testAuthority,
		Vote:      &types.Vote{},
	}

	resultTEE := simulateAddVoteTxResult(t, kTEE, ctxTEE, msg)
	resultNonTEE := simulateAddVoteTxResult(t, kNonTEE, ctxNonTEE, msg)

	hashTEE := computeLastResultsHash([]*abci.ExecTxResult{resultTEE})
	hashNonTEE := computeLastResultsHash([]*abci.ExecTxResult{resultNonTEE})

	require.Equal(t, hashTEE, hashNonTEE,
		"Empty vote: LastResultsHash must be identical")
}

// TestBeginBlocker_Determinism_DealingTransition verifies that the KV reads
// in BeginDealing (shouldDeal + getDKGRegistration for ensureSessionIndex)
// produce identical gas on TEE and non-TEE nodes.
func TestBeginBlocker_Determinism_DealingTransition(t *testing.T) {
	t.Parallel()

	kTEE, ctxTEE := setupGasTestKeeper(t)
	kNonTEE, ctxNonTEE := setupGasTestKeeper(t)

	kTEE.isDKGSvcEnabled = true
	kTEE.validatorEVMAddr = "0x1234567890abcdef1234567890abcdef12345678"
	kNonTEE.isDKGSvcEnabled = false

	network := &types.DKGNetwork{
		Round:            1,
		Total:            3,
		Threshold:        2,
		Stage:            types.DKGStageDealing,
		ActiveValSet:     []string{"0xaaaa", "0xbbbb", "0xcccc"},
		StartBlockHeight: 100,
	}

	// Simulate the unconditional KV reads that BeginDealing performs
	runKVReads := func(k *Keeper, ctx sdk.Context) uint64 {
		gasMeter := storetypes.NewInfiniteGasMeter()
		gasCtx := ctx.WithGasMeter(gasMeter)
		require.NoError(t, k.setDKGNetwork(gasCtx, network))

		// These are the reads moved outside isDKGSvcEnabled
		_, _ = k.shouldDeal(gasCtx, network)
		_, _ = k.getDKGRegistration(gasCtx, network.Round, common.HexToAddress(k.validatorEVMAddr))

		return gasMeter.GasConsumed()
	}

	gasTEE := runKVReads(kTEE, ctxTEE)
	gasNonTEE := runKVReads(kNonTEE, ctxNonTEE)

	t.Logf("BeginDealing KV reads gas — TEE: %d, NonTEE: %d", gasTEE, gasNonTEE)

	require.Equal(t, gasTEE, gasNonTEE,
		"BeginDealing: KV read gas must be identical between TEE and non-TEE nodes")
}

// TestBeginBlocker_Determinism_InitiateDKGRound verifies that InitiateDKGRound
// gas is identical, with isAlreadyRegistered and getOldCodeCommitment now
// pre-computed outside isDKGSvcEnabled.
func TestBeginBlocker_Determinism_InitiateDKGRound(t *testing.T) {
	t.Parallel()

	kTEE, ctxTEE := setupGasTestKeeper(t)
	kNonTEE, ctxNonTEE := setupGasTestKeeper(t)

	kTEE.isDKGSvcEnabled = true
	kTEE.validatorEVMAddr = "0x1234567890abcdef1234567890abcdef12345678"
	smTEE, smErr := NewStateManager(t.TempDir())
	require.NoError(t, smErr)
	kTEE.stateManager = smTEE
	kNonTEE.isDKGSvcEnabled = false

	// Mock staking keeper for GetActiveValidators
	skTEE := kTEE.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	skNonTEE := kNonTEE.stakingKeeper.(*dkgtestutil.MockStakingKeeper)
	skTEE.EXPECT().GetAllValidators(gomock.Any()).Return(nil, nil).Times(1)
	skNonTEE.EXPECT().GetAllValidators(gomock.Any()).Return(nil, nil).Times(1)

	runInitiate := func(k *Keeper, ctx sdk.Context) uint64 {
		gasMeter := storetypes.NewInfiniteGasMeter()
		gasCtx := ctx.WithGasMeter(gasMeter)
		_ = k.InitiateDKGRound(gasCtx, false)
		return gasMeter.GasConsumed()
	}

	gasTEE := runInitiate(kTEE, ctxTEE)
	gasNonTEE := runInitiate(kNonTEE, ctxNonTEE)

	t.Logf("InitiateDKGRound gas — TEE: %d, NonTEE: %d", gasTEE, gasNonTEE)

	require.Equal(t, gasTEE, gasNonTEE,
		"InitiateDKGRound: gas must be identical between TEE and non-TEE nodes")
}

// TestBeginBlocker_Determinism_ResumeDKGService verifies that the KV reads
// pre-computed for ResumeDKGService (isAlreadyRegistered, getOldCodeCommitment,
// shouldDeal, getDKGRegistration) produce identical gas on all nodes.
func TestBeginBlocker_Determinism_ResumeDKGService(t *testing.T) {
	t.Parallel()

	kTEE, ctxTEE := setupGasTestKeeper(t)
	kNonTEE, ctxNonTEE := setupGasTestKeeper(t)

	kTEE.isDKGSvcEnabled = true
	kTEE.validatorEVMAddr = "0x1234567890abcdef1234567890abcdef12345678"
	kNonTEE.isDKGSvcEnabled = false

	network := &types.DKGNetwork{
		Round:            2,
		Total:            3,
		Threshold:        2,
		Stage:            types.DKGStageDealing,
		ActiveValSet:     []string{"0xaaaa", "0xbbbb"},
		StartBlockHeight: 100,
	}
	require.NoError(t, kTEE.setDKGNetwork(ctxTEE, network))
	require.NoError(t, kNonTEE.setDKGNetwork(ctxNonTEE, network))

	// Simulate the pre-computation that BeginBlocker does before isDKGSvcEnabled
	runPreCompute := func(k *Keeper, ctx sdk.Context) uint64 {
		gasMeter := storetypes.NewInfiniteGasMeter()
		gasCtx := ctx.WithGasMeter(gasMeter)

		_ = k.isAlreadyRegistered(gasCtx, network.Round)
		_, _ = k.getOldCodeCommitment(gasCtx)
		_, _ = k.shouldDeal(gasCtx, network)
		_, _ = k.getDKGRegistration(gasCtx, network.Round, common.HexToAddress(k.validatorEVMAddr))

		return gasMeter.GasConsumed()
	}

	gasTEE := runPreCompute(kTEE, ctxTEE)
	gasNonTEE := runPreCompute(kNonTEE, ctxNonTEE)

	t.Logf("ResumeDKGService pre-compute gas — TEE: %d, NonTEE: %d", gasTEE, gasNonTEE)

	require.Equal(t, gasTEE, gasNonTEE,
		"ResumeDKGService pre-compute: gas must be identical between TEE and non-TEE nodes")
}
