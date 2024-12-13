// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
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
        expectRevertTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(IUBIPool.setUBIPercentage.selector, tooMuch),
            "UBIPool: percentage too high"
        );

        // Set percentage
        schedule(address(ubiPool), abi.encodeWithSelector(IUBIPool.setUBIPercentage.selector, 12_00));
        waitForTimelock();
        vm.expectEmit(true, true, true, true);
        emit IUBIPool.UBIPercentageSet(12_00);
        executeTimelocked(address(ubiPool), abi.encodeWithSelector(IUBIPool.setUBIPercentage.selector, 12_00));
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
        performTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(IUBIPool.setUBIDistribution.selector, totalAmount, validatorUncmpPubKeys, amounts)
        );
        assertEq(ubiPool.currentDistributionId(), 1);

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
        performTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(IUBIPool.setUBIDistribution.selector, totalAmount, validatorUncmpPubKeys, amounts),
            bytes32(keccak256("setUBIDistribution 2nd time"))
        );
        assertEq(ubiPool.currentDistributionId(), 2);

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
        expectRevertTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(IUBIPool.setUBIDistribution.selector, 100 ether, new bytes[](0), new uint256[](0)),
            "UBIPool: validatorUncmpPubKeys cannot be empty"
        );

        // Fail if validatorUncmpPubKeys and percentages do not match
        expectRevertTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(IUBIPool.setUBIDistribution.selector, 100 ether, new bytes[](1), new uint256[](0)),
            "UBIPool: length mismatch"
        );

        // Fail if UBIPool: not enough balance
        expectRevertTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(IUBIPool.setUBIDistribution.selector, 100 ether, new bytes[](1), new uint256[](1)),
            "UBIPool: not enough balance"
        );

        // Fail if amounts do not sum to totalUBI

        uint256[] memory amounts = new uint256[](1);
        bytes[] memory validatorUncmpPubKeys = new bytes[](1);
        validatorUncmpPubKeys[0] = validators[0].uncompressedHex;
        amounts[0] = 1 ether;
        vm.deal(address(ubiPool), 100 ether);
        expectRevertTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(IUBIPool.setUBIDistribution.selector, 100 ether, validatorUncmpPubKeys, amounts),
            "UBIPool: total amount mismatch"
        );

        // Fail if one amount is zero
        vm.deal(address(ubiPool), 100 ether);
        amounts = new uint256[](1);
        amounts[0] = 0;
        validatorUncmpPubKeys = new bytes[](1);
        validatorUncmpPubKeys[0] = validators[0].uncompressedHex;
        expectRevertTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(IUBIPool.setUBIDistribution.selector, 100 ether, validatorUncmpPubKeys, amounts),
            "UBIPool: amounts cannot be zero"
        );

        // Fail if pubkey is not valid
        vm.deal(address(ubiPool), 100 ether);
        amounts = new uint256[](1);
        amounts[0] = 100 ether;
        validatorUncmpPubKeys = new bytes[](1);
        // Invalid pubkey
        validatorUncmpPubKeys[
            0
        ] = hex"0482782124bc9cd03c38aa4cac234dc4e4e3cecf04d57914371baf7fa78ffb975f6d58e245bea952dd039f0fec4e9db418c3b00000"; // pragma: allowlist secret
        expectRevertTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(IUBIPool.setUBIDistribution.selector, 100 ether, validatorUncmpPubKeys, amounts),
            "PubKeyVerifier: Invalid pubkey length"
        );
    }

    function test_claimUBI_revert_respect_unclaimed_tokens() public {
        // Reverts if balance < totalPendingClaims + totalUBI
        // Set initial distribution and claim some
        uint256 totalAmount = 100 ether;
        uint256[] memory amounts = new uint256[](2);
        amounts[0] = 60 ether;
        amounts[1] = 40 ether;
        bytes[] memory validatorUncmpPubKeys = new bytes[](2);
        validatorUncmpPubKeys[0] = validators[0].uncompressedHex;
        validatorUncmpPubKeys[1] = validators[1].uncompressedHex;

        vm.deal(address(ubiPool), totalAmount);
        performTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(IUBIPool.setUBIDistribution.selector, totalAmount, validatorUncmpPubKeys, amounts),
            bytes32(keccak256("setUBIDistribution"))
        );

        // First validator claims their UBI
        vm.prank(validators[0].evmAddress);
        ubiPool.claimUBI(1, validatorUncmpPubKeys[0]);

        // Try to set new distribution with not enough balance
        // Only 40 ether left in contract but 60 ether still pending claims
        uint256 newTotalAmount = 1 ether;
        uint256[] memory newAmounts = new uint256[](1);
        newAmounts[0] = 1 ether;
        bytes[] memory newValidatorUncmpPubKeys = new bytes[](1);
        newValidatorUncmpPubKeys[0] = validators[0].uncompressedHex;

        expectRevertTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(
                IUBIPool.setUBIDistribution.selector,
                newTotalAmount,
                newValidatorUncmpPubKeys,
                newAmounts
            ),
            "UBIPool: not enough balance"
        );
    }

    function test_claimUBI_respect_unclaimed_tokens() public {
        // Set initial distribution
        uint256 totalAmount = 100 ether;
        uint256[] memory amounts = new uint256[](2);
        amounts[0] = 60 ether;
        amounts[1] = 40 ether;
        bytes[] memory validatorUncmpPubKeys = new bytes[](2);
        validatorUncmpPubKeys[0] = validators[0].uncompressedHex;
        validatorUncmpPubKeys[1] = validators[1].uncompressedHex;

        vm.deal(address(ubiPool), totalAmount + 10 ether); // Extra balance for next distribution
        performTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(IUBIPool.setUBIDistribution.selector, totalAmount, validatorUncmpPubKeys, amounts),
            bytes32(keccak256("setUBIDistribution"))
        );

        // First validator claims their UBI
        vm.prank(validators[0].evmAddress);
        ubiPool.claimUBI(1, validatorUncmpPubKeys[0]);

        // Set new distribution with enough balance
        // 40 ether pending claims + 10 ether extra balance = enough for 5 ether new distribution
        uint256 newTotalAmount = 5 ether;
        uint256[] memory newAmounts = new uint256[](1);
        newAmounts[0] = 5 ether;
        bytes[] memory newValidatorUncmpPubKeys = new bytes[](1);
        newValidatorUncmpPubKeys[0] = validators[0].uncompressedHex;

        performTimelocked(
            address(ubiPool),
            abi.encodeWithSelector(
                IUBIPool.setUBIDistribution.selector,
                newTotalAmount,
                newValidatorUncmpPubKeys,
                newAmounts
            ),
            bytes32(keccak256("setUBIDistribution"))
        );
    }
}
