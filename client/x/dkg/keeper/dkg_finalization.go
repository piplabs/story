package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func (k *Keeper) BeginFinalization(ctx context.Context, latestRound *types.DKGNetwork) error {
	if err := k.emitBeginDKGFinalization(ctx, latestRound); err != nil {
		return errors.Wrap(err, "failed to emit begin DKG finalization event")
	}

	if k.isDKGSvcEnabled {
		go k.handleDKGFinalization(ctx, latestRound)
	}

	return nil
}

func (k *Keeper) FinalizeDKGRound(ctx context.Context, latestRound *types.DKGNetwork) error {
	finalizedCount, err := k.countDKGRegistrationsByStatus(ctx, latestRound.Mrenclave, latestRound.Round, types.DKGRegStatusFinalized)
	if err != nil {
		return errors.Wrap(err, "failed to fetch DKG registrations in Finalized status")
	}

	if finalizedCount < latestRound.Threshold {
		log.Info(ctx, "The number of DKG registrations in Finalized status is smaller than the threshold. Skipping current round.", "current", latestRound.Round, "next", latestRound.Round+1)

		return k.SkipToNextRound(ctx, latestRound)
	}

	if err := k.emitDKGFinalized(ctx, latestRound); err != nil {
		return errors.Wrap(err, "failed to emit DKG finalized event")
	}

	if k.isDKGSvcEnabled {
		go k.handleDKGComplete(ctx, latestRound)
	}

	log.Info(ctx, "DKG network setup completed", "round", latestRound.Round, "mrenclave", hex.EncodeToString(latestRound.Mrenclave))

	return nil
}
