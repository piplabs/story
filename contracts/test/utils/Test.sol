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
    address internal admin = address(0x123);

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
                1 gwei, // stakingRounding
                1000, // defaultCommissionRate, 10%
                5000, // defaultMaxCommissionRate, 50%
                500 // defaultMaxCommissionChangeRate, 5%
            )
        );
        bytes memory initializer = abi.encodeCall(
            IPTokenStaking.initialize,
            (admin, 1 ether, 1 ether, 1 ether, 7 days)
        );
        ipTokenStaking = IPTokenStaking(address(new ERC1967Proxy(impl, initializer)));
    }

    function setSlashing() internal {
        require(address(ipTokenStaking) != address(0), "ipTokenStaking not set");

        address impl = address(new IPTokenSlashing(address(ipTokenStaking)));

        bytes memory initializer = abi.encodeCall(IPTokenSlashing.initialize, (admin, 1 ether));
        ipTokenSlashing = IPTokenSlashing(address(new ERC1967Proxy(impl, initializer)));

        console2.log("unjailFee:", ipTokenSlashing.unjailFee());
    }

    function setUpgrade() internal {
        address impl = address(new UpgradeEntrypoint());

        bytes memory initializer = abi.encodeWithSignature("initialize(address)", admin);
        upgradeEntrypoint = UpgradeEntrypoint(address(new ERC1967Proxy(impl, initializer)));
    }
}
