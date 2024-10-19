// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

/**
 * @title Predeploys
 * @notice Predeploy addresses (match story/genutil/evm/predeploys.go)
 */
library Predeploys {
    address internal constant Namespace = 0xCCcCCc0000000000000000000000000000000000;
    uint256 internal constant NamespaceSize = 1024;

    /// @notice Predeploys
    address internal constant WIP = 0x1513000000000000000000000000000000000000;
    address internal constant Staking = 0xCCcCcC0000000000000000000000000000000001;
    address internal constant UBIPool = 0xCccCCC0000000000000000000000000000000002;
    address internal constant Upgrades = 0xccCCcc0000000000000000000000000000000003;
    address internal constant Create3 = 0xCcCcCC0000000000000000000000000000000004;

    address internal constant Timelock = address(uint160(Namespace) - 1);

    /// @notice Return true if `addr` is not proxied
    function notProxied(address addr) internal pure returns (bool) {
        return addr == WIP;
    }

    /// @notice Return implementation address for a proxied predeploy
    function getImplAddress(address addr) internal pure returns (address) {
        require(isPredeploy(addr), "Predeploys: not a predeploy");
        require(!notProxied(addr), "Predeploys: not proxied");

        // max uint160 is odd, which gives us unique implementation for each predeploy
        return address(type(uint160).max - uint160(addr));
    }

    /// @notice Return true if `addr` is an active predeploy
    function isActivePredeploy(address addr) internal pure returns (bool) {
        return addr == WIP || addr == Staking || addr == UBIPool || addr == Upgrades;
    }

    /// @notice Return true if `addr` is in some predeploy namespace
    function isPredeploy(address addr) internal pure returns (bool) {
        return (uint160(addr) >> 11 == uint160(Namespace) >> 11);
    }
}
