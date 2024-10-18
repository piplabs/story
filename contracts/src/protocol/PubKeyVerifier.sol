// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

/**
 * @title PubKeyVerifier
 * @notice Utility functions for pubkey verification
 */
abstract contract PubKeyVerifier {
    /// @notice Verifies that the syntax of the given public key is a 65 byte uncompressed secp256k1 public key.
    modifier verifyUncmpPubkey(bytes calldata uncmpPubkey) {
        _verifyUncmpPubkey(uncmpPubkey);
        _;
    }

    /// @notice Verifies that the given 65 byte uncompressed secp256k1 public key (with 0x04 prefix) is valid and
    /// matches the expected EVM address.
    modifier verifyUncmpPubkeyWithExpectedAddress(bytes calldata uncmpPubkey, address expectedAddress) {
        _verifyUncmpPubkey(uncmpPubkey);
        require(
            _uncmpPubkeyToAddress(uncmpPubkey) == expectedAddress,
            "PubKeyVerifier: Invalid pubkey derived address"
        );
        _;
    }

    /// @notice Verifies that the syntax of the given public key is a 65 byte uncompressed secp256k1 public key.
    function _verifyUncmpPubkey(bytes calldata uncmpPubkey) internal pure {
        require(uncmpPubkey.length == 65, "PubKeyVerifier: Invalid pubkey length");
        require(uncmpPubkey[0] == 0x04, "PubKeyVerifier: Invalid pubkey prefix");
    }

    /// @notice Converts the given public key to an EVM address.
    /// @dev Assume all calls to this function passes in the uncompressed public key.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key, with prefix 04.
    /// @return address The EVM address derived from the public key.
    function _uncmpPubkeyToAddress(bytes calldata uncmpPubkey) internal pure returns (address) {
        return address(uint160(uint256(keccak256(uncmpPubkey[1:]))));
    }
}
