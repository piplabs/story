// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

interface IAttestationReportValidator {
    // validateReport accepts a remote attestation report (e.g. rawQuote for SGX)
    // and validates it following these steps:
    // 1. validate the report
    //   1.1 checks the report size and format
    //   1.2 checks the signature included in the report
    //   1.3 checks if the signature is from an authorized authority
    // 2. extracts the code commitment from the report (e.g. MRENCLAVE in SGX)
    //    and compares it with the expected value. This ensures the code and
    //    init data for loading the enclave was correct and untampered.
    // 3. extracts the data commitment from the report
    //    (e.g. first 32 bytes of the REPORT_DATA) and compares it with the
    //    expected value. This is instance specific, e.g. hash of node info.
    function validateReport(
        bytes32 expectedCodeCommitment,
        bytes32 expectedDataCommitment,
        bytes calldata enclaveReport,
        bytes calldata validationContext
    ) external returns (bool);
}
