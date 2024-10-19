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

// IIPTokenStakingInitializerArgs is an auto generated low-level Go binding around an user-defined struct.
type IIPTokenStakingInitializerArgs struct {
	Owner             common.Address
	MinStakeAmount    *big.Int
	MinUnstakeAmount  *big.Int
	MinCommissionRate *big.Int
	Fee               *big.Int
}

// IPTokenStakingMetaData contains all meta data concerning the IPTokenStaking contract.
var IPTokenStakingMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"stakingRounding\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"defaultMinFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_MIN_FEE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"STAKE_ROUNDING\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addOperator\",\"inputs\":[{\"name\":\"uncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"createValidator\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"moniker\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionChangeRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"supportsUnlocked\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"createValidatorOnBehalf\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"moniker\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionChangeRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"supportsUnlocked\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"args\",\"type\":\"tuple\",\"internalType\":\"structIIPTokenStaking.InitializerArgs\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minUnstakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minCommissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minCommissionRate\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minStakeAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minUnstakeAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"redelegate\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpSrcPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpDstPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"redelegateOnBehalf\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpSrcPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpDstPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"removeOperator\",\"inputs\":[{\"name\":\"uncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"roundedStakeAmount\",\"inputs\":[{\"name\":\"rawAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainder\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setFee\",\"inputs\":[{\"name\":\"newFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinCommissionRate\",\"inputs\":[{\"name\":\"newValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinStakeAmount\",\"inputs\":[{\"name\":\"newMinStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinUnstakeAmount\",\"inputs\":[{\"name\":\"newMinUnstakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRewardsAddress\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"newRewardsAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setWithdrawalAddress\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"newWithdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"stakingPeriod\",\"type\":\"uint8\",\"internalType\":\"enumIIPTokenStaking.StakingPeriod\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeOnBehalf\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"stakingPeriod\",\"type\":\"uint8\",\"internalType\":\"enumIIPTokenStaking.StakingPeriod\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unjail\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unjailOnBehalf\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unstake\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstakeOnBehalf\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateValidatorCommission\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AddOperator\",\"inputs\":[{\"name\":\"uncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CreateValidator\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"moniker\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"stakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"maxCommissionRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"maxCommissionChangeRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"supportsUnlocked\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorUnCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"stakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"stakingPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeSet\",\"inputs\":[{\"name\":\"newFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinCommissionRateChanged\",\"inputs\":[{\"name\":\"minCommissionRate\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinStakeAmountSet\",\"inputs\":[{\"name\":\"minStakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinUnstakeAmountSet\",\"inputs\":[{\"name\":\"minUnstakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Redelegate\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpSrcPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpDstPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoveOperator\",\"inputs\":[{\"name\":\"uncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetRewardAddress\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"executionAddress\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetWithdrawalAddress\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"executionAddress\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unjail\",\"inputs\":[{\"name\":\"unjailer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdateValidatorCommssion\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorUnCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"stakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"IPTokenStaking__CommissionRateOverMax\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__CommissionRateUnderMin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__FailedRemainerRefund\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidDefaultMinFee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidDelegationId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidFeeAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidMinFee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__LowUnstakeAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__RedelegatingToSameValidator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__StakeAmountUnderMin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroMinCommissionRate\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroMinStakeAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroMinUnstakeAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroStakingRounding\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PubKeyVerifier__InvalidPubkeyDerivedAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubKeyVerifier__InvalidPubkeyLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PubKeyVerifier__InvalidPubkeyPrefix\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60c0346200016c5762002930906001600160401b0390601f38849003908101601f1916820190838211838310176200017157808391604096879485528339810103126200016c57602081519101519080156200015b57608052633b9aca0081106200014a5760a0527ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff82851c1662000139578080831603620000f4575b83516127a890816200018882396080518181816105f60152818161087801528181611c7a01528181611ec2015281816121ae015281816124b6015261265f015260a051818181610c9901526124040152f35b6001600160401b0319909116811790915581519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a1388080620000a2565b835163f92ee8a960e01b8152600490fd5b8251630f8cc23360e21b8152600490fd5b835163591eebd160e11b8152600490fd5b600080fd5b634e487b7160e01b600052604160045260246000fdfe6040608081526004908136101561001557600080fd5b600091823560e01c8063014e817814611694578063057b9296146115095780631487153e146114ec57806317e42e12146113ad57806339ec4df91461138e5780633dd9fb9a1461124657806369fe0e2d146112215780636ea3a228146111fc578063715018a614611135578063787f82c81461104957806379ba509714610fbf57806386eb5e4814610f975780638740597a14610e4e5780638da5cb5b14610dfa5780638ed65fbc14610cbc57806394fd0fe014610c815780639d04b12114610b045780639d9d293f14610a60578063a0284f1614610a08578063ab8870f6146109e3578063b2bc29ef1461089b578063bda16b1514610860578063c582db4414610632578063d2e1f5b8146105d9578063ddca3f43146105bc578063e30c397814610568578063eb4af04514610543578063ec21dac214610510578063f1887684146104f1578063f2fde38b1461041f578063f9550a8d146103c75763fce5dc8c1461018157600080fd5b346103c35760a06003193601126103c3577ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff82851c16159167ffffffffffffffff8116801590816103bb575b60011490816103b1575b1590816103a8575b50610380578260017fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000831617855561034b575b50610221612719565b610229612719565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580359073ffffffffffffffffffffffffffffffffffffffff821680830361034757610276612719565b61027e612719565b15610318575061028d90612518565b61029860243561262d565b6102a3604435612484565b6102ae6064356125cc565b6102b9608435612402565b6102c1578280f35b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d291817fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff602093541690555160018152a138808280f35b602490868651917f1e4fbdf7000000000000000000000000000000000000000000000000000000008352820152fd5b8680fd5b7fffffffffffffffffffffffffffffffffffffffffffffff000000000000000000166801000000000000000117835538610218565b5083517ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b905015386101e5565b303b1591506101dd565b8491506101d3565b8280fd5b836103f86103d436611866565b986103eb89829a939a999499989598979697612001565b6103f3612138565b611c5e565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580f35b8382346104ed5760206003193601126104ed573573ffffffffffffffffffffffffffffffffffffffff8082168092036103c35761045a612392565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227008380a380f35b5080fd5b5050346104ed57816003193601126104ed576020906001549051908152f35b8361054061051d36611903565b966105318782989398979497969596612001565b61053b8484612001565b611e80565b80f35b8382346104ed5760206003193601126104ed5761054090610562612392565b3561262d565b5050346104ed57816003193601126104ed5760209073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0054169051908152f35b50346103c357826003193601126103c35760209250549051908152f35b50823461062f57602060031936011261062f57503561062261061b7f000000000000000000000000000000000000000000000000000000000000000083611f8c565b8092611fc5565b9082519182526020820152f35b80fd5b509190806003193601126104ed57823567ffffffffffffffff81116103c35761065e90369085016116bf565b9190936024359063ffffffff821680920361085c5760418403610835578315610809577f04000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000873516036107e2573373ffffffffffffffffffffffffffffffffffffffff6106eb86896120de565b16036107bb5780543403610794578480341561078b575b81808092813491f115610781578454821061075a575061074e7f202c9aad6965f28c0ce1cd00460c1adfa2c90277f4f0a7abb813e2f04cecd70b94958351948486958652850191611a59565b9060208301520390a180f35b82517f183785b6000000000000000000000000000000000000000000000000000000008152fd5b82513d86823e3d90fd5b506108fc610702565b82517f5097ac51000000000000000000000000000000000000000000000000000000008152fd5b82517ff68f8e95000000000000000000000000000000000000000000000000000000008152fd5b82517f759898d3000000000000000000000000000000000000000000000000000000008152fd5b8460326024927f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b82517f23679d9f000000000000000000000000000000000000000000000000000000008152fd5b8480fd5b5050346104ed57816003193601126104ed57602090517f00000000000000000000000000000000000000000000000000000000000000008152f35b509190346104ed576108ac366116f2565b976041879397989298969496036109bb57821561098f577f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000083351603610967573373ffffffffffffffffffffffffffffffffffffffff61092d85856120de565b160361093f5750610540989950611971565b8a90517ff68f8e95000000000000000000000000000000000000000000000000000000008152fd5b8a90517f759898d3000000000000000000000000000000000000000000000000000000008152fd5b60248a60328d7f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b8a90517f23679d9f000000000000000000000000000000000000000000000000000000008152fd5b8382346104ed5760206003193601126104ed5761054090610a02612392565b356125cc565b602083610a35610a17366117b0565b95610a288682979397969496612001565b610a30612138565b612192565b9060017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005551908152f35b509190610a6c36611903565b976041879397989298969496036109bb57821561098f577f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000083351603610967573373ffffffffffffffffffffffffffffffffffffffff610aed85856120de565b160361093f575061054098995061053b8484612001565b5090610b0f36611760565b9092919360418403610835578315610c55577f04000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000863516036107e25773ffffffffffffffffffffffffffffffffffffffff903382610b8d87896120de565b1603610c2e5780543403610c07575085803415610bfe575b81808092813491f115610bf457610be77f28c0529db8cf660d5b4c1e4b9313683fa7241c3fc49452e7d0ebae215a5f84b2958451958587968752860191611a59565b911660208301520390a180f35b82513d87823e3d90fd5b506108fc610ba5565b83517f5097ac51000000000000000000000000000000000000000000000000000000008152fd5b83517ff68f8e95000000000000000000000000000000000000000000000000000000008152fd5b8560326024927f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b5050346104ed57816003193601126104ed57602090517f00000000000000000000000000000000000000000000000000000000000000008152f35b509190610cc83661181f565b9360418394929403610dd2578215610da6577f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000083351603610d7e573373ffffffffffffffffffffffffffffffffffffffff610d4485856120de565b1603610d565750610540949550611acd565b8690517ff68f8e95000000000000000000000000000000000000000000000000000000008152fd5b8690517f759898d3000000000000000000000000000000000000000000000000000000008152fd5b6024866032897f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b8690517f23679d9f000000000000000000000000000000000000000000000000000000008152fd5b5050346104ed57816003193601126104ed5760209073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054169051908152f35b509190610e5a36611866565b9960418993999a929a98949897959703610f6f578215610f43577f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000083351603610f1b573373ffffffffffffffffffffffffffffffffffffffff610ede85856120de565b1603610ef357506103f89a9b506103f3612138565b8c90517ff68f8e95000000000000000000000000000000000000000000000000000000008152fd5b8c90517f759898d3000000000000000000000000000000000000000000000000000000008152fd5b60248c60328f7f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b8c90517f23679d9f000000000000000000000000000000000000000000000000000000008152fd5b836103f8610fa43661181f565b92610fb0929192612138565b610fba8282612001565b611acd565b5090346103c357826003193601126103c3573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00541603611019578261054033612518565b6024925051907f118cdaa70000000000000000000000000000000000000000000000000000000082523390820152fd5b509061105436611760565b9092919360418403610835578315610c55577f04000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000863516036107e25773ffffffffffffffffffffffffffffffffffffffff9033826110d287896120de565b1603610c2e5780543403610c0757508580341561112c575b81808092813491f115610bf457610be77f9f7f04f688298f474ed4c786abb29e0ca0173d70516d55d9eac515609b45fbca958451958587968752860191611a59565b506108fc6110ea565b833461062f578060031936011261062f5761114e612392565b8073ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b8382346104ed5760206003193601126104ed576105409061121b612392565b35612484565b8382346104ed5760206003193601126104ed5761054090611240612392565b35612402565b5091611251366117b0565b98909660418694969793970361136657831561133b57507f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000083351603611314573373ffffffffffffffffffffffffffffffffffffffff6112d285856120de565b16036112ed575091610a359593916020989593610a30612138565b87517ff68f8e95000000000000000000000000000000000000000000000000000000008152fd5b87517f759898d3000000000000000000000000000000000000000000000000000000008152fd5b60326024927f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b5087517f23679d9f000000000000000000000000000000000000000000000000000000008152fd5b5050346104ed57816003193601126104ed576020906002549051908152f35b508290346104ed576113be36611760565b949092604184036114c5578315610809577f04000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000008435160361149e573373ffffffffffffffffffffffffffffffffffffffff61143986866120de565b160361147757507f65729f64aec4981a7e5cedc9abbed98ce4ee8a5c6ecefc35e32d646d517180429394611471915193849384611a98565b0390a180f35b90517ff68f8e95000000000000000000000000000000000000000000000000000000008152fd5b90517f759898d3000000000000000000000000000000000000000000000000000000008152fd5b90517f23679d9f000000000000000000000000000000000000000000000000000000008152fd5b5050346104ed57816003193601126104ed57602091549051908152f35b5061151336611760565b90926041840361166d578315610c55577f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000084351603611646573373ffffffffffffffffffffffffffffffffffffffff61158d86866120de565b160361161f57805434036115f85750848034156115ef575b81808092813491f1156115e257611471907f6ac365cf05479bb8a295fbf9637875411d6d6f2a0ac7c4b1f560cedcf1a33081945193849384611a98565b50505051903d90823e3d90fd5b506108fc6115a5565b84517f5097ac51000000000000000000000000000000000000000000000000000000008152fd5b84517ff68f8e95000000000000000000000000000000000000000000000000000000008152fd5b84517f759898d3000000000000000000000000000000000000000000000000000000008152fd5b84517f23679d9f000000000000000000000000000000000000000000000000000000008152fd5b833461062f576105406116a6366116f2565b966116ba8782989398979497969596612001565b611971565b9181601f840112156116ed5782359167ffffffffffffffff83116116ed57602083818601950101116116ed57565b600080fd5b60a06003198201126116ed5767ffffffffffffffff906004358281116116ed578161171f916004016116bf565b939093926024358181116116ed578361173a916004016116bf565b9390939260443592606435926084359182116116ed5761175c916004016116bf565b9091565b60406003198201126116ed576004359067ffffffffffffffff82116116ed5761178b916004016116bf565b909160243573ffffffffffffffffffffffffffffffffffffffff811681036116ed5790565b60806003198201126116ed5767ffffffffffffffff916004358381116116ed57826117dd916004016116bf565b939093926024358281116116ed57816117f8916004016116bf565b9390939260443560048110156116ed57926064359182116116ed5761175c916004016116bf565b60406003198201126116ed5767ffffffffffffffff916004358381116116ed578261184c916004016116bf565b939093926024359182116116ed5761175c916004016116bf565b9060e06003198301126116ed5767ffffffffffffffff916004358381116116ed5781611894916004016116bf565b939093926024358281116116ed57836118af916004016116bf565b9093909263ffffffff9160443583811681036116ed579360643584811681036116ed579360843590811681036116ed579260a43580151581036116ed579260c4359182116116ed5761175c916004016116bf565b9060a06003198301126116ed5767ffffffffffffffff6004358181116116ed5783611930916004016116bf565b939093926024358381116116ed578261194b916004016116bf565b939093926044359182116116ed57611965916004016116bf565b90916064359060843590565b959094929691936119828886612001565b6003548111611a2f576002548410611a05576119e2611a00957fac41e6ee15d2d0047feb1ea8aba74b92c0334cd3e78024a5ad679d7d08b8fbc5996119d46040519a8b9a60c08c5260c08c0191611a59565b9189830360208b0152611a59565b936040870152606086015233608086015284830360a0860152611a59565b0390a1565b60046040517f23870ab9000000000000000000000000000000000000000000000000000000008152fd5b60046040517fc7617a88000000000000000000000000000000000000000000000000000000008152fd5b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b91611ac660209273ffffffffffffffffffffffffffffffffffffffff92969596604086526040860191611a59565b9416910152565b91926004543403611b565760003415611b4d575b600080808093813491f115611b41577f026c2e156478ec2a25ccebac97a338d301f69b6d5aeec39c578b28a95e11820193611a0091611b33604051958695338752606060208801526060870191611a59565b918483036040860152611a59565b6040513d6000823e3d90fd5b506108fc611ae1565b60046040517f5097ac51000000000000000000000000000000000000000000000000000000008152fd5b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f604051930116820182811067ffffffffffffffff821117611bc457604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff8111611bc457601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b929192611c41611c3c83611bf3565b611b80565b93828552828201116116ed57816000926020928387013784010152565b9291611c74919798969594929a999a3691611c2d565b93611c9f7f000000000000000000000000000000000000000000000000000000000000000034611f8c565b97611caa8934611fc5565b946001548610611e565760009687549163ffffffff809616928310611e2c57851692838311611e02578880898015611df8575b82809291818093f115611ded5715611de357611d096001965b6040519b8c6101208091528d0191611a59565b906020988b83038a8d0152815191828452815b838110611dd0575050937f65bfc2fa1cd4c6f50f60983ad1cf1cb4bff5ee6570428254dfce41b085ef6d149c9d9e9793837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8f9e9c999560ff99611dba9f82819e9a0101520116019660408d015260608c015260808b01521660a08901521660c08701523360e087015281868203016101008701520191611a59565b0390a180611dc55750565b611dce906126bb565b565b8181018c01518582018d01528b01611d1c565b611d098896611cf6565b6040513d8a823e3d90fd5b6108fc9150611cdd565b60046040517f809afa64000000000000000000000000000000000000000000000000000000008152fd5b60046040517f183785b6000000000000000000000000000000000000000000000000000000008152fd5b60046040517fda15b66c000000000000000000000000000000000000000000000000000000008152fd5b95939094969296611e918183612001565b611e9c368985611c2d565b60208151910120611eae368385611c2d565b6020815191012014611f6257611eed611ee77f000000000000000000000000000000000000000000000000000000000000000034611f8c565b34611fc5565b60015411611e56576003548511611a2f57611f3e611f4c937f210091050fbe3add6ade45436b6c7aed210ef28fc37e1a1775970fc391272fe8996119d46040519a8b9a60c08c5260c08c0191611a59565b918683036040880152611a59565b91606084015233608084015260a08301520390a1565b60046040517f43df0a36000000000000000000000000000000000000000000000000000000008152fd5b8115611f96570690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b91908203918211611fd257565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b90604181036120b45715612085577fff000000000000000000000000000000000000000000000000000000000000007f04000000000000000000000000000000000000000000000000000000000000009135160361205b57565b60046040517f759898d3000000000000000000000000000000000000000000000000000000008152fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60046040517f23679d9f000000000000000000000000000000000000000000000000000000008152fd5b816001116116ed5773ffffffffffffffffffffffffffffffffffffffff9161212d9160017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff3693019101611c2d565b602081519101201690565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0060028154146121685760029055565b60046040517f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b929593909193600482101561236357600382116122df576121d37f000000000000000000000000000000000000000000000000000000000000000034611f8c565b956121de8734611fc5565b956001548710611e565760009884612298575b946122676000989495899893967f269a32ff589c9b701f49ab6aa532ee8f55901df71a7fca2d70dc9f45314f1be39560ff6122418c9b9a8c9b6119d46040519a8b9a60e08c5260e08c0191611a59565b938960408801521660608601528d60808601523360a086015284830360c0860152611a59565b0390a181811561228f575b8290f115611b415780612283575090565b61228c906126bb565b90565b506108fc612272565b91949850929591946003547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114611fd257600101806003559894919693909592966121f1565b60846040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4950546f6b656e5374616b696e673a20496e76616c6964207374616b696e672060448201527f706572696f6400000000000000000000000000000000000000000000000000006064820152fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300541633036123d257565b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b7f0000000000000000000000000000000000000000000000000000000000000000811061245a576020817f20461e09b8e557b77e107939f9ce6544698123aad0fc964ac5cc59b7df2e608f92600455604051908152a1565b60046040517f7840bc30000000000000000000000000000000000000000000000000000000008152fd5b80156124ee5760206124e17ff93d77980ae5a1ddd008d6a7f02cbee5af2a4fcea850c4b55828de4f644e589f926124db7f000000000000000000000000000000000000000000000000000000000000000082611f8c565b90611fc5565b80600255604051908152a1565b60046040517f8d04d544000000000000000000000000000000000000000000000000000000008152fd5b7fffffffffffffffffffffffff0000000000000000000000000000000000000000907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008281541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549073ffffffffffffffffffffffffffffffffffffffff80931680948316179055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3565b8015612603576020817f4167b1de65292a9ff628c9136823791a1de701e1fbdda4863ce22a1cfaf4d0f792600055604051908152a1565b60046040517f23cf9ec0000000000000000000000000000000000000000000000000000000008152fd5b80156126915760206126847fea095c2fea861b87f0fd54d0d4453358692a527e120df22b62c71696247dfb9f926124db7f000000000000000000000000000000000000000000000000000000000000000082611f8c565b80600155604051908152a1565b60046040517ff4d335c6000000000000000000000000000000000000000000000000000000008152fd5b600080808093335af13d15612714573d6126d7611c3c82611bf3565b908152600060203d92013e5b156126ea57565b60046040517ffc0ea4f4000000000000000000000000000000000000000000000000000000008152fd5b6126e3565b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c161561274857565b60046040517fd7e6bcf8000000000000000000000000000000000000000000000000000000008152fdfea26469706673582212204fb961a7fa17dabd05d4bd81fa8217c1a6372f7a6eb89bc59468308ae2458ff264736f6c63430008170033",
}

// IPTokenStakingABI is the input ABI used to generate the binding from.
// Deprecated: Use IPTokenStakingMetaData.ABI instead.
var IPTokenStakingABI = IPTokenStakingMetaData.ABI

// IPTokenStakingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use IPTokenStakingMetaData.Bin instead.
var IPTokenStakingBin = IPTokenStakingMetaData.Bin

// DeployIPTokenStaking deploys a new Ethereum contract, binding an instance of IPTokenStaking to it.
func DeployIPTokenStaking(auth *bind.TransactOpts, backend bind.ContractBackend, stakingRounding *big.Int, defaultMinFee *big.Int) (common.Address, *types.Transaction, *IPTokenStaking, error) {
	parsed, err := IPTokenStakingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(IPTokenStakingBin), backend, stakingRounding, defaultMinFee)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IPTokenStaking{IPTokenStakingCaller: IPTokenStakingCaller{contract: contract}, IPTokenStakingTransactor: IPTokenStakingTransactor{contract: contract}, IPTokenStakingFilterer: IPTokenStakingFilterer{contract: contract}}, nil
}

// IPTokenStaking is an auto generated Go binding around an Ethereum contract.
type IPTokenStaking struct {
	IPTokenStakingCaller     // Read-only binding to the contract
	IPTokenStakingTransactor // Write-only binding to the contract
	IPTokenStakingFilterer   // Log filterer for contract events
}

// IPTokenStakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type IPTokenStakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPTokenStakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IPTokenStakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPTokenStakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IPTokenStakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPTokenStakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IPTokenStakingSession struct {
	Contract     *IPTokenStaking   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IPTokenStakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IPTokenStakingCallerSession struct {
	Contract *IPTokenStakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// IPTokenStakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IPTokenStakingTransactorSession struct {
	Contract     *IPTokenStakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IPTokenStakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type IPTokenStakingRaw struct {
	Contract *IPTokenStaking // Generic contract binding to access the raw methods on
}

// IPTokenStakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IPTokenStakingCallerRaw struct {
	Contract *IPTokenStakingCaller // Generic read-only contract binding to access the raw methods on
}

// IPTokenStakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IPTokenStakingTransactorRaw struct {
	Contract *IPTokenStakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIPTokenStaking creates a new instance of IPTokenStaking, bound to a specific deployed contract.
func NewIPTokenStaking(address common.Address, backend bind.ContractBackend) (*IPTokenStaking, error) {
	contract, err := bindIPTokenStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IPTokenStaking{IPTokenStakingCaller: IPTokenStakingCaller{contract: contract}, IPTokenStakingTransactor: IPTokenStakingTransactor{contract: contract}, IPTokenStakingFilterer: IPTokenStakingFilterer{contract: contract}}, nil
}

// NewIPTokenStakingCaller creates a new read-only instance of IPTokenStaking, bound to a specific deployed contract.
func NewIPTokenStakingCaller(address common.Address, caller bind.ContractCaller) (*IPTokenStakingCaller, error) {
	contract, err := bindIPTokenStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingCaller{contract: contract}, nil
}

// NewIPTokenStakingTransactor creates a new write-only instance of IPTokenStaking, bound to a specific deployed contract.
func NewIPTokenStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*IPTokenStakingTransactor, error) {
	contract, err := bindIPTokenStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingTransactor{contract: contract}, nil
}

// NewIPTokenStakingFilterer creates a new log filterer instance of IPTokenStaking, bound to a specific deployed contract.
func NewIPTokenStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*IPTokenStakingFilterer, error) {
	contract, err := bindIPTokenStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingFilterer{contract: contract}, nil
}

// bindIPTokenStaking binds a generic wrapper to an already deployed contract.
func bindIPTokenStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IPTokenStakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPTokenStaking *IPTokenStakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPTokenStaking.Contract.IPTokenStakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPTokenStaking *IPTokenStakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.IPTokenStakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPTokenStaking *IPTokenStakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.IPTokenStakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPTokenStaking *IPTokenStakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPTokenStaking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPTokenStaking *IPTokenStakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPTokenStaking *IPTokenStakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTMINFEE is a free data retrieval call binding the contract method 0x94fd0fe0.
//
// Solidity: function DEFAULT_MIN_FEE() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) DEFAULTMINFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "DEFAULT_MIN_FEE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEFAULTMINFEE is a free data retrieval call binding the contract method 0x94fd0fe0.
//
// Solidity: function DEFAULT_MIN_FEE() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) DEFAULTMINFEE() (*big.Int, error) {
	return _IPTokenStaking.Contract.DEFAULTMINFEE(&_IPTokenStaking.CallOpts)
}

// DEFAULTMINFEE is a free data retrieval call binding the contract method 0x94fd0fe0.
//
// Solidity: function DEFAULT_MIN_FEE() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) DEFAULTMINFEE() (*big.Int, error) {
	return _IPTokenStaking.Contract.DEFAULTMINFEE(&_IPTokenStaking.CallOpts)
}

// STAKEROUNDING is a free data retrieval call binding the contract method 0xbda16b15.
//
// Solidity: function STAKE_ROUNDING() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) STAKEROUNDING(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "STAKE_ROUNDING")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// STAKEROUNDING is a free data retrieval call binding the contract method 0xbda16b15.
//
// Solidity: function STAKE_ROUNDING() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) STAKEROUNDING() (*big.Int, error) {
	return _IPTokenStaking.Contract.STAKEROUNDING(&_IPTokenStaking.CallOpts)
}

// STAKEROUNDING is a free data retrieval call binding the contract method 0xbda16b15.
//
// Solidity: function STAKE_ROUNDING() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) STAKEROUNDING() (*big.Int, error) {
	return _IPTokenStaking.Contract.STAKEROUNDING(&_IPTokenStaking.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) Fee() (*big.Int, error) {
	return _IPTokenStaking.Contract.Fee(&_IPTokenStaking.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) Fee() (*big.Int, error) {
	return _IPTokenStaking.Contract.Fee(&_IPTokenStaking.CallOpts)
}

// MinCommissionRate is a free data retrieval call binding the contract method 0x1487153e.
//
// Solidity: function minCommissionRate() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) MinCommissionRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "minCommissionRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinCommissionRate is a free data retrieval call binding the contract method 0x1487153e.
//
// Solidity: function minCommissionRate() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) MinCommissionRate() (*big.Int, error) {
	return _IPTokenStaking.Contract.MinCommissionRate(&_IPTokenStaking.CallOpts)
}

// MinCommissionRate is a free data retrieval call binding the contract method 0x1487153e.
//
// Solidity: function minCommissionRate() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) MinCommissionRate() (*big.Int, error) {
	return _IPTokenStaking.Contract.MinCommissionRate(&_IPTokenStaking.CallOpts)
}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) MinStakeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "minStakeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) MinStakeAmount() (*big.Int, error) {
	return _IPTokenStaking.Contract.MinStakeAmount(&_IPTokenStaking.CallOpts)
}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) MinStakeAmount() (*big.Int, error) {
	return _IPTokenStaking.Contract.MinStakeAmount(&_IPTokenStaking.CallOpts)
}

// MinUnstakeAmount is a free data retrieval call binding the contract method 0x39ec4df9.
//
// Solidity: function minUnstakeAmount() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) MinUnstakeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "minUnstakeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinUnstakeAmount is a free data retrieval call binding the contract method 0x39ec4df9.
//
// Solidity: function minUnstakeAmount() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) MinUnstakeAmount() (*big.Int, error) {
	return _IPTokenStaking.Contract.MinUnstakeAmount(&_IPTokenStaking.CallOpts)
}

// MinUnstakeAmount is a free data retrieval call binding the contract method 0x39ec4df9.
//
// Solidity: function minUnstakeAmount() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) MinUnstakeAmount() (*big.Int, error) {
	return _IPTokenStaking.Contract.MinUnstakeAmount(&_IPTokenStaking.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IPTokenStaking *IPTokenStakingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IPTokenStaking *IPTokenStakingSession) Owner() (common.Address, error) {
	return _IPTokenStaking.Contract.Owner(&_IPTokenStaking.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IPTokenStaking *IPTokenStakingCallerSession) Owner() (common.Address, error) {
	return _IPTokenStaking.Contract.Owner(&_IPTokenStaking.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_IPTokenStaking *IPTokenStakingCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_IPTokenStaking *IPTokenStakingSession) PendingOwner() (common.Address, error) {
	return _IPTokenStaking.Contract.PendingOwner(&_IPTokenStaking.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_IPTokenStaking *IPTokenStakingCallerSession) PendingOwner() (common.Address, error) {
	return _IPTokenStaking.Contract.PendingOwner(&_IPTokenStaking.CallOpts)
}

// RoundedStakeAmount is a free data retrieval call binding the contract method 0xd2e1f5b8.
//
// Solidity: function roundedStakeAmount(uint256 rawAmount) view returns(uint256 amount, uint256 remainder)
func (_IPTokenStaking *IPTokenStakingCaller) RoundedStakeAmount(opts *bind.CallOpts, rawAmount *big.Int) (struct {
	Amount    *big.Int
	Remainder *big.Int
}, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "roundedStakeAmount", rawAmount)

	outstruct := new(struct {
		Amount    *big.Int
		Remainder *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Remainder = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// RoundedStakeAmount is a free data retrieval call binding the contract method 0xd2e1f5b8.
//
// Solidity: function roundedStakeAmount(uint256 rawAmount) view returns(uint256 amount, uint256 remainder)
func (_IPTokenStaking *IPTokenStakingSession) RoundedStakeAmount(rawAmount *big.Int) (struct {
	Amount    *big.Int
	Remainder *big.Int
}, error) {
	return _IPTokenStaking.Contract.RoundedStakeAmount(&_IPTokenStaking.CallOpts, rawAmount)
}

// RoundedStakeAmount is a free data retrieval call binding the contract method 0xd2e1f5b8.
//
// Solidity: function roundedStakeAmount(uint256 rawAmount) view returns(uint256 amount, uint256 remainder)
func (_IPTokenStaking *IPTokenStakingCallerSession) RoundedStakeAmount(rawAmount *big.Int) (struct {
	Amount    *big.Int
	Remainder *big.Int
}, error) {
	return _IPTokenStaking.Contract.RoundedStakeAmount(&_IPTokenStaking.CallOpts, rawAmount)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_IPTokenStaking *IPTokenStakingTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_IPTokenStaking *IPTokenStakingSession) AcceptOwnership() (*types.Transaction, error) {
	return _IPTokenStaking.Contract.AcceptOwnership(&_IPTokenStaking.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _IPTokenStaking.Contract.AcceptOwnership(&_IPTokenStaking.TransactOpts)
}

// AddOperator is a paid mutator transaction binding the contract method 0x057b9296.
//
// Solidity: function addOperator(bytes uncmpPubkey, address operator) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) AddOperator(opts *bind.TransactOpts, uncmpPubkey []byte, operator common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "addOperator", uncmpPubkey, operator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x057b9296.
//
// Solidity: function addOperator(bytes uncmpPubkey, address operator) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) AddOperator(uncmpPubkey []byte, operator common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.AddOperator(&_IPTokenStaking.TransactOpts, uncmpPubkey, operator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x057b9296.
//
// Solidity: function addOperator(bytes uncmpPubkey, address operator) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) AddOperator(uncmpPubkey []byte, operator common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.AddOperator(&_IPTokenStaking.TransactOpts, uncmpPubkey, operator)
}

// CreateValidator is a paid mutator transaction binding the contract method 0x8740597a.
//
// Solidity: function createValidator(bytes validatorUncmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool supportsUnlocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) CreateValidator(opts *bind.TransactOpts, validatorUncmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, supportsUnlocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "createValidator", validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, supportsUnlocked, data)
}

// CreateValidator is a paid mutator transaction binding the contract method 0x8740597a.
//
// Solidity: function createValidator(bytes validatorUncmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool supportsUnlocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) CreateValidator(validatorUncmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, supportsUnlocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.CreateValidator(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, supportsUnlocked, data)
}

// CreateValidator is a paid mutator transaction binding the contract method 0x8740597a.
//
// Solidity: function createValidator(bytes validatorUncmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool supportsUnlocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) CreateValidator(validatorUncmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, supportsUnlocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.CreateValidator(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, supportsUnlocked, data)
}

// CreateValidatorOnBehalf is a paid mutator transaction binding the contract method 0xf9550a8d.
//
// Solidity: function createValidatorOnBehalf(bytes validatorUncmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool supportsUnlocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) CreateValidatorOnBehalf(opts *bind.TransactOpts, validatorUncmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, supportsUnlocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "createValidatorOnBehalf", validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, supportsUnlocked, data)
}

// CreateValidatorOnBehalf is a paid mutator transaction binding the contract method 0xf9550a8d.
//
// Solidity: function createValidatorOnBehalf(bytes validatorUncmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool supportsUnlocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) CreateValidatorOnBehalf(validatorUncmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, supportsUnlocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.CreateValidatorOnBehalf(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, supportsUnlocked, data)
}

// CreateValidatorOnBehalf is a paid mutator transaction binding the contract method 0xf9550a8d.
//
// Solidity: function createValidatorOnBehalf(bytes validatorUncmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool supportsUnlocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) CreateValidatorOnBehalf(validatorUncmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, supportsUnlocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.CreateValidatorOnBehalf(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, supportsUnlocked, data)
}

// Initialize is a paid mutator transaction binding the contract method 0xfce5dc8c.
//
// Solidity: function initialize((address,uint256,uint256,uint256,uint256) args) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) Initialize(opts *bind.TransactOpts, args IIPTokenStakingInitializerArgs) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "initialize", args)
}

// Initialize is a paid mutator transaction binding the contract method 0xfce5dc8c.
//
// Solidity: function initialize((address,uint256,uint256,uint256,uint256) args) returns()
func (_IPTokenStaking *IPTokenStakingSession) Initialize(args IIPTokenStakingInitializerArgs) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Initialize(&_IPTokenStaking.TransactOpts, args)
}

// Initialize is a paid mutator transaction binding the contract method 0xfce5dc8c.
//
// Solidity: function initialize((address,uint256,uint256,uint256,uint256) args) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) Initialize(args IIPTokenStakingInitializerArgs) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Initialize(&_IPTokenStaking.TransactOpts, args)
}

// Redelegate is a paid mutator transaction binding the contract method 0x9d9d293f.
//
// Solidity: function redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 delegationId, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) Redelegate(opts *bind.TransactOpts, delegatorUncmpPubkey []byte, validatorUncmpSrcPubkey []byte, validatorUncmpDstPubkey []byte, delegationId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "redelegate", delegatorUncmpPubkey, validatorUncmpSrcPubkey, validatorUncmpDstPubkey, delegationId, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x9d9d293f.
//
// Solidity: function redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 delegationId, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) Redelegate(delegatorUncmpPubkey []byte, validatorUncmpSrcPubkey []byte, validatorUncmpDstPubkey []byte, delegationId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Redelegate(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpSrcPubkey, validatorUncmpDstPubkey, delegationId, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x9d9d293f.
//
// Solidity: function redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 delegationId, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) Redelegate(delegatorUncmpPubkey []byte, validatorUncmpSrcPubkey []byte, validatorUncmpDstPubkey []byte, delegationId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Redelegate(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpSrcPubkey, validatorUncmpDstPubkey, delegationId, amount)
}

// RedelegateOnBehalf is a paid mutator transaction binding the contract method 0xec21dac2.
//
// Solidity: function redelegateOnBehalf(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 delegationId, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) RedelegateOnBehalf(opts *bind.TransactOpts, delegatorUncmpPubkey []byte, validatorUncmpSrcPubkey []byte, validatorUncmpDstPubkey []byte, delegationId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "redelegateOnBehalf", delegatorUncmpPubkey, validatorUncmpSrcPubkey, validatorUncmpDstPubkey, delegationId, amount)
}

// RedelegateOnBehalf is a paid mutator transaction binding the contract method 0xec21dac2.
//
// Solidity: function redelegateOnBehalf(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 delegationId, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) RedelegateOnBehalf(delegatorUncmpPubkey []byte, validatorUncmpSrcPubkey []byte, validatorUncmpDstPubkey []byte, delegationId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.RedelegateOnBehalf(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpSrcPubkey, validatorUncmpDstPubkey, delegationId, amount)
}

// RedelegateOnBehalf is a paid mutator transaction binding the contract method 0xec21dac2.
//
// Solidity: function redelegateOnBehalf(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 delegationId, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) RedelegateOnBehalf(delegatorUncmpPubkey []byte, validatorUncmpSrcPubkey []byte, validatorUncmpDstPubkey []byte, delegationId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.RedelegateOnBehalf(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpSrcPubkey, validatorUncmpDstPubkey, delegationId, amount)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0x17e42e12.
//
// Solidity: function removeOperator(bytes uncmpPubkey, address operator) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) RemoveOperator(opts *bind.TransactOpts, uncmpPubkey []byte, operator common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "removeOperator", uncmpPubkey, operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0x17e42e12.
//
// Solidity: function removeOperator(bytes uncmpPubkey, address operator) returns()
func (_IPTokenStaking *IPTokenStakingSession) RemoveOperator(uncmpPubkey []byte, operator common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.RemoveOperator(&_IPTokenStaking.TransactOpts, uncmpPubkey, operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0x17e42e12.
//
// Solidity: function removeOperator(bytes uncmpPubkey, address operator) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) RemoveOperator(uncmpPubkey []byte, operator common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.RemoveOperator(&_IPTokenStaking.TransactOpts, uncmpPubkey, operator)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IPTokenStaking *IPTokenStakingTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IPTokenStaking *IPTokenStakingSession) RenounceOwnership() (*types.Transaction, error) {
	return _IPTokenStaking.Contract.RenounceOwnership(&_IPTokenStaking.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _IPTokenStaking.Contract.RenounceOwnership(&_IPTokenStaking.TransactOpts)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 newFee) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetFee(opts *bind.TransactOpts, newFee *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setFee", newFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 newFee) returns()
func (_IPTokenStaking *IPTokenStakingSession) SetFee(newFee *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetFee(&_IPTokenStaking.TransactOpts, newFee)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 newFee) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetFee(newFee *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetFee(&_IPTokenStaking.TransactOpts, newFee)
}

// SetMinCommissionRate is a paid mutator transaction binding the contract method 0xab8870f6.
//
// Solidity: function setMinCommissionRate(uint256 newValue) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetMinCommissionRate(opts *bind.TransactOpts, newValue *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setMinCommissionRate", newValue)
}

// SetMinCommissionRate is a paid mutator transaction binding the contract method 0xab8870f6.
//
// Solidity: function setMinCommissionRate(uint256 newValue) returns()
func (_IPTokenStaking *IPTokenStakingSession) SetMinCommissionRate(newValue *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetMinCommissionRate(&_IPTokenStaking.TransactOpts, newValue)
}

// SetMinCommissionRate is a paid mutator transaction binding the contract method 0xab8870f6.
//
// Solidity: function setMinCommissionRate(uint256 newValue) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetMinCommissionRate(newValue *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetMinCommissionRate(&_IPTokenStaking.TransactOpts, newValue)
}

// SetMinStakeAmount is a paid mutator transaction binding the contract method 0xeb4af045.
//
// Solidity: function setMinStakeAmount(uint256 newMinStakeAmount) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetMinStakeAmount(opts *bind.TransactOpts, newMinStakeAmount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setMinStakeAmount", newMinStakeAmount)
}

// SetMinStakeAmount is a paid mutator transaction binding the contract method 0xeb4af045.
//
// Solidity: function setMinStakeAmount(uint256 newMinStakeAmount) returns()
func (_IPTokenStaking *IPTokenStakingSession) SetMinStakeAmount(newMinStakeAmount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetMinStakeAmount(&_IPTokenStaking.TransactOpts, newMinStakeAmount)
}

// SetMinStakeAmount is a paid mutator transaction binding the contract method 0xeb4af045.
//
// Solidity: function setMinStakeAmount(uint256 newMinStakeAmount) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetMinStakeAmount(newMinStakeAmount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetMinStakeAmount(&_IPTokenStaking.TransactOpts, newMinStakeAmount)
}

// SetMinUnstakeAmount is a paid mutator transaction binding the contract method 0x6ea3a228.
//
// Solidity: function setMinUnstakeAmount(uint256 newMinUnstakeAmount) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetMinUnstakeAmount(opts *bind.TransactOpts, newMinUnstakeAmount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setMinUnstakeAmount", newMinUnstakeAmount)
}

// SetMinUnstakeAmount is a paid mutator transaction binding the contract method 0x6ea3a228.
//
// Solidity: function setMinUnstakeAmount(uint256 newMinUnstakeAmount) returns()
func (_IPTokenStaking *IPTokenStakingSession) SetMinUnstakeAmount(newMinUnstakeAmount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetMinUnstakeAmount(&_IPTokenStaking.TransactOpts, newMinUnstakeAmount)
}

// SetMinUnstakeAmount is a paid mutator transaction binding the contract method 0x6ea3a228.
//
// Solidity: function setMinUnstakeAmount(uint256 newMinUnstakeAmount) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetMinUnstakeAmount(newMinUnstakeAmount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetMinUnstakeAmount(&_IPTokenStaking.TransactOpts, newMinUnstakeAmount)
}

// SetRewardsAddress is a paid mutator transaction binding the contract method 0x9d04b121.
//
// Solidity: function setRewardsAddress(bytes delegatorUncmpPubkey, address newRewardsAddress) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetRewardsAddress(opts *bind.TransactOpts, delegatorUncmpPubkey []byte, newRewardsAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setRewardsAddress", delegatorUncmpPubkey, newRewardsAddress)
}

// SetRewardsAddress is a paid mutator transaction binding the contract method 0x9d04b121.
//
// Solidity: function setRewardsAddress(bytes delegatorUncmpPubkey, address newRewardsAddress) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) SetRewardsAddress(delegatorUncmpPubkey []byte, newRewardsAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetRewardsAddress(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, newRewardsAddress)
}

// SetRewardsAddress is a paid mutator transaction binding the contract method 0x9d04b121.
//
// Solidity: function setRewardsAddress(bytes delegatorUncmpPubkey, address newRewardsAddress) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetRewardsAddress(delegatorUncmpPubkey []byte, newRewardsAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetRewardsAddress(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, newRewardsAddress)
}

// SetWithdrawalAddress is a paid mutator transaction binding the contract method 0x787f82c8.
//
// Solidity: function setWithdrawalAddress(bytes delegatorUncmpPubkey, address newWithdrawalAddress) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetWithdrawalAddress(opts *bind.TransactOpts, delegatorUncmpPubkey []byte, newWithdrawalAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setWithdrawalAddress", delegatorUncmpPubkey, newWithdrawalAddress)
}

// SetWithdrawalAddress is a paid mutator transaction binding the contract method 0x787f82c8.
//
// Solidity: function setWithdrawalAddress(bytes delegatorUncmpPubkey, address newWithdrawalAddress) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) SetWithdrawalAddress(delegatorUncmpPubkey []byte, newWithdrawalAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetWithdrawalAddress(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, newWithdrawalAddress)
}

// SetWithdrawalAddress is a paid mutator transaction binding the contract method 0x787f82c8.
//
// Solidity: function setWithdrawalAddress(bytes delegatorUncmpPubkey, address newWithdrawalAddress) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetWithdrawalAddress(delegatorUncmpPubkey []byte, newWithdrawalAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetWithdrawalAddress(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, newWithdrawalAddress)
}

// Stake is a paid mutator transaction binding the contract method 0x3dd9fb9a.
//
// Solidity: function stake(bytes delegatorUncmpPubkey, bytes validatorUncmpPubkey, uint8 stakingPeriod, bytes data) payable returns(uint256 delegationId)
func (_IPTokenStaking *IPTokenStakingTransactor) Stake(opts *bind.TransactOpts, delegatorUncmpPubkey []byte, validatorUncmpPubkey []byte, stakingPeriod uint8, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "stake", delegatorUncmpPubkey, validatorUncmpPubkey, stakingPeriod, data)
}

// Stake is a paid mutator transaction binding the contract method 0x3dd9fb9a.
//
// Solidity: function stake(bytes delegatorUncmpPubkey, bytes validatorUncmpPubkey, uint8 stakingPeriod, bytes data) payable returns(uint256 delegationId)
func (_IPTokenStaking *IPTokenStakingSession) Stake(delegatorUncmpPubkey []byte, validatorUncmpPubkey []byte, stakingPeriod uint8, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Stake(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpPubkey, stakingPeriod, data)
}

// Stake is a paid mutator transaction binding the contract method 0x3dd9fb9a.
//
// Solidity: function stake(bytes delegatorUncmpPubkey, bytes validatorUncmpPubkey, uint8 stakingPeriod, bytes data) payable returns(uint256 delegationId)
func (_IPTokenStaking *IPTokenStakingTransactorSession) Stake(delegatorUncmpPubkey []byte, validatorUncmpPubkey []byte, stakingPeriod uint8, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Stake(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpPubkey, stakingPeriod, data)
}

// StakeOnBehalf is a paid mutator transaction binding the contract method 0xa0284f16.
//
// Solidity: function stakeOnBehalf(bytes delegatorUncmpPubkey, bytes validatorUncmpPubkey, uint8 stakingPeriod, bytes data) payable returns(uint256 delegationId)
func (_IPTokenStaking *IPTokenStakingTransactor) StakeOnBehalf(opts *bind.TransactOpts, delegatorUncmpPubkey []byte, validatorUncmpPubkey []byte, stakingPeriod uint8, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "stakeOnBehalf", delegatorUncmpPubkey, validatorUncmpPubkey, stakingPeriod, data)
}

// StakeOnBehalf is a paid mutator transaction binding the contract method 0xa0284f16.
//
// Solidity: function stakeOnBehalf(bytes delegatorUncmpPubkey, bytes validatorUncmpPubkey, uint8 stakingPeriod, bytes data) payable returns(uint256 delegationId)
func (_IPTokenStaking *IPTokenStakingSession) StakeOnBehalf(delegatorUncmpPubkey []byte, validatorUncmpPubkey []byte, stakingPeriod uint8, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.StakeOnBehalf(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpPubkey, stakingPeriod, data)
}

// StakeOnBehalf is a paid mutator transaction binding the contract method 0xa0284f16.
//
// Solidity: function stakeOnBehalf(bytes delegatorUncmpPubkey, bytes validatorUncmpPubkey, uint8 stakingPeriod, bytes data) payable returns(uint256 delegationId)
func (_IPTokenStaking *IPTokenStakingTransactorSession) StakeOnBehalf(delegatorUncmpPubkey []byte, validatorUncmpPubkey []byte, stakingPeriod uint8, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.StakeOnBehalf(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpPubkey, stakingPeriod, data)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IPTokenStaking *IPTokenStakingSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.TransferOwnership(&_IPTokenStaking.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.TransferOwnership(&_IPTokenStaking.TransactOpts, newOwner)
}

// Unjail is a paid mutator transaction binding the contract method 0x8ed65fbc.
//
// Solidity: function unjail(bytes validatorUncmpPubkey, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) Unjail(opts *bind.TransactOpts, validatorUncmpPubkey []byte, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "unjail", validatorUncmpPubkey, data)
}

// Unjail is a paid mutator transaction binding the contract method 0x8ed65fbc.
//
// Solidity: function unjail(bytes validatorUncmpPubkey, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) Unjail(validatorUncmpPubkey []byte, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Unjail(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, data)
}

// Unjail is a paid mutator transaction binding the contract method 0x8ed65fbc.
//
// Solidity: function unjail(bytes validatorUncmpPubkey, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) Unjail(validatorUncmpPubkey []byte, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Unjail(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, data)
}

// UnjailOnBehalf is a paid mutator transaction binding the contract method 0x86eb5e48.
//
// Solidity: function unjailOnBehalf(bytes validatorUncmpPubkey, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) UnjailOnBehalf(opts *bind.TransactOpts, validatorUncmpPubkey []byte, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "unjailOnBehalf", validatorUncmpPubkey, data)
}

// UnjailOnBehalf is a paid mutator transaction binding the contract method 0x86eb5e48.
//
// Solidity: function unjailOnBehalf(bytes validatorUncmpPubkey, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) UnjailOnBehalf(validatorUncmpPubkey []byte, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UnjailOnBehalf(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, data)
}

// UnjailOnBehalf is a paid mutator transaction binding the contract method 0x86eb5e48.
//
// Solidity: function unjailOnBehalf(bytes validatorUncmpPubkey, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) UnjailOnBehalf(validatorUncmpPubkey []byte, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UnjailOnBehalf(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, data)
}

// Unstake is a paid mutator transaction binding the contract method 0xb2bc29ef.
//
// Solidity: function unstake(bytes delegatorUncmpPubkey, bytes validatorUncmpPubkey, uint256 delegationId, uint256 amount, bytes data) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) Unstake(opts *bind.TransactOpts, delegatorUncmpPubkey []byte, validatorUncmpPubkey []byte, delegationId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "unstake", delegatorUncmpPubkey, validatorUncmpPubkey, delegationId, amount, data)
}

// Unstake is a paid mutator transaction binding the contract method 0xb2bc29ef.
//
// Solidity: function unstake(bytes delegatorUncmpPubkey, bytes validatorUncmpPubkey, uint256 delegationId, uint256 amount, bytes data) returns()
func (_IPTokenStaking *IPTokenStakingSession) Unstake(delegatorUncmpPubkey []byte, validatorUncmpPubkey []byte, delegationId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Unstake(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpPubkey, delegationId, amount, data)
}

// Unstake is a paid mutator transaction binding the contract method 0xb2bc29ef.
//
// Solidity: function unstake(bytes delegatorUncmpPubkey, bytes validatorUncmpPubkey, uint256 delegationId, uint256 amount, bytes data) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) Unstake(delegatorUncmpPubkey []byte, validatorUncmpPubkey []byte, delegationId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Unstake(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpPubkey, delegationId, amount, data)
}

// UnstakeOnBehalf is a paid mutator transaction binding the contract method 0x014e8178.
//
// Solidity: function unstakeOnBehalf(bytes delegatorUncmpPubkey, bytes validatorUncmpPubkey, uint256 delegationId, uint256 amount, bytes data) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) UnstakeOnBehalf(opts *bind.TransactOpts, delegatorUncmpPubkey []byte, validatorUncmpPubkey []byte, delegationId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "unstakeOnBehalf", delegatorUncmpPubkey, validatorUncmpPubkey, delegationId, amount, data)
}

// UnstakeOnBehalf is a paid mutator transaction binding the contract method 0x014e8178.
//
// Solidity: function unstakeOnBehalf(bytes delegatorUncmpPubkey, bytes validatorUncmpPubkey, uint256 delegationId, uint256 amount, bytes data) returns()
func (_IPTokenStaking *IPTokenStakingSession) UnstakeOnBehalf(delegatorUncmpPubkey []byte, validatorUncmpPubkey []byte, delegationId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UnstakeOnBehalf(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpPubkey, delegationId, amount, data)
}

// UnstakeOnBehalf is a paid mutator transaction binding the contract method 0x014e8178.
//
// Solidity: function unstakeOnBehalf(bytes delegatorUncmpPubkey, bytes validatorUncmpPubkey, uint256 delegationId, uint256 amount, bytes data) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) UnstakeOnBehalf(delegatorUncmpPubkey []byte, validatorUncmpPubkey []byte, delegationId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UnstakeOnBehalf(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpPubkey, delegationId, amount, data)
}

// UpdateValidatorCommission is a paid mutator transaction binding the contract method 0xc582db44.
//
// Solidity: function updateValidatorCommission(bytes validatorUncmpPubkey, uint32 commissionRate) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) UpdateValidatorCommission(opts *bind.TransactOpts, validatorUncmpPubkey []byte, commissionRate uint32) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "updateValidatorCommission", validatorUncmpPubkey, commissionRate)
}

// UpdateValidatorCommission is a paid mutator transaction binding the contract method 0xc582db44.
//
// Solidity: function updateValidatorCommission(bytes validatorUncmpPubkey, uint32 commissionRate) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) UpdateValidatorCommission(validatorUncmpPubkey []byte, commissionRate uint32) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UpdateValidatorCommission(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, commissionRate)
}

// UpdateValidatorCommission is a paid mutator transaction binding the contract method 0xc582db44.
//
// Solidity: function updateValidatorCommission(bytes validatorUncmpPubkey, uint32 commissionRate) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) UpdateValidatorCommission(validatorUncmpPubkey []byte, commissionRate uint32) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UpdateValidatorCommission(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, commissionRate)
}

// IPTokenStakingAddOperatorIterator is returned from FilterAddOperator and is used to iterate over the raw logs and unpacked data for AddOperator events raised by the IPTokenStaking contract.
type IPTokenStakingAddOperatorIterator struct {
	Event *IPTokenStakingAddOperator // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingAddOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingAddOperator)
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
		it.Event = new(IPTokenStakingAddOperator)
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
func (it *IPTokenStakingAddOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingAddOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingAddOperator represents a AddOperator event raised by the IPTokenStaking contract.
type IPTokenStakingAddOperator struct {
	UncmpPubkey []byte
	Operator    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAddOperator is a free log retrieval operation binding the contract event 0x6ac365cf05479bb8a295fbf9637875411d6d6f2a0ac7c4b1f560cedcf1a33081.
//
// Solidity: event AddOperator(bytes uncmpPubkey, address operator)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterAddOperator(opts *bind.FilterOpts) (*IPTokenStakingAddOperatorIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "AddOperator")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingAddOperatorIterator{contract: _IPTokenStaking.contract, event: "AddOperator", logs: logs, sub: sub}, nil
}

// WatchAddOperator is a free log subscription operation binding the contract event 0x6ac365cf05479bb8a295fbf9637875411d6d6f2a0ac7c4b1f560cedcf1a33081.
//
// Solidity: event AddOperator(bytes uncmpPubkey, address operator)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchAddOperator(opts *bind.WatchOpts, sink chan<- *IPTokenStakingAddOperator) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "AddOperator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingAddOperator)
				if err := _IPTokenStaking.contract.UnpackLog(event, "AddOperator", log); err != nil {
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

// ParseAddOperator is a log parse operation binding the contract event 0x6ac365cf05479bb8a295fbf9637875411d6d6f2a0ac7c4b1f560cedcf1a33081.
//
// Solidity: event AddOperator(bytes uncmpPubkey, address operator)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseAddOperator(log types.Log) (*IPTokenStakingAddOperator, error) {
	event := new(IPTokenStakingAddOperator)
	if err := _IPTokenStaking.contract.UnpackLog(event, "AddOperator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingCreateValidatorIterator is returned from FilterCreateValidator and is used to iterate over the raw logs and unpacked data for CreateValidator events raised by the IPTokenStaking contract.
type IPTokenStakingCreateValidatorIterator struct {
	Event *IPTokenStakingCreateValidator // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingCreateValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingCreateValidator)
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
		it.Event = new(IPTokenStakingCreateValidator)
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
func (it *IPTokenStakingCreateValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingCreateValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingCreateValidator represents a CreateValidator event raised by the IPTokenStaking contract.
type IPTokenStakingCreateValidator struct {
	ValidatorUncmpPubkey    []byte
	Moniker                 string
	StakeAmount             *big.Int
	CommissionRate          uint32
	MaxCommissionRate       uint32
	MaxCommissionChangeRate uint32
	SupportsUnlocked        uint8
	OperatorAddress         common.Address
	Data                    []byte
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterCreateValidator is a free log retrieval operation binding the contract event 0x65bfc2fa1cd4c6f50f60983ad1cf1cb4bff5ee6570428254dfce41b085ef6d14.
//
// Solidity: event CreateValidator(bytes validatorUncmpPubkey, string moniker, uint256 stakeAmount, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, uint8 supportsUnlocked, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterCreateValidator(opts *bind.FilterOpts) (*IPTokenStakingCreateValidatorIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "CreateValidator")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingCreateValidatorIterator{contract: _IPTokenStaking.contract, event: "CreateValidator", logs: logs, sub: sub}, nil
}

// WatchCreateValidator is a free log subscription operation binding the contract event 0x65bfc2fa1cd4c6f50f60983ad1cf1cb4bff5ee6570428254dfce41b085ef6d14.
//
// Solidity: event CreateValidator(bytes validatorUncmpPubkey, string moniker, uint256 stakeAmount, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, uint8 supportsUnlocked, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchCreateValidator(opts *bind.WatchOpts, sink chan<- *IPTokenStakingCreateValidator) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "CreateValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingCreateValidator)
				if err := _IPTokenStaking.contract.UnpackLog(event, "CreateValidator", log); err != nil {
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

// ParseCreateValidator is a log parse operation binding the contract event 0x65bfc2fa1cd4c6f50f60983ad1cf1cb4bff5ee6570428254dfce41b085ef6d14.
//
// Solidity: event CreateValidator(bytes validatorUncmpPubkey, string moniker, uint256 stakeAmount, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, uint8 supportsUnlocked, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseCreateValidator(log types.Log) (*IPTokenStakingCreateValidator, error) {
	event := new(IPTokenStakingCreateValidator)
	if err := _IPTokenStaking.contract.UnpackLog(event, "CreateValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the IPTokenStaking contract.
type IPTokenStakingDepositIterator struct {
	Event *IPTokenStakingDeposit // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingDeposit)
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
		it.Event = new(IPTokenStakingDeposit)
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
func (it *IPTokenStakingDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingDeposit represents a Deposit event raised by the IPTokenStaking contract.
type IPTokenStakingDeposit struct {
	DelegatorUncmpPubkey []byte
	ValidatorUnCmpPubkey []byte
	StakeAmount          *big.Int
	StakingPeriod        *big.Int
	DelegationId         *big.Int
	OperatorAddress      common.Address
	Data                 []byte
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x269a32ff589c9b701f49ab6aa532ee8f55901df71a7fca2d70dc9f45314f1be3.
//
// Solidity: event Deposit(bytes delegatorUncmpPubkey, bytes validatorUnCmpPubkey, uint256 stakeAmount, uint256 stakingPeriod, uint256 delegationId, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterDeposit(opts *bind.FilterOpts) (*IPTokenStakingDepositIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingDepositIterator{contract: _IPTokenStaking.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x269a32ff589c9b701f49ab6aa532ee8f55901df71a7fca2d70dc9f45314f1be3.
//
// Solidity: event Deposit(bytes delegatorUncmpPubkey, bytes validatorUnCmpPubkey, uint256 stakeAmount, uint256 stakingPeriod, uint256 delegationId, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *IPTokenStakingDeposit) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingDeposit)
				if err := _IPTokenStaking.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0x269a32ff589c9b701f49ab6aa532ee8f55901df71a7fca2d70dc9f45314f1be3.
//
// Solidity: event Deposit(bytes delegatorUncmpPubkey, bytes validatorUnCmpPubkey, uint256 stakeAmount, uint256 stakingPeriod, uint256 delegationId, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseDeposit(log types.Log) (*IPTokenStakingDeposit, error) {
	event := new(IPTokenStakingDeposit)
	if err := _IPTokenStaking.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingFeeSetIterator is returned from FilterFeeSet and is used to iterate over the raw logs and unpacked data for FeeSet events raised by the IPTokenStaking contract.
type IPTokenStakingFeeSetIterator struct {
	Event *IPTokenStakingFeeSet // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingFeeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingFeeSet)
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
		it.Event = new(IPTokenStakingFeeSet)
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
func (it *IPTokenStakingFeeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingFeeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingFeeSet represents a FeeSet event raised by the IPTokenStaking contract.
type IPTokenStakingFeeSet struct {
	NewFee *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeeSet is a free log retrieval operation binding the contract event 0x20461e09b8e557b77e107939f9ce6544698123aad0fc964ac5cc59b7df2e608f.
//
// Solidity: event FeeSet(uint256 newFee)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterFeeSet(opts *bind.FilterOpts) (*IPTokenStakingFeeSetIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "FeeSet")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingFeeSetIterator{contract: _IPTokenStaking.contract, event: "FeeSet", logs: logs, sub: sub}, nil
}

// WatchFeeSet is a free log subscription operation binding the contract event 0x20461e09b8e557b77e107939f9ce6544698123aad0fc964ac5cc59b7df2e608f.
//
// Solidity: event FeeSet(uint256 newFee)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchFeeSet(opts *bind.WatchOpts, sink chan<- *IPTokenStakingFeeSet) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "FeeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingFeeSet)
				if err := _IPTokenStaking.contract.UnpackLog(event, "FeeSet", log); err != nil {
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

// ParseFeeSet is a log parse operation binding the contract event 0x20461e09b8e557b77e107939f9ce6544698123aad0fc964ac5cc59b7df2e608f.
//
// Solidity: event FeeSet(uint256 newFee)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseFeeSet(log types.Log) (*IPTokenStakingFeeSet, error) {
	event := new(IPTokenStakingFeeSet)
	if err := _IPTokenStaking.contract.UnpackLog(event, "FeeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the IPTokenStaking contract.
type IPTokenStakingInitializedIterator struct {
	Event *IPTokenStakingInitialized // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingInitialized)
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
		it.Event = new(IPTokenStakingInitialized)
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
func (it *IPTokenStakingInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingInitialized represents a Initialized event raised by the IPTokenStaking contract.
type IPTokenStakingInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterInitialized(opts *bind.FilterOpts) (*IPTokenStakingInitializedIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingInitializedIterator{contract: _IPTokenStaking.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *IPTokenStakingInitialized) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingInitialized)
				if err := _IPTokenStaking.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_IPTokenStaking *IPTokenStakingFilterer) ParseInitialized(log types.Log) (*IPTokenStakingInitialized, error) {
	event := new(IPTokenStakingInitialized)
	if err := _IPTokenStaking.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingMinCommissionRateChangedIterator is returned from FilterMinCommissionRateChanged and is used to iterate over the raw logs and unpacked data for MinCommissionRateChanged events raised by the IPTokenStaking contract.
type IPTokenStakingMinCommissionRateChangedIterator struct {
	Event *IPTokenStakingMinCommissionRateChanged // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingMinCommissionRateChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingMinCommissionRateChanged)
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
		it.Event = new(IPTokenStakingMinCommissionRateChanged)
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
func (it *IPTokenStakingMinCommissionRateChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingMinCommissionRateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingMinCommissionRateChanged represents a MinCommissionRateChanged event raised by the IPTokenStaking contract.
type IPTokenStakingMinCommissionRateChanged struct {
	MinCommissionRate *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterMinCommissionRateChanged is a free log retrieval operation binding the contract event 0x4167b1de65292a9ff628c9136823791a1de701e1fbdda4863ce22a1cfaf4d0f7.
//
// Solidity: event MinCommissionRateChanged(uint256 minCommissionRate)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterMinCommissionRateChanged(opts *bind.FilterOpts) (*IPTokenStakingMinCommissionRateChangedIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "MinCommissionRateChanged")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingMinCommissionRateChangedIterator{contract: _IPTokenStaking.contract, event: "MinCommissionRateChanged", logs: logs, sub: sub}, nil
}

// WatchMinCommissionRateChanged is a free log subscription operation binding the contract event 0x4167b1de65292a9ff628c9136823791a1de701e1fbdda4863ce22a1cfaf4d0f7.
//
// Solidity: event MinCommissionRateChanged(uint256 minCommissionRate)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchMinCommissionRateChanged(opts *bind.WatchOpts, sink chan<- *IPTokenStakingMinCommissionRateChanged) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "MinCommissionRateChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingMinCommissionRateChanged)
				if err := _IPTokenStaking.contract.UnpackLog(event, "MinCommissionRateChanged", log); err != nil {
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

// ParseMinCommissionRateChanged is a log parse operation binding the contract event 0x4167b1de65292a9ff628c9136823791a1de701e1fbdda4863ce22a1cfaf4d0f7.
//
// Solidity: event MinCommissionRateChanged(uint256 minCommissionRate)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseMinCommissionRateChanged(log types.Log) (*IPTokenStakingMinCommissionRateChanged, error) {
	event := new(IPTokenStakingMinCommissionRateChanged)
	if err := _IPTokenStaking.contract.UnpackLog(event, "MinCommissionRateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingMinStakeAmountSetIterator is returned from FilterMinStakeAmountSet and is used to iterate over the raw logs and unpacked data for MinStakeAmountSet events raised by the IPTokenStaking contract.
type IPTokenStakingMinStakeAmountSetIterator struct {
	Event *IPTokenStakingMinStakeAmountSet // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingMinStakeAmountSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingMinStakeAmountSet)
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
		it.Event = new(IPTokenStakingMinStakeAmountSet)
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
func (it *IPTokenStakingMinStakeAmountSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingMinStakeAmountSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingMinStakeAmountSet represents a MinStakeAmountSet event raised by the IPTokenStaking contract.
type IPTokenStakingMinStakeAmountSet struct {
	MinStakeAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterMinStakeAmountSet is a free log retrieval operation binding the contract event 0xea095c2fea861b87f0fd54d0d4453358692a527e120df22b62c71696247dfb9f.
//
// Solidity: event MinStakeAmountSet(uint256 minStakeAmount)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterMinStakeAmountSet(opts *bind.FilterOpts) (*IPTokenStakingMinStakeAmountSetIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "MinStakeAmountSet")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingMinStakeAmountSetIterator{contract: _IPTokenStaking.contract, event: "MinStakeAmountSet", logs: logs, sub: sub}, nil
}

// WatchMinStakeAmountSet is a free log subscription operation binding the contract event 0xea095c2fea861b87f0fd54d0d4453358692a527e120df22b62c71696247dfb9f.
//
// Solidity: event MinStakeAmountSet(uint256 minStakeAmount)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchMinStakeAmountSet(opts *bind.WatchOpts, sink chan<- *IPTokenStakingMinStakeAmountSet) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "MinStakeAmountSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingMinStakeAmountSet)
				if err := _IPTokenStaking.contract.UnpackLog(event, "MinStakeAmountSet", log); err != nil {
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

// ParseMinStakeAmountSet is a log parse operation binding the contract event 0xea095c2fea861b87f0fd54d0d4453358692a527e120df22b62c71696247dfb9f.
//
// Solidity: event MinStakeAmountSet(uint256 minStakeAmount)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseMinStakeAmountSet(log types.Log) (*IPTokenStakingMinStakeAmountSet, error) {
	event := new(IPTokenStakingMinStakeAmountSet)
	if err := _IPTokenStaking.contract.UnpackLog(event, "MinStakeAmountSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingMinUnstakeAmountSetIterator is returned from FilterMinUnstakeAmountSet and is used to iterate over the raw logs and unpacked data for MinUnstakeAmountSet events raised by the IPTokenStaking contract.
type IPTokenStakingMinUnstakeAmountSetIterator struct {
	Event *IPTokenStakingMinUnstakeAmountSet // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingMinUnstakeAmountSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingMinUnstakeAmountSet)
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
		it.Event = new(IPTokenStakingMinUnstakeAmountSet)
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
func (it *IPTokenStakingMinUnstakeAmountSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingMinUnstakeAmountSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingMinUnstakeAmountSet represents a MinUnstakeAmountSet event raised by the IPTokenStaking contract.
type IPTokenStakingMinUnstakeAmountSet struct {
	MinUnstakeAmount *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterMinUnstakeAmountSet is a free log retrieval operation binding the contract event 0xf93d77980ae5a1ddd008d6a7f02cbee5af2a4fcea850c4b55828de4f644e589f.
//
// Solidity: event MinUnstakeAmountSet(uint256 minUnstakeAmount)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterMinUnstakeAmountSet(opts *bind.FilterOpts) (*IPTokenStakingMinUnstakeAmountSetIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "MinUnstakeAmountSet")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingMinUnstakeAmountSetIterator{contract: _IPTokenStaking.contract, event: "MinUnstakeAmountSet", logs: logs, sub: sub}, nil
}

// WatchMinUnstakeAmountSet is a free log subscription operation binding the contract event 0xf93d77980ae5a1ddd008d6a7f02cbee5af2a4fcea850c4b55828de4f644e589f.
//
// Solidity: event MinUnstakeAmountSet(uint256 minUnstakeAmount)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchMinUnstakeAmountSet(opts *bind.WatchOpts, sink chan<- *IPTokenStakingMinUnstakeAmountSet) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "MinUnstakeAmountSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingMinUnstakeAmountSet)
				if err := _IPTokenStaking.contract.UnpackLog(event, "MinUnstakeAmountSet", log); err != nil {
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

// ParseMinUnstakeAmountSet is a log parse operation binding the contract event 0xf93d77980ae5a1ddd008d6a7f02cbee5af2a4fcea850c4b55828de4f644e589f.
//
// Solidity: event MinUnstakeAmountSet(uint256 minUnstakeAmount)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseMinUnstakeAmountSet(log types.Log) (*IPTokenStakingMinUnstakeAmountSet, error) {
	event := new(IPTokenStakingMinUnstakeAmountSet)
	if err := _IPTokenStaking.contract.UnpackLog(event, "MinUnstakeAmountSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the IPTokenStaking contract.
type IPTokenStakingOwnershipTransferStartedIterator struct {
	Event *IPTokenStakingOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingOwnershipTransferStarted)
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
		it.Event = new(IPTokenStakingOwnershipTransferStarted)
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
func (it *IPTokenStakingOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the IPTokenStaking contract.
type IPTokenStakingOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IPTokenStakingOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingOwnershipTransferStartedIterator{contract: _IPTokenStaking.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *IPTokenStakingOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingOwnershipTransferStarted)
				if err := _IPTokenStaking.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_IPTokenStaking *IPTokenStakingFilterer) ParseOwnershipTransferStarted(log types.Log) (*IPTokenStakingOwnershipTransferStarted, error) {
	event := new(IPTokenStakingOwnershipTransferStarted)
	if err := _IPTokenStaking.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the IPTokenStaking contract.
type IPTokenStakingOwnershipTransferredIterator struct {
	Event *IPTokenStakingOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingOwnershipTransferred)
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
		it.Event = new(IPTokenStakingOwnershipTransferred)
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
func (it *IPTokenStakingOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingOwnershipTransferred represents a OwnershipTransferred event raised by the IPTokenStaking contract.
type IPTokenStakingOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IPTokenStakingOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingOwnershipTransferredIterator{contract: _IPTokenStaking.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IPTokenStakingOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingOwnershipTransferred)
				if err := _IPTokenStaking.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_IPTokenStaking *IPTokenStakingFilterer) ParseOwnershipTransferred(log types.Log) (*IPTokenStakingOwnershipTransferred, error) {
	event := new(IPTokenStakingOwnershipTransferred)
	if err := _IPTokenStaking.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingRedelegateIterator is returned from FilterRedelegate and is used to iterate over the raw logs and unpacked data for Redelegate events raised by the IPTokenStaking contract.
type IPTokenStakingRedelegateIterator struct {
	Event *IPTokenStakingRedelegate // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingRedelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingRedelegate)
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
		it.Event = new(IPTokenStakingRedelegate)
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
func (it *IPTokenStakingRedelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingRedelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingRedelegate represents a Redelegate event raised by the IPTokenStaking contract.
type IPTokenStakingRedelegate struct {
	DelegatorUncmpPubkey    []byte
	ValidatorUncmpSrcPubkey []byte
	ValidatorUncmpDstPubkey []byte
	DelegationId            *big.Int
	OperatorAddress         common.Address
	Amount                  *big.Int
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterRedelegate is a free log retrieval operation binding the contract event 0x210091050fbe3add6ade45436b6c7aed210ef28fc37e1a1775970fc391272fe8.
//
// Solidity: event Redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 delegationId, address operatorAddress, uint256 amount)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterRedelegate(opts *bind.FilterOpts) (*IPTokenStakingRedelegateIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "Redelegate")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingRedelegateIterator{contract: _IPTokenStaking.contract, event: "Redelegate", logs: logs, sub: sub}, nil
}

// WatchRedelegate is a free log subscription operation binding the contract event 0x210091050fbe3add6ade45436b6c7aed210ef28fc37e1a1775970fc391272fe8.
//
// Solidity: event Redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 delegationId, address operatorAddress, uint256 amount)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchRedelegate(opts *bind.WatchOpts, sink chan<- *IPTokenStakingRedelegate) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "Redelegate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingRedelegate)
				if err := _IPTokenStaking.contract.UnpackLog(event, "Redelegate", log); err != nil {
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

// ParseRedelegate is a log parse operation binding the contract event 0x210091050fbe3add6ade45436b6c7aed210ef28fc37e1a1775970fc391272fe8.
//
// Solidity: event Redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 delegationId, address operatorAddress, uint256 amount)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseRedelegate(log types.Log) (*IPTokenStakingRedelegate, error) {
	event := new(IPTokenStakingRedelegate)
	if err := _IPTokenStaking.contract.UnpackLog(event, "Redelegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingRemoveOperatorIterator is returned from FilterRemoveOperator and is used to iterate over the raw logs and unpacked data for RemoveOperator events raised by the IPTokenStaking contract.
type IPTokenStakingRemoveOperatorIterator struct {
	Event *IPTokenStakingRemoveOperator // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingRemoveOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingRemoveOperator)
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
		it.Event = new(IPTokenStakingRemoveOperator)
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
func (it *IPTokenStakingRemoveOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingRemoveOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingRemoveOperator represents a RemoveOperator event raised by the IPTokenStaking contract.
type IPTokenStakingRemoveOperator struct {
	UncmpPubkey []byte
	Operator    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRemoveOperator is a free log retrieval operation binding the contract event 0x65729f64aec4981a7e5cedc9abbed98ce4ee8a5c6ecefc35e32d646d51718042.
//
// Solidity: event RemoveOperator(bytes uncmpPubkey, address operator)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterRemoveOperator(opts *bind.FilterOpts) (*IPTokenStakingRemoveOperatorIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "RemoveOperator")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingRemoveOperatorIterator{contract: _IPTokenStaking.contract, event: "RemoveOperator", logs: logs, sub: sub}, nil
}

// WatchRemoveOperator is a free log subscription operation binding the contract event 0x65729f64aec4981a7e5cedc9abbed98ce4ee8a5c6ecefc35e32d646d51718042.
//
// Solidity: event RemoveOperator(bytes uncmpPubkey, address operator)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchRemoveOperator(opts *bind.WatchOpts, sink chan<- *IPTokenStakingRemoveOperator) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "RemoveOperator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingRemoveOperator)
				if err := _IPTokenStaking.contract.UnpackLog(event, "RemoveOperator", log); err != nil {
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

// ParseRemoveOperator is a log parse operation binding the contract event 0x65729f64aec4981a7e5cedc9abbed98ce4ee8a5c6ecefc35e32d646d51718042.
//
// Solidity: event RemoveOperator(bytes uncmpPubkey, address operator)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseRemoveOperator(log types.Log) (*IPTokenStakingRemoveOperator, error) {
	event := new(IPTokenStakingRemoveOperator)
	if err := _IPTokenStaking.contract.UnpackLog(event, "RemoveOperator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingSetRewardAddressIterator is returned from FilterSetRewardAddress and is used to iterate over the raw logs and unpacked data for SetRewardAddress events raised by the IPTokenStaking contract.
type IPTokenStakingSetRewardAddressIterator struct {
	Event *IPTokenStakingSetRewardAddress // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingSetRewardAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingSetRewardAddress)
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
		it.Event = new(IPTokenStakingSetRewardAddress)
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
func (it *IPTokenStakingSetRewardAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingSetRewardAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingSetRewardAddress represents a SetRewardAddress event raised by the IPTokenStaking contract.
type IPTokenStakingSetRewardAddress struct {
	DelegatorUncmpPubkey []byte
	ExecutionAddress     [32]byte
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetRewardAddress is a free log retrieval operation binding the contract event 0x28c0529db8cf660d5b4c1e4b9313683fa7241c3fc49452e7d0ebae215a5f84b2.
//
// Solidity: event SetRewardAddress(bytes delegatorUncmpPubkey, bytes32 executionAddress)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterSetRewardAddress(opts *bind.FilterOpts) (*IPTokenStakingSetRewardAddressIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "SetRewardAddress")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingSetRewardAddressIterator{contract: _IPTokenStaking.contract, event: "SetRewardAddress", logs: logs, sub: sub}, nil
}

// WatchSetRewardAddress is a free log subscription operation binding the contract event 0x28c0529db8cf660d5b4c1e4b9313683fa7241c3fc49452e7d0ebae215a5f84b2.
//
// Solidity: event SetRewardAddress(bytes delegatorUncmpPubkey, bytes32 executionAddress)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchSetRewardAddress(opts *bind.WatchOpts, sink chan<- *IPTokenStakingSetRewardAddress) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "SetRewardAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingSetRewardAddress)
				if err := _IPTokenStaking.contract.UnpackLog(event, "SetRewardAddress", log); err != nil {
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

// ParseSetRewardAddress is a log parse operation binding the contract event 0x28c0529db8cf660d5b4c1e4b9313683fa7241c3fc49452e7d0ebae215a5f84b2.
//
// Solidity: event SetRewardAddress(bytes delegatorUncmpPubkey, bytes32 executionAddress)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseSetRewardAddress(log types.Log) (*IPTokenStakingSetRewardAddress, error) {
	event := new(IPTokenStakingSetRewardAddress)
	if err := _IPTokenStaking.contract.UnpackLog(event, "SetRewardAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingSetWithdrawalAddressIterator is returned from FilterSetWithdrawalAddress and is used to iterate over the raw logs and unpacked data for SetWithdrawalAddress events raised by the IPTokenStaking contract.
type IPTokenStakingSetWithdrawalAddressIterator struct {
	Event *IPTokenStakingSetWithdrawalAddress // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingSetWithdrawalAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingSetWithdrawalAddress)
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
		it.Event = new(IPTokenStakingSetWithdrawalAddress)
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
func (it *IPTokenStakingSetWithdrawalAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingSetWithdrawalAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingSetWithdrawalAddress represents a SetWithdrawalAddress event raised by the IPTokenStaking contract.
type IPTokenStakingSetWithdrawalAddress struct {
	DelegatorUncmpPubkey []byte
	ExecutionAddress     [32]byte
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetWithdrawalAddress is a free log retrieval operation binding the contract event 0x9f7f04f688298f474ed4c786abb29e0ca0173d70516d55d9eac515609b45fbca.
//
// Solidity: event SetWithdrawalAddress(bytes delegatorUncmpPubkey, bytes32 executionAddress)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterSetWithdrawalAddress(opts *bind.FilterOpts) (*IPTokenStakingSetWithdrawalAddressIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "SetWithdrawalAddress")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingSetWithdrawalAddressIterator{contract: _IPTokenStaking.contract, event: "SetWithdrawalAddress", logs: logs, sub: sub}, nil
}

// WatchSetWithdrawalAddress is a free log subscription operation binding the contract event 0x9f7f04f688298f474ed4c786abb29e0ca0173d70516d55d9eac515609b45fbca.
//
// Solidity: event SetWithdrawalAddress(bytes delegatorUncmpPubkey, bytes32 executionAddress)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchSetWithdrawalAddress(opts *bind.WatchOpts, sink chan<- *IPTokenStakingSetWithdrawalAddress) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "SetWithdrawalAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingSetWithdrawalAddress)
				if err := _IPTokenStaking.contract.UnpackLog(event, "SetWithdrawalAddress", log); err != nil {
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

// ParseSetWithdrawalAddress is a log parse operation binding the contract event 0x9f7f04f688298f474ed4c786abb29e0ca0173d70516d55d9eac515609b45fbca.
//
// Solidity: event SetWithdrawalAddress(bytes delegatorUncmpPubkey, bytes32 executionAddress)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseSetWithdrawalAddress(log types.Log) (*IPTokenStakingSetWithdrawalAddress, error) {
	event := new(IPTokenStakingSetWithdrawalAddress)
	if err := _IPTokenStaking.contract.UnpackLog(event, "SetWithdrawalAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingUnjailIterator is returned from FilterUnjail and is used to iterate over the raw logs and unpacked data for Unjail events raised by the IPTokenStaking contract.
type IPTokenStakingUnjailIterator struct {
	Event *IPTokenStakingUnjail // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingUnjailIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingUnjail)
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
		it.Event = new(IPTokenStakingUnjail)
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
func (it *IPTokenStakingUnjailIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingUnjailIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingUnjail represents a Unjail event raised by the IPTokenStaking contract.
type IPTokenStakingUnjail struct {
	Unjailer             common.Address
	ValidatorUncmpPubkey []byte
	Data                 []byte
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterUnjail is a free log retrieval operation binding the contract event 0x026c2e156478ec2a25ccebac97a338d301f69b6d5aeec39c578b28a95e118201.
//
// Solidity: event Unjail(address unjailer, bytes validatorUncmpPubkey, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterUnjail(opts *bind.FilterOpts) (*IPTokenStakingUnjailIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "Unjail")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingUnjailIterator{contract: _IPTokenStaking.contract, event: "Unjail", logs: logs, sub: sub}, nil
}

// WatchUnjail is a free log subscription operation binding the contract event 0x026c2e156478ec2a25ccebac97a338d301f69b6d5aeec39c578b28a95e118201.
//
// Solidity: event Unjail(address unjailer, bytes validatorUncmpPubkey, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchUnjail(opts *bind.WatchOpts, sink chan<- *IPTokenStakingUnjail) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "Unjail")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingUnjail)
				if err := _IPTokenStaking.contract.UnpackLog(event, "Unjail", log); err != nil {
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

// ParseUnjail is a log parse operation binding the contract event 0x026c2e156478ec2a25ccebac97a338d301f69b6d5aeec39c578b28a95e118201.
//
// Solidity: event Unjail(address unjailer, bytes validatorUncmpPubkey, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseUnjail(log types.Log) (*IPTokenStakingUnjail, error) {
	event := new(IPTokenStakingUnjail)
	if err := _IPTokenStaking.contract.UnpackLog(event, "Unjail", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingUpdateValidatorCommssionIterator is returned from FilterUpdateValidatorCommssion and is used to iterate over the raw logs and unpacked data for UpdateValidatorCommssion events raised by the IPTokenStaking contract.
type IPTokenStakingUpdateValidatorCommssionIterator struct {
	Event *IPTokenStakingUpdateValidatorCommssion // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingUpdateValidatorCommssionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingUpdateValidatorCommssion)
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
		it.Event = new(IPTokenStakingUpdateValidatorCommssion)
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
func (it *IPTokenStakingUpdateValidatorCommssionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingUpdateValidatorCommssionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingUpdateValidatorCommssion represents a UpdateValidatorCommssion event raised by the IPTokenStaking contract.
type IPTokenStakingUpdateValidatorCommssion struct {
	ValidatorUncmpPubkey []byte
	CommissionRate       uint32
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterUpdateValidatorCommssion is a free log retrieval operation binding the contract event 0x202c9aad6965f28c0ce1cd00460c1adfa2c90277f4f0a7abb813e2f04cecd70b.
//
// Solidity: event UpdateValidatorCommssion(bytes validatorUncmpPubkey, uint32 commissionRate)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterUpdateValidatorCommssion(opts *bind.FilterOpts) (*IPTokenStakingUpdateValidatorCommssionIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "UpdateValidatorCommssion")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingUpdateValidatorCommssionIterator{contract: _IPTokenStaking.contract, event: "UpdateValidatorCommssion", logs: logs, sub: sub}, nil
}

// WatchUpdateValidatorCommssion is a free log subscription operation binding the contract event 0x202c9aad6965f28c0ce1cd00460c1adfa2c90277f4f0a7abb813e2f04cecd70b.
//
// Solidity: event UpdateValidatorCommssion(bytes validatorUncmpPubkey, uint32 commissionRate)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchUpdateValidatorCommssion(opts *bind.WatchOpts, sink chan<- *IPTokenStakingUpdateValidatorCommssion) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "UpdateValidatorCommssion")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingUpdateValidatorCommssion)
				if err := _IPTokenStaking.contract.UnpackLog(event, "UpdateValidatorCommssion", log); err != nil {
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

// ParseUpdateValidatorCommssion is a log parse operation binding the contract event 0x202c9aad6965f28c0ce1cd00460c1adfa2c90277f4f0a7abb813e2f04cecd70b.
//
// Solidity: event UpdateValidatorCommssion(bytes validatorUncmpPubkey, uint32 commissionRate)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseUpdateValidatorCommssion(log types.Log) (*IPTokenStakingUpdateValidatorCommssion, error) {
	event := new(IPTokenStakingUpdateValidatorCommssion)
	if err := _IPTokenStaking.contract.UnpackLog(event, "UpdateValidatorCommssion", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the IPTokenStaking contract.
type IPTokenStakingWithdrawIterator struct {
	Event *IPTokenStakingWithdraw // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingWithdraw)
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
		it.Event = new(IPTokenStakingWithdraw)
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
func (it *IPTokenStakingWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingWithdraw represents a Withdraw event raised by the IPTokenStaking contract.
type IPTokenStakingWithdraw struct {
	DelegatorUncmpPubkey []byte
	ValidatorUnCmpPubkey []byte
	StakeAmount          *big.Int
	DelegationId         *big.Int
	OperatorAddress      common.Address
	Data                 []byte
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xac41e6ee15d2d0047feb1ea8aba74b92c0334cd3e78024a5ad679d7d08b8fbc5.
//
// Solidity: event Withdraw(bytes delegatorUncmpPubkey, bytes validatorUnCmpPubkey, uint256 stakeAmount, uint256 delegationId, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterWithdraw(opts *bind.FilterOpts) (*IPTokenStakingWithdrawIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingWithdrawIterator{contract: _IPTokenStaking.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xac41e6ee15d2d0047feb1ea8aba74b92c0334cd3e78024a5ad679d7d08b8fbc5.
//
// Solidity: event Withdraw(bytes delegatorUncmpPubkey, bytes validatorUnCmpPubkey, uint256 stakeAmount, uint256 delegationId, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *IPTokenStakingWithdraw) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingWithdraw)
				if err := _IPTokenStaking.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0xac41e6ee15d2d0047feb1ea8aba74b92c0334cd3e78024a5ad679d7d08b8fbc5.
//
// Solidity: event Withdraw(bytes delegatorUncmpPubkey, bytes validatorUnCmpPubkey, uint256 stakeAmount, uint256 delegationId, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseWithdraw(log types.Log) (*IPTokenStakingWithdraw, error) {
	event := new(IPTokenStakingWithdraw)
	if err := _IPTokenStaking.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
