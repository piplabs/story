package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func (*Keeper) emitBeginDKGInitialization(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginInitialization{
		Mrenclave:        dkgNetwork.Mrenclave,
		Round:            dkgNetwork.Round,
		ActiveValidators: dkgNetwork.ActiveValSet,
		Total:            uint32(len(dkgNetwork.ActiveValSet)),
		StartBlock:       uint32(dkgNetwork.StartBlock),
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit DKG initialization event")
	}

	log.Info(ctx, "Emitted DKG initialization signal",
		"mrenclave", dkgNetwork.Mrenclave,
		"round", dkgNetwork.Round,
		"total_validators", len(dkgNetwork.ActiveValSet),
		"height", dkgNetwork.StartBlock,
	)

	return nil
}

func (*Keeper) emitBeginChallengePeriod(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginChallengePeriod{
		Mrenclave: dkgNetwork.Mrenclave,
		Round:     dkgNetwork.Round,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit BeginChallengePeriod event")
	}

	log.Info(ctx, "Emitted BeginChallengePeriod event", "round", dkgNetwork.Round)

	return nil
}

func (*Keeper) emitBeginDKGNetworkSet(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginNetworkSet{
		Mrenclave: dkgNetwork.Mrenclave,
		Round:     dkgNetwork.Round,
		Total:     dkgNetwork.Total,
		Threshold: dkgNetwork.Threshold,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit BeginDKGNetworkSet event")
	}

	log.Info(ctx, "Emitted BeginDKGNetworkSet event", "round", dkgNetwork.Round)

	return nil
}

func (*Keeper) emitBeginDKGFinalization(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginFinalization{
		Mrenclave: dkgNetwork.Mrenclave,
		Round:     dkgNetwork.Round,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit BeginDKGFinalization event")
	}

	log.Info(ctx, "Emitted BeginDKGFinalization event", "round", dkgNetwork.Round)

	return nil
}

func (*Keeper) emitDKGFinalized(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventDKGFinalized{
		Mrenclave: dkgNetwork.Mrenclave,
		Round:     dkgNetwork.Round,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit DKGFinalized event")
	}

	log.Info(ctx, "Emitted DKGFinalized event", "round", dkgNetwork.Round)

	return nil
}
