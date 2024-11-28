/* solhint-disable no-console */
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import { console2 } from "forge-std/console2.sol";
import { Script } from "forge-std/Script.sol";

import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

// script
import { JSONBatchTxHelper } from "./JSONBatchTxHelper.s.sol";
import { StringUtil } from "./StringUtil.sol";

/**
 * @title TimelockedOperations
 * @notice Script to schedule, execute, or cancel upgrades for the protocol contracts
 */
abstract contract TimelockedOperations is Script, JSONBatchTxHelper {

    /// @notice Upgrade modes
    enum UpgradeModes {
        UNSET, // Unset mode
        SCHEDULE, // Schedule upgrades in AccessManager
        EXECUTE, // Execute scheduled upgrades
        CANCEL // Cancel scheduled upgrades
    }
    /// @notice End result of the script
    enum Output {
        UNSET, // Unset output
        TX_EXECUTION, // One action per operation
        BATCH_TX_EXECUTION, // Multiple actions in one tx
        BATCH_TX_JSON // Prepare raw bytes for multisig. Multisig may batch txs (e.g. Gnosis Safe JSON input in tx builder)
    }


    ///////// EDITABLE INPUT /////////
    UpgradeModes internal mode;
    Output internal outputType;

    /// @notice action acumulator for batch txs
    string action;
    TimelockController timelock;
    uint256 minDelay;
    /// @notice salt for the timelocked call
    bytes32 salt;
    /// @notice predecessor for the timelocked call
    bytes32 predecessor;

    constructor(string memory _action, UpgradeModes _mode, Output _outputType, address _timelock) JSONBatchTxHelper() {
        if (_mode == UpgradeModes.UNSET) {
            revert("Mode must be set");
        }
        mode = _mode;
        if (_outputType == Output.UNSET) {
            revert("Output type must be set");
        }
        outputType = _outputType;
        action = _action;
        address timelockAddress;
        if (_timelock == address(0)) {
            timelockAddress = Ownable2StepUpgradeable(Predeploys.Staking).owner();
        } else {
            timelockAddress = _timelock;
        }
        timelock = TimelockController(payable(timelockAddress));
    }

    function run() public virtual {        
        minDelay = timelock.getMinDelay();
        console2.log("Min delay: ", minDelay);
        _startOperation();
        // Read upgrade proposals file
        if (outputType == Output.BATCH_TX_JSON) {
            console2.log("Generating tx json...");
        }
       
        // If output is JSON, write actions to file
        if (outputType == Output.BATCH_TX_JSON) {    
            _writeBatchTxsOutput(string.concat(action, "-", _modeDescription())); // JsonBatchTxHelper.s.sol
        } else if (outputType == Output.BATCH_TX_EXECUTION) {
            _executeBatchedOperations();
        }
            
        // If output is TX_EXECUTION or BATCH_TX_EXECUTION, no further action is needed
    }

    function _startOperation() private {
        // Decide actions based on mode
        if (mode == UpgradeModes.SCHEDULE) {
            _scheduleActions();
        } else if (mode == UpgradeModes.EXECUTE) {
            _executeActions();
        } else if (mode == UpgradeModes.CANCEL) {
            _cancelActions();
        } else {
            revert("Invalid mode");
        }
    }

    function _scheduleActions() internal virtual;
    function _executeActions() internal virtual;
    function _cancelActions() internal virtual;

    function _scheduleAction(address target, uint256 value, bytes memory data, bytes32 _predecessor, bytes32 _salt) internal {
        if (outputType == Output.TX_EXECUTION) {
            timelock.schedule(target, value, data, _predecessor, _salt, minDelay);
        } else {
            bytes memory scheduleData = abi.encodeCall(TimelockController.schedule, (target, value, data, _predecessor, _salt, minDelay));
            _saveTx(address(timelock), 0, scheduleData, string.concat(action, "-", _modeDescription()));
        }
    }

    function _executeAction(address target, uint256 value, bytes memory data, bytes32 _predecessor, bytes32 _salt) internal {
        if (outputType == Output.TX_EXECUTION) {
            timelock.execute(target, value, data, _predecessor, _salt);
        } else {
            bytes memory executeData = abi.encodeCall(TimelockController.execute, (target, value, data, _predecessor, _salt));
            _saveTx(address(timelock), 0, executeData, string.concat(action, "-", _modeDescription()));
        }
    }

    function _cancelAction(address target, uint256 value, bytes memory data, bytes32 _predecessor, bytes32 _salt) internal {
        if (outputType == Output.TX_EXECUTION) {
            timelock.cancel(timelock.hashOperation(target, value, data, _predecessor, _salt));
        } else {
            bytes memory cancelData = abi.encodeCall(TimelockController.cancel, (timelock.hashOperation(target, value, data, _predecessor, _salt)));
            _saveTx(address(timelock), 0, cancelData, string.concat(action, "-", _modeDescription()));
        }
    }

    function _executeBatchedOperations() private {
        address[] memory targets = new address[](transactions.length);
        uint256[] memory values = new uint256[](transactions.length);
        bytes[] memory payloads = new bytes[](transactions.length);
        for (uint256 i = 0; i < transactions.length; i++) {
            targets[i] = transactions[i].to;
            values[i] = transactions[i].value;
            payloads[i] = transactions[i].data;
        }
        // Decide actions based on mode
        if (mode == UpgradeModes.SCHEDULE) {
            timelock.scheduleBatch(targets, values, payloads, predecessor, salt, minDelay);
        } else if (mode == UpgradeModes.EXECUTE) {
            timelock.executeBatch(targets, values, payloads, predecessor, salt);
        } else if (mode == UpgradeModes.CANCEL) {
            timelock.cancel(timelock.hashOperationBatch(targets, values, payloads, predecessor, salt));
        } else {
            revert("Invalid mode");
        }
    }

    function _modeDescription() internal view returns (string memory) {
        if (mode == UpgradeModes.SCHEDULE) {
            return "schedule";
        } else if (mode == UpgradeModes.EXECUTE) {
            return "execute";
        } else if (mode == UpgradeModes.CANCEL) {
            return "cancel";
        } else {
            revert("Invalid mode");
        }
    }
}
