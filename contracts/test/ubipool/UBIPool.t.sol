// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets
import { OwnableUpgradeable } from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

import { UBIPool } from "../../src/protocol/UBIPool.sol";
import { IUBIPool } from "../../src/interfaces/IUBIPool.sol";
import { Test } from "../utils/Test.sol";

contract UBIPoolTest is Test {
    function setUp() public virtual override {
        super.setUp();
    }

    function test_setUBIPercentage() public {
        // Fail if not protocol admin
        vm.expectRevert(abi.encodeWithSelector(OwnableUpgradeable.OwnableUnauthorizedAccount.selector, address(this)));
        ubiPool.setUBIPercentage(12 ether);

        // Fail if percentage too high
        vm.expectRevert("UBIPool: percentage too high");
        vm.prank(admin);
        ubiPool.setUBIPercentage(ubiPool.MAX_UBI_PERCENTAGE() + 1);

        // Set percentage
        vm.expectEmit(true, true, true, true);
        emit IUBIPool.UBIPercentageSet(22 ether);
        vm.prank(admin);
        ubiPool.setUBIPercentage(22 ether);
    }

    function test_setUBIDistribution() public {
        // Fail if not protocol admin
        vm.expectRevert(abi.encodeWithSelector(OwnableUpgradeable.OwnableUnauthorizedAccount.selector, address(this)));
        ubiPool.setUBIDistribution(1, 100 ether, new bytes[](0), new uint256[](0));

        // Fail if validatorUncmpPubKeys is empty
        vm.expectRevert("UBIPool: validatorUncmpPubKeys cannot be empty");
        vm.prank(admin);
        ubiPool.setUBIDistribution(1, 100 ether, new bytes[](0), new uint256[](0));

        // Fail if validatorUncmpPubKeys and percentages do not match
        vm.expectRevert("UBIPool: validatorUncmpPubKeys and percentages do not match");
        vm.prank(admin);
        ubiPool.setUBIDistribution(1, 100 ether, new bytes[](1), new uint256[](0));

        // Fail if percentages do not sum to 100%
        vm.expectRevert("UBIPool: percentages do not sum to 100%");
    }
}
