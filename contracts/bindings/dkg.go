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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"complainDeals\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"curMrenclave\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dealComplaints\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainant\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dkgNodeInfos\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"nodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.NodeStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"publicCoeffs\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getNodeInfo\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDKG.NodeInfo\",\"components\":[{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"nodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.NodeStatus\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestRemoteAttestationOnChain\",\"inputs\":[{\"name\":\"targetValidatorAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNetwork\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"DKGFinalized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"publicCoeffs\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGInitialized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGNetworkSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"total\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealComplaintsSubmitted\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"indexed\":false,\"internalType\":\"uint32[]\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealVerified\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"recipientIndex\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidDeal\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoteAttestationProcessedOnChain\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeScheduled\",\"inputs\":[{\"name\":\"activationHeight\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60803461005057601f61217a38819003918201601f19168301916001600160401b038311848410176100545780849260209460405283398101031261005057515f5560405161211190816100698239f35b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c806308ad63ac146100a4578063227cd9221461009f578063681a0fa81461009a5780639af5962c14610095578063a26f51a414610090578063aab066c61461008b578063b1888cd314610086578063dea942d9146100815763fa4e9f631461007c575f80fd5b610d44565b610c75565b6109de565b61091e565b610622565b610546565b61050c565b6102c2565b346101705760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610170576100db610179565b6100e361018c565b906044359167ffffffffffffffff908184116101705736602385011215610170578360040135918211610174578160051b93602094604051936101296020830186610244565b8452602460208501918301019136831161017057602401905b82821061015957610157606435868689610fa5565b005b868091610165846101b2565b815201910190610142565b5f80fd5b6101c3565b6004359063ffffffff8216820361017057565b6024359063ffffffff8216820361017057565b6044359063ffffffff8216820361017057565b359063ffffffff8216820361017057565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b60a0810190811067ffffffffffffffff82111761017457604052565b6040810190811067ffffffffffffffff82111761017457604052565b6060810190811067ffffffffffffffff82111761017457604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761017457604052565b60405190610292826101f0565b565b9181601f840112156101705782359167ffffffffffffffff8311610170576020838186019501011161017057565b346101705760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610170576102f9610179565b61030161018c565b9061030a61019f565b6064359160843567ffffffffffffffff81116101705761032e903690600401610294565b9060408051602090818101908882528281526103498161020c565b5190205f5483518381019182528381526103628161020c565b5190201461036f90610f40565b5f8781526001825282812063ffffffff87168252602090815260408083203384529091529020906103a2600183016107a8565b8351828101908a82527fffffffff00000000000000000000000000000000000000000000000000000000808a60e01b1687830152808d60e01b1660448301528a60e01b166048820152602c81526103f881610228565b51902061042c907f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f52601c52603c5f2090565b61043736888861158c565b61044091611e6d565b8151919092012073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166104a99173ffffffffffffffffffffffffffffffffffffffff1614611193565b60030180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166102001790555194859433976104e69587611272565b037fc7a37268197965e156b6d53085e9e20ba69f731868b09d00c2b2c3925f25f4f891a2005b34610170575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101705760205f54604051908152f35b346101705760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101705761057d610179565b67ffffffffffffffff6044358181116101705761059e903690600401610294565b60649391933591838311610170573660238401121561017057826004013591848311610170573660248460051b860101116101705760843594851161017057610157956105f16024963690600401610294565b9690950192602435906112a9565b6004359073ffffffffffffffffffffffffffffffffffffffff8216820361017057565b346101705760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101705761065961018c565b61066161019f565b906064359173ffffffffffffffffffffffffffffffffffffffff83168303610170576020926106c66106e9926106b360ff956004355f526002885260405f209063ffffffff165f5260205260405f2090565b9063ffffffff165f5260205260405f2090565b9073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b54166040519015158152f35b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc6060910112610170576004359060243563ffffffff81168103610170579060443573ffffffffffffffffffffffffffffffffffffffff811681036101705790565b90600182811c9216801561079e575b602083101461077157565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f1691610766565b9060405191825f82546107ba81610757565b908184526020946001916001811690815f1461082657506001146107e8575b50505061029292500383610244565b5f90815285812095935091905b81831061080e57505061029293508201015f80806107d9565b855488840185015294850194879450918301916107f5565b9150506102929593507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201015f80806107d9565b5f5b8381106108785750505f910152565b8181015183820152602001610869565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f6020936108c481518092818752878088019101610867565b0116010190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b6003111561090257565b6108cb565b9060038210156109025752565b6004111561090257565b34610170576109a26109556106c6610935366106f5565b92915f52600160205260405f209063ffffffff165f5260205260405f2090565b61095e816107a8565b9061096b600182016107a8565b6109cc6109be600361097f600286016107a8565b940154936109b060ff8660081c169460405198899860a08a5260a08a0190610888565b9088820360208a0152610888565b908682036040880152610888565b9260ff606086019116610907565b6109d581610914565b60808301520390f35b346101705760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261017057610a15610179565b60243567ffffffffffffffff60443581811161017057610a39903690600401610294565b93909160643581811161017057610a54903690600401610294565b9160843590811161017057610a6d903690600401610294565b92909160405197610ab06020998a8101908a82528b8152610a8d8161020c565b5190205f546040518c81019182528c8152610aa78161020c565b51902014610f40565b610add610abe36878761158c565b610ac936848b61158c565b88610ad536888861158c565b923390611d81565b15610b855790610b8094939291610b727f1bd0faa06edbfccdd0f51f46517f5bae23b4abca2dad81e938e89f3ddf7cab1d999a610b18610285565b90610b2436858d61158c565b8252610b3136878761158c565b90820152610b4036888861158c565b60408201525f606082015260016080820152610b6d8c6106c68b6106b333935f52600160205260405f2090565b6117ad565b604051978897339b89611928565b0390a2005b606489604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601a60248201527f496e76616c69642072656d6f7465206174746573746174696f6e0000000000006044820152fd5b6020815260a06080610c52610c03855184602087015260c0860190610888565b610c3d6020870151917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09283888303016040890152610888565b90604087015190868303016060870152610888565b93610c64606082015183860190610907565b015191610c7083610914565b015290565b3461017057610d406106c6610cd6610c8c366106f5565b9193906040945f60808751610ca0816101f0565b606081526060602082015260608982015282606082015201525f526001602052845f209063ffffffff165f5260205260405f2090565b9060ff6003825193610ce7856101f0565b610cf0816107a8565b8552610cfe600182016107a8565b6020860152610d0f600282016107a8565b848601520154610d24828216606086016115f0565b60081c16610d3181610914565b60808301525191829182610be3565b0390f35b346101705760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261017057610d7b6105ff565b610d8361018c565b60443591604051610dbd60209182810190868252838152610da38161020c565b5190205f54604051848101918252848152610aa78161020c565b835f5260018152610de3826106c68560405f209063ffffffff165f5260205260405f2090565b90610dee8254610757565b15610ee257507f54690f98c0ec0056e0e487f4fe5e8eea7bee88d2dbb7cc9ddca22981f06d9dbb93610ea282610e6c6003610eaf950191610e42610e33845460ff1690565b610e3c816108f8565b15611973565b610e4e600282016107a8565b908888610e666001610e5f856107a8565b94016107a8565b93611d81565b15610eb45780547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660021781555b5460ff1690565b93604051948594856119d8565b0390a1005b80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001178155610e9b565b606490604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601360248201527f4e6f646520646f6573206e6f74206578697374000000000000000000000000006044820152fd5b15610f4757565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f496e76616c6964206d72656e636c6176650000000000000000000000000000006044820152fd5b9392936040805190610fe060209283810190898252848152610fc68161020c565b5190205f54604051858101918252858152610aa78161020c565b865f5260019182916001825261102d61100a8660405f209063ffffffff165f5260205260405f2090565b3373ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b505f925b611072575b505050507fa89000d88bdc9c3e92c10abb67235241f8c6803723e88e1e2420533e8fe2b8d8939461106d9160405194859485611138565b0390a1565b8651831015611106575f8981526002835281812063ffffffff871682526020526040902087518410156111015784936110fa6110cf8a6106c6889563ffffffff8933948860051b0101511663ffffffff165f5260205260405f2090565b60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00825416179055565b0192611031565b61110b565b611036565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b91949392946080830163ffffffff8093168452602090608060208601528251809152602060a086019301915f5b82811061117d57505050509416604082015260600152565b8351861685529381019392810192600101611165565b1561119a57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f496e76616c696420736574206e6574776f726b207369676e61747572650000006044820152fd5b9061120281610914565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff0083549260081b169116179055565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe093818652868601375f8582860101520116010190565b9160a0936112a697959263ffffffff9283809216865216602085015216604083015260608201528160808201520191611234565b90565b96919495909392956112e86040516020810190878252602081526112cc8161020c565b5190205f546040516020810191825260208152610aa78161020c565b845f52600160205261130e61100a8960405f209063ffffffff165f5260205260405f2090565b93600385019760ff8954166003811015610902577f5d25bc3c675e166c30fb9ab70ab3f79501dd65672130ffdaf3bb412531d9fd829961137e61137988888f888f8f906113b89f918d9461137460018e9561136e6113aa9f8314156113bd565b016107a8565b611ab4565b611422565b6103007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff825416179055565b604051978897339b89611487565b0390a2565b156113c457565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4e6f64652077617320696e76616c6964617465640000000000000000000000006044820152fd5b1561142957565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f496e76616c69642066696e616c697a6174696f6e207369676e617475726500006044820152fd5b9363ffffffff6114b59394929a99979a9896981685526020938486015260a0604086015260a0850191611234565b828103606084015287815296600581901b8801820195915f818a01845b8483106114f2575050505050506112a69495506080818503910152611234565b9091929394987fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08c820301835289357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18336030181121561017057820185810191903567ffffffffffffffff81116101705780360383136101705761157c87928392600195611234565b9b019301930191949392906114d2565b92919267ffffffffffffffff821161017457604051916115d460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160184610244565b829481845281830111610170578281602093845f960137010152565b60038210156109025752565b601f821161160957505050565b5f5260205f20906020601f840160051c83019310611641575b601f0160051c01905b818110611636575050565b5f815560010161162b565b9091508190611622565b919091825167ffffffffffffffff8111610174576116738161166d8454610757565b846115fc565b602080601f83116001146116d2575081906116c39394955f926116c7575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9055565b015190505f80611691565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0831695611704855f5260205f2090565b925f905b88821061175e57505083600195969710611727575b505050811b019055565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690555f808061171d565b80600185968294968601518155019501930190611708565b9060038110156109025760ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008354169116179055565b9190805192835167ffffffffffffffff8111610174576117d18161166d8454610757565b602080601f83116001146118775750916118288260039360809561029298995f926116c75750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b81555b61183c60208501516001830161164b565b61184d60408501516002830161164b565b01916118666060820151611860816108f8565b84611776565b01519061187282610914565b6111f8565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08316966118a9855f5260205f2090565b925f905b8982106119105750509260809492600192600395836102929a9b106118da575b505050811b01815561182b565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f884881b161c191690555f80806118cd565b806001859682949686015181550195019301906118ad565b9694926119659463ffffffff6112a69a9894611957948b521660208a015260a060408a015260a0890191611234565b918683036060880152611234565b926080818503910152611234565b1561197a57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4e6f646520616c7265616479206368616c6c656e6765640000000000000000006044820152fd5b9094939263ffffffff90611a0f60609473ffffffffffffffffffffffffffffffffffffffff60808601991685526020850190610907565b1660408201520152565b91908110156111015760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18136030182121561017057019081359167ffffffffffffffff8311610170576020018236038113610170579190565b60209083610292939594956040519683611a9c8995518092888089019101610867565b84019185830137015f83820152038085520183610244565b979695939182611b1093604493957fffffffff00000000000000000000000000000000000000000000000000000000604051978895602087015260e01b1660408501528484013781015f83820152036024810184520182610244565b915f915b808310611bca57505050611bab611b8a73ffffffffffffffffffffffffffffffffffffffff94611b84611bc495611b7c866020611bab98519101207f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f52601c52603c5f2090565b92369161158c565b90611e6d565b946020815191012073ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1690565b91161490565b909192611be4600191611bde868587611a19565b91611a79565b93019190611b14565b15611bf457565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f496e76616c69642076616c696461746f722061646472657373000000000000006044820152fd5b15611c5957565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f496e76616c696420726f756e64206f6620444b470000000000000000000000006044820152fd5b15611cbe57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f496e76616c696420444b47207075626c6963206b6579000000000000000000006044820152fd5b15611d2357565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602060248201527f496e76616c696420636f6d6d756e69636174696f6e207075626c6963206b65796044820152fd5b9391929092604085511115611de9576101906112a695611db873ffffffffffffffffffffffffffffffffffffffff87161515611bed565b611dc963ffffffff84161515611c52565b611dd584511515611cb7565b611de185511515611d1c565b015193611e83565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f496e76616c6964207261772071756f74652c2071756f746520746f6f2073686f60448201527f72740000000000000000000000000000000000000000000000000000000000006064820152fd5b6112a691611e7a91611f30565b90929192611f6a565b92611f1f916038917fffffffff00000000000000000000000000000000000000000000000000000000946040519586937fffffffffffffffffffffffffffffffffffffffff000000000000000000000000602086019960601b16895260e01b166034840152611efb8151809260208787019101610867565b8201611f108251809360208785019101610867565b01036018810184520182610244565b51902003611f2c57600190565b5f90565b8151919060418303611f6057611f599250602082015190606060408401519301515f1a90612041565b9192909190565b50505f9160029190565b611f7381610914565b80611f7c575050565b611f8581610914565b60018103611fb75760046040517ff645eedf000000000000000000000000000000000000000000000000000000008152fd5b611fc081610914565b60028103611ffa576040517ffce698f700000000000000000000000000000000000000000000000000000000815260048101839052602490fd5b80612006600392610914565b1461200e5750565b6040517fd78bce0c0000000000000000000000000000000000000000000000000000000081526004810191909152602490fd5b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a084116120d0579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa156120c5575f5173ffffffffffffffffffffffffffffffffffffffff8116156120bb57905f905f90565b505f906001905f90565b6040513d5f823e3d90fd5b5050505f916003919056fea26469706673582212203a782eec16d222e7b78f5cf08bab209d652f10e9a365d1aefdb0243ab068a0dd64736f6c63430008170033",
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
