package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/cast"
	"github.com/piplabs/story/lib/errors"
)

// registrationListKey returns the key for the registration list: mrenclave_round
func registrationListKey(mrenclave [32]byte, round uint32) string {
	return fmt.Sprintf("%s_%d", hex.EncodeToString(mrenclave[:]), round)
}

// setDKGRegistration adds or updates a DKG registration in the list for the given mrenclave and round.
func (k *Keeper) setDKGRegistration(ctx context.Context, mrenclave [32]byte, validatorAddr common.Address, dkgReg *types.DKGRegistration) error {
	key := registrationListKey(mrenclave, dkgReg.Round)

	regList, err := k.DKGRegistrations.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			regList = types.DKGRegistrationList{Registrations: []types.DKGRegistration{}}
		} else {
			return errors.Wrap(err, "failed to get dkg registration list")
		}
	}

	// Check if registration already exists (update) or is new (append)
	validatorAddrLower := strings.ToLower(validatorAddr.Hex())
	found := false
	for i, reg := range regList.Registrations {
		if strings.ToLower(reg.MsgSender) == validatorAddrLower {
			regList.Registrations[i] = *dkgReg
			found = true
			break
		}
	}

	if !found {
		regList.Registrations = append(regList.Registrations, *dkgReg)
	}

	if err := k.DKGRegistrations.Set(ctx, key, regList); err != nil {
		return errors.Wrap(err, "failed to set dkg registration list")
	}

	return nil
}

// getDKGRegistration retrieves a DKG registration by mrenclave, round, and validator address.
func (k *Keeper) getDKGRegistration(ctx context.Context, mrenclave [32]byte, round uint32, validatorAddr common.Address) (*types.DKGRegistration, error) {
	key := registrationListKey(mrenclave, round)

	regList, err := k.DKGRegistrations.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errors.New("dkg registration not found")
		}
		return nil, errors.Wrap(err, "failed to get dkg registration list")
	}

	validatorAddrLower := strings.ToLower(validatorAddr.Hex())
	for _, reg := range regList.Registrations {
		if strings.ToLower(reg.MsgSender) == validatorAddrLower {
			return &reg, nil
		}
	}

	return nil, errors.New("dkg registration not found")
}

// getDKGRegistrationIndex returns the 1-based PID (party index) for a validator.
func (k *Keeper) getDKGRegistrationIndex(ctx context.Context, mrenclave [32]byte, round uint32, msgSender common.Address) (uint32, error) {
	reg, err := k.getDKGRegistration(ctx, mrenclave, round, msgSender)
	if err != nil {
		return 0, err
	}

	return reg.Index, nil
}

// getNextDKGRegistrationIndex gets the next DKG registration index for a specific mrenclave and round.
func (k *Keeper) getNextDKGRegistrationIndex(ctx context.Context, mrenclave [32]byte, round uint32) (int, error) {
	registrations, err := k.getDKGRegistrationsByRound(ctx, mrenclave, round)
	if err != nil {
		return 0, err
	}

	return len(registrations) + 1, nil
}

// getDKGRegistrationsByRound retrieves all DKG registrations for a specific mrenclave and round.
func (k *Keeper) getDKGRegistrationsByRound(ctx context.Context, mrenclave [32]byte, round uint32) ([]types.DKGRegistration, error) {
	key := registrationListKey(mrenclave, round)

	regList, err := k.DKGRegistrations.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return []types.DKGRegistration{}, nil
		}
		return nil, errors.Wrap(err, "failed to get dkg registration list")
	}

	return regList.Registrations, nil
}

// updateDKGRegistrationStatus updates the status of a specific DKG registration.
func (k *Keeper) updateDKGRegistrationStatus(ctx context.Context, mrenclave [32]byte, round uint32, validatorAddr common.Address, status types.DKGRegStatus) error {
	dkgReg, err := k.getDKGRegistration(ctx, mrenclave, round, validatorAddr)
	if err != nil {
		return err
	}

	dkgReg.Status = status

	return k.setDKGRegistration(ctx, mrenclave, validatorAddr, dkgReg)
}

// getDKGRegistrationsByStatus retrieves all DKG registrations with a specific status for a given mrenclave and round.
func (k *Keeper) getDKGRegistrationsByStatus(ctx context.Context, mrenclave [32]byte, round uint32, status types.DKGRegStatus) ([]types.DKGRegistration, error) {
	allRegs, err := k.getDKGRegistrationsByRound(ctx, mrenclave, round)
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
func (k *Keeper) countDKGRegistrationsByStatus(ctx context.Context, mrenclave []byte, round uint32, status types.DKGRegStatus) (uint32, error) {
	mrenclave32, err := cast.ToBytes32(mrenclave)
	if err != nil {
		return 0, errors.Wrap(err, "failed to cast to bytes32")
	}

	regs, err := k.getDKGRegistrationsByStatus(ctx, mrenclave32, round, status)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get registrations by status")
	}

	return uint32(len(regs)), nil
}
