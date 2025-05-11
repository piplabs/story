package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
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

func TestProcessUnjail(t *testing.T) {
	pubKeys, _, _ := createAddresses(1)

	valPubKey := pubKeys[0]
	evmAddr, err := keeper.CmpPubKeyToEVMAddress(valPubKey.Bytes())
	require.NoError(t, err)

	invalidPubKey := append([]byte{0x04}, valPubKey.Bytes()[1:]...)

	tcs := []struct {
		name        string
		setup       func(ctx sdk.Context, esk *keeper.Keeper)
		setupMock   func(c sdk.Context, slk *moduletestutil.MockSlashingKeeper)
		unjailEv    *bindings.IPTokenStakingUnjail
		expectedErr string
	}{
		{
			name: "fail: invalid validator pubkey - length",
			unjailEv: &bindings.IPTokenStakingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes()[10:],
				Unjailer:           evmAddr,
			},
			expectedErr: "validator pubkey to cosmos: invalid pubkey length",
		},
		{
			name: "fail: invalid validator pubkey - prefix",
			unjailEv: &bindings.IPTokenStakingUnjail{
				ValidatorCmpPubkey: invalidPubKey,
				Unjailer:           evmAddr,
			},
			expectedErr: "validator pubkey to evm address",
		},
		{
			name: "fail: different executor and operator not found",
			unjailEv: &bindings.IPTokenStakingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Unjailer:           common.MaxAddress,
			},
			expectedErr: "invalid unjailOnBehalf txn, no operator for delegator",
		},
		{
			name: "fail: different executor and not allowed operator",
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorOperatorAddress.Set(ctx, sdk.AccAddress(evmAddr.Bytes()).String(), evmAddr.String()))
			},
			unjailEv: &bindings.IPTokenStakingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Unjailer:           common.MaxAddress,
			},
			expectedErr: "invalid unjailOnBehalf txn, not from operator",
		},
		{
			name: "fail: unjail failure - validator not found",
			setupMock: func(c sdk.Context, slk *moduletestutil.MockSlashingKeeper) {
				slk.EXPECT().Unjail(gomock.Any(), gomock.Any()).Return(slashtypes.ErrNoValidatorForAddress)
			},
			unjailEv: &bindings.IPTokenStakingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Unjailer:           evmAddr,
			},
			expectedErr: slashtypes.ErrNoValidatorForAddress.Error(),
		},
		{
			name: "fail: unjail failure - missing self delegation",
			setupMock: func(c sdk.Context, slk *moduletestutil.MockSlashingKeeper) {
				slk.EXPECT().Unjail(gomock.Any(), gomock.Any()).Return(slashtypes.ErrMissingSelfDelegation)
			},
			unjailEv: &bindings.IPTokenStakingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Unjailer:           evmAddr,
			},
			expectedErr: slashtypes.ErrMissingSelfDelegation.Error(),
		},
		{
			name: "fail: unjail failure - validator not jailed",
			setupMock: func(c sdk.Context, slk *moduletestutil.MockSlashingKeeper) {
				slk.EXPECT().Unjail(gomock.Any(), gomock.Any()).Return(slashtypes.ErrValidatorNotJailed)
			},
			unjailEv: &bindings.IPTokenStakingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Unjailer:           evmAddr,
			},
			expectedErr: slashtypes.ErrValidatorNotJailed.Error(),
		},
		{
			name: "fail: unjail failure - validator still jailed",
			setupMock: func(c sdk.Context, slk *moduletestutil.MockSlashingKeeper) {
				slk.EXPECT().Unjail(gomock.Any(), gomock.Any()).Return(slashtypes.ErrValidatorJailed)
			},
			unjailEv: &bindings.IPTokenStakingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Unjailer:           evmAddr,
			},
			expectedErr: slashtypes.ErrValidatorJailed.Error(),
		},
		{
			name: "fail: unjail failure - unknown error",
			setupMock: func(c sdk.Context, slk *moduletestutil.MockSlashingKeeper) {
				slk.EXPECT().Unjail(gomock.Any(), gomock.Any()).Return(errors.New("unknown unjail error"))
			},
			unjailEv: &bindings.IPTokenStakingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Unjailer:           evmAddr,
			},
			expectedErr: "unjail: unknown unjail error",
		},
		{
			name: "pass: unjail validator from the same executor",
			setupMock: func(c sdk.Context, slk *moduletestutil.MockSlashingKeeper) {
				slk.EXPECT().Unjail(gomock.Any(), gomock.Any()).Return(nil)
			},
			unjailEv: &bindings.IPTokenStakingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Unjailer:           evmAddr,
			},
		},
		{
			name: "pass: unjail validator from allowed operator",
			setup: func(ctx sdk.Context, esk *keeper.Keeper) {
				require.NoError(t, esk.DelegatorOperatorAddress.Set(ctx, sdk.AccAddress(evmAddr.Bytes()).String(), common.MaxAddress.String()))
			},
			setupMock: func(c sdk.Context, slk *moduletestutil.MockSlashingKeeper) {
				slk.EXPECT().Unjail(gomock.Any(), gomock.Any()).Return(nil)
			},
			unjailEv: &bindings.IPTokenStakingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes(),
				Unjailer:           common.MaxAddress,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, _, slk, esk := createKeeperWithMockStaking(t)

			cachedCtx, _ := ctx.CacheContext()

			if tc.setupMock != nil {
				tc.setupMock(cachedCtx, slk)
			}

			if tc.setup != nil {
				tc.setup(cachedCtx, esk)
			}

			err := esk.ProcessUnjail(cachedCtx, tc.unjailEv)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParseUnjailLog(t *testing.T) {
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
				Topics: []common.Hash{types.UnjailEvent.ID},
			},
			expectErr: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			_, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

			_, err := esk.ParseUnjailLog(tc.log)
			if tc.expectErr {
				require.Error(t, err, "should return error for %s", tc.name)
			} else {
				require.NoError(t, err, "should not return error for %s", tc.name)
			}
		})
	}
}
