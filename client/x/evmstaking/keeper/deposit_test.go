package keeper_test

/*
import (
	"context"
	"math/big"
	"time"

	"cosmossdk.io/math"
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
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"

	"go.uber.org/mock/gomock"
)

// createValidator creates a validator.
func (s *TestSuite) createValidator(ctx context.Context, valPubKey crypto.PubKey, valAddr sdk.ValAddress) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	require := s.Require()
	bankKeeper, stakingKeeper := s.BankKeeper, s.StakingKeeper

	// Convert public key to cosmos format
	valCosmosPubKey, err := k1util.PubKeyToCosmos(valPubKey)
	require.NoError(err)

	// Create and update validator
	val := testutil.NewValidator(s.T(), valAddr, valCosmosPubKey)
	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
	validator, _, _ := val.AddTokensFromDel(valTokens, sdkmath.LegacyOneDec())
	bankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.NotBondedPoolName, stypes.BondedPoolName, gomock.Any())
	_ = skeeper.TestingUpdateValidator(stakingKeeper, sdkCtx, validator, true)
}

func (s *TestSuite) TestProcessDeposit() {
	require := s.Require()
	ctx, keeper, accountKeeper, bankKeeper, stakingKeeper := s.Ctx, s.EVMStakingKeeper, s.AccountKeeper, s.BankKeeper, s.StakingKeeper

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
	expectAccountMock := func(isNewAccount bool) {
		if isNewAccount {
			accountKeeper.EXPECT().HasAccount(gomock.Any(), delAddr).Return(false)
			accountKeeper.EXPECT().NewAccountWithAddress(gomock.Any(), delAddr).Return(nil)
			accountKeeper.EXPECT().SetAccount(gomock.Any(), gomock.Any())
		} else {
			accountKeeper.EXPECT().HasAccount(gomock.Any(), delAddr).Return(true)
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
		{
			name: "fail: corrupted delegator pubkey",
			deposit: &bindings.IPTokenStakingDeposit{
				DelegatorUncmpPubkey: createCorruptedPubKey(cmpToUncmp(delPubKey.Bytes())),
				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes()),
				StakeAmount:          new(big.Int).SetUint64(1),
				StakingPeriod:        big.NewInt(0),
				DelegationId:         big.NewInt(0),
				OperatorAddress:      cmpToEVM(delPubKey.Bytes()),
			},
			expectedErr: "invalid uncompressed public key length or format",
		},
		// {
		// 	name:        "fail: corrupted validator pubkey",
		// 	deposit:     createDeposit(delPubKey.Bytes(), createCorruptedPubKey(valPubKey.Bytes()), new(big.Int).SetUint64(1)),
		// 	expectedErr: "validator pubkey to evm address",
		// },
		{
			name: "fail: mint coins to existing delegator",
			settingMock: func() {
				accountKeeper.EXPECT().HasAccount(gomock.Any(), delAddr).Return(true)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(errors.New(""))
			},
			deposit:     createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(1)),
			expectedErr: "create stake coin for depositor: mint coins",
		},
		{
			name: "fail: mint coins to new delegator",
			settingMock: func() {
				expectAccountMock(true)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(errors.New(""))
			},
			deposit:     createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(1)),
			expectedErr: "create stake coin for depositor: mint coins",
		},
		{
			name: "fail: send coins from module to existing delegator",
			settingMock: func() {
				expectAccountMock(false)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, delAddr, gomock.Any()).Return(errors.New(""))
			},
			deposit:     createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(1)),
			expectedErr: "create stake coin for depositor: send coins",
		},
		{
			name: "fail: send coins from module to new delegator",
			settingMock: func() {
				expectAccountMock(true)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, delAddr, gomock.Any()).Return(errors.New(""))
			},
			deposit:     createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(1)),
			expectedErr: "create stake coin for depositor: send coins",
		},
		{
			name: "fail: delegate to existing delegator",
			settingMock: func() {
				expectAccountMock(false)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, delAddr, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), delAddr, stypes.BondedPoolName, gomock.Any()).Return(errors.New("failed to delegate"))
			},
			deposit:     createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(1)),
			expectedErr: "failed to delegate",
		},
		{
			name: "fail: delegate to new delegator",
			settingMock: func() {
				expectAccountMock(true)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, delAddr, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), delAddr, stypes.BondedPoolName, gomock.Any()).Return(errors.New("failed to delegate"))
			},
			deposit:     createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(1)),
			expectedErr: "failed to delegate",
		},
		{
			name: "pass: existing delegator",
			settingMock: func() {
				expectAccountMock(false)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, delAddr, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), delAddr, stypes.BondedPoolName, gomock.Any()).Return(nil)
			},
			deposit: createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(1)),
			expectedResult: stypes.Delegation{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Shares:           math.LegacyNewDecFromInt(math.NewInt(1)),
				RewardsShares:    math.LegacyNewDecFromInt(math.NewInt(1)),
				PeriodDelegations: map[string]*stypes.PeriodDelegation{
					stypes.FlexibleDelegationID: {
						PeriodDelegationId: stypes.FlexibleDelegationID,
						Shares:             math.LegacyNewDecFromInt(math.NewInt(1)),
						RewardsShares:      math.LegacyNewDecFromInt(math.NewInt(1)),
						EndTime:            time.Time{},
					},
				},
			},
		},
		{
			name: "pass: new delegator",
			settingMock: func() {
				expectAccountMock(true)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, delAddr, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), delAddr, stypes.BondedPoolName, gomock.Any()).Return(nil)
			},
			deposit: createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(1)),
			expectedResult: stypes.Delegation{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Shares:           math.LegacyNewDecFromInt(math.NewInt(1)),
				RewardsShares:    math.LegacyNewDecFromInt(math.NewInt(1)),
				PeriodDelegations: map[string]*stypes.PeriodDelegation{
					stypes.FlexibleDelegationID: {
						PeriodDelegationId: stypes.FlexibleDelegationID,
						Shares:             math.LegacyNewDecFromInt(math.NewInt(1)),
						RewardsShares:      math.LegacyNewDecFromInt(math.NewInt(1)),
						EndTime:            time.Time{},
					},
				},
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
				delegation.PeriodDelegations[stypes.FlexibleDelegationID].EndTime = time.Time{}
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
*/
