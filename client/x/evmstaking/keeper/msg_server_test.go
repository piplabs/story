package keeper_test

import (
	"context"

	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
)

func (s *TestSuite) TestAddWithdrawal() {
	require := s.Require()
	ctx, msgServer, keeper := s.Ctx, s.msgServer, s.EVMStakingKeeper

	pubKeys, accAddrs, valAddrs := createAddresses(3)
	delAddr := accAddrs[0]
	delPubKey := pubKeys[0]
	valAddr := valAddrs[1]
	valPubKey := pubKeys[1]

	// self delegation
	valTokens := s.StakingKeeper.TokensFromConsensusPower(ctx, 10)

	tcs := []struct {
		name          string
		preRun        func(c context.Context)
		msg           *types.MsgAddWithdrawal
		expectedError string
	}{
		{
			name: "fail: invalid validator address",
			msg: &types.MsgAddWithdrawal{
				Withdrawal: &types.Withdrawal{
					ValidatorAddress: "invalid",
					DelegatorAddress: delAddr.String(),
					Amount:           1,
				},
			},
			expectedError: "invalid validator address",
		},
		{
			name: "fail: invalid delegator address",
			msg: &types.MsgAddWithdrawal{
				Withdrawal: &types.Withdrawal{
					ValidatorAddress: valAddr.String(),
					DelegatorAddress: "invalid",
					Amount:           1,
				},
			},
			expectedError: "invalid delegator address",
		},
		{
			name: "fail: invalid amount (should be > 0)",
			msg: &types.MsgAddWithdrawal{
				Withdrawal: &types.Withdrawal{
					ValidatorAddress: valAddr.String(),
					DelegatorAddress: delAddr.String(),
					Amount:           0,
				},
			},
			expectedError: "invalid withdrawal amount",
		},
		{
			name: "fail: unknown validator",
			msg: &types.MsgAddWithdrawal{
				Withdrawal: &types.Withdrawal{
					ValidatorAddress: valAddr.String(),
					DelegatorAddress: delAddr.String(),
					Amount:           1,
				},
			},
			expectedError: "validator not found",
		},
		{
			name: "fail: jailed validator",
			preRun: func(c context.Context) {
				s.setupValidatorAndDelegation(c, valPubKey, delPubKey, valAddr, delAddr, valTokens)
				validator, err := s.StakingKeeper.GetValidator(c, valAddr)
				require.NoError(err)
				validator.Jailed = true
				require.NoError(s.StakingKeeper.SetValidator(c, validator))
			},
			msg: &types.MsgAddWithdrawal{
				Withdrawal: &types.Withdrawal{
					ValidatorAddress: valAddr.String(),
					DelegatorAddress: delAddr.String(),
					Amount:           1,
				},
			},
			expectedError: "validator is jailed",
		},
		{
			name: "fail: unbonded validator",
			preRun: func(c context.Context) {
				s.setupValidatorAndDelegation(c, valPubKey, delPubKey, valAddr, delAddr, valTokens)
				validator, err := s.StakingKeeper.GetValidator(c, valAddr)
				require.NoError(err)
				validator.Status = stypes.Unbonded
				require.NoError(s.StakingKeeper.SetValidator(c, validator))
			},
			msg: &types.MsgAddWithdrawal{
				Withdrawal: &types.Withdrawal{
					ValidatorAddress: valAddr.String(),
					DelegatorAddress: delAddr.String(),
					Amount:           1,
				},
			},
			expectedError: "validator is unbonded",
		},
		{
			name: "pass",
			preRun: func(c context.Context) {
				s.setupValidatorAndDelegation(c, valPubKey, delPubKey, valAddr, delAddr, valTokens)
			},
			msg: &types.MsgAddWithdrawal{
				Withdrawal: &types.Withdrawal{
					ValidatorAddress: valAddr.String(),
					DelegatorAddress: delAddr.String(),
					Amount:           1,
				},
			},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			if tc.preRun != nil {
				tc.preRun(ctx)
			}
			cachedCtx, _ := ctx.CacheContext()
			resp, err := msgServer.AddWithdrawal(cachedCtx, tc.msg)
			if tc.expectedError != "" {
				require.Error(err)
				require.Contains(err.Error(), tc.expectedError)
			} else {
				require.NoError(err)
				require.NotNil(resp)
				// check withdrawal is queued
				require.Equal(uint64(1), keeper.WithdrawalQueue.Len(cachedCtx))
			}
		})
	}
}
