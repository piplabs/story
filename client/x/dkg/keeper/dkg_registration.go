package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/piplabs/story/lib/cast"

	"cosmossdk.io/collections"

	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

// setDKGRegistration stores a DKG registration in the store using codeCommitment_round_index as the key.
func (k *Keeper) setDKGRegistration(ctx context.Context, codeCommitment [32]byte, validatorAddr common.Address, dkgReg *types.DKGRegistration) error {
	key := fmt.Sprintf("%s_%d_%s", hex.EncodeToString(codeCommitment[:]), dkgReg.Round, validatorAddr.Hex())
	if err := k.DKGRegistrations.Set(ctx, key, *dkgReg); err != nil {
		return errors.Wrap(err, "failed to set dkg registration")
	}

	return nil
}

// getDKGRegistration retrieves a DKG registration by code commitment, round, and index.
func (k *Keeper) getDKGRegistration(ctx context.Context, codeCommitment [32]byte, round uint32, validatorAddr common.Address) (*types.DKGRegistration, error) {
	key := fmt.Sprintf("%s_%d_%s", hex.EncodeToString(codeCommitment[:]), round, validatorAddr.Hex())
	dkgReg, err := k.DKGRegistrations.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errors.Wrap(err, "dkg registration not found")
		}

		return nil, errors.Wrap(err, "failed to get dkg registration")
	}

	return &dkgReg, nil
}

// getDKGRegistrationIndex gets the index of a specific DKG registration by code commitment, round, and msg sender.
//
// TODO: optimize since `getDKGRegistrationsByRound` walks all registrations.
func (k *Keeper) getDKGRegistrationIndex(ctx context.Context, codeCommitment [32]byte, round uint32, msgSender string) (uint32, error) {
	registrations, err := k.getDKGRegistrationsByRound(ctx, codeCommitment, round)
	if err != nil {
		return 0, err
	}

	for _, registration := range registrations {
		if registration.ValidatorAddr == msgSender {
			return registration.Index, nil
		}
	}

	return 0, errors.New("dkg registration not found")
}

// getNextDKGRegistrationIndex gets the next DKG registration index for a specific code commitment and round.
func (k *Keeper) getNextDKGRegistrationIndex(ctx context.Context, codeCommitment [32]byte, round uint32) (int, error) {
	registrations, err := k.getDKGRegistrationsByRound(ctx, codeCommitment, round)
	if err != nil {
		return 0, err
	}

	return len(registrations) + 1, nil
}

// getDKGRegistrationsByRound retrieves all DKG registrations for a specific code commitment and round.
func (k *Keeper) getDKGRegistrationsByRound(ctx context.Context, codeCommitment [32]byte, round uint32) ([]types.DKGRegistration, error) {
	var registrations []types.DKGRegistration
	prefix := fmt.Sprintf("%s_%d_", hex.EncodeToString(codeCommitment[:]), round)

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

// updateDKGRegistrationStatus updates the status of a specific DKG registration.
func (k *Keeper) updateDKGRegistrationStatus(ctx context.Context, codeCommitment [32]byte, round uint32, validatorAddr common.Address, status types.DKGRegStatus) error {
	dkgReg, err := k.getDKGRegistration(ctx, codeCommitment, round, validatorAddr)
	if err != nil {
		return err
	}

	dkgReg.Status = status

	return k.setDKGRegistration(ctx, codeCommitment, validatorAddr, dkgReg)
}

// getDKGRegistrationsByStatus retrieves all DKG registrations with a specific status for a given code commitment and round.
func (k *Keeper) getDKGRegistrationsByStatus(ctx context.Context, codeCommitment [32]byte, round uint32, status types.DKGRegStatus) ([]types.DKGRegistration, error) {
	allRegs, err := k.getDKGRegistrationsByRound(ctx, codeCommitment, round)
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

// countDKGRegistrationsByStatus returns the count of DKG registrations in the status
func (k *Keeper) countDKGRegistrationsByStatus(ctx context.Context, codeCommitment []byte, round uint32, status types.DKGRegStatus) (uint32, error) {
	codeCommitment32, err := cast.ToBytes32(codeCommitment)
	if err != nil {
		return 0, errors.Wrap(err, "failed to cast to bytes32")
	}

	// Get registrations with status VERIFIED
	regs, err := k.getDKGRegistrationsByStatus(ctx, codeCommitment32, round, status)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get verified registrations")
	}

	return uint32(len(regs)), nil
}
