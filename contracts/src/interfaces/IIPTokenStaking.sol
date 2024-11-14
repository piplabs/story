// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

/// @title IIPTokenStaking
/// @notice Interface for the IPTokenStaking contract
interface IIPTokenStaking {
    /// @notice Enum representing the different staking periods
    /// @dev FLEXIBLE is used for flexible staking, where the staking period is not fixed and can be changed by the user
    /// SHORT, MEDIUM, and LONG are used for staking with specific periods
    enum StakingPeriod {
        FLEXIBLE,
        SHORT,
        MEDIUM,
        LONG
    }

    /// @notice Struct for initialize method args
    /// @dev Contains various parameters for the contract's functionality
    /// @param owner The address of the admin addres
    /// @param minStakeAmount Global minimum amount required to stake
    /// @param minUnstakeAmount Global minimum amount required to unstake
    /// @param minCommissionRate Global minimum commission rate for validators
    /// @param fee The fee charged for adding to CL storage
    struct InitializerArgs {
        address owner;
        uint256 minStakeAmount;
        uint256 minUnstakeAmount;
        uint256 minCommissionRate;
        uint256 fee;
    }

    /// @notice Emitted when the fee charged for adding to CL storage is updated
    /// @param newFee The new fee
    event FeeSet(uint256 newFee);

    /// @notice Emitted when a new validator is created.
    /// @param validatorUncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param moniker The moniker of the validator.
    /// @param stakeAmount Token staked to the validator as self-delegation.
    /// @param commissionRate The commission rate of the validator.
    /// @param maxCommissionRate The maximum commission rate of the validator.
    /// @param maxCommissionChangeRate The maximum commission change rate of the validator.
    /// @param supportsUnlocked Whether the validator supports unlocked staking
    /// @param operatorAddress The caller's address
    /// @param data Additional data for the validator
    event CreateValidator(
        bytes validatorUncmpPubkey,
        string moniker,
        uint256 stakeAmount,
        uint32 commissionRate,
        uint32 maxCommissionRate,
        uint32 maxCommissionChangeRate,
        uint8 supportsUnlocked,
        address operatorAddress,
        bytes data
    );

    /// @notice Emitted when the withdrawal address is set/changed.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param executionAddress Left-padded 32 bytes of the EVM address to receive stake and reward withdrawals.
    event SetWithdrawalAddress(bytes delegatorUncmpPubkey, bytes32 executionAddress);

    /// @notice Emitted when the rewards address is set/changed.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param executionAddress Left-padded 32 bytes of the EVM address to receive stake and reward withdrawals.
    event SetRewardAddress(bytes delegatorUncmpPubkey, bytes32 executionAddress);

    /// @notice Emitted when the validator commission is updated
    /// @param validatorUncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param commissionRate The new commission rate of the validator.
    event UpdateValidatorCommssion(bytes validatorUncmpPubkey, uint32 commissionRate);

    /// @notice Emitted when a user deposits token into the contract.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpPubkey Validator's 65 bytes uncompressed secp256k1 public key.
    /// @param stakeAmount Token deposited.
    /// @param stakingPeriod of the deposit
    /// @param delegationId The ID of the delegation
    /// @param operatorAddress The caller's address
    /// @param data Additional data for the deposit
    event Deposit(
        bytes delegatorUncmpPubkey,
        bytes validatorUncmpPubkey,
        uint256 stakeAmount,
        uint256 stakingPeriod,
        uint256 delegationId,
        address operatorAddress,
        bytes data
    );

    /// @notice Emitted when a user withdraws her stake and starts the unbonding period.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpPubkey Validator's 65 bytes uncompressed secp256k1 public key.
    /// @param stakeAmount Token deposited.
    /// @param delegationId The ID of the delegation, 0 if flexible
    /// @param operatorAddress The caller's address
    /// @param data Additional data for the deposit
    event Withdraw(
        bytes delegatorUncmpPubkey,
        bytes validatorUncmpPubkey,
        uint256 stakeAmount,
        uint256 delegationId,
        address operatorAddress,
        bytes data
    );

    /// @notice Emitted when a user triggers redelegation of token from source validator to destination validator.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpSrcPubkey Source validator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpDstPubkey Destination validator's 65 bytes uncompressed secp256k1 public key.
    /// @param delegationId if delegation has staking period, 0 if flexible
    /// @param operatorAddress The caller's address
    /// @param amount Token redelegated.
    event Redelegate(
        bytes delegatorUncmpPubkey,
        bytes validatorUncmpSrcPubkey,
        bytes validatorUncmpDstPubkey,
        uint256 delegationId,
        address operatorAddress,
        uint256 amount
    );

    /// @notice Emitted to request adding an operator address to a validator
    /// @param uncmpPubkey delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param operator address
    event AddOperator(bytes uncmpPubkey, address operator);

    /// @notice Emitted to request removing an operator address to a validator
    /// @param uncmpPubkey delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param operator address
    event RemoveOperator(bytes uncmpPubkey, address operator);

    /// @notice Emitted when the minimum stake amount is set.
    /// @param minStakeAmount The new minimum stake amount.
    event MinStakeAmountSet(uint256 minStakeAmount);

    /// @notice Emitted when the minimum unstake amount is set.
    /// @param minUnstakeAmount The new minimum unstake amount.
    event MinUnstakeAmountSet(uint256 minUnstakeAmount);

    /// @notice Emitted when the global minimum commission rate is set.
    /// @param minCommissionRate The new global minimum commission rate.
    event MinCommissionRateChanged(uint256 minCommissionRate);

    /// @notice Emitted when a validator is unjailed.
    event Unjail(address unjailer, bytes validatorUncmpPubkey, bytes data);

    /// @notice Returns the rounded stake amount and the remainder.
    /// @param rawAmount The raw stake amount.
    /// @return amount The rounded stake amount.
    /// @return remainder The remainder of the stake amount.
    function roundedStakeAmount(uint256 rawAmount) external view returns (uint256 amount, uint256 remainder);

    /// @notice Adds an operator for a delegator.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param operator The operator address to add.
    function addOperator(bytes calldata uncmpPubkey, address operator) external payable;

    /// @notice Removes an operator for a delegator.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param operator The operator address to remove.
    function removeOperator(bytes calldata uncmpPubkey, address operator) external payable;

    /// @notice Set/Update the withdrawal address that receives the withdrawals.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param newWithdrawalAddress EVM address to receive the  withdrawals.
    function setWithdrawalAddress(bytes calldata delegatorUncmpPubkey, address newWithdrawalAddress) external payable;

    /// @notice Set/Update the withdrawal address that receives the stake and reward withdrawals.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param newRewardsAddress EVM address to receive the stake and reward withdrawals.
    function setRewardsAddress(bytes calldata delegatorUncmpPubkey, address newRewardsAddress) external payable;

    /// @notice Update the commission rate of a validator.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @param validatorUncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param commissionRate The new commission rate of the validator.
    function updateValidatorCommission(bytes calldata validatorUncmpPubkey, uint32 commissionRate) external payable;

    /// @notice Entry point for creating a new validator with self delegation.
    /// @dev The caller must provide the uncompressed public key that matches the expected EVM address.
    /// Use this method to make sure the caller is the owner of the validator.
    /// @param validatorUncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param moniker The moniker of the validator.
    /// @param commissionRate The commission rate of the validator.
    /// @param maxCommissionRate The maximum commission rate of the validator.
    /// @param maxCommissionChangeRate The maximum commission change rate of the validator.
    /// @param supportsUnlocked Whether the validator supports unlocked staking.
    /// @param data Additional data for the validator.
    function createValidator(
        bytes calldata validatorUncmpPubkey,
        string calldata moniker,
        uint32 commissionRate,
        uint32 maxCommissionRate,
        uint32 maxCommissionChangeRate,
        bool supportsUnlocked,
        bytes calldata data
    ) external payable;

    /// @notice Entry point for creating a new validator on behalf of someone else.
    /// WARNING: If validatorUncmpPubkey is wrong, the stake will go to an address that the sender
    /// won't be able to control and unstake from, funds will be lost. If you want to make sure the
    /// caller is the owner of the validator, use createValidator instead.
    /// @param validatorUncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param moniker The moniker of the validator.
    /// @param commissionRate The commission rate of the validator.
    /// @param maxCommissionRate The maximum commission rate of the validator.
    /// @param maxCommissionChangeRate The maximum commission change rate of the validator.
    /// @param supportsUnlocked Whether the validator supports unlocked staking.
    /// @param data Additional data for the validator.
    function createValidatorOnBehalf(
        bytes calldata validatorUncmpPubkey,
        string calldata moniker,
        uint32 commissionRate,
        uint32 maxCommissionRate,
        uint32 maxCommissionChangeRate,
        bool supportsUnlocked,
        bytes calldata data
    ) external payable;

    /// @notice Entry point to stake (delegate) to the given validator. The consensus client (CL) is notified of
    /// the deposit and manages the stake accounting and validator onboarding. Payer must be the delegator.
    /// @dev Staking burns tokens in Execution Layer (EL). Unstaking (withdrawal) will trigger minting through
    /// withdrawal queue.
    /// This method will revert if delegatorUncmpPubkey is not the sender of the validator.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpPubkey Validator's65 bytes uncompressed secp256k1 public key.
    /// @param stakingPeriod The staking period.
    /// @param data Additional data for the stake.
    /// @return delegationId The delegation ID, always 0 for flexible staking.
    function stake(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpPubkey,
        StakingPeriod stakingPeriod,
        bytes calldata data
    ) external payable returns (uint256 delegationId);

    /// @notice Entry point for staking IP token to stake to the given validator. The consensus chain is notified of
    /// the stake and manages the stake accounting and validator onboarding. Payer can stake on behalf of another user,
    /// who will be the beneficiary of the stake.
    /// @dev Staking burns tokens in Execution Layer (EL). Unstaking (withdrawal) will trigger minting through
    /// withdrawal queue.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpPubkey Validator's65 bytes uncompressed secp256k1 public key.
    /// @param stakingPeriod The staking period.
    /// @param data Additional data for the stake.
    /// @return delegationId The delegation ID, always 0 for flexible staking.
    function stakeOnBehalf(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpPubkey,
        IIPTokenStaking.StakingPeriod stakingPeriod,
        bytes calldata data
    ) external payable returns (uint256 delegationId);

    /// @notice Entry point for redelegating the stake to another validator.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @dev For non flexible staking, your staking period will continue as is.
    /// @dev For locked tokens, this will fail in CL if the validator doesn't support unlocked staking.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpSrcPubkey Validator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpDstPubkey Validator's 65 bytes uncompressed secp256k1 public key.
    /// @param delegationId The delegation ID, 0 for flexible staking.
    /// @param amount The amount of stake to redelegate.
    function redelegate(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpSrcPubkey,
        bytes calldata validatorUncmpDstPubkey,
        uint256 delegationId,
        uint256 amount
    ) external payable;

    /// @notice Entry point for redelegating the stake to another validator on behalf of the delegator.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @dev For non flexible staking, your staking period will continue as is.
    /// @dev For locked tokens, this will fail in CL if the validator doesn't support unlocked staking.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpSrcPubkey Validator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpDstPubkey Validator's 65 bytes uncompressed secp256k1 public key.
    /// @param delegationId The delegation ID, 0 for flexible staking.
    /// @param amount The amount of stake to redelegate.
    function redelegateOnBehalf(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpSrcPubkey,
        bytes calldata validatorUncmpDstPubkey,
        uint256 delegationId,
        uint256 amount
    ) external payable;

    /// @notice Entry point for unstaking the previously staked token.
    /// @dev Unstake (withdrawal) will trigger native minting, so token in this contract is considered as burned.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpPubkey Validator's65 bytes uncompressed secp256k1 public key.
    /// @param delegationId The delegation ID, 0 for flexible staking.
    /// @param amount Token amount to unstake.
    /// @param data Additional data for the unstake.
    function unstake(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpPubkey,
        uint256 delegationId,
        uint256 amount,
        bytes calldata data
    ) external payable;

    /// @notice Entry point for unstaking the previously staked token on behalf of the delegator.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @dev Must be an approved operator for the delegator.
    /// @param delegatorUncmpPubkey Delegator's65 bytes uncompressed secp256k1 public key.
    /// @param validatorUncmpPubkey Validator's65 bytes uncompressed secp256k1 public key.
    /// @param delegationId The delegation ID, 0 for flexible staking.
    /// @param amount Token amount to unstake.
    /// @param data Additional data for the unstake.
    function unstakeOnBehalf(
        bytes calldata delegatorUncmpPubkey,
        bytes calldata validatorUncmpPubkey,
        uint256 delegationId,
        uint256 amount,
        bytes calldata data
    ) external payable;

    /// @notice Requests to unjail the validator. Caller must pay a fee to prevent spamming.
    /// Fee must be exact amount.
    /// @param validatorUncmpPubkey The validator's 65-byte uncompressed Secp256k1 public key
    /// @param data Additional data for the unjail.
    function unjail(bytes calldata validatorUncmpPubkey, bytes calldata data) external payable;

    /// @notice Requests to unjail a validator on behalf. Caller must pay a fee to prevent spamming.
    /// Fee must be exact amount.
    /// @param validatorUncmpPubkey The validator's 65-byte uncompressed Secp256k1 public key
    /// @param data Additional data for the unjail.
    function unjailOnBehalf(bytes calldata validatorUncmpPubkey, bytes calldata data) external payable;
}
