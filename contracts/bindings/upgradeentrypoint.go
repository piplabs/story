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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"accessManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"planUpgrade\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"info\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SoftwareUpgrade\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"int64\",\"indexed\":false,\"internalType\":\"int64\"},{\"name\":\"info\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60a06040523060805234801561001457600080fd5b5061001d610022565b6100d4565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100725760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100d15780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b608051610d696100fd6000396000818161054a0152818161057301526106b90152610d696000f3fe6080604052600436106100915760003560e01c8063ad3cb1cc11610059578063ad3cb1cc1461012a578063c4d66de814610168578063e30c397814610188578063ef176e0e1461019d578063f2fde38b146101bd57600080fd5b80634f1ef2861461009657806352d1902d146100ab578063715018a6146100d357806379ba5097146100e85780638da5cb5b146100fd575b600080fd5b6100a96100a4366004610a70565b6101dd565b005b3480156100b757600080fd5b506100c06101fc565b6040519081526020015b60405180910390f35b3480156100df57600080fd5b506100a9610219565b3480156100f457600080fd5b506100a961022d565b34801561010957600080fd5b5061011261027a565b6040516001600160a01b0390911681526020016100ca565b34801561013657600080fd5b5061015b604051806040016040528060058152602001640352e302e360dc1b81525081565b6040516100ca9190610b56565b34801561017457600080fd5b506100a9610183366004610b89565b6102af565b34801561019457600080fd5b50610112610443565b3480156101a957600080fd5b506100a96101b8366004610bed565b61046c565b3480156101c957600080fd5b506100a96101d8366004610b89565b6104ba565b6101e561053f565b6101ee826105e4565b6101f882826105ec565b5050565b60006102066106ae565b50600080516020610d1483398151915290565b6102216106f7565b61022b6000610729565b565b3380610237610443565b6001600160a01b03161461026e5760405163118cdaa760e01b81526001600160a01b03821660048201526024015b60405180910390fd5b61027781610729565b50565b6000807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b546001600160a01b031692915050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008054600160401b810460ff16159067ffffffffffffffff166000811580156102f55750825b905060008267ffffffffffffffff1660011480156103125750303b155b905081158015610320575080155b1561033e5760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561036857845460ff60401b1916600160401b1785555b6001600160a01b0386166103e45760405162461bcd60e51b815260206004820152603760248201527f55706772616465456e747279706f696e743a206163636573734d616e6167657260448201527f2063616e6e6f74206265207a65726f20616464726573730000000000000000006064820152608401610265565b6103ec610761565b6103f586610769565b831561043b57845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b6000807f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0061029f565b6104746106f7565b7f112749e79b2098b58eab36c21f123b2883c3ecbbb4f41623a744fa6d9b3e37c685858585856040516104ab959493929190610ca2565b60405180910390a15050505050565b6104c26106f7565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0080546001600160a01b0319166001600160a01b038316908117825561050661027a565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614806105c657507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166105ba600080516020610d14833981519152546001600160a01b031690565b6001600160a01b031614155b1561022b5760405163703e46dd60e11b815260040160405180910390fd5b6102776106f7565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610646575060408051601f3d908101601f1916820190925261064391810190610cde565b60015b61066e57604051634c9c8ce360e01b81526001600160a01b0383166004820152602401610265565b600080516020610d14833981519152811461069f57604051632a87526960e21b815260048101829052602401610265565b6106a9838361077a565b505050565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461022b5760405163703e46dd60e11b815260040160405180910390fd5b3361070061027a565b6001600160a01b03161461022b5760405163118cdaa760e01b8152336004820152602401610265565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0080546001600160a01b03191681556101f8826107d0565b61022b610841565b610771610841565b6102778161088a565b610783826108bc565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a28051156107c8576106a98282610921565b6101f8610997565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054600160401b900460ff1661022b57604051631afcd79f60e31b815260040160405180910390fd5b610892610841565b6001600160a01b03811661026e57604051631e4fbdf760e01b815260006004820152602401610265565b806001600160a01b03163b6000036108f257604051634c9c8ce360e01b81526001600160a01b0382166004820152602401610265565b600080516020610d1483398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b03168460405161093e9190610cf7565b600060405180830381855af49150503d8060008114610979576040519150601f19603f3d011682016040523d82523d6000602084013e61097e565b606091505b509150915061098e8583836109b6565b95945050505050565b341561022b5760405163b398979f60e01b815260040160405180910390fd5b6060826109cb576109c682610a15565b610a0e565b81511580156109e257506001600160a01b0384163b155b15610a0b57604051639996b31560e01b81526001600160a01b0385166004820152602401610265565b50805b9392505050565b805115610a255780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b80356001600160a01b0381168114610a5557600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60008060408385031215610a8357600080fd5b610a8c83610a3e565b9150602083013567ffffffffffffffff80821115610aa957600080fd5b818501915085601f830112610abd57600080fd5b813581811115610acf57610acf610a5a565b604051601f8201601f19908116603f01168101908382118183101715610af757610af7610a5a565b81604052828152886020848701011115610b1057600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60005b83811015610b4d578181015183820152602001610b35565b50506000910152565b6020815260008251806020840152610b75816040850160208701610b32565b601f01601f19169190910160400192915050565b600060208284031215610b9b57600080fd5b610a0e82610a3e565b60008083601f840112610bb657600080fd5b50813567ffffffffffffffff811115610bce57600080fd5b602083019150836020828501011115610be657600080fd5b9250929050565b600080600080600060608688031215610c0557600080fd5b853567ffffffffffffffff80821115610c1d57600080fd5b610c2989838a01610ba4565b909750955060208801359150600782900b8214610c4557600080fd5b90935060408701359080821115610c5b57600080fd5b50610c6888828901610ba4565b969995985093965092949392505050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b606081526000610cb6606083018789610c79565b8560070b60208401528281036040840152610cd2818587610c79565b98975050505050505050565b600060208284031215610cf057600080fd5b5051919050565b60008251610d09818460208701610b32565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca26469706673582212202fed7feda2d375d5be4b2847badd7e61be35b824b0a1a393337ce76a48fe290564736f6c63430008180033",
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

// PlanUpgrade is a paid mutator transaction binding the contract method 0xef176e0e.
//
// Solidity: function planUpgrade(string name, int64 height, string info) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactor) PlanUpgrade(opts *bind.TransactOpts, name string, height int64, info string) (*types.Transaction, error) {
	return _UpgradeEntrypoint.contract.Transact(opts, "planUpgrade", name, height, info)
}

// PlanUpgrade is a paid mutator transaction binding the contract method 0xef176e0e.
//
// Solidity: function planUpgrade(string name, int64 height, string info) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointSession) PlanUpgrade(name string, height int64, info string) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.PlanUpgrade(&_UpgradeEntrypoint.TransactOpts, name, height, info)
}

// PlanUpgrade is a paid mutator transaction binding the contract method 0xef176e0e.
//
// Solidity: function planUpgrade(string name, int64 height, string info) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactorSession) PlanUpgrade(name string, height int64, info string) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.PlanUpgrade(&_UpgradeEntrypoint.TransactOpts, name, height, info)
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
	Name   string
	Height int64
	Info   string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSoftwareUpgrade is a free log retrieval operation binding the contract event 0x112749e79b2098b58eab36c21f123b2883c3ecbbb4f41623a744fa6d9b3e37c6.
//
// Solidity: event SoftwareUpgrade(string name, int64 height, string info)
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) FilterSoftwareUpgrade(opts *bind.FilterOpts) (*UpgradeEntrypointSoftwareUpgradeIterator, error) {

	logs, sub, err := _UpgradeEntrypoint.contract.FilterLogs(opts, "SoftwareUpgrade")
	if err != nil {
		return nil, err
	}
	return &UpgradeEntrypointSoftwareUpgradeIterator{contract: _UpgradeEntrypoint.contract, event: "SoftwareUpgrade", logs: logs, sub: sub}, nil
}

// WatchSoftwareUpgrade is a free log subscription operation binding the contract event 0x112749e79b2098b58eab36c21f123b2883c3ecbbb4f41623a744fa6d9b3e37c6.
//
// Solidity: event SoftwareUpgrade(string name, int64 height, string info)
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

// ParseSoftwareUpgrade is a log parse operation binding the contract event 0x112749e79b2098b58eab36c21f123b2883c3ecbbb4f41623a744fa6d9b3e37c6.
//
// Solidity: event SoftwareUpgrade(string name, int64 height, string info)
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
