# Agent Memory: Cosmos SDK & Upgrade Handler Auditor

## Codebase Knowledge
- DKG module at `client/x/dkg/` — keeper, types, module
- Upgrade handler at `client/x/dkg/module/` — v2.0.0 upgrade from v1.5.3
- Vote extensions: ExtendVote/VerifyVoteExtension in keeper, PrepareProposal aggregates
- Reward distribution: `dkg_rewards.go` (DKG rewards) and `dkg_cdr_fees.go` (CDR fees)
- `dkg_rewards.go` correctly sorts with `sort.Strings(memberAddrs)` before map iteration
- Params: dual registration issue — `ParamsStore` (collections.Item) vs raw KV store access

## Past False Positives
- COSMOS-006, COSMOS-007, COSMOS-008, COSMOS-010 from first audit — flagged issues that were design-documented behaviors

## Blind Spots Discovered
- **Map iteration determinism:** Found it in rewards but MISSED it in CDR fee distribution (`dkg_cdr_fees.go:182`). Lesson: check ALL map iterations, not just the obvious ones
- **Vote extension dedup ordering (STOR-7):** `aggregateVotes()` does first-wins dedup BEFORE DKG message authentication. Malicious earlier-sorted validator can erase honest messages. I completely missed adversarial ordering in vote extension processing
- **Stage transition boundary bugs (STOR-19, STOR-22, STOR-23):** The "last block" of one DKG stage vs "first block" of the next stage creates off-by-one windows where valid events are rejected. I need to check every stage transition boundary

## Effective Techniques
- Comparing parallel code paths (rewards vs fees) to find inconsistencies — this found C-01
- Checking all store access patterns for consistency

## Ineffective Techniques
- Only checking "does this code work correctly?" without asking "can this code be abused by someone who controls the input ordering?"
