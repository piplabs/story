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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"complainDeals\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"curMrenclave\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dealComplaints\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainant\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dkgNodeInfos\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"nodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.NodeStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"publicCoeffs\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getNodeInfo\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDKG.NodeInfo\",\"components\":[{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"nodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.NodeStatus\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"partialDecrypts\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"labelHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"pid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"partialDecryption\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestRemoteAttestationOnChain\",\"inputs\":[{\"name\":\"targetValidatorAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestThresholdDecryption\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNetwork\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitPartialDecryption\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"pid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ephemeralPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"DKGFinalized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"publicCoeffs\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGInitialized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGNetworkSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"total\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealComplaintsSubmitted\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"indexed\":false,\"internalType\":\"uint32[]\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealVerified\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"recipientIndex\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidDeal\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PartialDecryptionSubmitted\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"pid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ephemeralPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoteAttestationProcessedOnChain\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ThresholdDecryptRequested\",\"inputs\":[{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeScheduled\",\"inputs\":[{\"name\":\"activationHeight\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60803461005057601f612c2b38819003918201601f19168301916001600160401b038311848410176100545780849260209460405283398101031261005057515f55604051612bc290816100698239f35b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c806308ad63ac146100d4578063227cd922146100cf578063681a0fa8146100ca57806379b131a7146100c55780639af5962c146100c0578063a26f51a4146100bb578063aab066c6146100b6578063b1133cac146100b1578063b1888cd3146100ac578063dd7b0d8a146100a7578063dea942d9146100a25763fa4e9f631461009d575f80fd5b6110f1565b611026565b610eb4565b610ce8565b610b7c565b610a23565b610727565b61064b565b61058e565b610554565b61030a565b610210565b6004359063ffffffff821682036100ec57565b5f80fd5b6024359063ffffffff821682036100ec57565b6044359063ffffffff821682036100ec57565b6064359063ffffffff821682036100ec57565b359063ffffffff821682036100ec57565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b60a0810190811067ffffffffffffffff82111761018357604052565b61013a565b6040810190811067ffffffffffffffff82111761018357604052565b6060810190811067ffffffffffffffff82111761018357604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761018357604052565b6040519061020e82610167565b565b346100ec5760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec576102476100d9565b61024f6100f0565b906044359167ffffffffffffffff908184116100ec57366023850112156100ec578360040135918211610183578160051b936020946040519361029560208301866101c0565b845260246020850191830101913683116100ec57602401905b8282106102c5576102c3606435868689611338565b005b8680916102d184610129565b8152019101906102ae565b9181601f840112156100ec5782359167ffffffffffffffff83116100ec57602083818601950101116100ec57565b346100ec5760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec576103416100d9565b6103496100f0565b90610352610103565b6064359160843567ffffffffffffffff81116100ec576103769036906004016102dc565b90604080516020908181019088825282815261039181610188565b5190205f5483518381019182528381526103aa81610188565b519020146103b7906112d3565b5f8781526001825282812063ffffffff87168252602090815260408083203384529091529020906103ea600183016108ad565b8351828101908a82527fffffffff00000000000000000000000000000000000000000000000000000000808a60e01b1687830152808d60e01b1660448301528a60e01b166048820152602c8152610440816101a4565b519020610474907f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f52601c52603c5f2090565b61047f368888611a7a565b6104889161291e565b8151919092012073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166104f19173ffffffffffffffffffffffffffffffffffffffff1614611526565b60030180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff1661020017905551948594339761052e9587611605565b037fc7a37268197965e156b6d53085e9e20ba69f731868b09d00c2b2c3925f25f4f891a2005b346100ec575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec5760205f54604051908152f35b346100ec5760e07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec576105c56100d9565b6105cd610103565b67ffffffffffffffff916064358381116100ec576105ef9036906004016102dc565b906084358581116100ec576106089036906004016102dc565b9060a4358781116100ec576106219036906004016102dc565b94909360c4359889116100ec5761063f6102c39936906004016102dc565b9890976024359061163c565b346100ec5760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec576106826100d9565b67ffffffffffffffff6044358181116100ec576106a39036906004016102dc565b606493919335918383116100ec57366023840112156100ec578260040135918483116100ec573660248460051b860101116100ec576084359485116100ec576102c3956106f660249636906004016102dc565b969095019260243590611f09565b6004359073ffffffffffffffffffffffffffffffffffffffff821682036100ec57565b346100ec5760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec5761075e6100f0565b610766610103565b906064359173ffffffffffffffffffffffffffffffffffffffff831683036100ec576020926107cb6107ee926107b860ff956004355f526002885260405f209063ffffffff165f5260205260405f2090565b9063ffffffff165f5260205260405f2090565b9073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b54166040519015158152f35b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc60609101126100ec576004359060243563ffffffff811681036100ec579060443573ffffffffffffffffffffffffffffffffffffffff811681036100ec5790565b90600182811c921680156108a3575b602083101461087657565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f169161086b565b9060405191825f82546108bf8161085c565b908184526020946001916001811690815f1461092b57506001146108ed575b50505061020e925003836101c0565b5f90815285812095935091905b81831061091357505061020e93508201015f80806108de565b855488840185015294850194879450918301916108fa565b91505061020e9593507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201015f80806108de565b5f5b83811061097d5750505f910152565b818101518382015260200161096e565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f6020936109c98151809281875287808801910161096c565b0116010190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b60031115610a0757565b6109d0565b906003821015610a075752565b60041115610a0757565b346100ec57610aa7610a5a6107cb610a3a366107fa565b92915f52600160205260405f209063ffffffff165f5260205260405f2090565b610a63816108ad565b90610a70600182016108ad565b610ad1610ac36003610a84600286016108ad565b94015493610ab560ff8660081c169460405198899860a08a5260a08a019061098d565b9088820360208a015261098d565b90868203604088015261098d565b9260ff606086019116610a0c565b610ada81610a19565b60808301520390f35b9060a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8301126100ec5760043563ffffffff811681036100ec57916024359167ffffffffffffffff6044358181116100ec5783610b43916004016102dc565b939093926064358381116100ec5782610b5e916004016102dc565b939093926084359182116100ec57610b78916004016102dc565b9091565b346100ec57610b8a36610ae3565b969593604093919351610bcf60209182810190888252838152610bac81610188565b5190205f54604051848101918252848152610bc681610188565b519020146112d3565b610be063ffffffff8916151561181c565b8315610c8a5760418203610c2c575091610c2793917f5f9d4b68667a3b91f2fe8369c2ac9040ce4a68400aadbe076f9d109b13e09b61979893604051978897339b896121ec565b0390a2005b606490604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601860248201527f496e76616c696420726571756573746572207075626b657900000000000000006044820152fd5b606490604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601060248201527f456d7074792063697068657274657874000000000000000000000000000000006044820152fd5b346100ec57610cf636610ae3565b9260409794979291925197610d346020998a8101908a82528b8152610d1a81610188565b5190205f546040518c81019182528c8152610bc681610188565b610d61610d42368787611a7a565b610d4d36848b611a7a565b88610d59368888611a7a565b923390612832565b15610e045790610c2794939291610df67f1bd0faa06edbfccdd0f51f46517f5bae23b4abca2dad81e938e89f3ddf7cab1d999a610d9c610201565b90610da836858d611a7a565b8252610db5368787611a7a565b90820152610dc4368888611a7a565b60408201525f606082015260016080820152610df18c6107cb8b6107b833935f52600160205260405f2090565b61227a565b604051978897339b896123f5565b606489604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601a60248201527f496e76616c69642072656d6f7465206174746573746174696f6e0000000000006044820152fd5b9390608093610ab5610eac9473ffffffffffffffffffffffffffffffffffffffff610e9e949a999a16885260a0602089015260a088019061098d565b90848203606086015261098d565b931515910152565b346100ec5760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec57610f3a610eee6100f0565b610f1b610ef9610116565b916004355f52600360205260405f209063ffffffff165f5260205260405f2090565b6044355f5260205260405f209063ffffffff165f5260205260405f2090565b73ffffffffffffffffffffffffffffffffffffffff815416610f90610f61600184016108ad565b92610f6e600282016108ad565b9060ff6004610f7f600384016108ad565b920154169160405195869586610e62565b0390f35b6020815260a06080611003610fb4855184602087015260c086019061098d565b610fee6020870151917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0928388830301604089015261098d565b9060408701519086830301606087015261098d565b93611015606082015183860190610a0c565b01519161102183610a19565b015290565b346100ec57610f906107cb61108761103d366107fa565b9193906040945f6080875161105181610167565b606081526060602082015260608982015282606082015201525f526001602052845f209063ffffffff165f5260205260405f2090565b9060ff600382519361109885610167565b6110a1816108ad565b85526110af600182016108ad565b60208601526110c0600282016108ad565b8486015201546110d582821660608601612237565b60081c166110e281610a19565b60808301525191829182610f94565b346100ec5760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100ec57611128610704565b6111306100f0565b6044359160405161115060209182810190868252838152610bac81610188565b835f5260018152611176826107cb8560405f209063ffffffff165f5260205260405f2090565b90611181825461085c565b1561127557507f54690f98c0ec0056e0e487f4fe5e8eea7bee88d2dbb7cc9ddca22981f06d9dbb93611235826111ff60036112429501916111d56111c6845460ff1690565b6111cf816109fd565b15612424565b6111e1600282016108ad565b9088886111f960016111f2856108ad565b94016108ad565b93612832565b156112475780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660021781555b5460ff1690565b9360405194859485612489565b0390a1005b80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117815561122e565b606490604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601360248201527f4e6f646520646f6573206e6f74206578697374000000000000000000000000006044820152fd5b156112da57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f496e76616c6964206d72656e636c6176650000000000000000000000000000006044820152fd5b93929360408051906113736020928381019089825284815261135981610188565b5190205f54604051858101918252858152610bc681610188565b865f526001918291600182526113c061139d8660405f209063ffffffff165f5260205260405f2090565b3373ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b505f925b611405575b505050507fa89000d88bdc9c3e92c10abb67235241f8c6803723e88e1e2420533e8fe2b8d8939461140091604051948594856114cb565b0390a1565b8651831015611499575f8981526002835281812063ffffffff8716825260205260409020875184101561149457849361148d6114628a6107cb889563ffffffff8933948860051b0101511663ffffffff165f5260205260405f2090565b60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00825416179055565b01926113c4565b61149e565b6113c9565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b91949392946080830163ffffffff8093168452602090608060208601528251809152602060a086019301915f5b82811061151057505050509416604082015260600152565b83518616855293810193928101926001016114f8565b1561152d57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f496e76616c696420736574206e6574776f726b207369676e61747572650000006044820152fd5b9061159581610a19565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff0083549260081b169116179055565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe093818652868601375f8582860101520116010190565b9160a09361163997959263ffffffff92838092168652166020850152166040830152606082015281608082015201916115c7565b90565b9692989490999591979389898c8a604051916020928381019082825284815261166481610188565b5190205f5460405185810191825285815261167e81610188565b5190201461168b906112d3565b63ffffffff61169d838216151561181c565b841615156116aa90611881565b6116b58615156118e6565b6116c08a151561194b565b6116cb8c15156119b0565b6116d760418914611a15565b6116e2368d8d611a7a565b838151910120928484846116fe855f52600360205260405f2090565b90611716919063ffffffff165f5260205260405f2090565b5f9182526020526040902090611739919063ffffffff165f5260205260405f2090565b6004015460ff161561174a90611ade565b611752610201565b3381529587369061176292611a7a565b90860152611771368b8b611a7a565b6040860152611781368d8d611a7a565b60608601526001608086015261179f905f52600360205260405f2090565b906117b7919063ffffffff165f5260205260405f2090565b5f91825260205260409020906117da919063ffffffff165f5260205260405f2090565b906117e491611cbd565b604051998a99339c6117f69a8c611ea1565b037f835e2245f021610650983a80011abf0755d752f7ce7935d861f90dcfa4ad8db291a2565b1561182357565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f496e76616c696420726f756e64000000000000000000000000000000000000006044820152fd5b1561188857565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f496e76616c6964207069640000000000000000000000000000000000000000006044820152fd5b156118ed57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f456d707479207061727469616c000000000000000000000000000000000000006044820152fd5b1561195257565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f456d7074792070756253686172650000000000000000000000000000000000006044820152fd5b156119b757565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f456d707479206c6162656c0000000000000000000000000000000000000000006044820152fd5b15611a1c57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f496e76616c696420657068656d6572616c207075626b657900000000000000006044820152fd5b92919267ffffffffffffffff82116101835760405191611ac260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601846101c0565b8294818452818301116100ec578281602093845f960137010152565b15611ae557565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f5061727469616c20616c7265616479207375626d6974746564000000000000006044820152fd5b601f8211611b5057505050565b5f5260205f20906020601f840160051c83019310611b88575b601f0160051c01905b818110611b7d575050565b5f8155600101611b72565b9091508190611b69565b919091825167ffffffffffffffff811161018357611bba81611bb4845461085c565b84611b43565b602080601f8311600114611c1957508190611c0a9394955f92611c0e575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9055565b015190505f80611bd8565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0831695611c4b855f5260205f2090565b925f905b888210611ca557505083600195969710611c6e575b505050811b019055565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690555f8080611c64565b80600185968294968601518155019501930190611c4f565b9073ffffffffffffffffffffffffffffffffffffffff8151167fffffffffffffffffffffffff000000000000000000000000000000000000000083541617825560018083019060208084015180519267ffffffffffffffff841161018357611d2f84611d29875461085c565b87611b43565b602092601f8511600114611de957505093600493611d8b84611db79560809561020e9a995f92611c0e5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b611d9f604082015160028701611b92565b611db0606082015160038701611b92565b0151151590565b91019060ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0083541691151516179055565b9291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0851690611e1d875f5260205f2090565b945f915b838310611e8a575050508460809461020e99989460049894611db79860019510611e53575b505050811b019055611d8e565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690555f8080611e46565b848601518755958601959481019491810191611e21565b99979592936116399b9995611edf92611eed96611efb9a9560208f63ffffffff809516815201521660408d015260e060608d015260e08c01916115c7565b9189830360808b01526115c7565b9186830360a08801526115c7565b9260c08185039101526115c7565b9691949590939295611f48604051602081019087825260208152611f2c81610188565b5190205f546040516020810191825260208152610bc681610188565b845f526001602052611f6e61139d8960405f209063ffffffff165f5260205260405f2090565b93600385019760ff8954166003811015610a07577f5d25bc3c675e166c30fb9ab70ab3f79501dd65672130ffdaf3bb412531d9fd8299611fde611fd988888f888f8f906120189f918d94611fd460018e95611fce61200a9f83141561201d565b016108ad565b612565565b612082565b6103007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff825416179055565b604051978897339b896120e7565b0390a2565b1561202457565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4e6f64652077617320696e76616c6964617465640000000000000000000000006044820152fd5b1561208957565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f496e76616c69642066696e616c697a6174696f6e207369676e617475726500006044820152fd5b9363ffffffff6121159394929a99979a9896981685526020938486015260a0604086015260a08501916115c7565b828103606084015287815296600581901b8801820195915f818a01845b8483106121525750505050505061163994955060808185039101526115c7565b9091929394987fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08c820301835289357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1833603018112156100ec57820185810191903567ffffffffffffffff81116100ec5780360383136100ec576121dc879283926001956115c7565b9b01930193019194939290612132565b9694926122299461163999979363ffffffff61221b94168a5260208a015260a060408a015260a08901916115c7565b9186830360608801526115c7565b9260808185039101526115c7565b6003821015610a075752565b906003811015610a075760ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008354169116179055565b9190805192835167ffffffffffffffff81116101835761229e81611bb4845461085c565b602080601f83116001146123445750916122f58260039360809561020e98995f92611c0e5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b81555b612309602085015160018301611b92565b61231a604085015160028301611b92565b0191612333606082015161232d816109fd565b84612243565b01519061233f82610a19565b61158b565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0831696612376855f5260205f2090565b925f905b8982106123dd57505092608094926001926003958361020e9a9b106123a7575b505050811b0181556122f8565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f884881b161c191690555f808061239a565b8060018596829496860151815501950193019061237a565b9694926122299463ffffffff6116399a989461221b948b521660208a015260a060408a015260a08901916115c7565b1561242b57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4e6f646520616c7265616479206368616c6c656e6765640000000000000000006044820152fd5b9094939263ffffffff906124c060609473ffffffffffffffffffffffffffffffffffffffff60808601991685526020850190610a0c565b1660408201520152565b91908110156114945760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156100ec57019081359167ffffffffffffffff83116100ec5760200182360381136100ec579190565b6020908361020e93959495604051968361254d899551809288808901910161096c565b84019185830137015f838201520380855201836101c0565b9796959391826125c193604493957fffffffff00000000000000000000000000000000000000000000000000000000604051978895602087015260e01b1660408501528484013781015f838201520360248101845201826101c0565b915f915b80831061267b5750505061265c61263b73ffffffffffffffffffffffffffffffffffffffff946126356126759561262d86602061265c98519101207f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f52601c52603c5f2090565b923691611a7a565b9061291e565b946020815191012073ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1690565b91161490565b90919261269560019161268f8685876124ca565b9161252a565b930191906125c5565b156126a557565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f496e76616c69642076616c696461746f722061646472657373000000000000006044820152fd5b1561270a57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f496e76616c696420726f756e64206f6620444b470000000000000000000000006044820152fd5b1561276f57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f496e76616c696420444b47207075626c6963206b6579000000000000000000006044820152fd5b156127d457565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602060248201527f496e76616c696420636f6d6d756e69636174696f6e207075626c6963206b65796044820152fd5b939192909260408551111561289a576101906116399561286973ffffffffffffffffffffffffffffffffffffffff8716151561269e565b61287a63ffffffff84161515612703565b61288684511515612768565b612892855115156127cd565b015193612934565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f496e76616c6964207261772071756f74652c2071756f746520746f6f2073686f60448201527f72740000000000000000000000000000000000000000000000000000000000006064820152fd5b6116399161292b916129e1565b90929192612a1b565b926129d0916038917fffffffff00000000000000000000000000000000000000000000000000000000946040519586937fffffffffffffffffffffffffffffffffffffffff000000000000000000000000602086019960601b16895260e01b1660348401526129ac815180926020878701910161096c565b82016129c1825180936020878501910161096c565b010360188101845201826101c0565b519020036129dd57600190565b5f90565b8151919060418303612a1157612a0a9250602082015190606060408401519301515f1a90612af2565b9192909190565b50505f9160029190565b612a2481610a19565b80612a2d575050565b612a3681610a19565b60018103612a685760046040517ff645eedf000000000000000000000000000000000000000000000000000000008152fd5b612a7181610a19565b60028103612aab576040517ffce698f700000000000000000000000000000000000000000000000000000000815260048101839052602490fd5b80612ab7600392610a19565b14612abf5750565b6040517fd78bce0c0000000000000000000000000000000000000000000000000000000081526004810191909152602490fd5b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411612b81579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa15612b76575f5173ffffffffffffffffffffffffffffffffffffffff811615612b6c57905f905f90565b505f906001905f90565b6040513d5f823e3d90fd5b5050505f916003919056fea264697066735822122017586c51906413d1b625fbf614bf1129a0995e7aa04b6b7b051bda05894b8ee064736f6c63430008170033",
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

// PartialDecrypts is a free data retrieval call binding the contract method 0xdd7b0d8a.
//
// Solidity: function partialDecrypts(bytes32 mrenclave, uint32 round, bytes32 labelHash, uint32 pid) view returns(address validator, bytes partialDecryption, bytes pubShare, bytes label, bool exists)
func (_DKG *DKGCaller) PartialDecrypts(opts *bind.CallOpts, mrenclave [32]byte, round uint32, labelHash [32]byte, pid uint32) (struct {
	Validator         common.Address
	PartialDecryption []byte
	PubShare          []byte
	Label             []byte
	Exists            bool
}, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "partialDecrypts", mrenclave, round, labelHash, pid)

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

// PartialDecrypts is a free data retrieval call binding the contract method 0xdd7b0d8a.
//
// Solidity: function partialDecrypts(bytes32 mrenclave, uint32 round, bytes32 labelHash, uint32 pid) view returns(address validator, bytes partialDecryption, bytes pubShare, bytes label, bool exists)
func (_DKG *DKGSession) PartialDecrypts(mrenclave [32]byte, round uint32, labelHash [32]byte, pid uint32) (struct {
	Validator         common.Address
	PartialDecryption []byte
	PubShare          []byte
	Label             []byte
	Exists            bool
}, error) {
	return _DKG.Contract.PartialDecrypts(&_DKG.CallOpts, mrenclave, round, labelHash, pid)
}

// PartialDecrypts is a free data retrieval call binding the contract method 0xdd7b0d8a.
//
// Solidity: function partialDecrypts(bytes32 mrenclave, uint32 round, bytes32 labelHash, uint32 pid) view returns(address validator, bytes partialDecryption, bytes pubShare, bytes label, bool exists)
func (_DKG *DKGCallerSession) PartialDecrypts(mrenclave [32]byte, round uint32, labelHash [32]byte, pid uint32) (struct {
	Validator         common.Address
	PartialDecryption []byte
	PubShare          []byte
	Label             []byte
	Exists            bool
}, error) {
	return _DKG.Contract.PartialDecrypts(&_DKG.CallOpts, mrenclave, round, labelHash, pid)
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

// FinalizeDKG is a paid mutator transaction binding the contract method 0x9af5962c.
//
// Solidity: function finalizeDKG(uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes[] publicCoeffs, bytes signature) returns()
func (_DKG *DKGTransactor) FinalizeDKG(opts *bind.TransactOpts, round uint32, mrenclave [32]byte, globalPubKey []byte, publicCoeffs [][]byte, signature []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "finalizeDKG", round, mrenclave, globalPubKey, publicCoeffs, signature)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x9af5962c.
//
// Solidity: function finalizeDKG(uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes[] publicCoeffs, bytes signature) returns()
func (_DKG *DKGSession) FinalizeDKG(round uint32, mrenclave [32]byte, globalPubKey []byte, publicCoeffs [][]byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.FinalizeDKG(&_DKG.TransactOpts, round, mrenclave, globalPubKey, publicCoeffs, signature)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x9af5962c.
//
// Solidity: function finalizeDKG(uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes[] publicCoeffs, bytes signature) returns()
func (_DKG *DKGTransactorSession) FinalizeDKG(round uint32, mrenclave [32]byte, globalPubKey []byte, publicCoeffs [][]byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.FinalizeDKG(&_DKG.TransactOpts, round, mrenclave, globalPubKey, publicCoeffs, signature)
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

// SubmitPartialDecryption is a paid mutator transaction binding the contract method 0x79b131a7.
//
// Solidity: function submitPartialDecryption(uint32 round, bytes32 mrenclave, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label) returns()
func (_DKG *DKGTransactor) SubmitPartialDecryption(opts *bind.TransactOpts, round uint32, mrenclave [32]byte, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, label []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "submitPartialDecryption", round, mrenclave, pid, encryptedPartial, ephemeralPubKey, pubShare, label)
}

// SubmitPartialDecryption is a paid mutator transaction binding the contract method 0x79b131a7.
//
// Solidity: function submitPartialDecryption(uint32 round, bytes32 mrenclave, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label) returns()
func (_DKG *DKGSession) SubmitPartialDecryption(round uint32, mrenclave [32]byte, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, label []byte) (*types.Transaction, error) {
	return _DKG.Contract.SubmitPartialDecryption(&_DKG.TransactOpts, round, mrenclave, pid, encryptedPartial, ephemeralPubKey, pubShare, label)
}

// SubmitPartialDecryption is a paid mutator transaction binding the contract method 0x79b131a7.
//
// Solidity: function submitPartialDecryption(uint32 round, bytes32 mrenclave, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label) returns()
func (_DKG *DKGTransactorSession) SubmitPartialDecryption(round uint32, mrenclave [32]byte, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, label []byte) (*types.Transaction, error) {
	return _DKG.Contract.SubmitPartialDecryption(&_DKG.TransactOpts, round, mrenclave, pid, encryptedPartial, ephemeralPubKey, pubShare, label)
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
	PublicCoeffs [][]byte
	Signature    []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDKGFinalized is a free log retrieval operation binding the contract event 0x5d25bc3c675e166c30fb9ab70ab3f79501dd65672130ffdaf3bb412531d9fd82.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes[] publicCoeffs, bytes signature)
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

// WatchDKGFinalized is a free log subscription operation binding the contract event 0x5d25bc3c675e166c30fb9ab70ab3f79501dd65672130ffdaf3bb412531d9fd82.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes[] publicCoeffs, bytes signature)
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

// ParseDKGFinalized is a log parse operation binding the contract event 0x5d25bc3c675e166c30fb9ab70ab3f79501dd65672130ffdaf3bb412531d9fd82.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, bytes32 mrenclave, bytes globalPubKey, bytes[] publicCoeffs, bytes signature)
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
	Pid              uint32
	EncryptedPartial []byte
	EphemeralPubKey  []byte
	PubShare         []byte
	Label            []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterPartialDecryptionSubmitted is a free log retrieval operation binding the contract event 0x835e2245f021610650983a80011abf0755d752f7ce7935d861f90dcfa4ad8db2.
//
// Solidity: event PartialDecryptionSubmitted(address indexed validator, uint32 round, bytes32 mrenclave, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label)
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

// WatchPartialDecryptionSubmitted is a free log subscription operation binding the contract event 0x835e2245f021610650983a80011abf0755d752f7ce7935d861f90dcfa4ad8db2.
//
// Solidity: event PartialDecryptionSubmitted(address indexed validator, uint32 round, bytes32 mrenclave, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label)
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

// ParsePartialDecryptionSubmitted is a log parse operation binding the contract event 0x835e2245f021610650983a80011abf0755d752f7ce7935d861f90dcfa4ad8db2.
//
// Solidity: event PartialDecryptionSubmitted(address indexed validator, uint32 round, bytes32 mrenclave, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label)
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
