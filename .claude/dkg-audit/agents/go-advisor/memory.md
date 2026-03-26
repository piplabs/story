# Agent Memory: Go Architecture & Clean Code Advisor

## Codebase Knowledge
- NEW AGENT — No prior audit history.

## Known Structural Issues (from prior audits)
- **Package-level globals (H-05):** 8 global variables with 7 mutexes in `keeper.go:21-51`. Should be Keeper struct fields
- **Params dual registration (H-02):** collections.Item registered but raw KV used. Dead code path
- **dkg_service.go monolith:** Was split into per-RPC files (PR #6), good pattern to maintain

## Patterns to Evaluate
1. **Error handling consistency:** Are errors wrapped with context? Sentinel errors vs string matching?
2. **Concurrency model:** Package globals vs struct fields? Mutex granularity?
3. **Interface segregation:** Are interfaces minimal and consumer-defined?
4. **Testability:** Can services be tested without real gRPC/chain? Are dependencies injectable?
5. **Code duplication:** DKG vs CDR paths sharing similar logic?
6. **Context propagation:** Is context.Context threaded correctly? Any time.Sleep ignoring context?
7. **Naming conventions:** Consistent abbreviations? Clear intent?

## Story-Kernel Specific
- cb-mpc CGO integration: build complexity, linking patterns
- Gramine SGX constraints: no flock, limited syscalls
- gRPC service structure: single DKGServer with all RPCs

## Story CL Specific
- Cosmos SDK keeper pattern: keeper → msgserver → handler
- Vote extension handling: ABCI methods
- Module wiring: app.go registration
