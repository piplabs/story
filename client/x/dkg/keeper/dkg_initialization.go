package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

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
	// Flush stale deals/responses/justifications from previous rounds to prevent
	// them from contaminating the new round's vote extensions.
	k.FlushAllQueues()

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := k.distributeCDRFee(ctx); err != nil {
		return errors.Wrap(err, "failed to distribute CDR fee pool")
	}

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
	roundsTotal.WithLabelValues("initiated").Inc()

	if err := k.emitBeginDKGInitialization(ctx, &dkgNetwork); err != nil {
		return errors.Wrap(err, "failed to emit begin dkg initialization event")
	}

	if k.isDKGSvcEnabled {
		// Use a gasless context for KV reads inside isDKGSvcEnabled so that
		// DKG-enabled and DKG-disabled nodes produce identical GasUsed.
		gaslessCtx := gaslessSDKContext(ctx)

		// Pre-compute registration check while we still have SDK context.
		// The async goroutine uses context.Background() which cannot access
		// the Cosmos KV store.
		alreadyRegistered := k.isAlreadyRegistered(gaslessCtx, roundNum)
		if alreadyRegistered {
			return nil
		}

		// Pre-compute old code commitment while we still have SDK context.
		oldCC, _ := k.getOldCodeCommitment(gaslessCtx)

		asyncCtx, cancel := dkgAsyncContext()

		go func() {
			defer cancel()

			k.handleDKGRegistration(asyncCtx, &dkgNetwork, oldCC, alreadyRegistered)
		}()
	}

	return nil
}

// isAlreadyRegistered checks whether this validator has an existing on-chain
// registration for the given round. Used to prevent re-registration with
// potentially different keys after sealed_keys deletion or kernel restart.
func (k *Keeper) isAlreadyRegistered(ctx context.Context, round uint32) bool {
	addr := ethcommon.HexToAddress(k.validatorEVMAddr)
	reg, err := k.getDKGRegistration(ctx, round, addr)
	if err != nil || reg == nil {
		return false
	}

	if len(reg.DkgPubKey) > 0 {
		log.Info(ctx, "Validator already registered on-chain; skipping re-registration",
			"round", round,
			"status", reg.Status.String(),
		)

		return true
	}

	return false
}

func (k *Keeper) shouldReshare(ctx context.Context) (bool, error) {
	activeNetwork, err := k.getLatestActiveDKGNetwork(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get latest active DKG network")
	}

	return activeNetwork != nil, nil
}
