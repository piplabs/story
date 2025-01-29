// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;
/* solhint-disable no-console */
/* solhint-disable max-line-length */
/* solhint-disable function-state-mutability */
/// NOTE: pragma allowlist-secret must be inline (same line as the pubkey hex string) to avoid false positive
/// flag "Hex High Entropy String" in CI run detect-secrets

import { Test } from "forge-std/Test.sol";

import { Secp256k1Verifier } from "../../src/protocol/Secp256k1Verifier.sol";

contract Secp256k1VerifierHarness is Secp256k1Verifier {
    function uncompressPublicKey(bytes memory compressedKey) public pure returns (bytes memory) {
        return _uncompressPublicKey(compressedKey);
    }
}

contract Secp256k1VerifierTest is Test {
    Secp256k1VerifierHarness public verifier;

    function setUp() public {
        verifier = new Secp256k1VerifierHarness();
    }

    function testCompressPublicKey_validKey() public view {
        // prefix: 04
        bytes
            memory uncmpPubkey = hex"04ed58a9319aba87f60fe08e87bc31658dda6bfd7931686790a2ff803846d4e59c215b515f2acba1de2979c9b1376e088d4e48b20331a876ae4ed2c4f6bafc4016"; // pragma: allowlist-secret
        // prefix: 03
        bytes memory cmpPubkey = hex"02ed58a9319aba87f60fe08e87bc31658dda6bfd7931686790a2ff803846d4e59c"; // pragma: allowlist-secret

        bytes memory anotherCmpPubkey = hex"037ff1214f5af4b652bc6c352ecb1296791cf754a426b49a7c4a263124f7497e98"; // pragma: allowlist-secret

        vm.assertEq(verifier.uncompressPublicKey(cmpPubkey), uncmpPubkey);
        vm.assertNotEq(verifier.uncompressPublicKey(cmpPubkey), verifier.uncompressPublicKey(anotherCmpPubkey));
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
