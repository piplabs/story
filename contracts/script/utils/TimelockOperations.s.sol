/* solhint-disable no-console */
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import { console2 } from "forge-std/console2.sol";
import { Script } from "forge-std/Script.sol";

import { Predeploys } from "src/libraries/Predeploys.sol";
import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

// script
import { JSONTxWriter } from "./JSONTxWriter.s.sol";
import { StringUtil } from "./StringUtil.sol";

/**
 * @title TimelockOperations
 * @notice Script to generate tx json to schedule, execute, and cancel upgrades for the protocol through the TimelockController
 * via multisig.
 */
abstract contract TimelockOperations is Script, JSONTxWriter {

    /// @notice timelock controller
    TimelockController public timelock;

    /// @notice min delay for the timelock operation
    uint256 public minDelay;

    /// @notice constant to use the current timelock (owner of IPTokenStaking)
    address constant USE_CURRENT_TIMELOCK = 0x1111111111111111111111111111111111111111;

    modifier verifyFrom(address[] memory from) {
        require(from.length == 3, "TimelockOperations: from must be an array of 3 addresses");
        require(from[0] != address(0), "TimelockOperations: scheduler must be a valid address");
        require(from[1] != address(0), "TimelockOperations: executor must be a valid address");
        require(from[2] != address(0), "TimelockOperations: canceler must be a valid address");
        _;
    }

    /// @notice constructor
    /// @param _action the action name, will be used in the name of the tx json file
    constructor(string memory _action) JSONTxWriter(_action) {

    }

    /// @notice get the current timelock
    /// @return the address of the current timelock (owner of IPTokenStaking)
    function currentTimelock() public view returns (address) {
        return Ownable2StepUpgradeable(Predeploys.Staking).owner();
    }

    function run() public virtual {
        address _timelockAddress = _getTargetTimelock();
        if (_timelockAddress == USE_CURRENT_TIMELOCK) {
            timelock = TimelockController(payable(currentTimelock()));
        } else {
            if (_timelockAddress.code.length == 0) {
                revert("TimelockOperations: timelock address must be a valid contract");
            }
            timelock = TimelockController(payable(_timelockAddress));
        }

        minDelay = timelock.getMinDelay();
        console2.log("Timelock deployed at", address(timelock));
        console2.log("Generating actions...");
        _generate();
        console2.log("Writing tx json...");
        _writeFiles();
    }

    /// @notice provide the target timelock
    /// @return the address of the timelock that will be used for the operation
    function _getTargetTimelock() internal virtual returns (address);

    /// @notice generate the tx json. Should be overridden by the child contract
    /// @dev Child contract should call _generateAction for each action to be scheduled, executed, and cancelled
    function _generate() internal virtual;

    /// @notice generate the 3 JSON (schedule, execute, cancel) for a single target
    /// @param from The addresses of the sender for the schedule, execute, and cancel
    /// @param target The address of the contract to call
    /// @param value The value to send with the call
    /// @param data The encoded target method call
    /// @param predecessor The hash of the predecessor operation (optional)
    /// @param salt The salt for the timelock operation (optional, needed for calls with repeated `data`)
    /// @param delay The delay for the timelock operation. 0 is minimum delay
    function _generateAction(
        address[] memory from,
        address target,
        uint256 value,
        bytes memory data,
        bytes32 predecessor,
        bytes32 salt,
        uint256 delay
    ) internal verifyFrom(from) {
        _scheduleAction(from[0], target, value, data, predecessor, salt, delay);
        _executeAction(from[1], target, value, data, predecessor, salt);
        _cancelAction(from[2], target, value, data, predecessor, salt);
    }

    /// @notice Encodes the call to TimelockController.schedule
    /// @param from The address of the sender
    /// @param target The address of the contract to call
    /// @param value The value to send with the call
    /// @param data The encoded target method call
    /// @param predecessor The hash of the predecessor operation (optional)
    /// @param salt The salt for the timelock operation (optional, needed for calls with repeated `data`)
    /// @param delay The delay for the timelock operation. Must be >= minDelay
    function _scheduleAction(
        address from,
        address target,
        uint256 value,
        bytes memory data,
        bytes32 predecessor,
        bytes32 salt,
        uint256 delay
    ) internal {
        bytes memory _txData = abi.encodeWithSelector(
            TimelockController.schedule.selector,
            target,
            value,
            data,
            predecessor,
            salt,
            delay
        );
        _saveTx(TimelockOp.SCHEDULE, from, address(timelock), value, _txData, string.concat(action, "-schedule"));
    }

    /// @notice Encodes the call to TimelockController.execute
    /// @param target The address of the contract to call
    /// @param value The value to send with the call
    /// @param data The encoded target method call
    /// @param predecessor The hash of the predecessor operation (optional)
    /// @param salt The salt for the timelock operation (optional, needed for calls with repeated `data`)
    function _executeAction(
        address from,
        address target,
        uint256 value,
        bytes memory data,
        bytes32 predecessor,
        bytes32 salt
    ) internal {
        bytes memory _txData = abi.encodeWithSelector(
            TimelockController.execute.selector,
            target,
            value,
            data,
            predecessor,
            salt
        );
        _saveTx(TimelockOp.EXECUTE, from, address(timelock), value, _txData, string.concat(action, "-execute"));
    }

    /// @notice Encodes the call to TimelockController.cancel
    /// @param target The address of the contract to call
    /// @param value The value sent to scheduled call
    /// @param data The encoded target method call
    /// @param predecessor The hash of the predecessor operation (optional)
    /// @param salt The salt for the timelock operation (optional, needed for calls with repeated `data`)
    function _cancelAction(
        address from,
        address target,
        uint256 value,
        bytes memory data,
        bytes32 predecessor,
        bytes32 salt
    ) internal {
        bytes memory _txData = abi.encodeWithSelector(
            TimelockController.cancel.selector,
            timelock.hashOperation(target, value, data, predecessor, salt)
        );
        _saveTx(TimelockOp.CANCEL, from, address(timelock), value, _txData, string.concat(action, "-cancel"));
    }

    /// @notice Encodes calls to TimelockController.scheduleBatch, TimelockController.executeBatch, and TimelockController.cancel
    /// @param from The addresses of the sender for the schedule, execute, and cancel
    /// @param targets The addresses of the contracts to call
    /// @param values The values to send with the calls
    /// @param data The encoded target method calls
    /// @param predecessor The hash of the predecessor operation (optional)
    /// @param salt The salt for the timelock operation (optional, needed for calls with repeated `data`)
    /// @param delay The delay for the timelock operation. 0 is minimum delay
    function _generateBatchAction(
        address[] memory from,
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory data,
        bytes32 predecessor,
        bytes32 salt,
        uint256 delay
    ) internal verifyFrom(from) {
        _scheduleBatchAction(from[0], targets, values, data, predecessor, salt, delay);
        _executeBatchAction(from[1], targets, values, data, predecessor, salt);
        _cancelBatchAction(from[2], targets, values, data, predecessor, salt);
    }

    /// @notice Encodes the call to TimelockController.scheduleBatch
    /// @param from The address of the sender
    /// @param targets The addresses of the contracts to call
    /// @param values The values to send with the calls
    /// @param data The encoded target method calls
    /// @param predecessor The hash of the predecessor operation (optional)
    /// @param salt The salt for the timelock operation (optional, needed for calls with repeated `data`)
    /// @param delay The delay for the timelock operation. 0 is minimum delay
    function _scheduleBatchAction(
        address from,
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory data,
        bytes32 predecessor,
        bytes32 salt,
        uint256 delay
    ) internal {
        bytes memory _txData = abi.encodeWithSelector(
            TimelockController.scheduleBatch.selector,
            targets,
            values,
            data,
            predecessor,
            salt,
            delay
        );
        _saveTx(
            TimelockOp.SCHEDULE,
            from,
            address(timelock),
            _sumValues(values),
            _txData,
            string.concat(action, "-schedule")
        );
    }

    /// @notice Encodes the call to TimelockController.executeBatch
    /// @param from The address of the sender
    /// @param targets The addresses of the contracts to call
    /// @param values The values to send with the calls
    /// @param data The encoded target method calls
    /// @param predecessor The hash of the predecessor operation (optional)
    /// @param salt The salt for the timelock operation (optional, needed for calls with repeated `data`)
    function _executeBatchAction(
        address from,
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory data,
        bytes32 predecessor,
        bytes32 salt
    ) internal {
        bytes memory _txData = abi.encodeWithSelector(
            TimelockController.executeBatch.selector,
            targets,
            values,
            data,
            predecessor,
            salt
        );
        _saveTx(
            TimelockOp.EXECUTE,
            from,
            address(timelock),
            _sumValues(values),
            _txData,
            string.concat(action, "-execute")
        );
    }

    /// @notice Encodes the call to TimelockController.cancelBatch
    /// @param from The address of the sender
    /// @param targets The addresses of the contracts to call
    /// @param values The values to send with the calls
    /// @param data The encoded target method calls
    /// @param predecessor The hash of the predecessor operation (optional)
    /// @param salt The salt for the timelock operation (optional, needed for calls with repeated `data`)
    function _cancelBatchAction(
        address from,
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory data,
        bytes32 predecessor,
        bytes32 salt
    ) internal {
        bytes memory _txData = abi.encodeWithSelector(
            TimelockController.cancel.selector,
            timelock.hashOperationBatch(targets, values, data, predecessor, salt)
        );
        _saveTx(
            TimelockOp.CANCEL,
            from,
            address(timelock),
            _sumValues(values),
            _txData,
            string.concat(action, "-cancel")
        );
    }

    /// @notice Sums the values of an array to get the total value of the batch
    /// @param values The values to sum
    /// @return the sum of the values
    function _sumValues(uint256[] memory values) internal pure returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < values.length; i++) {
            sum += values[i];
        }
        return sum;
    }
}
