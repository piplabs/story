package keeper_test

/*
import (
	"context"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/lib/k1util"
)

func (s *TestSuite) TestInitGenesis() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper

	validMaxWithdrawalPerBlock := types.DefaultMaxWithdrawalPerBlock + 100
	validMaxsweepPerBlock := types.DefaultMaxSweepPerBlock + 100
	validMinPartialWithdrawalAmount := types.DefaultMinPartialWithdrawalAmount + 100
	validParams := types.Params{
		MaxWithdrawalPerBlock:      validMaxWithdrawalPerBlock,
		MaxSweepPerBlock:           validMaxsweepPerBlock,
		MinPartialWithdrawalAmount: validMinPartialWithdrawalAmount,
	}

	// setup addresses and keys for testing
	pubKeys, addrs, valAddrs := createAddresses(3)
	delAddr := addrs[0]
	delPubKey := pubKeys[0]
	valAddr1 := valAddrs[1]
	valAccAddr1 := addrs[1]
	valPubKey1 := pubKeys[1]
	valEvmAddr1, err := k1util.CosmosPubkeyToEVMAddress(valPubKey1.Bytes())
	require.NoError(err)
	valAddr2 := valAddrs[2]
	valAccAddr2 := addrs[2]
	valPubKey2 := pubKeys[2]
	valEvmAddr2, err := k1util.CosmosPubkeyToEVMAddress(valPubKey2.Bytes())
	require.NoError(err)
	valTokens := s.StakingKeeper.TokensFromConsensusPower(ctx, 10)

	tcs := []struct {
		name           string
		setup          func(c context.Context)
		gs             func() *types.GenesisState
		postStateCheck func(c context.Context)
		expectedError  string
	}{
		{
			name: "pass: no validators",
			gs: func() *types.GenesisState {
				return &types.GenesisState{
					Params: validParams,
				}
			},
			postStateCheck: func(c context.Context) {
				params, err := keeper.GetParams(c)
				require.NoError(err)
				require.Equal(validParams, params)
			},
		},
		{
			name: "pass: with validators",
			setup: func(c context.Context) {
				s.setupValidatorAndDelegation(c, valPubKey1, delPubKey, valAddr1, delAddr, valTokens)
				s.setupValidatorAndDelegation(c, valPubKey2, delPubKey, valAddr2, delAddr, valTokens)
			},
			gs: func() *types.GenesisState {
				return &types.GenesisState{
					Params: validParams,
				}
			},
			postStateCheck: func(c context.Context) {
				params, err := keeper.GetParams(c)
				require.NoError(err)
				require.Equal(validParams, params)

				// check delegator map
				evmAddr1, err := keeper.DelegatorMap.Get(c, valAccAddr1.String())
				require.NoError(err)
				require.Equal(valEvmAddr1.String(), evmAddr1)
				evmAddr2, err := keeper.DelegatorMap.Get(c, valAccAddr2.String())
				require.NoError(err)
				require.Equal(valEvmAddr2.String(), evmAddr2)
			},
		},
		{
			name: "fail: invalid params",
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
		s.Run(tc.name, func() {
			cachedCtx, _ := ctx.CacheContext()
			if tc.setup != nil {
				tc.setup(cachedCtx)
			}
			err := keeper.InitGenesis(cachedCtx, tc.gs())
			if tc.expectedError != "" {
				require.Error(err)
				require.Contains(err.Error(), tc.expectedError)
			} else {
				require.NoError(err)
				if tc.postStateCheck != nil {
					tc.postStateCheck(cachedCtx)
				}
			}
		})
	}
}

func (s *TestSuite) TestExportGenesis() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper

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
		setup           func(c context.Context)
		expectedGenesis *types.GenesisState
	}{
		{
			name: "pass: case1",
			setup: func(c context.Context) {
				cpy := validParams
				// modify params to test
				cpy.MaxWithdrawalPerBlock += 100
				cpy.MaxSweepPerBlock += 100
				cpy.MinPartialWithdrawalAmount += 100
				require.NoError(keeper.SetParams(c, cpy))
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
			setup: func(c context.Context) {
				cpy := validParams
				// modify params to test
				cpy.MaxWithdrawalPerBlock += 2
				cpy.MaxSweepPerBlock += 2
				cpy.MinPartialWithdrawalAmount += 2
				require.NoError(keeper.SetParams(c, cpy))
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
		s.Run(tc.name, func() {
			cachedCtx, _ := ctx.CacheContext()
			if tc.setup != nil {
				tc.setup(cachedCtx)
			}
			genesis := keeper.ExportGenesis(cachedCtx)
			require.Equal(tc.expectedGenesis, genesis)
		})
	}
}
*/
