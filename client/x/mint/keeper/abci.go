package keeper

import (
	"context"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/mint/types"
	"github.com/piplabs/story/lib/log"
)

// BeginBlocker mints new tokens for the previous block.
func (k Keeper) BeginBlocker(ctx context.Context, ic types.InflationCalculationFn) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	singularityHeight, err := k.stakingKeeper.GetSingularityHeight(ctx)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if sdkCtx.BlockHeight() < int64(singularityHeight) {
		log.Debug(ctx, "Skip minting during singularity")
		return nil
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	// mint coins, update supply
	mintedCoinAmt := ic(ctx, params, math.LegacyNewDec(0)) // NOTE: bondedRatio is not used in current implementation.
	mintedCoin := sdk.NewCoin(params.MintDenom, mintedCoinAmt.TruncateInt())

	mintedCoins := sdk.NewCoins(mintedCoin)

	if err := k.MintCoins(ctx, mintedCoins); err != nil {
		return err
	}

	// send the minted coins to the fee collector account
	if err := k.AddCollectedFees(ctx, mintedCoins); err != nil {
		return err
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)

	return nil
}
