// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { IAttestationReportValidator } from "./IAttestationReportValidator.sol";

interface ISGXValidationHook is IAttestationReportValidator {
    /// @notice Sets the address of the automata validation contract
    /// @param newAutomataValidationAddr The address of the automata validation contract
    function setAutomataValidationAddr(address newAutomataValidationAddr) external;

    /// @notice Sets the TCB evaluation data number
    /// @param newTcbEvaluationDataNumber The TCB evaluation data number
    function setTcbEvaluationDataNumber(uint32 newTcbEvaluationDataNumber) external;

    /// @notice Gets the address of the automata validation contract
    /// @return The address of the automata validation contract
    function automataValidationAddr() external view returns (address);

    /// @notice Gets the TCB evaluation data number
    /// @return The TCB evaluation data number
    function tcbEvaluationDataNumber() external view returns (uint32);
}