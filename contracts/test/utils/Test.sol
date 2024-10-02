// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { console2 } from "forge-std/console2.sol";
import { Test as ForgeTest } from "forge-std/Test.sol";
import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import { IPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { IPTokenSlashing } from "../../src/protocol/IPTokenSlashing.sol";
import { UpgradeEntrypoint } from "../../src/protocol/UpgradeEntrypoint.sol";
import { Predeploys } from "../../src/libraries/Predeploys.sol";

import { EtchInitialState } from "../../script/EtchInitialState.s.sol";

contract Test is ForgeTest {
    address internal admin = address(0x123);
    address internal upgradeAdmin = address(0x456);

    IPTokenStaking internal ipTokenStaking;
    IPTokenSlashing internal ipTokenSlashing;
    UpgradeEntrypoint internal upgradeEntrypoint;

    function setUp() virtual public {
        EtchInitialState initializer = new EtchInitialState();
        initializer.disableStateDump(); // Faster tests. Don't call to verify JSON output
        initializer.run();
        ipTokenStaking = IPTokenStaking(Predeploys.Staking);
        ipTokenSlashing = IPTokenSlashing(Predeploys.Slashing);
        upgradeEntrypoint = UpgradeEntrypoint(Predeploys.Upgrades);
    }
}
