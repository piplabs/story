/* solhint-disable no-console */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

import { console2 } from "forge-std/console2.sol";
/* solhint-disable max-line-length */

import { TimelockOperations } from "script/utils/TimelockOperations.s.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";
import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { EIP1967Helper } from "script/utils/EIP1967Helper.sol";
import { Create3 } from "src/deploy/Create3.sol";

/// @title BaseTransferOwnershipProxyAdmin
/// @notice Base contract to generate json files with the timelocked operations to transfer the ownership of half of the proxy admins to the new timelock.abi
/// We start with the last half of the proxy admins and move backwards, to test the migration in case of failure.
abstract contract BaseTransferOwnershipProxyAdmin is TimelockOperations {

    TimelockController public newTimelock;

    address[] public from;

    uint160 public fromIndex;
    uint160 public toIndex;

    string public message;

    constructor(string memory _message, uint160 _fromIndex, uint160 _toIndex) TimelockOperations(_message) {
        message = _message;
        require(_fromIndex < _toIndex, "fromIndex must be less than toIndex");
        require(_toIndex <= Predeploys.NamespaceSize, "toIndex must be less than or equal to Predeploys.NamespaceSize");
        fromIndex = _fromIndex;
        toIndex = _toIndex;
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
        uint256 targetsLength = toIndex - fromIndex + 1;
        console2.log("targetsLength", targetsLength);

        address[] memory targets = new address[](targetsLength);
        for (uint160 i = fromIndex; i <= toIndex; i++) {
            console2.log("i", i);
            // Get proxy admins for each predeploy with the EIP1967 helper
            address predeploy = address(uint160(Predeploys.Namespace) + i);
            console2.log("predeploy", predeploy);
            address proxyAdmin = EIP1967Helper.getAdmin(predeploy);
            console2.log("proxyAdmin", proxyAdmin);
            console2.log("fuck");
            uint160 targetIndex = i-fromIndex;
            console2.log(targetIndex);
            targets[targetIndex] = proxyAdmin;
        }

        bytes4 selector = Ownable.transferOwnership.selector;
        bytes[] memory data = new bytes[](targetsLength);
        for (uint160 i = 0; i < targetsLength; i++) {
            data[i] = abi.encodeWithSelector(selector, address(newTimelock));
        }
        uint256[] memory values = new uint256[](targetsLength);

        _generateBatchAction(from, targets, values, data, bytes32(0), bytes32(0), minDelay);
    }

}