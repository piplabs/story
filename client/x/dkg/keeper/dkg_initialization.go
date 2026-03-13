package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/client/server/utils"
	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

// GetActiveValidators returns the bonded validators' EVM addresses excluding jailed validators.
func (k *Keeper) GetActiveValidators(ctx context.Context) ([]string, error) {
	validators, err := k.stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all validators")
	}

	bondedValidators := make([]string, 0, len(validators))
	for _, val := range validators {
		if val.IsBonded() && !val.IsJailed() {
			evmOperatorAddress, err := utils.Bech32ValidatorAddressToEvmAddress(val.OperatorAddress)
			if err != nil {
				return nil, errors.Wrap(err, "failed to convert to evm address", "operator_addr", val.OperatorAddress)
			}

			bondedValidators = append(bondedValidators, strings.ToLower(evmOperatorAddress))
		}
	}

	return bondedValidators, nil
}

// InitiateDKGRound starts a new DKG round. If isUpgrade is true, the round is marked
// as an upgrade resharing round (IsResharing=true, IsUpgrade=true) and the current
// active round is NOT inactive.
func (k *Keeper) InitiateDKGRound(ctx context.Context, isUpgrade bool) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	activeValidators, err := k.GetActiveValidators(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get active validators")
	}

	roundNum := k.getNextRoundNumber(ctx)

	isResharing, err := k.shouldReshare(ctx)
	if err != nil {
		return err
	}

	dkgNetwork := types.DKGNetwork{
		Round:            roundNum,
		StartBlockHeight: sdkCtx.BlockHeight(),
		StartBlockHash:   sdkCtx.HeaderHash(),
		ActiveValSet:     activeValidators,
		Total:            0,
		Threshold:        0,
		Stage:            types.DKGStageRegistration,
		IsResharing:      isResharing,
		IsUpgrade:        isUpgrade,
	}

	if err := k.setDKGNetwork(ctx, &dkgNetwork); err != nil {
		return err
	}

	log.Info(ctx, "Initiated new DKG round",
		"round", roundNum,
		"start_block", sdkCtx.BlockHeight(),
		"is_upgrade", isUpgrade,
	)

	if err := k.emitBeginDKGInitialization(ctx, &dkgNetwork); err != nil {
		return errors.Wrap(err, "failed to emit begin dkg initialization event")
	}

	if k.isDKGSvcEnabled {
		// Pre-compute old code commitment while we still have SDK context.
		// The async goroutine uses context.Background() which cannot access
		// the Cosmos KV store.
		oldCC, _ := k.getOldCodeCommitment(ctx)

		asyncCtx, cancel := dkgAsyncContext()

		go func() {
			defer cancel()

			k.handleDKGRegistration(asyncCtx, &dkgNetwork, oldCC)
		}()
	}

	return nil
}

func (k *Keeper) shouldReshare(ctx context.Context) (bool, error) {
	activeNetwork, err := k.getLatestActiveDKGNetwork(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get latest active DKG network")
	}

	return activeNetwork != nil, nil
}
