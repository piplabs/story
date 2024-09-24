package keeper

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmengine/types"
)

func TestKeeper_InitGenesis(t *testing.T) {
	t.Parallel()
	dummyExecutionHead := common.HexToHash("0x047e24c3455107d87c68dffa307b3b7fa1877f3e9d7f30c7ee359f2eff3a75d9")
	validParams := types.NewParams(dummyExecutionHead.Bytes())

	tcs := []struct {
		name           string
		gs             func() *types.GenesisState
		setup          func(c context.Context, k *Keeper)
		postStateCheck func(c context.Context, k *Keeper)
		expectedError  string
		requirePanic   bool
	}{
		{
			name: "pass",
			gs: func() *types.GenesisState {
				return &types.GenesisState{
					Params: validParams,
				}
			},
			postStateCheck: func(c context.Context, k *Keeper) {
				params, err := k.GetParams(c)
				require.NoError(t, err)
				require.Equal(t, validParams, params)
			},
		},
		{
			name: "fail: invalid execution block hash",
			gs: func() *types.GenesisState {
				invalidParams := validParams
				invalidParams.ExecutionBlockHash = []byte("invalid")

				return &types.GenesisState{
					Params: invalidParams,
				}
			},
			expectedError: "invalid execution block hash length",
		},
		{
			name: "panic: execution head already exists",
			setup: func(c context.Context, k *Keeper) {
				require.NoError(t, k.InsertGenesisHead(c, dummyExecutionHead.Bytes()))
			},
			gs: func() *types.GenesisState {
				return &types.GenesisState{
					Params: validParams,
				}
			},
			expectedError: "insert genesis head: unexpected genesis head id",
			requirePanic:  true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, keeper := createTestKeeper(t)
			if tc.setup != nil {
				tc.setup(ctx, keeper)
			}
			if tc.requirePanic {
				require.PanicsWithError(t, tc.expectedError, func() {
					_ = keeper.InitGenesis(ctx, tc.gs())
				})
			} else {
				err := keeper.InitGenesis(ctx, tc.gs())
				if tc.expectedError != "" {
					require.Error(t, err)
					require.Contains(t, err.Error(), tc.expectedError)
				} else {
					require.NoError(t, err)
					tc.postStateCheck(ctx, keeper)
				}
			}
		})
	}
}

func TestKeeper_ExportGenesis(t *testing.T) {
	t.Parallel()
	dummyExecutionHead := common.HexToHash("0x047e24c3455107d87c68dffa307b3b7fa1877f3e9d7f30c7ee359f2eff3a75d9")
	validParams := types.NewParams(dummyExecutionHead.Bytes())

	tcs := []struct {
		name           string
		setup          func(c context.Context, k *Keeper)
		postStateCheck func(c context.Context, k *Keeper)
	}{
		{
			name: "pass",
			setup: func(c context.Context, k *Keeper) {
				require.NoError(t, k.SetParams(c, validParams))
			},
			postStateCheck: func(c context.Context, k *Keeper) {
				gs := k.ExportGenesis(sdk.UnwrapSDKContext(c))
				require.Equal(t, validParams, gs.Params)
			},
		},
		{
			name: "pass: default params",
			postStateCheck: func(c context.Context, k *Keeper) {
				gs := k.ExportGenesis(sdk.UnwrapSDKContext(c))
				require.Equal(t, types.DefaultParams(), gs.Params)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, keeper := createTestKeeper(t)
			if tc.setup != nil {
				tc.setup(ctx, keeper)
			}
			tc.postStateCheck(ctx, keeper)
		})
	}
}
