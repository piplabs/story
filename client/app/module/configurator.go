package module

import (
	"fmt"

	cosmosmsg "cosmossdk.io/api/cosmos/msg/v1"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/gogoproto/grpc"

	googlegrpc "google.golang.org/grpc"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Configurator implements the module.Configurator interface.
var _ module.Configurator = &configurator{}
var _ Configurator = &configurator{}

// Configurator is a struct used at startup to register all the message and
// query servers for all modules. It allows the module to register any migrations from
// one consensus version of the module to the next. Finally it maps all the messages
// to the app versions that they are accepted in. This then gets used in the antehandler
// to prevent users from submitting messages that can not yet be executed.
type Configurator interface {
	grpc.Server

	// Error returns the last error encountered during RegisterService.
	Error() error

	// MsgServer returns a grpc.Server instance which allows registering services
	// that will handle TxBody.messages in transactions. These Msg's WILL NOT
	// be exposed as gRPC services.
	MsgServer() grpc.Server

	// QueryServer returns a grpc.Server instance which allows registering services
	// that will be exposed as gRPC services as well as ABCI query handlers.
	QueryServer() grpc.Server

	// RegisterMigration registers an in-place store migration for a module. The
	// handler is a migration script to perform in-place migrations from version
	// `fromVersion` to version `fromVersion+1`.
	//
	// EACH TIME a module's ConsensusVersion increments, a new migration MUST
	// be registered using this function. If a migration handler is missing for
	// a particular function, the upgrade logic (see RunMigrations function)
	// will panic. If the ConsensusVersion bump does not introduce any store
	// changes, then a no-op function must be registered here.
	RegisterMigration(moduleName string, fromVersion uint64, handler module.MigrationHandler) error
}

type configurator struct {
	fromVersion, toVersion uint64
	cdc                    codec.Codec
	msgServer              grpc.Server
	queryServer            grpc.Server
	// acceptedMessages is a map from appVersion -> msgTypeURL -> struct{}.
	acceptedMessages map[uint64]map[string]struct{}
	// migrations is a map of moduleName -> fromVersion -> migration script handler.
	migrations map[string]map[uint64]module.MigrationHandler
	err        error
}

// RegisterService implements the grpc.Server interface.
func (c *configurator) RegisterService(sd *googlegrpc.ServiceDesc, ss interface{}) {
	desc, err := c.cdc.InterfaceRegistry().FindDescriptorByName(protoreflect.FullName(sd.ServiceName))
	if err != nil {
		c.err = err
		return
	}

	if protobuf.HasExtension(desc.Options(), cosmosmsg.E_Service) {
		c.msgServer.RegisterService(sd, ss)
	} else {
		c.queryServer.RegisterService(sd, ss)
	}
}

// Error returns the last error encountered during RegisterService.
func (c *configurator) Error() error {
	return c.err
}

// NewConfigurator returns a new Configurator instance.
func NewConfigurator(cdc codec.Codec, msgServer, queryServer grpc.Server) Configurator {
	return &configurator{
		cdc:              cdc,
		msgServer:        msgServer,
		queryServer:      queryServer,
		migrations:       map[string]map[uint64]module.MigrationHandler{},
		acceptedMessages: map[uint64]map[string]struct{}{},
	}
}

// MsgServer implements the Configurator.MsgServer method.
func (c configurator) MsgServer() grpc.Server {
	return &serverWrapper{
		addMessages: c.addMessages,
		msgServer:   c.msgServer,
	}
}

// QueryServer implements the Configurator.QueryServer method.
func (c *configurator) QueryServer() grpc.Server {
	return c.queryServer
}

// WithVersions(version uint64, version2 uint64) googlegrpc.ServiceRegistrar.
func (c *configurator) WithVersions(fromVersion, toVersion uint64) module.Configurator {
	c.fromVersion = fromVersion
	c.toVersion = toVersion

	return c
}

// GetAcceptedMessages returns the accepted messages for all versions.
// acceptedMessages is a map from appVersion -> msgTypeURL -> struct{}.
func (c *configurator) GetAcceptedMessages() map[uint64]map[string]struct{} {
	return c.acceptedMessages
}

// RegisterMigration implements the Configurator.RegisterMigration method.
func (c *configurator) RegisterMigration(moduleName string, fromVersion uint64, handler module.MigrationHandler) error {
	if fromVersion == 0 {
		return sdkerrors.ErrInvalidVersion.Wrap("module migration versions should start at 1")
	}

	if c.migrations[moduleName] == nil {
		c.migrations[moduleName] = map[uint64]module.MigrationHandler{}
	}

	if c.migrations[moduleName][fromVersion] != nil {
		return sdkerrors.ErrLogic.Wrapf("another migration for module %s and version %d already exists", moduleName, fromVersion)
	}

	c.migrations[moduleName][fromVersion] = handler

	return nil
}

func (c *configurator) addMessages(msgs []string) {
	for version := c.fromVersion; version <= c.toVersion; version++ {
		if _, exists := c.acceptedMessages[version]; !exists {
			c.acceptedMessages[version] = map[string]struct{}{}
		}
		for _, msg := range msgs {
			c.acceptedMessages[version][msg] = struct{}{}
		}
	}
}

// runModuleMigrations runs all in-place store migrations for one given module from a
// version to another version.
func (c *configurator) runModuleMigrations(ctx sdk.Context, moduleName string, fromVersion, toVersion uint64) error {
	// No-op if toVersion is the initial version or if the version is unchanged.
	if toVersion <= 1 || fromVersion == toVersion {
		return nil
	}

	moduleMigrationsMap, found := c.migrations[moduleName]
	if !found {
		return sdkerrors.ErrNotFound.Wrapf("no migrations found for module %s", moduleName)
	}

	// Run in-place migrations for the module sequentially until toVersion.
	for i := fromVersion; i < toVersion; i++ {
		migrateFn, found := moduleMigrationsMap[i]
		if !found {
			// no migrations needed
			continue
		}
		ctx.Logger().Info(fmt.Sprintf("migrating module %s from version %d to version %d", moduleName, i, i+1))

		err := migrateFn(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
