// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { Test as ForgeTest } from "forge-std/Test.sol";

import { IPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { UpgradeEntrypoint } from "../../src/protocol/UpgradeEntrypoint.sol";
import { UBIPool } from "../../src/protocol/UBIPool.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";

import { GenerateAlloc } from "../../script/GenerateAlloc.s.sol";

contract Test is ForgeTest {
    address internal admin = address(0x123);
    address internal upgradeAdmin = address(0x456);

    IPTokenStaking internal ipTokenStaking;
    UpgradeEntrypoint internal upgradeEntrypoint;
    UBIPool internal ubiPool;

    function setUp() public virtual {
        GenerateAlloc initializer = new GenerateAlloc();
        initializer.disableStateDump(); // Faster tests. Don't call to verify JSON output
        initializer.setAdminAddresses(upgradeAdmin, admin);
        initializer.run();
        ipTokenStaking = IPTokenStaking(Predeploys.Staking);
        upgradeEntrypoint = UpgradeEntrypoint(Predeploys.Upgrades);
        ubiPool = UBIPool(Predeploys.UBIPool);
    }
}
