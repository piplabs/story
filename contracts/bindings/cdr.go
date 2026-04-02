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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocate\",\"inputs\":[{\"name\":\"updatable\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"writeConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"readConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"writeConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"readConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"newVaultUuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"allocateFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"baseFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"baseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"writeFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"readFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allocateFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxEncryptedDataSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxEncryptedPartialSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxEncryptedDataSize\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxEncryptedPartialSize\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"read\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accessAuxData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"readFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAllocateFee\",\"inputs\":[{\"name\":\"newAllocateFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseFee\",\"inputs\":[{\"name\":\"newBaseFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMaxEncryptedDataSize\",\"inputs\":[{\"name\":\"newMaxEncryptedDataSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMaxEncryptedPartialSize\",\"inputs\":[{\"name\":\"newMaxEncryptedPartialSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReadFee\",\"inputs\":[{\"name\":\"newReadFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWriteFee\",\"inputs\":[{\"name\":\"newWriteFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitEncryptedPartialDecryption\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"pid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ephemeralPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"uuid\",\"inputs\":[],\"outputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vaults\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"vault\",\"type\":\"tuple\",\"internalType\":\"structICDR.Vault\",\"components\":[{\"name\":\"updatable\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"writeConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"readConditionAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"writeConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"readConditionData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"write\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"accessAuxData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"writeFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"EncryptedPartialDecryptionSubmitted\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"pid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ephemeralPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeCollected\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"feeType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumICDR.FeeType\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultAllocated\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"updatable\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"writeConditionAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"readConditionAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"writeConditionData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"readConditionData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultRead\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultWritten\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x608080604052346100b8577ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c166100a957506001600160401b036002600160401b031982821601610064575b60405161292d90816100bd8239f35b6001600160401b031990911681179091556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f8080610055565b63f92ee8a960e01b8152600490fd5b5f80fdfe60a06040526004361015610011575f80fd5b5f3560e01c80630dae5528146122b25780631cf9c8621461228e578063468606981461224e5780635c975abb1461220d57806361cfc52c14611ed85780636954b1cd14611e9c5780636ef25c3a14611e60578063715018a614611d9b5780637741958414611d5f57806379af44d81461166657806379ba5097146115df5780638142951a146112d15780638465ff72146112955780638da5cb5b146112435780638e3850c71461120357806390ab171e146111c75780639cadce2414610895578063a060fa1214610859578063b98a6c1a1461048b578063bb07e85d14610449578063dcaa51d414610423578063e30c3978146103d1578063e60adef114610391578063f2fde38b146102c55763f8b5ac351461012c575f80fd5b346102c15760206003193601126102c15760a06101aa61014a6122f2565b60405161015681612356565b5f81525f60208201525f604082015260609381858080940152826080820152015263ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60760205260405f2090565b6102bd6004604051926101bc84612356565b6102ac61029682549660ff88161515875273ffffffffffffffffffffffffffffffffffffffff9081602089019960081c1689528160018601541691604089019283526040519261021a846102138160028b016124a2565b0385612372565b828a0193845261025d6040519761023f8961023881600385016124a2565b038a612372565b60808c01988952610256604051809b8193016124a2565b0389612372565b60a08a01978852816040519b8c9b60208d5251151560208d0152511660408b01525116908801525160c0608088015260e08701906123db565b915191601f1992838783030160a08801526123db565b9151908483030160c08501526123db565b0390f35b5f80fd5b346102c15760206003193601126102c1576102de612333565b6102e66125a0565b73ffffffffffffffffffffffffffffffffffffffff809116907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227005f80a3005b346102c15760206003193601126102c1576103aa6125a0565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60455005b346102c1575f6003193601126102c157602073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416604051908152f35b346102c15760206003193601126102c15761043c6125a0565b6104476004356127f0565b005b346102c1575f6003193601126102c157602063ffffffff7f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf6005416604051908152f35b6003196060813601126102c1576104a06122f2565b602480359167ffffffffffffffff928381116102c157366023820112156102c1576104d49036908481600401359101612395565b926044359081116102c1576104ed903690600401612305565b9390946104f861268e565b6105006126e8565b6105378363ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60760205260405f2090565b9060048201926105478454612451565b156107fb57859073ffffffffffffffffffffffffffffffffffffffff60018501541690868233036106e1575b5050505050507f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60354928334036106835750825f811561067a575b5f808093818094f11561066f577fc42aa8cd96f3895f59a8e3e596ef5c3326d70ca274c0be540724fa4093a0a634936106366106469263ffffffff95604051908152600260208201527f93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f740031860403392a26040519586951685526060602086015260608501906124a2565b9083820360408501523396612419565b0390a260017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055005b6040513d5f823e3d90fd5b506108fc6105ad565b6064906017604051917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401528201527f4344523a20496e76616c69642066656520616d6f756e740000000000000000006044820152fd5b610745602095600360809861073563ffffffff966040519b8c9a8b998a997f8db3eb17000000000000000000000000000000000000000000000000000000008b521660048a015288015260848701906123db565b92858403016044860152016124a2565b33606483015203915afa90811561066f575f916107cc575b501561076e57858084818086610573565b606483601b604051917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401528201527f4344523a205265616420636f6e646974696f6e206e6f74206d657400000000006044820152fd5b6107ee915060203d6020116107f4575b6107e68183612372565b810190612439565b8661075d565b503d6107dc565b606486601e604051917f08c379a0000000000000000000000000000000000000000000000000000000008352602060048401528201527f4344523a205661756c7420686173206e6f206461746120746f207265616400006044820152fd5b346102c1575f6003193601126102c15760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60554604051908152f35b60a06003193601126102c15760043580151581036102c1576024359073ffffffffffffffffffffffffffffffffffffffff821682036102c1576044359073ffffffffffffffffffffffffffffffffffffffff821682036102c15760643567ffffffffffffffff81116102c15761090f903690600401612305565b919060843567ffffffffffffffff81116102c157610931903690600401612305565b91909361093c61268e565b6109446126e8565b73ffffffffffffffffffffffffffffffffffffffff87161515806111a8575b1561114a577f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf604548034036110ec57805f81156110e3575b5f808093818094f11561066f576040519081525f60208201527f93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f740031860403392a27f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf600549563ffffffff80881610156110855763ffffffff808816146110585763ffffffff600181891601167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000008816177f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf6005560405197610a7689612356565b851515895273ffffffffffffffffffffffffffffffffffffffff811660208a015273ffffffffffffffffffffffffffffffffffffffff821660408a0152610abe368486612395565b60608a0152610ace368689612395565b60808a015260405180602081011067ffffffffffffffff602083011117610efa57602081016040525f815260a08a0152610b3b63ffffffff891663ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60760205260405f2090565b98805115158a547fffffffffffffffffffffff00000000000000000000000000000000000000000060ff74ffffffffffffffffffffffffffffffffffffffff00602086015160081b169316911617178a5560018a0173ffffffffffffffffffffffffffffffffffffffff6040830151167fffffffffffffffffffffffff0000000000000000000000000000000000000000825416179055606081015180519067ffffffffffffffff8211610efa57610c03828d6002610bfc81830154612451565b9101612551565b602090601f8311600114610fcb57610c5092915f9183610fc0575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60028b01555b608081015180519067ffffffffffffffff8211610efa57610c80828d6003610bfc81830154612451565b602090601f8311600114610f32579180610cd09260a095945f92610f275750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60038c01555b01519586519967ffffffffffffffff8b11610efa578a8a98610d0a60209d610d016004860154612451565b60048601612551565b8c90601f8311600114610e2257610dd495610de299957ff1370099e56e061a64615edc07c1d17e36a20585cdc9b288bf0259a528365d0a9d99956004610da463ffffffff9f9e9b96978073ffffffffffffffffffffffffffffffffffffffff998a985f92610e175750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9101555b6040519c8d9c168c5215158f8c01521660408a015216606088015260c0608088015260c0870191612419565b9184830360a0860152612419565b0390a160017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005563ffffffff60405191168152f35b015190505f80610c1e565b90600484015f528d5f20915f5b601f1985168110610ee0575095610de299957ff1370099e56e061a64615edc07c1d17e36a20585cdc9b288bf0259a528365d0a9d99956004600163ffffffff9f9e9b969773ffffffffffffffffffffffffffffffffffffffff988997610dd49d83601f19811610610ea9575b505050811b01910155610da8565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690555f8080610e9b565b8282015184558e9c50600190930192918f01918f01610e2f565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b015190508e80610c1e565b9060038d015f5260205f20915f5b601f1985168110610fa8575091839160019383601f1960a098971610610f71575b505050811b0160038c0155610cd6565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558d8080610f61565b91926020600181928685015181550194019201610f40565b015190508d80610c1e565b919060028d015f5260205f20905f935b601f198416851061103d576001945083601f19811610611006575b505050811b0160028b0155610c56565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558c8080610ff6565b81810151835560209485019460019093019290910190610fdb565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f4344523a205661756c742055554944206f766572666c6f7700000000000000006044820152fd5b506108fc61099a565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4344523a20496e76616c69642066656520616d6f756e740000000000000000006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f496e76616c696420636f6e646974696f6e2061646472657373000000000000006044820152fd5b5073ffffffffffffffffffffffffffffffffffffffff86161515610963565b346102c1575f6003193601126102c15760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60354604051908152f35b346102c15760206003193601126102c15761121c6125a0565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60355005b346102c1575f6003193601126102c157602073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416604051908152f35b346102c1575f6003193601126102c15760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60654604051908152f35b346102c15760e06003193601126102c1576112ea612333565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c16159167ffffffffffffffff8116801590816115d7575b60011490816115cd575b1590816115c4575b5061159a578260017fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008316178555611565575b5061137a61289e565b61138261289e565b73ffffffffffffffffffffffffffffffffffffffff811615611535576113a79061273d565b6113af61289e565b6113b761289e565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00556113e361289e565b6113eb61289e565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0081541690556024357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf601556044357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf602556064357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf603556084357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf604556114d160a4356125e0565b6114dc60c4356127f0565b6114e257005b7fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff81541690557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160018152a1005b60246040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081525f6004820152fd5b7fffffffffffffffffffffffffffffffffffffffffffffff000000000000000000166801000000000000000117835583611371565b60046040517ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b9050158561133e565b303b159150611336565b84915061132c565b346102c1575f6003193601126102c1573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00541603611636576104473361273d565b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b6003196060813601126102c15761167b6122f2565b60249167ffffffffffffffff9083358281116102c15761169f903690600401612305565b92909460449586358381116102c1576116bc903690600401612305565b9490956116c761268e565b6116cf6126e8565b8515611cdc577f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf605548611611c5a576117348863ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60760205260405f2090565b9182549373ffffffffffffffffffffffffffffffffffffffff8560081c16928315611bd857918b8b92859489963303611aca575b50505050505050600461177c910154612451565b611a5f575b507f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf6025495863403611a025750855f81156119f9575b5f808093818094f11561066f57604051958652600190602096600160208201527f93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f740031860403392a260046118358763ffffffff165f527f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60760205260405f2090565b019284116119ce57506118528361184c8454612451565b84612551565b5f95601f841160011461191d575050918163ffffffff94936118cc826118e9957f68b7452742f8d8707af90ddff5ef1a7f8850f5724b804e72b4df1a37050a6355995f91611912575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b604051948594168452604060208501526040840191612419565b0390a160017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055005b90508501358a61189b565b5f838152602081209792601f19861692905b8383106119b757505050917f68b7452742f8d8707af90ddff5ef1a7f8850f5724b804e72b4df1a37050a63559684926118e99563ffffffff9897951061197f575b5050600182811b0190556118cf565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88560031b161c19908501351690558780611970565b878501358a5598890198938101939181019161192f565b7f4e487b71000000000000000000000000000000000000000000000000000000005f5260416004525ffd5b506108fc6117b6565b907f4344523a20496e76616c69642066656520616d6f756e740000000000000000006064926017604051937f08c379a000000000000000000000000000000000000000000000000000000000855260206004860152840152820152fd5b60ff1615611a6d5786611781565b6064907f4344523a205661756c74206973206e6f7420757064617461626c65000000000087601b604051937f08c379a000000000000000000000000000000000000000000000000000000000855260206004860152840152820152fd5b6020959363ffffffff93611b206080999794611b32946040519b8c9a8b998a997f5645dbbf000000000000000000000000000000000000000000000000000000008b521660048a01528801526084870191612419565b918483030190840152600289016124a2565b33606483015203915afa90811561066f575f91611bb9575b5015611b5c57828988818b8180611768565b6064837f4344523a20577269746520636f6e646974696f6e206e6f74206d6574000000008a601c604051937f08c379a000000000000000000000000000000000000000000000000000000000855260206004860152840152820152fd5b611bd2915060203d6020116107f4576107e68183612372565b89611b4a565b6084877f4344523a20577269746520636f6e646974696f6e2061646472657373206e6f748e604051927f08c379a000000000000000000000000000000000000000000000000000000000845260206004850152808401528201527f20736574000000000000000000000000000000000000000000000000000000006064820152fd5b6084847f4344523a20456e6372797074656420646174612065786365656473206d6178208b604051927f08c379a000000000000000000000000000000000000000000000000000000000845260206004850152808401528201527f73697a65000000000000000000000000000000000000000000000000000000006064820152fd5b6084847f4344523a20456e6372797074656420646174612063616e6e6f7420626520656d8b6023604051937f08c379a0000000000000000000000000000000000000000000000000000000008552602060048601528401528201527f70747900000000000000000000000000000000000000000000000000000000006064820152fd5b346102c1575f6003193601126102c15760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60454604051908152f35b346102c1575f6003193601126102c157611db36125a0565b5f73ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b346102c1575f6003193601126102c15760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60154604051908152f35b346102c1575f6003193601126102c15760207f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60254604051908152f35b610120806003193601126102c157611eee6122f2565b6024359063ffffffff821682036102c15767ffffffffffffffff6044358181116102c157611f20903690600401612305565b90926064358381116102c157611f3a903690600401612305565b9590946084358581116102c157611f55903690600401612305565b93909460a4358781116102c157611f70903690600401612305565b919060c4358981116102c157611f8a903690600401612305565b93909460e4359a63ffffffff8c168c036102c157610104359081116102c157611fb7903690600401612305565b9a9099611fc261268e565b611fca6126e8565b851515806121e2575b1561215e577f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf601549e8f34036110ec578f5f80808093818115612155575b8290f11561066f578f9e6120fd9f9161209b6120e29c6120d09a63ffffffff9f9863ffffffff6120be9a6120ad98604051908152600360208201527f93a12ee83f0a200e41b4335defc8068ad4673cbe46a7eb937ecfa7f3f740031860403392a260405160805281610140931660805152166020608051015280604060805101526080510191612419565b91608051830360606080510152612419565b916080518303608080510152612419565b91608051830360a06080510152612419565b91608051830360c06080510152612419565b931660e0608051015260805183036101006080510152612419565b9160805101527fe2d964436c9da9c84292bd1b67e737171e45dbb6ca28f34fcbfb5c9ce20c1dbc33916080519003608051a260017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055005b506108fc612010565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f4344523a20496e76616c696420656e63727970746564207061727469616c206c60448201527f656e6774680000000000000000000000000000000000000000000000000000006064820152fd5b507f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60654861115611fd3565b346102c1575f6003193601126102c157602060ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330054166040519015158152f35b346102c15760206003193601126102c1576122676125a0565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60155005b346102c15760206003193601126102c1576122a76125a0565b6104476004356125e0565b346102c15760206003193601126102c1576122cb6125a0565b6004357f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60255005b6004359063ffffffff821682036102c157565b9181601f840112156102c15782359167ffffffffffffffff83116102c157602083818601950101116102c157565b6004359073ffffffffffffffffffffffffffffffffffffffff821682036102c157565b60c0810190811067ffffffffffffffff821117610efa57604052565b90601f601f19910116810190811067ffffffffffffffff821117610efa57604052565b92919267ffffffffffffffff8211610efa57604051916123bf6020601f19601f8401160184612372565b8294818452818301116102c1578281602093845f960137010152565b91908251928382525f5b848110612405575050601f19601f845f6020809697860101520116010190565b6020818301810151848301820152016123e5565b601f8260209493601f1993818652868601375f8582860101520116010190565b908160209103126102c1575180151581036102c15790565b90600182811c92168015612498575b602083101461246b57565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f1691612460565b80545f93926124b082612451565b918282526020936001916001811690815f1461251457506001146124d6575b5050505050565b90939495505f92919252835f2092845f945b83861061250057505050500101905f808080806124cf565b8054858701830152940193859082016124e8565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168685015250505090151560051b010191505f808080806124cf565b601f821161255e57505050565b5f5260205f20906020601f840160051c83019310612596575b601f0160051c01905b81811061258b575050565b5f8155600101612580565b9091508190612577565b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416330361163657565b801561260a577f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60555565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f4344523a204d617820656e6372797074656420646174612073697a65206d757360448201527f74206265203e20300000000000000000000000000000000000000000000000006064820152fd5b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0060028154146126be5760029055565b60046040517f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b60ff7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300541661271357565b60046040517fd93c0665000000000000000000000000000000000000000000000000000000008152fd5b7fffffffffffffffffffffffff0000000000000000000000000000000000000000907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008281541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549073ffffffffffffffffffffffffffffffffffffffff80931680948316179055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3565b801561281a577f38eb98a52971d8773d43e336c762a70f1492f62ea143e494a29d8ec99eadf60655565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f4344523a204d617820656e63727970746564207061727469616c2073697a652060448201527f6d757374206265203e20300000000000000000000000000000000000000000006064820152fd5b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c16156128cd57565b60046040517fd7e6bcf8000000000000000000000000000000000000000000000000000000008152fdfea2646970667358221220fffdf11ac6b3488ac630f81c965c2fb06ed3b41d3589feceb05ba2c287e59e9064736f6c63430008170033",
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

// MaxEncryptedDataSize is a free data retrieval call binding the contract method 0xa060fa12.
//
// Solidity: function maxEncryptedDataSize() view returns(uint256)
func (_CDR *CDRCaller) MaxEncryptedDataSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "maxEncryptedDataSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxEncryptedDataSize is a free data retrieval call binding the contract method 0xa060fa12.
//
// Solidity: function maxEncryptedDataSize() view returns(uint256)
func (_CDR *CDRSession) MaxEncryptedDataSize() (*big.Int, error) {
	return _CDR.Contract.MaxEncryptedDataSize(&_CDR.CallOpts)
}

// MaxEncryptedDataSize is a free data retrieval call binding the contract method 0xa060fa12.
//
// Solidity: function maxEncryptedDataSize() view returns(uint256)
func (_CDR *CDRCallerSession) MaxEncryptedDataSize() (*big.Int, error) {
	return _CDR.Contract.MaxEncryptedDataSize(&_CDR.CallOpts)
}

// MaxEncryptedPartialSize is a free data retrieval call binding the contract method 0x8465ff72.
//
// Solidity: function maxEncryptedPartialSize() view returns(uint256)
func (_CDR *CDRCaller) MaxEncryptedPartialSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CDR.contract.Call(opts, &out, "maxEncryptedPartialSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxEncryptedPartialSize is a free data retrieval call binding the contract method 0x8465ff72.
//
// Solidity: function maxEncryptedPartialSize() view returns(uint256)
func (_CDR *CDRSession) MaxEncryptedPartialSize() (*big.Int, error) {
	return _CDR.Contract.MaxEncryptedPartialSize(&_CDR.CallOpts)
}

// MaxEncryptedPartialSize is a free data retrieval call binding the contract method 0x8465ff72.
//
// Solidity: function maxEncryptedPartialSize() view returns(uint256)
func (_CDR *CDRCallerSession) MaxEncryptedPartialSize() (*big.Int, error) {
	return _CDR.Contract.MaxEncryptedPartialSize(&_CDR.CallOpts)
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

// Initialize is a paid mutator transaction binding the contract method 0x8142951a.
//
// Solidity: function initialize(address owner, uint256 baseFee, uint256 writeFee, uint256 readFee, uint256 allocateFee, uint256 maxEncryptedDataSize, uint256 maxEncryptedPartialSize) returns()
func (_CDR *CDRTransactor) Initialize(opts *bind.TransactOpts, owner common.Address, baseFee *big.Int, writeFee *big.Int, readFee *big.Int, allocateFee *big.Int, maxEncryptedDataSize *big.Int, maxEncryptedPartialSize *big.Int) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "initialize", owner, baseFee, writeFee, readFee, allocateFee, maxEncryptedDataSize, maxEncryptedPartialSize)
}

// Initialize is a paid mutator transaction binding the contract method 0x8142951a.
//
// Solidity: function initialize(address owner, uint256 baseFee, uint256 writeFee, uint256 readFee, uint256 allocateFee, uint256 maxEncryptedDataSize, uint256 maxEncryptedPartialSize) returns()
func (_CDR *CDRSession) Initialize(owner common.Address, baseFee *big.Int, writeFee *big.Int, readFee *big.Int, allocateFee *big.Int, maxEncryptedDataSize *big.Int, maxEncryptedPartialSize *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.Initialize(&_CDR.TransactOpts, owner, baseFee, writeFee, readFee, allocateFee, maxEncryptedDataSize, maxEncryptedPartialSize)
}

// Initialize is a paid mutator transaction binding the contract method 0x8142951a.
//
// Solidity: function initialize(address owner, uint256 baseFee, uint256 writeFee, uint256 readFee, uint256 allocateFee, uint256 maxEncryptedDataSize, uint256 maxEncryptedPartialSize) returns()
func (_CDR *CDRTransactorSession) Initialize(owner common.Address, baseFee *big.Int, writeFee *big.Int, readFee *big.Int, allocateFee *big.Int, maxEncryptedDataSize *big.Int, maxEncryptedPartialSize *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.Initialize(&_CDR.TransactOpts, owner, baseFee, writeFee, readFee, allocateFee, maxEncryptedDataSize, maxEncryptedPartialSize)
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

// SetMaxEncryptedDataSize is a paid mutator transaction binding the contract method 0x1cf9c862.
//
// Solidity: function setMaxEncryptedDataSize(uint256 newMaxEncryptedDataSize) returns()
func (_CDR *CDRTransactor) SetMaxEncryptedDataSize(opts *bind.TransactOpts, newMaxEncryptedDataSize *big.Int) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "setMaxEncryptedDataSize", newMaxEncryptedDataSize)
}

// SetMaxEncryptedDataSize is a paid mutator transaction binding the contract method 0x1cf9c862.
//
// Solidity: function setMaxEncryptedDataSize(uint256 newMaxEncryptedDataSize) returns()
func (_CDR *CDRSession) SetMaxEncryptedDataSize(newMaxEncryptedDataSize *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.SetMaxEncryptedDataSize(&_CDR.TransactOpts, newMaxEncryptedDataSize)
}

// SetMaxEncryptedDataSize is a paid mutator transaction binding the contract method 0x1cf9c862.
//
// Solidity: function setMaxEncryptedDataSize(uint256 newMaxEncryptedDataSize) returns()
func (_CDR *CDRTransactorSession) SetMaxEncryptedDataSize(newMaxEncryptedDataSize *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.SetMaxEncryptedDataSize(&_CDR.TransactOpts, newMaxEncryptedDataSize)
}

// SetMaxEncryptedPartialSize is a paid mutator transaction binding the contract method 0xdcaa51d4.
//
// Solidity: function setMaxEncryptedPartialSize(uint256 newMaxEncryptedPartialSize) returns()
func (_CDR *CDRTransactor) SetMaxEncryptedPartialSize(opts *bind.TransactOpts, newMaxEncryptedPartialSize *big.Int) (*types.Transaction, error) {
	return _CDR.contract.Transact(opts, "setMaxEncryptedPartialSize", newMaxEncryptedPartialSize)
}

// SetMaxEncryptedPartialSize is a paid mutator transaction binding the contract method 0xdcaa51d4.
//
// Solidity: function setMaxEncryptedPartialSize(uint256 newMaxEncryptedPartialSize) returns()
func (_CDR *CDRSession) SetMaxEncryptedPartialSize(newMaxEncryptedPartialSize *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.SetMaxEncryptedPartialSize(&_CDR.TransactOpts, newMaxEncryptedPartialSize)
}

// SetMaxEncryptedPartialSize is a paid mutator transaction binding the contract method 0xdcaa51d4.
//
// Solidity: function setMaxEncryptedPartialSize(uint256 newMaxEncryptedPartialSize) returns()
func (_CDR *CDRTransactorSession) SetMaxEncryptedPartialSize(newMaxEncryptedPartialSize *big.Int) (*types.Transaction, error) {
	return _CDR.Contract.SetMaxEncryptedPartialSize(&_CDR.TransactOpts, newMaxEncryptedPartialSize)
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
