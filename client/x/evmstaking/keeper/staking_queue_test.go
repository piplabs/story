package keeper_test

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// setupMaturedUnbonding creates matured unbondings for testing.
func (s *TestSuite) setupMatureUnbonding(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amt string, duration time.Duration) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	pastHeader := sdkCtx.BlockHeader()
	pastHeader.Time = pastHeader.Time.Add(-duration).Add(-time.Minute)
	pastCtx := sdkCtx.WithBlockHeader(pastHeader)

	s.setupUnbonding(pastCtx, delAddr, valAddr, amt)
}

func (s *TestSuite) TestGetMatureUnbondedDelegations() {
	require := s.Require()
	ctx, keeper, stakingKeeper := s.Ctx, s.EVMStakingKeeper, s.StakingKeeper

	// Setup unbonding period
	params, err := stakingKeeper.GetParams(ctx)
	require.NoError(err)
	params.UnbondingTime, err = time.ParseDuration("3600s")
	require.NoError(err)
	require.NoError(stakingKeeper.SetParams(ctx, params))

	// Setup 2 validators
	pubKeys, accAddrs, valAddrs := createAddresses(3)
	delAddr := accAddrs[0]
	delPubKey := pubKeys[0]
	valPubKey1 := pubKeys[1]
	valAddr1 := valAddrs[1]
	valPubKey2 := pubKeys[2]
	valAddr2 := valAddrs[2]
	// self delegation
	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
	s.setupValidatorAndDelegation(ctx, valPubKey1, delPubKey, valAddr1, delAddr, valTokens)
	s.setupValidatorAndDelegation(ctx, valPubKey2, delPubKey, valAddr2, delAddr, valTokens)

	// set staking module's params
	ubdTime := time.Duration(3600) * time.Second
	stakingParams, err := stakingKeeper.GetParams(ctx)
	require.NoError(err)
	stakingParams.UnbondingTime = ubdTime
	require.NoError(err)
	require.NoError(stakingKeeper.SetParams(ctx, stakingParams))

	tcs := []struct {
		name           string
		preRun         func(c context.Context)
		setUnbondings  func(c context.Context)
		expectedResult []stypes.DVPair
		expectedError  string
	}{
		{
			name: "pass: no matured unbondings",
			preRun: func(c context.Context) {
				// s.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.NotBondedPoolName, stypes.BondedPoolName, gomock.Any())
			},
			setUnbondings: func(c context.Context) {
				s.setupUnbonding(c, delAddr, valAddr1, "100")
				s.setupUnbonding(c, delAddr, valAddr2, "100")
			},
			expectedResult: nil,
		},
		{
			name: "pass: one matured and one not matured unbonding",
			preRun: func(c context.Context) {
				// s.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.NotBondedPoolName, stypes.BondedPoolName, gomock.Any())
			},
			setUnbondings: func(c context.Context) {
				s.setupMatureUnbonding(c, delAddr, valAddr1, "100", ubdTime)
				s.setupUnbonding(c, delAddr, valAddr2, "100")
			},
			expectedResult: []stypes.DVPair{
				{DelegatorAddress: delAddr.String(), ValidatorAddress: valAddr1.String()},
			},
		},
		{
			name: "pass: two matured unbondings",
			preRun: func(c context.Context) {
				// s.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.NotBondedPoolName, stypes.BondedPoolName, gomock.Any())
			},
			setUnbondings: func(c context.Context) {
				s.setupMatureUnbonding(c, delAddr, valAddr1, "100", ubdTime)
				s.setupMatureUnbonding(c, delAddr, valAddr2, "100", ubdTime)
			},
			expectedResult: []stypes.DVPair{
				{DelegatorAddress: delAddr.String(), ValidatorAddress: valAddr1.String()},
				{DelegatorAddress: delAddr.String(), ValidatorAddress: valAddr2.String()},
			},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cachedCtx, _ := ctx.CacheContext()
			if tc.preRun != nil {
				tc.preRun(cachedCtx)
			}
			if tc.setUnbondings != nil {
				tc.setUnbondings(cachedCtx)
			}

			result, err := keeper.GetMatureUnbondedDelegations(cachedCtx)

			// Run the test
			if tc.expectedError != "" {
				require.ErrorContains(err, tc.expectedError)
			} else {
				require.NoError(err)
				require.Equal(tc.expectedResult, result)
			}
		})
	}
}
