// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

/// @title TPM2Attest
/// @notice Minimal parser for TPMS_ATTEST blobs emitted by TPM2_Quote.
///
/// @dev This library extracts only the field the TDX bundle path needs
///      to authorize: the `qualifyingData` (TPMS_ATTEST.extraData), which
///      is the user-supplied 32-byte commitment the TPM signed as part
///      of the TPM2_Quote ceremony.
///
///      Wire layout (TCG TPM 2.0 Library, Part 2 §10.12):
///
///          uint32  magic                    // TPM_GENERATED_VALUE
///          uint16  type                     // TPM_ST_ATTEST_QUOTE
///          TPM2B_NAME      qualifiedSigner  // uint16 size + bytes
///          TPM2B_DATA      extraData        // uint16 size + bytes  <-- qualifyingData
///          TPMS_CLOCK_INFO clockInfo        // 17 bytes
///          uint64          firmwareVersion
///          TPMS_QUOTE_INFO attested         // PCRSelections + digest
///
///      All multi-byte integers are network byte order (big-endian).
///      The on-chain parser stops after extracting `extraData`; the
///      parser does NOT validate the trailing fields because (a) they
///      are not load-bearing for our trust property, and (b) the TPM
///      itself signs the entire blob so any tampering after the prefix
///      would already have been caught at signature verification.
///
///      Constants pinned:
///      - magic == TPM_GENERATED_VALUE (0xFF544347): proves the blob
///        was constructed by a TPM with access to the AK private key,
///        not synthesized by software.
///      - type == TPM_ST_ATTEST_QUOTE (0x8018): rejects other attest
///        types (TPM_ST_ATTEST_CREATION, TPM_ST_ATTEST_CERTIFY, etc.)
///        which have different field layouts.
library TPM2Attest {
    /// @dev TPM_GENERATED_VALUE — the magic the TPM stamps on every
    ///      attestation structure to prove on-TPM origin. From TPM 2.0
    ///      spec §6.9.
    uint32 internal constant TPM_GENERATED_VALUE = 0xFF544347;

    /// @dev TPM_ST_ATTEST_QUOTE — structure tag for TPMS_ATTEST issued
    ///      by TPM2_Quote. From TPM 2.0 spec §6.9 and §18.4.
    uint16 internal constant TPM_ST_ATTEST_QUOTE = 0x8018;

    /// @dev Minimum length of a TPMS_ATTEST blob this library can parse:
    ///      magic (4) + type (2) + qualifiedSigner.size (2) + min name (0)
    ///      + extraData.size (2) + min data (0) = 10 bytes. Any shorter
    ///      blob cannot contain even an empty quote.
    uint256 internal constant MIN_ATTEST_LEN = 10;

    /// @dev Maximum length of a TPMS_ATTEST blob the parser accepts.
    ///      Mirrors the bundle's maxTPMAttestSize cap in
    ///      story-kernel/enclave/tdx/platform/bundle.go (4 KiB).
    uint256 internal constant MAX_ATTEST_LEN = 4 * 1024;

    /// @notice Extracts qualifyingData (TPMS_ATTEST.extraData) from a
    ///         TPM2_Quote attestation blob.
    /// @param attest The TPMS_ATTEST blob (no TPM2B_ATTEST size prefix).
    /// @return data The extraData bytes — the on-chain caller compares
    ///         this against the expected commitment.
    /// @dev Reverts with descriptive messages on any structural defect:
    ///      wrong magic, wrong type, oversize fields, or truncation.
    function extractQualifyingData(bytes memory attest) internal pure returns (bytes memory data) {
        require(attest.length >= MIN_ATTEST_LEN, "TPM2Attest: too short");
        require(attest.length <= MAX_ATTEST_LEN, "TPM2Attest: too long");

        // magic (4 bytes BE) at offset 0.
        uint32 magic = (uint32(uint8(attest[0])) << 24) | (uint32(uint8(attest[1])) << 16)
            | (uint32(uint8(attest[2])) << 8) | uint32(uint8(attest[3]));
        require(magic == TPM_GENERATED_VALUE, "TPM2Attest: bad magic");

        // type (2 bytes BE) at offset 4.
        uint16 attestType = (uint16(uint8(attest[4])) << 8) | uint16(uint8(attest[5]));
        require(attestType == TPM_ST_ATTEST_QUOTE, "TPM2Attest: not a quote");

        // qualifiedSigner is TPM2B_NAME at offset 6: uint16 size + name bytes.
        uint256 nameLenOff = 6;
        require(attest.length >= nameLenOff + 2, "TPM2Attest: truncated at signer len");
        uint256 nameLen = (uint256(uint8(attest[nameLenOff])) << 8) | uint256(uint8(attest[nameLenOff + 1]));
        // TPM2 names are bounded by the largest hash output (SHA-512 = 64
        // bytes) plus a 2-byte algorithm prefix; cap conservatively at
        // 80 bytes so we reject obvious overruns before slicing.
        require(nameLen <= 80, "TPM2Attest: signer name too long");

        uint256 dataLenOff = nameLenOff + 2 + nameLen;
        require(attest.length >= dataLenOff + 2, "TPM2Attest: truncated at data len");
        uint256 dataLen = (uint256(uint8(attest[dataLenOff])) << 8) | uint256(uint8(attest[dataLenOff + 1]));
        // TPM2B_DATA is bounded by the TPM's TPM2B_DIGEST size (typically
        // 64 bytes; our kernel uses up to that). 256 is a safe cap that
        // covers every implementation we have observed and matches the
        // kernel-side tpmQualifyingDataMax pre-check.
        require(dataLen <= 256, "TPM2Attest: extraData too long");

        uint256 dataStart = dataLenOff + 2;
        require(attest.length >= dataStart + dataLen, "TPM2Attest: truncated extraData");

        data = new bytes(dataLen);
        for (uint256 i = 0; i < dataLen; i++) {
            data[i] = attest[dataStart + i];
        }
    }
}
