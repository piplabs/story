// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// UpgradeEntrypointMetaData contains all meta data concerning the UpgradeEntrypoint contract.
var UpgradeEntrypointMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"accessManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"planUpgrade\",\"inputs\":[{\"name\":\"appVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"upgradeHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"info\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SoftwareUpgrade\",\"inputs\":[{\"name\":\"appVersion\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"upgradeHeight\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"info\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60a080604052346100cc57306080527ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c166100bd57506001600160401b036002600160401b031982821601610078575b604051610f8d90816100d282396080518181816107aa01526108c60152f35b6001600160401b031990911681179091556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a1388080610059565b63f92ee8a960e01b8152600490fd5b600080fdfe60406080815260048036101561001457600080fd5b600091823560e01c9182634065a81514610b605782634f1ef2861461081f57826352d1902d14610762578263715018a61461067d57826379ba5097146105d25782638da5cb5b14610560578263ad3cb1cc14610453578263c4d66de8146101f557508163e30c39781461017f575063f2fde38b1461009157600080fd5b3461017c5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261017c576100c8610c68565b6100d0610d3a565b73ffffffffffffffffffffffffffffffffffffffff809116907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227008380a380f35b80fd5b9050346101f157817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101f15760209073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0054169051908152f35b5080fd5b9091503461044f5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261044f5761022f610c68565b907ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009182549160ff83861c16159267ffffffffffffffff811680159081610447575b600114908161043d575b159081610434575b5061040c578360017fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000083161786556103d7575b5073ffffffffffffffffffffffffffffffffffffffff82161561035457506102f5906102e0610e5e565b6102e8610e5e565b6102f0610e5e565b610daa565b6102fd578280f35b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d291817fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff602093541690555160018152a138808280f35b60849060208651917f08c379a0000000000000000000000000000000000000000000000000000000008352820152603760248201527f55706772616465456e747279706f696e743a206163636573734d616e6167657260448201527f2063616e6e6f74206265207a65726f20616464726573730000000000000000006064820152fd5b7fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000001668010000000000000001178455386102b6565b5084517ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b90501538610283565b303b15915061027b565b859150610271565b8280fd5b9091503461044f57827ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261044f578151908282019082821067ffffffffffffffff83111761053457508252600581526020907f352e302e300000000000000000000000000000000000000000000000000000006020820152825193849260208452825192836020860152825b84811061051e57505050828201840152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168101030190f35b81810183015188820188015287955082016104e2565b8460416024927f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b8382346101f157817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101f15760209073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054169051908152f35b91503461044f57827ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261044f573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0054160361064d578261064a33610daa565b80f35b6024925051907f118cdaa70000000000000000000000000000000000000000000000000000000082523390820152fd5b833461017c57807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261017c576106b4610d3a565b8073ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b90833461017c57807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261017c575073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001630036107f957602090517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b517fe07c8dba000000000000000000000000000000000000000000000000000000008152fd5b9150807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261044f57610852610c68565b90602493843567ffffffffffffffff81116101f157366023820112156101f1578085013561087f81610d00565b9461088c85519687610c90565b81865260209182870193368a8383010111610b5c578186928b86930187378801015273ffffffffffffffffffffffffffffffffffffffff807f000000000000000000000000000000000000000000000000000000000000000016803014908115610b2e575b50610b06576108fe610d3a565b81169585517f52d1902d00000000000000000000000000000000000000000000000000000000815283818a818b5afa869181610ad3575b5061096a5750505050505051917f4c9c8ce3000000000000000000000000000000000000000000000000000000008352820152fd5b9088888894938c7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc91828103610aa65750853b15610a79575080547fffffffffffffffffffffffff000000000000000000000000000000000000000016821790558451889392917fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b8580a2825115610a42575050610a349582915190845af4913d15610a38573d610a26610a1d82610d00565b92519283610c90565b81528581943d92013e610eb7565b5080f35b5060609250610eb7565b955095505050505034610a5457505080f35b7fb398979f000000000000000000000000000000000000000000000000000000008152fd5b83838851917f4c9c8ce3000000000000000000000000000000000000000000000000000000008352820152fd5b84908851917faa1d49a4000000000000000000000000000000000000000000000000000000008352820152fd5b9091508481813d8311610aff575b610aeb8183610c90565b81010312610afb57519038610935565b8680fd5b503d610ae1565b8786517fe07c8dba000000000000000000000000000000000000000000000000000000008152fd5b9050817f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc54161415386108f1565b8580fd5b9091503461044f5760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261044f57803567ffffffffffffffff90818116809103610c645760243592828416809403610b5c5760443590838211610afb5736602383011215610afb57810135928311610b5c573660248483010111610b5c57601f8360809460247fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09460607f914a3dba53fada9763698f71d4e9e847abcc184646b9bc29c243cf00c183a6459a610c38610d3a565b80519a8b998a5260208a015288015282606088015201868601378785828601015201168101030190a180f35b8480fd5b6004359073ffffffffffffffffffffffffffffffffffffffff82168203610c8b57565b600080fd5b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff821117610cd157604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff8111610cd157601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054163303610d7a57565b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b7fffffffffffffffffffffffff0000000000000000000000000000000000000000907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008281541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549073ffffffffffffffffffffffffffffffffffffffff80931680948316179055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3565b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c1615610e8d57565b60046040517fd7e6bcf8000000000000000000000000000000000000000000000000000000008152fd5b90610ef65750805115610ecc57805190602001fd5b60046040517f1425ea42000000000000000000000000000000000000000000000000000000008152fd5b81511580610f4e575b610f07575090565b60249073ffffffffffffffffffffffffffffffffffffffff604051917f9996b315000000000000000000000000000000000000000000000000000000008352166004820152fd5b50803b15610eff56fea2646970667358221220769bd9a03478d2b6fd4343ebff70f214a346e5739c02df176c0007d61cac101664736f6c63430008170033",
}

// UpgradeEntrypointABI is the input ABI used to generate the binding from.
// Deprecated: Use UpgradeEntrypointMetaData.ABI instead.
var UpgradeEntrypointABI = UpgradeEntrypointMetaData.ABI

// UpgradeEntrypointBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UpgradeEntrypointMetaData.Bin instead.
var UpgradeEntrypointBin = UpgradeEntrypointMetaData.Bin

// DeployUpgradeEntrypoint deploys a new Ethereum contract, binding an instance of UpgradeEntrypoint to it.
func DeployUpgradeEntrypoint(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UpgradeEntrypoint, error) {
	parsed, err := UpgradeEntrypointMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpgradeEntrypointBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpgradeEntrypoint{UpgradeEntrypointCaller: UpgradeEntrypointCaller{contract: contract}, UpgradeEntrypointTransactor: UpgradeEntrypointTransactor{contract: contract}, UpgradeEntrypointFilterer: UpgradeEntrypointFilterer{contract: contract}}, nil
}

// UpgradeEntrypoint is an auto generated Go binding around an Ethereum contract.
type UpgradeEntrypoint struct {
	UpgradeEntrypointCaller     // Read-only binding to the contract
	UpgradeEntrypointTransactor // Write-only binding to the contract
	UpgradeEntrypointFilterer   // Log filterer for contract events
}

// UpgradeEntrypointCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpgradeEntrypointCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeEntrypointTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpgradeEntrypointTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeEntrypointFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpgradeEntrypointFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeEntrypointSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpgradeEntrypointSession struct {
	Contract     *UpgradeEntrypoint // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// UpgradeEntrypointCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpgradeEntrypointCallerSession struct {
	Contract *UpgradeEntrypointCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// UpgradeEntrypointTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpgradeEntrypointTransactorSession struct {
	Contract     *UpgradeEntrypointTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// UpgradeEntrypointRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpgradeEntrypointRaw struct {
	Contract *UpgradeEntrypoint // Generic contract binding to access the raw methods on
}

// UpgradeEntrypointCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpgradeEntrypointCallerRaw struct {
	Contract *UpgradeEntrypointCaller // Generic read-only contract binding to access the raw methods on
}

// UpgradeEntrypointTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpgradeEntrypointTransactorRaw struct {
	Contract *UpgradeEntrypointTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpgradeEntrypoint creates a new instance of UpgradeEntrypoint, bound to a specific deployed contract.
func NewUpgradeEntrypoint(address common.Address, backend bind.ContractBackend) (*UpgradeEntrypoint, error) {
	contract, err := bindUpgradeEntrypoint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpgradeEntrypoint{UpgradeEntrypointCaller: UpgradeEntrypointCaller{contract: contract}, UpgradeEntrypointTransactor: UpgradeEntrypointTransactor{contract: contract}, UpgradeEntrypointFilterer: UpgradeEntrypointFilterer{contract: contract}}, nil
}

// NewUpgradeEntrypointCaller creates a new read-only instance of UpgradeEntrypoint, bound to a specific deployed contract.
func NewUpgradeEntrypointCaller(address common.Address, caller bind.ContractCaller) (*UpgradeEntrypointCaller, error) {
	contract, err := bindUpgradeEntrypoint(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeEntrypointCaller{contract: contract}, nil
}

// NewUpgradeEntrypointTransactor creates a new write-only instance of UpgradeEntrypoint, bound to a specific deployed contract.
func NewUpgradeEntrypointTransactor(address common.Address, transactor bind.ContractTransactor) (*UpgradeEntrypointTransactor, error) {
	contract, err := bindUpgradeEntrypoint(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeEntrypointTransactor{contract: contract}, nil
}

// NewUpgradeEntrypointFilterer creates a new log filterer instance of UpgradeEntrypoint, bound to a specific deployed contract.
func NewUpgradeEntrypointFilterer(address common.Address, filterer bind.ContractFilterer) (*UpgradeEntrypointFilterer, error) {
	contract, err := bindUpgradeEntrypoint(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpgradeEntrypointFilterer{contract: contract}, nil
}

// bindUpgradeEntrypoint binds a generic wrapper to an already deployed contract.
func bindUpgradeEntrypoint(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UpgradeEntrypointMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpgradeEntrypoint *UpgradeEntrypointRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpgradeEntrypoint.Contract.UpgradeEntrypointCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpgradeEntrypoint *UpgradeEntrypointRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.UpgradeEntrypointTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpgradeEntrypoint *UpgradeEntrypointRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.UpgradeEntrypointTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpgradeEntrypoint *UpgradeEntrypointCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpgradeEntrypoint.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpgradeEntrypoint *UpgradeEntrypointTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpgradeEntrypoint *UpgradeEntrypointTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_UpgradeEntrypoint *UpgradeEntrypointCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UpgradeEntrypoint.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_UpgradeEntrypoint *UpgradeEntrypointSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _UpgradeEntrypoint.Contract.UPGRADEINTERFACEVERSION(&_UpgradeEntrypoint.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_UpgradeEntrypoint *UpgradeEntrypointCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _UpgradeEntrypoint.Contract.UPGRADEINTERFACEVERSION(&_UpgradeEntrypoint.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UpgradeEntrypoint *UpgradeEntrypointCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UpgradeEntrypoint.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UpgradeEntrypoint *UpgradeEntrypointSession) Owner() (common.Address, error) {
	return _UpgradeEntrypoint.Contract.Owner(&_UpgradeEntrypoint.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UpgradeEntrypoint *UpgradeEntrypointCallerSession) Owner() (common.Address, error) {
	return _UpgradeEntrypoint.Contract.Owner(&_UpgradeEntrypoint.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_UpgradeEntrypoint *UpgradeEntrypointCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UpgradeEntrypoint.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_UpgradeEntrypoint *UpgradeEntrypointSession) PendingOwner() (common.Address, error) {
	return _UpgradeEntrypoint.Contract.PendingOwner(&_UpgradeEntrypoint.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_UpgradeEntrypoint *UpgradeEntrypointCallerSession) PendingOwner() (common.Address, error) {
	return _UpgradeEntrypoint.Contract.PendingOwner(&_UpgradeEntrypoint.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_UpgradeEntrypoint *UpgradeEntrypointCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _UpgradeEntrypoint.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_UpgradeEntrypoint *UpgradeEntrypointSession) ProxiableUUID() ([32]byte, error) {
	return _UpgradeEntrypoint.Contract.ProxiableUUID(&_UpgradeEntrypoint.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_UpgradeEntrypoint *UpgradeEntrypointCallerSession) ProxiableUUID() ([32]byte, error) {
	return _UpgradeEntrypoint.Contract.ProxiableUUID(&_UpgradeEntrypoint.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeEntrypoint.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_UpgradeEntrypoint *UpgradeEntrypointSession) AcceptOwnership() (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.AcceptOwnership(&_UpgradeEntrypoint.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.AcceptOwnership(&_UpgradeEntrypoint.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address accessManager) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactor) Initialize(opts *bind.TransactOpts, accessManager common.Address) (*types.Transaction, error) {
	return _UpgradeEntrypoint.contract.Transact(opts, "initialize", accessManager)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address accessManager) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointSession) Initialize(accessManager common.Address) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.Initialize(&_UpgradeEntrypoint.TransactOpts, accessManager)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address accessManager) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactorSession) Initialize(accessManager common.Address) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.Initialize(&_UpgradeEntrypoint.TransactOpts, accessManager)
}

// PlanUpgrade is a paid mutator transaction binding the contract method 0x4065a815.
//
// Solidity: function planUpgrade(uint64 appVersion, uint64 upgradeHeight, string info) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactor) PlanUpgrade(opts *bind.TransactOpts, appVersion uint64, upgradeHeight uint64, info string) (*types.Transaction, error) {
	return _UpgradeEntrypoint.contract.Transact(opts, "planUpgrade", appVersion, upgradeHeight, info)
}

// PlanUpgrade is a paid mutator transaction binding the contract method 0x4065a815.
//
// Solidity: function planUpgrade(uint64 appVersion, uint64 upgradeHeight, string info) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointSession) PlanUpgrade(appVersion uint64, upgradeHeight uint64, info string) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.PlanUpgrade(&_UpgradeEntrypoint.TransactOpts, appVersion, upgradeHeight, info)
}

// PlanUpgrade is a paid mutator transaction binding the contract method 0x4065a815.
//
// Solidity: function planUpgrade(uint64 appVersion, uint64 upgradeHeight, string info) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactorSession) PlanUpgrade(appVersion uint64, upgradeHeight uint64, info string) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.PlanUpgrade(&_UpgradeEntrypoint.TransactOpts, appVersion, upgradeHeight, info)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeEntrypoint.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UpgradeEntrypoint *UpgradeEntrypointSession) RenounceOwnership() (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.RenounceOwnership(&_UpgradeEntrypoint.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.RenounceOwnership(&_UpgradeEntrypoint.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _UpgradeEntrypoint.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.TransferOwnership(&_UpgradeEntrypoint.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.TransferOwnership(&_UpgradeEntrypoint.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _UpgradeEntrypoint.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_UpgradeEntrypoint *UpgradeEntrypointSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.UpgradeToAndCall(&_UpgradeEntrypoint.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.UpgradeToAndCall(&_UpgradeEntrypoint.TransactOpts, newImplementation, data)
}

// UpgradeEntrypointInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the UpgradeEntrypoint contract.
type UpgradeEntrypointInitializedIterator struct {
	Event *UpgradeEntrypointInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpgradeEntrypointInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeEntrypointInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UpgradeEntrypointInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UpgradeEntrypointInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeEntrypointInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeEntrypointInitialized represents a Initialized event raised by the UpgradeEntrypoint contract.
type UpgradeEntrypointInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) FilterInitialized(opts *bind.FilterOpts) (*UpgradeEntrypointInitializedIterator, error) {

	logs, sub, err := _UpgradeEntrypoint.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &UpgradeEntrypointInitializedIterator{contract: _UpgradeEntrypoint.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *UpgradeEntrypointInitialized) (event.Subscription, error) {

	logs, sub, err := _UpgradeEntrypoint.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeEntrypointInitialized)
				if err := _UpgradeEntrypoint.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) ParseInitialized(log types.Log) (*UpgradeEntrypointInitialized, error) {
	event := new(UpgradeEntrypointInitialized)
	if err := _UpgradeEntrypoint.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeEntrypointOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the UpgradeEntrypoint contract.
type UpgradeEntrypointOwnershipTransferStartedIterator struct {
	Event *UpgradeEntrypointOwnershipTransferStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpgradeEntrypointOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeEntrypointOwnershipTransferStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UpgradeEntrypointOwnershipTransferStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UpgradeEntrypointOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeEntrypointOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeEntrypointOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the UpgradeEntrypoint contract.
type UpgradeEntrypointOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*UpgradeEntrypointOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UpgradeEntrypoint.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &UpgradeEntrypointOwnershipTransferStartedIterator{contract: _UpgradeEntrypoint.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *UpgradeEntrypointOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UpgradeEntrypoint.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeEntrypointOwnershipTransferStarted)
				if err := _UpgradeEntrypoint.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) ParseOwnershipTransferStarted(log types.Log) (*UpgradeEntrypointOwnershipTransferStarted, error) {
	event := new(UpgradeEntrypointOwnershipTransferStarted)
	if err := _UpgradeEntrypoint.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeEntrypointOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the UpgradeEntrypoint contract.
type UpgradeEntrypointOwnershipTransferredIterator struct {
	Event *UpgradeEntrypointOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpgradeEntrypointOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeEntrypointOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UpgradeEntrypointOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UpgradeEntrypointOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeEntrypointOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeEntrypointOwnershipTransferred represents a OwnershipTransferred event raised by the UpgradeEntrypoint contract.
type UpgradeEntrypointOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*UpgradeEntrypointOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UpgradeEntrypoint.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &UpgradeEntrypointOwnershipTransferredIterator{contract: _UpgradeEntrypoint.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *UpgradeEntrypointOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UpgradeEntrypoint.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeEntrypointOwnershipTransferred)
				if err := _UpgradeEntrypoint.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) ParseOwnershipTransferred(log types.Log) (*UpgradeEntrypointOwnershipTransferred, error) {
	event := new(UpgradeEntrypointOwnershipTransferred)
	if err := _UpgradeEntrypoint.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeEntrypointSoftwareUpgradeIterator is returned from FilterSoftwareUpgrade and is used to iterate over the raw logs and unpacked data for SoftwareUpgrade events raised by the UpgradeEntrypoint contract.
type UpgradeEntrypointSoftwareUpgradeIterator struct {
	Event *UpgradeEntrypointSoftwareUpgrade // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpgradeEntrypointSoftwareUpgradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeEntrypointSoftwareUpgrade)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UpgradeEntrypointSoftwareUpgrade)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UpgradeEntrypointSoftwareUpgradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeEntrypointSoftwareUpgradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeEntrypointSoftwareUpgrade represents a SoftwareUpgrade event raised by the UpgradeEntrypoint contract.
type UpgradeEntrypointSoftwareUpgrade struct {
	AppVersion    uint64
	UpgradeHeight uint64
	Info          string
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSoftwareUpgrade is a free log retrieval operation binding the contract event 0x914a3dba53fada9763698f71d4e9e847abcc184646b9bc29c243cf00c183a645.
//
// Solidity: event SoftwareUpgrade(uint64 appVersion, uint64 upgradeHeight, string info)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) FilterSoftwareUpgrade(opts *bind.FilterOpts) (*UpgradeEntrypointSoftwareUpgradeIterator, error) {

	logs, sub, err := _UpgradeEntrypoint.contract.FilterLogs(opts, "SoftwareUpgrade")
	if err != nil {
		return nil, err
	}
	return &UpgradeEntrypointSoftwareUpgradeIterator{contract: _UpgradeEntrypoint.contract, event: "SoftwareUpgrade", logs: logs, sub: sub}, nil
}

// WatchSoftwareUpgrade is a free log subscription operation binding the contract event 0x914a3dba53fada9763698f71d4e9e847abcc184646b9bc29c243cf00c183a645.
//
// Solidity: event SoftwareUpgrade(uint64 appVersion, uint64 upgradeHeight, string info)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) WatchSoftwareUpgrade(opts *bind.WatchOpts, sink chan<- *UpgradeEntrypointSoftwareUpgrade) (event.Subscription, error) {

	logs, sub, err := _UpgradeEntrypoint.contract.WatchLogs(opts, "SoftwareUpgrade")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeEntrypointSoftwareUpgrade)
				if err := _UpgradeEntrypoint.contract.UnpackLog(event, "SoftwareUpgrade", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSoftwareUpgrade is a log parse operation binding the contract event 0x914a3dba53fada9763698f71d4e9e847abcc184646b9bc29c243cf00c183a645.
//
// Solidity: event SoftwareUpgrade(uint64 appVersion, uint64 upgradeHeight, string info)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) ParseSoftwareUpgrade(log types.Log) (*UpgradeEntrypointSoftwareUpgrade, error) {
	event := new(UpgradeEntrypointSoftwareUpgrade)
	if err := _UpgradeEntrypoint.contract.UnpackLog(event, "SoftwareUpgrade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeEntrypointUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the UpgradeEntrypoint contract.
type UpgradeEntrypointUpgradedIterator struct {
	Event *UpgradeEntrypointUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpgradeEntrypointUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeEntrypointUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(UpgradeEntrypointUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *UpgradeEntrypointUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeEntrypointUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeEntrypointUpgraded represents a Upgraded event raised by the UpgradeEntrypoint contract.
type UpgradeEntrypointUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*UpgradeEntrypointUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _UpgradeEntrypoint.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &UpgradeEntrypointUpgradedIterator{contract: _UpgradeEntrypoint.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *UpgradeEntrypointUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _UpgradeEntrypoint.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeEntrypointUpgraded)
				if err := _UpgradeEntrypoint.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) ParseUpgraded(log types.Log) (*UpgradeEntrypointUpgraded, error) {
	event := new(UpgradeEntrypointUpgraded)
	if err := _UpgradeEntrypoint.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
