package keeper

import (
	"context"

	"github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
)

// This method performs no validation of the parameters.
func (k *Keeper) SetParams(ctx context.Context, params types.Params) error {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return errors.Wrap(err, "marshal params")
	}

	err = store.Set(types.ParamsKey, bz)
	if err != nil {
		return errors.Wrap(err, "set params")
	}

	return nil
}

func (k *Keeper) GetParams(ctx context.Context) (params types.Params, err error) {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := store.Get(types.ParamsKey)
	if err != nil {
		return params, errors.Wrap(err, "get params")
	}

	if bz == nil {
		return params, nil
	}

	err = k.cdc.Unmarshal(bz, &params)
	if err != nil {
		return params, errors.Wrap(err, "unmarshal params")
	}

	// TODO: remove hardcoded overrides after testing
	params.RegistrationPeriod = 50 // ~2.5 min at 3s blocks
	params.DealingPeriod = 30      // ~1.5 min
	params.FinalizationPeriod = 30 // ~1.5 min
	params.ActivePeriod = 1814400  // 21 days
	params.MinReqRegisteredParticipants = 1
	params.MinReqFinalizedParticipants = 1
	params.OperationalThreshold = 500 // 50%

	return params, nil
}

// SetMinReqRegisteredParticipants updates the min_req_registered_participants param.
// Called when a MinReqRegisteredParticipantsSet event is received from DKG.sol.
func (k *Keeper) SetMinReqRegisteredParticipants(ctx context.Context, value uint32) error {
	if err := types.ValidateMinReqRegisteredParticipants(value); err != nil {
		return err
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "get params")
	}

	params.MinReqRegisteredParticipants = value

	return k.SetParams(ctx, params)
}

// SetMinReqFinalizedParticipants updates the min_req_finalized_participants param.
// Called when a MinReqFinalizedParticipantsSet event is received from DKG.sol.
func (k *Keeper) SetMinReqFinalizedParticipants(ctx context.Context, value uint32) error {
	if err := types.ValidateMinReqFinalizedParticipants(value); err != nil {
		return err
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "get params")
	}

	params.MinReqFinalizedParticipants = value

	return k.SetParams(ctx, params)
}

// SetOperationalThreshold updates the operational_threshold param.
// Called when an OperationalThresholdSet event is received from DKG.sol.
func (k *Keeper) SetOperationalThreshold(ctx context.Context, value uint32) error {
	if err := types.ValidateOperationalThreshold(value); err != nil {
		return err
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "get params")
	}

	params.OperationalThreshold = value

	return k.SetParams(ctx, params)
}
