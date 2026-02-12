package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmstaking/keeper"
	moduletestutil "github.com/piplabs/story/client/x/evmstaking/testutil"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/k1util"

	"go.uber.org/mock/gomock"
)

func TestProcessCreateValidator(t *testing.T) {
	pubKeys, _, valAddrs := createAddresses(1)

	// validator
	valPubKey := pubKeys[0]
	valAddr := valAddrs[0]
	valEVMAddr, err := k1util.CosmosPubkeyToEVMAddress(valPubKey.Bytes())
	require.NoError(t, err)

	invalidPubKey := append([]byte{0x04}, valPubKey.Bytes()[1:]...)

	createValidatorEv := func(amount uint64, rate, maxRate, maxChangeRate uint32, supportedTokenType uint8) *bindings.IPTokenStakingCreateValidator {
		return &bindings.IPTokenStakingCreateValidator{
			ValidatorCmpPubkey:      valPubKey.Bytes(),
			Moniker:                 "",
			StakeAmount:             new(big.Int).SetUint64(amount),
			CommissionRate:          rate,
			MaxCommissionRate:       maxRate,
			MaxCommissionChangeRate: maxChangeRate,
			SupportsUnlocked:        supportedTokenType,
			OperatorAddress:         cmpToEVM(valPubKey.Bytes()),
		}
	}

	tcs := []struct {
		name              string
		setupMock         func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper)
		setup             func(ctx sdk.Context, sk *skeeper.Keeper) sdk.Context
		mockStakingKeeper bool
		createValidator   *bindings.IPTokenStakingCreateValidator
		expectedError     string
		expectedAmount    uint64
	}{
		{
			name: "fail: invalid validator public key - length",
			createValidator: &bindings.IPTokenStakingCreateValidator{
				ValidatorCmpPubkey: valPubKey.Bytes()[1:16],
			},
			expectedError: "validator pubkey to cosmos",
		},
		{
			name: "fail: invalid validator public key - prefix",
			createValidator: &bindings.IPTokenStakingCreateValidator{
				ValidatorCmpPubkey: invalidPubKey,
			},
			expectedError: "validator pubkey to evm address",
		},
		{
			name: "fail: type assertion to staking keeper fail",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
			},
			mockStakingKeeper: true,
			createValidator:   createValidatorEv(100, 1000, 5000, 500, 0),
			expectedError:     "type assertion failed",
		},
		{
			name: "fail: mint coins to evmstaking module",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to mint coins"))
			},
			createValidator: createValidatorEv(100, 1000, 5000, 500, 0),
			expectedError:   "create stake coin for depositor: mint coins",
		},
		{
			name: "fail: send coins from evmstaking module to account",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to send coins"))
			},
			createValidator: createValidatorEv(100, 1000, 5000, 500, 0),
			expectedError:   "create stake coin for depositor: send coins",
		},
		{
			name: "fail: smaller amount of self delegation than minimum self delegation",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper) sdk.Context {
				ctx = ctx.WithBlockHeight(1)

				return ctx
			},
			createValidator: createValidatorEv(1, 1000, 5000, 500, 0),
			expectedError:   "validator's self delegation must be greater than their minimum self delegation",
		},
		{
			name: "pass: smaller amount of delegation than minimum self delegation, but genesis validator",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			createValidator: createValidatorEv(1, 1000, 5000, 500, 0),
			expectedAmount:  1,
		},
		{
			name: "fail: validator already exist",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper) sdk.Context {
				createValidator(t, ctx, sk, valPubKey, valAddr, 0)

				return ctx
			},
			createValidator: createValidatorEv(100, 1000, 5000, 500, 0),
			expectedError:   stypes.ErrValidatorOwnerExists.Error(),
		},
		{
			name: "fail: invalid commission rate - less commission rate than minimum commission rate",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setup: func(ctx sdk.Context, sk *skeeper.Keeper) sdk.Context {
				newParams := stypes.DefaultParams()
				newParams.MinCommissionRate = math.LegacyNewDecWithPrec(100, 4)
				require.NoError(t, sk.SetParams(ctx, newParams))

				return ctx
			},
			createValidator: createValidatorEv(100, 1, 5000, 500, 0),
			expectedError:   stypes.ErrCommissionLTMinRate.Error(),
		},
		{
			name: "fail: invalid commission rate - greater max rate than one",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			createValidator: createValidatorEv(100, 1000, 10001, 500, 0),
			expectedError:   stypes.ErrCommissionHuge.Error(),
		},
		{
			name: "fail: invalid commission rate - greater commission rate than max rate",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			createValidator: createValidatorEv(100, 6000, 5000, 500, 0),
			expectedError:   stypes.ErrCommissionGTMaxRate.Error(),
		},
		{
			name: "fail: invalid commission rate - greater max change rate than max rate",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			createValidator: createValidatorEv(100, 1000, 5000, 5001, 0),
			expectedError:   stypes.ErrCommissionChangeRateGTMaxRate.Error(),
		},
		{
			name: "fail: invalid supported token type",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			createValidator: createValidatorEv(100, 1000, 5000, 500, 2),
			expectedError:   stypes.ErrNoTokenTypeFound.Error(),
		},
		{
			name: "pass: new validator with existing account",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			createValidator: createValidatorEv(100, 1000, 5000, 500, 0),
			expectedAmount:  100,
		},
		{
			name: "pass: new validator with new account",
			setupMock: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(false)
				ak.EXPECT().NewAccountWithAddress(gomock.Any(), gomock.Any()).Return(nil)
				ak.EXPECT().SetAccount(gomock.Any(), gomock.Any()).Return()
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			createValidator: createValidatorEv(100, 1000, 5000, 500, 0),
			expectedAmount:  100,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var (
				ctx sdk.Context
				ak  *moduletestutil.MockAccountKeeper
				bk  *moduletestutil.MockBankKeeper
				sk  *skeeper.Keeper
				esk *keeper.Keeper
			)

			if tc.mockStakingKeeper {
				ctx, ak, bk, _, _, _, _, esk = createKeeperWithMockStaking(t)
			} else {
				ctx, ak, bk, sk, _, esk = createKeeperWithRealStaking(t)
			}

			if tc.setupMock != nil {
				tc.setupMock(ak, bk)
			}

			if tc.setup != nil {
				ctx = tc.setup(ctx, sk)
			}

			cachedCtx, _ := ctx.CacheContext()

			err := esk.ProcessCreateValidator(cachedCtx, tc.createValidator)
			if tc.expectedError != "" {
				require.ErrorContains(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				validator, err := sk.Validator(cachedCtx, valAddr)
				require.NoError(t, err)
				require.NotNil(t, validator)
				require.Equal(t, tc.expectedAmount, validator.GetTokens().Uint64())

				withdrawAddr, err := esk.DelegatorWithdrawAddress.Get(cachedCtx, sdk.AccAddress(valEVMAddr.Bytes()).String())
				require.NoError(t, err)
				rewardAddr, err := esk.DelegatorRewardAddress.Get(cachedCtx, sdk.AccAddress(valEVMAddr.Bytes()).String())
				require.NoError(t, err)
				require.Equal(t, withdrawAddr, rewardAddr)
			}
		})
	}
}

func TestParseCreateValidatorLog(t *testing.T) {
	tcs := []struct {
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

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			_, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

			_, err := esk.ParseCreateValidatorLog(tc.log)
			if tc.expectErr {
				require.Error(t, err, "should return error for %s", tc.name)
			} else {
				require.NoError(t, err, "should not return error for %s", tc.name)
			}
		})
	}
}
