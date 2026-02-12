package keeper_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/k1util"
)

func TestProcessSetWithdrawalAddress(t *testing.T) {
	pubKeys, accAddrs, _ := createAddresses(2)
	delAddr := accAddrs[0]
	delPubKey := pubKeys[0]

	execAddr, err := k1util.CosmosPubkeyToEVMAddress(delPubKey.Bytes())
	require.NoError(t, err)
	anotherExecAddr, err := k1util.CosmosPubkeyToEVMAddress(pubKeys[1].Bytes())
	require.NoError(t, err)

	tcs := []struct {
		name     string
		sameAddr bool
		input    func(execAddr common.Address, addr common.Address) *bindings.IPTokenStakingSetWithdrawalAddress
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
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
			evmAddr := execAddr
			if !tc.sameAddr {
				evmAddr = anotherExecAddr
			}
			ev := tc.input(evmAddr, common.Address(accAddrs[0]))

			cachedCtx, _ := ctx.CacheContext()
			err := esk.ProcessSetWithdrawalAddress(cachedCtx, ev)
			require.NoError(t, err)

			// check result
			addr, err := esk.DelegatorWithdrawAddress.Get(cachedCtx, delAddr.String())
			require.NoError(t, err)
			require.Equal(t, evmAddr.String(), addr)
		})
	}
}

func TestProcessSetRewardAddress(t *testing.T) {
	pubKeys, accAddrs, _ := createAddresses(2)
	delAddr := accAddrs[0]
	delPubKey := pubKeys[0]

	execAddr, err := k1util.CosmosPubkeyToEVMAddress(delPubKey.Bytes())
	require.NoError(t, err)
	anotherExecAddr, err := k1util.CosmosPubkeyToEVMAddress(pubKeys[1].Bytes())
	require.NoError(t, err)

	tcs := []struct {
		name     string
		sameAddr bool
		input    func(execAddr common.Address, addr common.Address) *bindings.IPTokenStakingSetRewardAddress
	}{
		{
			name:     "pass: delegator and execution address are the same",
			sameAddr: true,
			input: func(execAddr common.Address, addr common.Address) *bindings.IPTokenStakingSetRewardAddress {
				paddedExecAddr := common.LeftPadBytes(execAddr.Bytes(), 32)
				return &bindings.IPTokenStakingSetRewardAddress{
					Delegator:        addr,
					ExecutionAddress: [32]byte(paddedExecAddr),
				}
			},
		},
		{
			name:     "pass: delegator and execution address are different",
			sameAddr: false,
			input: func(execAddr common.Address, addr common.Address) *bindings.IPTokenStakingSetRewardAddress {
				paddedExecAddr := common.LeftPadBytes(execAddr.Bytes(), 32)
				return &bindings.IPTokenStakingSetRewardAddress{
					Delegator:        addr,
					ExecutionAddress: [32]byte(paddedExecAddr),
				}
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)
			evmAddr := execAddr
			if !tc.sameAddr {
				evmAddr = anotherExecAddr
			}
			ev := tc.input(evmAddr, common.Address(accAddrs[0]))

			cachedCtx, _ := ctx.CacheContext()
			err := esk.ProcessSetRewardAddress(cachedCtx, ev)
			require.NoError(t, err)

			// check result
			addr, err := esk.DelegatorRewardAddress.Get(cachedCtx, delAddr.String())
			require.NoError(t, err)
			require.Equal(t, evmAddr.String(), addr)
		})
	}
}

func TestProcessSetOperatorAddress(t *testing.T) {
	pubKeys, accAddrs, _ := createAddresses(2)
	delAddr := accAddrs[0]
	delPubKey := pubKeys[0]

	execAddr, err := k1util.CosmosPubkeyToEVMAddress(delPubKey.Bytes())
	require.NoError(t, err)
	anotherExecAddr, err := k1util.CosmosPubkeyToEVMAddress(pubKeys[1].Bytes())
	require.NoError(t, err)

	tcs := []struct {
		name     string
		sameAddr bool
		input    func(execAddr common.Address, addr common.Address) *bindings.IPTokenStakingSetOperator
	}{
		{
			name:     "pass: delegator and execution address are the same",
			sameAddr: true,
			input: func(execAddr common.Address, addr common.Address) *bindings.IPTokenStakingSetOperator {
				return &bindings.IPTokenStakingSetOperator{
					Delegator: addr,
					Operator:  execAddr,
				}
			},
		},
		{
			name:     "pass: delegator and execution address are different",
			sameAddr: false,
			input: func(execAddr common.Address, addr common.Address) *bindings.IPTokenStakingSetOperator {
				return &bindings.IPTokenStakingSetOperator{
					Delegator: addr,
					Operator:  execAddr,
				}
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

			evmAddr := execAddr
			if !tc.sameAddr {
				evmAddr = anotherExecAddr
			}
			ev := tc.input(evmAddr, common.Address(accAddrs[0]))

			cachedCtx, _ := ctx.CacheContext()
			err := esk.ProcessSetOperator(cachedCtx, ev)
			require.NoError(t, err)

			// check result
			addr, err := esk.DelegatorOperatorAddress.Get(cachedCtx, delAddr.String())
			require.NoError(t, err)
			require.Equal(t, evmAddr.String(), addr)
		})
	}
}

func TestProcessUnsetOperatorAddress(t *testing.T) {
	pubKeys, accAddrs, _ := createAddresses(2)
	delAddr := accAddrs[0]

	anotherExecAddr, err := k1util.CosmosPubkeyToEVMAddress(pubKeys[1].Bytes())
	require.NoError(t, err)

	tcs := []struct {
		name     string
		existing bool
		input    func(execAddr common.Address) *bindings.IPTokenStakingUnsetOperator
	}{
		{
			name:     "pass: operator set previously",
			existing: true,
			input: func(delAddr common.Address) *bindings.IPTokenStakingUnsetOperator {
				return &bindings.IPTokenStakingUnsetOperator{
					Delegator: delAddr,
				}
			},
		},
		{
			name:     "pass: no operator set",
			existing: false,
			input: func(delAddr common.Address) *bindings.IPTokenStakingUnsetOperator {
				return &bindings.IPTokenStakingUnsetOperator{
					Delegator: delAddr,
				}
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _, _, _, _, _, _, esk := createKeeperWithMockStaking(t)

			ev := tc.input(common.Address(accAddrs[0]))
			cachedCtx, _ := ctx.CacheContext()

			if tc.existing {
				require.NoError(t, esk.DelegatorOperatorAddress.Set(cachedCtx, delAddr.String(), anotherExecAddr.String()))
			}

			err := esk.ProcessUnsetOperator(cachedCtx, ev)
			require.NoError(t, err)

			// check result
			has, err := esk.DelegatorOperatorAddress.Has(cachedCtx, delAddr.String())
			require.NoError(t, err)
			require.False(t, has)
		})
	}
}
