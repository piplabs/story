package keeper

import (
	"context"

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

func (*Keeper) emitBeginChallengePeriod(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginChallengePeriod{
		Mrenclave: dkgNetwork.Mrenclave,
		Round:     dkgNetwork.Round,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_begin_challenge_period event")
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
		return errors.Wrap(err, "failed to emit dkg_begin_network_set event")
	}

	log.Info(ctx, "Emitted BeginDKGNetworkSet event", "round", dkgNetwork.Round)

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

func (*Keeper) emitBeginDKGFinalization(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventBeginFinalization{
		Mrenclave: dkgNetwork.Mrenclave,
		Round:     dkgNetwork.Round,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_begin_finalization event")
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
		return errors.Wrap(err, "failed to emit dkg_finalized_event")
	}

	log.Info(ctx, "Emitted DKGFinalized event", "round", dkgNetwork.Round)

	return nil
}

//
// events for non-stages (ie. events that update states inside a stage)
//

func (*Keeper) emitDKGRegistrationInitialized(ctx context.Context, mrenclave []byte, round uint32, index uint32, dkgPubKey []byte, commPubKey []byte, remoteReport []byte) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventDKGRegistrationInitialized{
		Mrenclave:  mrenclave,
		Round:      round,
		Index:      index,
		DkgPubKey:  dkgPubKey,
		CommPubKey: commPubKey,
		RawQuote:   remoteReport,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_registration_initialized event")
	}

	log.Info(ctx, "Emitted DKGRegistrationInitialized event", "round", round, "index", index)

	return nil
}

func (*Keeper) emitDKGRegistrationCommitmentsUpdated(ctx context.Context, mrenclave []byte, round uint32, total uint32, threshold uint32, index uint32, commitments []byte, signature []byte) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := sdkCtx.EventManager().EmitTypedEvent(&types.EventDKGRegistrationCommitmentsUpdated{
		Mrenclave:   mrenclave,
		Round:       round,
		Total:       total,
		Threshold:   threshold,
		Index:       index,
		Commitments: commitments,
		Signature:   signature,
	})
	if err != nil {
		return errors.Wrap(err, "failed to emit dkg_registration_commitments_updated event")
	}

	log.Info(ctx, "Emitted DKGRegistrationCommitmentsUpdated event", "round", round, "index", index)

	return nil
}
