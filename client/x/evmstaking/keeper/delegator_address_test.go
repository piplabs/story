package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/k1util"
)

func (s *TestSuite) TestProcessSetWithdrawalAddress() {
	require := s.Require()
	ctx, keeper := s.Ctx, s.EVMStakingKeeper

	pubKeys, accAddrs, _ := createAddresses(2)
	delAddr := accAddrs[0]
	delPubKey := pubKeys[0]

	execAddr, err := k1util.CosmosPubkeyToEVMAddress(delPubKey.Bytes())
	require.NoError(err)
	anotherExecAddr, err := k1util.CosmosPubkeyToEVMAddress(pubKeys[1].Bytes())
	require.NoError(err)

	tcs := []struct {
		name          string
		sameAddr      bool
		input         func(execAddr common.Address, addr common.Address) *bindings.IPTokenStakingSetWithdrawalAddress
		expectedError string
	}{
		{
			name:     "pass: delegator and execution address are the same",
			sameAddr: true,
			input: func(execAddr common.Address, addr common.Address) *bindings.IPTokenStakingSetWithdrawalAddress {
				paddedExecAddr := common.LeftPadBytes(execAddr.Bytes(), 32)
				return &bindings.IPTokenStakingSetWithdrawalAddress{
					Delegator:        addr,
					ExecutionAddress: [32]byte(paddedExecAddr),
				}
			},
			expectedError: "",
		},
		{
			name:     "pass: delegator and execution address are different",
			sameAddr: false,
			input: func(execAddr common.Address, addr common.Address) *bindings.IPTokenStakingSetWithdrawalAddress {
				paddedExecAddr := common.LeftPadBytes(execAddr.Bytes(), 32)
				return &bindings.IPTokenStakingSetWithdrawalAddress{
					Delegator:        addr,
					ExecutionAddress: [32]byte(paddedExecAddr),
				}
			},
			expectedError: "",
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			evmAddr := execAddr
			if !tc.sameAddr {
				evmAddr = anotherExecAddr
			}
			ev := tc.input(evmAddr, common.Address(accAddrs[0]))

			cachedCtx, _ := ctx.CacheContext()
			err := keeper.ProcessSetWithdrawalAddress(cachedCtx, ev)
			if tc.expectedError != "" {
				require.Error(err)
				require.Contains(err.Error(), tc.expectedError)

				// Ensure no state change occurred
				addr, _ := keeper.DelegatorWithdrawAddress.Get(cachedCtx, delAddr.String())
				require.Empty(addr)
			} else {
				require.NoError(err)

				// check result
				addr, err := keeper.DelegatorWithdrawAddress.Get(cachedCtx, delAddr.String())
				require.NoError(err)
				require.Equal(evmAddr.String(), addr)
			}
		})
	}
}
