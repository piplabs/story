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

type redelegateResult struct {
	srcDelegation *stypes.Delegation
	dstDelegation *stypes.Delegation
}

func TestProcessRedelegate(t *testing.T) {
	pubKeys, accAddrs, valAddrs := createAddresses(3)

	// delegator
	delPubKey := pubKeys[0]
	delAddr := accAddrs[0]

	// validators
	valSrcPubKey := pubKeys[1]
	valSrcAddr := valAddrs[1]
	valDstPubKey := pubKeys[2]
	valDstAddr := valAddrs[2]
	invalidPubKey := append([]byte{0x04}, valSrcPubKey.Bytes()[1:]...)

	createRedelegate := func(valSrcPubKey, valDstPubKey []byte, amount *big.Int) *bindings.IPTokenStakingRedelegate {
		return &bindings.IPTokenStakingRedelegate{
			Delegator:             common.Address(delAddr),
			ValidatorSrcCmpPubkey: valSrcPubKey,
			ValidatorDstCmpPubkey: valDstPubKey,
			DelegationId:          big.NewInt(0),
			OperatorAddress:       common.Address(delAddr),
			Amount:                amount,
		}
	}

	tcs := []struct {
		name              string
		setup             func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context
		setupMocks        func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper)
		mockStakingKeeper bool
		redelegate        *bindings.IPTokenStakingRedelegate
		expectedResult    redelegateResult
		expectedErr       string
	}{
		{
			name:              "fail: get singularity height from staking keeper",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), errors.New("failed to get singularity height"))
				}
			},
			redelegate:  createRedelegate(valSrcPubKey.Bytes(), valDstPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr: "check if it is singularity",
		},
		{
			name:              "pass(skip): within the singularity",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(10), nil)
				}
			},
			redelegate: createRedelegate(valSrcPubKey.Bytes(), valDstPubKey.Bytes(), new(big.Int).SetUint64(2)),
		},
		{
			name:              "fail: invalid source validator public key - length",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				}
			},
			redelegate: &bindings.IPTokenStakingRedelegate{
				Delegator:             common.Address(delAddr),
				ValidatorSrcCmpPubkey: valSrcPubKey.Bytes()[1:16],
				ValidatorDstCmpPubkey: valDstPubKey.Bytes(),
				DelegationId:          big.NewInt(0),
				OperatorAddress:       common.Address(delAddr),
				Amount:                new(big.Int).SetUint64(2),
			},
			expectedErr: "src validator pubkey to cosmos",
		},
		{
			name:              "fail: invalid source validator public key - prefix",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				}
			},
			redelegate: &bindings.IPTokenStakingRedelegate{
				Delegator:             common.Address(delAddr),
				ValidatorSrcCmpPubkey: invalidPubKey,
				ValidatorDstCmpPubkey: valDstPubKey.Bytes(),
				DelegationId:          big.NewInt(0),
				OperatorAddress:       common.Address(delAddr),
				Amount:                new(big.Int).SetUint64(2),
			},
			expectedErr: "src validator pubkey to evm address",
		},
		{
			name:              "fail: invalid destination validator public key - length",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				}
			},
			redelegate: &bindings.IPTokenStakingRedelegate{
				Delegator:             common.Address(delAddr),
				ValidatorSrcCmpPubkey: valSrcPubKey.Bytes(),
				ValidatorDstCmpPubkey: valDstPubKey.Bytes()[1:16],
				DelegationId:          big.NewInt(0),
				OperatorAddress:       common.Address(delAddr),
				Amount:                new(big.Int).SetUint64(2),
			},
			expectedErr: "dst validator pubkey to cosmos",
		},
		{
			name:              "fail: invalid destination validator public key - prefix",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				}
			},
			redelegate: &bindings.IPTokenStakingRedelegate{
				Delegator:             common.Address(delAddr),
				ValidatorSrcCmpPubkey: valSrcPubKey.Bytes(),
				ValidatorDstCmpPubkey: invalidPubKey,
				DelegationId:          big.NewInt(0),
				OperatorAddress:       common.Address(delAddr),
				Amount:                new(big.Int).SetUint64(2),
			},
			expectedErr: "dst validator pubkey to evm address",
		},
		{
			name:              "fail: different executor and operator not found",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				}
			},
			redelegate: &bindings.IPTokenStakingRedelegate{
				Delegator:             common.Address(delAddr),
				ValidatorSrcCmpPubkey: valSrcPubKey.Bytes(),
				ValidatorDstCmpPubkey: valDstPubKey.Bytes(),
				DelegationId:          big.NewInt(0),
				OperatorAddress:       cmpToEVM(delPubKey.Bytes()),
				Amount:                new(big.Int).SetUint64(2),
			},
			expectedErr: "invalid redelegateOnBehalf txn, no operator",
		},
		{
			name: "fail: different executor and not allowed operator",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				require.NoError(t, esk.DelegatorOperatorAddress.Set(ctx, delAddr.String(), cmpToEVM(delPubKey.Bytes()).String()))

				return ctx
			},
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
				}
			},
			redelegate: &bindings.IPTokenStakingRedelegate{
				Delegator:             common.Address(delAddr),
				ValidatorSrcCmpPubkey: valSrcPubKey.Bytes(),
				ValidatorDstCmpPubkey: valDstPubKey.Bytes(),
				DelegationId:          big.NewInt(0),
				OperatorAddress:       cmpToEVM(valSrcPubKey.Bytes()), // different executor
				Amount:                new(big.Int).SetUint64(2),
			},
			expectedErr: "invalid redelegateOnBehalf txn, not from operator",
		},
		{
			name:              "fail: get locked token type",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
					sk.EXPECT().GetLockedTokenType(gomock.Any()).Return(int32(0), errors.New("failed to get locked token type"))
				}
			},
			redelegate:  createRedelegate(valSrcPubKey.Bytes(), valDstPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr: "get locked token type",
		},
		{
			name:              "fail: get source validator - not found",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
					sk.EXPECT().GetLockedTokenType(gomock.Any()).Return(int32(0), nil)
					sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(stypes.Validator{}, stypes.ErrNoValidatorFound)
				}
			},
			redelegate:  createRedelegate(valSrcPubKey.Bytes(), valDstPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr: stypes.ErrNoValidatorFound.Error(),
		},
		{
			name:              "fail: get source validator",
			mockStakingKeeper: true,
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				if sk != nil {
					sk.EXPECT().GetSingularityHeight(gomock.Any()).Return(uint64(0), nil)
					sk.EXPECT().GetLockedTokenType(gomock.Any()).Return(int32(0), nil)
					sk.EXPECT().GetValidator(gomock.Any(), gomock.Any()).Return(stypes.Validator{}, errors.New("failed to get validator"))
				}
			},
			redelegate:  createRedelegate(valSrcPubKey.Bytes(), valDstPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr: "get src validator failed",
		},
		{
			name: "fail: redelegate to the same validator",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// pass singularity
				ctx = ctx.WithBlockHeight(1000000)

				// create source validator
				createValidator(t, ctx, sk, valSrcPubKey, valSrcAddr, 0)

				return ctx
			},
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			redelegate:  createRedelegate(valSrcPubKey.Bytes(), valSrcPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr: stypes.ErrSelfRedelegation.Error(),
		},
		{
			name: "fail: no destination validator",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// pass singularity
				ctx = ctx.WithBlockHeight(1000000)

				// create source validator
				createValidator(t, ctx, sk, valSrcPubKey, valSrcAddr, 0)

				return ctx
			},
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			redelegate:  createRedelegate(valSrcPubKey.Bytes(), valDstPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr: stypes.ErrNoValidatorFound.Error(),
		},
		{
			name: "fail: different support token type between source and destination validators",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// pass singularity
				ctx = ctx.WithBlockHeight(1000000)

				// create source validator
				createValidator(t, ctx, sk, valSrcPubKey, valSrcAddr, 0)

				// create destination validator
				createValidator(t, ctx, sk, valDstPubKey, valDstAddr, 1)

				return ctx
			},
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			redelegate:  createRedelegate(valSrcPubKey.Bytes(), valDstPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr: stypes.ErrTokenTypeMismatch.Error(),
		},
		{
			name: "fail: no source delegation",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// pass singularity
				ctx = ctx.WithBlockHeight(1000000)

				// create source validator
				createValidator(t, ctx, sk, valSrcPubKey, valSrcAddr, 0)

				// create destination validator
				createValidator(t, ctx, sk, valDstPubKey, valDstAddr, 0)

				return ctx
			},
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			redelegate:  createRedelegate(valSrcPubKey.Bytes(), valDstPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr: stypes.ErrNoDelegation.Error(),
		},
		{
			name: "fail: no source period delegation",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// pass singularity
				ctx = ctx.WithBlockHeight(1000000)

				// create source validator
				createValidator(t, ctx, sk, valSrcPubKey, valSrcAddr, 0)

				// create destination validator
				createValidator(t, ctx, sk, valDstPubKey, valDstAddr, 0)

				// make delegation (not period delegation) to source validator
				createDelegation(t, ctx, sk, delAddr, valSrcAddr, 4, false)

				return ctx
			},
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			redelegate:  createRedelegate(valSrcPubKey.Bytes(), valDstPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedErr: stypes.ErrNoPeriodDelegation.Error(),
		},
		{
			name: "fail: redelegation amount is less than minimum undelegation amount",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// pass singularity
				ctx = ctx.WithBlockHeight(1000000)

				// create source validator
				createValidator(t, ctx, sk, valSrcPubKey, valSrcAddr, 0)

				// create destination validator
				createValidator(t, ctx, sk, valDstPubKey, valDstAddr, 0)

				return ctx
			},
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			redelegate:  createRedelegate(valSrcPubKey.Bytes(), valDstPubKey.Bytes(), new(big.Int).SetUint64(1)),
			expectedErr: "undelegation amount is less than the minimum undelegation amount",
		},
		{
			name: "pass: redelegate",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// pass singularity
				ctx = ctx.WithBlockHeight(1000000)

				// create source validator
				createValidator(t, ctx, sk, valSrcPubKey, valSrcAddr, 0)

				// create destination validator
				createValidator(t, ctx, sk, valDstPubKey, valDstAddr, 0)

				// make delegation to source validator
				createDelegation(t, ctx, sk, delAddr, valSrcAddr, 4, true)

				return ctx
			},
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			redelegate: createRedelegate(valSrcPubKey.Bytes(), valDstPubKey.Bytes(), new(big.Int).SetUint64(2)),
			expectedResult: redelegateResult{
				srcDelegation: &stypes.Delegation{
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: valSrcAddr.String(),
					Shares:           sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2)),
					RewardsShares:    sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2)).Quo(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2))),
				},
				dstDelegation: &stypes.Delegation{
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: valDstAddr.String(),
					Shares:           sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2)),
					RewardsShares:    sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2)).Quo(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2))),
				},
			},
		},
		{
			name: "pass: try to larger amount of delegation, but pass with the total amount of delegation",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// pass singularity
				ctx = ctx.WithBlockHeight(1000000)

				// create source validator
				createValidator(t, ctx, sk, valSrcPubKey, valSrcAddr, 0)

				// create destination validator
				createValidator(t, ctx, sk, valDstPubKey, valDstAddr, 0)

				// make delegation to source validator
				createDelegation(t, ctx, sk, delAddr, valSrcAddr, 2, true)

				return ctx
			},
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			redelegate: createRedelegate(valSrcPubKey.Bytes(), valDstPubKey.Bytes(), new(big.Int).SetUint64(3)),
			expectedResult: redelegateResult{
				dstDelegation: &stypes.Delegation{
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: valDstAddr.String(),
					Shares:           sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2)),
					RewardsShares:    sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2)).Quo(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2))),
				},
			},
		},
		{
			name: "pass: redelegate from the allowed operator",
			setup: func(ctx sdk.Context, sk *skeeper.Keeper, esk *keeper.Keeper) sdk.Context {
				// pass singularity
				ctx = ctx.WithBlockHeight(1000000)

				// set operator
				require.NoError(t, esk.DelegatorOperatorAddress.Set(ctx, delAddr.String(), cmpToEVM(delPubKey.Bytes()).String()))

				// create source validator
				createValidator(t, ctx, sk, valSrcPubKey, valSrcAddr, 0)

				// create destination validator
				createValidator(t, ctx, sk, valDstPubKey, valDstAddr, 0)

				// make delegation to source validator
				createDelegation(t, ctx, sk, delAddr, valSrcAddr, 2, true)

				return ctx
			},
			setupMocks: func(ak *moduletestutil.MockAccountKeeper, bk *moduletestutil.MockBankKeeper, sk *moduletestutil.MockStakingKeeper) {
				bk.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			redelegate: &bindings.IPTokenStakingRedelegate{
				Delegator:             common.Address(delAddr),
				ValidatorSrcCmpPubkey: valSrcPubKey.Bytes(),
				ValidatorDstCmpPubkey: valDstPubKey.Bytes(),
				DelegationId:          big.NewInt(0),
				OperatorAddress:       cmpToEVM(delPubKey.Bytes()),
				Amount:                new(big.Int).SetUint64(2),
			},
			expectedResult: redelegateResult{
				dstDelegation: &stypes.Delegation{
					DelegatorAddress: delAddr.String(),
					ValidatorAddress: valDstAddr.String(),
					Shares:           sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2)),
					RewardsShares:    sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2)).Quo(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(2))),
				},
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
				tc.setupMocks(ak, bk, mockSK)
			}

			cachedCtx, _ := ctx.CacheContext()

			if tc.setup != nil {
				cachedCtx = tc.setup(cachedCtx, sk, esk)
			}

			err := esk.ProcessRedelegate(cachedCtx, tc.redelegate)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)

				// check delegations
				if tc.expectedResult.srcDelegation != nil {
					srcDelegation, err := sk.GetDelegation(cachedCtx, delAddr, valSrcAddr)
					require.NoError(t, err)
					require.Equal(t, *tc.expectedResult.srcDelegation, srcDelegation)
				}
				if tc.expectedResult.dstDelegation != nil {
					dstDelegation, err := sk.GetDelegation(cachedCtx, delAddr, valDstAddr)
					require.NoError(t, err)
					require.Equal(t, *tc.expectedResult.dstDelegation, dstDelegation)
				}
			}
		})
	}
}

func TestParseRedelegateLog(t *testing.T) {
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
				Topics: []common.Hash{types.RedelegateEvent.ID},
			},
			expectErr: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			_, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
			_, err := esk.ParseRedelegateLog(tc.log)
			if tc.expectErr {
				require.Error(t, err, "should return error for %s", tc.name)
			} else {
				require.NoError(t, err, "should not return error for %s", tc.name)
			}
		})
	}
}
