package v_1_6_0

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	"github.com/piplabs/story/client/app/keepers"
	dkgtypes "github.com/piplabs/story/client/x/dkg/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		log.Info(ctx, "Start v1.6.0 upgrade — activating DKG module")

		// RunMigrations handles InitGenesis for truly new modules (not in vm).
		// On a fresh genesis chain, DKG is already in the vm so this is a no-op
		// for DKG. On a live chain upgraded from v1.5.3, DKG is new and gets
		// InitGenesis called automatically.
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		// Explicitly set DKG params to ensure they are initialized regardless
		// of whether RunMigrations called InitGenesis for DKG or not. This
		// covers the fresh-genesis case where DKG was in the module list but
		// had no genesis state in the static genesis.json.
		if err := keepers.DKGKeeper.SetParams(ctx, dkgtypes.DefaultParams()); err != nil {
			return newVM, errors.Wrap(err, "set DKG default params")
		}

		// Enable vote extensions at this upgrade height. The DKG module
		// requires vote extensions for aggregating DKG messages across
		// validators during consensus.
		if err := enableVoteExtensions(ctx, keepers, plan.Height); err != nil {
			return newVM, err
		}

		log.Info(ctx, "V1.6.0 upgrade complete — DKG module activated")

		return newVM, nil
	}
}

// enableVoteExtensions updates the consensus params to enable vote extensions
// starting from the upgrade height + 1 (must be a future height).
func enableVoteExtensions(ctx context.Context, keepers *keepers.Keepers, upgradeHeight int64) error {
	currentParams, err := keepers.ConsensusParamsKeeper.ParamsStore.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "get consensus params")
	}

	veHeight := upgradeHeight + 1

	// Vote extensions are already enabled — nothing to do.
	// This covers both:
	//   - Fresh genesis chains where VE is enabled from height 1
	//   - Re-runs where VE was already set to the correct upgrade height
	if currentParams.Abci != nil && currentParams.Abci.VoteExtensionsEnableHeight > 0 {
		existingHeight := currentParams.Abci.VoteExtensionsEnableHeight
		if existingHeight != veHeight {
			log.Warn(ctx, "Vote extensions enabled at unexpected height",
				errors.New("height mismatch"),
				"expected", veHeight,
				"actual", existingHeight,
			)
		}
		log.Info(ctx, "Vote extensions already enabled, skipping",
			"enabled_height", existingHeight,
		)

		return nil
	}
	abci := currentParams.Abci
	if abci == nil {
		abci = &cmtproto.ABCIParams{}
	}
	abci.VoteExtensionsEnableHeight = veHeight

	updateMsg := consensusparamtypes.MsgUpdateParams{
		Authority: keepers.ConsensusParamsKeeper.GetAuthority(),
		Block:     currentParams.Block,
		Evidence:  currentParams.Evidence,
		Validator: currentParams.Validator,
		Abci:      abci,
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if _, err := keepers.ConsensusParamsKeeper.UpdateParams(sdkCtx, &updateMsg); err != nil {
		return errors.Wrap(err, "enable vote extensions")
	}

	log.Info(ctx, "Enabled vote extensions", "height", veHeight)

	return nil
}
