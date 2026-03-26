# Agent Memory: TEE/SGX Full Audit (story-kernel)

## Codebase Knowledge
- gRPC server: `server/server.go` — default bind `:50051` (all interfaces), TLS optional
- DKG RPCs: GenerateAndSealKey, GenerateDeals, ProcessDeals, ProcessResponses, FinalizeDKG, PartialDecryptTDH2
- Sealing: `enclave/sealed_leveldb.go` — SGX sealed storage via Gramine
- Key management: `store/key_store.go` — Ed25519 (DKG) + Secp256k1 (communication)
- State persistence: `store/dkg_state.go` — JSON state file, polynomial coefficients sealed separately
- Light client: `story/query_client.go` — verified queries to story chain
- Manifest: `story-kernel.manifest.template` — Gramine SGX configuration

## Past False Positives
- KERNEL-002: flagged state.json as containing private keys — it contains EncryptedDeal (already encrypted) and responses (no secrets). Justification SecShare is protocol-public data
- KERNEL-009: HKDF nil salt — RFC 5869 compliant, not a vulnerability
- KERNEL-010: flagged reflection usage — no reflection in current code
- KERNEL-016: flagged nonexistent function

## Blind Spots Discovered
- **Kernel as decryption oracle (STOR-5):** Default config binds to all interfaces on `:50051` without TLS. GetCodeCommitment() is public and exposes the only precondition. PartialDecryptTDH2 has no on-chain proof verification (TODO in code). I focused on enclave internals but COMPLETELY missed the network attack surface
- **gRPC auth gap:** No authentication interceptor beyond panic recovery. I should have asked "what if someone other than the CL calls this?"
- **cachedLastBlockHeight race (H-06):** int64 read/written by concurrent goroutines without atomic/mutex

## Effective Techniques
- Analyzing sealed storage security model
- Checking manifest for debug mode

## Ineffective Techniques
- Trusting that gRPC callers are always the legitimate CL client
- Focusing on enclave-internal security without modeling the network boundary
- Not checking default configuration values for security implications
