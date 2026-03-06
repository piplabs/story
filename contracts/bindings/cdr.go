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
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocate\",\"inputs\":[{\"name\":\"updatable\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"writeConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"readConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"writeConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"readConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"newVaultUuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"allocateFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"baseFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"read\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accessAuxData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"recipientPublicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"readFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAllocateFee\",\"inputs\":[{\"name\":\"newAllocateFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseFee\",\"inputs\":[{\"name\":\"newBaseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReadFee\",\"inputs\":[{\"name\":\"newReadFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWriteFee\",\"inputs\":[{\"name\":\"newWriteFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitEncryptedPartialDecryption\",\"inputs\":[{\"name\":\"enclaveID\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"uuid\",\"inputs\":[],\"outputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vaults\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"vault\",\"type\":\"tuple\",\"internalType\":\"structICDR.Vault\",\"components\":[{\"name\":\"updatable\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"writeConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"readConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"writeConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"readConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"write\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accessAuxData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"writeFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"EncryptedPartialDecryptionSubmitted\",\"inputs\":[{\"name\":\"enclaveID\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultAllocated\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"updatable\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"writeConditionAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"readConditionAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"writeConditionData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"readConditionData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultRead\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"recipientPublicKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultWritten\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60a0806040523461002a573060805261240d908161002f8239608051818181611aab0152611b720152f35b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c80630dae552814611e215780634686069814611de15780634f1ef28614611b2557806352d1902d14611a845780635c975abb14611a435780636954b1cd14611a075780636ef25c3a146119cb578063715018a61461190657806377419584146118ca57806379af44d8146112a557806379ba5097146111a85780638da5cb5b146111565780638e3850c71461111657806390ab171e146110da5780639cadce2414610913578063ad3cb1cc14610876578063af26423b14610794578063b98a6c1a1461043f578063bb07e85d146103fd578063e30c3978146103ab578063e60adef11461036b578063f2fde38b1461029f5763f8b5ac3514610116575f80fd5b3461029b57602060031936011261029b5760a0610194610134611f6f565b60405161014081611e84565b5f81525f60208201525f604082015260609381858080940152826080820152015263ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b90610297604051916101a583611e84565b83549360ff85161515845261028661025273ffffffffffffffffffffffffffffffffffffffff80602088019860081c1688528060018501541694604088019586526101f26002860161205d565b91818901928352610218600461020a6003890161205d565b9760808c019889520161205d565b9660a08a01978852816040519b8c9b60208d5251151560208d0152511660408b01525116908801525160c0608088015260e0870190611fb0565b9151917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe092838783030160a0880152611fb0565b9151908483030160c0850152611fb0565b0390f35b5f80fd5b3461029b57602060031936011261029b576102b8611e61565b6102c06121c3565b73ffffffffffffffffffffffffffffffffffffffff809116907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227005f80a3005b3461029b57602060031936011261029b576103846121c3565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60455005b3461029b575f60031936011261029b57602073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416604051908152f35b3461029b575f60031936011261029b57602063ffffffff7f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf6005416604051908152f35b60031960608136011261029b57610454611f6f565b67ffffffffffffffff60243581811161029b57610475903690600401611f51565b9060443590811161029b5761048e903690600401611f82565b9091610498612203565b6104a061225d565b6104d78463ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b90604051916104e583611e84565b80549660ff88161515845273ffffffffffffffffffffffffffffffffffffffff8060209960081c1689860152806001840154166040860190815261052b6002850161205d565b606087015260a061055160046105436003880161205d565b9660808a019788520161205d565b960195808752511561073657511691823303610614575b7f38cf6a5ff02d5d4c4e1e1ef8e4d4c19035d060c11ecb6677827a85315365505863ffffffff896105eb8a8a6105dd8f60608d6105c57f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf603546122b2565b51916040519889981688528701526060860190611fb0565b918483036040860152612136565b0390a160017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055005b9163ffffffff8861067961066a96948c9651604051988997889687967f8db3eb17000000000000000000000000000000000000000000000000000000008852166004870152608060248701526084860190611fb0565b91848303016044850152611fb0565b33606483015203915afa90811561072b575f916106fe575b50156106a05785808080610568565b606485604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601b60248201527f4344523a205265616420636f6e646974696f6e206e6f74206d657400000000006044820152fd5b61071e9150863d8811610724575b6107168183611ea0565b81019061211e565b86610691565b503d61070c565b6040513d5f823e3d90fd5b60648a604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601e60248201527f4344523a205661756c7420686173206e6f206461746120746f207265616400006044820152fd5b606060031936011261029b576107a8611e61565b67ffffffffffffffff60243581811161029b576107c9903690600401611f82565b92909160443590811161029b577f280b85b1fa1bce22ee09c7a1dc856c913dc966030c38d851bb88b08b5c6e321f9373ffffffffffffffffffffffffffffffffffffffff936105dd610822610871943690600401611f82565b92909361082d61225d565b6108577f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf601546122b2565b604051978897168752606060208801526060870191612136565b0390a1005b3461029b575f60031936011261029b57604051604081019080821067ffffffffffffffff8311176108e65761029791604052600581527f352e302e300000000000000000000000000000000000000000000000000000006020820152604051918291602083526020830190611fb0565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b60a060031936011261029b57600435801515810361029b5760243573ffffffffffffffffffffffffffffffffffffffff8116810361029b576044359173ffffffffffffffffffffffffffffffffffffffff8316830361029b5760643567ffffffffffffffff811161029b5761098c903690600401611f82565b9060843567ffffffffffffffff811161029b576109ad903690600401611f82565b9190926109b861225d565b73ffffffffffffffffffffffffffffffffffffffff8616158015906110bb575b1561105d57610a077f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf604546122b2565b7f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf6009586549563ffffffff97888089161461103057886001818a1601167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000008916179055604051610a7581611e84565b811515815273ffffffffffffffffffffffffffffffffffffffff8316602082015273ffffffffffffffffffffffffffffffffffffffff8a166040820152610abd368587611f1b565b6060820152610acd368789611f1b565b608082015260405180602081011067ffffffffffffffff6020830111176108e657602081016040525f815260a0820152610b3689891663ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b8151151581547fffffffffffffffffffffff00000000000000000000000000000000000000000060ff74ffffffffffffffffffffffffffffffffffffffff00602087015160081b1693169116171781556001810173ffffffffffffffffffffffffffffffffffffffff6040840151167fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055606082015180519067ffffffffffffffff82116108e657610bfd82610bf4600286015461200c565b60028601612174565b602090601f8311600114610f8957610c4a92915f9183610f7e575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60028201555b608082015180519067ffffffffffffffff82116108e657610c8182610c78600386015461200c565b60038601612174565b602090601f8311600114610ed2578260049360a09593610cd3935f92610de05750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60038201555b0191015180519067ffffffffffffffff82116108e657610d0382610cfd855461200c565b85612174565b602090601f8311600114610deb5794610dc5969473ffffffffffffffffffffffffffffffffffffffff94610d957ff1370099e56e061a64615edc07c1d17e36a20585cdc9b288bf0259a528365d0a9f95808896610dd39f9e9d9b5f92610de05750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b6040519c8d9c168c52151560208c01521660408a015216606088015260c0608088015260c0870191612136565b9184830360a0860152612136565b0390a160206040515f8152f35b015190505f80610c18565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0831691845f5260205f20925f5b818110610eba57509460017ff1370099e56e061a64615edc07c1d17e36a20585cdc9b288bf0259a528365d0a9f9573ffffffffffffffffffffffffffffffffffffffff9795610dd39e9d9c9a95610dc59c9a95838b9910610e83575b505050811b019055610d98565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690555f8080610e76565b92936020600181928786015181550195019301610e1a565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0831691600385015f5260205f20925f5b818110610f66575092600192859260049660a0989610610f2f575b505050811b016003820155610cd9565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558f8080610f1f565b92936020600181928786015181550195019301610f04565b015190508e80610c18565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0831691600285015f5260205f20925f5b8181106110185750908460019594939210610fe1575b505050811b016002820155610c50565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558d8080610fd1565b92936020600181928786015181550195019301610fbb565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f496e76616c696420636f6e646974696f6e2061646472657373000000000000006044820152fd5b5073ffffffffffffffffffffffffffffffffffffffff871615156109d8565b3461029b575f60031936011261029b5760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60354604051908152f35b3461029b57602060031936011261029b5761112f6121c3565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60355005b3461029b575f60031936011261029b57602073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416604051908152f35b3461029b575f60031936011261029b577f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0080549073ffffffffffffffffffffffffffffffffffffffff903382841603611275577fffffffffffffffffffffffff000000000000000000000000000000000000000080931690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805492339084161790553391167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3005b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b606060031936011261029b576112b9611f6f565b602490813567ffffffffffffffff811161029b576112db903690600401611f82565b60443567ffffffffffffffff811161029b576112fb903690600401611f82565b919092611306612203565b61130e61225d565b82156118465761134b8563ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b906040519261135984611e84565b73ffffffffffffffffffffffffffffffffffffffff835460ff81161515865260081c16602085015273ffffffffffffffffffffffffffffffffffffffff60018401541660408501526113d160046113b26002860161205d565b94606087019586526113c66003820161205d565b60808801520161205d565b9260a0850193845273ffffffffffffffffffffffffffffffffffffffff602086015116156117c357889073ffffffffffffffffffffffffffffffffffffffff60208701511690898233036116b7575b505050505050515161164d575b506114587f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf602546122b2565b60046114918463ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60560205260405f2090565b019367ffffffffffffffff821161162257506114b7816114b1865461200c565b86612174565b5f93601f821160011461155757918163ffffffff949361152f826105eb957f68b7452742f8d8707af90ddff5ef1a7f8850f5724b804e72b4df1a37050a6355995f9161154c575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b604051948594168452604060208501526040840191612136565b90508501358a6114fe565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08216815f5260205f20905f5b81811061160a57509183917f68b7452742f8d8707af90ddff5ef1a7f8850f5724b804e72b4df1a37050a6355976105eb9563ffffffff989795106115d2575b5050600182811b019055611532565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c199085013516905587806115c3565b85880135835560209788019760019093019201611584565b7f4e487b71000000000000000000000000000000000000000000000000000000005f5260416004525ffd5b5115611659578461142d565b606484601b604051917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401528201527f4344523a205661756c74206973206e6f7420757064617461626c6500000000006044820152fd5b61171d60209561170b60809863ffffffff9551926040519a8b998a9889987f5645dbbf000000000000000000000000000000000000000000000000000000008a521660048901528701526084860191612136565b90600319848303016044850152611fb0565b33606483015203915afa90811561072b575f916117a4575b501561174657868087818089611420565b606486601c604051917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401528201527f4344523a20577269746520636f6e646974696f6e206e6f74206d6574000000006044820152fd5b6117bd915060203d602011610724576107168183611ea0565b87611735565b608489604051907f08c379a000000000000000000000000000000000000000000000000000000000825260206004830152808201527f4344523a20577269746520636f6e646974696f6e2061646472657373206e6f7460448201527f20736574000000000000000000000000000000000000000000000000000000006064820152fd5b6084866023604051917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401528201527f4344523a20456e6372797074656420646174612063616e6e6f7420626520656d60448201527f70747900000000000000000000000000000000000000000000000000000000006064820152fd5b3461029b575f60031936011261029b5760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60454604051908152f35b3461029b575f60031936011261029b5761191e6121c3565b5f73ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b3461029b575f60031936011261029b5760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60154604051908152f35b3461029b575f60031936011261029b5760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60254604051908152f35b3461029b575f60031936011261029b57602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b3461029b575f60031936011261029b5773ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163003611afb5760206040517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b60046040517fe07c8dba000000000000000000000000000000000000000000000000000000008152fd5b604060031936011261029b57611b39611e61565b60243567ffffffffffffffff811161029b57611b59903690600401611f51565b9073ffffffffffffffffffffffffffffffffffffffff807f000000000000000000000000000000000000000000000000000000000000000016803014908115611db3575b50611afb57611baa6121c3565b811690604051907f52d1902d0000000000000000000000000000000000000000000000000000000082526020918281600481875afa5f9181611d84575b50611c1d57602484604051907f4c9c8ce30000000000000000000000000000000000000000000000000000000082526004820152fd5b9284937f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc90818103611d535750823b15611d2257817fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055604051907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a2835115611cef57505f808484611ce496519101845af4903d15611ce6573d611cc881611ee1565b90611cd66040519283611ea0565b81525f81943d92013e612337565b005b60609250612337565b9250505034611cfa57005b807fb398979f0000000000000000000000000000000000000000000000000000000060049252fd5b602482604051907f4c9c8ce30000000000000000000000000000000000000000000000000000000082526004820152fd5b602490604051907faa1d49a40000000000000000000000000000000000000000000000000000000082526004820152fd5b9091508381813d8311611dac575b611d9c8183611ea0565b8101031261029b57519086611be7565b503d611d92565b9050817f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5416141584611b9d565b3461029b57602060031936011261029b57611dfa6121c3565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60155005b3461029b57602060031936011261029b57611e3a6121c3565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60255005b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361029b57565b60c0810190811067ffffffffffffffff8211176108e657604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176108e657604052565b67ffffffffffffffff81116108e657601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b929192611f2782611ee1565b91611f356040519384611ea0565b82948184528183011161029b578281602093845f960137010152565b9080601f8301121561029b57816020611f6c93359101611f1b565b90565b6004359063ffffffff8216820361029b57565b9181601f8401121561029b5782359167ffffffffffffffff831161029b576020838186019501011161029b57565b91908251928382525f5b848110611ff85750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f845f6020809697860101520116010190565b602081830181015184830182015201611fba565b90600182811c92168015612053575b602083101461202657565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f169161201b565b9060405191825f825461206f8161200c565b908184526020946001916001811690815f146120dd575060011461209f575b50505061209d92500383611ea0565b565b5f90815285812095935091905b8183106120c557505061209d93508201015f808061208e565b855488840185015294850194879450918301916120ac565b91505061209d9593507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201015f808061208e565b9081602091031261029b5751801515810361029b5790565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe093818652868601375f8582860101520116010190565b601f821161218157505050565b5f5260205f20906020601f840160051c830193106121b9575b601f0160051c01905b8181106121ae575050565b5f81556001016121a3565b909150819061219a565b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416330361127557565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0060028154146122335760029055565b60046040517f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300541661228857565b60046040517fd93c0665000000000000000000000000000000000000000000000000000000008152fd5b8034036122d9575f808080938181156122d0575b8290f11561072b57565b506108fc6122c6565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4344523a20496e76616c69642066656520616d6f756e740000000000000000006044820152fd5b90612376575080511561234c57805190602001fd5b60046040517f1425ea42000000000000000000000000000000000000000000000000000000008152fd5b815115806123ce575b612387575090565b60249073ffffffffffffffffffffffffffffffffffffffff604051917f9996b315000000000000000000000000000000000000000000000000000000008352166004820152fd5b50803b1561237f56fea2646970667358221220a27b5423438b1b4c7f72f6f4a266e90d4c914a12e05099848f9164545b820c7864736f6c63430008170033",
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

// Read is a paid mutator transaction binding the contract method 0xb98a6c1a.
//
// Solidity: function read(uint32 uuid, bytes accessAuxData, bytes recipientPublicKey) payable returns()
func (_CDR *CDRTransactor) Read(opts *bind.TransactOpts, uuid uint32, accessAuxData []byte, recipientPublicKey []byte) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "read", uuid, accessAuxData, recipientPublicKey)
}

// Read is a paid mutator transaction binding the contract method 0xb98a6c1a.
//
// Solidity: function read(uint32 uuid, bytes accessAuxData, bytes recipientPublicKey) payable returns()
func (_CDR *CDRSession) Read(uuid uint32, accessAuxData []byte, recipientPublicKey []byte) (*types.Transaction, error) {
	return _CDR.Contract.Read(&_CDR.TransactOpts, uuid, accessAuxData, recipientPublicKey)
}

// Read is a paid mutator transaction binding the contract method 0xb98a6c1a.
//
// Solidity: function read(uint32 uuid, bytes accessAuxData, bytes recipientPublicKey) payable returns()
func (_CDR *CDRTransactorSession) Read(uuid uint32, accessAuxData []byte, recipientPublicKey []byte) (*types.Transaction, error) {
	return _CDR.Contract.Read(&_CDR.TransactOpts, uuid, accessAuxData, recipientPublicKey)
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

// SubmitEncryptedPartialDecryption is a paid mutator transaction binding the contract method 0xaf26423b.
//
// Solidity: function submitEncryptedPartialDecryption(address enclaveID, bytes encryptedPartial, bytes signature) payable returns()
func (_CDR *CDRTransactor) SubmitEncryptedPartialDecryption(opts *bind.TransactOpts, enclaveID common.Address, encryptedPartial []byte, signature []byte) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "submitEncryptedPartialDecryption", enclaveID, encryptedPartial, signature)
}

// SubmitEncryptedPartialDecryption is a paid mutator transaction binding the contract method 0xaf26423b.
//
// Solidity: function submitEncryptedPartialDecryption(address enclaveID, bytes encryptedPartial, bytes signature) payable returns()
func (_CDR *CDRSession) SubmitEncryptedPartialDecryption(enclaveID common.Address, encryptedPartial []byte, signature []byte) (*types.Transaction, error) {
	return _CDR.Contract.SubmitEncryptedPartialDecryption(&_CDR.TransactOpts, enclaveID, encryptedPartial, signature)
}

// SubmitEncryptedPartialDecryption is a paid mutator transaction binding the contract method 0xaf26423b.
//
// Solidity: function submitEncryptedPartialDecryption(address enclaveID, bytes encryptedPartial, bytes signature) payable returns()
func (_CDR *CDRTransactorSession) SubmitEncryptedPartialDecryption(enclaveID common.Address, encryptedPartial []byte, signature []byte) (*types.Transaction, error) {
	return _CDR.Contract.SubmitEncryptedPartialDecryption(&_CDR.TransactOpts, enclaveID, encryptedPartial, signature)
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
	EnclaveID        common.Address
	EncryptedPartial []byte
	Signature        []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterEncryptedPartialDecryptionSubmitted is a free log retrieval operation binding the contract event 0x280b85b1fa1bce22ee09c7a1dc856c913dc966030c38d851bb88b08b5c6e321f.
//
// Solidity: event EncryptedPartialDecryptionSubmitted(address enclaveID, bytes encryptedPartial, bytes signature)
func (_CDR *CDRFilterer) FilterEncryptedPartialDecryptionSubmitted(opts *bind.FilterOpts) (*CDREncryptedPartialDecryptionSubmittedIterator, error) {

	logs, sub, err := _CDR.contract.FilterLogs(opts, "EncryptedPartialDecryptionSubmitted")
	if err != nil {
		return nil, err
	}
	return &CDREncryptedPartialDecryptionSubmittedIterator{contract: _CDR.contract, event: "EncryptedPartialDecryptionSubmitted", logs: logs, sub: sub}, nil
}

// WatchEncryptedPartialDecryptionSubmitted is a free log subscription operation binding the contract event 0x280b85b1fa1bce22ee09c7a1dc856c913dc966030c38d851bb88b08b5c6e321f.
//
// Solidity: event EncryptedPartialDecryptionSubmitted(address enclaveID, bytes encryptedPartial, bytes signature)
func (_CDR *CDRFilterer) WatchEncryptedPartialDecryptionSubmitted(opts *bind.WatchOpts, sink chan<- *CDREncryptedPartialDecryptionSubmitted) (event.Subscription, error) {

	logs, sub, err := _CDR.contract.WatchLogs(opts, "EncryptedPartialDecryptionSubmitted")
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

// ParseEncryptedPartialDecryptionSubmitted is a log parse operation binding the contract event 0x280b85b1fa1bce22ee09c7a1dc856c913dc966030c38d851bb88b08b5c6e321f.
//
// Solidity: event EncryptedPartialDecryptionSubmitted(address enclaveID, bytes encryptedPartial, bytes signature)
func (_CDR *CDRFilterer) ParseEncryptedPartialDecryptionSubmitted(log types.Log) (*CDREncryptedPartialDecryptionSubmitted, error) {
	event := new(CDREncryptedPartialDecryptionSubmitted)
	if err := _CDR.contract.UnpackLog(event, "EncryptedPartialDecryptionSubmitted", log); err != nil {
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
	Uuid               uint32
	EncryptedData      []byte
	RecipientPublicKey []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterVaultRead is a free log retrieval operation binding the contract event 0x38cf6a5ff02d5d4c4e1e1ef8e4d4c19035d060c11ecb6677827a853153655058.
//
// Solidity: event VaultRead(uint32 uuid, bytes encryptedData, bytes recipientPublicKey)
func (_CDR *CDRFilterer) FilterVaultRead(opts *bind.FilterOpts) (*CDRVaultReadIterator, error) {

	logs, sub, err := _CDR.contract.FilterLogs(opts, "VaultRead")
	if err != nil {
		return nil, err
	}
	return &CDRVaultReadIterator{contract: _CDR.contract, event: "VaultRead", logs: logs, sub: sub}, nil
}

// WatchVaultRead is a free log subscription operation binding the contract event 0x38cf6a5ff02d5d4c4e1e1ef8e4d4c19035d060c11ecb6677827a853153655058.
//
// Solidity: event VaultRead(uint32 uuid, bytes encryptedData, bytes recipientPublicKey)
func (_CDR *CDRFilterer) WatchVaultRead(opts *bind.WatchOpts, sink chan<- *CDRVaultRead) (event.Subscription, error) {

	logs, sub, err := _CDR.contract.WatchLogs(opts, "VaultRead")
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

// ParseVaultRead is a log parse operation binding the contract event 0x38cf6a5ff02d5d4c4e1e1ef8e4d4c19035d060c11ecb6677827a853153655058.
//
// Solidity: event VaultRead(uint32 uuid, bytes encryptedData, bytes recipientPublicKey)
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
