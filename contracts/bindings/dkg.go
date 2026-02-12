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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"complainDeals\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"curCodeCommitment\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dealComplaints\",\"inputs\":[{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainant\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dkgNodeInfos\",\"inputs\":[{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"nodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.NodeStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"participantsRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"publicCoeffs\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getNodeInfo\",\"inputs\":[{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDKG.NodeInfo\",\"components\":[{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"nodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.NodeStatus\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"startBlockHeight\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"partialDecrypts\",\"inputs\":[{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"labelHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"pid\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"partialDecryption\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestThresholdDecryption\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitPartialDecryption\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"pid\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ephemeralPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"DKGFinalized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"participantsRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"publicCoeffs\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGInitialized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"startBlockHeight\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"startBlockHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealComplaintsSubmitted\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"indexed\":false,\"internalType\":\"uint32[]\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealVerified\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"recipientIndex\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidDeal\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PartialDecryptionSubmitted\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"pid\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"encryptedPartial\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ephemeralPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"pubShare\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoteAttestationProcessedOnChain\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ThresholdDecryptRequested\",\"inputs\":[{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"requesterPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"label\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeScheduled\",\"inputs\":[{\"name\":\"activationHeight\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"codeCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60803461005057601f61280c38819003918201601f19168301916001600160401b038311848410176100545780849260209460405283398101031261005057515f556040516127a390816100698239f35b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c806308ad63ac146100b45780630a36f7d8146100af5780633995d3b8146100aa57806379b131a7146100a5578063a26f51a4146100a0578063aab066c61461009b578063b1133cac14610096578063dd7b0d8a14610091578063dea942d91461008c5763f0b98a2814610087575f80fd5b610caf565b610be2565b610a6f565b610823565b610762565b61047d565b6103c0565b610386565b6102c9565b346101805760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610180576100eb610189565b6100f361019c565b906044359167ffffffffffffffff908184116101805736602385011215610180578360040135918211610184578160051b9360209460405193610139602083018661024b565b8452602460208501918301019136831161018057602401905b82821061016957610167606435868689610dc0565b005b868091610175846101d5565b815201910190610152565b5f80fd5b6101e6565b6004359063ffffffff8216820361018057565b6024359063ffffffff8216820361018057565b6044359063ffffffff8216820361018057565b6064359063ffffffff8216820361018057565b359063ffffffff8216820361018057565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b60a0810190811067ffffffffffffffff82111761018457604052565b6040810190811067ffffffffffffffff82111761018457604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761018457604052565b6040519061029982610213565b565b9181601f840112156101805782359167ffffffffffffffff8311610180576020838186019501011161018057565b346101805760c07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261018057610300610189565b67ffffffffffffffff6064358181116101805761032190369060040161029b565b60849391933591838311610180573660238401121561018057826004013591848311610180573660248460051b860101116101805760a4359485116101805761016795610374602496369060040161029b565b96909501926044359060243590610fae565b34610180575f7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101805760205f54604051908152f35b346101805760e07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610180576103f7610189565b6103ff6101af565b67ffffffffffffffff916064358381116101805761042190369060040161029b565b906084358581116101805761043a90369060040161029b565b9060a4358781116101805761045390369060040161029b565b94909360c4359889116101805761047161016799369060040161029b565b98909760243590611316565b346101805760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc360112610180576104b461019c565b6104bc6101af565b906064359173ffffffffffffffffffffffffffffffffffffffff83168303610180576020926105216105449261050e60ff956004355f526002885260405f209063ffffffff165f5260205260405f2090565b9063ffffffff165f5260205260405f2090565b9073ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b54166040519015158152f35b7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc6060910112610180576004359060243563ffffffff81168103610180579060443573ffffffffffffffffffffffffffffffffffffffff811681036101805790565b90600182811c921680156105f9575b60208310146105cc57565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f16916105c1565b9060405191825f8254610615816105b2565b908184526020946001916001811690815f146106815750600114610643575b5050506102999250038361024b565b5f90815285812095935091905b81831061066957505061029993508201015f8080610634565b85548884018501529485019487945091830191610650565b9150506102999593507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201015f8080610634565b5f5b8381106106d35750505f910152565b81810151838201526020016106c4565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f60209361071f815180928187528780880191016106c2565b0116010190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b6003111561075d57565b610726565b34610180576107e661079961052161077936610550565b92915f52600160205260405f209063ffffffff165f5260205260405f2090565b6107a281610603565b906107af60018201610603565b61080260036107c060028501610603565b930154926107f460ff8086169560081c169360405197889760a0895260a08901906106e3565b9087820360208901526106e3565b9085820360408701526106e3565b9161080c81610753565b606084015261081a81610753565b60808301520390f35b346101805760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101805761085a610189565b67ffffffffffffffff906024356044358381116101805761087f90369060040161029b565b6064929192358581116101805761089a90369060040161029b565b9095608435908111610180576108b490369060040161029b565b9690926040516108f6602091828101908882528381526108d38161022f565b5190205f546040518481019182528481526108ed8161022f565b51902014610d5b565b61090763ffffffff891615156114f6565b83156109b1576041820361095357509161094e93917f5f9d4b68667a3b91f2fe8369c2ac9040ce4a68400aadbe076f9d109b13e09b61979893604051978897339b89611be3565b0390a2005b606490604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601860248201527f496e76616c696420726571756573746572207075626b657900000000000000006044820152fd5b606490604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601060248201527f456d7074792063697068657274657874000000000000000000000000000000006044820152fd5b9390608093610a4b610a679473ffffffffffffffffffffffffffffffffffffffff610a59949a999a16885260a0602089015260a08801906106e3565b9086820360408801526106e3565b9084820360608601526106e3565b931515910152565b346101805760807ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261018057610af5610aa961019c565b610ad6610ab46101c2565b916004355f52600360205260405f209063ffffffff165f5260205260405f2090565b6044355f5260205260405f209063ffffffff165f5260205260405f2090565b73ffffffffffffffffffffffffffffffffffffffff815416610b4b610b1c60018401610603565b92610b2960028201610603565b9060ff6004610b3a60038401610603565b920154169160405195869586610a0f565b0390f35b6020815260a06080610bbe610b6f855184602087015260c08601906106e3565b610ba96020870151917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe092838883030160408901526106e3565b906040870151908683030160608701526106e3565b936060810151610bcd81610753565b82850152015191610bdd83610753565b015290565b3461018057610b4b610521610c43610bf936610550565b9193906040945f60808751610c0d81610213565b606081526060602082015260608982015282606082015201525f526001602052845f209063ffffffff165f5260205260405f2090565b9060ff6003825193610c5485610213565b610c5d81610603565b8552610c6b60018201610603565b6020860152610c7c60028201610603565b848601520154818116610c8e81610753565b606085015260081c16610ca081610753565b60808301525191829182610b4f565b346101805760e07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261018057610ce6610189565b67ffffffffffffffff9060443582811681036101805760843583811161018057610d1490369060040161029b565b9060a43585811161018057610d2d90369060040161029b565b92909160c43596871161018057610d4b61016797369060040161029b565b9690956064359160243590611c2e565b15610d6257565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f496e76616c696420636f646520636f6d6d69746d656e740000000000000000006044820152fd5b9392936040805190610dfb60209283810190898252848152610de18161022f565b5190205f546040518581019182528581526108ed8161022f565b865f52600191829160018252610e48610e258660405f209063ffffffff165f5260205260405f2090565b3373ffffffffffffffffffffffffffffffffffffffff165f5260205260405f2090565b505f925b610e8d575b505050507fa89000d88bdc9c3e92c10abb67235241f8c6803723e88e1e2420533e8fe2b8d89394610e889160405194859485610f53565b0390a1565b8651831015610f21575f8981526002835281812063ffffffff87168252602052604090208751841015610f1c578493610f15610eea8a610521889563ffffffff8933948860051b0101511663ffffffff165f5260205260405f2090565b60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00825416179055565b0192610e4c565b610f26565b610e51565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b91949392946080830163ffffffff8093168452602090608060208601528251809152602060a086019301915f5b828110610f9857505050509416604082015260600152565b8351861685529381019392810192600101610f80565b969298979395919594909460405199610ff060209b8c8101908982528d8152610fd68161022f565b5190205f54604051808f019182528e81526108ed8161022f565b865f5260018b52611015610e258a60405f209063ffffffff165f5260205260405f2090565b6003810190600160ff83541661102a81610753565b146110c9578a9b9c50829161108a6110858a8a7f5ab3d263439d47a6c2c16ff562f09006fad2885b900ac3edf84f37e0693b28909f9a8f9a8f9a6110c49f9e9d9a839d839d839d61108060016110b69e01610603565b612135565b611127565b6102007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff825416179055565b604051988998339c8a611206565b0390a2565b60648d604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601460248201527f4e6f64652077617320696e76616c6964617465640000000000000000000000006044820152fd5b1561112e57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f496e76616c69642066696e616c697a6174696f6e207369676e617475726500006044820152fd5b9061119681610753565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff0083549260081b169116179055565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe093818652868601375f8582860101520116010190565b94909263ffffffff61123c949b9a989b99979995939516865260209485870152604086015260c0606086015260c08501916111c8565b828103608084015287815296600581901b8801820195915f818a01845b84831061127c5750505050505061127994955060a08185039101526111c8565b90565b9091929394987fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08c820301835289357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18336030181121561018057820185810191903567ffffffffffffffff811161018057803603831361018057611306879283926001956111c8565b9b01930193019194939290611259565b9692989490999591979389898c8a604051916020928381019082825284815261133e8161022f565b5190205f546040518581019182528581526113588161022f565b5190201461136590610d5b565b63ffffffff61137783821615156114f6565b841615156113849061155b565b61138f8615156115c0565b61139a8a1515611625565b6113a58c151561168a565b6113b1604189146116ef565b6113bc368d8d611754565b838151910120928484846113d8855f52600360205260405f2090565b906113f0919063ffffffff165f5260205260405f2090565b5f9182526020526040902090611413919063ffffffff165f5260205260405f2090565b6004015460ff1615611424906117b8565b61142c61028c565b3381529587369061143c92611754565b9086015261144b368b8b611754565b604086015261145b368d8d611754565b606086015260016080860152611479905f52600360205260405f2090565b90611491919063ffffffff165f5260205260405f2090565b5f91825260205260409020906114b4919063ffffffff165f5260205260405f2090565b906114be91611997565b604051998a99339c6114d09a8c611b7b565b037f835e2245f021610650983a80011abf0755d752f7ce7935d861f90dcfa4ad8db291a2565b156114fd57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f496e76616c696420726f756e64000000000000000000000000000000000000006044820152fd5b1561156257565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f496e76616c6964207069640000000000000000000000000000000000000000006044820152fd5b156115c757565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f456d707479207061727469616c000000000000000000000000000000000000006044820152fd5b1561162c57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f456d7074792070756253686172650000000000000000000000000000000000006044820152fd5b1561169157565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f456d707479206c6162656c0000000000000000000000000000000000000000006044820152fd5b156116f657565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f496e76616c696420657068656d6572616c207075626b657900000000000000006044820152fd5b92919267ffffffffffffffff8211610184576040519161179c60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116018461024b565b829481845281830111610180578281602093845f960137010152565b156117bf57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f5061727469616c20616c7265616479207375626d6974746564000000000000006044820152fd5b601f821161182a57505050565b5f5260205f20906020601f840160051c83019310611862575b601f0160051c01905b818110611857575050565b5f815560010161184c565b9091508190611843565b919091825167ffffffffffffffff8111610184576118948161188e84546105b2565b8461181d565b602080601f83116001146118f3575081906118e49394955f926118e8575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b9055565b015190505f806118b2565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0831695611925855f5260205f2090565b925f905b88821061197f57505083600195969710611948575b505050811b019055565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690555f808061193e565b80600185968294968601518155019501930190611929565b9073ffffffffffffffffffffffffffffffffffffffff8151167fffffffffffffffffffffffff000000000000000000000000000000000000000083541617825560018083019060208084015180519267ffffffffffffffff841161018457611a0984611a0387546105b2565b8761181d565b602092601f8511600114611ac357505093600493611a6584611a91956080956102999a995f926118e85750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b611a7960408201516002870161186c565b611a8a60608201516003870161186c565b0151151590565b91019060ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0083541691151516179055565b9291907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0851690611af7875f5260205f2090565b945f915b838310611b64575050508460809461029999989460049894611a919860019510611b2d575b505050811b019055611a68565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690555f8080611b20565b848601518755958601959481019491810191611afb565b99979592936112799b9995611bb992611bc796611bd59a9560208f63ffffffff809516815201521660408d015260e060608d015260e08c01916111c8565b9189830360808b01526111c8565b9186830360a08801526111c8565b9260c08185039101526111c8565b969492611c209461127999979363ffffffff611c1294168a5260208a015260a060408a015260a08901916111c8565b9186830360608801526111c8565b9260808185039101526111c8565b9591969299989490979360409a8b51611c718d8c602093848101918252848152611c578161022f565b519020905f5490518481019182528481526108ed8161022f565b8c611c7d368a8a611754565b611c88368787611754565b90611c94368a8a611754565b9281511115611db957958d9e9f927fe0df76790639f7d48a78e05bf8e2e1a8e5fdc63ee96ef359a4f04bfc509035619e6105218f8f906110c49f9e8f819f9e9d9a8f8f8a9f9e9d918f8f90611d60611dad9f8b611da89f611d5b92611d7096611d849e610190611d799c611d09331515612273565b611d1a63ffffffff881615156122d8565b611d268551151561233d565b611d32865115156123a2565b611d4767ffffffffffffffff84161515612407565b611d5284151561246c565b015194336124e7565b611e3d565b611d6861028c565b9d3691611754565b8c523691611754565b908901523691611754565b858801525f60608601526001608086015261050e33935f52600160205260405f2090565b611ed9565b51998a99339d8b612054565b8f8460849151907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152602260248201527f496e76616c6964207261772071756f74652c2071756f746520746f6f2073686f60448201527f72740000000000000000000000000000000000000000000000000000000000006064820152fd5b15611e4457565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f496e76616c69642072656d6f7465206174746573746174696f6e0000000000006044820152fd5b90611eac81610753565b60ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008354169116179055565b9190805192835167ffffffffffffffff811161018457611efd8161188e84546105b2565b602080601f8311600114611fa3575091611f548260039360809561029998995f926118e85750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b81555b611f6860208501516001830161186c565b611f7960408501516002830161186c565b0191611f926060820151611f8c81610753565b84611ea2565b015190611f9e82610753565b61118c565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0831696611fd5855f5260205f2090565b925f905b89821061203c5750509260809492600192600395836102999a9b10612006575b505050811b018155611f57565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f884881b161c191690555f8080611ff9565b80600185968294968601518155019501930190611fd9565b989694919367ffffffffffffffff611bc79463ffffffff6112799d9b97611bd59a958e521660208d01521660408b015260608a015260e060808a015260e08901916111c8565b9190811015610f1c5760051b810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18136030182121561018057019081359167ffffffffffffffff8311610180576020018236038113610180579190565b6020908361029993959495604051968361211d89955180928880890191016106c2565b84019185830137015f8382015203808552018361024b565b98979694806064939295612196957fffffffff00000000000000000000000000000000000000000000000000000000604051988996602088015260e01b16604086015260448501528484013781015f8382015203604481018452018261024b565b915f915b8083106122505750505061223161221073ffffffffffffffffffffffffffffffffffffffff9461220a61224a9561220286602061223198519101207f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f52601c52603c5f2090565b923691611754565b906124d1565b946020815191012073ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1690565b91161490565b90919261226a60019161226486858761209a565b916120fa565b9301919061219a565b1561227a57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f496e76616c69642076616c696461746f722061646472657373000000000000006044820152fd5b156122df57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f496e76616c696420726f756e64206f6620444b470000000000000000000000006044820152fd5b1561234457565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f496e76616c696420444b47207075626c6963206b6579000000000000000000006044820152fd5b156123a957565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602060248201527f496e76616c696420636f6d6d756e69636174696f6e207075626c6963206b65796044820152fd5b1561240e57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f496e76616c696420737461727420626c6f636b206865696768740000000000006044820152fd5b1561247357565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f496e76616c696420737461727420626c6f636b206861736800000000000000006044820152fd5b611279916124de916125b8565b909291926125fc565b947fffffffff00000000000000000000000000000000000000000000000000000000946125b1947fffffffffffffffff0000000000000000000000000000000000000000000000006060956040519889967fffffffffffffffffffffffffffffffffffffffff000000000000000000000000602089019c8a1b168c5260e01b16603487015260c01b166038850152604084015261258d81518092602087870191016106c2565b82016125a282518093602087850191016106c2565b0103604081018452018261024b565b5190201490565b81519190604183036125e8576125e19250602082015190606060408401519301515f1a906126d3565b9192909190565b50505f9160029190565b6004111561075d57565b612605816125f2565b8061260e575050565b612617816125f2565b600181036126495760046040517ff645eedf000000000000000000000000000000000000000000000000000000008152fd5b612652816125f2565b6002810361268c576040517ffce698f700000000000000000000000000000000000000000000000000000000815260048101839052602490fd5b806126986003926125f2565b146126a05750565b6040517fd78bce0c0000000000000000000000000000000000000000000000000000000081526004810191909152602490fd5b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411612762579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa15612757575f5173ffffffffffffffffffffffffffffffffffffffff81161561274d57905f905f90565b505f906001905f90565b6040513d5f823e3d90fd5b5050505f916003919056fea264697066735822122042f8a4d276ddf62f2cde11fb37b76f69b23d30910e00ef27f7de922a7d8fa13b64736f6c63430008170033",
}

// DKGABI is the input ABI used to generate the binding from.
// Deprecated: Use DKGMetaData.ABI instead.
var DKGABI = DKGMetaData.ABI

// DKGBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DKGMetaData.Bin instead.
var DKGBin = DKGMetaData.Bin

// DeployDKG deploys a new Ethereum contract, binding an instance of DKG to it.
func DeployDKG(auth *bind.TransactOpts, backend bind.ContractBackend, codeCommitment [32]byte) (common.Address, *types.Transaction, *DKG, error) {
	parsed, err := DKGMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DKGBin), backend, codeCommitment)
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

// CurCodeCommitment is a free data retrieval call binding the contract method 0x3995d3b8.
//
// Solidity: function curCodeCommitment() view returns(bytes32)
func (_DKG *DKGCaller) CurCodeCommitment(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "curCodeCommitment")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CurCodeCommitment is a free data retrieval call binding the contract method 0x3995d3b8.
//
// Solidity: function curCodeCommitment() view returns(bytes32)
func (_DKG *DKGSession) CurCodeCommitment() ([32]byte, error) {
	return _DKG.Contract.CurCodeCommitment(&_DKG.CallOpts)
}

// CurCodeCommitment is a free data retrieval call binding the contract method 0x3995d3b8.
//
// Solidity: function curCodeCommitment() view returns(bytes32)
func (_DKG *DKGCallerSession) CurCodeCommitment() ([32]byte, error) {
	return _DKG.Contract.CurCodeCommitment(&_DKG.CallOpts)
}

// DealComplaints is a free data retrieval call binding the contract method 0xa26f51a4.
//
// Solidity: function dealComplaints(bytes32 codeCommitment, uint32 round, uint32 index, address complainant) view returns(bool)
func (_DKG *DKGCaller) DealComplaints(opts *bind.CallOpts, codeCommitment [32]byte, round uint32, index uint32, complainant common.Address) (bool, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "dealComplaints", codeCommitment, round, index, complainant)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DealComplaints is a free data retrieval call binding the contract method 0xa26f51a4.
//
// Solidity: function dealComplaints(bytes32 codeCommitment, uint32 round, uint32 index, address complainant) view returns(bool)
func (_DKG *DKGSession) DealComplaints(codeCommitment [32]byte, round uint32, index uint32, complainant common.Address) (bool, error) {
	return _DKG.Contract.DealComplaints(&_DKG.CallOpts, codeCommitment, round, index, complainant)
}

// DealComplaints is a free data retrieval call binding the contract method 0xa26f51a4.
//
// Solidity: function dealComplaints(bytes32 codeCommitment, uint32 round, uint32 index, address complainant) view returns(bool)
func (_DKG *DKGCallerSession) DealComplaints(codeCommitment [32]byte, round uint32, index uint32, complainant common.Address) (bool, error) {
	return _DKG.Contract.DealComplaints(&_DKG.CallOpts, codeCommitment, round, index, complainant)
}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xaab066c6.
//
// Solidity: function dkgNodeInfos(bytes32 codeCommitment, uint32 round, address validator) view returns(bytes dkgPubKey, bytes commPubKey, bytes rawQuote, uint8 chalStatus, uint8 nodeStatus)
func (_DKG *DKGCaller) DkgNodeInfos(opts *bind.CallOpts, codeCommitment [32]byte, round uint32, validator common.Address) (struct {
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	NodeStatus uint8
}, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "dkgNodeInfos", codeCommitment, round, validator)

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
// Solidity: function dkgNodeInfos(bytes32 codeCommitment, uint32 round, address validator) view returns(bytes dkgPubKey, bytes commPubKey, bytes rawQuote, uint8 chalStatus, uint8 nodeStatus)
func (_DKG *DKGSession) DkgNodeInfos(codeCommitment [32]byte, round uint32, validator common.Address) (struct {
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	NodeStatus uint8
}, error) {
	return _DKG.Contract.DkgNodeInfos(&_DKG.CallOpts, codeCommitment, round, validator)
}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xaab066c6.
//
// Solidity: function dkgNodeInfos(bytes32 codeCommitment, uint32 round, address validator) view returns(bytes dkgPubKey, bytes commPubKey, bytes rawQuote, uint8 chalStatus, uint8 nodeStatus)
func (_DKG *DKGCallerSession) DkgNodeInfos(codeCommitment [32]byte, round uint32, validator common.Address) (struct {
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	NodeStatus uint8
}, error) {
	return _DKG.Contract.DkgNodeInfos(&_DKG.CallOpts, codeCommitment, round, validator)
}

// GetNodeInfo is a free data retrieval call binding the contract method 0xdea942d9.
//
// Solidity: function getNodeInfo(bytes32 codeCommitment, uint32 round, address validator) view returns((bytes,bytes,bytes,uint8,uint8))
func (_DKG *DKGCaller) GetNodeInfo(opts *bind.CallOpts, codeCommitment [32]byte, round uint32, validator common.Address) (IDKGNodeInfo, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "getNodeInfo", codeCommitment, round, validator)

	if err != nil {
		return *new(IDKGNodeInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IDKGNodeInfo)).(*IDKGNodeInfo)

	return out0, err

}

// GetNodeInfo is a free data retrieval call binding the contract method 0xdea942d9.
//
// Solidity: function getNodeInfo(bytes32 codeCommitment, uint32 round, address validator) view returns((bytes,bytes,bytes,uint8,uint8))
func (_DKG *DKGSession) GetNodeInfo(codeCommitment [32]byte, round uint32, validator common.Address) (IDKGNodeInfo, error) {
	return _DKG.Contract.GetNodeInfo(&_DKG.CallOpts, codeCommitment, round, validator)
}

// GetNodeInfo is a free data retrieval call binding the contract method 0xdea942d9.
//
// Solidity: function getNodeInfo(bytes32 codeCommitment, uint32 round, address validator) view returns((bytes,bytes,bytes,uint8,uint8))
func (_DKG *DKGCallerSession) GetNodeInfo(codeCommitment [32]byte, round uint32, validator common.Address) (IDKGNodeInfo, error) {
	return _DKG.Contract.GetNodeInfo(&_DKG.CallOpts, codeCommitment, round, validator)
}

// PartialDecrypts is a free data retrieval call binding the contract method 0xdd7b0d8a.
//
// Solidity: function partialDecrypts(bytes32 codeCommitment, uint32 round, bytes32 labelHash, uint32 pid) view returns(address validator, bytes partialDecryption, bytes pubShare, bytes label, bool exists)
func (_DKG *DKGCaller) PartialDecrypts(opts *bind.CallOpts, codeCommitment [32]byte, round uint32, labelHash [32]byte, pid uint32) (struct {
	Validator         common.Address
	PartialDecryption []byte
	PubShare          []byte
	Label             []byte
	Exists            bool
}, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "partialDecrypts", codeCommitment, round, labelHash, pid)

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
// Solidity: function partialDecrypts(bytes32 codeCommitment, uint32 round, bytes32 labelHash, uint32 pid) view returns(address validator, bytes partialDecryption, bytes pubShare, bytes label, bool exists)
func (_DKG *DKGSession) PartialDecrypts(codeCommitment [32]byte, round uint32, labelHash [32]byte, pid uint32) (struct {
	Validator         common.Address
	PartialDecryption []byte
	PubShare          []byte
	Label             []byte
	Exists            bool
}, error) {
	return _DKG.Contract.PartialDecrypts(&_DKG.CallOpts, codeCommitment, round, labelHash, pid)
}

// PartialDecrypts is a free data retrieval call binding the contract method 0xdd7b0d8a.
//
// Solidity: function partialDecrypts(bytes32 codeCommitment, uint32 round, bytes32 labelHash, uint32 pid) view returns(address validator, bytes partialDecryption, bytes pubShare, bytes label, bool exists)
func (_DKG *DKGCallerSession) PartialDecrypts(codeCommitment [32]byte, round uint32, labelHash [32]byte, pid uint32) (struct {
	Validator         common.Address
	PartialDecryption []byte
	PubShare          []byte
	Label             []byte
	Exists            bool
}, error) {
	return _DKG.Contract.PartialDecrypts(&_DKG.CallOpts, codeCommitment, round, labelHash, pid)
}

// ComplainDeals is a paid mutator transaction binding the contract method 0x08ad63ac.
//
// Solidity: function complainDeals(uint32 round, uint32 index, uint32[] complainIndexes, bytes32 codeCommitment) returns()
func (_DKG *DKGTransactor) ComplainDeals(opts *bind.TransactOpts, round uint32, index uint32, complainIndexes []uint32, codeCommitment [32]byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "complainDeals", round, index, complainIndexes, codeCommitment)
}

// ComplainDeals is a paid mutator transaction binding the contract method 0x08ad63ac.
//
// Solidity: function complainDeals(uint32 round, uint32 index, uint32[] complainIndexes, bytes32 codeCommitment) returns()
func (_DKG *DKGSession) ComplainDeals(round uint32, index uint32, complainIndexes []uint32, codeCommitment [32]byte) (*types.Transaction, error) {
	return _DKG.Contract.ComplainDeals(&_DKG.TransactOpts, round, index, complainIndexes, codeCommitment)
}

// ComplainDeals is a paid mutator transaction binding the contract method 0x08ad63ac.
//
// Solidity: function complainDeals(uint32 round, uint32 index, uint32[] complainIndexes, bytes32 codeCommitment) returns()
func (_DKG *DKGTransactorSession) ComplainDeals(round uint32, index uint32, complainIndexes []uint32, codeCommitment [32]byte) (*types.Transaction, error) {
	return _DKG.Contract.ComplainDeals(&_DKG.TransactOpts, round, index, complainIndexes, codeCommitment)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x0a36f7d8.
//
// Solidity: function finalizeDKG(uint32 round, bytes32 codeCommitment, bytes32 participantsRoot, bytes globalPubKey, bytes[] publicCoeffs, bytes signature) returns()
func (_DKG *DKGTransactor) FinalizeDKG(opts *bind.TransactOpts, round uint32, codeCommitment [32]byte, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, signature []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "finalizeDKG", round, codeCommitment, participantsRoot, globalPubKey, publicCoeffs, signature)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x0a36f7d8.
//
// Solidity: function finalizeDKG(uint32 round, bytes32 codeCommitment, bytes32 participantsRoot, bytes globalPubKey, bytes[] publicCoeffs, bytes signature) returns()
func (_DKG *DKGSession) FinalizeDKG(round uint32, codeCommitment [32]byte, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.FinalizeDKG(&_DKG.TransactOpts, round, codeCommitment, participantsRoot, globalPubKey, publicCoeffs, signature)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x0a36f7d8.
//
// Solidity: function finalizeDKG(uint32 round, bytes32 codeCommitment, bytes32 participantsRoot, bytes globalPubKey, bytes[] publicCoeffs, bytes signature) returns()
func (_DKG *DKGTransactorSession) FinalizeDKG(round uint32, codeCommitment [32]byte, participantsRoot [32]byte, globalPubKey []byte, publicCoeffs [][]byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.FinalizeDKG(&_DKG.TransactOpts, round, codeCommitment, participantsRoot, globalPubKey, publicCoeffs, signature)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0xf0b98a28.
//
// Solidity: function initializeDKG(uint32 round, bytes32 codeCommitment, uint64 startBlockHeight, bytes32 startBlockHash, bytes dkgPubKey, bytes commPubKey, bytes rawQuote) returns()
func (_DKG *DKGTransactor) InitializeDKG(opts *bind.TransactOpts, round uint32, codeCommitment [32]byte, startBlockHeight uint64, startBlockHash [32]byte, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "initializeDKG", round, codeCommitment, startBlockHeight, startBlockHash, dkgPubKey, commPubKey, rawQuote)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0xf0b98a28.
//
// Solidity: function initializeDKG(uint32 round, bytes32 codeCommitment, uint64 startBlockHeight, bytes32 startBlockHash, bytes dkgPubKey, bytes commPubKey, bytes rawQuote) returns()
func (_DKG *DKGSession) InitializeDKG(round uint32, codeCommitment [32]byte, startBlockHeight uint64, startBlockHash [32]byte, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) (*types.Transaction, error) {
	return _DKG.Contract.InitializeDKG(&_DKG.TransactOpts, round, codeCommitment, startBlockHeight, startBlockHash, dkgPubKey, commPubKey, rawQuote)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0xf0b98a28.
//
// Solidity: function initializeDKG(uint32 round, bytes32 codeCommitment, uint64 startBlockHeight, bytes32 startBlockHash, bytes dkgPubKey, bytes commPubKey, bytes rawQuote) returns()
func (_DKG *DKGTransactorSession) InitializeDKG(round uint32, codeCommitment [32]byte, startBlockHeight uint64, startBlockHash [32]byte, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) (*types.Transaction, error) {
	return _DKG.Contract.InitializeDKG(&_DKG.TransactOpts, round, codeCommitment, startBlockHeight, startBlockHash, dkgPubKey, commPubKey, rawQuote)
}

// RequestThresholdDecryption is a paid mutator transaction binding the contract method 0xb1133cac.
//
// Solidity: function requestThresholdDecryption(uint32 round, bytes32 codeCommitment, bytes requesterPubKey, bytes ciphertext, bytes label) returns()
func (_DKG *DKGTransactor) RequestThresholdDecryption(opts *bind.TransactOpts, round uint32, codeCommitment [32]byte, requesterPubKey []byte, ciphertext []byte, label []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "requestThresholdDecryption", round, codeCommitment, requesterPubKey, ciphertext, label)
}

// RequestThresholdDecryption is a paid mutator transaction binding the contract method 0xb1133cac.
//
// Solidity: function requestThresholdDecryption(uint32 round, bytes32 codeCommitment, bytes requesterPubKey, bytes ciphertext, bytes label) returns()
func (_DKG *DKGSession) RequestThresholdDecryption(round uint32, codeCommitment [32]byte, requesterPubKey []byte, ciphertext []byte, label []byte) (*types.Transaction, error) {
	return _DKG.Contract.RequestThresholdDecryption(&_DKG.TransactOpts, round, codeCommitment, requesterPubKey, ciphertext, label)
}

// RequestThresholdDecryption is a paid mutator transaction binding the contract method 0xb1133cac.
//
// Solidity: function requestThresholdDecryption(uint32 round, bytes32 codeCommitment, bytes requesterPubKey, bytes ciphertext, bytes label) returns()
func (_DKG *DKGTransactorSession) RequestThresholdDecryption(round uint32, codeCommitment [32]byte, requesterPubKey []byte, ciphertext []byte, label []byte) (*types.Transaction, error) {
	return _DKG.Contract.RequestThresholdDecryption(&_DKG.TransactOpts, round, codeCommitment, requesterPubKey, ciphertext, label)
}

// SubmitPartialDecryption is a paid mutator transaction binding the contract method 0x79b131a7.
//
// Solidity: function submitPartialDecryption(uint32 round, bytes32 codeCommitment, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label) returns()
func (_DKG *DKGTransactor) SubmitPartialDecryption(opts *bind.TransactOpts, round uint32, codeCommitment [32]byte, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, label []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "submitPartialDecryption", round, codeCommitment, pid, encryptedPartial, ephemeralPubKey, pubShare, label)
}

// SubmitPartialDecryption is a paid mutator transaction binding the contract method 0x79b131a7.
//
// Solidity: function submitPartialDecryption(uint32 round, bytes32 codeCommitment, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label) returns()
func (_DKG *DKGSession) SubmitPartialDecryption(round uint32, codeCommitment [32]byte, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, label []byte) (*types.Transaction, error) {
	return _DKG.Contract.SubmitPartialDecryption(&_DKG.TransactOpts, round, codeCommitment, pid, encryptedPartial, ephemeralPubKey, pubShare, label)
}

// SubmitPartialDecryption is a paid mutator transaction binding the contract method 0x79b131a7.
//
// Solidity: function submitPartialDecryption(uint32 round, bytes32 codeCommitment, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label) returns()
func (_DKG *DKGTransactorSession) SubmitPartialDecryption(round uint32, codeCommitment [32]byte, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, label []byte) (*types.Transaction, error) {
	return _DKG.Contract.SubmitPartialDecryption(&_DKG.TransactOpts, round, codeCommitment, pid, encryptedPartial, ephemeralPubKey, pubShare, label)
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
	MsgSender        common.Address
	Round            uint32
	CodeCommitment   [32]byte
	ParticipantsRoot [32]byte
	GlobalPubKey     []byte
	PublicCoeffs     [][]byte
	Signature        []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterDKGFinalized is a free log retrieval operation binding the contract event 0x5ab3d263439d47a6c2c16ff562f09006fad2885b900ac3edf84f37e0693b2890.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, bytes32 codeCommitment, bytes32 participantsRoot, bytes globalPubKey, bytes[] publicCoeffs, bytes signature)
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

// WatchDKGFinalized is a free log subscription operation binding the contract event 0x5ab3d263439d47a6c2c16ff562f09006fad2885b900ac3edf84f37e0693b2890.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, bytes32 codeCommitment, bytes32 participantsRoot, bytes globalPubKey, bytes[] publicCoeffs, bytes signature)
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

// ParseDKGFinalized is a log parse operation binding the contract event 0x5ab3d263439d47a6c2c16ff562f09006fad2885b900ac3edf84f37e0693b2890.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, bytes32 codeCommitment, bytes32 participantsRoot, bytes globalPubKey, bytes[] publicCoeffs, bytes signature)
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
	MsgSender        common.Address
	CodeCommitment   [32]byte
	Round            uint32
	StartBlockHeight uint64
	StartBlockHash   [32]byte
	DkgPubKey        []byte
	CommPubKey       []byte
	RawQuote         []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterDKGInitialized is a free log retrieval operation binding the contract event 0xe0df76790639f7d48a78e05bf8e2e1a8e5fdc63ee96ef359a4f04bfc50903561.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes32 codeCommitment, uint32 round, uint64 startBlockHeight, bytes32 startBlockHash, bytes dkgPubKey, bytes commPubKey, bytes rawQuote)
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

// WatchDKGInitialized is a free log subscription operation binding the contract event 0xe0df76790639f7d48a78e05bf8e2e1a8e5fdc63ee96ef359a4f04bfc50903561.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes32 codeCommitment, uint32 round, uint64 startBlockHeight, bytes32 startBlockHash, bytes dkgPubKey, bytes commPubKey, bytes rawQuote)
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

// ParseDKGInitialized is a log parse operation binding the contract event 0xe0df76790639f7d48a78e05bf8e2e1a8e5fdc63ee96ef359a4f04bfc50903561.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes32 codeCommitment, uint32 round, uint64 startBlockHeight, bytes32 startBlockHash, bytes dkgPubKey, bytes commPubKey, bytes rawQuote)
func (_DKG *DKGFilterer) ParseDKGInitialized(log types.Log) (*DKGDKGInitialized, error) {
	event := new(DKGDKGInitialized)
	if err := _DKG.contract.UnpackLog(event, "DKGInitialized", log); err != nil {
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
	CodeCommitment  [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterDealComplaintsSubmitted is a free log retrieval operation binding the contract event 0xa89000d88bdc9c3e92c10abb67235241f8c6803723e88e1e2420533e8fe2b8d8.
//
// Solidity: event DealComplaintsSubmitted(uint32 index, uint32[] complainIndexes, uint32 round, bytes32 codeCommitment)
func (_DKG *DKGFilterer) FilterDealComplaintsSubmitted(opts *bind.FilterOpts) (*DKGDealComplaintsSubmittedIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "DealComplaintsSubmitted")
	if err != nil {
		return nil, err
	}
	return &DKGDealComplaintsSubmittedIterator{contract: _DKG.contract, event: "DealComplaintsSubmitted", logs: logs, sub: sub}, nil
}

// WatchDealComplaintsSubmitted is a free log subscription operation binding the contract event 0xa89000d88bdc9c3e92c10abb67235241f8c6803723e88e1e2420533e8fe2b8d8.
//
// Solidity: event DealComplaintsSubmitted(uint32 index, uint32[] complainIndexes, uint32 round, bytes32 codeCommitment)
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
// Solidity: event DealComplaintsSubmitted(uint32 index, uint32[] complainIndexes, uint32 round, bytes32 codeCommitment)
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
	CodeCommitment [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDealVerified is a free log retrieval operation binding the contract event 0x1a8f868a6f5b289bec8c24a7c28727ddc869adf4a3c3c0ae8a2d41b9afc345bb.
//
// Solidity: event DealVerified(uint32 index, uint32 recipientIndex, uint32 round, bytes32 codeCommitment)
func (_DKG *DKGFilterer) FilterDealVerified(opts *bind.FilterOpts) (*DKGDealVerifiedIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "DealVerified")
	if err != nil {
		return nil, err
	}
	return &DKGDealVerifiedIterator{contract: _DKG.contract, event: "DealVerified", logs: logs, sub: sub}, nil
}

// WatchDealVerified is a free log subscription operation binding the contract event 0x1a8f868a6f5b289bec8c24a7c28727ddc869adf4a3c3c0ae8a2d41b9afc345bb.
//
// Solidity: event DealVerified(uint32 index, uint32 recipientIndex, uint32 round, bytes32 codeCommitment)
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
// Solidity: event DealVerified(uint32 index, uint32 recipientIndex, uint32 round, bytes32 codeCommitment)
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
	Index          uint32
	Round          uint32
	CodeCommitment [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterInvalidDeal is a free log retrieval operation binding the contract event 0x90395d01853e3de18e643761e8429ec973c5b4843dbf47451c4e90f37c3447ca.
//
// Solidity: event InvalidDeal(uint32 index, uint32 round, bytes32 codeCommitment)
func (_DKG *DKGFilterer) FilterInvalidDeal(opts *bind.FilterOpts) (*DKGInvalidDealIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "InvalidDeal")
	if err != nil {
		return nil, err
	}
	return &DKGInvalidDealIterator{contract: _DKG.contract, event: "InvalidDeal", logs: logs, sub: sub}, nil
}

// WatchInvalidDeal is a free log subscription operation binding the contract event 0x90395d01853e3de18e643761e8429ec973c5b4843dbf47451c4e90f37c3447ca.
//
// Solidity: event InvalidDeal(uint32 index, uint32 round, bytes32 codeCommitment)
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
// Solidity: event InvalidDeal(uint32 index, uint32 round, bytes32 codeCommitment)
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
	CodeCommitment   [32]byte
	Pid              uint32
	EncryptedPartial []byte
	EphemeralPubKey  []byte
	PubShare         []byte
	Label            []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterPartialDecryptionSubmitted is a free log retrieval operation binding the contract event 0x835e2245f021610650983a80011abf0755d752f7ce7935d861f90dcfa4ad8db2.
//
// Solidity: event PartialDecryptionSubmitted(address indexed validator, uint32 round, bytes32 codeCommitment, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label)
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
// Solidity: event PartialDecryptionSubmitted(address indexed validator, uint32 round, bytes32 codeCommitment, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label)
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
// Solidity: event PartialDecryptionSubmitted(address indexed validator, uint32 round, bytes32 codeCommitment, uint32 pid, bytes encryptedPartial, bytes ephemeralPubKey, bytes pubShare, bytes label)
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
	Validator      common.Address
	ChalStatus     uint8
	Round          uint32
	CodeCommitment [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRemoteAttestationProcessedOnChain is a free log retrieval operation binding the contract event 0x54690f98c0ec0056e0e487f4fe5e8eea7bee88d2dbb7cc9ddca22981f06d9dbb.
//
// Solidity: event RemoteAttestationProcessedOnChain(address validator, uint8 chalStatus, uint32 round, bytes32 codeCommitment)
func (_DKG *DKGFilterer) FilterRemoteAttestationProcessedOnChain(opts *bind.FilterOpts) (*DKGRemoteAttestationProcessedOnChainIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "RemoteAttestationProcessedOnChain")
	if err != nil {
		return nil, err
	}
	return &DKGRemoteAttestationProcessedOnChainIterator{contract: _DKG.contract, event: "RemoteAttestationProcessedOnChain", logs: logs, sub: sub}, nil
}

// WatchRemoteAttestationProcessedOnChain is a free log subscription operation binding the contract event 0x54690f98c0ec0056e0e487f4fe5e8eea7bee88d2dbb7cc9ddca22981f06d9dbb.
//
// Solidity: event RemoteAttestationProcessedOnChain(address validator, uint8 chalStatus, uint32 round, bytes32 codeCommitment)
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
// Solidity: event RemoteAttestationProcessedOnChain(address validator, uint8 chalStatus, uint32 round, bytes32 codeCommitment)
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
	CodeCommitment  [32]byte
	RequesterPubKey []byte
	Ciphertext      []byte
	Label           []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterThresholdDecryptRequested is a free log retrieval operation binding the contract event 0x5f9d4b68667a3b91f2fe8369c2ac9040ce4a68400aadbe076f9d109b13e09b61.
//
// Solidity: event ThresholdDecryptRequested(address indexed requester, uint32 round, bytes32 codeCommitment, bytes requesterPubKey, bytes ciphertext, bytes label)
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
// Solidity: event ThresholdDecryptRequested(address indexed requester, uint32 round, bytes32 codeCommitment, bytes requesterPubKey, bytes ciphertext, bytes label)
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
// Solidity: event ThresholdDecryptRequested(address indexed requester, uint32 round, bytes32 codeCommitment, bytes requesterPubKey, bytes ciphertext, bytes label)
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
	CodeCommitment   [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpgradeScheduled is a free log retrieval operation binding the contract event 0xba889db5cdb62d54ea6ab3c85ea27c52b2710b39cb8ecddf1a360a51cbb40110.
//
// Solidity: event UpgradeScheduled(uint32 activationHeight, bytes32 codeCommitment)
func (_DKG *DKGFilterer) FilterUpgradeScheduled(opts *bind.FilterOpts) (*DKGUpgradeScheduledIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "UpgradeScheduled")
	if err != nil {
		return nil, err
	}
	return &DKGUpgradeScheduledIterator{contract: _DKG.contract, event: "UpgradeScheduled", logs: logs, sub: sub}, nil
}

// WatchUpgradeScheduled is a free log subscription operation binding the contract event 0xba889db5cdb62d54ea6ab3c85ea27c52b2710b39cb8ecddf1a360a51cbb40110.
//
// Solidity: event UpgradeScheduled(uint32 activationHeight, bytes32 codeCommitment)
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
// Solidity: event UpgradeScheduled(uint32 activationHeight, bytes32 codeCommitment)
func (_DKG *DKGFilterer) ParseUpgradeScheduled(log types.Log) (*DKGUpgradeScheduled, error) {
	event := new(DKGUpgradeScheduled)
	if err := _DKG.contract.UnpackLog(event, "UpgradeScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
