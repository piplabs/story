// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { Test } from "../utils/Test.sol";
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets

contract ERC6551RegistryTest is Test {
    function test_erc6551Registry() public {
        address account = erc6551Registry.createAccount(
            address(this),
            0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef,
            block.chainid,
            address(this),
            1
        );

        assertNotEq(account, address(0));

        assertEq(
            account,
            erc6551Registry.account(
                address(this),
                0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef,
                block.chainid,
                address(this),
                1
            )
        );
    }
}
