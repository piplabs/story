// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { IAttestationReportValidator } from "./IAttestationReportValidator.sol";

interface ITDXValidationHook is IAttestationReportValidator {
    /// @notice Sets the address of the Automata DCAP attestation contract.
    function setAutomataValidationAddr(address newAutomataValidationAddr) external;

    /// @notice Approves a platform identity tuple keyed by keccak256(MRTD || RTMR0 || RTMR1).
    /// @dev Hybrid identity model: binary identity is enforced via DKG.codeCommitment, platform
    ///      identity via this whitelist. See TDXValidationHook for details.
    function approvePlatform(bytes32 platformCommitment, string calldata label) external;

    /// @notice Revokes a previously approved platform identity tuple.
    function revokePlatform(bytes32 platformCommitment) external;

    /// @notice Returns whether the given platform commitment is approved.
    function isPlatformApproved(bytes32 platformCommitment) external view returns (bool);

    /// @notice Returns the address of the Automata DCAP attestation contract.
    function automataValidationAddr() external view returns (address);
}
