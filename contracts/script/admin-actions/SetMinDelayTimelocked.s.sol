// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import { TimelockedOperations } from "../utils/TimelockedOperations.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { ChainIds } from "../utils/ChainIds.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";

contract SetMinDelayTimelocked is TimelockedOperations {
    constructor() TimelockedOperations(
        "test-set-min-delay-timelock",
        UpgradeModes.CANCEL,
        Output.BATCH_TX_JSON,
        0x9DAE0C1E36b65F9e5ABd5998ce83a64706d206E1 // New timelock
    ) {}

    function _scheduleActions() internal virtual override {
        _scheduleAction(
            address(timelock),
            0,
            abi.encodeCall(TimelockController.updateDelay, (4 minutes)),
            bytes32(0),
            bytes32(keccak256("test-set-min-delay-timelock-schedule"))
        );
    }
    function _executeActions() internal virtual override {
        _executeAction(
            address(timelock),
            0,
            abi.encodeCall(TimelockController.updateDelay, (4 minutes)),
            bytes32(0),
            bytes32(keccak256("test-set-min-delay-timelock-execute"))
        );
    }

    function _cancelActions() internal virtual override {
        _cancelAction(
            address(timelock),
            0,
            abi.encodeCall(TimelockController.updateDelay, (4 minutes)),
            bytes32(0),
            bytes32(keccak256("test-set-min-delay-timelock-cancel"))
        );
    }
}
