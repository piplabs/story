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
import { ERC6551Registry } from "erc6551/ERC6551Registry.sol";

contract Test is ForgeTest {
    address internal admin = address(0x123);
    address internal executor = address(0x456);
    address internal guardian = address(0x789);

    address internal deployer = address(0xDDdDddDdDdddDDddDDddDDDDdDdDDdDDdDDDDDDd);

    IPTokenStaking internal ipTokenStaking;
    UpgradeEntrypoint internal upgradeEntrypoint;
    UBIPool internal ubiPool;
    Create3 internal create3;
    ERC6551Registry internal erc6551Registry;
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
        erc6551Registry = ERC6551Registry(Predeploys.ERC6551Registry);
        address timelockAddress = create3.getDeployed(deployer, keccak256("STORY_TIMELOCK_CONTROLLER"));
        timelock = TimelockController(payable(timelockAddress));
        require(timelockAddress.code.length > 0, "Timelock not deployed");
    }

    /// @notice schedules, waits for timelock and executes a timelocked call
    /// @param target The address to call
    /// @param data The data to call with
    function performTimelocked(address target, bytes memory data) internal {
        uint256 minDelay = timelock.getMinDelay();
        vm.prank(admin);
        timelock.schedule(
            target,
            0, // value
            data,
            bytes32(0), // predecessor: Non Zero if order must be respected
            bytes32(keccak256("SALT")), // salt
            minDelay
        );
        vm.warp(block.timestamp + minDelay + 1);
        vm.prank(executor);
        timelock.execute(
            target,
            0, // value
            data,
            bytes32(0), // predecessor: Non Zero if order must be respected
            bytes32(keccak256("SALT")) // salt
        );
    }

    /// @notice schedules, waits for timelock and executes a timelocked call, with a custom salt
    /// @dev This is to be used if we want to call the same target with the same data multiple times
    /// @param target The address to call
    /// @param data The data to call with
    /// @param salt The salt to use for the timelock
    function performTimelocked(address target, bytes memory data, bytes32 salt) internal {
        uint256 minDelay = timelock.getMinDelay();
        vm.prank(admin);
        timelock.schedule(
            target,
            0, // value
            data,
            bytes32(0), // predecessor: Non Zero if order must be respected
            bytes32(salt), // salt
            minDelay
        );
        vm.warp(block.timestamp + minDelay + 1);
        vm.prank(executor);
        timelock.execute(
            target,
            0, // value
            data,
            bytes32(0), // predecessor: Non Zero if order must be respected
            bytes32(salt) // salt
        );
    }

    /// @notice schedules a timelocked call
    /// @param target The address to call
    /// @param data The data to call with
    function schedule(address target, bytes memory data) internal {
        uint256 minDelay = timelock.getMinDelay();
        vm.prank(admin);
        timelock.schedule(
            target,
            0, // value
            data,
            bytes32(0), // predecessor: Non Zero if order must be respected
            bytes32(keccak256("SALT")), // salt
            minDelay
        );
    }

    /// @notice waits for timelock (minDelay)
    /// @dev This is to be called after schedule()
    /// If the scheduled time > minDelay, this wait time won't be enough
    /// and the test will revert
    function waitForTimelock() internal {
        uint256 minDelay = timelock.getMinDelay();
        vm.warp(block.timestamp + minDelay + 1);
    }

    /// @notice executes a timelocked call
    /// @param target The address to call
    /// @param data The data to call with
    function executeTimelocked(address target, bytes memory data) internal {
        vm.prank(executor);
        timelock.execute(
            target,
            0, // value
            data,
            bytes32(0), // predecessor: Non Zero if order must be respected
            bytes32(keccak256("SALT")) // salt
        );
    }

    /// @notice schedules, waits for timelock and executes a timelocked call that is expected to revert
    /// @param target The address to call
    /// @param data The data to call with
    /// @param reason The expected revert reason
    function expectRevertTimelocked(address target, bytes memory data, string memory reason) internal {
        uint256 minDelay = timelock.getMinDelay();
        vm.prank(admin);
        timelock.schedule(
            target,
            0, // value
            data,
            bytes32(0), // predecessor: Non Zero if order must be respected
            bytes32(keccak256("SALT")), // salt
            minDelay
        );
        vm.warp(block.timestamp + minDelay + 1);
        vm.prank(executor);
        vm.expectRevert(bytes(reason));
        timelock.execute(
            target,
            0, // value
            data,
            bytes32(0), // predecessor: Non Zero if order must be respected
            bytes32(keccak256("SALT")) // salt
        );
        // Cancel the scheduled call to clean the hash in the timelock
        bytes32 id = timelock.hashOperation(target, 0, data, bytes32(0), bytes32(keccak256("SALT")));
        vm.prank(admin);
        timelock.cancel(id);
    }
}
