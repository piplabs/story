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

// ICDRVault is an auto generated low-level Go binding around an user-defined struct.
type ICDRVault struct {
	Updatable          bool
	WriteConditionAddr common.Address
	ReadConditionAddr  common.Address
	WriteConditionData []byte
	ReadConditionData  []byte
	EncryptedData      []byte
	Label              []byte
}

// CDRMetaData contains all meta data concerning the CDR contract.
var CDRMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocate\",\"inputs\":[{\"name\":\"updatable\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"writeConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"readConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"writeConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"readConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"newVaultUuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"allocateFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"baseFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"baseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"writeFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"readFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allocateFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"read\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accessAuxData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"readFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAllocateFee\",\"inputs\":[{\"name\":\"newAllocateFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseFee\",\"inputs\":[{\"name\":\"newBaseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReadFee\",\"inputs\":[{\"name\":\"newReadFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWriteFee\",\"inputs\":[{\"name\":\"newWriteFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitEncryptedPartialDecryption\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"pid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ephemeralPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"uuid\",\"inputs\":[],\"outputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vaults\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"vault\",\"type\":\"tuple\",\"internalType\":\"structICDR.Vault\",\"components\":[{\"name\":\"updatable\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"writeConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"readConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"writeConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"readConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"write\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accessAuxData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"writeFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"EncryptedPartialDecryptionSubmitted\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"pid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ephemeralPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeCollected\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"feeType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumICDR.FeeType\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultAllocated\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"updatable\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"writeConditionAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"readConditionAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"writeConditionData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"readConditionData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultRead\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultWritten\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60a080604052346100cd57306080527ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c166100be57506001600160401b036002600160401b031982821601610079575b604051612b049081620000d2823960805181818161217901526122400152f35b6001600160401b031990911681179091556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f8080610059565b63f92ee8a960e01b8152600490fd5b5f80fdfe60a06040526004361015610011575f80fd5b5f3560e01c80630dae5528146124ed57806346860698146124ad5780634f1ef286146121f357806352d1902d146121525780635c975abb1461211157806361cfc52c14611ec55780636954b1cd14611e895780636ef25c3a14611e4d578063715018a614611d885780637741958414611d4c57806379ba509714611cc35780638da5cb5b14611c715780638e3850c714611c3157806390ab171e14611bf55780639cadce24146113e2578063ad3cb1cc14611372578063b98a6c1a14610f53578063bb07e85d14610f11578063e30c397814610ebf578063e60adef114610e7f578063eaf537b814610688578063f2fde38b146105bc578063f8b5ac35146104255763f92ad21914610121575f80fd5b346104215760a06003193601126104215761013a61252d565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c16159167ffffffffffffffff811680159081610419575b600114908161040f575b159081610406575b506103dc578260017fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000083161785556103a7575b506101ca6129d5565b6101d26129d5565b73ffffffffffffffffffffffffffffffffffffffff811615610377576101f7906128c8565b6101ff6129d5565b6102076129d5565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00556102336129d5565b61023b6129d5565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00815416905561028a6129d5565b6024357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf601556044357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf602556064357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf603556084357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf6045561032457005b7fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff81541690557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160018152a1005b60246040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081525f6004820152fd5b7fffffffffffffffffffffffffffffffffffffffffffffff00000000000000000016680100000000000000011783555f6101c1565b60046040517ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b9050155f61018e565b303b159150610186565b84915061017c565b5f80fd5b346104215760206003193601126104215761059660206105b861058060c06105a76104b661045161261b565b60405161045d81612550565b5f81525f888201525f6040820152606094818680809401528260808201528260a0820152015263ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b916040519687948580936104c982612550565b865460ff8116151583528773ffffffffffffffffffffffffffffffffffffffff8092818e87019160081c16815260c061054a60056105296003866001890154169760408c019889528a61051e6002830161275a565b9c019b8c520161275a565b9d60808d019e8f5260a061053f6004830161275a565b9d019c8d520161275a565b9c019b8c526040519e8f9e8f91818352511515910152511660408d01525116908a01525160e060808a015261010089019061265c565b935193601f1994858983030160a08a015261265c565b9051838783030160c088015261265c565b9151908483030160e085015261265c565b0390f35b34610421576020600319360112610421576105d561252d565b6105dd612833565b73ffffffffffffffffffffffffffffffffffffffff809116907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227005f80a3005b60806003193601126104215761069c61261b565b60243567ffffffffffffffff8111610421576106bc90369060040161262e565b919060443567ffffffffffffffff8111610421576106de90369060040161262e565b939060643567ffffffffffffffff81116104215761070090369060040161262e565b92909361070b61297b565b610713612873565b8615610dfb576107508663ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b916107e460056040519461076386612550565b73ffffffffffffffffffffffffffffffffffffffff815460ff81161515885260081c16602087015273ffffffffffffffffffffffffffffffffffffffff60018201541660408701526107b76002820161275a565b60608701526107c86003820161275a565b60808701526107d96004820161275a565b60a08701520161275a565b60c084015273ffffffffffffffffffffffffffffffffffffffff60208401511615610d785773ffffffffffffffffffffffffffffffffffffffff60208401511690813303610c5b575b50505060a081015151610bf1575b507f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60254803403610b9357805f8115610b8a575b5f808093818094f115610b7f57604051908152600160208201527f93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f740031860403392a260046108e78563ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b0167ffffffffffffffff8611610ae65761090b8661090583546126ba565b8361270b565b5f86601f8111600114610b1e5780610936925f91610b13575b505f198260011b9260031b1c19161790565b90555b60056109728563ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b019467ffffffffffffffff8311610ae6576109978361099188546126ba565b8861270b565b5f95601f8411600114610a535763ffffffff9594928492610a11926109f485610a1f987fef3894f4f6e35b318095909888a488a7b5835123cf07c51beccdb8cef00d7bdb9c5f91610a4857505f198260011b9260031b1c19161790565b90555b60405197889716875260606020880152606087019161269a565b91848303604086015261269a565b0390a160017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055005b90508801358d610924565b601f198416815f5260205f20905f5b818110610ace575092610a119263ffffffff989795927fef3894f4f6e35b318095909888a488a7b5835123cf07c51beccdb8cef00d7bdb9a88610a1f999710610ab5575b5050600185811b0190556109f7565b5f1960f88860031b161c19908801351690558a80610aa6565b878a013583556020998a019960019093019201610a62565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b905084013589610924565b601f191690825f528760205f20925f5b818110610b64575010610b4b575b5050600186811b019055610939565b5f1960f88960031b161c19908401351690558680610b3c565b8684013585556001909401936020938401938b935001610b2e565b6040513d5f823e3d90fd5b506108fc61086e565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4344523a20496e76616c69642066656520616d6f756e740000000000000000006044820152fd5b5115610bfd578561083b565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4344523a205661756c74206973206e6f7420757064617461626c6500000000006044820152fd5b60209163ffffffff89610cc56060880151610cb3604051988997889687967f5645dbbf00000000000000000000000000000000000000000000000000000000885216600487015260806024870152608486019161269a565b9060031984830301604485015261265c565b33606483015203915afa908115610b7f575f91610d49575b5015610ceb5786808061082d565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f4344523a20577269746520636f6e646974696f6e206e6f74206d6574000000006044820152fd5b610d6b915060203d602011610d71575b610d638183612588565b81019061281b565b87610cdd565b503d610d59565b60846040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f4344523a20577269746520636f6e646974696f6e2061646472657373206e6f7460448201527f20736574000000000000000000000000000000000000000000000000000000006064820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f4344523a20456e6372797074656420646174612063616e6e6f7420626520656d60448201527f70747900000000000000000000000000000000000000000000000000000000006064820152fd5b3461042157602060031936011261042157610e98612833565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60455005b34610421575f60031936011261042157602073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416604051908152f35b34610421575f60031936011261042157602063ffffffff7f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf6005416604051908152f35b60031960608136011261042157610f6861261b565b67ffffffffffffffff9160243583811161042157610f8a9036906004016125fd565b9260443590811161042157610fa390369060040161262e565b929091610fae61297b565b610fb6612873565b610fed8263ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b9260405191610ffb83612550565b84549660ff88161515845273ffffffffffffffffffffffffffffffffffffffff918260209960081c16898601528260018801541692604086019384526110436002890161275a565b60608701526110546003890161275a565b936080870194855260c061107d600561106f60048d0161275a565b9b60a08b019c8d520161275a565b97019687528851511561131457908a9392915116908782330361120f575b5050505050507f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60354908134036111b157815f81156111a8575b5f808093818094f115610b7f577f193a00215cd09d5e194005f27dbce2e464d9737adbb59c47310e176dce083f7a9561116261117f94608063ffffffff986111709660405190815260028d8201527f93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f740031860403392a25191519a604051998a99168952880152608087019061265c565b91858303604087015261269a565b8281036060840152339561265c565b0390a260017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055005b506108fc6110d4565b606487604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601760248201527f4344523a20496e76616c69642066656520616d6f756e740000000000000000006044820152fd5b61127063ffffffff926112619751604051988997889687967f8db3eb1700000000000000000000000000000000000000000000000000000000885216600487015260806024870152608486019061265c565b9184830301604485015261265c565b33606483015203915afa908115610b7f575f916112f7575b50156112995786868180808761109b565b606486604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601b60248201527f4344523a205265616420636f6e646974696f6e206e6f74206d657400000000006044820152fd5b61130e9150873d8911610d7157610d638183612588565b87611288565b60648b604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601e60248201527f4344523a205661756c7420686173206e6f206461746120746f207265616400006044820152fd5b34610421575f60031936011261042157604051604081019080821067ffffffffffffffff831117610ae6576105b891604052600581527f352e302e30000000000000000000000000000000000000000000000000000000602082015260405191829160208352602083019061265c565b60a0600319360112610421576004358015158103610421576024359073ffffffffffffffffffffffffffffffffffffffff82168203610421576044359073ffffffffffffffffffffffffffffffffffffffff821682036104215760643567ffffffffffffffff81116104215761145c90369060040161262e565b919060843567ffffffffffffffff81116104215761147e90369060040161262e565b919093611489612873565b73ffffffffffffffffffffffffffffffffffffffff871615801590611bd6575b15611b78577f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60454803403610b9357805f8115611b6f575b5f808093818094f115610b7f576040519081525f60208201527f93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f740031860403392a27f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf600549563ffffffff80881614611b425763ffffffff600181891601167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000008816177f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60055604051976115ae89612550565b8515158952602089019873ffffffffffffffffffffffffffffffffffffffff82168a526040810173ffffffffffffffffffffffffffffffffffffffff841681526115f93686886125c7565b906060830191825261160c36898c6125c7565b608084015260405161161d8161256c565b5f815260a08401526040516116318161256c565b5f815260c084015261167663ffffffff8c1663ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b809c7fffffffffffffffffffffff00000000000000000000000000000000000000000060ff74ffffffffffffffffffffffffffffffffffffffff00875115159454935160081b169316911617178c5573ffffffffffffffffffffffffffffffffffffffff60018d019151167fffffffffffffffffffffffff000000000000000000000000000000000000000082541617905551805167ffffffffffffffff8111610ae657611734818d600261172d818301546126ba565b910161270b565b60208c601f8311600114611ad857508190611763935f92611a63575b50505f198260011b9260031b1c19161790565b60028b01555b6080810151805167ffffffffffffffff8111610ae657611792818d600361172d818301546126ba565b60208c601f8311600114611a6e575081906117c0935f92611a635750505f198260011b9260031b1c19161790565b60038b01555b60a081015180519067ffffffffffffffff8211610ae6576117f0828d600461172d818301546126ba565b6020908c601f84116001146119f6575091806118239260c095945f926119eb5750505f198260011b9260031b1c19161790565b60048c01555b01519586519967ffffffffffffffff8b11610ae6578a8a9861185d60209d61185460058601546126ba565b6005860161270b565b8c90601f8311600114611932576119089561191699957ff1370099e56e061a64615edc07c1d17e36a20585cdc9b288bf0259a528365d0a9d999560056118d863ffffffff9f9e9b96978073ffffffffffffffffffffffffffffffffffffffff998a985f926119275750505f198260011b9260031b1c19161790565b9101555b6040519c8d9c168c5215158f8c01521660408a015216606088015260c0608088015260c087019161269a565b9184830360a086015261269a565b0390a163ffffffff60405191168152f35b015190505f80611750565b90600584015f528d5f20915f5b601f19851681106119d157509561191699957ff1370099e56e061a64615edc07c1d17e36a20585cdc9b288bf0259a528365d0a9d99956005600163ffffffff9f9e9b969773ffffffffffffffffffffffffffffffffffffffff9889976119089d83601f198116106119b9575b505050811b019101556118dc565b01515f1960f88460031b161c191690555f80806119ab565b8282015184558e9c50600190930192918f01918f0161193f565b015190508e80611750565b91906004601f19851693015f5260205f20925f5b818110611a4b575091600193918560c097969410611a33575b505050811b0160048c0155611829565b01515f1960f88460031b161c191690558d8080611a23565b92936020600181928786015181550195019301611a0a565b015190508d80611750565b91926003601f19851693015f5260205f20925f5b818110611ac05750908460019594939210611aa8575b505050811b0160038b01556117c6565b01515f1960f88460031b161c191690558c8080611a98565b92936020600181928786015181550195019301611a82565b91926002601f19851693015f5260205f20925f5b818110611b2a5750908460019594939210611b12575b505050811b0160028b0155611769565b01515f1960f88460031b161c191690558c8080611b02565b92936020600181928786015181550195019301611aec565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b506108fc6114e0565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f496e76616c696420636f6e646974696f6e2061646472657373000000000000006044820152fd5b5073ffffffffffffffffffffffffffffffffffffffff861615156114a9565b34610421575f6003193601126104215760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60354604051908152f35b3461042157602060031936011261042157611c4a612833565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60355005b34610421575f60031936011261042157602073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416604051908152f35b34610421575f600319360112610421573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00541603611d1c57611d1a336128c8565b005b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b34610421575f6003193601126104215760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60454604051908152f35b34610421575f60031936011261042157611da0612833565b5f73ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b34610421575f6003193601126104215760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60154604051908152f35b34610421575f6003193601126104215760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60254604051908152f35b6101208060031936011261042157611edb61261b565b6024359063ffffffff821682036104215767ffffffffffffffff60443581811161042157611f0d90369060040161262e565b909260643583811161042157611f2790369060040161262e565b95909460843585811161042157611f4290369060040161262e565b93909460a43587811161042157611f5d90369060040161262e565b919060c43589811161042157611f7790369060040161262e565b93909460e4359a63ffffffff8c168c03610421576101043590811161042157611fa490369060040161262e565b9a9099611faf612873565b7f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf601549e8f3403610b93578f5f80808093818115612108575b8290f115610b7f578f9e6120d49f916120726120b99c6120a79a63ffffffff9f9863ffffffff6120959a61208498604051908152600360208201527f93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f740031860403392a26040516080528161014093166080515216602060805101528060406080510152608051019161269a565b9160805183036060608051015261269a565b91608051830360808051015261269a565b91608051830360a0608051015261269a565b91608051830360c0608051015261269a565b931660e060805101526080518303610100608051015261269a565b9160805101527fe2d964436c9da9c84292bd1b67e737171e45dbb6ca28f34fcbfb5c9ce20c1dbc33916080519003608051a2005b506108fc611fe7565b34610421575f60031936011261042157602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b34610421575f6003193601126104215773ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001630036121c95760206040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b60046040517fe07c8dba000000000000000000000000000000000000000000000000000000008152fd5b60406003193601126104215761220761252d565b60243567ffffffffffffffff8111610421576122279036906004016125fd565b9073ffffffffffffffffffffffffffffffffffffffff807f00000000000000000000000000000000000000000000000000000000000000001680301490811561247f575b506121c957612278612833565b811690604051907f52d1902d0000000000000000000000000000000000000000000000000000000082526020918281600481875afa5f9181612450575b506122eb57602484604051907f4c9c8ce30000000000000000000000000000000000000000000000000000000082526004820152fd5b9284937f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc9081810361241f5750823b156123ee57817fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055604051907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a28351156123bb57505f808484611d1a96519101845af4903d156123b2573d612396816125ab565b906123a46040519283612588565b81525f81943d92013e612a2e565b60609250612a2e565b92505050346123c657005b807fb398979f0000000000000000000000000000000000000000000000000000000060049252fd5b602482604051907f4c9c8ce30000000000000000000000000000000000000000000000000000000082526004820152fd5b602490604051907faa1d49a40000000000000000000000000000000000000000000000000000000082526004820152fd5b9091508381813d8311612478575b6124688183612588565b81010312610421575190866122b5565b503d61245e565b9050817f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc541614158461226b565b34610421576020600319360112610421576124c6612833565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60155005b3461042157602060031936011261042157612506612833565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60255005b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361042157565b60e0810190811067ffffffffffffffff821117610ae657604052565b6020810190811067ffffffffffffffff821117610ae657604052565b90601f601f19910116810190811067ffffffffffffffff821117610ae657604052565b67ffffffffffffffff8111610ae657601f01601f191660200190565b9291926125d3826125ab565b916125e16040519384612588565b829481845281830111610421578281602093845f960137010152565b9080601f8301121561042157816020612618933591016125c7565b90565b6004359063ffffffff8216820361042157565b9181601f840112156104215782359167ffffffffffffffff8311610421576020838186019501011161042157565b91908251928382525f5b848110612686575050601f19601f845f6020809697860101520116010190565b602081830181015184830182015201612666565b601f8260209493601f1993818652868601375f8582860101520116010190565b90600182811c92168015612701575b60208310146126d457565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f16916126c9565b601f821161271857505050565b5f5260205f20906020601f840160051c83019310612750575b601f0160051c01905b818110612745575050565b5f815560010161273a565b9091508190612731565b9060405191825f825461276c816126ba565b908184526020946001916001811690815f146127da575060011461279c575b50505061279a92500383612588565b565b5f90815285812095935091905b8183106127c257505061279a93508201015f808061278b565b855488840185015294850194879450918301916127a9565b91505061279a9593507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201015f808061278b565b90816020910312610421575180151581036104215790565b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054163303611d1c57565b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300541661289e57565b60046040517fd93c0665000000000000000000000000000000000000000000000000000000008152fd5b7fffffffffffffffffffffffff0000000000000000000000000000000000000000907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008281541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549073ffffffffffffffffffffffffffffffffffffffff80931680948316179055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0060028154146129ab5760029055565b60046040517f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c1615612a0457565b60046040517fd7e6bcf8000000000000000000000000000000000000000000000000000000008152fd5b90612a6d5750805115612a4357805190602001fd5b60046040517f1425ea42000000000000000000000000000000000000000000000000000000008152fd5b81511580612ac5575b612a7e575090565b60249073ffffffffffffffffffffffffffffffffffffffff604051917f9996b315000000000000000000000000000000000000000000000000000000008352166004820152fd5b50803b15612a7656fea2646970667358221220b05c1a682202ef2d289ad1362aa4c1035ed2c46c135f075af5c5401c4695e6d564736f6c63430008170033",
}

// CDRABI is the input ABI used to generate the binding from.
// Deprecated: Use CDRMetaData.ABI instead.
var CDRABI = CDRMetaData.ABI

// CDRBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CDRMetaData.Bin instead.
var CDRBin = CDRMetaData.Bin

// DeployCDR deploys a new Ethereum contract, binding an instance of CDR to it.
func DeployCDR(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CDR, error) {
	parsed, err := CDRMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CDRBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CDR{CDRCaller: CDRCaller{contract: contract}, CDRTransactor: CDRTransactor{contract: contract}, CDRFilterer: CDRFilterer{contract: contract}}, nil
}

// CDR is an auto generated Go binding around an Ethereum contract.
type CDR struct {
	CDRCaller     // Read-only binding to the contract
	CDRTransactor // Write-only binding to the contract
	CDRFilterer   // Log filterer for contract events
}

// CDRCaller is an auto generated read-only Go binding around an Ethereum contract.
type CDRCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CDRTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CDRTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CDRFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CDRFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CDRSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CDRSession struct {
	Contract     *CDR              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CDRCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CDRCallerSession struct {
	Contract *CDRCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// CDRTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CDRTransactorSession struct {
	Contract     *CDRTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CDRRaw is an auto generated low-level Go binding around an Ethereum contract.
type CDRRaw struct {
	Contract *CDR // Generic contract binding to access the raw methods on
}

// CDRCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CDRCallerRaw struct {
	Contract *CDRCaller // Generic read-only contract binding to access the raw methods on
}

// CDRTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CDRTransactorRaw struct {
	Contract *CDRTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCDR creates a new instance of CDR, bound to a specific deployed contract.
func NewCDR(address common.Address, backend bind.ContractBackend) (*CDR, error) {
	contract, err := bindCDR(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CDR{CDRCaller: CDRCaller{contract: contract}, CDRTransactor: CDRTransactor{contract: contract}, CDRFilterer: CDRFilterer{contract: contract}}, nil
}

// NewCDRCaller creates a new read-only instance of CDR, bound to a specific deployed contract.
func NewCDRCaller(address common.Address, caller bind.ContractCaller) (*CDRCaller, error) {
	contract, err := bindCDR(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CDRCaller{contract: contract}, nil
}

// NewCDRTransactor creates a new write-only instance of CDR, bound to a specific deployed contract.
func NewCDRTransactor(address common.Address, transactor bind.ContractTransactor) (*CDRTransactor, error) {
	contract, err := bindCDR(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CDRTransactor{contract: contract}, nil
}

// NewCDRFilterer creates a new log filterer instance of CDR, bound to a specific deployed contract.
func NewCDRFilterer(address common.Address, filterer bind.ContractFilterer) (*CDRFilterer, error) {
	contract, err := bindCDR(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CDRFilterer{contract: contract}, nil
}

// bindCDR binds a generic wrapper to an already deployed contract.
func bindCDR(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CDRMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CDR *CDRRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CDR.Contract.CDRCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CDR *CDRRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CDR.Contract.CDRTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CDR *CDRRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CDR.Contract.CDRTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CDR *CDRCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CDR.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CDR *CDRTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CDR.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CDR *CDRTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CDR.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_CDR *CDRCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_CDR *CDRSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _CDR.Contract.UPGRADEINTERFACEVERSION(&_CDR.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_CDR *CDRCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _CDR.Contract.UPGRADEINTERFACEVERSION(&_CDR.CallOpts)
}

// AllocateFee is a free data retrieval call binding the contract method 0x77419584.
//
// Solidity: function allocateFee() view returns(uint256)
func (_CDR *CDRCaller) AllocateFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "allocateFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AllocateFee is a free data retrieval call binding the contract method 0x77419584.
//
// Solidity: function allocateFee() view returns(uint256)
func (_CDR *CDRSession) AllocateFee() (*big.Int, error) {
	return _CDR.Contract.AllocateFee(&_CDR.CallOpts)
}

// AllocateFee is a free data retrieval call binding the contract method 0x77419584.
//
// Solidity: function allocateFee() view returns(uint256)
func (_CDR *CDRCallerSession) AllocateFee() (*big.Int, error) {
	return _CDR.Contract.AllocateFee(&_CDR.CallOpts)
}

// BaseFee is a free data retrieval call binding the contract method 0x6ef25c3a.
//
// Solidity: function baseFee() view returns(uint256)
func (_CDR *CDRCaller) BaseFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "baseFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseFee is a free data retrieval call binding the contract method 0x6ef25c3a.
//
// Solidity: function baseFee() view returns(uint256)
func (_CDR *CDRSession) BaseFee() (*big.Int, error) {
	return _CDR.Contract.BaseFee(&_CDR.CallOpts)
}

// BaseFee is a free data retrieval call binding the contract method 0x6ef25c3a.
//
// Solidity: function baseFee() view returns(uint256)
func (_CDR *CDRCallerSession) BaseFee() (*big.Int, error) {
	return _CDR.Contract.BaseFee(&_CDR.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CDR *CDRCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CDR *CDRSession) Owner() (common.Address, error) {
	return _CDR.Contract.Owner(&_CDR.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CDR *CDRCallerSession) Owner() (common.Address, error) {
	return _CDR.Contract.Owner(&_CDR.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CDR *CDRCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CDR *CDRSession) Paused() (bool, error) {
	return _CDR.Contract.Paused(&_CDR.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CDR *CDRCallerSession) Paused() (bool, error) {
	return _CDR.Contract.Paused(&_CDR.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_CDR *CDRCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_CDR *CDRSession) PendingOwner() (common.Address, error) {
	return _CDR.Contract.PendingOwner(&_CDR.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_CDR *CDRCallerSession) PendingOwner() (common.Address, error) {
	return _CDR.Contract.PendingOwner(&_CDR.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_CDR *CDRCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_CDR *CDRSession) ProxiableUUID() ([32]byte, error) {
	return _CDR.Contract.ProxiableUUID(&_CDR.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_CDR *CDRCallerSession) ProxiableUUID() ([32]byte, error) {
	return _CDR.Contract.ProxiableUUID(&_CDR.CallOpts)
}

// ReadFee is a free data retrieval call binding the contract method 0x90ab171e.
//
// Solidity: function readFee() view returns(uint256)
func (_CDR *CDRCaller) ReadFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "readFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ReadFee is a free data retrieval call binding the contract method 0x90ab171e.
//
// Solidity: function readFee() view returns(uint256)
func (_CDR *CDRSession) ReadFee() (*big.Int, error) {
	return _CDR.Contract.ReadFee(&_CDR.CallOpts)
}

// ReadFee is a free data retrieval call binding the contract method 0x90ab171e.
//
// Solidity: function readFee() view returns(uint256)
func (_CDR *CDRCallerSession) ReadFee() (*big.Int, error) {
	return _CDR.Contract.ReadFee(&_CDR.CallOpts)
}

// Uuid is a free data retrieval call binding the contract method 0xbb07e85d.
//
// Solidity: function uuid() view returns(uint32 uuid)
func (_CDR *CDRCaller) Uuid(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "uuid")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// Uuid is a free data retrieval call binding the contract method 0xbb07e85d.
//
// Solidity: function uuid() view returns(uint32 uuid)
func (_CDR *CDRSession) Uuid() (uint32, error) {
	return _CDR.Contract.Uuid(&_CDR.CallOpts)
}

// Uuid is a free data retrieval call binding the contract method 0xbb07e85d.
//
// Solidity: function uuid() view returns(uint32 uuid)
func (_CDR *CDRCallerSession) Uuid() (uint32, error) {
	return _CDR.Contract.Uuid(&_CDR.CallOpts)
}

// Vaults is a free data retrieval call binding the contract method 0xf8b5ac35.
//
// Solidity: function vaults(uint32 uuid) view returns((bool,address,address,bytes,bytes,bytes,bytes) vault)
func (_CDR *CDRCaller) Vaults(opts *bind.CallOpts, uuid uint32) (ICDRVault, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "vaults", uuid)

	if err != nil {
		return *new(ICDRVault), err
	}

	out0 := *abi.ConvertType(out[0], new(ICDRVault)).(*ICDRVault)

	return out0, err

}

// Vaults is a free data retrieval call binding the contract method 0xf8b5ac35.
//
// Solidity: function vaults(uint32 uuid) view returns((bool,address,address,bytes,bytes,bytes,bytes) vault)
func (_CDR *CDRSession) Vaults(uuid uint32) (ICDRVault, error) {
	return _CDR.Contract.Vaults(&_CDR.CallOpts, uuid)
}

// Vaults is a free data retrieval call binding the contract method 0xf8b5ac35.
//
// Solidity: function vaults(uint32 uuid) view returns((bool,address,address,bytes,bytes,bytes,bytes) vault)
func (_CDR *CDRCallerSession) Vaults(uuid uint32) (ICDRVault, error) {
	return _CDR.Contract.Vaults(&_CDR.CallOpts, uuid)
}

// WriteFee is a free data retrieval call binding the contract method 0x6954b1cd.
//
// Solidity: function writeFee() view returns(uint256)
func (_CDR *CDRCaller) WriteFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "writeFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WriteFee is a free data retrieval call binding the contract method 0x6954b1cd.
//
// Solidity: function writeFee() view returns(uint256)
func (_CDR *CDRSession) WriteFee() (*big.Int, error) {
	return _CDR.Contract.WriteFee(&_CDR.CallOpts)
}

// WriteFee is a free data retrieval call binding the contract method 0x6954b1cd.
//
// Solidity: function writeFee() view returns(uint256)
func (_CDR *CDRCallerSession) WriteFee() (*big.Int, error) {
	return _CDR.Contract.WriteFee(&_CDR.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_CDR *CDRTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_CDR *CDRSession) AcceptOwnership() (*types.Transaction, error) {
	return _CDR.Contract.AcceptOwnership(&_CDR.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_CDR *CDRTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _CDR.Contract.AcceptOwnership(&_CDR.TransactOpts)
}

// Allocate is a paid mutator transaction binding the contract method 0x9cadce24.
//
// Solidity: function allocate(bool updatable, address writeConditionAddr, address readConditionAddr, bytes writeConditionData, bytes readConditionData) payable returns(uint32 newVaultUuid)
func (_CDR *CDRTransactor) Allocate(opts *bind.TransactOpts, updatable bool, writeConditionAddr common.Address, readConditionAddr common.Address, writeConditionData []byte, readConditionData []byte) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "allocate", updatable, writeConditionAddr, readConditionAddr, writeConditionData, readConditionData)
}

// Allocate is a paid mutator transaction binding the contract method 0x9cadce24.
//
// Solidity: function allocate(bool updatable, address writeConditionAddr, address readConditionAddr, bytes writeConditionData, bytes readConditionData) payable returns(uint32 newVaultUuid)
func (_CDR *CDRSession) Allocate(updatable bool, writeConditionAddr common.Address, readConditionAddr common.Address, writeConditionData []byte, readConditionData []byte) (*types.Transaction, error) {
	return _CDR.Contract.Allocate(&_CDR.TransactOpts, updatable, writeConditionAddr, readConditionAddr, writeConditionData, readConditionData)
}

// Allocate is a paid mutator transaction binding the contract method 0x9cadce24.
//
// Solidity: function allocate(bool updatable, address writeConditionAddr, address readConditionAddr, bytes writeConditionData, bytes readConditionData) payable returns(uint32 newVaultUuid)
func (_CDR *CDRTransactorSession) Allocate(updatable bool, writeConditionAddr common.Address, readConditionAddr common.Address, writeConditionData []byte, readConditionData []byte) (*types.Transaction, error) {
	return _CDR.Contract.Allocate(&_CDR.TransactOpts, updatable, writeConditionAddr, readConditionAddr, writeConditionData, readConditionData)
}

// Initialize is a paid mutator transaction binding the contract method 0xf92ad219.
//
// Solidity: function initialize(address owner, uint256 baseFee, uint256 writeFee, uint256 readFee, uint256 allocateFee) returns()
func (_CDR *CDRTransactor) Initialize(opts *bind.TransactOpts, owner common.Address, baseFee *big.Int, writeFee *big.Int, readFee *big.Int, allocateFee *big.Int) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "initialize", owner, baseFee, writeFee, readFee, allocateFee)
}

// Initialize is a paid mutator transaction binding the contract method 0xf92ad219.
//
// Solidity: function initialize(address owner, uint256 baseFee, uint256 writeFee, uint256 readFee, uint256 allocateFee) returns()
func (_CDR *CDRSession) Initialize(owner common.Address, baseFee *big.Int, writeFee *big.Int, readFee *big.Int, allocateFee *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.Initialize(&_CDR.TransactOpts, owner, baseFee, writeFee, readFee, allocateFee)
}

// Initialize is a paid mutator transaction binding the contract method 0xf92ad219.
//
// Solidity: function initialize(address owner, uint256 baseFee, uint256 writeFee, uint256 readFee, uint256 allocateFee) returns()
func (_CDR *CDRTransactorSession) Initialize(owner common.Address, baseFee *big.Int, writeFee *big.Int, readFee *big.Int, allocateFee *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.Initialize(&_CDR.TransactOpts, owner, baseFee, writeFee, readFee, allocateFee)
}

// Read is a paid mutator transaction binding the contract method 0xb98a6c1a.
//
// Solidity: function read(uint32 uuid, bytes accessAuxData, bytes requesterPubKey) payable returns()
func (_CDR *CDRTransactor) Read(opts *bind.TransactOpts, uuid uint32, accessAuxData []byte, requesterPubKey []byte) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "read", uuid, accessAuxData, requesterPubKey)
}

// Read is a paid mutator transaction binding the contract method 0xb98a6c1a.
//
// Solidity: function read(uint32 uuid, bytes accessAuxData, bytes requesterPubKey) payable returns()
func (_CDR *CDRSession) Read(uuid uint32, accessAuxData []byte, requesterPubKey []byte) (*types.Transaction, error) {
	return _CDR.Contract.Read(&_CDR.TransactOpts, uuid, accessAuxData, requesterPubKey)
}

// Read is a paid mutator transaction binding the contract method 0xb98a6c1a.
//
// Solidity: function read(uint32 uuid, bytes accessAuxData, bytes requesterPubKey) payable returns()
func (_CDR *CDRTransactorSession) Read(uuid uint32, accessAuxData []byte, requesterPubKey []byte) (*types.Transaction, error) {
	return _CDR.Contract.Read(&_CDR.TransactOpts, uuid, accessAuxData, requesterPubKey)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CDR *CDRTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CDR *CDRSession) RenounceOwnership() (*types.Transaction, error) {
	return _CDR.Contract.RenounceOwnership(&_CDR.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CDR *CDRTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CDR.Contract.RenounceOwnership(&_CDR.TransactOpts)
}

// SetAllocateFee is a paid mutator transaction binding the contract method 0xe60adef1.
//
// Solidity: function setAllocateFee(uint256 newAllocateFee) returns()
func (_CDR *CDRTransactor) SetAllocateFee(opts *bind.TransactOpts, newAllocateFee *big.Int) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "setAllocateFee", newAllocateFee)
}

// SetAllocateFee is a paid mutator transaction binding the contract method 0xe60adef1.
//
// Solidity: function setAllocateFee(uint256 newAllocateFee) returns()
func (_CDR *CDRSession) SetAllocateFee(newAllocateFee *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.SetAllocateFee(&_CDR.TransactOpts, newAllocateFee)
}

// SetAllocateFee is a paid mutator transaction binding the contract method 0xe60adef1.
//
// Solidity: function setAllocateFee(uint256 newAllocateFee) returns()
func (_CDR *CDRTransactorSession) SetAllocateFee(newAllocateFee *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.SetAllocateFee(&_CDR.TransactOpts, newAllocateFee)
}

// SetBaseFee is a paid mutator transaction binding the contract method 0x46860698.
//
// Solidity: function setBaseFee(uint256 newBaseFee) returns()
func (_CDR *CDRTransactor) SetBaseFee(opts *bind.TransactOpts, newBaseFee *big.Int) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "setBaseFee", newBaseFee)
}

// SetBaseFee is a paid mutator transaction binding the contract method 0x46860698.
//
// Solidity: function setBaseFee(uint256 newBaseFee) returns()
func (_CDR *CDRSession) SetBaseFee(newBaseFee *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.SetBaseFee(&_CDR.TransactOpts, newBaseFee)
}

// SetBaseFee is a paid mutator transaction binding the contract method 0x46860698.
//
// Solidity: function setBaseFee(uint256 newBaseFee) returns()
func (_CDR *CDRTransactorSession) SetBaseFee(newBaseFee *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.SetBaseFee(&_CDR.TransactOpts, newBaseFee)
}

// SetReadFee is a paid mutator transaction binding the contract method 0x8e3850c7.
//
// Solidity: function setReadFee(uint256 newReadFee) returns()
func (_CDR *CDRTransactor) SetReadFee(opts *bind.TransactOpts, newReadFee *big.Int) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "setReadFee", newReadFee)
}

// SetReadFee is a paid mutator transaction binding the contract method 0x8e3850c7.
//
// Solidity: function setReadFee(uint256 newReadFee) returns()
func (_CDR *CDRSession) SetReadFee(newReadFee *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.SetReadFee(&_CDR.TransactOpts, newReadFee)
}

// SetReadFee is a paid mutator transaction binding the contract method 0x8e3850c7.
//
// Solidity: function setReadFee(uint256 newReadFee) returns()
func (_CDR *CDRTransactorSession) SetReadFee(newReadFee *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.SetReadFee(&_CDR.TransactOpts, newReadFee)
}

// SetWriteFee is a paid mutator transaction binding the contract method 0x0dae5528.
//
// Solidity: function setWriteFee(uint256 newWriteFee) returns()
func (_CDR *CDRTransactor) SetWriteFee(opts *bind.TransactOpts, newWriteFee *big.Int) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "setWriteFee", newWriteFee)
}

// SetWriteFee is a paid mutator transaction binding the contract method 0x0dae5528.
//
// Solidity: function setWriteFee(uint256 newWriteFee) returns()
func (_CDR *CDRSession) SetWriteFee(newWriteFee *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.SetWriteFee(&_CDR.TransactOpts, newWriteFee)
}

// SetWriteFee is a paid mutator transaction binding the contract method 0x0dae5528.
//
// Solidity: function setWriteFee(uint256 newWriteFee) returns()
func (_CDR *CDRTransactorSession) SetWriteFee(newWriteFee *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.SetWriteFee(&_CDR.TransactOpts, newWriteFee)
}

// SubmitEncryptedPartialDecryption is a paid mutator transaction binding the contract method 0x61cfc52c.
//
// Solidity: function submitEncryptedPartialDecryption(uint32 round, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes requesterPubKey, bytes ciphertext, uint32 uuid, bytes signature) payable returns()
func (_CDR *CDRTransactor) SubmitEncryptedPartialDecryption(opts *bind.TransactOpts, round uint32, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, requesterPubKey []byte, ciphertext []byte, uuid uint32, signature []byte) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "submitEncryptedPartialDecryption", round, pid, encryptedPartial, ephemeralPubKey, pubShare, requesterPubKey, ciphertext, uuid, signature)
}

// SubmitEncryptedPartialDecryption is a paid mutator transaction binding the contract method 0x61cfc52c.
//
// Solidity: function submitEncryptedPartialDecryption(uint32 round, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes requesterPubKey, bytes ciphertext, uint32 uuid, bytes signature) payable returns()
func (_CDR *CDRSession) SubmitEncryptedPartialDecryption(round uint32, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, requesterPubKey []byte, ciphertext []byte, uuid uint32, signature []byte) (*types.Transaction, error) {
	return _CDR.Contract.SubmitEncryptedPartialDecryption(&_CDR.TransactOpts, round, pid, encryptedPartial, ephemeralPubKey, pubShare, requesterPubKey, ciphertext, uuid, signature)
}

// SubmitEncryptedPartialDecryption is a paid mutator transaction binding the contract method 0x61cfc52c.
//
// Solidity: function submitEncryptedPartialDecryption(uint32 round, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes requesterPubKey, bytes ciphertext, uint32 uuid, bytes signature) payable returns()
func (_CDR *CDRTransactorSession) SubmitEncryptedPartialDecryption(round uint32, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, requesterPubKey []byte, ciphertext []byte, uuid uint32, signature []byte) (*types.Transaction, error) {
	return _CDR.Contract.SubmitEncryptedPartialDecryption(&_CDR.TransactOpts, round, pid, encryptedPartial, ephemeralPubKey, pubShare, requesterPubKey, ciphertext, uuid, signature)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CDR *CDRTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CDR *CDRSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CDR.Contract.TransferOwnership(&_CDR.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CDR *CDRTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CDR.Contract.TransferOwnership(&_CDR.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_CDR *CDRTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_CDR *CDRSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _CDR.Contract.UpgradeToAndCall(&_CDR.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_CDR *CDRTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _CDR.Contract.UpgradeToAndCall(&_CDR.TransactOpts, newImplementation, data)
}

// Write is a paid mutator transaction binding the contract method 0xeaf537b8.
//
// Solidity: function write(uint32 uuid, bytes accessAuxData, bytes encryptedData, bytes label) payable returns()
func (_CDR *CDRTransactor) Write(opts *bind.TransactOpts, uuid uint32, accessAuxData []byte, encryptedData []byte, label []byte) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "write", uuid, accessAuxData, encryptedData, label)
}

// Write is a paid mutator transaction binding the contract method 0xeaf537b8.
//
// Solidity: function write(uint32 uuid, bytes accessAuxData, bytes encryptedData, bytes label) payable returns()
func (_CDR *CDRSession) Write(uuid uint32, accessAuxData []byte, encryptedData []byte, label []byte) (*types.Transaction, error) {
	return _CDR.Contract.Write(&_CDR.TransactOpts, uuid, accessAuxData, encryptedData, label)
}

// Write is a paid mutator transaction binding the contract method 0xeaf537b8.
//
// Solidity: function write(uint32 uuid, bytes accessAuxData, bytes encryptedData, bytes label) payable returns()
func (_CDR *CDRTransactorSession) Write(uuid uint32, accessAuxData []byte, encryptedData []byte, label []byte) (*types.Transaction, error) {
	return _CDR.Contract.Write(&_CDR.TransactOpts, uuid, accessAuxData, encryptedData, label)
}

// CDREncryptedPartialDecryptionSubmittedIterator is returned from FilterEncryptedPartialDecryptionSubmitted and is used to iterate over the raw logs and unpacked data for EncryptedPartialDecryptionSubmitted events raised by the CDR contract.
type CDREncryptedPartialDecryptionSubmittedIterator struct {
	Event *CDREncryptedPartialDecryptionSubmitted // Event containing the contract specifics and raw log

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
func (it *CDREncryptedPartialDecryptionSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CDREncryptedPartialDecryptionSubmitted)
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
		it.Event = new(CDREncryptedPartialDecryptionSubmitted)
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
func (it *CDREncryptedPartialDecryptionSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CDREncryptedPartialDecryptionSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CDREncryptedPartialDecryptionSubmitted represents a EncryptedPartialDecryptionSubmitted event raised by the CDR contract.
type CDREncryptedPartialDecryptionSubmitted struct {
	Validator        common.Address
	Round            uint32
	Pid              uint32
	EncryptedPartial []byte
	EphemeralPubKey  []byte
	PubShare         []byte
	RequesterPubKey  []byte
	Ciphertext       []byte
	Uuid             uint32
	Signature        []byte
	Fee              *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterEncryptedPartialDecryptionSubmitted is a free log retrieval operation binding the contract event 0xe2d964436c9da9c84292bd1b67e737171e45dbb6ca28f34fcbfb5c9ce20c1dbc.
//
// Solidity: event EncryptedPartialDecryptionSubmitted(address indexed validator, uint32 round, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes requesterPubKey, bytes ciphertext, uint32 uuid, bytes signature, uint256 fee)
func (_CDR *CDRFilterer) FilterEncryptedPartialDecryptionSubmitted(opts *bind.FilterOpts, validator []common.Address) (*CDREncryptedPartialDecryptionSubmittedIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _CDR.contract.FilterLogs(opts, "EncryptedPartialDecryptionSubmitted", validatorRule)
	if err != nil {
		return nil, err
	}
	return &CDREncryptedPartialDecryptionSubmittedIterator{contract: _CDR.contract, event: "EncryptedPartialDecryptionSubmitted", logs: logs, sub: sub}, nil
}

// WatchEncryptedPartialDecryptionSubmitted is a free log subscription operation binding the contract event 0xe2d964436c9da9c84292bd1b67e737171e45dbb6ca28f34fcbfb5c9ce20c1dbc.
//
// Solidity: event EncryptedPartialDecryptionSubmitted(address indexed validator, uint32 round, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes requesterPubKey, bytes ciphertext, uint32 uuid, bytes signature, uint256 fee)
func (_CDR *CDRFilterer) WatchEncryptedPartialDecryptionSubmitted(opts *bind.WatchOpts, sink chan<- *CDREncryptedPartialDecryptionSubmitted, validator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _CDR.contract.WatchLogs(opts, "EncryptedPartialDecryptionSubmitted", validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CDREncryptedPartialDecryptionSubmitted)
				if err := _CDR.contract.UnpackLog(event, "EncryptedPartialDecryptionSubmitted", log); err != nil {
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

// ParseEncryptedPartialDecryptionSubmitted is a log parse operation binding the contract event 0xe2d964436c9da9c84292bd1b67e737171e45dbb6ca28f34fcbfb5c9ce20c1dbc.
//
// Solidity: event EncryptedPartialDecryptionSubmitted(address indexed validator, uint32 round, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes requesterPubKey, bytes ciphertext, uint32 uuid, bytes signature, uint256 fee)
func (_CDR *CDRFilterer) ParseEncryptedPartialDecryptionSubmitted(log types.Log) (*CDREncryptedPartialDecryptionSubmitted, error) {
	event := new(CDREncryptedPartialDecryptionSubmitted)
	if err := _CDR.contract.UnpackLog(event, "EncryptedPartialDecryptionSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CDRFeeCollectedIterator is returned from FilterFeeCollected and is used to iterate over the raw logs and unpacked data for FeeCollected events raised by the CDR contract.
type CDRFeeCollectedIterator struct {
	Event *CDRFeeCollected // Event containing the contract specifics and raw log

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
func (it *CDRFeeCollectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CDRFeeCollected)
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
		it.Event = new(CDRFeeCollected)
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
func (it *CDRFeeCollectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CDRFeeCollectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CDRFeeCollected represents a FeeCollected event raised by the CDR contract.
type CDRFeeCollected struct {
	Payer   common.Address
	Amount  *big.Int
	FeeType uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterFeeCollected is a free log retrieval operation binding the contract event 0x93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f7400318.
//
// Solidity: event FeeCollected(address indexed payer, uint256 amount, uint8 feeType)
func (_CDR *CDRFilterer) FilterFeeCollected(opts *bind.FilterOpts, payer []common.Address) (*CDRFeeCollectedIterator, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _CDR.contract.FilterLogs(opts, "FeeCollected", payerRule)
	if err != nil {
		return nil, err
	}
	return &CDRFeeCollectedIterator{contract: _CDR.contract, event: "FeeCollected", logs: logs, sub: sub}, nil
}

// WatchFeeCollected is a free log subscription operation binding the contract event 0x93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f7400318.
//
// Solidity: event FeeCollected(address indexed payer, uint256 amount, uint8 feeType)
func (_CDR *CDRFilterer) WatchFeeCollected(opts *bind.WatchOpts, sink chan<- *CDRFeeCollected, payer []common.Address) (event.Subscription, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _CDR.contract.WatchLogs(opts, "FeeCollected", payerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CDRFeeCollected)
				if err := _CDR.contract.UnpackLog(event, "FeeCollected", log); err != nil {
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

// ParseFeeCollected is a log parse operation binding the contract event 0x93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f7400318.
//
// Solidity: event FeeCollected(address indexed payer, uint256 amount, uint8 feeType)
func (_CDR *CDRFilterer) ParseFeeCollected(log types.Log) (*CDRFeeCollected, error) {
	event := new(CDRFeeCollected)
	if err := _CDR.contract.UnpackLog(event, "FeeCollected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CDRInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the CDR contract.
type CDRInitializedIterator struct {
	Event *CDRInitialized // Event containing the contract specifics and raw log

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
func (it *CDRInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CDRInitialized)
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
		it.Event = new(CDRInitialized)
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
func (it *CDRInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CDRInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CDRInitialized represents a Initialized event raised by the CDR contract.
type CDRInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CDR *CDRFilterer) FilterInitialized(opts *bind.FilterOpts) (*CDRInitializedIterator, error) {

	logs, sub, err := _CDR.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &CDRInitializedIterator{contract: _CDR.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CDR *CDRFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *CDRInitialized) (event.Subscription, error) {

	logs, sub, err := _CDR.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CDRInitialized)
				if err := _CDR.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_CDR *CDRFilterer) ParseInitialized(log types.Log) (*CDRInitialized, error) {
	event := new(CDRInitialized)
	if err := _CDR.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CDROwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the CDR contract.
type CDROwnershipTransferStartedIterator struct {
	Event *CDROwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *CDROwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CDROwnershipTransferStarted)
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
		it.Event = new(CDROwnershipTransferStarted)
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
func (it *CDROwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CDROwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CDROwnershipTransferStarted represents a OwnershipTransferStarted event raised by the CDR contract.
type CDROwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_CDR *CDRFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CDROwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CDR.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CDROwnershipTransferStartedIterator{contract: _CDR.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_CDR *CDRFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *CDROwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CDR.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CDROwnershipTransferStarted)
				if err := _CDR.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_CDR *CDRFilterer) ParseOwnershipTransferStarted(log types.Log) (*CDROwnershipTransferStarted, error) {
	event := new(CDROwnershipTransferStarted)
	if err := _CDR.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CDROwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CDR contract.
type CDROwnershipTransferredIterator struct {
	Event *CDROwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CDROwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CDROwnershipTransferred)
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
		it.Event = new(CDROwnershipTransferred)
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
func (it *CDROwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CDROwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CDROwnershipTransferred represents a OwnershipTransferred event raised by the CDR contract.
type CDROwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CDR *CDRFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CDROwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CDR.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CDROwnershipTransferredIterator{contract: _CDR.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CDR *CDRFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CDROwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CDR.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CDROwnershipTransferred)
				if err := _CDR.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_CDR *CDRFilterer) ParseOwnershipTransferred(log types.Log) (*CDROwnershipTransferred, error) {
	event := new(CDROwnershipTransferred)
	if err := _CDR.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CDRPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the CDR contract.
type CDRPausedIterator struct {
	Event *CDRPaused // Event containing the contract specifics and raw log

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
func (it *CDRPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CDRPaused)
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
		it.Event = new(CDRPaused)
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
func (it *CDRPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CDRPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CDRPaused represents a Paused event raised by the CDR contract.
type CDRPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_CDR *CDRFilterer) FilterPaused(opts *bind.FilterOpts) (*CDRPausedIterator, error) {

	logs, sub, err := _CDR.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &CDRPausedIterator{contract: _CDR.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_CDR *CDRFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *CDRPaused) (event.Subscription, error) {

	logs, sub, err := _CDR.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CDRPaused)
				if err := _CDR.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_CDR *CDRFilterer) ParsePaused(log types.Log) (*CDRPaused, error) {
	event := new(CDRPaused)
	if err := _CDR.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CDRUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the CDR contract.
type CDRUnpausedIterator struct {
	Event *CDRUnpaused // Event containing the contract specifics and raw log

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
func (it *CDRUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CDRUnpaused)
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
		it.Event = new(CDRUnpaused)
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
func (it *CDRUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CDRUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CDRUnpaused represents a Unpaused event raised by the CDR contract.
type CDRUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_CDR *CDRFilterer) FilterUnpaused(opts *bind.FilterOpts) (*CDRUnpausedIterator, error) {

	logs, sub, err := _CDR.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &CDRUnpausedIterator{contract: _CDR.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_CDR *CDRFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *CDRUnpaused) (event.Subscription, error) {

	logs, sub, err := _CDR.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CDRUnpaused)
				if err := _CDR.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_CDR *CDRFilterer) ParseUnpaused(log types.Log) (*CDRUnpaused, error) {
	event := new(CDRUnpaused)
	if err := _CDR.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CDRUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the CDR contract.
type CDRUpgradedIterator struct {
	Event *CDRUpgraded // Event containing the contract specifics and raw log

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
func (it *CDRUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CDRUpgraded)
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
		it.Event = new(CDRUpgraded)
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
func (it *CDRUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CDRUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CDRUpgraded represents a Upgraded event raised by the CDR contract.
type CDRUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_CDR *CDRFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*CDRUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _CDR.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &CDRUpgradedIterator{contract: _CDR.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_CDR *CDRFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *CDRUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _CDR.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CDRUpgraded)
				if err := _CDR.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_CDR *CDRFilterer) ParseUpgraded(log types.Log) (*CDRUpgraded, error) {
	event := new(CDRUpgraded)
	if err := _CDR.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CDRVaultAllocatedIterator is returned from FilterVaultAllocated and is used to iterate over the raw logs and unpacked data for VaultAllocated events raised by the CDR contract.
type CDRVaultAllocatedIterator struct {
	Event *CDRVaultAllocated // Event containing the contract specifics and raw log

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
func (it *CDRVaultAllocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CDRVaultAllocated)
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
		it.Event = new(CDRVaultAllocated)
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
func (it *CDRVaultAllocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CDRVaultAllocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CDRVaultAllocated represents a VaultAllocated event raised by the CDR contract.
type CDRVaultAllocated struct {
	Uuid               uint32
	Updatable          bool
	WriteConditionAddr common.Address
	ReadConditionAddr  common.Address
	WriteConditionData []byte
	ReadConditionData  []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterVaultAllocated is a free log retrieval operation binding the contract event 0xf1370099e56e061a64615edc07c1d17e36a20585cdc9b288bf0259a528365d0a.
//
// Solidity: event VaultAllocated(uint32 uuid, bool updatable, address writeConditionAddr, address readConditionAddr, bytes writeConditionData, bytes readConditionData)
func (_CDR *CDRFilterer) FilterVaultAllocated(opts *bind.FilterOpts) (*CDRVaultAllocatedIterator, error) {

	logs, sub, err := _CDR.contract.FilterLogs(opts, "VaultAllocated")
	if err != nil {
		return nil, err
	}
	return &CDRVaultAllocatedIterator{contract: _CDR.contract, event: "VaultAllocated", logs: logs, sub: sub}, nil
}

// WatchVaultAllocated is a free log subscription operation binding the contract event 0xf1370099e56e061a64615edc07c1d17e36a20585cdc9b288bf0259a528365d0a.
//
// Solidity: event VaultAllocated(uint32 uuid, bool updatable, address writeConditionAddr, address readConditionAddr, bytes writeConditionData, bytes readConditionData)
func (_CDR *CDRFilterer) WatchVaultAllocated(opts *bind.WatchOpts, sink chan<- *CDRVaultAllocated) (event.Subscription, error) {

	logs, sub, err := _CDR.contract.WatchLogs(opts, "VaultAllocated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CDRVaultAllocated)
				if err := _CDR.contract.UnpackLog(event, "VaultAllocated", log); err != nil {
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

// ParseVaultAllocated is a log parse operation binding the contract event 0xf1370099e56e061a64615edc07c1d17e36a20585cdc9b288bf0259a528365d0a.
//
// Solidity: event VaultAllocated(uint32 uuid, bool updatable, address writeConditionAddr, address readConditionAddr, bytes writeConditionData, bytes readConditionData)
func (_CDR *CDRFilterer) ParseVaultAllocated(log types.Log) (*CDRVaultAllocated, error) {
	event := new(CDRVaultAllocated)
	if err := _CDR.contract.UnpackLog(event, "VaultAllocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CDRVaultReadIterator is returned from FilterVaultRead and is used to iterate over the raw logs and unpacked data for VaultRead events raised by the CDR contract.
type CDRVaultReadIterator struct {
	Event *CDRVaultRead // Event containing the contract specifics and raw log

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
func (it *CDRVaultReadIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CDRVaultRead)
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
		it.Event = new(CDRVaultRead)
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
func (it *CDRVaultReadIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CDRVaultReadIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CDRVaultRead represents a VaultRead event raised by the CDR contract.
type CDRVaultRead struct {
	Uuid            uint32
	Requester       common.Address
	Ciphertext      []byte
	RequesterPubKey []byte
	Label           []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterVaultRead is a free log retrieval operation binding the contract event 0x193a00215cd09d5e194005f27dbce2e464d9737adbb59c47310e176dce083f7a.
//
// Solidity: event VaultRead(uint32 uuid, address indexed requester, bytes ciphertext, bytes requesterPubKey, bytes label)
func (_CDR *CDRFilterer) FilterVaultRead(opts *bind.FilterOpts, requester []common.Address) (*CDRVaultReadIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _CDR.contract.FilterLogs(opts, "VaultRead", requesterRule)
	if err != nil {
		return nil, err
	}
	return &CDRVaultReadIterator{contract: _CDR.contract, event: "VaultRead", logs: logs, sub: sub}, nil
}

// WatchVaultRead is a free log subscription operation binding the contract event 0x193a00215cd09d5e194005f27dbce2e464d9737adbb59c47310e176dce083f7a.
//
// Solidity: event VaultRead(uint32 uuid, address indexed requester, bytes ciphertext, bytes requesterPubKey, bytes label)
func (_CDR *CDRFilterer) WatchVaultRead(opts *bind.WatchOpts, sink chan<- *CDRVaultRead, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _CDR.contract.WatchLogs(opts, "VaultRead", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CDRVaultRead)
				if err := _CDR.contract.UnpackLog(event, "VaultRead", log); err != nil {
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

// ParseVaultRead is a log parse operation binding the contract event 0x193a00215cd09d5e194005f27dbce2e464d9737adbb59c47310e176dce083f7a.
//
// Solidity: event VaultRead(uint32 uuid, address indexed requester, bytes ciphertext, bytes requesterPubKey, bytes label)
func (_CDR *CDRFilterer) ParseVaultRead(log types.Log) (*CDRVaultRead, error) {
	event := new(CDRVaultRead)
	if err := _CDR.contract.UnpackLog(event, "VaultRead", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CDRVaultWrittenIterator is returned from FilterVaultWritten and is used to iterate over the raw logs and unpacked data for VaultWritten events raised by the CDR contract.
type CDRVaultWrittenIterator struct {
	Event *CDRVaultWritten // Event containing the contract specifics and raw log

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
func (it *CDRVaultWrittenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CDRVaultWritten)
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
		it.Event = new(CDRVaultWritten)
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
func (it *CDRVaultWrittenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CDRVaultWrittenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CDRVaultWritten represents a VaultWritten event raised by the CDR contract.
type CDRVaultWritten struct {
	Uuid          uint32
	EncryptedData []byte
	Label         []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterVaultWritten is a free log retrieval operation binding the contract event 0xef3894f4f6e35b318095909888a488a7b5835123cf07c51beccdb8cef00d7bdb.
//
// Solidity: event VaultWritten(uint32 uuid, bytes encryptedData, bytes label)
func (_CDR *CDRFilterer) FilterVaultWritten(opts *bind.FilterOpts) (*CDRVaultWrittenIterator, error) {

	logs, sub, err := _CDR.contract.FilterLogs(opts, "VaultWritten")
	if err != nil {
		return nil, err
	}
	return &CDRVaultWrittenIterator{contract: _CDR.contract, event: "VaultWritten", logs: logs, sub: sub}, nil
}

// WatchVaultWritten is a free log subscription operation binding the contract event 0xef3894f4f6e35b318095909888a488a7b5835123cf07c51beccdb8cef00d7bdb.
//
// Solidity: event VaultWritten(uint32 uuid, bytes encryptedData, bytes label)
func (_CDR *CDRFilterer) WatchVaultWritten(opts *bind.WatchOpts, sink chan<- *CDRVaultWritten) (event.Subscription, error) {

	logs, sub, err := _CDR.contract.WatchLogs(opts, "VaultWritten")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CDRVaultWritten)
				if err := _CDR.contract.UnpackLog(event, "VaultWritten", log); err != nil {
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

// ParseVaultWritten is a log parse operation binding the contract event 0xef3894f4f6e35b318095909888a488a7b5835123cf07c51beccdb8cef00d7bdb.
//
// Solidity: event VaultWritten(uint32 uuid, bytes encryptedData, bytes label)
func (_CDR *CDRFilterer) ParseVaultWritten(log types.Log) (*CDRVaultWritten, error) {
	event := new(CDRVaultWritten)
	if err := _CDR.contract.UnpackLog(event, "VaultWritten", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
