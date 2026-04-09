package module_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/require"

	dkgmodule "github.com/piplabs/story/client/x/dkg/module"
	"github.com/piplabs/story/client/x/dkg/types"
)

func TestAppModuleBasic_Name(t *testing.T) {
	t.Parallel()

	m := dkgmodule.NewAppModuleBasic(nil)
	require.Equal(t, types.ModuleName, m.Name())
	require.Equal(t, "dkg", m.Name())
}

func TestAppModule_ConsensusVersion(t *testing.T) {
	t.Parallel()

	m := dkgmodule.NewAppModule(nil, nil)
	require.Equal(t, uint64(dkgmodule.ConsensusVersion), m.ConsensusVersion())
	require.Equal(t, uint64(1), m.ConsensusVersion())
}

func TestAppModule_IsOnePerModuleType(t *testing.T) {
	t.Parallel()

	m := dkgmodule.AppModule{}
	// Should not panic
	require.NotPanics(t, func() {
		m.IsOnePerModuleType()
	})
}

func TestAppModule_IsAppModule(t *testing.T) {
	t.Parallel()

	m := dkgmodule.AppModule{}
	// Should not panic
	require.NotPanics(t, func() {
		m.IsAppModule()
	})
}

func TestAppModuleBasic_RegisterLegacyAminoCodec(t *testing.T) {
	t.Parallel()

	m := dkgmodule.NewAppModuleBasic(nil)
	// Should not panic with nil codec
	require.NotPanics(t, func() {
		m.RegisterLegacyAminoCodec(nil)
	})
}

func TestAppModuleBasic_RegisterGRPCGatewayRoutes(t *testing.T) {
	t.Parallel()

	m := dkgmodule.NewAppModuleBasic(nil)
	// Should not panic with zero-value context and nil mux (no-op implementation)
	require.NotPanics(t, func() {
		m.RegisterGRPCGatewayRoutes(client.Context{}, nil)
	})
}
