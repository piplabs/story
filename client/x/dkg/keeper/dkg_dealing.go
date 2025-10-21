package keeper

import (
	"context"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

func (k *Keeper) BeginDealing(ctx context.Context, latestRound *types.DKGNetwork) error {
	if err := k.emitBeginDKGDealing(ctx, latestRound); err != nil {
		return errors.Wrap(err, "failed to emit begin DKG dealing event")
	}

	if k.isDKGSvcEnabled {
		go k.handleDKGDealing(ctx, latestRound)
	}

	return nil
}

func (k *Keeper) ProcessDeals(ctx context.Context, latestRound *types.DKGNetwork, deals []types.Deal) error {
	if err := k.emitBeginProcessDeals(ctx, latestRound, deals); err != nil {
		return errors.Wrap(err, "failed to emit begin process deals event")
	}

	if k.isDKGSvcEnabled {
		go k.handleDKGProcessDeals(ctx, latestRound, deals)
	}

	return nil
}

func (k *Keeper) ProcessResponses(ctx context.Context, latestRound *types.DKGNetwork, responses []types.Response) error {
	if err := k.emitBeginProcessResponses(ctx, latestRound, responses); err != nil {
		return errors.Wrap(err, "failed to emit begin process responses event")
	}

	if k.isDKGSvcEnabled {
		go k.handleDKGProcessResponses(ctx, latestRound, responses)
	}

	return nil
}
