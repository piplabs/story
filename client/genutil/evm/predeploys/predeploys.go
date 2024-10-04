//nolint:unused // fix with proper predeploy script
package predeploys

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/piplabs/story/client/genutil/evm/state"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/netconf"
	"github.com/piplabs/story/lib/solc"
)

const (
	// StoryNamespace is namespace of for story specific predeploys.
	StoryNamespace = "0x121E240000000000000000000000000000000000"

	// IPTokenNamespace is namespace of for IP Token specific predeploys.
	IPTokenNamespace = "0xcccccc0000000000000000000000000000000000"

	// NamespaceSize is the number of proxies to deploy per namespace.
	NamespaceSize = 2048

	// IP Token Predeploys.
	IPTokenStaking    = "0xcccccc0000000000000000000000000000000001"
	IPTokenSlashing   = "0xa39241Eb9Ff830178339D1E6aD38EfB160Ee9ab1"
	UpgradeEntrypoint = "0xcccccc0000000000000000000000000000000003"

	Secp256k1 = "0x00000000000000000000000000000000000256f1"

	// TransparentUpgradeableProxy storage slots.
	// ProxyImplementationSlot = "0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"
	// ProxyAdminSlot          = "0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103".
)

var (
	// Namespace big.Ints.
	storyNamespace   = common.HexToAddress(StoryNamespace).Big()
	ipTokenNamespace = common.HexToAddress(IPTokenNamespace).Big()

	// Predeploy addresses.
	ipTokenStaking    = common.HexToAddress(IPTokenStaking)
	ipTokenSlashing   = common.HexToAddress(IPTokenSlashing)
	upgradeEntrypoint = common.HexToAddress(UpgradeEntrypoint)

	// Predeploy bytecodes.
	ipTokenStakingCode = hexutil.MustDecode(bindings.IPTokenStakingDeployedBytecode)
)

// Alloc returns the genesis allocs for the predeployed contracts, initializing code and storage.
func Alloc(_ netconf.ID) (types.GenesisAlloc, error) {
	emptyGenesis := &core.Genesis{Alloc: types.GenesisAlloc{}}

	db := state.NewMemDB(emptyGenesis)

	// setProxies(db)

	// admin, err := eoa.Admin(network)
	// if err != nil {
	//	return nil, errors.Wrap(err, "network admin")
	//}

	if err := setStaking(db); err != nil {
		return nil, errors.Wrap(err, "set staking")
	}

	if err := setSlashing(db); err != nil {
		return nil, errors.Wrap(err, "set slashing")
	}

	return db.Genesis().Alloc, nil
}

// setStaking sets the Staking predeploy.
func setStaking(db *state.MemDB) error {
	storage := state.StorageValues{}

	return setPredeploy(db, ipTokenStaking, ipTokenStakingCode, bindings.IPTokenStakingStorageLayout, storage)
}

// setSlashing sets the Slashing predeploy.
// TODO: Slashing design.
func setSlashing(db *state.MemDB) error {
	storage := state.StorageValues{}

	return setPredeploy(db, ipTokenSlashing, ipTokenStakingCode, bindings.IPTokenStakingStorageLayout, storage)
}

// setPredeploy sets the implementation code and proxy storage for the given predeploy.
func setPredeploy(db *state.MemDB, proxy common.Address, code []byte, layout solc.StorageLayout, storage state.StorageValues) error {
	impl := impl(proxy)
	// setProxyImplementation(db, proxy, impl)
	db.SetCode(impl, code)

	return setStorage(db, proxy, layout, storage)
}

// setStorage sets the code and storage for the given predeploy.
func setStorage(db *state.MemDB, addr common.Address, layout solc.StorageLayout, storage state.StorageValues) error {
	slots, err := state.EncodeStorage(layout, storage)
	if err != nil {
		return errors.Wrap(err, "encode storage", "addr", addr)
	}

	for _, slot := range slots {
		db.SetState(addr, slot.Key, slot.Value)
	}

	return nil
}

// setProxyImplementation sets the implementation address for the given proxy address.
// func setProxyImplementation(db *state.MemDB, proxy, impl common.Address) {
//	db.SetState(proxy, common.HexToHash(ProxyImplementationSlot), common.HexToHash(impl.Hex()))
//}

// impl returns the implementation address for the given proxy address.
func impl(addr common.Address) common.Address {
	// To get a unique implementation per each proxy address, we subtract the address from the max address.
	// Max address is odd, so the result will be unique.
	maxAddr := common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffff").Big()
	return common.BigToAddress(new(big.Int).Sub(maxAddr, addr.Big()))
}

//nolint:unused // address returns the address at the given index in the namespace.
func address(namespace *big.Int, i int) common.Address {
	return common.BigToAddress(new(big.Int).Add(namespace, big.NewInt(int64(i))))
}
