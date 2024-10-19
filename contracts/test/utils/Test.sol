// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Test as ForgeTest } from "forge-std/Test.sol";

import { IPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { UpgradeEntrypoint } from "../../src/protocol/UpgradeEntrypoint.sol";
import { UBIPool } from "../../src/protocol/UBIPool.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";
import { Create3 } from "../../src/deploy/Create3.sol";
import { GenerateAlloc } from "../../script/GenerateAlloc.s.sol";
import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";

import { console2 } from "forge-std/console2.sol";

contract Test is ForgeTest {
    address internal admin = address(0x123);
    address internal executor = address(0x456);
    address internal guardian = address(0x789);

    IPTokenStaking internal ipTokenStaking;
    UpgradeEntrypoint internal upgradeEntrypoint;
    UBIPool internal ubiPool;
    Create3 internal create3;
    TimelockController internal timelock;
    function setUp() public virtual {
        GenerateAlloc initializer = new GenerateAlloc();
        initializer.disableStateDump(); // Faster tests. Don't call to verify JSON output
        initializer.setAdminAddresses(admin, executor, guardian);
        initializer.run();
        ipTokenStaking = IPTokenStaking(Predeploys.Staking);
        upgradeEntrypoint = UpgradeEntrypoint(Predeploys.Upgrades);
        ubiPool = UBIPool(Predeploys.UBIPool);
        create3 = Create3(Predeploys.Create3);
        timelock = TimelockController(payable(Predeploys.Timelock));
    }

    function performTimelocked(address target, bytes memory data) internal {
        vm.prank(admin);
        console2.log("target");
        console2.log(target);
        console2.log(address(timelock));
        timelock.schedule(
            target,
            0, // value
            data,
            bytes32(0), // predecessor: Non Zero if order must be respected
            bytes32(keccak256("SALT")), // salt
            0 // delay: 0 to use minimum delay
        );
        vm.warp(block.timestamp + timelock.getMinDelay() + 1);
        vm.prank(executor);
        timelock.execute(
            target,
            0, // value
            data,
            bytes32(0), // predecessor: Non Zero if order must be respected
            bytes32(keccak256("SALT")) // salt
        );
    }

    function expectRevertTimelocked(address target, bytes memory data, string memory reason) internal {
        vm.prank(admin);
        timelock.schedule(
            target,
            0, // value
            data,
            bytes32(0), // predecessor: Non Zero if order must be respected
            bytes32(keccak256("SALT")), // salt
            0 // delay: 0 to use minimum delay
        );
        vm.warp(block.timestamp + timelock.getMinDelay() + 1);
        vm.expectRevert(bytes(reason));
        timelock.execute(
            target,
            0, // value
            data,
            bytes32(0), // predecessor: Non Zero if order must be respected
            bytes32(keccak256("SALT")) // salt
        );
    }
}
