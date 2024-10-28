package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/k1util"
)

func (s *TestSuite) TestProcessCreateValidator() {
	require := s.Require()
	ctx, eskeeper, stakingKeeper := s.Ctx, s.EVMStakingKeeper, s.StakingKeeper

	pubKeys, addrs, valAddrs := createAddresses(3)

	uncmpPubKey0 := cmpToUncmp(pubKeys[0].Bytes())
	uncmpPubKey1 := cmpToUncmp(pubKeys[1].Bytes())
	uncmpPubKey2 := cmpToUncmp(pubKeys[2].Bytes())

	corruptedPubKey := append([]byte{}, uncmpPubKey0...)
	corruptedPubKey[0] = 0x04
	corruptedPubKey[1] = 0xFF

	tokens10 := stakingKeeper.TokensFromConsensusPower(ctx, 10)

	// checkDelegatorMapAndValidator checks if the delegator map and validator are created
	checkDelegatorMapAndValidator := func(c sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, delEvmAddr common.Address, _ math.Int) {
		val, err := eskeeper.DelegatorWithdrawAddress.Get(c, delAddr.String())
		require.NoError(err)
		require.Equal(delEvmAddr.String(), val)
		// check validator is created
		_, err = stakingKeeper.GetValidator(c, valAddr)
		require.NoError(err)
	}
	// checkDelegatorMapAndValTokens checks if the delegator map and validator tokens are added
	checkDelegatorMapAndValTokens := func(c sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, delEvmAddr common.Address, previousValTokens math.Int) {
		val, err := eskeeper.DelegatorWithdrawAddress.Get(c, delAddr.String())
		require.NoError(err)
		require.Equal(delEvmAddr.String(), val)
		// check validator tokens are added
		validator, err := stakingKeeper.GetValidator(c, valAddr)
		require.NoError(err)
		require.True(validator.Tokens.GT(previousValTokens))
	}

	tcs := []struct {
		name           string
		valDelAddr     sdk.AccAddress
		valAddr        sdk.ValAddress
		valPubKey      crypto.PubKey
		valPubKeyBytes []byte
		valUncmpPubKey []byte
		valTokens      math.Int
		moniker        string
		preRun         func(t *testing.T, c sdk.Context, valDelAddr sdk.AccAddress, valPubKey crypto.PubKey, valTokens math.Int)
		postCheck      func(c sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, delEvmAddr common.Address, previousTokens math.Int)
		expectedError  string
	}{
		{
			name:           "fail: nil validator pubkey",
			valUncmpPubKey: nil,
			expectedError:  "validator pubkey to cosmos",
		},
		{
			name:           "fail: invalid validator pubkey",
			valUncmpPubKey: uncmpPubKey0[1:],
			expectedError:  "validator pubkey to cosmos",
		},
		{
			name:           "fail: corrupted validator pubkey",
			valUncmpPubKey: corruptedPubKey,
			expectedError:  "validator pubkey to evm address",
		},
		{
			name:           "fail: mint coins",
			valDelAddr:     addrs[0],
			valUncmpPubKey: uncmpPubKey0,
			preRun: func(_ *testing.T, _ sdk.Context, valDelAddr sdk.AccAddress, _ crypto.PubKey, _ math.Int) {
				// accountKeeper.EXPECT().HasAccount(gomock.Any(), valDelAddr).Return(true)
				// bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(errors.New("mint coins"))
			},
			expectedError: "create stake coin for depositor: mint coins",
		},
		{
			name:           "fail: send coins from module to account",
			valDelAddr:     addrs[0],
			valUncmpPubKey: uncmpPubKey0,
			preRun: func(_ *testing.T, _ sdk.Context, valDelAddr sdk.AccAddress, _ crypto.PubKey, _ math.Int) {
				// accountKeeper.EXPECT().HasAccount(gomock.Any(), valDelAddr).Return(true)
				// bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				// bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, valDelAddr, gomock.Any()).Return(errors.New("send coins"))
			},
			expectedError: "create stake coin for depositor: mint coins",
		},
		{
			name:           "pass: new validator & existing delegator",
			valDelAddr:     addrs[2],
			valAddr:        valAddrs[2],
			valUncmpPubKey: uncmpPubKey2,
			valPubKey:      pubKeys[2],
			preRun: func(_ *testing.T, _ sdk.Context, valDelAddr sdk.AccAddress, _ crypto.PubKey, _ math.Int) {
				// accountKeeper.EXPECT().HasAccount(gomock.Any(), valDelAddr).Return(true)
				// bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				// bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, valDelAddr, gomock.Any()).Return(nil)
				// bankKeeper.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), valDelAddr, gomock.Any(), gomock.Any()).Return(nil)
			},
			postCheck: checkDelegatorMapAndValidator,
		},
		{
			name:           "pass: new validator & existing delegator & default moniker",
			valDelAddr:     addrs[2],
			valAddr:        valAddrs[2],
			valUncmpPubKey: uncmpPubKey2,
			valPubKey:      pubKeys[2],
			moniker:        "validator",
			preRun: func(_ *testing.T, _ sdk.Context, valDelAddr sdk.AccAddress, _ crypto.PubKey, _ math.Int) {
				// accountKeeper.EXPECT().HasAccount(gomock.Any(), valDelAddr).Return(true)
				// bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				// bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, valDelAddr, gomock.Any()).Return(nil)
				// bankKeeper.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), valDelAddr, gomock.Any(), gomock.Any()).Return(nil)
			},
			postCheck: checkDelegatorMapAndValidator,
		},
		{
			name:           "pass: new validator & new delegator",
			valDelAddr:     addrs[1],
			valAddr:        valAddrs[1],
			valUncmpPubKey: uncmpPubKey1,
			valPubKey:      pubKeys[1],
			preRun: func(_ *testing.T, _ sdk.Context, valDelAddr sdk.AccAddress, _ crypto.PubKey, _ math.Int) {
				// accountKeeper.EXPECT().HasAccount(gomock.Any(), valDelAddr).Return(false)
				// accountKeeper.EXPECT().NewAccountWithAddress(gomock.Any(), valDelAddr).Return(nil)
				// accountKeeper.EXPECT().SetAccount(gomock.Any(), gomock.Any())
				// bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				// bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, valDelAddr, gomock.Any()).Return(nil)
				//bankKeeper.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), valDelAddr, gomock.Any(), gomock.Any()).Return(nil)
			},
			postCheck: checkDelegatorMapAndValidator,
		},
		{
			name:           "pass: existing validator & delegator",
			valDelAddr:     addrs[1],
			valAddr:        valAddrs[1],
			valUncmpPubKey: uncmpPubKey1,
			valPubKey:      pubKeys[1],
			valTokens:      tokens10,
			preRun: func(t *testing.T, c sdk.Context, valDelAddr sdk.AccAddress, valPubKey crypto.PubKey, _ math.Int) {
				t.Helper()
				// create a validator with valTokens
				valAddr := sdk.ValAddress(valPubKey.Address().Bytes())
				pubKey, err := k1util.PubKeyToCosmos(valPubKey)
				require.NoError(err)
				val := testutil.NewValidator(t, valAddr, pubKey)
				validator, _, _ := val.AddTokensFromDel(tokens10, math.LegacyOneDec())
				// s.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.NotBondedPoolName, stypes.BondedPoolName, gomock.Any())
				_ = skeeper.TestingUpdateValidator(stakingKeeper, c, validator, true)

				// accountKeeper.EXPECT().HasAccount(gomock.Any(), valDelAddr).Return(true)
				// bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				// bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, valDelAddr, gomock.Any()).Return(nil)
				// bankKeeper.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), valDelAddr, gomock.Any(), gomock.Any()).Return(nil)
			},
			postCheck: checkDelegatorMapAndValTokens,
		},
		{
			name:           "pass: existing validator & new delegator",
			valDelAddr:     addrs[1],
			valAddr:        valAddrs[1],
			valUncmpPubKey: uncmpPubKey1,
			valPubKey:      pubKeys[1],
			valTokens:      tokens10,
			preRun: func(t *testing.T, c sdk.Context, valDelAddr sdk.AccAddress, valPubKey crypto.PubKey, valTokens math.Int) {
				t.Helper()
				// create a validator
				valAddr := sdk.ValAddress(valPubKey.Address().Bytes())
				pubKey, err := k1util.PubKeyToCosmos(valPubKey)
				require.NoError(err)
				val := testutil.NewValidator(t, valAddr, pubKey)
				validator, _, _ := val.AddTokensFromDel(valTokens, math.LegacyOneDec())
				// s.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), stypes.NotBondedPoolName, stypes.BondedPoolName, gomock.Any())
				_ = skeeper.TestingUpdateValidator(stakingKeeper, c, validator, true)

				// accountKeeper.EXPECT().HasAccount(gomock.Any(), valDelAddr).Return(false)
				// accountKeeper.EXPECT().NewAccountWithAddress(gomock.Any(), valDelAddr).Return(nil)
				// accountKeeper.EXPECT().SetAccount(gomock.Any(), gomock.Any())
				// bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				// bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, valDelAddr, gomock.Any()).Return(nil)
				//bankKeeper.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), valDelAddr, gomock.Any(), gomock.Any()).Return(nil)
			},
			postCheck: checkDelegatorMapAndValTokens,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cachedCtx, _ := ctx.CacheContext()
			// create valPubKey using tc.valPubKeyBytes
			if tc.preRun != nil {
				tc.preRun(s.T(), cachedCtx, tc.valDelAddr, tc.valPubKey, tc.valTokens)
			}
			moniker := tc.moniker
			if moniker == "" {
				moniker = "testing"
			}
			err := eskeeper.ProcessCreateValidator(cachedCtx, &bindings.IPTokenStakingCreateValidator{
				ValidatorUncmpPubkey:    tc.valUncmpPubKey,
				Moniker:                 moniker,
				StakeAmount:             new(big.Int).SetUint64(100),
				CommissionRate:          1000, // 10%
				MaxCommissionRate:       5000, // 50%
				MaxCommissionChangeRate: 500,  // 5%
				SupportsUnlocked:        uint8(0),
				Raw:                     gethtypes.Log{},
			})
			if tc.expectedError != "" {
				require.Error(err, tc.expectedError)
			} else {
				require.NoError(err)
				delEvmAddr, err := k1util.CosmosPubkeyToEVMAddress(tc.valPubKeyBytes)
				require.NoError(err)
				tc.postCheck(cachedCtx, tc.valDelAddr, tc.valAddr, delEvmAddr, tc.valTokens)
			}
		})
	}
}

func (s *TestSuite) TestParseCreateValidatorLog() {
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
				Topics: []common.Hash{types.CreateValidatorEvent.ID},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := keeper.ParseCreateValidatorLog(tc.log)
			if tc.expectErr {
				require.Error(err, "should return error for %s", tc.name)
			} else {
				require.NoError(err, "should not return error for %s", tc.name)
			}
		})
	}
}
