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

// IDKGEnclaveInstanceData is an auto generated low-level Go binding around an user-defined struct.
type IDKGEnclaveInstanceData struct {
	Round          uint32
	ValidatorAddr  common.Address
	EnclaveType    [32]byte
	EnclaveCommKey []byte
	DkgPubKey      []byte
}

// IDKGEnclaveTypeData is an auto generated low-level Go binding around an user-defined struct.
type IDKGEnclaveTypeData struct {
	CodeCommitment     [32]byte
	ValidationHookAddr common.Address
}

// DKGMetaData contains all meta data concerning the DKG contract.
var DKGMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authenticateEnclaveReport\",\"inputs\":[{\"name\":\"enclaveReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"enclaveInstanceData\",\"type\":\"tuple\",\"internalType\":\"structIDKG.EnclaveInstanceData\",\"components\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validatorAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"enclaveType\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"enclaveCommKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"validationContext\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"enclaveTypeData\",\"inputs\":[{\"name\":\"enclaveType\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDKG.EnclaveTypeData\",\"components\":[{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validationHookAddr\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalize\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validatorAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"enclaveType\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"participantsRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"publicCoeffs\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"pubKeyShare\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minReqRegisteredParticipants\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minReqFinalizedParticipants\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"operationalThreshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isEnclaveTypeWhitelisted\",\"inputs\":[{\"name\":\"enclaveType\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minReqFinalizedParticipants\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minReqRegisteredParticipants\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operationalThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operationalThresholdBasis\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"register\",\"inputs\":[{\"name\":\"enclaveReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"enclaveInstanceData\",\"type\":\"tuple\",\"internalType\":\"structIDKG.EnclaveInstanceData\",\"components\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validatorAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"enclaveType\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"enclaveCommKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"startBlockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validationContext\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFee\",\"inputs\":[{\"name\":\"newFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinReqFinalizedParticipants\",\"inputs\":[{\"name\":\"newMinReqFinalizedParticipants\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinReqRegisteredParticipants\",\"inputs\":[{\"name\":\"newMinReqRegisteredParticipants\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOperationalThreshold\",\"inputs\":[{\"name\":\"newOperationalThreshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"whitelistEnclaveType\",\"inputs\":[{\"name\":\"enclaveType\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"enclaveTypeData\",\"type\":\"tuple\",\"internalType\":\"structIDKG.EnclaveTypeData\",\"components\":[{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validationHookAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"isWhitelisted\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"EnclaveTypeWhitelisted\",\"inputs\":[{\"name\":\"enclaveType\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"validationHookAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"isWhitelisted\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeSet\",\"inputs\":[{\"name\":\"newFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Finalized\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"validatorAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"enclaveType\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"participantsRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"publicCoeffs\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"},{\"name\":\"pubKeyShare\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinReqFinalizedParticipantsSet\",\"inputs\":[{\"name\":\"newMinReqFinalizedParticipants\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinReqRegisteredParticipantsSet\",\"inputs\":[{\"name\":\"newMinReqRegisteredParticipants\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperationalThresholdSet\",\"inputs\":[{\"name\":\"newOperationalThreshold\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Registered\",\"inputs\":[{\"name\":\"enclaveReport\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"validatorAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"enclaveType\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"enclaveCommKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"startBlockHeight\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"startBlockHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"validationContext\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60a080604052346100cc57306080527ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c166100bd57506001600160401b036002600160401b031982821601610078575b6040516126a390816100d1823960805181818161152501526116280152f35b6001600160401b031990911681179091556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f8080610059565b63f92ee8a960e01b8152600490fd5b5f80fdfe6080806040526004361015610012575f80fd5b5f3560e01c9081631b0d13ff146118eb575080632adbc9b6146118af578063427c294f1461188b5780634f1ef2861461159f57806352d1902d146114fe5780635c975abb146114bd5780635fbbc42d1461149957806369fe0e2d14611475578063715018a6146113b057806379ba509714611329578063804faf5b14610f595780638da5cb5b14610f075780639f5fc11414610e45578063a80e9f4214610951578063a84565e714610935578063ad3cb1cc14610881578063b72d322714610845578063d3d3409b146107f7578063d78a5416146107d1578063ddca3f4314610795578063e30c397814610743578063ed3e270d14610707578063f2fde38b1461063b578063f92ad219146103d75763fb3dde661461012f575f80fd5b346103d35760806003193601126103d35760043560407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdc3601126103d3576040519061017a826119a2565b602435825273ffffffffffffffffffffffffffffffffffffffff9160443583811681036103d357602093848301918252606435928315158094036103d3576101c0611e0f565b6101cb851515611c1b565b8051156103505781835116156102cd57917f0be17b82f2eb1d1de6349e90606850d99d513c3a3c418ea90b7f0cd0a543863095918594936080965f527f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f048452600160405f208251815501828451167fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f05845260405f207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0081541660ff8716179055519151169160405194855284015260408301526060820152a1005b608486604051907f08c379a000000000000000000000000000000000000000000000000000000000825260048201526024808201527f444b473a2056616c69646174696f6e20686f6f6b2063616e6e6f74206265206560448201527f6d707479000000000000000000000000000000000000000000000000000000006064820152fd5b608486604051907f08c379a000000000000000000000000000000000000000000000000000000000825260048201526024808201527f444b473a20436f646520636f6d6d69746d656e742063616e6e6f74206265206560448201527f6d707479000000000000000000000000000000000000000000000000000000006064820152fd5b5f80fd5b346103d35760a06003193601126103d3576103f061197f565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c16159167ffffffffffffffff811680159081610633575b6001149081610629575b159081610620575b506105f6578260017fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000083161785556105c1575b50610480612574565b610488612574565b73ffffffffffffffffffffffffffffffffffffffff811615610591576104ad906120e0565b6104b5612574565b6104bd612574565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00815416905561050c612574565b61051760243561249a565b610522604435611e4f565b61052d606435611f29565b610538608435612090565b61053e57005b7fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff81541690557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160018152a1005b60246040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081525f6004820152fd5b7fffffffffffffffffffffffffffffffffffffffffffffff000000000000000000166801000000000000000117835583610477565b60046040517ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b90501585610444565b303b15915061043c565b849150610432565b346103d35760206003193601126103d35761065461197f565b61065c611e0f565b73ffffffffffffffffffffffffffffffffffffffff809116907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227005f80a3005b346103d3575f6003193601126103d35760207f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f0054604051908152f35b346103d3575f6003193601126103d357602073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416604051908152f35b346103d3575f6003193601126103d35760207f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f0354604051908152f35b346103d35760206003193601126103d3576107ea611e0f565b6107f560043561249a565b005b346103d35760206003193601126103d3576004355f527f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f05602052602060ff60405f2054166040519015158152f35b346103d3575f6003193601126103d35760207f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f0154604051908152f35b346103d3575f6003193601126103d35760405161089d816119a2565b600581526020907f352e302e3000000000000000000000000000000000000000000000000000000060208201526040518092602082528251928360208401525f5b84811061091e575050507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f835f604080968601015201168101030190f35b8181018301518682016040015285935082016108de565b346103d3575f6003193601126103d35760206040516103e88152f35b6101006003193601126103d35760043563ffffffff811681036103d3576024359073ffffffffffffffffffffffffffffffffffffffff821682036103d35767ffffffffffffffff916084358381116103d3576109b1903690600401611a66565b909360a435908082116103d357366023830112156103d357808260040135116103d357366024836004013560051b840101116103d35760c4358181116103d3576109ff903690600401611a66565b9160e4359081116103d357610a18903690600401611a66565b949097610a477f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f03543414611a94565b5f3415610e3c575b5f80808093813491f115610e3157610a65612193565b610a7663ffffffff89161515611b0a565b610a9773ffffffffffffffffffffffffffffffffffffffff88161515611b90565b6044355f527f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f05602052610ad060ff60405f205416611d35565b60643515610dad578115610d2957846004013515610ca5578515610c4757610b54916044355f527f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f0460205260405f205463ffffffff6040519a168a5260443560208b015260408a015260643560608a015261010060808a0152610100890191611cf7565b86810360a08801528360040135815260208101906020856004013560051b8201019460248101925f915b80600401358310610bf057508a7feeeaac510415cc754b05e40e2e701b7d22d7d764d2a0b232c26d9eea3b62565d8b80610beb8d8d73ffffffffffffffffffffffffffffffffffffffff610bdb8f8f8f88830360c08a0152611cf7565b9285840360e08701521696611cf7565b0390a2005b9091929396602080610c38837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08960019603018752610c328c60248801611dbf565b90611cf7565b99019301930191939290610b7e565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f444b473a205369676e61747572652063616e6e6f7420626520656d70747900006044820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f444b473a205075626c696320636f656666696369656e74732063616e6e6f742060448201527f626520656d7074790000000000000000000000000000000000000000000000006064820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f444b473a20476c6f62616c207075626c6963206b65792063616e6e6f7420626560448201527f20656d70747900000000000000000000000000000000000000000000000000006064820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f444b473a205061727469636970616e747320726f6f742063616e6e6f7420626560448201527f20656d70747900000000000000000000000000000000000000000000000000006064820152fd5b6040513d5f823e3d90fd5b506108fc610a4f565b6003196060813601126103d35767ffffffffffffffff6004358181116103d357610e73903690600401611a66565b9091602435938185116103d35760a09085360301126103d3576044359081116103d357610ea4903690600401611a66565b929091610ed37f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f03543414611a94565b5f3415610efe575b5f80808093813491f115610e31576107f594610ef5612193565b600401916121e8565b506108fc610edb565b346103d3575f6003193601126103d357602073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416604051908152f35b60031960a0813601126103d35767ffffffffffffffff906004358281116103d357610f88903690600401611a66565b9091602435938085116103d35760a0856004019286360301126103d3576084359081116103d357610fbd903690600401611a66565b949091610fec7f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f03543414611a94565b5f3415611320575b5f80808093813491f115610e315761100a612193565b831561129c5763ffffffff9161102b8361102384611af9565b161515611b0a565b73ffffffffffffffffffffffffffffffffffffffff916024820161105a8461105283611b6f565b161515611b90565b60448301359161106b831515611c1b565b606484019361107a8583611ca6565b905015611218576084019161108f8383611ca6565b905015611194576110a38b89848c8e6121e8565b6110ac82611af9565b906110b690611b6f565b946110c19083611ca6565b936110cc9193611ca6565b939094805f527f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f0460205260405f2054986040519c8d9c8d61012090818152019061111592611cf7565b931660208c015260408b015289820360608b015261113292611cf7565b90878203608089015261114492611cf7565b9260a086015260443560c086015260643560e0860152848303610100860152169461116e92611cf7565b037f66979019c139c21e4120d46d1cb1ae87145337ca67db0c998266f9f58ce8f2be91a2005b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f444b473a20444b47207075626c6963206b65792063616e6e6f7420626520656d60448201527f70747900000000000000000000000000000000000000000000000000000000006064820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f444b473a20456e636c61766520636f6d6d756e69636174696f6e206b6579206360448201527f616e6e6f7420626520656d7074790000000000000000000000000000000000006064820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f444b473a20456e636c617665207265706f72742063616e6e6f7420626520656d60448201527f70747900000000000000000000000000000000000000000000000000000000006064820152fd5b506108fc610ff4565b346103d3575f6003193601126103d3573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00541603611380576107f5336120e0565b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b346103d3575f6003193601126103d3576113c8611e0f565b5f73ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b346103d35760206003193601126103d35761148e611e0f565b6107f5600435612090565b346103d35760206003193601126103d3576114b2611e0f565b6107f5600435611f29565b346103d3575f6003193601126103d357602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b346103d3575f6003193601126103d35773ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001630036115755760206040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b60046040517fe07c8dba000000000000000000000000000000000000000000000000000000008152fd5b60406003193601126103d3576115b361197f565b602490813567ffffffffffffffff81116103d357366023820112156103d35780600401356115e081611a2c565b926115ee60405194856119eb565b81845260209182850193368783830101116103d357815f928886930187378601015273ffffffffffffffffffffffffffffffffffffffff807f00000000000000000000000000000000000000000000000000000000000000001680301490811561185d575b5061157557611660611e0f565b8116936040517f52d1902d0000000000000000000000000000000000000000000000000000000081528381600481895afa5f918161182e575b506116ce578686604051907f4c9c8ce30000000000000000000000000000000000000000000000000000000082526004820152fd5b8590877f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc918281036118005750843b156117d15750817fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055604051907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a281511561179d57505f6107f59481925190845af4903d15611794573d61177881611a2c565b9061178660405192836119eb565b81525f81943d92013e6125cd565b606092506125cd565b9350505050346117a957005b807fb398979f0000000000000000000000000000000000000000000000000000000060049252fd5b82604051907f4c9c8ce30000000000000000000000000000000000000000000000000000000082526004820152fd5b604051907faa1d49a40000000000000000000000000000000000000000000000000000000082526004820152fd5b9091508481813d8311611856575b61184681836119eb565b810103126103d357519088611699565b503d61183c565b9050817f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5416141587611653565b346103d35760206003193601126103d3576118a4611e0f565b6107f5600435611e4f565b346103d3575f6003193601126103d35760207f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f0254604051908152f35b346103d35760206003193601126103d3576020816119095f936119a2565b82815201526004355f527f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f046020526040805f2060018251611949816119a2565b602083549384835273ffffffffffffffffffffffffffffffffffffffff9384910154169101908152835192835251166020820152f35b6004359073ffffffffffffffffffffffffffffffffffffffff821682036103d357565b6040810190811067ffffffffffffffff8211176119be57604052565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176119be57604052565b67ffffffffffffffff81116119be57601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b9181601f840112156103d35782359167ffffffffffffffff83116103d357602083818601950101116103d357565b15611a9b57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f444b473a20496e76616c69642066656520616d6f756e740000000000000000006044820152fd5b3563ffffffff811681036103d35790565b15611b1157565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f444b473a20526f756e642063616e6e6f74206265207a65726f000000000000006044820152fd5b3573ffffffffffffffffffffffffffffffffffffffff811681036103d35790565b15611b9757565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f444b473a2056616c696461746f7220616464726573732063616e6e6f7420626560448201527f20656d70747900000000000000000000000000000000000000000000000000006064820152fd5b15611c2257565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f444b473a20456e636c61766520747970652063616e6e6f7420626520656d707460448201527f79000000000000000000000000000000000000000000000000000000000000006064820152fd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156103d3570180359067ffffffffffffffff82116103d3576020019181360383136103d357565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe093818652868601375f8582860101520116010190565b15611d3c57565b60846040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f444b473a20456e636c6176652074797065206973206e6f742077686974656c6960448201527f73746564000000000000000000000000000000000000000000000000000000006064820152fd5b90357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1823603018112156103d357016020813591019167ffffffffffffffff82116103d35781360383136103d357565b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416330361138057565b8015611ea5576020817fced3da4b96d055042d4df66edffabf4cdaa5909bab55d737f34d71e50a07ae58927f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f0155604051908152a1565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602f60248201527f444b473a204d696e52657146696e616c697a65645061727469636970616e747360448201527f2063616e6e6f74206265207a65726f00000000000000000000000000000000006064820152fd5b801561200c576103e88111611f88576020817f90701eb5a3be7e3396be0be01f464764263a71e87a65140746a0e909384e990b927f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f0255604051908152a1565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f444b473a204f7065726174696f6e616c207468726573686f6c642063616e6e6f60448201527f742062652067726561746572207468616e2031303030000000000000000000006064820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602960248201527f444b473a204f7065726174696f6e616c207468726573686f6c642063616e6e6f60448201527f74206265207a65726f00000000000000000000000000000000000000000000006064820152fd5b60207f20461e09b8e557b77e107939f9ce6544698123aad0fc964ac5cc59b7df2e608f91807f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f0355604051908152a1565b7fffffffffffffffffffffffff0000000000000000000000000000000000000000907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008281541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549073ffffffffffffffffffffffffffffffffffffffff80931680948316179055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3565b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166121be57565b60046040517fd93c0665000000000000000000000000000000000000000000000000000000008152fd5b91906040938483013595865f526020967f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f048852865f2092600188519461222d866119a2565b8054865273ffffffffffffffffffffffffffffffffffffffff809281920154168b8701908152845f527f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f058c5261228860ff8c5f205416611d35565b511694519189518b8101908c8252893563ffffffff81168091036103d3578c8201528c8a01359384168094036103d35761235b8161232f8f9c5f996123bb98606085015260808401526122ff6122f46122e46060840184611dbf565b60a08088015260e0870191611cf7565b916080810190611dbf565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08584030160c0860152611cf7565b037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081018352826119eb565b519020946123a98b519a8b998a9889977f793e6ce500000000000000000000000000000000000000000000000000000000895260048901526024880152608060448801526084870191611cf7565b91600319858403016064860152611cf7565b03925af1908115612490575f9161245a575b50156123d7575050565b6084925051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152602260248201527f444b473a20456e636c6176652061757468656e7469636174696f6e206661696c60448201527f65640000000000000000000000000000000000000000000000000000000000006064820152fd5b90508281813d8311612489575b61247181836119eb565b810103126103d3575180151581036103d3575f6123cd565b503d612467565b82513d5f823e3d90fd5b80156124f0576020817fcc68805432194238c4be6e6bec7f8c92631b27c0abd3e73f7f6467f41f9dfebd927f12adbc8310743862abf90938062c5fa55b79e1f01d48ab009dbc12b76d617f0055604051908152a1565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603060248201527f444b473a204d696e526571526567697374657265645061727469636970616e7460448201527f732063616e6e6f74206265207a65726f000000000000000000000000000000006064820152fd5b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c16156125a357565b60046040517fd7e6bcf8000000000000000000000000000000000000000000000000000000008152fd5b9061260c57508051156125e257805190602001fd5b60046040517f1425ea42000000000000000000000000000000000000000000000000000000008152fd5b81511580612664575b61261d575090565b60249073ffffffffffffffffffffffffffffffffffffffff604051917f9996b315000000000000000000000000000000000000000000000000000000008352166004820152fd5b50803b1561261556fea264697066735822122008f0954d70a869b1b80971a88478a544a197ab56486635e5100738de47d952fd64736f6c63430008170033",
}

// DKGABI is the input ABI used to generate the binding from.
// Deprecated: Use DKGMetaData.ABI instead.
var DKGABI = DKGMetaData.ABI

// DKGBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DKGMetaData.Bin instead.
var DKGBin = DKGMetaData.Bin

// DeployDKG deploys a new Ethereum contract, binding an instance of DKG to it.
func DeployDKG(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DKG, error) {
	parsed, err := DKGMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DKGBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DKG{DKGCaller: DKGCaller{contract: contract}, DKGTransactor: DKGTransactor{contract: contract}, DKGFilterer: DKGFilterer{contract: contract}}, nil
}

// DKG is an auto generated Go binding around an Ethereum contract.
type DKG struct {
	DKGCaller     // Read-only binding to the contract
	DKGTransactor // Write-only binding to the contract
	DKGFilterer   // Log filterer for contract events
}

// DKGCaller is an auto generated read-only Go binding around an Ethereum contract.
type DKGCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DKGTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DKGTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DKGFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DKGFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DKGSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DKGSession struct {
	Contract     *DKG              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DKGCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DKGCallerSession struct {
	Contract *DKGCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DKGTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DKGTransactorSession struct {
	Contract     *DKGTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DKGRaw is an auto generated low-level Go binding around an Ethereum contract.
type DKGRaw struct {
	Contract *DKG // Generic contract binding to access the raw methods on
}

// DKGCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DKGCallerRaw struct {
	Contract *DKGCaller // Generic read-only contract binding to access the raw methods on
}

// DKGTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DKGTransactorRaw struct {
	Contract *DKGTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDKG creates a new instance of DKG, bound to a specific deployed contract.
func NewDKG(address common.Address, backend bind.ContractBackend) (*DKG, error) {
	contract, err := bindDKG(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DKG{DKGCaller: DKGCaller{contract: contract}, DKGTransactor: DKGTransactor{contract: contract}, DKGFilterer: DKGFilterer{contract: contract}}, nil
}

// NewDKGCaller creates a new read-only instance of DKG, bound to a specific deployed contract.
func NewDKGCaller(address common.Address, caller bind.ContractCaller) (*DKGCaller, error) {
	contract, err := bindDKG(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DKGCaller{contract: contract}, nil
}

// NewDKGTransactor creates a new write-only instance of DKG, bound to a specific deployed contract.
func NewDKGTransactor(address common.Address, transactor bind.ContractTransactor) (*DKGTransactor, error) {
	contract, err := bindDKG(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DKGTransactor{contract: contract}, nil
}

// NewDKGFilterer creates a new log filterer instance of DKG, bound to a specific deployed contract.
func NewDKGFilterer(address common.Address, filterer bind.ContractFilterer) (*DKGFilterer, error) {
	contract, err := bindDKG(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DKGFilterer{contract: contract}, nil
}

// bindDKG binds a generic wrapper to an already deployed contract.
func bindDKG(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DKGMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DKG *DKGRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DKG.Contract.DKGCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DKG *DKGRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DKG.Contract.DKGTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DKG *DKGRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DKG.Contract.DKGTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DKG *DKGCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DKG.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DKG *DKGTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DKG.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DKG *DKGTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DKG.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_DKG *DKGCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_DKG *DKGSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _DKG.Contract.UPGRADEINTERFACEVERSION(&_DKG.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_DKG *DKGCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _DKG.Contract.UPGRADEINTERFACEVERSION(&_DKG.CallOpts)
}

// EnclaveTypeData is a free data retrieval call binding the contract method 0x1b0d13ff.
//
// Solidity: function enclaveTypeData(bytes32 enclaveType) view returns((bytes32,address))
func (_DKG *DKGCaller) EnclaveTypeData(opts *bind.CallOpts, enclaveType [32]byte) (IDKGEnclaveTypeData, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "enclaveTypeData", enclaveType)

	if err != nil {
		return *new(IDKGEnclaveTypeData), err
	}

	out0 := *abi.ConvertType(out[0], new(IDKGEnclaveTypeData)).(*IDKGEnclaveTypeData)

	return out0, err

}

// EnclaveTypeData is a free data retrieval call binding the contract method 0x1b0d13ff.
//
// Solidity: function enclaveTypeData(bytes32 enclaveType) view returns((bytes32,address))
func (_DKG *DKGSession) EnclaveTypeData(enclaveType [32]byte) (IDKGEnclaveTypeData, error) {
	return _DKG.Contract.EnclaveTypeData(&_DKG.CallOpts, enclaveType)
}

// EnclaveTypeData is a free data retrieval call binding the contract method 0x1b0d13ff.
//
// Solidity: function enclaveTypeData(bytes32 enclaveType) view returns((bytes32,address))
func (_DKG *DKGCallerSession) EnclaveTypeData(enclaveType [32]byte) (IDKGEnclaveTypeData, error) {
	return _DKG.Contract.EnclaveTypeData(&_DKG.CallOpts, enclaveType)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_DKG *DKGCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_DKG *DKGSession) Fee() (*big.Int, error) {
	return _DKG.Contract.Fee(&_DKG.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_DKG *DKGCallerSession) Fee() (*big.Int, error) {
	return _DKG.Contract.Fee(&_DKG.CallOpts)
}

// IsEnclaveTypeWhitelisted is a free data retrieval call binding the contract method 0xd3d3409b.
//
// Solidity: function isEnclaveTypeWhitelisted(bytes32 enclaveType) view returns(bool)
func (_DKG *DKGCaller) IsEnclaveTypeWhitelisted(opts *bind.CallOpts, enclaveType [32]byte) (bool, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "isEnclaveTypeWhitelisted", enclaveType)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEnclaveTypeWhitelisted is a free data retrieval call binding the contract method 0xd3d3409b.
//
// Solidity: function isEnclaveTypeWhitelisted(bytes32 enclaveType) view returns(bool)
func (_DKG *DKGSession) IsEnclaveTypeWhitelisted(enclaveType [32]byte) (bool, error) {
	return _DKG.Contract.IsEnclaveTypeWhitelisted(&_DKG.CallOpts, enclaveType)
}

// IsEnclaveTypeWhitelisted is a free data retrieval call binding the contract method 0xd3d3409b.
//
// Solidity: function isEnclaveTypeWhitelisted(bytes32 enclaveType) view returns(bool)
func (_DKG *DKGCallerSession) IsEnclaveTypeWhitelisted(enclaveType [32]byte) (bool, error) {
	return _DKG.Contract.IsEnclaveTypeWhitelisted(&_DKG.CallOpts, enclaveType)
}

// MinReqFinalizedParticipants is a free data retrieval call binding the contract method 0xb72d3227.
//
// Solidity: function minReqFinalizedParticipants() view returns(uint256)
func (_DKG *DKGCaller) MinReqFinalizedParticipants(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "minReqFinalizedParticipants")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinReqFinalizedParticipants is a free data retrieval call binding the contract method 0xb72d3227.
//
// Solidity: function minReqFinalizedParticipants() view returns(uint256)
func (_DKG *DKGSession) MinReqFinalizedParticipants() (*big.Int, error) {
	return _DKG.Contract.MinReqFinalizedParticipants(&_DKG.CallOpts)
}

// MinReqFinalizedParticipants is a free data retrieval call binding the contract method 0xb72d3227.
//
// Solidity: function minReqFinalizedParticipants() view returns(uint256)
func (_DKG *DKGCallerSession) MinReqFinalizedParticipants() (*big.Int, error) {
	return _DKG.Contract.MinReqFinalizedParticipants(&_DKG.CallOpts)
}

// MinReqRegisteredParticipants is a free data retrieval call binding the contract method 0xed3e270d.
//
// Solidity: function minReqRegisteredParticipants() view returns(uint256)
func (_DKG *DKGCaller) MinReqRegisteredParticipants(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "minReqRegisteredParticipants")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinReqRegisteredParticipants is a free data retrieval call binding the contract method 0xed3e270d.
//
// Solidity: function minReqRegisteredParticipants() view returns(uint256)
func (_DKG *DKGSession) MinReqRegisteredParticipants() (*big.Int, error) {
	return _DKG.Contract.MinReqRegisteredParticipants(&_DKG.CallOpts)
}

// MinReqRegisteredParticipants is a free data retrieval call binding the contract method 0xed3e270d.
//
// Solidity: function minReqRegisteredParticipants() view returns(uint256)
func (_DKG *DKGCallerSession) MinReqRegisteredParticipants() (*big.Int, error) {
	return _DKG.Contract.MinReqRegisteredParticipants(&_DKG.CallOpts)
}

// OperationalThreshold is a free data retrieval call binding the contract method 0x2adbc9b6.
//
// Solidity: function operationalThreshold() view returns(uint256)
func (_DKG *DKGCaller) OperationalThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "operationalThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OperationalThreshold is a free data retrieval call binding the contract method 0x2adbc9b6.
//
// Solidity: function operationalThreshold() view returns(uint256)
func (_DKG *DKGSession) OperationalThreshold() (*big.Int, error) {
	return _DKG.Contract.OperationalThreshold(&_DKG.CallOpts)
}

// OperationalThreshold is a free data retrieval call binding the contract method 0x2adbc9b6.
//
// Solidity: function operationalThreshold() view returns(uint256)
func (_DKG *DKGCallerSession) OperationalThreshold() (*big.Int, error) {
	return _DKG.Contract.OperationalThreshold(&_DKG.CallOpts)
}

// OperationalThresholdBasis is a free data retrieval call binding the contract method 0xa84565e7.
//
// Solidity: function operationalThresholdBasis() view returns(uint256)
func (_DKG *DKGCaller) OperationalThresholdBasis(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "operationalThresholdBasis")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OperationalThresholdBasis is a free data retrieval call binding the contract method 0xa84565e7.
//
// Solidity: function operationalThresholdBasis() view returns(uint256)
func (_DKG *DKGSession) OperationalThresholdBasis() (*big.Int, error) {
	return _DKG.Contract.OperationalThresholdBasis(&_DKG.CallOpts)
}

// OperationalThresholdBasis is a free data retrieval call binding the contract method 0xa84565e7.
//
// Solidity: function operationalThresholdBasis() view returns(uint256)
func (_DKG *DKGCallerSession) OperationalThresholdBasis() (*big.Int, error) {
	return _DKG.Contract.OperationalThresholdBasis(&_DKG.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DKG *DKGCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DKG *DKGSession) Owner() (common.Address, error) {
	return _DKG.Contract.Owner(&_DKG.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DKG *DKGCallerSession) Owner() (common.Address, error) {
	return _DKG.Contract.Owner(&_DKG.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_DKG *DKGCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_DKG *DKGSession) Paused() (bool, error) {
	return _DKG.Contract.Paused(&_DKG.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_DKG *DKGCallerSession) Paused() (bool, error) {
	return _DKG.Contract.Paused(&_DKG.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_DKG *DKGCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_DKG *DKGSession) PendingOwner() (common.Address, error) {
	return _DKG.Contract.PendingOwner(&_DKG.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_DKG *DKGCallerSession) PendingOwner() (common.Address, error) {
	return _DKG.Contract.PendingOwner(&_DKG.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_DKG *DKGCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_DKG *DKGSession) ProxiableUUID() ([32]byte, error) {
	return _DKG.Contract.ProxiableUUID(&_DKG.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_DKG *DKGCallerSession) ProxiableUUID() ([32]byte, error) {
	return _DKG.Contract.ProxiableUUID(&_DKG.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DKG *DKGTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DKG *DKGSession) AcceptOwnership() (*types.Transaction, error) {
	return _DKG.Contract.AcceptOwnership(&_DKG.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DKG *DKGTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DKG.Contract.AcceptOwnership(&_DKG.TransactOpts)
}

// AuthenticateEnclaveReport is a paid mutator transaction binding the contract method 0x9f5fc114.
//
// Solidity: function authenticateEnclaveReport(bytes enclaveReport, (uint32,address,bytes32,bytes,bytes) enclaveInstanceData, bytes validationContext) payable returns()
func (_DKG *DKGTransactor) AuthenticateEnclaveReport(opts *bind.TransactOpts, enclaveReport []byte, enclaveInstanceData IDKGEnclaveInstanceData, validationContext []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "authenticateEnclaveReport", enclaveReport, enclaveInstanceData, validationContext)
}

// AuthenticateEnclaveReport is a paid mutator transaction binding the contract method 0x9f5fc114.
//
// Solidity: function authenticateEnclaveReport(bytes enclaveReport, (uint32,address,bytes32,bytes,bytes) enclaveInstanceData, bytes validationContext) payable returns()
func (_DKG *DKGSession) AuthenticateEnclaveReport(enclaveReport []byte, enclaveInstanceData IDKGEnclaveInstanceData, validationContext []byte) (*types.Transaction, error) {
	return _DKG.Contract.AuthenticateEnclaveReport(&_DKG.TransactOpts, enclaveReport, enclaveInstanceData, validationContext)
}

// AuthenticateEnclaveReport is a paid mutator transaction binding the contract method 0x9f5fc114.
//
// Solidity: function authenticateEnclaveReport(bytes enclaveReport, (uint32,address,bytes32,bytes,bytes) enclaveInstanceData, bytes validationContext) payable returns()
func (_DKG *DKGTransactorSession) AuthenticateEnclaveReport(enclaveReport []byte, enclaveInstanceData IDKGEnclaveInstanceData, validationContext []byte) (*types.Transaction, error) {
	return _DKG.Contract.AuthenticateEnclaveReport(&_DKG.TransactOpts, enclaveReport, enclaveInstanceData, validationContext)
}

// Finalize is a paid mutator transaction binding the contract method 0xa80e9f42.
//
// Solidity: function finalize(uint32 round, address validatorAddr, bytes32 enclaveType, bytes32 participantsRoot, bytes globalPubKey, bytes[] publicCoeffs, bytes pubKeyShare, bytes signature) payable returns()
func (_DKG *DKGTransactor) Finalize(opts *bind.TransactOpts, round uint32, validatorAddr common.Address, enclaveType [32]byte, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "finalize", round, validatorAddr, enclaveType, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare, signature)
}

// Finalize is a paid mutator transaction binding the contract method 0xa80e9f42.
//
// Solidity: function finalize(uint32 round, address validatorAddr, bytes32 enclaveType, bytes32 participantsRoot, bytes globalPubKey, bytes[] publicCoeffs, bytes pubKeyShare, bytes signature) payable returns()
func (_DKG *DKGSession) Finalize(round uint32, validatorAddr common.Address, enclaveType [32]byte, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.Finalize(&_DKG.TransactOpts, round, validatorAddr, enclaveType, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare, signature)
}

// Finalize is a paid mutator transaction binding the contract method 0xa80e9f42.
//
// Solidity: function finalize(uint32 round, address validatorAddr, bytes32 enclaveType, bytes32 participantsRoot, bytes globalPubKey, bytes[] publicCoeffs, bytes pubKeyShare, bytes signature) payable returns()
func (_DKG *DKGTransactorSession) Finalize(round uint32, validatorAddr common.Address, enclaveType [32]byte, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.Finalize(&_DKG.TransactOpts, round, validatorAddr, enclaveType, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare, signature)
}

// Initialize is a paid mutator transaction binding the contract method 0xf92ad219.
//
// Solidity: function initialize(address owner, uint256 minReqRegisteredParticipants, uint256 minReqFinalizedParticipants, uint256 operationalThreshold, uint256 fee) returns()
func (_DKG *DKGTransactor) Initialize(opts *bind.TransactOpts, owner common.Address, minReqRegisteredParticipants *big.Int, minReqFinalizedParticipants *big.Int, operationalThreshold *big.Int, fee *big.Int) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "initialize", owner, minReqRegisteredParticipants, minReqFinalizedParticipants, operationalThreshold, fee)
}

// Initialize is a paid mutator transaction binding the contract method 0xf92ad219.
//
// Solidity: function initialize(address owner, uint256 minReqRegisteredParticipants, uint256 minReqFinalizedParticipants, uint256 operationalThreshold, uint256 fee) returns()
func (_DKG *DKGSession) Initialize(owner common.Address, minReqRegisteredParticipants *big.Int, minReqFinalizedParticipants *big.Int, operationalThreshold *big.Int, fee *big.Int) (*types.Transaction, error) {
	return _DKG.Contract.Initialize(&_DKG.TransactOpts, owner, minReqRegisteredParticipants, minReqFinalizedParticipants, operationalThreshold, fee)
}

// Initialize is a paid mutator transaction binding the contract method 0xf92ad219.
//
// Solidity: function initialize(address owner, uint256 minReqRegisteredParticipants, uint256 minReqFinalizedParticipants, uint256 operationalThreshold, uint256 fee) returns()
func (_DKG *DKGTransactorSession) Initialize(owner common.Address, minReqRegisteredParticipants *big.Int, minReqFinalizedParticipants *big.Int, operationalThreshold *big.Int, fee *big.Int) (*types.Transaction, error) {
	return _DKG.Contract.Initialize(&_DKG.TransactOpts, owner, minReqRegisteredParticipants, minReqFinalizedParticipants, operationalThreshold, fee)
}

// Register is a paid mutator transaction binding the contract method 0x804faf5b.
//
// Solidity: function register(bytes enclaveReport, (uint32,address,bytes32,bytes,bytes) enclaveInstanceData, uint256 startBlockHeight, bytes32 startBlockHash, bytes validationContext) payable returns()
func (_DKG *DKGTransactor) Register(opts *bind.TransactOpts, enclaveReport []byte, enclaveInstanceData IDKGEnclaveInstanceData, startBlockHeight *big.Int, startBlockHash [32]byte, validationContext []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "register", enclaveReport, enclaveInstanceData, startBlockHeight, startBlockHash, validationContext)
}

// Register is a paid mutator transaction binding the contract method 0x804faf5b.
//
// Solidity: function register(bytes enclaveReport, (uint32,address,bytes32,bytes,bytes) enclaveInstanceData, uint256 startBlockHeight, bytes32 startBlockHash, bytes validationContext) payable returns()
func (_DKG *DKGSession) Register(enclaveReport []byte, enclaveInstanceData IDKGEnclaveInstanceData, startBlockHeight *big.Int, startBlockHash [32]byte, validationContext []byte) (*types.Transaction, error) {
	return _DKG.Contract.Register(&_DKG.TransactOpts, enclaveReport, enclaveInstanceData, startBlockHeight, startBlockHash, validationContext)
}

// Register is a paid mutator transaction binding the contract method 0x804faf5b.
//
// Solidity: function register(bytes enclaveReport, (uint32,address,bytes32,bytes,bytes) enclaveInstanceData, uint256 startBlockHeight, bytes32 startBlockHash, bytes validationContext) payable returns()
func (_DKG *DKGTransactorSession) Register(enclaveReport []byte, enclaveInstanceData IDKGEnclaveInstanceData, startBlockHeight *big.Int, startBlockHash [32]byte, validationContext []byte) (*types.Transaction, error) {
	return _DKG.Contract.Register(&_DKG.TransactOpts, enclaveReport, enclaveInstanceData, startBlockHeight, startBlockHash, validationContext)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DKG *DKGTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DKG *DKGSession) RenounceOwnership() (*types.Transaction, error) {
	return _DKG.Contract.RenounceOwnership(&_DKG.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DKG *DKGTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DKG.Contract.RenounceOwnership(&_DKG.TransactOpts)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 newFee) returns()
func (_DKG *DKGTransactor) SetFee(opts *bind.TransactOpts, newFee *big.Int) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "setFee", newFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 newFee) returns()
func (_DKG *DKGSession) SetFee(newFee *big.Int) (*types.Transaction, error) {
	return _DKG.Contract.SetFee(&_DKG.TransactOpts, newFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 newFee) returns()
func (_DKG *DKGTransactorSession) SetFee(newFee *big.Int) (*types.Transaction, error) {
	return _DKG.Contract.SetFee(&_DKG.TransactOpts, newFee)
}

// SetMinReqFinalizedParticipants is a paid mutator transaction binding the contract method 0x427c294f.
//
// Solidity: function setMinReqFinalizedParticipants(uint256 newMinReqFinalizedParticipants) returns()
func (_DKG *DKGTransactor) SetMinReqFinalizedParticipants(opts *bind.TransactOpts, newMinReqFinalizedParticipants *big.Int) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "setMinReqFinalizedParticipants", newMinReqFinalizedParticipants)
}

// SetMinReqFinalizedParticipants is a paid mutator transaction binding the contract method 0x427c294f.
//
// Solidity: function setMinReqFinalizedParticipants(uint256 newMinReqFinalizedParticipants) returns()
func (_DKG *DKGSession) SetMinReqFinalizedParticipants(newMinReqFinalizedParticipants *big.Int) (*types.Transaction, error) {
	return _DKG.Contract.SetMinReqFinalizedParticipants(&_DKG.TransactOpts, newMinReqFinalizedParticipants)
}

// SetMinReqFinalizedParticipants is a paid mutator transaction binding the contract method 0x427c294f.
//
// Solidity: function setMinReqFinalizedParticipants(uint256 newMinReqFinalizedParticipants) returns()
func (_DKG *DKGTransactorSession) SetMinReqFinalizedParticipants(newMinReqFinalizedParticipants *big.Int) (*types.Transaction, error) {
	return _DKG.Contract.SetMinReqFinalizedParticipants(&_DKG.TransactOpts, newMinReqFinalizedParticipants)
}

// SetMinReqRegisteredParticipants is a paid mutator transaction binding the contract method 0xd78a5416.
//
// Solidity: function setMinReqRegisteredParticipants(uint256 newMinReqRegisteredParticipants) returns()
func (_DKG *DKGTransactor) SetMinReqRegisteredParticipants(opts *bind.TransactOpts, newMinReqRegisteredParticipants *big.Int) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "setMinReqRegisteredParticipants", newMinReqRegisteredParticipants)
}

// SetMinReqRegisteredParticipants is a paid mutator transaction binding the contract method 0xd78a5416.
//
// Solidity: function setMinReqRegisteredParticipants(uint256 newMinReqRegisteredParticipants) returns()
func (_DKG *DKGSession) SetMinReqRegisteredParticipants(newMinReqRegisteredParticipants *big.Int) (*types.Transaction, error) {
	return _DKG.Contract.SetMinReqRegisteredParticipants(&_DKG.TransactOpts, newMinReqRegisteredParticipants)
}

// SetMinReqRegisteredParticipants is a paid mutator transaction binding the contract method 0xd78a5416.
//
// Solidity: function setMinReqRegisteredParticipants(uint256 newMinReqRegisteredParticipants) returns()
func (_DKG *DKGTransactorSession) SetMinReqRegisteredParticipants(newMinReqRegisteredParticipants *big.Int) (*types.Transaction, error) {
	return _DKG.Contract.SetMinReqRegisteredParticipants(&_DKG.TransactOpts, newMinReqRegisteredParticipants)
}

// SetOperationalThreshold is a paid mutator transaction binding the contract method 0x5fbbc42d.
//
// Solidity: function setOperationalThreshold(uint256 newOperationalThreshold) returns()
func (_DKG *DKGTransactor) SetOperationalThreshold(opts *bind.TransactOpts, newOperationalThreshold *big.Int) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "setOperationalThreshold", newOperationalThreshold)
}

// SetOperationalThreshold is a paid mutator transaction binding the contract method 0x5fbbc42d.
//
// Solidity: function setOperationalThreshold(uint256 newOperationalThreshold) returns()
func (_DKG *DKGSession) SetOperationalThreshold(newOperationalThreshold *big.Int) (*types.Transaction, error) {
	return _DKG.Contract.SetOperationalThreshold(&_DKG.TransactOpts, newOperationalThreshold)
}

// SetOperationalThreshold is a paid mutator transaction binding the contract method 0x5fbbc42d.
//
// Solidity: function setOperationalThreshold(uint256 newOperationalThreshold) returns()
func (_DKG *DKGTransactorSession) SetOperationalThreshold(newOperationalThreshold *big.Int) (*types.Transaction, error) {
	return _DKG.Contract.SetOperationalThreshold(&_DKG.TransactOpts, newOperationalThreshold)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DKG *DKGTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DKG *DKGSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DKG.Contract.TransferOwnership(&_DKG.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DKG *DKGTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DKG.Contract.TransferOwnership(&_DKG.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_DKG *DKGTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_DKG *DKGSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _DKG.Contract.UpgradeToAndCall(&_DKG.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_DKG *DKGTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _DKG.Contract.UpgradeToAndCall(&_DKG.TransactOpts, newImplementation, data)
}

// WhitelistEnclaveType is a paid mutator transaction binding the contract method 0xfb3dde66.
//
// Solidity: function whitelistEnclaveType(bytes32 enclaveType, (bytes32,address) enclaveTypeData, bool isWhitelisted) returns()
func (_DKG *DKGTransactor) WhitelistEnclaveType(opts *bind.TransactOpts, enclaveType [32]byte, enclaveTypeData IDKGEnclaveTypeData, isWhitelisted bool) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "whitelistEnclaveType", enclaveType, enclaveTypeData, isWhitelisted)
}

// WhitelistEnclaveType is a paid mutator transaction binding the contract method 0xfb3dde66.
//
// Solidity: function whitelistEnclaveType(bytes32 enclaveType, (bytes32,address) enclaveTypeData, bool isWhitelisted) returns()
func (_DKG *DKGSession) WhitelistEnclaveType(enclaveType [32]byte, enclaveTypeData IDKGEnclaveTypeData, isWhitelisted bool) (*types.Transaction, error) {
	return _DKG.Contract.WhitelistEnclaveType(&_DKG.TransactOpts, enclaveType, enclaveTypeData, isWhitelisted)
}

// WhitelistEnclaveType is a paid mutator transaction binding the contract method 0xfb3dde66.
//
// Solidity: function whitelistEnclaveType(bytes32 enclaveType, (bytes32,address) enclaveTypeData, bool isWhitelisted) returns()
func (_DKG *DKGTransactorSession) WhitelistEnclaveType(enclaveType [32]byte, enclaveTypeData IDKGEnclaveTypeData, isWhitelisted bool) (*types.Transaction, error) {
	return _DKG.Contract.WhitelistEnclaveType(&_DKG.TransactOpts, enclaveType, enclaveTypeData, isWhitelisted)
}

// DKGEnclaveTypeWhitelistedIterator is returned from FilterEnclaveTypeWhitelisted and is used to iterate over the raw logs and unpacked data for EnclaveTypeWhitelisted events raised by the DKG contract.
type DKGEnclaveTypeWhitelistedIterator struct {
	Event *DKGEnclaveTypeWhitelisted // Event containing the contract specifics and raw log

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
func (it *DKGEnclaveTypeWhitelistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGEnclaveTypeWhitelisted)
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
		it.Event = new(DKGEnclaveTypeWhitelisted)
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
func (it *DKGEnclaveTypeWhitelistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGEnclaveTypeWhitelistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGEnclaveTypeWhitelisted represents a EnclaveTypeWhitelisted event raised by the DKG contract.
type DKGEnclaveTypeWhitelisted struct {
	EnclaveType        [32]byte
	CodeCommitment     [32]byte
	ValidationHookAddr common.Address
	IsWhitelisted      bool
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterEnclaveTypeWhitelisted is a free log retrieval operation binding the contract event 0x0be17b82f2eb1d1de6349e90606850d99d513c3a3c418ea90b7f0cd0a5438630.
//
// Solidity: event EnclaveTypeWhitelisted(bytes32 enclaveType, bytes32 codeCommitment, address validationHookAddr, bool isWhitelisted)
func (_DKG *DKGFilterer) FilterEnclaveTypeWhitelisted(opts *bind.FilterOpts) (*DKGEnclaveTypeWhitelistedIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "EnclaveTypeWhitelisted")
	if err != nil {
		return nil, err
	}
	return &DKGEnclaveTypeWhitelistedIterator{contract: _DKG.contract, event: "EnclaveTypeWhitelisted", logs: logs, sub: sub}, nil
}

// WatchEnclaveTypeWhitelisted is a free log subscription operation binding the contract event 0x0be17b82f2eb1d1de6349e90606850d99d513c3a3c418ea90b7f0cd0a5438630.
//
// Solidity: event EnclaveTypeWhitelisted(bytes32 enclaveType, bytes32 codeCommitment, address validationHookAddr, bool isWhitelisted)
func (_DKG *DKGFilterer) WatchEnclaveTypeWhitelisted(opts *bind.WatchOpts, sink chan<- *DKGEnclaveTypeWhitelisted) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "EnclaveTypeWhitelisted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGEnclaveTypeWhitelisted)
				if err := _DKG.contract.UnpackLog(event, "EnclaveTypeWhitelisted", log); err != nil {
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

// ParseEnclaveTypeWhitelisted is a log parse operation binding the contract event 0x0be17b82f2eb1d1de6349e90606850d99d513c3a3c418ea90b7f0cd0a5438630.
//
// Solidity: event EnclaveTypeWhitelisted(bytes32 enclaveType, bytes32 codeCommitment, address validationHookAddr, bool isWhitelisted)
func (_DKG *DKGFilterer) ParseEnclaveTypeWhitelisted(log types.Log) (*DKGEnclaveTypeWhitelisted, error) {
	event := new(DKGEnclaveTypeWhitelisted)
	if err := _DKG.contract.UnpackLog(event, "EnclaveTypeWhitelisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGFeeSetIterator is returned from FilterFeeSet and is used to iterate over the raw logs and unpacked data for FeeSet events raised by the DKG contract.
type DKGFeeSetIterator struct {
	Event *DKGFeeSet // Event containing the contract specifics and raw log

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
func (it *DKGFeeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGFeeSet)
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
		it.Event = new(DKGFeeSet)
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
func (it *DKGFeeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGFeeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGFeeSet represents a FeeSet event raised by the DKG contract.
type DKGFeeSet struct {
	NewFee *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeeSet is a free log retrieval operation binding the contract event 0x20461e09b8e557b77e107939f9ce6544698123aad0fc964ac5cc59b7df2e608f.
//
// Solidity: event FeeSet(uint256 newFee)
func (_DKG *DKGFilterer) FilterFeeSet(opts *bind.FilterOpts) (*DKGFeeSetIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "FeeSet")
	if err != nil {
		return nil, err
	}
	return &DKGFeeSetIterator{contract: _DKG.contract, event: "FeeSet", logs: logs, sub: sub}, nil
}

// WatchFeeSet is a free log subscription operation binding the contract event 0x20461e09b8e557b77e107939f9ce6544698123aad0fc964ac5cc59b7df2e608f.
//
// Solidity: event FeeSet(uint256 newFee)
func (_DKG *DKGFilterer) WatchFeeSet(opts *bind.WatchOpts, sink chan<- *DKGFeeSet) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "FeeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGFeeSet)
				if err := _DKG.contract.UnpackLog(event, "FeeSet", log); err != nil {
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

// ParseFeeSet is a log parse operation binding the contract event 0x20461e09b8e557b77e107939f9ce6544698123aad0fc964ac5cc59b7df2e608f.
//
// Solidity: event FeeSet(uint256 newFee)
func (_DKG *DKGFilterer) ParseFeeSet(log types.Log) (*DKGFeeSet, error) {
	event := new(DKGFeeSet)
	if err := _DKG.contract.UnpackLog(event, "FeeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGFinalizedIterator is returned from FilterFinalized and is used to iterate over the raw logs and unpacked data for Finalized events raised by the DKG contract.
type DKGFinalizedIterator struct {
	Event *DKGFinalized // Event containing the contract specifics and raw log

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
func (it *DKGFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGFinalized)
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
		it.Event = new(DKGFinalized)
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
func (it *DKGFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGFinalized represents a Finalized event raised by the DKG contract.
type DKGFinalized struct {
	Round            uint32
	ValidatorAddr    common.Address
	EnclaveType      [32]byte
	CodeCommitment   [32]byte
	ParticipantsRoot [32]byte
	GlobalPubKey     []byte
	PublicCoeffs     [][]byte
	PubKeyShare      []byte
	Signature        []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterFinalized is a free log retrieval operation binding the contract event 0xeeeaac510415cc754b05e40e2e701b7d22d7d764d2a0b232c26d9eea3b62565d.
//
// Solidity: event Finalized(uint32 round, address indexed validatorAddr, bytes32 enclaveType, bytes32 codeCommitment, bytes32 participantsRoot, bytes globalPubKey, bytes[] publicCoeffs, bytes pubKeyShare, bytes signature)
func (_DKG *DKGFilterer) FilterFinalized(opts *bind.FilterOpts, validatorAddr []common.Address) (*DKGFinalizedIterator, error) {

	var validatorAddrRule []interface{}
	for _, validatorAddrItem := range validatorAddr {
		validatorAddrRule = append(validatorAddrRule, validatorAddrItem)
	}

	logs, sub, err := _DKG.contract.FilterLogs(opts, "Finalized", validatorAddrRule)
	if err != nil {
		return nil, err
	}
	return &DKGFinalizedIterator{contract: _DKG.contract, event: "Finalized", logs: logs, sub: sub}, nil
}

// WatchFinalized is a free log subscription operation binding the contract event 0xeeeaac510415cc754b05e40e2e701b7d22d7d764d2a0b232c26d9eea3b62565d.
//
// Solidity: event Finalized(uint32 round, address indexed validatorAddr, bytes32 enclaveType, bytes32 codeCommitment, bytes32 participantsRoot, bytes globalPubKey, bytes[] publicCoeffs, bytes pubKeyShare, bytes signature)
func (_DKG *DKGFilterer) WatchFinalized(opts *bind.WatchOpts, sink chan<- *DKGFinalized, validatorAddr []common.Address) (event.Subscription, error) {

	var validatorAddrRule []interface{}
	for _, validatorAddrItem := range validatorAddr {
		validatorAddrRule = append(validatorAddrRule, validatorAddrItem)
	}

	logs, sub, err := _DKG.contract.WatchLogs(opts, "Finalized", validatorAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGFinalized)
				if err := _DKG.contract.UnpackLog(event, "Finalized", log); err != nil {
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

// ParseFinalized is a log parse operation binding the contract event 0xeeeaac510415cc754b05e40e2e701b7d22d7d764d2a0b232c26d9eea3b62565d.
//
// Solidity: event Finalized(uint32 round, address indexed validatorAddr, bytes32 enclaveType, bytes32 codeCommitment, bytes32 participantsRoot, bytes globalPubKey, bytes[] publicCoeffs, bytes pubKeyShare, bytes signature)
func (_DKG *DKGFilterer) ParseFinalized(log types.Log) (*DKGFinalized, error) {
	event := new(DKGFinalized)
	if err := _DKG.contract.UnpackLog(event, "Finalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the DKG contract.
type DKGInitializedIterator struct {
	Event *DKGInitialized // Event containing the contract specifics and raw log

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
func (it *DKGInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGInitialized)
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
		it.Event = new(DKGInitialized)
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
func (it *DKGInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGInitialized represents a Initialized event raised by the DKG contract.
type DKGInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DKG *DKGFilterer) FilterInitialized(opts *bind.FilterOpts) (*DKGInitializedIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DKGInitializedIterator{contract: _DKG.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DKG *DKGFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DKGInitialized) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGInitialized)
				if err := _DKG.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_DKG *DKGFilterer) ParseInitialized(log types.Log) (*DKGInitialized, error) {
	event := new(DKGInitialized)
	if err := _DKG.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGMinReqFinalizedParticipantsSetIterator is returned from FilterMinReqFinalizedParticipantsSet and is used to iterate over the raw logs and unpacked data for MinReqFinalizedParticipantsSet events raised by the DKG contract.
type DKGMinReqFinalizedParticipantsSetIterator struct {
	Event *DKGMinReqFinalizedParticipantsSet // Event containing the contract specifics and raw log

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
func (it *DKGMinReqFinalizedParticipantsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGMinReqFinalizedParticipantsSet)
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
		it.Event = new(DKGMinReqFinalizedParticipantsSet)
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
func (it *DKGMinReqFinalizedParticipantsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGMinReqFinalizedParticipantsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGMinReqFinalizedParticipantsSet represents a MinReqFinalizedParticipantsSet event raised by the DKG contract.
type DKGMinReqFinalizedParticipantsSet struct {
	NewMinReqFinalizedParticipants *big.Int
	Raw                            types.Log // Blockchain specific contextual infos
}

// FilterMinReqFinalizedParticipantsSet is a free log retrieval operation binding the contract event 0xced3da4b96d055042d4df66edffabf4cdaa5909bab55d737f34d71e50a07ae58.
//
// Solidity: event MinReqFinalizedParticipantsSet(uint256 newMinReqFinalizedParticipants)
func (_DKG *DKGFilterer) FilterMinReqFinalizedParticipantsSet(opts *bind.FilterOpts) (*DKGMinReqFinalizedParticipantsSetIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "MinReqFinalizedParticipantsSet")
	if err != nil {
		return nil, err
	}
	return &DKGMinReqFinalizedParticipantsSetIterator{contract: _DKG.contract, event: "MinReqFinalizedParticipantsSet", logs: logs, sub: sub}, nil
}

// WatchMinReqFinalizedParticipantsSet is a free log subscription operation binding the contract event 0xced3da4b96d055042d4df66edffabf4cdaa5909bab55d737f34d71e50a07ae58.
//
// Solidity: event MinReqFinalizedParticipantsSet(uint256 newMinReqFinalizedParticipants)
func (_DKG *DKGFilterer) WatchMinReqFinalizedParticipantsSet(opts *bind.WatchOpts, sink chan<- *DKGMinReqFinalizedParticipantsSet) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "MinReqFinalizedParticipantsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGMinReqFinalizedParticipantsSet)
				if err := _DKG.contract.UnpackLog(event, "MinReqFinalizedParticipantsSet", log); err != nil {
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

// ParseMinReqFinalizedParticipantsSet is a log parse operation binding the contract event 0xced3da4b96d055042d4df66edffabf4cdaa5909bab55d737f34d71e50a07ae58.
//
// Solidity: event MinReqFinalizedParticipantsSet(uint256 newMinReqFinalizedParticipants)
func (_DKG *DKGFilterer) ParseMinReqFinalizedParticipantsSet(log types.Log) (*DKGMinReqFinalizedParticipantsSet, error) {
	event := new(DKGMinReqFinalizedParticipantsSet)
	if err := _DKG.contract.UnpackLog(event, "MinReqFinalizedParticipantsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGMinReqRegisteredParticipantsSetIterator is returned from FilterMinReqRegisteredParticipantsSet and is used to iterate over the raw logs and unpacked data for MinReqRegisteredParticipantsSet events raised by the DKG contract.
type DKGMinReqRegisteredParticipantsSetIterator struct {
	Event *DKGMinReqRegisteredParticipantsSet // Event containing the contract specifics and raw log

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
func (it *DKGMinReqRegisteredParticipantsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGMinReqRegisteredParticipantsSet)
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
		it.Event = new(DKGMinReqRegisteredParticipantsSet)
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
func (it *DKGMinReqRegisteredParticipantsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGMinReqRegisteredParticipantsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGMinReqRegisteredParticipantsSet represents a MinReqRegisteredParticipantsSet event raised by the DKG contract.
type DKGMinReqRegisteredParticipantsSet struct {
	NewMinReqRegisteredParticipants *big.Int
	Raw                             types.Log // Blockchain specific contextual infos
}

// FilterMinReqRegisteredParticipantsSet is a free log retrieval operation binding the contract event 0xcc68805432194238c4be6e6bec7f8c92631b27c0abd3e73f7f6467f41f9dfebd.
//
// Solidity: event MinReqRegisteredParticipantsSet(uint256 newMinReqRegisteredParticipants)
func (_DKG *DKGFilterer) FilterMinReqRegisteredParticipantsSet(opts *bind.FilterOpts) (*DKGMinReqRegisteredParticipantsSetIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "MinReqRegisteredParticipantsSet")
	if err != nil {
		return nil, err
	}
	return &DKGMinReqRegisteredParticipantsSetIterator{contract: _DKG.contract, event: "MinReqRegisteredParticipantsSet", logs: logs, sub: sub}, nil
}

// WatchMinReqRegisteredParticipantsSet is a free log subscription operation binding the contract event 0xcc68805432194238c4be6e6bec7f8c92631b27c0abd3e73f7f6467f41f9dfebd.
//
// Solidity: event MinReqRegisteredParticipantsSet(uint256 newMinReqRegisteredParticipants)
func (_DKG *DKGFilterer) WatchMinReqRegisteredParticipantsSet(opts *bind.WatchOpts, sink chan<- *DKGMinReqRegisteredParticipantsSet) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "MinReqRegisteredParticipantsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGMinReqRegisteredParticipantsSet)
				if err := _DKG.contract.UnpackLog(event, "MinReqRegisteredParticipantsSet", log); err != nil {
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

// ParseMinReqRegisteredParticipantsSet is a log parse operation binding the contract event 0xcc68805432194238c4be6e6bec7f8c92631b27c0abd3e73f7f6467f41f9dfebd.
//
// Solidity: event MinReqRegisteredParticipantsSet(uint256 newMinReqRegisteredParticipants)
func (_DKG *DKGFilterer) ParseMinReqRegisteredParticipantsSet(log types.Log) (*DKGMinReqRegisteredParticipantsSet, error) {
	event := new(DKGMinReqRegisteredParticipantsSet)
	if err := _DKG.contract.UnpackLog(event, "MinReqRegisteredParticipantsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGOperationalThresholdSetIterator is returned from FilterOperationalThresholdSet and is used to iterate over the raw logs and unpacked data for OperationalThresholdSet events raised by the DKG contract.
type DKGOperationalThresholdSetIterator struct {
	Event *DKGOperationalThresholdSet // Event containing the contract specifics and raw log

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
func (it *DKGOperationalThresholdSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGOperationalThresholdSet)
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
		it.Event = new(DKGOperationalThresholdSet)
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
func (it *DKGOperationalThresholdSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGOperationalThresholdSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGOperationalThresholdSet represents a OperationalThresholdSet event raised by the DKG contract.
type DKGOperationalThresholdSet struct {
	NewOperationalThreshold *big.Int
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterOperationalThresholdSet is a free log retrieval operation binding the contract event 0x90701eb5a3be7e3396be0be01f464764263a71e87a65140746a0e909384e990b.
//
// Solidity: event OperationalThresholdSet(uint256 newOperationalThreshold)
func (_DKG *DKGFilterer) FilterOperationalThresholdSet(opts *bind.FilterOpts) (*DKGOperationalThresholdSetIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "OperationalThresholdSet")
	if err != nil {
		return nil, err
	}
	return &DKGOperationalThresholdSetIterator{contract: _DKG.contract, event: "OperationalThresholdSet", logs: logs, sub: sub}, nil
}

// WatchOperationalThresholdSet is a free log subscription operation binding the contract event 0x90701eb5a3be7e3396be0be01f464764263a71e87a65140746a0e909384e990b.
//
// Solidity: event OperationalThresholdSet(uint256 newOperationalThreshold)
func (_DKG *DKGFilterer) WatchOperationalThresholdSet(opts *bind.WatchOpts, sink chan<- *DKGOperationalThresholdSet) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "OperationalThresholdSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGOperationalThresholdSet)
				if err := _DKG.contract.UnpackLog(event, "OperationalThresholdSet", log); err != nil {
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

// ParseOperationalThresholdSet is a log parse operation binding the contract event 0x90701eb5a3be7e3396be0be01f464764263a71e87a65140746a0e909384e990b.
//
// Solidity: event OperationalThresholdSet(uint256 newOperationalThreshold)
func (_DKG *DKGFilterer) ParseOperationalThresholdSet(log types.Log) (*DKGOperationalThresholdSet, error) {
	event := new(DKGOperationalThresholdSet)
	if err := _DKG.contract.UnpackLog(event, "OperationalThresholdSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the DKG contract.
type DKGOwnershipTransferStartedIterator struct {
	Event *DKGOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *DKGOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGOwnershipTransferStarted)
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
		it.Event = new(DKGOwnershipTransferStarted)
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
func (it *DKGOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the DKG contract.
type DKGOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_DKG *DKGFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DKGOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DKG.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DKGOwnershipTransferStartedIterator{contract: _DKG.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_DKG *DKGFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *DKGOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DKG.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGOwnershipTransferStarted)
				if err := _DKG.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_DKG *DKGFilterer) ParseOwnershipTransferStarted(log types.Log) (*DKGOwnershipTransferStarted, error) {
	event := new(DKGOwnershipTransferStarted)
	if err := _DKG.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DKG contract.
type DKGOwnershipTransferredIterator struct {
	Event *DKGOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DKGOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGOwnershipTransferred)
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
		it.Event = new(DKGOwnershipTransferred)
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
func (it *DKGOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGOwnershipTransferred represents a OwnershipTransferred event raised by the DKG contract.
type DKGOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DKG *DKGFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DKGOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DKG.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DKGOwnershipTransferredIterator{contract: _DKG.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DKG *DKGFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DKGOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DKG.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGOwnershipTransferred)
				if err := _DKG.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_DKG *DKGFilterer) ParseOwnershipTransferred(log types.Log) (*DKGOwnershipTransferred, error) {
	event := new(DKGOwnershipTransferred)
	if err := _DKG.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the DKG contract.
type DKGPausedIterator struct {
	Event *DKGPaused // Event containing the contract specifics and raw log

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
func (it *DKGPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGPaused)
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
		it.Event = new(DKGPaused)
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
func (it *DKGPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGPaused represents a Paused event raised by the DKG contract.
type DKGPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_DKG *DKGFilterer) FilterPaused(opts *bind.FilterOpts) (*DKGPausedIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &DKGPausedIterator{contract: _DKG.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_DKG *DKGFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *DKGPaused) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGPaused)
				if err := _DKG.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_DKG *DKGFilterer) ParsePaused(log types.Log) (*DKGPaused, error) {
	event := new(DKGPaused)
	if err := _DKG.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGRegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the DKG contract.
type DKGRegisteredIterator struct {
	Event *DKGRegistered // Event containing the contract specifics and raw log

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
func (it *DKGRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGRegistered)
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
		it.Event = new(DKGRegistered)
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
func (it *DKGRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGRegistered represents a Registered event raised by the DKG contract.
type DKGRegistered struct {
	EnclaveReport     []byte
	Round             uint32
	ValidatorAddr     common.Address
	EnclaveType       [32]byte
	EnclaveCommKey    []byte
	DkgPubKey         []byte
	CodeCommitment    [32]byte
	StartBlockHeight  *big.Int
	StartBlockHash    [32]byte
	ValidationContext []byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x66979019c139c21e4120d46d1cb1ae87145337ca67db0c998266f9f58ce8f2be.
//
// Solidity: event Registered(bytes enclaveReport, uint32 round, address indexed validatorAddr, bytes32 enclaveType, bytes enclaveCommKey, bytes dkgPubKey, bytes32 codeCommitment, uint256 startBlockHeight, bytes32 startBlockHash, bytes validationContext)
func (_DKG *DKGFilterer) FilterRegistered(opts *bind.FilterOpts, validatorAddr []common.Address) (*DKGRegisteredIterator, error) {

	var validatorAddrRule []interface{}
	for _, validatorAddrItem := range validatorAddr {
		validatorAddrRule = append(validatorAddrRule, validatorAddrItem)
	}

	logs, sub, err := _DKG.contract.FilterLogs(opts, "Registered", validatorAddrRule)
	if err != nil {
		return nil, err
	}
	return &DKGRegisteredIterator{contract: _DKG.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x66979019c139c21e4120d46d1cb1ae87145337ca67db0c998266f9f58ce8f2be.
//
// Solidity: event Registered(bytes enclaveReport, uint32 round, address indexed validatorAddr, bytes32 enclaveType, bytes enclaveCommKey, bytes dkgPubKey, bytes32 codeCommitment, uint256 startBlockHeight, bytes32 startBlockHash, bytes validationContext)
func (_DKG *DKGFilterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *DKGRegistered, validatorAddr []common.Address) (event.Subscription, error) {

	var validatorAddrRule []interface{}
	for _, validatorAddrItem := range validatorAddr {
		validatorAddrRule = append(validatorAddrRule, validatorAddrItem)
	}

	logs, sub, err := _DKG.contract.WatchLogs(opts, "Registered", validatorAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGRegistered)
				if err := _DKG.contract.UnpackLog(event, "Registered", log); err != nil {
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

// ParseRegistered is a log parse operation binding the contract event 0x66979019c139c21e4120d46d1cb1ae87145337ca67db0c998266f9f58ce8f2be.
//
// Solidity: event Registered(bytes enclaveReport, uint32 round, address indexed validatorAddr, bytes32 enclaveType, bytes enclaveCommKey, bytes dkgPubKey, bytes32 codeCommitment, uint256 startBlockHeight, bytes32 startBlockHash, bytes validationContext)
func (_DKG *DKGFilterer) ParseRegistered(log types.Log) (*DKGRegistered, error) {
	event := new(DKGRegistered)
	if err := _DKG.contract.UnpackLog(event, "Registered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the DKG contract.
type DKGUnpausedIterator struct {
	Event *DKGUnpaused // Event containing the contract specifics and raw log

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
func (it *DKGUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGUnpaused)
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
		it.Event = new(DKGUnpaused)
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
func (it *DKGUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGUnpaused represents a Unpaused event raised by the DKG contract.
type DKGUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_DKG *DKGFilterer) FilterUnpaused(opts *bind.FilterOpts) (*DKGUnpausedIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &DKGUnpausedIterator{contract: _DKG.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_DKG *DKGFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *DKGUnpaused) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGUnpaused)
				if err := _DKG.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_DKG *DKGFilterer) ParseUnpaused(log types.Log) (*DKGUnpaused, error) {
	event := new(DKGUnpaused)
	if err := _DKG.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the DKG contract.
type DKGUpgradedIterator struct {
	Event *DKGUpgraded // Event containing the contract specifics and raw log

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
func (it *DKGUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGUpgraded)
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
		it.Event = new(DKGUpgraded)
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
func (it *DKGUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGUpgraded represents a Upgraded event raised by the DKG contract.
type DKGUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DKG *DKGFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*DKGUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _DKG.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &DKGUpgradedIterator{contract: _DKG.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DKG *DKGFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *DKGUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _DKG.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGUpgraded)
				if err := _DKG.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_DKG *DKGFilterer) ParseUpgraded(log types.Log) (*DKGUpgraded, error) {
	event := new(DKGUpgraded)
	if err := _DKG.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
