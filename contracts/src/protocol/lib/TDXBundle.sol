// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

/// @title TDXBundle
/// @notice Parser for the Path-B "STBN" TDX attestation bundle.
///
/// @dev Mirrors the wire-format spec defined in
///      story-kernel/enclave/tdx/platform/bundle.go. The kernel emits
///      this bundle on vendors where V4.report_data is not under guest
///      control (notably Azure CVM TDX, where the OpenHCL paravisor
///      locks report_data at boot to hash(VariableData)). The bundle
///      composes a TDX V4 quote with a TPM2_Quote signed by the AK so
///      the on-chain hook can bind arbitrary user_data via the TPM
///      qualifyingData channel.
///
///      Wire layout (locked, big-endian length fields):
///
///        offset            size  field
///        ------            ----  -----
///        0                 4     magic = "STBN" (0x5354424E)
///        4                 1     version = 0x01
///        5                 1     flags  (bit0 TPM_PRESENT; bit1..7 MUST be 0)
///        6                 2     vendor_tag (0x0000 direct, 0x0001 azure, 0xFFFF test)
///        8                 4     tdx_v4_len (BE) = N
///        12                N     tdx_v4 bytes
///        12+N              4     tpm_attest_len (BE) = M
///        16+N              M     tpm_attest (TPMS_ATTEST blob, no TPM2B prefix)
///        16+N+M            4     tpm_sig_len (BE) = K
///        20+N+M            K     tpm_sig (TPMT_SIGNATURE: algId + hashAlg + sig)
///        20+N+M+K          4     ak_pub_len (BE) = L
///        24+N+M+K          L     ak_pub DER SubjectPublicKeyInfo (RSA-2048)
///        24+N+M+K+L        4     runtime_data_len (BE) = R
///        28+N+M+K+L        R     runtime_data (Azure: VariableData JSON; direct: empty)
///
///      Length caps mirror the kernel-side enforcement so on-chain DoS
///      surfaces are bounded.
library TDXBundle {
    /*//////////////////////////////////////////////////////////////////////////
    //                              Wire constants                            //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev "STBN" — 4-byte magic at offset 0. We compare via word load
    ///      and bit-mask, so the constant is encoded as the high four
    ///      bytes of a bytes32.
    bytes4 internal constant BUNDLE_MAGIC = 0x5354424E;

    /// @dev Currently the only accepted bundle version. Forward-
    ///      compatible slots (post-quantum AKs, extra evidence) bump
    ///      this byte.
    uint8 internal constant BUNDLE_VERSION = 0x01;

    /// @dev TPM_PRESENT — bit0 of the flags byte. The current bundle
    ///      always has this set; bundles without TPM evidence are
    ///      reserved for a future bare-metal-no-TPM vendor.
    uint8 internal constant FLAG_TPM_PRESENT = 0x01;

    /// @dev Mask of reserved flag bits. Any bit set here triggers a
    ///      reject — defense against silent feature smuggling.
    uint8 internal constant FLAGS_RESERVED_MASK = 0xFE;

    /// @dev Vendor-tag values. Used for vendor-aware AK binding (direct
    ///      vendor binds via SHA256(AK_pub); Azure binds via
    ///      SHA256(runtime_data)). On-chain trust does NOT depend on
    ///      these values — they merely select which binding equation
    ///      the verifier evaluates.
    uint16 internal constant VENDOR_TAG_DIRECT = 0x0000;
    uint16 internal constant VENDOR_TAG_AZURE = 0x0001;
    uint16 internal constant VENDOR_TAG_TEST = 0xFFFF;

    /// @dev Fixed prefix size up to and including the tdx_v4_len field.
    ///      Inputs shorter than this fail fast before allocation.
    uint256 internal constant BUNDLE_HEADER_SIZE = 12;

    /*//////////////////////////////////////////////////////////////////////////
    //                              Length caps                               //
    //                  (mirror story-kernel/enclave/tdx/platform/bundle.go)  //
    //////////////////////////////////////////////////////////////////////////*/

    uint256 internal constant MAX_BUNDLE_SIZE = 64 * 1024;
    uint256 internal constant MAX_TDX_V4_SIZE = 16 * 1024;
    uint256 internal constant MAX_TPM_ATTEST_SIZE = 4 * 1024;
    uint256 internal constant MAX_TPM_SIG_SIZE = 1 * 1024;
    uint256 internal constant MAX_AK_PUB_SIZE = 2 * 1024;
    uint256 internal constant MAX_RUNTIME_DATA_SIZE = 4 * 1024;

    /*//////////////////////////////////////////////////////////////////////////
    //                              Parsed view                               //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Parsed bundle. `tdxV4`, `tpmAttest`, `tpmSig`, `akPub`,
    ///      `runtimeData` are fresh memory copies — calldata-aliased
    ///      slices would be cheaper but Solidity does not support
    ///      taking a `bytes calldata` slice into a `bytes memory`
    ///      without copy when crossing the function boundary, so the
    ///      tradeoff is paid up front for ergonomic library use.
    struct ParsedBundle {
        uint8 version;
        uint8 flags;
        uint16 vendorTag;
        bytes tdxV4;
        bytes tpmAttest;
        bytes tpmSig;
        bytes akPub;
        bytes runtimeData;
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              Public API                                //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Returns true iff `buf` begins with the "STBN" magic.
    /// @dev Used by the on-chain hook to dispatch raw-V4 vs bundle
    ///      without a full parse. Constant-time at the word level.
    function hasMagic(bytes calldata buf) internal pure returns (bool) {
        if (buf.length < 4) return false;
        return bytes4(buf[:4]) == BUNDLE_MAGIC;
    }

    /// @notice Parses an STBN-prefixed bundle.
    /// @param input The bundle bytes (calldata).
    /// @return parsed Decoded fields. Reverts with descriptive errors
    ///         on any malformation.
    /// @dev This function is fail-closed: every length-cap violation,
    ///      reserved-flag bit, or unknown vendor_tag aborts before any
    ///      cryptographic verification is attempted. The kernel's
    ///      MarshalBundle enforces the same caps so a bundle the
    ///      kernel produced will round-trip cleanly.
    function parse(bytes calldata input) internal pure returns (ParsedBundle memory parsed) {
        require(input.length <= MAX_BUNDLE_SIZE, "TDXBundle: bundle too large");
        require(input.length >= BUNDLE_HEADER_SIZE, "TDXBundle: bundle too short");

        // Magic.
        require(bytes4(input[:4]) == BUNDLE_MAGIC, "TDXBundle: bad magic");

        // Version.
        uint8 version = uint8(input[4]);
        require(version == BUNDLE_VERSION, "TDXBundle: unsupported version");

        // Flags. Reserved bits MUST be zero. TPM_PRESENT MUST be set —
        // we do not yet support the bare-metal-no-TPM vendor; admitting
        // that path here would let an attacker omit the TPM2_Quote and
        // bypass qualifyingData binding.
        uint8 flags = uint8(input[5]);
        require(flags & FLAGS_RESERVED_MASK == 0, "TDXBundle: reserved flag bits set");
        require(flags & FLAG_TPM_PRESENT != 0, "TDXBundle: missing TPM_PRESENT flag");

        // Vendor tag.
        uint16 vendorTag = (uint16(uint8(input[6])) << 8) | uint16(uint8(input[7]));
        require(
            vendorTag == VENDOR_TAG_DIRECT || vendorTag == VENDOR_TAG_AZURE || vendorTag == VENDOR_TAG_TEST,
            "TDXBundle: unknown vendor tag"
        );

        // tdx_v4 section.
        uint256 tdxLen = _readBE32(input, 8);
        require(tdxLen <= MAX_TDX_V4_SIZE, "TDXBundle: tdx_v4 too large");
        uint256 off = BUNDLE_HEADER_SIZE;
        require(off + tdxLen + 4 <= input.length, "TDXBundle: truncated at attest_len");
        bytes memory tdxV4 = _copy(input, off, tdxLen);
        off += tdxLen;

        // tpm_attest section.
        uint256 attestLen = _readBE32(input, off);
        require(attestLen <= MAX_TPM_ATTEST_SIZE, "TDXBundle: tpm_attest too large");
        off += 4;
        require(off + attestLen + 4 <= input.length, "TDXBundle: truncated at sig_len");
        bytes memory tpmAttest = _copy(input, off, attestLen);
        off += attestLen;

        // tpm_sig section.
        uint256 sigLen = _readBE32(input, off);
        require(sigLen <= MAX_TPM_SIG_SIZE, "TDXBundle: tpm_sig too large");
        off += 4;
        require(off + sigLen + 4 <= input.length, "TDXBundle: truncated at ak_pub_len");
        bytes memory tpmSig = _copy(input, off, sigLen);
        off += sigLen;

        // ak_pub section.
        uint256 akLen = _readBE32(input, off);
        require(akLen <= MAX_AK_PUB_SIZE, "TDXBundle: ak_pub too large");
        off += 4;
        require(off + akLen + 4 <= input.length, "TDXBundle: truncated at runtime_data_len");
        bytes memory akPub = _copy(input, off, akLen);
        off += akLen;

        // runtime_data section.
        uint256 rtLen = _readBE32(input, off);
        require(rtLen <= MAX_RUNTIME_DATA_SIZE, "TDXBundle: runtime_data too large");
        off += 4;
        require(off + rtLen <= input.length, "TDXBundle: truncated runtime_data");
        bytes memory runtimeData = _copy(input, off, rtLen);
        off += rtLen;

        // Trailing bytes are NOT silently tolerated — any leftover
        // indicates malformation and we reject. Mirrors the kernel-side
        // exact-fit check in UnmarshalBundle.
        require(off == input.length, "TDXBundle: trailing bytes");

        // Vendor-specific RuntimeData invariants. Mirrors MarshalBundle
        // / UnmarshalBundle in the kernel.
        if (vendorTag == VENDOR_TAG_DIRECT) {
            require(rtLen == 0, "TDXBundle: direct must omit runtime_data");
        } else if (vendorTag == VENDOR_TAG_AZURE) {
            require(rtLen != 0, "TDXBundle: azure missing runtime_data");
        }

        parsed = ParsedBundle({
            version: version,
            flags: flags,
            vendorTag: vendorTag,
            tdxV4: tdxV4,
            tpmAttest: tpmAttest,
            tpmSig: tpmSig,
            akPub: akPub,
            runtimeData: runtimeData
        });
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                          Internal helpers                              //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Reads a big-endian uint32 from `input` at `offset`.
    ///      Caller must have already bounds-checked offset + 4.
    function _readBE32(bytes calldata input, uint256 offset) private pure returns (uint256) {
        return (uint256(uint8(input[offset])) << 24) | (uint256(uint8(input[offset + 1])) << 16)
            | (uint256(uint8(input[offset + 2])) << 8) | uint256(uint8(input[offset + 3]));
    }

    /// @dev Copies `len` bytes from `input` at `offset` into a fresh
    ///      memory buffer. Caller must have bounds-checked.
    function _copy(bytes calldata input, uint256 offset, uint256 len) private pure returns (bytes memory out) {
        out = new bytes(len);
        for (uint256 i = 0; i < len; i++) {
            out[i] = input[offset + i];
        }
    }
}
