# Agent Memory Guide

## Purpose

Each audit agent maintains persistent memory across audit sessions. This allows agents to:
- Avoid repeating false positives from past audits
- Build deeper understanding of the codebase over time
- Track which areas have changed since last audit
- Learn from external audit findings (APEX, CDR report) that exposed gaps
- Refine their approach based on what worked and what didn't

## Directory Structure

```
agents/
├── AGENT_MEMORY_GUIDE.md        # This file
├── cosmos-sdk/
│   ├── memory.md                # Agent's persistent memory
│   └── findings-log.md          # Historical findings + outcomes
├── dkg-lifecycle/
├── smart-contract/
├── tee-sgx/
├── cryptography/
├── cross-cutting/
├── economic/
├── adversarial/
└── go-advisor/
```

## Per-Agent Files

### memory.md — Agent's Persistent Knowledge

Updated by the agent AFTER each audit. Contains:

```markdown
# Agent Memory: {Agent Name}

## Codebase Knowledge
Key architectural facts learned. File locations, patterns, invariants.
Only things that are NOT obvious from reading the code.

## Past False Positives
Things I flagged incorrectly and WHY they were wrong.
Prevents repeating the same mistake.

## Blind Spots Discovered
Areas I missed that external auditors found.
What I should look harder at next time.

## Effective Techniques
Approaches that produced real findings.
What to keep doing.

## Ineffective Techniques
Approaches that produced noise or missed real bugs.
What to stop doing.

## Codebase Change Watch
Files/patterns that were buggy before — check if they've been fixed
or if the same pattern has been introduced elsewhere.
```

### findings-log.md — Historical Findings + Outcomes

Append-only log. Each audit session appends its findings with eventual outcome:

```markdown
## Audit Session: {date} — {branches}

### {Finding ID}: {Title}
- **Severity:** {level}
- **Status:** TRUE_POSITIVE / FALSE_POSITIVE / DUPLICATE / FIXED_SINCE
- **Outcome notes:** What happened. If FP, why. If fixed, how.
```

## Agent Lifecycle Per Audit

1. **Before audit:** Agent reads its `memory.md` and `findings-log.md`
2. **During audit:** Agent uses memory to avoid known FPs, focus on known blind spots
3. **After audit:** Agent updates `memory.md` with new learnings, appends to `findings-log.md`

## Cross-Agent Learning

After re-validation phase, the orchestrator shares:
- Which findings from each agent were TRUE vs FALSE POSITIVE
- Which APEX/external findings map to each agent's scope
- This feedback is written into the relevant agent's `memory.md` under "Blind Spots"
