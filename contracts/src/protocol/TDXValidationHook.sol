// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { PausableUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { UUPSUpgradeable } from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

import { ITDXValidationHook } from "../interfaces/ITDXValidationHook.sol";
import { IAutomataDcapAttestationFee } from "../interfaces/external/IAutomataDcapAttestationFee.sol";
import { TDXBundle } from "./lib/TDXBundle.sol";
import { TPM2Attest } from "./lib/TPM2Attest.sol";
import { RSASSAVerify } from "./lib/RSASSAVerify.sol";

/// @title TDXValidationHook
/// @notice On-chain validator for Intel TDX attestation quotes used by the DKG flow.
/// @dev Mirrors SGXValidationHook in structure (UUPS proxy, owner-only admin, pausable),
///      and replaces SGX-specific quote field offsets with the TDX V4/V5 layout.
///
///      Identity model:
///      - The kernel's TDX backend treats CodeCommitment as the native concatenation
///        MRTD || RTMR0 || RTMR1 || RTMR2 || RTMR3 (240 bytes). On-chain we cannot store
///        240 bytes as a bytes32, so this hook compresses by computing
///        keccak256(MRTD || RTMR0 || RTMR1 || RTMR2 || RTMR3) and compares against the
///        whitelisted bytes32 value in DKG.enclaveTypeData[enclaveType].codeCommitment.
///        Operators must whitelist the matching keccak digest.
///
///      Quote layout assumptions (Intel TDX DCAP V4 and V5):
///      - 48-byte header. byte[0..1] little-endian uint16 version, byte[4..7] little-endian
///        uint32 tee_type. TDX tee_type == 0x00000081 (wire bytes 0x81 0x00 0x00 0x00).
///      - V4 body is 584 bytes; V5 body is 648 bytes (V5 adds a trailing TEE_TCB_SVN_2 field).
///      - Every measurement field this contract reads (MRTD, RTMR0..3, ReportData) sits at
///        the same absolute byte offset in V4 and V5, so the offset constants are
///        version-independent. Only the quote-length minimum varies per version.
contract TDXValidationHook is ITDXValidationHook, Ownable2StepUpgradeable, PausableUpgradeable, UUPSUpgradeable {
    /// @dev Storage structure for the TDXValidationHook
    /// @param automataValidationAddr The address of the automata validation contract
    /// @param tcbEvaluationDataNumber The tcb evaluation data number
    /// @custom:storage-location erc7201:story.TDXValidationHook
    struct TDXValidationHookStorage {
        address automataValidationAddr;
        uint32 tcbEvaluationDataNumber;
    }

    address public immutable DKG;

    // keccak256(abi.encode(uint256(keccak256("story.TDXValidationHook")) - 1)) & ~bytes32(uint256(0xff));
    bytes32 private constant TDXValidationHookStorageLocation =
        0xc82cee3b70f9b4f764f18b66d0052cef718142a50149d16a1c68ba0074816f00;

    /*//////////////////////////////////////////////////////////////////////////
    //                          Quote layout constants                        //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev TDX quote header size (V4 and V5 share the same header layout).
    uint256 private constant QUOTE_HEADER_SIZE = 48;

    /// @dev Per-measurement field size: MRTD, each RTMR, MRSEAM, etc. are all 48 bytes.
    uint256 private constant MEASUREMENT_SIZE = 48;

    /// @dev Number of contiguous RTMRs we hash (RTMR0..3).
    uint256 private constant RTMR_COUNT = 4;

    /// @dev Body length for TD10 (V4) quotes. Total min quote size = header + body.
    uint256 private constant V4_BODY_SIZE = 584;
    /// @dev Body length for TD15 (V5) quotes. V5 appends a 16-byte TEE_TCB_SVN_2 field.
    uint256 private constant V5_BODY_SIZE = 648;

    /// @dev Minimum total quote length for V4 (header + body, exclusive of auth_data trailer).
    uint256 private constant MIN_V4_QUOTE_SIZE = QUOTE_HEADER_SIZE + V4_BODY_SIZE; // 632
    /// @dev Minimum total quote length for V5 (header + body, exclusive of auth_data trailer).
    uint256 private constant MIN_V5_QUOTE_SIZE = QUOTE_HEADER_SIZE + V5_BODY_SIZE; // 696

    /// @dev TDX TEE-type discriminator. Wire bytes at offset 4 are [0x81, 0x00, 0x00, 0x00]
    ///      (little-endian uint32 = 0x00000081). SGX has tee_type == 0; we explicitly
    ///      reject any non-TDX value before extracting at TDX offsets.
    uint8 private constant TEE_TYPE_TDX_BYTE0 = 0x81;

    /// @dev Absolute offsets into the raw quote (header + body). Identical for V4 and V5.
    ///      MRTD: body offset 136 + 48-byte header = 184.
    uint256 private constant OFFSET_MRTD = 184;
    /// @dev RTMR0 absolute offset. RTMR0..3 are contiguous (4 * 48 = 192 bytes ending at 568).
    uint256 private constant OFFSET_RTMR0 = 376;
    /// @dev REPORT_DATA absolute offset. Total length 64 bytes; we compare only the first 32.
    uint256 private constant OFFSET_REPORT_DATA = 568;

    constructor(address dkg) {
        require(dkg != address(0), "TDXValidationHook: DKG cannot be empty");
        DKG = dkg;
        _disableInitializers();
    }

    /// @notice Initializes the contract
    /// @param owner The address of the owner of the contract
    /// @param automataValidationAddr The address of the automata validation contract
    /// @param tcbEvaluationDataNumber The TCB evaluation data number
    function initialize(
        address owner,
        address automataValidationAddr,
        uint32 tcbEvaluationDataNumber
    ) external initializer {
        __Ownable_init(owner);
        __Pausable_init();
        __UUPSUpgradeable_init();

        _setAutomataValidationAddr(automataValidationAddr);
        _setTcbEvaluationDataNumber(tcbEvaluationDataNumber);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Admin Setters                              //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Sets the address of the automata validation contract
    /// @param newAutomataValidationAddr The address of the automata validation contract
    function setAutomataValidationAddr(address newAutomataValidationAddr) external onlyOwner {
        _setAutomataValidationAddr(newAutomataValidationAddr);
    }

    /// @notice Sets the TCB evaluation data number
    /// @param newTcbEvaluationDataNumber The TCB evaluation data number
    function setTcbEvaluationDataNumber(uint32 newTcbEvaluationDataNumber) external onlyOwner {
        _setTcbEvaluationDataNumber(newTcbEvaluationDataNumber);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Authentication Logic                         //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Validates a TDX attestation quote.
    ///
    ///      Dispatches on the first four bytes of `enclaveReport`:
    ///      - "STBN" (0x5354424E) → Path-B bundle (paravisor-mediated
    ///        TDX guests and any future bundle-emitting vendor). The
    ///        bundle composes a V4 quote with a TPM2_Quote signed by
    ///        the AK; binding the kernel-controlled qualifyingData
    ///        against `expectedDataCommitment` is what authorizes the
    ///        registration.
    ///      - any other prefix → legacy raw V4/V5 path (direct vendor on
    ///        GCP/bare-metal). V4.report_data is under guest control so
    ///        `expectedDataCommitment` is read directly from there,
    ///        mirroring SGX semantics.
    ///
    ///      Both paths converge on the same return shape and same
    ///      require-fail semantics so callers (DKG.register) need not
    ///      know which vendor produced the quote.
    /// @param expectedCodeCommitment keccak256(MRTD || RTMR0 || RTMR1 || RTMR2 || RTMR3)
    /// @param expectedDataCommitment First 32 bytes of REPORT_DATA (raw V4 path) or
    ///        TPMS_ATTEST.qualifyingData (bundle path), set by the kernel to
    ///        keccak256(validatorAddr || round || startBlockHeight || startBlockHash ||
    ///                  dkgPubKey || enclaveCommKey).
    /// @param enclaveReport Raw TDX V4/V5 quote OR an STBN-prefixed bundle.
    /// @param validationContext Reserved for future use; not consumed here.
    function validateReport(
        bytes32 expectedCodeCommitment,
        bytes32 expectedDataCommitment,
        bytes calldata enclaveReport,
        bytes calldata validationContext
    ) external override returns (bool) {
        require(msg.sender == DKG, "TDXValidationHook: Only DKG can call this function");
        require(enclaveReport.length > 0, "TDXValidationHook: Empty enclave report");
        require(expectedCodeCommitment != bytes32(0), "TDXValidationHook: Zero code commitment");
        require(expectedDataCommitment != bytes32(0), "TDXValidationHook: Zero data commitment");

        // Dispatch on the STBN magic. Bundle path piggy-backs on the
        // direct path's DCAP verification, code-commitment extraction,
        // and Automata trust anchor — the only added work is parsing
        // the bundle envelope, verifying the TPM2_Quote signature, and
        // checking AK-binding and qualifyingData.
        if (TDXBundle.hasMagic(enclaveReport)) {
            _validateBundle(expectedCodeCommitment, expectedDataCommitment, enclaveReport);
            return true;
        }

        // Defense-in-depth: assert quote header shape and TEE type BEFORE invoking
        // Automata. This prevents an attacker from supplying a (cryptographically valid)
        // SGX quote against this hook to confuse the offset-based extraction below.
        // Automata's verifyAndAttestOnChain dispatches on tee_type internally, but the
        // explicit gate keeps this contract self-contained against future Automata
        // contract changes that might soften that dispatch.
        require(enclaveReport.length >= QUOTE_HEADER_SIZE, "TDXValidationHook: Quote too short for header");
        require(
            uint8(enclaveReport[4]) == TEE_TYPE_TDX_BYTE0 &&
                uint8(enclaveReport[5]) == 0 &&
                uint8(enclaveReport[6]) == 0 &&
                uint8(enclaveReport[7]) == 0,
            "TDXValidationHook: Not a TDX quote"
        );

        // Body-size minimum varies per version. Read version as little-endian uint16
        // from offset 0..1.
        uint16 version = uint16(uint8(enclaveReport[0])) | (uint16(uint8(enclaveReport[1])) << 8);
        if (version == 4) {
            require(enclaveReport.length >= MIN_V4_QUOTE_SIZE, "TDXValidationHook: Quote too short for V4 body");
        } else if (version == 5) {
            require(enclaveReport.length >= MIN_V5_QUOTE_SIZE, "TDXValidationHook: Quote too short for V5 body");
        } else {
            revert("TDXValidationHook: Unsupported quote version");
        }

        // Cryptographic chain of trust: Automata DCAP verifies signatures and TCB.
        // Automata routes V4 and V5 internally based on header.version + tee_type.
        TDXValidationHookStorage storage $ = _getTDXValidationHookStorage();
        (bool success, ) = IAutomataDcapAttestationFee($.automataValidationAddr).verifyAndAttestOnChain(
            enclaveReport,
            $.tcbEvaluationDataNumber
        );
        require(success, "TDXValidationHook: Attestation failed");

        // Identity match: hash MRTD || RTMR0..3 and compare against the whitelisted digest.
        require(
            _extractReportCodeCommitment(enclaveReport) == expectedCodeCommitment,
            "TDXValidationHook: Code commitment does not match"
        );

        // Instance-data match: first 32 bytes of REPORT_DATA must equal the kernel-bound
        // attestation payload computed by DKG.register from EnclaveInstanceData.
        require(
            _extractReportInstanceDataCommitment(enclaveReport) == expectedDataCommitment,
            "TDXValidationHook: Data commitment does not match"
        );

        return true;
    }

    /// @dev Bundle-path validation ladder. See the per-step comments for
    ///      the security invariant each check enforces. Total estimated
    ///      gas ≤ ~6 M (modexp ~1.6k under EIP-2565, plus DCAP verify
    ///      which dominates and is identical to the raw-V4 path).
    ///
    ///      Steps:
    ///        1. Parse bundle (length caps, magic, version, flags, vendor_tag).
    ///        2. Validate inner V4: TEE-type byte and quote-version length.
    ///        3. DCAP verify the inner V4 (same trust anchor as raw path).
    ///        4. Code commitment: keccak256(MRTD||RTMR0..3) == expected.
    ///        5. AK binding (vendor-aware): direct binds SHA256(AK_pub)
    ///           into V4.report_data[0:32]; paravisor binds
    ///           SHA256(runtime_data) there; report_data[32:64] MUST be
    ///           zero either way.
    ///        6. RSASSA-PKCS#1 v1.5 verify TPM2_Quote signature with AK pub.
    ///        7. Extract qualifyingData (TPMS_ATTEST.extraData).
    ///        8. qualifyingData == expectedDataCommitment (32-byte equality).
    function _validateBundle(
        bytes32 expectedCodeCommitment,
        bytes32 expectedDataCommitment,
        bytes calldata enclaveReport
    ) internal {
        // Step 1: parse the bundle envelope.
        TDXBundle.ParsedBundle memory bundle = TDXBundle.parse(enclaveReport);

        // Step 2: validate the inner V4 has TDX tee_type and a known
        // version length. Mirrors the raw-V4 path's gate.
        require(bundle.tdxV4.length >= QUOTE_HEADER_SIZE, "TDXValidationHook: Quote too short for header");
        require(
            uint8(bundle.tdxV4[4]) == TEE_TYPE_TDX_BYTE0 &&
                uint8(bundle.tdxV4[5]) == 0 &&
                uint8(bundle.tdxV4[6]) == 0 &&
                uint8(bundle.tdxV4[7]) == 0,
            "TDXValidationHook: Not a TDX quote"
        );
        uint16 version = uint16(uint8(bundle.tdxV4[0])) | (uint16(uint8(bundle.tdxV4[1])) << 8);
        if (version == 4) {
            require(bundle.tdxV4.length >= MIN_V4_QUOTE_SIZE, "TDXValidationHook: Quote too short for V4 body");
        } else if (version == 5) {
            require(bundle.tdxV4.length >= MIN_V5_QUOTE_SIZE, "TDXValidationHook: Quote too short for V5 body");
        } else {
            revert("TDXValidationHook: Unsupported quote version");
        }

        // Step 3: DCAP verify the inner V4. Identical trust anchor as the
        // raw-V4 path; Automata routes by tee_type + version internally.
        TDXValidationHookStorage storage $ = _getTDXValidationHookStorage();
        (bool success, ) = IAutomataDcapAttestationFee($.automataValidationAddr).verifyAndAttestOnChain(
            bundle.tdxV4,
            $.tcbEvaluationDataNumber
        );
        require(success, "TDXValidationHook: Attestation failed");

        // Step 4: code commitment match. Reuses the same MRTD/RTMR
        // extraction as the raw-V4 path against the inner V4.
        require(
            _extractCodeCommitmentMemory(bundle.tdxV4) == expectedCodeCommitment,
            "TDXValidationHook: Code commitment does not match"
        );

        // Step 5: vendor-aware AK binding plus reserved-zero check on
        // V4.report_data[32:64]. The kernel-side mirror is in
        // story-kernel/enclave/tdx/selfcheck.go verifyBundleAKBinding.
        bytes32 reportDataLeft = _extractReportDataLeft(bundle.tdxV4);
        bytes32 reportDataRight = _extractReportDataRight(bundle.tdxV4);
        require(reportDataRight == bytes32(0), "TDXValidationHook: report_data[32:64] must be zero");

        if (bundle.vendorTag == TDXBundle.VENDOR_TAG_DIRECT) {
            // Direct vendor: report_data[0:32] == sha256(AK_pub_DER).
            require(sha256(bundle.akPub) == reportDataLeft, "TDXValidationHook: AK binding (direct) mismatch");
        } else if (
            bundle.vendorTag == TDXBundle.VENDOR_TAG_PARAVISOR || bundle.vendorTag == TDXBundle.VENDOR_TAG_TEST
        ) {
            // Paravisor-mediated (or test): report_data[0:32] == sha256(runtime_data).
            // Trust model A: we intentionally do NOT parse the JWK on
            // chain — the kernel selfcheck enforces HCLAkPub.n match
            // locally. The chain trusts that the AK in the bundle is
            // the one whose pub-hash sits inside runtime_data because
            // any divergence would be caught by the selfcheck refusing
            // to start.
            require(sha256(bundle.runtimeData) == reportDataLeft, "TDXValidationHook: AK binding (paravisor) mismatch");
        } else {
            // Defense in depth: TDXBundle.parse already rejects unknown
            // vendor tags, so this branch is unreachable in well-formed
            // input. Explicit revert keeps the dispatch exhaustive against
            // future vendor tag additions.
            revert("TDXValidationHook: unsupported vendor tag");
        }

        // Step 6: parse TPMT_SIGNATURE, pin algorithm + hash, RSASSA-
        // PKCS#1 v1.5 verify over sha256(tpm_attest).
        _verifyTPMSignature(bundle.tpmSig, bundle.tpmAttest, bundle.akPub);

        // Step 7+8: extract qualifyingData and compare against expected.
        bytes memory qualifyingData = TPM2Attest.extractQualifyingData(bundle.tpmAttest);
        require(qualifyingData.length == 32, "TDXValidationHook: qualifyingData not 32 bytes");
        require(
            _bytes32From(qualifyingData) == expectedDataCommitment,
            "TDXValidationHook: Data commitment does not match"
        );
    }

    /// @dev Parses TPMT_SIGNATURE and verifies the RSASSA-PKCS#1 v1.5
    ///      signature over sha256(tpmAttest) using akPubDER.
    ///      Pinned alg/hash so a downgrade to PSS or SHA-1 is rejected.
    function _verifyTPMSignature(bytes memory tpmSig, bytes memory tpmAttest, bytes memory akPubDER) internal view {
        // TPMT_SIGNATURE wire layout for RSASSA:
        //   algId       2 bytes BE  = 0x0014 (TPM_ALG_RSASSA)
        //   hashAlg     2 bytes BE  = 0x000B (TPM_ALG_SHA256)
        //   sigLen      2 bytes BE  = 256 (RSA-2048)
        //   sigBytes    256 bytes
        // Total = 262 bytes.
        require(tpmSig.length == 262, "TDXValidationHook: bad TPMT_SIGNATURE length");
        uint16 algId = (uint16(uint8(tpmSig[0])) << 8) | uint16(uint8(tpmSig[1]));
        require(algId == 0x0014, "TDXValidationHook: TPM sig alg != RSASSA");
        uint16 hashAlg = (uint16(uint8(tpmSig[2])) << 8) | uint16(uint8(tpmSig[3]));
        require(hashAlg == 0x000B, "TDXValidationHook: TPM hash alg != SHA-256");
        uint16 sigLen = (uint16(uint8(tpmSig[4])) << 8) | uint16(uint8(tpmSig[5]));
        require(sigLen == 256, "TDXValidationHook: TPM sig length != 256");

        // Slice the 256-byte signature blob.
        bytes memory sig = new bytes(256);
        for (uint256 i = 0; i < 256; i++) {
            sig[i] = tpmSig[6 + i];
        }

        bytes memory modulus = RSASSAVerify.extractRSA2048Modulus(akPubDER);
        bytes32 digest = sha256(tpmAttest);
        require(RSASSAVerify.verify(digest, sig, modulus), "TDXValidationHook: TPM sig verify failed");
    }

    /// @dev Memory-version of MRTD||RTMR0..3 keccak. Mirrors the
    ///      calldata version used on the raw-V4 path, but operates on
    ///      `bytes memory` because the inner V4 inside a bundle is a
    ///      copied buffer rather than a calldata slice.
    function _extractCodeCommitmentMemory(bytes memory v4) internal pure returns (bytes32) {
        bytes memory ident = new bytes(MEASUREMENT_SIZE + RTMR_COUNT * MEASUREMENT_SIZE);
        for (uint256 i = 0; i < MEASUREMENT_SIZE; i++) {
            ident[i] = v4[OFFSET_MRTD + i];
        }
        for (uint256 i = 0; i < RTMR_COUNT * MEASUREMENT_SIZE; i++) {
            ident[MEASUREMENT_SIZE + i] = v4[OFFSET_RTMR0 + i];
        }
        return keccak256(ident);
    }

    /// @dev Reads V4.report_data[0:32] from a memory-buffered V4 quote.
    function _extractReportDataLeft(bytes memory v4) internal pure returns (bytes32 r) {
        // bytes memory: skip 32-byte length prefix, then offset.
        assembly {
            r := mload(add(add(v4, 32), OFFSET_REPORT_DATA))
        }
    }

    /// @dev Reads V4.report_data[32:64] from a memory-buffered V4 quote.
    function _extractReportDataRight(bytes memory v4) internal pure returns (bytes32 r) {
        assembly {
            r := mload(add(add(v4, 32), add(OFFSET_REPORT_DATA, 32)))
        }
    }

    /// @dev Loads a 32-byte value from `data[0:32]`. Caller has bounds-
    ///      checked length.
    function _bytes32From(bytes memory data) internal pure returns (bytes32 r) {
        assembly {
            r := mload(add(data, 32))
        }
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              Get Functions                             //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Gets the address of the automata validation contract
    /// @return The address of the automata validation contract
    function automataValidationAddr() external view returns (address) {
        return _getTDXValidationHookStorage().automataValidationAddr;
    }

    /// @notice Gets the TCB evaluation data number
    /// @return The TCB evaluation data number
    function tcbEvaluationDataNumber() external view returns (uint32) {
        return _getTDXValidationHookStorage().tcbEvaluationDataNumber;
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Internal Functions                           //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Sets the address of the automata validation contract
    /// @param newAutomataValidationAddr The address of the automata validation contract
    function _setAutomataValidationAddr(address newAutomataValidationAddr) internal {
        require(newAutomataValidationAddr != address(0), "TDXValidationHook: Automata Validation cannot be empty");
        _getTDXValidationHookStorage().automataValidationAddr = newAutomataValidationAddr;
    }

    /// @notice Sets the TCB evaluation data number
    /// @param newTcbEvaluationDataNumber The TCB evaluation data number
    function _setTcbEvaluationDataNumber(uint32 newTcbEvaluationDataNumber) internal {
        _getTDXValidationHookStorage().tcbEvaluationDataNumber = newTcbEvaluationDataNumber;
    }

    /// @dev Extracts the code commitment from a raw TDX quote.
    ///      Computes keccak256(MRTD || RTMR0 || RTMR1 || RTMR2 || RTMR3).
    ///      RTMR0..3 are contiguous in the body (4 * 48 = 192 bytes), so we hash
    ///      across two calldata slices: [MRTD] and [RTMR0..RTMR3].
    /// @param enclaveReport The raw TDX quote (header + body + auth_data)
    /// @return The compressed code commitment (32 bytes)
    function _extractReportCodeCommitment(bytes calldata enclaveReport) internal pure returns (bytes32) {
        return
            keccak256(
                abi.encodePacked(
                    enclaveReport[OFFSET_MRTD:OFFSET_MRTD + MEASUREMENT_SIZE],
                    enclaveReport[OFFSET_RTMR0:OFFSET_RTMR0 + RTMR_COUNT * MEASUREMENT_SIZE]
                )
            );
    }

    /// @dev Extracts the first 32 bytes of REPORT_DATA from a raw TDX quote.
    ///      REPORT_DATA is 64 bytes; the kernel populates the first 32 bytes with a
    ///      keccak256 commitment over the registration payload. The remaining 32 bytes
    ///      are zero padding and are not consumed on-chain.
    /// @param enclaveReport The raw TDX quote (header + body + auth_data)
    /// @return result The first 32 bytes of REPORT_DATA as bytes32
    function _extractReportInstanceDataCommitment(bytes calldata enclaveReport) internal pure returns (bytes32 result) {
        // calldataload reads exactly 32 bytes; safe because validateReport asserts
        // enclaveReport.length >= MIN_V4_QUOTE_SIZE (= 632 > OFFSET_REPORT_DATA + 32).
        assembly {
            result := calldataload(add(enclaveReport.offset, OFFSET_REPORT_DATA))
        }
    }

    /// @dev Hook to authorize the upgrade according to UUPSUpgradeable
    /// @param newImplementation The address of the new implementation
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    /// @dev Returns the storage struct of TDXValidationHook.
    function _getTDXValidationHookStorage() private pure returns (TDXValidationHookStorage storage $) {
        assembly {
            $.slot := TDXValidationHookStorageLocation
        }
    }
}
