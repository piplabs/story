// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;

interface IIPTokenStaking {
    /// @param exists Set to true once the validator is created.
    /// @param moniker The moniker of the validator.
    /// @param totalStake Total amount of stake for the validator.
    /// @param commissionRate The commission rate of the validator.
    /// @param maxCommissionRate The maximum commission rate of the validator.
    /// @param maxCommissionChangeRate The maximum commission change rate of the validator.
    struct ValidatorMetadata {
        bool exists;
        string moniker;
        uint256 totalStake;
        uint32 commissionRate;
        uint32 maxCommissionRate;
        uint32 maxCommissionChangeRate;
    }

    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorCmpSrcPubkey Source validator's 33 bytes compressed secp256k1 public key.
    /// @param validatorCmpDstPubkey Destination validator's 33 bytes compressed secp256k1 public key.
    /// @param amount Token amount to redelegate.
    struct RedelegateParams {
        bytes delegatorUncmpPubkey;
        bytes validatorCmpSrcPubkey;
        bytes validatorCmpDstPubkey;
        uint256 amount;
    }

    /// @notice Emitted when a new validator is created.
    /// @param validatorUncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    /// @param moniker The moniker of the validator.
    /// @param stakeAmount Token staked to the validator as self-delegation.
    /// @param commissionRate The commission rate of the validator.
    /// @param maxCommissionRate The maximum commission rate of the validator.
    /// @param maxCommissionChangeRate The maximum commission change rate of the validator.
    event CreateValidator(
        bytes validatorUncmpPubkey,
        bytes validatorCmpPubkey,
        string moniker,
        uint256 stakeAmount,
        uint32 commissionRate,
        uint32 maxCommissionRate,
        uint32 maxCommissionChangeRate
    );

    /// @notice Emitted when the withdrawal address is set/changed.
    /// @param delegatorCmpPubkey Delegator's 33 bytes compressed secp256k1 public key.
    /// @param executionAddress Left-padded 32 bytes of the EVM address to receive stake and reward withdrawals.
    event SetWithdrawalAddress(bytes delegatorCmpPubkey, bytes32 executionAddress);

    /// @notice Emitted when a user deposits token into the contract.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param delegatorCmpPubkey Delegator's 33 bytes compressed secp256k1 public key.
    /// @param validatorCmpPubkey Validator's 33 bytes compressed secp256k1 public key.
    /// @param amount Token deposited.
    event Deposit(bytes delegatorUncmpPubkey, bytes delegatorCmpPubkey, bytes validatorCmpPubkey, uint256 amount);

    /// @notice Emitted when a user triggers redelegation of token from source validator to destination validator.
    /// @param delegatorCmpPubkey Delegator's 33 bytes compressed secp256k1 public key.
    /// @param validatorSrcPubkey Source validator's 33 bytes compressed secp256k1 public key.
    /// @param validatorDstPubkey Destination validator's 33 bytes compressed secp256k1 public key.
    /// @param amount Token redelegated.
    event Redelegate(bytes delegatorCmpPubkey, bytes validatorSrcPubkey, bytes validatorDstPubkey, uint256 amount);

    /// @notice Emitted when a user withdraws her stake and starts the unbonding period.
    /// @param delegatorCmpPubkey Delegator's 33 bytes compressed secp256k1 public key.
    /// @param validatorCmpPubkey Validator's 33 bytes compressed secp256k1 public key.
    /// @param amount Token deposited.
    event Withdraw(bytes delegatorCmpPubkey, bytes validatorCmpPubkey, uint256 amount);

    /// @notice Emitted when the minimum stake amount is set.
    /// @param minStakeAmount The new minimum stake amount.
    event MinStakeAmountSet(uint256 minStakeAmount);

    /// @notice Emitted when the minimum unstake amount is set.
    /// @param minUnstakeAmount The new minimum unstake amount.
    event MinUnstakeAmountSet(uint256 minUnstakeAmount);

    /// @notice Emitted when the minimum redelegation amount is set.
    /// @param minRedelegateAmount The new minimum redelegation amount.
    event MinRedelegateAmountSet(uint256 minRedelegateAmount);

    /// @notice Emitted when the unbonding period is set.
    /// @param newInterval The new unbonding period.
    event WithdrawalAddressChangeIntervalSet(uint256 newInterval);

    /// @notice Returns the rounded stake amount and the remainder.
    /// @param rawAmount The raw stake amount.
    /// @return amount The rounded stake amount.
    /// @return remainder The remainder of the stake amount.
    function roundedStakeAmount(uint256 rawAmount) external view returns (uint256 amount, uint256 remainder);

    /// @notice Returns the operators for the delegator.
    /// @param pubkey 33 bytes compressed secp256k1 public key.
    function getOperators(bytes calldata pubkey) external view returns (address[] memory);

    /// @notice Adds an operator for the delegator.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param operator The operator address to add.
    function addOperator(bytes calldata uncmpPubkey, address operator) external;

    /// @notice Removes an operator for the delegator.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key.
    /// @param operator The operator address to remove.
    function removeOperator(bytes calldata uncmpPubkey, address operator) external;

    /// @notice Set/Update the withdrawal address that receives the stake and reward withdrawals.
    /// @dev To prevent spam, only delegators with stake can call this function with cool-down time.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param newWithdrawalAddress EVM address to receive the stake and reward withdrawals.
    function setWithdrawalAddress(bytes calldata delegatorUncmpPubkey, address newWithdrawalAddress) external;

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
    ) external payable;

    /// @notice Entry point for creating a new validator with self delegation on behalf of the validator.
    /// @dev There's no minimum amount required to stake when creating a new validator.
    /// @param validatorUncmpPubkey 65 bytes uncompressed secp256k1 public key.
    function createValidatorOnBehalf(bytes calldata validatorUncmpPubkey) external payable;

    /// @notice Entry point for staking IP token to stake to the given validator. The consensus chain is notified of
    /// the deposit and manages the stake accounting and validator onboarding. Payer must be the delegator.
    /// @dev When staking, consider it as BURNING. Unstaking (withdrawal) will trigger native minting.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorCmpPubkey Validator's 33 bytes compressed secp256k1 public key.
    function stake(bytes calldata delegatorUncmpPubkey, bytes calldata validatorCmpPubkey) external payable;

    /// @notice Entry point for staking IP token to stake to the given validator. The consensus chain is notified of
    /// the stake and manages the stake accounting and validator onboarding. Payer can stake on behalf of another user,
    /// who will be the beneficiary of the stake.
    /// @dev When staking, consider it as BURNING. Unstaking (withdrawal) will trigger native minting.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorCmpPubkey Validator's 33 bytes compressed secp256k1 public key.
    function stakeOnBehalf(bytes calldata delegatorUncmpPubkey, bytes calldata validatorCmpPubkey) external payable;

    // TODO: Redelegation also requires unbonding period to be executed. Should we separate storage for this for el?
    /// @notice Entry point for redelegating the staked token.
    /// @dev Redelegateion redelegates staked token from src validator to dst validator (x/staking.MsgBeginRedelegate)
    /// @param p See RedelegateParams
    function redelegate(RedelegateParams calldata p) external;

    /// @notice Entry point for redelegating the staked token on behalf of the delegator.
    /// @dev Redelegateion redelegates staked token from src validator to dst validator (x/staking.MsgBeginRedelegate)
    /// @param p See RedelegateParams
    function redelegateOnBehalf(RedelegateParams calldata p) external;

    /// @notice Entry point for unstaking the previously staked token.
    /// @dev Unstake (withdrawal) will trigger native minting, so token in this contract is considered as burned.
    /// @param delegatorUncmpPubkey Delegator's 65 bytes uncompressed secp256k1 public key.
    /// @param validatorCmpPubkey Validator's 33 bytes compressed secp256k1 public key.
    /// @param amount Token amount to unstake.
    function unstake(bytes calldata delegatorUncmpPubkey, bytes calldata validatorCmpPubkey, uint256 amount) external;

    /// @notice Entry point for unstaking the previously staked token on behalf of the delegator.
    /// @dev Must be an approved operator for the delegator.
    /// @param delegatorCmpPubkey Delegator's 33 bytes compressed secp256k1 public key.
    /// @param validatorCmpPubkey Validator's 33 bytes compressed secp256k1 public key.
    /// @param amount Token amount to unstake.
    function unstakeOnBehalf(
        bytes calldata delegatorCmpPubkey,
        bytes calldata validatorCmpPubkey,
        uint256 amount
    ) external;
}
