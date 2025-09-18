package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

// SetDKGRegistration stores a DKG registration in the store using mrenclave_round_index as the key.
func (k *Keeper) SetDKGRegistration(ctx context.Context, mrenclave []byte, dkgReg *types.DKGRegistration) error {
	key := fmt.Sprintf("%s_%d_%d", string(mrenclave), dkgReg.Round, dkgReg.Index)
	if err := k.DKGRegistrations.Set(ctx, key, *dkgReg); err != nil {
		return errors.Wrap(err, "failed to set dkg registration")
	}

	return nil
}

// GetDKGRegistration retrieves a DKG registration by unique ID, round, and index.
func (k *Keeper) GetDKGRegistration(ctx context.Context, mrenclave []byte, round, index uint32) (*types.DKGRegistration, error) {
	key := fmt.Sprintf("%s_%d_%d", string(mrenclave), round, index)
	dkgReg, err := k.DKGRegistrations.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errors.Wrap(err, "dkg registration not found")
		}

		return nil, errors.Wrap(err, "failed to get dkg registration")
	}

	return &dkgReg, nil
}

// GetDKGRegistrationsByRound retrieves all DKG registrations for a specific unique ID and round.
func (k *Keeper) GetDKGRegistrationsByRound(ctx context.Context, mrenclave []byte, round uint32) ([]types.DKGRegistration, error) {
	var registrations []types.DKGRegistration
	prefix := fmt.Sprintf("%s_%d_", string(mrenclave), round)

	err := k.DKGRegistrations.Walk(ctx, nil, func(key string, reg types.DKGRegistration) (bool, error) {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			registrations = append(registrations, reg)
		}

		return false, nil // Continue iteration
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate dkg registrations")
	}

	return registrations, nil
}

// GetAllDKGRegistrations retrieves all DKG registrations.
func (k *Keeper) GetAllDKGRegistrations(ctx context.Context) ([]types.DKGRegistration, error) {
	var registrations []types.DKGRegistration

	err := k.DKGRegistrations.Walk(ctx, nil, func(_ string, reg types.DKGRegistration) (bool, error) {
		registrations = append(registrations, reg)

		return false, nil // Continue iteration
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate dkg registrations")
	}

	return registrations, nil
}

// DeleteDKGRegistration removes a DKG registration from the store.
func (k *Keeper) DeleteDKGRegistration(ctx context.Context, mrenclave []byte, round, index uint32) error {
	key := fmt.Sprintf("%s_%d_%d", string(mrenclave), round, index)
	if err := k.DKGRegistrations.Remove(ctx, key); err != nil {
		return errors.Wrap(err, "failed to delete dkg registration")
	}

	return nil
}

// UpdateDKGRegistrationStatus updates the status of a specific DKG registration.
func (k *Keeper) UpdateDKGRegistrationStatus(ctx context.Context, mrenclave []byte, round, index uint32, status types.DKGRegStatus) error {
	dkgReg, err := k.GetDKGRegistration(ctx, mrenclave, round, index)
	if err != nil {
		return err
	}

	dkgReg.Status = status

	return k.SetDKGRegistration(ctx, mrenclave, dkgReg)
}

// GetDKGRegistrationsByStatus retrieves all DKG registrations with a specific status for a given unique ID and round.
func (k *Keeper) GetDKGRegistrationsByStatus(ctx context.Context, mrenclave []byte, round uint32, status types.DKGRegStatus) ([]types.DKGRegistration, error) {
	allRegs, err := k.GetDKGRegistrationsByRound(ctx, mrenclave, round)
	if err != nil {
		return nil, err
	}

	var filteredRegs []types.DKGRegistration
	for _, reg := range allRegs {
		if reg.Status == status {
			filteredRegs = append(filteredRegs, reg)
		}
	}

	return filteredRegs, nil
}
