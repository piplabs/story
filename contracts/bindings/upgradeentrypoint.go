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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"planUpgrade\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"info\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SoftwareUpgrade\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"int64\",\"indexed\":false,\"internalType\":\"int64\"},{\"name\":\"info\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161053f38038061053f83398101604081905261002f916100da565b806001600160a01b03811661005e57604051631e4fbdf760e01b81526000600482015260240160405180910390fd5b6100678161006e565b505061010a565b600180546001600160a01b03191690556100878161008a565b50565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6000602082840312156100ec57600080fd5b81516001600160a01b038116811461010357600080fd5b9392505050565b610426806101196000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c8063715018a61461006757806379ba5097146100715780638da5cb5b14610079578063e30c3978146100a2578063ef176e0e146100b3578063f2fde38b146100c6575b600080fd5b61006f6100d9565b005b61006f6100ed565b6000546001600160a01b03165b6040516001600160a01b03909116815260200160405180910390f35b6001546001600160a01b0316610086565b61006f6100c13660046102cf565b610136565b61006f6100d436600461035b565b610184565b6100e16101f5565b6100eb6000610222565b565b60015433906001600160a01b0316811461012a5760405163118cdaa760e01b81526001600160a01b03821660048201526024015b60405180910390fd5b61013381610222565b50565b61013e6101f5565b7f112749e79b2098b58eab36c21f123b2883c3ecbbb4f41623a744fa6d9b3e37c685858585856040516101759594939291906103b4565b60405180910390a15050505050565b61018c6101f5565b600180546001600160a01b0383166001600160a01b031990911681179091556101bd6000546001600160a01b031690565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a350565b6000546001600160a01b031633146100eb5760405163118cdaa760e01b8152336004820152602401610121565b600180546001600160a01b031916905561013381600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b60008083601f84011261029857600080fd5b50813567ffffffffffffffff8111156102b057600080fd5b6020830191508360208285010111156102c857600080fd5b9250929050565b6000806000806000606086880312156102e757600080fd5b853567ffffffffffffffff808211156102ff57600080fd5b61030b89838a01610286565b909750955060208801359150600782900b821461032757600080fd5b9093506040870135908082111561033d57600080fd5b5061034a88828901610286565b969995985093965092949392505050565b60006020828403121561036d57600080fd5b81356001600160a01b038116811461038457600080fd5b9392505050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6060815260006103c860608301878961038b565b8560070b602084015282810360408401526103e481858761038b565b9897505050505050505056fea2646970667358221220f12e8408869738cf77e756c57f3dbd8115f9eea757a287364347b7f1ff36591964736f6c63430008180033",
}

// UpgradeEntrypointABI is the input ABI used to generate the binding from.
// Deprecated: Use UpgradeEntrypointMetaData.ABI instead.
var UpgradeEntrypointABI = UpgradeEntrypointMetaData.ABI

// UpgradeEntrypointBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UpgradeEntrypointMetaData.Bin instead.
var UpgradeEntrypointBin = UpgradeEntrypointMetaData.Bin

// DeployUpgradeEntrypoint deploys a new Ethereum contract, binding an instance of UpgradeEntrypoint to it.
func DeployUpgradeEntrypoint(auth *bind.TransactOpts, backend bind.ContractBackend, newOwner common.Address) (common.Address, *types.Transaction, *UpgradeEntrypoint, error) {
	parsed, err := UpgradeEntrypointMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpgradeEntrypointBin), backend, newOwner)
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
