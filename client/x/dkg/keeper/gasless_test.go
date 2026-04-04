package keeper

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// TestGaslessSDKContext_DoesNotAffectCallerGas verifies that KV reads on a
// gasless context do not increment the original context's gas meter.
func TestGaslessSDKContext_DoesNotAffectCallerGas(t *testing.T) {
	t.Parallel()

	// 1. Create a real KV store and SDK context with a tracked gas meter.
	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	kvStore := runtime.NewKVStoreService(key)

	gasMeter := storetypes.NewGasMeter(1_000_000)
	sdkCtx := testCtx.Ctx.WithGasMeter(gasMeter)

	// 2. Write a key-value pair using the original context to have something to read.
	store := kvStore.OpenKVStore(sdkCtx)
	require.NoError(t, store.Set([]byte("test-key"), []byte("test-value")))
	gasAfterWrite := gasMeter.GasConsumed()
	require.Greater(t, gasAfterWrite, uint64(0), "writing should consume gas")

	// 3. Create a gasless context and read the key on it.
	gCtx := gaslessSDKContext(sdkCtx)
	gaslessStore := kvStore.OpenKVStore(gCtx)
	val, err := gaslessStore.Get([]byte("test-key"))
	require.NoError(t, err)
	require.Equal(t, []byte("test-value"), val, "gasless context should read correct data")

	// 4. Verify the original context's gas meter has NOT changed.
	gasAfterGaslessRead := gasMeter.GasConsumed()
	require.Equal(t, gasAfterWrite, gasAfterGaslessRead,
		"KV read on gasless context should not affect the original gas meter")

	// 5. Verify that a normal read on the original context DOES consume gas.
	_, err = store.Get([]byte("test-key"))
	require.NoError(t, err)
	gasAfterNormalRead := gasMeter.GasConsumed()
	require.Greater(t, gasAfterNormalRead, gasAfterGaslessRead,
		"KV read on original context should consume gas")
}

// TestGaslessSDKContext_InfiniteGasMeter verifies the gasless context uses
// an InfiniteGasMeter that never panics on gas consumption.
func TestGaslessSDKContext_InfiniteGasMeter(t *testing.T) {
	t.Parallel()

	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	sdkCtx := testCtx.Ctx.WithGasMeter(storetypes.NewGasMeter(1)) // very small meter

	gCtx := gaslessSDKContext(sdkCtx)
	gaslessSdkCtx := sdk.UnwrapSDKContext(gCtx)

	// The gasless context should have an infinite gas meter that does not panic.
	require.True(t, gaslessSdkCtx.GasMeter().IsOutOfGas() == false,
		"gasless context should have infinite gas meter")

	// Consume a large amount of gas on the gasless context — should not panic.
	require.NotPanics(t, func() {
		gaslessSdkCtx.GasMeter().ConsumeGas(1_000_000_000, "test")
	}, "infinite gas meter should not panic on large consumption")
}
