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
	AccessManager       common.Address
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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"stakingRounding\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"defaultMinUnjailFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_MIN_UNJAIL_FEE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"STAKE_ROUNDING\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addOperator\",\"inputs\":[{\"name\":\"uncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createValidator\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"moniker\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionChangeRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isLocked\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"createValidatorOnBehalf\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"moniker\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionChangeRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"isLocked\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"args\",\"type\":\"tuple\",\"internalType\":\"structIIPTokenStaking.InitializerArgs\",\"components\":[{\"name\":\"accessManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minUnstakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minCommissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"shortStakingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"mediumStakingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"longStakingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unjailFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minCommissionRate\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minStakeAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minUnstakeAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"redelegate\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpSrcPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpDstPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"removeOperator\",\"inputs\":[{\"name\":\"uncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"roundedStakeAmount\",\"inputs\":[{\"name\":\"rawAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainder\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setMinCommissionRate\",\"inputs\":[{\"name\":\"newValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinStakeAmount\",\"inputs\":[{\"name\":\"newMinStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinUnstakeAmount\",\"inputs\":[{\"name\":\"newMinUnstakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStakingPeriods\",\"inputs\":[{\"name\":\"short\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"medium\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"long\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnjailFee\",\"inputs\":[{\"name\":\"newUnjailFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWithdrawalAddress\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"newWithdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"stakingPeriod\",\"type\":\"uint8\",\"internalType\":\"enumIIPTokenStaking.StakingPeriod\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeOnBehalf\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"stakingPeriod\",\"type\":\"uint8\",\"internalType\":\"enumIIPTokenStaking.StakingPeriod\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unjail\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unjailFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unjailOnBehalf\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unstake\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstakeOnBehalf\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AddOperator\",\"inputs\":[{\"name\":\"uncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CreateValidator\",\"inputs\":[{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"moniker\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"stakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"maxCommissionRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"maxCommissionChangeRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"isLocked\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorUnCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"stakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"stakingPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinCommissionRateChanged\",\"inputs\":[{\"name\":\"minCommissionRate\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinStakeAmountSet\",\"inputs\":[{\"name\":\"minStakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinUnstakeAmountSet\",\"inputs\":[{\"name\":\"minUnstakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Redelegate\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpSrcPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorUncmpDstPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoveOperator\",\"inputs\":[{\"name\":\"uncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetWithdrawalAddress\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"executionAddress\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakingPeriodsChanged\",\"inputs\":[{\"name\":\"short\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"medium\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"long\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unjail\",\"inputs\":[{\"name\":\"unjailer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"validatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnjailFeeSet\",\"inputs\":[{\"name\":\"newUnjailFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"delegatorUncmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorUnCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"stakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalAddressChangeIntervalSet\",\"inputs\":[{\"name\":\"newInterval\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"IPTokenStaking__CommissionRateOverMax\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__CommissionRateUnderMin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__FailedRemainerRefund\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InsufficientFee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidDefaultMinUnjailFee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidMinUnjailFee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidPubkeyDerivedAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidPubkeyLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__InvalidPubkeyPrefix\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__LowUnstakeAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__MediumLongerThanLong\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ShortPeriodLongerThanMedium\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__StakeAmountUnderMin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroMinCommissionRate\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroMinStakeAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroMinUnstakeAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroShortPeriodDuration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IPTokenStaking__ZeroStakingRounding\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162003ab538038062003ab583398181016040528101906200003791906200024b565b6000820362000072576040517fb23dd7a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8160808181525050633b9aca00811015620000b9576040517f622370a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060a08181525050620000d1620000d960201b60201c565b5050620002d4565b6000620000eb620001e360201b60201c565b90508060000160089054906101000a900460ff161562000137576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff80168160000160009054906101000a900467ffffffffffffffff1667ffffffffffffffff1614620001e05767ffffffffffffffff8160000160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d267ffffffffffffffff604051620001d79190620002b7565b60405180910390a15b50565b60007ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b600080fd5b6000819050919050565b620002258162000210565b81146200023157600080fd5b50565b60008151905062000245816200021a565b92915050565b600080604083850312156200026557620002646200020b565b5b6000620002758582860162000234565b9250506020620002888582860162000234565b9150509250929050565b600067ffffffffffffffff82169050919050565b620002b18162000292565b82525050565b6000602082019050620002ce6000830184620002a6565b92915050565b60805160a05161379f62000316600039600081816118c80152611ff00152600081816118100152818161183701528181611c3b0152611cee015261379f6000f3fe6080604052600436106101cd5760003560e01c80638740597a116100f7578063d2e1f5b811610095578063eb4af04511610064578063eb4af045146105d5578063f1887684146105fe578063f2fde38b14610629578063f9550a8d14610652576101cd565b8063d2e1f5b814610518578063d6d7566014610556578063e30c39781461057f578063ead71c10146105aa576101cd565b8063a0284f16116100d1578063a0284f161461046b578063ab8870f61461049b578063b2bc29ef146104c4578063bda16b15146104ed576101cd565b80638740597a146104085780638da5cb5b146104245780638ed65fbc1461044f576101cd565b806339ec4df91161016f578063715018a61161013e578063715018a614610395578063787f82c8146103ac57806379ba5097146103d557806386eb5e48146103ec576101cd565b806339ec4df9146102f55780633dd9fb9a1461032057806365ffee06146103505780636ea3a2281461036c576101cd565b80630c863f77116101ab5780630c863f771461024d5780631487153e1461027657806317e42e12146102a15780632801f1ec146102ca576101cd565b8063014e8178146101d2578063057b9296146101fb5780630745031a14610224575b600080fd5b3480156101de57600080fd5b506101f960048036038101906101f49190612947565b61066e565b005b34801561020757600080fd5b50610222600480360381019061021d9190612a81565b610745565b005b34801561023057600080fd5b5061024b60048036038101906102469190612b06565b6108b2565b005b34801561025957600080fd5b50610274600480360381019061026f9190612b34565b610ad0565b005b34801561028257600080fd5b5061028b610ae4565b6040516102989190612b70565b60405180910390f35b3480156102ad57600080fd5b506102c860048036038101906102c39190612a81565b610aea565b005b3480156102d657600080fd5b506102df610c57565b6040516102ec9190612b70565b60405180910390f35b34801561030157600080fd5b5061030a610c5d565b6040516103179190612b70565b60405180910390f35b61033a60048036038101906103359190612bb0565b610c63565b6040516103479190612b70565b60405180910390f35b61036a60048036038101906103659190612c79565b610dbe565b005b34801561037857600080fd5b50610393600480360381019061038e9190612b34565b610f85565b005b3480156103a157600080fd5b506103aa610f99565b005b3480156103b857600080fd5b506103d360048036038101906103ce9190612a81565b610fad565b005b3480156103e157600080fd5b506103ea611133565b005b61040660048036038101906104019190612d42565b6111c2565b005b610422600480360381019061041d9190612e8d565b6112a2565b005b34801561043057600080fd5b50610439611440565b6040516104469190612f9f565b60405180910390f35b61046960048036038101906104649190612d42565b611478565b005b61048560048036038101906104809190612bb0565b6115c8565b6040516104929190612b70565b60405180910390f35b3480156104a757600080fd5b506104c260048036038101906104bd9190612b34565b6116b3565b005b3480156104d057600080fd5b506104eb60048036038101906104e69190612947565b6116c7565b005b3480156104f957600080fd5b5061050261180e565b60405161050f9190612b70565b60405180910390f35b34801561052457600080fd5b5061053f600480360381019061053a9190612b34565b611832565b60405161054d929190612fba565b60405180910390f35b34801561056257600080fd5b5061057d60048036038101906105789190612fe3565b611876565b005b34801561058b57600080fd5b5061059461188e565b6040516105a19190612f9f565b60405180910390f35b3480156105b657600080fd5b506105bf6118c6565b6040516105cc9190612b70565b60405180910390f35b3480156105e157600080fd5b506105fc60048036038101906105f79190612b34565b6118ea565b005b34801561060a57600080fd5b506106136118fe565b6040516106209190612b70565b60405180910390f35b34801561063557600080fd5b50610650600480360381019061064b9190613036565b611904565b005b61066c60048036038101906106679190612e8d565b6119c0565b005b8787604182829050036106ad576040517fffc9acd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460f81b828260008181106106c6576106c5613063565b5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603610729576040517f395e38cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107398a8a8a8a8a8a8a8a611aee565b50505050505050505050565b82823360418383905003610785576040517fffc9acd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460f81b8383600081811061079e5761079d613063565b5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603610801576040517f395e38cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166108228484611b7b565b73ffffffffffffffffffffffffffffffffffffffff160361086f576040517ff78f17c100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f6ac365cf05479bb8a295fbf9637875411d6d6f2a0ac7c4b1f560cedcf1a330818686866040516108a2939291906130f0565b60405180910390a1505050505050565b60006108bc611bb1565b905060008160000160089054906101000a900460ff1615905060008260000160009054906101000a900467ffffffffffffffff1690506000808267ffffffffffffffff1614801561090a5750825b9050600060018367ffffffffffffffff1614801561093f575060003073ffffffffffffffffffffffffffffffffffffffff163b145b90508115801561094d575080155b15610984576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018560000160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555083156109d45760018560000160086101000a81548160ff0219169083151502179055505b6109dc611bd9565b6109f78660000160208101906109f29190613036565b611beb565b610a048660200135611bff565b610a118660400135611cb2565b610a1e8660600135611d65565b610a5f866080016020810190610a349190613122565b8760a0016020810190610a479190613122565b8860c0016020810190610a5a9190613122565b611de0565b610a6c8660e00135611fee565b8315610ac85760008560000160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d26001604051610abf91906131a8565b60405180910390a15b505050505050565b610ad8612089565b610ae181611fee565b50565b60005481565b82823360418383905003610b2a576040517fffc9acd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460f81b83836000818110610b4357610b42613063565b5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603610ba6576040517f395e38cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff16610bc78484611b7b565b73ffffffffffffffffffffffffffffffffffffffff1603610c14576040517ff78f17c100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f65729f64aec4981a7e5cedc9abbed98ce4ee8a5c6ecefc35e32d646d51718042868686604051610c47939291906130f0565b60405180910390a1505050505050565b60045481565b60025481565b600087873360418383905003610ca5576040517fffc9acd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460f81b83836000818110610cbe57610cbd613063565b5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603610d21576040517f395e38cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff16610d428484611b7b565b73ffffffffffffffffffffffffffffffffffffffff1603610d8f576040517ff78f17c100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610d97612110565b610da68b8b8b8b8b8b8b612167565b9350610db06122f2565b505050979650505050505050565b86863360418383905003610dfe576040517fffc9acd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460f81b83836000818110610e1757610e16613063565b5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603610e7a576040517f395e38cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff16610e9b8484611b7b565b73ffffffffffffffffffffffffffffffffffffffff1603610ee8576040517ff78f17c100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080610ef434611832565b91509150600154821015610f34576040517fda15b66c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7fb025fa2a574dd306182c6ac63bf7b05482b99680c1b38a42d8401a0adfd3775a8c8c8c8c8c8c8c604051610f6f97969594939291906131c3565b60405180910390a1505050505050505050505050565b610f8d612089565b610f9681611cb2565b50565b610fa1612089565b610fab600061230b565b565b82823360418383905003610fed576040517fffc9acd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460f81b8383600081811061100657611005613063565b5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603611069576040517f395e38cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff1661108a8484611b7b565b73ffffffffffffffffffffffffffffffffffffffff16036110d7576040517ff78f17c100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f9f7f04f688298f474ed4c786abb29e0ca0173d70516d55d9eac515609b45fbca86868673ffffffffffffffffffffffffffffffffffffffff1660001b6040516111239392919061323c565b60405180910390a1505050505050565b600061113d61234b565b90508073ffffffffffffffffffffffffffffffffffffffff1661115e61188e565b73ffffffffffffffffffffffffffffffffffffffff16146111b657806040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016111ad9190612f9f565b60405180910390fd5b6111bf8161230b565b50565b6111ca612110565b838360418282905003611209576040517fffc9acd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460f81b8282600081811061122257611221613063565b5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603611285576040517f395e38cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6112923487878787612353565b505061129c6122f2565b50505050565b898933604183839050036112e2576040517fffc9acd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460f81b838360008181106112fb576112fa613063565b5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19160361135e576040517f395e38cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff1661137f8484611b7b565b73ffffffffffffffffffffffffffffffffffffffff16036113cc576040517ff78f17c100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6113d4612110565b6114298d8d8d8d8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050508c8c8c8c8c8c61241c565b6114316122f2565b50505050505050505050505050565b60008061144b6125a5565b90508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b838333604183839050036114b8576040517fffc9acd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460f81b838360008181106114d1576114d0613063565b5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603611534576040517f395e38cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166115558484611b7b565b73ffffffffffffffffffffffffffffffffffffffff16036115a2576040517ff78f17c100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6115aa612110565b6115b73488888888612353565b6115bf6122f2565b50505050505050565b6000878760418282905003611609576040517fffc9acd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460f81b8282600081811061162257611621613063565b5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603611685576040517f395e38cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61168d612110565b61169c8a8a8a8a8a8a8a612167565b92506116a66122f2565b5050979650505050505050565b6116bb612089565b6116c481611d65565b50565b87873360418383905003611707576040517fffc9acd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460f81b838360008181106117205761171f613063565b5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603611783576040517f395e38cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166117a48484611b7b565b73ffffffffffffffffffffffffffffffffffffffff16036117f1576040517ff78f17c100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6118018b8b8b8b8b8b8b8b611aee565b5050505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000807f000000000000000000000000000000000000000000000000000000000000000083611861919061329d565b9050808361186f91906132fd565b9150915091565b61187e612089565b611889838383611de0565b505050565b6000806118996125cd565b90508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b7f000000000000000000000000000000000000000000000000000000000000000081565b6118f2612089565b6118fb81611bff565b50565b60015481565b61190c612089565b60006119166125cd565b9050818160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff1661197a611440565b73ffffffffffffffffffffffffffffffffffffffff167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b8989604182829050036119ff576040517fffc9acd800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460f81b82826000818110611a1857611a17613063565b5b9050013560f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191603611a7b576040517f395e38cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611a83612110565b611ad88c8c8c8c8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050508b8b8b8b8b8b61241c565b611ae06122f2565b505050505050505050505050565b600254831015611b2a576040517f23870ab900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7fac41e6ee15d2d0047feb1ea8aba74b92c0334cd3e78024a5ad679d7d08b8fbc5888888888789338989604051611b6999989796959493929190613331565b60405180910390a15050505050505050565b600082826001908092611b90939291906133b7565b604051611b9e929190613422565b604051809103902060001c905092915050565b60007ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b611be16125f5565b611be9612635565b565b611bf36125f5565b611bfc81612656565b50565b60008103611c39576040517ff4d335c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000081611c65919061329d565b81611c7091906132fd565b6001819055507fea095c2fea861b87f0fd54d0d4453358692a527e120df22b62c71696247dfb9f600154604051611ca79190612b70565b60405180910390a150565b60008103611cec576040517f8d04d54400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000081611d18919061329d565b81611d2391906132fd565b6002819055507ff93d77980ae5a1ddd008d6a7f02cbee5af2a4fcea850c4b55828de4f644e589f600254604051611d5a9190612b70565b60405180910390a150565b60008103611d9f576040517f23cf9ec000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806000819055507f4167b1de65292a9ff628c9136823791a1de701e1fbdda4863ce22a1cfaf4d0f781604051611dd59190612b70565b60405180910390a150565b60008363ffffffff1603611e20576040517fd8daa8cc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8163ffffffff168363ffffffff1610611e65576040517fc5c0381600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8063ffffffff168263ffffffff1610611eaa576040517fb8e74f7800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826005600060016003811115611ec357611ec261343b565b5b6003811115611ed557611ed461343b565b5b815260200190815260200160002060006101000a81548163ffffffff021916908363ffffffff160217905550816005600060026003811115611f1a57611f1961343b565b5b6003811115611f2c57611f2b61343b565b5b815260200190815260200160002060006101000a81548163ffffffff021916908363ffffffff1602179055508060056000600380811115611f7057611f6f61343b565b5b6003811115611f8257611f8161343b565b5b815260200190815260200160002060006101000a81548163ffffffff021916908363ffffffff1602179055507fa5790d6f3c39faf4bb9bf83076f4b9aeb8c509b3892a128081246ab871e6de06838383604051611fe193929190613479565b60405180910390a1505050565b7f0000000000000000000000000000000000000000000000000000000000000000811015612048576040517f53c11b3b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806004819055507feac81de2f20162b0540ca5d3f43896af15b471a55729ff0c000e611d8b2723638160405161207e9190612b70565b60405180910390a150565b61209161234b565b73ffffffffffffffffffffffffffffffffffffffff166120af611440565b73ffffffffffffffffffffffffffffffffffffffff161461210e576120d261234b565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016121059190612f9f565b60405180910390fd5b565b600061211a6126dc565b9050600281600001540361215a576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002816000018190555050565b600080600061217534611832565b915091506001548210156121b5576040517fda15b66c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600060038111156121cc576121cb61343b565b5b8860038111156121df576121de61343b565b5b14612246576003600081546121f3906134b0565b9190508190559050600560008960038111156122125761221161343b565b5b60038111156122245761222361343b565b5b815260200190815260200160002060009054906101000a900463ffffffff1691505b7f269a32ff589c9b701f49ab6aa532ee8f55901df71a7fca2d70dc9f45314f1be38c8c8c8c888787338f8f6040516122879a99989796959493929190613529565b60405180910390a1600073ffffffffffffffffffffffffffffffffffffffff166108fc859081150290604051600060405180830381858888f193505050501580156122d6573d6000803e3d6000fd5b506122e083612704565b80945050505050979650505050505050565b60006122fc6126dc565b90506001816000018190555050565b60006123156125cd565b90508060000160006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055612347826127ab565b5050565b600033905090565b600454851461238e576040517fb887b2ec00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff166108fc869081150290604051600060405180830381858888f193505050501580156123d5573d6000803e3d6000fd5b507f026c2e156478ec2a25ccebac97a338d301f69b6d5aeec39c578b28a95e118201338585858560405161240d9594939291906135b3565b60405180910390a15050505050565b60008061242834611832565b91509150600154821015612468576040517fda15b66c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000548863ffffffff1610156124aa576040517f183785b600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8663ffffffff168863ffffffff1611156124f0576040517f809afa6400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f19350505050158015612537573d6000803e3d6000fd5b507f6eed9f531c68bead1bd40fc83b45bcfc0dafe6c487ce6e3870025761fe7e811d8b8b8b858c8c8c8c61256c57600061256f565b60015b8c8c6040516125879a99989796959493929190613697565b60405180910390a161259881612704565b5050505050505050505050565b60007f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300905090565b60007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00905090565b6125fd612882565b612633576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b61263d6125f5565b60006126476126dc565b90506001816000018190555050565b61265e6125f5565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036126d05760006040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016126c79190612f9f565b60405180910390fd5b6126d98161230b565b50565b60007f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00905090565b60003373ffffffffffffffffffffffffffffffffffffffff168260405161272a90613754565b60006040518083038185875af1925050503d8060008114612767576040519150601f19603f3d011682016040523d82523d6000602084013e61276c565b606091505b50509050806127a7576040517ffc0ea4f400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60006127b56125a5565b905060008160000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050828260000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3505050565b600061288c611bb1565b60000160089054906101000a900460ff16905090565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f8401126128d1576128d06128ac565b5b8235905067ffffffffffffffff8111156128ee576128ed6128b1565b5b60208301915083600182028301111561290a576129096128b6565b5b9250929050565b6000819050919050565b61292481612911565b811461292f57600080fd5b50565b6000813590506129418161291b565b92915050565b60008060008060008060008060a0898b031215612967576129666128a2565b5b600089013567ffffffffffffffff811115612985576129846128a7565b5b6129918b828c016128bb565b9850985050602089013567ffffffffffffffff8111156129b4576129b36128a7565b5b6129c08b828c016128bb565b965096505060406129d38b828c01612932565b94505060606129e48b828c01612932565b935050608089013567ffffffffffffffff811115612a0557612a046128a7565b5b612a118b828c016128bb565b92509250509295985092959890939650565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000612a4e82612a23565b9050919050565b612a5e81612a43565b8114612a6957600080fd5b50565b600081359050612a7b81612a55565b92915050565b600080600060408486031215612a9a57612a996128a2565b5b600084013567ffffffffffffffff811115612ab857612ab76128a7565b5b612ac4868287016128bb565b93509350506020612ad786828701612a6c565b9150509250925092565b600080fd5b60006101008284031215612afd57612afc612ae1565b5b81905092915050565b60006101008284031215612b1d57612b1c6128a2565b5b6000612b2b84828501612ae6565b91505092915050565b600060208284031215612b4a57612b496128a2565b5b6000612b5884828501612932565b91505092915050565b612b6a81612911565b82525050565b6000602082019050612b856000830184612b61565b92915050565b60048110612b9857600080fd5b50565b600081359050612baa81612b8b565b92915050565b60008060008060008060006080888a031215612bcf57612bce6128a2565b5b600088013567ffffffffffffffff811115612bed57612bec6128a7565b5b612bf98a828b016128bb565b9750975050602088013567ffffffffffffffff811115612c1c57612c1b6128a7565b5b612c288a828b016128bb565b95509550506040612c3b8a828b01612b9b565b935050606088013567ffffffffffffffff811115612c5c57612c5b6128a7565b5b612c688a828b016128bb565b925092505092959891949750929550565b60008060008060008060006080888a031215612c9857612c976128a2565b5b600088013567ffffffffffffffff811115612cb657612cb56128a7565b5b612cc28a828b016128bb565b9750975050602088013567ffffffffffffffff811115612ce557612ce46128a7565b5b612cf18a828b016128bb565b9550955050604088013567ffffffffffffffff811115612d1457612d136128a7565b5b612d208a828b016128bb565b93509350506060612d338a828b01612932565b91505092959891949750929550565b60008060008060408587031215612d5c57612d5b6128a2565b5b600085013567ffffffffffffffff811115612d7a57612d796128a7565b5b612d86878288016128bb565b9450945050602085013567ffffffffffffffff811115612da957612da86128a7565b5b612db5878288016128bb565b925092505092959194509250565b60008083601f840112612dd957612dd86128ac565b5b8235905067ffffffffffffffff811115612df657612df56128b1565b5b602083019150836001820283011115612e1257612e116128b6565b5b9250929050565b600063ffffffff82169050919050565b612e3281612e19565b8114612e3d57600080fd5b50565b600081359050612e4f81612e29565b92915050565b60008115159050919050565b612e6a81612e55565b8114612e7557600080fd5b50565b600081359050612e8781612e61565b92915050565b60008060008060008060008060008060e08b8d031215612eb057612eaf6128a2565b5b60008b013567ffffffffffffffff811115612ece57612ecd6128a7565b5b612eda8d828e016128bb565b9a509a505060208b013567ffffffffffffffff811115612efd57612efc6128a7565b5b612f098d828e01612dc3565b98509850506040612f1c8d828e01612e40565b9650506060612f2d8d828e01612e40565b9550506080612f3e8d828e01612e40565b94505060a0612f4f8d828e01612e78565b93505060c08b013567ffffffffffffffff811115612f7057612f6f6128a7565b5b612f7c8d828e016128bb565b92509250509295989b9194979a5092959850565b612f9981612a43565b82525050565b6000602082019050612fb46000830184612f90565b92915050565b6000604082019050612fcf6000830185612b61565b612fdc6020830184612b61565b9392505050565b600080600060608486031215612ffc57612ffb6128a2565b5b600061300a86828701612e40565b935050602061301b86828701612e40565b925050604061302c86828701612e40565b9150509250925092565b60006020828403121561304c5761304b6128a2565b5b600061305a84828501612a6c565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082825260208201905092915050565b82818337600083830152505050565b6000601f19601f8301169050919050565b60006130cf8385613092565b93506130dc8385846130a3565b6130e5836130b2565b840190509392505050565b6000604082019050818103600083015261310b8185876130c3565b905061311a6020830184612f90565b949350505050565b600060208284031215613138576131376128a2565b5b600061314684828501612e40565b91505092915050565b6000819050919050565b600067ffffffffffffffff82169050919050565b6000819050919050565b600061319261318d6131888461314f565b61316d565b613159565b9050919050565b6131a281613177565b82525050565b60006020820190506131bd6000830184613199565b92915050565b600060808201905081810360008301526131de81898b6130c3565b905081810360208301526131f38187896130c3565b905081810360408301526132088185876130c3565b90506132176060830184612b61565b98975050505050505050565b6000819050919050565b61323681613223565b82525050565b600060408201905081810360008301526132578185876130c3565b9050613266602083018461322d565b949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006132a882612911565b91506132b383612911565b9250826132c3576132c261326e565b5b828206905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061330882612911565b915061331383612911565b925082820390508181111561332b5761332a6132ce565b5b92915050565b600060c082019050818103600083015261334c818b8d6130c3565b9050818103602083015261336181898b6130c3565b90506133706040830188612b61565b61337d6060830187612b61565b61338a6080830186612f90565b81810360a083015261339d8184866130c3565b90509a9950505050505050505050565b600080fd5b600080fd5b600080858511156133cb576133ca6133ad565b5b838611156133dc576133db6133b2565b5b6001850283019150848603905094509492505050565b600081905092915050565b600061340983856133f2565b93506134168385846130a3565b82840190509392505050565b600061342f8284866133fd565b91508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61347381612e19565b82525050565b600060608201905061348e600083018661346a565b61349b602083018561346a565b6134a8604083018461346a565b949350505050565b60006134bb82612911565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036134ed576134ec6132ce565b5b600182019050919050565b600061351361350e61350984612e19565b61316d565b612911565b9050919050565b613523816134f8565b82525050565b600060e0820190508181036000830152613544818c8e6130c3565b90508181036020830152613559818a8c6130c3565b90506135686040830189612b61565b613575606083018861351a565b6135826080830187612b61565b61358f60a0830186612f90565b81810360c08301526135a28184866130c3565b90509b9a5050505050505050505050565b60006060820190506135c86000830188612f90565b81810360208301526135db8186886130c3565b905081810360408301526135f08184866130c3565b90509695505050505050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561363657808201518184015260208101905061361b565b60008484015250505050565b600061364d826135fc565b6136578185613607565b9350613667818560208601613618565b613670816130b2565b840191505092915050565b600060ff82169050919050565b6136918161367b565b82525050565b60006101008201905081810360008301526136b3818c8e6130c3565b905081810360208301526136c7818b613642565b90506136d6604083018a612b61565b6136e3606083018961346a565b6136f0608083018861346a565b6136fd60a083018761346a565b61370a60c0830186613688565b81810360e083015261371d8184866130c3565b90509b9a5050505050505050505050565b50565b600061373e6000836133f2565b91506137498261372e565b600082019050919050565b600061375f82613731565b915081905091905056fea2646970667358221220423dc2f75949268a2b2d73de9db72d3d414414b4783b4e679926ffeec7f5b21064736f6c63430008170033",
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
// Solidity: function createValidator(bytes validatorUncmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool isLocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) CreateValidator(opts *bind.TransactOpts, validatorUncmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, isLocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "createValidator", validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, isLocked, data)
}

// CreateValidator is a paid mutator transaction binding the contract method 0x8740597a.
//
// Solidity: function createValidator(bytes validatorUncmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool isLocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) CreateValidator(validatorUncmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, isLocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.CreateValidator(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, isLocked, data)
}

// CreateValidator is a paid mutator transaction binding the contract method 0x8740597a.
//
// Solidity: function createValidator(bytes validatorUncmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool isLocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) CreateValidator(validatorUncmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, isLocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.CreateValidator(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, isLocked, data)
}

// CreateValidatorOnBehalf is a paid mutator transaction binding the contract method 0xf9550a8d.
//
// Solidity: function createValidatorOnBehalf(bytes validatorUncmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool isLocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) CreateValidatorOnBehalf(opts *bind.TransactOpts, validatorUncmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, isLocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "createValidatorOnBehalf", validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, isLocked, data)
}

// CreateValidatorOnBehalf is a paid mutator transaction binding the contract method 0xf9550a8d.
//
// Solidity: function createValidatorOnBehalf(bytes validatorUncmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool isLocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) CreateValidatorOnBehalf(validatorUncmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, isLocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.CreateValidatorOnBehalf(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, isLocked, data)
}

// CreateValidatorOnBehalf is a paid mutator transaction binding the contract method 0xf9550a8d.
//
// Solidity: function createValidatorOnBehalf(bytes validatorUncmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool isLocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) CreateValidatorOnBehalf(validatorUncmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, isLocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.CreateValidatorOnBehalf(&_IPTokenStaking.TransactOpts, validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, isLocked, data)
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

// Redelegate is a paid mutator transaction binding the contract method 0x65ffee06.
//
// Solidity: function redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) Redelegate(opts *bind.TransactOpts, delegatorUncmpPubkey []byte, validatorUncmpSrcPubkey []byte, validatorUncmpDstPubkey []byte, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "redelegate", delegatorUncmpPubkey, validatorUncmpSrcPubkey, validatorUncmpDstPubkey, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x65ffee06.
//
// Solidity: function redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) Redelegate(delegatorUncmpPubkey []byte, validatorUncmpSrcPubkey []byte, validatorUncmpDstPubkey []byte, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Redelegate(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpSrcPubkey, validatorUncmpDstPubkey, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x65ffee06.
//
// Solidity: function redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) Redelegate(delegatorUncmpPubkey []byte, validatorUncmpSrcPubkey []byte, validatorUncmpDstPubkey []byte, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Redelegate(&_IPTokenStaking.TransactOpts, delegatorUncmpPubkey, validatorUncmpSrcPubkey, validatorUncmpDstPubkey, amount)
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
	IsLocked                uint8
	Data                    []byte
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterCreateValidator is a free log retrieval operation binding the contract event 0x6eed9f531c68bead1bd40fc83b45bcfc0dafe6c487ce6e3870025761fe7e811d.
//
// Solidity: event CreateValidator(bytes validatorUncmpPubkey, string moniker, uint256 stakeAmount, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, uint8 isLocked, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterCreateValidator(opts *bind.FilterOpts) (*IPTokenStakingCreateValidatorIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "CreateValidator")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingCreateValidatorIterator{contract: _IPTokenStaking.contract, event: "CreateValidator", logs: logs, sub: sub}, nil
}

// WatchCreateValidator is a free log subscription operation binding the contract event 0x6eed9f531c68bead1bd40fc83b45bcfc0dafe6c487ce6e3870025761fe7e811d.
//
// Solidity: event CreateValidator(bytes validatorUncmpPubkey, string moniker, uint256 stakeAmount, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, uint8 isLocked, bytes data)
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

// ParseCreateValidator is a log parse operation binding the contract event 0x6eed9f531c68bead1bd40fc83b45bcfc0dafe6c487ce6e3870025761fe7e811d.
//
// Solidity: event CreateValidator(bytes validatorUncmpPubkey, string moniker, uint256 stakeAmount, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, uint8 isLocked, bytes data)
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
	Amount                  *big.Int
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterRedelegate is a free log retrieval operation binding the contract event 0xb025fa2a574dd306182c6ac63bf7b05482b99680c1b38a42d8401a0adfd3775a.
//
// Solidity: event Redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 amount)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterRedelegate(opts *bind.FilterOpts) (*IPTokenStakingRedelegateIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "Redelegate")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingRedelegateIterator{contract: _IPTokenStaking.contract, event: "Redelegate", logs: logs, sub: sub}, nil
}

// WatchRedelegate is a free log subscription operation binding the contract event 0xb025fa2a574dd306182c6ac63bf7b05482b99680c1b38a42d8401a0adfd3775a.
//
// Solidity: event Redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 amount)
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

// ParseRedelegate is a log parse operation binding the contract event 0xb025fa2a574dd306182c6ac63bf7b05482b99680c1b38a42d8401a0adfd3775a.
//
// Solidity: event Redelegate(bytes delegatorUncmpPubkey, bytes validatorUncmpSrcPubkey, bytes validatorUncmpDstPubkey, uint256 amount)
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
