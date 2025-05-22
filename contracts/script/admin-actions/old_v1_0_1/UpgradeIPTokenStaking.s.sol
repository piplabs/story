// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

import { UpgradeTransparentProxy } from "../../utils/UpgradeTransparentProxy.s.sol";

/// @notice Script to upgrade the IPTokenStaking contract through a timelock
contract UpgradeIpTokenStaking is UpgradeTransparentProxy {
    constructor()
        UpgradeTransparentProxy(
            "upgrade-staking-v1_0_1", // file name
            vm.envAddress("OLD_TIMELOCK_PROPOSER"),
            vm.envAddress("OLD_TIMELOCK_EXECUTOR"),
            vm.envAddress("OLD_TIMELOCK_CANCELLER"),
            bytes32(0) // salt
        )
    {}

    function _getAddresses()
        internal
        view
        override
        returns (
            address[] memory proxyAddresses,
            address[] memory newImplementationAddresses,
            address[] memory proxyAdminAddresses
        )
    {
        proxyAddresses = new address[](1);
        proxyAddresses[0] = vm.envAddress("PROXY_ADDRESS");
        newImplementationAddresses = new address[](1);
        newImplementationAddresses[0] = vm.envAddress("NEW_IMPLEMENTATION_ADDRESS");
        proxyAdminAddresses = new address[](1);
        proxyAdminAddresses[0] = vm.envAddress("PROXY_ADMIN_ADDRESS");
        return (proxyAddresses, newImplementationAddresses, proxyAdminAddresses);
    }
}
