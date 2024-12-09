package keeper_test

import (
	"context"
	"math/big"

	sdkmath "cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/k1util"
)

// createValidator creates a validator.
func (s *TestSuite) createValidator(ctx context.Context, valPubKey crypto.PubKey, valAddr sdk.ValAddress) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	require := s.Require()
	stakingKeeper := s.StakingKeeper

	// Convert public key to cosmos format
	valCosmosPubKey, err := k1util.PubKeyToCosmos(valPubKey)
	require.NoError(err)

	// Create and update validator
	val := testutil.NewValidator(s.T(), valAddr, valCosmosPubKey)
	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
	validator, _, _ := val.AddTokensFromDel(valTokens, sdkmath.LegacyOneDec())
	// bankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.NotBondedPoolName, stypes.BondedPoolName, gomock.Any())
	require.NoError(s.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, valTokens))))
	require.NoError(s.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, stypes.NotBondedPoolName, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, valTokens))))
	_ = skeeper.TestingUpdateValidator(stakingKeeper, sdkCtx, validator, true)
}

func (s *TestSuite) TestProcessDeposit() {
	require := s.Require()
	ctx, keeper, stakingKeeper := s.Ctx, s.EVMStakingKeeper, s.StakingKeeper

	pubKeys, accAddrs, valAddrs := createAddresses(2)
	// delegator
	delPubKey := pubKeys[0]
	delAddr := accAddrs[0]
	// validator
	valPubKey := pubKeys[1]
	valAddr := valAddrs[1]
	s.createValidator(ctx, valPubKey, valAddr)

	createDeposit := func(delPubKey, valPubKey []byte, amount *big.Int) *bindings.IPTokenStakingDeposit {
		return &bindings.IPTokenStakingDeposit{
			DelegatorUncmpPubkey: cmpToUncmp(delPubKey),
			ValidatorUncmpPubkey: cmpToUncmp(valPubKey),
			StakeAmount:          amount,
			StakingPeriod:        big.NewInt(0),
			DelegationId:         big.NewInt(0),
			OperatorAddress:      cmpToEVM(delPubKey),
		}
	}

	tcs := []struct {
		name           string
		settingMock    func()
		deposit        *bindings.IPTokenStakingDeposit
		expectedResult stypes.Delegation
		expectedErr    string
	}{
		{
			name: "fail: invalid delegator pubkey",
			deposit: &bindings.IPTokenStakingDeposit{
				DelegatorUncmpPubkey: cmpToUncmp(delPubKey.Bytes())[:16],
				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes()),
				StakeAmount:          new(big.Int).SetUint64(1),
				StakingPeriod:        big.NewInt(0),
				DelegationId:         big.NewInt(0),
				OperatorAddress:      cmpToEVM(delPubKey.Bytes()),
			},
			expectedErr: "invalid uncompressed public key length or format",
		},
		{
			name: "fail: invalid validator pubkey",
			deposit: &bindings.IPTokenStakingDeposit{
				DelegatorUncmpPubkey: cmpToUncmp(delPubKey.Bytes()),
				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes())[:16],
				StakeAmount:          new(big.Int).SetUint64(1),
				StakingPeriod:        big.NewInt(0),
				DelegationId:         big.NewInt(0),
				OperatorAddress:      cmpToEVM(delPubKey.Bytes()),
			},
			expectedErr: "invalid uncompressed public key length or format",
		},
		// TODO: corrupted delegator and validator pubkey
		{
			name:    "pass: existing delegator",
			deposit: createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(1)),
			expectedResult: stypes.Delegation{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Shares:           sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(1)),
				RewardsShares:    sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(1)).Quo(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2))),
			},
		},
		{
			name:    "pass: new delegator",
			deposit: createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(1)),
			expectedResult: stypes.Delegation{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Shares:           sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(1)),
				RewardsShares:    sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(1)).Quo(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2))),
			},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			if tc.settingMock != nil {
				tc.settingMock()
			}
			cachedCtx, _ := ctx.CacheContext()
			err := keeper.ProcessDeposit(cachedCtx, tc.deposit)
			if tc.expectedErr != "" {
				require.ErrorContains(err, tc.expectedErr)
			} else {
				require.NoError(err)
				// check delegation
				delegation, err := stakingKeeper.GetDelegation(cachedCtx, delAddr, valAddr)
				require.NoError(err)
				require.Equal(tc.expectedResult, delegation)
			}
		})
	}
}

func (s *TestSuite) TestParseDepositLog() {
	require := s.Require()
	keeper := s.EVMStakingKeeper

	testCases := []struct {
		name      string
		log       gethtypes.Log
		expectErr bool
	}{
		{
			name: "Unknown Topic",
			log: gethtypes.Log{
				Topics: []common.Hash{common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")},
			},
			expectErr: true,
		},
		{
			name: "Valid Topic",
			log: gethtypes.Log{
				Topics: []common.Hash{types.DepositEvent.ID},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := keeper.ParseDepositLog(tc.log)
			if tc.expectErr {
				require.Error(err, "should return error for %s", tc.name)
			} else {
				require.NoError(err, "should not return error for %s", tc.name)
			}
		})
	}
}
