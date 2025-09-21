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
	Index      uint32
	Validator  common.Address
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	Finalized  bool
}

// DKGMetaData contains all meta data concerning the DKG contract.
var DKGMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"complainDeals\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dealComplaints\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainant\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dkgNodeInfos\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"finalized\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"finalized\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getGlobalPubKey\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeCount\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeInfo\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDKG.NodeInfo\",\"components\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"finalized\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isActiveValidator\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nodeCount\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"nodeCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestRemoteAttestationOnChain\",\"inputs\":[{\"name\":\"targetIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"roundInfo\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"submitActiveValSet\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"valSet\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"valSets\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"votes\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"globalPubKeyCandidates\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"votes\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DKGFinalized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"finalized\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"globalPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGInitialized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"rawQuote\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealComplaintsSubmitted\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"indexed\":false,\"internalType\":\"uint32[]\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealVerified\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"recipientIndex\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidDKGInitialization\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidDeal\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RegistrationChallenged\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"challenger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoteAttestationProcessedOnChain\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeScheduled\",\"inputs\":[{\"name\":\"activationHeight\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x6080806040523461001657612210908161001b8239f35b5f80fdfe6080806040526004361015610012575f80fd5b5f3560e01c90816308e48422146119cb575080630c5a402014611814578063159726be146117dc578063285f0e63146117845780632e266b4c146116055780633c3acdb4146115535780633ce1c1e0146114d65780634f1bf88d1461143257806378e510d2146113f15780637a5747651461114e578063a27c2f6c14611025578063c2009b9214610f6f578063c3c77ff914610810578063fb5a783a146107785763fe9dee68146100c1575f80fd5b346107745760c0600319360112610774576100da611c1e565b6100e2611c31565b906044351515604435036107745760643567ffffffffffffffff811161077457610110903690600401611c57565b919060843567ffffffffffffffff811161077457610132903690600401611c57565b60a43567ffffffffffffffff811161077457610152903690600401611c57565b939095604051818382376020818381015f81520301902063ffffffff87165f5260205260405f2063ffffffff89165f5260205260405f2073ffffffffffffffffffffffffffffffffffffffff815460201c1633036107165760ff60048201541660038110156106e95760011461068b576102ad73ffffffffffffffffffffffffffffffffffffffff80878c8761026160298e8c6101f160028c01611e4b565b968c60405196879460208601997fffffffff00000000000000000000000000000000000000000000000000000000809260e01b168b5260e01b166024860152604435151560f81b602886015285850137818d8401918583015f815237015f83820152036009810184520182611cf4565b5190207f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f52601c526102a48c61029e603c5f20918d3691611d17565b9061204e565b90959195612088565b602081519101201691160361062d5760040180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16604435151560081b61ff00161790556044356102fb57005b60405181838237602081838101600481520301902063ffffffff87165f5260205261032a60405f208486611f2f565b80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000063ffffffff61035d818416611f68565b16911617905560405181838237602081838101600481520301902063ffffffff87165f5260205263ffffffff61039760405f208587611f2f565b541660405182848237602081848101600581520301902063ffffffff88165f5260205263ffffffff60405f205460201c161115610454575b9263ffffffff95949261043161043f9361044f96897f694566864b727668476d766533aac9ac9c8f2a98c4c58322f1cc6b01258a5c489b9c6040519b8c9b168b521660208a0152604435151560408a015260c060608a015260c0890191611f48565b918683036080880152611f48565b9083820360a08501523396611f48565b0390a2005b60405181838237602081838101600581520301902063ffffffff87165f52602052600160405f20019767ffffffffffffffff8411610600576104a08461049a8b54611dfa565b8b611faa565b5f601f85116001146105445761043f9361044f969363ffffffff87947f694566864b727668476d766533aac9ac9c8f2a98c4c58322f1cc6b01258a5c489c9d61052387849e9d9b610431985f91610539575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b9c9b5050939650935050929495506103cf565b90508a01355f6104f2565b98805f5260205f205f5b601f19871681106105e857509361044f969363ffffffff87946104319461043f987f694566864b727668476d766533aac9ac9c8f2a98c4c58322f1cc6b01258a5c489e9f88859f9e9c601f1916106105b0575b5050600187811b019055610526565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88a60031b161c19908a01351690555f806105a1565b878c0135825560209b8c019b6001909201910161054e565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f496e76616c69642066696e616c697a6174696f6e207369676e617475726500006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4e6f64652077617320696e76616c6964617465640000000000000000000000006044820152fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f496e76616c69642073656e6465720000000000000000000000000000000000006044820152fd5b5f80fd5b346107745760606003193601126107745760043567ffffffffffffffff8111610774576107a9903690600401611c57565b63ffffffff6107b6611c31565b9160206107c1611f0c565b948260405193849283378101600181520301902091165f5260205273ffffffffffffffffffffffffffffffffffffffff60405f2091165f52602052602060ff60405f2054166040519015158152f35b346107745760a060031936011261077457610829611c1e565b60243567ffffffffffffffff811161077457610849903690600401611c57565b60449291923567ffffffffffffffff81116107745761086c903690600401611c57565b9060643567ffffffffffffffff81116107745761088d903690600401611c57565b9560843567ffffffffffffffff8111610774576108ae903690600401611c57565b92909760405187848237602081898101600181520301902063ffffffff89165f5260205260405f20335f5260205260ff60405f2054168015610f67575b15610f09576109126108fe36868c611d17565b8961090a368a8a611d17565b913390611ff9565b15610eab57604051878482376003888201526020818981010301902063ffffffff89165f5260205263ffffffff60405f20541697604051888582376020818a8101600381520301902063ffffffff82165f5260205260405f2080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000063ffffffff61099e818416611f68565b169116179055604051906109b182611cd8565b8982523360208301526109c5368989611d17565b60408301526109d5368486611d17565b60608301526109e536878d611d17565b60808301525f60a08301525f60c0830152604051898682376020818b81015f81520301902063ffffffff82165f5260205260405f208a5f5260205260405f2063ffffffff8351168154907fffffffffffffffff00000000000000000000000000000000000000000000000077ffffffffffffffffffffffffffffffffffffffff00000000602087015160201b16921617178155604083015180519067ffffffffffffffff821161060057610aa982610aa06001860154611dfa565b60018601611faa565b602090601f8311600114610e1e57610af692915f9183610d86575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60018201555b606083015180519067ffffffffffffffff821161060057610b2d82610b246002860154611dfa565b60028601611faa565b602090601f8311600114610d9157610b7992915f9183610d865750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60028201555b608083015180519067ffffffffffffffff821161060057610bb082610ba76003860154611dfa565b60038601611faa565b602090601f8311600114610cf8579180610c0092600495945f92610ced5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60038201555b019760a08301519760038910156106e9579061043192918d9c63ffffffff610cd461044f9b9a7f5c5aa077fdba8d39c55de6079c75de88ec2031742c08609237382f11414de6169f9e8060c09f60c0909b61043f9e9d9c60ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00610cc39654169116178355015115157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff00835492151560081b169116179055565b6040519e8f9e8f8181520191611f48565b941660208c015260408b015289830360608b0152611f48565b015190505f80610ac4565b90600384015f5260205f20915f5b601f1985168110610d6e575091839160019383601f19600498971610610d37575b505050811b016003820155610c06565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558f8080610d27565b91926020600181928685015181550194019201610d06565b015190508f80610ac4565b9190600284015f5260205f20905f935b601f1984168510610e03576001945083601f19811610610dcc575b505050811b016002820155610b7f565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e8080610dbc565b81810151835560209485019460019093019290910190610da1565b9190600184015f5260205f20905f935b601f1984168510610e90576001945083601f19811610610e59575b505050811b016001820155610afc565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e8080610e49565b81810151835560209485019460019093019290910190610e2e565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f496e76616c69642072656d6f7465206174746573746174696f6e0000000000006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f56616c696461746f72206e6f7420696e206163746976652073657400000000006044820152fd5b5060016108eb565b346107745760806003193601126107745760043567ffffffffffffffff811161077457610fa0903690600401611d5d565b610fa8611c31565b610fb0611c44565b6064359173ffffffffffffffffffffffffffffffffffffffff8316809303610774576040518481809651610fea816020998a809601611c85565b8101600281520301902063ffffffff8092165f52845260405f2091165f52825260405f20905f52815260ff60405f2054166040519015158152f35b346107745760606003193601126107745760043567ffffffffffffffff811161077457611056903690600401611d5d565b61105e611c31565b611066611c44565b604051838180955161107e8160209889809601611c85565b81015f81520301902063ffffffff8093165f5283528160405f2091165f52825260405f209160ff61112f8454946111216110ba60018301611e4b565b6111136110c960028501611e4b565b9160046110d860038701611e4b565b9501549773ffffffffffffffffffffffffffffffffffffffff6040519b8c9b81168c52821c16908a015260e060408a015260e0890190611ca6565b908782036060890152611ca6565b908582036080870152611ca6565b9161113f60a08501838316611ccb565b60081c16151560c08301520390f35b3461077457608060031936011261077457611167611c1e565b61116f611c31565b60443567ffffffffffffffff808211610774573660238301121561077457816004013593602491808611610600576005908660051b90602095604051986111b88885018b611cf4565b8952868901602481948301019136831161077457602401905b8282106113d557505050606435908111610774576111f3903690600401611c57565b906040518282823787818481015f8152030190209263ffffffff80961693845f5288528560405f20991698895f52885273ffffffffffffffffffffffffffffffffffffffff60405f2054891c163303611377575f5b8a518110156112f0576040518484823789818681016002815203019020855f52895260405f208b518210156112c45790600191888b8e848b1b010151165f528a5260405f20335f528a5260405f20827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082541617905501611248565b887f4e487b71000000000000000000000000000000000000000000000000000000005f5260326004525ffd5b509350889285898960405196608088019288526080828901525180925260a0870197925f905b83821061135e577f1c2112af5fd37661e3dd248d701decebf291f13de6a411176337ef21a7b1a63089806113598d8c8c8c60408601528483036060860152611f48565b0390a1005b845181168a529882019893820193600190910190611316565b606488604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601360248201527f496e76616c696420636f6d706c61696e616e74000000000000000000000000006044820152fd5b813563ffffffff811681036107745781529088019088016111d1565b346107745760208061140236611dba565b9290918260405193849283378101600381520301902063ffffffff8092165f52825260405f205416604051908152f35b346107745760606003193601126107745760043567ffffffffffffffff811161077457611463903690600401611d5d565b61146b611c31565b63ffffffff611491602061147d611f0c565b948160405193828580945193849201611c85565b8101600181520301902091165f5260205273ffffffffffffffffffffffffffffffffffffffff60405f2091165f52602052602060ff60405f2054166040519015158152f35b346107745760206114fe6114e936611d7b565b92908160405193828580945193849201611c85565b8101600581520301902063ffffffff8092165f5260205260405f2061154f61152a600183549301611e4b565b604051938381869516855260201c166020840152606060408401526060830190611ca6565b0390f35b346107745760606003193601126107745767ffffffffffffffff60043581811161077457611585903690600401611d5d565b9061158e611c31565b90604435908111610774576115aa6115f2913690600401611d5d565b604051908185519581876115c56020998a9788809601611c85565b8101600481520301902063ffffffff8095165f52825260405f208260405194838680955193849201611c85565b8201908152030190205416604051908152f35b346107745760606003193601126107745761161e611c1e565b67ffffffffffffffff6024803582811161077457611640903690600401611c57565b9290916044358281116107745736602382011215610774578060040135928311610774576005923660248260051b8401011161077457945f9063ffffffff809816915b87811061168c57005b6040519082888337818381015f8152602093849103019020845f52825260405f208a82165f52825260ff600460405f2001541660038110156117585760018091036116dd575b506001915001611683565b604051848a8237838186810184815203019020855f52835260405f2087838a1b880101359373ffffffffffffffffffffffffffffffffffffffff8516809503610774576001945f525260405f20907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008254161790558a6116d2565b867f4e487b71000000000000000000000000000000000000000000000000000000005f5260216004525ffd5b3461077457602063ffffffff61179936611dba565b9390918260405193849283378101600581520301902091165f5260205261154f6117c8600160405f2001611e4b565b604051918291602083526020830190611ca6565b34610774576020806117f06114e936611d7b565b8101600381520301902063ffffffff8092165f52825260405f205416604051908152f35b3461077457606080600319360112610774576004359067ffffffffffffffff82116107745761184960c0923690600401611c57565b9290611853611c31565b9361185c611c44565b916040519061186a82611cd8565b5f82525f6020958382888096015288604082015288808201528860808201528260a0820152015282604051938492833781015f81520301902063ffffffff8095165f5282528360405f2091165f52815260405f20604051926118cb84611cd8565b815491858316855273ffffffffffffffffffffffffffffffffffffffff918285870194861c1684526118ff60018301611e4b565b6040870190815261191260028401611e4b565b93828801948552600461192760038601611e4b565b9460808a0195865201549560ff87169760a08a019860038110156106e957895260c08a019760081c60ff16151588526040519a8b9a828c525116908a0152511660408801525190860160e09052610100860161198291611ca6565b9151601f19928387820301608088015261199b91611ca6565b9051918582030160a08601526119b091611ca6565b915160c084016119bf91611ccb565b51151560e08301520390f35b34610774576060600319360112610774576119e4611c1e565b906119ed611c31565b9160443567ffffffffffffffff811161077457611a0e903690600401611c57565b909381858537838281015f81526020958691030190209063ffffffff80821692835f52865260405f20941693845f52855260405f2073ffffffffffffffffffffffffffffffffffffffff90818154881c16928315611bc057600482019384549160ff831660038110156106e957611b62579383611359989693611ad360ff947f5625cad747d0f660a99f4fe677a17d22db2a71e28505cd2148f61ad3cefc4f1d9e9a98611ac06003611b209a01611e4b565b91611acd60018701611e4b565b92611ff9565b15611b36575060027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008554161784555b548a1c169154169060405198899889528801526040870190611ccb565b606085015260a0608085015260a0840191611f48565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001178455611b03565b60648a604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601760248201527f4e6f646520616c7265616479206368616c6c656e6765640000000000000000006044820152fd5b606488604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601360248201527f4e6f646520646f6573206e6f74206578697374000000000000000000000000006044820152fd5b6004359063ffffffff8216820361077457565b6024359063ffffffff8216820361077457565b6044359063ffffffff8216820361077457565b9181601f840112156107745782359167ffffffffffffffff8311610774576020838186019501011161077457565b5f5b838110611c965750505f910152565b8181015183820152602001611c87565b90601f19601f602093611cc481518092818752878088019101611c85565b0116010190565b9060038210156106e95752565b60e0810190811067ffffffffffffffff82111761060057604052565b90601f601f19910116810190811067ffffffffffffffff82111761060057604052565b92919267ffffffffffffffff82116106005760405191611d416020601f19601f8401160184611cf4565b829481845281830111610774578281602093845f960137010152565b9080601f8301121561077457816020611d7893359101611d17565b90565b6040600319820112610774576004359067ffffffffffffffff821161077457611da691600401611d5d565b9060243563ffffffff811681036107745790565b6040600319820112610774576004359067ffffffffffffffff821161077457611de591600401611c57565b909160243563ffffffff811681036107745790565b90600182811c92168015611e41575b6020831014611e1457565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f1691611e09565b9060405191825f8254611e5d81611dfa565b908184526020946001916001811690815f14611ecb5750600114611e8d575b505050611e8b92500383611cf4565b565b5f90815285812095935091905b818310611eb3575050611e8b93508201015f8080611e7c565b85548884018501529485019487945091830191611e9a565b915050611e8b9593507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201015f8080611e7c565b6044359073ffffffffffffffffffffffffffffffffffffffff8216820361077457565b6020919283604051948593843782019081520301902090565b601f8260209493601f1993818652868601375f8582860101520116010190565b63ffffffff809116908114611f7d5760010190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b601f8211611fb757505050565b5f5260205f20906020601f840160051c83019310611fef575b601f0160051c01905b818110611fe4575050565b5f8155600101611fd9565b9091508190611fd0565b5115159291908361202e575b508261201e575b5081612016575090565b905051151590565b63ffffffff16151591505f61200c565b73ffffffffffffffffffffffffffffffffffffffff16151592505f612005565b815191906041830361207e576120779250602082015190606060408401519301515f1a90612140565b9192909190565b50505f9160029190565b60048110156106e9578061209a575050565b600181036120cc5760046040517ff645eedf000000000000000000000000000000000000000000000000000000008152fd5b6002810361210557602482604051907ffce698f70000000000000000000000000000000000000000000000000000000082526004820152fd5b60031461210f5750565b602490604051907fd78bce0c0000000000000000000000000000000000000000000000000000000082526004820152fd5b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a084116121cf579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa156121c4575f5173ffffffffffffffffffffffffffffffffffffffff8116156121ba57905f905f90565b505f906001905f90565b6040513d5f823e3d90fd5b5050505f916003919056fea26469706673582212201519ab4a7549720292a430109a35b6bfdba52f9b1a206a50216976eb3e54e85e64736f6c63430008170033",
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

// DkgNodeInfos is a free data retrieval call binding the contract method 0xa27c2f6c.
//
// Solidity: function dkgNodeInfos(bytes mrenclave, uint32 round, uint32 index) view returns(uint32 index, address validator, bytes dkgPubKey, bytes commPubKey, bytes rawQuote, uint8 chalStatus, bool finalized)
func (_DKG *DKGCaller) DkgNodeInfos(opts *bind.CallOpts, mrenclave []byte, round uint32, index uint32) (struct {
	Index      uint32
	Validator  common.Address
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	Finalized  bool
}, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "dkgNodeInfos", mrenclave, round, index)

	outstruct := new(struct {
		Index      uint32
		Validator  common.Address
		DkgPubKey  []byte
		CommPubKey []byte
		RawQuote   []byte
		ChalStatus uint8
		Finalized  bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Index = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Validator = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.DkgPubKey = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.CommPubKey = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.RawQuote = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.ChalStatus = *abi.ConvertType(out[5], new(uint8)).(*uint8)
	outstruct.Finalized = *abi.ConvertType(out[6], new(bool)).(*bool)

	return *outstruct, err

}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xa27c2f6c.
//
// Solidity: function dkgNodeInfos(bytes mrenclave, uint32 round, uint32 index) view returns(uint32 index, address validator, bytes dkgPubKey, bytes commPubKey, bytes rawQuote, uint8 chalStatus, bool finalized)
func (_DKG *DKGSession) DkgNodeInfos(mrenclave []byte, round uint32, index uint32) (struct {
	Index      uint32
	Validator  common.Address
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	Finalized  bool
}, error) {
	return _DKG.Contract.DkgNodeInfos(&_DKG.CallOpts, mrenclave, round, index)
}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xa27c2f6c.
//
// Solidity: function dkgNodeInfos(bytes mrenclave, uint32 round, uint32 index) view returns(uint32 index, address validator, bytes dkgPubKey, bytes commPubKey, bytes rawQuote, uint8 chalStatus, bool finalized)
func (_DKG *DKGCallerSession) DkgNodeInfos(mrenclave []byte, round uint32, index uint32) (struct {
	Index      uint32
	Validator  common.Address
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	ChalStatus uint8
	Finalized  bool
}, error) {
	return _DKG.Contract.DkgNodeInfos(&_DKG.CallOpts, mrenclave, round, index)
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

// GetNodeCount is a free data retrieval call binding the contract method 0x78e510d2.
//
// Solidity: function getNodeCount(bytes mrenclave, uint32 round) view returns(uint32)
func (_DKG *DKGCaller) GetNodeCount(opts *bind.CallOpts, mrenclave []byte, round uint32) (uint32, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "getNodeCount", mrenclave, round)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetNodeCount is a free data retrieval call binding the contract method 0x78e510d2.
//
// Solidity: function getNodeCount(bytes mrenclave, uint32 round) view returns(uint32)
func (_DKG *DKGSession) GetNodeCount(mrenclave []byte, round uint32) (uint32, error) {
	return _DKG.Contract.GetNodeCount(&_DKG.CallOpts, mrenclave, round)
}

// GetNodeCount is a free data retrieval call binding the contract method 0x78e510d2.
//
// Solidity: function getNodeCount(bytes mrenclave, uint32 round) view returns(uint32)
func (_DKG *DKGCallerSession) GetNodeCount(mrenclave []byte, round uint32) (uint32, error) {
	return _DKG.Contract.GetNodeCount(&_DKG.CallOpts, mrenclave, round)
}

// GetNodeInfo is a free data retrieval call binding the contract method 0x0c5a4020.
//
// Solidity: function getNodeInfo(bytes mrenclave, uint32 round, uint32 index) view returns((uint32,address,bytes,bytes,bytes,uint8,bool))
func (_DKG *DKGCaller) GetNodeInfo(opts *bind.CallOpts, mrenclave []byte, round uint32, index uint32) (IDKGNodeInfo, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "getNodeInfo", mrenclave, round, index)

	if err != nil {
		return *new(IDKGNodeInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IDKGNodeInfo)).(*IDKGNodeInfo)

	return out0, err

}

// GetNodeInfo is a free data retrieval call binding the contract method 0x0c5a4020.
//
// Solidity: function getNodeInfo(bytes mrenclave, uint32 round, uint32 index) view returns((uint32,address,bytes,bytes,bytes,uint8,bool))
func (_DKG *DKGSession) GetNodeInfo(mrenclave []byte, round uint32, index uint32) (IDKGNodeInfo, error) {
	return _DKG.Contract.GetNodeInfo(&_DKG.CallOpts, mrenclave, round, index)
}

// GetNodeInfo is a free data retrieval call binding the contract method 0x0c5a4020.
//
// Solidity: function getNodeInfo(bytes mrenclave, uint32 round, uint32 index) view returns((uint32,address,bytes,bytes,bytes,uint8,bool))
func (_DKG *DKGCallerSession) GetNodeInfo(mrenclave []byte, round uint32, index uint32) (IDKGNodeInfo, error) {
	return _DKG.Contract.GetNodeInfo(&_DKG.CallOpts, mrenclave, round, index)
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

// NodeCount is a free data retrieval call binding the contract method 0x159726be.
//
// Solidity: function nodeCount(bytes mrenclave, uint32 round) view returns(uint32 nodeCount)
func (_DKG *DKGCaller) NodeCount(opts *bind.CallOpts, mrenclave []byte, round uint32) (uint32, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "nodeCount", mrenclave, round)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// NodeCount is a free data retrieval call binding the contract method 0x159726be.
//
// Solidity: function nodeCount(bytes mrenclave, uint32 round) view returns(uint32 nodeCount)
func (_DKG *DKGSession) NodeCount(mrenclave []byte, round uint32) (uint32, error) {
	return _DKG.Contract.NodeCount(&_DKG.CallOpts, mrenclave, round)
}

// NodeCount is a free data retrieval call binding the contract method 0x159726be.
//
// Solidity: function nodeCount(bytes mrenclave, uint32 round) view returns(uint32 nodeCount)
func (_DKG *DKGCallerSession) NodeCount(mrenclave []byte, round uint32) (uint32, error) {
	return _DKG.Contract.NodeCount(&_DKG.CallOpts, mrenclave, round)
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

// FinalizeDKG is a paid mutator transaction binding the contract method 0xfe9dee68.
//
// Solidity: function finalizeDKG(uint32 round, uint32 index, bool finalized, bytes mrenclave, bytes globalPubKey, bytes signature) returns()
func (_DKG *DKGTransactor) FinalizeDKG(opts *bind.TransactOpts, round uint32, index uint32, finalized bool, mrenclave []byte, globalPubKey []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "finalizeDKG", round, index, finalized, mrenclave, globalPubKey, signature)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0xfe9dee68.
//
// Solidity: function finalizeDKG(uint32 round, uint32 index, bool finalized, bytes mrenclave, bytes globalPubKey, bytes signature) returns()
func (_DKG *DKGSession) FinalizeDKG(round uint32, index uint32, finalized bool, mrenclave []byte, globalPubKey []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.FinalizeDKG(&_DKG.TransactOpts, round, index, finalized, mrenclave, globalPubKey, signature)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0xfe9dee68.
//
// Solidity: function finalizeDKG(uint32 round, uint32 index, bool finalized, bytes mrenclave, bytes globalPubKey, bytes signature) returns()
func (_DKG *DKGTransactorSession) FinalizeDKG(round uint32, index uint32, finalized bool, mrenclave []byte, globalPubKey []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.FinalizeDKG(&_DKG.TransactOpts, round, index, finalized, mrenclave, globalPubKey, signature)
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

// RequestRemoteAttestationOnChain is a paid mutator transaction binding the contract method 0x08e48422.
//
// Solidity: function requestRemoteAttestationOnChain(uint32 targetIndex, uint32 round, bytes mrenclave) returns()
func (_DKG *DKGTransactor) RequestRemoteAttestationOnChain(opts *bind.TransactOpts, targetIndex uint32, round uint32, mrenclave []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "requestRemoteAttestationOnChain", targetIndex, round, mrenclave)
}

// RequestRemoteAttestationOnChain is a paid mutator transaction binding the contract method 0x08e48422.
//
// Solidity: function requestRemoteAttestationOnChain(uint32 targetIndex, uint32 round, bytes mrenclave) returns()
func (_DKG *DKGSession) RequestRemoteAttestationOnChain(targetIndex uint32, round uint32, mrenclave []byte) (*types.Transaction, error) {
	return _DKG.Contract.RequestRemoteAttestationOnChain(&_DKG.TransactOpts, targetIndex, round, mrenclave)
}

// RequestRemoteAttestationOnChain is a paid mutator transaction binding the contract method 0x08e48422.
//
// Solidity: function requestRemoteAttestationOnChain(uint32 targetIndex, uint32 round, bytes mrenclave) returns()
func (_DKG *DKGTransactorSession) RequestRemoteAttestationOnChain(targetIndex uint32, round uint32, mrenclave []byte) (*types.Transaction, error) {
	return _DKG.Contract.RequestRemoteAttestationOnChain(&_DKG.TransactOpts, targetIndex, round, mrenclave)
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
	Index        uint32
	Finalized    bool
	Mrenclave    []byte
	GlobalPubKey []byte
	Signature    []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDKGFinalized is a free log retrieval operation binding the contract event 0x694566864b727668476d766533aac9ac9c8f2a98c4c58322f1cc6b01258a5c48.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, uint32 index, bool finalized, bytes mrenclave, bytes globalPubKey, bytes signature)
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

// WatchDKGFinalized is a free log subscription operation binding the contract event 0x694566864b727668476d766533aac9ac9c8f2a98c4c58322f1cc6b01258a5c48.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, uint32 index, bool finalized, bytes mrenclave, bytes globalPubKey, bytes signature)
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

// ParseDKGFinalized is a log parse operation binding the contract event 0x694566864b727668476d766533aac9ac9c8f2a98c4c58322f1cc6b01258a5c48.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, uint32 index, bool finalized, bytes mrenclave, bytes globalPubKey, bytes signature)
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
	Index      uint32
	DkgPubKey  []byte
	CommPubKey []byte
	RawQuote   []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDKGInitialized is a free log retrieval operation binding the contract event 0x5c5aa077fdba8d39c55de6079c75de88ec2031742c08609237382f11414de616.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes mrenclave, uint32 round, uint32 index, bytes dkgPubKey, bytes commPubKey, bytes rawQuote)
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

// WatchDKGInitialized is a free log subscription operation binding the contract event 0x5c5aa077fdba8d39c55de6079c75de88ec2031742c08609237382f11414de616.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes mrenclave, uint32 round, uint32 index, bytes dkgPubKey, bytes commPubKey, bytes rawQuote)
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

// ParseDKGInitialized is a log parse operation binding the contract event 0x5c5aa077fdba8d39c55de6079c75de88ec2031742c08609237382f11414de616.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes mrenclave, uint32 round, uint32 index, bytes dkgPubKey, bytes commPubKey, bytes rawQuote)
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

// DKGInvalidDKGInitializationIterator is returned from FilterInvalidDKGInitialization and is used to iterate over the raw logs and unpacked data for InvalidDKGInitialization events raised by the DKG contract.
type DKGInvalidDKGInitializationIterator struct {
	Event *DKGInvalidDKGInitialization // Event containing the contract specifics and raw log

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
func (it *DKGInvalidDKGInitializationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGInvalidDKGInitialization)
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
		it.Event = new(DKGInvalidDKGInitialization)
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
func (it *DKGInvalidDKGInitializationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGInvalidDKGInitializationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGInvalidDKGInitialization represents a InvalidDKGInitialization event raised by the DKG contract.
type DKGInvalidDKGInitialization struct {
	Round     uint32
	Index     uint32
	Validator common.Address
	Mrenclave []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterInvalidDKGInitialization is a free log retrieval operation binding the contract event 0xe6250d83180fe43a8d437f7314068cdc648a6a64ddf77cde43af1c6a25725eb3.
//
// Solidity: event InvalidDKGInitialization(uint32 round, uint32 index, address validator, bytes mrenclave)
func (_DKG *DKGFilterer) FilterInvalidDKGInitialization(opts *bind.FilterOpts) (*DKGInvalidDKGInitializationIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "InvalidDKGInitialization")
	if err != nil {
		return nil, err
	}
	return &DKGInvalidDKGInitializationIterator{contract: _DKG.contract, event: "InvalidDKGInitialization", logs: logs, sub: sub}, nil
}

// WatchInvalidDKGInitialization is a free log subscription operation binding the contract event 0xe6250d83180fe43a8d437f7314068cdc648a6a64ddf77cde43af1c6a25725eb3.
//
// Solidity: event InvalidDKGInitialization(uint32 round, uint32 index, address validator, bytes mrenclave)
func (_DKG *DKGFilterer) WatchInvalidDKGInitialization(opts *bind.WatchOpts, sink chan<- *DKGInvalidDKGInitialization) (event.Subscription, error) {

	logs, sub, err := _DKG.contract.WatchLogs(opts, "InvalidDKGInitialization")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGInvalidDKGInitialization)
				if err := _DKG.contract.UnpackLog(event, "InvalidDKGInitialization", log); err != nil {
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

// ParseInvalidDKGInitialization is a log parse operation binding the contract event 0xe6250d83180fe43a8d437f7314068cdc648a6a64ddf77cde43af1c6a25725eb3.
//
// Solidity: event InvalidDKGInitialization(uint32 round, uint32 index, address validator, bytes mrenclave)
func (_DKG *DKGFilterer) ParseInvalidDKGInitialization(log types.Log) (*DKGInvalidDKGInitialization, error) {
	event := new(DKGInvalidDKGInitialization)
	if err := _DKG.contract.UnpackLog(event, "InvalidDKGInitialization", log); err != nil {
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

// DKGRegistrationChallengedIterator is returned from FilterRegistrationChallenged and is used to iterate over the raw logs and unpacked data for RegistrationChallenged events raised by the DKG contract.
type DKGRegistrationChallengedIterator struct {
	Event *DKGRegistrationChallenged // Event containing the contract specifics and raw log

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
func (it *DKGRegistrationChallengedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGRegistrationChallenged)
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
		it.Event = new(DKGRegistrationChallenged)
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
func (it *DKGRegistrationChallengedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGRegistrationChallengedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGRegistrationChallenged represents a RegistrationChallenged event raised by the DKG contract.
type DKGRegistrationChallenged struct {
	Round      uint32
	Mrenclave  []byte
	Challenger common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRegistrationChallenged is a free log retrieval operation binding the contract event 0x8e6f8f23941b429128eca7a26ddf8eb34d0c1f64f6fc2d4a4625ec0222bc1d67.
//
// Solidity: event RegistrationChallenged(uint32 round, bytes mrenclave, address indexed challenger)
func (_DKG *DKGFilterer) FilterRegistrationChallenged(opts *bind.FilterOpts, challenger []common.Address) (*DKGRegistrationChallengedIterator, error) {

	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _DKG.contract.FilterLogs(opts, "RegistrationChallenged", challengerRule)
	if err != nil {
		return nil, err
	}
	return &DKGRegistrationChallengedIterator{contract: _DKG.contract, event: "RegistrationChallenged", logs: logs, sub: sub}, nil
}

// WatchRegistrationChallenged is a free log subscription operation binding the contract event 0x8e6f8f23941b429128eca7a26ddf8eb34d0c1f64f6fc2d4a4625ec0222bc1d67.
//
// Solidity: event RegistrationChallenged(uint32 round, bytes mrenclave, address indexed challenger)
func (_DKG *DKGFilterer) WatchRegistrationChallenged(opts *bind.WatchOpts, sink chan<- *DKGRegistrationChallenged, challenger []common.Address) (event.Subscription, error) {

	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _DKG.contract.WatchLogs(opts, "RegistrationChallenged", challengerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGRegistrationChallenged)
				if err := _DKG.contract.UnpackLog(event, "RegistrationChallenged", log); err != nil {
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

// ParseRegistrationChallenged is a log parse operation binding the contract event 0x8e6f8f23941b429128eca7a26ddf8eb34d0c1f64f6fc2d4a4625ec0222bc1d67.
//
// Solidity: event RegistrationChallenged(uint32 round, bytes mrenclave, address indexed challenger)
func (_DKG *DKGFilterer) ParseRegistrationChallenged(log types.Log) (*DKGRegistrationChallenged, error) {
	event := new(DKGRegistrationChallenged)
	if err := _DKG.contract.UnpackLog(event, "RegistrationChallenged", log); err != nil {
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
	Index      uint32
	Validator  common.Address
	ChalStatus uint8
	Round      uint32
	Mrenclave  []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRemoteAttestationProcessedOnChain is a free log retrieval operation binding the contract event 0x5625cad747d0f660a99f4fe677a17d22db2a71e28505cd2148f61ad3cefc4f1d.
//
// Solidity: event RemoteAttestationProcessedOnChain(uint32 index, address validator, uint8 chalStatus, uint32 round, bytes mrenclave)
func (_DKG *DKGFilterer) FilterRemoteAttestationProcessedOnChain(opts *bind.FilterOpts) (*DKGRemoteAttestationProcessedOnChainIterator, error) {

	logs, sub, err := _DKG.contract.FilterLogs(opts, "RemoteAttestationProcessedOnChain")
	if err != nil {
		return nil, err
	}
	return &DKGRemoteAttestationProcessedOnChainIterator{contract: _DKG.contract, event: "RemoteAttestationProcessedOnChain", logs: logs, sub: sub}, nil
}

// WatchRemoteAttestationProcessedOnChain is a free log subscription operation binding the contract event 0x5625cad747d0f660a99f4fe677a17d22db2a71e28505cd2148f61ad3cefc4f1d.
//
// Solidity: event RemoteAttestationProcessedOnChain(uint32 index, address validator, uint8 chalStatus, uint32 round, bytes mrenclave)
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

// ParseRemoteAttestationProcessedOnChain is a log parse operation binding the contract event 0x5625cad747d0f660a99f4fe677a17d22db2a71e28505cd2148f61ad3cefc4f1d.
//
// Solidity: event RemoteAttestationProcessedOnChain(uint32 index, address validator, uint8 chalStatus, uint32 round, bytes mrenclave)
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
