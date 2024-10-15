// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Test as ForgeTest } from "forge-std/Test.sol";

import { IPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { UpgradeEntrypoint } from "../../src/protocol/UpgradeEntrypoint.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";

import { GenerateAlloc } from "../../script/GenerateAlloc.s.sol";

contract Test is ForgeTest {
    address internal admin = address(0x123);
    address internal upgradeAdmin = address(0x456);

    IPTokenStaking internal ipTokenStaking;
    UpgradeEntrypoint internal upgradeEntrypoint;

    function setUp() public virtual {
        GenerateAlloc initializer = new GenerateAlloc();
        initializer.disableStateDump(); // Faster tests. Don't call to verify JSON output
        initializer.setAdminAddresses(upgradeAdmin, admin);
        initializer.run();
        ipTokenStaking = IPTokenStaking(Predeploys.Staking);
        upgradeEntrypoint = UpgradeEntrypoint(Predeploys.Upgrades);
    }
}
