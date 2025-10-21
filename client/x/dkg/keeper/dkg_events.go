package keeper

import (
	"context"
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

//
// events for stages (ie. events that drive the stage transitions)
//

func (*Keeper) emitBeginDKGInitialization(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginInitialization{
		Mrenclave:        dkgNetwork.Mrenclave,
		Round:            dkgNetwork.Round,
		ActiveValidators: dkgNetwork.ActiveValSet,
		StartBlock:       uint32(dkgNetwork.StartBlock),
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_begin_initialization event")
	}

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
		return errors.Wrap(err, "failed to emit dkg_begin_network_set event")
	}

	return nil
}

func (*Keeper) emitBeginDKGDealing(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginDealing{
		Mrenclave: dkgNetwork.Mrenclave,
		Round:     dkgNetwork.Round,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_begin_dealing event")
	}

	log.Info(ctx, "Emitted BeginDKGDealing event", "round", dkgNetwork.Round)

	return nil
}

func (*Keeper) emitBeginProcessDeals(ctx context.Context, dkgNetwork *types.DKGNetwork, deals []types.Deal) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginProcessDeals{
		Mrenclave: dkgNetwork.Mrenclave,
		Round:     dkgNetwork.Round,
		Deals:     deals,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_begin_process_deals event")
	}

	log.Info(ctx, "Emitted BeginProcessDeals event", "round", dkgNetwork.Round, "mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave), "num_deals", len(deals))

	return nil
}

func (*Keeper) emitBeginProcessResponses(ctx context.Context, dkgNetwork *types.DKGNetwork, responses []types.Response) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginProcessResponses{
		Mrenclave: dkgNetwork.Mrenclave,
		Round:     dkgNetwork.Round,
		Responses: responses,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_begin_process_responses event")
	}

	log.Info(ctx, "Emitted BeginProcessResponses event", "round", dkgNetwork.Round, "mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave), "num_responses", len(responses))

	return nil
}

func (*Keeper) emitBeginDKGFinalization(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginFinalization{
		Mrenclave: dkgNetwork.Mrenclave,
		Round:     dkgNetwork.Round,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_begin_finalization event")
	}

	log.Info(ctx, "Emitted BeginDKGFinalization event", "round", dkgNetwork.Round, "mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave))

	return nil
}

func (*Keeper) emitDKGFinalized(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventDKGFinalized{
		Mrenclave: dkgNetwork.Mrenclave,
		Round:     dkgNetwork.Round,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_finalized_event")
	}

	log.Info(ctx, "Emitted DKGFinalized event", "round", dkgNetwork.Round, "mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave))

	return nil
}
