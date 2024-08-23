// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets

import { UpgradeEntrypoint, IUpgradeEntrypoint } from "../../src/protocol/UpgradeEntrypoint.sol";
import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { Test } from "../utils/Test.sol";
import { MockUpgradeEntryPointV2 } from "../utils/Mocks.sol";
import { EIP1967Helper } from "../../script/utils/EIP1967Helper.sol";

contract UpgradeEntrypointTest is Test {
    function setUp() public override {
        address impl = address(new UpgradeEntrypoint());
        bytes memory initializer = abi.encodeCall(UpgradeEntrypoint.initialize, (admin));
        upgradeEntrypoint = UpgradeEntrypoint(address(new ERC1967Proxy(impl, initializer)));
    }

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

    function testUpgradeEntrypoint_OwnershipTransfer() public {
        address newOwner = address(0x444);
        vm.prank(admin);
        upgradeEntrypoint.transferOwnership(newOwner);
        assertEq(upgradeEntrypoint.pendingOwner(), newOwner);

        // contract shall not allow not new owner to accept the ownership.
        vm.expectRevert();
        upgradeEntrypoint.acceptOwnership();
        assertEq(upgradeEntrypoint.pendingOwner(), newOwner);

        // contract shall allow the new owner to accept the ownership.
        vm.prank(newOwner);
        upgradeEntrypoint.acceptOwnership();
        assertEq(upgradeEntrypoint.pendingOwner(), address(0));
        assertEq(upgradeEntrypoint.owner(), newOwner);

        // contract shall not allow non-owner to transfer the ownership.
        vm.expectRevert();
        upgradeEntrypoint.transferOwnership(newOwner);
        assertEq(upgradeEntrypoint.pendingOwner(), address(0));
    }

    function testUpgradeEntrypoint_testUpgradeabilityACL() public {
        address newImpl = address(new MockUpgradeEntryPointV2());
        // Network shall not allow non-owner to upgrade the contract.
        vm.expectRevert();
        upgradeEntrypoint.upgradeToAndCall(newImpl, "");
        assertTrue(EIP1967Helper.getImplementation(address(upgradeEntrypoint)) != newImpl);

        // Network shall not allow non-owner to disable the upgradeability.
        vm.expectRevert();
        upgradeEntrypoint.disableUpgradeability();
        assertFalse(upgradeEntrypoint.upgradeabilityDisabled());

        // Network shall allow the owner to disable the upgradeability.
        vm.prank(admin);
        upgradeEntrypoint.disableUpgradeability();
        assertTrue(upgradeEntrypoint.upgradeabilityDisabled());

        // Network shall not allow owner to upgrade the contract after the upgradeability is disabled.
        vm.expectRevert();
        upgradeEntrypoint.upgradeToAndCall(newImpl, "");
        assertTrue(EIP1967Helper.getImplementation(address(upgradeEntrypoint)) != newImpl);
    }

    function testUpgradeEntrypoint_testUpgradeability() public {
        address newImpl = address(new MockUpgradeEntryPointV2());
        // Network shall allow the owner to upgrade the contract.
        vm.prank(admin);
        upgradeEntrypoint.upgradeToAndCall(newImpl, "");
        assertTrue(EIP1967Helper.getImplementation(address(upgradeEntrypoint)) == newImpl);
        assertTrue(MockUpgradeEntryPointV2(address(upgradeEntrypoint)).upgraded());
    }
}
