// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

interface IAutomataDcapAttestationFee {
    function verifyAndAttestOnChain(
        bytes calldata rawQuote,
        uint32 tcbEvaluationDataNumber
    ) external payable returns (bool success, bytes memory output);
}
