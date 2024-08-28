package keeper_test

import (
	"github.com/cometbft/cometbft/crypto"
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
	invalidPubKey := delPubKey.Bytes()[0:20]

	execAddr, err := k1util.CosmosPubkeyToEVMAddress(delPubKey.Bytes())
	require.NoError(err)
	anotherExecAddr, err := k1util.CosmosPubkeyToEVMAddress(pubKeys[1].Bytes())
	require.NoError(err)

	tcs := []struct {
		name          string
		sameAddr      bool
		input         func(execAddr common.Address, pubKey crypto.PubKey) *bindings.IPTokenStakingSetWithdrawalAddress
		expectedError string
	}{
		{
			name:     "pass: same delegator and execution address",
			sameAddr: true,
			input: func(execAddr common.Address, pubKey crypto.PubKey) *bindings.IPTokenStakingSetWithdrawalAddress {
				paddedExecAddr := common.LeftPadBytes(execAddr.Bytes(), 32)
				return &bindings.IPTokenStakingSetWithdrawalAddress{
					DelegatorCmpPubkey: pubKey.Bytes(),
					ExecutionAddress:   [32]byte(paddedExecAddr),
				}
			},
			expectedError: "",
		},
		{
			name:     "pass: same delegator and different execution address",
			sameAddr: false,
			input: func(execAddr common.Address, pubKey crypto.PubKey) *bindings.IPTokenStakingSetWithdrawalAddress {
				paddedExecAddr := common.LeftPadBytes(execAddr.Bytes(), 32)
				return &bindings.IPTokenStakingSetWithdrawalAddress{
					DelegatorCmpPubkey: pubKey.Bytes(),
					ExecutionAddress:   [32]byte(paddedExecAddr),
				}
			},
			expectedError: "",
		},
		{
			name:     "fail: delegator key is corrupted",
			sameAddr: false,
			input: func(execAddr common.Address, pubKey crypto.PubKey) *bindings.IPTokenStakingSetWithdrawalAddress {
				paddedExecAddr := common.LeftPadBytes(execAddr.Bytes(), 32)
				return &bindings.IPTokenStakingSetWithdrawalAddress{
					DelegatorCmpPubkey: invalidPubKey,
					ExecutionAddress:   [32]byte(paddedExecAddr),
				}
			},
			expectedError: "depositor pubkey to cosmos",
		},
	}

	for _, tc := range tcs {
		tc := tc
		s.Run(tc.name, func() {
			var ev *bindings.IPTokenStakingSetWithdrawalAddress
			evmAddr := execAddr
			if !tc.sameAddr {
			} else {
				evmAddr = anotherExecAddr
			}
			ev = tc.input(evmAddr, pubKeys[0])

			// Process the event
			err := keeper.ProcessSetWithdrawalAddress(ctx, ev)
			if tc.expectedError != "" {
				require.Error(err)
				require.Contains(err.Error(), tc.expectedError)
			} else {
				require.NoError(err)
				// check result
				addr, err := keeper.DelegatorMap.Get(ctx, delAddr.String())
				require.NoError(err)
				require.Equal(evmAddr.String(), addr)
			}
		})
	}
}
