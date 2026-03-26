# TEE Client Upgrade Design (Comprehensive - Mar 2026)

## Core Constraint
- `enclave/quote.go:ValidateCodeCommitment` compares request codeCommitment vs binary MRENCLAVE
- Single story-kernel binary can ONLY handle requests for its own code commitment
- CL has single `teeClient` in Keeper, single `TEEEndpoint` in DKGConfig
- UniqueKey (MRENCLAVE) sealing: new binary CANNOT unseal old binary's sealed data

## Dual Client Solution: TEEClientRouter
- CL-level router implements `types.TEEClient` interface
- Routes requests based on `CodeCommitment` to primary (new) or legacy (old) client
- **Port strategy**: old binary stays on 50051, new binary on 50052 (no port swapping)
- DKGConfig extended with `TEEUpgradeEndpoint` for new binary
- No story-kernel changes needed (its security model preserved)

## Upgrade Flow (Revised)
1. New binary built, MRENCLAVE determined
2. Governance: `whitelistEnclaveType` on DKG.sol with new code commitment
3. Governance: `scheduleUpgrade(activationHeight)` -- no code_commitment param
4. Validators: start new binary on 50052, configure upgrade endpoint
5. At activation height: BeginBlocker auto-starts upgrade resharing round
6. Router auto-switches: old CC -> 50051 (legacy), new CC -> 50052 (primary)
7. Old active round stays active (serves decrypt) during resharing
8. When new round reaches Active: old round deactivated, legacy client disconnected
9. Validators: stop old binary on 50051

## Key Design Decisions
- **Parallel rounds**: old active round NOT ended during upgrade resharing
- **Port assignment**: binary-based (50051=old, 50052=new), NOT role-based
- Router ActivateUpgrade() called from BeginBlocker (consensus height trigger)
- Router CompleteUpgrade() called from FinalizeDKGRound when upgrade round succeeds
- `scheduleUpgrade` has no code_commitment param -- whitelist managed separately
- `scheduleUpgrade` works even when paused (emergency upgrade scenario)
- Upgrade resharing failure -> SkipToNextRound with retry (reuses SkipToNextRound pattern)

## Code Commitment Decoupling from Rounds
- Remove `code_commitment` from DKGNetwork (deprecated field, keep field number)
- Add `code_commitment` to DKGRegistration (per-validator, per-registration)
- Collection key changes: `{ccHex}_{round}` -> `{round}` for DKGNetworks, DKGRegistrations, GlobalPubKeyVotes
- TEEUpgradeInfos key unchanged (keyed by code commitment)
- **REQUIRES UPGRADE HANDLER** for collection key migration

## Contract Changes (DKG.sol)
- `scheduleUpgrade(uint256 activationHeight)` -- onlyOwner, no whenNotPaused
- `cancelUpgrade()` -- onlyOwner, before activation
- `UpgradeScheduled(uint256 activationHeight)` event
- `UpgradeCancelled(uint256 cancelledHeight)` event
- Storage: `pendingActivationHeight`, `upgradeActive`

## DKGNetwork Proto Addition
- `bool is_upgrade_round = 12` -- marks upgrade-triggered resharing rounds
- Additive field, no upgrade handler needed for this alone

## Security Notes
- Dual client preserves UniqueKey sealing security (t+ collusion defense)
- Legacy crash during resharing -> round fails -> retry via SkipToNextRound
- 1/3+ validator non-participation blocks resharing -> governance intervention needed
- Activation is deterministic (BeginBlocker) -> no race conditions between nodes
- ThresholdDecrypt continues on old round until new round Active

## Edge Cases
- Resharing failure: SkipToNextRound with IsUpgradeRound flag preserved for retry
- Validator doesn't upgrade: excluded from new round, threshold still met if enough participate
- Previous binary crash: that validator's deals missing, threshold covers it
- Emergency: pause() + scheduleUpgrade() both callable, pause blocks register/finalize -> round fails -> retry after unpause
