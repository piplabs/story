package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/log"
)

func (k *Keeper) emitBeginDKGInitialization(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	validatorsStr := ""
	for i, val := range dkgNetwork.ActiveValSet {
		if i > 0 && i < len(dkgNetwork.ActiveValSet)-1 {
			validatorsStr += ","
		}
		validatorsStr += val
	}

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent("dkg_begin_initialization",
			sdk.NewAttribute("active_validator_set", validatorsStr),
			sdk.NewAttribute("mrenclave", string(dkgNetwork.Mrenclave)),
			sdk.NewAttribute("round", fmt.Sprintf("%d", dkgNetwork.Round)),
			sdk.NewAttribute("total", fmt.Sprintf("%d", len(dkgNetwork.ActiveValSet))),
			sdk.NewAttribute("start_block", fmt.Sprintf("%d", dkgNetwork.StartBlock)),
		),
	})

	log.Info(ctx, "Emitted DKG initialization signal",
		"mrenclave", dkgNetwork.Mrenclave,
		"round", dkgNetwork.Round,
		"total_validators", len(dkgNetwork.ActiveValSet),
		"height", dkgNetwork.StartBlock,
	)

	return nil
}

func (k *Keeper) emitBeginChallengePeriod(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent("dkg_begin_challenge_period",
			sdk.NewAttribute("round", fmt.Sprintf("%d", dkgNetwork.Round)),
			sdk.NewAttribute("mrenclave", string(dkgNetwork.Mrenclave)),
		),
	})

	log.Info(ctx, "Emitted BeginChallengePeriod event", "round", dkgNetwork.Round)
	return nil
}

func (k *Keeper) emitBeginDKGNetworkSet(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent("dkg_begin_network_set",
			sdk.NewAttribute("round", fmt.Sprintf("%d", dkgNetwork.Round)),
			sdk.NewAttribute("mrenclave", string(dkgNetwork.Mrenclave)),
			sdk.NewAttribute("total", fmt.Sprintf("%d", dkgNetwork.Total)),
			sdk.NewAttribute("threshold", fmt.Sprintf("%d", dkgNetwork.Threshold)),
		),
	})

	log.Info(ctx, "Emitted BeginDKGNetworkSet event", "round", dkgNetwork.Round)
	return nil
}

func (k *Keeper) emitBeginDKGFinalization(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent("dkg_begin_finalization",
			sdk.NewAttribute("round", fmt.Sprintf("%d", dkgNetwork.Round)),
			sdk.NewAttribute("mrenclave", string(dkgNetwork.Mrenclave)),
		),
	})

	log.Info(ctx, "Emitted BeginDKGFinalization event", "round", dkgNetwork.Round)
	return nil
}

func (k *Keeper) emitDKGCompleted(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent("dkg_completed",
			sdk.NewAttribute("round", fmt.Sprintf("%d", dkgNetwork.Round)),
			sdk.NewAttribute("mrenclave", string(dkgNetwork.Mrenclave)),
		),
	})

	log.Info(ctx, "Emitted DKGCompleted event", "round", dkgNetwork.Round)
	return nil
}
