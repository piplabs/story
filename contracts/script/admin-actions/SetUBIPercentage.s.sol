// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

import { TimelockOperations } from "script/utils/TimelockOperations.s.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { ChainIds } from "script/utils/ChainIds.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { IUBIPool } from "src/interfaces/IUBIPool.sol";
import { console2 } from "forge-std/console2.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

/// @notice Helper script that generates a json file with the timelocked operation to set the UBI percentage
/// @dev Set in the constructor Modes.SCHEDULE to run _scheduleActions, Modes.EXECUTE to run _executeActions
/// or Modes.CANCEL to run _cancelActions
contract SetUBIPercentage is TimelockOperations {
    address from;
    uint32 public percentage = 500; // 5.00%

    constructor() TimelockOperations("set-ubi-percentage-5%") {
        from = vm.envAddress("ADMIN_ADDRESS");
        console2.log("from", from);
    }

    function _getTargetTimelock() internal view virtual override returns (address) {
        return USE_CURRENT_TIMELOCK;
    }

    function _generate() internal virtual override {
        _generateAction(
            from,
            Predeploys.UBIPool,
            uint256(0),
            abi.encodeWithSelector(IUBIPool.setUBIPercentage.selector, percentage),
            bytes32(0),
            bytes32(0),
            minDelay
        );
    }
}
