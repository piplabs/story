# Story Test Engineer Memory

## DKG Keeper Package (`client/x/dkg/keeper/`)

### Test Setup Helpers
- `setupDKGKeeper(t)` → returns `(*Keeper, context.Context)` — uses `setupDKGKeeperWithMocks` internally
- `setupDKGKeeperWithMocks(t)` → also returns mock BankKeeper and DistributionKeeper for fine-grained control
- Mock controller is created via `go.uber.org/mock/gomock`; mocks live in `client/x/dkg/testutil/`
- `setupDealerRegistrationWithKey(t, k, ctx, codeCommitment, round, addr, index, dkgPubKeyBytes)` — registers a dealer with a real Edwards25519 pubkey

### Justification Test Helpers (`dkg_dealing_test.go`)
- `newDealerTestContext(t, n, threshold)` — generates Edwards25519 longterm key + VSS polynomial for n participants
- `dtc.makeSignedJustification(t, dealerIndex, kyberRecipientIndex)` — builds a valid Schnorr-signed justification (0-based kyber indices)
- `dtc.makeInvalidDealJustification(t, dealerIndex)` — valid signature but wrong share (for VSS failure testing)
- `dealingTestThreshold = 2` is the standard threshold constant used across justification tests

### Key Functions Under Test (`dkg_justification.go`)
- `verifyJustificationSignature(suite, j, dealerPubKeys)` — Schnorr sig check; all indices 0-based
- `deduplicateJustifications([]Justification)` — deduplicates by (dealerIndex, recipientIndex)
- `verifyJustification(network, j)` — Pedersen VSS check; internally converts 0-based `SecShare.I` to 1-based for `VerifyPedersenVSS`
- `buildDealerPubKeyMap(ctx, network, suite)` — builds `map[uint32]kyber.Point` from on-chain registrations
- `MaxJustificationsPerBlock = 10` constant in `dkg_justification.go`

### Pipeline Ordering Invariant
The three-step pipeline in `handleDKGProcessJustifications` MUST be:
1. Schnorr signature verification (before dedup — prevents attacker preempting valid justification with unsigned duplicate)
2. Deduplication by (dealerIndex, recipientIndex)
3. Pedersen VSS verification

### Coverage Notes (as of 2026-03-22, branch: dkg/add-unit-tests)
- **Overall package coverage: 83.8%** (improved from 80.2%)
- Functions at 100%: `callContractRegister`, `callContractFinalizeDKG`, `resolveRegistrationKernelClient`, `getRegistrationKernelClient`
- Functions with network deps at 0%: `ConnectAndDiscover`, `NewContractClient`, `CreateKernelClient` — require live gRPC/ETH RPC, skip unit test
- `PrepareVotes` 31.6% — main body needs `baseapp.ValidateVoteExtensions` with real valStore; test helpers directly instead
- `handleDKGRegistration` 78.6% — `session.Phase != PhaseInitializing` branch is dead code (NewDKGSession always returns PhaseInitializing)
- `GetCDRPartials` (query.go:137): tested with `setPartialDecryptionSubmission` direct call; success path covers grouped results and ThresholdMet flag
- `buildDealerPubKeyMap`: tested with empty DkgPubKey (skip branch) and invalid bytes (unmarshal-fail/skip branch)
- `InitiateDKGRound` with isDKGSvcEnabled + not registered: requires `initTestStateManager(t, k)` to avoid goroutine nil-stateManager panic
- `handleDKGProcessJustifications` itself stays at 0% — async goroutine handler requiring stateManager + teeClient mocks; pipeline logic tested via individual function tests
- Remaining low-coverage functions (< 80%): `PrepareVotes` 31.6%, `hasPendingUpgradeActivation` 75.0%, `distributeCDRRewardPool` 75.4%, `shouldReshare` 75.0%

### Critical Test Patterns (dkg_handler_threshold_test.go)
- Timeout path: `sdk.UnwrapSDKContext(baseCtx).WithBlockHeight(300)` — adjusts height while reusing same KV store
- uint64 underflow risk: if `current_height < request_height`, `current - request` wraps to huge number → always exceeds timeout
- `PartialDecryptionTimeoutBlocks = 200` — requests stored at height=H are cleaned up when `currentHeight - H > 200`
- Ciphertext key includes hash: "ciphertext mismatch" is actually a not-found path (different ciphertext → different key)
- 64-byte zero commPubKey + 65-byte zero signature → valid length triggers ECDSA verify path, fails correctly

### State Manager Initialization Pattern
- `stateManager` is nil by default in `setupDKGKeeperWithMocks` — only set by `InitDKGService`
- Use `initTestStateManager(t, k)` helper (creates `NewStateManager(t.TempDir())`) for tests needing session state
- `CreateSession(ctx, session) error` — returns error only (not session+error)
- After CreateSession, use `GetSession(ctx, round)` to retrieve the created session

### Testing Patterns Confirmed
- Table-driven tests used throughout `dkg_handler_internal_test.go`
- `t.Run(...)` sub-tests used to group pipeline scenarios in a single `Test*` function
- All test files are in `package keeper` (internal/white-box testing)
- `require.*` used (not `assert.*`) for all assertions

### Mock gRPC Client Patterns
- `MockKernelServiceClient.ProcessDeals/ProcessResponses/ProcessJustification/FinalizeDKG/GenerateAndSealKey` use variadic `opts ...grpc.CallOption` — pass 2 matchers `(gomock.Any(), gomock.Any())` without 3rd (variadic recorder accepts it)
- `DoAndReturn` for variadic gRPC mock: `func(_ context.Context, req *types.Foo, _ ...grpc.CallOption)` — must import `google.golang.org/grpc` in test file
- `retryAttemts = 3` (note: typo in source) — use this constant in `Times(retryAttemts)` to expect all retry attempts when kernel fails

### Race Condition Patterns
- `resumeFailedSession` spawns goroutines that modify `session` pointer and `dkgSvcRound` atomic
- Tests that call `resumeFailedSession` for DKGStageActive/DKGStageFinalization MUST NOT be `t.Parallel()` AND must use `resetDKGSvcRound()` + `defer resetDKGSvcRound()`
- For ActiveStage: goroutine calls `handleDKGComplete` which calls `tryAcquireDKGSvc` — poll `dkgSvcRound.Load() == 0` to wait for goroutine completion before inspecting session state
- StateManager returns cached pointers; concurrent goroutine modification of session is a real race; always wait for goroutine to finish before reading session state

### types.Justification Struct Fields
- `Signature` field is on `VSSJustification`, NOT on `Justification` directly
- `SecShare.V` is `*types.Scalar`, not `[]byte`
- `PlainDeal` has `Commitments []*Point`, not `Commits [][]byte`
