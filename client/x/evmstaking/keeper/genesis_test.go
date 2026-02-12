package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sttestutil "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmstaking/keeper"
	estestutil "github.com/piplabs/story/client/x/evmstaking/testutil"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"

	"go.uber.org/mock/gomock"
)

func TestInitGenesis(t *testing.T) {
	validMaxWithdrawalPerBlock := types.DefaultMaxWithdrawalPerBlock + 100
	validMaxsweepPerBlock := types.DefaultMaxSweepPerBlock + 100
	validMinPartialWithdrawalAmount := types.DefaultMinPartialWithdrawalAmount + 100
	validParams := types.Params{
		MaxWithdrawalPerBlock:      validMaxWithdrawalPerBlock,
		MaxSweepPerBlock:           validMaxsweepPerBlock,
		MinPartialWithdrawalAmount: validMinPartialWithdrawalAmount,
	}

	// setup addresses and keys for testing
	pubKeys, _, valAddrs := createAddresses(1)
	valAddr := valAddrs[0]
	valPubKey := pubKeys[0]
	valCosmosPubKey, err := k1util.PubKeyToCosmos(valPubKey)
	require.NoError(t, err)
	valEvmAddr1, err := k1util.CosmosPubkeyToEVMAddress(valPubKey.Bytes())
	require.NoError(t, err)

	tcs := []struct {
		name           string
		setupMock      func(c sdk.Context, sk *estestutil.MockStakingKeeper)
		gs             func() *types.GenesisState
		postStateCheck func(c sdk.Context, esk *keeper.Keeper)
		expectedError  string
	}{
		{
			name: "pass: no validators",
			setupMock: func(c sdk.Context, sk *estestutil.MockStakingKeeper) {
				sk.EXPECT().GetAllValidators(gomock.Any()).Return(nil, nil)
			},
			gs: func() *types.GenesisState {
				return &types.GenesisState{
					Params: validParams,
				}
			},
			postStateCheck: func(c sdk.Context, esk *keeper.Keeper) {
				params, err := esk.GetParams(c)
				require.NoError(t, err)
				require.Equal(t, validParams, params)
			},
		},
		{
			name: "pass: with validators",
			setupMock: func(c sdk.Context, sk *estestutil.MockStakingKeeper) {
				sk.EXPECT().GetAllValidators(gomock.All()).Return([]stypes.Validator{
					sttestutil.NewValidator(t, valAddr, valCosmosPubKey),
				}, nil)
			},
			gs: func() *types.GenesisState {
				return &types.GenesisState{
					Params: validParams,
				}
			},
			postStateCheck: func(c sdk.Context, esk *keeper.Keeper) {
				params, err := esk.GetParams(c)
				require.NoError(t, err)
				require.Equal(t, validParams, params)

				// check withdraw and reward map
				evmAddrWithdraw, err := esk.DelegatorWithdrawAddress.Get(c, sdk.AccAddress(valEvmAddr1.Bytes()).String())
				require.NoError(t, err)
				require.Equal(t, valEvmAddr1.String(), evmAddrWithdraw)
				evmAddrReward, err := esk.DelegatorRewardAddress.Get(c, sdk.AccAddress(valEvmAddr1.Bytes()).String())
				require.NoError(t, err)
				require.Equal(t, valEvmAddr1.String(), evmAddrReward)
			},
		},
		{
			name: "fail: invalid params",
			setupMock: func(c sdk.Context, sk *estestutil.MockStakingKeeper) {
				sk.EXPECT().GetAllValidators(gomock.Any()).Return(nil, errors.New("failed to get all validators"))
			},
			gs: func() *types.GenesisState {
				return &types.GenesisState{
					Params: validParams,
				}
			},
			expectedError: "failed to get all validators",
		},
		{
			name: "fail: get all validators from staking keeper",
			gs: func() *types.GenesisState {
				invalidParams := validParams
				// make params invalid
				invalidParams.MaxWithdrawalPerBlock = 0

				return &types.GenesisState{
					Params: invalidParams,
				}
			},
			expectedError: "max withdrawal per block must be positive",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, sk, _, _, esk := createKeeperWithMockStaking(t)

			cachedCtx, _ := ctx.CacheContext()

			if tc.setupMock != nil {
				tc.setupMock(ctx, sk)
			}

			err := esk.InitGenesis(cachedCtx, tc.gs())
			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)
				if tc.postStateCheck != nil {
					tc.postStateCheck(cachedCtx, esk)
				}
			}
		})
	}
}

func TestExportGenesis(t *testing.T) {
	validMaxWithdrawalPerBlock := types.DefaultMaxWithdrawalPerBlock + 100
	validMaxsweepPerBlock := types.DefaultMaxSweepPerBlock + 100
	validMinPartialWithdrawalAmount := types.DefaultMinPartialWithdrawalAmount + 100
	validParams := types.Params{
		MaxWithdrawalPerBlock:      validMaxWithdrawalPerBlock,
		MaxSweepPerBlock:           validMaxsweepPerBlock,
		MinPartialWithdrawalAmount: validMinPartialWithdrawalAmount,
	}

	tcs := []struct {
		name            string
		setup           func(c sdk.Context, esk *keeper.Keeper)
		expectedGenesis *types.GenesisState
	}{
		{
			name: "pass: case1",
			setup: func(c sdk.Context, esk *keeper.Keeper) {
				cpy := validParams
				// modify params to test
				cpy.MaxWithdrawalPerBlock += 100
				cpy.MaxSweepPerBlock += 100
				cpy.MinPartialWithdrawalAmount += 100
				require.NoError(t, esk.SetParams(c, cpy))
			},
			expectedGenesis: &types.GenesisState{
				Params: types.Params{
					MaxWithdrawalPerBlock:      validParams.MaxWithdrawalPerBlock + 100,
					MaxSweepPerBlock:           validParams.MaxSweepPerBlock + 100,
					MinPartialWithdrawalAmount: validParams.MinPartialWithdrawalAmount + 100,
				},
			},
		},
		{
			name: "pass: case2",
			setup: func(c sdk.Context, esk *keeper.Keeper) {
				cpy := validParams
				// modify params to test
				cpy.MaxWithdrawalPerBlock += 2
				cpy.MaxSweepPerBlock += 2
				cpy.MinPartialWithdrawalAmount += 2
				require.NoError(t, esk.SetParams(c, cpy))
			},
			expectedGenesis: &types.GenesisState{
				Params: types.Params{
					MaxWithdrawalPerBlock:      validParams.MaxWithdrawalPerBlock + 2,
					MaxSweepPerBlock:           validParams.MaxSweepPerBlock + 2,
					MinPartialWithdrawalAmount: validParams.MinPartialWithdrawalAmount + 2,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

			cachedCtx, _ := ctx.CacheContext()
			if tc.setup != nil {
				tc.setup(cachedCtx, esk)
			}
			genesis := esk.ExportGenesis(cachedCtx)
			require.Equal(t, tc.expectedGenesis, genesis)
		})
	}
}
