package keeper_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"

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

	"go.uber.org/mock/gomock"
)

func TestProcessDeposit(t *testing.T) {
	pubKeys, accAddrs, valAddrs := createAddresses(2)

	// delegator
	delPubKey := pubKeys[0]
	delAddr := accAddrs[0]

	// validator
	valPubKey := pubKeys[1]
	valAddr := valAddrs[1]
	invalidPubKey := append([]byte{0x04}, valPubKey.Bytes()[1:]...)

	createDeposit := func(delPubKey, valPubKey []byte, amount *big.Int) *bindings.IPTokenStakingDeposit {
		return &bindings.IPTokenStakingDeposit{
			Delegator:          common.Address(delAddr),
			ValidatorCmpPubkey: valPubKey,
			StakeAmount:        amount,
			StakingPeriod:      big.NewInt(0),
			DelegationId:       big.NewInt(0),
			OperatorAddress:    cmpToEVM(delPubKey),
		}
	}

	tcs := []struct {
		name       string
		setupMocks func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper)
		// mockStakingKeeper determines whether to mock the staking keeper or use the real implementation
		mockStakingKeeper bool
		// if mockStakingKeeper is true, setup mock functions for mocked staking keeper
		setupStakingKeeperMock func(sk *moduletestutil.MockStakingKeeper)
		deposit                *bindings.IPTokenStakingDeposit
		validatorNotExist      bool
		expectedResult         stypes.Delegation
		expectedErr            string
	}{
		{
			name: "fail: invalid length of validator pubkey",
			deposit: &bindings.IPTokenStakingDeposit{
				Delegator:          common.Address(delAddr),
				ValidatorCmpPubkey: valPubKey.Bytes()[:16],
				StakeAmount:        new(big.Int).SetUint64(1),
				StakingPeriod:      big.NewInt(0),
				DelegationId:       big.NewInt(0),
				OperatorAddress:    cmpToEVM(delPubKey.Bytes()),
			},
			validatorNotExist: true,
			expectedErr:       "validator pubkey to cosmos: invalid pubkey length",
		},
		{
			name: "fail: invalid compressed validator pubkey",
			deposit: &bindings.IPTokenStakingDeposit{
				Delegator:          common.Address(delAddr),
				ValidatorCmpPubkey: invalidPubKey,
				StakeAmount:        new(big.Int).SetUint64(1),
				StakingPeriod:      big.NewInt(0),
				DelegationId:       big.NewInt(0),
				OperatorAddress:    cmpToEVM(delPubKey.Bytes()),
			},
			validatorNotExist: true,
			expectedErr:       "validator pubkey to evm address",
		},
		{
			name: "fail: get locked token type",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
			},
			mockStakingKeeper: true,
			setupStakingKeeperMock: func(sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetLockedTokenType(gomock.Any()).Return(int32(0), errors.New("failed to get locked token type"))
			},
			deposit:           createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(2)),
			validatorNotExist: true,
			expectedErr:       "get locked token type",
		},
		{
			name: "fail: validator not found",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
			},
			validatorNotExist: true,
			deposit:           createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr:       "validator does not exist",
		},
		{
			name: "fail: get validator error",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
			},
			mockStakingKeeper: true,
			setupStakingKeeperMock: func(sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetLockedTokenType(gomock.Any()).Return(int32(0), nil)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(stypes.Validator{}, errors.New("failed to get validator"))
			},
			validatorNotExist: true,
			deposit:           createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr:       "get validator failed",
		},
		{
			name: "fail: get flexible period type",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
			},
			mockStakingKeeper: true,
			setupStakingKeeperMock: func(sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetLockedTokenType(gomock.Any()).Return(int32(0), nil)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(stypes.Validator{
					SupportTokenType: int32(0),
				}, nil)
				sk.EXPECT().GetFlexiblePeriodType(gomock.Any()).Return(int32(0), errors.New("failed to get flexible period type"))
			},
			validatorNotExist: true,
			deposit:           createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr:       "get flexible period type",
		},
		{
			name: "fail: type assertion to staking keeper fail",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
			},
			mockStakingKeeper: true,
			setupStakingKeeperMock: func(sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetLockedTokenType(gomock.Any()).Return(int32(0), nil)
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(stypes.Validator{
					SupportTokenType: int32(0),
				}, nil)
				sk.EXPECT().GetFlexiblePeriodType(gomock.Any()).Return(int32(0), nil)
			},
			validatorNotExist: true,
			deposit:           createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr:       "type assertion failed",
		},
		{
			name: "fail: create stake coin - failed to mint coins",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to mint coins"))
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			deposit:     createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr: "create stake coin for depositor: mint coins",
		},
		{
			name: "fail: create stake coin - failed to send coins from module to account",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to send coins from module to account"))
			},
			deposit:     createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr: "create stake coin for depositor: send coins",
		},
		{
			name: "fail: delegate - less than min amount of delegation",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			deposit:     createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(1)),
			expectedErr: "delegation amount must be greater than or equal to minimum delegation",
		},
		{
			name: "fail: delegate - failed to delegate coins from account to module",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to delegate"))
			},
			deposit:     createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr: "delegate",
		},
		{
			name: "pass: existing delegator",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(true)
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			deposit: createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedResult: stypes.Delegation{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Shares:           sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2)),
				RewardsShares:    sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2)).Quo(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2))),
			},
		},
		{
			name: "pass: new delegator",
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper) {
				ak.EXPECT().HasAccount(gomock.Any(), gomock.Any()).Return(false)
				ak.EXPECT().NewAccountWithAddress(gomock.Any(), gomock.Any()).Return(nil)
				ak.EXPECT().SetAccount(gomock.Any(), gomock.Any()).Return()
				bk.EXPECT().MintCoins(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				bk.EXPECT().DelegateCoinsFromAccountToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			deposit: createDeposit(delPubKey.Bytes(), valPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedResult: stypes.Delegation{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Shares:           sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2)),
				RewardsShares:    sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2)).Quo(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2))),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var (
				ctx    sdk.Context
				ak     *moduletestutil.MockAccountKeeper
				bk     *moduletestutil.MockBankKeeper
				sk     *skeeper.Keeper
				mockSK *moduletestutil.MockStakingKeeper
				esk    *keeper.Keeper
			)

			if tc.mockStakingKeeper {
				ctx, ak, bk, _, mockSK, _, _, esk = createKeeperWithMockStaking(t)
			} else {
				ctx, ak, bk, sk, _, esk = createKeeperWithRealStaking(t)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(ak, bk)
			}

			if mockSK != nil && tc.setupStakingKeeperMock != nil {
				tc.setupStakingKeeperMock(mockSK)
			}

			cachedCtx, _ := ctx.CacheContext()

			// create a validator
			if !tc.validatorNotExist {
				createValidator(t, cachedCtx, sk, valPubKey, valAddr, 0)
			}

			err := esk.ProcessDeposit(cachedCtx, tc.deposit)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				// check delegation
				delegation, err := sk.GetDelegation(cachedCtx, delAddr, valAddr)
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, delegation)
			}
		})
	}
}

func TestParseDepositLog(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			_, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
			_, err := esk.ParseDepositLog(tc.log)
			if tc.expectErr {
				require.Error(t, err, "should return error for %s", tc.name)
			} else {
				require.NoError(t, err, "should not return error for %s", tc.name)
			}
		})
	}
}
