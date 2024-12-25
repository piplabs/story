// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

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
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param moniker The moniker of the validator.
    /// @param stakeAmount Token staked to the validator as self-delegation.
    /// @param commissionRate The commission rate of the validator.
    /// @param maxCommissionRate The maximum commission rate of the validator.
    /// @param maxCommissionChangeRate The maximum commission change rate of the validator.
    /// @param supportsUnlocked Whether the validator supports unlocked staking
    /// @param operatorAddress The caller's address
    /// @param data Additional data for the validator
    event CreateValidator(
        bytes validatorCmpPubkey,
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
    /// @param delegator The delegator's address
    /// @param executionAddress Left-padded 32 bytes of the EVM address to receive stake and reward withdrawals.
    event SetWithdrawalAddress(address delegator, bytes32 executionAddress);

    /// @notice Emitted when the rewards address is set/changed.
    /// @param delegator The delegator's address
    /// @param executionAddress Left-padded 32 bytes of the EVM address to receive stake and reward withdrawals.
    event SetRewardAddress(address delegator, bytes32 executionAddress);

    /// @notice Emitted when the validator commission is updated
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param commissionRate The new commission rate of the validator.
    event UpdateValidatorCommission(bytes validatorCmpPubkey, uint32 commissionRate);

    /// @notice Emitted when a user deposits token into the contract.
    /// @param delegator The delegator's address
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param stakeAmount Token deposited
    /// @param stakingPeriod The staking period of the deposit
    /// @param delegationId The ID of the delegation
    /// @param operatorAddress The caller's address
    /// @param data Additional data for the deposit
    event Deposit(
        address delegator,
        bytes validatorCmpPubkey,
        uint256 stakeAmount,
        uint256 stakingPeriod,
        uint256 delegationId,
        address operatorAddress,
        bytes data
    );

    /// @notice Emitted when a user withdraws her stake and starts the unbonding period.
    /// @param delegator The delegator's address
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param stakeAmount Token deposited.
    /// @param delegationId The ID of the delegation, 0 if flexible
    /// @param operatorAddress The caller's address
    /// @param data Additional data for the deposit
    event Withdraw(
        address delegator,
        bytes validatorCmpPubkey,
        uint256 stakeAmount,
        uint256 delegationId,
        address operatorAddress,
        bytes data
    );

    /// @notice Emitted when a user triggers redelegation of token from source validator to destination validator.
    /// @param delegator The delegator's address
    /// @param validatorSrcCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param validatorDstCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param delegationId if delegation has staking period, 0 if flexible
    /// @param operatorAddress The caller's address
    /// @param amount Token redelegated.
    event Redelegate(
        address delegator,
        bytes validatorSrcCmpPubkey,
        bytes validatorDstCmpPubkey,
        uint256 delegationId,
        address operatorAddress,
        uint256 amount
    );

    /// @notice Emitted to request setting an operator address to a delegator
    /// @param delegator The delegator's address
    /// @param operator The operator's address
    event SetOperator(address delegator, address operator);

    /// @notice Emitted to request removing the operator address to a delegator
    /// @param delegator The delegator's address
    event UnsetOperator(address delegator);

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
    /// @param unjailer The unjailer's address
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param data Additional data for the unjail.
    event Unjail(address unjailer, bytes validatorCmpPubkey, bytes data);

    /// @notice Returns the rounded stake amount and the remainder.
    /// @param rawAmount The raw stake amount.
    /// @return amount The rounded stake amount.
    /// @return remainder The remainder of the stake amount.
    function roundedStakeAmount(uint256 rawAmount) external view returns (uint256 amount, uint256 remainder);

    /// @notice Sets an operator for a delegator.
    /// Calling this method will override any existing operator.
    /// @param operator The operator address to add.
    function setOperator(address operator) external payable;

    /// @notice Removes current operator for a delegator.
    function unsetOperator() external payable;

    /// @notice Set/Update the withdrawal address that receives the withdrawals.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @param newWithdrawalAddress EVM address to receive the  withdrawals.
    function setWithdrawalAddress(address newWithdrawalAddress) external payable;

    /// @notice Set/Update the withdrawal address that receives the stake and reward withdrawals.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @param newRewardsAddress EVM address to receive the stake and reward withdrawals.
    function setRewardsAddress(address newRewardsAddress) external payable;

    /// @notice Update the commission rate of a validator.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param commissionRate The new commission rate of the validator.
    function updateValidatorCommission(bytes calldata validatorCmpPubkey, uint32 commissionRate) external payable;

    /// @notice Entry point for creating a new validator with self delegation.
    /// @dev The caller must provide the compressed public key that matches the expected EVM address.
    /// Use this method to make sure the caller is the owner of the validator.
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
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
    ) external payable;

    /// @notice Entry point to stake (delegate) to the given validator. The consensus client (CL) is notified of
    /// the deposit and manages the stake accounting and validator onboarding. Payer must be the delegator.
    /// @dev Staking burns tokens in Execution Layer (EL). Unstaking (withdrawal) will trigger minting through
    /// withdrawal queue.
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param stakingPeriod The staking period.
    /// @param data Additional data for the stake.
    /// @return delegationId The delegation ID, always 0 for flexible staking.
    function stake(
        bytes calldata validatorCmpPubkey,
        StakingPeriod stakingPeriod,
        bytes calldata data
    ) external payable returns (uint256 delegationId);

    /// @notice Entry point for staking IP token to stake to the given validator. The consensus chain is notified of
    /// the stake and manages the stake accounting and validator onboarding. Payer can stake on behalf of another user,
    /// who will be the beneficiary of the stake.
    /// @dev Staking burns tokens in Execution Layer (EL). Unstaking (withdrawal) will trigger minting through
    /// withdrawal queue.
    /// @param delegator The delegator's address
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param stakingPeriod The staking period.
    /// @param data Additional data for the stake.
    /// @return delegationId The delegation ID, always 0 for flexible staking.
    function stakeOnBehalf(
        address delegator,
        bytes calldata validatorCmpPubkey,
        StakingPeriod stakingPeriod,
        bytes calldata data
    ) external payable returns (uint256 delegationId);

    /// @notice Entry point for redelegating the stake to another validator.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @dev For non flexible staking, your staking period will continue as is.
    /// @dev For locked tokens, this will fail in CL if the validator doesn't support unlocked staking.
    /// @param validatorSrcCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param validatorDstCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param delegationId The delegation ID, 0 for flexible staking.
    /// @param amount The amount of stake to redelegate.
    function redelegate(
        bytes calldata validatorSrcCmpPubkey,
        bytes calldata validatorDstCmpPubkey,
        uint256 delegationId,
        uint256 amount
    ) external payable;

    /// @notice Entry point for redelegating the stake to another validator on behalf of the delegator.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @dev For non flexible staking, your staking period will continue as is.
    /// @dev For locked tokens, this will fail in CL if the validator doesn't support unlocked staking.
    /// @dev Caller must be the operator for the delegator, set via `setOperator`. The operator check is done in CL, so
    /// this method will succeed even if the caller is not the operator (but will fail in CL).
    /// @param delegator The delegator's address
    /// @param validatorSrcCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param validatorDstCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param delegationId The delegation ID, 0 for flexible staking.
    /// @param amount The amount of stake to redelegate.
    function redelegateOnBehalf(
        address delegator,
        bytes calldata validatorSrcCmpPubkey,
        bytes calldata validatorDstCmpPubkey,
        uint256 delegationId,
        uint256 amount
    ) external payable;

    /// @notice Entry point for unstaking the previously staked token.
    /// @dev Unstake (withdrawal) will trigger native minting, so token in this contract is considered as burned.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param delegationId The delegation ID, 0 for flexible staking.
    /// @param amount Token amount to unstake.
    /// @param data Additional data for the unstake.
    function unstake(
        bytes calldata validatorCmpPubkey,
        uint256 delegationId,
        uint256 amount,
        bytes calldata data
    ) external payable;

    /// @notice Entry point for unstaking the previously staked token on behalf of the delegator.
    /// Charges fee (CL spam prevention). Must be exact amount.
    /// @dev Caller must be the operator for the delegator, set via `setOperator`. The operator check is done in CL, so
    /// this method will succeed even if the caller is not the operator (but will fail in CL).
    /// @param delegator The delegator's address
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param delegationId The delegation ID, 0 for flexible staking.
    /// @param amount Token amount to unstake.
    /// @param data Additional data for the unstake.
    function unstakeOnBehalf(
        address delegator,
        bytes calldata validatorCmpPubkey,
        uint256 delegationId,
        uint256 amount,
        bytes calldata data
    ) external payable;

    /// @notice Requests to unjail the validator. Caller must pay a fee to prevent spamming.
    /// Fee must be exact amount.
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param data Additional data for the unjail.
    function unjail(bytes calldata validatorCmpPubkey, bytes calldata data) external payable;
}
