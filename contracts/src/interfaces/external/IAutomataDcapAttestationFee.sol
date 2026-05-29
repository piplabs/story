// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

interface IAutomataDcapAttestationFee {
    /// @notice Verify a raw SGX/TDX quote using the *standard* TCB Evaluation Data Number resolved
    ///         on-chain by the Automata PCCS Router per Intel's TCB Recovery policy (the highest
    ///         evaluation data number whose recovery event date is ≥ 12 months before
    ///         `block.timestamp`). Use this overload to follow Intel-defined standard transitions
    ///         automatically rather than pinning a specific evaluation data number.
    function verifyAndAttestOnChain(
        bytes calldata rawQuote
    ) external payable returns (bool success, bytes memory output);

    /// @notice Verify a raw SGX/TDX quote against a caller-supplied TCB Evaluation Data Number.
    /// @dev Passing `0` is equivalent to the no-arg overload (falls back to the standard version).
    function verifyAndAttestOnChain(
        bytes calldata rawQuote,
        uint32 tcbEvaluationDataNumber
    ) external payable returns (bool success, bytes memory output);
}
