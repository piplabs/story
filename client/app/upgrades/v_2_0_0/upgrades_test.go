package v_2_0_0

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/core/event"
	storetypes "cosmossdk.io/store/types"

	"google.golang.org/protobuf/runtime/protoiface"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/app/keepers"
)

// noopEventManager satisfies event.Manager for tests.
type noopEventManager struct{}

func (noopEventManager) Emit(_ context.Context, _ protoiface.MessageV1) error { return nil }
func (noopEventManager) EmitKV(_ context.Context, _ string, _ ...event.Attribute) error {
	return nil
}
func (noopEventManager) EmitNonConsensus(_ context.Context, _ protoiface.MessageV1) error {
	return nil
}

// noopEventService satisfies event.Service for tests.
type noopEventService struct{}

func (noopEventService) EventManager(_ context.Context) event.Manager {
	return noopEventManager{}
}

// newTestConsensusKeeper creates a real consensus keeper backed by in-memory storage.
func newTestConsensusKeeper(t *testing.T) (consensuskeeper.Keeper, sdk.Context) {
	t.Helper()

	key := storetypes.NewKVStoreKey("consensus")
	ctx := sdktestutil.DefaultContext(key, storetypes.NewTransientStoreKey("consensus_transient"))

	ir := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(ir)

	authority := authtypes.NewModuleAddress("gov").String()
	keeper := consensuskeeper.NewKeeper(cdc, runtime.NewKVStoreService(key), authority, noopEventService{})

	return keeper, ctx
}

func TestEnableVoteExtensions_NoExistingAbci(t *testing.T) {
	cpKeeper, ctx := newTestConsensusKeeper(t)

	// Set initial consensus params without ABCI section
	require.NoError(t, cpKeeper.ParamsStore.Set(ctx, cmtproto.ConsensusParams{
		Block: &cmtproto.BlockParams{
			MaxBytes: 1048576,
			MaxGas:   -1,
		},
		Evidence:  &cmtproto.EvidenceParams{MaxAgeNumBlocks: 100000, MaxAgeDuration: 48 * time.Hour, MaxBytes: 1048576},
		Validator: &cmtproto.ValidatorParams{PubKeyTypes: []string{"ed25519"}},
	}))

	keepers := &keepers.Keepers{ConsensusParamsKeeper: cpKeeper}

	require.NoError(t, enableVoteExtensions(ctx, keepers, 100))

	// Verify vote extensions enabled at upgradeHeight + 1
	params, err := cpKeeper.ParamsStore.Get(ctx)
	require.NoError(t, err)
	require.NotNil(t, params.Abci)
	require.Equal(t, int64(101), params.Abci.VoteExtensionsEnableHeight)
}

func TestEnableVoteExtensions_AbciExistsButVEDisabled(t *testing.T) {
	cpKeeper, ctx := newTestConsensusKeeper(t)

	// ABCI section exists but VE height is 0 (disabled)
	require.NoError(t, cpKeeper.ParamsStore.Set(ctx, cmtproto.ConsensusParams{
		Block:     &cmtproto.BlockParams{MaxBytes: 1048576, MaxGas: -1},
		Evidence:  &cmtproto.EvidenceParams{MaxAgeNumBlocks: 100000, MaxAgeDuration: 48 * time.Hour, MaxBytes: 1048576},
		Validator: &cmtproto.ValidatorParams{PubKeyTypes: []string{"ed25519"}},
		Abci:      &cmtproto.ABCIParams{VoteExtensionsEnableHeight: 0},
	}))

	keepers := &keepers.Keepers{ConsensusParamsKeeper: cpKeeper}

	require.NoError(t, enableVoteExtensions(ctx, keepers, 200))

	params, err := cpKeeper.ParamsStore.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(201), params.Abci.VoteExtensionsEnableHeight)
}

func TestEnableVoteExtensions_AlreadyEnabled_SameHeight(t *testing.T) {
	cpKeeper, ctx := newTestConsensusKeeper(t)

	// VE already enabled at the expected height (upgradeHeight+1 = 101)
	require.NoError(t, cpKeeper.ParamsStore.Set(ctx, cmtproto.ConsensusParams{
		Block:     &cmtproto.BlockParams{MaxBytes: 1048576, MaxGas: -1},
		Evidence:  &cmtproto.EvidenceParams{MaxAgeNumBlocks: 100000, MaxAgeDuration: 48 * time.Hour, MaxBytes: 1048576},
		Validator: &cmtproto.ValidatorParams{PubKeyTypes: []string{"ed25519"}},
		Abci:      &cmtproto.ABCIParams{VoteExtensionsEnableHeight: 101},
	}))

	keepers := &keepers.Keepers{ConsensusParamsKeeper: cpKeeper}

	// Should skip without error
	require.NoError(t, enableVoteExtensions(ctx, keepers, 100))

	// Height unchanged
	params, err := cpKeeper.ParamsStore.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(101), params.Abci.VoteExtensionsEnableHeight)
}

func TestEnableVoteExtensions_AlreadyEnabled_DifferentHeight(t *testing.T) {
	cpKeeper, ctx := newTestConsensusKeeper(t)

	// VE enabled at height 1 (fresh genesis), but upgrade at height 100
	require.NoError(t, cpKeeper.ParamsStore.Set(ctx, cmtproto.ConsensusParams{
		Block:     &cmtproto.BlockParams{MaxBytes: 1048576, MaxGas: -1},
		Evidence:  &cmtproto.EvidenceParams{MaxAgeNumBlocks: 100000, MaxAgeDuration: 48 * time.Hour, MaxBytes: 1048576},
		Validator: &cmtproto.ValidatorParams{PubKeyTypes: []string{"ed25519"}},
		Abci:      &cmtproto.ABCIParams{VoteExtensionsEnableHeight: 1},
	}))

	keepers := &keepers.Keepers{ConsensusParamsKeeper: cpKeeper}

	// Should skip without error (warns about mismatch)
	require.NoError(t, enableVoteExtensions(ctx, keepers, 100))

	// Height unchanged at original value
	params, err := cpKeeper.ParamsStore.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(1), params.Abci.VoteExtensionsEnableHeight)
}

func TestEnableVoteExtensions_Idempotent(t *testing.T) {
	cpKeeper, ctx := newTestConsensusKeeper(t)

	require.NoError(t, cpKeeper.ParamsStore.Set(ctx, cmtproto.ConsensusParams{
		Block:     &cmtproto.BlockParams{MaxBytes: 1048576, MaxGas: -1},
		Evidence:  &cmtproto.EvidenceParams{MaxAgeNumBlocks: 100000, MaxAgeDuration: 48 * time.Hour, MaxBytes: 1048576},
		Validator: &cmtproto.ValidatorParams{PubKeyTypes: []string{"ed25519"}},
	}))

	keepers := &keepers.Keepers{ConsensusParamsKeeper: cpKeeper}

	// First call enables
	require.NoError(t, enableVoteExtensions(ctx, keepers, 100))

	params, err := cpKeeper.ParamsStore.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(101), params.Abci.VoteExtensionsEnableHeight)

	// Second call is a no-op
	require.NoError(t, enableVoteExtensions(ctx, keepers, 100))

	params, err = cpKeeper.ParamsStore.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(101), params.Abci.VoteExtensionsEnableHeight)
}
