# Agent Memory: Cryptography Auditor

## Codebase Knowledge
- DKG: Pedersen DKG via kyber v4 (Edwards25519)
- TDH2: cb-mpc library for threshold decryption
- Polynomial persistence: `dkgutil/poly_persist.go` — reflect-based access to kyber internals
- Signature scheme: secp256k1 ECDSA for kernel communication keys
- HKDF: Used for AES key derivation in partial decrypt encryption layer

## Past False Positives
- CRYPTO-004, CRYPTO-007: flagged issues that were not real vulnerabilities

## Blind Spots Discovered
- **Signature preimage mismatch (STOR-1, CDR-004):** Kernel signs `round || ciphertext || encryptedPartial || ephPubKey || pubShare` but omits `requesterPubKey` and `label`. CL verifier also omits them. This means partials from different requests with same ciphertext/round are interchangeable. I checked math correctness but missed signature binding completeness
- **pubKeyShare format mismatch (STOR-13):** Finalization stores raw pubKeyShare from kernel. Partial decrypt submits cb-mpc prefixed pubShare. Format differs at the CL verification point. I checked crypto operations in isolation but NEVER compared serialization across boundaries
- **EIP-191 prefix inconsistency (CDR-004):** verifyFinalizationSignature uses Ethereum signed message prefix. verifyPartialDecryptionSignature does NOT. I should compare ALL signature paths for convention consistency

## Effective Techniques
- Verifying polynomial evaluation math
- Checking commitment scheme correctness

## Ineffective Techniques
- Checking crypto operations in isolation per-layer without comparing serialization/convention across layers
- Assuming "if the math is correct, the protocol is secure" — binding and format matter as much as math
