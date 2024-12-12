// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { EllipticCurve } from "elliptic-curve-solidity/contracts/EllipticCurve.sol";

/**
 * @title PubKeyVerifier
 * @notice Utility functions for pubkey verification
 */
abstract contract PubKeyVerifier {
    /// @notice Curve parameter a
    uint256 public constant AA = 0;

    /// @notice Curve parameter b
    uint256 public constant BB = 7;

    /// @notice Prime field modulus
    uint256 public constant PP = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F;

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

    /// @notice Verifies that the given public key is a 65 byte uncompressed secp256k1 public key on the curve.
    /// @param uncmpPubkey The uncompressed 65-byte public key to validate
    function _verifyUncmpPubkey(bytes calldata uncmpPubkey) internal pure {
        require(uncmpPubkey.length == 65, "PubKeyVerifier: Invalid pubkey length");
        require(uncmpPubkey[0] == 0x04, "PubKeyVerifier: Invalid pubkey prefix");

        // Extract x and y coordinates
        uint256 x;
        uint256 y;
        assembly {
            let xPtr := add(uncmpPubkey.offset, 1)
            let yPtr := add(uncmpPubkey.offset, 33)
            x := calldataload(xPtr)
            y := calldataload(yPtr)
        }

        // Verify the derived point lies on the curve
        require(EllipticCurve.isOnCurve(x, y, AA, BB, PP), "PubKeyVerifier: Invalid pubkey on curve");
    }

    /// @notice Converts the given public key to an EVM address.
    /// @dev Assume all calls to this function passes in the uncompressed public key.
    /// @param uncmpPubkey 65 bytes uncompressed secp256k1 public key, with prefix 04.
    /// @return address The EVM address derived from the public key.
    function _uncmpPubkeyToAddress(bytes calldata uncmpPubkey) internal pure returns (address) {
        return address(uint160(uint256(keccak256(uncmpPubkey[1:]))));
    }
}
