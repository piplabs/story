# DKG Module Patterns

## Store Keys (prefix numbers)
- 0: Params
- 1: DKGNetworks (key: `{codeCommitmentHex}_{round}`) -- planned migration to `{round}` only
- 2: LatestDKGNetwork (pointer to DKGNetworks key)
- 3: DKGRegistrations (key: `{codeCommitmentHex}_{round}_{validatorHex}`) -- planned migration to `{round}_{addr}`
- 4: LatestActiveRound (pointer to DKGNetworks key)
- 5: GlobalPubKeyVotes (key includes codeCommitment prefix -- planned migration)
- 6: TEEUpgradeInfos (key: codeCommitment -- keeping as-is)
- 7: SettlementBalance

## DKG Reward Distribution Flow
1. `BeginBlocker` detects stage transition to `DKGStageActive`
2. Calls `FinalizeDKGRound(ctx, latestRound)` where latestRound is the NEWLY finalized round
3. Inside, calls `distributeDKGCommitteeRewards(ctx)` BEFORE updating active round pointer
4. `getLatestActiveDKGNetwork` returns the PREVIOUS active round (the one whose committee served)
5. Gets finalized registrations from that previous round
6. Withdraws ALL UBI to DKG module, splits between committee and evmstaking

## Protobuf Params
- Field 9: `dkg_committee_reward_portion` (cosmos.Dec, non-nullable)
- Default: "0.10" (10%)
- Validation: [0, 1.0] inclusive

## DKG Registration Statuses
- 0: Unspecified
- 1: Verified (passed TEE attestation)
- 2: Finalized (completed DKG protocol)
