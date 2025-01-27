// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { EllipticCurve } from "elliptic-curve-solidity/contracts/EllipticCurve.sol";

/**
 * @title Secp256k1Verifier
 * @notice Utility functions for secp256k1 public key verification
 */
abstract contract Secp256k1Verifier {
    /// @notice Curve parameter a
    uint256 public constant AA = 0;

    /// @notice Curve parameter b
    uint256 public constant BB = 7;

    /// @notice Prime field modulus
    uint256 public constant PP = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F;

    /// @notice Verifies that the syntax of the given public key is a 33 byte compressed secp256k1 public key.
    modifier verifyCmpPubkey(bytes calldata cmpPubkey) {
        bytes memory uncmpPubkey = _uncompressPublicKey(cmpPubkey);
        _verifyUncmpPubkey(uncmpPubkey);
        _;
    }

    /// @notice Verifies that the given 33 byte compressed secp256k1 public key is valid and
    /// matches the expected EVM address.
    modifier verifyCmpPubkeyWithExpectedAddress(bytes calldata cmpPubkey, address expectedAddress) {
        bytes memory uncmpPubkey = _uncompressPublicKey(cmpPubkey);
        _verifyUncmpPubkey(uncmpPubkey);
        require(
            _uncmpPubkeyToAddress(uncmpPubkey) == expectedAddress,
            "Secp256k1Verifier: Invalid pubkey derived address"
        );
        _;
    }

    /// @notice Verifies that the given public key is a 33 byte compressed secp256k1 public key on the curve.
    /// @param cmpPubkey The compressed 33-byte public key to validate
    function _verifyCmpPubkey(bytes memory cmpPubkey) internal pure {
        bytes memory uncmpPubkey = _uncompressPublicKey(cmpPubkey);
        _verifyUncmpPubkey(uncmpPubkey);
    }

    /// @notice Uncompress a compressed 33-byte Secp256k1 public key.
    /// @dev Uses EllipticCurve.deriveY to recover the Y coordinate
    function _uncompressPublicKey(bytes memory cmpPubkey) internal pure returns (bytes memory) {
        require(cmpPubkey.length == 33, "Secp256k1Verifier: Invalid cmp pubkey length");
        require(cmpPubkey[0] == 0x02 || cmpPubkey[0] == 0x03, "Secp256k1Verifier: Invalid cmp pubkey prefix");

        // Extract X coordinate
        uint256 x;
        assembly {
            x := mload(add(cmpPubkey, 0x21))
        }
        uint8 prefix = uint8(cmpPubkey[0]);
        // Derive Y coordinate
        uint256 y = EllipticCurve.deriveY(prefix, x, AA, BB, PP);

        // Construct uncompressed key
        bytes memory uncmpPubkey = new bytes(65);
        uncmpPubkey[0] = 0x04;
        assembly {
            mstore(add(uncmpPubkey, 0x21), x)
            mstore(add(uncmpPubkey, 0x41), y)
        }
        return uncmpPubkey;
    }

    /// @notice Verifies that the given public key is a 65 byte uncompressed secp256k1 public key on the curve.
    /// @param uncmpPubkey The uncompressed 65-byte public key to validate
    function _verifyUncmpPubkey(bytes memory uncmpPubkey) internal pure {
        require(uncmpPubkey.length == 65, "Secp256k1Verifier: Invalid uncmp pubkey length");
        require(uncmpPubkey[0] == 0x04, "Secp256k1Verifier: Invalid uncmp pubkey prefix");

        // Extract x and y coordinates
        uint256 x;
        uint256 y;
        assembly {
            x := mload(add(uncmpPubkey, 0x21))
            y := mload(add(uncmpPubkey, 0x41))
        }

        // Verify the derived point lies on the curve
        require(EllipticCurve.isOnCurve(x, y, AA, BB, PP), "Secp256k1Verifier: pubkey not on curve");
    }

    /// @notice Converts the given public key to an EVM address.
    /// @dev Assume all calls to this function passes in the uncompressed public key.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key, with prefix 04.
    /// @return address The EVM address derived from the public key.
    function _uncmpPubkeyToAddress(bytes memory uncmpPubkey) internal pure returns (address) {
        // Create a new bytes memory array with length 64 (65-1 to skip prefix)
        bytes memory pubkeyNoPrefix = new bytes(64);

        // Copy bytes after prefix using assembly
        assembly {
            // Copy 64 bytes starting from position 1 of input
            // to position 0 of output
            let srcPtr := add(add(uncmpPubkey, 0x20), 1) // Skip first byte
            let destPtr := add(pubkeyNoPrefix, 0x20)
            mstore(destPtr, mload(srcPtr))
            mstore(add(destPtr, 0x20), mload(add(srcPtr, 0x20)))
        }

        return address(uint160(uint256(keccak256(pubkeyNoPrefix))));
    }
}
