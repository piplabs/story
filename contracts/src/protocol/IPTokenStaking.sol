// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { ReentrancyGuardUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import { EnumerableSet } from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

import { IIPTokenStaking } from "../interfaces/IIPTokenStaking.sol";
import { Secp256k1 } from "../libraries/Secp256k1.sol";

/**
 * @title IPTokenStaking
 * @notice The deposit contract for IP token staked validators.
 */
contract IPTokenStaking is IIPTokenStaking, Ownable2StepUpgradeable, ReentrancyGuardUpgradeable {
    using EnumerableSet for EnumerableSet.AddressSet;

    /// @notice Default commission rate for a validator. Out of 100%, or 10_000.
    uint32 public immutable DEFAULT_COMMISSION_RATE;

    /// @notice Default maximum commission rate for a validator. Out of 100%, or 10_000.
    uint32 public immutable DEFAULT_MAX_COMMISSION_RATE;

    /// @notice Default maximum commission change rate for a validator. Out of 100%, or 10_000.
    uint32 public immutable DEFAULT_MAX_COMMISSION_CHANGE_RATE;

    /// @notice Stake amount increments, 1 ether => e.g. 1 ether, 2 ether, 5 ether etc.
    uint256 public immutable STAKE_ROUNDING;

    uint256 public immutable DEFAULT_MIN_UNJAIL_FEE;

    /// @notice Minimum amount required to stake.
    uint256 public minStakeAmount;

    /// @notice Minimum amount required to unstake.
    uint256 public minUnstakeAmount;

    /// @notice The interval between changing the withdrawal address for each delegator
    uint256 public withdrawalAddressChangeInterval; // TODO: Remove


    /// @notice Approved operators for delegators.
    /// @dev Delegator public key is a 33-byte compressed Secp256k1 public key.
    mapping(bytes validatorUncmpPubkey => EnumerableSet.AddressSet operators) private delegatorOperators; // TODO: Remove

    /// @notice The timestamp of last withdrawal address changes for each delegator.
    /// @dev Delegator public key is a 33-byte compressed Secp256k1 public key.
    mapping(bytes validatorUncmpPubkey => uint256 lastChange) public withdrawalAddressChange; // TODO: Remove

    /// @notice The timestamp of last reward address changes for each delegator.
    /// @dev Delegator public key is a 33-byte compressed Secp256k1 public key.
    mapping(bytes validatorUncmpPubkey => uint256 lastChange) public rewardAddressChange; // TODO: Remove

    /// @notice Counter to generate delegationIds for delegations with period.
    /// @dev Starts in 1, since 0 is reserved for flexible delegations.
    uint256 private _delegationIdCounter;

    /// @notice The fee paid to unjail a validator.
    uint256 public unjailFee;

    mapping(IIPTokenStaking.StakingPeriod period => uint32 duration) stakingDurations;

    constructor(
        uint256 stakingRounding,
        uint32 defaultCommissionRate,
        uint32 defaultMaxCommissionRate,
        uint32 defaultMaxCommissionChangeRate,
        uint256 defaultMinUnjailFee
    ) {
        STAKE_ROUNDING = stakingRounding; // Recommended: 1 gwei (10^9)
        require(defaultCommissionRate <= 10_000, "IPTokenStaking: Invalid default commission rate");
        DEFAULT_COMMISSION_RATE = defaultCommissionRate; // Recommended: 10%, or 1_000 / 10_000

        require(
            defaultMaxCommissionRate >= DEFAULT_COMMISSION_RATE && defaultMaxCommissionRate <= 10_000,
            "IPTokenStaking: Invalid default max commission rate"
        );
        DEFAULT_MAX_COMMISSION_RATE = defaultMaxCommissionRate; // Recommended: 50%, or 5_000 / 10_000

        require(defaultMaxCommissionChangeRate <= 10_000, "IPTokenStaking: Invalid default max commission change rate");
        DEFAULT_MAX_COMMISSION_CHANGE_RATE = defaultMaxCommissionChangeRate; // Recommended: 5%, or 500 / 10_000

        require(defaultMinUnjailFee >= 1 gwei, "IPTokenStaking: Invalid default min unjail fee");
        DEFAULT_MIN_UNJAIL_FEE = defaultMinUnjailFee;

        _disableInitializers();
    }

    /// @notice Initializes the contract.
    function initialize(IIPTokenStaking.InitializerArgs calldata args) public initializer {
        __ReentrancyGuard_init();
        __Ownable_init(args.accessManager);
        _setMinStakeAmount(args.minStakeAmount);
        _setMinUnstakeAmount(args.minUnstakeAmount);
        _setWithdrawalAddressChangeInterval(args.withdrawalAddressChangeInterval);
        _setStakingPeriods(args.shortStakingPeriod, args.mediumStakingPeriod, args.longStakingPeriod);
        _setUnjailFee(args.unjailFee);
    }

    /// @notice Verifies that the syntax of the given public key is a 65 byte uncompressed secp256k1 public key.
    modifier verifyUncmpPubkey(bytes calldata uncmpPubkey) {
        require(uncmpPubkey.length == 65, "IPTokenStaking: Invalid pubkey length");
        require(uncmpPubkey[0] == 0x04, "IPTokenStaking: Invalid pubkey prefix");
        _;
    }

    /// @notice Verifies that the given 65 byte uncompressed secp256k1 public key (with 0x04 prefix) is valid and
    /// matches the expected EVM address.
    modifier verifyUncmpPubkeyWithExpectedAddress(bytes calldata uncmpPubkey, address expectedAddress) {
        require(uncmpPubkey.length == 65, "IPTokenStaking: Invalid pubkey length");
        require(uncmpPubkey[0] == 0x04, "IPTokenStaking: Invalid pubkey prefix");
        require(
            _uncmpPubkeyToAddress(uncmpPubkey) == expectedAddress,
            "IPTokenStaking: Invalid pubkey derived address"
        );
        _;
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

    /// @dev Sets the interval for updating the withdrawal address.
    /// @param newWithdrawalAddressChangeInterval The interval between updating the withdrawal address.
    function setWithdrawalAddressChangeInterval(uint256 newWithdrawalAddressChangeInterval) external onlyOwner {
        _setWithdrawalAddressChangeInterval(newWithdrawalAddressChangeInterval);
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
            revert("IPTokenStaking: Invalid min unjail fee");
        }
        unjailFee = newUnjailFee;
        emit UnjailFeeSet(newUnjailFee);
    }

    function _setMinStakeAmount(uint256 newMinStakeAmount) private {
        require(newMinStakeAmount > 0, "IPTokenStaking: minStakeAmount cannot be 0");
        minStakeAmount = newMinStakeAmount - (newMinStakeAmount % STAKE_ROUNDING);
        emit MinStakeAmountSet(minStakeAmount);
    }

    /// @dev Sets the minimum amount required to withdraw.
    /// @param newMinUnstakeAmount The minimum amount required to stake.
    function _setMinUnstakeAmount(uint256 newMinUnstakeAmount) private {
        require(newMinUnstakeAmount > 0, "IPTokenStaking: minUnstakeAmount cannot be 0");
        minUnstakeAmount = newMinUnstakeAmount - (newMinUnstakeAmount % STAKE_ROUNDING);
        emit MinUnstakeAmountSet(minUnstakeAmount);
    }

    /// @dev Sets the interval for updating the withdrawal address.
    /// @param newWithdrawalAddressChangeInterval The interval between updating the withdrawal address.
    function _setWithdrawalAddressChangeInterval(uint256 newWithdrawalAddressChangeInterval) private {
        require(
            newWithdrawalAddressChangeInterval > 0,
            "IPTokenStaking: newWithdrawalAddressChangeInterval cannot be 0"
        );
        withdrawalAddressChangeInterval = newWithdrawalAddressChangeInterval;
        emit WithdrawalAddressChangeIntervalSet(newWithdrawalAddressChangeInterval);
    }

    function _setStakingPeriods(uint32 short, uint32 medium, uint32 long) private {
        if (short == 0) {
            revert("IPTokenStaking: short == 0");
        }
        if (short >= medium) {
            revert("IPTokenStaking: short >= medium");
        }
        if (medium >= long) {
            revert("IPTokenStaking: medium >= long");
        }
        stakingDurations[IIPTokenStaking.StakingPeriod.SHORT] = short;
        stakingDurations[IIPTokenStaking.StakingPeriod.MEDIUM] = medium;
        stakingDurations[IIPTokenStaking.StakingPeriod.LONG] = long;
        emit StakingPeriodsChanged(short, medium, long);
    }

    /// @notice Converts the given public key to an EVM address.
    /// @dev Assume all calls to this function passes in the uncompressed public key.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key, with prefix 04.
    /// @return address The EVM address derived from the public key.
    function _uncmpPubkeyToAddress(bytes calldata uncmpPubkey) internal pure returns (address) {
        return address(uint160(uint256(keccak256(uncmpPubkey[1:]))));
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                            Operator functions                          //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Returns the operators for the delegator.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key.
    function getOperators(bytes calldata uncmpPubkey) external view returns (address[] memory) { // TODO: Remove
        return delegatorOperators[uncmpPubkey].values();
    }

    /// @notice Adds an operator for the delegator.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param operator The operator address to add.
    function addOperator(
        bytes calldata uncmpPubkey,
        address operator
    ) external verifyUncmpPubkeyWithExpectedAddress(uncmpPubkey, msg.sender) {
        require(delegatorOperators[uncmpPubkey].add(operator), "IPTokenStaking: Operator already exists");
    }

    /// @notice Removes an operator for the delegator.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param operator The operator address to remove.
    function removeOperator(
        bytes calldata uncmpPubkey,
        address operator
    ) external verifyUncmpPubkeyWithExpectedAddress(uncmpPubkey, msg.sender) {
        require(delegatorOperators[uncmpPubkey].remove(operator), "IPTokenStaking: Operator not found");
    }

    /// @dev Verifies that the caller is an operator for the delegator. This check will revert if the caller is the
    /// delegator, which is intentional as this function should only be called in `onBehalf` functions.
    function _verifyCallerIsOperator(bytes memory delegatorUncmpPubkey, address caller) internal view {
        require(delegatorOperators[delegatorUncmpPubkey].contains(caller), "IPTokenStaking: Caller is not an operator");
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
        require(
            withdrawalAddressChange[delegatorUncmpPubkey] + withdrawalAddressChangeInterval < block.timestamp,
            "IPTokenStaking: Withdrawal address change cool-down"
        );
        withdrawalAddressChange[delegatorUncmpPubkey] = block.timestamp;

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
    // TODO: MIN STAKE 1024
    // TODO: Burn tokens
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
    // TODO: MIN STAKE, NON DEFAULT 1024
    // TODO: Explain reasoning
    // TODO: fee? 100 IP // Set fee
    // TODO: Burn tokens
    function createValidatorOnBehalf(
        bytes calldata validatorUncmpPubkey,
        bool isLocked,
        bytes calldata data
    ) external payable verifyUncmpPubkey(validatorUncmpPubkey) nonReentrant {
        _createValidator(
            validatorUncmpPubkey,
            "validator",
            DEFAULT_COMMISSION_RATE,
            DEFAULT_MAX_COMMISSION_RATE,
            DEFAULT_MAX_COMMISSION_CHANGE_RATE,
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
        // NOTE: We intentionally disable this check to circumvent the unsycned validator state between CL and EL.
        // In particular, when a validator gets its stake withdrawn to below the minimum power-reduction, it will get
        // removed from the validator set. If a user attempts to stake to that validator after it's removed from the
        // set, the EL contract will allow for staking but CL will fail because the validator info doesn't exist, thus
        // the user would lose the fund. Hence, by removing this check, user can and must first call createValidator
        // again for a removed validator before staking to it.
        // require(!validatorMetadata[validatorCmpPubkey].exists, "IPTokenStaking: Validator already exists");

        (uint256 stakeAmount, uint256 remainder) = roundedStakeAmount(msg.value);
        require(stakeAmount > 0, "IPTokenStaking: Stake amount too low");

        emit CreateValidator(
            validatorUncmpPubkey,
            moniker,
            stakeAmount,
            commissionRate,
            maxCommissionRate,
            maxCommissionChangeRate,
            isLocked? 1 : 0,
            data
        );

        _refundRemainder(remainder);
    }

    // TODO: ADD REDELEGATION BACK
    // TODO: MIN REDELEGATION

    // TODO: update validator method
    

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
    ) external payable verifyUncmpPubkeyWithExpectedAddress(delegatorUncmpPubkey, msg.sender) nonReentrant returns (uint256 delegationId) {
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
        require(stakeAmount >= minStakeAmount, "IPTokenStaking: Stake amount too low");
        
        uint32 duration = 0;
        if (stakingPeriod != IIPTokenStaking.StakingPeriod.FLEXIBLE) {
            _delegationIdCounter++;
            duration = stakingDurations[stakingPeriod];
        }
        emit Deposit(
            delegatorUncmpPubkey,
            validatorUncmpPubkey,
            stakeAmount,
            duration,
            _delegationIdCounter, // TODO: FLEXIBLE is 0
            msg.sender,
            data
        );
        // We burn staked
        payable(address(0x0)).transfer(stakeAmount);

        _refundRemainder(remainder);

        return _delegationIdCounter;// TODO: FLEXIBLE is 0
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
        _verifyCallerIsOperator(delegatorUncmpPubkey, msg.sender);
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
            revert();
        }
        emit Withdraw(delegatorUncmpPubkey, validatorUncmpPubkey, amount, delegationId, msg.sender, data);
    }

    /// @dev Refunds the remainder of the stake amount to the msg sender.
    /// @param remainder The remainder of the stake amount.
    function _refundRemainder(uint256 remainder) internal {
        (bool success, ) = msg.sender.call{ value: remainder }("");
        require(success, "IPTokenStaking: Failed to refund remainder");
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
    /// @param validatorUncmpPubkey The validator's 33-byte compressed Secp256k1 public key
    function unjailOnBehalf(bytes calldata validatorUncmpPubkey, bytes calldata data) external payable nonReentrant {
        _unjail(msg.value, validatorUncmpPubkey, data);
    }

    /// @dev Emits the Unjail event after burning the fee.
    function _unjail(uint256 fee, bytes calldata validatorCmpPubkey, bytes calldata data) private {
        require(fee == unjailFee, "IPTokenSlashing: Insufficient fee");
        payable(address(0x0)).transfer(fee);
        emit Unjail(msg.sender, validatorCmpPubkey, data);
    }
}
