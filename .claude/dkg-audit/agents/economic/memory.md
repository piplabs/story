# Agent Memory: Economic & Value Flow Auditor

## Codebase Knowledge
- NEW AGENT — No prior audit history. Bootstrapped from APEX findings analysis.

## Key Value Flows to Track
- **Staking bridge:** EL wei → CL stake (normalized by /1e9 gwei). Located in `evmstaking/keeper/keeper.go`
- **CDR fee bridge:** EL wei → CL stake via `ProcessCDRFeeCollected` → `AddCDRFeeToPool`. MISSING normalization (STOR-15)
- **CDR fee refund:** `RefundCDRFee` sends from cdr-fee-pool to validator account
- **CDR reward distribution:** `distributeCDRRewardPool` pays proportional to `CDRPartialSubmitCount`
- **DKG rewards:** `distributeRewardsFromModule` — correctly sorted

## APEX Findings in My Scope (reference material)
- **STOR-15:** CDR fee wei/gwei 1e9 amplification — unbacked stake minting
- **STOR-4:** Unknown/expired partials return nil → fee refund + reward increment
- **STOR-3:** Duplicate partial replays accrue fresh refunds and rewards
- **STOR-6:** Arbitrary submissions from registered validators rewarded without TDH2 verification

## Patterns to Always Check
1. Every EL→CL value crossing: is it normalized to gwei?
2. Every fee refund: what conditions? Can nil/not-found/expired trigger?
3. Every reward increment: can non-committee members trigger it?
4. Every payable function: permissionless drain test
5. Every nil-return in fee/reward path: trace through caller to economic consequence

## Blind Spots (from APEX)
- This agent didn't exist before. All APEX economic findings were gaps in our coverage.
