/* solhint-disable no-console */
/* solhint-disable max-line-length */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

import { TimelockOperations } from "script/utils/TimelockOperations.s.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { AccessControl } from "@openzeppelin/contracts/access/AccessControl.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { Create3 } from "src/deploy/Create3.sol";

/// @title RenounceGovernanceRoles
/// @notice Generates json files with the timelocked operations to renounce roles from the old timelock
contract RenounceGovernanceRoles is TimelockOperations {
    TimelockController public newTimelock;

    address[] public from;

    constructor() TimelockOperations("safe-migr-renounce-gov-roles") {
        from = new address[](3);
        from[0] = vm.envAddress("OLD_TIMELOCK_PROPOSER");
        from[1] = vm.envAddress("OLD_TIMELOCK_EXECUTOR");
        from[2] = vm.envAddress("OLD_TIMELOCK_CANCELLER");
    }

    /// @dev target timelock is the newer timelock
    function _getTargetTimelock() internal virtual override returns (address) {
        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER_SAFE");
        address newTimelockAddress = Create3(Predeploys.Create3).predictDeterministicAddress(salt);
        newTimelock = TimelockController(payable(newTimelockAddress));
        require(address(newTimelock) != address(0), "Timelock not deployed");
        return address(newTimelock);
    }

    function _generate() internal virtual override {
        require(address(newTimelock) != address(0), "Timelock not deployed");
        require(address(newTimelock) == address(currentTimelock()), "Renouncing before timelock migration");
        // Just in case
        for (uint256 i = 0; i < from.length; i++) {
            require(from[i] != vm.envAddress("SAFE_TIMELOCK_PROPOSER"), "From address is Safe proposer");
            require(from[i] != vm.envAddress("SAFE_TIMELOCK_EXECUTOR"), "From address is Safe executor");
            require(from[i] != vm.envAddress("SAFE_TIMELOCK_CANCELLER"), "From address is Safe guardian");
        }

        bytes4 selector = AccessControl.renounceRole.selector;

        // Renounce proposer wallet
        _saveTx(
            Operation.REGULAR_TX,
            from[0],
            address(newTimelock),
            0,
            abi.encodeWithSelector(selector, newTimelock.PROPOSER_ROLE(), from[0]),
            string.concat("Renounce proposer role old multisig")
        );
        // Canceller role is set in TimelockController constructor
        _saveTx(
            Operation.REGULAR_TX,
            from[0],
            address(newTimelock),
            0,
            abi.encodeWithSelector(selector, newTimelock.CANCELLER_ROLE(), from[0]),
            string.concat("Renounce canceller role old multisig")
        );
        _saveTx(
            Operation.REGULAR_TX,
            from[0],
            address(newTimelock),
            0,
            abi.encodeWithSelector(selector, newTimelock.DEFAULT_ADMIN_ROLE(), from[0]),
            string.concat("Renounce default admin role old multisig")
        );

        // Renounce executor role
        _saveTx(
            Operation.REGULAR_TX,
            from[1],
            address(newTimelock),
            0,
            abi.encodeWithSelector(selector, newTimelock.EXECUTOR_ROLE(), from[1]),
            string.concat("Renounce executor role old multisig")
        );
        // Renounce canceller role
        _saveTx(
            Operation.REGULAR_TX,
            from[2],
            address(newTimelock),
            0,
            abi.encodeWithSelector(selector, newTimelock.CANCELLER_ROLE(), from[2]),
            string.concat("Renounce canceller role old guardian multisig")
        );
    }
}
