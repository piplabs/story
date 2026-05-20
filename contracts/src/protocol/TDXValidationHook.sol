// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { PausableUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

import { ITDXValidationHook } from "../interfaces/ITDXValidationHook.sol";
import { IAutomataDcapAttestationFee } from "../interfaces/external/IAutomataDcapAttestationFee.sol";

/// @title TDXValidationHook
/// @notice On-chain validator for Intel TDX V4/V5 attestation quotes used by the DKG flow.
/// @dev Hybrid identity model:
///      - Binary identity: keccak256(RTMR2), matched against DKG.codeCommitment. RTMR2 is the
///        TDVF-measured digest of initrd + cmdline; since story-kernel runs as PID 1 from initrd
///        and TDVF measures it before any user-space code executes, this hardware-binds the
///        binary at the same DKG storage slot SGX uses for MRENCLAVE.
///      - Platform identity: keccak256(MRTD || RTMR0 || RTMR1), enforced via `approvedPlatforms`
///        whitelist. MRTD + RTMR0 pin cloud SKU/TDVF; RTMR1 pins kernel image. Decouples
///        cloud-vendor lifecycle from binary-release lifecycle without changing DKG storage.
///
///      V4 and V5 share header layout and the absolute byte offsets for every field this
///      contract reads (Automata routes per-version internally). The TEE-type guard below is
///      defense-in-depth against an SGX quote being routed into the TDX hook; Automata is the
///      authoritative shape validator.
contract TDXValidationHook is ITDXValidationHook, Ownable2StepUpgradeable, PausableUpgradeable {
    /*//////////////////////////////////////////////////////////////////////////
    //                              Events                                    //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Emitted when the Automata DCAP attestation contract address is (re)set.
    event AutomataValidationAddrSet(address indexed newAutomataValidationAddr);

    /// @notice Emitted when a platform identity tuple is approved.
    event PlatformApproved(bytes32 indexed platformCommitment, string label);

    /// @notice Emitted when a platform identity tuple is revoked.
    event PlatformRevoked(bytes32 indexed platformCommitment);

    /// @dev Maps keccak256(MRTD || RTMR0 || RTMR1) to approval status.
    /// @custom:storage-location erc7201:story.TDXValidationHook
    struct TDXValidationHookStorage {
        address automataValidationAddr;
        mapping(bytes32 => bool) approvedPlatforms;
    }

    address public immutable DKG;

    // keccak256(abi.encode(uint256(keccak256("story.TDXValidationHook")) - 1)) & ~bytes32(uint256(0xff));
    bytes32 private constant TDXValidationHookStorageLocation =
        0xc82cee3b70f9b4f764f18b66d0052cef718142a50149d16a1c68ba0074816f00;

    /*//////////////////////////////////////////////////////////////////////////
    //                          Quote layout constants                        //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev TDX quote header size; identical in V4 and V5.
    uint256 private constant QUOTE_HEADER_SIZE = 48;

    /// @dev MRTD and each RTMR are 48 bytes (SHA-384).
    uint256 private constant MEASUREMENT_SIZE = 48;

    /// @dev TDX tee_type wire byte at offset 4 (little-endian uint32 = 0x00000081). SGX has 0.
    uint8 private constant TEE_TYPE_TDX_BYTE0 = 0x81;

    // Absolute offsets into the raw quote (header + body). Identical for V4 and V5.
    uint256 private constant OFFSET_MRTD = 184;
    uint256 private constant OFFSET_RTMR0 = 376;
    uint256 private constant OFFSET_RTMR1 = 424;
    /// @dev RTMR2 carries the binary identity (initrd + cmdline measurement).
    uint256 private constant OFFSET_RTMR2 = 472;
    /// @dev REPORT_DATA is 64 bytes; we consume only the first 32.
    uint256 private constant OFFSET_REPORT_DATA = 568;

    /// @dev Minimum quote length defense-in-depth: the highest field we read is the first 32
    ///      bytes of REPORT_DATA. A real V4 TD10 quote is ≥ 632 bytes and a V5 TD15 quote is
    ///      ≥ 696 bytes (header + body), so this floor is well below either. Automata is the
    ///      authoritative shape validator (it parses header + body + signature/cert chain) but
    ///      hardening this floor lets the assembly extractors below treat their calldata
    ///      offsets as in-bounds without trusting the Automata side-effect order.
    uint256 private constant MIN_QUOTE_SIZE = OFFSET_REPORT_DATA + 32; // 600

    constructor(address dkg) {
        require(dkg != address(0), "TDXValidationHook: DKG cannot be empty");
        DKG = dkg;
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    function initialize(address owner, address automataValidationAddr) external initializer {
        require(owner != address(0), "TDXValidationHook: owner cannot be empty");
        __Ownable_init(owner);
        __Pausable_init();

        _setAutomataValidationAddr(automataValidationAddr);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Admin Setters                              //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Sets the address of the Automata DCAP attestation contract.
    function setAutomataValidationAddr(address newAutomataValidationAddr) external onlyOwner {
        _setAutomataValidationAddr(newAutomataValidationAddr);
    }

    /// @notice Approves a platform identity tuple keyed by keccak256(MRTD || RTMR0 || RTMR1).
    /// @dev The key is computed from the raw 48-byte SHA-384 measurements concatenated in
    ///      MRTD || RTMR0 || RTMR1 order, matching what the hook derives from incoming quotes.
    function approvePlatform(bytes32 platformCommitment, string calldata label) external override onlyOwner {
        require(platformCommitment != bytes32(0), "TDXValidationHook: platform commitment cannot be empty");
        _getTDXValidationHookStorage().approvedPlatforms[platformCommitment] = true;
        emit PlatformApproved(platformCommitment, label);
    }

    /// @notice Revokes a previously approved platform identity tuple.
    function revokePlatform(bytes32 platformCommitment) external override onlyOwner {
        require(platformCommitment != bytes32(0), "TDXValidationHook: platform commitment cannot be empty");
        delete _getTDXValidationHookStorage().approvedPlatforms[platformCommitment];
        emit PlatformRevoked(platformCommitment);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Authentication Logic                         //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Validates a TDX attestation quote against the DKG-supplied binary commitment and
    ///      the kernel-bound data commitment carried in REPORT_DATA[0:32].
    function validateReport(
        bytes32 expectedCodeCommitment,
        bytes32 expectedDataCommitment,
        bytes calldata enclaveReport,
        bytes calldata /* validationContext */
    ) external override returns (bool) {
        require(msg.sender == DKG, "TDXValidationHook: Only DKG can call this function");
        require(enclaveReport.length > 0, "TDXValidationHook: Empty enclave report");
        require(expectedCodeCommitment != bytes32(0), "TDXValidationHook: Zero code commitment");
        require(expectedDataCommitment != bytes32(0), "TDXValidationHook: Zero data commitment");

        // Reject anything obviously not a TDX quote before reading at TDX-specific offsets.
        // Automata is the authoritative shape validator; this is defense-in-depth only.
        // The length floor below ensures every assembly extractor reads in-bounds calldata
        // without trusting Automata to have done so first.
        require(enclaveReport.length >= QUOTE_HEADER_SIZE, "TDXValidationHook: Quote too short for header");
        require(enclaveReport.length >= MIN_QUOTE_SIZE, "TDXValidationHook: Quote too short for body");
        require(
            uint8(enclaveReport[4]) == TEE_TYPE_TDX_BYTE0 &&
                uint8(enclaveReport[5]) == 0 &&
                uint8(enclaveReport[6]) == 0 &&
                uint8(enclaveReport[7]) == 0,
            "TDXValidationHook: Not a TDX quote"
        );

        // No-arg overload: Automata resolves the standard TCB Evaluation Data Number from the
        // PCCS Router per Intel's TCB Recovery policy (mirrors PR #816 for the SGX hook).
        TDXValidationHookStorage storage $ = _getTDXValidationHookStorage();
        (bool success, ) = IAutomataDcapAttestationFee($.automataValidationAddr).verifyAndAttestOnChain(enclaveReport);
        require(success, "TDXValidationHook: Attestation failed");

        require(
            _computeBinaryCommitment(enclaveReport) == expectedCodeCommitment,
            "TDXValidationHook: unapproved binary"
        );
        require(
            $.approvedPlatforms[_computePlatformCommitment(enclaveReport)],
            "TDXValidationHook: unapproved platform"
        );
        require(
            _extractReportInstanceDataCommitment(enclaveReport) == expectedDataCommitment,
            "TDXValidationHook: Data commitment does not match"
        );

        return true;
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              Get Functions                             //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Returns the address of the Automata DCAP attestation contract.
    function automataValidationAddr() external view override returns (address) {
        return _getTDXValidationHookStorage().automataValidationAddr;
    }

    /// @notice Returns whether the given platform commitment is approved.
    function isPlatformApproved(bytes32 platformCommitment) external view override returns (bool) {
        return _getTDXValidationHookStorage().approvedPlatforms[platformCommitment];
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Internal Functions                           //
    //////////////////////////////////////////////////////////////////////////*/

    function _setAutomataValidationAddr(address newAutomataValidationAddr) internal {
        require(newAutomataValidationAddr != address(0), "TDXValidationHook: Automata Validation cannot be empty");
        _getTDXValidationHookStorage().automataValidationAddr = newAutomataValidationAddr;
        emit AutomataValidationAddrSet(newAutomataValidationAddr);
    }

    /// @dev keccak256(RTMR2). Matched against DKG.codeCommitment as the binary identity.
    function _computeBinaryCommitment(bytes calldata enclaveReport) internal pure returns (bytes32 result) {
        assembly {
            let ptr := mload(0x40)
            calldatacopy(ptr, add(enclaveReport.offset, OFFSET_RTMR2), MEASUREMENT_SIZE)
            result := keccak256(ptr, MEASUREMENT_SIZE)
        }
    }

    /// @dev keccak256(MRTD || RTMR0 || RTMR1). MRTD (offset 184) and RTMR0/RTMR1 (offset 376,
    ///      contiguous) are not adjacent in the quote, so we assemble the 144-byte preimage
    ///      in scratch memory via two calldatacopy calls.
    function _computePlatformCommitment(bytes calldata enclaveReport) internal pure returns (bytes32 result) {
        assembly {
            let ptr := mload(0x40)
            calldatacopy(ptr, add(enclaveReport.offset, OFFSET_MRTD), MEASUREMENT_SIZE)
            calldatacopy(add(ptr, MEASUREMENT_SIZE), add(enclaveReport.offset, OFFSET_RTMR0), mul(MEASUREMENT_SIZE, 2))
            result := keccak256(ptr, mul(MEASUREMENT_SIZE, 3))
        }
    }

    /// @dev First 32 bytes of REPORT_DATA. Automata guarantees the quote is long enough.
    function _extractReportInstanceDataCommitment(bytes calldata enclaveReport) internal pure returns (bytes32 result) {
        assembly {
            result := calldataload(add(enclaveReport.offset, OFFSET_REPORT_DATA))
        }
    }

    function _getTDXValidationHookStorage() private pure returns (TDXValidationHookStorage storage $) {
        assembly {
            $.slot := TDXValidationHookStorageLocation
        }
    }
}
