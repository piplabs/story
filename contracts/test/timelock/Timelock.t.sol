// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

import { Test } from "../utils/Test.sol";
import { IPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";

contract TimelockTest is Test {
    function setUp() public override {
        super.setUp();
    }

    function testCancelBeforeExecute() public {
        // Prepare a sample operation
        address target = address(ipTokenStaking);
        uint256 value = 0;
        bytes memory data = abi.encodeWithSelector(IPTokenStaking.setFee.selector, 2 ether);
        bytes32 predecessor = bytes32(0);
        bytes32 salt = keccak256("TEST_SALT");
        uint256 delay = timelock.getMinDelay();

        // Schedule the operation
        vm.prank(admin);
        timelock.schedule(target, value, data, predecessor, salt, delay);

        // Ensure the operation is pending
        bytes32 operationId = timelock.hashOperation(target, value, data, predecessor, salt);
        assertTrue(timelock.isOperationPending(operationId));

        // Cancel the operation
        vm.prank(admin);
        timelock.cancel(operationId);

        // Ensure the operation is no longer pending
        assertFalse(timelock.isOperationPending(operationId));

        // Wait for the delay to pass
        vm.warp(block.timestamp + delay + 1);

        // Try to execute the cancelled operation
        vm.prank(executor);
        vm.expectRevert(
            abi.encodeWithSelector(
                TimelockController.TimelockUnexpectedOperationState.selector,
                operationId,
                bytes32(1 << uint8(TimelockController.OperationState.Ready))
            )
        );

        timelock.execute(target, value, data, predecessor, salt);

        // Verify that the fee wasn't changed
        assertEq(ipTokenStaking.fee(), 1 ether);
    }

    function testExecuteSequenceWithPredecessors() public {
        // Prepare sample operations
        address target = address(ipTokenStaking);
        uint256 value = 0;
        bytes memory data1 = abi.encodeWithSelector(IPTokenStaking.setFee.selector, 2 ether);
        bytes memory data2 = abi.encodeWithSelector(IPTokenStaking.setFee.selector, 3 ether);
        bytes memory data3 = abi.encodeWithSelector(IPTokenStaking.setFee.selector, 4 ether);
        bytes32 salt1 = keccak256("SALT_1");
        bytes32 salt2 = keccak256("SALT_2");
        bytes32 salt3 = keccak256("SALT_3");
        uint256 delay = timelock.getMinDelay();

        // Schedule the first operation
        vm.prank(admin);
        timelock.schedule(target, value, data1, bytes32(0), salt1, delay);
        bytes32 id1 = timelock.hashOperation(target, value, data1, bytes32(0), salt1);

        // Schedule the second operation with the first as predecessor
        vm.prank(admin);
        timelock.schedule(target, value, data2, id1, salt2, delay);
        bytes32 id2 = timelock.hashOperation(target, value, data2, id1, salt2);

        // Schedule the third operation with the second as predecessor
        vm.prank(admin);
        timelock.schedule(target, value, data3, id2, salt3, delay);

        // Wait for the delay to pass
        vm.warp(block.timestamp + delay + 1);

        // Execute the first operation
        vm.prank(executor);
        timelock.execute(target, value, data1, bytes32(0), salt1);
        assertEq(ipTokenStaking.fee(), 2 ether);

        // Try to execute the third operation (should fail due to unexecuted predecessor)
        vm.prank(executor);
        vm.expectRevert(abi.encodeWithSelector(TimelockController.TimelockUnexecutedPredecessor.selector, id2));
        timelock.execute(target, value, data3, id2, salt3);

        // Execute the second operation
        vm.prank(executor);
        timelock.execute(target, value, data2, id1, salt2);
        assertEq(ipTokenStaking.fee(), 3 ether);

        // Finally, execute the third operation
        vm.prank(executor);
        timelock.execute(target, value, data3, id2, salt3);
        assertEq(ipTokenStaking.fee(), 4 ether);
    }
}
