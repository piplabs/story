// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { ReentrancyGuardUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import { EnumerableSet } from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import { IIPTokenStaking } from "../interfaces/IIPTokenStaking.sol";
import { Secp256k1Verifier } from "./Secp256k1Verifier.sol";

/**
 * @title IPTokenStaking
 * @notice The deposit contract for IP token staked validators.
 * @dev This contract is a sort of "bridge" to request validator related actions on the consensus chain.
 * The response will happen on the consensus chain.
 * Since most of the validator related actions are executed on the consensus chain, the methods in this contract
 * must be considered requests and not final actions, a successful transaction here does not guarantee the success
 * of the transaction on the consensus chain.
 * NOTE: All $IP tokens staked to this contract will be burned (transferred to the zero address).
 * The flow is as follows:
 * 1. User calls a method in this contract, which will emit an event if checks pass.
 * 2. Modules on the consensus chain are listening for these events and execute the corresponding logic
 * (e.g. staking, create validator, etc.), minting tokens in CL if needed.
 * 3. If the action fails in CL, for example staking on a validator that doesn't exist, the deposited $IP tokens will
 * not be refunded to the user. Remember that the EL transaction of step 2 would not have reverted. So please be
 * cautious when making transactions with this contract.
 */
contract IPTokenStaking is IIPTokenStaking, Ownable2StepUpgradeable, ReentrancyGuardUpgradeable, Secp256k1Verifier {
    using EnumerableSet for EnumerableSet.AddressSet;

    /// @notice Maximum length of the validator moniker, in bytes.
    uint256 public constant MAX_MONIKER_LENGTH = 70;

    /// @notice Stake amount increments. Consensus Layer requires staking in increments of 1 gwei.
    uint256 public constant STAKE_ROUNDING = 1 gwei;

    /// @notice Default minimum fee charged for adding to CL storage
    uint256 public immutable DEFAULT_MIN_FEE;

    /// @notice The maximum size of the data field in the event logs.
    uint256 public immutable MAX_DATA_LENGTH;

    /// @notice Global minimum commission rate for validators
    uint256 public minCommissionRate;

    /// @notice Minimum amount required to stake.
    uint256 public minStakeAmount;

    /// @notice Minimum amount required to unstake.
    uint256 public minUnstakeAmount;

    /// @notice Counter to generate delegationIds for delegations with period.
    /// @dev Starts in 1, since 0 is reserved for flexible delegations.
    uint256 private _delegationIdCounter;

    /// @notice The fee paid to update a validator (unjail, commission update, etc.)
    uint256 public fee;

    modifier chargesFee() {
        require(msg.value == fee, "IPTokenStaking: Invalid fee amount");
        payable(address(0x0)).transfer(msg.value);
        _;
    }

    constructor(uint256 defaultMinFee, uint256 maxDataLength) {
        require(defaultMinFee >= 1 gwei, "IPTokenStaking: Invalid default min fee");
        DEFAULT_MIN_FEE = defaultMinFee;
        MAX_DATA_LENGTH = maxDataLength;
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    /// @dev Only callable once at proxy deployment.
    /// @param args The initializer arguments.
    function initialize(IIPTokenStaking.InitializerArgs calldata args) public initializer {
        __ReentrancyGuard_init();
        __Ownable_init(args.owner);
        _setMinStakeAmount(args.minStakeAmount);
        _setMinUnstakeAmount(args.minUnstakeAmount);
        _setMinCommissionRate(args.minCommissionRate);
        _setFee(args.fee);
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

    /// @notice Sets the fee charged for adding to CL storage.
    /// @param newFee The new fee
    function setFee(uint256 newFee) external onlyOwner {
        _setFee(newFee);
    }

    /// @notice Sets the global minimum commission rate for validators.
    /// @param newValue The new minimum commission rate.
    function setMinCommissionRate(uint256 newValue) external onlyOwner {
        _setMinCommissionRate(newValue);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                            Internal setters                            //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Sets the fee charged for adding to CL storage.
    function _setFee(uint256 newFee) private {
        require(newFee >= DEFAULT_MIN_FEE, "IPTokenStaking: Invalid min fee");
        fee = newFee;
        emit FeeSet(newFee);
    }

    /// @dev Sets the minimum amount required to stake.
    /// @param newMinStakeAmount The minimum amount required to stake.
    function _setMinStakeAmount(uint256 newMinStakeAmount) private {
        minStakeAmount = newMinStakeAmount - (newMinStakeAmount % STAKE_ROUNDING);
        require(minStakeAmount > 0, "IPTokenStaking: Zero min stake amount");
        emit MinStakeAmountSet(minStakeAmount);
    }

    /// @dev Sets the minimum amount required to withdraw.
    /// @param newMinUnstakeAmount The minimum amount required to stake.
    function _setMinUnstakeAmount(uint256 newMinUnstakeAmount) private {
        minUnstakeAmount = newMinUnstakeAmount - (newMinUnstakeAmount % STAKE_ROUNDING);
        require(minUnstakeAmount > 0, "IPTokenStaking: Zero min unstake amount");
        emit MinUnstakeAmountSet(minUnstakeAmount);
    }

    /// @dev Sets the minimum global commission rate for validators.
    /// @param newValue The new minimum commission rate.
    function _setMinCommissionRate(uint256 newValue) private {
        require(newValue > 0, "IPTokenStaking: Zero min commission rate");
        minCommissionRate = newValue;
        emit MinCommissionRateChanged(newValue);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                            Operator functions                          //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Sets an operator for a delegator.
    /// Calling this method will override any existing operator.
    /// @param operator The operator address to add.
    function setOperator(address operator) external payable chargesFee {
        // Use unsetOperator to remove an operator, not setting to zero address
        require(operator != address(0), "IPTokenStaking: zero input address");
        emit SetOperator(msg.sender, operator);
    }

    /// @notice Removes current operator for a delegator.
    function unsetOperator() external payable chargesFee {
        emit UnsetOperator(msg.sender);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                     Staking Configuration functions                    //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Set/Update the withdrawal address that receives the withdrawals.
    /// @param newWithdrawalAddress EVM address to receive the  withdrawals.
    function setWithdrawalAddress(address newWithdrawalAddress) external payable chargesFee {
        require(newWithdrawalAddress != address(0), "IPTokenStaking: zero input address");
        emit SetWithdrawalAddress({
            delegator: msg.sender,
            executionAddress: bytes32(uint256(uint160(newWithdrawalAddress))) // left-padded bytes32 of the address
        });
    }

    /// @notice Set/Update the withdrawal address that receives the stake and reward withdrawals.
    /// @dev To prevent spam, only delegators with stake can call this function with cool-down time.
    /// @param newRewardsAddress EVM address to receive the stake and reward withdrawals.
    function setRewardsAddress(address newRewardsAddress) external payable chargesFee {
        require(newRewardsAddress != address(0), "IPTokenStaking: zero input address");
        emit SetRewardAddress({
            delegator: msg.sender,
            executionAddress: bytes32(uint256(uint160(newRewardsAddress))) // left-padded bytes32 of the address
        });
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                          Validator Creation                            //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Entry point for creating a new validator with self delegation.
    /// @dev The caller must provide the compressed public key that matches the expected EVM address.
    /// Use this method to make sure the caller is the owner of the validator.
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key of validator.
    /// @param moniker The moniker of the validator.
    /// @param commissionRate The commission rate of the validator.
    /// @param maxCommissionRate The maximum commission rate of the validator.
    /// @param maxCommissionChangeRate The maximum commission change rate of the validator.
    /// @param supportsUnlocked Whether the validator supports unlocked staking.
    /// @param data Additional data for the validator.
    function createValidator(
        bytes calldata validatorCmpPubkey,
        string calldata moniker,
        uint32 commissionRate,
        uint32 maxCommissionRate,
        uint32 maxCommissionChangeRate,
        bool supportsUnlocked,
        bytes calldata data
    ) external payable verifyCmpPubkeyWithExpectedAddress(validatorCmpPubkey, msg.sender) nonReentrant {
        _createValidator(
            validatorCmpPubkey,
            moniker,
            commissionRate,
            maxCommissionRate,
            maxCommissionChangeRate,
            supportsUnlocked,
            data
        );
    }

    /// @dev Validator is the delegator when creating a new validator (self-delegation).
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key of validator.
    /// @param moniker The moniker of the validator.
    /// @param commissionRate The commission rate of the validator.
    /// @param maxCommissionRate The maximum commission rate of the validator.
    /// @param maxCommissionChangeRate The maximum commission change rate of the validator.
    /// @param supportsUnlocked Whether the validator supports unlocked staking.
    /// @param data Additional data for the validator.
    function _createValidator(
        bytes calldata validatorCmpPubkey,
        string memory moniker,
        uint32 commissionRate,
        uint32 maxCommissionRate,
        uint32 maxCommissionChangeRate,
        bool supportsUnlocked,
        bytes calldata data
    ) internal {
        (uint256 stakeAmount, uint256 remainder) = roundedStakeAmount(msg.value);
        require(stakeAmount >= minStakeAmount, "IPTokenStaking: Stake amount under min");
        require(commissionRate >= minCommissionRate, "IPTokenStaking: Commission rate under min");
        require(commissionRate <= maxCommissionRate, "IPTokenStaking: Commission rate over max");
        require(bytes(moniker).length <= MAX_MONIKER_LENGTH, "IPTokenStaking: Moniker length over max");
        require(data.length <= MAX_DATA_LENGTH, "IPTokenStaking: Data length over max");

        payable(address(0)).transfer(stakeAmount);
        emit CreateValidator(
            validatorCmpPubkey,
            moniker,
            stakeAmount,
            commissionRate,
            maxCommissionRate,
            maxCommissionChangeRate,
            supportsUnlocked ? 1 : 0,
            msg.sender,
            data
        );
        if (remainder > 0) {
            _refundRemainder(remainder);
        }
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Validator Config                             //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Update the commission rate of a validator.
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key of validator.
    /// @param commissionRate The new commission rate of the validator.
    function updateValidatorCommission(
        bytes calldata validatorCmpPubkey,
        uint32 commissionRate
    ) external payable chargesFee verifyCmpPubkeyWithExpectedAddress(validatorCmpPubkey, msg.sender) {
        require(commissionRate >= minCommissionRate, "IPTokenStaking: Commission rate under min");
        emit UpdateValidatorCommission(validatorCmpPubkey, commissionRate);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Token Staking                              //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Entry point to stake (delegate) to the given validator. The consensus client (CL) is notified of
    /// the deposit and manages the stake accounting and validator onboarding. Payer must be the delegator.
    /// @dev Staking burns tokens in Execution Layer (EL). Unstaking (withdrawal) will trigger minting through
    /// withdrawal queue.
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key of validator.
    /// @param stakingPeriod The staking period.
    /// @param data Additional data for the stake.
    /// @return delegationId The delegation ID, always 0 for flexible staking.
    function stake(
        bytes calldata validatorCmpPubkey,
        IIPTokenStaking.StakingPeriod stakingPeriod,
        bytes calldata data
    ) external payable nonReentrant returns (uint256 delegationId) {
        return _stake(msg.sender, validatorCmpPubkey, stakingPeriod, data);
    }

    /// @notice Entry point for staking IP token to stake to the given validator. The consensus chain is notified of
    /// the stake and manages the stake accounting and validator onboarding. Payer can stake on behalf of another user,
    /// who will be the beneficiary of the stake.
    /// @dev Staking burns tokens in Execution Layer (EL). Unstaking (withdrawal) will trigger minting through
    /// withdrawal queue.
    /// @param delegator The delegator's address
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key of validator.
    /// @param stakingPeriod The staking period.
    /// @param data Additional data for the stake.
    /// @return delegationId The delegation ID, always 0 for flexible staking.
    function stakeOnBehalf(
        address delegator,
        bytes calldata validatorCmpPubkey,
        IIPTokenStaking.StakingPeriod stakingPeriod,
        bytes calldata data
    ) external payable nonReentrant returns (uint256 delegationId) {
        return _stake(delegator, validatorCmpPubkey, stakingPeriod, data);
    }

    /// @dev Creates a validator (x/staking.MsgCreateValidator) if it does not exist. Then delegates the stake to the
    /// validator (x/staking.MsgDelegate).
    /// @param delegator The delegator's address
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key of validator.
    /// @param stakingPeriod The staking period.
    /// @param data Additional data for the stake.
    /// @return delegationId The delegation ID, always 0 for flexible staking.
    function _stake(
        address delegator,
        bytes calldata validatorCmpPubkey,
        IIPTokenStaking.StakingPeriod stakingPeriod,
        bytes calldata data
    ) internal verifyCmpPubkey(validatorCmpPubkey) returns (uint256) {
        require(delegator != address(0), "IPTokenStaking: invalid delegator");
        require(data.length <= MAX_DATA_LENGTH, "IPTokenStaking: Data length over max");
        // This can't be tested from Foundry (Solidity), but can be triggered from js/rpc
        require(stakingPeriod <= IIPTokenStaking.StakingPeriod.LONG, "IPTokenStaking: Invalid staking period");
        (uint256 stakeAmount, uint256 remainder) = roundedStakeAmount(msg.value);
        require(stakeAmount >= minStakeAmount, "IPTokenStaking: Stake amount under min");

        uint256 delegationId = 0;
        if (stakingPeriod != IIPTokenStaking.StakingPeriod.FLEXIBLE) {
            delegationId = ++_delegationIdCounter;
        }
        emit Deposit(delegator, validatorCmpPubkey, stakeAmount, uint8(stakingPeriod), delegationId, msg.sender, data);
        // We burn staked tokens
        payable(address(0)).transfer(stakeAmount);

        if (remainder > 0) {
            _refundRemainder(remainder);
        }

        return delegationId;
    }

    /// @notice Entry point for redelegating the stake to another validator.
    /// @dev For non flexible staking, your staking period will continue as is.
    /// @dev For locked tokens, this will fail in CL if the validator doesn't support unlocked staking.
    /// @param validatorSrcCmpPubkey 33 bytes compressed secp256k1 public key of source validator.
    /// @param validatorDstCmpPubkey 33 bytes compressed secp256k1 public key of destination validator.
    /// @param delegationId The delegation ID, 0 for flexible staking.
    /// @param amount The amount of stake to redelegate.
    function redelegate(
        bytes calldata validatorSrcCmpPubkey,
        bytes calldata validatorDstCmpPubkey,
        uint256 delegationId,
        uint256 amount
    ) external payable chargesFee {
        _redelegate(msg.sender, validatorSrcCmpPubkey, validatorDstCmpPubkey, delegationId, amount);
    }

    /// @notice Entry point for redelegating the stake to another validator on behalf of the delegator.
    /// @dev For non flexible staking, your staking period will continue as is.
    /// @dev For locked tokens, this will fail in CL if the validator doesn't support unlocked staking.
    /// @dev Caller must be the operator for the delegator, set via `setOperator`. The operator check is done in CL, so
    /// this method will succeed even if the caller is not the operator (but will fail in CL).
    /// @param delegator The delegator's address
    /// @param validatorSrcCmpPubkey 33 bytes compressed secp256k1 public key of source validator.
    /// @param validatorDstCmpPubkey 33 bytes compressed secp256k1 public key of destination validator.
    /// @param delegationId The delegation ID, 0 for flexible staking.
    /// @param amount The amount of stake to redelegate.
    function redelegateOnBehalf(
        address delegator,
        bytes calldata validatorSrcCmpPubkey,
        bytes calldata validatorDstCmpPubkey,
        uint256 delegationId,
        uint256 amount
    ) external payable chargesFee {
        _redelegate(delegator, validatorSrcCmpPubkey, validatorDstCmpPubkey, delegationId, amount);
    }

    function _redelegate(
        address delegator,
        bytes calldata validatorSrcCmpPubkey,
        bytes calldata validatorDstCmpPubkey,
        uint256 delegationId,
        uint256 amount
    ) private verifyCmpPubkey(validatorSrcCmpPubkey) verifyCmpPubkey(validatorDstCmpPubkey) {
        require(delegator != address(0), "IPTokenStaking: Invalid delegator");
        require(
            keccak256(validatorSrcCmpPubkey) != keccak256(validatorDstCmpPubkey),
            "IPTokenStaking: Redelegating to same validator"
        );
        require(delegationId <= _delegationIdCounter, "IPTokenStaking: Invalid delegation id");
        (uint256 stakeAmount, ) = roundedStakeAmount(amount);
        require(stakeAmount >= minStakeAmount, "IPTokenStaking: Stake amount under min");

        emit Redelegate(delegator, validatorSrcCmpPubkey, validatorDstCmpPubkey, delegationId, msg.sender, stakeAmount);
    }

    /// @notice Returns the rounded stake amount and the remainder.
    /// @param rawAmount The raw stake amount.
    /// @return amount The rounded stake amount.
    /// @return remainder The remainder of the stake amount.
    function roundedStakeAmount(uint256 rawAmount) public pure returns (uint256 amount, uint256 remainder) {
        remainder = rawAmount % STAKE_ROUNDING;
        amount = rawAmount - remainder;
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Unstake                                    //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Entry point for unstaking the previously staked token.
    /// @dev Unstake (withdrawal) will trigger native minting, so token in this contract is considered as burned.
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key of validator.
    /// @param delegationId The delegation ID, 0 for flexible staking.
    /// @param amount Token amount to unstake.
    /// @param data Additional data for the unstake.
    function unstake(
        bytes calldata validatorCmpPubkey,
        uint256 delegationId,
        uint256 amount,
        bytes calldata data
    ) external payable chargesFee {
        _unstake(msg.sender, validatorCmpPubkey, delegationId, amount, data);
    }

    /// @notice Entry point for unstaking the previously staked token on behalf of the delegator.
    /// NOTE: If the amount is not divisible by STAKE_ROUNDING, it will be rounded down.
    /// @dev Caller must be the operator for the delegator, set via `setOperator`. The operator check is done in CL, so
    /// this method will succeed even if the caller is not the operator (but will fail in CL)
    /// @param delegator The delegator's address
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key of validator.
    /// @param delegationId The delegation ID, 0 for flexible staking.
    /// @param amount Token amount to unstake. This amount will be rounded to STAKE_ROUNDING.
    /// @param data Additional data for the unstake.
    function unstakeOnBehalf(
        address delegator,
        bytes calldata validatorCmpPubkey,
        uint256 delegationId,
        uint256 amount,
        bytes calldata data
    ) external payable chargesFee {
        _unstake(delegator, validatorCmpPubkey, delegationId, amount, data);
    }

    function _unstake(
        address delegator,
        bytes calldata validatorCmpPubkey,
        uint256 delegationId,
        uint256 amount,
        bytes calldata data
    ) private verifyCmpPubkey(validatorCmpPubkey) {
        require(delegationId <= _delegationIdCounter, "IPTokenStaking: Invalid delegation id");
        (uint256 unstakeAmount, ) = roundedStakeAmount(amount);
        require(unstakeAmount >= minUnstakeAmount, "IPTokenStaking: Unstake amount under min");
        require(data.length <= MAX_DATA_LENGTH, "IPTokenStaking: Data length over max");

        emit Withdraw(delegator, validatorCmpPubkey, unstakeAmount, delegationId, msg.sender, data);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Unjail                                    //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Requests to unjail the validator. Caller must pay a fee to prevent spamming.
    /// Fee must be exact amount.
    /// @param data Additional data for the unjail.
    function unjail(
        bytes calldata validatorCmpPubkey,
        bytes calldata data
    ) external payable chargesFee verifyCmpPubkeyWithExpectedAddress(validatorCmpPubkey, msg.sender) {
        require(data.length <= MAX_DATA_LENGTH, "IPTokenStaking: Data length over max");
        emit Unjail(msg.sender, validatorCmpPubkey, data);
    }

    /// @notice Requests to unjail the validator on behalf of the delegator.
    /// @dev Must be an approved operator for the delegator.
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key of validator.
    /// @param data Additional data for the unjail.
    function unjailOnBehalf(
        bytes calldata validatorCmpPubkey,
        bytes calldata data
    ) external payable chargesFee verifyCmpPubkey(validatorCmpPubkey) {
        require(data.length <= MAX_DATA_LENGTH, "IPTokenStaking: Data length over max");
        emit Unjail(msg.sender, validatorCmpPubkey, data);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Helpers                                    //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Refunds the remainder of the stake amount to the msg sender.
    /// WARNING: Methods using this function should have nonReentrant modifier
    /// to prevent potential reentrancy attacks.
    /// @param remainder The remainder of the stake amount.
    function _refundRemainder(uint256 remainder) private {
        (bool success, ) = msg.sender.call{ value: remainder }("");
        require(success, "IPTokenStaking: Failed to refund remainder");
    }
}
