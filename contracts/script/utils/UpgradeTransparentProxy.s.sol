// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import { TimelockOperations } from "script/utils/TimelockOperations.s.sol";
import { console2 } from "forge-std/console2.sol";
import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { EIP1967Helper } from "../utils/EIP1967Helper.sol";

/// @title UpgradeTransparentProxy
/// @notice Helper script that generates a json file with the timelocked operations to upgrade a TransparentUpgradeableProxy
abstract contract UpgradeTransparentProxy is TimelockOperations {
    // The address of the sender
    address public from;
    // The operation salt for deterministic operation ID
    bytes32 public salt;

    constructor(string memory message, address _from, bytes32 _salt) TimelockOperations(message) {
        from = _from;
        console2.log("From:", from);
        salt = _salt;
        console2.logBytes32(salt);
    }

    /// @notice Get the addresses of the proxy, new implementation, and proxy admin
    /// @dev Should be overridden by the child contract
    /// @return proxyAddresses The addresses of the proxies
    /// @return newImplementationAddresses The addresses of the new implementations
    /// @return proxyAdminAddresses The addresses of the proxy admins
    function _getAddresses() internal view virtual returns (address[] memory, address[] memory, address[] memory);

    /// @notice Provide the target timelock to TimelockOperations
    /// @return currentTimelock
    function _getTargetTimelock() internal view virtual override returns (address) {
        return USE_CURRENT_TIMELOCK;
    }

    /// @notice Generate the timelocked operations
    function _generate() internal virtual override {
        (
            address[] memory proxyAddresses,
            address[] memory newImplementationAddresses,
            address[] memory proxyAdminAddresses
        ) = _getAddresses();
        require(proxyAddresses.length == newImplementationAddresses.length, "Proxy and new implementation addresses must be the same length");
        require(proxyAddresses.length > 0, "At least one proxy address must be provided");
        require(proxyAddresses.length == newImplementationAddresses.length, "Proxy and new implementation addresses must be the same length");

        // Get the proxy admin address from the proxy
        proxyAdminAddresses = new address[](proxyAddresses.length);
        for (uint256 i = 0; i < proxyAddresses.length; i++) {
            console2.log("Proxy address:", proxyAddresses[i]);
            console2.log("New implementation:", newImplementationAddresses[i]);
            proxyAdminAddresses[i] = EIP1967Helper.getAdmin(proxyAddresses[i]);
        }
        require(proxyAddresses.length == proxyAdminAddresses.length, "Proxy and proxy admin addresses must be the same length");

        // If there is only one proxy, we use the method for single operation
        if (proxyAddresses.length == 1) {
            _generateAction(
                from,
                proxyAdminAddresses[0],
                uint256(0),
                abi.encodeWithSignature(
                    "upgradeAndCall(address,address,bytes)",
                    proxyAddresses[0],
                    newImplementationAddresses[0],
                    ""
                ),
                bytes32(0),
                salt,
                minDelay
            );
        } else {
            // If there is more than one proxy, we use the method for batch operation
            uint256[] memory values = new uint256[](proxyAddresses.length);
            bytes[] memory data = new bytes[](proxyAddresses.length);
            for (uint256 i = 0; i < proxyAddresses.length; i++) {
                data[i] = abi.encodeWithSignature(
                    "upgradeAndCall(address,address,bytes)",
                    proxyAddresses[i],
                    newImplementationAddresses[i],
                    ""
                );
            }
            _generateBatchAction(
                from,
                proxyAdminAddresses,
                values,
                data,
                bytes32(0),
                salt,
                minDelay
            );
        }
    }
}
