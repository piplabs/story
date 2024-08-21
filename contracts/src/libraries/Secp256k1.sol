// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

library Secp256k1 {
    /// @notice Compress an uncompressed 65-byte Secp256k1 public key.
    /// @dev Assumes that the input is a valid 65-byte uncompressed public key.
    function compressPublicKey(bytes memory uncompressedKey) public pure returns (bytes memory) {
        require(uncompressedKey.length == 65, "Invalid uncompressed public key length");

        // Extract the X and Y coordinates
        bytes32 x;
        bytes32 y;

        assembly {
            x := mload(add(uncompressedKey, 0x21))
            y := mload(add(uncompressedKey, 0x41))
        }

        // Determine the prefix (0x02 for even Y, 0x03 for odd Y)
        bytes1 prefix = (uint8(y[31]) % 2 == 0) ? bytes1(0x02) : bytes1(0x03);

        // Concatenate the prefix and X coordinate
        bytes memory compressedKey = new bytes(33);
        compressedKey[0] = prefix;

        for (uint256 i = 0; i < 32; i++) {
            compressedKey[i + 1] = x[i];
        }

        return compressedKey;
    }
}
