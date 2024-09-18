package keepers

// Keepers have been moved to a separate package to ensure all keepers are accessible in `upgrades.Upgrade.CreateUpgradeHandler`.
// This allows for passing all keepers into the upgrade handler and accessing/changing blockchain state across all modules.
// When performing `ignite scaffold` the keeper will be added to `app.go`. Please move them here.

import (
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	evmengkeeper "github.com/piplabs/story/client/x/evmengine/keeper"
	evmstakingkeeper "github.com/piplabs/story/client/x/evmstaking/keeper"
	mintkeeper "github.com/piplabs/story/client/x/mint/keeper"
	signalkeeper "github.com/piplabs/story/client/x/signal/keeper"
)

// Keepers includes all possible keepers. We separated it into a separate struct to make it easier to scaffold upgrades.
type Keepers struct {
	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	SignalKeeper          signalkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	ConsensusParamsKeeper consensuskeeper.Keeper
	GovKeeper             *govkeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper

	// Story
	EvmStakingKeeper *evmstakingkeeper.Keeper
	EVMEngKeeper     *evmengkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
}
