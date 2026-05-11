// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { IAttestationReportValidator } from "./IAttestationReportValidator.sol";

interface ITDXValidationHook is IAttestationReportValidator {
    /// @notice Sets the address of the automata validation contract
    /// @param newAutomataValidationAddr The address of the automata validation contract
    function setAutomataValidationAddr(address newAutomataValidationAddr) external;

    /// @notice Sets the TCB evaluation data number
    /// @param newTcbEvaluationDataNumber The TCB evaluation data number
    function setTcbEvaluationDataNumber(uint32 newTcbEvaluationDataNumber) external;

    /// @notice Approves a cloud platform identified by (MRTD, RTMR0).
    /// @dev Both inputs MUST be 48-byte SHA-384 measurements; the on-chain key is
    ///      keccak256(mrtd48 || rtmr048). See contract docs for the Option D model.
    /// @param mrtd48 The 48-byte MRTD measurement
    /// @param rtmr048 The 48-byte RTMR0 measurement
    /// @param label Free-form label for governance bookkeeping (e.g. "GCP c3-standard-4 / TDVF v1.5")
    function approveCloudPlatform(bytes calldata mrtd48, bytes calldata rtmr048, string calldata label) external;

    /// @notice Revokes a previously approved cloud platform.
    /// @param mrtd48 The 48-byte MRTD measurement
    /// @param rtmr048 The 48-byte RTMR0 measurement
    function revokeCloudPlatform(bytes calldata mrtd48, bytes calldata rtmr048) external;

    /// @notice Approves a binary release identified by (RTMR1, RTMR2).
    /// @dev Both inputs MUST be 48-byte SHA-384 measurements; the on-chain key is
    ///      keccak256(rtmr148 || rtmr248).
    /// @param rtmr148 The 48-byte RTMR1 measurement (kernel image)
    /// @param rtmr248 The 48-byte RTMR2 measurement (initrd + cmdline / OS boot)
    /// @param version Free-form label for governance bookkeeping (e.g. "story-kernel v1.7.0")
    function approveBinaryRelease(bytes calldata rtmr148, bytes calldata rtmr248, string calldata version) external;

    /// @notice Revokes a previously approved binary release.
    /// @param rtmr148 The 48-byte RTMR1 measurement
    /// @param rtmr248 The 48-byte RTMR2 measurement
    function revokeBinaryRelease(bytes calldata rtmr148, bytes calldata rtmr248) external;

    /// @notice Gets the address of the automata validation contract
    /// @return The address of the automata validation contract
    function automataValidationAddr() external view returns (address);

    /// @notice Gets the TCB evaluation data number
    /// @return The TCB evaluation data number
    function tcbEvaluationDataNumber() external view returns (uint32);

    /// @notice Returns whether the given (MRTD, RTMR0) tuple is approved.
    /// @param mrtd48 The 48-byte MRTD measurement
    /// @param rtmr048 The 48-byte RTMR0 measurement
    function isCloudPlatformApproved(bytes calldata mrtd48, bytes calldata rtmr048) external view returns (bool);

    /// @notice Returns whether the given (RTMR1, RTMR2) tuple is approved.
    /// @param rtmr148 The 48-byte RTMR1 measurement
    /// @param rtmr248 The 48-byte RTMR2 measurement
    function isBinaryReleaseApproved(bytes calldata rtmr148, bytes calldata rtmr248) external view returns (bool);
}
