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
	Index        uint32
	Validator    common.Address
	DkgPubKey    []byte
	CommPubKey   []byte
	RemoteReport []byte
	Commitments  []byte
	ChalStatus   uint8
	Finalized    bool
}

// DKGMetaData contains all meta data concerning the DKG contract.
var DKGMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"complainDeals\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dealComplaints\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainant\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dkgNodeInfos\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitments\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"finalized\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"finalized\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getNodeCount\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeInfo\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDKG.NodeInfo\",\"components\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitments\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"finalized\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isActiveValidator\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nodeCount\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"nodeCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestRemoteAttestationOnChain\",\"inputs\":[{\"name\":\"targetIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitActiveValSet\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"valSet\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateDKGCommitments\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitments\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"valSets\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DKGCommitmentsUpdated\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"total\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"commitments\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGFinalized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"finalized\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGInitialized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"dkgPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"commPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"remoteReport\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealComplaintsSubmitted\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"indexed\":false,\"internalType\":\"uint32[]\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealVerified\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"recipientIndex\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidDKGInitialization\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidDeal\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RegistrationChallenged\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"challenger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoteAttestationProcessedOnChain\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeScheduled\",\"inputs\":[{\"name\":\"activationHeight\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
	Bin: "0x6080806040523461001657611f0b908161001b8239f35b5f80fdfe6080806040526004361015610012575f80fd5b5f3560e01c90816308e4842214611826575080630c5a40201461163d578063159726be146115c65780632e266b4c14611448578063496c634c146111585780634d71b63414610fb35780634f1bf88d14610f0f57806378e510d214610ea55780637a57476514610c02578063a27c2f6c14610ab9578063c2009b9214610a03578063c3c77ff9146101475763fb5a783a146100ab575f80fd5b346101435760606003193601126101435760043567ffffffffffffffff8111610143576100dc903690600401611add565b63ffffffff6100e9611ab7565b9160206100f4611c02565b948260405193849283378101600181520301902091165f5260205273ffffffffffffffffffffffffffffffffffffffff60405f2091165f52602052602060ff60405f2054166040519015158152f35b5f80fd5b346101435760a060031936011261014357610160611aa4565b60243567ffffffffffffffff811161014357610180903690600401611add565b60449291923567ffffffffffffffff8111610143576101a3903690600401611add565b9060643567ffffffffffffffff8111610143576101c4903690600401611add565b9560843567ffffffffffffffff8111610143576101e5903690600401611add565b92909760405187848237602081898101600181520301902063ffffffff89165f5260205260405f20335f5260205260ff60405f20541680156109fb575b1561099d5760405187848237602081898101600381520301902063ffffffff89165f5260205263ffffffff60405f20541697604051888582376020818a8101600381520301902063ffffffff82165f5260205260405f20805463ffffffff80821614610970577fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000063ffffffff600181841601169116179055604051906102c782611b5e565b8982523360208301526102db368989611b9e565b60408301526102eb368486611b9e565b60608301526102fb36878d611b9e565b608083015260405180602081011067ffffffffffffffff60208301111761079157602081016040525f815260a08301525f60c08301525f60e0830152604051898682376020818b81015f81520301902063ffffffff82165f5260205260405f208a5f5260205260405f2063ffffffff8351168154907fffffffffffffffff00000000000000000000000000000000000000000000000077ffffffffffffffffffffffffffffffffffffffff00000000602087015160201b16921617178155604083015180519067ffffffffffffffff8211610791576103ea826103e16001860154611c25565b60018601611e86565b602090601f83116001146108e35761043792915f91836107be575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60018201555b606083015180519067ffffffffffffffff82116107915761046e826104656002860154611c25565b60028601611e86565b602090601f8311600114610856576104ba92915f91836107be5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60028201555b608083015180519067ffffffffffffffff8211610791576104f1826104e86003860154611c25565b60038601611e86565b602090601f83116001146107c95761053d92915f91836107be5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60038201555b60a083015180519067ffffffffffffffff8211610791576105748261056b6004860154611c25565b60048601611e86565b602090601f83116001146107035791806105c492600595945f926106f85750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60048201555b019760c08301519760038910156106cb577f5c5aa077fdba8d39c55de6079c75de88ec2031742c08609237382f11414de6169b63ffffffff61068f6106b69961067e8e60c09f6106c69f9a6106a89b60ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0060e09454169116178355015115157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff00835492151560081b169116179055565b6040519e8f9e8f8181520191611d37565b941660208c015260408b015289830360608b0152611d37565b918683036080880152611d37565b9083820360a08501523396611d37565b0390a2005b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b015190505f80610405565b90600484015f5260205f20915f5b601f1985168110610779575091839160019383601f19600598971610610742575b505050811b0160048201556105ca565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558f8080610732565b91926020600181928685015181550194019201610711565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b015190508f80610405565b9190600384015f5260205f20905f935b601f198416851061083b576001945083601f19811610610804575b505050811b016003820155610543565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e80806107f4565b818101518355602094850194600190930192909101906107d9565b9190600284015f5260205f20905f935b601f19841685106108c8576001945083601f19811610610891575b505050811b0160028201556104c0565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e8080610881565b81810151835560209485019460019093019290910190610866565b9190600184015f5260205f20905f935b601f1984168510610955576001945083601f1981161061091e575b505050811b01600182015561043d565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558e808061090e565b818101518355602094850194600190930192909101906108f3565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f56616c696461746f72206e6f7420696e206163746976652073657400000000006044820152fd5b506001610222565b346101435760806003193601126101435760043567ffffffffffffffff811161014357610a34903690600401611be4565b610a3c611ab7565b610a44611aca565b6064359173ffffffffffffffffffffffffffffffffffffffff8316809303610143576040518481809651610a7e816020998a809601611b0b565b8101600281520301902063ffffffff8092165f52845260405f2091165f52825260405f20905f52815260ff60405f2054166040519015158152f35b346101435760606003193601126101435760043567ffffffffffffffff811161014357610aff73ffffffffffffffffffffffffffffffffffffffff913690600401611be4565b610b07611ab7565b90610b10611aca565b6040518281809451610b288160209788809601611b0b565b81015f81520301902063ffffffff8094165f5282528260405f2091165f52815260ff610be360405f20610bd5815495610bc7610b6660018501611c76565b610bb9610b7560028701611c76565b91610b8260038801611c76565b946005610b9160048a01611c76565b9801549a6040519d8e9d8e6101009482169052821c16908d01528060408d01528b0190611b2c565b9089820360608b0152611b2c565b908782036080890152611b2c565b9085820360a0870152611b2c565b91610bf360c08501838316611b51565b60081c16151560e08301520390f35b3461014357608060031936011261014357610c1b611aa4565b610c23611ab7565b60443567ffffffffffffffff808211610143573660238301121561014357816004013593602491808611610791576005908660051b9060209560405198610c6c8885018b611b7b565b8952868901602481948301019136831161014357602401905b828210610e895750505060643590811161014357610ca7903690600401611add565b906040518282823787818481015f8152030190209263ffffffff80961693845f5288528560405f20991698895f52885273ffffffffffffffffffffffffffffffffffffffff60405f2054891c163303610e2b575f5b8a51811015610da4576040518484823789818681016002815203019020855f52895260405f208b51821015610d785790600191888b8e848b1b010151165f528a5260405f20335f528a5260405f20827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0082541617905501610cfc565b887f4e487b71000000000000000000000000000000000000000000000000000000005f5260326004525ffd5b509350889285898960405196608088019288526080828901525180925260a0870197925f905b838210610e12577f1c2112af5fd37661e3dd248d701decebf291f13de6a411176337ef21a7b1a6308980610e0d8d8c8c8c60408601528483036060860152611d37565b0390a1005b845181168a529882019893820193600190910190610dca565b606488604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601360248201527f496e76616c696420636f6d706c61696e616e74000000000000000000000000006044820152fd5b813563ffffffff81168103610143578152908801908801610c85565b346101435760406003193601126101435760043567ffffffffffffffff811161014357610ed86020913690600401611add565b82610ee1611ab7565b928260405193849283378101600381520301902063ffffffff8092165f52825260405f205416604051908152f35b346101435760606003193601126101435760043567ffffffffffffffff811161014357610f40903690600401611be4565b610f48611ab7565b63ffffffff610f6e6020610f5a611c02565b948160405193828580945193849201611b0b565b8101600181520301902091165f5260205273ffffffffffffffffffffffffffffffffffffffff60405f2091165f52602052602060ff60405f2054166040519015158152f35b346101435760a060031936011261014357610fcc611aa4565b610fd4611ab7565b60443590811515928383036101435767ffffffffffffffff9360643585811161014357611005903690600401611add565b90956084359081116101435761101f903690600401611add565b919096604051828282376020818481015f8152030190209463ffffffff80911695865f5260205260405f20961695865f5260205260405f209161107d73ffffffffffffffffffffffffffffffffffffffff845460201c163314611d57565b600583019460ff86541660038110156106cb577fbeb92472b2143160d5e8c3fd25327aabddd7a0f4b80ade42661dd5b4e5ff4f90996106c6976110d5600261113f986110cf600161111c971415611dbc565b01611c76565b5115158061114f575b6110e790611e21565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff00835492151560081b169116179055565b60405197889788526020880152604087015260a0606087015260a0860191611d37565b9083820360808501523396611d37565b508715156110de565b346101435760e060031936011261014357611171611aa4565b611179611ab7565b90611182611aca565b916064359163ffffffff938484168094036101435767ffffffffffffffff90608435828111610143576111b9903690600401611add565b96909560a435848111610143576111d4903690600401611add565b91909260c43594868611610143576111f18b963690600401611add565b969097838c83604051948592833781015f81526020948591030190209a16998a5f52825260405f20855f52825260405f2061124673ffffffffffffffffffffffffffffffffffffffff8254851c163314611d57565b60ff6005820154169060038210156106cb57600490611269600180941415611dbc565b61127560028201611c76565b5115158061143f575b80611436575b61128d90611e21565b01918711610791576112a9876112a38454611c25565b84611e86565b5f90601f8811600114611380575093611365979387937f1af28d6ae67078289128cdbaca9c65cda5fda4d71a49bb5e2d917ef3bf5f44dd9e9f9a979361132a866106c69e9b6113579a5f91611375575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b826040519e8f9e8f5216908d01521660408b015260608a015260e060808a015260e0890191611d37565b9186830360a0880152611d37565b9083820360c08501523396611d37565b90508901355f6112f9565b90601f198816835f52845f20925f905b82821061141f575050937f1af28d6ae67078289128cdbaca9c65cda5fda4d71a49bb5e2d917ef3bf5f44dd9e9f9a979361135797936113659b97938b6106c69f9c98106113e7575b5050600186811b01905561132d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88960031b161c19908901351690555f806113d8565b8b8401358555938401939286019290860190611390565b508a1515611284565b5088151561127e565b3461014357606060031936011261014357611461611aa4565b67ffffffffffffffff6024803582811161014357611483903690600401611add565b9290916044358281116101435736602382011215610143578060040135928311610143576005923660248260051b8401011161014357945f9063ffffffff809816915b8781106114cf57005b6040519082888337818381015f8152602093849103019020845f52825260405f208a82165f52825260ff8760405f20015416600381101561159a57600180910361151f575b5060019150016114c6565b604051848a8237838186810184815203019020855f52835260405f2087838a1b880101359373ffffffffffffffffffffffffffffffffffffffff8516809503610143576001945f525260405f20907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008254161790558a611514565b867f4e487b71000000000000000000000000000000000000000000000000000000005f5260216004525ffd5b346101435760406003193601126101435760043567ffffffffffffffff8111610143576115f96020913690600401611be4565b61161982611605611ab7565b928160405193828580945193849201611b0b565b8101600381520301902063ffffffff8092165f52825260405f205416604051908152f35b3461014357606080600319360112610143576004359067ffffffffffffffff82116101435761167260e0923690600401611add565b929061167c611ab7565b93611685611aca565b916040519061169382611b5e565b5f82525f6020958382888096015288604082015288808201528860808201528860a08201528260c0820152015282604051938492833781015f81520301902063ffffffff8095165f5282528360405f2091165f52815260405f2091604051926116fb84611b5e565b805490858216855273ffffffffffffffffffffffffffffffffffffffff908185870193861c16835261172f60018201611c76565b916040870192835261174360028301611c76565b9285880193845261175660038401611c76565b9160808901928352600561176c60048601611c76565b9460a08b0195865201549560ff87169860c08b019960038110156106cb578a5260e08b019760081c60ff16151588526040519b8c9b828d525116908b01525116604089015251610100809689015261012088016117c891611b2c565b9251601f1993848982030160808a01526117e191611b2c565b905190838882030160a08901526117f791611b2c565b9051918682030160c087015261180c91611b2c565b925160e0850161181b91611b51565b511515908301520390f35b346101435760606003193601126101435761183f611aa4565b90611848611ab7565b9160443567ffffffffffffffff811161014357611869903690600401611add565b909381858537838281015f81526020958691030190209063ffffffff80911691825f52855260405f20931692835f52845260405f2073ffffffffffffffffffffffffffffffffffffffff808254871c16159182611a4657600581019283549060ff821660038110156106cb576119e8577f5625cad747d0f660a99f4fe677a17d22db2a71e28505cd2148f61ad3cefc4f1d99959383610e0d98969360ff9361191660036119889801611c76565b61192260028501611c76565b9051151591826119df575b50816119d5575b816119ca575b501561199e575060027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008554161784555b548a1c169154169060405198899889528801526040870190611b51565b606085015260a0608085015260a0840191611d37565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117845561196b565b90505115158e61193a565b8815159150611934565b1591508f61192d565b606489604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601760248201527f4e6f646520616c7265616479206368616c6c656e6765640000000000000000006044820152fd5b606487604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601360248201527f4e6f646520646f6573206e6f74206578697374000000000000000000000000006044820152fd5b6004359063ffffffff8216820361014357565b6024359063ffffffff8216820361014357565b6044359063ffffffff8216820361014357565b9181601f840112156101435782359167ffffffffffffffff8311610143576020838186019501011161014357565b5f5b838110611b1c5750505f910152565b8181015183820152602001611b0d565b90601f19601f602093611b4a81518092818752878088019101611b0b565b0116010190565b9060038210156106cb5752565b610100810190811067ffffffffffffffff82111761079157604052565b90601f601f19910116810190811067ffffffffffffffff82111761079157604052565b92919267ffffffffffffffff82116107915760405191611bc86020601f19601f8401160184611b7b565b829481845281830111610143578281602093845f960137010152565b9080601f8301121561014357816020611bff93359101611b9e565b90565b6044359073ffffffffffffffffffffffffffffffffffffffff8216820361014357565b90600182811c92168015611c6c575b6020831014611c3f57565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f1691611c34565b9060405191825f8254611c8881611c25565b908184526020946001916001811690815f14611cf65750600114611cb8575b505050611cb692500383611b7b565b565b5f90815285812095935091905b818310611cde575050611cb693508201015f8080611ca7565b85548884018501529485019487945091830191611cc5565b915050611cb69593507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201015f8080611ca7565b601f8260209493601f1993818652868601375f8582860101520116010190565b15611d5e57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f496e76616c69642076616c696461746f720000000000000000000000000000006044820152fd5b15611dc357565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4e6f64652077617320696e76616c6964617465640000000000000000000000006044820152fd5b15611e2857565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f496e76616c6964207369676e61747572650000000000000000000000000000006044820152fd5b601f8211611e9357505050565b5f5260205f20906020601f840160051c83019310611ecb575b601f0160051c01905b818110611ec0575050565b5f8155600101611eb5565b9091508190611eac56fea2646970667358221220b9ecadd13cfc7842c9e23cae251a3c1ed476640461b4b9b4e19bec51c53bc0a164736f6c63430008170033",
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
// Solidity: function dkgNodeInfos(bytes mrenclave, uint32 round, uint32 index) view returns(uint32 index, address validator, bytes dkgPubKey, bytes commPubKey, bytes remoteReport, bytes commitments, uint8 chalStatus, bool finalized)
func (_DKG *DKGCaller) DkgNodeInfos(opts *bind.CallOpts, mrenclave []byte, round uint32, index uint32) (struct {
	Index        uint32
	Validator    common.Address
	DkgPubKey    []byte
	CommPubKey   []byte
	RemoteReport []byte
	Commitments  []byte
	ChalStatus   uint8
	Finalized    bool
}, error) {
	var out []interface{}
	err := _DKG.contract.Call(opts, &out, "dkgNodeInfos", mrenclave, round, index)

	outstruct := new(struct {
		Index        uint32
		Validator    common.Address
		DkgPubKey    []byte
		CommPubKey   []byte
		RemoteReport []byte
		Commitments  []byte
		ChalStatus   uint8
		Finalized    bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Index = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Validator = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.DkgPubKey = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.CommPubKey = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.RemoteReport = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.Commitments = *abi.ConvertType(out[5], new([]byte)).(*[]byte)
	outstruct.ChalStatus = *abi.ConvertType(out[6], new(uint8)).(*uint8)
	outstruct.Finalized = *abi.ConvertType(out[7], new(bool)).(*bool)

	return *outstruct, err

}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xa27c2f6c.
//
// Solidity: function dkgNodeInfos(bytes mrenclave, uint32 round, uint32 index) view returns(uint32 index, address validator, bytes dkgPubKey, bytes commPubKey, bytes remoteReport, bytes commitments, uint8 chalStatus, bool finalized)
func (_DKG *DKGSession) DkgNodeInfos(mrenclave []byte, round uint32, index uint32) (struct {
	Index        uint32
	Validator    common.Address
	DkgPubKey    []byte
	CommPubKey   []byte
	RemoteReport []byte
	Commitments  []byte
	ChalStatus   uint8
	Finalized    bool
}, error) {
	return _DKG.Contract.DkgNodeInfos(&_DKG.CallOpts, mrenclave, round, index)
}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xa27c2f6c.
//
// Solidity: function dkgNodeInfos(bytes mrenclave, uint32 round, uint32 index) view returns(uint32 index, address validator, bytes dkgPubKey, bytes commPubKey, bytes remoteReport, bytes commitments, uint8 chalStatus, bool finalized)
func (_DKG *DKGCallerSession) DkgNodeInfos(mrenclave []byte, round uint32, index uint32) (struct {
	Index        uint32
	Validator    common.Address
	DkgPubKey    []byte
	CommPubKey   []byte
	RemoteReport []byte
	Commitments  []byte
	ChalStatus   uint8
	Finalized    bool
}, error) {
	return _DKG.Contract.DkgNodeInfos(&_DKG.CallOpts, mrenclave, round, index)
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
// Solidity: function getNodeInfo(bytes mrenclave, uint32 round, uint32 index) view returns((uint32,address,bytes,bytes,bytes,bytes,uint8,bool))
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
// Solidity: function getNodeInfo(bytes mrenclave, uint32 round, uint32 index) view returns((uint32,address,bytes,bytes,bytes,bytes,uint8,bool))
func (_DKG *DKGSession) GetNodeInfo(mrenclave []byte, round uint32, index uint32) (IDKGNodeInfo, error) {
	return _DKG.Contract.GetNodeInfo(&_DKG.CallOpts, mrenclave, round, index)
}

// GetNodeInfo is a free data retrieval call binding the contract method 0x0c5a4020.
//
// Solidity: function getNodeInfo(bytes mrenclave, uint32 round, uint32 index) view returns((uint32,address,bytes,bytes,bytes,bytes,uint8,bool))
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

// FinalizeDKG is a paid mutator transaction binding the contract method 0x4d71b634.
//
// Solidity: function finalizeDKG(uint32 round, uint32 index, bool finalized, bytes mrenclave, bytes signature) returns()
func (_DKG *DKGTransactor) FinalizeDKG(opts *bind.TransactOpts, round uint32, index uint32, finalized bool, mrenclave []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "finalizeDKG", round, index, finalized, mrenclave, signature)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x4d71b634.
//
// Solidity: function finalizeDKG(uint32 round, uint32 index, bool finalized, bytes mrenclave, bytes signature) returns()
func (_DKG *DKGSession) FinalizeDKG(round uint32, index uint32, finalized bool, mrenclave []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.FinalizeDKG(&_DKG.TransactOpts, round, index, finalized, mrenclave, signature)
}

// FinalizeDKG is a paid mutator transaction binding the contract method 0x4d71b634.
//
// Solidity: function finalizeDKG(uint32 round, uint32 index, bool finalized, bytes mrenclave, bytes signature) returns()
func (_DKG *DKGTransactorSession) FinalizeDKG(round uint32, index uint32, finalized bool, mrenclave []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.FinalizeDKG(&_DKG.TransactOpts, round, index, finalized, mrenclave, signature)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0xc3c77ff9.
//
// Solidity: function initializeDKG(uint32 round, bytes mrenclave, bytes dkgPubKey, bytes commPubKey, bytes remoteReport) returns()
func (_DKG *DKGTransactor) InitializeDKG(opts *bind.TransactOpts, round uint32, mrenclave []byte, dkgPubKey []byte, commPubKey []byte, remoteReport []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "initializeDKG", round, mrenclave, dkgPubKey, commPubKey, remoteReport)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0xc3c77ff9.
//
// Solidity: function initializeDKG(uint32 round, bytes mrenclave, bytes dkgPubKey, bytes commPubKey, bytes remoteReport) returns()
func (_DKG *DKGSession) InitializeDKG(round uint32, mrenclave []byte, dkgPubKey []byte, commPubKey []byte, remoteReport []byte) (*types.Transaction, error) {
	return _DKG.Contract.InitializeDKG(&_DKG.TransactOpts, round, mrenclave, dkgPubKey, commPubKey, remoteReport)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0xc3c77ff9.
//
// Solidity: function initializeDKG(uint32 round, bytes mrenclave, bytes dkgPubKey, bytes commPubKey, bytes remoteReport) returns()
func (_DKG *DKGTransactorSession) InitializeDKG(round uint32, mrenclave []byte, dkgPubKey []byte, commPubKey []byte, remoteReport []byte) (*types.Transaction, error) {
	return _DKG.Contract.InitializeDKG(&_DKG.TransactOpts, round, mrenclave, dkgPubKey, commPubKey, remoteReport)
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

// UpdateDKGCommitments is a paid mutator transaction binding the contract method 0x496c634c.
//
// Solidity: function updateDKGCommitments(uint32 round, uint32 total, uint32 threshold, uint32 index, bytes mrenclave, bytes commitments, bytes signature) returns()
func (_DKG *DKGTransactor) UpdateDKGCommitments(opts *bind.TransactOpts, round uint32, total uint32, threshold uint32, index uint32, mrenclave []byte, commitments []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "updateDKGCommitments", round, total, threshold, index, mrenclave, commitments, signature)
}

// UpdateDKGCommitments is a paid mutator transaction binding the contract method 0x496c634c.
//
// Solidity: function updateDKGCommitments(uint32 round, uint32 total, uint32 threshold, uint32 index, bytes mrenclave, bytes commitments, bytes signature) returns()
func (_DKG *DKGSession) UpdateDKGCommitments(round uint32, total uint32, threshold uint32, index uint32, mrenclave []byte, commitments []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.UpdateDKGCommitments(&_DKG.TransactOpts, round, total, threshold, index, mrenclave, commitments, signature)
}

// UpdateDKGCommitments is a paid mutator transaction binding the contract method 0x496c634c.
//
// Solidity: function updateDKGCommitments(uint32 round, uint32 total, uint32 threshold, uint32 index, bytes mrenclave, bytes commitments, bytes signature) returns()
func (_DKG *DKGTransactorSession) UpdateDKGCommitments(round uint32, total uint32, threshold uint32, index uint32, mrenclave []byte, commitments []byte, signature []byte) (*types.Transaction, error) {
	return _DKG.Contract.UpdateDKGCommitments(&_DKG.TransactOpts, round, total, threshold, index, mrenclave, commitments, signature)
}

// DKGDKGCommitmentsUpdatedIterator is returned from FilterDKGCommitmentsUpdated and is used to iterate over the raw logs and unpacked data for DKGCommitmentsUpdated events raised by the DKG contract.
type DKGDKGCommitmentsUpdatedIterator struct {
	Event *DKGDKGCommitmentsUpdated // Event containing the contract specifics and raw log

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
func (it *DKGDKGCommitmentsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DKGDKGCommitmentsUpdated)
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
		it.Event = new(DKGDKGCommitmentsUpdated)
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
func (it *DKGDKGCommitmentsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DKGDKGCommitmentsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DKGDKGCommitmentsUpdated represents a DKGCommitmentsUpdated event raised by the DKG contract.
type DKGDKGCommitmentsUpdated struct {
	MsgSender   common.Address
	Round       uint32
	Total       uint32
	Threshold   uint32
	Index       uint32
	Commitments []byte
	Signature   []byte
	Mrenclave   []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDKGCommitmentsUpdated is a free log retrieval operation binding the contract event 0x1af28d6ae67078289128cdbaca9c65cda5fda4d71a49bb5e2d917ef3bf5f44dd.
//
// Solidity: event DKGCommitmentsUpdated(address indexed msgSender, uint32 round, uint32 total, uint32 threshold, uint32 index, bytes commitments, bytes signature, bytes mrenclave)
func (_DKG *DKGFilterer) FilterDKGCommitmentsUpdated(opts *bind.FilterOpts, msgSender []common.Address) (*DKGDKGCommitmentsUpdatedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _DKG.contract.FilterLogs(opts, "DKGCommitmentsUpdated", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &DKGDKGCommitmentsUpdatedIterator{contract: _DKG.contract, event: "DKGCommitmentsUpdated", logs: logs, sub: sub}, nil
}

// WatchDKGCommitmentsUpdated is a free log subscription operation binding the contract event 0x1af28d6ae67078289128cdbaca9c65cda5fda4d71a49bb5e2d917ef3bf5f44dd.
//
// Solidity: event DKGCommitmentsUpdated(address indexed msgSender, uint32 round, uint32 total, uint32 threshold, uint32 index, bytes commitments, bytes signature, bytes mrenclave)
func (_DKG *DKGFilterer) WatchDKGCommitmentsUpdated(opts *bind.WatchOpts, sink chan<- *DKGDKGCommitmentsUpdated, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _DKG.contract.WatchLogs(opts, "DKGCommitmentsUpdated", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DKGDKGCommitmentsUpdated)
				if err := _DKG.contract.UnpackLog(event, "DKGCommitmentsUpdated", log); err != nil {
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

// ParseDKGCommitmentsUpdated is a log parse operation binding the contract event 0x1af28d6ae67078289128cdbaca9c65cda5fda4d71a49bb5e2d917ef3bf5f44dd.
//
// Solidity: event DKGCommitmentsUpdated(address indexed msgSender, uint32 round, uint32 total, uint32 threshold, uint32 index, bytes commitments, bytes signature, bytes mrenclave)
func (_DKG *DKGFilterer) ParseDKGCommitmentsUpdated(log types.Log) (*DKGDKGCommitmentsUpdated, error) {
	event := new(DKGDKGCommitmentsUpdated)
	if err := _DKG.contract.UnpackLog(event, "DKGCommitmentsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	MsgSender common.Address
	Round     uint32
	Index     uint32
	Finalized bool
	Mrenclave []byte
	Signature []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDKGFinalized is a free log retrieval operation binding the contract event 0xbeb92472b2143160d5e8c3fd25327aabddd7a0f4b80ade42661dd5b4e5ff4f90.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, uint32 index, bool finalized, bytes mrenclave, bytes signature)
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

// WatchDKGFinalized is a free log subscription operation binding the contract event 0xbeb92472b2143160d5e8c3fd25327aabddd7a0f4b80ade42661dd5b4e5ff4f90.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, uint32 index, bool finalized, bytes mrenclave, bytes signature)
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

// ParseDKGFinalized is a log parse operation binding the contract event 0xbeb92472b2143160d5e8c3fd25327aabddd7a0f4b80ade42661dd5b4e5ff4f90.
//
// Solidity: event DKGFinalized(address indexed msgSender, uint32 round, uint32 index, bool finalized, bytes mrenclave, bytes signature)
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
	MsgSender    common.Address
	Mrenclave    []byte
	Round        uint32
	Index        uint32
	DkgPubKey    []byte
	CommPubKey   []byte
	RemoteReport []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDKGInitialized is a free log retrieval operation binding the contract event 0x5c5aa077fdba8d39c55de6079c75de88ec2031742c08609237382f11414de616.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes mrenclave, uint32 round, uint32 index, bytes dkgPubKey, bytes commPubKey, bytes remoteReport)
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
// Solidity: event DKGInitialized(address indexed msgSender, bytes mrenclave, uint32 round, uint32 index, bytes dkgPubKey, bytes commPubKey, bytes remoteReport)
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
// Solidity: event DKGInitialized(address indexed msgSender, bytes mrenclave, uint32 round, uint32 index, bytes dkgPubKey, bytes commPubKey, bytes remoteReport)
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
