package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	skeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/evmstaking/keeper"
	moduletestutil "github.com/piplabs/story/client/x/evmstaking/testutil"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"

	"go.uber.org/mock/gomock"
)

func TestProcessUpdateValidatorCommission(t *testing.T) {
	pubKeys, _, valAddrs := createAddresses(1)

	valAddr := valAddrs[0]
	valPubKey := pubKeys[0]

	invalidPubKey := append([]byte{0x04}, valPubKey.Bytes()[1:]...)

	createUpdateValidatorCommision := func(newCommissionRate uint32) *bindings.IPTokenStakingUpdateValidatorCommission {
		return &bindings.IPTokenStakingUpdateValidatorCommission{
			ValidatorCmpPubkey: valPubKey.Bytes(),
			CommissionRate:     newCommissionRate,
		}
	}

	tcs := []struct {
		name                      string
		mockStakingKeeper         bool
		setupMock                 func(bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper)
		beforeTest                func(ctx sdk.Context, sk *skeeper.Keeper)
		setValidator              bool
		updateValidatorCommission *bindings.IPTokenStakingUpdateValidatorCommission
		expectedErr               string
		expectedResult            stypes.CommissionRates
	}{
		{
			name:              "fail: invalid validator public key - length",
			mockStakingKeeper: true,
			updateValidatorCommission: &bindings.IPTokenStakingUpdateValidatorCommission{
				ValidatorCmpPubkey: valPubKey.Bytes()[1:10],
			},
			expectedErr: "validator pubkey to cosmos",
		},
		{
			name:              "fail: invalid validator public key - prefix",
			mockStakingKeeper: true,
			updateValidatorCommission: &bindings.IPTokenStakingUpdateValidatorCommission{
				ValidatorCmpPubkey: invalidPubKey,
			},
			expectedErr: "validator pubkey to evm address",
		},
		{
			name:              "fail: get validator - not found",
			mockStakingKeeper: true,
			setupMock: func(bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(stypes.Validator{}, stypes.ErrNoValidatorFound)
			},
			updateValidatorCommission: createUpdateValidatorCommision(1),
			expectedErr:               stypes.ErrNoValidatorFound.Error(),
		},
		{
			name:              "fail: get validator - unknown error",
			mockStakingKeeper: true,
			setupMock: func(bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(stypes.Validator{}, errors.New("unknown error"))
			},
			updateValidatorCommission: createUpdateValidatorCommision(1),
			expectedErr:               "get validator: unknown error",
		},
		{
			name:              "fail: type assertion of staking keeper",
			mockStakingKeeper: true,
			setupMock: func(bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(stypes.Validator{
					MinSelfDelegation: math.NewInt(1),
				}, nil)
			},
			updateValidatorCommission: createUpdateValidatorCommision(1),
			expectedErr:               "type assertion failed",
		},
		{
			name: "fail: edit validator - greater commission rate than one",
			setupMock: func(bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setValidator:              true,
			updateValidatorCommission: createUpdateValidatorCommision(10001),
			expectedErr:               "commission rate must be between 0 and 1 (inclusive)",
		},
		{
			name: "fail: edit validator - greater commission rate than max rate",
			setupMock: func(bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setValidator:              true,
			updateValidatorCommission: createUpdateValidatorCommision(6000),
			expectedErr:               stypes.ErrCommissionGTMaxRate.Error(),
		},
		{
			name: "fail: edit validator - greater commission rate than max change rate",
			setupMock: func(bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setValidator:              true,
			updateValidatorCommission: createUpdateValidatorCommision(3000),
			expectedErr:               stypes.ErrCommissionGTMaxChangeRate.Error(),
		},
		{
			name: "fail: edit validator - try to change again within 24 hours after the previous change",
			setupMock: func(bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			beforeTest: func(ctx sdk.Context, sk *skeeper.Keeper) {
				validator, err := sk.GetValidator(ctx, valAddr)
				require.NoError(t, err)

				newCommission, err := sk.UpdateValidatorCommission(ctx, validator, math.LegacyNewDecWithPrec(1200, 4))
				require.NoError(t, err)

				validator.Commission = newCommission
				require.NoError(t, sk.SetValidator(ctx, validator))
			},
			setValidator:              true,
			updateValidatorCommission: createUpdateValidatorCommision(1500),
			expectedErr:               stypes.ErrCommissionUpdateTime.Error(),
		},
		{
			name: "pass: update commission rate",
			setupMock: func(bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			setValidator:              true,
			updateValidatorCommission: createUpdateValidatorCommision(1500),
			expectedResult: stypes.CommissionRates{
				Rate:          math.LegacyNewDecWithPrec(1500, 4),
				MaxRate:       math.LegacyNewDecWithPrec(5000, 4),
				MaxChangeRate: math.LegacyNewDecWithPrec(1000, 4),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var (
				ctx    sdk.Context
				bk     *moduletestutil.MockBankKeeper
				sk     *skeeper.Keeper
				mockSK *moduletestutil.MockStakingKeeper
				esk    *keeper.Keeper
			)

			if tc.mockStakingKeeper {
				ctx, _, bk, _, mockSK, _, _, esk = createKeeperWithMockStaking(t)
			} else {
				ctx, _, bk, sk, _, esk = createKeeperWithRealStaking(t)
			}

			if tc.setupMock != nil {
				tc.setupMock(bk, mockSK)
			}

			cachedCtx, _ := ctx.CacheContext()

			if tc.setValidator {
				createValidator(t, cachedCtx, sk, valPubKey, valAddr, 0)
			}

			if tc.beforeTest != nil {
				tc.beforeTest(cachedCtx, sk)
			}

			err := esk.ProcessUpdateValidatorCommission(cachedCtx, tc.updateValidatorCommission)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)

				validator, err := sk.GetValidator(cachedCtx, valAddr)
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, validator.Commission.CommissionRates)
				require.Equal(t, cachedCtx.BlockTime(), validator.Commission.UpdateTime)
			}
		})
	}
}
