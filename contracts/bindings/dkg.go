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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"complainDeals\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"curMrenclave\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dealComplaints\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainant\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dkgNodeInfos\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"nodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.NodeStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getGlobalPubKey\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeInfo\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDKG.NodeInfo\",\"components\":[{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"nodeStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.NodeStatus\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isActiveValidator\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestRemoteAttestationOnChain\",\"inputs\":[{\"name\":\"targetValidatorAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"roundInfo\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setNetwork\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitActiveValSet\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"valSet\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"valSets\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"votes\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"globalPubKeyCandidates\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"votes\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DKGFinalized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGInitialized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGNetworkSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"total\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealComplaintsSubmitted\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"indexed\":false,\"internalType\":\"uint32[]\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealVerified\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"recipientIndex\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidDeal\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoteAttestationProcessedOnChain\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeScheduled\",\"inputs\":[{\"name\":\"activationHeight\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x6080604052346200020357620024e5803803806200001d8162000207565b92833981019060208082840312620002035781516001600160401b039283821162000203570192601f9080828601121562000203578451848111620001da57601f1995620000718285018816860162000207565b92828452858383010111620002035784905f5b838110620001ee5750505f91830101528051938411620001da575f54926001938481811c91168015620001cf575b82821014620001bb5783811162000173575b50809285116001146200010e5750839450908392915f9462000102575b50501b915f199060031b1c1916175f555b6040516122b790816200022e8239f35b015192505f80620000e1565b9294849081165f8052845f20945f905b888383106200015857505050106200013f575b505050811b015f55620000f2565b01515f1960f88460031b161c191690555f808062000131565b8587015188559096019594850194879350908101906200011e565b5f8052815f208480880160051c820192848910620001b1575b0160051c019085905b828110620001a5575050620000c4565b5f815501859062000195565b925081926200018c565b634e487b7160e01b5f52602260045260245ffd5b90607f1690620000b2565b634e487b7160e01b5f52604160045260245ffd5b81810183015185820184015286920162000084565b5f80fd5b6040519190601f01601f191682016001600160401b03811183821017620001da5760405256fe60806040526004361015610011575f80fd5b5f3560e01c806302372b2e146119c5578063285f0e63146119555780632e266b4c146117995780632ef619b21461151c5780633090f41114610f635780633c3acdb414610eb15780633ce1c1e014610e0e5780634f1bf88d14610db1578063681a0fa814610d7e5780636ddc4a2514610c235780637a574765146109b9578063bf45f2c3146108c8578063c2009b9214610812578063c3c77ff91461012a5763fb5a783a146100be575f80fd5b346101265763ffffffff60206100d336611eae565b94928260409592955193849283378101600281520301902091165f5260205273ffffffffffffffffffffffffffffffffffffffff60405f2091165f52602052602060ff60405f2054166040519015158152f35b5f80fd5b346101265760a060031936011261012657610143611b0e565b60243567ffffffffffffffff811161012657610163903690600401611b47565b60449291923567ffffffffffffffff811161012657610186903690600401611b47565b909260643567ffffffffffffffff8111610126576101a8903690600401611b47565b9560843567ffffffffffffffff8111610126576101c9903690600401611b47565b9290976101f66101da368986611bfa565b602081519101206101e9611caf565b6020815191012014611f26565b60405187848237602081898101600281520301902063ffffffff86165f5260205260405f20335f5260205260ff60405f205416801561080a575b156107ac5761026261024336868c611bfa565b61024e36898c611bfa565b8761025a368688611bfa565b923390612071565b1561074e576040519461027486611bbb565b61027f36888b611bfa565b865261028c368385611bfa565b602087015261029c36868c611bfa565b60408701525f606087015260016080870152604051888582376020818a8101600181520301902063ffffffff82165f5260205260405f20335f5260205260405f20865180519067ffffffffffffffff821161060257610305826102ff8554611c5e565b85612022565b602090601f83116001146106c75761035292915f918361062f575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b81555b602087015180519067ffffffffffffffff8211610602576103868261037d6001860154611c5e565b60018601612022565b602090601f831160011461063a576103d292915f918361062f5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60018201555b604087015180519067ffffffffffffffff821161060257610409826104006002860154611c5e565b60028601612022565b602090601f831160011461057557918061045992600395945f9261056a5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60028201555b0190606087015197600389101561053d577fe8461cf9f1302c7b5cf4074f825c58d73aa03cca3c607cfcf9d90a137e5f2f0f9a63ffffffff610506610528988660a09d6105389d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000060ff61ff00608061051a9d549401516104df81611f1c565b6104e881611f1c565b60081b1693169116171790556040519d8d8f9e8f9081520191611fa4565b931660208b015289830360408b0152611fa4565b918683036060880152611fa4565b9083820360808501523396611fa4565b0390a2005b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b015190508f80610320565b90600284015f5260205f20915f5b601f19851681106105ea575091839160019383601f196003989716106105b4575b505050811b01600282015561045f565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f884881b161c191690558e80806105a4565b91926020600181928685015181550194019201610583565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b015190508e80610320565b9190600184015f5260205f20905f935b601f19841685106106ac576001945083601f19811610610675575b505050811b0160018201556103d8565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558d8080610665565b8181015183556020948501946001909301929091019061064a565b9190835f5260205f20905f935b601f1984168510610733576001945083601f198116106106fc575b505050811b018155610355565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558d80806106ef565b818101518355602094850194600190930192909101906106d4565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f496e76616c69642072656d6f7465206174746573746174696f6e0000000000006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f56616c696461746f72206e6f7420696e206163746976652073657400000000006044820152fd5b506001610230565b346101265760806003193601126101265760043567ffffffffffffffff811161012657610843903690600401611c40565b61084b611b21565b610853611b34565b6064359173ffffffffffffffffffffffffffffffffffffffff831680930361012657604051848180965161088d816020998a809601611b75565b8101600381520301902063ffffffff8092165f52845260405f2091165f52825260405f20905f52815260ff60405f2054166040519015158152f35b34610126576108f8602063ffffffff6108e036611e4e565b94918460409592955193828580945193849201611b75565b8101600181520301902091165f5260205273ffffffffffffffffffffffffffffffffffffffff60405f2091165f5260205261097d60405f2061093981611d8f565b9061094660018201611d8f565b6109a7610999600361095a60028601611d8f565b9401549361098b60ff8660081c169460405198899860a08a5260a08a0190611b96565b9088820360208a0152611b96565b908682036040880152611b96565b9260ff606086019116611f0f565b6109b081611f1c565b60808301520390f35b34610126576080600319360112610126576109d2611b0e565b6109da611b21565b6044359167ffffffffffffffff908184116101265736602385011215610126578360040135926024838511610602576005918560051b9460209760405197610a248a89018a611bd7565b8852888801602481988301019136831161012657602401905b828210610c075750505060643590811161012657610a649097949291973690600401611b47565b979094610a8f610a75368b89611bfa565b828151910120610a83611caf565b83815191012014611f26565b604051928987853789840193828160019660018152030190209563ffffffff80961696875f52835260405f20335f528352845f905b610b4b575b505050604051978460808a01931689526080828a01525180925260a0880196935f905b838210610b34577f1c2112af5fd37661e3dd248d701decebf291f13de6a411176337ef21a7b1a6308a80610b2f8c8f8d8d60408601528483036060860152611fa4565b0390a1005b855181168952978201979482019490840190610aec565b8a969394959651811015610bfd5760405185818e808d833781016003815203019020885f52855260405f208b51821015610bd15781831b8c0186015185165f90815290865260408082203383528752902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016881790559295949392850185610ac4565b837f4e487b71000000000000000000000000000000000000000000000000000000005f5260326004525ffd5b9594939295610ac9565b813563ffffffff81168103610126578152908a01908a01610a3d565b346101265763ffffffff6080610c3836611eae565b92949091604051610c4881611bbb565b6060968782525f602097838a8a809601528a6040820152828b82015201528260405193849283378101600181520301902091165f52825273ffffffffffffffffffffffffffffffffffffffff60405f2091165f52815260405f2060405191610caf83611bbb565b610cb882611d8f565b8352610cc660018301611d8f565b918184019283526003610cdb60028301611d8f565b9160408601928352015460ff81169186860194600384101561053d5760a097610d4f610d3960ff610d5f96610d6b988b5260081c169760808b0198610d1f81611f1c565b89526040519b8c9b828d5251918c015260c08b0190611b96565b925192601f1993848b83030160408c0152611b96565b9251918884030190880152611b96565b92516080850190611f0f565b51610d7581611f1c565b60a08301520390f35b34610126575f60031936011261012657610dad610d99611caf565b604051918291602083526020830190611b96565b0390f35b3461012657610dc9602063ffffffff6108e036611e4e565b8101600281520301902091165f5260205273ffffffffffffffffffffffffffffffffffffffff60405f2091165f52602052602060ff60405f2054166040519015158152f35b346101265760406003193601126101265760043567ffffffffffffffff811161012657610e3f903690600401611c40565b610e606020610e4c611b21565b928160405193828580945193849201611b75565b8101600581520301902063ffffffff8092165f5260205260405f20610dad610e8c600183549301611d8f565b604051938381869516855260201c166020840152606060408401526060830190611b96565b346101265760606003193601126101265767ffffffffffffffff60043581811161012657610ee3903690600401611c40565b90610eec611b21565b9060443590811161012657610f08610f50913690600401611c40565b60405190818551958187610f236020998a9788809601611b75565b8101600481520301902063ffffffff8095165f52825260405f208260405194838680955193849201611b75565b8201908152030190205416604051908152f35b3461012657608060031936011261012657610f7c611b0e565b60243567ffffffffffffffff811161012657610f9c903690600401611b47565b9160443567ffffffffffffffff811161012657610fbd903690600401611b47565b9060643567ffffffffffffffff811161012657610fde903690600401611b47565b929094610fef6101da368984611bfa565b60405187828237602081898101600181520301902063ffffffff86165f5260205260405f20335f5260205260405f20600381015460ff8116600381101561053d576001146114be576110f161104660018401611d8f565b888b6110a66024888a8c81866040519788957fffffffff00000000000000000000000000000000000000000000000000000000602088019b60e01b168b52878701378401918583015f815237015f83820152036004810184520182611bd7565b5190207f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f52601c526110e8603c5f206110e28c8b3691611bfa565b906120d9565b90939193612113565b73ffffffffffffffffffffffffffffffffffffffff81602082935191012016911603611460576103007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff600392161791015560405187828237602081898101600481520301902063ffffffff86165f5260205261117260405f208385611f8b565b805463ffffffff811663ffffffff81146114335763ffffffff60017fffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000920116911617905560405187828237602081898101600481520301902063ffffffff86165f5260205263ffffffff6111ea60405f208486611f8b565b5416604051888382376020818a8101600581520301902063ffffffff87165f5260205263ffffffff60405f205460201c16111561128f575b9261127f7f545dfb21205c3e096da434dbc9121fc3a052df066ae7bcdcd8e1d86e2ed5bdc195936105389361127163ffffffff978b604051998a99168952608060208a01526080890191611fa4565b918683036040880152611fa4565b9083820360608501523396611fa4565b60405187828237602081898101600581520301902063ffffffff86165f52602052600160405f20019667ffffffffffffffff8311610602576112db836112d58a54611c5e565b8a612022565b5f97601f841160011461137c579261127161127f93610538969363ffffffff99989661135d85807f545dfb21205c3e096da434dbc9121fc3a052df066ae7bcdcd8e1d86e2ed5bdc19e9f5f91611371575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b9a9950935050929495509250611222565b90508801355f61132c565b601f198416815f5260205f20905f5b81811061141b57509361053896938693611271937f545dfb21205c3e096da434dbc9121fc3a052df066ae7bcdcd8e1d86e2ed5bdc19c9d63ffffffff9c9b9961127f99106113e3575b5050600185811b019055611360565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19908801351690558d806113d4565b878c0135835560209b8c019b6001909301920161138b565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f496e76616c69642066696e616c697a6174696f6e207369676e617475726500006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4e6f64652077617320696e76616c6964617465640000000000000000000000006044820152fd5b346101265760606003193601126101265760043573ffffffffffffffffffffffffffffffffffffffff81169081810361012657611557611b21565b9160443567ffffffffffffffff811161012657611578903690600401611b47565b9093611585368387611bfa565b936115a58551602080970120611599611caf565b87815191012014611f26565b60405183878237858185810160018152030190209163ffffffff811692835f52865260405f20855f52865260405f20916115df8354611c5e565b1561173b57600383019283549260ff8416600381101561053d576116dd57937fdb5976211fb963579ffe7833a334dc83e95788b3c5612f860a8b31de484aedc799959361165961169b9484610b2f9a9861163d600260ff9801611d8f565b92611653600161164c85611d8f565b9401611d8f565b93612071565b156116b1575060027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008254161781555b54166040519788978852870190611f0f565b6040850152608060608501526080840191611fa4565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001178155611689565b606489604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601760248201527f4e6f646520616c7265616479206368616c6c656e6765640000000000000000006044820152fd5b606487604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601360248201527f4e6f646520646f6573206e6f74206578697374000000000000000000000000006044820152fd5b34610126576060600319360112610126576117b2611b0e565b67ffffffffffffffff9060248035838111610126576117d5903690600401611b47565b604435918583116101265736602384011215610126578260040135958611610126576024830192602436918860051b0101116101265793611817368387611bfa565b91611837835160208095012061182b611caf565b85815191012014611f26565b63ffffffff5f9216915b87811061184a57005b6040518288823782810190858160019384815203019020845f5285528160405f208a61189561189073ffffffffffffffffffffffffffffffffffffffff9485938c611fc4565b612001565b165f528652600360ff8160405f2001541690811015611929578291600194938c92036118c5575b50505001611841565b604051868c823788818881016002815203019020875f5288526118f06118908560405f20948c611fc4565b165f52865260405f20907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008254161790558989816118bc565b887f4e487b71000000000000000000000000000000000000000000000000000000005f5260216004525ffd5b346101265760406003193601126101265760043567ffffffffffffffff81116101265761198b63ffffffff913690600401611b47565b6020611998939293611b21565b938260405193849283378101600581520301902091165f52602052610dad610d99600160405f2001611d8f565b346101265760a0600319360112610126576119de611b0e565b6119e6611b21565b6119ee611b34565b67ffffffffffffffff919060643583811161012657611a11903690600401611b47565b6084359485116101265761052861053892611a517fde6a4d0e1f964755b0ea6e4653ff1b803cc161f611762ebebbcbc0def6c004f5973690600401611b47565b939098611a5f368385611bfa565b96611a7f88516020809a0120611a73611caf565b8a815191012014611f26565b60405183858237888185810160018152030190209763ffffffff80931698895f52815260405f20335f528152600360405f20611abd60018201611d8f565b50016102007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff825416179055826040519a8b9a8b52169089015216604087015260a0606087015260a0860191611fa4565b6004359063ffffffff8216820361012657565b6024359063ffffffff8216820361012657565b6044359063ffffffff8216820361012657565b9181601f840112156101265782359167ffffffffffffffff8311610126576020838186019501011161012657565b5f5b838110611b865750505f910152565b8181015183820152602001611b77565b90601f19601f602093611bb481518092818752878088019101611b75565b0116010190565b60a0810190811067ffffffffffffffff82111761060257604052565b90601f601f19910116810190811067ffffffffffffffff82111761060257604052565b92919267ffffffffffffffff82116106025760405191611c246020601f19601f8401160184611bd7565b829481845281830111610126578281602093845f960137010152565b9080601f8301121561012657816020611c5b93359101611bfa565b90565b90600182811c92168015611ca5575b6020831014611c7857565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f1691611c6d565b604051905f825f5491611cc183611c5e565b80835292602090600190818116908115611d4c5750600114611cee575b5050611cec92500383611bd7565b565b9150925f80527f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563935f925b828410611d345750611cec9450505081016020015f80611cde565b85548885018301529485019487945092810192611d19565b905060209350611cec9592507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201015f80611cde565b9060405191825f8254611da181611c5e565b908184526020946001916001811690815f14611e0d5750600114611dcf575b505050611cec92500383611bd7565b5f90815285812095935091905b818310611df5575050611cec93508201015f8080611dc0565b85548884018501529485019487945091830191611ddc565b915050611cec9593507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201015f8080611dc0565b6060600319820112610126576004359067ffffffffffffffff821161012657611e7991600401611c40565b9060243563ffffffff81168103610126579060443573ffffffffffffffffffffffffffffffffffffffff811681036101265790565b6060600319820112610126576004359067ffffffffffffffff821161012657611ed991600401611b47565b909160243563ffffffff81168103610126579060443573ffffffffffffffffffffffffffffffffffffffff811681036101265790565b90600382101561053d5752565b6004111561053d57565b15611f2d57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f496e76616c6964206d72656e636c6176650000000000000000000000000000006044820152fd5b6020919283604051948593843782019081520301902090565b601f8260209493601f1993818652868601375f8582860101520116010190565b9190811015611fd45760051b0190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b3573ffffffffffffffffffffffffffffffffffffffff811681036101265790565b601f821161202f57505050565b5f5260205f20906020601f840160051c83019310612067575b601f0160051c01905b81811061205c575050565b5f8155600101612051565b9091508190612048565b51151593929190846120b9575b50836120a9575b508261209e575b5081612096575090565b905051151590565b51151591505f61208c565b63ffffffff16151592505f612085565b73ffffffffffffffffffffffffffffffffffffffff16151593505f61207e565b8151919060418303612109576121029250602082015190606060408401519301515f1a906121e7565b9192909190565b50505f9160029190565b61211c81611f1c565b80612125575050565b61212e81611f1c565b600181036121605760046040517ff645eedf000000000000000000000000000000000000000000000000000000008152fd5b61216981611f1c565b600281036121a257602482604051907ffce698f70000000000000000000000000000000000000000000000000000000082526004820152fd5b806121ae600392611f1c565b146121b65750565b602490604051907fd78bce0c0000000000000000000000000000000000000000000000000000000082526004820152fd5b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411612276579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa1561226b575f5173ffffffffffffffffffffffffffffffffffffffff81161561226157905f905f90565b505f906001905f90565b6040513d5f823e3d90fd5b5050505f916003919056fea26469706673582212209fb238612241c7ddf97fb823e3bede7c321fe02b34de7542dca77f3fc91f035264736f6c63430008170033",
}

// DKGABI is the input ABI used to generate the binding from.
// Deprecated: Use DKGMetaData.ABI instead.
var DKGABI = DKGMetaData.ABI

// DKGBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DKGMetaData.Bin instead.
var DKGBin = DKGMetaData.Bin

// DeployDKG deploys a new Ethereum contract, binding an instance of DKG to it.
func DeployDKG(auth *bind.TransactOpts, backend bind.ContractBackend, mrenclave []byte) (common.Address, *types.Transaction, *DKG, error) {
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
// Solidity: function curMrenclave() view returns(bytes)
func (_DKG *DKGCaller) CurMrenclave(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "curMrenclave")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// CurMrenclave is a free data retrieval call binding the contract method 0x681a0fa8.
//
// Solidity: function curMrenclave() view returns(bytes)
func (_DKG *DKGSession) CurMrenclave() ([]byte, error) {
	return _DKG.Contract.CurMrenclave(&_DKG.CallOpts)
}

// CurMrenclave is a free data retrieval call binding the contract method 0x681a0fa8.
//
// Solidity: function curMrenclave() view returns(bytes)
func (_DKG *DKGCallerSession) CurMrenclave() ([]byte, error) {
	return _DKG.Contract.CurMrenclave(&_DKG.CallOpts)
}

// DealComplaints is a free data retrieval call binding the contract method 0xc2009b92.
//
// Solidity: function dealComplaints(bytes mrenclave, uint32 round, uint32 index, address complainant) view returns(bool)
func (_DKG *DKGCaller) DealComplaints(opts *bind.CallOpts, mrenclave []byte, round uint32, index uint32, complainant common.Address) (bool, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "dealComplaints", mrenclave, round, index, complainant)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DealComplaints is a free data retrieval call binding the contract method 0xc2009b92.
//
// Solidity: function dealComplaints(bytes mrenclave, uint32 round, uint32 index, address complainant) view returns(bool)
func (_DKG *DKGSession) DealComplaints(mrenclave []byte, round uint32, index uint32, complainant common.Address) (bool, error) {
	return _DKG.Contract.DealComplaints(&_DKG.CallOpts, mrenclave, round, index, complainant)
}

// DealComplaints is a free data retrieval call binding the contract method 0xc2009b92.
//
// Solidity: function dealComplaints(bytes mrenclave, uint32 round, uint32 index, address complainant) view returns(bool)
func (_DKG *DKGCallerSession) DealComplaints(mrenclave []byte, round uint32, index uint32, complainant common.Address) (bool, error) {
	return _DKG.Contract.DealComplaints(&_DKG.CallOpts, mrenclave, round, index, complainant)
}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xbf45f2c3.
//
// Solidity: function dkgNodeInfos(bytes mrenclave, uint32 round, address validator) view returns(bytes dkgPubKey, bytes commPubKey, bytes rawQuote, uint8 chalStatus, uint8 nodeStatus)
func (_DKG *DKGCaller) DkgNodeInfos(opts *bind.CallOpts, mrenclave []byte, round uint32, validator common.Address) (struct {
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

// DkgNodeInfos is a free data retrieval call binding the contract method 0xbf45f2c3.
//
// Solidity: function dkgNodeInfos(bytes mrenclave, uint32 round, address validator) view returns(bytes dkgPubKey, bytes commPubKey, bytes rawQuote, uint8 chalStatus, uint8 nodeStatus)
func (_DKG *DKGSession) DkgNodeInfos(mrenclave []byte, round uint32, validator common.Address) (struct {
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	NodeStatus uint8
}, error) {
	return _DKG.Contract.DkgNodeInfos(&_DKG.CallOpts, mrenclave, round, validator)
}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xbf45f2c3.
//
// Solidity: function dkgNodeInfos(bytes mrenclave, uint32 round, address validator) view returns(bytes dkgPubKey, bytes commPubKey, bytes rawQuote, uint8 chalStatus, uint8 nodeStatus)
func (_DKG *DKGCallerSession) DkgNodeInfos(mrenclave []byte, round uint32, validator common.Address) (struct {
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	NodeStatus uint8
}, error) {
	return _DKG.Contract.DkgNodeInfos(&_DKG.CallOpts, mrenclave, round, validator)
}

// GetGlobalPubKey is a free data retrieval call binding the contract method 0x285f0e63.
//
// Solidity: function getGlobalPubKey(bytes mrenclave, uint32 round) view returns(bytes)
func (_DKG *DKGCaller) GetGlobalPubKey(opts *bind.CallOpts, mrenclave []byte, round uint32) ([]byte, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "getGlobalPubKey", mrenclave, round)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetGlobalPubKey is a free data retrieval call binding the contract method 0x285f0e63.
//
// Solidity: function getGlobalPubKey(bytes mrenclave, uint32 round) view returns(bytes)
func (_DKG *DKGSession) GetGlobalPubKey(mrenclave []byte, round uint32) ([]byte, error) {
	return _DKG.Contract.GetGlobalPubKey(&_DKG.CallOpts, mrenclave, round)
}

// GetGlobalPubKey is a free data retrieval call binding the contract method 0x285f0e63.
//
// Solidity: function getGlobalPubKey(bytes mrenclave, uint32 round) view returns(bytes)
func (_DKG *DKGCallerSession) GetGlobalPubKey(mrenclave []byte, round uint32) ([]byte, error) {
	return _DKG.Contract.GetGlobalPubKey(&_DKG.CallOpts, mrenclave, round)
}

// GetNodeInfo is a free data retrieval call binding the contract method 0x6ddc4a25.
//
// Solidity: function getNodeInfo(bytes mrenclave, uint32 round, address validator) view returns((bytes,bytes,bytes,uint8,uint8))
func (_DKG *DKGCaller) GetNodeInfo(opts *bind.CallOpts, mrenclave []byte, round uint32, validator common.Address) (IDKGNodeInfo, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "getNodeInfo", mrenclave, round, validator)

	if err != nil {
		return *new(IDKGNodeInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IDKGNodeInfo)).(*IDKGNodeInfo)

	return out0, err

}

// GetNodeInfo is a free data retrieval call binding the contract method 0x6ddc4a25.
//
// Solidity: function getNodeInfo(bytes mrenclave, uint32 round, address validator) view returns((bytes,bytes,bytes,uint8,uint8))
func (_DKG *DKGSession) GetNodeInfo(mrenclave []byte, round uint32, validator common.Address) (IDKGNodeInfo, error) {
	return _DKG.Contract.GetNodeInfo(&_DKG.CallOpts, mrenclave, round, validator)
}

// GetNodeInfo is a free data retrieval call binding the contract method 0x6ddc4a25.
//
// Solidity: function getNodeInfo(bytes mrenclave, uint32 round, address validator) view returns((bytes,bytes,bytes,uint8,uint8))
func (_DKG *DKGCallerSession) GetNodeInfo(mrenclave []byte, round uint32, validator common.Address) (IDKGNodeInfo, error) {
	return _DKG.Contract.GetNodeInfo(&_DKG.CallOpts, mrenclave, round, validator)
}

// IsActiveValidator is a free data retrieval call binding the contract method 0xfb5a783a.
//
// Solidity: function isActiveValidator(bytes mrenclave, uint32 round, address validator) view returns(bool)
func (_DKG *DKGCaller) IsActiveValidator(opts *bind.CallOpts, mrenclave []byte, round uint32, validator common.Address) (bool, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "isActiveValidator", mrenclave, round, validator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsActiveValidator is a free data retrieval call binding the contract method 0xfb5a783a.
//
// Solidity: function isActiveValidator(bytes mrenclave, uint32 round, address validator) view returns(bool)
func (_DKG *DKGSession) IsActiveValidator(mrenclave []byte, round uint32, validator common.Address) (bool, error) {
	return _DKG.Contract.IsActiveValidator(&_DKG.CallOpts, mrenclave, round, validator)
}

// IsActiveValidator is a free data retrieval call binding the contract method 0xfb5a783a.
//
// Solidity: function isActiveValidator(bytes mrenclave, uint32 round, address validator) view returns(bool)
func (_DKG *DKGCallerSession) IsActiveValidator(mrenclave []byte, round uint32, validator common.Address) (bool, error) {
	return _DKG.Contract.IsActiveValidator(&_DKG.CallOpts, mrenclave, round, validator)
}

// RoundInfo is a free data retrieval call binding the contract method 0x3ce1c1e0.
//
// Solidity: function roundInfo(bytes mrenclave, uint32 round) view returns(uint32 total, uint32 threshold, bytes globalPubKey)
func (_DKG *DKGCaller) RoundInfo(opts *bind.CallOpts, mrenclave []byte, round uint32) (struct {
	Total        uint32
	Threshold    uint32
	GlobalPubKey []byte
}, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "roundInfo", mrenclave, round)

	outstruct := new(struct {
		Total        uint32
		Threshold    uint32
		GlobalPubKey []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Total = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Threshold = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.GlobalPubKey = *abi.ConvertType(out[2], new([]byte)).(*[]byte)

	return *outstruct, err

}

// RoundInfo is a free data retrieval call binding the contract method 0x3ce1c1e0.
//
// Solidity: function roundInfo(bytes mrenclave, uint32 round) view returns(uint32 total, uint32 threshold, bytes globalPubKey)
func (_DKG *DKGSession) RoundInfo(mrenclave []byte, round uint32) (struct {
	Total        uint32
	Threshold    uint32
	GlobalPubKey []byte
}, error) {
	return _DKG.Contract.RoundInfo(&_DKG.CallOpts, mrenclave, round)
}

// RoundInfo is a free data retrieval call binding the contract method 0x3ce1c1e0.
//
// Solidity: function roundInfo(bytes mrenclave, uint32 round) view returns(uint32 total, uint32 threshold, bytes globalPubKey)
func (_DKG *DKGCallerSession) RoundInfo(mrenclave []byte, round uint32) (struct {
	Total        uint32
	Threshold    uint32
	GlobalPubKey []byte
}, error) {
	return _DKG.Contract.RoundInfo(&_DKG.CallOpts, mrenclave, round)
}

// ValSets is a free data retrieval call binding the contract method 0x4f1bf88d.
//
// Solidity: function valSets(bytes mrenclave, uint32 round, address validator) view returns(bool)
func (_DKG *DKGCaller) ValSets(opts *bind.CallOpts, mrenclave []byte, round uint32, validator common.Address) (bool, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "valSets", mrenclave, round, validator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValSets is a free data retrieval call binding the contract method 0x4f1bf88d.
//
// Solidity: function valSets(bytes mrenclave, uint32 round, address validator) view returns(bool)
func (_DKG *DKGSession) ValSets(mrenclave []byte, round uint32, validator common.Address) (bool, error) {
	return _DKG.Contract.ValSets(&_DKG.CallOpts, mrenclave, round, validator)
}

// ValSets is a free data retrieval call binding the contract method 0x4f1bf88d.
//
// Solidity: function valSets(bytes mrenclave, uint32 round, address validator) view returns(bool)
func (_DKG *DKGCallerSession) ValSets(mrenclave []byte, round uint32, validator common.Address) (bool, error) {
	return _DKG.Contract.ValSets(&_DKG.CallOpts, mrenclave, round, validator)
}

// Votes is a free data retrieval call binding the contract method 0x3c3acdb4.
//
// Solidity: function votes(bytes mrenclave, uint32 round, bytes globalPubKeyCandidates) view returns(uint32 votes)
func (_DKG *DKGCaller) Votes(opts *bind.CallOpts, mrenclave []byte, round uint32, globalPubKeyCandidates []byte) (uint32, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "votes", mrenclave, round, globalPubKeyCandidates)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// Votes is a free data retrieval call binding the contract method 0x3c3acdb4.
//
// Solidity: function votes(bytes mrenclave, uint32 round, bytes globalPubKeyCandidates) view returns(uint32 votes)
func (_DKG *DKGSession) Votes(mrenclave []byte, round uint32, globalPubKeyCandidates []byte) (uint32, error) {
	return _DKG.Contract.Votes(&_DKG.CallOpts, mrenclave, round, globalPubKeyCandidates)
}

// Votes is a free data retrieval call binding the contract method 0x3c3acdb4.
//
// Solidity: function votes(bytes mrenclave, uint32 round, bytes globalPubKeyCandidates) view returns(uint32 votes)
func (_DKG *DKGCallerSession) Votes(mrenclave []byte, round uint32, globalPubKeyCandidates []byte) (uint32, error) {
	return _DKG.Contract.Votes(&_DKG.CallOpts, mrenclave, round, globalPubKeyCandidates)
}

// ComplainDeals is a paid mutator transaction binding the contract method 0x7a574765.
//
// Solidity: function complainDeals(uint32 round, uint32 index, uint32[] complainIndexes, bytes mrenclave) returns()
func (_DKG *DKGTransactor) ComplainDeals(opts *bind.TransactOpts, round uint32, index uint32, complainIndexes []uint32, mrenclave []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "complainDeals", round, index, complainIndexes, mrenclave)
}

// ComplainDeals is a paid mutator transaction binding the contract method 0x7a574765.
//
// Solidity: function complainDeals(uint32 round, uint32 index, uint32[] complainIndexes, bytes mrenclave) returns()
func (_DKG *DKGSession) ComplainDeals(round uint32, index uint32, complainIndexes []uint32, mrenclave []byte) (*types.Transaction, error) {
	return _DKG.Contract.ComplainDeals(&_DKG.TransactOpts, round, index, complainIndexes, mrenclave)
}

// ComplainDeals is a paid mutator transaction binding the contract method 0x7a574765.
//
// Solidity: function complainDeals(uint32 round, uint32 index, uint32[] complainIndexes, bytes mrenclave) returns()
func (_DKG *DKGTransactorSession) ComplainDeals(round uint32, index uint32, complainIndexes []uint32, mrenclave []byte) (*types.Transaction, error) {
	return _DKG.Contract.ComplainDeals(&_DKG.TransactOpts, round, index, complainIndexes, mrenclave)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x3090f411.
//
// Solidity: function finalizeDKG(uint32 round, bytes mrenclave, bytes globalPubKey, bytes signature) returns()
func (_DKG *DKGTransactor) FinalizeDKG(opts *bind.TransactOpts, round uint32, mrenclave []byte, globalPubKey []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "finalizeDKG", round, mrenclave, globalPubKey, signature)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x3090f411.
//
// Solidity: function finalizeDKG(uint32 round, bytes mrenclave, bytes globalPubKey, bytes signature) returns()
func (_DKG *DKGSession) FinalizeDKG(round uint32, mrenclave []byte, globalPubKey []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.FinalizeDKG(&_DKG.TransactOpts, round, mrenclave, globalPubKey, signature)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x3090f411.
//
// Solidity: function finalizeDKG(uint32 round, bytes mrenclave, bytes globalPubKey, bytes signature) returns()
func (_DKG *DKGTransactorSession) FinalizeDKG(round uint32, mrenclave []byte, globalPubKey []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.FinalizeDKG(&_DKG.TransactOpts, round, mrenclave, globalPubKey, signature)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0xc3c77ff9.
//
// Solidity: function initializeDKG(uint32 round, bytes mrenclave, bytes dkgPubKey, bytes commPubKey, bytes rawQuote) returns()
func (_DKG *DKGTransactor) InitializeDKG(opts *bind.TransactOpts, round uint32, mrenclave []byte, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "initializeDKG", round, mrenclave, dkgPubKey, commPubKey, rawQuote)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0xc3c77ff9.
//
// Solidity: function initializeDKG(uint32 round, bytes mrenclave, bytes dkgPubKey, bytes commPubKey, bytes rawQuote) returns()
func (_DKG *DKGSession) InitializeDKG(round uint32, mrenclave []byte, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) (*types.Transaction, error) {
	return _DKG.Contract.InitializeDKG(&_DKG.TransactOpts, round, mrenclave, dkgPubKey, commPubKey, rawQuote)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0xc3c77ff9.
//
// Solidity: function initializeDKG(uint32 round, bytes mrenclave, bytes dkgPubKey, bytes commPubKey, bytes rawQuote) returns()
func (_DKG *DKGTransactorSession) InitializeDKG(round uint32, mrenclave []byte, dkgPubKey []byte, commPubKey []byte, rawQuote []byte) (*types.Transaction, error) {
	return _DKG.Contract.InitializeDKG(&_DKG.TransactOpts, round, mrenclave, dkgPubKey, commPubKey, rawQuote)
}

// RequestRemoteAttestationOnChain is a paid mutator transaction binding the contract method 0x2ef619b2.
//
// Solidity: function requestRemoteAttestationOnChain(address targetValidatorAddr, uint32 round, bytes mrenclave) returns()
func (_DKG *DKGTransactor) RequestRemoteAttestationOnChain(opts *bind.TransactOpts, targetValidatorAddr common.Address, round uint32, mrenclave []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "requestRemoteAttestationOnChain", targetValidatorAddr, round, mrenclave)
}

// RequestRemoteAttestationOnChain is a paid mutator transaction binding the contract method 0x2ef619b2.
//
// Solidity: function requestRemoteAttestationOnChain(address targetValidatorAddr, uint32 round, bytes mrenclave) returns()
func (_DKG *DKGSession) RequestRemoteAttestationOnChain(targetValidatorAddr common.Address, round uint32, mrenclave []byte) (*types.Transaction, error) {
	return _DKG.Contract.RequestRemoteAttestationOnChain(&_DKG.TransactOpts, targetValidatorAddr, round, mrenclave)
}

// RequestRemoteAttestationOnChain is a paid mutator transaction binding the contract method 0x2ef619b2.
//
// Solidity: function requestRemoteAttestationOnChain(address targetValidatorAddr, uint32 round, bytes mrenclave) returns()
func (_DKG *DKGTransactorSession) RequestRemoteAttestationOnChain(targetValidatorAddr common.Address, round uint32, mrenclave []byte) (*types.Transaction, error) {
	return _DKG.Contract.RequestRemoteAttestationOnChain(&_DKG.TransactOpts, targetValidatorAddr, round, mrenclave)
}

// SetNetwork is a paid mutator transaction binding the contract method 0x02372b2e.
//
// Solidity: function setNetwork(uint32 round, uint32 total, uint32 threshold, bytes mrenclave, bytes signature) returns()
func (_DKG *DKGTransactor) SetNetwork(opts *bind.TransactOpts, round uint32, total uint32, threshold uint32, mrenclave []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "setNetwork", round, total, threshold, mrenclave, signature)
}

// SetNetwork is a paid mutator transaction binding the contract method 0x02372b2e.
//
// Solidity: function setNetwork(uint32 round, uint32 total, uint32 threshold, bytes mrenclave, bytes signature) returns()
func (_DKG *DKGSession) SetNetwork(round uint32, total uint32, threshold uint32, mrenclave []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.SetNetwork(&_DKG.TransactOpts, round, total, threshold, mrenclave, signature)
}

// SetNetwork is a paid mutator transaction binding the contract method 0x02372b2e.
//
// Solidity: function setNetwork(uint32 round, uint32 total, uint32 threshold, bytes mrenclave, bytes signature) returns()
func (_DKG *DKGTransactorSession) SetNetwork(round uint32, total uint32, threshold uint32, mrenclave []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.SetNetwork(&_DKG.TransactOpts, round, total, threshold, mrenclave, signature)
}

// SubmitActiveValSet is a paid mutator transaction binding the contract method 0x2e266b4c.
//
// Solidity: function submitActiveValSet(uint32 round, bytes mrenclave, address[] valSet) returns()
func (_DKG *DKGTransactor) SubmitActiveValSet(opts *bind.TransactOpts, round uint32, mrenclave []byte, valSet []common.Address) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "submitActiveValSet", round, mrenclave, valSet)
}

// SubmitActiveValSet is a paid mutator transaction binding the contract method 0x2e266b4c.
//
// Solidity: function submitActiveValSet(uint32 round, bytes mrenclave, address[] valSet) returns()
func (_DKG *DKGSession) SubmitActiveValSet(round uint32, mrenclave []byte, valSet []common.Address) (*types.Transaction, error) {
	return _DKG.Contract.SubmitActiveValSet(&_DKG.TransactOpts, round, mrenclave, valSet)
}

// SubmitActiveValSet is a paid mutator transaction binding the contract method 0x2e266b4c.
//
// Solidity: function submitActiveValSet(uint32 round, bytes mrenclave, address[] valSet) returns()
func (_DKG *DKGTransactorSession) SubmitActiveValSet(round uint32, mrenclave []byte, valSet []common.Address) (*types.Transaction, error) {
	return _DKG.Contract.SubmitActiveValSet(&_DKG.TransactOpts, round, mrenclave, valSet)
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
	Mrenclave    []byte
	GlobalPubKey []byte
	Signature    []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDKGFinalized is a free log retrieval operation binding the contract event 0x545dfb21205c3e096da434dbc9121fc3a052df066ae7bcdcd8e1d86e2ed5bdc1.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, bytes mrenclave, bytes globalPubKey, bytes signature)
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

// WatchDKGFinalized is a free log subscription operation binding the contract event 0x545dfb21205c3e096da434dbc9121fc3a052df066ae7bcdcd8e1d86e2ed5bdc1.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, bytes mrenclave, bytes globalPubKey, bytes signature)
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

// ParseDKGFinalized is a log parse operation binding the contract event 0x545dfb21205c3e096da434dbc9121fc3a052df066ae7bcdcd8e1d86e2ed5bdc1.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, bytes mrenclave, bytes globalPubKey, bytes signature)
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
	Mrenclave  []byte
	Round      uint32
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDKGInitialized is a free log retrieval operation binding the contract event 0xe8461cf9f1302c7b5cf4074f825c58d73aa03cca3c607cfcf9d90a137e5f2f0f.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes mrenclave, uint32 round, bytes dkgPubKey, bytes commPubKey, bytes rawQuote)
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

// WatchDKGInitialized is a free log subscription operation binding the contract event 0xe8461cf9f1302c7b5cf4074f825c58d73aa03cca3c607cfcf9d90a137e5f2f0f.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes mrenclave, uint32 round, bytes dkgPubKey, bytes commPubKey, bytes rawQuote)
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

// ParseDKGInitialized is a log parse operation binding the contract event 0xe8461cf9f1302c7b5cf4074f825c58d73aa03cca3c607cfcf9d90a137e5f2f0f.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes mrenclave, uint32 round, bytes dkgPubKey, bytes commPubKey, bytes rawQuote)
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
	Mrenclave []byte
	Signature []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDKGNetworkSet is a free log retrieval operation binding the contract event 0xde6a4d0e1f964755b0ea6e4653ff1b803cc161f611762ebebbcbc0def6c004f5.
//
// Solidity: event DKGNetworkSet(address indexed msgSender, uint32 round, uint32 total, uint32 threshold, bytes mrenclave, bytes signature)
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

// WatchDKGNetworkSet is a free log subscription operation binding the contract event 0xde6a4d0e1f964755b0ea6e4653ff1b803cc161f611762ebebbcbc0def6c004f5.
//
// Solidity: event DKGNetworkSet(address indexed msgSender, uint32 round, uint32 total, uint32 threshold, bytes mrenclave, bytes signature)
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

// ParseDKGNetworkSet is a log parse operation binding the contract event 0xde6a4d0e1f964755b0ea6e4653ff1b803cc161f611762ebebbcbc0def6c004f5.
//
// Solidity: event DKGNetworkSet(address indexed msgSender, uint32 round, uint32 total, uint32 threshold, bytes mrenclave, bytes signature)
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
	Mrenclave       []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterDealComplaintsSubmitted is a free log retrieval operation binding the contract event 0x1c2112af5fd37661e3dd248d701decebf291f13de6a411176337ef21a7b1a630.
//
// Solidity: event DealComplaintsSubmitted(uint32 index, uint32[] complainIndexes, uint32 round, bytes mrenclave)
func (_DKG *DKGFilterer) FilterDealComplaintsSubmitted(opts *bind.FilterOpts) (*DKGDealComplaintsSubmittedIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "DealComplaintsSubmitted")
	if err != nil {
		return nil, err
	}
	return &DKGDealComplaintsSubmittedIterator{contract: _DKG.contract, event: "DealComplaintsSubmitted", logs: logs, sub: sub}, nil
}

// WatchDealComplaintsSubmitted is a free log subscription operation binding the contract event 0x1c2112af5fd37661e3dd248d701decebf291f13de6a411176337ef21a7b1a630.
//
// Solidity: event DealComplaintsSubmitted(uint32 index, uint32[] complainIndexes, uint32 round, bytes mrenclave)
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

// ParseDealComplaintsSubmitted is a log parse operation binding the contract event 0x1c2112af5fd37661e3dd248d701decebf291f13de6a411176337ef21a7b1a630.
//
// Solidity: event DealComplaintsSubmitted(uint32 index, uint32[] complainIndexes, uint32 round, bytes mrenclave)
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
	Mrenclave      []byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDealVerified is a free log retrieval operation binding the contract event 0x47a7318130d57c2f660d7db8746be21176056c1548727535255ad2005a709277.
//
// Solidity: event DealVerified(uint32 index, uint32 recipientIndex, uint32 round, bytes mrenclave)
func (_DKG *DKGFilterer) FilterDealVerified(opts *bind.FilterOpts) (*DKGDealVerifiedIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "DealVerified")
	if err != nil {
		return nil, err
	}
	return &DKGDealVerifiedIterator{contract: _DKG.contract, event: "DealVerified", logs: logs, sub: sub}, nil
}

// WatchDealVerified is a free log subscription operation binding the contract event 0x47a7318130d57c2f660d7db8746be21176056c1548727535255ad2005a709277.
//
// Solidity: event DealVerified(uint32 index, uint32 recipientIndex, uint32 round, bytes mrenclave)
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

// ParseDealVerified is a log parse operation binding the contract event 0x47a7318130d57c2f660d7db8746be21176056c1548727535255ad2005a709277.
//
// Solidity: event DealVerified(uint32 index, uint32 recipientIndex, uint32 round, bytes mrenclave)
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
	Mrenclave []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterInvalidDeal is a free log retrieval operation binding the contract event 0x51c062375794d1d1717c2b487871bb5a145bb01911584ff161aa0a01e06b2b9d.
//
// Solidity: event InvalidDeal(uint32 index, uint32 round, bytes mrenclave)
func (_DKG *DKGFilterer) FilterInvalidDeal(opts *bind.FilterOpts) (*DKGInvalidDealIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "InvalidDeal")
	if err != nil {
		return nil, err
	}
	return &DKGInvalidDealIterator{contract: _DKG.contract, event: "InvalidDeal", logs: logs, sub: sub}, nil
}

// WatchInvalidDeal is a free log subscription operation binding the contract event 0x51c062375794d1d1717c2b487871bb5a145bb01911584ff161aa0a01e06b2b9d.
//
// Solidity: event InvalidDeal(uint32 index, uint32 round, bytes mrenclave)
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

// ParseInvalidDeal is a log parse operation binding the contract event 0x51c062375794d1d1717c2b487871bb5a145bb01911584ff161aa0a01e06b2b9d.
//
// Solidity: event InvalidDeal(uint32 index, uint32 round, bytes mrenclave)
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
	Mrenclave  []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRemoteAttestationProcessedOnChain is a free log retrieval operation binding the contract event 0xdb5976211fb963579ffe7833a334dc83e95788b3c5612f860a8b31de484aedc7.
//
// Solidity: event RemoteAttestationProcessedOnChain(address validator, uint8 chalStatus, uint32 round, bytes mrenclave)
func (_DKG *DKGFilterer) FilterRemoteAttestationProcessedOnChain(opts *bind.FilterOpts) (*DKGRemoteAttestationProcessedOnChainIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "RemoteAttestationProcessedOnChain")
	if err != nil {
		return nil, err
	}
	return &DKGRemoteAttestationProcessedOnChainIterator{contract: _DKG.contract, event: "RemoteAttestationProcessedOnChain", logs: logs, sub: sub}, nil
}

// WatchRemoteAttestationProcessedOnChain is a free log subscription operation binding the contract event 0xdb5976211fb963579ffe7833a334dc83e95788b3c5612f860a8b31de484aedc7.
//
// Solidity: event RemoteAttestationProcessedOnChain(address validator, uint8 chalStatus, uint32 round, bytes mrenclave)
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

// ParseRemoteAttestationProcessedOnChain is a log parse operation binding the contract event 0xdb5976211fb963579ffe7833a334dc83e95788b3c5612f860a8b31de484aedc7.
//
// Solidity: event RemoteAttestationProcessedOnChain(address validator, uint8 chalStatus, uint32 round, bytes mrenclave)
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
	Mrenclave        []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpgradeScheduled is a free log retrieval operation binding the contract event 0x1134b804c6fb57d1b1e7840459379c705537ad30ec193c6deb6cc7ee7709ed70.
//
// Solidity: event UpgradeScheduled(uint32 activationHeight, bytes mrenclave)
func (_DKG *DKGFilterer) FilterUpgradeScheduled(opts *bind.FilterOpts) (*DKGUpgradeScheduledIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "UpgradeScheduled")
	if err != nil {
		return nil, err
	}
	return &DKGUpgradeScheduledIterator{contract: _DKG.contract, event: "UpgradeScheduled", logs: logs, sub: sub}, nil
}

// WatchUpgradeScheduled is a free log subscription operation binding the contract event 0x1134b804c6fb57d1b1e7840459379c705537ad30ec193c6deb6cc7ee7709ed70.
//
// Solidity: event UpgradeScheduled(uint32 activationHeight, bytes mrenclave)
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

// ParseUpgradeScheduled is a log parse operation binding the contract event 0x1134b804c6fb57d1b1e7840459379c705537ad30ec193c6deb6cc7ee7709ed70.
//
// Solidity: event UpgradeScheduled(uint32 activationHeight, bytes mrenclave)
func (_DKG *DKGFilterer) ParseUpgradeScheduled(log types.Log) (*DKGUpgradeScheduled, error) {
	event := new(DKGUpgradeScheduled)
	if err := _DKG.contract.UnpackLog(event, "UpgradeScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
