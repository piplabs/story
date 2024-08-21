// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;

interface IIPTokenSlashing {
    /// @notice Emitted when a request to unjail a validator is made
    /// @param sender The address that sent the unjail request
    /// @param validatorCmpPubkey 33 bytes compressed secp256k1 public key.
    event Unjail(address indexed sender, bytes validatorCmpPubkey);

    /// @notice Emitted when the unjail fee is updated
    /// @param newUnjailFee The new unjail fee
    event UnjailFeeSet(uint256 newUnjailFee);

    /// @notice Requests to unjail the validator. Must pay fee on the execution side to prevent spamming.
    /// @param validatorUncmpPubkey The validator's 65-byte uncompressed Secp256k1 public key
    function unjail(bytes calldata validatorUncmpPubkey) external payable;

    /// @notice Requests to unjail a validator on behalf. Must pay fee on the execution side to prevent spamming.
    /// @param validatorUncmpPubkey The validator's 65-byte uncompressed Secp256k1 public key
    function unjailOnBehalf(bytes calldata validatorUncmpPubkey) external payable;
}
