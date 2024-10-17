// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/* solhint-disable function-state-mutability */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets

import { Test } from "forge-std/Test.sol";

import { Secp256k1 } from "../../src/libraries/Secp256k1.sol";

contract IPTokenStakingTest is Test {
    function setUp() public {}

    function testCompressPublicKey_validKey() public pure {
        // prefix: 04
        bytes
            memory uncmpPubkey = hex"04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe47b"; // pragma: allowlist-secret
        // prefix: 03
        bytes memory cmpPubkey = hex"03e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced1"; // pragma: allowlist-secret

        bytes
            memory anotherUncmpPubkey = hex"04e38d111122223333ce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe47b"; // pragma: allowlist-secret

        vm.assertEq(Secp256k1.compressPublicKey(uncmpPubkey), cmpPubkey);
        vm.assertNotEq(Secp256k1.compressPublicKey(uncmpPubkey), Secp256k1.compressPublicKey(anotherUncmpPubkey));
    }

    function testCompressPublicKey_deriveAddress() public pure {
        // prefix 04 sliced from `04e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe47b` // pragma: allowlist-secret
        bytes
            memory uncmpPubkeySliced = hex"e38d15ae6cc5d41cce27a2307903cb12a406cbf463fe5fef215bdf8aa988ced195e9327ac89cd362eaa0397f8d7f007c02b2a75642f174e455d339e4a1efe47b"; // pragma: allowlist-secret

        address expectedAddr = 0xf398C12A45Bc409b6C652E25bb0a3e702492A4ab;
        address derivedAddr = address(uint160(uint256(keccak256(uncmpPubkeySliced))));

        vm.assertEq(uncmpPubkeySliced.length, 64);
        vm.assertEq(derivedAddr, expectedAddr);
    }
}
