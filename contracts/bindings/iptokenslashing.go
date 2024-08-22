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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"ipTokenStaking\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newUnjailFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"IP_TOKEN_STAKING\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPTokenStaking\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"accessManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnjailFee\",\"inputs\":[{\"name\":\"newUnjailFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unjail\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unjailFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unjailOnBehalf\",\"inputs\":[{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unjail\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnjailFeeSet\",\"inputs\":[{\"name\":\"newUnjailFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60c0604052306080523480156200001557600080fd5b5060405162001794380380620017948339810160408190526200003891620001e2565b600081116200009a5760405162461bcd60e51b815260206004820152602360248201527f4950546f6b656e536c617368696e673a20496e76616c696420756e6a61696c2060448201526266656560e81b60648201526084015b60405180910390fd5b6001600160a01b0382166200010a5760405162461bcd60e51b815260206004820152602f60248201527f4950546f6b656e536c617368696e673a20496e76616c6964204950546f6b656e60448201526e5374616b696e67206164647265737360881b606482015260840162000091565b6001600160a01b03821660a0526000819055620001266200012e565b50506200021e565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156200017f5760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b0390811614620001df5780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b60008060408385031215620001f657600080fd5b82516001600160a01b03811681146200020e57600080fd5b6020939093015192949293505050565b60805160a05161153c620002586000396000818160f4015261088b015260008181610a4201528181610a6b0152610bb1015261153c6000f3fe6080604052600436106100dd5760003560e01c806379ba50971161007f578063c4d66de811610059578063c4d66de814610231578063e30c397814610251578063e4dfccd814610266578063f2fde38b1461027957600080fd5b806379ba5097146101c95780638da5cb5b146101de578063ad3cb1cc146101f357600080fd5b806340eda14a116100bb57806340eda14a146101795780634f1ef2861461018c57806352d1902d1461019f578063715018a6146101b457600080fd5b806304ff53ed146100e25780630c863f77146101335780632801f1ec14610155575b600080fd5b3480156100ee57600080fd5b506101167f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020015b60405180910390f35b34801561013f57600080fd5b5061015361014e36600461107f565b610299565b005b34801561016157600080fd5b5061016b60005481565b60405190815260200161012a565b610153610187366004611098565b6102dc565b61015361019a366004611195565b61035e565b3480156101ab57600080fd5b5061016b610379565b3480156101c057600080fd5b50610153610396565b3480156101d557600080fd5b506101536103aa565b3480156101ea57600080fd5b506101166103f7565b3480156101ff57600080fd5b50610224604051806040016040528060058152602001640352e302e360dc1b81525081565b60405161012a9190611276565b34801561023d57600080fd5b5061015361024c366004611289565b61042c565b34801561025d57600080fd5b506101166105b8565b610153610274366004611098565b6105e1565b34801561028557600080fd5b50610153610294366004611289565b610724565b6102a16107a9565b60008190556040518181527feac81de2f20162b0540ca5d3f43896af15b471a55729ff0c000e611d8b2723639060200160405180910390a150565b61031b82828080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506107db92505050565b61035a82828080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061096b92505050565b5050565b610366610a37565b61036f82610adc565b61035a8282610ae4565b6000610383610ba6565b506000805160206114e783398151915290565b61039e6107a9565b6103a86000610bef565b565b33806103b46105b8565b6001600160a01b0316146103eb5760405163118cdaa760e01b81526001600160a01b03821660048201526024015b60405180910390fd5b6103f481610bef565b50565b6000807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b546001600160a01b031692915050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008054600160401b810460ff16159067ffffffffffffffff166000811580156104725750825b905060008267ffffffffffffffff16600114801561048f5750303b155b90508115801561049d575080155b156104bb5760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156104e557845460ff60401b1916600160401b1785555b6001600160a01b0386166105595760405162461bcd60e51b815260206004820152603560248201527f4950546f6b656e536c617368696e673a206163636573734d616e616765722063604482015274616e6e6f74206265207a65726f206164647265737360581b60648201526084016103e2565b610561610c27565b61056a86610c2f565b83156105b057845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b6000807f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0061041c565b818133604182146106045760405162461bcd60e51b81526004016103e2906112a4565b82826000818110610617576106176112ea565b9050013560f81c60f81b6001600160f81b031916600460f81b1461064d5760405162461bcd60e51b81526004016103e290611300565b806001600160a01b03166106618484610c40565b6001600160a01b0316146106cf5760405162461bcd60e51b815260206004820152602f60248201527f4950546f6b656e536c617368696e673a20496e76616c6964207075626b65792060448201526e64657269766564206164647265737360881b60648201526084016103e2565b600061071086868080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610c6f92505050565b905061071b816107db565b6105b08161096b565b61072c6107a9565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0080546001600160a01b0319166001600160a01b03831690811782556107706103f7565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b336107b26103f7565b6001600160a01b0316146103a85760405163118cdaa760e01b81523360048201526024016103e2565b80516021146107fc5760405162461bcd60e51b81526004016103e2906112a4565b8060008151811061080f5761080f6112ea565b6020910101516001600160f81b031916600160f91b148061085557508060008151811061083e5761083e6112ea565b6020910101516001600160f81b031916600360f81b145b6108715760405162461bcd60e51b81526004016103e290611300565b604051638d3e1e4160e01b81526000906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690638d3e1e41906108c0908590600401611276565b600060405180830381865afa1580156108dd573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610905919081019061135a565b505050505090508061035a5760405162461bcd60e51b815260206004820152602960248201527f4950546f6b656e536c617368696e673a2056616c696461746f7220646f6573206044820152681b9bdd08195e1a5cdd60ba1b60648201526084016103e2565b60005434146109c65760405162461bcd60e51b815260206004820152602160248201527f4950546f6b656e536c617368696e673a20496e73756666696369656e742066656044820152606560f81b60648201526084016103e2565b6040516000903480156108fc029183818181858288f193505050501580156109f2573d6000803e3d6000fd5b50336001600160a01b03167f4a90ea32527ecacc0f4b32b31f99e4c633a2b4fe81ea7444989e2e68bc9ece3b82604051610a2c9190611276565b60405180910390a250565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480610abe57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610ab26000805160206114e7833981519152546001600160a01b031690565b6001600160a01b031614155b156103a85760405163703e46dd60e11b815260040160405180910390fd5b6103f46107a9565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610b3e575060408051601f3d908101601f19168201909252610b3b91810190611426565b60015b610b6657604051634c9c8ce360e01b81526001600160a01b03831660048201526024016103e2565b6000805160206114e78339815191528114610b9757604051632a87526960e21b8152600481018290526024016103e2565b610ba18383610dbb565b505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146103a85760405163703e46dd60e11b815260040160405180910390fd5b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0080546001600160a01b031916815561035a82610e11565b6103a8610e82565b610c37610e82565b6103f481610ecb565b6000610c4f826001818661143f565b604051610c5d929190611469565b60405190819003902090505b92915050565b60608151604114610cd15760405162461bcd60e51b815260206004820152602660248201527f496e76616c696420756e636f6d70726573736564207075626c6963206b6579206044820152650d8cadccee8d60d31b60648201526084016103e2565b602182015160418301516000610ceb600260ff8416611479565b60ff1615610cfd57600360f81b610d03565b600160f91b5b60408051602180825260608201909252919250600091906020820181803683370190505090508181600081518110610d3d57610d3d6112ea565b60200101906001600160f81b031916908160001a90535060005b6020811015610db157848160208110610d7257610d726112ea565b1a60f81b82610d828360016114a9565b81518110610d9257610d926112ea565b60200101906001600160f81b031916908160001a905350600101610d57565b5095945050505050565b610dc482610efd565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a2805115610e0957610ba18282610f62565b61035a610fd8565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054600160401b900460ff166103a857604051631afcd79f60e31b815260040160405180910390fd5b610ed3610e82565b6001600160a01b0381166103eb57604051631e4fbdf760e01b8152600060048201526024016103e2565b806001600160a01b03163b600003610f3357604051634c9c8ce360e01b81526001600160a01b03821660048201526024016103e2565b6000805160206114e783398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b031684604051610f7f91906114ca565b600060405180830381855af49150503d8060008114610fba576040519150601f19603f3d011682016040523d82523d6000602084013e610fbf565b606091505b5091509150610fcf858383610ff7565b95945050505050565b34156103a85760405163b398979f60e01b815260040160405180910390fd5b60608261100c5761100782611056565b61104f565b815115801561102357506001600160a01b0384163b155b1561104c57604051639996b31560e01b81526001600160a01b03851660048201526024016103e2565b50805b9392505050565b8051156110665780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b60006020828403121561109157600080fd5b5035919050565b600080602083850312156110ab57600080fd5b823567ffffffffffffffff808211156110c357600080fd5b818501915085601f8301126110d757600080fd5b8135818111156110e657600080fd5b8660208285010111156110f857600080fd5b60209290920196919550909350505050565b80356001600160a01b038116811461112157600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561116557611165611126565b604052919050565b600067ffffffffffffffff82111561118757611187611126565b50601f01601f191660200190565b600080604083850312156111a857600080fd5b6111b18361110a565b9150602083013567ffffffffffffffff8111156111cd57600080fd5b8301601f810185136111de57600080fd5b80356111f16111ec8261116d565b61113c565b81815286602083850101111561120657600080fd5b816020840160208301376000602083830101528093505050509250929050565b60005b83811015611241578181015183820152602001611229565b50506000910152565b60008151808452611262816020860160208601611226565b601f01601f19169290920160200192915050565b60208152600061104f602083018461124a565b60006020828403121561129b57600080fd5b61104f8261110a565b60208082526026908201527f4950546f6b656e536c617368696e673a20496e76616c6964207075626b6579206040820152650d8cadccee8d60d31b606082015260800190565b634e487b7160e01b600052603260045260246000fd5b60208082526026908201527f4950546f6b656e536c617368696e673a20496e76616c6964207075626b6579206040820152650e0e4caccd2f60d31b606082015260800190565b805163ffffffff8116811461112157600080fd5b60008060008060008060c0878903121561137357600080fd5b8651801515811461138357600080fd5b602088015190965067ffffffffffffffff8111156113a057600080fd5b8701601f810189136113b157600080fd5b80516113bf6111ec8261116d565b8181528a60208385010111156113d457600080fd5b6113e5826020830160208601611226565b809750505050604087015193506113fe60608801611346565b925061140c60808801611346565b915061141a60a08801611346565b90509295509295509295565b60006020828403121561143857600080fd5b5051919050565b6000808585111561144f57600080fd5b8386111561145c57600080fd5b5050820193919092039150565b8183823760009101908152919050565b600060ff83168061149a57634e487b7160e01b600052601260045260246000fd5b8060ff84160691505092915050565b80820180821115610c6957634e487b7160e01b600052601160045260246000fd5b600082516114dc818460208701611226565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca2646970667358221220144587c59cbacff3f2fd676bff35eb96ab661aedce15a1a06cdaf920eda0f2cc64736f6c63430008180033",
}

// IPTokenSlashingABI is the input ABI used to generate the binding from.
// Deprecated: Use IPTokenSlashingMetaData.ABI instead.
var IPTokenSlashingABI = IPTokenSlashingMetaData.ABI

// IPTokenSlashingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use IPTokenSlashingMetaData.Bin instead.
var IPTokenSlashingBin = IPTokenSlashingMetaData.Bin

// DeployIPTokenSlashing deploys a new Ethereum contract, binding an instance of IPTokenSlashing to it.
func DeployIPTokenSlashing(auth *bind.TransactOpts, backend bind.ContractBackend, ipTokenStaking common.Address, newUnjailFee *big.Int) (common.Address, *types.Transaction, *IPTokenSlashing, error) {
	parsed, err := IPTokenSlashingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(IPTokenSlashingBin), backend, ipTokenStaking, newUnjailFee)
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

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_IPTokenSlashing *IPTokenSlashingCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IPTokenSlashing.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_IPTokenSlashing *IPTokenSlashingSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _IPTokenSlashing.Contract.UPGRADEINTERFACEVERSION(&_IPTokenSlashing.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_IPTokenSlashing *IPTokenSlashingCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _IPTokenSlashing.Contract.UPGRADEINTERFACEVERSION(&_IPTokenSlashing.CallOpts)
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

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IPTokenSlashing *IPTokenSlashingCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IPTokenSlashing.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IPTokenSlashing *IPTokenSlashingSession) ProxiableUUID() ([32]byte, error) {
	return _IPTokenSlashing.Contract.ProxiableUUID(&_IPTokenSlashing.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IPTokenSlashing *IPTokenSlashingCallerSession) ProxiableUUID() ([32]byte, error) {
	return _IPTokenSlashing.Contract.ProxiableUUID(&_IPTokenSlashing.CallOpts)
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

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address accessManager) returns()
func (_IPTokenSlashing *IPTokenSlashingTransactor) Initialize(opts *bind.TransactOpts, accessManager common.Address) (*types.Transaction, error) {
	return _IPTokenSlashing.contract.Transact(opts, "initialize", accessManager)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address accessManager) returns()
func (_IPTokenSlashing *IPTokenSlashingSession) Initialize(accessManager common.Address) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.Initialize(&_IPTokenSlashing.TransactOpts, accessManager)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address accessManager) returns()
func (_IPTokenSlashing *IPTokenSlashingTransactorSession) Initialize(accessManager common.Address) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.Initialize(&_IPTokenSlashing.TransactOpts, accessManager)
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

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_IPTokenSlashing *IPTokenSlashingTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _IPTokenSlashing.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_IPTokenSlashing *IPTokenSlashingSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.UpgradeToAndCall(&_IPTokenSlashing.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_IPTokenSlashing *IPTokenSlashingTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _IPTokenSlashing.Contract.UpgradeToAndCall(&_IPTokenSlashing.TransactOpts, newImplementation, data)
}

// IPTokenSlashingInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the IPTokenSlashing contract.
type IPTokenSlashingInitializedIterator struct {
	Event *IPTokenSlashingInitialized // Event containing the contract specifics and raw log

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
func (it *IPTokenSlashingInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenSlashingInitialized)
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
		it.Event = new(IPTokenSlashingInitialized)
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
func (it *IPTokenSlashingInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenSlashingInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenSlashingInitialized represents a Initialized event raised by the IPTokenSlashing contract.
type IPTokenSlashingInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IPTokenSlashing *IPTokenSlashingFilterer) FilterInitialized(opts *bind.FilterOpts) (*IPTokenSlashingInitializedIterator, error) {

	logs, sub, err := _IPTokenSlashing.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &IPTokenSlashingInitializedIterator{contract: _IPTokenSlashing.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IPTokenSlashing *IPTokenSlashingFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *IPTokenSlashingInitialized) (event.Subscription, error) {

	logs, sub, err := _IPTokenSlashing.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenSlashingInitialized)
				if err := _IPTokenSlashing.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_IPTokenSlashing *IPTokenSlashingFilterer) ParseInitialized(log types.Log) (*IPTokenSlashingInitialized, error) {
	event := new(IPTokenSlashingInitialized)
	if err := _IPTokenSlashing.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

// IPTokenSlashingUnjailFeeSetIterator is returned from FilterUnjailFeeSet and is used to iterate over the raw logs and unpacked data for UnjailFeeSet events raised by the IPTokenSlashing contract.
type IPTokenSlashingUnjailFeeSetIterator struct {
	Event *IPTokenSlashingUnjailFeeSet // Event containing the contract specifics and raw log

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
func (it *IPTokenSlashingUnjailFeeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenSlashingUnjailFeeSet)
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
		it.Event = new(IPTokenSlashingUnjailFeeSet)
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
func (it *IPTokenSlashingUnjailFeeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenSlashingUnjailFeeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenSlashingUnjailFeeSet represents a UnjailFeeSet event raised by the IPTokenSlashing contract.
type IPTokenSlashingUnjailFeeSet struct {
	NewUnjailFee *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUnjailFeeSet is a free log retrieval operation binding the contract event 0xeac81de2f20162b0540ca5d3f43896af15b471a55729ff0c000e611d8b272363.
//
// Solidity: event UnjailFeeSet(uint256 newUnjailFee)
func (_IPTokenSlashing *IPTokenSlashingFilterer) FilterUnjailFeeSet(opts *bind.FilterOpts) (*IPTokenSlashingUnjailFeeSetIterator, error) {

	logs, sub, err := _IPTokenSlashing.contract.FilterLogs(opts, "UnjailFeeSet")
	if err != nil {
		return nil, err
	}
	return &IPTokenSlashingUnjailFeeSetIterator{contract: _IPTokenSlashing.contract, event: "UnjailFeeSet", logs: logs, sub: sub}, nil
}

// WatchUnjailFeeSet is a free log subscription operation binding the contract event 0xeac81de2f20162b0540ca5d3f43896af15b471a55729ff0c000e611d8b272363.
//
// Solidity: event UnjailFeeSet(uint256 newUnjailFee)
func (_IPTokenSlashing *IPTokenSlashingFilterer) WatchUnjailFeeSet(opts *bind.WatchOpts, sink chan<- *IPTokenSlashingUnjailFeeSet) (event.Subscription, error) {

	logs, sub, err := _IPTokenSlashing.contract.WatchLogs(opts, "UnjailFeeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenSlashingUnjailFeeSet)
				if err := _IPTokenSlashing.contract.UnpackLog(event, "UnjailFeeSet", log); err != nil {
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

// ParseUnjailFeeSet is a log parse operation binding the contract event 0xeac81de2f20162b0540ca5d3f43896af15b471a55729ff0c000e611d8b272363.
//
// Solidity: event UnjailFeeSet(uint256 newUnjailFee)
func (_IPTokenSlashing *IPTokenSlashingFilterer) ParseUnjailFeeSet(log types.Log) (*IPTokenSlashingUnjailFeeSet, error) {
	event := new(IPTokenSlashingUnjailFeeSet)
	if err := _IPTokenSlashing.contract.UnpackLog(event, "UnjailFeeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenSlashingUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the IPTokenSlashing contract.
type IPTokenSlashingUpgradedIterator struct {
	Event *IPTokenSlashingUpgraded // Event containing the contract specifics and raw log

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
func (it *IPTokenSlashingUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenSlashingUpgraded)
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
		it.Event = new(IPTokenSlashingUpgraded)
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
func (it *IPTokenSlashingUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenSlashingUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenSlashingUpgraded represents a Upgraded event raised by the IPTokenSlashing contract.
type IPTokenSlashingUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IPTokenSlashing *IPTokenSlashingFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*IPTokenSlashingUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IPTokenSlashing.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &IPTokenSlashingUpgradedIterator{contract: _IPTokenSlashing.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IPTokenSlashing *IPTokenSlashingFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *IPTokenSlashingUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IPTokenSlashing.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenSlashingUpgraded)
				if err := _IPTokenSlashing.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_IPTokenSlashing *IPTokenSlashingFilterer) ParseUpgraded(log types.Log) (*IPTokenSlashingUpgraded, error) {
	event := new(IPTokenSlashingUpgraded)
	if err := _IPTokenSlashing.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
