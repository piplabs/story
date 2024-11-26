// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

/**
 * @title Predeploys
 * @notice Predeploy addresses (match story/genutil/evm/predeploys.go)
 */
library Predeploys {
    address internal constant Namespace = 0xCCcCCc0000000000000000000000000000000000;
    uint256 internal constant NamespaceSize = 1024;

    /// @notice Predeploys
    address internal constant WIP = 0x1516000000000000000000000000000000000000;
    address internal constant Staking = 0xCCcCcC0000000000000000000000000000000001;
    address internal constant UBIPool = 0xCccCCC0000000000000000000000000000000002;
    address internal constant Upgrades = 0xccCCcc0000000000000000000000000000000003;

    /// @notice Create3 factory address
    /// @dev We maximize compatibility with the contracts deployed by ZeframLou
    address internal constant Create3 = 0x9fBB3DF7C40Da2e5A0dE984fFE2CCB7C47cd0ABf;

    /// @notice ERC6551Registry address
    /// @dev The common address for the ERC6551Registry across all chains defined by ERC-6551
    address internal constant ERC6551Registry = 0x000000006551c19487814612e58FE06813775758;

    /// @notice Return true if `addr` is proxied
    function proxied(address addr) internal pure returns (bool) {
        return addr > Namespace && addr <= address(uint160(Namespace) + uint160(NamespaceSize));
    }

    /// @notice Return implementation address for a proxied predeploy
    function getImplAddress(address proxyAddress) internal pure returns (address) {
        require(proxied(proxyAddress), "Predeploys: not proxied");

        // max uint160 is odd, which gives us unique implementation for each predeploy
        return address(type(uint160).max - uint160(proxyAddress));
    }
}
