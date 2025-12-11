package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func (k *Keeper) SetDKGNetwork(ctx context.Context, latestRound *types.DKGNetwork) error {
	params, err := k.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get dkg params")
	}

	verifiedCount, err := k.countDKGRegistrationsByStatus(ctx, latestRound.Mrenclave, latestRound.Round, types.DKGRegStatusVerified)
	if err != nil {
		return errors.Wrap(err, "failed to get verified DKG validators count")
	}

	if verifiedCount < params.MinCommitteeSize {
		log.Info(ctx, "The number of DKG registrations in Verified status is smaller than the minimum committee size. Skipping current round.", "current", latestRound.Round, "next", latestRound.Round+1)

		return k.SkipToNextRound(ctx, latestRound)
	}

	latestRound.Total = verifiedCount
	latestRound.Threshold = k.calculateThreshold(verifiedCount)

	if err := k.setDKGNetwork(ctx, latestRound); err != nil {
		return errors.Wrap(err, "failed to update total and threshold of the DKG network", "mrenclave", hex.EncodeToString(latestRound.Mrenclave), "round", latestRound.Round)
	}

	if err := k.emitBeginDKGNetworkSet(ctx, latestRound); err != nil {
		return errors.Wrap(err, "failed to emit begin DKG network set event")
	}

	if k.isDKGSvcEnabled {
		go k.handleDKGNetworkSet(ctx, latestRound)
	}

	return nil
}
