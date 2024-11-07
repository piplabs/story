package keeper_test

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmstaking/keeper"
	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
)

func (s *TestSuite) TestProcessUnjail() {
	require := s.Require()
	ctx, _, eskeeper := s.Ctx, s.SlashingKeeper, s.EVMStakingKeeper
	pubKeys, _, valAddrs := createAddresses(1)

	_ = valAddrs[0]
	valPubKey := pubKeys[0]
	valUncmpPubkey, err := keeper.CmpPubKeyToUncmpPubkey(valPubKey.Bytes())
	require.NoError(err)
	evmAddr, err := keeper.CmpPubKeyToEVMAddress(valPubKey.Bytes())
	require.NoError(err)

	tcs := []struct {
		name        string
		setupMock   func(c context.Context)
		unjailEv    *bindings.IPTokenStakingUnjail
		expectedErr string
	}{
		/*
			{
				name: "pass: valid unjail event",
				setupMock: func(c context.Context) {
					slashingKeeper.EXPECT().Unjail(c, valAddr).Return(nil)
				},
				unjailEv: &bindings.IPTokenStakingUnjail{
					ValidatorUncmpPubkey: valUncmpPubkey,
					Unjailer:             evmAddr,
				},
			},
		*/
		{
			name: "fail: invalid validator pubkey",
			unjailEv: &bindings.IPTokenStakingUnjail{
				ValidatorUncmpPubkey: valUncmpPubkey[10:],
				Unjailer:             evmAddr,
			},
			expectedErr: "invalid uncompressed public key length or format",
		},
		/*
			{
				name: "fail: validator not jailed",
				setupMock: func(c context.Context) {
					// MOCK Unjail to return error.
					slashingKeeper.EXPECT().Unjail(c, valAddr).Return(slashingtypes.ErrValidatorNotJailed)
				},
				unjailEv: &bindings.IPTokenStakingUnjail{
					ValidatorUncmpPubkey: valUncmpPubkey,
					Unjailer:             evmAddr,
				},
				expectedErr: slashingtypes.ErrValidatorNotJailed.Error(),
			},
		*/
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cachedCtx, _ := ctx.CacheContext()
			if tc.setupMock != nil {
				tc.setupMock(cachedCtx)
			}
			err := eskeeper.ProcessUnjail(cachedCtx, tc.unjailEv)
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
