// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets

import { UpgradeEntrypoint, IUpgradeEntrypoint } from "../../src/protocol/UpgradeEntrypoint.sol";
import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import { Test } from "../utils/Test.sol";

contract UpgradeEntrypointTest is Test {

    function testUpgradeEntrypoint_planUpgrade() public {
        // Network shall allow the protocol owner to submit an upgrade plan.
        string memory name = "upgrade";
        int64 height = 1;
        string memory info = "info";

        vm.expectEmit(address(upgradeEntrypoint));
        emit IUpgradeEntrypoint.SoftwareUpgrade(name, height, info);
        vm.prank(admin);
        upgradeEntrypoint.planUpgrade(name, height, info);

        // Network shall not allow non-protocol owner to submit an upgrade plan.
        address otherAddr = address(0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab);
        vm.prank(otherAddr);
        vm.expectRevert();
        upgradeEntrypoint.planUpgrade(name, height, info);
    }
}
