// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

/// @title RSASSAVerify
/// @notice Minimal RSASSA-PKCS#1 v1.5 / SHA-256 verifier for the TDX
///         Path-B bundle TPM2_Quote signature.
///
/// @dev This library implements *only* the verification surface required
///      by `TDXValidationHook.validateReport`'s STBN bundle path:
///      - RSA modulus length: 2048 bits (256 bytes).
///      - Public exponent: 65537 (0x010001), pinned at verification time.
///      - Signature length: 256 bytes (matches RSA-2048).
///      - Hash algorithm: SHA-256.
///      - Padding scheme: PKCS#1 v1.5 EMSA-PKCS1-v1_5.
///
///      Trust model:
///      - The modulus is extracted from a DER SubjectPublicKeyInfo
///        (RSA-2048) that sits *inside* the bundle, not from a chain-
///        managed registry. Cross-binding with the V4 quote (the
///        `report_data` AK-binding step in TDXValidationHook) is what
///        prevents an attacker from forging by swapping in their own
///        AK pub.
///      - The exponent is hard-coded to 65537. This matches every
///        OpenHCL-class paravisor AK we have observed and the
///        AKTemplate used by the kernel-side direct vendor (see
///        story-kernel/enclave/tdx/platform/tpmquote.go AKTemplate).
///        Bundles signed with a non-65537 exponent are rejected.
///
///      Algorithm:
///      - Compute m = signature^65537 mod modulus via the EIP-198
///        modular exponentiation precompile at address 0x05.
///      - Compare the resulting 256-byte buffer against the canonical
///        PKCS#1 v1.5 prefix for SHA-256 followed by the 32-byte
///        message digest. Constant-time per byte (keccak-equality).
///
///      Gas profile: the modexp precompile dominates cost. EIP-2565
///      makes RSA-2048 with e=65537 ~1.6k gas plus surrounding
///      memory shuffling — well below the 6 M gas budget the bundle
///      path targets.
library RSASSAVerify {
    /// @dev Length in bytes of an RSA-2048 public modulus and signature.
    uint256 internal constant RSA2048_BYTES = 256;

    /// @dev Length in bytes of a SHA-256 digest.
    uint256 internal constant SHA256_DIGEST_BYTES = 32;

    /// @dev DER PKIX SubjectPublicKeyInfo prefix length for an
    ///      RSA-2048 + e=65537 key. The full DER blob is exactly
    ///      33 (prefix) + 256 (modulus) + 5 (`02 03 01 00 01`
    ///      INTEGER trailer encoding the exponent 65537) = 294 bytes.
    ///      The on-chain extractor MUST validate the prefix bit-for-
    ///      bit so an attacker cannot smuggle a non-2048-bit key or
    ///      a non-65537 exponent past the verifier.
    uint256 internal constant SPKI_PREFIX_LEN = 33;

    /// @dev Total DER length the parser accepts for an RSA-2048
    ///      SubjectPublicKeyInfo. Anything else is rejected — we do
    ///      not support shorter or longer modulus encodings here.
    uint256 internal constant SPKI_RSA2048_LEN = 294;

    /*//////////////////////////////////////////////////////////////////////////
    //                            Public API                                  //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Verifies an RSASSA-PKCS#1 v1.5 / SHA-256 signature.
    /// @param messageDigest 32-byte SHA-256 digest of the signed message
    ///        (the caller hashes; this function does not).
    /// @param signature    256-byte RSA signature blob (big-endian).
    /// @param modulus      256-byte RSA-2048 modulus (big-endian).
    /// @return ok true if the signature verifies, false otherwise.
    /// @dev Constant-time at the block level — every branch leaks only
    ///      the public quantities (modulus, signature, prefix). Reverts
    ///      on malformed inputs (wrong lengths) so the caller can rely
    ///      on the boolean for cryptographic outcome only.
    function verify(
        bytes32 messageDigest,
        bytes memory signature,
        bytes memory modulus
    ) internal view returns (bool ok) {
        require(signature.length == RSA2048_BYTES, "RSASSAVerify: bad sig length");
        require(modulus.length == RSA2048_BYTES, "RSASSAVerify: bad modulus length");

        // EIP-198 modexp call: compute signature^65537 mod modulus.
        // Input layout (big-endian uint256s where required):
        //   [0..31]   base length   = 256
        //   [32..63]  exponent length = 3
        //   [64..95]  modulus length  = 256
        //   [96..]    base bytes      (signature)
        //   [..]      exponent bytes  (3 bytes: 0x01 0x00 0x01)
        //   [..]      modulus bytes
        bytes memory input = new bytes(96 + RSA2048_BYTES + 3 + RSA2048_BYTES);
        bytes memory output = new bytes(RSA2048_BYTES);

        assembly {
            // Layout offsets relative to data segment (skip 32B length prefix).
            let inPtr := add(input, 32)
            mstore(inPtr, RSA2048_BYTES) // base length
            mstore(add(inPtr, 32), 3) // exponent length
            mstore(add(inPtr, 64), RSA2048_BYTES) // modulus length
        }

        // Copy signature, exponent, modulus into the input buffer.
        _memcpy(_dataPtr(input, 96), _dataPtr(signature, 0), RSA2048_BYTES);
        // Exponent is 3 bytes: 0x01 0x00 0x01 (= 65537). We pin the
        // exponent here rather than reading it from the SPKI prefix
        // because every AK we accept (paravisor + kernel direct) uses
        // 65537 and admitting anything else would require the SPKI
        // exponent suffix to be variable-length.
        assembly {
            let expPtr := add(add(input, 32), add(96, 256))
            mstore8(expPtr, 0x01)
            mstore8(add(expPtr, 1), 0x00)
            mstore8(add(expPtr, 2), 0x01)
        }
        _memcpy(_dataPtr(input, 96 + RSA2048_BYTES + 3), _dataPtr(modulus, 0), RSA2048_BYTES);

        bool success;
        assembly {
            success := staticcall(
                gas(),
                0x05, // EIP-198 modexp precompile
                add(input, 32),
                mload(input),
                add(output, 32),
                RSA2048_BYTES
            )
        }
        require(success, "RSASSAVerify: modexp failed");

        // Compare output against the canonical PKCS#1 v1.5 envelope:
        //   0x00 0x01 (PS=0xFF * 202 bytes) 0x00 (DER ASN.1 prefix 19 bytes) (SHA-256 digest 32 bytes)
        // Total length is 256 bytes. The DER ASN.1 prefix encodes:
        //   30 31 30 0D 06 09 60 86 48 01 65 03 04 02 01 05 00 04 20
        // which is the SHA-256 DigestInfo with a 32-byte OCTET STRING.
        ok = _isValidPKCS1v15Envelope(output, messageDigest);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                       DER SPKI modulus extraction                      //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Extracts the 256-byte RSA modulus from a DER PKIX
    ///         SubjectPublicKeyInfo (RSA-2048, exponent 65537).
    /// @param der The DER-encoded SubjectPublicKeyInfo (must be exactly
    ///        294 bytes; prefix and exponent suffix are validated bit-
    ///        for-bit).
    /// @return modulus The 256-byte modulus (a fresh memory buffer).
    /// @dev Validating the full DER skeleton lets us treat the modulus
    ///      offset as a constant rather than re-parsing TLV. The cost
    ///      is rejecting any DER encoding that re-orders or pads the
    ///      structure — but x509.MarshalPKIXPublicKey produces a
    ///      canonical encoding that matches exactly.
    function extractRSA2048Modulus(bytes memory der) internal pure returns (bytes memory modulus) {
        require(der.length == SPKI_RSA2048_LEN, "RSASSAVerify: bad DER length");

        // Validate the canonical RSA-2048 SPKI prefix.
        // 30 82 01 22 -- SEQUENCE (0x122 = 290 bytes follow)
        // 30 0D       -- SEQUENCE (AlgorithmIdentifier, 13 bytes)
        // 06 09       -- OID, 9 bytes
        // 2A 86 48 86 F7 0D 01 01 01 -- 1.2.840.113549.1.1.1 (rsaEncryption)
        // 05 00       -- NULL (parameters)
        // 03 82 01 0F -- BIT STRING (271 bytes follow)
        // 00          -- unused bits = 0
        // 30 82 01 0A -- SEQUENCE (RSAPublicKey, 266 bytes)
        // 02 82 01 01 -- INTEGER (modulus, 257 bytes follow)
        // 00          -- INTEGER sign-bit prefix (modulus high bit set)
        // SEQUENCE 0x122B  30 82 01 22
        // SEQUENCE 0x0DB    30 0D
        // OID rsaEncryption 06 09 2A 86 48 86 F7 0D 01 01 01
        // NULL              05 00
        // BIT STRING 0x10FB 03 82 01 0F
        // unused-bits=0     00
        // SEQUENCE 0x10AB   30 82 01 0A   (RSAPublicKey)
        // INTEGER 0x101B    02 82 01 01   (modulus, leading 0x00 sign byte)
        // sign byte         00
        // (next 256 bytes are the modulus, then 02 03 01 00 01 = e=65537)
        bytes memory expectedPrefix = (hex"30820122300D06092A864886F70D01010105000382010F003082010A0282010100");
        require(expectedPrefix.length == SPKI_PREFIX_LEN, "RSASSAVerify: prefix length bug");
        for (uint256 i = 0; i < SPKI_PREFIX_LEN; i++) {
            require(der[i] == expectedPrefix[i], "RSASSAVerify: bad SPKI prefix");
        }

        // Validate exponent suffix: 02 03 01 00 01 (INTEGER, 3 bytes, 65537).
        require(
            der[SPKI_PREFIX_LEN + RSA2048_BYTES + 0] == 0x02 &&
                der[SPKI_PREFIX_LEN + RSA2048_BYTES + 1] == 0x03 &&
                der[SPKI_PREFIX_LEN + RSA2048_BYTES + 2] == 0x01 &&
                der[SPKI_PREFIX_LEN + RSA2048_BYTES + 3] == 0x00 &&
                der[SPKI_PREFIX_LEN + RSA2048_BYTES + 4] == 0x01,
            "RSASSAVerify: bad SPKI exponent"
        );

        // Extract the 256-byte modulus into a fresh buffer.
        modulus = new bytes(RSA2048_BYTES);
        for (uint256 i = 0; i < RSA2048_BYTES; i++) {
            modulus[i] = der[SPKI_PREFIX_LEN + i];
        }
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                          Internal helpers                              //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Verifies the modexp output is a well-formed PKCS#1 v1.5
    ///      EMSA envelope wrapping the supplied SHA-256 digest. The
    ///      envelope is fully determined by digest length (32) and
    ///      modulus length (256), so we reconstruct it byte-for-byte
    ///      and compare via keccak256 — naturally constant-time.
    function _isValidPKCS1v15Envelope(bytes memory candidate, bytes32 digest) private pure returns (bool) {
        if (candidate.length != RSA2048_BYTES) return false;

        // Build the expected envelope:
        //   0x00 0x01 0xFF * 202 0x00 (DigestInfo 19 bytes) (digest 32 bytes)
        // Total length = 2 + 202 + 1 + 19 + 32 = 256.
        bytes memory expected = new bytes(RSA2048_BYTES);
        expected[0] = 0x00;
        expected[1] = 0x01;
        for (uint256 i = 2; i < 204; i++) {
            expected[i] = 0xFF;
        }
        expected[204] = 0x00;
        // SHA-256 DigestInfo DER: 30 31 30 0D 06 09 60 86 48 01 65 03 04 02 01 05 00 04 20
        bytes19 digestInfoPrefix = 0x3031300D060960864801650304020105000420;
        for (uint256 i = 0; i < 19; i++) {
            expected[205 + i] = digestInfoPrefix[i];
        }
        // Append the digest.
        for (uint256 i = 0; i < SHA256_DIGEST_BYTES; i++) {
            expected[224 + i] = digest[i];
        }

        return keccak256(candidate) == keccak256(expected);
    }

    /// @dev Returns the address of `bytes memory data` element at
    ///      `offset` (skipping the 32-byte length prefix).
    function _dataPtr(bytes memory data, uint256 offset) private pure returns (uint256 ptr) {
        assembly {
            ptr := add(add(data, 32), offset)
        }
    }

    /// @dev Word-aware memory copy. Used for moving the signature and
    ///      modulus bytes into the modexp input buffer; safe because
    ///      both sources and the destination are freshly allocated.
    function _memcpy(uint256 dest, uint256 src, uint256 len) private pure {
        for (; len >= 32; len -= 32) {
            assembly {
                mstore(dest, mload(src))
            }
            dest += 32;
            src += 32;
        }
        if (len > 0) {
            uint256 mask = 256 ** (32 - len) - 1;
            assembly {
                let srcpart := and(mload(src), not(mask))
                let destpart := and(mload(dest), mask)
                mstore(dest, or(destpart, srcpart))
            }
        }
    }
}
