package keeper

import (
	"context"
	"encoding/hex"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/cast"
	"github.com/piplabs/story/lib/errors"
)

// getVerifiedDKGValidators returns the count of verified DKG validators (those who are participating and not invalidated).
func (k *Keeper) getVerifiedDKGValidators(ctx context.Context, mrenclave []byte, round uint32) (uint32, error) {
	mrenclave32, err := cast.ToBytes32(mrenclave)
	if err != nil {
		return 0, errors.Wrap(err, "failed to cast to bytes32")
	}

	// Get registrations with status VERIFIED
	verifiedRegs, err := k.getDKGRegistrationsByStatus(ctx, mrenclave32, round, types.DKGRegStatusVerified)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get verified registrations")
	}

	total := uint32(len(verifiedRegs))

	return total, nil
}

func (k *Keeper) updateDKGNetworkTotalAndThreshold(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	verifiedCount, err := k.getVerifiedDKGValidators(ctx, dkgNetwork.Mrenclave, dkgNetwork.Round)
	if err != nil {
		return errors.Wrap(err, "failed to get verified DKG validators count")
	}

	dkgNetwork.Total = verifiedCount
	dkgNetwork.Threshold = k.calculateThreshold(verifiedCount)

	return k.setDKGNetwork(ctx, dkgNetwork)
}

func (k *Keeper) SetDKGNetwork(ctx context.Context, latestRound *types.DKGNetwork) error {
	// TODO: check if there's enough number of registrations to set the network (and start dealing).
	// Use a DKG module parameter (minDKGMemberAmount). If the amount is less, we need to restart the round.

	if err := k.updateDKGNetworkTotalAndThreshold(ctx, latestRound); err != nil {
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
