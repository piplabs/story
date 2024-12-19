// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

import { JSONTimelockedOperations } from "../utils/JSONTimelockedOperations.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { ChainIds } from "../utils/ChainIds.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { IPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";

/// @notice Helper script that generates a json file with the timelocked operation to set the IPToken staking fee
/// @dev Set in the constructor UpgradeModes.SCHEDULE to run _scheduleActions, UpgradeModes.EXECUTE to run _executeActions
/// or UpgradeModes.CANCEL to run _cancelActions
contract SetIPTokenStakingFee is JSONTimelockedOperations {
    constructor() JSONTimelockedOperations(
        "set-ip-token-staking-fee-10IP",
        UpgradeModes.SCHEDULE
    ) {}

    uint256 public fee = 10 ether;

    function _scheduleActions() internal virtual override {
        _scheduleAction(
            Predeploys.Staking,
            uint256(0),
            abi.encodeWithSelector(IPTokenStaking.setFee.selector, fee),
            bytes32(0),
            bytes32(0),
            minDelay
        );
    }

    function _executeActions() internal virtual override {
        // Public execution, can be called directly or here with OUTPUT_TYPE.TX_EXECUTION
        _executeAction(
            Predeploys.Staking,
            uint256(0),
            abi.encodeWithSelector(IPTokenStaking.setFee.selector, fee),
            bytes32(0),
            bytes32(0)
        );
    }

    function _cancelActions() internal virtual override {
        _cancelAction(
            Predeploys.Staking,
            uint256(0),
            abi.encodeWithSelector(IPTokenStaking.setFee.selector, fee),
            bytes32(0),
            bytes32(0)
        );
    }
}
