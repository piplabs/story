// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

library BytesUtils {
    /// @dev Copies a substring into a new byte string
    /// @param self The byte string to copy from
    /// @param offset The offset to start copying at
    /// @param len The number of bytes to copy
    /// @return The new byte string
    function substring(bytes memory self, uint256 offset, uint256 len) internal pure returns (bytes memory) {
        require(offset + len <= self.length);

        bytes memory ret = new bytes(len);
        uint256 dest;
        uint256 src;

        assembly {
            dest := add(ret, 32)
            src := add(add(self, 32), offset)
        }
        memcpy(dest, src, len);

        return ret;
    }

    function memcpy(uint256 dest, uint256 src, uint256 len) private pure {
        // Copy word-length chunks while possible
        for (; len >= 32; len -= 32) {
            assembly {
                mstore(dest, mload(src))
            }
            dest += 32;
            src += 32;
        }

        // Copy remaining bytes
        uint256 mask;
        if (len == 0) {
            mask = type(uint256).max; // Set to maximum value of uint256
        } else {
            mask = 256 ** (32 - len) - 1;
        }

        assembly {
            let srcpart := and(mload(src), not(mask))
            let destpart := and(mload(dest), mask)
            mstore(dest, or(destpart, srcpart))
        }
    }
}
