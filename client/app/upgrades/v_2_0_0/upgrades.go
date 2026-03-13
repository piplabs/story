package v_2_0_0

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
	"github.com/piplabs/story/lib/netconf"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.Keepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		log.Info(ctx, "Start v2.0.0 upgrade — activating DKG module")

		// RunMigrations handles InitGenesis for truly new modules (not in vm).
		// On a fresh genesis chain, DKG is already in the vm so this is a no-op
		// for DKG. On a live chain upgraded from v1.5.3, DKG is new and gets
		// InitGenesis called automatically.
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		// Set DKG params, using shorter periods for devnet/test chains.
		dkgParams := dkgParamsForChain(ctx)
		if err := keepers.DKGKeeper.SetParams(ctx, dkgParams); err != nil {
			return newVM, errors.Wrap(err, "set DKG params")
		}

		// Enable vote extensions at this upgrade height. The DKG module
		// requires vote extensions for aggregating DKG messages across
		// validators during consensus.
		if err := enableVoteExtensions(ctx, keepers, plan.Height); err != nil {
			return newVM, err
		}

		log.Info(ctx, "V2.0.0 upgrade complete — DKG module activated")

		return newVM, nil
	}
}

// dkgParamsForChain returns DKG params appropriate for the current chain.
// Devnet/test chains use shorter stage periods for faster iteration.
func dkgParamsForChain(ctx context.Context) dkgtypes.Params {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	chainID := sdkCtx.ChainID()

	if chainID == netconf.DKGTestChainID || chainID == netconf.LocalChainID {
		log.Info(ctx, "Using devnet DKG params with short stage periods", "chain_id", chainID)

		return dkgtypes.NewParams(
			200, // registration: ~6.5 min at 2s blocks
			300, // dealing: ~10 min (kernel light client needs time to sync)
			300, // finalization: ~10 min
			600, // active: ~20 min
			dkgtypes.DefaultDkgCommitteeRewardPortion,
			dkgtypes.DefaultMinReqRegisteredParticipants,
			dkgtypes.DefaultMinReqFinalizedParticipants,
			dkgtypes.DefaultOperationalThreshold,
		)
	}

	return dkgtypes.DefaultParams()
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
