// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { ReentrancyGuardUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import { EnumerableSet } from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

import { IIPTokenStaking } from "../interfaces/IIPTokenStaking.sol";
import { Secp256k1 } from "../libraries/Secp256k1.sol";
import { UpgradeabilityFlag } from "./UpradeabilityFlag.sol";

/**
 * @title IPTokenStaking
 * @notice The deposit contract for IP token staked validators.
 */
contract IPTokenStaking is
    IIPTokenStaking,
    Ownable2StepUpgradeable,
    ReentrancyGuardUpgradeable,
    UUPSUpgradeable,
    UpgradeabilityFlag
{
    using EnumerableSet for EnumerableSet.AddressSet;

    /// @notice Default commission rate for a validator. Out of 100%, or 10_000.
    uint32 public immutable DEFAULT_COMMISSION_RATE;

    /// @notice Default maximum commission rate for a validator. Out of 100%, or 10_000.
    uint32 public immutable DEFAULT_MAX_COMMISSION_RATE;

    /// @notice Default maximum commission change rate for a validator. Out of 100%, or 10_000.
    uint32 public immutable DEFAULT_MAX_COMMISSION_CHANGE_RATE;

    /// @notice Stake amount increments, 1 ether => e.g. 1 ether, 2 ether, 5 ether etc.
    uint256 public immutable STAKE_ROUNDING;

    /// @notice Minimum amount required to stake.
    uint256 public minStakeAmount;

    /// @notice Minimum amount required to unstake.
    uint256 public minUnstakeAmount;

    /// @notice Minimum amount required to redelegate.
    uint256 public minRedelegateAmount;

    /// @notice The interval between changing the withdrawal address for each delegator
    uint256 public withdrawalAddressChangeInterval;

    /// @notice Validator's metadata.
    /// @dev Validator public key is a 33-byte compressed Secp256k1 public key.
    mapping(bytes validatorCmpPubkey => ValidatorMetadata metadata) public validatorMetadata;

    /// @notice Delegator's total staked amount.
    /// @dev Delegator public key is a 33-byte compressed Secp256k1 public key.
    mapping(bytes delegatorCmpPubkey => uint256 stakedAmount) public delegatorTotalStakes;

    /// @notice Delegator's staked amount for the given validator.
    /// @dev Delegator and validator public keys are 33-byte compressed Secp256k1 public keys.
    mapping(bytes delegatorCmpPubkey => mapping(bytes validatorCmpPubkey => uint256 stakedAmount))
        public delegatorValidatorStakes;

    /// @notice Approved operators for delegators.
    /// @dev Delegator public key is a 33-byte compressed Secp256k1 public key.
    mapping(bytes delegatorCmpPubkey => EnumerableSet.AddressSet operators) private delegatorOperators;

    /// @notice The timestamp of last withdrawal address changes for each delegator.
    /// @dev Delegator public key is a 33-byte compressed Secp256k1 public key.
    mapping(bytes delegatorCmpPubkey => uint256 lastChange) public withdrawalAddressChange;

    constructor(
        uint256 stakingRounding,
        uint32 defaultCommissionRate,
        uint32 defaultMaxCommissionRate,
        uint32 defaultMaxCommissionChangeRate
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

        _disableInitializers();
    }

    /// @notice Initializes the contract.
    function initialize(
        address accessManager,
        uint256 _minStakeAmount,
        uint256 _minUnstakeAmount,
        uint256 _minRedelegateAmount,
        uint256 _withdrawalAddressChangeInterval
    ) public initializer {
        __ReentrancyGuard_init();
        __UUPSUpgradeable_init();
        __Ownable_init(accessManager);
        _setMinStakeAmount(_minStakeAmount);
        _setMinUnstakeAmount(_minUnstakeAmount);
        _setMinRedelegateAmount(_minRedelegateAmount);
        _setWithdrawalAddressChangeInterval(_withdrawalAddressChangeInterval);
    }

    /// @notice Verifies that the syntax of the given public key is a 33 byte compressed secp256k1 public key.
    modifier verifyPubkey(bytes calldata pubkey) {
        require(pubkey.length == 33, "IPTokenStaking: Invalid pubkey length");
        require(pubkey[0] == 0x02 || pubkey[0] == 0x03, "IPTokenStaking: Invalid pubkey prefix");
        _;
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

    /// @notice Verifies that the validator with the given pubkey exists.
    modifier verifyExistingValidator(bytes calldata validatorCmpPubkey) {
        require(validatorCmpPubkey.length == 33, "IPTokenStaking: Invalid pubkey length");
        require(
            validatorCmpPubkey[0] == 0x02 || validatorCmpPubkey[0] == 0x03,
            "IPTokenStaking: Invalid pubkey prefix"
        );
        require(validatorMetadata[validatorCmpPubkey].exists, "IPTokenStaking: Validator does not exist");
        _;
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Setters/Getters                            //
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

    /// @dev Sets the minimum amount required to redelegate.
    /// @param newMinRedelegateAmount The minimum amount required to redelegate.
    function setMinRedelegateAmount(uint256 newMinRedelegateAmount) external onlyOwner {
        _setMinRedelegateAmount(newMinRedelegateAmount);
    }

    /// @dev Sets the interval for updating the withdrawal address.
    /// @param newWithdrawalAddressChangeInterval The interval between updating the withdrawal address.
    function setWithdrawalAddressChangeInterval(uint256 newWithdrawalAddressChangeInterval) external onlyOwner {
        _setWithdrawalAddressChangeInterval(newWithdrawalAddressChangeInterval);
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

    /// @dev Sets the minimum amount required to redelegate.
    /// @param newMinRedelegateAmount The minimum amount required to redelegate.
    function _setMinRedelegateAmount(uint256 newMinRedelegateAmount) private {
        require(newMinRedelegateAmount > 0, "IPTokenStaking: minRedelegateAmount cannot be 0");
        minRedelegateAmount = newMinRedelegateAmount - (newMinRedelegateAmount % STAKE_ROUNDING);
        emit MinRedelegateAmountSet(minRedelegateAmount);
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

    /// @notice Returns the rounded stake amount and the remainder.
    /// @param rawAmount The raw stake amount.
    /// @return amount The rounded stake amount.
    /// @return remainder The remainder of the stake amount.
    function roundedStakeAmount(uint256 rawAmount) public view returns (uint256 amount, uint256 remainder) {
        remainder = rawAmount % STAKE_ROUNDING;
        amount = rawAmount - remainder;
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
    /// @param pubkey 33 bytes compressed secp256k1 public key.
    function getOperators(bytes calldata pubkey) external view returns (address[] memory) {
        return delegatorOperators[pubkey].values();
    }

    /// @notice Adds an operator for the delegator.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param operator The operator address to add.
    function addOperator(
        bytes calldata uncmpPubkey,
        address operator
    ) external verifyUncmpPubkeyWithExpectedAddress(uncmpPubkey, msg.sender) {
        bytes memory delegatorCmpPubkey = Secp256k1.compressPublicKey(uncmpPubkey);
        require(delegatorOperators[delegatorCmpPubkey].add(operator), "IPTokenStaking: Operator already exists");
    }

    /// @notice Removes an operator for the delegator.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param operator The operator address to remove.
    function removeOperator(
        bytes calldata uncmpPubkey,
        address operator
    ) external verifyUncmpPubkeyWithExpectedAddress(uncmpPubkey, msg.sender) {
        bytes memory delegatorCmpPubkey = Secp256k1.compressPublicKey(uncmpPubkey);
        require(delegatorOperators[delegatorCmpPubkey].remove(operator), "IPTokenStaking: Operator not found");
    }

    /// @dev Verifies that the caller is an operator for the delegator. This check will revert if the caller is the
    /// delegator, which is intentional as this function should only be called in `onBehalf` functions.
    function _verifyCallerIsOperator(bytes memory delegatorCmpPubkey, address caller) internal view {
        require(delegatorOperators[delegatorCmpPubkey].contains(caller), "IPTokenStaking: Caller is not an operator");
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
        bytes memory delegatorCmpPubkey = Secp256k1.compressPublicKey(delegatorUncmpPubkey);
        require(delegatorTotalStakes[delegatorCmpPubkey] > 0, "IPTokenStaking: Delegator must have stake");

        require(
            withdrawalAddressChange[delegatorCmpPubkey] + withdrawalAddressChangeInterval < block.timestamp,
            "IPTokenStaking: Withdrawal address change cool-down"
        );
        withdrawalAddressChange[delegatorCmpPubkey] = block.timestamp;

        emit SetWithdrawalAddress({
            delegatorCmpPubkey: delegatorCmpPubkey,
            executionAddress: bytes32(uint256(uint160(newWithdrawalAddress))) // left-padded bytes32 of the address
        });
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                         Token Staking functions                        //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Entry point for creating a new validator with self delegation.
    /// @dev The caller must provide the uncompressed public key that matches the expected EVM address.
    /// @param validatorUncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param moniker The moniker of the validator.
    /// @param commissionRate The commission rate of the validator.
    /// @param maxCommissionRate The maximum commission rate of the validator.
    /// @param maxCommissionChangeRate The maximum commission change rate of the validator.
    function createValidator(
        bytes calldata validatorUncmpPubkey,
        string calldata moniker,
        uint32 commissionRate,
        uint32 maxCommissionRate,
        uint32 maxCommissionChangeRate
    ) external payable verifyUncmpPubkeyWithExpectedAddress(validatorUncmpPubkey, msg.sender) nonReentrant {
        _createValidator(validatorUncmpPubkey, moniker, commissionRate, maxCommissionRate, maxCommissionChangeRate);
    }

    /// @notice Entry point for creating a new validator with self delegation on behalf of the validator.
    /// @dev There's no minimum amount required to stake when creating a new validator.
    /// @param validatorUncmpPubkey 65 bytes uncompressed secp256k1 public key.
    function createValidatorOnBehalf(
        bytes calldata validatorUncmpPubkey
    ) external payable verifyUncmpPubkey(validatorUncmpPubkey) nonReentrant {
        _createValidator(
            validatorUncmpPubkey,
            "validator",
            DEFAULT_COMMISSION_RATE,
            DEFAULT_MAX_COMMISSION_RATE,
            DEFAULT_MAX_COMMISSION_CHANGE_RATE
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
        uint32 maxCommissionChangeRate
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

        bytes memory validatorCmpPubkey = Secp256k1.compressPublicKey(validatorUncmpPubkey);

        // Since users can call createValidator multiple times, we need to reuse the existing metadata if it exists.
        // The only data that monotonically increases is the totalStake. All others are selected based on whether the
        // validator data already exists or not.
        ValidatorMetadata storage vm = validatorMetadata[validatorCmpPubkey];
        bool exists = vm.exists; // get before updating
        vm.exists = true;
        vm.moniker = exists ? vm.moniker : moniker;
        vm.totalStake = vm.totalStake + stakeAmount;
        vm.commissionRate = exists ? vm.commissionRate : commissionRate;
        vm.maxCommissionRate = exists ? vm.maxCommissionRate : maxCommissionRate;
        vm.maxCommissionChangeRate = exists ? vm.maxCommissionChangeRate : maxCommissionChangeRate;

        delegatorTotalStakes[validatorCmpPubkey] += stakeAmount;
        delegatorValidatorStakes[validatorCmpPubkey][validatorCmpPubkey] += stakeAmount;

        _refundRemainder(remainder);

        emit CreateValidator({
            validatorUncmpPubkey: validatorUncmpPubkey,
            validatorCmpPubkey: validatorCmpPubkey,
            moniker: vm.moniker,
            stakeAmount: stakeAmount,
            commissionRate: vm.commissionRate,
            maxCommissionRate: vm.maxCommissionRate,
            maxCommissionChangeRate: vm.maxCommissionChangeRate
        });
    }

    /// @notice Entry point for staking IP token to stake to the given validator. The consensus chain is notified of
    /// the deposit and manages the stake accounting and validator onboarding. Payer must be the delegator.
    /// @dev When staking, consider it as BURNING. Unstaking (withdrawal) will trigger native minting.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorCmpPubkey Validator's 33 bytes compressed secp256k1 public key.
    function stake(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorCmpPubkey
    )
        external
        payable
        verifyUncmpPubkeyWithExpectedAddress(delegatorUncmpPubkey, msg.sender)
        verifyExistingValidator(validatorCmpPubkey)
        nonReentrant
    {
        _stake(delegatorUncmpPubkey, validatorCmpPubkey);
    }

    /// @notice Entry point for staking IP token to stake to the given validator. The consensus chain is notified of
    /// the stake and manages the stake accounting and validator onboarding. Payer can stake on behalf of another user,
    /// who will be the beneficiary of the stake.
    /// @dev When staking, consider it as BURNING. Unstaking (withdrawal) will trigger native minting.
    /// @param delegatorUncmpPubkey Delegator's 33 bytes compressed secp256k1 public key.
    /// @param validatorCmpPubkey Validator's 33 bytes compressed secp256k1 public key.
    function stakeOnBehalf(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorCmpPubkey
    )
        external
        payable
        verifyUncmpPubkey(delegatorUncmpPubkey)
        verifyExistingValidator(validatorCmpPubkey)
        nonReentrant
    {
        _stake(delegatorUncmpPubkey, validatorCmpPubkey);
    }

    /// @dev Creates a validator (x/staking.MsgCreateValidator) if it does not exist. Then delegates the stake to the
    /// validator (x/staking.MsgDelegate).
    /// @param delegatorUncmpPubkey Delegator's 65 byte uncompressed secp256k1 public key (no 0x04 prefix).
    /// @param validatorCmpPubkey 33 byte compressed secp256k1 public key (no 0x04 prefix).
    function _stake(bytes memory delegatorUncmpPubkey, bytes calldata validatorCmpPubkey) internal {
        (uint256 stakeAmount, uint256 remainder) = roundedStakeAmount(msg.value);
        require(stakeAmount >= minStakeAmount, "IPTokenStaking: Stake amount too low");

        bytes memory delegatorCmpPubkey = Secp256k1.compressPublicKey(delegatorUncmpPubkey);

        unchecked {
            validatorMetadata[validatorCmpPubkey].totalStake += stakeAmount;
            delegatorTotalStakes[delegatorCmpPubkey] += stakeAmount;
            delegatorValidatorStakes[delegatorCmpPubkey][validatorCmpPubkey] += stakeAmount;
        }

        _refundRemainder(remainder);

        emit Deposit(delegatorUncmpPubkey, delegatorCmpPubkey, validatorCmpPubkey, stakeAmount);
    }

    // TODO: Redelegation also requires unbonding period to be executed. Should we separate storage for this for el?
    /// @notice Entry point for redelegating the staked token.
    /// @dev Redelegateion redelegates staked token from src validator to dst validator (x/staking.MsgBeginRedelegate)
    /// @param p See RedelegateParams
    function redelegate(
        RedelegateParams calldata p
    )
        external
        verifyUncmpPubkeyWithExpectedAddress(p.delegatorUncmpPubkey, msg.sender)
        verifyExistingValidator(p.validatorCmpSrcPubkey)
        verifyExistingValidator(p.validatorCmpDstPubkey)
    {
        (uint256 stakeAmount, ) = roundedStakeAmount(p.amount);
        bytes memory delegatorCmpPubkey = Secp256k1.compressPublicKey(p.delegatorUncmpPubkey);

        _redelegate(delegatorCmpPubkey, p.validatorCmpSrcPubkey, p.validatorCmpDstPubkey, stakeAmount);
    }

    /// @notice Entry point for redelegating the staked token on behalf of the delegator.
    /// @dev Redelegateion redelegates staked token from src validator to dst validator (x/staking.MsgBeginRedelegate)
    /// @param p See RedelegateParams
    function redelegateOnBehalf(
        RedelegateParams calldata p
    )
        external
        verifyUncmpPubkey(p.delegatorUncmpPubkey)
        verifyExistingValidator(p.validatorCmpSrcPubkey)
        verifyExistingValidator(p.validatorCmpDstPubkey)
    {
        bytes memory delegatorCmpPubkey = Secp256k1.compressPublicKey(p.delegatorUncmpPubkey);

        (uint256 stakeAmount, ) = roundedStakeAmount(p.amount);

        _verifyCallerIsOperator(delegatorCmpPubkey, msg.sender);
        _redelegate(delegatorCmpPubkey, p.validatorCmpSrcPubkey, p.validatorCmpDstPubkey, stakeAmount);
    }

    /// @dev Redelegates the given amount from the source validator to the destination validator.
    /// @param delegatorCmpPubkey Delegator's 33 bytes compressed secp256k1 public key.
    /// @param validatorSrcPubkey Source validator's 33 bytes compressed secp256k1 public key.
    /// @param validatorDstPubkey Destination validator's 33 bytes compressed secp256k1 public key.
    /// @param amount Token amount to redelegate.
    function _redelegate(
        bytes memory delegatorCmpPubkey,
        bytes calldata validatorSrcPubkey,
        bytes calldata validatorDstPubkey,
        uint256 amount
    ) internal {
        require(
            delegatorValidatorStakes[delegatorCmpPubkey][validatorSrcPubkey] >= amount,
            "IPTokenStaking: Insufficient staked amount"
        );

        validatorMetadata[validatorSrcPubkey].totalStake -= amount;
        validatorMetadata[validatorDstPubkey].totalStake += amount;

        delegatorValidatorStakes[delegatorCmpPubkey][validatorSrcPubkey] -= amount;
        delegatorValidatorStakes[delegatorCmpPubkey][validatorDstPubkey] += amount;

        emit Redelegate(delegatorCmpPubkey, validatorSrcPubkey, validatorDstPubkey, amount);
    }

    /// @notice Entry point for unstaking the previously staked token.
    /// @dev Unstake (withdrawal) will trigger native minting, so token in this contract is considered as burned.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorCmpPubkey Validator's 33 bytes compressed secp256k1 public key.
    /// @param amount Token amount to unstake.
    function unstake(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorCmpPubkey,
        uint256 amount
    )
        external
        verifyUncmpPubkeyWithExpectedAddress(delegatorUncmpPubkey, msg.sender)
        verifyExistingValidator(validatorCmpPubkey)
    {
        bytes memory delegatorCmpPubkey = Secp256k1.compressPublicKey(delegatorUncmpPubkey);
        _unstake(delegatorCmpPubkey, validatorCmpPubkey, amount);
    }

    /// @notice Entry point for unstaking the previously staked token on behalf of the delegator.
    /// @dev Must be an approved operator for the delegator.
    /// @param delegatorCmpPubkey Delegator's 33 bytes compressed secp256k1 public key.
    /// @param validatorCmpPubkey Validator's 33 bytes compressed secp256k1 public key.
    /// @param amount Token amount to unstake.
    function unstakeOnBehalf(
        bytes calldata delegatorCmpPubkey,
        bytes calldata validatorCmpPubkey,
        uint256 amount
    ) external verifyPubkey(delegatorCmpPubkey) verifyExistingValidator(validatorCmpPubkey) {
        _verifyCallerIsOperator(delegatorCmpPubkey, msg.sender);
        _unstake(delegatorCmpPubkey, validatorCmpPubkey, amount);
    }

    /// @dev Unstakes the given amount from the validator for the delegator, where the amount is deposited to the
    /// execution address.
    function _unstake(bytes memory delegatorCmpPubkey, bytes calldata validatorCmpPubkey, uint256 amount) internal {
        require(
            delegatorValidatorStakes[delegatorCmpPubkey][validatorCmpPubkey] >= amount,
            "IPTokenStaking: Insufficient staked amount"
        );

        validatorMetadata[validatorCmpPubkey].totalStake -= amount;
        delegatorTotalStakes[delegatorCmpPubkey] -= amount;
        delegatorValidatorStakes[delegatorCmpPubkey][validatorCmpPubkey] -= amount;

        // If validator gets slashed and the total staked in CL is less than the total staked in EL, then the validator
        // might not exist in CL while still existing in EL.
        if (validatorMetadata[validatorCmpPubkey].totalStake == 0) {
            delete validatorMetadata[validatorCmpPubkey];
        }

        emit Withdraw(delegatorCmpPubkey, validatorCmpPubkey, amount);
    }

    /// @dev Refunds the remainder of the stake amount to the msg sender.
    /// @param remainder The remainder of the stake amount.
    function _refundRemainder(uint256 remainder) internal {
        (bool success, ) = msg.sender.call{ value: remainder }("");
        require(success, "IPTokenStaking: Failed to refund remainder");
    }

    //////////////////////////////// Upgradeability ////////////////////////////////////

    /// @notice Disables the upgradeability of the contract.
    /// @dev WARNING: This action is irreversible.
    function disableUpgradeability() external override onlyOwner {
        _disableUpgradeability();
    }

    /// @dev Hook to authorize the upgrade according to UUPSUpgradeable
    /// @param newImplementation The address of the new implementation
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner upgradeabilityEnabled {}
}
