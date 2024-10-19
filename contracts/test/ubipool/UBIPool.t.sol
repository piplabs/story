// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets
import { OwnableUpgradeable } from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

import { IUBIPool } from "../../src/interfaces/IUBIPool.sol";
import { Test } from "../utils/Test.sol";
import { ValidatorData } from "../data/ValidatorData.sol";

contract UBIPoolTest is Test, ValidatorData {
    function setUp() public virtual override {
        super.setUp();
    }

    function test_setUBIPercentage() public {
        // Fail if not protocol admin
        vm.expectRevert(abi.encodeWithSelector(OwnableUpgradeable.OwnableUnauthorizedAccount.selector, address(this)));
        ubiPool.setUBIPercentage(12_00);

        // Fail if percentage too high
        uint32 tooMuch = ubiPool.MAX_UBI_PERCENTAGE() + 1;
        vm.expectRevert("UBIPool: percentage too high");
        vm.prank(admin);
        ubiPool.setUBIPercentage(tooMuch);

        // Set percentage
        vm.expectEmit(true, true, true, true);
        emit IUBIPool.UBIPercentageSet(12_00);
        vm.prank(admin);
        ubiPool.setUBIPercentage(12_00);
    }

    function test_setUBIDistribution_claimUBI() public {
        uint256[] memory amounts = new uint256[](validators.length);
        bytes[] memory validatorUncmpPubKeys = new bytes[](validators.length);
        uint256 totalAmount = 0;
        for (uint256 i = 0; i < validators.length; i++) {
            amounts[i] = 100 ether + i * 10 ether;
            validatorUncmpPubKeys[i] = validators[i].uncompressedHex;
            totalAmount += amounts[i];
            vm.prank(validators[i].evmAddress);
            vm.expectRevert("UBIPool: no UBI to claim");
            ubiPool.claimUBI(1, validatorUncmpPubKeys[i]);
        }
        vm.deal(address(ubiPool), totalAmount);
        vm.prank(admin);
        uint256 distributionId = ubiPool.setUBIDistribution(totalAmount, validatorUncmpPubKeys, amounts);
        assertEq(distributionId, 1);
        for (uint256 i = 0; i < validators.length; i++) {
            uint256 amount = amounts[i];
            assertEq(ubiPool.validatorUBIAmounts(1, validatorUncmpPubKeys[i]), amount);
            vm.prank(validators[i].evmAddress);
            uint256 balanceBefore = address(validators[i].evmAddress).balance;
            uint256 poolBalanceBefore = address(ubiPool).balance;
            ubiPool.claimUBI(1, validatorUncmpPubKeys[i]);
            assertEq(address(validators[i].evmAddress).balance, balanceBefore + amount);
            assertEq(address(ubiPool).balance, poolBalanceBefore - amount);
            assertEq(ubiPool.validatorUBIAmounts(1, validatorUncmpPubKeys[i]), 0);
        }
        vm.deal(address(ubiPool), totalAmount);
        vm.prank(admin);
        assertEq(ubiPool.setUBIDistribution(totalAmount, validatorUncmpPubKeys, amounts), 2);
        for (uint256 i = 0; i < validators.length; i++) {
            uint256 amount = amounts[i];
            assertEq(ubiPool.validatorUBIAmounts(2, validatorUncmpPubKeys[i]), amount);
            vm.prank(validators[i].evmAddress);
            uint256 balanceBefore = address(validators[i].evmAddress).balance;
            uint256 poolBalanceBefore = address(ubiPool).balance;
            ubiPool.claimUBI(2, validatorUncmpPubKeys[i]);
            assertEq(address(validators[i].evmAddress).balance, balanceBefore + amount);
            assertEq(address(ubiPool).balance, poolBalanceBefore - amount);
            assertEq(ubiPool.validatorUBIAmounts(2, validatorUncmpPubKeys[i]), 0);
        }
    }

    function test_setUBIDistribution_reverts() public {
        // Fail if not protocol admin
        vm.expectRevert(abi.encodeWithSelector(OwnableUpgradeable.OwnableUnauthorizedAccount.selector, address(this)));
        ubiPool.setUBIDistribution(100 ether, new bytes[](0), new uint256[](0));

        // Fail if validatorUncmpPubKeys is empty
        vm.expectRevert("UBIPool: validatorUncmpPubKeys cannot be empty");
        vm.prank(admin);
        ubiPool.setUBIDistribution(100 ether, new bytes[](0), new uint256[](0));

        // Fail if validatorUncmpPubKeys and percentages do not match
        vm.expectRevert("UBIPool: length mismatch");
        vm.prank(admin);
        ubiPool.setUBIDistribution(100 ether, new bytes[](1), new uint256[](0));

        // Fail if UBIPool: not enough balance
        vm.expectRevert("UBIPool: not enough balance");
        vm.prank(admin);
        ubiPool.setUBIDistribution(100 ether, new bytes[](1), new uint256[](1));

        // Fail if amounts do not sum to totalUBI
        vm.expectRevert("UBIPool: total amount mismatch");
        uint256[] memory amounts = new uint256[](1);
        bytes[] memory validatorUncmpPubKeys = new bytes[](1);
        validatorUncmpPubKeys[0] = validators[0].uncompressedHex;
        amounts[0] = 1 ether;
        vm.deal(address(ubiPool), 100 ether);
        vm.prank(admin);
        ubiPool.setUBIDistribution(100 ether, validatorUncmpPubKeys, amounts);

        // Fail if one amount is zero
        vm.deal(address(ubiPool), 100 ether);
        amounts = new uint256[](1);
        amounts[0] = 0;
        validatorUncmpPubKeys = new bytes[](1);
        validatorUncmpPubKeys[0] = validators[0].uncompressedHex;
        vm.expectRevert("UBIPool: amounts cannot be zero");
        vm.prank(admin);
        ubiPool.setUBIDistribution(100 ether, validatorUncmpPubKeys, amounts);

        // Fail if pubkey is not valid
        vm.deal(address(ubiPool), 100 ether);
        amounts = new uint256[](1);
        amounts[0] = 100 ether;
        validatorUncmpPubKeys = new bytes[](1);
        // Invalid pubkey
        validatorUncmpPubKeys[
            0
        ] = hex"0482782124bc9cd03c38aa4cac234dc4e4e3cecf04d57914371baf7fa78ffb975f6d58e245bea952dd039f0fec4e9db418c3b00000"; // pragma: allowlist secret
        vm.expectRevert("PubKeyVerifier: Invalid pubkey length");
        vm.prank(admin);
        ubiPool.setUBIDistribution(100 ether, validatorUncmpPubKeys, amounts);
    }
}
