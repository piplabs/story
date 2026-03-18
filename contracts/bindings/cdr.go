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
}

// CDRMetaData contains all meta data concerning the CDR contract.
var CDRMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocate\",\"inputs\":[{\"name\":\"updatable\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"writeConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"readConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"writeConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"readConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"newVaultUuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"allocateFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"baseFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"baseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"writeFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"readFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allocateFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"read\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accessAuxData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"readFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAllocateFee\",\"inputs\":[{\"name\":\"newAllocateFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseFee\",\"inputs\":[{\"name\":\"newBaseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReadFee\",\"inputs\":[{\"name\":\"newReadFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWriteFee\",\"inputs\":[{\"name\":\"newWriteFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitEncryptedPartialDecryption\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"pid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ephemeralPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"uuid\",\"inputs\":[],\"outputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vaults\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"vault\",\"type\":\"tuple\",\"internalType\":\"structICDR.Vault\",\"components\":[{\"name\":\"updatable\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"writeConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"readConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"writeConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"readConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"write\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accessAuxData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"writeFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"EncryptedPartialDecryptionSubmitted\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"pid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ephemeralPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeCollected\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"feeType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumICDR.FeeType\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultAllocated\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"updatable\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"writeConditionAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"readConditionAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"writeConditionData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"readConditionData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultRead\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultWritten\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60a080604052346100cc57306080527ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c166100bd57506001600160401b036002600160401b031982821601610078575b604051612aaf90816100d182396080518181816120c8015261218f0152f35b6001600160401b031990911681179091556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f8080610059565b63f92ee8a960e01b8152600490fd5b5f80fdfe60a06040526004361015610011575f80fd5b5f3560e01c80630dae55281461243c57806346860698146123fc5780634f1ef2861461214257806352d1902d146120a15780635c975abb1461206057806361cfc52c14611e145780636954b1cd14611dd85780636ef25c3a14611d9c578063715018a614611cd75780637741958414611c9b57806379af44d81461162f57806379ba5097146115a65780638da5cb5b146115545780638e3850c71461151457806390ab171e146114d85780639cadce2414610bf6578063ad3cb1cc14610b59578063b98a6c1a1461074a578063bb07e85d14610708578063e30c3978146106b6578063e60adef114610676578063f2fde38b146105aa578063f8b5ac35146104255763f92ad21914610121575f80fd5b346104215760a06003193601126104215761013a61247c565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c16159167ffffffffffffffff811680159081610419575b600114908161040f575b159081610406575b506103dc578260017fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000083161785556103a7575b506101ca612980565b6101d2612980565b73ffffffffffffffffffffffffffffffffffffffff811615610377576101f790612873565b6101ff612980565b610207612980565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055610233612980565b61023b612980565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00815416905561028a612980565b6024357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf601556044357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf602556064357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf603556084357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf6045561032457005b7fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff81541690557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160018152a1005b60246040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081525f6004820152fd5b7fffffffffffffffffffffffffffffffffffffffffffffff00000000000000000016680100000000000000011783555f6101c1565b60046040517ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b9050155f61018e565b303b159150610186565b84915061017c565b5f80fd5b346104215760206003193601126104215760a06104a361044361258a565b60405161044f8161249f565b5f81525f60208201525f604082015260609381858080940152826080820152015263ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b906105a6604051916104b48361249f565b83549360ff85161515845261059561056173ffffffffffffffffffffffffffffffffffffffff80602088019860081c168852806001850154169460408801958652610501600286016126b6565b918189019283526105276004610519600389016126b6565b9760808c01988952016126b6565b9660a08a01978852816040519b8c9b60208d5251151560208d0152511660408b01525116908801525160c0608088015260e08701906125cb565b9151917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe092838783030160a08801526125cb565b9151908483030160c08501526125cb565b0390f35b34610421576020600319360112610421576105c361247c565b6105cb6127de565b73ffffffffffffffffffffffffffffffffffffffff809116907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227005f80a3005b346104215760206003193601126104215761068f6127de565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60455005b34610421575f60031936011261042157602073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416604051908152f35b34610421575f60031936011261042157602063ffffffff7f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf6005416604051908152f35b6003196060813601126104215761075f61258a565b9067ffffffffffffffff906024358281116104215761078290369060040161256c565b916044359081116104215761079b90369060040161259d565b6107a6949194612926565b6107ae61281e565b6107e58263ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b92604051946107f38661249f565b84549460ff86161515875273ffffffffffffffffffffffffffffffffffffffff8060209760081c16878901528060018301541660408901908152610839600284016126b6565b60608a015260a061085f6004610851600387016126b6565b9560808d01968752016126b6565b990198808a525115610afb5790879392915116908133036109e6575b50505050507f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf603549081340361098857815f811561097f575b5f808093818094f1156109745761093b63ffffffff9460607fc42aa8cd96f3895f59a8e3e596ef5c3326d70ca274c0be540724fa4093a0a6349761094b956040519081526002848201527f93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f740031860403392a2519160405197889716875286015260608501906125cb565b9083820360408501523396612627565b0390a260017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055005b6040513d5f823e3d90fd5b506108fc6108b3565b606484604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601760248201527f4344523a20496e76616c69642066656520616d6f756e740000000000000000006044820152fd5b610a39905192610a48604051968795869485947f8db3eb1700000000000000000000000000000000000000000000000000000000865263ffffffff8d1660048701526080602487015260848601906125cb565b918483030160448501526125cb565b33606483015203915afa908115610974575f91610ace575b5015610a7057858381808061087b565b606483604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601b60248201527f4344523a205265616420636f6e646974696f6e206e6f74206d657400000000006044820152fd5b610aee9150843d8611610af4575b610ae681836124bb565b810190612777565b86610a60565b503d610adc565b606488604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601e60248201527f4344523a205661756c7420686173206e6f206461746120746f207265616400006044820152fd5b34610421575f60031936011261042157604051604081019080821067ffffffffffffffff831117610bc9576105a691604052600581527f352e302e3000000000000000000000000000000000000000000000000000000060208201526040519182916020835260208301906125cb565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b60a0600319360112610421576004358015158103610421576024359073ffffffffffffffffffffffffffffffffffffffff82168203610421576044359073ffffffffffffffffffffffffffffffffffffffff821682036104215760643567ffffffffffffffff811161042157610c7090369060040161259d565b919060843567ffffffffffffffff811161042157610c9290369060040161259d565b919093610c9d61281e565b73ffffffffffffffffffffffffffffffffffffffff8716158015906114b9575b1561145b577f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf604548034036113fd57805f81156113f4575b5f808093818094f115610974576040519081525f60208201527f93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f740031860403392a27f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf600549563ffffffff808816146113c75763ffffffff600181891601167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000008816177f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf6005560405197610dc28961249f565b851515895273ffffffffffffffffffffffffffffffffffffffff811660208a015273ffffffffffffffffffffffffffffffffffffffff821660408a0152610e0a368486612536565b60608a0152610e1a368689612536565b60808a015260405180602081011067ffffffffffffffff602083011117610bc957602081016040525f815260a08a0152610e8763ffffffff891663ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b98805115158a547fffffffffffffffffffffff00000000000000000000000000000000000000000060ff74ffffffffffffffffffffffffffffffffffffffff00602086015160081b169316911617178a5560018a0173ffffffffffffffffffffffffffffffffffffffff6040830151167fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790556060810151805167ffffffffffffffff8111610bc957610f4e818d6002610f4781830154612665565b910161278f565b60208c601f831160011461132057508190610f9c935f92611315575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60028b01555b608081015180519067ffffffffffffffff8211610bc957610fcc828d6003610f4781830154612665565b6020908c601f841160011461126b5750918061101e9260a095945f926112605750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60038c01555b01519586519967ffffffffffffffff8b11610bc9578a8a9861105860209d61104f6004860154612665565b6004860161278f565b8c90601f831160011461114c576111229561113099957ff1370099e56e061a64615edc07c1d17e36a20585cdc9b288bf0259a528365d0a9d999560046110f263ffffffff9f9e9b96978073ffffffffffffffffffffffffffffffffffffffff998a985f926111415750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9101555b6040519c8d9c168c5215158f8c01521660408a015216606088015260c0608088015260c0870191612627565b9184830360a0860152612627565b0390a163ffffffff60405191168152f35b015190505f80610f6a565b90600484015f528d5f20915f5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08516811061124657509561113099957ff1370099e56e061a64615edc07c1d17e36a20585cdc9b288bf0259a528365d0a9d99956004600163ffffffff9f9e9b969773ffffffffffffffffffffffffffffffffffffffff9889976111229d837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081161061120f575b505050811b019101556110f6565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690555f8080611201565b8282015184558e9c50600190930192918f01918f01611159565b015190508e80610f6a565b919060037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0851693015f5260205f20925f5b8181106112fd575091600193918560a0979694106112c6575b505050811b0160038c0155611024565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558d80806112b6565b9293602060018192878601518155019501930161129d565b015190508d80610f6a565b919260027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0851693015f5260205f20925f5b8181106113af5750908460019594939210611378575b505050811b0160028b0155610fa2565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558c8080611368565b92936020600181928786015181550195019301611352565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b506108fc610cf4565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4344523a20496e76616c69642066656520616d6f756e740000000000000000006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f496e76616c696420636f6e646974696f6e2061646472657373000000000000006044820152fd5b5073ffffffffffffffffffffffffffffffffffffffff86161515610cbd565b34610421575f6003193601126104215760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60354604051908152f35b346104215760206003193601126104215761152d6127de565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60355005b34610421575f60031936011261042157602073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416604051908152f35b34610421575f600319360112610421573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416036115ff576115fd33612873565b005b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b60606003193601126104215761164361258a565b60243567ffffffffffffffff81116104215761166390369060040161259d565b919060443567ffffffffffffffff81116104215761168590369060040161259d565b939091611690612926565b61169861281e565b8415611c17576116d58463ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b916117586004604051946116e88661249f565b73ffffffffffffffffffffffffffffffffffffffff815460ff81161515885260081c16602087015273ffffffffffffffffffffffffffffffffffffffff600182015416604087015261173c600282016126b6565b606087015261174d600382016126b6565b6080870152016126b6565b60a084015273ffffffffffffffffffffffffffffffffffffffff60208401511615611b945773ffffffffffffffffffffffffffffffffffffffff60208401511690813303611a87575b50505060a081015151611a1d575b507f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf602548034036113fd57805f8115611a14575b5f808093818094f11561097457604051908152600160208201527f93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f740031860403392a2600461185b8363ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b019267ffffffffffffffff8111610bc9576118808161187a8654612665565b8661278f565b5f93601f821160011461194957918163ffffffff94936118f882611915957f68b7452742f8d8707af90ddff5ef1a7f8850f5724b804e72b4df1a37050a6355995f9161193e575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b604051948594168452604060208501526040840191612627565b0390a160017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055005b90508501358a6118c7565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08216815f5260205f20905f5b8181106119fc57509183917f68b7452742f8d8707af90ddff5ef1a7f8850f5724b804e72b4df1a37050a6355976119159563ffffffff989795106119c4575b5050600182811b0190556118fb565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c199085013516905587806119b5565b85880135835560209788019760019093019201611976565b506108fc6117e2565b5115611a2957836117af565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f4344523a205661756c74206973206e6f7420757064617461626c6500000000006044820152fd5b60209163ffffffff87611af16060880151611adf604051988997889687967f5645dbbf000000000000000000000000000000000000000000000000000000008852166004870152608060248701526084860191612627565b906003198483030160448501526125cb565b33606483015203915afa908115610974575f91611b75575b5015611b17578480806117a1565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f4344523a20577269746520636f6e646974696f6e206e6f74206d6574000000006044820152fd5b611b8e915060203d602011610af457610ae681836124bb565b85611b09565b60846040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f4344523a20577269746520636f6e646974696f6e2061646472657373206e6f7460448201527f20736574000000000000000000000000000000000000000000000000000000006064820152fd5b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f4344523a20456e6372797074656420646174612063616e6e6f7420626520656d60448201527f70747900000000000000000000000000000000000000000000000000000000006064820152fd5b34610421575f6003193601126104215760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60454604051908152f35b34610421575f60031936011261042157611cef6127de565b5f73ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b34610421575f6003193601126104215760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60154604051908152f35b34610421575f6003193601126104215760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60254604051908152f35b6101208060031936011261042157611e2a61258a565b6024359063ffffffff821682036104215767ffffffffffffffff60443581811161042157611e5c90369060040161259d565b909260643583811161042157611e7690369060040161259d565b95909460843585811161042157611e9190369060040161259d565b93909460a43587811161042157611eac90369060040161259d565b919060c43589811161042157611ec690369060040161259d565b93909460e4359a63ffffffff8c168c03610421576101043590811161042157611ef390369060040161259d565b9a9099611efe61281e565b7f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf601549e8f34036113fd578f5f80808093818115612057575b8290f115610974578f9e6120239f91611fc16120089c611ff69a63ffffffff9f9863ffffffff611fe49a611fd398604051908152600360208201527f93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f740031860403392a260405160805281610140931660805152166020608051015280604060805101526080510191612627565b91608051830360606080510152612627565b916080518303608080510152612627565b91608051830360a06080510152612627565b91608051830360c06080510152612627565b931660e0608051015260805183036101006080510152612627565b9160805101527fe2d964436c9da9c84292bd1b67e737171e45dbb6ca28f34fcbfb5c9ce20c1dbc33916080519003608051a2005b506108fc611f36565b34610421575f60031936011261042157602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b34610421575f6003193601126104215773ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001630036121185760206040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b60046040517fe07c8dba000000000000000000000000000000000000000000000000000000008152fd5b60406003193601126104215761215661247c565b60243567ffffffffffffffff81116104215761217690369060040161256c565b9073ffffffffffffffffffffffffffffffffffffffff807f0000000000000000000000000000000000000000000000000000000000000000168030149081156123ce575b50612118576121c76127de565b811690604051907f52d1902d0000000000000000000000000000000000000000000000000000000082526020918281600481875afa5f918161239f575b5061223a57602484604051907f4c9c8ce30000000000000000000000000000000000000000000000000000000082526004820152fd5b9284937f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc9081810361236e5750823b1561233d57817fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055604051907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a283511561230a57505f8084846115fd96519101845af4903d15612301573d6122e5816124fc565b906122f360405192836124bb565b81525f81943d92013e6129d9565b606092506129d9565b925050503461231557005b807fb398979f0000000000000000000000000000000000000000000000000000000060049252fd5b602482604051907f4c9c8ce30000000000000000000000000000000000000000000000000000000082526004820152fd5b602490604051907faa1d49a40000000000000000000000000000000000000000000000000000000082526004820152fd5b9091508381813d83116123c7575b6123b781836124bb565b8101031261042157519086612204565b503d6123ad565b9050817f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc54161415846121ba565b34610421576020600319360112610421576124156127de565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60155005b34610421576020600319360112610421576124556127de565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60255005b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361042157565b60c0810190811067ffffffffffffffff821117610bc957604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff821117610bc957604052565b67ffffffffffffffff8111610bc957601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b929192612542826124fc565b9161255060405193846124bb565b829481845281830111610421578281602093845f960137010152565b9080601f830112156104215781602061258793359101612536565b90565b6004359063ffffffff8216820361042157565b9181601f840112156104215782359167ffffffffffffffff8311610421576020838186019501011161042157565b91908251928382525f5b8481106126135750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f845f6020809697860101520116010190565b6020818301810151848301820152016125d5565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe093818652868601375f8582860101520116010190565b90600182811c921680156126ac575b602083101461267f57565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f1691612674565b9060405191825f82546126c881612665565b908184526020946001916001811690815f1461273657506001146126f8575b5050506126f6925003836124bb565b565b5f90815285812095935091905b81831061271e5750506126f693508201015f80806126e7565b85548884018501529485019487945091830191612705565b9150506126f69593507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201015f80806126e7565b90816020910312610421575180151581036104215790565b601f821161279c57505050565b5f5260205f20906020601f840160051c830193106127d4575b601f0160051c01905b8181106127c9575050565b5f81556001016127be565b90915081906127b5565b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300541633036115ff57565b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300541661284957565b60046040517fd93c0665000000000000000000000000000000000000000000000000000000008152fd5b7fffffffffffffffffffffffff0000000000000000000000000000000000000000907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008281541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549073ffffffffffffffffffffffffffffffffffffffff80931680948316179055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0060028154146129565760029055565b60046040517f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c16156129af57565b60046040517fd7e6bcf8000000000000000000000000000000000000000000000000000000008152fd5b90612a1857508051156129ee57805190602001fd5b60046040517f1425ea42000000000000000000000000000000000000000000000000000000008152fd5b81511580612a70575b612a29575090565b60249073ffffffffffffffffffffffffffffffffffffffff604051917f9996b315000000000000000000000000000000000000000000000000000000008352166004820152fd5b50803b15612a2156fea2646970667358221220d1917d22dfbbe36a5a9f70fb1c525ee9ca7cb26bbca6f49c8bc127ef22a10c4e64736f6c63430008170033",
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
// Solidity: function vaults(uint32 uuid) view returns((bool,address,address,bytes,bytes,bytes) vault)
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
// Solidity: function vaults(uint32 uuid) view returns((bool,address,address,bytes,bytes,bytes) vault)
func (_CDR *CDRSession) Vaults(uuid uint32) (ICDRVault, error) {
	return _CDR.Contract.Vaults(&_CDR.CallOpts, uuid)
}

// Vaults is a free data retrieval call binding the contract method 0xf8b5ac35.
//
// Solidity: function vaults(uint32 uuid) view returns((bool,address,address,bytes,bytes,bytes) vault)
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

// Write is a paid mutator transaction binding the contract method 0x79af44d8.
//
// Solidity: function write(uint32 uuid, bytes accessAuxData, bytes encryptedData) payable returns()
func (_CDR *CDRTransactor) Write(opts *bind.TransactOpts, uuid uint32, accessAuxData []byte, encryptedData []byte) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "write", uuid, accessAuxData, encryptedData)
}

// Write is a paid mutator transaction binding the contract method 0x79af44d8.
//
// Solidity: function write(uint32 uuid, bytes accessAuxData, bytes encryptedData) payable returns()
func (_CDR *CDRSession) Write(uuid uint32, accessAuxData []byte, encryptedData []byte) (*types.Transaction, error) {
	return _CDR.Contract.Write(&_CDR.TransactOpts, uuid, accessAuxData, encryptedData)
}

// Write is a paid mutator transaction binding the contract method 0x79af44d8.
//
// Solidity: function write(uint32 uuid, bytes accessAuxData, bytes encryptedData) payable returns()
func (_CDR *CDRTransactorSession) Write(uuid uint32, accessAuxData []byte, encryptedData []byte) (*types.Transaction, error) {
	return _CDR.Contract.Write(&_CDR.TransactOpts, uuid, accessAuxData, encryptedData)
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
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterVaultRead is a free log retrieval operation binding the contract event 0xc42aa8cd96f3895f59a8e3e596ef5c3326d70ca274c0be540724fa4093a0a634.
//
// Solidity: event VaultRead(uint32 uuid, address indexed requester, bytes ciphertext, bytes requesterPubKey)
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

// WatchVaultRead is a free log subscription operation binding the contract event 0xc42aa8cd96f3895f59a8e3e596ef5c3326d70ca274c0be540724fa4093a0a634.
//
// Solidity: event VaultRead(uint32 uuid, address indexed requester, bytes ciphertext, bytes requesterPubKey)
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

// ParseVaultRead is a log parse operation binding the contract event 0xc42aa8cd96f3895f59a8e3e596ef5c3326d70ca274c0be540724fa4093a0a634.
//
// Solidity: event VaultRead(uint32 uuid, address indexed requester, bytes ciphertext, bytes requesterPubKey)
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
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterVaultWritten is a free log retrieval operation binding the contract event 0x68b7452742f8d8707af90ddff5ef1a7f8850f5724b804e72b4df1a37050a6355.
//
// Solidity: event VaultWritten(uint32 uuid, bytes encryptedData)
func (_CDR *CDRFilterer) FilterVaultWritten(opts *bind.FilterOpts) (*CDRVaultWrittenIterator, error) {

	logs, sub, err := _CDR.contract.FilterLogs(opts, "VaultWritten")
	if err != nil {
		return nil, err
	}
	return &CDRVaultWrittenIterator{contract: _CDR.contract, event: "VaultWritten", logs: logs, sub: sub}, nil
}

// WatchVaultWritten is a free log subscription operation binding the contract event 0x68b7452742f8d8707af90ddff5ef1a7f8850f5724b804e72b4df1a37050a6355.
//
// Solidity: event VaultWritten(uint32 uuid, bytes encryptedData)
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

// ParseVaultWritten is a log parse operation binding the contract event 0x68b7452742f8d8707af90ddff5ef1a7f8850f5724b804e72b4df1a37050a6355.
//
// Solidity: event VaultWritten(uint32 uuid, bytes encryptedData)
func (_CDR *CDRFilterer) ParseVaultWritten(log types.Log) (*CDRVaultWritten, error) {
	event := new(CDRVaultWritten)
	if err := _CDR.contract.UnpackLog(event, "VaultWritten", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
