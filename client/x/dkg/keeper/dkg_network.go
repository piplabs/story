package keeper

import (
	"context"
	"strconv"

	"cosmossdk.io/collections"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

// setDKGNetwork stores a DKG network in the store using round as the key.
// If this DKG network is the latest DKG network (per `isLatestDKGNetwork`), it updates the latest pointer.
func (k *Keeper) setDKGNetwork(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	key := strconv.FormatUint(uint64(dkgNetwork.Round), 10)
	if err := k.DKGNetworks.Set(ctx, key, *dkgNetwork); err != nil {
		return err
	}

	shouldUpdateLatest, err := k.isLatestDKGNetwork(ctx, dkgNetwork)
	if err != nil {
		return errors.Wrap(err, "failed to check if DKG network is latest")
	}

	if shouldUpdateLatest {
		if err := k.LatestDKGNetwork.Set(ctx, key); err != nil {
			return errors.Wrap(err, "failed to update latest DKG network pointer")
		}
	}

	return nil
}

// getDKGNetwork retrieves a DKG network by round.
func (k *Keeper) getDKGNetwork(ctx context.Context, round uint32) (*types.DKGNetwork, error) {
	key := strconv.FormatUint(uint64(round), 10)

	dkgNetwork, err := k.DKGNetworks.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errors.Wrap(err, "dkg network not found")
		}

		return nil, errors.Wrap(err, "failed to get dkg network")
	}

	return &dkgNetwork, nil
}

func (k *Keeper) getLatestDKGNetwork(ctx context.Context) (*types.DKGNetwork, error) {
	latestKey, err := k.LatestDKGNetwork.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest DKG network key")
	}

	dkgNetwork, err := k.DKGNetworks.Get(ctx, latestKey)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errors.Wrap(err, "dkg network not found")
		}

		return nil, errors.Wrap(err, "failed to get dkg network")
	}

	return &dkgNetwork, nil
}

// GetLatestDKGRound retrieves the latest DKG network.
func (k *Keeper) GetLatestDKGRound(ctx context.Context) (*types.DKGNetwork, error) {
	latestKey, err := k.LatestDKGNetwork.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// No DKG network set yet
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get latest DKG network key")
	}

	dkgNetwork, err := k.DKGNetworks.Get(ctx, latestKey)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			// This shouldn't happen (pointer exists but network state doesn't)... reset the pointer and return nil
			_ = k.LatestDKGNetwork.Remove(ctx)

			return nil, errors.Wrap(err, "latest DKG network not found")
		}

		return nil, errors.Wrap(err, "failed to get latest DKG network")
	}

	return &dkgNetwork, nil
}

// GetDKGNetworksByRound retrieves one or many DKG networks by a specified round number.
func (k *Keeper) GetDKGNetworksByRound(ctx context.Context, round uint32) ([]types.DKGNetwork, error) {
	var foundNetworks []types.DKGNetwork

	// Iterate through all DKG networks to find the one with the specified round
	err := k.DKGNetworks.Walk(ctx, nil, func(_ string, dkgNetwork types.DKGNetwork) (bool, error) {
		if dkgNetwork.Round == round {
			foundNetworks = append(foundNetworks, dkgNetwork)
		} else if dkgNetwork.Round > round {
			return true, nil // Stop iteration
		}

		return false, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate DKG networks")
	}

	return foundNetworks, nil
}

// getAllDKGNetworks retrieves all DKG networks.
func (k *Keeper) getAllDKGNetworks(ctx context.Context) ([]types.DKGNetwork, error) {
	var networks []types.DKGNetwork

	err := k.DKGNetworks.Walk(ctx, nil, func(_ string, dkgNetwork types.DKGNetwork) (bool, error) {
		networks = append(networks, dkgNetwork)

		return false, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate DKG networks")
	}

	return networks, nil
}

// DeleteDKGNetwork removes a DKG network from the store.
func (k *Keeper) DeleteDKGNetwork(ctx context.Context, round uint32) error {
	key := strconv.FormatUint(uint64(round), 10)
	return k.DKGNetworks.Remove(ctx, key)
}

// isLatestDKGNetwork determines if the given DKG network should be the new latest
// Returns true if:
// - No current latest exists (collection not found), OR
// - This network has a higher round number, OR
// - Same round but newer start block (newer TEE binary for same round).
func (k *Keeper) isLatestDKGNetwork(ctx context.Context, dkgNetwork *types.DKGNetwork) (bool, error) {
	currentLatest, err := k.GetLatestDKGRound(ctx)
	if err != nil {
		// If no latest DKG network exists yet, return true
		if errors.Is(err, collections.ErrNotFound) {
			return true, nil
		}

		return false, err
	}

	if currentLatest == nil {
		return true, nil
	}

	if dkgNetwork.Round > currentLatest.Round {
		return true, nil
	}

	if dkgNetwork.Round == currentLatest.Round && dkgNetwork.StartBlockHeight > currentLatest.StartBlockHeight {
		return true, nil
	}

	return false, nil
}

func (k *Keeper) getNextRoundNumber(ctx context.Context) uint32 {
	latestNetwork, err := k.GetLatestDKGRound(ctx)
	if err != nil || latestNetwork == nil {
		return 1 // Start with round 1 if no previous rounds exist
	}

	return latestNetwork.Round + 1
}

func (k *Keeper) setLatestActiveRound(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	key := strconv.FormatUint(uint64(dkgNetwork.Round), 10)
	if err := k.LatestActiveRound.Set(ctx, key); err != nil {
		return errors.Wrap(err, "failed to update latest active round of DKG network pointer")
	}

	return nil
}

// GetLatestActiveRound retrieves the latest active DKG round.
// Returns nil, nil if no active round has been set yet.
func (k *Keeper) GetLatestActiveRound(ctx context.Context) (*types.DKGNetwork, error) {
	return k.getLatestActiveDKGNetwork(ctx)
}

func (k *Keeper) getLatestActiveDKGNetwork(ctx context.Context) (*types.DKGNetwork, error) {
	key, err := k.LatestActiveRound.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get latest active round of DKG network key")
	}

	dkgNetwork, err := k.DKGNetworks.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get latest active round of dkg network")
	}

	return &dkgNetwork, nil
}

// findPrevActiveRound returns the round number of the most recent successful
// (Active or Ended) DKG round strictly before beforeRound. Returns 0 if none found.
// DKG rounds are small in number so a full Walk is acceptable here.
func (k *Keeper) findPrevActiveRound(ctx context.Context, beforeRound uint32) (uint32, error) {
	var prev uint32
	err := k.DKGNetworks.Walk(ctx, nil, func(_ string, network types.DKGNetwork) (bool, error) {
		if network.Round >= beforeRound {
			return false, nil
		}
		if (network.Stage == types.DKGStageActive || network.Stage == types.DKGStageEnded) && network.Round > prev {
			prev = network.Round
		}
		return false, nil
	})
	return prev, errors.Wrap(err, "find previous active round")
}

// endPreviousActiveRound finds the previous active round and sets its stage to Ended.
// This prevents the previous round from continuing stage transitions after a new round
// has been finalized and becomes the active round.
func (k *Keeper) endPreviousActiveRound(ctx context.Context, currentRound uint32) error {
	prevActive, err := k.getLatestActiveDKGNetwork(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get previous active round")
	}

	if prevActive == nil || prevActive.Round == currentRound {
		return nil
	}

	if prevActive.Stage == types.DKGStageActive {
		prevActive.Stage = types.DKGStageEnded
		if err := k.setDKGNetwork(ctx, prevActive); err != nil {
			return errors.Wrap(err, "failed to set previous active round to ended")
		}
	}

	return nil
}
