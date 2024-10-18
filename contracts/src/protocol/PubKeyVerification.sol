// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

import { Errors } from "../libraries/Errors.sol";

/**
 * @title PubKeyVerification
 * @notice Utility functions for pubkey verification
 */
abstract contract PubKeyVerification {
    /// @notice Verifies that the syntax of the given public key is a 65 byte uncompressed secp256k1 public key.
    modifier verifyUncmpPubkey(bytes calldata uncmpPubkey) {
        _verifyUncmpPubkey(uncmpPubkey);
        _;
    }

    /// @notice Verifies that the given 65 byte uncompressed secp256k1 public key (with 0x04 prefix) is valid and
    /// matches the expected EVM address.
    modifier verifyUncmpPubkeyWithExpectedAddress(bytes calldata uncmpPubkey, address expectedAddress) {
        if (uncmpPubkey.length != 65) {
            revert Errors.PubKeyVerifier__InvalidPubkeyLength();
        }
        if (uncmpPubkey[0] != 0x04) {
            revert Errors.PubKeyVerifier__InvalidPubkeyPrefix();
        }
        if (_uncmpPubkeyToAddress(uncmpPubkey) != expectedAddress) {
            revert Errors.PubKeyVerifier__InvalidPubkeyDerivedAddress();
        }
        _;
    }

    /// @notice Verifies that the syntax of the given public key is a 65 byte uncompressed secp256k1 public key.
    function _verifyUncmpPubkey(bytes calldata uncmpPubkey) internal pure {
        if (uncmpPubkey.length != 65) {
            revert Errors.PubKeyVerifier__InvalidPubkeyLength();
        }
        if (uncmpPubkey[0] != 0x04) {
            revert Errors.PubKeyVerifier__InvalidPubkeyPrefix();
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
