package keeper

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// gaslessSDKContext returns a context with a fresh InfiniteGasMeter so that
// Cosmos KV store reads performed on it do not contribute to the caller's
// GasUsed. This is used inside isDKGSvcEnabled blocks where DKG-enabled
// nodes perform KV reads that DKG-disabled nodes skip. Without this, the
// different gas consumption produces different LastResultsHash and causes
// consensus failure.
//
// The underlying multi-store reference is unchanged: reads see the current
// block's latest uncommitted state (unlike ABCI queries which see N-1).
func gaslessSDKContext(ctx context.Context) context.Context {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.WithGasMeter(storetypes.NewInfiniteGasMeter())
}
