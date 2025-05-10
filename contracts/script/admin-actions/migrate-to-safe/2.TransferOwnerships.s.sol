/* solhint-disable no-console */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

import { console2 } from "forge-std/console2.sol";
/* solhint-disable max-line-length */

import { TimelockOperations } from "script/utils/TimelockOperations.s.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";

import { EIP1967Helper } from "script/utils/EIP1967Helper.sol";
import { Create3 } from "src/deploy/Create3.sol";

/// @title TransferOwnerships
/// @notice Generates json files with the timelocked operations to transfer the ownership of the contracts to the new timelock
contract TransferOwnerships is TimelockOperations {

    TimelockController public newTimelock;

    address public from;

    constructor() TimelockOperations("safe-migration-transfer-ownerships") {
        from = vm.envAddress("OLD_TIMELOCK_ADMIN_ADDRESS");
        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER_SAFE");
        address newTimelockAddress = Create3(Predeploys.Create3).predictDeterministicAddress(salt);
        newTimelock = TimelockController(payable(newTimelockAddress));
    }

    /// @dev the old timelock will execute the operations
    function _getTargetTimelock() internal virtual override returns (address) {
        return USE_CURRENT_TIMELOCK;
    }

    function _generate() internal virtual override {
        require(address(newTimelock) != address(0), "Timelock not deployed");
        require(address(newTimelock) != address(currentTimelock()), "Timelock already set");
        uint256 targetsLength = 3 + Predeploys.NamespaceSize;

        address[] memory targets = new address[](targetsLength);
        targets[0] = Predeploys.Staking;
        targets[1] = Predeploys.UBIPool;
        targets[2] = Predeploys.Upgrades;
        for (uint160 i = 1; i <= Predeploys.NamespaceSize; i++) {
            // Get proxy admins for each predeploy with the EIP1967 helper
            address proxyAdmin = EIP1967Helper.getAdmin(address(uint160(Predeploys.Namespace) + i));
            targets[2 + i] = proxyAdmin;
        }

        bytes4 selector = Ownable2StepUpgradeable.transferOwnership.selector;
        bytes[] memory data = new bytes[](targetsLength);
        for (uint160 i = 0; i < targetsLength; i++) {
            data[i] = abi.encodeWithSelector(selector, address(newTimelock));
        }
        uint256[] memory values = new uint256[](targetsLength);

        _generateBatchAction(from, targets, values, data, bytes32(0), bytes32(0), minDelay);
    }

}