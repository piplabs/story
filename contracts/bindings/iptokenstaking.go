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
	Owner                    common.Address
	MinCreateValidatorAmount *big.Int
	MinStakeAmount           *big.Int
	MinUnstakeAmount         *big.Int
	MinCommissionRate        *big.Int
	Fee                      *big.Int
}

// IPTokenStakingMetaData contains all meta data concerning the IPTokenStaking contract.
var IPTokenStakingMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"defaultMinFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxDataLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"AA\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"BB\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_MIN_FEE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_DATA_LENGTH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_MONIKER_LENGTH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"STAKE_ROUNDING\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createValidator\",\"inputs\":[{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"moniker\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxCommissionChangeRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"supportsUnlocked\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"args\",\"type\":\"tuple\",\"internalType\":\"structIIPTokenStaking.InitializerArgs\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minCreateValidatorAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minUnstakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minCommissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minCommissionRate\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minCreateValidatorAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minStakeAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minUnstakeAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"redelegate\",\"inputs\":[{\"name\":\"validatorSrcCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorDstCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"redelegateOnBehalf\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"validatorSrcCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validatorDstCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"roundedStakeAmount\",\"inputs\":[{\"name\":\"rawAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainder\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"setFee\",\"inputs\":[{\"name\":\"newFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinCommissionRate\",\"inputs\":[{\"name\":\"newValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinCreateValidatorAmount\",\"inputs\":[{\"name\":\"newMinCreateValidatorAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinStakeAmount\",\"inputs\":[{\"name\":\"newMinStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinUnstakeAmount\",\"inputs\":[{\"name\":\"newMinUnstakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setRewardsAddress\",\"inputs\":[{\"name\":\"newRewardsAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setWithdrawalAddress\",\"inputs\":[{\"name\":\"newWithdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"stakingPeriod\",\"type\":\"uint8\",\"internalType\":\"enumIIPTokenStaking.StakingPeriod\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeOnBehalf\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"stakingPeriod\",\"type\":\"uint8\",\"internalType\":\"enumIIPTokenStaking.StakingPeriod\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unjail\",\"inputs\":[{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unjailOnBehalf\",\"inputs\":[{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unsetOperator\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unstake\",\"inputs\":[{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unstakeOnBehalf\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"updateValidatorCommission\",\"inputs\":[{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"CreateValidator\",\"inputs\":[{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"moniker\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"stakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"maxCommissionRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"maxCommissionChangeRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"supportsUnlocked\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"stakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"stakingPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeSet\",\"inputs\":[{\"name\":\"newFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinCommissionRateChanged\",\"inputs\":[{\"name\":\"minCommissionRate\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinCreateValidatorAmountSet\",\"inputs\":[{\"name\":\"minCreateValidatorAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinStakeAmountSet\",\"inputs\":[{\"name\":\"minStakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinUnstakeAmountSet\",\"inputs\":[{\"name\":\"minUnstakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Redelegate\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"validatorSrcCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"validatorDstCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetOperator\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetRewardAddress\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"executionAddress\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetWithdrawalAddress\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"executionAddress\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unjail\",\"inputs\":[{\"name\":\"unjailer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnsetOperator\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdateValidatorCommission\",\"inputs\":[{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"commissionRate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"validatorCmpPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"stakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"delegationId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operatorAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
	Bin: "0x60c03462000190576200312a906001600160401b0390601f38849003908101601f19168201908382118383101762000194578083916040968794855283398101031262000190576020815191015190633b9aca0081106200013c5760805260a0527ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff82851c166200012b578080831603620000e6575b8351612f819081620001a98239608051818181610eaf0152611e95015260a0518181816106b401528181610fbe0152818161122a015281816118e501528181611a2b01526120ce0152f35b6001600160401b0319909116811790915581519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f80806200009b565b835163f92ee8a960e01b8152600490fd5b835162461bcd60e51b815260206004820152602760248201527f4950546f6b656e5374616b696e673a20496e76616c69642064656661756c74206044820152666d696e2066656560c81b6064820152608490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c806308a3002b1461027457806308e1b4ff1461026f5780631487153e1461026a57806321b8092e14610265578063396e1e471461026057806339ec4df91461025b578063420b72f814610256578063564bcc1f146102515780635727dc5c1461024c57806369fe0e2d146102475780636ea3a22814610242578063715018a61461023d57806377f8af381461023857806379ba50971461023357806386eb5e481461022e5780638740597a1461022957806388138319146102245780638906758d1461021f5780638da5cb5b1461021a5780638ed65fbc1461021557806394fd0fe014610210578063997da8d41461020b578063ab8870f614610206578063b3ab15fb14610201578063b3e85917146101fc578063bda16b15146101f7578063c582db44146101f2578063d2e1f5b8146101ed578063d2e9dedb146101e8578063d6f89acd146101e3578063ddca3f43146101de578063e30c3978146101d9578063e52da4fc146101d4578063e65f593f146101cf578063eb4af045146101ca578063eeeac01e146101c5578063f1887684146101c05763f2fde38b146101bb575f80fd5b611490565b611473565b611439565b611415565b6113f8565b6112d8565b611286565b611269565b611158565b6110f5565b6110cb565b610fff565b610fe1565b610fa7565b610f10565b610eec565b610ed2565b610e98565b610e1e565b610dcc565b610d42565b610b8c565b610afc565b610a60565b61098e565b61091a565b610855565b610831565b61080d565b6107f2565b61075c565b6105cf565b610584565b610569565b6104c2565b610488565b610462565b3461045e5760c060031936011261045e577ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005467ffffffffffffffff60ff8260401c1615911680159081610456575b600114908161044c575b159081610443575b50610419578061032a7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0060017fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000825416179055565b6103bd575b610337611560565b61033d57005b6103897ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a007fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff8154169055565b604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29080602081015b0390a1005b6104147ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00680100000000000000007fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff825416179055565b61032f565b60046040517ff92ee8a9000000000000000000000000000000000000000000000000000000008152fd5b9050155f6102d5565b303b1591506102cd565b8291506102c3565b5f80fd5b3461045e57602060031936011261045e5761047b611f2f565b610486600435611bdd565b005b3461045e575f60031936011261045e5760205f54604051908152f35b73ffffffffffffffffffffffffffffffffffffffff81160361045e57565b602060031936011261045e576004356104da816104a4565b6104e7600454341461163e565b5f3415610560575b5f80808093813491f11561055b576103b873ffffffffffffffffffffffffffffffffffffffff7f1717e61ad5ca09c8dbe5620a60e8907a3124e59063fbcc7f71912e536511b42192166105438115156116ba565b60408051338152602081019290925290918291820190565b6116af565b506108fc6104ef565b3461045e575f60031936011261045e57602060405160468152f35b3461045e575f60031936011261045e576020600254604051908152f35b9181601f8401121561045e5782359167ffffffffffffffff831161045e576020838186019501011161045e57565b60a060031936011261045e576004356105e7816104a4565b67ffffffffffffffff9060243582811161045e576106099036906004016105a1565b916044359360843590811161045e576106269036906004016105a1565b90610634600454341461163e565b5f3415610753575b5f80808093813491f11561055b5761066561066061065b3688886117db565b6125a9565b612793565b610673600354871115611f6f565b61067e606435611bc8565b509260025484106106e9577f89f97919c84baf7c0af918e2aae2581e51a70e1ce774948638769eac6b4d2d1c966103b8946106db7f000000000000000000000000000000000000000000000000000000000000000086111561180b565b604051978897339489611fe0565b608460405162461bcd60e51b815260206004820152602860248201527f4950546f6b656e5374616b696e673a20556e7374616b6520616d6f756e74207560448201527f6e646572206d696e0000000000000000000000000000000000000000000000006064820152fd5b506108fc61063c565b606060031936011261045e5767ffffffffffffffff60043581811161045e576107899036906004016105a1565b6024359291600484101561045e5760443592831161045e576020936107b56107c69436906004016105a1565b9390926107c061203b565b33612095565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055604051908152f35b3461045e575f60031936011261045e57602060405160078152f35b3461045e57602060031936011261045e57610826611f2f565b610486600435611e93565b3461045e57602060031936011261045e5761084a611f2f565b610486600435611d41565b3461045e575f60031936011261045e5761086d611f2f565b5f73ffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffff00000000000000000000000000000000000000007f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008181541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549182169055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a3005b608060031936011261045e57600435610932816104a4565b67ffffffffffffffff60243581811161045e576109539036906004016105a1565b92604435600481101561045e5760643593841161045e5760209461097e6107c69536906004016105a1565b94909361098961203b565b612095565b3461045e575f60031936011261045e573373ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416036109e5576104863361241f565b60246040517f118cdaa7000000000000000000000000000000000000000000000000000000008152336004820152fd5b604060031982011261045e5767ffffffffffffffff9160043583811161045e5782610a42916004016105a1565b9390939260243591821161045e57610a5c916004016105a1565b9091565b610a6936610a15565b91610a77600454341461163e565b5f3415610aa7575b5f80808093813491f11561055b5761048693610aa261066061065b3685856117db565b6118b9565b506108fc610a7f565b6044359063ffffffff8216820361045e57565b6064359063ffffffff8216820361045e57565b6084359063ffffffff8216820361045e57565b6024359063ffffffff8216820361045e57565b60e060031936011261045e5767ffffffffffffffff60043581811161045e57610b299036906004016105a1565b9060243583811161045e57610b429036906004016105a1565b610b4a610ab0565b610b52610ac3565b90610b5b610ad6565b9260a43594851515860361045e5760c43598891161045e57610b846104869936906004016105a1565b9890976119ad565b60a060031936011261045e57600435610ba4816104a4565b67ffffffffffffffff60243581811161045e57610bc59036906004016105a1565b9160443590811161045e57610bde9036906004016105a1565b9060643594610bf0600454341461163e565b5f3415610d39575b5f80808093813491f11561055b57610c1761066061065b3688886117db565b610c2861066061065b3686866117db565b73ffffffffffffffffffffffffffffffffffffffff811615610ccf577f08132c654ea3df406da083159e0fa3b4d3b8738bb81fc2dccefb22757868a122956103b893610c98610c783689896117db565b60208151910120610c8a3684886117db565b602081519101201415612b17565b610ca6600354831115611f6f565b610cb1608435611bc8565b5093610cc1600154861015612321565b604051978897339589612b88565b608460405162461bcd60e51b815260206004820152602160248201527f4950546f6b656e5374616b696e673a20496e76616c69642064656c656761746f60448201527f72000000000000000000000000000000000000000000000000000000000000006064820152fd5b506108fc610bf8565b602060031936011261045e57600435610d5a816104a4565b610d67600454341461163e565b5f3415610dc3575b5f80808093813491f11561055b576103b873ffffffffffffffffffffffffffffffffffffffff7f03564a99f640621f30a0be83f41fa9567acefa525e0132642d50d31c5b7d0e3992166105438115156116ba565b506108fc610d6f565b3461045e575f60031936011261045e57602073ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005416604051908152f35b610e2736610a15565b91610e35600454341461163e565b5f3415610e8f575b5f80808093813491f11561055b5761048693610aa2610e6061065b3685856117db565b610e6981612793565b73ffffffffffffffffffffffffffffffffffffffff610e8833926128de565b161461193c565b506108fc610e3d565b3461045e575f60031936011261045e5760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b3461045e575f60031936011261045e5760206040515f8152f35b3461045e57602060031936011261045e57610f05611f2f565b610486600435611df3565b602060031936011261045e57600435610f28816104a4565b610f35600454341461163e565b5f3415610f9e575b5f80808093813491f11561055b57604073ffffffffffffffffffffffffffffffffffffffff7f2709918445f306d3e94d280907c62c5d2525ac3192d2e544774c7f181d65af3e9216610f908115156116ba565b8151903382526020820152a1005b506108fc610f3d565b3461045e575f60031936011261045e5760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b3461045e575f60031936011261045e576020604051633b9aca008152f35b604060031936011261045e5760043567ffffffffffffffff811161045e5761102b9036906004016105a1565b90611034610ae9565b91611042600454341461163e565b5f34156110c2575b5f80808093813491f11561055b577fd87fbe52b4b01e284c00c7b9a719018ee55b3f823bd39556a77bd14829107d729261108b610e6061065b3685876117db565b6110a063ffffffff5f54921691821015611af7565b6110b760405193849360408552604085019161187b565b9060208301520390a1005b506108fc61104a565b3461045e57602060031936011261045e5760406110e9600435611bc8565b82519182526020820152f35b5f60031936011261045e5761110d600454341461163e565b5f341561114f575b5f80808093813491f11561055b577f5c6b8680e6b38b49c9b4293dc69be9b24e0ddcf054274619bb86e8e749f091a06020604051338152a1005b506108fc611115565b608060031936011261045e5767ffffffffffffffff60043581811161045e576111859036906004016105a1565b6024929192359160643590811161045e576111a49036906004016105a1565b6111b4600494929454341461163e565b5f3415611260575b5f80808093813491f11561055b576111db61066061065b3686896117db565b6111e9600354831115611f6f565b6111f4604435611bc8565b509360025485106106e9577f89f97919c84baf7c0af918e2aae2581e51a70e1ce774948638769eac6b4d2d1c956103b8936112517f000000000000000000000000000000000000000000000000000000000000000085111561180b565b60405196879633933389611fe0565b506108fc6111bc565b3461045e575f60031936011261045e576020600454604051908152f35b3461045e575f60031936011261045e57602073ffffffffffffffffffffffffffffffffffffffff7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c005416604051908152f35b608060031936011261045e5767ffffffffffffffff60043581811161045e576113059036906004016105a1565b909160243590811161045e5761131f9036906004016105a1565b60449291923590611333600454341461163e565b5f34156113ef575b5f80808093813491f11561055b5761135a61066061065b3686896117db565b61136b61066061065b3684886117db565b3315610ccf577f08132c654ea3df406da083159e0fa3b4d3b8738bb81fc2dccefb22757868a122946103b8926113b76113a53687856117db565b60208151910120610c8a36868a6117db565b6113c5600354821115611f6f565b6113d0606435611bc8565b50926113e0600154851015612321565b60405196879633943389612b88565b506108fc61133b565b3461045e575f60031936011261045e576020600554604051908152f35b3461045e57602060031936011261045e5761142e611f2f565b610486600435611c8f565b3461045e575f60031936011261045e5760206040517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f8152f35b3461045e575f60031936011261045e576020600154604051908152f35b3461045e57602060031936011261045e576004356114ad816104a4565b6114b5611f2f565b73ffffffffffffffffffffffffffffffffffffffff809116907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00827fffffffffffffffffffffffff00000000000000000000000000000000000000008254161790557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930054167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227005f80a3005b611568612be3565b611570612be3565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00556004356115a0816104a4565b6115a8612be3565b6115b0612be3565b73ffffffffffffffffffffffffffffffffffffffff81161561160e576115d59061241f565b6115e0602435611bdd565b6115eb604435611c8f565b6115f6606435611d41565b611601608435611df3565b61160c60a435611e93565b565b60246040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081525f6004820152fd5b1561164557565b608460405162461bcd60e51b815260206004820152602260248201527f4950546f6b656e5374616b696e673a20496e76616c69642066656520616d6f7560448201527f6e740000000000000000000000000000000000000000000000000000000000006064820152fd5b6040513d5f823e3d90fd5b156116c157565b608460405162461bcd60e51b815260206004820152602260248201527f4950546f6b656e5374616b696e673a207a65726f20696e70757420616464726560448201527f73730000000000000000000000000000000000000000000000000000000000006064820152fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f604051930116820182811067ffffffffffffffff82111761179c57604052565b61172b565b67ffffffffffffffff811161179c57601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b9291926117ef6117ea836117a1565b611758565b938285528282011161045e57815f926020928387013784010152565b1561181257565b608460405162461bcd60e51b8152602060048201526024808201527f4950546f6b656e5374616b696e673a2044617461206c656e677468206f76657260448201527f206d6178000000000000000000000000000000000000000000000000000000006064820152fd5b601f82602094937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe093818652868601375f8582860101520116010190565b91611937907f026c2e156478ec2a25ccebac97a338d301f69b6d5aeec39c578b28a95e1182019461190c7f000000000000000000000000000000000000000000000000000000000000000082111561180b565b61192960405195869533875260606020880152606087019161187b565b91848303604086015261187b565b0390a1565b1561194357565b608460405162461bcd60e51b815260206004820152603160248201527f536563703235366b3156657269666965723a20496e76616c6964207075626b6560448201527f79206465726976656420616464726573730000000000000000000000000000006064820152fd5b916119da9196999498939795976119cb610e6061065b368b886117db565b6119d361203b565b36916117db565b956119e434611bc8565b9990916119f5600554841015612321565b611a185f548b611a0f63ffffffff80921692831015611af7565b83161015612938565b611a2660468a5111156129a9565b611a527f000000000000000000000000000000000000000000000000000000000000000087111561180b565b825f8115611aee575b5f808093818094f11561055b577f65bfc2fa1cd4c6f50f60983ad1cf1cb4bff5ee6570428254dfce41b085ef6d1499611aa59715611ae7576001935b6040519a8b9a33978c612a1a565b0390a180611ad8575b5061160c60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b611ae190612c3c565b5f611aae565b5f93611a97565b506108fc611a5b565b15611afe57565b608460405162461bcd60e51b815260206004820152602960248201527f4950546f6b656e5374616b696e673a20436f6d6d697373696f6e20726174652060448201527f756e646572206d696e00000000000000000000000000000000000000000000006064820152fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f918203918211611bc357565b611b68565b633b9aca00810690818103908111611bc35791565b633b9aca00810680820391808311611bc3578260055514611c255760207fc58038d152352f59eeece288f00b152052e7ab8b5774c338a33ac4dffde1c37491604051908152a1565b608460405162461bcd60e51b815260206004820152603060248201527f4950546f6b656e5374616b696e673a205a65726f206d696e206372656174652060448201527f76616c696461746f7220616d6f756e74000000000000000000000000000000006064820152fd5b633b9aca00810680820391808311611bc3578260015514611cd75760207fea095c2fea861b87f0fd54d0d4453358692a527e120df22b62c71696247dfb9f91604051908152a1565b608460405162461bcd60e51b815260206004820152602560248201527f4950546f6b656e5374616b696e673a205a65726f206d696e207374616b65206160448201527f6d6f756e740000000000000000000000000000000000000000000000000000006064820152fd5b633b9aca00810680820391808311611bc3578260025514611d895760207ff93d77980ae5a1ddd008d6a7f02cbee5af2a4fcea850c4b55828de4f644e589f91604051908152a1565b608460405162461bcd60e51b815260206004820152602760248201527f4950546f6b656e5374616b696e673a205a65726f206d696e20756e7374616b6560448201527f20616d6f756e74000000000000000000000000000000000000000000000000006064820152fd5b8015611e29576020817f4167b1de65292a9ff628c9136823791a1de701e1fbdda4863ce22a1cfaf4d0f7925f55604051908152a1565b608460405162461bcd60e51b815260206004820152602860248201527f4950546f6b656e5374616b696e673a205a65726f206d696e20636f6d6d69737360448201527f696f6e20726174650000000000000000000000000000000000000000000000006064820152fd5b7f00000000000000000000000000000000000000000000000000000000000000008110611eeb576020817f20461e09b8e557b77e107939f9ce6544698123aad0fc964ac5cc59b7df2e608f92600455604051908152a1565b606460405162461bcd60e51b815260206004820152601f60248201527f4950546f6b656e5374616b696e673a20496e76616c6964206d696e20666565006044820152fd5b73ffffffffffffffffffffffffffffffffffffffff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300541633036109e557565b15611f7657565b608460405162461bcd60e51b815260206004820152602560248201527f4950546f6b656e5374616b696e673a20496e76616c69642064656c656761746960448201527f6f6e2069640000000000000000000000000000000000000000000000000000006064820152fd5b939694909261203898969261201a9173ffffffffffffffffffffffffffffffffffffffff809616875260c0602088015260c087019161187b565b966040850152606084015216608082015260a081850391015261187b565b90565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00600281541461206b5760029055565b60046040517f3ee5aeb5000000000000000000000000000000000000000000000000000000008152fd5b92946120a861066061065b3686866117db565b6120c973ffffffffffffffffffffffffffffffffffffffff85161515612208565b6120f57f000000000000000000000000000000000000000000000000000000000000000082111561180b565b6120fe86612279565b61210b60038711156122b0565b61211434611bc8565b959094612125600154871015612321565b5f9761213081612279565b806121ae575b945f96929461218088979588977ff3635b85ca76f364d2d87a3993c71d8461675e09f2cb9a04fba67aa22e48d916958d896121718c9b612279565b8960405198899833958a6123bf565b0390a18181156121a5575b8290f11561055b578061219c575090565b61203890612c3c565b506108fc61218b565b97505f95919486959461218087969587967ff3635b85ca76f364d2d87a3993c71d8461675e09f2cb9a04fba67aa22e48d916956121ec600354612392565b6121f581600355565b9d959a5095509599509550509450612136565b1561220f57565b608460405162461bcd60e51b815260206004820152602160248201527f4950546f6b656e5374616b696e673a20696e76616c69642064656c656761746f60448201527f72000000000000000000000000000000000000000000000000000000000000006064820152fd5b6004111561228357565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b156122b757565b608460405162461bcd60e51b815260206004820152602660248201527f4950546f6b656e5374616b696e673a20496e76616c6964207374616b696e672060448201527f706572696f6400000000000000000000000000000000000000000000000000006064820152fd5b1561232857565b608460405162461bcd60e51b815260206004820152602660248201527f4950546f6b656e5374616b696e673a205374616b6520616d6f756e7420756e6460448201527f6572206d696e00000000000000000000000000000000000000000000000000006064820152fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114611bc35760010190565b946123fb60ff92959997936120389b999573ffffffffffffffffffffffffffffffffffffffff809816895260e060208a015260e089019161187b565b98604087015216606085015260808401521660a082015260c081850391015261187b565b7fffffffffffffffffffffffff0000000000000000000000000000000000000000907f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c008281541690557f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080549073ffffffffffffffffffffffffffffffffffffffff80931680948316179055167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3565b8051156124df5760200190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b1561251357565b608460405162461bcd60e51b815260206004820152602c60248201527f536563703235366b3156657269666965723a20496e76616c696420636d70207060448201527f75626b65792070726566697800000000000000000000000000000000000000006064820152fd5b604051906080820182811067ffffffffffffffff82111761179c57604052604182526060366020840137565b60218151036126b85761265f61260b916126417fff000000000000000000000000000000000000000000000000000000000000007f02000000000000000000000000000000000000000000000000000000000000008161263161260b866124d2565b517fff000000000000000000000000000000000000000000000000000000000000001690565b1614908115612681575b5061250c565b61265a61265460218301519485936124d2565b60f81c90565b612ce5565b61266761257d565b916004612673846124d2565b536021830152604182015290565b7f030000000000000000000000000000000000000000000000000000000000000091506126b061260b856124d2565b16145f61263b565b608460405162461bcd60e51b815260206004820152602c60248201527f536563703235366b3156657269666965723a20496e76616c696420636d70207060448201527f75626b6579206c656e67746800000000000000000000000000000000000000006064820152fd5b1561272957565b608460405162461bcd60e51b815260206004820152602660248201527f536563703235366b3156657269666965723a207075626b6579206e6f74206f6e60448201527f20637572766500000000000000000000000000000000000000000000000000006064820152fd5b6041815103612874577f04000000000000000000000000000000000000000000000000000000000000007fff000000000000000000000000000000000000000000000000000000000000006127e7836124d2565b51160361280a57612805816041602161160c94015191015190612dcc565b612722565b608460405162461bcd60e51b815260206004820152602e60248201527f536563703235366b3156657269666965723a20496e76616c696420756e636d7060448201527f207075626b6579207072656669780000000000000000000000000000000000006064820152fd5b608460405162461bcd60e51b815260206004820152602e60248201527f536563703235366b3156657269666965723a20496e76616c696420756e636d7060448201527f207075626b6579206c656e6774680000000000000000000000000000000000006064820152fd5b60405190606082019180831067ffffffffffffffff84111761179c5773ffffffffffffffffffffffffffffffffffffffff926040526040815260416020820192604036853760218101518452015160408201525190201690565b1561293f57565b608460405162461bcd60e51b815260206004820152602860248201527f4950546f6b656e5374616b696e673a20436f6d6d697373696f6e20726174652060448201527f6f766572206d61780000000000000000000000000000000000000000000000006064820152fd5b156129b057565b608460405162461bcd60e51b815260206004820152602760248201527f4950546f6b656e5374616b696e673a204d6f6e696b6572206c656e677468206f60448201527f766572206d6178000000000000000000000000000000000000000000000000006064820152fd5b959190612a38919c9b9995989497939c61012080895288019161187b565b96602097868103898801528c518082525f5b818110612b0457506120389c9d50612af2969593612ab1612acb94847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8f965f612abe9882612ad59e9b0101520116019b60408c015260608b019063ffffffff169052565b63ffffffff166080890152565b63ffffffff1660a0870152565b60ff1660c0850152565b73ffffffffffffffffffffffffffffffffffffffff1660e0830152565b6101008382840301910152019161187b565b8e81018b01518382018c01528a01612a4a565b15612b1e57565b608460405162461bcd60e51b815260206004820152602e60248201527f4950546f6b656e5374616b696e673a20526564656c65676174696e6720746f2060448201527f73616d652076616c696461746f720000000000000000000000000000000000006064820152fd5b94612bc5612bd39360a0989b9a999596939673ffffffffffffffffffffffffffffffffffffffff809816895260c060208a015260c089019161187b565b91868303604088015261187b565b9660608401521660808201520152565b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c1615612c1257565b60046040517fd7e6bcf8000000000000000000000000000000000000000000000000000000008152fd5b5f80808093335af13d15612cd3573d612c576117ea826117a1565b9081525f60203d92013e5b15612c6957565b608460405162461bcd60e51b815260206004820152602a60248201527f4950546f6b656e5374616b696e673a204661696c656420746f20726566756e6460448201527f2072656d61696e646572000000000000000000000000000000000000000000006064820152fd5b612c62565b91908201809211611bc357565b60ff1690600282148015612dc2575b15612d5857612d42612d3b612d48927ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f9081806007815f8509089181818009900908612e7e565b9283612cd8565b60011690565b612d4f5790565b61203890611b95565b608460405162461bcd60e51b815260206004820152603160248201527f456c6c697074696343757276653a696e6e76616c696420636f6d70726573736560448201527f6420454320706f696e74207072656669780000000000000000000000000000006064820152fd5b5060038214612cf4565b80158015612e54575b8015612e4c575b8015612e22575b612e1c576007907ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f918282818181950909089180091490565b50505f90565b507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f821015612de3565b508115612ddc565b507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f811015612dd5565b8015612f46576001906001917f800000000000000000000000000000000000000000000000000000000000000091825b612eb85750505090565b9091927ffffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f90818080807f3fffffffffffffffffffffffffffffffffffffffffffffffffffffffbfffff0c94818a87161515890a918009098189891c86161515880a91800909818860021c85161515870a91800909918660031c161515840a918009099260041c919082612eae565b505f9056fea26469706673582212207f4946fa1cad86fb7d5c9c2ef7839da3eff0d121d5554fdcdf7d7a6d2f2ce1fa64736f6c63430008170033",
}

// IPTokenStakingABI is the input ABI used to generate the binding from.
// Deprecated: Use IPTokenStakingMetaData.ABI instead.
var IPTokenStakingABI = IPTokenStakingMetaData.ABI

// IPTokenStakingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use IPTokenStakingMetaData.Bin instead.
var IPTokenStakingBin = IPTokenStakingMetaData.Bin

// DeployIPTokenStaking deploys a new Ethereum contract, binding an instance of IPTokenStaking to it.
func DeployIPTokenStaking(auth *bind.TransactOpts, backend bind.ContractBackend, defaultMinFee *big.Int, maxDataLength *big.Int) (common.Address, *types.Transaction, *IPTokenStaking, error) {
	parsed, err := IPTokenStakingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(IPTokenStakingBin), backend, defaultMinFee, maxDataLength)
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

// AA is a free data retrieval call binding the contract method 0x997da8d4.
//
// Solidity: function AA() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) AA(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "AA")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AA is a free data retrieval call binding the contract method 0x997da8d4.
//
// Solidity: function AA() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) AA() (*big.Int, error) {
	return _IPTokenStaking.Contract.AA(&_IPTokenStaking.CallOpts)
}

// AA is a free data retrieval call binding the contract method 0x997da8d4.
//
// Solidity: function AA() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) AA() (*big.Int, error) {
	return _IPTokenStaking.Contract.AA(&_IPTokenStaking.CallOpts)
}

// BB is a free data retrieval call binding the contract method 0x5727dc5c.
//
// Solidity: function BB() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) BB(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "BB")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BB is a free data retrieval call binding the contract method 0x5727dc5c.
//
// Solidity: function BB() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) BB() (*big.Int, error) {
	return _IPTokenStaking.Contract.BB(&_IPTokenStaking.CallOpts)
}

// BB is a free data retrieval call binding the contract method 0x5727dc5c.
//
// Solidity: function BB() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) BB() (*big.Int, error) {
	return _IPTokenStaking.Contract.BB(&_IPTokenStaking.CallOpts)
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

// MAXDATALENGTH is a free data retrieval call binding the contract method 0xb3e85917.
//
// Solidity: function MAX_DATA_LENGTH() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) MAXDATALENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "MAX_DATA_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXDATALENGTH is a free data retrieval call binding the contract method 0xb3e85917.
//
// Solidity: function MAX_DATA_LENGTH() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) MAXDATALENGTH() (*big.Int, error) {
	return _IPTokenStaking.Contract.MAXDATALENGTH(&_IPTokenStaking.CallOpts)
}

// MAXDATALENGTH is a free data retrieval call binding the contract method 0xb3e85917.
//
// Solidity: function MAX_DATA_LENGTH() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) MAXDATALENGTH() (*big.Int, error) {
	return _IPTokenStaking.Contract.MAXDATALENGTH(&_IPTokenStaking.CallOpts)
}

// MAXMONIKERLENGTH is a free data retrieval call binding the contract method 0x396e1e47.
//
// Solidity: function MAX_MONIKER_LENGTH() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) MAXMONIKERLENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "MAX_MONIKER_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXMONIKERLENGTH is a free data retrieval call binding the contract method 0x396e1e47.
//
// Solidity: function MAX_MONIKER_LENGTH() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) MAXMONIKERLENGTH() (*big.Int, error) {
	return _IPTokenStaking.Contract.MAXMONIKERLENGTH(&_IPTokenStaking.CallOpts)
}

// MAXMONIKERLENGTH is a free data retrieval call binding the contract method 0x396e1e47.
//
// Solidity: function MAX_MONIKER_LENGTH() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) MAXMONIKERLENGTH() (*big.Int, error) {
	return _IPTokenStaking.Contract.MAXMONIKERLENGTH(&_IPTokenStaking.CallOpts)
}

// PP is a free data retrieval call binding the contract method 0xeeeac01e.
//
// Solidity: function PP() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) PP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "PP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PP is a free data retrieval call binding the contract method 0xeeeac01e.
//
// Solidity: function PP() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) PP() (*big.Int, error) {
	return _IPTokenStaking.Contract.PP(&_IPTokenStaking.CallOpts)
}

// PP is a free data retrieval call binding the contract method 0xeeeac01e.
//
// Solidity: function PP() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) PP() (*big.Int, error) {
	return _IPTokenStaking.Contract.PP(&_IPTokenStaking.CallOpts)
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

// MinCreateValidatorAmount is a free data retrieval call binding the contract method 0xe65f593f.
//
// Solidity: function minCreateValidatorAmount() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCaller) MinCreateValidatorAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPTokenStaking.contract.Call(opts, &out, "minCreateValidatorAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinCreateValidatorAmount is a free data retrieval call binding the contract method 0xe65f593f.
//
// Solidity: function minCreateValidatorAmount() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingSession) MinCreateValidatorAmount() (*big.Int, error) {
	return _IPTokenStaking.Contract.MinCreateValidatorAmount(&_IPTokenStaking.CallOpts)
}

// MinCreateValidatorAmount is a free data retrieval call binding the contract method 0xe65f593f.
//
// Solidity: function minCreateValidatorAmount() view returns(uint256)
func (_IPTokenStaking *IPTokenStakingCallerSession) MinCreateValidatorAmount() (*big.Int, error) {
	return _IPTokenStaking.Contract.MinCreateValidatorAmount(&_IPTokenStaking.CallOpts)
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
// Solidity: function roundedStakeAmount(uint256 rawAmount) pure returns(uint256 amount, uint256 remainder)
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
// Solidity: function roundedStakeAmount(uint256 rawAmount) pure returns(uint256 amount, uint256 remainder)
func (_IPTokenStaking *IPTokenStakingSession) RoundedStakeAmount(rawAmount *big.Int) (struct {
	Amount    *big.Int
	Remainder *big.Int
}, error) {
	return _IPTokenStaking.Contract.RoundedStakeAmount(&_IPTokenStaking.CallOpts, rawAmount)
}

// RoundedStakeAmount is a free data retrieval call binding the contract method 0xd2e1f5b8.
//
// Solidity: function roundedStakeAmount(uint256 rawAmount) pure returns(uint256 amount, uint256 remainder)
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

// CreateValidator is a paid mutator transaction binding the contract method 0x8740597a.
//
// Solidity: function createValidator(bytes validatorCmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool supportsUnlocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) CreateValidator(opts *bind.TransactOpts, validatorCmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, supportsUnlocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "createValidator", validatorCmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, supportsUnlocked, data)
}

// CreateValidator is a paid mutator transaction binding the contract method 0x8740597a.
//
// Solidity: function createValidator(bytes validatorCmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool supportsUnlocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) CreateValidator(validatorCmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, supportsUnlocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.CreateValidator(&_IPTokenStaking.TransactOpts, validatorCmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, supportsUnlocked, data)
}

// CreateValidator is a paid mutator transaction binding the contract method 0x8740597a.
//
// Solidity: function createValidator(bytes validatorCmpPubkey, string moniker, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, bool supportsUnlocked, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) CreateValidator(validatorCmpPubkey []byte, moniker string, commissionRate uint32, maxCommissionRate uint32, maxCommissionChangeRate uint32, supportsUnlocked bool, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.CreateValidator(&_IPTokenStaking.TransactOpts, validatorCmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate, supportsUnlocked, data)
}

// Initialize is a paid mutator transaction binding the contract method 0x08a3002b.
//
// Solidity: function initialize((address,uint256,uint256,uint256,uint256,uint256) args) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) Initialize(opts *bind.TransactOpts, args IIPTokenStakingInitializerArgs) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "initialize", args)
}

// Initialize is a paid mutator transaction binding the contract method 0x08a3002b.
//
// Solidity: function initialize((address,uint256,uint256,uint256,uint256,uint256) args) returns()
func (_IPTokenStaking *IPTokenStakingSession) Initialize(args IIPTokenStakingInitializerArgs) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Initialize(&_IPTokenStaking.TransactOpts, args)
}

// Initialize is a paid mutator transaction binding the contract method 0x08a3002b.
//
// Solidity: function initialize((address,uint256,uint256,uint256,uint256,uint256) args) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) Initialize(args IIPTokenStakingInitializerArgs) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Initialize(&_IPTokenStaking.TransactOpts, args)
}

// Redelegate is a paid mutator transaction binding the contract method 0xe52da4fc.
//
// Solidity: function redelegate(bytes validatorSrcCmpPubkey, bytes validatorDstCmpPubkey, uint256 delegationId, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) Redelegate(opts *bind.TransactOpts, validatorSrcCmpPubkey []byte, validatorDstCmpPubkey []byte, delegationId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "redelegate", validatorSrcCmpPubkey, validatorDstCmpPubkey, delegationId, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0xe52da4fc.
//
// Solidity: function redelegate(bytes validatorSrcCmpPubkey, bytes validatorDstCmpPubkey, uint256 delegationId, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) Redelegate(validatorSrcCmpPubkey []byte, validatorDstCmpPubkey []byte, delegationId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Redelegate(&_IPTokenStaking.TransactOpts, validatorSrcCmpPubkey, validatorDstCmpPubkey, delegationId, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0xe52da4fc.
//
// Solidity: function redelegate(bytes validatorSrcCmpPubkey, bytes validatorDstCmpPubkey, uint256 delegationId, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) Redelegate(validatorSrcCmpPubkey []byte, validatorDstCmpPubkey []byte, delegationId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Redelegate(&_IPTokenStaking.TransactOpts, validatorSrcCmpPubkey, validatorDstCmpPubkey, delegationId, amount)
}

// RedelegateOnBehalf is a paid mutator transaction binding the contract method 0x88138319.
//
// Solidity: function redelegateOnBehalf(address delegator, bytes validatorSrcCmpPubkey, bytes validatorDstCmpPubkey, uint256 delegationId, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) RedelegateOnBehalf(opts *bind.TransactOpts, delegator common.Address, validatorSrcCmpPubkey []byte, validatorDstCmpPubkey []byte, delegationId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "redelegateOnBehalf", delegator, validatorSrcCmpPubkey, validatorDstCmpPubkey, delegationId, amount)
}

// RedelegateOnBehalf is a paid mutator transaction binding the contract method 0x88138319.
//
// Solidity: function redelegateOnBehalf(address delegator, bytes validatorSrcCmpPubkey, bytes validatorDstCmpPubkey, uint256 delegationId, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) RedelegateOnBehalf(delegator common.Address, validatorSrcCmpPubkey []byte, validatorDstCmpPubkey []byte, delegationId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.RedelegateOnBehalf(&_IPTokenStaking.TransactOpts, delegator, validatorSrcCmpPubkey, validatorDstCmpPubkey, delegationId, amount)
}

// RedelegateOnBehalf is a paid mutator transaction binding the contract method 0x88138319.
//
// Solidity: function redelegateOnBehalf(address delegator, bytes validatorSrcCmpPubkey, bytes validatorDstCmpPubkey, uint256 delegationId, uint256 amount) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) RedelegateOnBehalf(delegator common.Address, validatorSrcCmpPubkey []byte, validatorDstCmpPubkey []byte, delegationId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.RedelegateOnBehalf(&_IPTokenStaking.TransactOpts, delegator, validatorSrcCmpPubkey, validatorDstCmpPubkey, delegationId, amount)
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

// SetMinCreateValidatorAmount is a paid mutator transaction binding the contract method 0x08e1b4ff.
//
// Solidity: function setMinCreateValidatorAmount(uint256 newMinCreateValidatorAmount) returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetMinCreateValidatorAmount(opts *bind.TransactOpts, newMinCreateValidatorAmount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setMinCreateValidatorAmount", newMinCreateValidatorAmount)
}

// SetMinCreateValidatorAmount is a paid mutator transaction binding the contract method 0x08e1b4ff.
//
// Solidity: function setMinCreateValidatorAmount(uint256 newMinCreateValidatorAmount) returns()
func (_IPTokenStaking *IPTokenStakingSession) SetMinCreateValidatorAmount(newMinCreateValidatorAmount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetMinCreateValidatorAmount(&_IPTokenStaking.TransactOpts, newMinCreateValidatorAmount)
}

// SetMinCreateValidatorAmount is a paid mutator transaction binding the contract method 0x08e1b4ff.
//
// Solidity: function setMinCreateValidatorAmount(uint256 newMinCreateValidatorAmount) returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetMinCreateValidatorAmount(newMinCreateValidatorAmount *big.Int) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetMinCreateValidatorAmount(&_IPTokenStaking.TransactOpts, newMinCreateValidatorAmount)
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

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setOperator", operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetOperator(&_IPTokenStaking.TransactOpts, operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetOperator(&_IPTokenStaking.TransactOpts, operator)
}

// SetRewardsAddress is a paid mutator transaction binding the contract method 0x8906758d.
//
// Solidity: function setRewardsAddress(address newRewardsAddress) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetRewardsAddress(opts *bind.TransactOpts, newRewardsAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setRewardsAddress", newRewardsAddress)
}

// SetRewardsAddress is a paid mutator transaction binding the contract method 0x8906758d.
//
// Solidity: function setRewardsAddress(address newRewardsAddress) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) SetRewardsAddress(newRewardsAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetRewardsAddress(&_IPTokenStaking.TransactOpts, newRewardsAddress)
}

// SetRewardsAddress is a paid mutator transaction binding the contract method 0x8906758d.
//
// Solidity: function setRewardsAddress(address newRewardsAddress) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetRewardsAddress(newRewardsAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetRewardsAddress(&_IPTokenStaking.TransactOpts, newRewardsAddress)
}

// SetWithdrawalAddress is a paid mutator transaction binding the contract method 0x21b8092e.
//
// Solidity: function setWithdrawalAddress(address newWithdrawalAddress) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) SetWithdrawalAddress(opts *bind.TransactOpts, newWithdrawalAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "setWithdrawalAddress", newWithdrawalAddress)
}

// SetWithdrawalAddress is a paid mutator transaction binding the contract method 0x21b8092e.
//
// Solidity: function setWithdrawalAddress(address newWithdrawalAddress) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) SetWithdrawalAddress(newWithdrawalAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetWithdrawalAddress(&_IPTokenStaking.TransactOpts, newWithdrawalAddress)
}

// SetWithdrawalAddress is a paid mutator transaction binding the contract method 0x21b8092e.
//
// Solidity: function setWithdrawalAddress(address newWithdrawalAddress) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) SetWithdrawalAddress(newWithdrawalAddress common.Address) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.SetWithdrawalAddress(&_IPTokenStaking.TransactOpts, newWithdrawalAddress)
}

// Stake is a paid mutator transaction binding the contract method 0x564bcc1f.
//
// Solidity: function stake(bytes validatorCmpPubkey, uint8 stakingPeriod, bytes data) payable returns(uint256 delegationId)
func (_IPTokenStaking *IPTokenStakingTransactor) Stake(opts *bind.TransactOpts, validatorCmpPubkey []byte, stakingPeriod uint8, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "stake", validatorCmpPubkey, stakingPeriod, data)
}

// Stake is a paid mutator transaction binding the contract method 0x564bcc1f.
//
// Solidity: function stake(bytes validatorCmpPubkey, uint8 stakingPeriod, bytes data) payable returns(uint256 delegationId)
func (_IPTokenStaking *IPTokenStakingSession) Stake(validatorCmpPubkey []byte, stakingPeriod uint8, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Stake(&_IPTokenStaking.TransactOpts, validatorCmpPubkey, stakingPeriod, data)
}

// Stake is a paid mutator transaction binding the contract method 0x564bcc1f.
//
// Solidity: function stake(bytes validatorCmpPubkey, uint8 stakingPeriod, bytes data) payable returns(uint256 delegationId)
func (_IPTokenStaking *IPTokenStakingTransactorSession) Stake(validatorCmpPubkey []byte, stakingPeriod uint8, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Stake(&_IPTokenStaking.TransactOpts, validatorCmpPubkey, stakingPeriod, data)
}

// StakeOnBehalf is a paid mutator transaction binding the contract method 0x77f8af38.
//
// Solidity: function stakeOnBehalf(address delegator, bytes validatorCmpPubkey, uint8 stakingPeriod, bytes data) payable returns(uint256 delegationId)
func (_IPTokenStaking *IPTokenStakingTransactor) StakeOnBehalf(opts *bind.TransactOpts, delegator common.Address, validatorCmpPubkey []byte, stakingPeriod uint8, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "stakeOnBehalf", delegator, validatorCmpPubkey, stakingPeriod, data)
}

// StakeOnBehalf is a paid mutator transaction binding the contract method 0x77f8af38.
//
// Solidity: function stakeOnBehalf(address delegator, bytes validatorCmpPubkey, uint8 stakingPeriod, bytes data) payable returns(uint256 delegationId)
func (_IPTokenStaking *IPTokenStakingSession) StakeOnBehalf(delegator common.Address, validatorCmpPubkey []byte, stakingPeriod uint8, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.StakeOnBehalf(&_IPTokenStaking.TransactOpts, delegator, validatorCmpPubkey, stakingPeriod, data)
}

// StakeOnBehalf is a paid mutator transaction binding the contract method 0x77f8af38.
//
// Solidity: function stakeOnBehalf(address delegator, bytes validatorCmpPubkey, uint8 stakingPeriod, bytes data) payable returns(uint256 delegationId)
func (_IPTokenStaking *IPTokenStakingTransactorSession) StakeOnBehalf(delegator common.Address, validatorCmpPubkey []byte, stakingPeriod uint8, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.StakeOnBehalf(&_IPTokenStaking.TransactOpts, delegator, validatorCmpPubkey, stakingPeriod, data)
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
// Solidity: function unjail(bytes validatorCmpPubkey, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) Unjail(opts *bind.TransactOpts, validatorCmpPubkey []byte, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "unjail", validatorCmpPubkey, data)
}

// Unjail is a paid mutator transaction binding the contract method 0x8ed65fbc.
//
// Solidity: function unjail(bytes validatorCmpPubkey, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) Unjail(validatorCmpPubkey []byte, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Unjail(&_IPTokenStaking.TransactOpts, validatorCmpPubkey, data)
}

// Unjail is a paid mutator transaction binding the contract method 0x8ed65fbc.
//
// Solidity: function unjail(bytes validatorCmpPubkey, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) Unjail(validatorCmpPubkey []byte, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Unjail(&_IPTokenStaking.TransactOpts, validatorCmpPubkey, data)
}

// UnjailOnBehalf is a paid mutator transaction binding the contract method 0x86eb5e48.
//
// Solidity: function unjailOnBehalf(bytes validatorCmpPubkey, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) UnjailOnBehalf(opts *bind.TransactOpts, validatorCmpPubkey []byte, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "unjailOnBehalf", validatorCmpPubkey, data)
}

// UnjailOnBehalf is a paid mutator transaction binding the contract method 0x86eb5e48.
//
// Solidity: function unjailOnBehalf(bytes validatorCmpPubkey, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) UnjailOnBehalf(validatorCmpPubkey []byte, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UnjailOnBehalf(&_IPTokenStaking.TransactOpts, validatorCmpPubkey, data)
}

// UnjailOnBehalf is a paid mutator transaction binding the contract method 0x86eb5e48.
//
// Solidity: function unjailOnBehalf(bytes validatorCmpPubkey, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) UnjailOnBehalf(validatorCmpPubkey []byte, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UnjailOnBehalf(&_IPTokenStaking.TransactOpts, validatorCmpPubkey, data)
}

// UnsetOperator is a paid mutator transaction binding the contract method 0xd2e9dedb.
//
// Solidity: function unsetOperator() payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) UnsetOperator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "unsetOperator")
}

// UnsetOperator is a paid mutator transaction binding the contract method 0xd2e9dedb.
//
// Solidity: function unsetOperator() payable returns()
func (_IPTokenStaking *IPTokenStakingSession) UnsetOperator() (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UnsetOperator(&_IPTokenStaking.TransactOpts)
}

// UnsetOperator is a paid mutator transaction binding the contract method 0xd2e9dedb.
//
// Solidity: function unsetOperator() payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) UnsetOperator() (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UnsetOperator(&_IPTokenStaking.TransactOpts)
}

// Unstake is a paid mutator transaction binding the contract method 0xd6f89acd.
//
// Solidity: function unstake(bytes validatorCmpPubkey, uint256 delegationId, uint256 amount, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) Unstake(opts *bind.TransactOpts, validatorCmpPubkey []byte, delegationId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "unstake", validatorCmpPubkey, delegationId, amount, data)
}

// Unstake is a paid mutator transaction binding the contract method 0xd6f89acd.
//
// Solidity: function unstake(bytes validatorCmpPubkey, uint256 delegationId, uint256 amount, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) Unstake(validatorCmpPubkey []byte, delegationId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Unstake(&_IPTokenStaking.TransactOpts, validatorCmpPubkey, delegationId, amount, data)
}

// Unstake is a paid mutator transaction binding the contract method 0xd6f89acd.
//
// Solidity: function unstake(bytes validatorCmpPubkey, uint256 delegationId, uint256 amount, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) Unstake(validatorCmpPubkey []byte, delegationId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.Unstake(&_IPTokenStaking.TransactOpts, validatorCmpPubkey, delegationId, amount, data)
}

// UnstakeOnBehalf is a paid mutator transaction binding the contract method 0x420b72f8.
//
// Solidity: function unstakeOnBehalf(address delegator, bytes validatorCmpPubkey, uint256 delegationId, uint256 amount, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) UnstakeOnBehalf(opts *bind.TransactOpts, delegator common.Address, validatorCmpPubkey []byte, delegationId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "unstakeOnBehalf", delegator, validatorCmpPubkey, delegationId, amount, data)
}

// UnstakeOnBehalf is a paid mutator transaction binding the contract method 0x420b72f8.
//
// Solidity: function unstakeOnBehalf(address delegator, bytes validatorCmpPubkey, uint256 delegationId, uint256 amount, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) UnstakeOnBehalf(delegator common.Address, validatorCmpPubkey []byte, delegationId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UnstakeOnBehalf(&_IPTokenStaking.TransactOpts, delegator, validatorCmpPubkey, delegationId, amount, data)
}

// UnstakeOnBehalf is a paid mutator transaction binding the contract method 0x420b72f8.
//
// Solidity: function unstakeOnBehalf(address delegator, bytes validatorCmpPubkey, uint256 delegationId, uint256 amount, bytes data) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) UnstakeOnBehalf(delegator common.Address, validatorCmpPubkey []byte, delegationId *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UnstakeOnBehalf(&_IPTokenStaking.TransactOpts, delegator, validatorCmpPubkey, delegationId, amount, data)
}

// UpdateValidatorCommission is a paid mutator transaction binding the contract method 0xc582db44.
//
// Solidity: function updateValidatorCommission(bytes validatorCmpPubkey, uint32 commissionRate) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactor) UpdateValidatorCommission(opts *bind.TransactOpts, validatorCmpPubkey []byte, commissionRate uint32) (*types.Transaction, error) {
	return _IPTokenStaking.contract.Transact(opts, "updateValidatorCommission", validatorCmpPubkey, commissionRate)
}

// UpdateValidatorCommission is a paid mutator transaction binding the contract method 0xc582db44.
//
// Solidity: function updateValidatorCommission(bytes validatorCmpPubkey, uint32 commissionRate) payable returns()
func (_IPTokenStaking *IPTokenStakingSession) UpdateValidatorCommission(validatorCmpPubkey []byte, commissionRate uint32) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UpdateValidatorCommission(&_IPTokenStaking.TransactOpts, validatorCmpPubkey, commissionRate)
}

// UpdateValidatorCommission is a paid mutator transaction binding the contract method 0xc582db44.
//
// Solidity: function updateValidatorCommission(bytes validatorCmpPubkey, uint32 commissionRate) payable returns()
func (_IPTokenStaking *IPTokenStakingTransactorSession) UpdateValidatorCommission(validatorCmpPubkey []byte, commissionRate uint32) (*types.Transaction, error) {
	return _IPTokenStaking.Contract.UpdateValidatorCommission(&_IPTokenStaking.TransactOpts, validatorCmpPubkey, commissionRate)
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
	ValidatorCmpPubkey      []byte
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
// Solidity: event CreateValidator(bytes validatorCmpPubkey, string moniker, uint256 stakeAmount, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, uint8 supportsUnlocked, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterCreateValidator(opts *bind.FilterOpts) (*IPTokenStakingCreateValidatorIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "CreateValidator")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingCreateValidatorIterator{contract: _IPTokenStaking.contract, event: "CreateValidator", logs: logs, sub: sub}, nil
}

// WatchCreateValidator is a free log subscription operation binding the contract event 0x65bfc2fa1cd4c6f50f60983ad1cf1cb4bff5ee6570428254dfce41b085ef6d14.
//
// Solidity: event CreateValidator(bytes validatorCmpPubkey, string moniker, uint256 stakeAmount, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, uint8 supportsUnlocked, address operatorAddress, bytes data)
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
// Solidity: event CreateValidator(bytes validatorCmpPubkey, string moniker, uint256 stakeAmount, uint32 commissionRate, uint32 maxCommissionRate, uint32 maxCommissionChangeRate, uint8 supportsUnlocked, address operatorAddress, bytes data)
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
	Delegator          common.Address
	ValidatorCmpPubkey []byte
	StakeAmount        *big.Int
	StakingPeriod      *big.Int
	DelegationId       *big.Int
	OperatorAddress    common.Address
	Data               []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xf3635b85ca76f364d2d87a3993c71d8461675e09f2cb9a04fba67aa22e48d916.
//
// Solidity: event Deposit(address delegator, bytes validatorCmpPubkey, uint256 stakeAmount, uint256 stakingPeriod, uint256 delegationId, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterDeposit(opts *bind.FilterOpts) (*IPTokenStakingDepositIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingDepositIterator{contract: _IPTokenStaking.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xf3635b85ca76f364d2d87a3993c71d8461675e09f2cb9a04fba67aa22e48d916.
//
// Solidity: event Deposit(address delegator, bytes validatorCmpPubkey, uint256 stakeAmount, uint256 stakingPeriod, uint256 delegationId, address operatorAddress, bytes data)
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

// ParseDeposit is a log parse operation binding the contract event 0xf3635b85ca76f364d2d87a3993c71d8461675e09f2cb9a04fba67aa22e48d916.
//
// Solidity: event Deposit(address delegator, bytes validatorCmpPubkey, uint256 stakeAmount, uint256 stakingPeriod, uint256 delegationId, address operatorAddress, bytes data)
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

// IPTokenStakingMinCreateValidatorAmountSetIterator is returned from FilterMinCreateValidatorAmountSet and is used to iterate over the raw logs and unpacked data for MinCreateValidatorAmountSet events raised by the IPTokenStaking contract.
type IPTokenStakingMinCreateValidatorAmountSetIterator struct {
	Event *IPTokenStakingMinCreateValidatorAmountSet // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingMinCreateValidatorAmountSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingMinCreateValidatorAmountSet)
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
		it.Event = new(IPTokenStakingMinCreateValidatorAmountSet)
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
func (it *IPTokenStakingMinCreateValidatorAmountSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingMinCreateValidatorAmountSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingMinCreateValidatorAmountSet represents a MinCreateValidatorAmountSet event raised by the IPTokenStaking contract.
type IPTokenStakingMinCreateValidatorAmountSet struct {
	MinCreateValidatorAmount *big.Int
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterMinCreateValidatorAmountSet is a free log retrieval operation binding the contract event 0xc58038d152352f59eeece288f00b152052e7ab8b5774c338a33ac4dffde1c374.
//
// Solidity: event MinCreateValidatorAmountSet(uint256 minCreateValidatorAmount)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterMinCreateValidatorAmountSet(opts *bind.FilterOpts) (*IPTokenStakingMinCreateValidatorAmountSetIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "MinCreateValidatorAmountSet")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingMinCreateValidatorAmountSetIterator{contract: _IPTokenStaking.contract, event: "MinCreateValidatorAmountSet", logs: logs, sub: sub}, nil
}

// WatchMinCreateValidatorAmountSet is a free log subscription operation binding the contract event 0xc58038d152352f59eeece288f00b152052e7ab8b5774c338a33ac4dffde1c374.
//
// Solidity: event MinCreateValidatorAmountSet(uint256 minCreateValidatorAmount)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchMinCreateValidatorAmountSet(opts *bind.WatchOpts, sink chan<- *IPTokenStakingMinCreateValidatorAmountSet) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "MinCreateValidatorAmountSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingMinCreateValidatorAmountSet)
				if err := _IPTokenStaking.contract.UnpackLog(event, "MinCreateValidatorAmountSet", log); err != nil {
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

// ParseMinCreateValidatorAmountSet is a log parse operation binding the contract event 0xc58038d152352f59eeece288f00b152052e7ab8b5774c338a33ac4dffde1c374.
//
// Solidity: event MinCreateValidatorAmountSet(uint256 minCreateValidatorAmount)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseMinCreateValidatorAmountSet(log types.Log) (*IPTokenStakingMinCreateValidatorAmountSet, error) {
	event := new(IPTokenStakingMinCreateValidatorAmountSet)
	if err := _IPTokenStaking.contract.UnpackLog(event, "MinCreateValidatorAmountSet", log); err != nil {
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
	Delegator             common.Address
	ValidatorSrcCmpPubkey []byte
	ValidatorDstCmpPubkey []byte
	DelegationId          *big.Int
	OperatorAddress       common.Address
	Amount                *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterRedelegate is a free log retrieval operation binding the contract event 0x08132c654ea3df406da083159e0fa3b4d3b8738bb81fc2dccefb22757868a122.
//
// Solidity: event Redelegate(address delegator, bytes validatorSrcCmpPubkey, bytes validatorDstCmpPubkey, uint256 delegationId, address operatorAddress, uint256 amount)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterRedelegate(opts *bind.FilterOpts) (*IPTokenStakingRedelegateIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "Redelegate")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingRedelegateIterator{contract: _IPTokenStaking.contract, event: "Redelegate", logs: logs, sub: sub}, nil
}

// WatchRedelegate is a free log subscription operation binding the contract event 0x08132c654ea3df406da083159e0fa3b4d3b8738bb81fc2dccefb22757868a122.
//
// Solidity: event Redelegate(address delegator, bytes validatorSrcCmpPubkey, bytes validatorDstCmpPubkey, uint256 delegationId, address operatorAddress, uint256 amount)
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

// ParseRedelegate is a log parse operation binding the contract event 0x08132c654ea3df406da083159e0fa3b4d3b8738bb81fc2dccefb22757868a122.
//
// Solidity: event Redelegate(address delegator, bytes validatorSrcCmpPubkey, bytes validatorDstCmpPubkey, uint256 delegationId, address operatorAddress, uint256 amount)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseRedelegate(log types.Log) (*IPTokenStakingRedelegate, error) {
	event := new(IPTokenStakingRedelegate)
	if err := _IPTokenStaking.contract.UnpackLog(event, "Redelegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingSetOperatorIterator is returned from FilterSetOperator and is used to iterate over the raw logs and unpacked data for SetOperator events raised by the IPTokenStaking contract.
type IPTokenStakingSetOperatorIterator struct {
	Event *IPTokenStakingSetOperator // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingSetOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingSetOperator)
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
		it.Event = new(IPTokenStakingSetOperator)
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
func (it *IPTokenStakingSetOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingSetOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingSetOperator represents a SetOperator event raised by the IPTokenStaking contract.
type IPTokenStakingSetOperator struct {
	Delegator common.Address
	Operator  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSetOperator is a free log retrieval operation binding the contract event 0x2709918445f306d3e94d280907c62c5d2525ac3192d2e544774c7f181d65af3e.
//
// Solidity: event SetOperator(address delegator, address operator)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterSetOperator(opts *bind.FilterOpts) (*IPTokenStakingSetOperatorIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "SetOperator")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingSetOperatorIterator{contract: _IPTokenStaking.contract, event: "SetOperator", logs: logs, sub: sub}, nil
}

// WatchSetOperator is a free log subscription operation binding the contract event 0x2709918445f306d3e94d280907c62c5d2525ac3192d2e544774c7f181d65af3e.
//
// Solidity: event SetOperator(address delegator, address operator)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchSetOperator(opts *bind.WatchOpts, sink chan<- *IPTokenStakingSetOperator) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "SetOperator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingSetOperator)
				if err := _IPTokenStaking.contract.UnpackLog(event, "SetOperator", log); err != nil {
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

// ParseSetOperator is a log parse operation binding the contract event 0x2709918445f306d3e94d280907c62c5d2525ac3192d2e544774c7f181d65af3e.
//
// Solidity: event SetOperator(address delegator, address operator)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseSetOperator(log types.Log) (*IPTokenStakingSetOperator, error) {
	event := new(IPTokenStakingSetOperator)
	if err := _IPTokenStaking.contract.UnpackLog(event, "SetOperator", log); err != nil {
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
	Delegator        common.Address
	ExecutionAddress [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSetRewardAddress is a free log retrieval operation binding the contract event 0x03564a99f640621f30a0be83f41fa9567acefa525e0132642d50d31c5b7d0e39.
//
// Solidity: event SetRewardAddress(address delegator, bytes32 executionAddress)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterSetRewardAddress(opts *bind.FilterOpts) (*IPTokenStakingSetRewardAddressIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "SetRewardAddress")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingSetRewardAddressIterator{contract: _IPTokenStaking.contract, event: "SetRewardAddress", logs: logs, sub: sub}, nil
}

// WatchSetRewardAddress is a free log subscription operation binding the contract event 0x03564a99f640621f30a0be83f41fa9567acefa525e0132642d50d31c5b7d0e39.
//
// Solidity: event SetRewardAddress(address delegator, bytes32 executionAddress)
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

// ParseSetRewardAddress is a log parse operation binding the contract event 0x03564a99f640621f30a0be83f41fa9567acefa525e0132642d50d31c5b7d0e39.
//
// Solidity: event SetRewardAddress(address delegator, bytes32 executionAddress)
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
	Delegator        common.Address
	ExecutionAddress [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSetWithdrawalAddress is a free log retrieval operation binding the contract event 0x1717e61ad5ca09c8dbe5620a60e8907a3124e59063fbcc7f71912e536511b421.
//
// Solidity: event SetWithdrawalAddress(address delegator, bytes32 executionAddress)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterSetWithdrawalAddress(opts *bind.FilterOpts) (*IPTokenStakingSetWithdrawalAddressIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "SetWithdrawalAddress")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingSetWithdrawalAddressIterator{contract: _IPTokenStaking.contract, event: "SetWithdrawalAddress", logs: logs, sub: sub}, nil
}

// WatchSetWithdrawalAddress is a free log subscription operation binding the contract event 0x1717e61ad5ca09c8dbe5620a60e8907a3124e59063fbcc7f71912e536511b421.
//
// Solidity: event SetWithdrawalAddress(address delegator, bytes32 executionAddress)
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

// ParseSetWithdrawalAddress is a log parse operation binding the contract event 0x1717e61ad5ca09c8dbe5620a60e8907a3124e59063fbcc7f71912e536511b421.
//
// Solidity: event SetWithdrawalAddress(address delegator, bytes32 executionAddress)
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
	Unjailer           common.Address
	ValidatorCmpPubkey []byte
	Data               []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterUnjail is a free log retrieval operation binding the contract event 0x026c2e156478ec2a25ccebac97a338d301f69b6d5aeec39c578b28a95e118201.
//
// Solidity: event Unjail(address unjailer, bytes validatorCmpPubkey, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterUnjail(opts *bind.FilterOpts) (*IPTokenStakingUnjailIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "Unjail")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingUnjailIterator{contract: _IPTokenStaking.contract, event: "Unjail", logs: logs, sub: sub}, nil
}

// WatchUnjail is a free log subscription operation binding the contract event 0x026c2e156478ec2a25ccebac97a338d301f69b6d5aeec39c578b28a95e118201.
//
// Solidity: event Unjail(address unjailer, bytes validatorCmpPubkey, bytes data)
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
// Solidity: event Unjail(address unjailer, bytes validatorCmpPubkey, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseUnjail(log types.Log) (*IPTokenStakingUnjail, error) {
	event := new(IPTokenStakingUnjail)
	if err := _IPTokenStaking.contract.UnpackLog(event, "Unjail", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingUnsetOperatorIterator is returned from FilterUnsetOperator and is used to iterate over the raw logs and unpacked data for UnsetOperator events raised by the IPTokenStaking contract.
type IPTokenStakingUnsetOperatorIterator struct {
	Event *IPTokenStakingUnsetOperator // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingUnsetOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingUnsetOperator)
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
		it.Event = new(IPTokenStakingUnsetOperator)
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
func (it *IPTokenStakingUnsetOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingUnsetOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingUnsetOperator represents a UnsetOperator event raised by the IPTokenStaking contract.
type IPTokenStakingUnsetOperator struct {
	Delegator common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUnsetOperator is a free log retrieval operation binding the contract event 0x5c6b8680e6b38b49c9b4293dc69be9b24e0ddcf054274619bb86e8e749f091a0.
//
// Solidity: event UnsetOperator(address delegator)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterUnsetOperator(opts *bind.FilterOpts) (*IPTokenStakingUnsetOperatorIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "UnsetOperator")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingUnsetOperatorIterator{contract: _IPTokenStaking.contract, event: "UnsetOperator", logs: logs, sub: sub}, nil
}

// WatchUnsetOperator is a free log subscription operation binding the contract event 0x5c6b8680e6b38b49c9b4293dc69be9b24e0ddcf054274619bb86e8e749f091a0.
//
// Solidity: event UnsetOperator(address delegator)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchUnsetOperator(opts *bind.WatchOpts, sink chan<- *IPTokenStakingUnsetOperator) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "UnsetOperator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingUnsetOperator)
				if err := _IPTokenStaking.contract.UnpackLog(event, "UnsetOperator", log); err != nil {
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

// ParseUnsetOperator is a log parse operation binding the contract event 0x5c6b8680e6b38b49c9b4293dc69be9b24e0ddcf054274619bb86e8e749f091a0.
//
// Solidity: event UnsetOperator(address delegator)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseUnsetOperator(log types.Log) (*IPTokenStakingUnsetOperator, error) {
	event := new(IPTokenStakingUnsetOperator)
	if err := _IPTokenStaking.contract.UnpackLog(event, "UnsetOperator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPTokenStakingUpdateValidatorCommissionIterator is returned from FilterUpdateValidatorCommission and is used to iterate over the raw logs and unpacked data for UpdateValidatorCommission events raised by the IPTokenStaking contract.
type IPTokenStakingUpdateValidatorCommissionIterator struct {
	Event *IPTokenStakingUpdateValidatorCommission // Event containing the contract specifics and raw log

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
func (it *IPTokenStakingUpdateValidatorCommissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPTokenStakingUpdateValidatorCommission)
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
		it.Event = new(IPTokenStakingUpdateValidatorCommission)
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
func (it *IPTokenStakingUpdateValidatorCommissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPTokenStakingUpdateValidatorCommissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPTokenStakingUpdateValidatorCommission represents a UpdateValidatorCommission event raised by the IPTokenStaking contract.
type IPTokenStakingUpdateValidatorCommission struct {
	ValidatorCmpPubkey []byte
	CommissionRate     uint32
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterUpdateValidatorCommission is a free log retrieval operation binding the contract event 0xd87fbe52b4b01e284c00c7b9a719018ee55b3f823bd39556a77bd14829107d72.
//
// Solidity: event UpdateValidatorCommission(bytes validatorCmpPubkey, uint32 commissionRate)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterUpdateValidatorCommission(opts *bind.FilterOpts) (*IPTokenStakingUpdateValidatorCommissionIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "UpdateValidatorCommission")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingUpdateValidatorCommissionIterator{contract: _IPTokenStaking.contract, event: "UpdateValidatorCommission", logs: logs, sub: sub}, nil
}

// WatchUpdateValidatorCommission is a free log subscription operation binding the contract event 0xd87fbe52b4b01e284c00c7b9a719018ee55b3f823bd39556a77bd14829107d72.
//
// Solidity: event UpdateValidatorCommission(bytes validatorCmpPubkey, uint32 commissionRate)
func (_IPTokenStaking *IPTokenStakingFilterer) WatchUpdateValidatorCommission(opts *bind.WatchOpts, sink chan<- *IPTokenStakingUpdateValidatorCommission) (event.Subscription, error) {

	logs, sub, err := _IPTokenStaking.contract.WatchLogs(opts, "UpdateValidatorCommission")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPTokenStakingUpdateValidatorCommission)
				if err := _IPTokenStaking.contract.UnpackLog(event, "UpdateValidatorCommission", log); err != nil {
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

// ParseUpdateValidatorCommission is a log parse operation binding the contract event 0xd87fbe52b4b01e284c00c7b9a719018ee55b3f823bd39556a77bd14829107d72.
//
// Solidity: event UpdateValidatorCommission(bytes validatorCmpPubkey, uint32 commissionRate)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseUpdateValidatorCommission(log types.Log) (*IPTokenStakingUpdateValidatorCommission, error) {
	event := new(IPTokenStakingUpdateValidatorCommission)
	if err := _IPTokenStaking.contract.UnpackLog(event, "UpdateValidatorCommission", log); err != nil {
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
	Delegator          common.Address
	ValidatorCmpPubkey []byte
	StakeAmount        *big.Int
	DelegationId       *big.Int
	OperatorAddress    common.Address
	Data               []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x89f97919c84baf7c0af918e2aae2581e51a70e1ce774948638769eac6b4d2d1c.
//
// Solidity: event Withdraw(address delegator, bytes validatorCmpPubkey, uint256 stakeAmount, uint256 delegationId, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) FilterWithdraw(opts *bind.FilterOpts) (*IPTokenStakingWithdrawIterator, error) {

	logs, sub, err := _IPTokenStaking.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &IPTokenStakingWithdrawIterator{contract: _IPTokenStaking.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x89f97919c84baf7c0af918e2aae2581e51a70e1ce774948638769eac6b4d2d1c.
//
// Solidity: event Withdraw(address delegator, bytes validatorCmpPubkey, uint256 stakeAmount, uint256 delegationId, address operatorAddress, bytes data)
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

// ParseWithdraw is a log parse operation binding the contract event 0x89f97919c84baf7c0af918e2aae2581e51a70e1ce774948638769eac6b4d2d1c.
//
// Solidity: event Withdraw(address delegator, bytes validatorCmpPubkey, uint256 stakeAmount, uint256 delegationId, address operatorAddress, bytes data)
func (_IPTokenStaking *IPTokenStakingFilterer) ParseWithdraw(log types.Log) (*IPTokenStakingWithdraw, error) {
	event := new(IPTokenStakingWithdraw)
	if err := _IPTokenStaking.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
