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
    /*//////////////////////////////////////////////////////////////////////////
    //                              Events                                    //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Emitted when a (MRTD, RTMR0) cloud-platform tuple is approved.
    /// @param key keccak256(mrtd48 || rtmr048)
    /// @param label Free-form governance label
    event CloudPlatformApproved(bytes32 indexed key, string label);

    /// @notice Emitted when a (MRTD, RTMR0) cloud-platform tuple is revoked.
    /// @param key keccak256(mrtd48 || rtmr048)
    event CloudPlatformRevoked(bytes32 indexed key);

    /// @notice Emitted when a (RTMR1, RTMR2) binary-release tuple is approved.
    /// @param key keccak256(rtmr148 || rtmr248)
    /// @param version Free-form governance label
    event BinaryReleaseApproved(bytes32 indexed key, string version);

    /// @notice Emitted when a (RTMR1, RTMR2) binary-release tuple is revoked.
    /// @param key keccak256(rtmr148 || rtmr248)
    event BinaryReleaseRevoked(bytes32 indexed key);

    /// @dev Storage structure for the TDXValidationHook.
    ///      Option D (decomposed whitelist): we split the on-chain identity check into
    ///      two independent tables so cloud-vendor lifecycle (MRTD/RTMR0 changes) and
    ///      Story binary release lifecycle (RTMR1/RTMR2 changes) can be governed
    ///      independently. Governance cost becomes N+M instead of N*M.
    /// @param automataValidationAddr The address of the automata validation contract
    /// @param tcbEvaluationDataNumber The tcb evaluation data number
    /// @param approvedCloudPlatforms keccak256(MRTD || RTMR0) => approved flag
    /// @param approvedBinaryReleases keccak256(RTMR1 || RTMR2) => approved flag
    /// @custom:storage-location erc7201:story.TDXValidationHook
    struct TDXValidationHookStorage {
        address automataValidationAddr;
        uint32 tcbEvaluationDataNumber;
        mapping(bytes32 => bool) approvedCloudPlatforms;
        mapping(bytes32 => bool) approvedBinaryReleases;
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
    /// @dev RTMR1 absolute offset (= OFFSET_RTMR0 + MEASUREMENT_SIZE). Option D binary identity.
    uint256 private constant OFFSET_RTMR1 = 424;
    /// @dev RTMR2 absolute offset (= OFFSET_RTMR0 + 2 * MEASUREMENT_SIZE). Option D binary identity.
    uint256 private constant OFFSET_RTMR2 = 472;
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

    /// @notice Approves a (MRTD, RTMR0) cloud platform tuple (Option D).
    /// @dev On-chain key is keccak256(mrtd48 || rtmr048). Both arguments MUST be the
    ///      raw 48-byte SHA-384 measurements from the quote (NOT pre-hashed).
    function approveCloudPlatform(
        bytes calldata mrtd48,
        bytes calldata rtmr048,
        string calldata label
    ) external override onlyOwner {
        require(mrtd48.length == MEASUREMENT_SIZE, "TDXValidationHook: MRTD must be 48 bytes");
        require(rtmr048.length == MEASUREMENT_SIZE, "TDXValidationHook: RTMR0 must be 48 bytes");
        bytes32 key = keccak256(abi.encodePacked(mrtd48, rtmr048));
        _getTDXValidationHookStorage().approvedCloudPlatforms[key] = true;
        emit CloudPlatformApproved(key, label);
    }

    /// @notice Revokes a previously approved (MRTD, RTMR0) cloud platform tuple.
    function revokeCloudPlatform(bytes calldata mrtd48, bytes calldata rtmr048) external override onlyOwner {
        require(mrtd48.length == MEASUREMENT_SIZE, "TDXValidationHook: MRTD must be 48 bytes");
        require(rtmr048.length == MEASUREMENT_SIZE, "TDXValidationHook: RTMR0 must be 48 bytes");
        bytes32 key = keccak256(abi.encodePacked(mrtd48, rtmr048));
        delete _getTDXValidationHookStorage().approvedCloudPlatforms[key];
        emit CloudPlatformRevoked(key);
    }

    /// @notice Approves a (RTMR1, RTMR2) binary release tuple (Option D).
    /// @dev On-chain key is keccak256(rtmr148 || rtmr248). RTMR1 is the TDVF-anchored
    ///      kernel measurement; RTMR2 is the initrd + cmdline measurement. Together they
    ///      bind the user-space binary IF the binary is baked into the initrd as PID 1
    ///      (TDVF measures kernel/initrd before any user-space code executes).
    function approveBinaryRelease(
        bytes calldata rtmr148,
        bytes calldata rtmr248,
        string calldata version
    ) external override onlyOwner {
        require(rtmr148.length == MEASUREMENT_SIZE, "TDXValidationHook: RTMR1 must be 48 bytes");
        require(rtmr248.length == MEASUREMENT_SIZE, "TDXValidationHook: RTMR2 must be 48 bytes");
        bytes32 key = keccak256(abi.encodePacked(rtmr148, rtmr248));
        _getTDXValidationHookStorage().approvedBinaryReleases[key] = true;
        emit BinaryReleaseApproved(key, version);
    }

    /// @notice Revokes a previously approved (RTMR1, RTMR2) binary release tuple.
    function revokeBinaryRelease(bytes calldata rtmr148, bytes calldata rtmr248) external override onlyOwner {
        require(rtmr148.length == MEASUREMENT_SIZE, "TDXValidationHook: RTMR1 must be 48 bytes");
        require(rtmr248.length == MEASUREMENT_SIZE, "TDXValidationHook: RTMR2 must be 48 bytes");
        bytes32 key = keccak256(abi.encodePacked(rtmr148, rtmr248));
        delete _getTDXValidationHookStorage().approvedBinaryReleases[key];
        emit BinaryReleaseRevoked(key);
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
    /// @param expectedCodeCommitment Retained for ABI compatibility with the DKG caller and
    ///        SGX hook. **No longer consulted** under Option D — code identity is now enforced
    ///        by the two decomposed whitelists `approvedCloudPlatforms` (MRTD, RTMR0) and
    ///        `approvedBinaryReleases` (RTMR1, RTMR2). Callers MUST still pass a non-zero
    ///        value (kept as a sanity-check input gate).
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

        // Option D — decomposed identity match.
        //
        // (1) Cloud platform: keccak256(MRTD || RTMR0) must be in the approved set.
        //     MRTD anchors the cloud's TDVF + machine type; RTMR0 anchors the cloud's
        //     TDVF configuration. Together they identify the cloud SKU + TDVF version.
        // (2) Binary release: keccak256(RTMR1 || RTMR2) must be in the approved set.
        //     RTMR1 captures the kernel image (TDVF-measured); RTMR2 captures the
        //     initrd + cmdline. With the story-kernel binary baked into the initrd as
        //     PID 1 init, these tuples hardware-bind the user-space binary identity
        //     (TDVF measures kernel/initrd before user-space code can run, so an
        //     attacker cannot fake matching values without actually loading our binary).
        bytes32 cloudKey;
        bytes32 binaryKey;
        assembly {
            // Allocate a 96-byte scratch buffer at the free memory pointer; copy MRTD
            // (48 bytes) from the calldata quote, then RTMR0 (48 bytes) directly after.
            // This avoids the overhead of abi.encodePacked / new bytes allocation.
            let ptr := mload(0x40)
            calldatacopy(ptr, add(enclaveReport.offset, OFFSET_MRTD), MEASUREMENT_SIZE)
            calldatacopy(add(ptr, MEASUREMENT_SIZE), add(enclaveReport.offset, OFFSET_RTMR0), MEASUREMENT_SIZE)
            cloudKey := keccak256(ptr, mul(MEASUREMENT_SIZE, 2))
            // Reuse the same scratch buffer for the (RTMR1, RTMR2) tuple.
            calldatacopy(ptr, add(enclaveReport.offset, OFFSET_RTMR1), MEASUREMENT_SIZE)
            calldatacopy(add(ptr, MEASUREMENT_SIZE), add(enclaveReport.offset, OFFSET_RTMR2), MEASUREMENT_SIZE)
            binaryKey := keccak256(ptr, mul(MEASUREMENT_SIZE, 2))
        }
        require($.approvedCloudPlatforms[cloudKey], "TDXValidationHook: unapproved cloud platform");
        require($.approvedBinaryReleases[binaryKey], "TDXValidationHook: unapproved binary release");

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

        // Step 4: Option D decomposed identity match against the inner V4.
        // Caveat (Azure paravisor): paravisor may extend RTMRs differently from
        // standard TDVF — operators MUST register the paravisor-specific
        // (MRTD, RTMR0) and (RTMR1, RTMR2) tuples observed under that vendor.
        // The on-chain check is the same shape; only the registered values differ.
        (bytes32 cloudKey, bytes32 binaryKey) = _bundleDecomposedKeys(bundle.tdxV4);
        require($.approvedCloudPlatforms[cloudKey], "TDXValidationHook: unapproved cloud platform");
        require($.approvedBinaryReleases[binaryKey], "TDXValidationHook: unapproved binary release");

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

    /// @dev Option D — memory-buffer version of the decomposed key extraction.
    ///      Used by the bundle path because the inner V4 is held in `bytes memory`
    ///      after envelope parsing.
    /// @return cloudKey  keccak256(MRTD || RTMR0)
    /// @return binaryKey keccak256(RTMR1 || RTMR2)
    function _bundleDecomposedKeys(bytes memory v4) internal pure returns (bytes32 cloudKey, bytes32 binaryKey) {
        // 96-byte scratch buffer for the two paired hashes.
        bytes memory cloudBuf = new bytes(MEASUREMENT_SIZE * 2);
        bytes memory binaryBuf = new bytes(MEASUREMENT_SIZE * 2);
        for (uint256 i = 0; i < MEASUREMENT_SIZE; i++) {
            cloudBuf[i] = v4[OFFSET_MRTD + i];
            cloudBuf[MEASUREMENT_SIZE + i] = v4[OFFSET_RTMR0 + i];
            binaryBuf[i] = v4[OFFSET_RTMR1 + i];
            binaryBuf[MEASUREMENT_SIZE + i] = v4[OFFSET_RTMR2 + i];
        }
        cloudKey = keccak256(cloudBuf);
        binaryKey = keccak256(binaryBuf);
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

    /// @notice Returns whether the given (MRTD, RTMR0) tuple is approved.
    function isCloudPlatformApproved(
        bytes calldata mrtd48,
        bytes calldata rtmr048
    ) external view override returns (bool) {
        if (mrtd48.length != MEASUREMENT_SIZE || rtmr048.length != MEASUREMENT_SIZE) {
            return false;
        }
        bytes32 key = keccak256(abi.encodePacked(mrtd48, rtmr048));
        return _getTDXValidationHookStorage().approvedCloudPlatforms[key];
    }

    /// @notice Returns whether the given (RTMR1, RTMR2) tuple is approved.
    function isBinaryReleaseApproved(
        bytes calldata rtmr148,
        bytes calldata rtmr248
    ) external view override returns (bool) {
        if (rtmr148.length != MEASUREMENT_SIZE || rtmr248.length != MEASUREMENT_SIZE) {
            return false;
        }
        bytes32 key = keccak256(abi.encodePacked(rtmr148, rtmr248));
        return _getTDXValidationHookStorage().approvedBinaryReleases[key];
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
