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

// Create3MetaData contains all meta data concerning the Create3 contract.
var Create3MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"deploy\",\"inputs\":[{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"creationCode\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"deployed\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"deployDeterministic\",\"inputs\":[{\"name\":\"creationCode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"deployed\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getDeployed\",\"inputs\":[{\"name\":\"deployer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"deployed\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"predictDeterministicAddress\",\"inputs\":[{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"deployed\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x6080806040523461001657610667908161001b8239f35b5f80fdfe6080604081815260049182361015610015575f80fd5b5f3560e01c90816350f1c464146101f0575080635414dff01461019f5780639881d195146101335763cdcb760a1461004b575f80fd5b807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012f576024359167ffffffffffffffff831161012f5761012773ffffffffffffffffffffffffffffffffffffffff916100af60209536908301610358565b84513360601b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000168782019081529235601484015290919061011c81603484015b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081018352826102dd565b5190209034916104ea565b915191168152f35b5f80fd5b50807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012f57813567ffffffffffffffff811161012f5761012761019560209473ffffffffffffffffffffffffffffffffffffffff93369101610358565b34906024356104ea565b503461012f5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012f5773ffffffffffffffffffffffffffffffffffffffff610127602093356103d7565b82843461012f57817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261012f57359173ffffffffffffffffffffffffffffffffffffffff90818416840361012f5760609390931b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166020848101918252602435603486015293610127919061028c81605481016100f0565b5190206103d7565b6040810190811067ffffffffffffffff8211176102b057604052565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176102b057604052565b67ffffffffffffffff81116102b057601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b81601f8201121561012f5780359061036f8261031e565b9261037d60405194856102dd565b8284526020838301011161012f57815f926020809301838601378301015290565b604051906103ab82610294565b601082527f67363d3d37363d34f03d5260086018f3000000000000000000000000000000006020830152565b6103df61039e565b602081519101206040519060208201927fff0000000000000000000000000000000000000000000000000000000000000084523060601b602184015260358301526055820152605581526080810181811067ffffffffffffffff8211176102b0577f010000000000000000000000000000000000000000000000000000000000000060b673ffffffffffffffffffffffffffffffffffffffff948360405284519020937fffffffffffffffffffffffffffffffffffffffff00000000000000000000000060a08201957fd694000000000000000000000000000000000000000000000000000000000000875260601b1660a28201520152601781526104e381610294565b5190201690565b929192806104f661039e565b6020815191015ff59173ffffffffffffffffffffffffffffffffffffffff8316156105d3575f9261052784936103d7565b95602083519301915af13d156105ce573d6105418161031e565b9061054f60405192836102dd565b81525f60203d92013e5b806105c4575b1561056657565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f494e495449414c495a4154494f4e5f4641494c454400000000000000000000006044820152fd5b50813b151561055f565b610559565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f4445504c4f594d454e545f4641494c45440000000000000000000000000000006044820152fdfea26469706673582212203f5ac15a5d7253024be8653eb0cf2d7aa21c6b76cfa6a8d2a5f4f13f82e1cec864736f6c63430008170033",
}

// Create3ABI is the input ABI used to generate the binding from.
// Deprecated: Use Create3MetaData.ABI instead.
var Create3ABI = Create3MetaData.ABI

// Create3Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Create3MetaData.Bin instead.
var Create3Bin = Create3MetaData.Bin

// DeployCreate3 deploys a new Ethereum contract, binding an instance of Create3 to it.
func DeployCreate3(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Create3, error) {
	parsed, err := Create3MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Create3Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Create3{Create3Caller: Create3Caller{contract: contract}, Create3Transactor: Create3Transactor{contract: contract}, Create3Filterer: Create3Filterer{contract: contract}}, nil
}

// Create3 is an auto generated Go binding around an Ethereum contract.
type Create3 struct {
	Create3Caller     // Read-only binding to the contract
	Create3Transactor // Write-only binding to the contract
	Create3Filterer   // Log filterer for contract events
}

// Create3Caller is an auto generated read-only Go binding around an Ethereum contract.
type Create3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Create3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Create3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Create3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Create3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Create3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Create3Session struct {
	Contract     *Create3          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Create3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Create3CallerSession struct {
	Contract *Create3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// Create3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Create3TransactorSession struct {
	Contract     *Create3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// Create3Raw is an auto generated low-level Go binding around an Ethereum contract.
type Create3Raw struct {
	Contract *Create3 // Generic contract binding to access the raw methods on
}

// Create3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Create3CallerRaw struct {
	Contract *Create3Caller // Generic read-only contract binding to access the raw methods on
}

// Create3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Create3TransactorRaw struct {
	Contract *Create3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewCreate3 creates a new instance of Create3, bound to a specific deployed contract.
func NewCreate3(address common.Address, backend bind.ContractBackend) (*Create3, error) {
	contract, err := bindCreate3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Create3{Create3Caller: Create3Caller{contract: contract}, Create3Transactor: Create3Transactor{contract: contract}, Create3Filterer: Create3Filterer{contract: contract}}, nil
}

// NewCreate3Caller creates a new read-only instance of Create3, bound to a specific deployed contract.
func NewCreate3Caller(address common.Address, caller bind.ContractCaller) (*Create3Caller, error) {
	contract, err := bindCreate3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Create3Caller{contract: contract}, nil
}

// NewCreate3Transactor creates a new write-only instance of Create3, bound to a specific deployed contract.
func NewCreate3Transactor(address common.Address, transactor bind.ContractTransactor) (*Create3Transactor, error) {
	contract, err := bindCreate3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Create3Transactor{contract: contract}, nil
}

// NewCreate3Filterer creates a new log filterer instance of Create3, bound to a specific deployed contract.
func NewCreate3Filterer(address common.Address, filterer bind.ContractFilterer) (*Create3Filterer, error) {
	contract, err := bindCreate3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Create3Filterer{contract: contract}, nil
}

// bindCreate3 binds a generic wrapper to an already deployed contract.
func bindCreate3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Create3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Create3 *Create3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Create3.Contract.Create3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Create3 *Create3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Create3.Contract.Create3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Create3 *Create3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Create3.Contract.Create3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Create3 *Create3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Create3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Create3 *Create3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Create3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Create3 *Create3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Create3.Contract.contract.Transact(opts, method, params...)
}

// GetDeployed is a free data retrieval call binding the contract method 0x50f1c464.
//
// Solidity: function getDeployed(address deployer, bytes32 salt) view returns(address deployed)
func (_Create3 *Create3Caller) GetDeployed(opts *bind.CallOpts, deployer common.Address, salt [32]byte) (common.Address, error) {
	var out []interface{}
	err := _Create3.contract.Call(opts, &out, "getDeployed", deployer, salt)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDeployed is a free data retrieval call binding the contract method 0x50f1c464.
//
// Solidity: function getDeployed(address deployer, bytes32 salt) view returns(address deployed)
func (_Create3 *Create3Session) GetDeployed(deployer common.Address, salt [32]byte) (common.Address, error) {
	return _Create3.Contract.GetDeployed(&_Create3.CallOpts, deployer, salt)
}

// GetDeployed is a free data retrieval call binding the contract method 0x50f1c464.
//
// Solidity: function getDeployed(address deployer, bytes32 salt) view returns(address deployed)
func (_Create3 *Create3CallerSession) GetDeployed(deployer common.Address, salt [32]byte) (common.Address, error) {
	return _Create3.Contract.GetDeployed(&_Create3.CallOpts, deployer, salt)
}

// PredictDeterministicAddress is a free data retrieval call binding the contract method 0x5414dff0.
//
// Solidity: function predictDeterministicAddress(bytes32 salt) view returns(address deployed)
func (_Create3 *Create3Caller) PredictDeterministicAddress(opts *bind.CallOpts, salt [32]byte) (common.Address, error) {
	var out []interface{}
	err := _Create3.contract.Call(opts, &out, "predictDeterministicAddress", salt)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PredictDeterministicAddress is a free data retrieval call binding the contract method 0x5414dff0.
//
// Solidity: function predictDeterministicAddress(bytes32 salt) view returns(address deployed)
func (_Create3 *Create3Session) PredictDeterministicAddress(salt [32]byte) (common.Address, error) {
	return _Create3.Contract.PredictDeterministicAddress(&_Create3.CallOpts, salt)
}

// PredictDeterministicAddress is a free data retrieval call binding the contract method 0x5414dff0.
//
// Solidity: function predictDeterministicAddress(bytes32 salt) view returns(address deployed)
func (_Create3 *Create3CallerSession) PredictDeterministicAddress(salt [32]byte) (common.Address, error) {
	return _Create3.Contract.PredictDeterministicAddress(&_Create3.CallOpts, salt)
}

// Deploy is a paid mutator transaction binding the contract method 0xcdcb760a.
//
// Solidity: function deploy(bytes32 salt, bytes creationCode) payable returns(address deployed)
func (_Create3 *Create3Transactor) Deploy(opts *bind.TransactOpts, salt [32]byte, creationCode []byte) (*types.Transaction, error) {
	return _Create3.contract.Transact(opts, "deploy", salt, creationCode)
}

// Deploy is a paid mutator transaction binding the contract method 0xcdcb760a.
//
// Solidity: function deploy(bytes32 salt, bytes creationCode) payable returns(address deployed)
func (_Create3 *Create3Session) Deploy(salt [32]byte, creationCode []byte) (*types.Transaction, error) {
	return _Create3.Contract.Deploy(&_Create3.TransactOpts, salt, creationCode)
}

// Deploy is a paid mutator transaction binding the contract method 0xcdcb760a.
//
// Solidity: function deploy(bytes32 salt, bytes creationCode) payable returns(address deployed)
func (_Create3 *Create3TransactorSession) Deploy(salt [32]byte, creationCode []byte) (*types.Transaction, error) {
	return _Create3.Contract.Deploy(&_Create3.TransactOpts, salt, creationCode)
}

// DeployDeterministic is a paid mutator transaction binding the contract method 0x9881d195.
//
// Solidity: function deployDeterministic(bytes creationCode, bytes32 salt) payable returns(address deployed)
func (_Create3 *Create3Transactor) DeployDeterministic(opts *bind.TransactOpts, creationCode []byte, salt [32]byte) (*types.Transaction, error) {
	return _Create3.contract.Transact(opts, "deployDeterministic", creationCode, salt)
}

// DeployDeterministic is a paid mutator transaction binding the contract method 0x9881d195.
//
// Solidity: function deployDeterministic(bytes creationCode, bytes32 salt) payable returns(address deployed)
func (_Create3 *Create3Session) DeployDeterministic(creationCode []byte, salt [32]byte) (*types.Transaction, error) {
	return _Create3.Contract.DeployDeterministic(&_Create3.TransactOpts, creationCode, salt)
}

// DeployDeterministic is a paid mutator transaction binding the contract method 0x9881d195.
//
// Solidity: function deployDeterministic(bytes creationCode, bytes32 salt) payable returns(address deployed)
func (_Create3 *Create3TransactorSession) DeployDeterministic(creationCode []byte, salt [32]byte) (*types.Transaction, error) {
	return _Create3.Contract.DeployDeterministic(&_Create3.TransactOpts, creationCode, salt)
}
