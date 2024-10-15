// Modifications made from https://github.com/cosmos/cosmos-sdk/blob/main/types/module/module.go
// to accommodate the versioned module mapping in the modified Manager.
package module

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sort"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/genesis"
	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkmodule "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/piplabs/story/lib/errors"
)

// AppModuleBasic is the standard form for basic non-dependant elements of an application module.
type AppModuleBasic interface {
	HasName
	HasConsensusVersion

	RegisterLegacyAminoCodec(*codec.LegacyAmino)
	RegisterInterfaces(types.InterfaceRegistry)
	RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux)
}

// HasName allows the module to provide its own name for legacy purposes.
// Newer apps should specify the name for their modules using a map
// using NewManagerFromMap.
type HasName interface {
	Name() string
}

// HasGenesisBasics is the legacy interface for stateless genesis methods.
type HasGenesisBasics interface {
	DefaultGenesis(codec.JSONCodec) json.RawMessage
	ValidateGenesis(codec.JSONCodec, client.TxEncodingConfig, json.RawMessage) error
}

// BasicManager is a collection of AppModuleBasic.
type BasicManager map[string]AppModuleBasic

// NewBasicManager creates a new BasicManager object.
func NewBasicManager(modules ...AppModuleBasic) BasicManager {
	moduleMap := make(map[string]AppModuleBasic)
	for _, module := range modules {
		moduleMap[module.Name()] = module
	}
	return moduleMap
}

// NewBasicManagerFromManager creates a new BasicManager from a Manager
// The BasicManager will contain all AppModuleBasic from the AppModule Manager
// Module's AppModuleBasic can be overridden by passing a custom AppModuleBasic map.
func NewBasicManagerFromManager(manager *Manager, version uint64, customModuleBasics map[string]AppModuleBasic) BasicManager {
	moduleMap := make(map[string]AppModuleBasic)
	for name, module := range manager.versionedModules[version] {
		if customBasicMod, ok := customModuleBasics[name]; ok {
			moduleMap[name] = customBasicMod
			continue
		}

		if appModule, ok := module.(appmodule.AppModule); ok {
			moduleMap[name] = sdkmodule.CoreAppModuleBasicAdaptor(name, appModule)
			continue
		}

		if basicMod, ok := module.(AppModuleBasic); ok {
			moduleMap[name] = basicMod
		}
	}

	return moduleMap
}

// RegisterLegacyAminoCodec registers all module codecs.
func (bm BasicManager) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	for _, b := range bm {
		b.RegisterLegacyAminoCodec(cdc)
	}
}

// RegisterInterfaces registers all module interface types.
func (bm BasicManager) RegisterInterfaces(registry types.InterfaceRegistry) {
	for _, m := range bm {
		m.RegisterInterfaces(registry)
	}
}

// DefaultGenesis provides default genesis information for all modules.
func (bm BasicManager) DefaultGenesis(cdc codec.JSONCodec) map[string]json.RawMessage {
	genesisData := make(map[string]json.RawMessage)
	for _, b := range bm {
		if mod, ok := b.(HasGenesisBasics); ok {
			genesisData[b.Name()] = mod.DefaultGenesis(cdc)
		}
	}

	return genesisData
}

// ValidateGenesis performs genesis state validation for all modules.
func (bm BasicManager) ValidateGenesis(cdc codec.JSONCodec, txEncCfg client.TxEncodingConfig, genesisData map[string]json.RawMessage) error {
	for _, b := range bm {
		// first check if the module is an adapted Core API Module
		if mod, ok := b.(HasGenesisBasics); ok {
			if err := mod.ValidateGenesis(cdc, txEncCfg, genesisData[b.Name()]); err != nil {
				return err
			}
		}
	}

	return nil
}

// RegisterGRPCGatewayRoutes registers all module rest routes.
func (bm BasicManager) RegisterGRPCGatewayRoutes(clientCtx client.Context, rtr *runtime.ServeMux) {
	for _, b := range bm {
		b.RegisterGRPCGatewayRoutes(clientCtx, rtr)
	}
}

// AddTxCommands adds all tx commands to the rootTxCmd.
func (bm BasicManager) AddTxCommands(rootTxCmd *cobra.Command) {
	for _, b := range bm {
		if mod, ok := b.(interface {
			GetTxCmd() *cobra.Command
		}); ok {
			if cmd := mod.GetTxCmd(); cmd != nil {
				rootTxCmd.AddCommand(cmd)
			}
		}
	}
}

// AddQueryCommands adds all query commands to the rootQueryCmd.
func (bm BasicManager) AddQueryCommands(rootQueryCmd *cobra.Command) {
	for _, b := range bm {
		if mod, ok := b.(interface {
			GetQueryCmd() *cobra.Command
		}); ok {
			if cmd := mod.GetQueryCmd(); cmd != nil {
				rootQueryCmd.AddCommand(cmd)
			}
		}
	}
}

// HasGenesis is the extension interface for stateful genesis methods.
type HasGenesis interface {
	HasGenesisBasics
	InitGenesis(sdk.Context, codec.JSONCodec, json.RawMessage)
	ExportGenesis(sdk.Context, codec.JSONCodec) json.RawMessage
}

// HasABCIGenesis is the extension interface for stateful genesis methods which returns validator updates.
type HasABCIGenesis interface {
	HasGenesisBasics
	InitGenesis(sdk.Context, codec.JSONCodec, json.RawMessage) []abci.ValidatorUpdate
	ExportGenesis(sdk.Context, codec.JSONCodec) json.RawMessage
}

// AppModule is the form for an application module. Most of
// its functionality has been moved to extension interfaces.
// Deprecated: use appmodule.AppModule with a combination of extension interfaes interfaces instead.
type AppModule interface {
	appmodule.AppModule

	AppModuleBasic
}

// HasInvariants is the interface for registering invariants.
type HasInvariants interface {
	// RegisterInvariants registers module invariants.
	RegisterInvariants(sdk.InvariantRegistry)
}

// HasServices is the interface for modules to register services.
type HasServices interface {
	// RegisterServices allows a module to register services.
	RegisterServices(Configurator)
}

// HasConsensusVersion is the interface for declaring a module consensus version.
type HasConsensusVersion interface {
	// ConsensusVersion is a sequence number for state-breaking change of the
	// module. It should be incremented on each consensus-breaking change
	// introduced by the module. To avoid wrong/empty versions, the initial version
	// should be set to 1.
	ConsensusVersion() uint64
}

// HasABCIEndblock is a released typo of HasABCIEndBlock.
// Deprecated: use HasABCIEndBlock instead.
type HasABCIEndblock HasABCIEndBlock

// HasABCIEndBlock is the interface for modules that need to run code at the end of the block.
type HasABCIEndBlock interface {
	AppModule
	EndBlock(context.Context) ([]abci.ValidatorUpdate, error)
}

var (
	_ appmodule.AppModule = (*GenesisOnlyAppModule)(nil)
	_ AppModuleBasic      = (*GenesisOnlyAppModule)(nil)
)

// AppModuleGenesis is the standard form for an application module genesis functions.
type AppModuleGenesis interface {
	AppModuleBasic
	HasABCIGenesis
}

// GenesisOnlyAppModule is an AppModule that only has import/export functionality.
type GenesisOnlyAppModule struct {
	AppModuleGenesis
}

// NewGenesisOnlyAppModule creates a new GenesisOnlyAppModule object.
func NewGenesisOnlyAppModule(amg AppModuleGenesis) GenesisOnlyAppModule {
	return GenesisOnlyAppModule{
		AppModuleGenesis: amg,
	}
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (GenesisOnlyAppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (GenesisOnlyAppModule) IsAppModule() {}

// RegisterInvariants is a placeholder function register no invariants.
func (GenesisOnlyAppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (gam GenesisOnlyAppModule) ConsensusVersion() uint64 { return 1 }

// Manager defines a module manager that provides the high level utility for
// managing and executing operations for a group of modules. This implementation
// was originally inspired by the module manager defined in Cosmos SDK but this
// implementation maps the state machine version to different versions of the
// module. It also provides a way to run migrations between different versions
// of a module.
type Manager struct {
	// versionedModules is a map from app version -> module name -> module.
	versionedModules map[uint64]map[string]AppModule
	// uniqueModuleVersions is a mapping of module name -> module consensus
	// version -> the range of app versions this particular module operates
	// over. The first element in the array represent the fromVersion and the
	// last the toVersion (this is inclusive).
	uniqueModuleVersions map[string]map[uint64][2]uint64
	allModules           []AppModule
	// firstVersion is the lowest app version supported.
	firstVersion uint64
	// lastVersion is the highest app version supported.
	lastVersion uint64

	OrderInitGenesis         []string
	OrderExportGenesis       []string
	OrderPreBlockers         []string
	OrderBeginBlockers       []string
	OrderEndBlockers         []string
	OrderPrepareCheckStaters []string
	OrderPrecommiters        []string
	OrderMigrations          []string
}

// NewManager returns a new Manager object.
func NewManager(modules []VersionedModule) (*Manager, error) {
	versionedModules := make(map[uint64]map[string]AppModule)
	allModules := make([]AppModule, len(modules))
	modulesStr := make([]string, 0, len(modules))
	preBlockModulesStr := make([]string, 0)
	uniqueModuleVersions := make(map[string]map[uint64][2]uint64)
	for idx, module := range modules {
		name := module.Name()
		moduleVersion := module.ConsensusVersion()

		if _, ok := module.Module.(sdkmodule.AppModule); !ok {
			panic(fmt.Sprintf("module %s does not implement sdkmodule.AppModule", name))
		}

		if module.FromVersion == 0 {
			return nil, sdkerrors.ErrInvalidVersion.Wrapf("v0 is not a valid version for module %s", module.Module.Name())
		}
		if module.FromVersion > module.ToVersion {
			return nil, sdkerrors.ErrLogic.Wrapf("FromVersion cannot be greater than ToVersion for module %s", module.Module.Name())
		}

		for version := module.FromVersion; version <= module.ToVersion; version++ {
			if versionedModules[version] == nil {
				versionedModules[version] = make(map[string]AppModule)
			}
			if _, exists := versionedModules[version][name]; exists {
				return nil, sdkerrors.ErrLogic.Wrapf("Two different modules with domain %s are registered with the same version %d", name, version)
			}
			versionedModules[version][name] = module.Module
		}

		allModules[idx] = module.Module
		modulesStr = append(modulesStr, name)

		if _, ok := module.Module.(appmodule.HasPreBlocker); ok {
			preBlockModulesStr = append(preBlockModulesStr, name)
		}

		if _, exists := uniqueModuleVersions[name]; !exists {
			uniqueModuleVersions[name] = make(map[uint64][2]uint64)
		}
		uniqueModuleVersions[name][moduleVersion] = [2]uint64{module.FromVersion, module.ToVersion}
	}
	firstVersion := slices.Min(getKeys(versionedModules))
	lastVersion := slices.Max(getKeys(versionedModules))

	m := &Manager{
		versionedModules:     versionedModules,
		uniqueModuleVersions: uniqueModuleVersions,
		allModules:           allModules,
		firstVersion:         firstVersion,
		lastVersion:          lastVersion,
		OrderInitGenesis:     modulesStr,
		OrderExportGenesis:   modulesStr,
		OrderBeginBlockers:   modulesStr,
		OrderEndBlockers:     modulesStr,
	}
	if err := m.checkUpgradeSchedule(); err != nil {
		return nil, err
	}
	return m, nil
}

// NewManagerFromMap creates a new Manager object from a map of module names to module implementations.
// This method should be used for apps and modules which have migrated to the cosmossdk.io/core.appmodule.AppModule API.
func NewManagerFromMap(moduleMap map[string]appmodule.AppModule) *Manager {
	simpleModuleMap := make(map[string]interface{})
	modulesStr := make([]string, 0, len(simpleModuleMap))
	preBlockModulesStr := make([]string, 0)
	for name, module := range moduleMap {
		simpleModuleMap[name] = module
		modulesStr = append(modulesStr, name)
		if _, ok := module.(appmodule.HasPreBlocker); ok {
			preBlockModulesStr = append(preBlockModulesStr, name)
		}
	}

	// Sort the modules by name. Given that we are using a map above we can't guarantee the order.
	sort.Strings(modulesStr)

	return &Manager{
		Modules:                  simpleModuleMap,
		OrderInitGenesis:         modulesStr,
		OrderExportGenesis:       modulesStr,
		OrderPreBlockers:         preBlockModulesStr,
		OrderBeginBlockers:       modulesStr,
		OrderEndBlockers:         modulesStr,
		OrderPrecommiters:        modulesStr,
		OrderPrepareCheckStaters: modulesStr,
	}
}

// SetOrderInitGenesis sets the order of init genesis calls.
func (m *Manager) SetOrderInitGenesis(moduleNames ...string) {
	m.assertNoForgottenModules("SetOrderInitGenesis", moduleNames, func(moduleName string) bool {
		module, found := m.FindModule(moduleName)
		if !found {
			return false
		}

		if _, hasGenesis := module.(appmodule.HasGenesis); hasGenesis {
			return !hasGenesis
		}

		if _, hasABCIGenesis := module.(HasABCIGenesis); hasABCIGenesis {
			return !hasABCIGenesis
		}

		_, hasGenesis := module.(HasGenesis)
		return !hasGenesis
	})
	m.OrderInitGenesis = moduleNames
}

// SetOrderExportGenesis sets the order of export genesis calls.
func (m *Manager) SetOrderExportGenesis(moduleNames ...string) {
	m.assertNoForgottenModules("SetOrderExportGenesis", moduleNames, func(moduleName string) bool {
		module, found := m.FindModule(moduleName)
		if !found {
			return false
		}

		if _, hasGenesis := module.(appmodule.HasGenesis); hasGenesis {
			return !hasGenesis
		}

		if _, hasABCIGenesis := module.(HasABCIGenesis); hasABCIGenesis {
			return !hasABCIGenesis
		}

		_, hasGenesis := module.(HasGenesis)
		return !hasGenesis
	})
	m.OrderExportGenesis = moduleNames
}

// SetOrderPreBlockers sets the order of set pre-blocker calls.
func (m *Manager) SetOrderPreBlockers(moduleNames ...string) {
	m.assertNoForgottenModules("SetOrderPreBlockers", moduleNames,
		func(moduleName string) bool {
			module, found := m.FindModule(moduleName)
			if !found {
				return false
			}

			_, hasBlock := module.(appmodule.HasPreBlocker)
			return !hasBlock
		})
	m.OrderPreBlockers = moduleNames
}

// SetOrderBeginBlockers sets the order of set begin-blocker calls.
func (m *Manager) SetOrderBeginBlockers(moduleNames ...string) {
	m.assertNoForgottenModules("SetOrderBeginBlockers", moduleNames,
		func(moduleName string) bool {
			module, found := m.FindModule(moduleName)
			if !found {
				return false
			}

			_, hasBeginBlock := module.(appmodule.HasBeginBlocker)
			return !hasBeginBlock
		})
	m.OrderBeginBlockers = moduleNames
}

// SetOrderEndBlockers sets the order of set end-blocker calls.
func (m *Manager) SetOrderEndBlockers(moduleNames ...string) {
	m.assertNoForgottenModules("SetOrderEndBlockers", moduleNames,
		func(moduleName string) bool {
			module, found := m.FindModule(moduleName)
			if !found {
				return false
			}

			if _, hasEndBlock := module.(appmodule.HasEndBlocker); hasEndBlock {
				return !hasEndBlock
			}

			_, hasABCIEndBlock := module.(HasABCIEndBlock)
			return !hasABCIEndBlock
		})
	m.OrderEndBlockers = moduleNames
}

// SetOrderPrepareCheckStaters sets the order of set prepare-check-stater calls.
func (m *Manager) SetOrderPrepareCheckStaters(moduleNames ...string) {
	m.assertNoForgottenModules("SetOrderPrepareCheckStaters", moduleNames,
		func(moduleName string) bool {
			module, found := m.FindModule(moduleName)
			if !found {
				return false
			}

			_, hasPrepareCheckState := module.(appmodule.HasPrepareCheckState)
			return !hasPrepareCheckState
		})
	m.OrderPrepareCheckStaters = moduleNames
}

// SetOrderPrecommiters sets the order of set precommiter calls.
func (m *Manager) SetOrderPrecommiters(moduleNames ...string) {
	m.assertNoForgottenModules("SetOrderPrecommiters", moduleNames,
		func(moduleName string) bool {
			module, found := m.FindModule(moduleName)
			if !found {
				return false
			}

			_, hasPrecommit := module.(appmodule.HasPrecommit)
			return !hasPrecommit
		})
	m.OrderPrecommiters = moduleNames
}

// SetOrderMigrations sets the order of migrations to be run. If not set
// then migrations will be run with an order defined in `DefaultMigrationsOrder`.
func (m *Manager) SetOrderMigrations(moduleNames ...string) {
	m.assertNoForgottenModules("SetOrderMigrations", moduleNames, nil)
	m.OrderMigrations = moduleNames
}

// RegisterInvariants registers all module invariants.
func (m *Manager) RegisterInvariants(ir sdk.InvariantRegistry) {
	for _, module := range m.allModules {
		if module, ok := module.(HasInvariants); ok {
			module.RegisterInvariants(ir)
		}
	}
}

// RegisterServices registers all module services.
func (m *Manager) RegisterServices(cfg Configurator) error {
	for _, module := range m.allModules {
		if module, ok := module.(HasServices); ok {
			module.RegisterServices(cfg)
		}

		fromVersion, toVersion := m.getAppVersionsForModule(module.Name(), module.ConsensusVersion())

		if module, ok := module.(appmodule.HasServices); ok {
			err := module.RegisterServices(cfg.WithVersions(fromVersion, toVersion))
			if err != nil {
				return err
			}
		}

		if cfg.Error() != nil {
			return cfg.Error()
		}
	}

	return nil
}

func (m *Manager) getAppVersionsForModule(moduleName string, moduleVersion uint64) (uint64, uint64) {
	return m.uniqueModuleVersions[moduleName][moduleVersion][0], m.uniqueModuleVersions[moduleName][moduleVersion][1]
}

// InitGenesis performs init genesis functionality for modules. Exactly one
// module must return a non-empty validator set update to correctly initialize
// the chain.
func (m *Manager) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, genesisData map[string]json.RawMessage) (*abci.ResponseInitChain, error) {
	var validatorUpdates []abci.ValidatorUpdate
	ctx.Logger().Info("initializing blockchain state from genesis.json")

	appVersion := ctx.BlockHeader().Version.App
	modules, versionSupported := m.versionedModules[appVersion]
	if !versionSupported {
		panic(fmt.Sprintf("version %d not supported", appVersion))
	}

	for _, moduleName := range m.OrderInitGenesis {
		if genesisData[moduleName] == nil {
			continue
		}

		mod := modules[moduleName]
		// we might get an adapted module, a native core API module or a legacy module
		if module, ok := mod.(appmodule.HasGenesis); ok {
			ctx.Logger().Debug("running initialization for module", "module", moduleName)
			// core API genesis
			source, err := genesis.SourceFromRawJSON(genesisData[moduleName])
			if err != nil {
				return &abci.ResponseInitChain{}, err
			}

			err = module.InitGenesis(ctx, source)
			if err != nil {
				return &abci.ResponseInitChain{}, err
			}
		} else if module, ok := mod.(HasGenesis); ok {
			ctx.Logger().Debug("running initialization for module", "module", moduleName)
			module.InitGenesis(ctx, cdc, genesisData[moduleName])
		} else if module, ok := mod.(HasABCIGenesis); ok {
			ctx.Logger().Debug("running initialization for module", "module", moduleName)
			moduleValUpdates := module.InitGenesis(ctx, cdc, genesisData[moduleName])

			// use these validator updates if provided, the module manager assumes
			// only one module will update the validator set
			if len(moduleValUpdates) > 0 {
				if len(validatorUpdates) > 0 {
					return &abci.ResponseInitChain{}, errors.New("validator InitGenesis updates already set by a previous module")
				}
				validatorUpdates = moduleValUpdates
			}
		}
	}

	// a chain must initialize with a non-empty validator set
	if len(validatorUpdates) == 0 {
		return &abci.ResponseInitChain{}, fmt.Errorf("validator set is empty after InitGenesis, please ensure at least one validator is initialized with a delegation greater than or equal to the DefaultPowerReduction (%d)", sdk.DefaultPowerReduction)
	}

	return &abci.ResponseInitChain{
		Validators: validatorUpdates,
	}, nil
}

// ExportGenesis performs export genesis functionality for the modules supported
// in a particular version.
func (m *Manager) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec, version uint64) map[string]json.RawMessage {
	genesisData := make(map[string]json.RawMessage)
	modules := m.versionedModules[version]
	moduleNamesForVersion := m.ModuleNames(version)
	moduleNamesToExport := filter(m.OrderExportGenesis, func(moduleName string) bool {
		// filter out modules that are not supported by this version
		return slices.Contains(moduleNamesForVersion, moduleName)
	})
	for _, moduleName := range moduleNamesToExport {
		if module, ok := modules[moduleName].(HasGenesis); ok {
			genesisData[moduleName] = module.ExportGenesis(ctx, cdc)
		}
	}

	return genesisData
}

// ExportGenesisForModules performs export genesis functionality for modules.
func (m *Manager) ExportGenesisForModules(ctx sdk.Context, cdc codec.JSONCodec, modulesToExport []string) (map[string]json.RawMessage, error) {
	modules := m.versionedModules[m.lastVersion]

	if len(modulesToExport) == 0 {
		modulesToExport = m.OrderExportGenesis
	}
	// verify modules exists in app, so that we don't panic in the middle of an export
	if err := m.checkModulesExists(modulesToExport); err != nil {
		return nil, err
	}

	type genesisResult struct {
		bz  json.RawMessage
		err error
	}

	channels := make(map[string]chan genesisResult)
	for _, moduleName := range modulesToExport {
		mod := modules[moduleName]
		if module, ok := mod.(appmodule.HasGenesis); ok {
			// core API genesis
			channels[moduleName] = make(chan genesisResult)
			go func(module appmodule.HasGenesis, ch chan genesisResult) {
				ctx := ctx.WithGasMeter(storetypes.NewInfiniteGasMeter()) // avoid race conditions
				target := genesis.RawJSONTarget{}
				err := module.ExportGenesis(ctx, target.Target())
				if err != nil {
					ch <- genesisResult{nil, err}
					return
				}

				rawJSON, err := target.JSON()
				if err != nil {
					ch <- genesisResult{nil, err}
					return
				}

				ch <- genesisResult{rawJSON, nil}
			}(module, channels[moduleName])
		} else if module, ok := mod.(HasGenesis); ok {
			channels[moduleName] = make(chan genesisResult)
			go func(module HasGenesis, ch chan genesisResult) {
				ctx := ctx.WithGasMeter(storetypes.NewInfiniteGasMeter()) // avoid race conditions
				ch <- genesisResult{module.ExportGenesis(ctx, cdc), nil}
			}(module, channels[moduleName])
		} else if module, ok := mod.(HasABCIGenesis); ok {
			channels[moduleName] = make(chan genesisResult)
			go func(module HasABCIGenesis, ch chan genesisResult) {
				ctx := ctx.WithGasMeter(storetypes.NewInfiniteGasMeter()) // avoid race conditions
				ch <- genesisResult{module.ExportGenesis(ctx, cdc), nil}
			}(module, channels[moduleName])
		}
	}

	genesisData := make(map[string]json.RawMessage)
	for moduleName := range channels {
		res := <-channels[moduleName]
		if res.err != nil {
			return nil, fmt.Errorf("genesis export error in %s: %w", moduleName, res.err)
		}

		genesisData[moduleName] = res.bz
	}

	return genesisData, nil
}

// checkModulesExists verifies that all modules in the list exist in the app.
func (m *Manager) checkModulesExists(moduleName []string) error {
	modules := m.versionedModules[m.lastVersion]
	for _, name := range moduleName {
		if _, ok := modules[name]; !ok {
			return fmt.Errorf("module %s does not exist", name)
		}
	}

	return nil
}

// assertNoForgottenModules checks that we didn't forget any modules in the SetOrder* functions.
// `pass` is a closure which allows one to omit modules from `moduleNames`.
// If you provide non-nil `pass` and it returns true, the module would not be subject of the assertion.
func (m *Manager) assertNoForgottenModules(setOrderFnName string, moduleNames []string, pass func(moduleName string) bool) {
	ms := make(map[string]bool)
	for _, m := range moduleNames {
		ms[m] = true
	}
	var missing []string
	for m := range m.uniqueModuleVersions {
		m := m
		if pass != nil && pass(m) {
			continue
		}

		if !ms[m] {
			missing = append(missing, m)
		}
	}
	if len(missing) != 0 {
		sort.Strings(missing)
		panic(fmt.Sprintf(
			"all modules must be defined when setting %s, missing: %v", setOrderFnName, missing))
	}
}

// MigrationHandler is the migration function that each module registers.
type MigrationHandler func(sdk.Context) error

// VersionMap is a map of moduleName -> version.
type VersionMap map[string]uint64

// RunMigrations performs in-place store migrations for all modules. This
// function MUST be called when the state machine changes appVersion.
func (m Manager) RunMigrations(ctx sdk.Context, cfg sdkmodule.Configurator, fromVersion, toVersion uint64) error {
	c, ok := cfg.(Configurator)
	if !ok {
		return sdkerrors.ErrInvalidType.Wrapf("expected %T, got %T", Configurator{}, cfg)
	}
	modules := m.OrderMigrations
	if modules == nil {
		modules = defaultMigrationsOrder(m.ModuleNames(toVersion))
	}
	currentVersionModules, exists := m.versionedModules[fromVersion]
	if !exists {
		return sdkerrors.ErrInvalidVersion.Wrapf("fromVersion %d not supported", fromVersion)
	}
	nextVersionModules, exists := m.versionedModules[toVersion]
	if !exists {
		return sdkerrors.ErrInvalidVersion.Wrapf("toVersion %d not supported", toVersion)
	}

	for _, moduleName := range modules {
		currentModule, currentModuleExists := currentVersionModules[moduleName]
		nextModule, nextModuleExists := nextVersionModules[moduleName]

		// if the module exists for both upgrades
		if currentModuleExists && nextModuleExists {
			// by using consensus version instead of app version we support the SDK's legacy method
			// of migrating modules which were made of several versions and consisted of a mapping of
			// app version to module version. Now, using go.mod, each module will have only a single
			// consensus version and each breaking upgrade will result in a new module and a new consensus
			// version.
			fromModuleVersion := currentModule.ConsensusVersion()
			toModuleVersion := nextModule.ConsensusVersion()
			err := c.runModuleMigrations(ctx, moduleName, fromModuleVersion, toModuleVersion)
			if err != nil {
				return err
			}
		} else if !currentModuleExists && nextModuleExists {
			ctx.Logger().Info(fmt.Sprintf("adding a new module: %s", moduleName))

			if module, ok := nextModule.(HasGenesis); ok {
				module.InitGenesis(ctx, c.cdc, module.DefaultGenesis(c.cdc))
			}

			if module, ok := nextModule.(HasABCIGenesis); ok {
				moduleValUpdates := module.InitGenesis(ctx, c.cdc, module.DefaultGenesis(c.cdc))
				// The module manager assumes only one module will update the
				// validator set, and it can't be a new module.
				if len(moduleValUpdates) > 0 {
					return errorsmod.Wrapf(sdkerrors.ErrLogic, "validator InitGenesis update is already set by another module")
				}
			}
		}
		// TODO: handle the case where a module is no longer supported (i.e. removed from the state machine)
	}

	return nil
}

// PreBlock performs begin block functionality for upgrade module.
// It takes the current context as a parameter and returns a boolean value
// indicating whether the migration was successfully executed or not.
func (m *Manager) PreBlock(ctx sdk.Context) (*sdk.ResponsePreBlock, error) {
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	modules := m.MustGetCurrentVersionModules(ctx)
	paramsChanged := false
	for _, moduleName := range m.OrderPreBlockers {
		if module, ok := modules[moduleName].(appmodule.HasPreBlocker); ok {
			rsp, err := module.PreBlock(ctx)
			if err != nil {
				return nil, err
			}
			if rsp.IsConsensusParamsChanged() {
				paramsChanged = true
			}
		}
	}
	return &sdk.ResponsePreBlock{
		ConsensusParamsChanged: paramsChanged,
	}, nil
}

// BeginBlock performs begin block functionality for all modules. It creates a
// child context with an event manager to aggregate events emitted from all
// modules.
func (m *Manager) BeginBlock(ctx sdk.Context) (sdk.BeginBlock, error) {
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	modules := m.MustGetCurrentVersionModules(ctx)
	for _, moduleName := range m.OrderBeginBlockers {
		if module, ok := modules[moduleName].(appmodule.HasBeginBlocker); ok {
			if err := module.BeginBlock(ctx); err != nil {
				return sdk.BeginBlock{}, err
			}
		}
	}

	return sdk.BeginBlock{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

// EndBlock performs end block functionality for all modules. It creates a
// child context with an event manager to aggregate events emitted from all
// modules.
func (m *Manager) EndBlock(ctx sdk.Context) (sdk.EndBlock, error) {
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	validatorUpdates := []abci.ValidatorUpdate{}

	modules := m.MustGetCurrentVersionModules(ctx)
	for _, moduleName := range m.OrderEndBlockers {
		if module, ok := modules[moduleName].(appmodule.HasEndBlocker); ok {
			err := module.EndBlock(ctx)
			if err != nil {
				return sdk.EndBlock{}, err
			}
		} else if module, ok := modules[moduleName].(HasABCIEndBlock); ok {
			moduleValUpdates, err := module.EndBlock(ctx)
			if err != nil {
				return sdk.EndBlock{}, err
			}
			// use these validator updates if provided, the module manager assumes
			// only one module will update the validator set
			if len(moduleValUpdates) > 0 {
				if len(validatorUpdates) > 0 {
					return sdk.EndBlock{}, errors.New("validator EndBlock updates already set by a previous module")
				}

				for _, updates := range moduleValUpdates {
					validatorUpdates = append(validatorUpdates, abci.ValidatorUpdate{PubKey: updates.PubKey, Power: updates.Power})
				}
			}
		} else {
			continue
		}
	}

	return sdk.EndBlock{
		ValidatorUpdates: validatorUpdates,
		Events:           ctx.EventManager().ABCIEvents(),
	}, nil
}

// Precommit performs precommit functionality for all modules.
func (m *Manager) Precommit(ctx sdk.Context) error {
	modules := m.MustGetCurrentVersionModules(ctx)
	for _, moduleName := range m.OrderPrecommiters {
		module, ok := modules[moduleName].(appmodule.HasPrecommit)
		if !ok {
			continue
		}
		if err := module.Precommit(ctx); err != nil {
			return err
		}
	}
	return nil
}

// PrepareCheckState performs functionality for preparing the check state for all modules.
func (m *Manager) PrepareCheckState(ctx sdk.Context) error {
	modules := m.MustGetCurrentVersionModules(ctx)
	for _, moduleName := range m.OrderPrepareCheckStaters {
		module, ok := modules[moduleName].(appmodule.HasPrepareCheckState)
		if !ok {
			continue
		}
		if err := module.PrepareCheckState(ctx); err != nil {
			return err
		}
	}
	return nil
}

// GetVersionMap gets consensus version from all modules.
func (m *Manager) GetVersionMap(version uint64) sdkmodule.VersionMap {
	vermap := make(sdkmodule.VersionMap)
	if version > m.lastVersion || version < m.firstVersion {
		return vermap
	}

	for _, v := range m.versionedModules[version] {
		version := v.ConsensusVersion()
		name := v.Name()
		vermap[name] = version
	}

	return vermap
}

// ModuleNames returns the list of module names that are supported for a
// particular version in no particular order.
func (m *Manager) ModuleNames(version uint64) []string {
	modules, ok := m.versionedModules[version]
	if !ok {
		return []string{}
	}

	names := make([]string, 0, len(modules))
	for name := range modules {
		names = append(names, name)
	}
	return names
}

// DefaultMigrationsOrder returns a default migrations order: ascending alphabetical by module name,
// except x/auth which will run last, see:
// https://github.com/cosmos/cosmos-sdk/issues/10591
func DefaultMigrationsOrder(modules []string) []string {
	const authName = "auth"
	out := make([]string, 0, len(modules))
	hasAuth := false
	for _, m := range modules {
		if m == authName {
			hasAuth = true
		} else {
			out = append(out, m)
		}
	}
	sort.Strings(out)
	if hasAuth {
		out = append(out, authName)
	}
	return out
}

// SupportedVersions returns all the supported versions for the module manager.
func (m *Manager) SupportedVersions() []uint64 {
	return getKeys(m.versionedModules)
}

// checkUpgradeSchedule performs a dry run of all the upgrades in all versions and asserts that the consensus version
// for a module domain i.e. auth, always increments for each module that uses the auth domain name.
func (m *Manager) checkUpgradeSchedule() error {
	if m.firstVersion == m.lastVersion {
		// there are no upgrades to check
		return nil
	}
	for _, moduleName := range m.OrderInitGenesis {
		lastConsensusVersion := uint64(0)
		for appVersion := m.firstVersion; appVersion <= m.lastVersion; appVersion++ {
			module, exists := m.versionedModules[appVersion][moduleName]
			if !exists {
				continue
			}
			moduleVersion := module.ConsensusVersion()
			if moduleVersion < lastConsensusVersion {
				return fmt.Errorf("error: module %s in appVersion %d goes from moduleVersion %d to %d", moduleName, appVersion, lastConsensusVersion, moduleVersion)
			}
			lastConsensusVersion = moduleVersion
		}
	}
	return nil
}

// assertMatchingModules performs a sanity check that the basic module manager
// contains all the same modules present in the module manager.
func (m *Manager) AssertMatchingModules(basicModuleManager sdkmodule.BasicManager) error {
	for _, module := range m.allModules {
		if _, exists := basicModuleManager[module.Name()]; !exists {
			return fmt.Errorf("module %s not found in basic module manager", module.Name())
		}
	}
	return nil
}

func (m *Manager) MustGetCurrentVersionModules(ctx sdk.Context) map[string]AppModule {
	modules := m.versionedModules[ctx.BlockHeader().Version.App]
	if modules == nil {
		panic(fmt.Sprintf("no modules for version %d", ctx.BlockHeader().Version.App))
	}

	return modules
}

func (m *Manager) FindModule(moduleName string) (AppModule, bool) {
	for _, module := range m.allModules {
		if module.Name() == moduleName {
			return module, true
		}
	}

	return nil, false
}
