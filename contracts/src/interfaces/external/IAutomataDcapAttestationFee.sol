// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

interface IAutomataDcapAttestationFee {
    /// @notice Verify a raw SGX/TDX quote using the standard TCB Evaluation Data Number
    ///         resolved on-chain by the PCCS Router per Intel's TCB Recovery policy.
    function verifyAndAttestOnChain(
        bytes calldata rawQuote
    ) external payable returns (bool success, bytes memory output);

    /// @notice Verify a raw SGX/TDX quote against a caller-supplied TCB Evaluation Data Number.
    /// @dev Passing `0` falls back to the standard resolution (equivalent to the no-arg overload).
    function verifyAndAttestOnChain(
        bytes calldata rawQuote,
        uint32 tcbEvaluationDataNumber
    ) external payable returns (bool success, bytes memory output);
}
