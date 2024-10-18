package keepers

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/piplabs/story/client/app/encoding"

	evmenginekeeper "github.com/piplabs/story/client/x/evmengine/keeper"
	evmenginetypes "github.com/piplabs/story/client/x/evmengine/types"
	evmstakingkeeper "github.com/piplabs/story/client/x/evmstaking/keeper"
	evmstakingtypes "github.com/piplabs/story/client/x/evmstaking/types"
	mintkeeper "github.com/piplabs/story/client/x/mint/keeper"
	minttypes "github.com/piplabs/story/client/x/mint/types"
	signalkeeper "github.com/piplabs/story/client/x/signal/keeper"
	signaltypes "github.com/piplabs/story/client/x/signal/types"
	"github.com/piplabs/story/lib/ethclient"

	storetypes "cosmossdk.io/store/types"
)

type Keepers struct {
	// keepers, by order of initialization
	// "Special" keepers
	ParamsKeeper *paramskeeper.Keeper
	//CrisisKeeper          *crisiskeeper.Keeper
	ConsensusParamsKeeper *consensusparamkeeper.Keeper

	// "Normal" keepers
	AccountKeeper *authkeeper.AccountKeeper
	// don't pass as pointer because bank module's legacy RegisterServices casts to (keeper.BaseKeeper) which will throw
	// panic if this is a pointer to bank keeper:
	// interface conversion: keeper.Keeper is *keeper.BaseKeeper, not keeper.BaseKeeper
	BankKeeper     bankkeeper.BaseKeeper
	StakingKeeper  *stakingkeeper.Keeper
	DistrKeeper    *distrkeeper.Keeper
	SlashingKeeper *slashingkeeper.Keeper
	GovKeeper      *govkeeper.Keeper

	// Normal Story keepers
	EvmStakingKeeper *evmstakingkeeper.Keeper
	EVMEngKeeper     *evmenginekeeper.Keeper
	MintKeeper       *mintkeeper.Keeper
	SignalKeeper     *signalkeeper.Keeper
}

// InitNormalKeepers initializes all 'normal' keepers (account, app, bank, auth, staking, distribution, slashing, transfer, gamm, IBC router, pool incentives, governance, mint, txfees keepers).
func (keepers *Keepers) InitNormalKeepers(
	appCodec codec.Codec,
	encodingConfig encoding.EncodingConfig,
	bApp *baseapp.BaseApp,
	maccPerms map[string][]string,
	blockedAddress map[string]bool,
	txConfig client.TxConfig,
	engineCl ethclient.EngineClient,
	keys map[string]*storetypes.KVStoreKey,
) {
	legacyAmino := encodingConfig.Amino
	// Add 'normal' keepers
	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		encoding.AccountAddressPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keepers.AccountKeeper = &accountKeeper
	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		keepers.AccountKeeper,
		blockedAddress,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		bApp.Logger(),
	)
	keepers.BankKeeper = bankKeeper

	govModuleAddr := keepers.AccountKeeper.GetModuleAddress(govtypes.ModuleName)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		keepers.AccountKeeper,
		keepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
	keepers.StakingKeeper = stakingKeeper

	distrKeeper := distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		keepers.AccountKeeper,
		keepers.BankKeeper,
		keepers.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keepers.DistrKeeper = &distrKeeper

	slashingKeeper := slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		keepers.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keepers.SlashingKeeper = &slashingKeeper

	mintKeeper := mintkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[minttypes.StoreKey]),
		keepers.StakingKeeper,
		keepers.AccountKeeper,
		keepers.BankKeeper,
		authtypes.FeeCollectorName,
	)
	keepers.MintKeeper = &mintKeeper

	// register the proposal types
	govRouter := govtypesv1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypesv1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(*keepers.ParamsKeeper))

	govConfig := govtypes.DefaultConfig()
	// Set the maximum metadata length for government-related configurations to 10,200, deviating from the default value of 256.
	govConfig.MaxMetadataLen = 10200
	govKeeper := govkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(keys[govtypes.StoreKey]),
		keepers.AccountKeeper, keepers.BankKeeper, keepers.StakingKeeper, bApp.MsgServiceRouter(),
		govConfig, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	keepers.GovKeeper = govKeeper
	keepers.GovKeeper.SetLegacyRouter(govRouter)

	// Add Story keepers
	evmStakingKeeper := evmstakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evmstakingtypes.StoreKey]),
		keepers.AccountKeeper,
		keepers.BankKeeper,
		keepers.SlashingKeeper,
		keepers.StakingKeeper,
		keepers.DistrKeeper,
		govModuleAddr.String(),
		engineCl,
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
	)
	keepers.EvmStakingKeeper = evmStakingKeeper

	signalKeeper := signalkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[signaltypes.StoreKey]),
		keepers.AccountKeeper,
		keepers.StakingKeeper,
		govModuleAddr.String(),
	)
	keepers.SignalKeeper = signalKeeper

	evmEngineKeeper, err := evmenginekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evmenginetypes.StoreKey]),
		engineCl,
		engineCl,
		txConfig,
		keepers.AccountKeeper,
		keepers.EvmStakingKeeper,
		keepers.MintKeeper,
		keepers.SignalKeeper,
	)
	if err != nil {
		panic(err)
	}
	keepers.EVMEngKeeper = evmEngineKeeper
}

// InitSpecialKeepers initiates special keepers (crisis appkeeper, params keeper)
func (keepers *Keepers) InitSpecialKeepers(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	cdc *codec.LegacyAmino,
	keys map[string]*storetypes.KVStoreKey,
	tkeys map[string]*storetypes.TransientStoreKey,
) {
	paramsKeeper := keepers.initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	keepers.ParamsKeeper = &paramsKeeper

	// set the BaseApp's parameter store
	consensusParamsKeeper := consensusparamkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]), authtypes.NewModuleAddress(govtypes.ModuleName).String(), runtime.EventService{})
	keepers.ConsensusParamsKeeper = &consensusParamsKeeper
	bApp.SetParamStore(keepers.ConsensusParamsKeeper.ParamsStore)

	// TODO: Make a SetInvCheckPeriod fn on CrisisKeeper.
	// IMO, its bad design atm that it requires this in state machine initialization
	//crisisKeeper := crisiskeeper.NewKeeper(
	//	appCodec, runtime.NewKVStoreService(keys[crisistypes.StoreKey]), invCheckPeriod, keepers.BankKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String(), addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()))
	//keepers.CrisisKeeper = crisisKeeper
}

// initParamsKeeper init params keeper and its subspaces.
func (keepers *Keepers) initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)

	paramsKeeper.Subspace(evmenginetypes.ModuleName)
	paramsKeeper.Subspace(evmstakingtypes.ModuleName)
	paramsKeeper.Subspace(signaltypes.ModuleName)

	return paramsKeeper
}

// SetupHooks sets up hooks for modules.
func (keepers *Keepers) SetupHooks() {
	// For every module that has hooks set on it,
	// you must check InitNormalKeepers to ensure that its not passed by de-reference
	// e.g. *app.StakingKeeper doesn't appear

	keepers.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			keepers.DistrKeeper.Hooks(),
			keepers.SlashingKeeper.Hooks(),
		),
	)

	keepers.GovKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// insert governance hooks receivers here
		),
	)
}
