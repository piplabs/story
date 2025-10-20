# Linear Ticket: Validator Jailing Investigation

## Ticket Details

**Title:** Investigate: Malicious validator not being jailed despite proposing invalid blocks

**Assignee:** Hans

**Priority:** Medium

**Type:** Bug

**Labels:** validator, consensus, devnet, investigation

---

## Description

### Issue
A malicious validator (Validator 1, EVM addr `0xC5c0BEEAC8B37eD52F6A675eE2154D926a88E3ec`) is no longer being selected as a proposer but has not been jailed.

### Context
- **Branch:** hans/malicious-empty-tx
- **Related PR:** #144
- **Environment:** internal-devnet
- **Other nodes:** Upgraded to release/1.3 (commit f190708)

### Observed Behavior
1. Validator 1 runs malicious code and continuously prints: `err="runtime error: invalid memory address or nil pointer dereference"`
2. Other nodes (RPC and bootnode) print: `err="validate tx: invalid fee in tx"`
3. Validator 1 is no longer selected as a proposer
   - Block finalization takes longer when the malicious validator is the proposer (gets rejected, system selects alternative)
4. **Issue:** Validator 1 is NOT jailed and still appears in validators API: https://rpc.devnet.storyrpc.io/validators?per_page=200
5. Validator 1 continues to receive rewards: https://devnet.storyscan.io/address/0xC5c0BEEAC8B37eD52F6A675eE2154D926a88E3ec?tab=index

### Current Understanding
- The malicious validator proposes invalid blocks but votes correctly on valid blocks
- Current jailing mechanism only covers: **downtime** and **double signing**
- This case doesn't trigger either condition, so no jailing occurs
- Validator remains in active validator list with voting power
- The network considers the validator as having missed its proposing chance, but not as malicious

### Expected Behavior
**This behavior does NOT meet expectations** - need to investigate:
1. Why the validator is not being selected as proposer
2. Whether additional jailing mechanisms are needed for validators that consistently propose invalid blocks
3. Whether validators proposing invalid blocks should have rewards reduced/eliminated

### Notes
- This is not the first time this behavior has been observed in testing
- The malicious validator mocks only the proposing part, acting honestly when processing valid blocks
- Related to Cantina-131 case, but different behavior (131 causes downtime, 144 does not)

### Questions to Answer
1. Should proposing invalid blocks be considered a jailable offense?
2. Should reward distribution be modified for validators that don't contribute properly?
3. Is the current proposer selection mechanism working as intended when a validator consistently proposes invalid blocks?

---

## References
- Devnet Validator API: https://rpc.devnet.storyrpc.io/validators?per_page=200
- Malicious Validator Address: https://devnet.storyscan.io/address/0xC5c0BEEAC8B37eD52F6A675eE2154D926a88E3ec?tab=index
- Slack Thread: Discussion between Yao and Hans on 10/20/2025
