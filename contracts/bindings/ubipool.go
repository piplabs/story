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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"maxUBIPercentage\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"AA\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"BB\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_UBI_PERCENTAGE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimUBI\",\"inputs\":[{\"name\":\"distributionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentDistributionId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUBIDistribution\",\"inputs\":[{\"name\":\"totalUBI\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorUncmpPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUBIPercentage\",\"inputs\":[{\"name\":\"percentage\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalPendingClaims\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validatorUBIAmounts\",\"inputs\":[{\"name\":\"distributionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UBIDistributionSet\",\"inputs\":[{\"name\":\"month\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalUBI\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"validatorUncmpPubKeys\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UBIPercentageSet\",\"inputs\":[{\"name\":\"percentage\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60a034610101576001600160401b0390601f611a3a38819003918201601f19168301918483118484101761010557808492602094604052833981010312610101575163ffffffff81168103610101576080527ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c166100ef5780808316036100aa575b604051611920908161011a823960805181818161024e0152610c1f0152f35b6001600160401b031990911681179091556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f808061008b565b60405163f92ee8a960e01b8152600490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60406080815260049081361015610014575f80fd5b5f3560e01c9081631971f77314610cce57816347564aa014610bf15781635727dc5c14610bd757816370bf238114610bba578163715018a614610af5578163747c4ef714610821578163780069e01461080557816379ba50971461077c5781638da5cb5b1461072a578163997da8d414610711578163ac9650d814610517578163c20c14721461049f578163c4d66de814610272578163d5077f4014610233578163e30c3978146101e1578163eeeac01e146101a8575063f2fde38b146100d9575f80fd5b346101a45760206003193601126101a4573573ffffffffffffffffffffffffffffffffffffffff8082168092036101a45761011261141e565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227005f80a3005b5f80fd5b346101a4575f6003193601126101a457602090517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8152f35b346101a4575f6003193601126101a45760209073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0054169051908152f35b346101a4575f6003193601126101a4576020905163ffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b82346101a45760206003193601126101a457803573ffffffffffffffffffffffffffffffffffffffff8116908181036101a4577ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009283549260ff84871c16159367ffffffffffffffff811680159081610497575b600114908161048d575b159081610484575b5061045c578460017fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008316178755610427575b50156103a4575061034b9061033e611891565b610346611891565b61168c565b61035157005b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d291817fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff602093541690555160018152a1005b60849060208651917f08c379a0000000000000000000000000000000000000000000000000000000008352820152602560248201527f554249506f6f6c3a206f776e65722063616e6e6f74206265207a65726f20616460448201527f64726573730000000000000000000000000000000000000000000000000000006064820152fd5b7fffffffffffffffffffffffffffffffffffffffffffffff00000000000000000016680100000000000000011785558661032b565b8287517ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b905015886102f8565b303b1591506102f0565b8691506102e6565b9050346101a457806003193601126101a4576024359167ffffffffffffffff83116101a457366023840112156101a4576020610506916104e982953690602481850135910161129d565b90355f5260018252835f20828551948386809551938492016111d4565b820190815203019020549051908152f35b82346101a457602090816003193601126101a45767ffffffffffffffff81358181116101a45761054a90369084016111a3565b9185519385850191858310908311176106e5575080869396525f845261056f866113f2565b9361057c845195866111f5565b868552610588876113f2565b967fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08098015f5b8181106106d65750505f5b81811061065557505050505080519280840190808552835180925280838601938360051b8701019401925f965b8388106105f45786860387f35b9091929394838080837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08b6001960301875285601f8b51610640815180928187528780880191016111d4565b011601019701930197019690939291936105e7565b806106b75f80896106976106a39b9e9b8e8b8b6106768f9b8d60019d61134a565b9290965195838794868601998a37840191858301938a8552519384916111d4565b010380845201826111f5565b5190305af46106b06113c3565b90306117f1565b6106c1828a61140a565b526106cc818961140a565b50019794976105ba565b60608882018a015288016105af565b6041907f4e487b71000000000000000000000000000000000000000000000000000000005f525260245ffd5b346101a4575f6003193601126101a457602090515f8152f35b346101a4575f6003193601126101a45760209073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054169051908152f35b82346101a4575f6003193601126101a4573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416036107d6576107d43361168c565b005b60249151907f118cdaa70000000000000000000000000000000000000000000000000000000082523390820152fd5b346101a4575f6003193601126101a4576020905f549051908152f35b9050346101a457806003193601126101a457813560249283359267ffffffffffffffff908185116101a457366023860112156101a457848301359182116101a457858501368784880101116101a4577f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00956002875414610acd57600287556108a9848361148e565b836001116101a4576108e390369060257fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8701910161129d565b9283519373ffffffffffffffffffffffffffffffffffffffff60209586339301201603610a4c57855f526001845261091e835f2082846113aa565b549586156109f1579161093c915f93845260018652848420916113aa565b555f80808087335af161094d6113c3565b501561099857505060025491820391821161096d57600282905560018355005b601184917f4e487b71000000000000000000000000000000000000000000000000000000005f52525ffd5b517f08c379a000000000000000000000000000000000000000000000000000000000815291820152601b818501527f554249506f6f6c3a206661696c656420746f2073656e642055424900000000006044820152606490fd5b60648660188b888851937f08c379a00000000000000000000000000000000000000000000000000000000085528401528201527f554249506f6f6c3a206e6f2055424920746f20636c61696d00000000000000006044820152fd5b608485602e8a878751937f08c379a00000000000000000000000000000000000000000000000000000000085528401528201527f5075624b657956657269666965723a20496e76616c6964207075626b6579206460448201527f65726976656420616464726573730000000000000000000000000000000000006064820152fd5b8483517f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b346101a4575f6003193601126101a457610b0d61141e565b5f73ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b346101a4575f6003193601126101a4576020906002549051908152f35b346101a4575f6003193601126101a4576020905160078152f35b82346101a45760206003193601126101a45780359063ffffffff8083168093036101a457610c1d61141e565b7f0000000000000000000000000000000000000000000000000000000000000000168211610c71577f6c6483041303ba314f169eb2d2af177b4f497324ccf0f3c1e68c2100f76c49296020838551908152a1005b60649060208451917f08c379a0000000000000000000000000000000000000000000000000000000008352820152601c60248201527f554249506f6f6c3a2070657263656e7461676520746f6f2068696768000000006044820152fd5b82346101a45760606003193601126101a457803591602480359267ffffffffffffffff938481116101a457610d0690369083016111a3565b94909160449384358381116101a457610d2290369084016111a3565b929095610d2d61141e565b8815611123578389036110c957600254610d47818c6112d3565b471061106e578a610d57916112d3565b6002555f925f54937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8514611043578b9060018096015f555f5b878110610f6e575003610f13575050509592915f549686519480608087018a885260209b8c89015260808a8901525260a086019460a08260051b8801019581945f925b848410610e4e575050505050505082820360608401528082527f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81116101a4577f1cc6f356308c8399caa490706b01fb9d52cdc87cdf639e66c3da7d4ce2db161c938792849260051b80928583013701030190a151908152f35b909192939495977fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff608a82030184528835907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1833603018212156101a457908201808f019190358481116101a45780360383136101a4578f827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84808e98869897879852868601375f858286010152011601019a0194019401929594939190610dd4565b60649291601e7f554249506f6f6c3a20746f74616c20616d6f756e74206d69736d6174636800009260208c51957f08c379a0000000000000000000000000000000000000000000000000000000008752860152840152820152fd5b9150610f7b82888c61130d565b3515610fe857610fde8691610fd7848f8f8f908f808f948681610fae610fa883610fcb95610fd19861134a565b9061148e565b610fb982898961130d565b35955f545f528d6020525f209361134a565b906113aa565b5561130d565b35906112d3565b9101908c91610d91565b8a517f08c379a0000000000000000000000000000000000000000000000000000000008152602081870152601f818501527f554249506f6f6c3a20616d6f756e74732063616e6e6f74206265207a65726f0081860152606490fd5b506011837f4e487b71000000000000000000000000000000000000000000000000000000005f52525ffd5b507f554249506f6f6c3a206e6f7420656e6f7567682062616c616e6365000000000090601b60649460208b51957f08c379a0000000000000000000000000000000000000000000000000000000008752860152840152820152fd5b7f554249506f6f6c3a206c656e677468206d69736d61746368000000000000000090601860649460208b51957f08c379a0000000000000000000000000000000000000000000000000000000008752860152840152820152fd5b7f554249506f6f6c3a2076616c696461746f72556e636d705075624b657973206390602e60849460208b51957f08c379a00000000000000000000000000000000000000000000000000000000087528601528401528201527f616e6e6f7420626520656d7074790000000000000000000000000000000000006064820152fd5b9181601f840112156101a45782359167ffffffffffffffff83116101a4576020808501948460051b0101116101a457565b5f5b8381106111e55750505f910152565b81810151838201526020016111d6565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761123657604052565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b67ffffffffffffffff811161123657601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b9291926112a982611263565b916112b760405193846111f5565b8294818452818301116101a4578281602093845f960137010152565b919082018092116112e057565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b919081101561131d5760051b0190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b919081101561131d5760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156101a457019081359167ffffffffffffffff83116101a45760200182360381136101a4579190565b6020919283604051948593843782019081520301902090565b3d156113ed573d906113d482611263565b916113e260405193846111f5565b82523d5f602084013e565b606090565b67ffffffffffffffff81116112365760051b60200190565b805182101561131d5760209160051b010190565b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416330361145e57565b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b9060418103611608571561131d577f04000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000823516036115845780600160216114f993013591013561173f565b1561150057565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602760248201527f5075624b657956657269666965723a20496e76616c6964207075626b6579206f60448201527f6e206375727665000000000000000000000000000000000000000000000000006064820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f5075624b657956657269666965723a20496e76616c6964207075626b6579207060448201527f72656669780000000000000000000000000000000000000000000000000000006064820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f5075624b657956657269666965723a20496e76616c6964207075626b6579206c60448201527f656e6774680000000000000000000000000000000000000000000000000000006064820152fd5b7fffffffffffffffffffffffff0000000000000000000000000000000000000000907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008281541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549073ffffffffffffffffffffffffffffffffffffffff80931680948316179055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3565b801580156117c7575b80156117bf575b8015611795575b61178f576007907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f918282818181950909089180091490565b50505f90565b507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f821015611756565b50811561174f565b507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f811015611748565b90611830575080511561180657805190602001fd5b60046040517f1425ea42000000000000000000000000000000000000000000000000000000008152fd5b81511580611888575b611841575090565b60249073ffffffffffffffffffffffffffffffffffffffff604051917f9996b315000000000000000000000000000000000000000000000000000000008352166004820152fd5b50803b15611839565b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c16156118c057565b60046040517fd7e6bcf8000000000000000000000000000000000000000000000000000000008152fdfea26469706673582212200361a4c84bf04a9cd9e78ad3dbadd7dbc05757c1987d87a95f29b643bf4f5c3e64736f6c63430008170033",
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

// AA is a free data retrieval call binding the contract method 0x997da8d4.
//
// Solidity: function AA() view returns(uint256)
func (_UBIPool *UBIPoolCaller) AA(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UBIPool.contract.Call(opts, &out, "AA")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AA is a free data retrieval call binding the contract method 0x997da8d4.
//
// Solidity: function AA() view returns(uint256)
func (_UBIPool *UBIPoolSession) AA() (*big.Int, error) {
	return _UBIPool.Contract.AA(&_UBIPool.CallOpts)
}

// AA is a free data retrieval call binding the contract method 0x997da8d4.
//
// Solidity: function AA() view returns(uint256)
func (_UBIPool *UBIPoolCallerSession) AA() (*big.Int, error) {
	return _UBIPool.Contract.AA(&_UBIPool.CallOpts)
}

// BB is a free data retrieval call binding the contract method 0x5727dc5c.
//
// Solidity: function BB() view returns(uint256)
func (_UBIPool *UBIPoolCaller) BB(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UBIPool.contract.Call(opts, &out, "BB")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BB is a free data retrieval call binding the contract method 0x5727dc5c.
//
// Solidity: function BB() view returns(uint256)
func (_UBIPool *UBIPoolSession) BB() (*big.Int, error) {
	return _UBIPool.Contract.BB(&_UBIPool.CallOpts)
}

// BB is a free data retrieval call binding the contract method 0x5727dc5c.
//
// Solidity: function BB() view returns(uint256)
func (_UBIPool *UBIPoolCallerSession) BB() (*big.Int, error) {
	return _UBIPool.Contract.BB(&_UBIPool.CallOpts)
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

// PP is a free data retrieval call binding the contract method 0xeeeac01e.
//
// Solidity: function PP() view returns(uint256)
func (_UBIPool *UBIPoolCaller) PP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UBIPool.contract.Call(opts, &out, "PP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PP is a free data retrieval call binding the contract method 0xeeeac01e.
//
// Solidity: function PP() view returns(uint256)
func (_UBIPool *UBIPoolSession) PP() (*big.Int, error) {
	return _UBIPool.Contract.PP(&_UBIPool.CallOpts)
}

// PP is a free data retrieval call binding the contract method 0xeeeac01e.
//
// Solidity: function PP() view returns(uint256)
func (_UBIPool *UBIPoolCallerSession) PP() (*big.Int, error) {
	return _UBIPool.Contract.PP(&_UBIPool.CallOpts)
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

// TotalPendingClaims is a free data retrieval call binding the contract method 0x70bf2381.
//
// Solidity: function totalPendingClaims() view returns(uint256)
func (_UBIPool *UBIPoolCaller) TotalPendingClaims(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UBIPool.contract.Call(opts, &out, "totalPendingClaims")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalPendingClaims is a free data retrieval call binding the contract method 0x70bf2381.
//
// Solidity: function totalPendingClaims() view returns(uint256)
func (_UBIPool *UBIPoolSession) TotalPendingClaims() (*big.Int, error) {
	return _UBIPool.Contract.TotalPendingClaims(&_UBIPool.CallOpts)
}

// TotalPendingClaims is a free data retrieval call binding the contract method 0x70bf2381.
//
// Solidity: function totalPendingClaims() view returns(uint256)
func (_UBIPool *UBIPoolCallerSession) TotalPendingClaims() (*big.Int, error) {
	return _UBIPool.Contract.TotalPendingClaims(&_UBIPool.CallOpts)
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
