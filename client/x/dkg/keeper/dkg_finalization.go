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

func (k *Keeper) FinalizeDKGRound(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	// TODO: check if enough number of validators submit the finalizeDKG tx
	if err := k.emitDKGFinalized(ctx, dkgNetwork); err != nil {
		return errors.Wrap(err, "failed to emit DKG finalized event")
	}

	if k.isDKGSvcEnabled {
		go k.handleDKGComplete(ctx, dkgNetwork)
	}

	log.Info(ctx, "DKG network setup completed", "round", dkgNetwork.Round, "mrenclave", hex.EncodeToString(dkgNetwork.Mrenclave))

	return nil
}
