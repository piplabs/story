package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

// SetDKGNetwork stores a DKG network in the store using the mrenclave as the key
// If this DKG network is the latest DKG network (per `isLatestDKGNetwork`), it updates the latest pointer.
func (k *Keeper) SetDKGNetwork(ctx context.Context, dkgNetwork *types.DKGNetwork) error {
	key := fmt.Sprintf("%s_%d", string(dkgNetwork.Mrenclave), dkgNetwork.Round)
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

// GetDKGNetworkByKey retrieves a DKG network by mrenclave.
func (k *Keeper) getDKGNetwork(ctx context.Context, mrenclave []byte, round uint32) (*types.DKGNetwork, error) {
	key := fmt.Sprintf("%s_%d", string(mrenclave), round)
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
			return nil, errors.Wrap(err, "latest DKG network not set")
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
func (k *Keeper) DeleteDKGNetwork(ctx context.Context, mrenclave []byte, round uint32) error {
	key := fmt.Sprintf("%s_%d", string(mrenclave), round)
	return k.DKGNetworks.Remove(ctx, key)
}

// isLatestDKGNetwork determines if the given DKG network should be the new latest
// Returns true if:
// - No current latest exists, OR
// - This network has a higher round number, OR
// - Same round but newer start block (newer TEE binary for same round).
func (k *Keeper) isLatestDKGNetwork(ctx context.Context, dkgNetwork *types.DKGNetwork) (bool, error) {
	currentLatest, err := k.GetLatestDKGRound(ctx)
	if err != nil {
		return false, err
	}
	if currentLatest == nil {
		return true, nil
	}
	if dkgNetwork.Round > currentLatest.Round {
		return true, nil
	}
	if dkgNetwork.Round == currentLatest.Round && dkgNetwork.StartBlock > currentLatest.StartBlock {
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

//nolint:unused // ignore unused error
func (*Keeper) calculateThreshold(total uint32) uint32 {
	threshold := (total * 2) / 3
	if threshold*3 < total*2 {
		threshold++
	}

	return threshold + 1 // 2/3 + 1 threshold
}
