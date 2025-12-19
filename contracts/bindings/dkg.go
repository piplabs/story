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

// IDKGNodeInfo is an auto generated low-level Go binding around an user-defined struct.
type IDKGNodeInfo struct {
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	NodeStatus uint8
}

// DKGMetaData contains all meta data concerning the DKG contract.
var DKGMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"complainDeals\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"curMrenclave\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dealComplaints\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainant\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dkgNodeInfos\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"nodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.NodeStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getNodeInfo\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDKG.NodeInfo\",\"components\":[{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"nodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.NodeStatus\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"partialDecrypts\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"labelHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"partialDecryption\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestRemoteAttestationOnChain\",\"inputs\":[{\"name\":\"targetValidatorAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestThresholdDecryption\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNetwork\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitPartialDecryption\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ephemeralPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"DKGFinalized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGInitialized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGNetworkSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"total\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealComplaintsSubmitted\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"indexed\":false,\"internalType\":\"uint32[]\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealVerified\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"recipientIndex\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidDeal\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PartialDecryptionSubmitted\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ephemeralPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoteAttestationProcessedOnChain\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ThresholdDecryptRequested\",\"inputs\":[{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeScheduled\",\"inputs\":[{\"name\":\"activationHeight\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60803461005057601f612b1038819003918201601f19168301916001600160401b038311848410176100545780849260209460405283398101031261005057515f55604051612aa790816100698239f35b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c806308ad63ac146100d4578063227cd922146100cf57806367295350146100ca578063681a0fa8146100c55780637fb93a3a146100c0578063a26f51a4146100bb578063aab066c6146100b6578063b1133cac146100b1578063b1888cd3146100ac578063dea942d9146100a7578063fa4e9f63146100a25763fdbc551f1461009d575f80fd5b6112f7565b6110c3565b610ff4565b610dbc565b610c5e565b610b05565b61082a565b61072b565b6106f1565b610541565b6102f7565b6101fd565b6004359063ffffffff821682036100ec57565b5f80fd5b6024359063ffffffff821682036100ec57565b6044359063ffffffff821682036100ec57565b359063ffffffff821682036100ec57565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b60a0810190811067ffffffffffffffff82111761017057604052565b610127565b6040810190811067ffffffffffffffff82111761017057604052565b6060810190811067ffffffffffffffff82111761017057604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761017057604052565b604051906101fb82610154565b565b346100ec5760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec576102346100d9565b61023c6100f0565b906044359167ffffffffffffffff908184116100ec57366023850112156100ec578360040135918211610170578160051b936020946040519361028260208301866101ad565b845260246020850191830101913683116100ec57602401905b8282106102b2576102b0606435868689611448565b005b8680916102be84610116565b81520191019061029b565b9181601f840112156100ec5782359167ffffffffffffffff83116100ec57602083818601950101116100ec57565b346100ec5760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec5761032e6100d9565b6103366100f0565b9061033f610103565b6064359160843567ffffffffffffffff81116100ec576103639036906004016102c9565b90604080516020908181019088825282815261037e81610175565b5190205f54835183810191825283815261039781610175565b519020146103a4906113e3565b5f8781526001825282812063ffffffff87168252602090815260408083203384529091529020906103d760018301610994565b8351828101908a82527fffffffff00000000000000000000000000000000000000000000000000000000808a60e01b1687830152808d60e01b1660448301528a60e01b166048820152602c815261042d81610191565b519020610461907f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f52601c52603c5f2090565b61046c368888611ce5565b61047591612803565b8151919092012073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166104de9173ffffffffffffffffffffffffffffffffffffffff161461160e565b60030180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661020017905551948594339761051b95876116ed565b037fc7a37268197965e156b6d53085e9e20ba69f731868b09d00c2b2c3925f25f4f891a2005b346100ec5760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec576105786100d9565b60243567ffffffffffffffff6044358181116100ec5761059c9036906004016102c9565b9290916064359081116100ec576105b79036906004016102c9565b6105f76040516020810190858252602081526105d281610175565b5190205f5460405160208101918252602081526105ee81610175565b519020146113e3565b825f52600160205261064061061d8760405f209063ffffffff165f5260205260405f2090565b3373ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b91600383019560ff8754169360038510156106ec577fe7419c96e4837a0c8c3c13342ccea4095f978269a67a3ce5dcfeac5664240205976106ad6106a88686868c8f8d906106a360016106e79f9b61069d826106d99e1415611724565b01610994565b6124ae565b611789565b6103007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff825416179055565b6040519586953399876117ee565b0390a2005b610ab7565b346100ec575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec5760205f54604051908152f35b346100ec5760c07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec576107626100d9565b67ffffffffffffffff906044358281116100ec576107849036906004016102c9565b6064929192358481116100ec5761079f9036906004016102c9565b6084929192358681116100ec576107ba9036906004016102c9565b93909260a4359788116100ec576107d86102b09836906004016102c9565b97909660243590611827565b6064359073ffffffffffffffffffffffffffffffffffffffff821682036100ec57565b6004359073ffffffffffffffffffffffffffffffffffffffff821682036100ec57565b346100ec5760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec57602060ff6108d56108686100f0565b6108b2610873610103565b61089f61087e6107e4565b936004355f526002885260405f209063ffffffff165f5260205260405f2090565b9063ffffffff165f5260205260405f2090565b9073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b54166040519015158152f35b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc60609101126100ec576004359060243563ffffffff811681036100ec579060443573ffffffffffffffffffffffffffffffffffffffff811681036100ec5790565b90600182811c9216801561098a575b602083101461095d57565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f1691610952565b9060405191825f82546109a681610943565b908184526020946001916001811690815f14610a1257506001146109d4575b5050506101fb925003836101ad565b5f90815285812095935091905b8183106109fa5750506101fb93508201015f80806109c5565b855488840185015294850194879450918301916109e1565b9150506101fb9593507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201015f80806109c5565b5f5b838110610a645750505f910152565b8181015183820152602001610a55565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f602093610ab081518092818752878088019101610a53565b0116010190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b600311156106ec57565b9060038210156106ec5752565b600411156106ec57565b346100ec57610b89610b3c6108b2610b1c366108e1565b92915f52600160205260405f209063ffffffff165f5260205260405f2090565b610b4581610994565b90610b5260018201610994565b610bb3610ba56003610b6660028601610994565b94015493610b9760ff8660081c169460405198899860a08a5260a08a0190610a74565b9088820360208a0152610a74565b908682036040880152610a74565b9260ff606086019116610aee565b610bbc81610afb565b60808301520390f35b9060a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8301126100ec5760043563ffffffff811681036100ec57916024359167ffffffffffffffff6044358181116100ec5783610c25916004016102c9565b939093926064358381116100ec5782610c40916004016102c9565b939093926084359182116100ec57610c5a916004016102c9565b9091565b346100ec57610c6c36610bc5565b969593604093919351610ca860209182810190888252838152610c8e81610175565b5190205f546040518481019182528481526105ee81610175565b610cb963ffffffff89161515611a7d565b8315610d5e5760418203610d005750916106e793917f5f9d4b68667a3b91f2fe8369c2ac9040ce4a68400aadbe076f9d109b13e09b61979893604051978897339b8961216b565b606490604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601860248201527f496e76616c696420726571756573746572207075626b657900000000000000006044820152fd5b606490604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601060248201527f456d7074792063697068657274657874000000000000000000000000000000006044820152fd5b346100ec57610dca36610bc5565b9260409794979291925197610e086020998a8101908a82528b8152610dee81610175565b5190205f546040518c81019182528c81526105ee81610175565b610e35610e16368787611ce5565b610e2136848b611ce5565b88610e2d368888611ce5565b923390612717565b15610f0457906106e794939291610ef6899a610e94610e8e610e897f1bd0faa06edbfccdd0f51f46517f5bae23b4abca2dad81e938e89f3ddf7cab1d9d6108b28d61089f33935f52600160205260405f2090565b611c76565b156121b6565b610e9c6101ee565b90610ea836858d611ce5565b8252610eb5368787611ce5565b90820152610ec4368888611ce5565b60408201525f606082015260016080820152610ef18c6108b28b61089f33935f52600160205260405f2090565b61225e565b604051978897339b896123d9565b606489604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601a60248201527f496e76616c69642072656d6f7465206174746573746174696f6e0000000000006044820152fd5b6020815260a06080610fd1610f82855184602087015260c0860190610a74565b610fbc6020870151917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09283888303016040890152610a74565b90604087015190868303016060870152610a74565b93610fe3606082015183860190610aee565b015191610fef83610afb565b015290565b346100ec576110bf6108b261105561100b366108e1565b9193906040945f6080875161101f81610154565b606081526060602082015260608982015282606082015201525f526001602052845f209063ffffffff165f5260205260405f2090565b9060ff600382519361106685610154565b61106f81610994565b855261107d60018201610994565b602086015261108e60028201610994565b8486015201546110a38282166060860161221b565b60081c166110b081610afb565b60808301525191829182610f62565b0390f35b346100ec5760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec576110fa610807565b6111026100f0565b6044359160405161112260209182810190868252838152610c8e81610175565b835f5260018152611148826108b28560405f209063ffffffff165f5260205260405f2090565b906111538254610943565b1561124757507f54690f98c0ec0056e0e487f4fe5e8eea7bee88d2dbb7cc9ddca22981f06d9dbb93611207826111d160036112149501916111a7611198845460ff1690565b6111a181610ae4565b15612408565b6111b360028201610994565b9088886111cb60016111c485610994565b9401610994565b93612717565b156112195780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660021781555b5460ff1690565b936040519485948561246d565b0390a1005b80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001178155611200565b606490604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601360248201527f4e6f646520646f6573206e6f74206578697374000000000000000000000000006044820152fd5b9390608093610b976112ef9473ffffffffffffffffffffffffffffffffffffffff6112e1949a999a16885260a0602089015260a0880190610a74565b908482036060860152610a74565b931515910152565b346100ec5760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec5761138d6113316100f0565b61135e61133c6107e4565b916004355f52600360205260405f209063ffffffff165f5260205260405f2090565b6044355f5260205260405f209073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b73ffffffffffffffffffffffffffffffffffffffff8154166110bf6113b460018401610994565b926113c160028201610994565b9060ff60046113d260038401610994565b9201541691604051958695866112a5565b156113ea57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f496e76616c6964206d72656e636c6176650000000000000000000000000000006044820152fd5b93929360408051906114836020928381019089825284815261146981610175565b5190205f546040518581019182528581526105ee81610175565b865f526001918291600182526114ad61061d8660405f209063ffffffff165f5260205260405f2090565b505f925b6114f2575b505050507fa89000d88bdc9c3e92c10abb67235241f8c6803723e88e1e2420533e8fe2b8d893946114ed91604051948594856115b3565b0390a1565b86518310156115ae575f8981526002835281812063ffffffff8716825260205260409020875184101561158157849361157a61154f8a6108b2889563ffffffff8933948860051b0101511663ffffffff165f5260205260405f2090565b60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00825416179055565b01926114b1565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b6114b6565b91949392946080830163ffffffff8093168452602090608060208601528251809152602060a086019301915f5b8281106115f857505050509416604082015260600152565b83518616855293810193928101926001016115e0565b1561161557565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f496e76616c696420736574206e6574776f726b207369676e61747572650000006044820152fd5b9061167d81610afb565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff0083549260081b169116179055565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe093818652868601375f8582860101520116010190565b9160a09361172197959263ffffffff92838092168652166020850152166040830152606082015281608082015201916116af565b90565b1561172b57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4e6f64652077617320696e76616c6964617465640000000000000000000000006044820152fd5b1561179057565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f496e76616c69642066696e616c697a6174696f6e207369676e617475726500006044820152fd5b949290936118199263ffffffff611721989616875260208701526080604087015260808601916116af565b9260608185039101526116af565b95919692979398949098604051602090818101908c825282815261184a81610175565b5190205f5460405183810191825283815261186481610175565b51902014611871906113e3565b61188263ffffffff89161515611a7d565b61188d8a1515611ae2565b611898851515611b47565b6118a3871515611bac565b6118af60418414611c11565b876118c28c5f52600160205260405f2090565b906118da919063ffffffff165f5260205260405f2090565b33611902919073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b61190b90611c76565b151561191690611c80565b8a611922368989611ce5565b82815191012033818b61193d855f52600360205260405f2090565b90611955919063ffffffff165f5260205260405f2090565b5f9182526020526040902090611988919073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b6004015460ff161561199990611d49565b6119a16101ee565b338152926119b0368e8e611ce5565b908401526119bf368888611ce5565b60408401526119cf368a8a611ce5565b6060840152600160808401528933926119f0905f52600360205260405f2090565b90611a08919063ffffffff165f5260205260405f2090565b5f9182526020526040902090611a3b919073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b90611a4591611f28565b604051988998339b611a57998b61210c565b037f2d05ab82940dc8b5196f5969dc54167bb42e65a5018b12a99c8a7f20e1dc246491a2565b15611a8457565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f496e76616c696420726f756e64000000000000000000000000000000000000006044820152fd5b15611ae957565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f456d707479207061727469616c000000000000000000000000000000000000006044820152fd5b15611b4e57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f456d7074792070756253686172650000000000000000000000000000000000006044820152fd5b15611bb357565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f456d707479206c6162656c0000000000000000000000000000000000000000006044820152fd5b15611c1857565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f496e76616c696420657068656d6572616c207075626b657900000000000000006044820152fd5b6117219054610943565b15611c8757565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f53656e646572206e6f74207265676973746572656400000000000000000000006044820152fd5b92919267ffffffffffffffff82116101705760405191611d2d60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601846101ad565b8294818452818301116100ec578281602093845f960137010152565b15611d5057565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f5061727469616c20616c7265616479207375626d6974746564000000000000006044820152fd5b601f8211611dbb57505050565b5f5260205f20906020601f840160051c83019310611df3575b601f0160051c01905b818110611de8575050565b5f8155600101611ddd565b9091508190611dd4565b919091825167ffffffffffffffff811161017057611e2581611e1f8454610943565b84611dae565b602080601f8311600114611e8457508190611e759394955f92611e79575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9055565b015190505f80611e43565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0831695611eb6855f5260205f2090565b925f905b888210611f1057505083600195969710611ed9575b505050811b019055565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690555f8080611ecf565b80600185968294968601518155019501930190611eba565b9073ffffffffffffffffffffffffffffffffffffffff8151167fffffffffffffffffffffffff000000000000000000000000000000000000000083541617825560018083019060208084015180519267ffffffffffffffff841161017057611f9a84611f948754610943565b87611dae565b602092601f851160011461205457505093600493611ff684612022956080956101fb9a995f92611e795750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b61200a604082015160028701611dfd565b61201b606082015160038701611dfd565b0151151590565b91019060ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0083541691151516179055565b9291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0851690612088875f5260205f2090565b945f915b8383106120f557505050846080946101fb9998946004989461202298600195106120be575b505050811b019055611ff9565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690555f80806120b1565b84860151875595860195948101949181019161208c565b9896949161214f936117219b99956121419263ffffffff61215d9a95168d5260208d015260c060408d015260c08c01916116af565b9189830360608b01526116af565b9186830360808801526116af565b9260a08185039101526116af565b9694926121a89461172199979363ffffffff61219a94168a5260208a015260a060408a015260a08901916116af565b9186830360608801526116af565b9260808185039101526116af565b156121bd57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f416c7265616479207265676973746572656400000000000000000000000000006044820152fd5b60038210156106ec5752565b9060038110156106ec5760ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008354169116179055565b9190805192835167ffffffffffffffff81116101705761228281611e1f8454610943565b602080601f83116001146123285750916122d9826003936080956101fb98995f92611e795750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b81555b6122ed602085015160018301611dfd565b6122fe604085015160028301611dfd565b0191612317606082015161231181610ae4565b84612227565b01519061232382610afb565b611673565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe083169661235a855f5260205f2090565b925f905b8982106123c15750509260809492600192600395836101fb9a9b1061238b575b505050811b0181556122dc565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f884881b161c191690555f808061237e565b8060018596829496860151815501950193019061235e565b9694926121a89463ffffffff6117219a989461219a948b521660208a015260a060408a015260a08901916116af565b1561240f57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4e6f646520616c7265616479206368616c6c656e6765640000000000000000006044820152fd5b9094939263ffffffff906124a460609473ffffffffffffffffffffffffffffffffffffffff60808601991685526020850190610aee565b1660408201520152565b949061254f927fffffffff00000000000000000000000000000000000000000000000000000000612517604461254794896125559a61255e9d9a60405196879460208601998a5260e01b1660408501528484013781015f838201520360248101845201826101ad565b5190207f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f52601c52603c5f2090565b923691611ce5565b906128c6565b90939193612900565b73ffffffffffffffffffffffffffffffffffffffff8160208293519101201691161490565b1561258a57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f496e76616c69642076616c696461746f722061646472657373000000000000006044820152fd5b156125ef57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f496e76616c696420726f756e64206f6620444b470000000000000000000000006044820152fd5b1561265457565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f496e76616c696420444b47207075626c6963206b6579000000000000000000006044820152fd5b156126b957565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602060248201527f496e76616c696420636f6d6d756e69636174696f6e207075626c6963206b65796044820152fd5b939192909260408551111561277f576101906117219561274e73ffffffffffffffffffffffffffffffffffffffff87161515612583565b61275f63ffffffff841615156125e8565b61276b8451151561264d565b612777855115156126b2565b015193612819565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f496e76616c6964207261772071756f74652c2071756f746520746f6f2073686f60448201527f72740000000000000000000000000000000000000000000000000000000000006064820152fd5b61172191612810916128c6565b90929192612900565b926128b5916038917fffffffff00000000000000000000000000000000000000000000000000000000946040519586937fffffffffffffffffffffffffffffffffffffffff000000000000000000000000602086019960601b16895260e01b1660348401526128918151809260208787019101610a53565b82016128a68251809360208785019101610a53565b010360188101845201826101ad565b519020036128c257600190565b5f90565b81519190604183036128f6576128ef9250602082015190606060408401519301515f1a906129d7565b9192909190565b50505f9160029190565b61290981610afb565b80612912575050565b61291b81610afb565b6001810361294d5760046040517ff645eedf000000000000000000000000000000000000000000000000000000008152fd5b61295681610afb565b60028103612990576040517ffce698f700000000000000000000000000000000000000000000000000000000815260048101839052602490fd5b8061299c600392610afb565b146129a45750565b6040517fd78bce0c0000000000000000000000000000000000000000000000000000000081526004810191909152602490fd5b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411612a66579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa15612a5b575f5173ffffffffffffffffffffffffffffffffffffffff811615612a5157905f905f90565b505f906001905f90565b6040513d5f823e3d90fd5b5050505f916003919056fea264697066735822122025c837d660c4c0b57cd5e77e4ec4c6af765b9f71b7de2737c1f05b300e1617fb64736f6c63430008170033",
}

// DKGABI is the input ABI used to generate the binding from.
// Deprecated: Use DKGMetaData.ABI instead.
var DKGABI = DKGMetaData.ABI

// DKGBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DKGMetaData.Bin instead.
var DKGBin = DKGMetaData.Bin

// DeployDKG deploys a new Ethereum contract, binding an instance of DKG to it.
func DeployDKG(auth *bind.TransactOpts, backend bind.ContractBackend, mrenclave [32]byte) (common.Address, *types.Transaction, *DKG, error) {
	parsed, err := DKGMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DKGBin), backend, mrenclave)
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

// CurMrenclave is a free data retrieval call binding the contract method 0x681a0fa8.
//
// Solidity: function curMrenclave() view returns(bytes32)
func (_DKG *DKGCaller) CurMrenclave(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "curMrenclave")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CurMrenclave is a free data retrieval call binding the contract method 0x681a0fa8.
//
// Solidity: function curMrenclave() view returns(bytes32)
func (_DKG *DKGSession) CurMrenclave() ([32]byte, error) {
	return _DKG.Contract.CurMrenclave(&_DKG.CallOpts)
}

// CurMrenclave is a free data retrieval call binding the contract method 0x681a0fa8.
//
// Solidity: function curMrenclave() view returns(bytes32)
func (_DKG *DKGCallerSession) CurMrenclave() ([32]byte, error) {
	return _DKG.Contract.CurMrenclave(&_DKG.CallOpts)
}

// DealComplaints is a free data retrieval call binding the contract method 0xa26f51a4.
//
// Solidity: function dealComplaints(bytes32 mrenclave, uint32 round, uint32 index, address complainant) view returns(bool)
func (_DKG *DKGCaller) DealComplaints(opts *bind.CallOpts, mrenclave [32]byte, round uint32, index uint32, complainant common.Address) (bool, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "dealComplaints", mrenclave, round, index, complainant)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DealComplaints is a free data retrieval call binding the contract method 0xa26f51a4.
//
// Solidity: function dealComplaints(bytes32 mrenclave, uint32 round, uint32 index, address complainant) view returns(bool)
func (_DKG *DKGSession) DealComplaints(mrenclave [32]byte, round uint32, index uint32, complainant common.Address) (bool, error) {
	return _DKG.Contract.DealComplaints(&_DKG.CallOpts, mrenclave, round, index, complainant)
}

// DealComplaints is a free data retrieval call binding the contract method 0xa26f51a4.
//
// Solidity: function dealComplaints(bytes32 mrenclave, uint32 round, uint32 index, address complainant) view returns(bool)
func (_DKG *DKGCallerSession) DealComplaints(mrenclave [32]byte, round uint32, index uint32, complainant common.Address) (bool, error) {
	return _DKG.Contract.DealComplaints(&_DKG.CallOpts, mrenclave, round, index, complainant)
}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xaab066c6.
//
// Solidity: function dkgNodeInfos(bytes32 mrenclave, uint32 round, address validator) view returns(bytes dkgPubKey, bytes commPubKey, bytes rawQuote, uint8 chalStatus, uint8 nodeStatus)
func (_DKG *DKGCaller) DkgNodeInfos(opts *bind.CallOpts, mrenclave [32]byte, round uint32, validator common.Address) (struct {
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	NodeStatus uint8
}, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "dkgNodeInfos", mrenclave, round, validator)

	outstruct := new(struct {
		DkgPubKey  []byte
		CommPubKey []byte
		RawQuote   []byte
		ChalStatus uint8
		NodeStatus uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DkgPubKey = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.CommPubKey = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.RawQuote = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.ChalStatus = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.NodeStatus = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xaab066c6.
//
// Solidity: function dkgNodeInfos(bytes32 mrenclave, uint32 round, address validator) view returns(bytes dkgPubKey, bytes commPubKey, bytes rawQuote, uint8 chalStatus, uint8 nodeStatus)
func (_DKG *DKGSession) DkgNodeInfos(mrenclave [32]byte, round uint32, validator common.Address) (struct {
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	NodeStatus uint8
}, error) {
	return _DKG.Contract.DkgNodeInfos(&_DKG.CallOpts, mrenclave, round, validator)
}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xaab066c6.
//
// Solidity: function dkgNodeInfos(bytes32 mrenclave, uint32 round, address validator) view returns(bytes dkgPubKey, bytes commPubKey, bytes rawQuote, uint8 chalStatus, uint8 nodeStatus)
func (_DKG *DKGCallerSession) DkgNodeInfos(mrenclave [32]byte, round uint32, validator common.Address) (struct {
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	NodeStatus uint8
}, error) {
	return _DKG.Contract.DkgNodeInfos(&_DKG.CallOpts, mrenclave, round, validator)
}

// GetNodeInfo is a free data retrieval call binding the contract method 0xdea942d9.
//
// Solidity: function getNodeInfo(bytes32 mrenclave, uint32 round, address validator) view returns((bytes,bytes,bytes,uint8,uint8))
func (_DKG *DKGCaller) GetNodeInfo(opts *bind.CallOpts, mrenclave [32]byte, round uint32, validator common.Address) (IDKGNodeInfo, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "getNodeInfo", mrenclave, round, validator)

	if err != nil {
		return *new(IDKGNodeInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IDKGNodeInfo)).(*IDKGNodeInfo)

	return out0, err

}

// GetNodeInfo is a free data retrieval call binding the contract method 0xdea942d9.
//
// Solidity: function getNodeInfo(bytes32 mrenclave, uint32 round, address validator) view returns((bytes,bytes,bytes,uint8,uint8))
func (_DKG *DKGSession) GetNodeInfo(mrenclave [32]byte, round uint32, validator common.Address) (IDKGNodeInfo, error) {
	return _DKG.Contract.GetNodeInfo(&_DKG.CallOpts, mrenclave, round, validator)
}

// GetNodeInfo is a free data retrieval call binding the contract method 0xdea942d9.
//
// Solidity: function getNodeInfo(bytes32 mrenclave, uint32 round, address validator) view returns((bytes,bytes,bytes,uint8,uint8))
func (_DKG *DKGCallerSession) GetNodeInfo(mrenclave [32]byte, round uint32, validator common.Address) (IDKGNodeInfo, error) {
	return _DKG.Contract.GetNodeInfo(&_DKG.CallOpts, mrenclave, round, validator)
}

// PartialDecrypts is a free data retrieval call binding the contract method 0xfdbc551f.
//
// Solidity: function partialDecrypts(bytes32 mrenclave, uint32 round, bytes32 labelHash, address validator) view returns(address validator, bytes partialDecryption, bytes pubShare, bytes label, bool exists)
func (_DKG *DKGCaller) PartialDecrypts(opts *bind.CallOpts, mrenclave [32]byte, round uint32, labelHash [32]byte, validator common.Address) (struct {
	Validator         common.Address
	PartialDecryption []byte
	PubShare          []byte
	Label             []byte
	Exists            bool
}, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "partialDecrypts", mrenclave, round, labelHash, validator)

	outstruct := new(struct {
		Validator         common.Address
		PartialDecryption []byte
		PubShare          []byte
		Label             []byte
		Exists            bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.PartialDecryption = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.PubShare = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.Label = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.Exists = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// PartialDecrypts is a free data retrieval call binding the contract method 0xfdbc551f.
//
// Solidity: function partialDecrypts(bytes32 mrenclave, uint32 round, bytes32 labelHash, address validator) view returns(address validator, bytes partialDecryption, bytes pubShare, bytes label, bool exists)
func (_DKG *DKGSession) PartialDecrypts(mrenclave [32]byte, round uint32, labelHash [32]byte, validator common.Address) (struct {
	Validator         common.Address
	PartialDecryption []byte
	PubShare          []byte
	Label             []byte
	Exists            bool
}, error) {
	return _DKG.Contract.PartialDecrypts(&_DKG.CallOpts, mrenclave, round, labelHash, validator)
}

// PartialDecrypts is a free data retrieval call binding the contract method 0xfdbc551f.
//
// Solidity: function partialDecrypts(bytes32 mrenclave, uint32 round, bytes32 labelHash, address validator) view returns(address validator, bytes partialDecryption, bytes pubShare, bytes label, bool exists)
func (_DKG *DKGCallerSession) PartialDecrypts(mrenclave [32]byte, round uint32, labelHash [32]byte, validator common.Address) (struct {
	Validator         common.Address
	PartialDecryption []byte
	PubShare          []byte
	Label             []byte
	Exists            bool
}, error) {
	return _DKG.Contract.PartialDecrypts(&_DKG.CallOpts, mrenclave, round, labelHash, validator)
}

// ComplainDeals is a paid mutator transaction binding the contract method 0x08ad63ac.
//
// Solidity: function complainDeals(uint32 round, uint32 index, uint32[] complainIndexes, bytes32 mrenclave) returns()
func (_DKG *DKGTransactor) ComplainDeals(opts *bind.TransactOpts, round uint32, index uint32, complainIndexes []uint32, mrenclave [32]byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "complainDeals", round, index, complainIndexes, mrenclave)
}

// ComplainDeals is a paid mutator transaction binding the contract method 0x08ad63ac.
//
// Solidity: function complainDeals(uint32 round, uint32 index, uint32[] complainIndexes, bytes32 mrenclave) returns()
func (_DKG *DKGSession) ComplainDeals(round uint32, index uint32, complainIndexes []uint32, mrenclave [32]byte) (*types.Transaction, error) {
	return _DKG.Contract.ComplainDeals(&_DKG.TransactOpts, round, index, complainIndexes, mrenclave)
}

// ComplainDeals is a paid mutator transaction binding the contract method 0x08ad63ac.
//
// Solidity: function complainDeals(uint32 round, uint32 index, uint32[] complainIndexes, bytes32 mrenclave) returns()
func (_DKG *DKGTransactorSession) ComplainDeals(round uint32, index uint32, complainIndexes []uint32, mrenclave [32]byte) (*types.Transaction, error) {
	return _DKG.Contract.ComplainDeals(&_DKG.TransactOpts, round, index, complainIndexes, mrenclave)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x67295350.
//
// Solidity: function finalizeDKG(uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes signature) returns()
func (_DKG *DKGTransactor) FinalizeDKG(opts *bind.TransactOpts, round uint32, mrenclave [32]byte, globalPubKey []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "finalizeDKG", round, mrenclave, globalPubKey, signature)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x67295350.
//
// Solidity: function finalizeDKG(uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes signature) returns()
func (_DKG *DKGSession) FinalizeDKG(round uint32, mrenclave [32]byte, globalPubKey []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.FinalizeDKG(&_DKG.TransactOpts, round, mrenclave, globalPubKey, signature)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x67295350.
//
// Solidity: function finalizeDKG(uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes signature) returns()
func (_DKG *DKGTransactorSession) FinalizeDKG(round uint32, mrenclave [32]byte, globalPubKey []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.FinalizeDKG(&_DKG.TransactOpts, round, mrenclave, globalPubKey, signature)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0xb1888cd3.
//
// Solidity: function initializeDKG(uint32 round, bytes32 mrenclave, bytes dkgPubKey, bytes commPubKey, bytes rawQuote) returns()
func (_DKG *DKGTransactor) InitializeDKG(opts *bind.TransactOpts, round uint32, mrenclave [32]byte, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "initializeDKG", round, mrenclave, dkgPubKey, commPubKey, rawQuote)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0xb1888cd3.
//
// Solidity: function initializeDKG(uint32 round, bytes32 mrenclave, bytes dkgPubKey, bytes commPubKey, bytes rawQuote) returns()
func (_DKG *DKGSession) InitializeDKG(round uint32, mrenclave [32]byte, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) (*types.Transaction, error) {
	return _DKG.Contract.InitializeDKG(&_DKG.TransactOpts, round, mrenclave, dkgPubKey, commPubKey, rawQuote)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0xb1888cd3.
//
// Solidity: function initializeDKG(uint32 round, bytes32 mrenclave, bytes dkgPubKey, bytes commPubKey, bytes rawQuote) returns()
func (_DKG *DKGTransactorSession) InitializeDKG(round uint32, mrenclave [32]byte, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) (*types.Transaction, error) {
	return _DKG.Contract.InitializeDKG(&_DKG.TransactOpts, round, mrenclave, dkgPubKey, commPubKey, rawQuote)
}

// RequestRemoteAttestationOnChain is a paid mutator transaction binding the contract method 0xfa4e9f63.
//
// Solidity: function requestRemoteAttestationOnChain(address targetValidatorAddr, uint32 round, bytes32 mrenclave) returns()
func (_DKG *DKGTransactor) RequestRemoteAttestationOnChain(opts *bind.TransactOpts, targetValidatorAddr common.Address, round uint32, mrenclave [32]byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "requestRemoteAttestationOnChain", targetValidatorAddr, round, mrenclave)
}

// RequestRemoteAttestationOnChain is a paid mutator transaction binding the contract method 0xfa4e9f63.
//
// Solidity: function requestRemoteAttestationOnChain(address targetValidatorAddr, uint32 round, bytes32 mrenclave) returns()
func (_DKG *DKGSession) RequestRemoteAttestationOnChain(targetValidatorAddr common.Address, round uint32, mrenclave [32]byte) (*types.Transaction, error) {
	return _DKG.Contract.RequestRemoteAttestationOnChain(&_DKG.TransactOpts, targetValidatorAddr, round, mrenclave)
}

// RequestRemoteAttestationOnChain is a paid mutator transaction binding the contract method 0xfa4e9f63.
//
// Solidity: function requestRemoteAttestationOnChain(address targetValidatorAddr, uint32 round, bytes32 mrenclave) returns()
func (_DKG *DKGTransactorSession) RequestRemoteAttestationOnChain(targetValidatorAddr common.Address, round uint32, mrenclave [32]byte) (*types.Transaction, error) {
	return _DKG.Contract.RequestRemoteAttestationOnChain(&_DKG.TransactOpts, targetValidatorAddr, round, mrenclave)
}

// RequestThresholdDecryption is a paid mutator transaction binding the contract method 0xb1133cac.
//
// Solidity: function requestThresholdDecryption(uint32 round, bytes32 mrenclave, bytes requesterPubKey, bytes ciphertext, bytes label) returns()
func (_DKG *DKGTransactor) RequestThresholdDecryption(opts *bind.TransactOpts, round uint32, mrenclave [32]byte, requesterPubKey []byte, ciphertext []byte, label []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "requestThresholdDecryption", round, mrenclave, requesterPubKey, ciphertext, label)
}

// RequestThresholdDecryption is a paid mutator transaction binding the contract method 0xb1133cac.
//
// Solidity: function requestThresholdDecryption(uint32 round, bytes32 mrenclave, bytes requesterPubKey, bytes ciphertext, bytes label) returns()
func (_DKG *DKGSession) RequestThresholdDecryption(round uint32, mrenclave [32]byte, requesterPubKey []byte, ciphertext []byte, label []byte) (*types.Transaction, error) {
	return _DKG.Contract.RequestThresholdDecryption(&_DKG.TransactOpts, round, mrenclave, requesterPubKey, ciphertext, label)
}

// RequestThresholdDecryption is a paid mutator transaction binding the contract method 0xb1133cac.
//
// Solidity: function requestThresholdDecryption(uint32 round, bytes32 mrenclave, bytes requesterPubKey, bytes ciphertext, bytes label) returns()
func (_DKG *DKGTransactorSession) RequestThresholdDecryption(round uint32, mrenclave [32]byte, requesterPubKey []byte, ciphertext []byte, label []byte) (*types.Transaction, error) {
	return _DKG.Contract.RequestThresholdDecryption(&_DKG.TransactOpts, round, mrenclave, requesterPubKey, ciphertext, label)
}

// SetNetwork is a paid mutator transaction binding the contract method 0x227cd922.
//
// Solidity: function setNetwork(uint32 round, uint32 total, uint32 threshold, bytes32 mrenclave, bytes signature) returns()
func (_DKG *DKGTransactor) SetNetwork(opts *bind.TransactOpts, round uint32, total uint32, threshold uint32, mrenclave [32]byte, signature []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "setNetwork", round, total, threshold, mrenclave, signature)
}

// SetNetwork is a paid mutator transaction binding the contract method 0x227cd922.
//
// Solidity: function setNetwork(uint32 round, uint32 total, uint32 threshold, bytes32 mrenclave, bytes signature) returns()
func (_DKG *DKGSession) SetNetwork(round uint32, total uint32, threshold uint32, mrenclave [32]byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.SetNetwork(&_DKG.TransactOpts, round, total, threshold, mrenclave, signature)
}

// SetNetwork is a paid mutator transaction binding the contract method 0x227cd922.
//
// Solidity: function setNetwork(uint32 round, uint32 total, uint32 threshold, bytes32 mrenclave, bytes signature) returns()
func (_DKG *DKGTransactorSession) SetNetwork(round uint32, total uint32, threshold uint32, mrenclave [32]byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.SetNetwork(&_DKG.TransactOpts, round, total, threshold, mrenclave, signature)
}

// SubmitPartialDecryption is a paid mutator transaction binding the contract method 0x7fb93a3a.
//
// Solidity: function submitPartialDecryption(uint32 round, bytes32 mrenclave, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label) returns()
func (_DKG *DKGTransactor) SubmitPartialDecryption(opts *bind.TransactOpts, round uint32, mrenclave [32]byte, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, label []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "submitPartialDecryption", round, mrenclave, encryptedPartial, ephemeralPubKey, pubShare, label)
}

// SubmitPartialDecryption is a paid mutator transaction binding the contract method 0x7fb93a3a.
//
// Solidity: function submitPartialDecryption(uint32 round, bytes32 mrenclave, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label) returns()
func (_DKG *DKGSession) SubmitPartialDecryption(round uint32, mrenclave [32]byte, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, label []byte) (*types.Transaction, error) {
	return _DKG.Contract.SubmitPartialDecryption(&_DKG.TransactOpts, round, mrenclave, encryptedPartial, ephemeralPubKey, pubShare, label)
}

// SubmitPartialDecryption is a paid mutator transaction binding the contract method 0x7fb93a3a.
//
// Solidity: function submitPartialDecryption(uint32 round, bytes32 mrenclave, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label) returns()
func (_DKG *DKGTransactorSession) SubmitPartialDecryption(round uint32, mrenclave [32]byte, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, label []byte) (*types.Transaction, error) {
	return _DKG.Contract.SubmitPartialDecryption(&_DKG.TransactOpts, round, mrenclave, encryptedPartial, ephemeralPubKey, pubShare, label)
}

// DKGDKGFinalizedIterator is returned from FilterDKGFinalized and is used to iterate over the raw logs and unpacked data for DKGFinalized events raised by the DKG contract.
type DKGDKGFinalizedIterator struct {
	Event *DKGDKGFinalized // Event containing the contract specifics and raw log

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
func (it *DKGDKGFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGDKGFinalized)
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
		it.Event = new(DKGDKGFinalized)
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
func (it *DKGDKGFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGDKGFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGDKGFinalized represents a DKGFinalized event raised by the DKG contract.
type DKGDKGFinalized struct {
	MsgSender    common.Address
	Round        uint32
	Mrenclave    [32]byte
	GlobalPubKey []byte
	Signature    []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDKGFinalized is a free log retrieval operation binding the contract event 0xe7419c96e4837a0c8c3c13342ccea4095f978269a67a3ce5dcfeac5664240205.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes signature)
func (_DKG *DKGFilterer) FilterDKGFinalized(opts *bind.FilterOpts, msgSender []common.Address) (*DKGDKGFinalizedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _DKG.contract.FilterLogs(opts, "DKGFinalized", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &DKGDKGFinalizedIterator{contract: _DKG.contract, event: "DKGFinalized", logs: logs, sub: sub}, nil
}

// WatchDKGFinalized is a free log subscription operation binding the contract event 0xe7419c96e4837a0c8c3c13342ccea4095f978269a67a3ce5dcfeac5664240205.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes signature)
func (_DKG *DKGFilterer) WatchDKGFinalized(opts *bind.WatchOpts, sink chan<- *DKGDKGFinalized, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _DKG.contract.WatchLogs(opts, "DKGFinalized", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGDKGFinalized)
				if err := _DKG.contract.UnpackLog(event, "DKGFinalized", log); err != nil {
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

// ParseDKGFinalized is a log parse operation binding the contract event 0xe7419c96e4837a0c8c3c13342ccea4095f978269a67a3ce5dcfeac5664240205.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes signature)
func (_DKG *DKGFilterer) ParseDKGFinalized(log types.Log) (*DKGDKGFinalized, error) {
	event := new(DKGDKGFinalized)
	if err := _DKG.contract.UnpackLog(event, "DKGFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGDKGInitializedIterator is returned from FilterDKGInitialized and is used to iterate over the raw logs and unpacked data for DKGInitialized events raised by the DKG contract.
type DKGDKGInitializedIterator struct {
	Event *DKGDKGInitialized // Event containing the contract specifics and raw log

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
func (it *DKGDKGInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGDKGInitialized)
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
		it.Event = new(DKGDKGInitialized)
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
func (it *DKGDKGInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGDKGInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGDKGInitialized represents a DKGInitialized event raised by the DKG contract.
type DKGDKGInitialized struct {
	MsgSender  common.Address
	Mrenclave  [32]byte
	Round      uint32
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDKGInitialized is a free log retrieval operation binding the contract event 0x1bd0faa06edbfccdd0f51f46517f5bae23b4abca2dad81e938e89f3ddf7cab1d.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes32 mrenclave, uint32 round, bytes dkgPubKey, bytes commPubKey, bytes rawQuote)
func (_DKG *DKGFilterer) FilterDKGInitialized(opts *bind.FilterOpts, msgSender []common.Address) (*DKGDKGInitializedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _DKG.contract.FilterLogs(opts, "DKGInitialized", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &DKGDKGInitializedIterator{contract: _DKG.contract, event: "DKGInitialized", logs: logs, sub: sub}, nil
}

// WatchDKGInitialized is a free log subscription operation binding the contract event 0x1bd0faa06edbfccdd0f51f46517f5bae23b4abca2dad81e938e89f3ddf7cab1d.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes32 mrenclave, uint32 round, bytes dkgPubKey, bytes commPubKey, bytes rawQuote)
func (_DKG *DKGFilterer) WatchDKGInitialized(opts *bind.WatchOpts, sink chan<- *DKGDKGInitialized, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _DKG.contract.WatchLogs(opts, "DKGInitialized", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGDKGInitialized)
				if err := _DKG.contract.UnpackLog(event, "DKGInitialized", log); err != nil {
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

// ParseDKGInitialized is a log parse operation binding the contract event 0x1bd0faa06edbfccdd0f51f46517f5bae23b4abca2dad81e938e89f3ddf7cab1d.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes32 mrenclave, uint32 round, bytes dkgPubKey, bytes commPubKey, bytes rawQuote)
func (_DKG *DKGFilterer) ParseDKGInitialized(log types.Log) (*DKGDKGInitialized, error) {
	event := new(DKGDKGInitialized)
	if err := _DKG.contract.UnpackLog(event, "DKGInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGDKGNetworkSetIterator is returned from FilterDKGNetworkSet and is used to iterate over the raw logs and unpacked data for DKGNetworkSet events raised by the DKG contract.
type DKGDKGNetworkSetIterator struct {
	Event *DKGDKGNetworkSet // Event containing the contract specifics and raw log

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
func (it *DKGDKGNetworkSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGDKGNetworkSet)
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
		it.Event = new(DKGDKGNetworkSet)
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
func (it *DKGDKGNetworkSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGDKGNetworkSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGDKGNetworkSet represents a DKGNetworkSet event raised by the DKG contract.
type DKGDKGNetworkSet struct {
	MsgSender common.Address
	Round     uint32
	Total     uint32
	Threshold uint32
	Mrenclave [32]byte
	Signature []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDKGNetworkSet is a free log retrieval operation binding the contract event 0xc7a37268197965e156b6d53085e9e20ba69f731868b09d00c2b2c3925f25f4f8.
//
// Solidity: event DKGNetworkSet(address indexed msgSender, uint32 round, uint32 total, uint32 threshold, bytes32 mrenclave, bytes signature)
func (_DKG *DKGFilterer) FilterDKGNetworkSet(opts *bind.FilterOpts, msgSender []common.Address) (*DKGDKGNetworkSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _DKG.contract.FilterLogs(opts, "DKGNetworkSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &DKGDKGNetworkSetIterator{contract: _DKG.contract, event: "DKGNetworkSet", logs: logs, sub: sub}, nil
}

// WatchDKGNetworkSet is a free log subscription operation binding the contract event 0xc7a37268197965e156b6d53085e9e20ba69f731868b09d00c2b2c3925f25f4f8.
//
// Solidity: event DKGNetworkSet(address indexed msgSender, uint32 round, uint32 total, uint32 threshold, bytes32 mrenclave, bytes signature)
func (_DKG *DKGFilterer) WatchDKGNetworkSet(opts *bind.WatchOpts, sink chan<- *DKGDKGNetworkSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _DKG.contract.WatchLogs(opts, "DKGNetworkSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGDKGNetworkSet)
				if err := _DKG.contract.UnpackLog(event, "DKGNetworkSet", log); err != nil {
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

// ParseDKGNetworkSet is a log parse operation binding the contract event 0xc7a37268197965e156b6d53085e9e20ba69f731868b09d00c2b2c3925f25f4f8.
//
// Solidity: event DKGNetworkSet(address indexed msgSender, uint32 round, uint32 total, uint32 threshold, bytes32 mrenclave, bytes signature)
func (_DKG *DKGFilterer) ParseDKGNetworkSet(log types.Log) (*DKGDKGNetworkSet, error) {
	event := new(DKGDKGNetworkSet)
	if err := _DKG.contract.UnpackLog(event, "DKGNetworkSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGDealComplaintsSubmittedIterator is returned from FilterDealComplaintsSubmitted and is used to iterate over the raw logs and unpacked data for DealComplaintsSubmitted events raised by the DKG contract.
type DKGDealComplaintsSubmittedIterator struct {
	Event *DKGDealComplaintsSubmitted // Event containing the contract specifics and raw log

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
func (it *DKGDealComplaintsSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGDealComplaintsSubmitted)
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
		it.Event = new(DKGDealComplaintsSubmitted)
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
func (it *DKGDealComplaintsSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGDealComplaintsSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGDealComplaintsSubmitted represents a DealComplaintsSubmitted event raised by the DKG contract.
type DKGDealComplaintsSubmitted struct {
	Index           uint32
	ComplainIndexes []uint32
	Round           uint32
	Mrenclave       [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterDealComplaintsSubmitted is a free log retrieval operation binding the contract event 0xa89000d88bdc9c3e92c10abb67235241f8c6803723e88e1e2420533e8fe2b8d8.
//
// Solidity: event DealComplaintsSubmitted(uint32 index, uint32[] complainIndexes, uint32 round, bytes32 mrenclave)
func (_DKG *DKGFilterer) FilterDealComplaintsSubmitted(opts *bind.FilterOpts) (*DKGDealComplaintsSubmittedIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "DealComplaintsSubmitted")
	if err != nil {
		return nil, err
	}
	return &DKGDealComplaintsSubmittedIterator{contract: _DKG.contract, event: "DealComplaintsSubmitted", logs: logs, sub: sub}, nil
}

// WatchDealComplaintsSubmitted is a free log subscription operation binding the contract event 0xa89000d88bdc9c3e92c10abb67235241f8c6803723e88e1e2420533e8fe2b8d8.
//
// Solidity: event DealComplaintsSubmitted(uint32 index, uint32[] complainIndexes, uint32 round, bytes32 mrenclave)
func (_DKG *DKGFilterer) WatchDealComplaintsSubmitted(opts *bind.WatchOpts, sink chan<- *DKGDealComplaintsSubmitted) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "DealComplaintsSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGDealComplaintsSubmitted)
				if err := _DKG.contract.UnpackLog(event, "DealComplaintsSubmitted", log); err != nil {
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

// ParseDealComplaintsSubmitted is a log parse operation binding the contract event 0xa89000d88bdc9c3e92c10abb67235241f8c6803723e88e1e2420533e8fe2b8d8.
//
// Solidity: event DealComplaintsSubmitted(uint32 index, uint32[] complainIndexes, uint32 round, bytes32 mrenclave)
func (_DKG *DKGFilterer) ParseDealComplaintsSubmitted(log types.Log) (*DKGDealComplaintsSubmitted, error) {
	event := new(DKGDealComplaintsSubmitted)
	if err := _DKG.contract.UnpackLog(event, "DealComplaintsSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGDealVerifiedIterator is returned from FilterDealVerified and is used to iterate over the raw logs and unpacked data for DealVerified events raised by the DKG contract.
type DKGDealVerifiedIterator struct {
	Event *DKGDealVerified // Event containing the contract specifics and raw log

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
func (it *DKGDealVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGDealVerified)
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
		it.Event = new(DKGDealVerified)
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
func (it *DKGDealVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGDealVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGDealVerified represents a DealVerified event raised by the DKG contract.
type DKGDealVerified struct {
	Index          uint32
	RecipientIndex uint32
	Round          uint32
	Mrenclave      [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDealVerified is a free log retrieval operation binding the contract event 0x1a8f868a6f5b289bec8c24a7c28727ddc869adf4a3c3c0ae8a2d41b9afc345bb.
//
// Solidity: event DealVerified(uint32 index, uint32 recipientIndex, uint32 round, bytes32 mrenclave)
func (_DKG *DKGFilterer) FilterDealVerified(opts *bind.FilterOpts) (*DKGDealVerifiedIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "DealVerified")
	if err != nil {
		return nil, err
	}
	return &DKGDealVerifiedIterator{contract: _DKG.contract, event: "DealVerified", logs: logs, sub: sub}, nil
}

// WatchDealVerified is a free log subscription operation binding the contract event 0x1a8f868a6f5b289bec8c24a7c28727ddc869adf4a3c3c0ae8a2d41b9afc345bb.
//
// Solidity: event DealVerified(uint32 index, uint32 recipientIndex, uint32 round, bytes32 mrenclave)
func (_DKG *DKGFilterer) WatchDealVerified(opts *bind.WatchOpts, sink chan<- *DKGDealVerified) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "DealVerified")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGDealVerified)
				if err := _DKG.contract.UnpackLog(event, "DealVerified", log); err != nil {
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

// ParseDealVerified is a log parse operation binding the contract event 0x1a8f868a6f5b289bec8c24a7c28727ddc869adf4a3c3c0ae8a2d41b9afc345bb.
//
// Solidity: event DealVerified(uint32 index, uint32 recipientIndex, uint32 round, bytes32 mrenclave)
func (_DKG *DKGFilterer) ParseDealVerified(log types.Log) (*DKGDealVerified, error) {
	event := new(DKGDealVerified)
	if err := _DKG.contract.UnpackLog(event, "DealVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGInvalidDealIterator is returned from FilterInvalidDeal and is used to iterate over the raw logs and unpacked data for InvalidDeal events raised by the DKG contract.
type DKGInvalidDealIterator struct {
	Event *DKGInvalidDeal // Event containing the contract specifics and raw log

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
func (it *DKGInvalidDealIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGInvalidDeal)
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
		it.Event = new(DKGInvalidDeal)
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
func (it *DKGInvalidDealIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGInvalidDealIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGInvalidDeal represents a InvalidDeal event raised by the DKG contract.
type DKGInvalidDeal struct {
	Index     uint32
	Round     uint32
	Mrenclave [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterInvalidDeal is a free log retrieval operation binding the contract event 0x90395d01853e3de18e643761e8429ec973c5b4843dbf47451c4e90f37c3447ca.
//
// Solidity: event InvalidDeal(uint32 index, uint32 round, bytes32 mrenclave)
func (_DKG *DKGFilterer) FilterInvalidDeal(opts *bind.FilterOpts) (*DKGInvalidDealIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "InvalidDeal")
	if err != nil {
		return nil, err
	}
	return &DKGInvalidDealIterator{contract: _DKG.contract, event: "InvalidDeal", logs: logs, sub: sub}, nil
}

// WatchInvalidDeal is a free log subscription operation binding the contract event 0x90395d01853e3de18e643761e8429ec973c5b4843dbf47451c4e90f37c3447ca.
//
// Solidity: event InvalidDeal(uint32 index, uint32 round, bytes32 mrenclave)
func (_DKG *DKGFilterer) WatchInvalidDeal(opts *bind.WatchOpts, sink chan<- *DKGInvalidDeal) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "InvalidDeal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGInvalidDeal)
				if err := _DKG.contract.UnpackLog(event, "InvalidDeal", log); err != nil {
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

// ParseInvalidDeal is a log parse operation binding the contract event 0x90395d01853e3de18e643761e8429ec973c5b4843dbf47451c4e90f37c3447ca.
//
// Solidity: event InvalidDeal(uint32 index, uint32 round, bytes32 mrenclave)
func (_DKG *DKGFilterer) ParseInvalidDeal(log types.Log) (*DKGInvalidDeal, error) {
	event := new(DKGInvalidDeal)
	if err := _DKG.contract.UnpackLog(event, "InvalidDeal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGPartialDecryptionSubmittedIterator is returned from FilterPartialDecryptionSubmitted and is used to iterate over the raw logs and unpacked data for PartialDecryptionSubmitted events raised by the DKG contract.
type DKGPartialDecryptionSubmittedIterator struct {
	Event *DKGPartialDecryptionSubmitted // Event containing the contract specifics and raw log

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
func (it *DKGPartialDecryptionSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGPartialDecryptionSubmitted)
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
		it.Event = new(DKGPartialDecryptionSubmitted)
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
func (it *DKGPartialDecryptionSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGPartialDecryptionSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGPartialDecryptionSubmitted represents a PartialDecryptionSubmitted event raised by the DKG contract.
type DKGPartialDecryptionSubmitted struct {
	Validator        common.Address
	Round            uint32
	Mrenclave        [32]byte
	EncryptedPartial []byte
	EphemeralPubKey  []byte
	PubShare         []byte
	Label            []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterPartialDecryptionSubmitted is a free log retrieval operation binding the contract event 0x2d05ab82940dc8b5196f5969dc54167bb42e65a5018b12a99c8a7f20e1dc2464.
//
// Solidity: event PartialDecryptionSubmitted(address indexed validator, uint32 round, bytes32 mrenclave, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label)
func (_DKG *DKGFilterer) FilterPartialDecryptionSubmitted(opts *bind.FilterOpts, validator []common.Address) (*DKGPartialDecryptionSubmittedIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _DKG.contract.FilterLogs(opts, "PartialDecryptionSubmitted", validatorRule)
	if err != nil {
		return nil, err
	}
	return &DKGPartialDecryptionSubmittedIterator{contract: _DKG.contract, event: "PartialDecryptionSubmitted", logs: logs, sub: sub}, nil
}

// WatchPartialDecryptionSubmitted is a free log subscription operation binding the contract event 0x2d05ab82940dc8b5196f5969dc54167bb42e65a5018b12a99c8a7f20e1dc2464.
//
// Solidity: event PartialDecryptionSubmitted(address indexed validator, uint32 round, bytes32 mrenclave, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label)
func (_DKG *DKGFilterer) WatchPartialDecryptionSubmitted(opts *bind.WatchOpts, sink chan<- *DKGPartialDecryptionSubmitted, validator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _DKG.contract.WatchLogs(opts, "PartialDecryptionSubmitted", validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGPartialDecryptionSubmitted)
				if err := _DKG.contract.UnpackLog(event, "PartialDecryptionSubmitted", log); err != nil {
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

// ParsePartialDecryptionSubmitted is a log parse operation binding the contract event 0x2d05ab82940dc8b5196f5969dc54167bb42e65a5018b12a99c8a7f20e1dc2464.
//
// Solidity: event PartialDecryptionSubmitted(address indexed validator, uint32 round, bytes32 mrenclave, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label)
func (_DKG *DKGFilterer) ParsePartialDecryptionSubmitted(log types.Log) (*DKGPartialDecryptionSubmitted, error) {
	event := new(DKGPartialDecryptionSubmitted)
	if err := _DKG.contract.UnpackLog(event, "PartialDecryptionSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGRemoteAttestationProcessedOnChainIterator is returned from FilterRemoteAttestationProcessedOnChain and is used to iterate over the raw logs and unpacked data for RemoteAttestationProcessedOnChain events raised by the DKG contract.
type DKGRemoteAttestationProcessedOnChainIterator struct {
	Event *DKGRemoteAttestationProcessedOnChain // Event containing the contract specifics and raw log

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
func (it *DKGRemoteAttestationProcessedOnChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGRemoteAttestationProcessedOnChain)
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
		it.Event = new(DKGRemoteAttestationProcessedOnChain)
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
func (it *DKGRemoteAttestationProcessedOnChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGRemoteAttestationProcessedOnChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGRemoteAttestationProcessedOnChain represents a RemoteAttestationProcessedOnChain event raised by the DKG contract.
type DKGRemoteAttestationProcessedOnChain struct {
	Validator  common.Address
	ChalStatus uint8
	Round      uint32
	Mrenclave  [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRemoteAttestationProcessedOnChain is a free log retrieval operation binding the contract event 0x54690f98c0ec0056e0e487f4fe5e8eea7bee88d2dbb7cc9ddca22981f06d9dbb.
//
// Solidity: event RemoteAttestationProcessedOnChain(address validator, uint8 chalStatus, uint32 round, bytes32 mrenclave)
func (_DKG *DKGFilterer) FilterRemoteAttestationProcessedOnChain(opts *bind.FilterOpts) (*DKGRemoteAttestationProcessedOnChainIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "RemoteAttestationProcessedOnChain")
	if err != nil {
		return nil, err
	}
	return &DKGRemoteAttestationProcessedOnChainIterator{contract: _DKG.contract, event: "RemoteAttestationProcessedOnChain", logs: logs, sub: sub}, nil
}

// WatchRemoteAttestationProcessedOnChain is a free log subscription operation binding the contract event 0x54690f98c0ec0056e0e487f4fe5e8eea7bee88d2dbb7cc9ddca22981f06d9dbb.
//
// Solidity: event RemoteAttestationProcessedOnChain(address validator, uint8 chalStatus, uint32 round, bytes32 mrenclave)
func (_DKG *DKGFilterer) WatchRemoteAttestationProcessedOnChain(opts *bind.WatchOpts, sink chan<- *DKGRemoteAttestationProcessedOnChain) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "RemoteAttestationProcessedOnChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGRemoteAttestationProcessedOnChain)
				if err := _DKG.contract.UnpackLog(event, "RemoteAttestationProcessedOnChain", log); err != nil {
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

// ParseRemoteAttestationProcessedOnChain is a log parse operation binding the contract event 0x54690f98c0ec0056e0e487f4fe5e8eea7bee88d2dbb7cc9ddca22981f06d9dbb.
//
// Solidity: event RemoteAttestationProcessedOnChain(address validator, uint8 chalStatus, uint32 round, bytes32 mrenclave)
func (_DKG *DKGFilterer) ParseRemoteAttestationProcessedOnChain(log types.Log) (*DKGRemoteAttestationProcessedOnChain, error) {
	event := new(DKGRemoteAttestationProcessedOnChain)
	if err := _DKG.contract.UnpackLog(event, "RemoteAttestationProcessedOnChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGThresholdDecryptRequestedIterator is returned from FilterThresholdDecryptRequested and is used to iterate over the raw logs and unpacked data for ThresholdDecryptRequested events raised by the DKG contract.
type DKGThresholdDecryptRequestedIterator struct {
	Event *DKGThresholdDecryptRequested // Event containing the contract specifics and raw log

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
func (it *DKGThresholdDecryptRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGThresholdDecryptRequested)
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
		it.Event = new(DKGThresholdDecryptRequested)
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
func (it *DKGThresholdDecryptRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGThresholdDecryptRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGThresholdDecryptRequested represents a ThresholdDecryptRequested event raised by the DKG contract.
type DKGThresholdDecryptRequested struct {
	Requester       common.Address
	Round           uint32
	Mrenclave       [32]byte
	RequesterPubKey []byte
	Ciphertext      []byte
	Label           []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterThresholdDecryptRequested is a free log retrieval operation binding the contract event 0x5f9d4b68667a3b91f2fe8369c2ac9040ce4a68400aadbe076f9d109b13e09b61.
//
// Solidity: event ThresholdDecryptRequested(address indexed requester, uint32 round, bytes32 mrenclave, bytes requesterPubKey, bytes ciphertext, bytes label)
func (_DKG *DKGFilterer) FilterThresholdDecryptRequested(opts *bind.FilterOpts, requester []common.Address) (*DKGThresholdDecryptRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _DKG.contract.FilterLogs(opts, "ThresholdDecryptRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &DKGThresholdDecryptRequestedIterator{contract: _DKG.contract, event: "ThresholdDecryptRequested", logs: logs, sub: sub}, nil
}

// WatchThresholdDecryptRequested is a free log subscription operation binding the contract event 0x5f9d4b68667a3b91f2fe8369c2ac9040ce4a68400aadbe076f9d109b13e09b61.
//
// Solidity: event ThresholdDecryptRequested(address indexed requester, uint32 round, bytes32 mrenclave, bytes requesterPubKey, bytes ciphertext, bytes label)
func (_DKG *DKGFilterer) WatchThresholdDecryptRequested(opts *bind.WatchOpts, sink chan<- *DKGThresholdDecryptRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _DKG.contract.WatchLogs(opts, "ThresholdDecryptRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGThresholdDecryptRequested)
				if err := _DKG.contract.UnpackLog(event, "ThresholdDecryptRequested", log); err != nil {
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

// ParseThresholdDecryptRequested is a log parse operation binding the contract event 0x5f9d4b68667a3b91f2fe8369c2ac9040ce4a68400aadbe076f9d109b13e09b61.
//
// Solidity: event ThresholdDecryptRequested(address indexed requester, uint32 round, bytes32 mrenclave, bytes requesterPubKey, bytes ciphertext, bytes label)
func (_DKG *DKGFilterer) ParseThresholdDecryptRequested(log types.Log) (*DKGThresholdDecryptRequested, error) {
	event := new(DKGThresholdDecryptRequested)
	if err := _DKG.contract.UnpackLog(event, "ThresholdDecryptRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DKGUpgradeScheduledIterator is returned from FilterUpgradeScheduled and is used to iterate over the raw logs and unpacked data for UpgradeScheduled events raised by the DKG contract.
type DKGUpgradeScheduledIterator struct {
	Event *DKGUpgradeScheduled // Event containing the contract specifics and raw log

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
func (it *DKGUpgradeScheduledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGUpgradeScheduled)
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
		it.Event = new(DKGUpgradeScheduled)
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
func (it *DKGUpgradeScheduledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGUpgradeScheduledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGUpgradeScheduled represents a UpgradeScheduled event raised by the DKG contract.
type DKGUpgradeScheduled struct {
	ActivationHeight uint32
	Mrenclave        [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpgradeScheduled is a free log retrieval operation binding the contract event 0xba889db5cdb62d54ea6ab3c85ea27c52b2710b39cb8ecddf1a360a51cbb40110.
//
// Solidity: event UpgradeScheduled(uint32 activationHeight, bytes32 mrenclave)
func (_DKG *DKGFilterer) FilterUpgradeScheduled(opts *bind.FilterOpts) (*DKGUpgradeScheduledIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "UpgradeScheduled")
	if err != nil {
		return nil, err
	}
	return &DKGUpgradeScheduledIterator{contract: _DKG.contract, event: "UpgradeScheduled", logs: logs, sub: sub}, nil
}

// WatchUpgradeScheduled is a free log subscription operation binding the contract event 0xba889db5cdb62d54ea6ab3c85ea27c52b2710b39cb8ecddf1a360a51cbb40110.
//
// Solidity: event UpgradeScheduled(uint32 activationHeight, bytes32 mrenclave)
func (_DKG *DKGFilterer) WatchUpgradeScheduled(opts *bind.WatchOpts, sink chan<- *DKGUpgradeScheduled) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "UpgradeScheduled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGUpgradeScheduled)
				if err := _DKG.contract.UnpackLog(event, "UpgradeScheduled", log); err != nil {
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

// ParseUpgradeScheduled is a log parse operation binding the contract event 0xba889db5cdb62d54ea6ab3c85ea27c52b2710b39cb8ecddf1a360a51cbb40110.
//
// Solidity: event UpgradeScheduled(uint32 activationHeight, bytes32 mrenclave)
func (_DKG *DKGFilterer) ParseUpgradeScheduled(log types.Log) (*DKGUpgradeScheduled, error) {
	event := new(DKGUpgradeScheduled)
	if err := _DKG.contract.UnpackLog(event, "UpgradeScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
