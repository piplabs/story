// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

interface IAttestationReportValidator {
    // validateReport accepts a remote attestation report (e.g. rawQuote for SGX) and validates a remote attestation report following these steps
    // 1. validate the report
    //   1.1 checks the report size and format
    //   1.2 checks the signature included in the report
    //   1.3 checks if the signature is from an authorized authority (e.g. valid certificate etc)
    // 2. extracts the code commitment from the report (e.g. MRENCLAVE in SGX) and compares it with the expected value
    // this ensures the code and init data for loading the enclave was correct and untamperred
    // 3. extracts the data commitment from the report (e.g. first 32 bytes of the REPORT_DATA) and compares it with the expected value
    // this ensures the data part is correct, this is instance specific for example this is hash of node info values.
    function validateReport(
        bytes32 expectedCodeCommitment,
        bytes32 expectedDataCommitment,
        bytes calldata enclaveReport,
        bytes calldata validationContext
    ) external returns (bool);
}
