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
	Owner               common.Address
	MinStakeAmount      *big.Int
	MinUnstakeAmount    *big.Int
	MinCommissionRate   *big.Int
	ShortStakingPeriod  uint32
	MediumStakingPeriod uint32
	LongStakingPeriod   uint32
	UnjailFee           *big.Int
}

// IPTokenStakingMetaData contains all meta data concerning the IPTokenStaking contract.
var IPTokenStakingMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"stakingRounding\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"defaultMinUnjailFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_MIN_UNJAIL_FEE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"STAKE_ROUNDING\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addOperator\",\"inputs\":[{\"name\":\"uncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createValidator\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"moniker\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionChangeRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"supportsUnlocked\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"createValidatorOnBehalf\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"moniker\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionChangeRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"supportsUnlocked\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"args\",\"type\":\"tuple\",\"internalType\":\"structIIPTokenStaking.InitializerArgs\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minUnstakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minCommissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"shortStakingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mediumStakingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"longStakingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unjailFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minCommissionRate\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minStakeAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minUnstakeAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"redelegate\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpSrcPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpDstPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"removeOperator\",\"inputs\":[{\"name\":\"uncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"roundedStakeAmount\",\"inputs\":[{\"name\":\"rawAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainder\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setMinCommissionRate\",\"inputs\":[{\"name\":\"newValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinStakeAmount\",\"inputs\":[{\"name\":\"newMinStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinUnstakeAmount\",\"inputs\":[{\"name\":\"newMinUnstakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRewardsAddress\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"newRewardsAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStakingPeriods\",\"inputs\":[{\"name\":\"short\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"medium\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"long\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnjailFee\",\"inputs\":[{\"name\":\"newUnjailFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWithdrawalAddress\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"newWithdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"stakingPeriod\",\"type\":\"uint8\",\"internalType\":\"enumIIPTokenStaking.StakingPeriod\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeOnBehalf\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"stakingPeriod\",\"type\":\"uint8\",\"internalType\":\"enumIIPTokenStaking.StakingPeriod\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakingDurations\",\"inputs\":[{\"name\":\"period\",\"type\":\"uint8\",\"internalType\":\"enumIIPTokenStaking.StakingPeriod\"}],\"outputs\":[{\"name\":\"duration\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unjail\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unjailFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unjailOnBehalf\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unstake\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstakeOnBehalf\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AddOperator\",\"inputs\":[{\"name\":\"uncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CreateValidator\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"moniker\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"stakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"maxCommissionRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"maxCommissionChangeRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"supportsUnlocked\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorUnCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"stakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"stakingPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinCommissionRateChanged\",\"inputs\":[{\"name\":\"minCommissionRate\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinStakeAmountSet\",\"inputs\":[{\"name\":\"minStakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinUnstakeAmountSet\",\"inputs\":[{\"name\":\"minUnstakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Redelegate\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpSrcPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpDstPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoveOperator\",\"inputs\":[{\"name\":\"uncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetRewardAddress\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"executionAddress\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetWithdrawalAddress\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"executionAddress\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakingPeriodsChanged\",\"inputs\":[{\"name\":\"short\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"medium\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"long\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unjail\",\"inputs\":[{\"name\":\"unjailer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnjailFeeSet\",\"inputs\":[{\"name\":\"newUnjailFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorUnCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"stakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalAddressChangeIntervalSet\",\"inputs\":[{\"name\":\"newInterval\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"IPTokenStaking__CommissionRateOverMax\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__CommissionRateUnderMin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__FailedRemainerRefund\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidDefaultMinUnjailFee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidFeeAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidMinUnjailFee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidPubkeyDerivedAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidPubkeyLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidPubkeyPrefix\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__LowUnstakeAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__MediumLongerThanLong\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__RedelegatingToSameValidator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ShortPeriodLongerThanMedium\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__StakeAmountUnderMin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroMinCommissionRate\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroMinStakeAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroMinUnstakeAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroShortPeriodDuration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroStakingRounding\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60c0346200016c57620029e6906001600160401b0390601f38849003908101601f1916820190838211838310176200017157808391604096879485528339810103126200016c57602081519101519080156200015b57608052633b9aca0081106200014a5760a0527ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff82851c1662000139578080831603620000f4575b835161285e90816200018882396080518181816104bf01528181610512015281816109a401528181611dd9015281816120d90152818161216d015261250b015260a0518181816103dc01526123b70152f35b6001600160401b0319909116811790915581519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a1388080620000a2565b835163f92ee8a960e01b8152600490fd5b825163622370a960e01b8152600490fd5b835163591eebd160e11b8152600490fd5b600080fd5b634e487b7160e01b600052604160045260246000fdfe60806040908082526004918236101561001757600080fd5b600091823560e01c908163014e81781461183d57508063057b92961461177a5780630745031a146115075780630c863f77146114e25780631487153e146114c657806317e42e121461138e5780632801f1ec14611370578063346cc7271461133957806339ec4df91461131b5780633dd9fb9a1461124d5780636ea3a22814611228578063715018a614611161578063787f82c81461109757806379ba50971461100c57806386eb5e4814610f925780638740597a14610ec75780638da5cb5b14610e745780638ed65fbc14610d385780639d04b12114610bc05780639d9d293f14610863578063a0284f1614610738578063ab8870f614610713578063b2bc29ef14610535578063bda16b15146104fb578063d2e1f5b8146104a1578063d6d7566014610452578063e30c3978146103ff578063ead71c10146103c5578063eb4af0451461039d578063f18876841461037f578063f2fde38b146102a95763f9550a8d1461018557600080fd5b61018e36611b16565b9960418993999a929a98949897959703610281578215610255577f04000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000008335160361022d57506102069a9b506102016124a7565b611dbd565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580f35b8c90517f395e38cb000000000000000000000000000000000000000000000000000000008152fd5b60248c60328f7f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b8c90517fffc9acd8000000000000000000000000000000000000000000000000000000008152fd5b82843461037b57602060031936011261037b573573ffffffffffffffffffffffffffffffffffffffff808216809203610377576102e4612437565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227008380a380f35b8280fd5b5080fd5b503461037b578160031936011261037b576020906001549051908152f35b82843461037b57602060031936011261037b576103c2906103bc612437565b356120a7565b80f35b503461037b578160031936011261037b57602090517f00000000000000000000000000000000000000000000000000000000000000008152f35b503461037b578160031936011261037b5760209073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0054169051908152f35b82843461037b57606060031936011261037b573563ffffffff808216820361037757602435818116810361049d57604435918216820361049d576103c292610498612437565b61222a565b8380fd5b509190346104f85760206003193601126104f85750356104eb6104e47f000000000000000000000000000000000000000000000000000000000000000083611fd8565b8092612011565b9082519182526020820152f35b80fd5b503461037b578160031936011261037b57602090517f00000000000000000000000000000000000000000000000000000000000000008152f35b50913461037b576105453661195a565b959697604188969394959296036106ec5787156106c0577f04000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000008a351603610699573373ffffffffffffffffffffffffffffffffffffffff6105c68a8c61204d565b160361067257600254841061064b5750916106296106459694928b969461061b7fac41e6ee15d2d0047feb1ea8aba74b92c0334cd3e78024a5ad679d7d08b8fbc59c9d519b8c9b60c08d5260c08d0191611bb3565b918a830360208c0152611bb3565b94870152606086015233608086015284830360a0860152611bb3565b0390a180f35b8a517f23870ab9000000000000000000000000000000000000000000000000000000008152fd5b8a517ff78f17c1000000000000000000000000000000000000000000000000000000008152fd5b8a517f395e38cb000000000000000000000000000000000000000000000000000000008152fd5b8960326024927f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b8a517fffc9acd8000000000000000000000000000000000000000000000000000000008152fd5b82843461037b57602060031936011261037b576103c290610732612437565b356121c9565b509061074336611a60565b98909660418694969793970361083b57831561081057507f04000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000833516036107e95750916107be95939160209895936107b96124a7565b612501565b9060017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005551908152f35b87517f395e38cb000000000000000000000000000000000000000000000000000000008152fd5b60326024927f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b5087517fffc9acd8000000000000000000000000000000000000000000000000000000008152fd5b509160a060031936011261037b5767ffffffffffffffff813581811161049d576108909036908401611927565b919092602435828111610bbc576108aa9036908301611927565b92604435908111610bb8576108c29036908401611927565b93909260418603610b91578515610b65577f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000081818a351603610b15573373ffffffffffffffffffffffffffffffffffffffff61093f8a8c61204d565b1603610b3d5760418403610aed578315610ac157818186351603610b155760418703610aed578615610ac15785351603610a9a5761097e368385611d8c565b60208151910120610990368787611d8c565b6020815191012014610a73576109cf6109c97f000000000000000000000000000000000000000000000000000000000000000034611fd8565b34612011565b60015411610a4c575091610a29610a3694928994610a1b7fafafb42f2318707386b88d22641806f098c4094bc7a6a8386d7ab5f7beebba019a9b51998a9960a08b5260a08b0191611bb3565b9188830360208a0152611bb3565b9285840390860152611bb3565b606435606083015260843560808301520390a180f35b88517fda15b66c000000000000000000000000000000000000000000000000000000008152fd5b88517f43df0a36000000000000000000000000000000000000000000000000000000008152fd5b88517f395e38cb000000000000000000000000000000000000000000000000000000008152fd5b60248a6032857f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b828b517fffc9acd8000000000000000000000000000000000000000000000000000000008152fd5b828b517f395e38cb000000000000000000000000000000000000000000000000000000008152fd5b828b517ff78f17c1000000000000000000000000000000000000000000000000000000008152fd5b8760326024927f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b88517fffc9acd8000000000000000000000000000000000000000000000000000000008152fd5b8680fd5b8580fd5b503461037b57610bcf366119c8565b9092919460418403610d11578315610ce5577f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000087351603610cbe5773ffffffffffffffffffffffffffffffffffffffff903382610c4d878a61204d565b1603610c975750610c8a7f28c0529db8cf660d5b4c1e4b9313683fa7241c3fc49452e7d0ebae215a5f84b295968451958587968752860191611bb3565b911660208301520390a180f35b83517ff78f17c1000000000000000000000000000000000000000000000000000000008152fd5b82517f395e38cb000000000000000000000000000000000000000000000000000000008152fd5b8460326024927f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b82517fffc9acd8000000000000000000000000000000000000000000000000000000008152fd5b50610d4236611acf565b9360418394929403610e4c578215610e20577f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000083351603610df8573373ffffffffffffffffffffffffffffffffffffffff610dbe858561204d565b1603610dd057506103c2949550611c27565b8690517ff78f17c1000000000000000000000000000000000000000000000000000000008152fd5b8690517f395e38cb000000000000000000000000000000000000000000000000000000008152fd5b6024866032897f4e487b7100000000000000000000000000000000000000000000000000000000835252fd5b8690517fffc9acd8000000000000000000000000000000000000000000000000000000008152fd5b503461037b578160031936011261037b5760209073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054169051908152f35b50610ed136611b16565b9960418993999a929a98949897959703610281578215610255577f04000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000008335160361022d573373ffffffffffffffffffffffffffffffffffffffff610f55858561204d565b1603610f6a57506102069a9b506102016124a7565b8c90517ff78f17c1000000000000000000000000000000000000000000000000000000008152fd5b50610f9c36611acf565b93610fa89391936124a7565b60418303610e4c578215610e20577f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000083351603610df85750610206949550611c27565b508290346103775782600319360112610377573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0054160361106757826103c2336126bd565b6024925051907f118cdaa70000000000000000000000000000000000000000000000000000000082523390820152fd5b503461037b576110a6366119c8565b9092919460418403610d11578315610ce5577f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000087351603610cbe5773ffffffffffffffffffffffffffffffffffffffff903382611124878a61204d565b1603610c975750610c8a7f9f7f04f688298f474ed4c786abb29e0ca0173d70516d55d9eac515609b45fbca95968451958587968752860191611bb3565b82346104f857806003193601126104f85761117a612437565b8073ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b82843461037b57602060031936011261037b576103c290611247612437565b3561213b565b509061125836611a60565b98909660418694969793970361083b57831561081057507f04000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000833516036107e9573373ffffffffffffffffffffffffffffffffffffffff6112d9858561204d565b16036112f45750916107be95939160209895936107b96124a7565b87517ff78f17c1000000000000000000000000000000000000000000000000000000008152fd5b503461037b578160031936011261037b576020906002549051908152f35b5090346104f85760206003193601126104f8578235928310156104f8575063ffffffff611367602093611a18565b54169051908152f35b50823461037757826003193601126103775760209250549051908152f35b50913461037b5761139e366119c8565b9490926041840361149f578315610ce5577f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000084351603611478573373ffffffffffffffffffffffffffffffffffffffff611419868661204d565b160361145157507f65729f64aec4981a7e5cedc9abbed98ce4ee8a5c6ecefc35e32d646d517180429394610645915193849384611bf2565b90517ff78f17c1000000000000000000000000000000000000000000000000000000008152fd5b90517f395e38cb000000000000000000000000000000000000000000000000000000008152fd5b90517fffc9acd8000000000000000000000000000000000000000000000000000000008152fd5b503461037b578160031936011261037b57602091549051908152f35b82843461037b57602060031936011261037b576103c290611501612437565b356123b5565b50823461037757610100600319360112610377577ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff82851c16159167ffffffffffffffff811680159081611772575b6001149081611768575b15908161175f575b50611737578260017fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008316178555611702575b506115aa612771565b6115b2612771565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005580359073ffffffffffffffffffffffffffffffffffffffff8216808303610bb8576115ff612771565b611607612771565b156116d35750611616906126bd565b6116216024356120a7565b61162c60443561213b565b6116376064356121c9565b60843563ffffffff8082168203610bbc5760a4358181168103610bb85760c4359182168203610bb8576116699261222a565b61167460e4356123b5565b61167c578280f35b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d291817fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff602093541690555160018152a181808280f35b602490868651917f1e4fbdf7000000000000000000000000000000000000000000000000000000008352820152fd5b7fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000001668010000000000000001178355856115a1565b5083517ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b9050158761156e565b303b159150611566565b84915061155c565b50913461037b5761178a366119c8565b9490926041840361149f578315610ce5577f04000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000084351603611478573373ffffffffffffffffffffffffffffffffffffffff611805868661204d565b160361145157507f6ac365cf05479bb8a295fbf9637875411d6d6f2a0ac7c4b1f560cedcf1a330819394610645915193849384611bf2565b83858492346103775761184f3661195a565b969890959760418996949695939503611902575087156106c0577f04000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000008a35160361069957600254841061064b5750916106296106459694928b969461061b7fac41e6ee15d2d0047feb1ea8aba74b92c0334cd3e78024a5ad679d7d08b8fbc59c9d519b8c9b60c08d5260c08d0191611bb3565b7fffc9acd8000000000000000000000000000000000000000000000000000000008152fd5b9181601f840112156119555782359167ffffffffffffffff8311611955576020838186019501011161195557565b600080fd5b60a06003198201126119555767ffffffffffffffff90600435828111611955578161198791600401611927565b9390939260243581811161195557836119a291600401611927565b939093926044359260643592608435918211611955576119c491600401611927565b9091565b6040600319820112611955576004359067ffffffffffffffff8211611955576119f391600401611927565b909160243573ffffffffffffffffffffffffffffffffffffffff811681036119555790565b6004811015611a31576000526005602052604060002090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60806003198201126119555767ffffffffffffffff916004358381116119555782611a8d91600401611927565b939093926024358281116119555781611aa891600401611927565b9390939260443560048110156119555792606435918211611955576119c491600401611927565b60406003198201126119555767ffffffffffffffff916004358381116119555782611afc91600401611927565b93909392602435918211611955576119c491600401611927565b9060e06003198301126119555767ffffffffffffffff916004358381116119555781611b4491600401611927565b939093926024358281116119555783611b5f91600401611927565b9093909263ffffffff9160443583811681036119555793606435848116810361195557936084359081168103611955579260a4358015158103611955579260c435918211611955576119c491600401611927565b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0938186528686013760008582860101520116010190565b91611c2060209273ffffffffffffffffffffffffffffffffffffffff92969596604086526040860191611bb3565b9416910152565b91926004543403611cb55760003415611cac575b600080808093813491f115611ca0577f026c2e156478ec2a25ccebac97a338d301f69b6d5aeec39c578b28a95e11820193611c9b91611c8d604051958695338752606060208801526060870191611bb3565b918483036040860152611bb3565b0390a1565b6040513d6000823e3d90fd5b506108fc611c3b565b60046040517f5097ac51000000000000000000000000000000000000000000000000000000008152fd5b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f604051930116820182811067ffffffffffffffff821117611d2357604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b67ffffffffffffffff8111611d2357601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b929192611da0611d9b83611d52565b611cdf565b938285528282011161195557816000926020928387013784010152565b9998959291611dd3919594929897983691611d8c565b91611dfe7f000000000000000000000000000000000000000000000000000000000000000034611fd8565b97611e098934612011565b946001548610611fae57600080549163ffffffff809616928310611f8457851692838311611f5a578180898015611f50575b82809291818093f115611f435715611f3c576001955b611e686040519e8f9a610120808d528c0191611bb3565b916020998a818503910152815191828452815b838110611f2957505094611f1f9894611f279e9f9b979294847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe07f65bfc2fa1cd4c6f50f60983ad1cf1cb4bff5ee6570428254dfce41b085ef6d149f9b968f9e9a9860ff998c82601f940101520116019660408d015260608c015260808b01521660a08901521660c08701523360e087015281868203016101008701520191611bb3565b0390a16127ca565b565b8181018c01518582018d01528b01611e7b565b8095611e51565b50604051903d90823e3d90fd5b6108fc9150611e3b565b60046040517f809afa64000000000000000000000000000000000000000000000000000000008152fd5b60046040517f183785b6000000000000000000000000000000000000000000000000000000008152fd5b60046040517fda15b66c000000000000000000000000000000000000000000000000000000008152fd5b8115611fe2570690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b9190820391821161201e57565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b816001116119555773ffffffffffffffffffffffffffffffffffffffff9161209c9160017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff3693019101611d8c565b602081519101201690565b80156121115760206121047fea095c2fea861b87f0fd54d0d4453358692a527e120df22b62c71696247dfb9f926120fe7f000000000000000000000000000000000000000000000000000000000000000082611fd8565b90612011565b80600155604051908152a1565b60046040517ff4d335c6000000000000000000000000000000000000000000000000000000008152fd5b801561219f5760206121927ff93d77980ae5a1ddd008d6a7f02cbee5af2a4fcea850c4b55828de4f644e589f926120fe7f000000000000000000000000000000000000000000000000000000000000000082611fd8565b80600255604051908152a1565b60046040517f8d04d544000000000000000000000000000000000000000000000000000000008152fd5b8015612200576020817f4167b1de65292a9ff628c9136823791a1de701e1fbdda4863ce22a1cfaf4d0f792600055604051908152a1565b60046040517f23cf9ec0000000000000000000000000000000000000000000000000000000008152fd5b63ffffffff9081169291831561238b57811691828410156123615716918282101561233757600560209081527f1471eb6eb2c5e789fc3de43f8ce62938c7d1836ec861730447e2ada8fd81017b80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000090811684179091557f89832631fb3c3307a103ba2c84ab569c64d6182a18893dcd163f0f1c2090733a805482168517905560036000527fa9bc9a3a348c357ba16b37005d7e6b3236198c0e939f4af8c5f19b8deeb8ebc08054909116851790556040805192835290820192909252908101919091527fa5790d6f3c39faf4bb9bf83076f4b9aeb8c509b3892a128081246ab871e6de0690606090a1565b60046040517fb8e74f78000000000000000000000000000000000000000000000000000000008152fd5b60046040517fc5c03816000000000000000000000000000000000000000000000000000000008152fd5b60046040517fd8daa8cc000000000000000000000000000000000000000000000000000000008152fd5b7f0000000000000000000000000000000000000000000000000000000000000000811061240d576020817feac81de2f20162b0540ca5d3f43896af15b471a55729ff0c000e611d8b27236392600455604051908152a1565b60046040517f53c11b3b000000000000000000000000000000000000000000000000000000008152fd5b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416330361247757565b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0060028154146124d75760029055565b60046040517f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b92909193956125307f000000000000000000000000000000000000000000000000000000000000000034611fd8565b9561253b8734612011565b956001548710611fae57600093849960048110156126905780612609575b50946125e16000989495899893967f269a32ff589c9b701f49ab6aa532ee8f55901df71a7fca2d70dc9f45314f1be39563ffffffff6125bb8c9b9a8c9b6125ad6040519a8b9a60e08c5260e08c0191611bb3565b9189830360208b0152611bb3565b938960408801521660608601528d60808601523360a086015284830360c0860152611bb3565b0390a1818115612600575b8290f115611ca0576125fd906127ca565b90565b506108fc6125ec565b995091949692959093600354907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82146126635750600101806003559861264f90611a18565b5463ffffffff169390959296949138612559565b807f4e487b7100000000000000000000000000000000000000000000000000000000602492526011600452fd5b6024867f4e487b710000000000000000000000000000000000000000000000000000000081526021600452fd5b7fffffffffffffffffffffffff0000000000000000000000000000000000000000907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008281541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549073ffffffffffffffffffffffffffffffffffffffff80931680948316179055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0600080a3565b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c16156127a057565b60046040517fd7e6bcf8000000000000000000000000000000000000000000000000000000008152fd5b600080808093335af13d15612823573d6127e6611d9b82611d52565b908152600060203d92013e5b156127f957565b60046040517ffc0ea4f4000000000000000000000000000000000000000000000000000000008152fd5b6127f256fea2646970667358221220bef586b565ee91ba56ae83b5d516b8e7b2ed68434b99c8812c81b3ef0cff956a64736f6c63430008170033",
}

// IPTokenStakingABI is the input ABI used to generate the binding from.
// Deprecated: Use IPTokenStakingMetaData.ABI instead.
var IPTokenStakingABI = IPTokenStakingMetaData.ABI

// IPTokenStakingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use IPTokenStakingMetaData.Bin instead.
var IPTokenStakingBin = IPTokenStakingMetaData.Bin

// DeployIPTokenStaking deploys a new Ethereum contract, binding an instance of IPTokenStaking to it.
func DeployIPTokenStaking(auth *bind.TransactOpts, backend bind.ContractBackend, stakingRounding *big.Int, defaultMinUnjailFee *big.Int) (common.Address, *types.Transaction, *IPTokenStaking, error) {
	parsed, err := IPTokenStakingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(IPTokenStakingBin), backend, stakingRounding, defaultMinUnjailFee)
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

// DEFAULTMINUNJAILFEE is a free data retrieval call binding the contract method 0xead71c10.
//
// Solidity: function DEFAULT_MIN_UNJAIL_FEE() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) DEFAULTMINUNJAILFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "DEFAULT_MIN_UNJAIL_FEE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEFAULTMINUNJAILFEE is a free data retrieval call binding the contract method 0xead71c10.
//
// Solidity: function DEFAULT_MIN_UNJAIL_FEE() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) DEFAULTMINUNJAILFEE() (*big.Int, error) {
	return _IPTokenStaking.Contract.DEFAULTMINUNJAILFEE(&_IPTokenStaking.CallOpts)
}

// DEFAULTMINUNJAILFEE is a free data retrieval call binding the contract method 0xead71c10.
//
// Solidity: function DEFAULT_MIN_UNJAIL_FEE() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) DEFAULTMINUNJAILFEE() (*big.Int, error) {
	return _IPTokenStaking.Contract.DEFAULTMINUNJAILFEE(&_IPTokenStaking.CallOpts)
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

// StakingDurations is a free data retrieval call binding the contract method 0x346cc727.
//
// Solidity: function stakingDurations(uint8 period) view returns(uint32 duration)
func (_IPTokenStaking *IPTokenStakingCaller) StakingDurations(opts *bind.CallOpts, period uint8) (uint32, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "stakingDurations", period)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// StakingDurations is a free data retrieval call binding the contract method 0x346cc727.
//
// Solidity: function stakingDurations(uint8 period) view returns(uint32 duration)
func (_IPTokenStaking *IPTokenStakingSession) StakingDurations(period uint8) (uint32, error) {
	return _IPTokenStaking.Contract.StakingDurations(&_IPTokenStaking.CallOpts, period)
}

// StakingDurations is a free data retrieval call binding the contract method 0x346cc727.
//
// Solidity: function stakingDurations(uint8 period) view returns(uint32 duration)
func (_IPTokenStaking *IPTokenStakingCallerSession) StakingDurations(period uint8) (uint32, error) {
	return _IPTokenStaking.Contract.StakingDurations(&_IPTokenStaking.CallOpts, period)
}

// UnjailFee is a free data retrieval call binding the contract method 0x2801f1ec.
//
// Solidity: function unjailFee() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) UnjailFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "unjailFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnjailFee is a free data retrieval call binding the contract method 0x2801f1ec.
//
// Solidity: function unjailFee() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) UnjailFee() (*big.Int, error) {
	return _IPTokenStaking.Contract.UnjailFee(&_IPTokenStaking.CallOpts)
}

// UnjailFee is a free data retrieval call binding the contract method 0x2801f1ec.
//
// Solidity: function unjailFee() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) UnjailFee() (*big.Int, error) {
	return _IPTokenStaking.Contract.UnjailFee(&_IPTokenStaking.CallOpts)
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
// Solidity: function addOperator(bytes uncmpPubkey, address operator) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) AddOperator(opts *bind.TransactOpts, uncmpPubkey []byte, operator common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "addOperator", uncmpPubkey, operator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x057b9296.
//
// Solidity: function addOperator(bytes uncmpPubkey, address operator) returns()
func (_IPTokenStaking *IPTokenStakingSession) AddOperator(uncmpPubkey []byte, operator common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.AddOperator(&_IPTokenStaking.TransactOpts, uncmpPubkey, operator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x057b9296.
//
// Solidity: function addOperator(bytes uncmpPubkey, address operator) returns()
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

// Initialize is a paid mutator transaction binding the contract method 0x0745031a.
//
// Solidity: function initialize((address,uint256,uint256,uint256,uint32,uint32,uint32,uint256) args) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) Initialize(opts *bind.TransactOpts, args IIPTokenStakingInitializerArgs) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "initialize", args)
}

// Initialize is a paid mutator transaction binding the contract method 0x0745031a.
//
// Solidity: function initialize((address,uint256,uint256,uint256,uint32,uint32,uint32,uint256) args) returns()
func (_IPTokenStaking *IPTokenStakingSession) Initialize(args IIPTokenStakingInitializerArgs) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Initialize(&_IPTokenStaking.TransactOpts, args)
}

// Initialize is a paid mutator transaction binding the contract method 0x0745031a.
//
// Solidity: function initialize((address,uint256,uint256,uint256,uint32,uint32,uint32,uint256) args) returns()
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
// Solidity: function setRewardsAddress(bytes delegatorUncmpPubkey, address newRewardsAddress) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetRewardsAddress(opts *bind.TransactOpts, delegatorUncmpPubkey []byte, newRewardsAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setRewardsAddress", delegatorUncmpPubkey, newRewardsAddress)
}

// SetRewardsAddress is a paid mutator transaction binding the contract method 0x9d04b121.
//
// Solidity: function setRewardsAddress(bytes delegatorUncmpPubkey, address newRewardsAddress) returns()
func (_IPTokenStaking *IPTokenStakingSession) SetRewardsAddress(delegatorUncmpPubkey []byte, newRewardsAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetRewardsAddress(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, newRewardsAddress)
}

// SetRewardsAddress is a paid mutator transaction binding the contract method 0x9d04b121.
//
// Solidity: function setRewardsAddress(bytes delegatorUncmpPubkey, address newRewardsAddress) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetRewardsAddress(delegatorUncmpPubkey []byte, newRewardsAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetRewardsAddress(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, newRewardsAddress)
}

// SetStakingPeriods is a paid mutator transaction binding the contract method 0xd6d75660.
//
// Solidity: function setStakingPeriods(uint32 short, uint32 medium, uint32 long) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetStakingPeriods(opts *bind.TransactOpts, short uint32, medium uint32, long uint32) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setStakingPeriods", short, medium, long)
}

// SetStakingPeriods is a paid mutator transaction binding the contract method 0xd6d75660.
//
// Solidity: function setStakingPeriods(uint32 short, uint32 medium, uint32 long) returns()
func (_IPTokenStaking *IPTokenStakingSession) SetStakingPeriods(short uint32, medium uint32, long uint32) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetStakingPeriods(&_IPTokenStaking.TransactOpts, short, medium, long)
}

// SetStakingPeriods is a paid mutator transaction binding the contract method 0xd6d75660.
//
// Solidity: function setStakingPeriods(uint32 short, uint32 medium, uint32 long) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetStakingPeriods(short uint32, medium uint32, long uint32) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetStakingPeriods(&_IPTokenStaking.TransactOpts, short, medium, long)
}

// SetUnjailFee is a paid mutator transaction binding the contract method 0x0c863f77.
//
// Solidity: function setUnjailFee(uint256 newUnjailFee) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetUnjailFee(opts *bind.TransactOpts, newUnjailFee *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setUnjailFee", newUnjailFee)
}

// SetUnjailFee is a paid mutator transaction binding the contract method 0x0c863f77.
//
// Solidity: function setUnjailFee(uint256 newUnjailFee) returns()
func (_IPTokenStaking *IPTokenStakingSession) SetUnjailFee(newUnjailFee *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetUnjailFee(&_IPTokenStaking.TransactOpts, newUnjailFee)
}

// SetUnjailFee is a paid mutator transaction binding the contract method 0x0c863f77.
//
// Solidity: function setUnjailFee(uint256 newUnjailFee) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetUnjailFee(newUnjailFee *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetUnjailFee(&_IPTokenStaking.TransactOpts, newUnjailFee)
}

// SetWithdrawalAddress is a paid mutator transaction binding the contract method 0x787f82c8.
//
// Solidity: function setWithdrawalAddress(bytes delegatorUncmpPubkey, address newWithdrawalAddress) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetWithdrawalAddress(opts *bind.TransactOpts, delegatorUncmpPubkey []byte, newWithdrawalAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setWithdrawalAddress", delegatorUncmpPubkey, newWithdrawalAddress)
}

// SetWithdrawalAddress is a paid mutator transaction binding the contract method 0x787f82c8.
//
// Solidity: function setWithdrawalAddress(bytes delegatorUncmpPubkey, address newWithdrawalAddress) returns()
func (_IPTokenStaking *IPTokenStakingSession) SetWithdrawalAddress(delegatorUncmpPubkey []byte, newWithdrawalAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetWithdrawalAddress(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, newWithdrawalAddress)
}

// SetWithdrawalAddress is a paid mutator transaction binding the contract method 0x787f82c8.
//
// Solidity: function setWithdrawalAddress(bytes delegatorUncmpPubkey, address newWithdrawalAddress) returns()
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
	Amount                  *big.Int
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterRedelegate is a free log retrieval operation binding the contract event 0xafafb42f2318707386b88d22641806f098c4094bc7a6a8386d7ab5f7beebba01.
//
// Solidity: event Redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 delegationId, uint256 amount)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterRedelegate(opts *bind.FilterOpts) (*IPTokenStakingRedelegateIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "Redelegate")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingRedelegateIterator{contract: _IPTokenStaking.contract, event: "Redelegate", logs: logs, sub: sub}, nil
}

// WatchRedelegate is a free log subscription operation binding the contract event 0xafafb42f2318707386b88d22641806f098c4094bc7a6a8386d7ab5f7beebba01.
//
// Solidity: event Redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 delegationId, uint256 amount)
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

// ParseRedelegate is a log parse operation binding the contract event 0xafafb42f2318707386b88d22641806f098c4094bc7a6a8386d7ab5f7beebba01.
//
// Solidity: event Redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 delegationId, uint256 amount)
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

// IPTokenStakingStakingPeriodsChangedIterator is returned from FilterStakingPeriodsChanged and is used to iterate over the raw logs and unpacked data for StakingPeriodsChanged events raised by the IPTokenStaking contract.
type IPTokenStakingStakingPeriodsChangedIterator struct {
	Event *IPTokenStakingStakingPeriodsChanged // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingStakingPeriodsChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingStakingPeriodsChanged)
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
		it.Event = new(IPTokenStakingStakingPeriodsChanged)
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
func (it *IPTokenStakingStakingPeriodsChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingStakingPeriodsChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingStakingPeriodsChanged represents a StakingPeriodsChanged event raised by the IPTokenStaking contract.
type IPTokenStakingStakingPeriodsChanged struct {
	Short  uint32
	Medium uint32
	Long   uint32
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterStakingPeriodsChanged is a free log retrieval operation binding the contract event 0xa5790d6f3c39faf4bb9bf83076f4b9aeb8c509b3892a128081246ab871e6de06.
//
// Solidity: event StakingPeriodsChanged(uint32 short, uint32 medium, uint32 long)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterStakingPeriodsChanged(opts *bind.FilterOpts) (*IPTokenStakingStakingPeriodsChangedIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "StakingPeriodsChanged")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingStakingPeriodsChangedIterator{contract: _IPTokenStaking.contract, event: "StakingPeriodsChanged", logs: logs, sub: sub}, nil
}

// WatchStakingPeriodsChanged is a free log subscription operation binding the contract event 0xa5790d6f3c39faf4bb9bf83076f4b9aeb8c509b3892a128081246ab871e6de06.
//
// Solidity: event StakingPeriodsChanged(uint32 short, uint32 medium, uint32 long)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchStakingPeriodsChanged(opts *bind.WatchOpts, sink chan<- *IPTokenStakingStakingPeriodsChanged) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "StakingPeriodsChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingStakingPeriodsChanged)
				if err := _IPTokenStaking.contract.UnpackLog(event, "StakingPeriodsChanged", log); err != nil {
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

// ParseStakingPeriodsChanged is a log parse operation binding the contract event 0xa5790d6f3c39faf4bb9bf83076f4b9aeb8c509b3892a128081246ab871e6de06.
//
// Solidity: event StakingPeriodsChanged(uint32 short, uint32 medium, uint32 long)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseStakingPeriodsChanged(log types.Log) (*IPTokenStakingStakingPeriodsChanged, error) {
	event := new(IPTokenStakingStakingPeriodsChanged)
	if err := _IPTokenStaking.contract.UnpackLog(event, "StakingPeriodsChanged", log); err != nil {
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

// IPTokenStakingUnjailFeeSetIterator is returned from FilterUnjailFeeSet and is used to iterate over the raw logs and unpacked data for UnjailFeeSet events raised by the IPTokenStaking contract.
type IPTokenStakingUnjailFeeSetIterator struct {
	Event *IPTokenStakingUnjailFeeSet // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingUnjailFeeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingUnjailFeeSet)
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
		it.Event = new(IPTokenStakingUnjailFeeSet)
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
func (it *IPTokenStakingUnjailFeeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingUnjailFeeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingUnjailFeeSet represents a UnjailFeeSet event raised by the IPTokenStaking contract.
type IPTokenStakingUnjailFeeSet struct {
	NewUnjailFee *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUnjailFeeSet is a free log retrieval operation binding the contract event 0xeac81de2f20162b0540ca5d3f43896af15b471a55729ff0c000e611d8b272363.
//
// Solidity: event UnjailFeeSet(uint256 newUnjailFee)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterUnjailFeeSet(opts *bind.FilterOpts) (*IPTokenStakingUnjailFeeSetIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "UnjailFeeSet")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingUnjailFeeSetIterator{contract: _IPTokenStaking.contract, event: "UnjailFeeSet", logs: logs, sub: sub}, nil
}

// WatchUnjailFeeSet is a free log subscription operation binding the contract event 0xeac81de2f20162b0540ca5d3f43896af15b471a55729ff0c000e611d8b272363.
//
// Solidity: event UnjailFeeSet(uint256 newUnjailFee)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchUnjailFeeSet(opts *bind.WatchOpts, sink chan<- *IPTokenStakingUnjailFeeSet) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "UnjailFeeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingUnjailFeeSet)
				if err := _IPTokenStaking.contract.UnpackLog(event, "UnjailFeeSet", log); err != nil {
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

// ParseUnjailFeeSet is a log parse operation binding the contract event 0xeac81de2f20162b0540ca5d3f43896af15b471a55729ff0c000e611d8b272363.
//
// Solidity: event UnjailFeeSet(uint256 newUnjailFee)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseUnjailFeeSet(log types.Log) (*IPTokenStakingUnjailFeeSet, error) {
	event := new(IPTokenStakingUnjailFeeSet)
	if err := _IPTokenStaking.contract.UnpackLog(event, "UnjailFeeSet", log); err != nil {
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

// IPTokenStakingWithdrawalAddressChangeIntervalSetIterator is returned from FilterWithdrawalAddressChangeIntervalSet and is used to iterate over the raw logs and unpacked data for WithdrawalAddressChangeIntervalSet events raised by the IPTokenStaking contract.
type IPTokenStakingWithdrawalAddressChangeIntervalSetIterator struct {
	Event *IPTokenStakingWithdrawalAddressChangeIntervalSet // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingWithdrawalAddressChangeIntervalSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingWithdrawalAddressChangeIntervalSet)
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
		it.Event = new(IPTokenStakingWithdrawalAddressChangeIntervalSet)
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
func (it *IPTokenStakingWithdrawalAddressChangeIntervalSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingWithdrawalAddressChangeIntervalSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingWithdrawalAddressChangeIntervalSet represents a WithdrawalAddressChangeIntervalSet event raised by the IPTokenStaking contract.
type IPTokenStakingWithdrawalAddressChangeIntervalSet struct {
	NewInterval *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalAddressChangeIntervalSet is a free log retrieval operation binding the contract event 0xbed33ba1e6aacc702f8e48397b388e43ca92a8898ed8bdb389fd8b18af95d32c.
//
// Solidity: event WithdrawalAddressChangeIntervalSet(uint256 newInterval)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterWithdrawalAddressChangeIntervalSet(opts *bind.FilterOpts) (*IPTokenStakingWithdrawalAddressChangeIntervalSetIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "WithdrawalAddressChangeIntervalSet")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingWithdrawalAddressChangeIntervalSetIterator{contract: _IPTokenStaking.contract, event: "WithdrawalAddressChangeIntervalSet", logs: logs, sub: sub}, nil
}

// WatchWithdrawalAddressChangeIntervalSet is a free log subscription operation binding the contract event 0xbed33ba1e6aacc702f8e48397b388e43ca92a8898ed8bdb389fd8b18af95d32c.
//
// Solidity: event WithdrawalAddressChangeIntervalSet(uint256 newInterval)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchWithdrawalAddressChangeIntervalSet(opts *bind.WatchOpts, sink chan<- *IPTokenStakingWithdrawalAddressChangeIntervalSet) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "WithdrawalAddressChangeIntervalSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingWithdrawalAddressChangeIntervalSet)
				if err := _IPTokenStaking.contract.UnpackLog(event, "WithdrawalAddressChangeIntervalSet", log); err != nil {
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

// ParseWithdrawalAddressChangeIntervalSet is a log parse operation binding the contract event 0xbed33ba1e6aacc702f8e48397b388e43ca92a8898ed8bdb389fd8b18af95d32c.
//
// Solidity: event WithdrawalAddressChangeIntervalSet(uint256 newInterval)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseWithdrawalAddressChangeIntervalSet(log types.Log) (*IPTokenStakingWithdrawalAddressChangeIntervalSet, error) {
	event := new(IPTokenStakingWithdrawalAddressChangeIntervalSet)
	if err := _IPTokenStaking.contract.UnpackLog(event, "WithdrawalAddressChangeIntervalSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
