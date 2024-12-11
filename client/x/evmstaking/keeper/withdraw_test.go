package keeper_test

import (
	"math"
	"math/big"

	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/x/evmstaking/types"
	"github.com/piplabs/story/contracts/bindings"
)

// TODO: Test for ExpectedUnbondingWithdrawals
// TODO: Test for ProcessUnbondingWithdrawals
// TODO: Test for ProcessRewardWithdrawals
// TODO: Test for ProcessEligibleRewardWithdrawal
// TODO: Test for EnqueueRewardWithdrawal

func (s *TestSuite) TestProcessWithdraw() {
	require := s.Require()
	ctx, keeper, stakingKeeper := s.Ctx, s.EVMStakingKeeper, s.StakingKeeper

	pubKeys, accAddrs, valAddrs := createAddresses(4)
	// delegator-1
	delPubKey1 := pubKeys[0]
	delAddr1 := accAddrs[0]
	// validator
	valPubKey := pubKeys[2]
	valAddr := valAddrs[2]
	// unknown pubkey
	unknownPubKey := pubKeys[3]

	valTokens := stakingKeeper.TokensFromConsensusPower(ctx, 10)
	s.setupValidatorAndDelegation(ctx, valPubKey, delPubKey1, valAddr, delAddr1, valTokens)

	singularityHeight, err := stakingKeeper.GetSingularityHeight(ctx)
	require.NoError(err)

	tcs := []struct {
		name        string
		withdraw    *bindings.IPTokenStakingWithdraw
		expectedErr string
	}{
		// TODO: test cases before and after singularity height
		{
			name: "pass: valid input",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorUncmpPubkey: cmpToUncmp(delPubKey1.Bytes()),
				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes()),
				StakeAmount:          new(big.Int).SetUint64(1),
				DelegationId:         big.NewInt(0),
				OperatorAddress:      cmpToEVM(delPubKey1.Bytes()),
			},
		},
		{
			name: "fail: invalid delegator pubkey",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorUncmpPubkey: cmpToUncmp(delPubKey1.Bytes())[:16],
				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes()),
				StakeAmount:          new(big.Int).SetUint64(1),
				DelegationId:         big.NewInt(0),
				OperatorAddress:      cmpToEVM(delPubKey1.Bytes()),
			},
			expectedErr: "invalid uncompressed public key length or format",
		},
		{
			name: "fail: invalid validator pubkey",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorUncmpPubkey: cmpToUncmp(delPubKey1.Bytes()),
				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes())[:16],
				StakeAmount:          new(big.Int).SetUint64(1),
				DelegationId:         big.NewInt(0),
				OperatorAddress:      cmpToEVM(delPubKey1.Bytes()),
			},
			expectedErr: "invalid uncompressed public key length or format",
		},
		// TODO corrupted delegator and validator pubkey
		{
			name: "fail: unknown depositor",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorUncmpPubkey: cmpToUncmp(unknownPubKey.Bytes()),
				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes()),
				StakeAmount:          new(big.Int).SetUint64(1),
				DelegationId:         big.NewInt(0),
				OperatorAddress:      cmpToEVM(unknownPubKey.Bytes()),
			},
			expectedErr: "depositor account not found",
		},
		// Intuitive behavior is to fail but instead, max share is withdrawn (100% of delegator's stake amount)
		{
			name: "pass: withdraw amount exceeds the delegation amount, results in max share withdrawal",
			withdraw: &bindings.IPTokenStakingWithdraw{
				DelegatorUncmpPubkey: cmpToUncmp(delPubKey1.Bytes()),
				ValidatorUncmpPubkey: cmpToUncmp(valPubKey.Bytes()),
				StakeAmount:          new(big.Int).SetUint64(math.MaxUint64),
				DelegationId:         big.NewInt(0),
				OperatorAddress:      cmpToEVM(delPubKey1.Bytes()),
			},
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			cachedCtx, _ := ctx.CacheContext()
			cachedCtx = cachedCtx.WithBlockHeight(int64(singularityHeight) + 1)
			require.Equal(singularityHeight+1, uint64(cachedCtx.BlockHeight()))

			// check undelegation does not exist
			_, err := s.StakingKeeper.GetUnbondingDelegation(cachedCtx, delAddr1, valAddr)
			require.ErrorContains(err, stypes.ErrNoUnbondingDelegation.Error())

			err = keeper.ProcessWithdraw(cachedCtx, tc.withdraw)
			if tc.expectedErr != "" {
				require.ErrorContains(err, tc.expectedErr)
			} else {
				require.NoError(err)
				// check undelegation exists
				ubd, err := s.StakingKeeper.GetUnbondingDelegation(cachedCtx, delAddr1, valAddr)
				require.NoError(err)
				require.NotNil(ubd)
			}
		})
	}
}

func (s *TestSuite) TestParseWithdraw() {
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
				Topics: []common.Hash{types.WithdrawEvent.ID},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := keeper.ParseWithdrawLog(tc.log)
			if tc.expectErr {
				require.Error(err, "should return error for %s", tc.name)
			} else {
				require.NoError(err, "should not return error for %s", tc.name)
			}
		})
	}
}
