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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"maxUBIPercentage\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"AA\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"BB\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_UBI_PERCENTAGE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimUBI\",\"inputs\":[{\"name\":\"distributionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentDistributionId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUBIDistribution\",\"inputs\":[{\"name\":\"totalUBI\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorCmpPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUBIPercentage\",\"inputs\":[{\"name\":\"percentage\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalPendingClaims\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validatorUBIAmounts\",\"inputs\":[{\"name\":\"distributionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UBIDistributionSet\",\"inputs\":[{\"name\":\"month\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalUBI\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"validatorCmpPubKeys\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UBIPercentageSet\",\"inputs\":[{\"name\":\"percentage\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60a034610103576001600160401b0390601f6200206838819003918201601f19168301918483118484101761010757808492602094604052833981010312610103575163ffffffff81168103610103576080527ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c166100f15780808316036100ac575b604051611f4c90816200011c823960805181818161031f0152610c970152f35b6001600160401b031990911681179091556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f808061008c565b60405163f92ee8a960e01b8152600490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c80631971f7731461012457806347564aa01461011f5780635727dc5c1461011a57806370bf238114610115578063715018a614610110578063747c4ef71461010b578063780069e01461010657806379ba5097146101015780638da5cb5b146100fc578063997da8d4146100f7578063ac9650d8146100f2578063c20c1472146100ed578063c4d66de8146100e8578063d5077f40146100e3578063e30c3978146100de578063eeeac01e146100d95763f2fde38b146100d4575f80fd5b610d47565b610d0d565b610cbb565b610c7b565b610a98565b6109e5565b61086a565b61077b565b610729565b6106a2565b610686565b6104cd565b610408565b6103eb565b6103d0565b6102f2565b61015e565b9181601f8401121561015a5782359167ffffffffffffffff831161015a576020808501948460051b01011161015a57565b5f80fd5b3461015a57606060031936011261015a5767ffffffffffffffff60043560243582811161015a57610193903690600401610129565b91909260443590811161015a576101ae903690600401610129565b906101b7611647565b6101c2841515610e08565b6101cd828514610e93565b6101f66101f1846002546101ec6101e48284610f25565b471015610f37565b610f25565b600255565b5f6102015f54610f9c565b5f555f5b83811061026a5750610266959284926102437f1cc6f356308c8399caa490706b01fb9d52cdc87cdf639e66c3da7d4ce2db161c9661025394146110e9565b5f549660405195869589876111c9565b0390a16040519081529081906020820190565b0390f35b906102eb60019161028761027f858888610ff6565b35151561100b565b6102a461029f610298868b8d611070565b36916109af565b611687565b6102af848787610ff6565b356102d88a6102d2878c6102cc5f545f52600160205260405f2090565b93611070565b906110d0565b556102e4848787610ff6565b3590610f25565b9101610205565b3461015a57602060031936011261015a5760043563ffffffff80821680920361015a5761031d611647565b7f00000000000000000000000000000000000000000000000000000000000000001681116103725760207f6c6483041303ba314f169eb2d2af177b4f497324ccf0f3c1e68c2100f76c492991604051908152a1005b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f554249506f6f6c3a2070657263656e7461676520746f6f2068696768000000006044820152fd5b3461015a575f60031936011261015a57602060405160078152f35b3461015a575f60031936011261015a576020600254604051908152f35b3461015a575f60031936011261015a57610420611647565b5f73ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b3461015a57604060031936011261015a5760243567ffffffffffffffff80821161015a573660238301121561015a57816004013590811161015a576024820191602482369201011161015a577f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00600281541461065c576002905561055a6105553683856109af565b6117e3565b61056381611a01565b73ffffffffffffffffffffffffffffffffffffffff61058061146b565b91604160208401916021810151835201516040840152339251902016036105d8576105ad91600435611322565b6105d660017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b005b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603160248201527f536563703235366b3156657269666965723a20496e76616c6964207075626b6560448201527f79206465726976656420616464726573730000000000000000000000000000006064820152fd5b60046040517f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b3461015a575f60031936011261015a5760205f54604051908152f35b3461015a575f60031936011261015a573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416036106f9576105d633611698565b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b3461015a575f60031936011261015a57602073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416604051908152f35b3461015a575f60031936011261015a5760206040515f8152f35b5f5b8381106107a65750505f910152565b8181015183820152602001610797565b602080820190808352835180925260408301928160408460051b8301019501935f915b8483106107e95750505050505090565b909192939495848080837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc086600196030187527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8c5161085681518092818752878088019101610795565b0116010198019301930191949392906107d9565b3461015a57602060031936011261015a5760043567ffffffffffffffff811161015a5761089b903690600401610129565b906108a4611444565b6108ad836114dc565b925f5b8181106108c5576040518061026687826107b6565b806108e66108e0856108da600195878a611070565b90611542565b30611b80565b6108f0828861157a565b526108fb818761157a565b50016108b0565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761097057604052565b610902565b67ffffffffffffffff811161097057601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b9291926109bb82610975565b916109c9604051938461092f565b82948184528183011161015a578281602093845f960137010152565b3461015a57604060031936011261015a5760243567ffffffffffffffff811161015a573660238201121561015a57610a4f6020610a2f6102669336906024816004013591016109af565b6004355f526001825260405f208260405194838680955193849201610795565b820190815203019020546040519081529081906020820190565b600319602091011261015a5760043573ffffffffffffffffffffffffffffffffffffffff8116810361015a5790565b3461015a57610aa636610a69565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00549067ffffffffffffffff60ff8360401c1615921680159081610c73575b6001149081610c69575b159081610c60575b50610c3657610b599082610b507ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0060017fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000825416179055565b610bda5761158e565b610b5f57005b610bab7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a007fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff8154169055565b604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a1005b610c317ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00680100000000000000007fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff825416179055565b61158e565b60046040517ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b9050155f610af7565b303b159150610aef565b839150610ae5565b3461015a575f60031936011261015a57602060405163ffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b3461015a575f60031936011261015a57602073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416604051908152f35b3461015a575f60031936011261015a5760206040517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8152f35b3461015a57610d5536610a69565b610d5d611647565b73ffffffffffffffffffffffffffffffffffffffff809116907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227005f80a3005b15610e0f57565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f554249506f6f6c3a2076616c696461746f72436d705075624b6579732063616e60448201527f6e6f7420626520656d70747900000000000000000000000000000000000000006064820152fd5b15610e9a57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f554249506f6f6c3a206c656e677468206d69736d6174636800000000000000006044820152fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b91908201809211610f3257565b610ef8565b15610f3e57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f554249506f6f6c3a206e6f7420656e6f7567682062616c616e636500000000006044820152fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610f325760010190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b91908110156110065760051b0190565b610fc9565b1561101257565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f554249506f6f6c3a20616d6f756e74732063616e6e6f74206265207a65726f006044820152fd5b91908110156110065760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18136030182121561015a57019081359167ffffffffffffffff831161015a57602001823603811361015a579190565b6020919283604051948593843782019081520301902090565b156110f057565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f554249506f6f6c3a20746f74616c20616d6f756e74206d69736d6174636800006044820152fd5b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe093818652868601375f8582860101520116010190565b90918281527f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831161015a5760209260051b809284830137010190565b96959380919593956080890190895260209360208a0152608060408a01525260a087019160a08260051b8901019580935f915b84831061121e5750505050505084606061121b9596850391015261118c565b90565b9091929394977fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff608b820301835288357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18336030181121561015a57820185810191903567ffffffffffffffff811161015a57803603831361015a576112a88792839260019561114e565b9a019301930191949392906111fc565b3d156112e2573d906112c982610975565b916112d7604051938461092f565b82523d5f602084013e565b606090565b907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f918203918211610f3257565b91908203918211610f3257565b919091805f52600160205261133b60405f2083856110d0565b549283156113e6575f9261135a611360935f52600160205260405f2090565b916110d0565b555f80808084335af16113716112b8565b5015611388576101f161138691600254611315565b565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f554249506f6f6c3a206661696c656420746f2073656e642055424900000000006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f554249506f6f6c3a206e6f2055424920746f20636c61696d00000000000000006044820152fd5b6040516020810181811067ffffffffffffffff821117610970576040525f8152905f368137565b604051906060820182811067ffffffffffffffff8211176109705760405260408252604082602036910137565b604051906080820182811067ffffffffffffffff82111761097057604052604182526060366020840137565b67ffffffffffffffff81116109705760051b60200190565b906114e6826114c4565b6114f3604051918261092f565b8281527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe061152182946114c4565b01905f5b82811061153157505050565b806060602080938501015201611525565b6020908161138693959461156e87604051988995858701378401918383015f8152815194859201610795565b0103808552018361092f565b80518210156110065760209160051b010190565b73ffffffffffffffffffffffffffffffffffffffff8116156115c357611386906115b6611df0565b6115be611df0565b611698565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f554249506f6f6c3a206f776e65722063616e6e6f74206265207a65726f20616460448201527f64726573730000000000000000000000000000000000000000000000000000006064820152fd5b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300541633036106f957565b611693611386916117e3565b611a01565b7fffffffffffffffffffffffff0000000000000000000000000000000000000000907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008281541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549073ffffffffffffffffffffffffffffffffffffffff80931680948316179055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3565b8051156110065760200190565b1561175f57565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f536563703235366b3156657269666965723a20496e76616c696420636d70207060448201527f75626b65792070726566697800000000000000000000000000000000000000006064820152fd5b60218151036118f2576118996118459161187b7fff000000000000000000000000000000000000000000000000000000000000007f02000000000000000000000000000000000000000000000000000000000000008161186b6118458661174b565b517fff000000000000000000000000000000000000000000000000000000000000001690565b16149081156118bb575b50611758565b61189461188e602183015194859361174b565b60f81c90565b611b9d565b6118a1611498565b9160046118ad8461174b565b536021830152604182015290565b7f030000000000000000000000000000000000000000000000000000000000000091506118ea6118458561174b565b16145f611875565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f536563703235366b3156657269666965723a20496e76616c696420636d70207060448201527f75626b6579206c656e67746800000000000000000000000000000000000000006064820152fd5b1561197d57565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f536563703235366b3156657269666965723a207075626b6579206e6f74206f6e60448201527f20637572766500000000000000000000000000000000000000000000000000006064820152fd5b6041815103611afc577f04000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000611a558361174b565b511603611a7857611a73816041602161138694015191015190611c9e565b611976565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f536563703235366b3156657269666965723a20496e76616c696420756e636d7060448201527f207075626b6579207072656669780000000000000000000000000000000000006064820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f536563703235366b3156657269666965723a20496e76616c696420756e636d7060448201527f207075626b6579206c656e6774680000000000000000000000000000000000006064820152fd5b5f8061121b93602081519101845af4611b976112b8565b91611d50565b60ff1690600282148015611c94575b15611c1057611bfa611bf3611c00927ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f9081806007815f8509089181818009900908611e49565b9283610f25565b60011690565b611c075790565b61121b906112e7565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603160248201527f456c6c697074696343757276653a696e6e76616c696420636f6d70726573736560448201527f6420454320706f696e74207072656669780000000000000000000000000000006064820152fd5b5060038214611bac565b80158015611d26575b8015611d1e575b8015611cf4575b611cee576007907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f918282818181950909089180091490565b50505f90565b507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f821015611cb5565b508115611cae565b507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f811015611ca7565b90611d8f5750805115611d6557805190602001fd5b60046040517f1425ea42000000000000000000000000000000000000000000000000000000008152fd5b81511580611de7575b611da0575090565b60249073ffffffffffffffffffffffffffffffffffffffff604051917f9996b315000000000000000000000000000000000000000000000000000000008352166004820152fd5b50803b15611d98565b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c1615611e1f57565b60046040517fd7e6bcf8000000000000000000000000000000000000000000000000000000008152fd5b8015611f11576001906001917f800000000000000000000000000000000000000000000000000000000000000091825b611e835750505090565b9091927ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f90818080807f3fffffffffffffffffffffffffffffffffffffffffffffffffffffffbfffff0c94818a87161515890a918009098189891c86161515880a91800909818860021c85161515870a91800909918660031c161515840a918009099260041c919082611e79565b505f9056fea2646970667358221220ed71f5269d3862ee0c130cf8ce47629173f59ccd8aa1559fd59f567eeee1a01564736f6c63430008170033",
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
// Solidity: function validatorUBIAmounts(uint256 distributionId, bytes validatorCmpPubkey) view returns(uint256 amount)
func (_UBIPool *UBIPoolCaller) ValidatorUBIAmounts(opts *bind.CallOpts, distributionId *big.Int, validatorCmpPubkey []byte) (*big.Int, error) {
	var out []interface{}
	err := _UBIPool.contract.Call(opts, &out, "validatorUBIAmounts", distributionId, validatorCmpPubkey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorUBIAmounts is a free data retrieval call binding the contract method 0xc20c1472.
//
// Solidity: function validatorUBIAmounts(uint256 distributionId, bytes validatorCmpPubkey) view returns(uint256 amount)
func (_UBIPool *UBIPoolSession) ValidatorUBIAmounts(distributionId *big.Int, validatorCmpPubkey []byte) (*big.Int, error) {
	return _UBIPool.Contract.ValidatorUBIAmounts(&_UBIPool.CallOpts, distributionId, validatorCmpPubkey)
}

// ValidatorUBIAmounts is a free data retrieval call binding the contract method 0xc20c1472.
//
// Solidity: function validatorUBIAmounts(uint256 distributionId, bytes validatorCmpPubkey) view returns(uint256 amount)
func (_UBIPool *UBIPoolCallerSession) ValidatorUBIAmounts(distributionId *big.Int, validatorCmpPubkey []byte) (*big.Int, error) {
	return _UBIPool.Contract.ValidatorUBIAmounts(&_UBIPool.CallOpts, distributionId, validatorCmpPubkey)
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
// Solidity: function claimUBI(uint256 distributionId, bytes validatorCmpPubkey) returns()
func (_UBIPool *UBIPoolTransactor) ClaimUBI(opts *bind.TransactOpts, distributionId *big.Int, validatorCmpPubkey []byte) (*types.Transaction, error) {
	return _UBIPool.contract.Transact(opts, "claimUBI", distributionId, validatorCmpPubkey)
}

// ClaimUBI is a paid mutator transaction binding the contract method 0x747c4ef7.
//
// Solidity: function claimUBI(uint256 distributionId, bytes validatorCmpPubkey) returns()
func (_UBIPool *UBIPoolSession) ClaimUBI(distributionId *big.Int, validatorCmpPubkey []byte) (*types.Transaction, error) {
	return _UBIPool.Contract.ClaimUBI(&_UBIPool.TransactOpts, distributionId, validatorCmpPubkey)
}

// ClaimUBI is a paid mutator transaction binding the contract method 0x747c4ef7.
//
// Solidity: function claimUBI(uint256 distributionId, bytes validatorCmpPubkey) returns()
func (_UBIPool *UBIPoolTransactorSession) ClaimUBI(distributionId *big.Int, validatorCmpPubkey []byte) (*types.Transaction, error) {
	return _UBIPool.Contract.ClaimUBI(&_UBIPool.TransactOpts, distributionId, validatorCmpPubkey)
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
// Solidity: function setUBIDistribution(uint256 totalUBI, bytes[] validatorCmpPubKeys, uint256[] amounts) returns(uint256)
func (_UBIPool *UBIPoolTransactor) SetUBIDistribution(opts *bind.TransactOpts, totalUBI *big.Int, validatorCmpPubKeys [][]byte, amounts []*big.Int) (*types.Transaction, error) {
	return _UBIPool.contract.Transact(opts, "setUBIDistribution", totalUBI, validatorCmpPubKeys, amounts)
}

// SetUBIDistribution is a paid mutator transaction binding the contract method 0x1971f773.
//
// Solidity: function setUBIDistribution(uint256 totalUBI, bytes[] validatorCmpPubKeys, uint256[] amounts) returns(uint256)
func (_UBIPool *UBIPoolSession) SetUBIDistribution(totalUBI *big.Int, validatorCmpPubKeys [][]byte, amounts []*big.Int) (*types.Transaction, error) {
	return _UBIPool.Contract.SetUBIDistribution(&_UBIPool.TransactOpts, totalUBI, validatorCmpPubKeys, amounts)
}

// SetUBIDistribution is a paid mutator transaction binding the contract method 0x1971f773.
//
// Solidity: function setUBIDistribution(uint256 totalUBI, bytes[] validatorCmpPubKeys, uint256[] amounts) returns(uint256)
func (_UBIPool *UBIPoolTransactorSession) SetUBIDistribution(totalUBI *big.Int, validatorCmpPubKeys [][]byte, amounts []*big.Int) (*types.Transaction, error) {
	return _UBIPool.Contract.SetUBIDistribution(&_UBIPool.TransactOpts, totalUBI, validatorCmpPubKeys, amounts)
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
	Month               *big.Int
	TotalUBI            *big.Int
	ValidatorCmpPubKeys [][]byte
	Amounts             []*big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterUBIDistributionSet is a free log retrieval operation binding the contract event 0x1cc6f356308c8399caa490706b01fb9d52cdc87cdf639e66c3da7d4ce2db161c.
//
// Solidity: event UBIDistributionSet(uint256 month, uint256 totalUBI, bytes[] validatorCmpPubKeys, uint256[] amounts)
func (_UBIPool *UBIPoolFilterer) FilterUBIDistributionSet(opts *bind.FilterOpts) (*UBIPoolUBIDistributionSetIterator, error) {

	logs, sub, err := _UBIPool.contract.FilterLogs(opts, "UBIDistributionSet")
	if err != nil {
		return nil, err
	}
	return &UBIPoolUBIDistributionSetIterator{contract: _UBIPool.contract, event: "UBIDistributionSet", logs: logs, sub: sub}, nil
}

// WatchUBIDistributionSet is a free log subscription operation binding the contract event 0x1cc6f356308c8399caa490706b01fb9d52cdc87cdf639e66c3da7d4ce2db161c.
//
// Solidity: event UBIDistributionSet(uint256 month, uint256 totalUBI, bytes[] validatorCmpPubKeys, uint256[] amounts)
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
// Solidity: event UBIDistributionSet(uint256 month, uint256 totalUBI, bytes[] validatorCmpPubKeys, uint256[] amounts)
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
