package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/mint/types"
)

var minterAcc = authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter)

func TestImportExportGenesis(t *testing.T) {
	ctx, ak, _, _, mk := createKeeper(t)

	ak.EXPECT().GetModuleAccount(ctx, minterAcc.Name).Return(minterAcc)

	genesisState := types.DefaultGenesisState()
	genesisState.Params = types.NewParams(
		"test",
		math.LegacyNewDec(24625000000000000.000000000000000000),
		uint64(60*60*8766/5),
	)

	mk.InitGenesis(ctx, ak, genesisState)

	params, err := mk.Params.Get(ctx)
	require.Equal(t, genesisState.Params, params)
	require.NoError(t, err)

	genesisState2 := mk.ExportGenesis(ctx)
	require.Equal(t, genesisState, genesisState2)
}
