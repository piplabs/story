package keeper

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/collections"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

func dkgRegistrationKey(round uint32, validatorAddr common.Address) string {
	return fmt.Sprintf("%d_%s", round, strings.ToLower(validatorAddr.Hex()))
}

// setDKGRegistration stores a DKG registration in the store using round_address as the key.
// The address is lowercased to match the format used in DKGNetwork.ActiveValSet and story-kernel queries.
func (k *Keeper) setDKGRegistration(ctx context.Context, validatorAddr common.Address, dkgReg *types.DKGRegistration) error {
	key := dkgRegistrationKey(dkgReg.Round, validatorAddr)
	if err := k.DKGRegistrations.Set(ctx, key, *dkgReg); err != nil {
		return errors.Wrap(err, "failed to set dkg registration")
	}

	return nil
}

// getDKGRegistration retrieves a DKG registration by round and validator address.
func (k *Keeper) getDKGRegistration(ctx context.Context, round uint32, validatorAddr common.Address) (*types.DKGRegistration, error) {
	key := dkgRegistrationKey(round, validatorAddr)

	dkgReg, err := k.DKGRegistrations.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errors.Wrap(err, "dkg registration not found")
		}

		return nil, errors.Wrap(err, "failed to get dkg registration")
	}

	return &dkgReg, nil
}

// hasDKGRegistration checks if a DKG registration exists for a given round and validator address.
func (k *Keeper) hasDKGRegistration(ctx context.Context, round uint32, validatorAddr common.Address) (bool, error) {
	key := dkgRegistrationKey(round, validatorAddr)

	exists, err := k.DKGRegistrations.Has(ctx, key)
	if err != nil {
		return false, errors.Wrap(err, "failed to check dkg registration existence")
	}

	return exists, nil
}

// getNextDKGRegistrationIndex gets the next DKG registration index for a specific round.
// The returned index is 1-based: the first participant gets index 1, the second gets index 2, etc.
// This 1-based convention is used throughout the on-chain DKG registration and deal/response
// routing logic. When interfacing with kyber's DKG library, which uses 0-based PIDs
// (participant IDs), the caller must convert by subtracting 1 (i.e., kyberPID = index - 1).
func (k *Keeper) getNextDKGRegistrationIndex(ctx context.Context, round uint32) (int, error) {
	registrations, err := k.getDKGRegistrationsByRound(ctx, round)
	if err != nil {
		return 0, err
	}

	return len(registrations) + 1, nil
}

// getDKGRegistrationsByRound retrieves all DKG registrations for a specific round.
// Uses a prefix range to iterate only keys matching the round, avoiding a full-table scan.
func (k *Keeper) getDKGRegistrationsByRound(ctx context.Context, round uint32) ([]types.DKGRegistration, error) {
	var registrations []types.DKGRegistration

	prefix := fmt.Sprintf("%d_", round)
	rng := (&collections.Range[string]{}).Prefix(prefix)

	err := k.DKGRegistrations.Walk(ctx, rng, func(_ string, reg types.DKGRegistration) (bool, error) {
		registrations = append(registrations, reg)

		return false, nil // Continue iteration
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate dkg registrations")
	}

	return registrations, nil
}

// finalizeDKGRegistration updates the status of a specific DKG registration and sets the pub key share.
func (k *Keeper) finalizeDKGRegistration(ctx context.Context, round uint32, validatorAddr common.Address, pubKeyShare []byte) error {
	dkgReg, err := k.getDKGRegistration(ctx, round, validatorAddr)
	if err != nil {
		return err
	}

	dkgReg.PubKeyShare = pubKeyShare
	dkgReg.Status = types.DKGRegStatusFinalized

	return k.setDKGRegistration(ctx, validatorAddr, dkgReg)
}

// getDKGRegistrationsByStatus retrieves all DKG registrations with a specific status for a given round.
func (k *Keeper) getDKGRegistrationsByStatus(ctx context.Context, round uint32, status types.DKGRegStatus) ([]types.DKGRegistration, error) {
	allRegs, err := k.getDKGRegistrationsByRound(ctx, round)
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

// HasFinalizedRegistration checks if a validator has a finalized DKG registration for a given round.
func (k *Keeper) HasFinalizedRegistration(ctx context.Context, round uint32, validatorAddr common.Address) (bool, error) {
	reg, err := k.getDKGRegistration(ctx, round, validatorAddr)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	return reg.Status == types.DKGRegStatusFinalized, nil
}

// countDKGRegistrationsByStatus returns the count of DKG registrations in the status.
func (k *Keeper) countDKGRegistrationsByStatus(ctx context.Context, round uint32, status types.DKGRegStatus) (uint32, error) {
	// Get registrations with status of registration
	regs, err := k.getDKGRegistrationsByStatus(ctx, round, status)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get verified registrations")
	}

	return uint32(len(regs)), nil
}
