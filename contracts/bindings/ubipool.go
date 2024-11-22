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

// UBIPoolMetaData contains all meta data concerning the UBIPool contract.
var UBIPoolMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"maxUBIPercentage\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_UBI_PERCENTAGE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimUBI\",\"inputs\":[{\"name\":\"distributionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentDistributionId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUBIDistribution\",\"inputs\":[{\"name\":\"totalUBI\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorUncmpPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUBIPercentage\",\"inputs\":[{\"name\":\"percentage\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validatorUBIAmounts\",\"inputs\":[{\"name\":\"distributionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UBIDistributionSet\",\"inputs\":[{\"name\":\"month\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalUBI\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"validatorUncmpPubKeys\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UBIPercentageSet\",\"inputs\":[{\"name\":\"percentage\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60a034610101576001600160401b0390601f6117ee38819003918201601f19168301918483118484101761010657808492602094604052833981010312610101575163ffffffff81168103610101576080527ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c166100ef5780808316036100aa575b6040516116d1908161011d82396080518181816101f60152610b690152f35b6001600160401b031990911681179091556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a138808061008b565b60405163f92ee8a960e01b8152600490fd5b600080fd5b634e487b7160e01b600052604160045260246000fdfe60406080815260048036101561001457600080fd5b600091823560e01c9081631971f77314610c1957816347564aa014610b3a578163715018a614610a70578163747c4ef7146107cd578163780069e0146107b057816379ba5097146107235781638da5cb5b146106cf578163ac9650d8146104cf578163c20c147214610453578163c4d66de81461021a578163d5077f40146101d9578163e30c397814610185575063f2fde38b146100b157600080fd5b34610181576020600319360112610181573573ffffffffffffffffffffffffffffffffffffffff80821680920361017d576100ea61131c565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227008380a380f35b8280fd5b5080fd5b83903461018157816003193601126101815760209073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0054169051908152f35b8390346101815781600319360112610181576020905163ffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b9190503461017d57602060031936011261017d57803573ffffffffffffffffffffffffffffffffffffffff81169081810361044f577ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009283549260ff84871c16159367ffffffffffffffff811680159081610447575b600114908161043d575b159081610434575b5061040c578460017fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000083161787556103d7575b501561035457506102f5906102e8611642565b6102f0611642565b6114ee565b6102fd578280f35b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d291817fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff602093541690555160018152a138808280f35b60849060208651917f08c379a0000000000000000000000000000000000000000000000000000000008352820152602560248201527f554249506f6f6c3a206f776e65722063616e6e6f74206265207a65726f20616460448201527f64726573730000000000000000000000000000000000000000000000000000006064820152fd5b7fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000001668010000000000000001178555386102d5565b8287517ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b905015386102a2565b303b15915061029a565b869150610290565b8480fd5b9190503461017d578160031936011261017d5760243567ffffffffffffffff81116104cb57366023820112156104cb576104ba918360209561049f8794369060248187013591016111d1565b92358152600184522082855194838680955193849201611104565b820190815203019020549051908152f35b8380fd5b8383346101815760208060031936011261017d5767ffffffffffffffff90823582811161044f5761050390369085016110ce565b9286519483860191868310908311176106a357508087939694975283855261052a876112f0565b9461053784519687611127565b878652610543886112f0565b977fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0809901865b818110610694575050855b81811061060f57505050505080519380850191818652845180935281818701918460051b880101950193965b8388106105ae5786860387f35b9091929394838080837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08b6001960301875285601f8b516105fa81518092818752878088019101611104565b011601019701930197019690939291936105a1565b806106738880896106538e9b9f9c9e61065f908b8b6106328f9b8d60019d611247565b9290965195838794868601998a37840191858301938a855251938491611104565b01038084520182611127565b5190305af461066c6112c0565b90306115a2565b61067d828b611308565b52610688818a611308565b50019894979598610575565b60608982018b0152890161056a565b8660416024927f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b83903461018157816003193601126101815760209073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054169051908152f35b90503461017d578260031936011261017d573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00541603610780578261077d336114ee565b80f35b6024925051907f118cdaa70000000000000000000000000000000000000000000000000000000082523390820152fd5b839034610181578160031936011261018157602091549051908152f35b9190503461017d578160031936011261017d57602480359267ffffffffffffffff908335828611610a6c5736602387011215610a6c5785850135928311610a6c5783860136858589010111610a68577f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00966002885414610a405760028855610855858361138c565b84600111610a3c5761088f90369060257fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff880191016111d1565b9384519473ffffffffffffffffffffffffffffffffffffffff602096873393012016036109bb57828952600185526108ca848a2082846112a7565b54928315610960578993846108ed819582958395845260018b52898420916112a7565b55335af16108f96112c0565b501561090757856001865580f35b517f08c379a000000000000000000000000000000000000000000000000000000000815292830152601b908201527f554249506f6f6c3a206661696c656420746f2073656e642055424900000000006044820152606490fd5b606488601889898951937f08c379a00000000000000000000000000000000000000000000000000000000085528401528201527f554249506f6f6c3a206e6f2055424920746f20636c61696d00000000000000006044820152fd5b608487602e88888851937f08c379a00000000000000000000000000000000000000000000000000000000085528401528201527f5075624b657956657269666965723a20496e76616c6964207075626b6579206460448201527f65726976656420616464726573730000000000000000000000000000000000006064820152fd5b8880fd5b8684517f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b8780fd5b8680fd5b8334610b375780600319360112610b3757610a8961131c565b8073ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b80fd5b90503461017d57602060031936011261017d5781359163ffffffff80841680940361044f57610b6761131c565b7f0000000000000000000000000000000000000000000000000000000000000000168311610bbd57507f6c6483041303ba314f169eb2d2af177b4f497324ccf0f3c1e68c2100f76c49299160209151908152a180f35b602060649251917f08c379a0000000000000000000000000000000000000000000000000000000008352820152601c60248201527f554249506f6f6c3a2070657263656e7461676520746f6f2068696768000000006044820152fd5b90503461017d57606060031936011261017d57813592602467ffffffffffffffff81358181116104cb57610c5090369087016110ce565b959093604491823584811161017d57610c6c90369083016110ce565b939095610c7761131c565b891561104d57848a03610ff257478b11610f97578384547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610f115760010185558b90855b878110610e99575003610e3e575050508096949296549686519580608088018a895260209b8c8a015260808a8a01525260a087019460a08260051b89010195819385925b848410610d7d5750505050505084830360608601528183527f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8211610b37575092809287927f1cc6f356308c8399caa490706b01fb9d52cdc87cdf639e66c3da7d4ce2db161c9560051b80928583013701030190a151908152f35b9091929394977fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff608b820301845288357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe183360301811215610a685782018035908f01848211610a3c578136038113610a3c578f837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8580859796869760019a52868601378d858286010152011601019a01940194019294939190610d02565b60649291601e7f554249506f6f6c3a20746f74616c20616d6f756e74206d69736d6174636800009260208c51957f08c379a0000000000000000000000000000000000000000000000000000000008752860152840152820152fd5b9150610ea682888b611208565b3515610f3c578a610ef18b610eeb8f80610ed88f8e8a9491610ed3610ecd8780958b611247565b9061138c565b611208565b35958c548d5260016020528c2093611247565b906112a7565b55610efd82888b611208565b358101809111610f11578c91600101610cbe565b82866011877f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b8a517f08c379a0000000000000000000000000000000000000000000000000000000008152602081870152601f818501527f554249506f6f6c3a20616d6f756e74732063616e6e6f74206265207a65726f0081860152606490fd5b60649291601b7f554249506f6f6c3a206e6f7420656e6f7567682062616c616e636500000000009260208c51957f08c379a0000000000000000000000000000000000000000000000000000000008752860152840152820152fd5b6064929160187f554249506f6f6c3a206c656e677468206d69736d6174636800000000000000009260208c51957f08c379a0000000000000000000000000000000000000000000000000000000008752860152840152820152fd5b60849291602e7f554249506f6f6c3a2076616c696461746f72556e636d705075624b65797320639260208c51957f08c379a00000000000000000000000000000000000000000000000000000000087528601528401528201527f616e6e6f7420626520656d7074790000000000000000000000000000000000006064820152fd5b9181601f840112156110ff5782359167ffffffffffffffff83116110ff576020808501948460051b0101116110ff57565b600080fd5b60005b8381106111175750506000910152565b8181015183820152602001611107565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761116857604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff811161116857601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b9291926111dd82611197565b916111eb6040519384611127565b8294818452818301116110ff578281602093846000960137010152565b91908110156112185760051b0190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b91908110156112185760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156110ff57019081359167ffffffffffffffff83116110ff5760200182360381136110ff579190565b6020919283604051948593843782019081520301902090565b3d156112eb573d906112d182611197565b916112df6040519384611127565b82523d6000602084013e565b606090565b67ffffffffffffffff81116111685760051b60200190565b80518210156112185760209160051b010190565b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416330361135c57565b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b906041810361146a5715611218577fff000000000000000000000000000000000000000000000000000000000000007f0400000000000000000000000000000000000000000000000000000000000000913516036113e657565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f5075624b657956657269666965723a20496e76616c6964207075626b6579207060448201527f72656669780000000000000000000000000000000000000000000000000000006064820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f5075624b657956657269666965723a20496e76616c6964207075626b6579206c60448201527f656e6774680000000000000000000000000000000000000000000000000000006064820152fd5b7fffffffffffffffffffffffff0000000000000000000000000000000000000000907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008281541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549073ffffffffffffffffffffffffffffffffffffffff80931680948316179055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3565b906115e157508051156115b757805190602001fd5b60046040517f1425ea42000000000000000000000000000000000000000000000000000000008152fd5b81511580611639575b6115f2575090565b60249073ffffffffffffffffffffffffffffffffffffffff604051917f9996b315000000000000000000000000000000000000000000000000000000008352166004820152fd5b50803b156115ea565b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c161561167157565b60046040517fd7e6bcf8000000000000000000000000000000000000000000000000000000008152fdfea26469706673582212207d7436b9e8aacd238b08265fc62bf6def4c93be8e5cb90404948960faf838f9b64736f6c63430008170033",
}

// UBIPoolABI is the input ABI used to generate the binding from.
// Deprecated: Use UBIPoolMetaData.ABI instead.
var UBIPoolABI = UBIPoolMetaData.ABI

// UBIPoolBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UBIPoolMetaData.Bin instead.
var UBIPoolBin = UBIPoolMetaData.Bin

// DeployUBIPool deploys a new Ethereum contract, binding an instance of UBIPool to it.
func DeployUBIPool(auth *bind.TransactOpts, backend bind.ContractBackend, maxUBIPercentage uint32) (common.Address, *types.Transaction, *UBIPool, error) {
	parsed, err := UBIPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UBIPoolBin), backend, maxUBIPercentage)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UBIPool{UBIPoolCaller: UBIPoolCaller{contract: contract}, UBIPoolTransactor: UBIPoolTransactor{contract: contract}, UBIPoolFilterer: UBIPoolFilterer{contract: contract}}, nil
}

// UBIPool is an auto generated Go binding around an Ethereum contract.
type UBIPool struct {
	UBIPoolCaller     // Read-only binding to the contract
	UBIPoolTransactor // Write-only binding to the contract
	UBIPoolFilterer   // Log filterer for contract events
}

// UBIPoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type UBIPoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UBIPoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UBIPoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UBIPoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UBIPoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UBIPoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UBIPoolSession struct {
	Contract     *UBIPool          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UBIPoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UBIPoolCallerSession struct {
	Contract *UBIPoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// UBIPoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UBIPoolTransactorSession struct {
	Contract     *UBIPoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// UBIPoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type UBIPoolRaw struct {
	Contract *UBIPool // Generic contract binding to access the raw methods on
}

// UBIPoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UBIPoolCallerRaw struct {
	Contract *UBIPoolCaller // Generic read-only contract binding to access the raw methods on
}

// UBIPoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UBIPoolTransactorRaw struct {
	Contract *UBIPoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUBIPool creates a new instance of UBIPool, bound to a specific deployed contract.
func NewUBIPool(address common.Address, backend bind.ContractBackend) (*UBIPool, error) {
	contract, err := bindUBIPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UBIPool{UBIPoolCaller: UBIPoolCaller{contract: contract}, UBIPoolTransactor: UBIPoolTransactor{contract: contract}, UBIPoolFilterer: UBIPoolFilterer{contract: contract}}, nil
}

// NewUBIPoolCaller creates a new read-only instance of UBIPool, bound to a specific deployed contract.
func NewUBIPoolCaller(address common.Address, caller bind.ContractCaller) (*UBIPoolCaller, error) {
	contract, err := bindUBIPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UBIPoolCaller{contract: contract}, nil
}

// NewUBIPoolTransactor creates a new write-only instance of UBIPool, bound to a specific deployed contract.
func NewUBIPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*UBIPoolTransactor, error) {
	contract, err := bindUBIPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UBIPoolTransactor{contract: contract}, nil
}

// NewUBIPoolFilterer creates a new log filterer instance of UBIPool, bound to a specific deployed contract.
func NewUBIPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*UBIPoolFilterer, error) {
	contract, err := bindUBIPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UBIPoolFilterer{contract: contract}, nil
}

// bindUBIPool binds a generic wrapper to an already deployed contract.
func bindUBIPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UBIPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UBIPool *UBIPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UBIPool.Contract.UBIPoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UBIPool *UBIPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UBIPool.Contract.UBIPoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UBIPool *UBIPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UBIPool.Contract.UBIPoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UBIPool *UBIPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UBIPool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UBIPool *UBIPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UBIPool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UBIPool *UBIPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UBIPool.Contract.contract.Transact(opts, method, params...)
}

// MAXUBIPERCENTAGE is a free data retrieval call binding the contract method 0xd5077f40.
//
// Solidity: function MAX_UBI_PERCENTAGE() view returns(uint32)
func (_UBIPool *UBIPoolCaller) MAXUBIPERCENTAGE(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _UBIPool.contract.Call(opts, &out, "MAX_UBI_PERCENTAGE")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// MAXUBIPERCENTAGE is a free data retrieval call binding the contract method 0xd5077f40.
//
// Solidity: function MAX_UBI_PERCENTAGE() view returns(uint32)
func (_UBIPool *UBIPoolSession) MAXUBIPERCENTAGE() (uint32, error) {
	return _UBIPool.Contract.MAXUBIPERCENTAGE(&_UBIPool.CallOpts)
}

// MAXUBIPERCENTAGE is a free data retrieval call binding the contract method 0xd5077f40.
//
// Solidity: function MAX_UBI_PERCENTAGE() view returns(uint32)
func (_UBIPool *UBIPoolCallerSession) MAXUBIPERCENTAGE() (uint32, error) {
	return _UBIPool.Contract.MAXUBIPERCENTAGE(&_UBIPool.CallOpts)
}

// CurrentDistributionId is a free data retrieval call binding the contract method 0x780069e0.
//
// Solidity: function currentDistributionId() view returns(uint256)
func (_UBIPool *UBIPoolCaller) CurrentDistributionId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UBIPool.contract.Call(opts, &out, "currentDistributionId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentDistributionId is a free data retrieval call binding the contract method 0x780069e0.
//
// Solidity: function currentDistributionId() view returns(uint256)
func (_UBIPool *UBIPoolSession) CurrentDistributionId() (*big.Int, error) {
	return _UBIPool.Contract.CurrentDistributionId(&_UBIPool.CallOpts)
}

// CurrentDistributionId is a free data retrieval call binding the contract method 0x780069e0.
//
// Solidity: function currentDistributionId() view returns(uint256)
func (_UBIPool *UBIPoolCallerSession) CurrentDistributionId() (*big.Int, error) {
	return _UBIPool.Contract.CurrentDistributionId(&_UBIPool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UBIPool *UBIPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UBIPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UBIPool *UBIPoolSession) Owner() (common.Address, error) {
	return _UBIPool.Contract.Owner(&_UBIPool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UBIPool *UBIPoolCallerSession) Owner() (common.Address, error) {
	return _UBIPool.Contract.Owner(&_UBIPool.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_UBIPool *UBIPoolCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UBIPool.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_UBIPool *UBIPoolSession) PendingOwner() (common.Address, error) {
	return _UBIPool.Contract.PendingOwner(&_UBIPool.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_UBIPool *UBIPoolCallerSession) PendingOwner() (common.Address, error) {
	return _UBIPool.Contract.PendingOwner(&_UBIPool.CallOpts)
}

// ValidatorUBIAmounts is a free data retrieval call binding the contract method 0xc20c1472.
//
// Solidity: function validatorUBIAmounts(uint256 distributionId, bytes validatorUncmpPubkey) view returns(uint256 amount)
func (_UBIPool *UBIPoolCaller) ValidatorUBIAmounts(opts *bind.CallOpts, distributionId *big.Int, validatorUncmpPubkey []byte) (*big.Int, error) {
	var out []interface{}
	err := _UBIPool.contract.Call(opts, &out, "validatorUBIAmounts", distributionId, validatorUncmpPubkey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorUBIAmounts is a free data retrieval call binding the contract method 0xc20c1472.
//
// Solidity: function validatorUBIAmounts(uint256 distributionId, bytes validatorUncmpPubkey) view returns(uint256 amount)
func (_UBIPool *UBIPoolSession) ValidatorUBIAmounts(distributionId *big.Int, validatorUncmpPubkey []byte) (*big.Int, error) {
	return _UBIPool.Contract.ValidatorUBIAmounts(&_UBIPool.CallOpts, distributionId, validatorUncmpPubkey)
}

// ValidatorUBIAmounts is a free data retrieval call binding the contract method 0xc20c1472.
//
// Solidity: function validatorUBIAmounts(uint256 distributionId, bytes validatorUncmpPubkey) view returns(uint256 amount)
func (_UBIPool *UBIPoolCallerSession) ValidatorUBIAmounts(distributionId *big.Int, validatorUncmpPubkey []byte) (*big.Int, error) {
	return _UBIPool.Contract.ValidatorUBIAmounts(&_UBIPool.CallOpts, distributionId, validatorUncmpPubkey)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_UBIPool *UBIPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UBIPool.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_UBIPool *UBIPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _UBIPool.Contract.AcceptOwnership(&_UBIPool.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_UBIPool *UBIPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _UBIPool.Contract.AcceptOwnership(&_UBIPool.TransactOpts)
}

// ClaimUBI is a paid mutator transaction binding the contract method 0x747c4ef7.
//
// Solidity: function claimUBI(uint256 distributionId, bytes validatorUncmpPubkey) returns()
func (_UBIPool *UBIPoolTransactor) ClaimUBI(opts *bind.TransactOpts, distributionId *big.Int, validatorUncmpPubkey []byte) (*types.Transaction, error) {
	return _UBIPool.contract.Transact(opts, "claimUBI", distributionId, validatorUncmpPubkey)
}

// ClaimUBI is a paid mutator transaction binding the contract method 0x747c4ef7.
//
// Solidity: function claimUBI(uint256 distributionId, bytes validatorUncmpPubkey) returns()
func (_UBIPool *UBIPoolSession) ClaimUBI(distributionId *big.Int, validatorUncmpPubkey []byte) (*types.Transaction, error) {
	return _UBIPool.Contract.ClaimUBI(&_UBIPool.TransactOpts, distributionId, validatorUncmpPubkey)
}

// ClaimUBI is a paid mutator transaction binding the contract method 0x747c4ef7.
//
// Solidity: function claimUBI(uint256 distributionId, bytes validatorUncmpPubkey) returns()
func (_UBIPool *UBIPoolTransactorSession) ClaimUBI(distributionId *big.Int, validatorUncmpPubkey []byte) (*types.Transaction, error) {
	return _UBIPool.Contract.ClaimUBI(&_UBIPool.TransactOpts, distributionId, validatorUncmpPubkey)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner) returns()
func (_UBIPool *UBIPoolTransactor) Initialize(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _UBIPool.contract.Transact(opts, "initialize", owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner) returns()
func (_UBIPool *UBIPoolSession) Initialize(owner common.Address) (*types.Transaction, error) {
	return _UBIPool.Contract.Initialize(&_UBIPool.TransactOpts, owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner) returns()
func (_UBIPool *UBIPoolTransactorSession) Initialize(owner common.Address) (*types.Transaction, error) {
	return _UBIPool.Contract.Initialize(&_UBIPool.TransactOpts, owner)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_UBIPool *UBIPoolTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _UBIPool.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_UBIPool *UBIPoolSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _UBIPool.Contract.Multicall(&_UBIPool.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_UBIPool *UBIPoolTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _UBIPool.Contract.Multicall(&_UBIPool.TransactOpts, data)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UBIPool *UBIPoolTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UBIPool.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UBIPool *UBIPoolSession) RenounceOwnership() (*types.Transaction, error) {
	return _UBIPool.Contract.RenounceOwnership(&_UBIPool.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UBIPool *UBIPoolTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _UBIPool.Contract.RenounceOwnership(&_UBIPool.TransactOpts)
}

// SetUBIDistribution is a paid mutator transaction binding the contract method 0x1971f773.
//
// Solidity: function setUBIDistribution(uint256 totalUBI, bytes[] validatorUncmpPubKeys, uint256[] amounts) returns(uint256)
func (_UBIPool *UBIPoolTransactor) SetUBIDistribution(opts *bind.TransactOpts, totalUBI *big.Int, validatorUncmpPubKeys [][]byte, amounts []*big.Int) (*types.Transaction, error) {
	return _UBIPool.contract.Transact(opts, "setUBIDistribution", totalUBI, validatorUncmpPubKeys, amounts)
}

// SetUBIDistribution is a paid mutator transaction binding the contract method 0x1971f773.
//
// Solidity: function setUBIDistribution(uint256 totalUBI, bytes[] validatorUncmpPubKeys, uint256[] amounts) returns(uint256)
func (_UBIPool *UBIPoolSession) SetUBIDistribution(totalUBI *big.Int, validatorUncmpPubKeys [][]byte, amounts []*big.Int) (*types.Transaction, error) {
	return _UBIPool.Contract.SetUBIDistribution(&_UBIPool.TransactOpts, totalUBI, validatorUncmpPubKeys, amounts)
}

// SetUBIDistribution is a paid mutator transaction binding the contract method 0x1971f773.
//
// Solidity: function setUBIDistribution(uint256 totalUBI, bytes[] validatorUncmpPubKeys, uint256[] amounts) returns(uint256)
func (_UBIPool *UBIPoolTransactorSession) SetUBIDistribution(totalUBI *big.Int, validatorUncmpPubKeys [][]byte, amounts []*big.Int) (*types.Transaction, error) {
	return _UBIPool.Contract.SetUBIDistribution(&_UBIPool.TransactOpts, totalUBI, validatorUncmpPubKeys, amounts)
}

// SetUBIPercentage is a paid mutator transaction binding the contract method 0x47564aa0.
//
// Solidity: function setUBIPercentage(uint32 percentage) returns()
func (_UBIPool *UBIPoolTransactor) SetUBIPercentage(opts *bind.TransactOpts, percentage uint32) (*types.Transaction, error) {
	return _UBIPool.contract.Transact(opts, "setUBIPercentage", percentage)
}

// SetUBIPercentage is a paid mutator transaction binding the contract method 0x47564aa0.
//
// Solidity: function setUBIPercentage(uint32 percentage) returns()
func (_UBIPool *UBIPoolSession) SetUBIPercentage(percentage uint32) (*types.Transaction, error) {
	return _UBIPool.Contract.SetUBIPercentage(&_UBIPool.TransactOpts, percentage)
}

// SetUBIPercentage is a paid mutator transaction binding the contract method 0x47564aa0.
//
// Solidity: function setUBIPercentage(uint32 percentage) returns()
func (_UBIPool *UBIPoolTransactorSession) SetUBIPercentage(percentage uint32) (*types.Transaction, error) {
	return _UBIPool.Contract.SetUBIPercentage(&_UBIPool.TransactOpts, percentage)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UBIPool *UBIPoolTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _UBIPool.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UBIPool *UBIPoolSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UBIPool.Contract.TransferOwnership(&_UBIPool.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UBIPool *UBIPoolTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UBIPool.Contract.TransferOwnership(&_UBIPool.TransactOpts, newOwner)
}

// UBIPoolInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the UBIPool contract.
type UBIPoolInitializedIterator struct {
	Event *UBIPoolInitialized // Event containing the contract specifics and raw log

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
func (it *UBIPoolInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UBIPoolInitialized)
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
		it.Event = new(UBIPoolInitialized)
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
func (it *UBIPoolInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UBIPoolInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UBIPoolInitialized represents a Initialized event raised by the UBIPool contract.
type UBIPoolInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_UBIPool *UBIPoolFilterer) FilterInitialized(opts *bind.FilterOpts) (*UBIPoolInitializedIterator, error) {

	logs, sub, err := _UBIPool.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &UBIPoolInitializedIterator{contract: _UBIPool.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_UBIPool *UBIPoolFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *UBIPoolInitialized) (event.Subscription, error) {

	logs, sub, err := _UBIPool.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UBIPoolInitialized)
				if err := _UBIPool.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_UBIPool *UBIPoolFilterer) ParseInitialized(log types.Log) (*UBIPoolInitialized, error) {
	event := new(UBIPoolInitialized)
	if err := _UBIPool.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UBIPoolOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the UBIPool contract.
type UBIPoolOwnershipTransferStartedIterator struct {
	Event *UBIPoolOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *UBIPoolOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UBIPoolOwnershipTransferStarted)
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
		it.Event = new(UBIPoolOwnershipTransferStarted)
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
func (it *UBIPoolOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UBIPoolOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UBIPoolOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the UBIPool contract.
type UBIPoolOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_UBIPool *UBIPoolFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*UBIPoolOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UBIPool.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &UBIPoolOwnershipTransferStartedIterator{contract: _UBIPool.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_UBIPool *UBIPoolFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *UBIPoolOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UBIPool.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UBIPoolOwnershipTransferStarted)
				if err := _UBIPool.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_UBIPool *UBIPoolFilterer) ParseOwnershipTransferStarted(log types.Log) (*UBIPoolOwnershipTransferStarted, error) {
	event := new(UBIPoolOwnershipTransferStarted)
	if err := _UBIPool.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UBIPoolOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the UBIPool contract.
type UBIPoolOwnershipTransferredIterator struct {
	Event *UBIPoolOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *UBIPoolOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UBIPoolOwnershipTransferred)
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
		it.Event = new(UBIPoolOwnershipTransferred)
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
func (it *UBIPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UBIPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UBIPoolOwnershipTransferred represents a OwnershipTransferred event raised by the UBIPool contract.
type UBIPoolOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UBIPool *UBIPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*UBIPoolOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UBIPool.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &UBIPoolOwnershipTransferredIterator{contract: _UBIPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UBIPool *UBIPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *UBIPoolOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UBIPool.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UBIPoolOwnershipTransferred)
				if err := _UBIPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_UBIPool *UBIPoolFilterer) ParseOwnershipTransferred(log types.Log) (*UBIPoolOwnershipTransferred, error) {
	event := new(UBIPoolOwnershipTransferred)
	if err := _UBIPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UBIPoolUBIDistributionSetIterator is returned from FilterUBIDistributionSet and is used to iterate over the raw logs and unpacked data for UBIDistributionSet events raised by the UBIPool contract.
type UBIPoolUBIDistributionSetIterator struct {
	Event *UBIPoolUBIDistributionSet // Event containing the contract specifics and raw log

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
func (it *UBIPoolUBIDistributionSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UBIPoolUBIDistributionSet)
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
		it.Event = new(UBIPoolUBIDistributionSet)
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
func (it *UBIPoolUBIDistributionSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UBIPoolUBIDistributionSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UBIPoolUBIDistributionSet represents a UBIDistributionSet event raised by the UBIPool contract.
type UBIPoolUBIDistributionSet struct {
	Month                 *big.Int
	TotalUBI              *big.Int
	ValidatorUncmpPubKeys [][]byte
	Amounts               []*big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterUBIDistributionSet is a free log retrieval operation binding the contract event 0x1cc6f356308c8399caa490706b01fb9d52cdc87cdf639e66c3da7d4ce2db161c.
//
// Solidity: event UBIDistributionSet(uint256 month, uint256 totalUBI, bytes[] validatorUncmpPubKeys, uint256[] amounts)
func (_UBIPool *UBIPoolFilterer) FilterUBIDistributionSet(opts *bind.FilterOpts) (*UBIPoolUBIDistributionSetIterator, error) {

	logs, sub, err := _UBIPool.contract.FilterLogs(opts, "UBIDistributionSet")
	if err != nil {
		return nil, err
	}
	return &UBIPoolUBIDistributionSetIterator{contract: _UBIPool.contract, event: "UBIDistributionSet", logs: logs, sub: sub}, nil
}

// WatchUBIDistributionSet is a free log subscription operation binding the contract event 0x1cc6f356308c8399caa490706b01fb9d52cdc87cdf639e66c3da7d4ce2db161c.
//
// Solidity: event UBIDistributionSet(uint256 month, uint256 totalUBI, bytes[] validatorUncmpPubKeys, uint256[] amounts)
func (_UBIPool *UBIPoolFilterer) WatchUBIDistributionSet(opts *bind.WatchOpts, sink chan<- *UBIPoolUBIDistributionSet) (event.Subscription, error) {

	logs, sub, err := _UBIPool.contract.WatchLogs(opts, "UBIDistributionSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UBIPoolUBIDistributionSet)
				if err := _UBIPool.contract.UnpackLog(event, "UBIDistributionSet", log); err != nil {
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

// ParseUBIDistributionSet is a log parse operation binding the contract event 0x1cc6f356308c8399caa490706b01fb9d52cdc87cdf639e66c3da7d4ce2db161c.
//
// Solidity: event UBIDistributionSet(uint256 month, uint256 totalUBI, bytes[] validatorUncmpPubKeys, uint256[] amounts)
func (_UBIPool *UBIPoolFilterer) ParseUBIDistributionSet(log types.Log) (*UBIPoolUBIDistributionSet, error) {
	event := new(UBIPoolUBIDistributionSet)
	if err := _UBIPool.contract.UnpackLog(event, "UBIDistributionSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UBIPoolUBIPercentageSetIterator is returned from FilterUBIPercentageSet and is used to iterate over the raw logs and unpacked data for UBIPercentageSet events raised by the UBIPool contract.
type UBIPoolUBIPercentageSetIterator struct {
	Event *UBIPoolUBIPercentageSet // Event containing the contract specifics and raw log

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
func (it *UBIPoolUBIPercentageSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UBIPoolUBIPercentageSet)
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
		it.Event = new(UBIPoolUBIPercentageSet)
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
func (it *UBIPoolUBIPercentageSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UBIPoolUBIPercentageSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UBIPoolUBIPercentageSet represents a UBIPercentageSet event raised by the UBIPool contract.
type UBIPoolUBIPercentageSet struct {
	Percentage uint32
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUBIPercentageSet is a free log retrieval operation binding the contract event 0x6c6483041303ba314f169eb2d2af177b4f497324ccf0f3c1e68c2100f76c4929.
//
// Solidity: event UBIPercentageSet(uint32 percentage)
func (_UBIPool *UBIPoolFilterer) FilterUBIPercentageSet(opts *bind.FilterOpts) (*UBIPoolUBIPercentageSetIterator, error) {

	logs, sub, err := _UBIPool.contract.FilterLogs(opts, "UBIPercentageSet")
	if err != nil {
		return nil, err
	}
	return &UBIPoolUBIPercentageSetIterator{contract: _UBIPool.contract, event: "UBIPercentageSet", logs: logs, sub: sub}, nil
}

// WatchUBIPercentageSet is a free log subscription operation binding the contract event 0x6c6483041303ba314f169eb2d2af177b4f497324ccf0f3c1e68c2100f76c4929.
//
// Solidity: event UBIPercentageSet(uint32 percentage)
func (_UBIPool *UBIPoolFilterer) WatchUBIPercentageSet(opts *bind.WatchOpts, sink chan<- *UBIPoolUBIPercentageSet) (event.Subscription, error) {

	logs, sub, err := _UBIPool.contract.WatchLogs(opts, "UBIPercentageSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UBIPoolUBIPercentageSet)
				if err := _UBIPool.contract.UnpackLog(event, "UBIPercentageSet", log); err != nil {
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

// ParseUBIPercentageSet is a log parse operation binding the contract event 0x6c6483041303ba314f169eb2d2af177b4f497324ccf0f3c1e68c2100f76c4929.
//
// Solidity: event UBIPercentageSet(uint32 percentage)
func (_UBIPool *UBIPoolFilterer) ParseUBIPercentageSet(log types.Log) (*UBIPoolUBIPercentageSet, error) {
	event := new(UBIPoolUBIPercentageSet)
	if err := _UBIPool.contract.UnpackLog(event, "UBIPercentageSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
