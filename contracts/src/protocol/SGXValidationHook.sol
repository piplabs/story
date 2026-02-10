// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { AccessManaged } from "@openzeppelin/contracts/access/manager/AccessManaged.sol";

import { IAttestationReportValidator } from "../interfaces/IAttestationReportValidator.sol";
import { IAutomataDcapAttestationFee } from "../interfaces/external/IAutomataDcapAttestationFee.sol";
import { BytesUtils } from "../libraries/BytesUtils.sol";

/// TODO: consider making this contract upgradeable
contract SGXValidationHook is IAttestationReportValidator, AccessManaged {
    using BytesUtils for bytes;

    uint32 tcbEvaluationDataNumber;
    address automataValidationAddr;

    constructor(address accessManager, uint32 tcbEvalNumber) AccessManaged(accessManager) {
        tcbEvaluationDataNumber = tcbEvalNumber;
    }

    /*//////////////////////////////////////////////////////////////////////////
    //                             Admin Setters                              //
    //////////////////////////////////////////////////////////////////////////*/

    /// @notice Sets the address of the automata validation contract
    /// @param newAutomataValidationAddr The address of the automata validation contract
    function setAutomataValidationAddr(address newAutomataValidationAddr) external restricted {
        require(newAutomataValidationAddr != address(0), "SGXValidationHook: Automata Validation cannot be empty");
        automataValidationAddr = newAutomataValidationAddr;
    }

    /// @notice Sets the TCB evaluation data number
    /// @param newTcbEvaluationDataNumber The TCB evaluation data number
    function setTcbEvaluationDataNumber(uint32 newTcbEvaluationDataNumber) external restricted {
        tcbEvaluationDataNumber = newTcbEvaluationDataNumber;
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
        // see verifyAndAttestOnChain  https://github.com/automata-network/automata-dcap-attestation/blob/4e7ab275ca8c358895a83fb6d51c9bd40ba1cf68/evm/contracts/AutomataDcapAttestationFee.sol#L23
        (bool success, bytes memory output) = IAutomataDcapAttestationFee(automataValidationAddr)
            .verifyAndAttestOnChain(enclaveReport, tcbEvaluationDataNumber);
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
    //                           Internal Functions                           //
    //////////////////////////////////////////////////////////////////////////*/

    /// @dev Extracts the code commitment from the enclave report
    /// @param enclaveReport The enclave report
    /// @return The code commitment
    function _extractReportCodeCommitment(bytes calldata enclaveReport) internal returns (bytes32) {
        return bytes32(enclaveReport.substring(64, 32));
    }

    /// @dev Extracts the instance data commitment from the enclave report
    /// @param enclaveReport The enclave report
    /// @return The instance data commitment
    function _extractReportInstanceDataCommitment(bytes memory enclaveReport) internal returns (bytes32) {
        // According to Intel’s SGX quote structure:
        // - The SGX quote header is 48 bytes in size
        // - The enclave report body is 384 bytes long
        // - The last 64 bytes of the enclave report body are reserved for report_data
        // Therefore, the starting offset for report_data is: 48 (quote header) + 320 = 368
        // https://github.com/intel/SGX-TDX-DCAP-QuoteVerificationLibrary/blob/16b7291a7a86e486fdfcf1dfb4be885c0cc00b4e/Src/AttestationLibrary/src/QuoteVerification/QuoteConstants.h
        uint256 start = 368;
        bytes32 first32;
        assembly {
            first32 := mload(add(add(enclaveReport, 32), start))
        }
        return first32;
    }
}
