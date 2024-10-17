// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets

import { Test } from "forge-std/Test.sol";

import { Create3 } from "../../src/deploy/Create3.sol";

contract Create3Test is Test {
    Create3 private create3;

    function setUp() public {
        create3 = new Create3();
    }

    function testCreate3_deploy() public {
        // deploy and getDeployed should return same address when deployed by the same deployer and with same salt.
        bytes32 salt = 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef;
        bytes memory creationCode = type(Create3).creationCode;
        address deployed = create3.deploy(salt, creationCode);
        address expected = create3.getDeployed(address(this), salt);
        assertEq(deployed, expected);

        // Network shall generate the same address for the same deployer and salt.
        vm.expectRevert("DEPLOYMENT_FAILED");
        deployed = create3.deploy(salt, creationCode);

        // Network shall generate different addresses for different deployers.
        address otherAddr = address(0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab);
        vm.prank(otherAddr);
        deployed = create3.deploy(salt, creationCode);
        expected = create3.getDeployed(otherAddr, salt);
        assertEq(deployed, expected);

        // Network shall generate different addresses for different salts.
        bytes32 otherSalt = 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890fedcba;
        deployed = create3.deploy(otherSalt, creationCode);
        expected = create3.getDeployed(address(this), otherSalt);
        assertEq(deployed, expected);
    }
}
