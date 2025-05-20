/* solhint-disable no-console */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

/* solhint-disable max-line-length */

import { TimelockOperations } from "script/utils/TimelockOperations.s.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { Create3 } from "src/deploy/Create3.sol";

/// @title TransferOwnershipsUpgradesEntrypoint
/// @notice Generates json files with the timelocked operations to transfer the ownership of the other half of the proxy admins to the new timelock
/// We start with UpgradesEntrypoint and move backwards, to test the migration in case of failure.
/// This contract is Ownable2StepUpgradeable, so we need to accept ownership transfer from the new timelock in the next step.
contract TransferOwnershipsUpgradesEntrypoint is TimelockOperations {
    TimelockController public newTimelock;

    address[] public from;

    constructor() TimelockOperations("safe-migr-transfer-ownerships-upgrades-entrypoint") {
        from = new address[](3);
        from[0] = vm.envAddress("OLD_TIMELOCK_PROPOSER");
        from[1] = vm.envAddress("OLD_TIMELOCK_EXECUTOR");
        from[2] = vm.envAddress("OLD_TIMELOCK_GUARDIAN");
        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER_SAFE");
        address newTimelockAddress = Create3(Predeploys.Create3).predictDeterministicAddress(salt);
        newTimelock = TimelockController(payable(newTimelockAddress));
    }

    /// @dev the old timelock will execute the operations
    function _getTargetTimelock() internal virtual override returns (address) {
        return vm.envAddress("OLD_TIMELOCK_ADDRESS");
    }

    function _generate() internal virtual override {
        require(address(newTimelock) != address(0), "Timelock not deployed");
        require(address(newTimelock) != address(currentTimelock()), "Timelock already set");

        bytes4 selector = Ownable2StepUpgradeable.transferOwnership.selector;
        _generateAction(
            from,
            address(Predeploys.Upgrades),
            0,
            abi.encodeWithSelector(selector, address(newTimelock)),
            bytes32(0),
            bytes32(0),
            minDelay
        );
    }
}
