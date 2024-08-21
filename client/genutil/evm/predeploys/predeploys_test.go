//nolint:unused // fix with proper predeploy test
package predeploys_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/storyprotocol/iliad/client/genutil/evm/predeploys"
	"github.com/storyprotocol/iliad/contracts/bindings"
)

var (
	// Namespace big.Ints.
	mainNamespace = addrToBig(common.HexToAddress(predeploys.IliadNamespace))

	// Predeploy addresses.
	ipTokenStaking = common.HexToAddress(predeploys.IPTokenStaking)

	// Predeploy bytecodes.
	ipTokenStakingCode  = mustDecodeHex(bindings.IPTokenStakingDeployedBytecode)
	ipTokenSlashingCode = mustDecodeHex(bindings.IPTokenSlashingDeployedBytecode)
)

func TestAlloc(t *testing.T) {
	t.Parallel()

	// predeps, err := predeploys.Alloc(netconf.Staging)
	// require.NoError(t, err)

	// Check namespace filled with proxies.
	// for i := 0; i < predeploys.NamespaceSize; i++ {
	//	// add one, so that we don't set the namespace zero address
	//	addr := namespaceAddr(mainNamespace, i+1)
	//
	//	// proxy admin is not a proxy
	//	if addr == proxyAdmin {
	//		continue
	//	}
	//
	//	proxyAloc, ok := predeps[addr]
	//	require.Truef(t, ok, "proxy not found")
	//	require.Equalf(t,
	//		hexutil.Encode(proxyCode),
	//		hexutil.Encode(proxyAloc.Code),
	//		"proxy code mismatch")
	//	require.Equalf(t,
	//		common.HexToHash(predeploys.ProxyAdmin),
	//		proxyAloc.Storage[common.HexToHash(predeploys.ProxyAdminSlot)],
	//		"proxy admin slot empty")
	//}
	//
	//// check ProxyAdmin
	// proxyAdminAlloc, ok := predeps[proxyAdmin]
	// require.True(t, ok, "proxy admin not found")
	// require.Equal(t,
	//	proxyAdminCode,
	//	proxyAdminAlloc.Code,
	//	"proxy admin code mismatch")
	//
	//// check IPTokenStaking proxy
	// ipTokenStakingAlloc, ok := predeps[ipTokenStaking]
	// require.True(t, ok, "iliad stake not found")
	// require.Equal(t,
	//	proxyCode,
	//	ipTokenStakingAlloc.Code,
	//	"iliad stake code mismatch")
	//
	//// check IPTokenStaking implementation
	// ipTokenStakingImpl, ok := ipTokenStakingAlloc.Storage[common.HexToHash(predeploys.ProxyImplementationSlot)]
	// require.True(t, ok, "iliad stake implementation not found")
	//
	// ipTokenStakingImplAlloc, ok := predeps[common.BytesToAddress(ipTokenStakingImpl.Bytes())]
	// require.True(t, ok, "iliad stake implementation not found")
	//
	// require.Equal(t,
	//	ipTokenStakingCode,
	//	ipTokenStakingImplAlloc.Code,
	//	"iliad stake implementation mismatch")
}

// namespaceAddr returns the address at the given index in the namespace.
func namespaceAddr(namespace *big.Int, i int) common.Address {
	return common.BigToAddress(new(big.Int).Add(namespace, big.NewInt(int64(i))))
}

// addrToBig converts an address to a big.Int.
func addrToBig(addr common.Address) *big.Int {
	return new(big.Int).SetBytes(addr.Bytes())
}

// mustDecodeHex decodes the given hex string. It panics on error.
func mustDecodeHex(hex string) []byte {
	b, err := hexutil.Decode(hex)
	if err != nil {
		panic(err)
	}

	return b
}
