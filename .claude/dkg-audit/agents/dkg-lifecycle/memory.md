# Agent Memory: DKG Service & Lifecycle Auditor

## Codebase Knowledge
- DKG service entry: `client/x/dkg/keeper/dkg_svc.go` — retry logic, goroutine management
- Kernel communication: `kernel_router.go` — gRPC connection management, reconnection
- CDR decrypt flow: `dkg_handler.go` — PartialDecryptionSubmitted, decrypt request registry
- CDR fee flow: `dkg_cdr_fees.go` — AddCDRFeeToPool, RefundCDRFee, distributeCDRRewardPool
- State machine: registration → dealing → response → finalization → active → (resharing)
- Vote extension queues: package-level globals in `keeper.go:21-51`

## Past False Positives
- DKG-SVC-002: flagged as issue but was design-documented behavior
- DKG-SVC-012: flagged dedup key as wrong but it was correct

## Blind Spots Discovered
- **nil-return = success to caller (STOR-4):** `PartialDecryptionSubmitted` returns nil for not-found/expired requests. Caller `ProcessDKGPartialDecryptionSubmitted` interprets nil as success → refunds fee + increments reward count. I checked the function's error handling but NEVER traced what the CALLER does with nil
- **CDR dead windows (STOR-9):** CDR read path is dead during non-active DKG phases. Decrypt worker is killed during DKG rounds, and resharing rejects request queueing. I never mapped CDR availability against DKG state
- **PIDCache lifecycle (STOR-8):** PIDCache not populated during resharing → post-resharing CDR permanently broken. I checked caches but not their population across round transitions
- **Duplicate replay rewards (STOR-3):** Same partial can be submitted multiple times and accrue rewards each time

## Effective Techniques
- Tracing goroutine lifecycle and context propagation — found C-04 (context leak)
- Checking kernel reconnection behavior

## Ineffective Techniques
- Checking only "does this error path handle errors?" without asking "what economic consequence does this error path have?"
- Assuming callers of internal functions validate their inputs
