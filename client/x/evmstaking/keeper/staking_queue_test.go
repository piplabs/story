package keeper_test

import (
	"context"
	"time"

	"cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/piplabs/story/lib/k1util"

	"go.uber.org/mock/gomock"
)

// setupValidatorAndDelegation creates a validator and delegation for testing.
func (s *TestSuite) setupValidatorAndDelegation(ctx sdk.Context, valPubKey, delPubKey crypto.PubKey, valAddr sdk.ValAddress, delAddr sdk.AccAddress) {
	require := s.Require()
	bankKeeper, stakingKeeper, keeper := s.BankKeeper, s.StakingKeeper, s.EVMStakingKeeper

	// Convert public key to cosmos format
	valCosmosPubKey, err := k1util.PubKeyToCosmos(valPubKey)
	require.NoError(err)

	// Create and update validator
	val := testutil.NewValidator(s.T(), valAddr, valCosmosPubKey)
	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
	validator, _ := val.AddTokensFromDel(valTokens)
	bankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.NotBondedPoolName, stypes.BondedPoolName, gomock.Any())
	_ = skeeper.TestingUpdateValidator(stakingKeeper, ctx, validator, true)

	// Create and set delegation
	delAmt := stakingKeeper.TokensFromConsensusPower(ctx, 100).ToLegacyDec()
	delegation := stypes.NewDelegation(delAddr.String(), valAddr.String(), delAmt)
	require.NoError(stakingKeeper.SetDelegation(ctx, delegation))

	// Map delegator to EVM address
	delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(delPubKey.Bytes())
	require.NoError(err)
	require.NoError(keeper.DelegatorMap.Set(ctx, delAddr.String(), delEvmAddr.String()))

	// Ensure delegation is set correctly
	delegation, err = stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	require.NoError(err)
	require.Equal(delAmt, delegation.GetShares())
}

// setupUnbonding creates unbondings for testing
func (s *TestSuite) setupUnbonding(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount string) {
	require := s.Require()
	bankKeeper, stakingKeeper := s.BankKeeper, s.StakingKeeper

	bankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.BondedPoolName, stypes.NotBondedPoolName, gomock.Any())
	_, _, err := stakingKeeper.Undelegate(ctx, delAddr, valAddr, math.LegacyMustNewDecFromStr(amount))
	require.NoError(err)
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
	valPubKey := pubKeys[1]
	valAddr := valAddrs[1]
	valPubKey2 := pubKeys[2]
	valAddr2 := valAddrs[2]
	s.setupValidatorAndDelegation(ctx, valPubKey, delPubKey, valAddr, delAddr)
	s.setupValidatorAndDelegation(ctx, valPubKey2, delPubKey, valAddr2, delAddr)

	tcs := []struct {
		name           string
		setUnbondings  func(c context.Context)
		expectedResult []stypes.DVPair
		expectedError  string
	}{
		{
			name:           "pass: no unbondings",
			expectedResult: nil,
		},
		{
			name: "pass: one unbonding",
			setUnbondings: func(c context.Context) {
				s.setupUnbonding(c, delAddr, valAddr, "100")
			},
			expectedResult: []stypes.DVPair{
				{delAddr.String(), valAddr.String()},
			},
		},
		{
			name: "pass: multiple unbondings",
			setUnbondings: func(c context.Context) {
				// Set up multiple unbondings
				s.setupUnbonding(c, delAddr, valAddr, "50")
				s.setupUnbonding(c, delAddr, valAddr2, "50")
			},
			expectedResult: []stypes.DVPair{
				{delAddr.String(), valAddr.String()},
				{delAddr.String(), valAddr2.String()},
			},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cachedCtx, _ := ctx.CacheContext()
			if tc.setUnbondings != nil {
				tc.setUnbondings(cachedCtx)
			}

			// Set block time to be after unbonding period
			header := cachedCtx.BlockHeader()
			header.Time = header.Time.Add(params.UnbondingTime).Add(24 * time.Hour)
			result, err := keeper.GetMatureUnbondedDelegations(cachedCtx.WithBlockHeader(header))

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
