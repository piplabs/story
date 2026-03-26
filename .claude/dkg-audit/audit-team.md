# DKG/CDR Audit Team Configuration v3

## Audit Methodology — Core Principles

These are NOT per-agent checklists. Every agent MUST internalize these as their default thinking mode.

### Principle 1: Think Like an Attacker, Not a Reviewer

**Wrong approach:** "Does this function handle errors correctly?"
**Right approach:** "If I'm a malicious validator / random EOA / network attacker, how do I profit from this function?"

Every function, every return path, every event emission — ask: **who benefits if this goes wrong, and can they make it go wrong on purpose?**

### Principle 2: Follow the Money Across Boundaries

Never analyze a value transfer in isolation. Trace it end-to-end:
- Where does the value originate? (EL wei? CL stake? Fee pool?)
- Does it cross a unit boundary? (wei↔gwei, EL↔CL)
- Who receives it? Are recipients validated?
- What conditions gate the transfer? Can those conditions be spoofed?
- What happens to the value if the operation partially fails?

**The bug is almost never in the function that moves money. It's in the function that DECIDES to move money.**

### Principle 3: nil Is Not Innocuous

In Go, `return nil` means "no error" to the caller. Every `nil` return is a **positive assertion of success**. Before writing `return nil`, ask:
- Does the caller reward/refund/credit based on this nil?
- Am I returning nil because the input was invalid, missing, or expired?
- Should the caller distinguish "success" from "nothing happened"?

### Principle 4: Every `external`/`public` Function Is Called by an Attacker

Solidity contracts are permissionless by default. For every function:
- Remove the mental label of "this is called by validators"
- Assume `msg.sender` is a bot with no stake, no TEE, no reputation
- Trace the emitted event through CL processing — does the CL validate the sender?

### Principle 5: Cross-Layer Serialization Breaks Silently

When a value crosses Contract → CL → Kernel → CL:
- Format may change (raw bytes vs hex, compressed vs uncompressed, prefixed vs raw)
- Field numbering may drift (proto field renumbering)
- Signature preimages may cover different fields on each side
- **Compare the exact bytes at each boundary.** If they differ, the system is broken even if each layer "works correctly" in isolation.

### Principle 6: State Machine Transitions Create Dead Windows

When a system has states (DKG: registration→dealing→finalization→active→resharing), ask:
- What functionality is UNAVAILABLE during each state?
- Is this documented? Does the user/requester know?
- Can an attacker force the system into a state that permanently disables functionality?
- What state must carry over across transitions? What happens if it doesn't?

---

## Audit Team Composition v3 (9 Agents)

### 1. Cosmos SDK & Upgrade Handler Auditor
- **Agent type:** `story-cosmos-engineer`
- **Scope:** Upgrade handler, DKG module architecture, ABCI integration, state management, module registration, rewards & fee distribution determinism
- **Key focus:** Consensus-critical determinism, store patterns, genesis export/import
- **Must-do:**
  - Every `map` iteration in FinalizeBlock: sorted or flagged
  - Every dedup: authenticated BEFORE or AFTER? If before, attacker wins first-seen
  - Every vote extension aggregation: model malicious earlier-sorted validator

### 2. DKG Service & Lifecycle Auditor
- **Agent type:** `story-cosmos-engineer`
- **Scope:** DKG service lifecycle, protocol flow, kernel gRPC communication, vote extensions, caching, error recovery, CDR flow
- **Key focus:** Race conditions, goroutine leaks, deadlocks, state lifecycle across round transitions
- **Must-do:**
  - Map CDR read availability against every DKG state — identify dead windows
  - Trace every cache (PIDCache, RoundContextCache) population across round transitions
  - Every nil-return: what does the CALLER do? Does nil = "success" = reward?

### 3. Smart Contract & Deploy Script Auditor
- **Agent type:** `blockchain-security-reviewer`
- **Scope:** DKG.sol, CDR.sol, SGXValidationHook.sol, deploy scripts, precompiles, ABI bindings
- **Key focus:** Reentrancy, access control, storage layout, proxy safety, deploy correctness
- **Must-do:**
  - Every `external`/`public` function: what if called by non-validator EOA with garbage data?
  - Every event: trace through CL processing — does CL trust unvalidated fields?
  - Compare access control across similar functions (register vs finalize)

### 4. TEE/SGX Full Audit (story-kernel)
- **Agent type:** `story-tee-engineer`
- **Scope:** Entire story-kernel — enclave security, sealing, attestation, key management, DKG protocol, network/API, persistence, light client
- **Key focus:** Enclave boundary, secret material lifecycle, network exposure
- **Must-do:**
  - Every gRPC method: what if non-CL caller invokes it? What auth exists?
  - Default bind address + TLS: is kernel a remote oracle out of the box?
  - Does kernel verify on-chain proof before touching sealed key material?

### 5. Cryptography Auditor
- **Agent type:** `cryptography-reviewer`
- **Scope:** Both repos — Pedersen DKG, VSS, threshold crypto, elliptic curve ops, randomness, serialization, TDH2
- **Key focus:** Mathematical correctness, protocol security, cryptographic best practices
- **Must-do:**
  - Trace serialization of every crypto value across ALL layers — flag format mismatches
  - Signature preimage: does it include ALL fields needed for binding? (requester, label, pid, codeCommitment)
  - Signature convention: EIP-191 prefix vs raw keccak256 — consistent across sign/verify paths?

### 6. Cross-Cutting Security & Workflow Auditor
- **Agent type:** `blockchain-security-reviewer`
- **Scope:** Both repos — Upgrade workflow, DKG ceremony lifecycle, CDR workflow, trust boundaries, economic security, network partition, DoS
- **Key focus:** System-level attack scenarios, malicious validator behavior, availability analysis
- **Must-do:**
  - Malicious earlier-sorted validator censoring honest messages via dedup
  - Cross-layer handoff: does receiving layer validate independently?
  - Resharing→active transition: what state carries over? What breaks if it doesn't?

### 7. Economic & Value Flow Auditor
- **Agent type:** `blockchain-security-reviewer`
- **Scope:** Both repos + contracts — Fee collection, fee bridging (EL→CL), reward distribution, refund paths, staking value flow, CDR fee pool
- **Key focus:** Value conservation, unit mismatches, reward manipulation, permissionless extraction
- **Must-do:**
  - **Unit trace:** Every value crossing EL↔CL boundary — normalized? Compare with staking bridge pattern
  - **Refund gates:** Every path that refunds/rewards — ALL conditions? Can non-validator trigger? Can nil trigger?
  - **Reward fairness:** Can attacker inflate count? Can non-committee members receive?
  - **Permissionless drain:** Every `payable` function called 1000x by bot — protocol loses money?
  - **Error-path economics:** Every `return nil` that caller interprets as success → economic consequence

### 8. Adversarial Integration Tester
- **Agent type:** `blockchain-security-reviewer`
- **Scope:** Both repos — End-to-end attack scenarios spanning multiple components
- **Key focus:** Compose single-component weaknesses into multi-step exploits
- **Mandatory attack scenarios:**
  1. Garbage partial decryptions from malicious validator → rewarded?
  2. Non-validator EOA calling CDR submit → fee pool drain?
  3. Network attacker with gRPC access → decrypt vaults without on-chain read?
  4. Replayed old partials for new requests → CL detects?
  5. Inflated reward counts via expired/missing request windows
  6. Earlier-sorted validator censoring honest DKG messages via dedup collision
  7. Resharing leaving CDR permanently broken (empty PIDCache, dead read path)

### 9. **[NEW] Go Architecture & Clean Code Advisor**
- **Agent type:** `story-cosmos-engineer`
- **Scope:** Both repos — Code structure, package organization, API design, error handling patterns, concurrency patterns, testability, idiomatic Go
- **Key focus:** Clean code, structural improvements, refactoring direction
- **This agent produces a SEPARATE section** in the final report: "Architecture & Code Quality Recommendations"
- **Review areas:**
  - **Package structure:** Is responsibility clearly separated? Are there circular dependencies or God packages?
  - **Error handling:** Consistent error wrapping? Sentinel errors vs string matching? Are errors actionable to callers?
  - **Concurrency:** Mutex placement (struct field vs package global?), channel usage, goroutine lifecycle, context propagation
  - **Interface design:** Are interfaces defined where they're used (consumer side), not where they're implemented? Are they minimal?
  - **Testability:** Can components be tested in isolation? Are dependencies injectable? Do tests share global state?
  - **Naming:** Do function/variable names communicate intent? Are abbreviations consistent?
  - **Code duplication:** Same logic repeated across DKG/CDR paths? Extract shared patterns?
  - **API boundaries:** Are gRPC service interfaces clean? Do they leak internal implementation details?
  - **State management:** Is mutable state minimized? Are state transitions explicit? Is the state machine documented in code?
  - **Go idioms:** Proper use of `context.Context`, `error` wrapping, `sync` primitives, `io.Reader`/`io.Writer` patterns, table-driven tests
- **Output format:**

```markdown
## Architecture & Code Quality Recommendations

### Structural Issues (refactoring direction)
Ordered by impact. Each item includes: current state, problem, recommended direction, effort estimate.

### Pattern Violations
Go idioms / best practices that are violated consistently. Not one-off nits — patterns that should be fixed project-wide.

### Positive Patterns
What the codebase does well that should be preserved and extended.
```

---

## Agent Persistent Memory

Each agent maintains its own persistent memory at `story/.claude/dkg-audit/agents/{agent-name}/`:
- `memory.md` — Codebase knowledge, past false positives, blind spots, effective/ineffective techniques
- `findings-log.md` — Append-only history of findings + outcomes (TRUE_POSITIVE / FALSE_POSITIVE / FIXED)

See `agents/AGENT_MEMORY_GUIDE.md` for full specification.

### Agent Memory Lifecycle

**Before each audit:**
1. Agent reads its `memory.md` — recalls blind spots, false positives, effective techniques
2. Agent reads its `findings-log.md` — knows what was found before and what was missed
3. Agent uses memory to AVOID known FPs and FOCUS on known blind spots

**During audit:**
4. Agent applies lessons from memory (e.g., "last time I missed nil-return economics, check those first")
5. Agent notes new learnings as it works

**After audit (Phase 4):**
6. Agent updates `memory.md` with new knowledge, blind spots discovered, technique effectiveness
7. Agent appends new findings to `findings-log.md` with severity and initial status
8. After re-validation, orchestrator updates each agent's `findings-log.md` with TRUE/FALSE POSITIVE outcome

**Cross-agent learning:**
9. After re-validation, orchestrator shares APEX/external findings that each agent SHOULD have caught
10. Each agent's `memory.md` "Blind Spots Discovered" section is updated with these misses

This creates a **feedback loop**: each audit makes every agent better at its specific domain.

---

## Design Documents (MUST provide to all agents)
- **DKG:** `story/docs/design/DKG.md`
- **CDR:** `story/docs/design/CDR.md`
- **Whitepaper CDR section** (if available): for stated security properties and threat model
- All agents MUST read these docs before auditing to avoid false positives

## Execution Pattern

### Phase 1: Initial Audit (Parallel)
- All 9 agents run **in parallel** (`run_in_background: true`)
- Each agent reads its `memory.md` + `findings-log.md` BEFORE starting
- Each agent works in a **git worktree** (to avoid disturbing active branches)
- All agents are **READ-ONLY** (no file modifications except their own memory)
- Agent #9 (Go Advisor) runs in parallel but produces a separate output section

### Phase 2: Cross-Layer Trace (Serial, after Phase 1)
- 2 additional trace agents:
  - **EL→CL Value Trace:** Takes Phase 1 findings, traces every value/fee/reward flow end-to-end
  - **Error-Path Consequence:** Takes all nil-return findings, traces caller interpretation, identifies silent-success-on-failure

### Phase 3: Design Doc Re-validation (Parallel, after Phase 2)
- 4 re-validation agents (grouped by domain)
- Each receives: design docs + all Phase 1+2 findings
- Classify: TRUE ISSUE / FALSE POSITIVE / SEVERITY ADJUSTED

### Phase 4: Report Assembly + Memory Update
- Deduplicated security findings (APEX-quality format)
- Separate "Architecture & Code Quality" section (Agent #9)
- Remediation priority matrix
- **Each agent updates its memory.md and findings-log.md**

---

## Finding Format (APEX-quality standard)

Every security finding MUST include ALL of the following:

```markdown
### ID: Title

### Executive Summary
2-3 sentences. Vulnerability, root cause, worst-case impact.

### Details
Full code trace with snippets across ALL relevant layers.

### Impact Cascade
Bullet list: immediate → systemic consequences.

### Assumptions and Uncertainties
Numbered list. What must be true for the exploit to work.

### How will the bug recipient respond?
Anticipate the most likely dismissal and explain why it's insufficient.

### Why did tests miss this issue?
Specific gap in test coverage.

### Recommendation
Concrete code fix with snippets + "Additional hardening" section.

### References
Numbered list of exact file:line locations across all repos.
```

---

## Severity Classification

- **CRITICAL:** Chain halt, consensus failure, fund loss/inflation, key compromise, state corruption, value amplification
- **HIGH:** Exploitable under realistic conditions (permissionless abuse, reward manipulation, confidentiality breach)
- **MEDIUM:** Edge cases, defense-in-depth gaps, availability degradation
- **LOW:** Minor quality, documentation, theoretical long-timeline issues
- **INFO:** Observations, positive notes, suggestions
