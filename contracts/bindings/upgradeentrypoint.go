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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelUpgrade\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"planUpgrade\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"info\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CancelUpgrade\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SoftwareUpgrade\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"height\",\"type\":\"int64\",\"indexed\":false,\"internalType\":\"int64\"},{\"name\":\"info\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608080604052346100b8577ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c166100a957506001600160401b036002600160401b031982821601610064575b60405161096690816100bd8239f35b6001600160401b031990911681179091556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f8080610055565b63f92ee8a960e01b8152600490fd5b5f80fdfe60406080815260049081361015610014575f80fd5b5f3560e01c90816355f29166146106ec578163715018a61461060957816379ba5097146105625781638da5cb5b146104f2578163c4d66de8146102a7578163e30c397814610237578163ef176e0e14610163575063f2fde38b14610076575f80fd5b3461015f5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015f573573ffffffffffffffffffffffffffffffffffffffff80821680920361015f576100cd6107b4565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227005f80a3005b5f80fd5b90503461015f5760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015f5767ffffffffffffffff91803583811161015f576101b49036908301610748565b919092602435908160070b80920361015f5760443595861161015f57610201610232937f112749e79b2098b58eab36c21f123b2883c3ecbbb4f41623a744fa6d9b3e37c697369101610748565b9161020a6107b4565b6102208151978897606089526060890191610776565b93602087015285840390860152610776565b0390a1005b3461015f575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015f5760209073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0054169051908152f35b823461015f5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015f57803573ffffffffffffffffffffffffffffffffffffffff81169081810361015f577ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009283549260ff84871c16159367ffffffffffffffff8116801590816104ea575b60011490816104e0575b1590816104d7575b506104af578460017fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000831617875561047a575b50156103f7575061039e906103916108d7565b6103996108d7565b610824565b6103a457005b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d291817fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff602093541690555160018152a1005b60849060208651917f08c379a0000000000000000000000000000000000000000000000000000000008352820152602f60248201527f55706772616465456e747279706f696e743a206f776e65722063616e6e6f742060448201527f6265207a65726f206164647265737300000000000000000000000000000000006064820152fd5b7fffffffffffffffffffffffffffffffffffffffffffffff00000000000000000016680100000000000000011785558661037e565b8287517ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b9050158861034b565b303b159150610343565b869150610339565b3461015f575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015f5760209073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054169051908152f35b823461015f575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015f573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416036105da576105d833610824565b005b60249151907f118cdaa70000000000000000000000000000000000000000000000000000000082523390820152fd5b3461015f575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015f5761063f6107b4565b5f73ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b3461015f575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261015f576107226107b4565b7f812c36a273ff85c1871fc7c629fa4c010821a53f3a2492dcc0ea00a396b6a64f5f80a1005b9181601f8401121561015f5782359167ffffffffffffffff831161015f576020838186019501011161015f57565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe093818652868601375f8582860101520116010190565b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300541633036107f457565b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b7fffffffffffffffffffffffff0000000000000000000000000000000000000000907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008281541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549073ffffffffffffffffffffffffffffffffffffffff80931680948316179055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3565b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c161561090657565b60046040517fd7e6bcf8000000000000000000000000000000000000000000000000000000008152fdfea2646970667358221220aaa587b5c6ac27bc988f3b521a565399d9e24777b52755132da9ef779433814b64736f6c63430008170033",
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

// CancelUpgrade is a paid mutator transaction binding the contract method 0x55f29166.
//
// Solidity: function cancelUpgrade() returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactor) CancelUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeEntrypoint.contract.Transact(opts, "cancelUpgrade")
}

// CancelUpgrade is a paid mutator transaction binding the contract method 0x55f29166.
//
// Solidity: function cancelUpgrade() returns()
func (_UpgradeEntrypoint *UpgradeEntrypointSession) CancelUpgrade() (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.CancelUpgrade(&_UpgradeEntrypoint.TransactOpts)
}

// CancelUpgrade is a paid mutator transaction binding the contract method 0x55f29166.
//
// Solidity: function cancelUpgrade() returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactorSession) CancelUpgrade() (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.CancelUpgrade(&_UpgradeEntrypoint.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactor) Initialize(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _UpgradeEntrypoint.contract.Transact(opts, "initialize", owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointSession) Initialize(owner common.Address) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.Initialize(&_UpgradeEntrypoint.TransactOpts, owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner) returns()
func (_UpgradeEntrypoint *UpgradeEntrypointTransactorSession) Initialize(owner common.Address) (*types.Transaction, error) {
	return _UpgradeEntrypoint.Contract.Initialize(&_UpgradeEntrypoint.TransactOpts, owner)
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

// UpgradeEntrypointCancelUpgradeIterator is returned from FilterCancelUpgrade and is used to iterate over the raw logs and unpacked data for CancelUpgrade events raised by the UpgradeEntrypoint contract.
type UpgradeEntrypointCancelUpgradeIterator struct {
	Event *UpgradeEntrypointCancelUpgrade // Event containing the contract specifics and raw log

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
func (it *UpgradeEntrypointCancelUpgradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeEntrypointCancelUpgrade)
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
		it.Event = new(UpgradeEntrypointCancelUpgrade)
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
func (it *UpgradeEntrypointCancelUpgradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeEntrypointCancelUpgradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeEntrypointCancelUpgrade represents a CancelUpgrade event raised by the UpgradeEntrypoint contract.
type UpgradeEntrypointCancelUpgrade struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCancelUpgrade is a free log retrieval operation binding the contract event 0x812c36a273ff85c1871fc7c629fa4c010821a53f3a2492dcc0ea00a396b6a64f.
//
// Solidity: event CancelUpgrade()
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) FilterCancelUpgrade(opts *bind.FilterOpts) (*UpgradeEntrypointCancelUpgradeIterator, error) {

	logs, sub, err := _UpgradeEntrypoint.contract.FilterLogs(opts, "CancelUpgrade")
	if err != nil {
		return nil, err
	}
	return &UpgradeEntrypointCancelUpgradeIterator{contract: _UpgradeEntrypoint.contract, event: "CancelUpgrade", logs: logs, sub: sub}, nil
}

// WatchCancelUpgrade is a free log subscription operation binding the contract event 0x812c36a273ff85c1871fc7c629fa4c010821a53f3a2492dcc0ea00a396b6a64f.
//
// Solidity: event CancelUpgrade()
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) WatchCancelUpgrade(opts *bind.WatchOpts, sink chan<- *UpgradeEntrypointCancelUpgrade) (event.Subscription, error) {

	logs, sub, err := _UpgradeEntrypoint.contract.WatchLogs(opts, "CancelUpgrade")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeEntrypointCancelUpgrade)
				if err := _UpgradeEntrypoint.contract.UnpackLog(event, "CancelUpgrade", log); err != nil {
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

// ParseCancelUpgrade is a log parse operation binding the contract event 0x812c36a273ff85c1871fc7c629fa4c010821a53f3a2492dcc0ea00a396b6a64f.
//
// Solidity: event CancelUpgrade()
func (_UpgradeEntrypoint *UpgradeEntrypointFilterer) ParseCancelUpgrade(log types.Log) (*UpgradeEntrypointCancelUpgrade, error) {
	event := new(UpgradeEntrypointCancelUpgrade)
	if err := _UpgradeEntrypoint.contract.UnpackLog(event, "CancelUpgrade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
