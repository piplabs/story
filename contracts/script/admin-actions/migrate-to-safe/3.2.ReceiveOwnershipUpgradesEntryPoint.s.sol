/* solhint-disable no-console */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

import { console2 } from "forge-std/console2.sol";
/* solhint-disable max-line-length */

import { TimelockOperations } from "script/utils/TimelockOperations.s.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";

import { EIP1967Helper } from "script/utils/EIP1967Helper.sol";
import { Create3 } from "src/deploy/Create3.sol";

/// @title ReceiveOwnershipUpgradesEntryPoint
/// @notice Generates json files with the timelocked operations to receive the ownership of the contracts from the old timelock
contract ReceiveOwnershipUpgradesEntryPoint is TimelockOperations {

    TimelockController public newTimelock;

    address public from;

    constructor() TimelockOperations("safe-migr-receive-ownerships-upgrades-entrypoint") {
        from = vm.envAddress("SAFE_TIMELOCK_PROPOSER");
    }

    /// @dev target timelock is the newer timelock
    function _getTargetTimelock() internal virtual override returns (address) {
        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER_SAFE");
        address newTimelockAddress = Create3(Predeploys.Create3).predictDeterministicAddress(salt);
        newTimelock = TimelockController(payable(newTimelockAddress));
        require(address(newTimelock) != address(0), "Timelock not deployed");
        require(address(newTimelock) != address(currentTimelock()), "Timelock already set");
        return address(newTimelock);
    }

    function _generate() internal virtual override {
        bytes4 selector = Ownable2StepUpgradeable.acceptOwnership.selector;
        _generateAction(from, address(Predeploys.Upgrades), 0, abi.encodeWithSelector(selector), bytes32(0), bytes32(0), minDelay);
    }

}