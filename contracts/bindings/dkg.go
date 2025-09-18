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
	PubKey       []byte
	RemoteReport []byte
	Commitments  []byte
	ChalStatus   uint8
	Finalized    bool
}

// DKGMetaData contains all meta data concerning the DKG contract.
var DKGMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"complainDeals\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dealComplaints\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"complainant\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dkgNodeInfos\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitments\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"finalized\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"finalized\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getNodeCount\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeInfo\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDKG.NodeInfo\",\"components\":[{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitments\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"finalized\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeDKG\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"remoteReport\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isActiveValidator\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nodeCount\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"nodeCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestRemoteAttestationOnChain\",\"inputs\":[{\"name\":\"targetIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitActiveValSet\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"valSet\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateDKGCommitments\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"total\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitments\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"valSets\",\"inputs\":[{\"name\":\"mrenclave\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DKGCommitmentsUpdated\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"total\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"threshold\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"commitments\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGFinalized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"finalized\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DKGInitialized\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"pubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"remoteReport\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealComplaintsSubmitted\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"complainIndexes\",\"type\":\"uint32[]\",\"indexed\":false,\"internalType\":\"uint32[]\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DealVerified\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"recipientIndex\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidDKGInitialization\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InvalidDeal\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RegistrationChallenged\",\"inputs\":[{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"challenger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoteAttestationProcessedOnChain\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"validator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"chalStatus\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIDKG.ChallengeStatus\"},{\"name\":\"round\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpgradeScheduled\",\"inputs\":[{\"name\":\"activationHeight\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"mrenclave\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
	Bin: "0x6080806040523461001657611d1b908161001b8239f35b5f80fdfe6080806040526004361015610012575f80fd5b60e05f35811c91826308e4842214611637575081630c5a40201461147f578163159726be146114085781632e266b4c14611289578163496c634c14610fcf5781634d71b63414610e3b5781634f1bf88d14610d975781636bad55fa1461065f57816378e510d2146105f55781637a57476514610325578163a27c2f6c146101ff57508063c2009b92146101495763fb5a783a146100ad575f80fd5b346101455760606003193601126101455760043567ffffffffffffffff8111610145576100de9036906004016118ee565b63ffffffff6100eb6118c8565b9160206100f6611a12565b948260405193849283378101600181520301902091165f5260205273ffffffffffffffffffffffffffffffffffffffff60405f2091165f52602052602060ff60405f2054166040519015158152f35b5f80fd5b346101455760806003193601126101455760043567ffffffffffffffff81116101455761017a9036906004016119f4565b6101826118c8565b61018a6118db565b6064359173ffffffffffffffffffffffffffffffffffffffff83168093036101455760405184818096516101c4816020998a80960161191c565b8101600281520301902063ffffffff8092165f52845260405f2091165f52825260405f20905f52815260ff60405f2054166040519015158152f35b346101455760606003193601126101455760043567ffffffffffffffff8111610145576102309036906004016119f4565b6102386118c8565b906102416118db565b6040518281809451610259816020978880960161191c565b81015f81520301902063ffffffff8094165f5282528260405f2091165f52815260ff61030660405f206102f88154966102ea61029760018501611a86565b6102a360028601611a86565b9260046102b260038801611a86565b9601549873ffffffffffffffffffffffffffffffffffffffff6040519c8d9c81168d52821c16908b01528060408b015289019061193d565b90878203606089015261193d565b90858203608087015261193d565b9161031660a08501838316611962565b60081c16151560c08301520390f35b346101455760806003193601126101455761033e6118b5565b6103466118c8565b60443567ffffffffffffffff8082116101455736602383011215610145578160040135936024918086116105c8576005908660051b906020956040519861038f8885018b61198b565b8952868901602481948301019136831161014557602401905b8282106105ac57505050606435908111610145576103ca9036906004016118ee565b906040518282823787818481015f8152030190209263ffffffff80961693845f5288528560405f20991698895f52885273ffffffffffffffffffffffffffffffffffffffff60405f2054891c16330361054e575f5b8a518110156104c7576040518484823789818681016002815203019020855f52895260405f208b5182101561049b5790600191888b8e848b1b010151165f528a5260405f20335f528a5260405f20827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008254161790550161041f565b887f4e487b71000000000000000000000000000000000000000000000000000000005f5260326004525ffd5b509350889285898960405196608088019288526080828901525180925260a0870197925f905b838210610535577f1c2112af5fd37661e3dd248d701decebf291f13de6a411176337ef21a7b1a63089806105308d8c8c8c60408601528483036060860152611b47565b0390a1005b845181168a5298820198938201936001909101906104ed565b606488604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601360248201527f496e76616c696420636f6d706c61696e616e74000000000000000000000000006044820152fd5b813563ffffffff811681036101455781529088019088016103a8565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b346101455760406003193601126101455760043567ffffffffffffffff81116101455761062860209136906004016118ee565b826106316118c8565b928260405193849283378101600381520301902063ffffffff8092165f52825260405f205416604051908152f35b34610145576080600319360112610145576106786118b5565b60243567ffffffffffffffff8111610145576106989036906004016118ee565b9060443567ffffffffffffffff8111610145576106b99036906004016118ee565b929060643567ffffffffffffffff8111610145576106db9036906004016118ee565b909560405184868237602081868101600181520301902063ffffffff82165f5260205260405f20335f5260205260ff60405f2054168015610d8f575b15610d315760405184868237602081868101600381520301902063ffffffff82165f5260205263ffffffff60405f2054169060405185878237602081878101600381520301902063ffffffff82165f5260205260405f20805463ffffffff80821614610d04577fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000063ffffffff600181841601169116179055604051966107bc8861196f565b8288523360208901526107d03682876119ae565b60408901526107e036858b6119ae565b606089015260405180602081011067ffffffffffffffff6020830111176105c857602081016040525f815260808901525f60a08901525f60c0890152604051868882376020818881015f81520301902063ffffffff83165f5260205260405f20835f5260205260405f2063ffffffff8951168154907fffffffffffffffff00000000000000000000000000000000000000000000000077ffffffffffffffffffffffffffffffffffffffff0000000060208d015160201b16921617178155604089015180519067ffffffffffffffff82116105c8576108cf826108c66001860154611a35565b60018601611c96565b602090601f8311600114610c775761091c92915f9183610bdf575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60018201555b606089015180519067ffffffffffffffff82116105c8576109538261094a6002860154611a35565b60028601611c96565b602090601f8311600114610bea5761099f92915f9183610bdf5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60028201555b608089015180519067ffffffffffffffff82116105c8576109d6826109cd6003860154611a35565b60038601611c96565b602090601f8311600114610b51579180610a2692600495945f92610b465750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b60038201555b019460a0890151936003851015610b1957610acf8760c07f365f737a37008e835fc0efda066a512c68f4d81401aa358182aa7999b9f15ec39c610b049860ff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00610b149d54169116178355015115157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff00835492151560081b169116179055565b63ffffffff610aeb6040519a8b9a60a08c5260a08c0191611b47565b9416602089015260408801528683036060880152611b47565b9083820360808501523396611b47565b0390a2005b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b015190508e806108ea565b90600384015f5260205f20915f5b601f1985168110610bc7575091839160019383601f19600498971610610b90575b505050811b016003820155610a2c565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558d8080610b80565b91926020600181928685015181550194019201610b5f565b015190508d806108ea565b9190600284015f5260205f20905f935b601f1984168510610c5c576001945083601f19811610610c25575b505050811b0160028201556109a5565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558c8080610c15565b81810151835560209485019460019093019290910190610bfa565b9190600184015f5260205f20905f935b601f1984168510610ce9576001945083601f19811610610cb2575b505050811b016001820155610922565b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c191690558c8080610ca2565b81810151835560209485019460019093019290910190610c87565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f56616c696461746f72206e6f7420696e206163746976652073657400000000006044820152fd5b506001610717565b346101455760606003193601126101455760043567ffffffffffffffff811161014557610dc89036906004016119f4565b610dd06118c8565b63ffffffff610df66020610de2611a12565b94816040519382858094519384920161191c565b8101600181520301902091165f5260205273ffffffffffffffffffffffffffffffffffffffff60405f2091165f52602052602060ff60405f2054166040519015158152f35b346101455760a060031936011261014557610e546118b5565b610e5c6118c8565b60443590811515928383036101455767ffffffffffffffff9360643585811161014557610e8d9036906004016118ee565b909560843590811161014557610ea79036906004016118ee565b919096604051828282376020818481015f8152030190209463ffffffff80911695865f5260205260405f20961695865f5260205260405f2091610f0573ffffffffffffffffffffffffffffffffffffffff845460201c163314611b67565b600483019460ff8654166003811015610b19577fbeb92472b2143160d5e8c3fd25327aabddd7a0f4b80ade42661dd5b4e5ff4f9099610b1497610f5c6001610b0498610f5682610fa3971415611bcc565b01611a86565b51151580610fc6575b610f6e90611c31565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff61ff00835492151560081b169116179055565b60405197889788526020880152604087015260a0606087015260a0860191611b47565b50871515610f65565b34610145578060031936011261014557610fe76118b5565b610fef6118c8565b90610ff86118db565b926064359363ffffffff908186168096036101455767ffffffffffffffff916084358381116101455761102f9036906004016118ee565b97909460a4358581116101455761104a9036906004016118ee565b92909360c43595878711610145576110678c9736906004016118ee565b979098838b83604051948592833781015f81526020948591030190209c169b8c5f52825260405f20855f52825260405f206110bc73ffffffffffffffffffffffffffffffffffffffff8254851c163314611b67565b60ff600482015416906003821015610b19576003906110df600180941415611bcc565b6110eb60018201611a86565b51151580611280575b80611277575b61110390611c31565b019188116105c85787906111218261111b8554611a35565b85611c96565b5f90601f83116001146111f15750908061116d925f916111e6575b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c19161790565b90555b826040519c8d9c8d5216908b0152166040890152606088015280608088015286019061119b92611b47565b9084820360a08601526111ad92611b47565b82810360c084015233946111c092611b47565b037f1af28d6ae67078289128cdbaca9c65cda5fda4d71a49bb5e2d917ef3bf5f44dd91a2005b90508a01355f61113c565b91601f198116845f528b865f2094875f925b84841061125b575050505010611223575b5050600187811b019055611170565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88a60031b161c19908a01351690558e80611214565b860135875595810195948501948d94509190910190878e611203565b508b15156110fa565b508915156110f4565b34610145576060600319360112610145576112a26118b5565b67ffffffffffffffff60248035828111610145576112c49036906004016118ee565b9290916044358281116101455736602382011215610145578060040135928311610145576005923660248260051b8401011161014557945f9063ffffffff809816915b87811061131057005b6040519082888337818381015f8152602093849103019020845f52825260405f208a82165f52825260ff600460405f2001541660038110156113dc576001809103611361575b506001915001611307565b604051848a8237838186810184815203019020855f52835260405f2087838a1b880101359373ffffffffffffffffffffffffffffffffffffffff8516809503610145576001945f525260405f20907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008254161790558a611356565b867f4e487b71000000000000000000000000000000000000000000000000000000005f5260216004525ffd5b346101455760406003193601126101455760043567ffffffffffffffff81116101455761143b60209136906004016119f4565b61145b826114476118c8565b92816040519382858094519384920161191c565b8101600381520301902063ffffffff8092165f52825260405f205416604051908152f35b3461014557606080600319360112610145576004359067ffffffffffffffff8211610145576114b460c09236906004016118ee565b93906114be6118c8565b946114c76118db565b91604051906114d58261196f565b5f82525f60209783828a8096015288604082015288808201528860808201528260a0820152015282604051938492833781015f81520301902063ffffffff8096165f5284528460405f2091165f52835260405f2091604051946115378661196f565b835493818516875273ffffffffffffffffffffffffffffffffffffffff928387890196881c16865261156b60018301611a86565b926040890193845261157f60028401611a86565b93828a01948552600461159460038601611a86565b9460808c0195865201549560ff87169860a08c01996003811015610b195789948d8d928d5260c0019960081c60ff1615158a526040519d8e9d8e525116908c0152511660408a0152519188015261010087016115ef9161193d565b9151601f1992838882030160808901526116089161193d565b9051918682030160a087015261161d9161193d565b925160c0850161162c91611962565b511515908301520390f35b34610145576060600319360112610145576116506118b5565b906116596118c8565b9160443567ffffffffffffffff81116101455761167a9036906004016118ee565b909381858537838281015f81526020958691030190209063ffffffff80911691825f52855260405f20931692835f52845260405f2073ffffffffffffffffffffffffffffffffffffffff808254871c1615918261185757600481019283549060ff82166003811015610b19576117f9577f5625cad747d0f660a99f4fe677a17d22db2a71e28505cd2148f61ad3cefc4f1d9995938361053098969360ff9361172760026117999801611a86565b61173360018501611a86565b9051151591826117f0575b50816117e6575b816117db575b50156117af575060027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008554161784555b548a1c169154169060405198899889528801526040870190611962565b606085015260a0608085015260a0840191611b47565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117845561177c565b90505115158e61174b565b8815159150611745565b1591508f61173e565b606489604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601760248201527f4e6f646520616c7265616479206368616c6c656e6765640000000000000000006044820152fd5b606487604051907f08c379a00000000000000000000000000000000000000000000000000000000082526004820152601360248201527f4e6f646520646f6573206e6f74206578697374000000000000000000000000006044820152fd5b6004359063ffffffff8216820361014557565b6024359063ffffffff8216820361014557565b6044359063ffffffff8216820361014557565b9181601f840112156101455782359167ffffffffffffffff8311610145576020838186019501011161014557565b5f5b83811061192d5750505f910152565b818101518382015260200161191e565b90601f19601f60209361195b8151809281875287808801910161191c565b0116010190565b906003821015610b195752565b60e0810190811067ffffffffffffffff8211176105c857604052565b90601f601f19910116810190811067ffffffffffffffff8211176105c857604052565b92919267ffffffffffffffff82116105c857604051916119d86020601f19601f840116018461198b565b829481845281830111610145578281602093845f960137010152565b9080601f8301121561014557816020611a0f933591016119ae565b90565b6044359073ffffffffffffffffffffffffffffffffffffffff8216820361014557565b90600182811c92168015611a7c575b6020831014611a4f57565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b91607f1691611a44565b9060405191825f8254611a9881611a35565b908184526020946001916001811690815f14611b065750600114611ac8575b505050611ac69250038361198b565b565b5f90815285812095935091905b818310611aee575050611ac693508201015f8080611ab7565b85548884018501529485019487945091830191611ad5565b915050611ac69593507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201015f8080611ab7565b601f8260209493601f1993818652868601375f8582860101520116010190565b15611b6e57565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f496e76616c69642076616c696461746f720000000000000000000000000000006044820152fd5b15611bd357565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f4e6f64652077617320696e76616c6964617465640000000000000000000000006044820152fd5b15611c3857565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f496e76616c6964207369676e61747572650000000000000000000000000000006044820152fd5b601f8211611ca357505050565b5f5260205f20906020601f840160051c83019310611cdb575b601f0160051c01905b818110611cd0575050565b5f8155600101611cc5565b9091508190611cbc56fea26469706673582212204b6e9c8acf9e3a61ad65cb99c34d3144316b76e2cdc625351f9068eaff3c43ea64736f6c63430008170033",
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
// Solidity: function dkgNodeInfos(bytes mrenclave, uint32 round, uint32 index) view returns(uint32 index, address validator, bytes pubKey, bytes remoteReport, bytes commitments, uint8 chalStatus, bool finalized)
func (_DKG *DKGCaller) DkgNodeInfos(opts *bind.CallOpts, mrenclave []byte, round uint32, index uint32) (struct {
	Index        uint32
	Validator    common.Address
	PubKey       []byte
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
		PubKey       []byte
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
	outstruct.PubKey = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.RemoteReport = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.Commitments = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.ChalStatus = *abi.ConvertType(out[5], new(uint8)).(*uint8)
	outstruct.Finalized = *abi.ConvertType(out[6], new(bool)).(*bool)

	return *outstruct, err

}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xa27c2f6c.
//
// Solidity: function dkgNodeInfos(bytes mrenclave, uint32 round, uint32 index) view returns(uint32 index, address validator, bytes pubKey, bytes remoteReport, bytes commitments, uint8 chalStatus, bool finalized)
func (_DKG *DKGSession) DkgNodeInfos(mrenclave []byte, round uint32, index uint32) (struct {
	Index        uint32
	Validator    common.Address
	PubKey       []byte
	RemoteReport []byte
	Commitments  []byte
	ChalStatus   uint8
	Finalized    bool
}, error) {
	return _DKG.Contract.DkgNodeInfos(&_DKG.CallOpts, mrenclave, round, index)
}

// DkgNodeInfos is a free data retrieval call binding the contract method 0xa27c2f6c.
//
// Solidity: function dkgNodeInfos(bytes mrenclave, uint32 round, uint32 index) view returns(uint32 index, address validator, bytes pubKey, bytes remoteReport, bytes commitments, uint8 chalStatus, bool finalized)
func (_DKG *DKGCallerSession) DkgNodeInfos(mrenclave []byte, round uint32, index uint32) (struct {
	Index        uint32
	Validator    common.Address
	PubKey       []byte
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

// InitializeDKG is a paid mutator transaction binding the contract method 0x6bad55fa.
//
// Solidity: function initializeDKG(uint32 round, bytes mrenclave, bytes pubKey, bytes remoteReport) returns()
func (_DKG *DKGTransactor) InitializeDKG(opts *bind.TransactOpts, round uint32, mrenclave []byte, pubKey []byte, remoteReport []byte) (*types.Transaction, error) {
	return _DKG.contract.Transact(opts, "initializeDKG", round, mrenclave, pubKey, remoteReport)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0x6bad55fa.
//
// Solidity: function initializeDKG(uint32 round, bytes mrenclave, bytes pubKey, bytes remoteReport) returns()
func (_DKG *DKGSession) InitializeDKG(round uint32, mrenclave []byte, pubKey []byte, remoteReport []byte) (*types.Transaction, error) {
	return _DKG.Contract.InitializeDKG(&_DKG.TransactOpts, round, mrenclave, pubKey, remoteReport)
}

// InitializeDKG is a paid mutator transaction binding the contract method 0x6bad55fa.
//
// Solidity: function initializeDKG(uint32 round, bytes mrenclave, bytes pubKey, bytes remoteReport) returns()
func (_DKG *DKGTransactorSession) InitializeDKG(round uint32, mrenclave []byte, pubKey []byte, remoteReport []byte) (*types.Transaction, error) {
	return _DKG.Contract.InitializeDKG(&_DKG.TransactOpts, round, mrenclave, pubKey, remoteReport)
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
	PubKey       []byte
	RemoteReport []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDKGInitialized is a free log retrieval operation binding the contract event 0x365f737a37008e835fc0efda066a512c68f4d81401aa358182aa7999b9f15ec3.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes mrenclave, uint32 round, uint32 index, bytes pubKey, bytes remoteReport)
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

// WatchDKGInitialized is a free log subscription operation binding the contract event 0x365f737a37008e835fc0efda066a512c68f4d81401aa358182aa7999b9f15ec3.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes mrenclave, uint32 round, uint32 index, bytes pubKey, bytes remoteReport)
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

// ParseDKGInitialized is a log parse operation binding the contract event 0x365f737a37008e835fc0efda066a512c68f4d81401aa358182aa7999b9f15ec3.
//
// Solidity: event DKGInitialized(address indexed msgSender, bytes mrenclave, uint32 round, uint32 index, bytes pubKey, bytes remoteReport)
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
