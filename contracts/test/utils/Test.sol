// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */

import { console2 } from "forge-std/console2.sol";
import { Test as ForgeTest } from "forge-std/Test.sol";
import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { IPTokenStaking } from "../../src/protocol/IPTokenStaking.sol";
import { IPTokenSlashing } from "../../src/protocol/IPTokenSlashing.sol";
import { UpgradeEntrypoint } from "../../src/protocol/UpgradeEntrypoint.sol";

contract Test is ForgeTest {
    address private admin = address(this);

    IPTokenStaking internal ipTokenStaking;
    IPTokenSlashing internal ipTokenSlashing;
    UpgradeEntrypoint internal upgradeEntrypoint;

    function setUp() public virtual {
        setStaking();
        setSlashing();
        // setUpgrade();
    }

    function setStaking() internal {
        address impl = address(
            new IPTokenStaking(
                1 ether, // minStakeAmount
                1 ether, // minUnstakeAmount
                1 ether, // minRedelegateAmount
                1 gwei, // stakingRounding
                7 days, // withdrawalAddressChangeInterval
                1000, // defaultCommissionRate, 10%
                5000, // defaultMaxCommissionRate, 50%
                500 // defaultMaxCommissionChangeRate, 5%
            )
        );

        ipTokenStaking = IPTokenStaking(address(new ERC1967Proxy(impl, "")));
        ipTokenStaking.initialize(admin);
    }

    function setSlashing() internal {
        require(address(ipTokenStaking) != address(0), "ipTokenStaking not set");

        address impl = address(
            new IPTokenSlashing(
                address(ipTokenStaking),
                1 ether // unjailFee
            )
        );

        ipTokenSlashing = IPTokenSlashing(address(new ERC1967Proxy(impl, "")));
        ipTokenSlashing.initialize(admin);

        console2.log("unjailFee:", ipTokenSlashing.unjailFee());
    }

    function setUpgrade() internal {
        address impl = address(new UpgradeEntrypoint());

        upgradeEntrypoint = UpgradeEntrypoint(address(new ERC1967Proxy(impl, "")));
        upgradeEntrypoint.initialize(admin);
    }
}
