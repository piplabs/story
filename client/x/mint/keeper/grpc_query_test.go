package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/mint/keeper"
	mintmodule "github.com/piplabs/story/client/x/mint/module"
	"github.com/piplabs/story/client/x/mint/types"
)

func TestGRPCParams(t *testing.T) {
	//nolint:dogsled // This is common helper function
	ctx, _, _, _, mk := createKeeper(t)

	encCfg := moduletestutil.MakeTestEncodingConfig(mintmodule.AppModuleBasic{})
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, keeper.NewQuerier(*mk))

	queryClient := types.NewQueryClient(queryHelper)

	params, err := queryClient.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	got, err := mk.Params.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, params.Params, got)
}
