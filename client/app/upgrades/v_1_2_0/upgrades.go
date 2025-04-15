//nolint:revive,stylecheck // versioning
package v_1_2_0

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	"github.com/piplabs/story/client/app/keepers"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func CreateUpgradeHandler(
	_ *module.Manager,
	_ module.Configurator,
	keepers *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		log.Info(ctx, "Start upgrade v1.2.0")

		// update consensus params
		oldConsParams, err := keepers.ConsensusParamsKeeper.ParamsStore.Get(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get consensus params")
		}

		updateConsParamsMsg := consensusparamtypes.MsgUpdateParams{
			Authority: keepers.ConsensusParamsKeeper.GetAuthority(),
			Block:     oldConsParams.Block,
			Evidence:  oldConsParams.Evidence,
			Validator: oldConsParams.Validator,
			Abci:      oldConsParams.Abci,
		}

		updateConsParamsMsg.Block.MaxBytes = newMaxBytes

		if _, err = keepers.ConsensusParamsKeeper.UpdateParams(ctx, &updateConsParamsMsg); err != nil {
			return vm, errors.Wrap(err, "failed to update consensus params")
		}

		newConsParams, err := keepers.ConsensusParamsKeeper.ParamsStore.Get(ctx)
		if err != nil {
			return vm, errors.Wrap(err, "failed to get updated consensus params")
		}

		if !oldConsParams.Evidence.Equal(newConsParams.Evidence) {
			return vm, errors.New("mismatch evidence in consensus params", "old_evidence", oldConsParams.Evidence, "new_evidence", newConsParams.Evidence)
		}

		if !oldConsParams.Validator.Equal(newConsParams.Validator) {
			return vm, errors.New("mismatch validator in consensus params", "old_validator", oldConsParams.Validator, "new_validator", newConsParams.Validator)
		}

		if !oldConsParams.Abci.Equal(newConsParams.Abci) {
			return vm, errors.New("mismatch ABCI in consensus params", "old_abci", oldConsParams.Abci, "new_abci", newConsParams.Abci)
		}

		if oldConsParams.Block.MaxGas != newConsParams.Block.MaxGas {
			return vm, errors.New("mismatch Block in consensus params", "old_max_gas", oldConsParams.Block.MaxGas, "new_max_gas", newConsParams.Block.MaxGas)
		}

		if newConsParams.Block.MaxBytes != newMaxBytes {
			return vm, errors.New("wrong updated max bytes in consensus params", "expected_max_bytes", newMaxBytes, "actual_max_bytes", newConsParams.Block.MaxBytes)
		}

		if oldConsParams.Version.App != newConsParams.Version.App {
			return vm, errors.New("mismatch App version in consensus params")
		}

		log.Info(ctx, "Upgrade v1.2.0 complete")

		return vm, nil
	}
}
