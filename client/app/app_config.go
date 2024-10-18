package app

import (
	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/piplabs/story/client/app/encoding"

	evmenginetypes "github.com/piplabs/story/client/x/evmengine/types"
	evmstakingtypes "github.com/piplabs/story/client/x/evmstaking/types"
	minttypes "github.com/piplabs/story/client/x/mint/types"
	signaltypes "github.com/piplabs/story/client/x/signal/types"
)

const (
	defaultPruningKeep     = 72_000 // Keep 1 day's of application state by default (given period of 1.2s).
	defaultPruningInterval = 300    // Prune every 5 minutes or so.
)

// init initializes the Cosmos SDK configuration.
//
//nolint:gochecknoinits // Cosmos-style
func init() {
	// Set and seal config
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(encoding.AccountAddressPrefix, encoding.AccountPubKeyPrefix)
	cfg.SetBech32PrefixForValidator(encoding.ValidatorAddressPrefix, encoding.ValidatorPubKeyPrefix)
	cfg.SetBech32PrefixForConsensusNode(encoding.ConsNodeAddressPrefix, encoding.ConsNodePubKeyPrefix)
	cfg.Seal()
}

var (
	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
	genesisModuleOrder = []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		evmstakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		consensusparamtypes.ModuleName,
		//crisistypes.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		signaltypes.ModuleName,
		evmenginetypes.ModuleName,
	}

	// TODO: signaltypes.ModuleName should be in preblocker like the new upgrade module in 0.50.x

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0.
	beginBlockers = []string{
		minttypes.ModuleName,
		distrtypes.ModuleName, // Note: slashing happens after distr.BeginBlocker
		slashingtypes.ModuleName,
		stakingtypes.ModuleName,
		signaltypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		//crisistypes.ModuleName,
		govtypes.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		signaltypes.ModuleName,
	}

	endBlockers = []string{
		govtypes.ModuleName,
		signaltypes.ModuleName,
		//crisistypes.ModuleName,
		govtypes.ModuleName,
		evmstakingtypes.ModuleName, // Must be before staking module removes mature unbonding delegations & validators.
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		signaltypes.ModuleName,
	}

	moduleAccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter, authtypes.Burner},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		evmstakingtypes.ModuleName:     {authtypes.Burner, authtypes.Minter},
		govtypes.ModuleName:            {authtypes.Burner},
		signaltypes.ModuleName:         nil,
	}

	overrideStoreKeys = []*runtimev1alpha1.StoreKeyConfig{
		{
			ModuleName: authtypes.ModuleName,
			KvStoreKey: "acc",
		},
	}
)
