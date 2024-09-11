package keeper_test

import (
	"context"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
)

func (s *TestSuite) TestProcessUnjail() {
	require := s.Require()
	ctx, slashingKeeper, keeper := s.Ctx, s.SlashingKeeper, s.EVMStakingKeeper
	pubKeys, _, valAddrs := createAddresses(1)
	valAddr := valAddrs[0]
	valPubKey := pubKeys[0]

	tcs := []struct {
		name        string
		setupMock   func(c context.Context)
		unjailEv    *bindings.IPTokenSlashingUnjail
		expectedErr string
	}{
		{
			name: "pass: valid unjail event",
			setupMock: func(c context.Context) {
				slashingKeeper.EXPECT().Unjail(c, valAddr).Return(nil)
			},
			unjailEv: &bindings.IPTokenSlashingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes(),
			},
		},
		{
			name: "fail: invalid validator pubkey",
			unjailEv: &bindings.IPTokenSlashingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes()[10:],
			},
			expectedErr: "validator pubkey to cosmos: invalid pubkey length",
		},
		{
			name: "fail: validator not jailed",
			setupMock: func(c context.Context) {
				// MOCK Unjail to return error.
				slashingKeeper.EXPECT().Unjail(c, valAddr).Return(slashingtypes.ErrValidatorNotJailed)
			},
			unjailEv: &bindings.IPTokenSlashingUnjail{
				ValidatorCmpPubkey: valPubKey.Bytes(),
			},
			expectedErr: slashingtypes.ErrValidatorNotJailed.Error(),
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cachedCtx, _ := ctx.CacheContext()
			if tc.setupMock != nil {
				tc.setupMock(cachedCtx)
			}
			err := keeper.ProcessUnjail(cachedCtx, tc.unjailEv)
			if tc.expectedErr != "" {
				require.ErrorContains(err, tc.expectedErr)
			} else {
				require.NoError(err)
			}
		})
	}
}

func (s *TestSuite) TestParseUnjailLog() {
	require := s.Require()
	keeper := s.EVMStakingKeeper

	dummyEthAddr := common.HexToAddress("0x1")

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
				Topics: []common.Hash{types.UnjailEvent.ID, common.BytesToHash(dummyEthAddr.Bytes())},
			},
			expectErr: false,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			_, err := keeper.ParseUnjailLog(tc.log)
			if tc.expectErr {
				require.Error(err, "should return error for %s", tc.name)
			} else {
				require.NoError(err, "should not return error for %s", tc.name)
			}
		})
	}
}
