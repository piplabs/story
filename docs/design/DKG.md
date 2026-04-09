# Story DKG System Architecture

## Glossary

| Term | Definition |
|------|-----------|
| **DKG** | Distributed Key Generation — protocol where N parties jointly compute a shared public key and individual private key shares, without any single party learning the full private key |
| **TDH2** | Threshold Decryption based on Diffie-Hellman (2-round variant) — the threshold encryption/decryption scheme used by cb-mpc |
| **TEE** | Trusted Execution Environment — hardware-isolated execution context (Intel SGX in Story's case) |
| **SGX** | Software Guard Extensions — Intel's TEE technology providing memory encryption and remote attestation |
| **MRENCLAVE** | A 32-byte hash of the SGX enclave's code and initial data; acts as the "code commitment" |
| **VSS** | Verifiable Secret Sharing — Pedersen VSS is used to share secrets with cryptographic proofs of correctness |
| **Resharing** | Process of transitioning the DKG key shares to a new committee while preserving the same global public key |
| **cb-mpc** | Coinbase's MPC (Multi-Party Computation) library — provides TDH2 encryption/decryption primitives |
| **kyber** | Go cryptographic library (DEDIS) — provides Edwards25519 curve operations, Pedersen DKG, and Schnorr signatures |

## Table of Contents

1. [System Overview](#1-system-overview)
2. [DKG Lifecycle](#2-dkg-lifecycle)
3. [Reward Distribution](#3-reward-distribution)
4. [Kernel Upgrade Workflow](#4-kernel-upgrade-workflow)
5. [Key Technical Details](#5-key-technical-details)
6. [story-kernel (TEE) Architecture](#6-story-kernel-tee-architecture)

> For CDR (Confidential Data Rails) documentation, see [CDR_DESIGN_DOC.md](CDR).

---

## 1. System Overview

### 1.1 Purpose

The DKG (Distributed Key Generation) and CDR (Confidential Data Rails) system enables Story blockchain validators to collectively generate a shared cryptographic key within TEE (Trusted Execution Environment) enclaves. This shared key enables threshold decryption of confidential data without any single validator having access to the full key.

### 1.2 High-Level Architecture

```
+----------------------------------------------------------------------+
|                         External Users                               |
|   (Allocate vault, encrypt with GlobalPubKey, write, read to decrypt)|
+-------------------+-----------------------------+--------------------+
                    |                             |
                    v                             v
+------------------------+----+    +--------------+----------------+
|    DKG Contract             |    |      CDR Contract             |
|    (predeploy)              |    |      (predeploy)              |
|                             |    |                               |
| Primary: event emission     |    | - allocate() -> VaultAllocated|
| for deterministic CL        |    | - write()    -> VaultWritten  |
| processing                  |    | - read()     -> VaultRead     |
|                             |    | - submitEncrypted-            |
| Validates: code commitment  |    |   PartialDecryption()         |
| via enclave report auth     |    |                               |
|                             |    |                               |
| - register()  (+ SGX auth) |    |                               |
| - finalize()  (+ whitelist) |    |                               |
| - scheduleUpgrade()        |    |                               |
| - cancelUpgrade()          |    |                               |
+----------+------------------+    +----------+-------------------+
           |                                   |
           |  Events: Registered,              |  Events: VaultAllocated,
           |  Finalized, UpgradeScheduled,     |  VaultRead,
           |  UpgradeCancelled                 |  EncryptedPartialDecryption-
           |                                   |  Submitted
           v                                    v
+----------+------------------------------------+-----------+
|                 evmengine Module                          |
|                                                           |
|  ProcessDKGEvents()      ProcessCDREvents()               |
|  - ParseRegistered       - ParseVaultRead                 |
|  - ParseFinalized        - ParsePartialDecryptSubmitted   |
|  - ParseUpgradeScheduled                                  |
+-----------------------+-----------------------------------+
                        |
                        v
+-----------------------+-----------------------------------+
|                   DKG Module (x/dkg)                      |
|                                                           |
|  BeginBlocker:                                            |
|  - Stage transitions (Reg -> Deal -> Final -> Active)     |
|  - Kernel upgrade activation                              |
|  - Decrypt request pruning                                |
|                                                           |
|  Vote Extensions:                                         |
|  - ExtendVote (dequeue deals/responses/justifications)    |
|  - VerifyVoteExtension (parse + size checks)              |
|  - PrepareVotes -> MsgAddDkgVote                          |
|  - AddVote (process deals/responses/justifications)       |
|                                                           |
|  On-chain State:                                          |
|  - DKGNetworks, LatestDKGNetwork, LatestActiveRound       |
|  - DKGRegistrations, GlobalPubKeyVotes                    |
|  - KernelUpgradeInfos, DecryptRequestRegistry             |
|  - DKGPartialDecrypt, SettlementBalance                   |
|                                                           |
|  Off-chain State (StateManager):                          |
|  - DKGSessions (JSON on disk), session phases             |
+-----------------------+-----------------------------------+
                        |
                        |  gRPC (per-validator)
                        v
+-----------------------+-----------------------------------+
|              story-kernel (TEE/SGX Enclave)               |
|                                                           |
|  gRPC Service:                                            |
|  - GetCodeCommitment                                      |
|  - GenerateAndSealKey                                     |
|  - GenerateDeals                                          |
|  - ProcessDeals -> responses                              |
|  - ProcessResponses -> justifications                     |
|  - ProcessJustification                                   |
|  - FinalizeDKG -> global_pub_key + signature              |
|  - PartialDecryptTDH2 -> encrypted_partial_decryption     |
|                                                           |
|  Security Boundary:                                       |
|  - SGX sealing (keys never leave enclave plaintext)       |
|  - MRENCLAVE (code commitment) verification               |
|  - SGX remote attestation (quotes)                        |
+-----------------------------------------------------------+
```

### 1.3 Component Roles

| Component | Role                                                                                                                                                                                                                                                    | Trust Level |
|-----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------|
| **DKG Contract** | Entrypoint whose primary role is **event emission** so CL can process state transitions deterministically. Validates enclave reports against whitelisted code commitments via `IAttestationReportValidator`; all other DKG state logic is handled by CL | On-chain (deterministic, publicly verifiable) |
| **CDR Contract** | Manages encrypted vault storage and read access control; emits events for CL processing. Partial decryptions are stored on CL (queried via `GetCDRPartials` RPC), and TDH2 combination is performed off-chain by the requester                          | On-chain (deterministic) |
| **evmengine Module** | Bridges EVM log events to Cosmos SDK module calls                                                                                                                                                                                                       | On-chain (consensus-critical) |
| **DKG Module (x/dkg)** | Orchestrates the DKG lifecycle via BeginBlocker stage transitions, vote extensions for data propagation, and on-chain state management                                                                                                                  | On-chain (consensus-critical) |
| **StateManager** | Manages off-chain DKG session state (JSON files on disk); NOT part of consensus                                                                                                                                                                         | Off-chain (per-validator) |
| **KernelRouter** | Routes gRPC calls to the correct story-kernel binary by code commitment; supports up to 2 endpoints (old + new for upgrades)                                                                                                                            | Off-chain (per-validator) |
| **story-kernel** | Runs inside SGX enclave; generates keys, deals, responses; performs threshold decryption; all private keys sealed by SGX                                                                                                                                | TEE-protected |

### 1.4 Interface Reference

#### story-kernel gRPC Service (`KernelService` in `proto/tee.proto`)

| RPC Method | Request | Response | Description |
|------------|---------|----------|-------------|
| `GetCodeCommitment` | `GetCodeCommitmentRequest` | `GetCodeCommitmentResponse` | Return enclave's MRENCLAVE |
| `GenerateAndSealKey` | `GenerateAndSealKeyRequest` | `GenerateAndSealKeyResponse` | Generate + SGX-seal Ed25519 & secp256k1 key pairs, return pub keys + SGX quote |
| `GenerateDeals` | `GenerateDealsRequest` | `GenerateDealsResponse` | Generate Pedersen VSS deals for each recipient |
| `ProcessDeals` | `ProcessDealsRequest` | `ProcessDealsResponse` | Verify and process received deals, return responses (approve/complaint) |
| `ProcessResponses` | `ProcessResponsesRequest` | `ProcessResponsesResponse` | Process responses, return justifications if complaints exist |
| `ProcessJustification` | `ProcessJustificationRequest` | `ProcessJustificationResponse` | Process justifications to resolve complaints |
| `FinalizeDKG` | `FinalizeDKGRequest` | `FinalizeDKGResponse` | Compute DistKeyShare, return globalPubKey + pubKeyShare + ECDSA signature |
| `PartialDecryptTDH2` | `PartialDecryptTDH2Request` | `PartialDecryptTDH2Response` | Perform partial decryption, encrypt result to requester via ECDH+AES-GCM |

#### DKG Contract (predeploy)

> The DKG contract's primary purpose is **event emission** — it emits structured events so that the CL (DKG module) can process DKG state transitions deterministically across all validators. The contract itself does NOT manage DKG rounds, stages, or participant state. However, it **does validate code commitments**: `register()` authenticates the SGX enclave report against the whitelisted `EnclaveTypeData.codeCommitment` via `IAttestationReportValidator`, and `finalize()` verifies the enclave type is whitelisted.

| Type | Name | Contract Validation | Event Emitted |
|------|------|---------------------|---------------|
| Method | `register()` | Input validation + `_authenticateEnclaveReport()`: verifies SGX quote against whitelisted `codeCommitment` via `IAttestationReportValidator.validateReport()` | `Registered` (enclaveReport, round, validator, enclaveType, codeCommitment, keys, startBlock) |
| Method | `finalize()` | Input validation + `isEnclaveTypeWhitelisted[enclaveType]` check | `Finalized` (round, validator, enclaveType, codeCommitment, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare, signature) |
| Method | `scheduleUpgrade()` | `activationHeight > block.number` | `UpgradeScheduled` (activationHeight, upgradeVersion) |
| Method | `cancelUpgrade()` | upgradeVersion non-empty | `UpgradeCancelled` (upgradeVersion) |
| Admin | `whitelistEnclaveType()` | Owner-only; sets `EnclaveTypeData{codeCommitment, validationHookAddr}` | — |
| Query | `EnclaveTypeData()` | — | — |
| Query | `Fee()` | — | — |

#### CDR Contract (predeploy)

| Type | Name | Description |
|------|------|-------------|
| Method | `allocate()` | Allocate a new vault, returns `uuid`; collects `allocateFee` |
| Method | `write()` | Store encrypted data in an allocated vault; collects `writeFee` |
| Method | `read()` | Request decryption of vault data; collects `readFee` |
| Method | `submitEncryptedPartialDecryption()` | Submit encrypted partial decryption with signature; collects `baseFee` (refunded to submitter on successful CL verification) |
| Event | `VaultAllocated` | Emitted on vault allocation (uuid, conditions) |
| Event | `VaultWritten` | Emitted on vault write |
| Event | `VaultRead` | Emitted on vault read (triggers decryption flow) |
| Event | `EncryptedPartialDecryptionSubmitted` | Emitted when partial decryption is submitted |

#### Story CL — DKG Keeper Service Handlers

| Handler | Calls Kernel | Calls Contract | Phase Transition |
|---------|-------------|----------------|-----------------|
| `handleDKGRegistration()` | `GenerateAndSealKey` | `Register()` (queries `Fee()`, sends value) | → PhaseInitialized |
| `handleDKGDealing()` | `GenerateDeals` | — | EnqueueDeals |
| `handleDKGProcessDeals()` | `ProcessDeals` | — | EnqueueResponses |
| `handleDKGProcessResponses()` | `ProcessResponses` | — | EnqueueJustifications |
| `handleDKGProcessJustifications()` | `ProcessJustification` | — | — |
| `handleDKGFinalization()` | `FinalizeDKG` | `Finalize()` (queries `Fee()`, sends value) | → PhaseFinalized |
| `handleDKGComplete()` | — | — | → PhaseCompleted |
| `handleDecryptRequest()` | `PartialDecryptTDH2` | `SubmitEncryptedPartialDecryption()` | — |

#### KernelRouter

| Method | Description |
|--------|-------------|
| `ConnectAndDiscover()` | Connect to endpoint, call `GetCodeCommitment`, register client |
| `GetClient(codeCommitmentHex)` | Get `KernelServiceClient` by code commitment |
| `GetAllCodeCommitments()` | List all connected code commitments |
| `Disconnect(codeCommitmentHex)` | Close and remove a client connection |

### 1.5 Verification Layers

```
+-------------------------------------------------------------------+
|  CONSENSUS LAYER (deterministic, all nodes agree)                 |
|                                                                   |
|  On-chain state: DKGNetworks, Registrations, GlobalPubKeyVotes,  |
|  KernelUpgradeInfos, DecryptRequestRegistry, Params              |
|                                                                   |
|  Stage transitions in BeginBlocker                                |
|  Event processing via evmengine (cached context)                  |
|  Vote extension aggregation (MsgAddDkgVote)                       |
+-----------------------------------+-------------------------------+
                                    |
              ---- verification boundary ----
                                    |
+-----------------------------------+-------------------------------+
|  OFF-CHAIN LAYER (per-validator, may differ)                      |
|                                                                   |
|  StateManager: session phases, decrypt request queue              |
|  KernelRouter: gRPC connections to story-kernel                   |
|  Async goroutines: handleDKGRegistration, handleDKGDealing, etc.  |
+-----------------------------------+-------------------------------+
                                    |
               ---- TEE boundary ----
                                    |
+-----------------------------------+-------------------------------+
|  TEE LAYER (hardware-protected)                                   |
|                                                                   |
|  SGX-sealed Ed25519 keys (DKG)                                    |
|  SGX-sealed secp256k1 keys (communication/signing)                |
|  SGX-sealed DistKeyShare (threshold decryption)                   |
|  MRENCLAVE-bound code commitment                                  |
+-------------------------------------------------------------------+
```

---

## 2. DKG Lifecycle

### 2.1 Round Stage State Machine

Each DKG round progresses through a sequence of time-bounded stages driven by `BeginBlocker`. Stage transitions are determined by block height elapsed since the round's `StartBlockHeight`.

```
                             elapsed >= RegistrationPeriod
  +-----------------+  ---------------------------------------->  +-------------+
  |  Registration   |                                             |   Dealing   |
  |  (Stage 1)      |                                             |  (Stage 2)  |
  +-----------------+                                             +------+------+
         ^                                                               |
         |                        elapsed >= DealingPeriod               |
         |  SkipToNextRound       <--------------------------------------+
         |  (on failure)                                                 |
         |                                                               v
         |                                                       +------+--------+
         +-------------------------------------------------------|  Finalization  |
         |                                                       |   (Stage 3)    |
         |                   elapsed >= FinalizationPeriod        +------+---------+
         |                                                               |
         |                                                               v
         |                                                       +------+------+
         |                                                       |   Active    |
         |                                                       |  (Stage 4)  |
         |                                                       +------+------+
         |                                                               |
         |                   elapsed >= ActivePeriod                     |
         +---------------------------------------------------------------+
                    (initiates next round with resharing)

  Active round N stays Active until round N+1 completes:
  FinalizeDKGRound(N+1) -> endPreviousActiveRound(N) -> Stage=Ended


  Special stages:
  +--------+     +-------+
  | Failed |     | Ended |  (set on previous round when NEXT round completes)
  +--------+     +-------+
```

**Stage Timing** (`shouldTransitionStage`):
```
elapsed = currentHeight - dkgNetwork.StartBlockHeight

Registration ends at:    elapsed >= RegistrationPeriod
Dealing ends at:         elapsed >= RegistrationPeriod + DealingPeriod
Finalization ends at:    elapsed >= RegistrationPeriod + DealingPeriod + FinalizationPeriod

Note: Active -> Ended is NOT a self-transition. A round stays Active until the
NEXT round reaches FinalizeDKGRound(), which calls endPreviousActiveRound() to
set the previous round's stage to Ended. In other words, round N's stage becomes
Ended only when round N+1 successfully completes.
```

### 2.2 Session Phase State Machine (Off-Chain)

Each validator maintains a local `DKGSession` with its own phase tracking. This is independent of the on-chain stage and tracks the validator's own progress through TEE operations.

```
  PhaseInitializing --> PhaseInitialized --> PhaseDealing --> PhaseFinalized --> PhaseCompleted
        |                    |                  |                 |
        +--------------------+------------------+-----------------+
                             |
                             v
                       PhaseFailed  -----> (ResumeDKGService recovers)
```

### 2.3 V1.6.0 Activation & BeginBlocker

The DKG module activates at the v1.6.0 upgrade height. The upgrade handler registers the DKG store, and `BeginBlocker` starts managing DKG rounds from that point.

**BeginBlocker flow**:

```
BeginBlocker(ctx)
  |
  +-- GetLatestDKGRound()
       |     |
       |     +-- nil? --> InitiateDKGRound(isUpgrade=false)  [first round ever]
       |
       +-- hasPendingUpgradeActivation(currentHeight)?
       |     |
       |     +-- Yes --> mark IsActivated=true, InitiateDKGRound(isUpgrade=true)
       |
       +-- pruneTimedOutDecryptRequests (every N blocks)
       |
       +-- shouldTransitionStage(currentHeight, latestRound, params)?
             |
             +-- Registration --> InitiateDKGRound(isUpgrade=false)  [resharing]
             +-- Dealing       --> BeginDealing(latestRound)
             +-- Finalization  --> BeginFinalization(latestRound)
             +-- Active        --> FinalizeDKGRound(latestRound)
```

### 2.4 Round Initiation

**`InitiateDKGRound(ctx, isUpgrade)`** creates a new `DKGNetwork` in state:

1. Fetch active (bonded, non-jailed) validators as EVM addresses
2. Determine next round number (`latestRound + 1`, or `1` if first)
3. Check `shouldReshare`: true if a previous active round exists
4. Create `DKGNetwork{round, start_block_height, start_block_hash, active_val_set, stage: DKG_STAGE_REGISTRATION, is_resharing, is_upgrade}`
5. Store via `setDKGNetwork` (updates `LatestDKGNetwork` pointer)
6. Emit `EventBeginInitialization`
7. If `isDKGSvcEnabled`, launch async `handleDKGRegistration` goroutine

**First round vs. resharing vs. upgrade**:

| Scenario | IsResharing | IsUpgrade | Who Deals | Who Processes |
|----------|-------------|-----------|-----------|---------------|
| First round (round 1) | false | false | Current set | Current set |
| Periodic resharing | true | false | Previous active set | Both sets |
| Kernel upgrade | true | true | Previous active set (old binary) | Both sets (both binaries) |

### 2.5 Registration Phase

The registration phase is where each validator generates cryptographic keys inside their TEE and registers them on-chain.

```
                                    Story CL (Validator)
                                          |
                                    handleDKGRegistration()
                                          |
                           +--------------+--------------+
                           |                             |
                    1. Create DKG Session         (old-only members:
                       (StateManager)              validators in prev set
                           |                       but NOT in new set —
                           |                       skip key gen, only
                           |                       generate deals during
                           |                       Dealing phase using
                           |                       their existing key share)
                           |
                    2. getRegistrationKernelClient()
                       -> returns (client, clientCC, error)
                       -> session.CodeCommitment = clientCC
                           |
                    3. callTEEGenerateAndSealKey()
                           |
                           v
                    +------+------+
                    | story-kernel |
                    | (SGX enclave)|
                    +------+------+
                           |
                    GenerateAndSealKey():
                    - Generate Ed25519 key pair (DKG)
                    - Generate secp256k1 key pair (signing)
                    - Seal both keys via SGX
                    - Verify DKG start block on canonical chain
                    - Generate SGX remote attestation quote
                    - Return: dkgPubKey, commPubKey, enclaveReport
                           |
                           v
                    4. callContractRegister() (queries Fee(), sends value)
                       DKG.register(enclaveReport, instanceData, startBlock)
                           |
                           v
                    +------+------+
                    | DKG Contract |  (validates code commitment, emits event)
                    +------+------+
                           |
                    Contract validation:
                    - Input validation (non-empty fields, round != 0)
                    - _authenticateEnclaveReport():
                      Lookup whitelisted EnclaveTypeData by enclaveType
                      -> IAttestationReportValidator.validateReport(
                           expectedCodeCommitment,  // from whitelist
                           keccak256(instanceData), // data commitment
                           enclaveReport,           // SGX quote
                           validationContext)
                      -> Verifies SGX quote signature + code commitment match
                           |
                    Emits: Registered event (includes codeCommitment from whitelist)
                           |
                           v
                    evmengine.ProcessDKGRegistered()
                           |
                           v
                    dkgKeeper.Registered()  (CL handles DKG state)
                    - Validate round, stage, startBlock match
                    - Verify validator is in ActiveValSet
                    - Assign 1-based index
                    - Store DKGRegistration{status: DKG_REG_STATUS_VERIFIED}
```

**Key details**:
- DKG public key: Ed25519 (Edwards25519, used in kyber Pedersen DKG)
- Communication public key: secp256k1 ECDSA (used for finalization/partial-decrypt signatures)
- Enclave report: SGX remote attestation quote binding `(address, round, edPub, secpPub, startBlockHeight, startBlockHash)`
- Registration index: 1-based on-chain. When interfacing with kyber (0-based PIDs), subtract 1.

### 2.6 Dealing Phase

When `elapsed >= RegistrationPeriod`, the stage transitions to Dealing.

**`BeginDealing(ctx, latestRound)`**:

1. Count verified registrations; if `< MinReqRegisteredParticipants`, skip to next round
2. Set `Total` and `Threshold` on the `DKGNetwork` based on actual verified count and `OperationalThreshold`
3. Pre-compute `shouldDeal` and `ensureSessionIndex` (requires SDK context)
4. Emit `EventBeginDealing`
5. Launch async `handleDKGDealing` goroutine (with pre-computed values)

**Deal generation flow**:

```
handleDKGDealing()
  |
  +-- tryAcquireDKGSvc (prevent duplicate goroutines)
  +-- Check stage == Dealing
  +-- Check shouldDeal (first round: current set; resharing: previous set)
  +-- Get session, verify phase == Initialized
  |
  +-- For upgrade: use old binary's CC for dealing
  |   (old key shares sealed by old binary)
  |
  +-- kernelRouter.GetClient(dealerCC)
  +-- client.GenerateDeals(codeCommitment, round, isResharing)
       |
       v
  story-kernel.GenerateDeals():
    - GetOrLoadRoundContext (network + registrations)
    - CachePID (for non-resharing rounds)
    - GetInitDKG or GetResharingPrevDKG (create/load DistKeyGenerator)
    - distKeyGen.Deals() -> encrypted Pedersen VSS deals
       |
       v
  EnqueueDeals(deals) --> deals queue --> ExtendVote --> broadcast
```

**Deal processing flow** (received via vote extensions):

```
MsgAddDkgVote.AddVote()  [stage == Dealing only]
  |
  +-- ProcessDeals(latestRound, deals)
  |     +-- Pre-compute shouldDeal (SDK context)
  |     +-- Pre-compute ensureSessionIndex (SDK context)
  |     +-- Emit EventBeginProcessDeals
  |     +-- async: handleDKGProcessDeals()
  |           +-- Accepts phase == PhaseInitialized OR PhaseDealing
  |           +-- dkgKernelMu.Lock() (serialize kernel ops; locked in handler)
  |           +-- Filter deals where RecipientIndex == session.Index - 1
  |           +-- client.ProcessDeals(filteredDeals, isResharing)
  |           +-- EnqueueResponses(responses) --> broadcast
  |
  +-- ProcessResponses(latestRound, responses)
  |     +-- Pre-compute shouldProcessResponses (SDK context)
  |     +-- async: handleDKGProcessResponses()
  |           +-- Accepts phase == PhaseInitialized OR PhaseDealing
  |           +-- dkgKernelMu.Lock() (locked in handler)
  |           +-- Filter out self-responses
  |           +-- For upgrade: send to BOTH old and new binary
  |           +-- client.ProcessResponses(filteredResponses)
  |           +-- EnqueueJustifications(justifications) --> broadcast
  |
  +-- ProcessJustifications(latestRound, justifications)
        +-- Pre-compute buildDealerPubKeyMap (SDK context)
        +-- async: handleDKGProcessJustifications()
              +-- Accepts phase == PhaseInitialized OR PhaseDealing
              +-- dkgKernelMu.Lock() (locked in handler)
              +-- 1. Schnorr signature verification
              +-- 2. Deduplication by (dealerIndex, recipientIndex)
              +-- 3. Pedersen VSS verification
              +-- 4. Forward valid justifications to story-kernel
              +-- For upgrade: send to BOTH old and new binary
```

### 2.7 Vote Extensions

DKG data (deals, responses, justifications) propagates between validators through CometBFT vote extensions. This is necessary because deal data is too large for standard Cosmos SDK transactions and must be broadcast alongside consensus votes.

```
                    Block N-1 (Proposer P)
                           |
                    ExtendVote() [each validator]
                    - DequeueDeals(maxItemsPerVote=80)
                    - DequeueResponses(80)
                    - DequeueJustifications(80)
                    - Marshal Vote proto
                    - Return as VoteExtension bytes
                           |
                    VerifyVoteExtension() [each validator]
                    - Size check (<= 256KB)
                    - Proto unmarshal
                    - Item count checks (<= 80 each)
                    - Return ACCEPT or REJECT
                           |
                    Block N (Proposer Q)
                    PrepareProposal:
                      PrepareVotes(LocalLastCommit)
                      - Skip first 2 blocks after v1.6.0
                      - ValidateVoteExtensions
                      - Parse, verify, discard invalid
                      - aggregateVotes (merge + deduplicate)
                      - Return MsgAddDkgVote{votes}
                           |
                    ProcessProposal:
                      AddVote(MsgAddDkgVote)
                      - Check authority
                      - If stage == Dealing:
                        - ProcessDeals
                        - ProcessResponses
                        - ProcessJustifications
```

**Deduplication keys**:
- Deals: `(dealerIndex, recipientIndex)`
- Responses: `(responderIndex, dealerIndex)`
- Justifications: `(dealerIndex, recipientIndex)`

### 2.8 Finalization Phase

When `elapsed >= DealingPeriod`, the stage transitions to Finalization.

**`BeginFinalization(ctx, latestRound)`** triggers each validator to finalize its DKG locally:

```
handleDKGFinalization()
  |
  +-- Check stage == Finalization, validator in current set
  +-- Get session, verify phase == Dealing
  |
  +-- callTEEFinalizeDKG():
  |     client.FinalizeDKG(codeCommitment, round, isResharing)
  |       |
  |       v
  |     story-kernel.FinalizeDKG():
  |     - distKeyGen.DistKeyShare() -> private share + global public key
  |     - SealAndStoreDistKeyShare (SGX sealed)
  |     - Cache DistKeyShare
  |     - Calculate participantsRoot = keccak256(sorted validator addresses)
  |     - Hash finalization data:
  |       hash = keccak256(codeCommitment || round(4B) || participantsRoot ||
  |                        globalPubKey || publicCoeffs... || pubKeyShare)
  |     - ethHash = keccak256("\x19Ethereum Signed Message:\n32" || hash)
  |     - Sign with sealed secp256k1 key
  |     - Return: participantsRoot, globalPubKey, publicCoeffs, pubKeyShare, signature
  |
  +-- callContractFinalizeDKG() (queries Fee(), sends value):
        DKG.finalize(round, validator, enclaveType, participantsRoot,
                     globalPubKey, publicCoeffs, pubKeyShare, signature)
          |
          v
        +------+------+
        | DKG Contract |  (validates enclave type whitelist, emits event)
        +------+------+
          |
        Contract validation:
        - Input validation (non-empty fields)
        - isEnclaveTypeWhitelisted[enclaveType] check
        - Extracts codeCommitment from EnclaveTypeData
          |
        Emits: Finalized event (includes codeCommitment from whitelist)
          |
          v
        evmengine.ProcessDKGFinalized()
          |
          v
        dkgKeeper.Finalized():  (CL handles DKG state)
        - Validate round, stage, registration status
        - Prevent double-finalization
        - Reject invalidated dealers
        - validateParticipantsRoot:
            includes both DKG_REG_STATUS_VERIFIED and DKG_REG_STATUS_FINALIZED
            registrations (keccak256 of sorted addresses)
        - verifyFinalizationSignature (ECDSA recovery against commPubKey)
        - AddGlobalPubKeyVote(round, globalPubKey, publicCoeffs)
        - If voteCount >= threshold AND globalPubKey not yet set:
            Set DKGNetwork.global_public_key = global_pub_key
        - finalizeDKGRegistration(round, validator, pubKeyShare)
            status: DKG_REG_STATUS_VERIFIED -> DKG_REG_STATUS_FINALIZED
```

**GlobalPubKey threshold voting**: Each validator independently computes the GlobalPubKey via kyber's `DistKeyShare.Public()`. They submit it on-chain via the finalize transaction. The DKG module counts votes keyed by `round_globalPubKeyHex_hash(publicCoeffs)`. When votes reach the threshold, the GlobalPubKey is committed to the DKGNetwork state.

### 2.9 Active Phase

When `elapsed >= FinalizationPeriod`, the stage transitions to Active via `FinalizeDKGRound`.

**`FinalizeDKGRound(ctx, latestRound)`**:

1. Count finalized registrations; skip to next round if below minimums (`MinReqFinalizedParticipants` or `Threshold`)
2. Emit `EventDKGFinalized`
3. `settleRewardsForPreviousCommittee()` - withdraw all UBI from distribution module to DKG module, distribute `DkgCommitteeRewardPortion` to previous committee members, store remainder as `SettlementBalance` (later withdrawn to `UbiWithdrawAddress` via evmstaking `ProcessUbiWithdrawal`)
4. `endPreviousActiveRound()` - set previous active round's stage to `Ended`
5. `setLatestActiveRound(latestRound)` - update active round pointer
6. If upgrade round: `deleteActivatedUpgradeInfo()`
7. Launch async `handleDKGComplete()`:
   - Mark session as `PhaseCompleted`, `IsFinalized = true`
   - `StartDecryptWorker()` - begin processing CDR decryption requests

### 2.10 Failed Rounds

**`SkipToNextRound(ctx, currentRound)`**:
1. `FlushAllQueues()` - clear all deal/response/justification queues
2. Preserve `isUpgrade` flag for retry
3. Mark current round as `DKG_STAGE_FAILED`
4. `InitiateDKGRound(ctx, isUpgrade)` - start a new round

**Failure triggers**:
- Verified registration count < `MinReqRegisteredParticipants` (at BeginDealing)
- Finalized count < `MinReqFinalizedParticipants` (at FinalizeDKGRound)
- Finalized count < `Threshold` (at FinalizeDKGRound)

**Recovery via `ResumeDKGService`**:
When a session enters `PhaseFailed`, `ResumeDKGService` can restart it by rewinding the phase to the appropriate state for the current on-chain stage:
- Registration stage: rewind to `PhaseInitializing`, retry `handleDKGRegistration`
- Dealing stage: rewind to `PhaseInitialized`, retry `handleDKGDealing`
- Finalization stage: rewind to `PhaseDealing`, retry `handleDKGFinalization`
- Active stage: rewind to `PhaseFinalized`, retry `handleDKGComplete`

---

## 3. Reward Distribution

### 3.1 DKG Committee Rewards from UBI Pool

DKG committee members receive a portion of UBI (Universal Basic Income) pool rewards. There are two distribution paths depending on timing:

**Path 1: Periodic distribution to CURRENT committee** (via evmstaking EndBlock):

Every UBI withdrawal cycle, evmstaking distributes a portion of the withdrawn UBI to the currently active DKG committee:

```
evmstaking.EndBlock()
  |
  +-- ProcessUbiWithdrawal()
        |
        +-- ClaimSettlementBalance() --> settlementAmount
        |     (claim leftover from previous round transition)
        |
        +-- WithdrawUbiByDenomToModule() --> withdrawnAmount
        |
        +-- DistributeRewardsToActiveCommittee(evmstaking, withdrawnAmount)
        |     |
        |     +-- Get latest active DKG network
        |     +-- Get finalized committee members
        |     +-- dkgReward = DkgCommitteeRewardPortion * withdrawnAmount
        |     +-- perMemberReward = dkgReward / memberCount
        |     +-- Transfer from evmstaking -> DKG module -> each member
        |     +-- Return totalDistributed
        |
        +-- withdrawnAmount -= totalDistributed
        +-- totalToWithdraw = withdrawnAmount + settlementAmount
        +-- BurnCoins(totalToWithdraw)
        +-- AddWithdrawalToQueue(UbiWithdrawAddress, totalToWithdraw)
```

**Path 2: Settlement at round transition for PREVIOUS committee** (during FinalizeDKGRound):

When a new round finalizes, the outgoing committee's accumulated UBI is settled. The DKG portion goes to previous members; the remainder is stored as `SettlementBalance` for evmstaking to claim in the next EndBlock:

```
FinalizeDKGRound()
  |
  +-- settleRewardsForPreviousCommittee()
        |
        +-- Get previous active DKG network (still pointed to before update)
        +-- WithdrawAllUBI -> DKG module account
        +-- dkgReward = DkgCommitteeRewardPortion * withdrawnAmount
        +-- perMemberReward = dkgReward / memberCount
        +-- Send to each finalized member of previous committee
        +-- remaining = withdrawnAmount - totalDistributed
        +-- SettlementBalance.Set(remaining)
```

### 3.2 Settlement Balance

The settlement balance holds the non-DKG portion of UBI that was withdrawn during `settleRewardsForPreviousCommittee`. It is claimed by evmstaking in the next EndBlock and withdrawn to `UbiWithdrawAddress`:

```
evmstaking.ProcessUbiWithdrawal()
  |
  +-- ClaimSettlementBalance(evmstakingModule)
        |
        +-- Read SettlementBalance
        +-- Transfer from DKG module -> evmstaking
        +-- Remove SettlementBalance entry
        +-- Return claimed amount (added to totalToWithdraw)
```

### 3.3 Reward Distribution Summary

```
UBI Pool (distribution module)
  |
  +-- Path 1: EndBlock periodic distribution (current committee)
  |     +-- DkgCommitteeRewardPortion -> active committee members
  |     +-- Remainder -> withdrawn to UbiWithdrawAddress
  |
  +-- Path 2: Round transition settlement (previous committee)
        +-- DkgCommitteeRewardPortion -> previous committee members
        +-- Remainder -> SettlementBalance
              -> claimed by evmstaking next EndBlock
              -> withdrawn to UbiWithdrawAddress
```

---

## 4. Kernel Upgrade Workflow

### 4.1 Overview

Kernel upgrades allow the story-kernel binary (TEE enclave) to be updated without disrupting the DKG committee. The process involves resharing the existing key material from the old binary to the new binary.

### 4.2 Scheduling

```
Governance/Admin calls DKG.scheduleUpgrade(activationHeight, upgradeVersion)
  |
  Emits: UpgradeScheduled event
  |
  v
evmengine.ProcessDKGUpgradeScheduled()
  |
  +-- Validate: activationHeight > currentBlock, fits int64
  |
  v
dkgKeeper.UpgradeScheduled(activationHeight, upgradeVersion)
  |
  +-- Reject if pending (non-activated) upgrade already exists
  +-- Store KernelUpgradeInfo{upgradeVersion, activationHeight, IsActivated: false}
```

### 4.3 Activation

```
BeginBlocker (at height >= activationHeight)
  |
  +-- hasPendingUpgradeActivation():
  |     Walk KernelUpgradeInfos, find non-activated entry
  |     where currentHeight >= activationHeight
  |
  +-- Mark IsActivated = true (not deleted yet)
  +-- InitiateDKGRound(isUpgrade = true)
        |
        Creates DKGNetwork with IsUpgrade=true, IsResharing=true
```

### 4.4 Upgrade Resharing

During upgrade resharing, the `KernelRouter` has TWO connected endpoints:
- Old binary (old code commitment / MRENCLAVE)
- New binary (new code commitment / MRENCLAVE)

```
KernelRouter
  |
  +-- Endpoint 1: old-binary (CC_old)
  +-- Endpoint 2: new-binary (CC_new)

Registration:
  - New binary: GenerateAndSealKey (CC_new)
  - Old-only members: session.CodeCommitment = CC_old

Dealing:
  - Dealers use OLD binary (key shares sealed by old binary)
  - dealerCC = session.OldCodeCommitment

Processing Deals/Responses/Justifications:
  - New members: send to CC_new
  - For upgrade: also send to CC_old (both binaries maintain DKG state)

Finalization:
  - New binary: FinalizeDKG (CC_new)
```

**The code commitment routing logic**:

| Operation | Non-upgrade | Upgrade |
|-----------|-------------|---------|
| `GenerateAndSealKey` | Previous active round's CC | New binary's CC (differs from oldCC) |
| `GenerateDeals` | Session CC | Old binary's CC (old key shares) |
| `ProcessDeals` | Session CC | Session CC (new binary) |
| `ProcessResponses` | Session CC | Both CC_new AND CC_old |
| `ProcessJustification` | Session CC | Both CC_new AND CC_old |
| `FinalizeDKG` | Session CC | Session CC (new binary) |

### 4.5 Completion

```
FinalizeDKGRound()
  |
  +-- If latestRound.IsUpgrade:
        deleteActivatedUpgradeInfo()
        - Walk KernelUpgradeInfos
        - Delete all entries where IsActivated == true
```

### 4.6 Cancellation

```
DKG.cancelUpgrade(upgradeVersion)
  |
  Emits: UpgradeCancelled event
  |
  v
dkgKeeper.UpgradeCancelled(upgradeVersion)
  |
  +-- Find upgrade info by version
  +-- DeleteKernelUpgradeInfo(upgradeVersion)
```

---

## 5. Key Technical Details

### 5.1 Vote Extensions

Vote extensions are the mechanism by which DKG data propagates through CometBFT consensus.

**Timing after v1.6.0 upgrade**:
- Height H: Upgrade handler sets `vote_extensions_enable_height = H+1`
- Height H+1: CometBFT starts collecting VEs from validators
- Height H+2: VEs appear in `LocalLastCommit`, `PrepareVotes` starts producing `MsgAddDkgVote`

**Size limits**:
- `maxVoteExtensionSize`: 256 KB per vote extension
- `maxItemsPerVote`: 80 deals, 80 responses, or 80 justifications per VE

**Malformation handling**:
- `VerifyVoteExtension` returns `REJECT` (not Go error) for malformed VEs
- Go errors from verification are treated as application bugs by CometBFT

### 5.2 Async Context

DKG operations that communicate with story-kernel (gRPC) run in async goroutines, NOT within the CometBFT consensus context.

```go
const dkgAsyncTimeout = 1 * time.Minute

func dkgAsyncContext() (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), dkgAsyncTimeout)
}
```

**Why**: The CometBFT consensus context is canceled when block processing completes. story-kernel gRPC calls (especially `GenerateAndSealKey` with SGX quote generation) can take 10-30 seconds. Using the consensus context would abort these in-flight calls.

**Consequence**: Async goroutines cannot access the Cosmos SDK KV store (it requires the SDK context). Any data needed from on-chain state must be pre-computed before launching the goroutine. All pre-computations happen in the caller (which has SDK context), and results are passed as parameters to the async handler:

| Pre-computed Value | Used By | Purpose |
|-------------------|---------|---------|
| `oldCC` | `handleDKGRegistration` | Previous round's code commitment for upgrade resharing |
| `shouldDeal` | `handleDKGDealing` | Whether this validator should generate deals |
| `ensureSessionIndex` | `handleDKGDealing`, `handleDKGProcessDeals` | Set session.Index from on-chain registration |
| `shouldProcessResponses` | `handleDKGProcessResponses` | Whether this validator should process responses |
| `buildDealerPubKeyMap` | `handleDKGProcessJustifications` | Dealer DKG public keys for Schnorr verification |

### 5.3 Concurrency Model

```
Mutexes:
  dealsMu          - protects deals queue (EnqueueDeals/DequeueDeals)
  responsesMu      - protects responses queue
  justificationsMu - protects justifications queue
  dkgKernelMu      - package-level mutex; serializes all kernel DKG mutations
                     Locked directly inside each handler:
                     handleDKGProcessDeals, handleDKGProcessResponses,
                     handleDKGProcessJustifications

Atomic flags:
  dkgSvcRound (atomic.Uint64)
    - Tracks which round's goroutine is running
    - tryAcquireDKGSvc: newer round always supersedes older
    - Prevents duplicate goroutines for same round
    - Prevents cross-round blocking

  decryptWorkerRunning (atomic.Bool)
    - Ensures only one decrypt worker loop runs
```

**`dkgKernelMu` is critical**: Without it, concurrent goroutines from `ProcessDeals`, `ProcessResponses`, and `ProcessJustifications` would corrupt the cached `DistKeyGenerator` in story-kernel, causing "different number of coefficients" errors during finalization.

### 5.4 Session State Machine (Off-Chain)

| Phase | Entered When | Exited When |
|-------|-------------|-------------|
| `PhaseInitializing` | `CreateSession` called | `GenerateAndSealKey` + `Register` succeed |
| `PhaseInitialized` | Registration complete | `GenerateDeals` succeeds |
| `PhaseDealing` | Deals generated | `FinalizeDKG` succeeds |
| `PhaseFinalized` | TEE finalization + contract finalize | `handleDKGComplete` marks complete |
| `PhaseCompleted` | Session fully complete | Terminal (starts DecryptWorker) |
| `PhaseFailed` | Any error in above steps | `ResumeDKGService` rewinds phase |

Session persistence: JSON files in `{dataDir}/session_{round}.json`

### 5.5 Security Properties

**SGX Sealing**:
- All private keys (Ed25519 DKG, secp256k1 communication, DistKeyShare) are sealed to the enclave
- Sealed data can only be unsealed by the same MRENCLAVE (same binary build)
- Sealing uses Intel SGX `sgx_seal_data` / Gramine's `/dev/attestation` interface

**Code Commitment (MRENCLAVE)** — validated at two layers:
- **DKG Contract (on-chain)**: `register()` calls `_authenticateEnclaveReport()` → `IAttestationReportValidator.validateReport(expectedCodeCommitment, ...)`, which verifies the SGX quote's MRENCLAVE matches the whitelisted `EnclaveTypeData.codeCommitment`. This is the trust anchor — only enclave binaries with approved code commitments can register.
- **story-kernel (TEE)**: `ValidateCodeCommitment()` verifies each gRPC request's `code_commitment` matches the binary's own MRENCLAVE, preventing cross-enclave request injection.
- `GetCodeCommitment()` returns the running binary's MRENCLAVE for KernelRouter discovery.
- Admin manages the whitelist via `DKG.whitelistEnclaveType(enclaveType, {codeCommitment, validationHookAddr})`.

**Signature Verification**:
- Finalization: ECDSA signature over `(codeCommitment, round, participantsRoot, globalPubKey, publicCoeffs, pubKeyShare)`, verified by recovering signer from `ethHash` and comparing to `address(keccak256(commPubKey))`
- Partial decryption: ECDSA signature over `(round, encryptedPartial, ephemeralPubKey, pubShare)`, verified similarly
- Justifications: Schnorr signature over kyber's canonical `Justification.Hash()`, verified using the dealer's Ed25519 DKG public key

**ParticipantsRoot Validation**:
- `participantsRoot = keccak256(sorted concatenated validator addresses)`
- Both story-kernel and the DKG module independently compute this from verified/finalized registrations
- Mismatch detection prevents a rogue validator from claiming a different participant set

**Replay Protection** (for justifications):
- **Cross-round**: kyber's `SessionID` is derived from `(dealer pubkey + verifiers + commitments + threshold)`. Fresh keys per round produce unique SessionIDs, so `Justification.Hash()` differs across rounds, and Schnorr verification rejects replays.
- **Stage gating**: justifications only accepted during `DKG_STAGE_DEALING`
- **Within-block**: deduplication by `(dealerIndex, recipientIndex)`

---

## 6. story-kernel (TEE) Architecture

### 6.1 gRPC Service Methods

| Method | Input | Output | Purpose |
|--------|-------|--------|---------|
| `GetCodeCommitment` | (empty) | `code_commitment` | Returns this binary's MRENCLAVE for KernelRouter discovery |
| `GenerateAndSealKey` | `code_commitment, round, address` | `code_commitment, round, start_block_height, start_block_hash, dkg_pub_key, comm_pub_key, enclave_report` | Generate + SGX-seal Ed25519 and secp256k1 key pairs; produce attestation quote |
| `GenerateDeals` | `code_commitment, round, is_resharing` | `code_commitment, round, deals[]` | Create Pedersen VSS encrypted deals for all recipients |
| `ProcessDeals` | `code_commitment, round, deals[], is_resharing` | `code_commitment, round, responses[]` | Process received deals, produce responses (approve/complain) |
| `ProcessResponses` | `code_commitment, round, responses[], is_resharing` | `justifications[]` | Process responses; emit justification if a complaint was invalid |
| `ProcessJustification` | `code_commitment, round, justifications[], is_resharing` | (empty) | Restore DKG state for deals that had valid justifications |
| `FinalizeDKG` | `code_commitment, round, is_resharing` | `code_commitment, round, participants_root, global_pub_key, public_coeffs[], pub_key_share, signature` | Compute DistKeyShare, seal it, sign finalization data |
| `PartialDecryptTDH2` | `code_commitment, round, ciphertext, label, global_pub_key, requester_pub_key` | `encrypted_partial_decryption, ephemeral_pub_key, pub_share, signature` | TDH2 partial decrypt + encrypt result to requester |

### 6.2 Store Layer

**DKGStore** (`store/dkg_store.go`):
- Manages SGX-sealed key storage and DKG state persistence
- Key directory structure: `{keyDir}/{round}/{codeCommitmentHex}/`
  - `ed25519_priv.sealed` - SGX-sealed Ed25519 DKG private key
  - `secp256k1_priv.sealed` - SGX-sealed secp256k1 communication key
- State directory structure: `{stateDir}/{round}/{codeCommitmentHex}/`
  - `state.json` - DKG state (pub keys, threshold, deals, responses, justifications)
  - `dist_key_share.sealed` - SGX-sealed distributed key share
- Thread safety: `stateMu` mutex protects concurrent state file access
- `LoadOrGenerate{Ed25519,Secp256k1}Key`: idempotent key creation

**DKGState** structure:
```go
type DKGState struct {
    PubKeys        []kyber.Point       // participant DKG public keys
    Threshold      uint32              // required threshold
    Deals          []dkg.Deal          // processed deals
    Responses      []dkg.Response      // processed responses
    Justifications []dkg.Justification // processed justifications
    PublicCoeffs   []kyber.Point       // polynomial commitments
    FromRound      uint32              // source round (for resharing)
}
```

**KeyStore** (`store/key_store.go`):
- `SealAndStore{Ed25519,Secp256k1}Key`: seal key bytes via SGX, write to file
- `LoadSealed{Ed25519,Secp256k1}Key`: unseal from file via SGX
- Keys are organized by `{round}/{codeCommitmentHex}/`

### 6.3 Cache Layer

| Cache | Type | Key | Purpose |
|-------|------|-----|---------|
| `RoundContextCache` | `map[uint32]*RoundContext` | round | Network + registrations + sorted pub keys for a round |
| `DKGCache` (InitDKG) | `map[string]*dkg.DistKeyGenerator` | round | Cached initial-round DKG state machine |
| `ResharingCache` | `map[string]*dkg.DistKeyGenerator` | `{fromRound}_{toRound}` | Cached resharing DKG state machine |
| `DistKeyShareCache` | `map[uint32]*dkg.DistKeyShare` | round | Cached finalized key share |
| `PIDCache` | `map[uint32]uint32` | round | 1-based participant ID for each round |

All caches use `sync.RWMutex` for thread safety.

**`RoundContextCache` retry logic**: When threshold is 0 (registration phase still in progress), `GetOrLoadRoundContext` retries up to 5 times with 2-second delays, allowing the light client to catch up to the block where `BeginDealing` set the threshold.

> **Light Client in story-kernel**: story-kernel uses a CometBFT light client (`VerifiedQueryClient`) to read on-chain state (DKG networks, registrations) with ICS23 Merkle proof verification. This ensures the TEE cannot be tricked with forged chain data. The light client is configured with a trusted height/hash and syncs headers to verify query responses.

### 6.4 Enclave Layer

- `enclave.GetSelfCodeCommitment()`: reads the binary's own MRENCLAVE via Gramine's `/dev/attestation/my_target_info`
- `enclave.ValidateCodeCommitment(cc)`: compares provided CC against self CC
- `enclave.GetRemoteQuote(reportData)`: generates SGX remote attestation quote via `/dev/attestation/quote`
- `enclave.SealToFile(data, path)`: SGX seal + file write
- `enclave.UnsealFromFile(path)`: file read + SGX unseal

### 6.5 Cryptographic Operations

| Operation | Library | Curve | Usage |
|-----------|---------|-------|-------|
| DKG Key Generation | kyber | Edwards25519 | Ed25519 key pairs for Pedersen DKG |
| VSS (Verifiable Secret Sharing) | kyber (Pedersen) | Edwards25519 | Deal generation, verification, justification |
| Schnorr Signatures | kyber | Edwards25519 | Justification signing and verification |
| Communication Signing | go-ethereum/crypto | secp256k1 | Finalization + partial decryption signatures |
| TDH2 Partial Decrypt | cb-mpc | Edwards25519 (custom curve ID 0x3f) | Threshold decryption of CDR vault data |
| Partial Encryption | ECIES + AES-GCM | secp256k1 | Encrypting partial decryptions to requester's public key |
| Participants Root | go-ethereum/crypto | - | Keccak256 hash of sorted validator addresses |
| Report Data | crypto/sha256 | - | Hash of (address, round, pubkeys, startBlock) for attestation |

---

## Appendix A: On-Chain State Schema

| Collection | Key Format | Value | Module |
|------------|-----------|-------|--------|
| `DKGNetworks` | `"{round}"` | `DKGNetwork` proto | dkg |
| `LatestDKGNetwork` | (singleton) | `"{round}"` string | dkg |
| `LatestActiveRound` | (singleton) | `"{round}"` string | dkg |
| `DKGRegistrations` | `"{round}_{address}"` | `DKGRegistration` proto | dkg |
| `GlobalPubKeyVotes` | `"{round}_{pubKeyHex}_{coeffHash}"` | `uint32` count | dkg |
| `SettlementBalance` | (singleton) | amount string | dkg |
| `KernelUpgradeInfos` | `"{upgradeVersion}"` | `KernelUpgradeInfo` proto | dkg |
| `DKGPartialDecrypt` | `"{requesterHash}_{labelHex}_{round}_{validator}"` | JSON bytes | dkg |
| `DecryptRequestRegistry` | `"{requesterHash}_{labelHex}"` | `uint64` blockHeight | dkg |

## Appendix B: DKG Network Stages

Defined in `tee.proto` as `DKGStage` enum:

| Stage | Value | Description |
|-------|-------|-------------|
| `DKG_STAGE_UNSPECIFIED` | 0 | Should not occur |
| `DKG_STAGE_REGISTRATION` | 1 | Validators registering keys |
| `DKG_STAGE_DEALING` | 2 | Validators exchanging encrypted deals + responses + justifications |
| `DKG_STAGE_FINALIZATION` | 3 | Validators computing + submitting distributed key shares |
| `DKG_STAGE_ACTIVE` | 4 | DKG round is active for threshold decryption |
| `DKG_STAGE_FAILED` | 5 | Round failed, new round initiated |

> Note: `DKGStageEnded` (value 6) is defined in the Story CL code only, not in `tee.proto`. It is set by the DKG module when a previous active round is superseded by a new one.

## Appendix C: DKG Registration Statuses

Defined in `tee.proto` as `DKGRegStatus` enum:

| Status | Value | Description |
|--------|-------|-------------|
| `DKG_REG_STATUS_UNSPECIFIED` | 0 | Should not occur |
| `DKG_REG_STATUS_VERIFIED` | 1 | Registration verified and stored |
| `DKG_REG_STATUS_FINALIZED` | 2 | Validator has finalized for this round |
| `DKG_REG_STATUS_INVALIDATED` | 3 | Dealer invalidated (invalid deal detected via VSS) |

## Appendix D: Index Convention

The system uses two indexing conventions that must be carefully managed:

| Context | Convention | Example |
|---------|-----------|---------|
| On-chain registration index | **1-based** | First participant = 1 |
| Kyber DKG PID | **0-based** | First participant = 0 |
| Deal.RecipientIndex | **0-based** (kyber) | Converted from on-chain: `session.Index - 1` |
| Response.VssResponse.Index | **0-based** (kyber) | Filtered by `resp.Index != session.Index - 1` |
| PIDCache | **1-based** (on-chain) | Used directly in `TDH2PartialDecrypt` |
| PriShare.I (DistKeyShare) | **0-based** (kyber) | Internal kyber convention |

## Appendix E: Complete Lifecycle Sequence Diagram

```
Height H: v1.6.0 upgrade
Height H+1: VE collection starts
Height H+2: First MsgAddDkgVote

Phase 1 - Registration (RegistrationPeriod blocks):
  BeginBlocker: InitiateDKGRound(isUpgrade=false)
  Each validator:
    TEE.GenerateAndSealKey -> (dkgPub, commPub, quote)
    DKG.register(quote, instanceData, startBlock)
    Event: Registered -> dkgKeeper.Registered -> DKGRegistration{status: DKG_REG_STATUS_VERIFIED}

Phase 2 - Dealing (DealingPeriod blocks):
  BeginBlocker: BeginDealing
    Count verified regs >= MinReqRegistered? Yes -> continue; No -> SkipToNextRound
    Set Total + Threshold
  Each dealer (first round: all; resharing: previous set):
    TEE.GenerateDeals -> deals
    Enqueue -> ExtendVote -> VoteExtension -> PrepareVotes -> MsgAddDkgVote
    AddVote:
      ProcessDeals -> TEE.ProcessDeals -> responses -> Enqueue
      ProcessResponses -> TEE.ProcessResponses -> justifications -> Enqueue
      ProcessJustifications -> verify + TEE.ProcessJustification

Phase 3 - Finalization (FinalizationPeriod blocks):
  BeginBlocker: BeginFinalization
  Each current-set validator:
    TEE.FinalizeDKG -> (globalPub, coeffs, pubShare, signature)
    DKG.finalize(round, participantsRoot, globalPub, coeffs, pubShare, sig)
    Event: Finalized -> dkgKeeper.Finalized:
      verifySignature, validateParticipantsRoot
      AddGlobalPubKeyVote -> if votes >= threshold: set global_public_key
      status: DKG_REG_STATUS_VERIFIED -> DKG_REG_STATUS_FINALIZED

Phase 4 - Active (ActivePeriod blocks):
  BeginBlocker: FinalizeDKGRound
    Count finalized >= MinReqFinalized AND >= Threshold? Yes -> continue
    settleRewardsForPreviousCommittee
    endPreviousActiveRound -> stage=DKGStageEnded (CL-only, value 6)
    setLatestActiveRound -> new round is THE active round
    handleDKGComplete -> session=Completed -> StartDecryptWorker

  During Active phase (CDR flow):
    CDR.allocate(conditions) -> VaultAllocated event -> uuid assigned
    CDR.write(uuid, ciphertext) -> VaultWritten event -> encrypted data stored
    CDR.read(uuid) -> VaultRead event
    -> ThresholdDecryptRequested -> queue decrypt request
    -> DecryptWorker: TEE.PartialDecryptTDH2 -> CDR.submitEncryptedPartialDecryption
    -> PartialDecryptionSubmitted -> verify signature + pubShare + timeout
    -> Partials stored in CL -> GetCDRPartials -> TDH2Combine -> plaintext

Phase 5 - Resharing (after ActivePeriod):
  BeginBlocker: shouldTransitionStage -> Registration
  -> InitiateDKGRound(isUpgrade=false) [new round, IsResharing=true]
  -> Cycle repeats
```
