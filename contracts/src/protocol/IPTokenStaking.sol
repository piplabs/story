// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { ReentrancyGuardUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import { EnumerableSet } from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

import { IIPTokenStaking } from "../interfaces/IIPTokenStaking.sol";
import { Secp256k1 } from "../libraries/Secp256k1.sol";
import { Errors } from "../libraries/Errors.sol";

/**
 * @title IPTokenStaking
 * @notice The deposit contract for IP token staked validators.
 */
contract IPTokenStaking is IIPTokenStaking, Ownable2StepUpgradeable, ReentrancyGuardUpgradeable {
    using EnumerableSet for EnumerableSet.AddressSet;

    /// @notice Stake amount increments, 1 ether => e.g. 1 ether, 2 ether, 5 ether etc.
    uint256 public immutable STAKE_ROUNDING;

    uint256 public immutable DEFAULT_MIN_UNJAIL_FEE;

    uint256 public minCommissionRate;

    /// @notice Minimum amount required to stake.
    uint256 public minStakeAmount;

    /// @notice Minimum amount required to unstake.
    uint256 public minUnstakeAmount;

    /// @notice Counter to generate delegationIds for delegations with period.
    /// @dev Starts in 1, since 0 is reserved for flexible delegations.
    uint256 private _delegationIdCounter;

    /// @notice The fee paid to unjail a validator.
    uint256 public unjailFee;

    mapping(IIPTokenStaking.StakingPeriod period => uint32 duration) stakingDurations;

    /// @notice Verifies that the syntax of the given public key is a 65 byte uncompressed secp256k1 public key.
    modifier verifyUncmpPubkey(bytes calldata uncmpPubkey) {
        if (uncmpPubkey.length == 65) {
            revert Errors.IPTokenStaking__InvalidPubkeyLength();
        }
        if (uncmpPubkey[0] == 0x04) {
            revert Errors.IPTokenStaking__InvalidPubkeyPrefix();
        }
        _;
    }

    /// @notice Verifies that the given 65 byte uncompressed secp256k1 public key (with 0x04 prefix) is valid and
    /// matches the expected EVM address.
    modifier verifyUncmpPubkeyWithExpectedAddress(bytes calldata uncmpPubkey, address expectedAddress) {
        if (uncmpPubkey.length == 65) {
            revert Errors.IPTokenStaking__InvalidPubkeyLength();
        }
        if (uncmpPubkey[0] == 0x04) {
            revert Errors.IPTokenStaking__InvalidPubkeyPrefix();
        }
        if (_uncmpPubkeyToAddress(uncmpPubkey) == expectedAddress) {
            revert Errors.IPTokenStaking__InvalidPubkeyDerivedAddress();
        }
        _;
    }

    constructor(uint256 stakingRounding, uint256 defaultMinUnjailFee) {
        if (stakingRounding == 0) {
            revert Errors.IPTokenStaking__ZeroStakingRounding();
        }
        STAKE_ROUNDING = stakingRounding; // Recommended: 1 gwei (10^9)
        if (defaultMinUnjailFee < 1 gwei) {
            revert Errors.IPTokenStaking__InvalidDefaultMinUnjailFee();
        }
        DEFAULT_MIN_UNJAIL_FEE = defaultMinUnjailFee;

        _disableInitializers();
    }

    /// @notice Initializes the contract.
    function initialize(IIPTokenStaking.InitializerArgs calldata args) public initializer {
        __ReentrancyGuard_init();
        __Ownable_init(args.accessManager);
        _setMinStakeAmount(args.minStakeAmount);
        _setMinUnstakeAmount(args.minUnstakeAmount);
        _setMinCommissionRate(args.minCommissionRate);
        _setStakingPeriods(args.shortStakingPeriod, args.mediumStakingPeriod, args.longStakingPeriod);
        _setUnjailFee(args.unjailFee);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                       Admin Setters/Getters                            //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Sets the minimum amount required to stake.
    /// @param newMinStakeAmount The minimum amount required to stake.
    function setMinStakeAmount(uint256 newMinStakeAmount) external onlyOwner {
        _setMinStakeAmount(newMinStakeAmount);
    }

    /// @dev Sets the minimum amount required to withdraw.
    /// @param newMinUnstakeAmount The minimum amount required to stake.
    function setMinUnstakeAmount(uint256 newMinUnstakeAmount) external onlyOwner {
        _setMinUnstakeAmount(newMinUnstakeAmount);
    }

    function setStakingPeriods(uint32 short, uint32 medium, uint32 long) external onlyOwner {
        _setStakingPeriods(short, medium, long);
    }

    /// @notice Sets the unjail fee.
    /// @param newUnjailFee The new unjail fee.
    function setUnjailFee(uint256 newUnjailFee) external onlyOwner {
        _setUnjailFee(newUnjailFee);
    }

    function _setUnjailFee(uint256 newUnjailFee) private {
        if (newUnjailFee < DEFAULT_MIN_UNJAIL_FEE) {
            revert Errors.IPTokenStaking__InvalidMinUnjailFee();
        }
        unjailFee = newUnjailFee;
        emit UnjailFeeSet(newUnjailFee);
    }

    function _setMinStakeAmount(uint256 newMinStakeAmount) private {
        if (newMinStakeAmount == 0) {
            revert Errors.IPTokenStaking__ZeroMinStakeAmount();
        }
        minStakeAmount = newMinStakeAmount - (newMinStakeAmount % STAKE_ROUNDING);
        emit MinStakeAmountSet(minStakeAmount);
    }

    /// @dev Sets the minimum amount required to withdraw.
    /// @param newMinUnstakeAmount The minimum amount required to stake.
    function _setMinUnstakeAmount(uint256 newMinUnstakeAmount) private {
        if (newMinUnstakeAmount == 0) {
            revert Errors.IPTokenStaking__ZeroMinUnstakeAmount();
        }
        minUnstakeAmount = newMinUnstakeAmount - (newMinUnstakeAmount % STAKE_ROUNDING);
        emit MinUnstakeAmountSet(minUnstakeAmount);
    }

    function _setStakingPeriods(uint32 short, uint32 medium, uint32 long) private {
        if (short == 0) {
            revert Errors.IPTokenStaking__ZeroShortPeriodDuration();
        }
        if (short >= medium) {
            revert Errors.IPTokenStaking__ShortPeriodLongerThanMedium();
        }
        if (medium >= long) {
            revert Errors.IPTokenStaking__MediumLongerThanLong();
        }
        stakingDurations[IIPTokenStaking.StakingPeriod.SHORT] = short;
        stakingDurations[IIPTokenStaking.StakingPeriod.MEDIUM] = medium;
        stakingDurations[IIPTokenStaking.StakingPeriod.LONG] = long;
        emit StakingPeriodsChanged(short, medium, long);
    }

    function setMinCommissionRate(uint256 newValue) external onlyOwner {
        _setMinCommissionRate(newValue);
    }

    function _setMinCommissionRate(uint256 newValue) private {
        if (newValue == 0) {
            revert Errors.IPTokenStaking__ZeroMinCommissionRate();
        }
        minCommissionRate = newValue;
        emit MinCommissionRateChanged(newValue);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                            Operator functions                          //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Adds an operator for the delegator.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param operator The operator address to add.
    function addOperator(
        bytes calldata uncmpPubkey,
        address operator
    ) external verifyUncmpPubkeyWithExpectedAddress(uncmpPubkey, msg.sender) {
        emit AddOperator(uncmpPubkey, operator);
    }

    /// @notice Removes an operator for the delegator.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param operator The operator address to remove.
    function removeOperator(
        bytes calldata uncmpPubkey,
        address operator
    ) external verifyUncmpPubkeyWithExpectedAddress(uncmpPubkey, msg.sender) {
        emit RemoveOperator(uncmpPubkey, operator);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                     Staking Configuration functions                    //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Set/Update the withdrawal address that receives the stake and reward withdrawals.
    /// @dev To prevent spam, only delegators with stake can call this function with cool-down time.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param newWithdrawalAddress EVM address to receive the stake and reward withdrawals.
    function setWithdrawalAddress(
        bytes calldata delegatorUncmpPubkey,
        address newWithdrawalAddress
    ) external verifyUncmpPubkeyWithExpectedAddress(delegatorUncmpPubkey, msg.sender) {
        emit SetWithdrawalAddress({
            delegatorUncmpPubkey: delegatorUncmpPubkey,
            executionAddress: bytes32(uint256(uint160(newWithdrawalAddress))) // left-padded bytes32 of the address
        });
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                          Validator Creation                            //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Entry point for creating a new validator with self delegation.
    /// @dev The caller must provide the uncompressed public key that matches the expected EVM address.
    /// @param validatorUncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param moniker The moniker of the validator.
    /// @param commissionRate The commission rate of the validator.
    /// @param maxCommissionRate The maximum commission rate of the validator.
    /// @param maxCommissionChangeRate The maximum commission change rate of the validator.
    // TODO: explain refund if fail
    function createValidator(
        bytes calldata validatorUncmpPubkey,
        string calldata moniker,
        uint32 commissionRate,
        uint32 maxCommissionRate,
        uint32 maxCommissionChangeRate,
        bool isLocked,
        bytes calldata data
    ) external payable verifyUncmpPubkeyWithExpectedAddress(validatorUncmpPubkey, msg.sender) nonReentrant {
        _createValidator(
            validatorUncmpPubkey,
            moniker,
            commissionRate,
            maxCommissionRate,
            maxCommissionChangeRate,
            isLocked,
            data
        );
    }

    /// @notice Entry point for creating a new validator with self delegation on behalf of the validator.
    /// @dev There's no minimum amount required to stake when creating a new validator.
    /// @param validatorUncmpPubkey 65 bytes uncompressed secp256k1 public key.
    // TODO: Explain reasoning
    function createValidatorOnBehalf(
        bytes calldata validatorUncmpPubkey,
        string calldata moniker,
        uint32 commissionRate,
        uint32 maxCommissionRate,
        uint32 maxCommissionChangeRate,
        bool isLocked,
        bytes calldata data
    ) external payable verifyUncmpPubkey(validatorUncmpPubkey) nonReentrant {
        _createValidator(
            validatorUncmpPubkey,
            moniker,
            commissionRate,
            maxCommissionRate,
            maxCommissionChangeRate,
            isLocked,
            data
        );
    }

    /// @dev Validator is the delegator when creating a new validator (self-delegation).
    /// @param validatorUncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param moniker The moniker of the validator.
    /// @param commissionRate The commission rate of the validator.
    /// @param maxCommissionRate The maximum commission rate of the validator.
    /// @param maxCommissionChangeRate The maximum commission change rate of the validator.
    function _createValidator(
        bytes calldata validatorUncmpPubkey,
        string memory moniker,
        uint32 commissionRate,
        uint32 maxCommissionRate,
        uint32 maxCommissionChangeRate,
        bool isLocked,
        bytes calldata data
    ) internal {
        (uint256 stakeAmount, uint256 remainder) = roundedStakeAmount(msg.value);
        if (stakeAmount < minStakeAmount) {
            revert Errors.IPTokenStaking__StakeAmountUnderMin();
        }
        if (commissionRate < minCommissionRate) {
            revert Errors.IPTokenStaking__CommissionRateUnderMin();
        }
        if (commissionRate > maxCommissionRate) {
            revert Errors.IPTokenStaking__CommissionRateOverMax();
        }
        payable(address(0)).transfer(stakeAmount);
        emit CreateValidator(
            validatorUncmpPubkey,
            moniker,
            stakeAmount,
            commissionRate,
            maxCommissionRate,
            maxCommissionChangeRate,
            isLocked ? 1 : 0,
            data
        );
        _refundRemainder(remainder);
    }

    function redelegate(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpSrcPubkey,
        bytes calldata validatorUncmpDstPubkey,
        uint256 amount
    ) external payable verifyUncmpPubkeyWithExpectedAddress(delegatorUncmpPubkey, msg.sender) {
        (uint256 stakeAmount, uint256 remainder) = roundedStakeAmount(msg.value);
        if (stakeAmount < minStakeAmount) {
            revert Errors.IPTokenStaking__StakeAmountUnderMin();
        }

        emit Redelegate(delegatorUncmpPubkey, validatorUncmpSrcPubkey, validatorUncmpDstPubkey, amount);
    }

    // TODO: update validator method (next version)

    /*//////////////////////////////////////////////////////////////////////////
    //                             Token Staking                              //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Entry point for staking IP token to stake to the given validator. The consensus chain is notified of
    /// the deposit and manages the stake accounting and validator onboarding. Payer must be the delegator.
    /// @dev When staking, consider it as BURNING. Unstaking (withdrawal) will trigger native minting.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpPubkey Validator's65 bytes uncompressed secp256k1 public key.
    function stake(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpPubkey,
        IIPTokenStaking.StakingPeriod stakingPeriod,
        bytes calldata data
    )
        external
        payable
        override
        verifyUncmpPubkeyWithExpectedAddress(delegatorUncmpPubkey, msg.sender)
        nonReentrant
        returns (uint256 delegationId)
    {
        return _stake(delegatorUncmpPubkey, validatorUncmpPubkey, stakingPeriod, data);
    }

    /// @notice Entry point for staking IP token to stake to the given validator. The consensus chain is notified of
    /// the stake and manages the stake accounting and validator onboarding. Payer can stake on behalf of another user,
    /// who will be the beneficiary of the stake.
    /// @dev When staking, consider it as BURNING. Unstaking (withdrawal) will trigger native minting.
    /// @param delegatorUncmpPubkey Delegator's65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpPubkey Validator's65 bytes uncompressed secp256k1 public key.
    function stakeOnBehalf(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpPubkey,
        IIPTokenStaking.StakingPeriod stakingPeriod,
        bytes calldata data
    ) external payable verifyUncmpPubkey(delegatorUncmpPubkey) nonReentrant returns (uint256 delegationId) {
        return _stake(delegatorUncmpPubkey, validatorUncmpPubkey, stakingPeriod, data);
    }

    /// @dev Creates a validator (x/staking.MsgCreateValidator) if it does not exist. Then delegates the stake to the
    /// validator (x/staking.MsgDelegate).
    /// @param delegatorUncmpPubkey Delegator's 65 byte uncompressed secp256k1 public key (no 0x04 prefix).
    /// @param validatorUncmpPubkey 33 byte compressed secp256k1 public key (no 0x04 prefix).
    function _stake(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpPubkey,
        IIPTokenStaking.StakingPeriod stakingPeriod,
        bytes calldata data
    ) internal returns (uint256) {
        (uint256 stakeAmount, uint256 remainder) = roundedStakeAmount(msg.value);
        if (stakeAmount < minStakeAmount) {
            revert Errors.IPTokenStaking__StakeAmountUnderMin();
        }
        uint32 duration = 0;
        uint256 delegationId = 0;
        if (stakingPeriod != IIPTokenStaking.StakingPeriod.FLEXIBLE) {
            delegationId = ++_delegationIdCounter;
            duration = stakingDurations[stakingPeriod];
        }
        emit Deposit(delegatorUncmpPubkey, validatorUncmpPubkey, stakeAmount, duration, delegationId, msg.sender, data);
        // We burn staked tokens
        payable(address(0)).transfer(stakeAmount);

        _refundRemainder(remainder);

        return delegationId;
    }

    /// @notice Returns the rounded stake amount and the remainder.
    /// @param rawAmount The raw stake amount.
    /// @return amount The rounded stake amount.
    /// @return remainder The remainder of the stake amount.
    function roundedStakeAmount(uint256 rawAmount) public view returns (uint256 amount, uint256 remainder) {
        remainder = rawAmount % STAKE_ROUNDING;
        amount = rawAmount - remainder;
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Unstake                                    //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Entry point for unstaking the previously staked token.
    /// @dev Unstake (withdrawal) will trigger native minting, so token in this contract is considered as burned.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpPubkey Validator's65 bytes uncompressed secp256k1 public key.
    /// @param amount Token amount to unstake.
    function unstake(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpPubkey,
        uint256 delegationId,
        uint256 amount,
        bytes calldata data
    ) external verifyUncmpPubkeyWithExpectedAddress(delegatorUncmpPubkey, msg.sender) {
        _unstake(delegatorUncmpPubkey, validatorUncmpPubkey, delegationId, amount, data);
    }

    /// @notice Entry point for unstaking the previously staked token on behalf of the delegator.
    /// @dev Must be an approved operator for the delegator.
    /// @param delegatorUncmpPubkey Delegator's65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpPubkey Validator's65 bytes uncompressed secp256k1 public key.
    /// @param amount Token amount to unstake.
    // TODO: explain operator needed in CL
    function unstakeOnBehalf(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpPubkey,
        uint256 delegationId,
        uint256 amount,
        bytes calldata data
    ) external verifyUncmpPubkey(delegatorUncmpPubkey) {
        _unstake(delegatorUncmpPubkey, validatorUncmpPubkey, delegationId, amount, data);
    }

    function _unstake(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpPubkey,
        uint256 delegationId,
        uint256 amount,
        bytes calldata data
    ) private {
        if (amount < minUnstakeAmount) {
            revert Errors.IPTokenStaking__LowUnstakeAmount();
        }
        emit Withdraw(delegatorUncmpPubkey, validatorUncmpPubkey, amount, delegationId, msg.sender, data);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Unjail                                    //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Requests to unjail the validator. Must pay fee on the execution side to prevent spamming.
    function unjail(
        bytes calldata validatorUncmpPubkey,
        bytes calldata data
    ) external payable verifyUncmpPubkeyWithExpectedAddress(validatorUncmpPubkey, msg.sender) nonReentrant {
        _unjail(msg.value, validatorUncmpPubkey, data);
    }

    /// @notice Requests to unjail a validator on behalf. Must pay fee on the execution side to prevent spamming.
    /// @param validatorUncmpPubkey The validator's 65-byte uncompressed Secp256k1 public key
    function unjailOnBehalf(
        bytes calldata validatorUncmpPubkey,
        bytes calldata data
    ) external payable nonReentrant verifyUncmpPubkey(validatorUncmpPubkey) {
        _unjail(msg.value, validatorUncmpPubkey, data);
    }

    /// @dev Emits the Unjail event after burning the fee.
    function _unjail(uint256 fee, bytes calldata validatorUncmpPubkey, bytes calldata data) private {
        if (fee != unjailFee) {
            revert Errors.IPTokenStaking__InsufficientFee();
        }
        payable(address(0x0)).transfer(fee);
        emit Unjail(msg.sender, validatorUncmpPubkey, data);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Helpers                                    //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Refunds the remainder of the stake amount to the msg sender.
    /// @param remainder The remainder of the stake amount.
    function _refundRemainder(uint256 remainder) private {
        (bool success, ) = msg.sender.call{ value: remainder }("");
        if (!success) {
            revert Errors.IPTokenStaking__FailedRemainerRefund();
        }
    }

    /// @notice Converts the given public key to an EVM address.
    /// @dev Assume all calls to this function passes in the uncompressed public key.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key, with prefix 04.
    /// @return address The EVM address derived from the public key.
    function _uncmpPubkeyToAddress(bytes calldata uncmpPubkey) internal pure returns (address) {
        return address(uint160(uint256(keccak256(uncmpPubkey[1:]))));
    }
}
