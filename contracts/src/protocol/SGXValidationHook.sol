// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { PausableUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

import { ISGXValidationHook } from "../interfaces/ISGXValidationHook.sol";
import { IAutomataDcapAttestationFee } from "../interfaces/external/IAutomataDcapAttestationFee.sol";
import { BytesUtils } from "../libraries/BytesUtils.sol";

contract SGXValidationHook is ISGXValidationHook, Ownable2StepUpgradeable, PausableUpgradeable {
    using BytesUtils for bytes;

    /// @dev Storage structure for the SGXValidationHook
    /// @param automataValidationAddr The address of the automata validation contract
    /// @param tcbEvaluationDataNumber The TCB evaluation data number. Retained for storage and
    ///        interface backward compatibility; no longer consumed by `validateReport`, which now
    ///        calls the parameterless `verifyAndAttestOnChain` overload so the verifier resolves
    ///        the standard TCB evaluation data number on-chain per Intel's TCB Recovery policy.
    /// @custom:storage-location erc7201:story.SGXValidationHook
    struct SGXValidationHookStorage {
        address automataValidationAddr;
        uint32 tcbEvaluationDataNumber;
    }

    address public immutable DKG;

    // keccak256(abi.encode(uint256(keccak256("story.SGXValidationHook")) - 1)) & ~bytes32(uint256(0xff));
    bytes32 private constant SGXValidationHookStorageLocation =
        0xb6733d04ab09a9ab7321af14605111bb04c61e96f75dd35de2adb36bd07c7a00;

    constructor(address dkg) {
        require(dkg != address(0), "SGXValidationHook: DKG cannot be empty");
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
    /// @dev Retained for backward compatibility. The stored value is no longer read by
    ///      `validateReport`; the standard TCB evaluation data number is resolved by the verifier
    ///      on every call instead.
    function setTcbEvaluationDataNumber(uint32 newTcbEvaluationDataNumber) external onlyOwner {
        _setTcbEvaluationDataNumber(newTcbEvaluationDataNumber);
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Authentication Logic                         //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Validates the enclave report
    /// @param expectedCodeCommitment The expected code commitment
    /// @param expectedDataCommitment The expected data commitment
    /// @param enclaveReport The enclave report
    /// @param validationContext The validation context
    /// @return True if the report is valid, false otherwise
    function validateReport(
        bytes32 expectedCodeCommitment,
        bytes32 expectedDataCommitment,
        bytes calldata enclaveReport,
        bytes calldata validationContext
    ) external override returns (bool) {
        require(msg.sender == DKG, "SGXValidationHook: Only DKG can call this function");
        require(enclaveReport.length > 0, "SGXValidationHook: Empty enclave report");
        require(expectedCodeCommitment != bytes32(0), "SGXValidationHook: Zero code commitment");
        require(expectedDataCommitment != bytes32(0), "SGXValidationHook: Zero data commitment");
        SGXValidationHookStorage storage $ = _getSGXValidationHookStorage();

        // Use the parameterless verifyAndAttestOnChain overload. The Automata verifier resolves
        // the standard TCB Evaluation Data Number on-chain per Intel's TCB Recovery policy
        // (highest evaluation data number whose recovery event date is ≥ 12 months before
        // block.timestamp), so version transitions take effect automatically without any owner
        // action. See AutomataDcapAttestationFee.sol#L23 (1-arg overload delegates to
        // `_verifyAndAttestOnChain(rawQuote, 0)`, and the `0` triggers the standard lookup).
        (bool success, ) = IAutomataDcapAttestationFee($.automataValidationAddr).verifyAndAttestOnChain(enclaveReport);
        require(success, "SGXAttestationReportValidator: Attestation failed");

        // check report code commitment against expectedCodeCommitment
        require(
            _extractReportCodeCommitment(enclaveReport) == expectedCodeCommitment,
            "SGXAttestationReportValidator: Code commitment does not match"
        );

        // check report data commitment against expectedDataCommitment
        require(
            _extractReportInstanceDataCommitment(enclaveReport) == expectedDataCommitment,
            "SGXAttestationReportValidator: Data commitment does not match"
        );

        return true;
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                              Get Functions                             //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Gets the address of the automata validation contract
    /// @return The address of the automata validation contract
    function automataValidationAddr() external view returns (address) {
        return _getSGXValidationHookStorage().automataValidationAddr;
    }

    /// @notice Gets the TCB evaluation data number stored at initialization time.
    /// @dev Retained for backward compatibility. This value is no longer used by `validateReport`.
    /// @return The TCB evaluation data number
    function tcbEvaluationDataNumber() external view returns (uint32) {
        return _getSGXValidationHookStorage().tcbEvaluationDataNumber;
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                           Internal Functions                           //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Sets the address of the automata validation contract
    /// @param newAutomataValidationAddr The address of the automata validation contract
    function _setAutomataValidationAddr(address newAutomataValidationAddr) internal {
        require(newAutomataValidationAddr != address(0), "SGXValidationHook: Automata Validation cannot be empty");
        _getSGXValidationHookStorage().automataValidationAddr = newAutomataValidationAddr;
    }

    /// @notice Sets the TCB evaluation data number
    /// @param newTcbEvaluationDataNumber The TCB evaluation data number
    function _setTcbEvaluationDataNumber(uint32 newTcbEvaluationDataNumber) internal {
        _getSGXValidationHookStorage().tcbEvaluationDataNumber = newTcbEvaluationDataNumber;
    }

    /// @dev Extracts the code commitment (MRENCLAVE) from the raw SGX quote
    /// @param enclaveReport The raw SGX quote (header + report body + auth data)
    /// @return The code commitment (MRENCLAVE)
    function _extractReportCodeCommitment(bytes calldata enclaveReport) internal pure returns (bytes32) {
        // SGX quote header is 48 bytes, MRENCLAVE is at offset 64 within the report body
        // Total offset from raw quote start: 48 (header) + 64 (MRENCLAVE in body) = 112
        return bytes32(enclaveReport.substring(112, 32));
    }

    /// @dev Extracts the instance data commitment from the enclave report
    /// @param enclaveReport The enclave report
    /// @return The instance data commitment
    function _extractReportInstanceDataCommitment(bytes memory enclaveReport) internal pure returns (bytes32) {
        // According to Intel’s SGX quote structure:
        // - The SGX quote header is 48 bytes in size
        // - The enclave report body is 384 bytes long
        // - The last 64 bytes of the enclave report body are reserved for report_data
        // Therefore, the starting offset for report_data is: 48 (quote header) + 320 = 368
        // See SGX-TDX-DCAP-QuoteVerificationLibrary QuoteConstants.h
        uint256 start = 368;
        bytes32 first32;
        assembly {
            first32 := mload(add(add(enclaveReport, 32), start))
        }
        return first32;
    }

    /// @dev Returns the storage struct of SGXValidationHook.
    function _getSGXValidationHookStorage() private pure returns (SGXValidationHookStorage storage $) {
        assembly {
            $.slot := SGXValidationHookStorageLocation
        }
    }
}
