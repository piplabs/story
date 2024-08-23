// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets

import { IIPTokenSlashing } from "../../src/protocol/IPTokenSlashing.sol";

import { Test } from "../utils/Test.sol";
import { MockIPTokenSlashingV2 } from "../utils/Mocks.sol";
import { EIP1967Helper } from "../../script/utils/EIP1967Helper.sol";

contract IPTokenSlashingTest is Test {
    bytes private delegatorUncmpPubkey =
        hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe47b"; // pragma: allowlist-secret
    // Address matching delegatorUncmpPubkey
    bytes private delegatorCmpPubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced1"; // pragma: allowlist-secret
    // Address matching delegatorCmpPubkey
    address private delegatorAddr = address(0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab);

    event Received(address, uint256);

    // For some tests, we need to receive the native token to this contract
    receive() external payable {
        emit Received(msg.sender, msg.value);
    }

    function setUp() public override {
        setStaking();
        setSlashing();
    }

    function testIPTokenSlashing_Parameters() public view {
        assertEq(ipTokenSlashing.unjailFee(), 1 ether);
    }

    function createDefaultValidator() private {
        vm.deal(delegatorAddr, 1 ether);
        vm.prank(delegatorAddr);
        ipTokenStaking.createValidator{ value: 1 ether }({
            validatorUncmpPubkey: delegatorUncmpPubkey,
            moniker: "delegator's validator",
            commissionRate: 1000,
            maxCommissionRate: 5000,
            maxCommissionChangeRate: 100
        });
    }

    function testIPTokenSlashing_Unjail() public {
        // Network shall not allow anyone to unjail an non-existing validator.
        uint256 feeAmount = 1 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenSlashing: Validator does not exist");
        ipTokenSlashing.unjail{ value: feeAmount }(delegatorUncmpPubkey);

        // Network shall not allow anyone to unjail a validator if it is not the validator itself.
        createDefaultValidator();
        address otherAddress = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenSlashing: Invalid pubkey derived address");
        ipTokenSlashing.unjail(delegatorUncmpPubkey);

        // Network shall not allow anyone to unjail a validator if the fee is not paid.
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenSlashing: Insufficient fee");
        ipTokenSlashing.unjail(delegatorUncmpPubkey);

        // Network shall not allow anyone to unjail a validator if the fee is not sufficient.
        feeAmount = 0.9 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenSlashing: Insufficient fee");
        ipTokenSlashing.unjail{ value: feeAmount }(delegatorUncmpPubkey);

        // Network shall allow anyone to unjail a validator if the fee is paid.
        feeAmount = 1 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectEmit(address(ipTokenSlashing));
        emit IIPTokenSlashing.Unjail(delegatorAddr, delegatorCmpPubkey);
        ipTokenSlashing.unjail{ value: feeAmount }(delegatorUncmpPubkey);

        // Network shall not allow anyone to unjail a validator if the fee is over.
        feeAmount = 1.1 ether;
        vm.deal(delegatorAddr, feeAmount);
        vm.prank(delegatorAddr);
        vm.expectRevert("IPTokenSlashing: Insufficient fee");
        ipTokenSlashing.unjail{ value: feeAmount }(delegatorUncmpPubkey);
    }

    function testIPTokenSlashing_UnjailOnBehalf() public {
        address otherAddress = address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA);

        // Network shall not allow anyone to unjail an non-existing validator.
        uint256 feeAmount = 1 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenSlashing: Validator does not exist");
        ipTokenSlashing.unjailOnBehalf{ value: feeAmount }(delegatorCmpPubkey);

        // Network shall not allow anyone to unjail with compressed pubkey of incorrect length.
        bytes memory delegatorCmpPubkeyShortLen = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ce"; // pragma: allowlist secret
        feeAmount = 1 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenSlashing: Invalid pubkey length");
        ipTokenSlashing.unjailOnBehalf{ value: feeAmount }(delegatorCmpPubkeyShortLen);

        // Network shall not allow anyone to unjail with compressed pubkey of incorrect prefix.
        bytes
            memory delegatorCmpPubkeyWrongPrefix = hex"05e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced1"; // pragma: allowlist secret
        feeAmount = 1 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenSlashing: Invalid pubkey prefix");
        ipTokenSlashing.unjailOnBehalf{ value: feeAmount }(delegatorCmpPubkeyWrongPrefix);

        // Network shall not allow anyone to unjail a validator if the fee is not paid.
        createDefaultValidator();
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenSlashing: Insufficient fee");
        ipTokenSlashing.unjailOnBehalf(delegatorCmpPubkey);

        // Network shall not allow anyone to unjail a validator if the fee is not sufficient.
        feeAmount = 0.9 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenSlashing: Insufficient fee");
        ipTokenSlashing.unjailOnBehalf{ value: feeAmount }(delegatorCmpPubkey);

        // Network shall allow anyone to unjail a validator on behalf if the fee is paid.
        feeAmount = 1 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectEmit(address(ipTokenSlashing));
        emit IIPTokenSlashing.Unjail(otherAddress, delegatorCmpPubkey);
        ipTokenSlashing.unjailOnBehalf{ value: feeAmount }(delegatorCmpPubkey);

        // Network shall not allow anyone to unjail a validator if the fee is over.
        feeAmount = 1.1 ether;
        vm.deal(otherAddress, feeAmount);
        vm.prank(otherAddress);
        vm.expectRevert("IPTokenSlashing: Insufficient fee");
        ipTokenSlashing.unjailOnBehalf{ value: feeAmount }(delegatorCmpPubkey);
    }

    function testIPTokenSlashing_SetUnjailFee() public {
        // Network shall allow the owner to set the unjail fee.
        uint256 newUnjailFee = 2 ether;
        vm.expectEmit(address(ipTokenSlashing));
        emit IIPTokenSlashing.UnjailFeeSet(newUnjailFee);
        vm.prank(admin);
        ipTokenSlashing.setUnjailFee(newUnjailFee);
        assertEq(ipTokenSlashing.unjailFee(), newUnjailFee);

        // Network shall not allow non-owner to set the unjail fee.
        vm.prank(address(0xf398c12A45BC409b6C652e25bb0A3e702492A4AA));
        vm.expectRevert();
        ipTokenSlashing.setUnjailFee(1 ether);
        assertEq(ipTokenSlashing.unjailFee(), newUnjailFee);
    }

    function testIPTokenSlashing_testOwnershipTransfer() public {
        address newOwner = address(0x444);
        vm.prank(admin);
        ipTokenSlashing.transferOwnership(newOwner);
        assertEq(ipTokenSlashing.pendingOwner(), newOwner);

        // Network shall not allow not new owner to accept the ownership.
        vm.expectRevert();
        ipTokenSlashing.acceptOwnership();
        assertEq(ipTokenSlashing.pendingOwner(), newOwner);

        // Network shall allow the new owner to accept the ownership.
        vm.prank(newOwner);
        ipTokenSlashing.acceptOwnership();
        assertEq(ipTokenSlashing.pendingOwner(), address(0));
        assertEq(ipTokenSlashing.owner(), newOwner);

        // Network shall not allow non-owner to transfer the ownership.
        vm.expectRevert();
        ipTokenSlashing.transferOwnership(newOwner);
        assertEq(ipTokenSlashing.pendingOwner(), address(0));
    }

    function testIPTokenSlashing_testUpgradeabilityACL() public {
        address newImpl = address(new MockIPTokenSlashingV2(address(ipTokenStaking)));
        // Network shall not allow non-owner to upgrade the contract.
        vm.expectRevert();
        ipTokenSlashing.upgradeToAndCall(newImpl, "");
        assertTrue(EIP1967Helper.getImplementation(address(ipTokenSlashing)) != newImpl);

        // Network shall not allow non-owner to disable the upgradeability.
        vm.expectRevert();
        ipTokenSlashing.disableUpgradeability();
        assertFalse(ipTokenSlashing.upgradeabilityDisabled());

        // Network shall allow the owner to disable the upgradeability.
        vm.prank(admin);
        ipTokenSlashing.disableUpgradeability();
        assertTrue(ipTokenSlashing.upgradeabilityDisabled());

        // Network shall not allow owner to upgrade the contract after the upgradeability is disabled.
        vm.expectRevert();
        ipTokenSlashing.upgradeToAndCall(newImpl, "");
        assertTrue(EIP1967Helper.getImplementation(address(ipTokenSlashing)) != newImpl);
    }

    function testIPTokenSlashing_testUpgradeability() public {
        address newImpl = address(new MockIPTokenSlashingV2(address(ipTokenStaking)));
        // Network shall allow the owner to upgrade the contract.
        vm.prank(admin);
        ipTokenSlashing.upgradeToAndCall(newImpl, "");
        assertTrue(EIP1967Helper.getImplementation(address(ipTokenSlashing)) == newImpl);
        assertTrue(MockIPTokenSlashingV2(address(ipTokenSlashing)).upgraded());
    }
}
