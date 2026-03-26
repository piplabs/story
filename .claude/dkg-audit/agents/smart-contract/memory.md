# Agent Memory: Smart Contract & Deploy Script Auditor

## Codebase Knowledge
- DKG.sol: register(), finalize(), scheduleUpgrade() — UUPS proxy, Ownable2Step
- CDR.sol: allocate(), write(), read(), submitEncryptedPartialDecryption() — ERC-7201 storage
- SGXValidationHook.sol: validates SGX attestation via Automata DCAP
- Deploy scripts: GenerateAlloc.s.sol generates genesis allocation
- CDR fee: _collectFee burns msg.value via transfer to address(0)

## Past False Positives
- None identified yet

## Blind Spots Discovered
- **Permissionless submit (STOR-4, STOR-6):** `submitEncryptedPartialDecryption` is `external payable` with NO sender validation. Anyone can call with garbage data, and CL rewards it. I focused on reentrancy/storage but missed "who can call this?"
- **finalize() missing msg.sender check (CDR-005):** register() checks `validatorAddr == msg.sender` but finalize() does NOT. I should have compared access control across all similar functions
- **CDR condition bypass (CDR-006):** If msg.sender IS the condition contract, the condition check is skipped entirely. Compromised condition contract gets unrestricted access
- **Fee unit mismatch (STOR-15):** CDR fee events emit raw wei, but CL mints as gwei-scaled stake. I never compared CDR fee path with staking bridge normalization pattern

## Effective Techniques
- UUPS proxy safety checks, ERC-7201 storage layout verification

## Ineffective Techniques
- Focusing on Solidity-specific vulnerabilities (reentrancy, overflow) without tracing event processing on CL side
- Assuming contract functions are only called by intended callers
