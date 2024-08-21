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

// IPTokenSlashingMetaData contains all meta data concerning the IPTokenSlashing contract.
var IPTokenSlashingMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"ipTokenStaking\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newUnjailFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"IP_TOKEN_STAKING\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPTokenStaking\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnjailFee\",\"inputs\":[{\"name\":\"newUnjailFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unjail\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unjailFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unjailOnBehalf\",\"inputs\":[{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unjail\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610d2f380380610d2f83398101604081905261002f91610107565b826001600160a01b03811661005e57604051631e4fbdf760e01b81526000600482015260240160405180910390fd5b6100678161007f565b506001600160a01b0390911660805260025550610143565b600180546001600160a01b03191690556100988161009b565b50565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b80516001600160a01b038116811461010257600080fd5b919050565b60008060006060848603121561011c57600080fd5b610125846100eb565b9250610133602085016100eb565b9150604084015190509250925092565b608051610bcb6101646000396000818160a8015261059e0152610bcb6000f3fe6080604052600436106100915760003560e01c806379ba50971161005957806379ba5097146101555780638da5cb5b1461016a578063e30c397814610188578063e4dfccd8146101a6578063f2fde38b146101b957600080fd5b806304ff53ed146100965780630c863f77146100e75780632801f1ec1461010957806340eda14a1461012d578063715018a614610140575b600080fd5b3480156100a257600080fd5b506100ca7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020015b60405180910390f35b3480156100f357600080fd5b506101076101023660046107e0565b6101d9565b005b34801561011557600080fd5b5061011f60025481565b6040519081526020016100de565b61010761013b3660046107f9565b6101e6565b34801561014c57600080fd5b50610107610268565b34801561016157600080fd5b5061010761027c565b34801561017657600080fd5b506000546001600160a01b03166100ca565b34801561019457600080fd5b506001546001600160a01b03166100ca565b6101076101b43660046107f9565b6102c5565b3480156101c557600080fd5b506101076101d436600461086b565b610450565b6101e16104c1565b600255565b61022582828080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506104ee92505050565b61026482828080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061067e92505050565b5050565b6102706104c1565b61027a600061074a565b565b60015433906001600160a01b031681146102b95760405163118cdaa760e01b81526001600160a01b03821660048201526024015b60405180910390fd5b6102c28161074a565b50565b818133604182146102e85760405162461bcd60e51b81526004016102b09061089b565b828260008181106102fb576102fb6108e1565b9050013560f81c60f81b6001600160f81b031916600460f81b146103315760405162461bcd60e51b81526004016102b0906108f7565b806001600160a01b03166103458484610763565b6001600160a01b0316146103b35760405162461bcd60e51b815260206004820152602f60248201527f4950546f6b656e536c617368696e673a20496e76616c6964207075626b65792060448201526e64657269766564206164647265737360881b60648201526084016102b0565b604051633444d8b760e11b815260009073__$070efe90de6222b6182e3f0710b89d2262$__90636889b16e906103ef908990899060040161093d565b600060405180830381865af415801561040c573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526104349190810190610a1a565b905061043f816104ee565b6104488161067e565b505050505050565b6104586104c1565b600180546001600160a01b0383166001600160a01b031990911681179091556104896000546001600160a01b031690565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a350565b6000546001600160a01b0316331461027a5760405163118cdaa760e01b81523360048201526024016102b0565b805160211461050f5760405162461bcd60e51b81526004016102b09061089b565b80600081518110610522576105226108e1565b6020910101516001600160f81b031916600160f91b1480610568575080600081518110610551576105516108e1565b6020910101516001600160f81b031916600360f81b145b6105845760405162461bcd60e51b81526004016102b0906108f7565b604051638d3e1e4160e01b81526000906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690638d3e1e41906105d3908590600401610a6b565b600060405180830381865afa1580156105f0573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526106189190810190610ab7565b50505050509050806102645760405162461bcd60e51b815260206004820152602960248201527f4950546f6b656e536c617368696e673a2056616c696461746f7220646f6573206044820152681b9bdd08195e1a5cdd60ba1b60648201526084016102b0565b60025434146106d95760405162461bcd60e51b815260206004820152602160248201527f4950546f6b656e536c617368696e673a20496e73756666696369656e742066656044820152606560f81b60648201526084016102b0565b6040516000903480156108fc029183818181858288f19350505050158015610705573d6000803e3d6000fd5b50336001600160a01b03167f4a90ea32527ecacc0f4b32b31f99e4c633a2b4fe81ea7444989e2e68bc9ece3b8260405161073f9190610a6b565b60405180910390a250565b600180546001600160a01b03191690556102c281610790565b60006107728260018186610b5b565b604051610780929190610b85565b6040519081900390209392505050565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6000602082840312156107f257600080fd5b5035919050565b6000806020838503121561080c57600080fd5b823567ffffffffffffffff8082111561082457600080fd5b818501915085601f83011261083857600080fd5b81358181111561084757600080fd5b86602082850101111561085957600080fd5b60209290920196919550909350505050565b60006020828403121561087d57600080fd5b81356001600160a01b038116811461089457600080fd5b9392505050565b60208082526026908201527f4950546f6b656e536c617368696e673a20496e76616c6964207075626b6579206040820152650d8cadccee8d60d31b606082015260800190565b634e487b7160e01b600052603260045260246000fd5b60208082526026908201527f4950546f6b656e536c617368696e673a20496e76616c6964207075626b6579206040820152650e0e4caccd2f60d31b606082015260800190565b60208152816020820152818360408301376000818301604090810191909152601f909201601f19160101919050565b634e487b7160e01b600052604160045260246000fd5b60005b8381101561099d578181015183820152602001610985565b50506000910152565b600067ffffffffffffffff808411156109c1576109c161096c565b604051601f8501601f19908116603f011681019082821181831017156109e9576109e961096c565b81604052809350858152868686011115610a0257600080fd5b610a10866020830187610982565b5050509392505050565b600060208284031215610a2c57600080fd5b815167ffffffffffffffff811115610a4357600080fd5b8201601f81018413610a5457600080fd5b610a63848251602084016109a6565b949350505050565b6020815260008251806020840152610a8a816040850160208701610982565b601f01601f19169190910160400192915050565b805163ffffffff81168114610ab257600080fd5b919050565b60008060008060008060c08789031215610ad057600080fd5b86518015158114610ae057600080fd5b602088015190965067ffffffffffffffff811115610afd57600080fd5b8701601f81018913610b0e57600080fd5b610b1d898251602084016109a6565b95505060408701519350610b3360608801610a9e565b9250610b4160808801610a9e565b9150610b4f60a08801610a9e565b90509295509295509295565b60008085851115610b6b57600080fd5b83861115610b7857600080fd5b5050820193919092039150565b818382376000910190815291905056fea2646970667358221220ef1caa5e2adb7ee57adc7890bbe897a407f3d4c2912ee3a3a663fc33d6345b8c64736f6c63430008180033",
}

// IPTokenSlashingABI is the input ABI used to generate the binding from.
// Deprecated: Use IPTokenSlashingMetaData.ABI instead.
var IPTokenSlashingABI = IPTokenSlashingMetaData.ABI

// IPTokenSlashingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use IPTokenSlashingMetaData.Bin instead.
var IPTokenSlashingBin = IPTokenSlashingMetaData.Bin

// DeployIPTokenSlashing deploys a new Ethereum contract, binding an instance of IPTokenSlashing to it.
func DeployIPTokenSlashing(auth *bind.TransactOpts, backend bind.ContractBackend, newOwner common.Address, ipTokenStaking common.Address, newUnjailFee *big.Int) (common.Address, *types.Transaction, *IPTokenSlashing, error) {
	parsed, err := IPTokenSlashingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(IPTokenSlashingBin), backend, newOwner, ipTokenStaking, newUnjailFee)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IPTokenSlashing{IPTokenSlashingCaller: IPTokenSlashingCaller{contract: contract}, IPTokenSlashingTransactor: IPTokenSlashingTransactor{contract: contract}, IPTokenSlashingFilterer: IPTokenSlashingFilterer{contract: contract}}, nil
}

// IPTokenSlashing is an auto generated Go binding around an Ethereum contract.
type IPTokenSlashing struct {
	IPTokenSlashingCaller     // Read-only binding to the contract
	IPTokenSlashingTransactor // Write-only binding to the contract
	IPTokenSlashingFilterer   // Log filterer for contract events
}

// IPTokenSlashingCaller is an auto generated read-only Go binding around an Ethereum contract.
type IPTokenSlashingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPTokenSlashingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IPTokenSlashingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPTokenSlashingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IPTokenSlashingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPTokenSlashingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IPTokenSlashingSession struct {
	Contract     *IPTokenSlashing  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IPTokenSlashingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IPTokenSlashingCallerSession struct {
	Contract *IPTokenSlashingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IPTokenSlashingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IPTokenSlashingTransactorSession struct {
	Contract     *IPTokenSlashingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IPTokenSlashingRaw is an auto generated low-level Go binding around an Ethereum contract.
type IPTokenSlashingRaw struct {
	Contract *IPTokenSlashing // Generic contract binding to access the raw methods on
}

// IPTokenSlashingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IPTokenSlashingCallerRaw struct {
	Contract *IPTokenSlashingCaller // Generic read-only contract binding to access the raw methods on
}

// IPTokenSlashingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IPTokenSlashingTransactorRaw struct {
	Contract *IPTokenSlashingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIPTokenSlashing creates a new instance of IPTokenSlashing, bound to a specific deployed contract.
func NewIPTokenSlashing(address common.Address, backend bind.ContractBackend) (*IPTokenSlashing, error) {
	contract, err := bindIPTokenSlashing(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IPTokenSlashing{IPTokenSlashingCaller: IPTokenSlashingCaller{contract: contract}, IPTokenSlashingTransactor: IPTokenSlashingTransactor{contract: contract}, IPTokenSlashingFilterer: IPTokenSlashingFilterer{contract: contract}}, nil
}

// NewIPTokenSlashingCaller creates a new read-only instance of IPTokenSlashing, bound to a specific deployed contract.
func NewIPTokenSlashingCaller(address common.Address, caller bind.ContractCaller) (*IPTokenSlashingCaller, error) {
	contract, err := bindIPTokenSlashing(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IPTokenSlashingCaller{contract: contract}, nil
}

// NewIPTokenSlashingTransactor creates a new write-only instance of IPTokenSlashing, bound to a specific deployed contract.
func NewIPTokenSlashingTransactor(address common.Address, transactor bind.ContractTransactor) (*IPTokenSlashingTransactor, error) {
	contract, err := bindIPTokenSlashing(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IPTokenSlashingTransactor{contract: contract}, nil
}

// NewIPTokenSlashingFilterer creates a new log filterer instance of IPTokenSlashing, bound to a specific deployed contract.
func NewIPTokenSlashingFilterer(address common.Address, filterer bind.ContractFilterer) (*IPTokenSlashingFilterer, error) {
	contract, err := bindIPTokenSlashing(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IPTokenSlashingFilterer{contract: contract}, nil
}

// bindIPTokenSlashing binds a generic wrapper to an already deployed contract.
func bindIPTokenSlashing(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IPTokenSlashingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPTokenSlashing *IPTokenSlashingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPTokenSlashing.Contract.IPTokenSlashingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPTokenSlashing *IPTokenSlashingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.IPTokenSlashingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPTokenSlashing *IPTokenSlashingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.IPTokenSlashingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPTokenSlashing *IPTokenSlashingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPTokenSlashing.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPTokenSlashing *IPTokenSlashingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPTokenSlashing *IPTokenSlashingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.contract.Transact(opts, method, params...)
}

// IPTOKENSTAKING is a free data retrieval call binding the contract method 0x04ff53ed.
//
// Solidity: function IP_TOKEN_STAKING() view returns(address)
func (_IPTokenSlashing *IPTokenSlashingCaller) IPTOKENSTAKING(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPTokenSlashing.contract.Call(opts, &out, "IP_TOKEN_STAKING")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// IPTOKENSTAKING is a free data retrieval call binding the contract method 0x04ff53ed.
//
// Solidity: function IP_TOKEN_STAKING() view returns(address)
func (_IPTokenSlashing *IPTokenSlashingSession) IPTOKENSTAKING() (common.Address, error) {
	return _IPTokenSlashing.Contract.IPTOKENSTAKING(&_IPTokenSlashing.CallOpts)
}

// IPTOKENSTAKING is a free data retrieval call binding the contract method 0x04ff53ed.
//
// Solidity: function IP_TOKEN_STAKING() view returns(address)
func (_IPTokenSlashing *IPTokenSlashingCallerSession) IPTOKENSTAKING() (common.Address, error) {
	return _IPTokenSlashing.Contract.IPTOKENSTAKING(&_IPTokenSlashing.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IPTokenSlashing *IPTokenSlashingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPTokenSlashing.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IPTokenSlashing *IPTokenSlashingSession) Owner() (common.Address, error) {
	return _IPTokenSlashing.Contract.Owner(&_IPTokenSlashing.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IPTokenSlashing *IPTokenSlashingCallerSession) Owner() (common.Address, error) {
	return _IPTokenSlashing.Contract.Owner(&_IPTokenSlashing.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_IPTokenSlashing *IPTokenSlashingCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPTokenSlashing.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_IPTokenSlashing *IPTokenSlashingSession) PendingOwner() (common.Address, error) {
	return _IPTokenSlashing.Contract.PendingOwner(&_IPTokenSlashing.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_IPTokenSlashing *IPTokenSlashingCallerSession) PendingOwner() (common.Address, error) {
	return _IPTokenSlashing.Contract.PendingOwner(&_IPTokenSlashing.CallOpts)
}

// UnjailFee is a free data retrieval call binding the contract method 0x2801f1ec.
//
// Solidity: function unjailFee() view returns(uint256)
func (_IPTokenSlashing *IPTokenSlashingCaller) UnjailFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenSlashing.contract.Call(opts, &out, "unjailFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnjailFee is a free data retrieval call binding the contract method 0x2801f1ec.
//
// Solidity: function unjailFee() view returns(uint256)
func (_IPTokenSlashing *IPTokenSlashingSession) UnjailFee() (*big.Int, error) {
	return _IPTokenSlashing.Contract.UnjailFee(&_IPTokenSlashing.CallOpts)
}

// UnjailFee is a free data retrieval call binding the contract method 0x2801f1ec.
//
// Solidity: function unjailFee() view returns(uint256)
func (_IPTokenSlashing *IPTokenSlashingCallerSession) UnjailFee() (*big.Int, error) {
	return _IPTokenSlashing.Contract.UnjailFee(&_IPTokenSlashing.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_IPTokenSlashing *IPTokenSlashingTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPTokenSlashing.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_IPTokenSlashing *IPTokenSlashingSession) AcceptOwnership() (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.AcceptOwnership(&_IPTokenSlashing.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_IPTokenSlashing *IPTokenSlashingTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.AcceptOwnership(&_IPTokenSlashing.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IPTokenSlashing *IPTokenSlashingTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPTokenSlashing.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IPTokenSlashing *IPTokenSlashingSession) RenounceOwnership() (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.RenounceOwnership(&_IPTokenSlashing.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IPTokenSlashing *IPTokenSlashingTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.RenounceOwnership(&_IPTokenSlashing.TransactOpts)
}

// SetUnjailFee is a paid mutator transaction binding the contract method 0x0c863f77.
//
// Solidity: function setUnjailFee(uint256 newUnjailFee) returns()
func (_IPTokenSlashing *IPTokenSlashingTransactor) SetUnjailFee(opts *bind.TransactOpts, newUnjailFee *big.Int) (*types.Transaction, error) {
	return _IPTokenSlashing.contract.Transact(opts, "setUnjailFee", newUnjailFee)
}

// SetUnjailFee is a paid mutator transaction binding the contract method 0x0c863f77.
//
// Solidity: function setUnjailFee(uint256 newUnjailFee) returns()
func (_IPTokenSlashing *IPTokenSlashingSession) SetUnjailFee(newUnjailFee *big.Int) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.SetUnjailFee(&_IPTokenSlashing.TransactOpts, newUnjailFee)
}

// SetUnjailFee is a paid mutator transaction binding the contract method 0x0c863f77.
//
// Solidity: function setUnjailFee(uint256 newUnjailFee) returns()
func (_IPTokenSlashing *IPTokenSlashingTransactorSession) SetUnjailFee(newUnjailFee *big.Int) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.SetUnjailFee(&_IPTokenSlashing.TransactOpts, newUnjailFee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IPTokenSlashing *IPTokenSlashingTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _IPTokenSlashing.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IPTokenSlashing *IPTokenSlashingSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.TransferOwnership(&_IPTokenSlashing.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IPTokenSlashing *IPTokenSlashingTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.TransferOwnership(&_IPTokenSlashing.TransactOpts, newOwner)
}

// Unjail is a paid mutator transaction binding the contract method 0xe4dfccd8.
//
// Solidity: function unjail(bytes validatorUncmpPubkey) payable returns()
func (_IPTokenSlashing *IPTokenSlashingTransactor) Unjail(opts *bind.TransactOpts, validatorUncmpPubkey []byte) (*types.Transaction, error) {
	return _IPTokenSlashing.contract.Transact(opts, "unjail", validatorUncmpPubkey)
}

// Unjail is a paid mutator transaction binding the contract method 0xe4dfccd8.
//
// Solidity: function unjail(bytes validatorUncmpPubkey) payable returns()
func (_IPTokenSlashing *IPTokenSlashingSession) Unjail(validatorUncmpPubkey []byte) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.Unjail(&_IPTokenSlashing.TransactOpts, validatorUncmpPubkey)
}

// Unjail is a paid mutator transaction binding the contract method 0xe4dfccd8.
//
// Solidity: function unjail(bytes validatorUncmpPubkey) payable returns()
func (_IPTokenSlashing *IPTokenSlashingTransactorSession) Unjail(validatorUncmpPubkey []byte) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.Unjail(&_IPTokenSlashing.TransactOpts, validatorUncmpPubkey)
}

// UnjailOnBehalf is a paid mutator transaction binding the contract method 0x40eda14a.
//
// Solidity: function unjailOnBehalf(bytes validatorCmpPubkey) payable returns()
func (_IPTokenSlashing *IPTokenSlashingTransactor) UnjailOnBehalf(opts *bind.TransactOpts, validatorCmpPubkey []byte) (*types.Transaction, error) {
	return _IPTokenSlashing.contract.Transact(opts, "unjailOnBehalf", validatorCmpPubkey)
}

// UnjailOnBehalf is a paid mutator transaction binding the contract method 0x40eda14a.
//
// Solidity: function unjailOnBehalf(bytes validatorCmpPubkey) payable returns()
func (_IPTokenSlashing *IPTokenSlashingSession) UnjailOnBehalf(validatorCmpPubkey []byte) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.UnjailOnBehalf(&_IPTokenSlashing.TransactOpts, validatorCmpPubkey)
}

// UnjailOnBehalf is a paid mutator transaction binding the contract method 0x40eda14a.
//
// Solidity: function unjailOnBehalf(bytes validatorCmpPubkey) payable returns()
func (_IPTokenSlashing *IPTokenSlashingTransactorSession) UnjailOnBehalf(validatorCmpPubkey []byte) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.UnjailOnBehalf(&_IPTokenSlashing.TransactOpts, validatorCmpPubkey)
}

// IPTokenSlashingOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the IPTokenSlashing contract.
type IPTokenSlashingOwnershipTransferStartedIterator struct {
	Event *IPTokenSlashingOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *IPTokenSlashingOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenSlashingOwnershipTransferStarted)
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
		it.Event = new(IPTokenSlashingOwnershipTransferStarted)
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
func (it *IPTokenSlashingOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenSlashingOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenSlashingOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the IPTokenSlashing contract.
type IPTokenSlashingOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_IPTokenSlashing *IPTokenSlashingFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IPTokenSlashingOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IPTokenSlashing.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IPTokenSlashingOwnershipTransferStartedIterator{contract: _IPTokenSlashing.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_IPTokenSlashing *IPTokenSlashingFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *IPTokenSlashingOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IPTokenSlashing.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenSlashingOwnershipTransferStarted)
				if err := _IPTokenSlashing.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_IPTokenSlashing *IPTokenSlashingFilterer) ParseOwnershipTransferStarted(log types.Log) (*IPTokenSlashingOwnershipTransferStarted, error) {
	event := new(IPTokenSlashingOwnershipTransferStarted)
	if err := _IPTokenSlashing.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenSlashingOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the IPTokenSlashing contract.
type IPTokenSlashingOwnershipTransferredIterator struct {
	Event *IPTokenSlashingOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *IPTokenSlashingOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenSlashingOwnershipTransferred)
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
		it.Event = new(IPTokenSlashingOwnershipTransferred)
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
func (it *IPTokenSlashingOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenSlashingOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenSlashingOwnershipTransferred represents a OwnershipTransferred event raised by the IPTokenSlashing contract.
type IPTokenSlashingOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IPTokenSlashing *IPTokenSlashingFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IPTokenSlashingOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IPTokenSlashing.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IPTokenSlashingOwnershipTransferredIterator{contract: _IPTokenSlashing.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IPTokenSlashing *IPTokenSlashingFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IPTokenSlashingOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IPTokenSlashing.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenSlashingOwnershipTransferred)
				if err := _IPTokenSlashing.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_IPTokenSlashing *IPTokenSlashingFilterer) ParseOwnershipTransferred(log types.Log) (*IPTokenSlashingOwnershipTransferred, error) {
	event := new(IPTokenSlashingOwnershipTransferred)
	if err := _IPTokenSlashing.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenSlashingUnjailIterator is returned from FilterUnjail and is used to iterate over the raw logs and unpacked data for Unjail events raised by the IPTokenSlashing contract.
type IPTokenSlashingUnjailIterator struct {
	Event *IPTokenSlashingUnjail // Event containing the contract specifics and raw log

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
func (it *IPTokenSlashingUnjailIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenSlashingUnjail)
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
		it.Event = new(IPTokenSlashingUnjail)
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
func (it *IPTokenSlashingUnjailIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenSlashingUnjailIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenSlashingUnjail represents a Unjail event raised by the IPTokenSlashing contract.
type IPTokenSlashingUnjail struct {
	Sender             common.Address
	ValidatorCmpPubkey []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterUnjail is a free log retrieval operation binding the contract event 0x4a90ea32527ecacc0f4b32b31f99e4c633a2b4fe81ea7444989e2e68bc9ece3b.
//
// Solidity: event Unjail(address indexed sender, bytes validatorCmpPubkey)
func (_IPTokenSlashing *IPTokenSlashingFilterer) FilterUnjail(opts *bind.FilterOpts, sender []common.Address) (*IPTokenSlashingUnjailIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IPTokenSlashing.contract.FilterLogs(opts, "Unjail", senderRule)
	if err != nil {
		return nil, err
	}
	return &IPTokenSlashingUnjailIterator{contract: _IPTokenSlashing.contract, event: "Unjail", logs: logs, sub: sub}, nil
}

// WatchUnjail is a free log subscription operation binding the contract event 0x4a90ea32527ecacc0f4b32b31f99e4c633a2b4fe81ea7444989e2e68bc9ece3b.
//
// Solidity: event Unjail(address indexed sender, bytes validatorCmpPubkey)
func (_IPTokenSlashing *IPTokenSlashingFilterer) WatchUnjail(opts *bind.WatchOpts, sink chan<- *IPTokenSlashingUnjail, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IPTokenSlashing.contract.WatchLogs(opts, "Unjail", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenSlashingUnjail)
				if err := _IPTokenSlashing.contract.UnpackLog(event, "Unjail", log); err != nil {
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

// ParseUnjail is a log parse operation binding the contract event 0x4a90ea32527ecacc0f4b32b31f99e4c633a2b4fe81ea7444989e2e68bc9ece3b.
//
// Solidity: event Unjail(address indexed sender, bytes validatorCmpPubkey)
func (_IPTokenSlashing *IPTokenSlashingFilterer) ParseUnjail(log types.Log) (*IPTokenSlashingUnjail, error) {
	event := new(IPTokenSlashingUnjail)
	if err := _IPTokenSlashing.contract.UnpackLog(event, "Unjail", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
